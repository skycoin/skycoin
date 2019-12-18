// Package integration_test implements API integration tests
package integration_test

import (
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/andreyvit/diff"
	"github.com/stretchr/testify/require"

	"github.com/SkycoinProject/skycoin/src/api"
	"github.com/SkycoinProject/skycoin/src/cipher"
	"github.com/SkycoinProject/skycoin/src/coin"
	"github.com/SkycoinProject/skycoin/src/daemon"
	"github.com/SkycoinProject/skycoin/src/readable"
	"github.com/SkycoinProject/skycoin/src/testutil"
	"github.com/SkycoinProject/skycoin/src/transaction"
	"github.com/SkycoinProject/skycoin/src/util/droplet"
	"github.com/SkycoinProject/skycoin/src/util/mathutil"
	"github.com/SkycoinProject/skycoin/src/util/useragent"
	"github.com/SkycoinProject/skycoin/src/visor"
)

/* Runs HTTP API tests against a running skycoin node

Set envvar SKYCOIN_INTEGRATION_TESTS=1 to enable them
Set SKYCOIN_NODE_HOST to the node's address (defaults to http://127.0.0.1:6420)
Set SKYCOIN_INTEGRATION_TEST_MODE to either "stable" or "live" (defaults to "stable")

Each test has two modes:
    1. against a stable, pinned blockchain
    2. against a live, active blockchain

When running mode 1, API responses do not change. The exact responses are compared to saved responses on disk.
Make sure the skycoin node is running against the pinned blockchain data provided in this package's folder.

When running mode 2, API responses may change (such as /coinSupply). The exact responses are not compared,
but the response is checked to be unmarshallable to a known JSON object.

When update flag is set to true all tests pass
*/

const (
	testModeStable           = "stable"
	testModeLive             = "live"
	testModeDisableWalletAPI = "disable-wallet-api"
	testModeEnableSeedAPI    = "enable-seed-api"
	testModeDisableGUI       = "disable-gui"

	testFixturesDir = "testdata"
)

type TestData struct {
	actual   interface{}
	expected interface{}
}

var update = flag.Bool("update", false, "update golden files")
var testLiveWallet = flag.Bool("test-live-wallet", false, "run live wallet tests, requires wallet envvars set")

func nodeAddress() string {
	addr := os.Getenv("SKYCOIN_NODE_HOST")
	if addr == "" {
		return "http://127.0.0.1:6420"
	}
	return addr
}

func nodeUsername() string {
	return os.Getenv("SKYCOIN_NODE_USERNAME")
}

func nodePassword() string {
	return os.Getenv("SKYCOIN_NODE_PASSWORD")
}

func newClient() *api.Client {
	c := api.NewClient(nodeAddress())
	c.SetAuth(nodeUsername(), nodePassword())
	return c
}

func mode(t *testing.T) string {
	mode := os.Getenv("SKYCOIN_INTEGRATION_TEST_MODE")
	switch mode {
	case "":
		mode = testModeStable
	case testModeLive,
		testModeStable,
		testModeDisableWalletAPI,
		testModeEnableSeedAPI,
		testModeDisableGUI:
	default:
		t.Fatal("Invalid test mode, must be stable, live or disable-wallet-api")
	}
	return mode
}

func enabled() bool {
	return os.Getenv("SKYCOIN_INTEGRATION_TESTS") == "1"
}

func useCSRF(t *testing.T) bool {
	x := os.Getenv("USE_CSRF")
	if x == "" {
		return false
	}

	useCSRF, err := strconv.ParseBool(x)
	require.NoError(t, err)
	return useCSRF
}

func doHeaderCheck(t *testing.T) bool {
	x := os.Getenv("HEADER_CHECK")
	if x == "" {
		return false
	}

	doHeaderCheck, err := strconv.ParseBool(x)
	require.NoError(t, err)
	return doHeaderCheck
}

func doStable(t *testing.T) bool {
	if enabled() && mode(t) == testModeStable {
		return true
	}

	t.Skip("Stable tests disabled")
	return false
}

func doLive(t *testing.T) bool {
	if enabled() && mode(t) == testModeLive {
		return true
	}

	t.Skip("Live tests disabled")
	return false
}

func doDisableWalletAPI(t *testing.T) bool {
	if enabled() && mode(t) == testModeDisableWalletAPI {
		return true
	}

	t.Skip("DisableWalletApi tests disabled")
	return false
}

func doEnableSeedAPI(t *testing.T) bool {
	if enabled() && mode(t) == testModeEnableSeedAPI {
		return true
	}

	t.Skip("EnableSeedAPI tests disabled")
	return false
}

func doDisableGUI(t *testing.T) bool {
	if enabled() && mode(t) == testModeDisableGUI {
		return true
	}

	t.Skip("DisableGUIAPI tests disabled")
	return false
}

func doLiveOrStable(t *testing.T) bool {
	if enabled() {
		switch mode(t) {
		case testModeStable, testModeLive:
			return true
		}
	}

	t.Skip("Live and stable tests disabled")
	return false
}

func doLiveWallet(t *testing.T) bool {
	if *testLiveWallet {
		return true
	}

	t.Skip("Tests requiring wallet envvars are disabled")
	return false
}

func envParseBool(t *testing.T, key string) bool {
	x := os.Getenv(key)
	if x == "" {
		return false
	}

	v, err := strconv.ParseBool(x)
	require.NoError(t, err)
	return v
}

func dbNoUnconfirmed(t *testing.T) bool {
	return envParseBool(t, "DB_NO_UNCONFIRMED")
}

func liveDisableNetworking(t *testing.T) bool {
	return envParseBool(t, "LIVE_DISABLE_NETWORKING")
}

func loadGoldenFile(t *testing.T, filename string, testData TestData) {
	require.NotEmpty(t, filename, "loadGoldenFile golden filename missing")

	goldenFile := filepath.Join(testFixturesDir, filename)

	if *update {
		updateGoldenFile(t, goldenFile, testData.actual)
	}

	f, err := os.Open(goldenFile)
	require.NoError(t, err)
	defer f.Close()

	d := json.NewDecoder(f)
	d.DisallowUnknownFields()

	err = d.Decode(testData.expected)
	require.NoError(t, err, filename)
}

func updateGoldenFile(t *testing.T, filename string, content interface{}) {
	contentJSON, err := json.MarshalIndent(content, "", "\t")
	require.NoError(t, err)
	contentJSON = append(contentJSON, '\n')
	err = ioutil.WriteFile(filename, contentJSON, 0644)
	require.NoError(t, err)
}

func checkGoldenFile(t *testing.T, goldenFile string, td TestData) {
	loadGoldenFile(t, goldenFile, td)
	require.Equal(t, reflect.Indirect(reflect.ValueOf(td.expected)).Interface(), td.actual)

	// Serialize expected to JSON and compare to the goldenFile's contents
	// This will detect field changes that could be missed otherwise
	b, err := json.MarshalIndent(td.expected, "", "\t")
	require.NoError(t, err)

	goldenFile = filepath.Join(testFixturesDir, goldenFile)

	f, err := os.Open(goldenFile)
	require.NoError(t, err)
	defer f.Close()

	c, err := ioutil.ReadAll(f)
	require.NoError(t, err)

	sc := string(c)
	sb := string(b) + "\n"

	require.Equal(t, sc, sb, "JSON struct output differs from golden file, was a field added to the struct?\nDiff:\n"+diff.LineDiff(sc, sb))
}

func assertResponseError(t *testing.T, err error, errCode int, errMsg string) {
	require.Error(t, err)
	require.IsType(t, api.ClientError{}, err)
	require.Equal(t, errCode, err.(api.ClientError).StatusCode, err.(api.ClientError).Message)
	require.Equal(t, errMsg, err.(api.ClientError).Message)
}

func TestStableCoinSupply(t *testing.T) {
	if !doStable(t) {
		return
	}

	c := newClient()

	cs, err := c.CoinSupply()
	require.NoError(t, err)

	var expected api.CoinSupply
	checkGoldenFile(t, "coinsupply.golden", TestData{*cs, &expected})
}

func TestLiveCoinSupply(t *testing.T) {
	if !doLive(t) {
		return
	}

	c := newClient()

	cs, err := c.CoinSupply()
	require.NoError(t, err)

	require.NotEmpty(t, cs.CurrentSupply)
	require.NotEmpty(t, cs.TotalSupply)
	require.NotEmpty(t, cs.MaxSupply)
	require.Equal(t, "100000000.000000", cs.MaxSupply)
	require.NotEmpty(t, cs.CurrentCoinHourSupply)
	require.NotEmpty(t, cs.TotalCoinHourSupply)
	require.Equal(t, 100, len(cs.UnlockedAddresses)+len(cs.LockedAddresses))
}

func TestVersion(t *testing.T) {
	if !doLiveOrStable(t) {
		return
	}

	c := newClient()

	v, err := c.Version()
	require.NoError(t, err)

	require.NotEmpty(t, v.Version)
}

func TestVerifyAddress(t *testing.T) {
	if !doLiveOrStable(t) {
		return
	}

	c := newClient()

	cases := []struct {
		name    string
		golden  string
		addr    string
		errCode int
		errMsg  string
	}{
		{
			name:   "valid address",
			golden: "verify-address.golden",
			addr:   "7cpQ7t3PZZXvjTst8G7Uvs7XH4LeM8fBPD",
		},

		{
			name:    "invalid address",
			addr:    "7apQ7t3PZZXvjTst8G7Uvs7XH4LeM8fBPD",
			errCode: http.StatusUnprocessableEntity,
			errMsg:  "Invalid checksum",
		},

		{
			name:    "missing address",
			addr:    "",
			errCode: http.StatusBadRequest,
			errMsg:  "address is required",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := c.VerifyAddress(tc.addr)

			if tc.errCode != 0 && tc.errCode != http.StatusOK {
				assertResponseError(t, err, tc.errCode, tc.errMsg)
				return
			}

			require.NoError(t, err)

			var expected api.VerifyAddressResponse
			checkGoldenFile(t, tc.golden, TestData{*resp, &expected})
		})
	}
}

func TestStableVerifyTransaction(t *testing.T) {
	if !doStable(t) {
		return
	}

	c := newClient()

	badSigStr := "71f2c01516fe696328e79bcf464eb0db374b63d494f7a307d1e77114f18581d7a81eed5275a9e04a336292dd2fd16977d9bef2a54ea3161d0876603d00c53bc9dd"
	badSigBytes, err := hex.DecodeString(badSigStr)
	require.NoError(t, err)
	badSig := cipher.MustNewSig(badSigBytes)

	inputHash := "75692aeff988ce0da734c474dbef3a1ce19a5a6823bbcd36acb856c83262261e"
	input := testutil.SHA256FromHex(t, inputHash)

	destAddrStr := "7cpQ7t3PZZXvjTst8G7Uvs7XH4LeM8fBPD"
	destAddr, err := cipher.DecodeBase58Address(destAddrStr)
	require.NoError(t, err)

	inputAddrStr := "qxmeHkwgAMfwXyaQrwv9jq3qt228xMuoT5"
	inputAddr, err := cipher.DecodeBase58Address(inputAddrStr)
	require.NoError(t, err)

	badSignatureTxn := coin.Transaction{
		Sigs: []cipher.Sig{badSig},
		In:   []cipher.SHA256{input},
		Out: []coin.TransactionOutput{
			{
				Address: destAddr,
				Coins:   1e3,
				Hours:   10,
			},
			{
				Address: inputAddr,
				Coins:   22100e6 - 1e3,
				Hours:   188761,
			},
		},
	}
	err = badSignatureTxn.UpdateHeader()
	require.NoError(t, err)

	cases := []struct {
		name     string
		golden   string
		txn      coin.Transaction
		unsigned bool
		errCode  int
		errMsg   string
	}{
		{
			name:    "unsigned=false invalid transaction empty",
			txn:     coin.Transaction{},
			golden:  "verify-transaction-invalid-empty.golden",
			errCode: http.StatusUnprocessableEntity,
			errMsg:  "Transaction violates soft constraint: Transaction has zero coinhour fee",
		},

		{
			name:    "unsigned=false invalid transaction bad signature",
			txn:     badSignatureTxn,
			golden:  "verify-transaction-invalid-bad-sig.golden",
			errCode: http.StatusUnprocessableEntity,
			errMsg:  "Transaction violates hard constraint: Signature not valid for hash",
		},

		{
			name:     "unsigned=true invalid transaction empty",
			txn:      coin.Transaction{},
			unsigned: true,
			golden:   "verify-transaction-invalid-empty.golden",
			errCode:  http.StatusUnprocessableEntity,
			errMsg:   "Transaction violates soft constraint: Transaction has zero coinhour fee",
		},

		{
			name:     "unsigned=true invalid transaction bad signature",
			txn:      badSignatureTxn,
			unsigned: true,
			golden:   "verify-transaction-invalid-bad-sig.golden",
			errCode:  http.StatusUnprocessableEntity,
			errMsg:   "Transaction violates hard constraint: Signature not valid for hash",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := c.VerifyTransaction(api.VerifyTransactionRequest{
				EncodedTransaction: tc.txn.MustSerializeHex(),
				Unsigned:           tc.unsigned,
			})

			if tc.errCode != 0 && tc.errCode != http.StatusOK {
				assertResponseError(t, err, tc.errCode, tc.errMsg)
				if tc.errCode != http.StatusUnprocessableEntity {
					return
				}
			}

			if tc.errCode != http.StatusUnprocessableEntity {
				require.NoError(t, err)
			}

			var expected api.VerifyTransactionResponse
			checkGoldenFile(t, tc.golden, TestData{*resp, &expected})
		})
	}

}

func TestStableNoUnconfirmedOutputs(t *testing.T) {
	if !doStable(t) || !dbNoUnconfirmed(t) {
		return
	}

	c := newClient()

	cases := []struct {
		name    string
		golden  string
		addrs   []string
		hashes  []string
		errCode int
		errMsg  string
	}{
		{
			name:   "no addrs or hashes",
			golden: "no-unconfirmed-outputs-noargs.golden",
		},
		{
			name: "only addrs",
			addrs: []string{
				"ALJVNKYL7WGxFBSriiZuwZKWD4b7fbV1od",
				"2THDupTBEo7UqB6dsVizkYUvkKq82Qn4gjf",
				"qxmeHkwgAMfwXyaQrwv9jq3qt228xMuoT5",
			},
			golden: "no-unconfirmed-outputs-addrs.golden",
		},
		{
			name: "only hashes",
			hashes: []string{
				"9e53268a18f8d32a44b4fb183033b49bebfe9d0da3bf3ef2ad1d560500aa54c6",
				"d91e07318227651129b715d2db448ae245b442acd08c8b4525a934f0e87efce9",
				"01f9c1d6c83dbc1c993357436cdf7f214acd0bfa107ff7f1466d1b18ec03563e",
				"fe6762d753d626115c8dd3a053b5fb75d6d419a8d0fb1478c5fffc1fe41c5f20",
			},
			golden: "no-unconfirmed-outputs-hashes.golden",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			require.False(t, tc.addrs != nil && tc.hashes != nil)

			var outputs *readable.UnspentOutputsSummary
			var err error
			switch {
			case tc.addrs == nil && tc.hashes == nil:
				outputs, err = c.Outputs()
			case tc.addrs != nil:
				outputs, err = c.OutputsForAddresses(tc.addrs)
			case tc.hashes != nil:
				outputs, err = c.OutputsForHashes(tc.hashes)
			}

			if tc.errCode != 0 && tc.errCode != http.StatusOK {
				assertResponseError(t, err, tc.errCode, tc.errMsg)
				return
			}

			require.NoError(t, err)

			var expected readable.UnspentOutputsSummary
			checkGoldenFile(t, tc.golden, TestData{*outputs, &expected})

			require.Equal(t, len(expected.HeadOutputs), len(outputs.HeadOutputs))
			require.Equal(t, len(expected.OutgoingOutputs), len(outputs.OutgoingOutputs))
			require.Equal(t, len(expected.IncomingOutputs), len(outputs.IncomingOutputs))

			for i, o := range expected.HeadOutputs {
				require.Equal(t, o, outputs.HeadOutputs[i], "mismatch at index %d", i)
			}
		})
	}
}

func TestStableOutputs(t *testing.T) {
	if !doStable(t) || dbNoUnconfirmed(t) {
		return
	}

	c := newClient()

	cases := []struct {
		name    string
		golden  string
		addrs   []string
		hashes  []string
		errCode int
		errMsg  string
	}{
		{
			name:   "no addrs or hashes",
			golden: "outputs-noargs.golden",
		},
		{
			name: "only addrs",
			addrs: []string{
				"ALJVNKYL7WGxFBSriiZuwZKWD4b7fbV1od",
				"2THDupTBEo7UqB6dsVizkYUvkKq82Qn4gjf",
				"qxmeHkwgAMfwXyaQrwv9jq3qt228xMuoT5",
				"212mwY3Dmey6vwnWpiph99zzCmopXTqeVEN",
			},
			golden: "outputs-addrs.golden",
		},
		{
			name: "only hashes",
			hashes: []string{
				"9e53268a18f8d32a44b4fb183033b49bebfe9d0da3bf3ef2ad1d560500aa54c6",
				"d91e07318227651129b715d2db448ae245b442acd08c8b4525a934f0e87efce9",
				"01f9c1d6c83dbc1c993357436cdf7f214acd0bfa107ff7f1466d1b18ec03563e",
				"fe6762d753d626115c8dd3a053b5fb75d6d419a8d0fb1478c5fffc1fe41c5f20",
				"701d23fd513bad325938ba56869f9faba19384a8ec3dd41833aff147eac53947",
				"540582ee4128b733f810f149e908d984a5f403ad2865108e6c1c5423aeefc759",
			},
			golden: "outputs-hashes.golden",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			require.False(t, tc.addrs != nil && tc.hashes != nil)

			var outputs *readable.UnspentOutputsSummary
			var err error
			switch {
			case tc.addrs == nil && tc.hashes == nil:
				outputs, err = c.Outputs()
			case tc.addrs != nil:
				outputs, err = c.OutputsForAddresses(tc.addrs)
			case tc.hashes != nil:
				outputs, err = c.OutputsForHashes(tc.hashes)
			}

			if tc.errCode != 0 && tc.errCode != http.StatusOK {
				assertResponseError(t, err, tc.errCode, tc.errMsg)
				return
			}

			require.NoError(t, err)

			var expected readable.UnspentOutputsSummary
			checkGoldenFile(t, tc.golden, TestData{*outputs, &expected})

			require.Equal(t, len(expected.HeadOutputs), len(outputs.HeadOutputs))
			require.Equal(t, len(expected.OutgoingOutputs), len(outputs.OutgoingOutputs))
			require.Equal(t, len(expected.IncomingOutputs), len(outputs.IncomingOutputs))

			for i, o := range expected.HeadOutputs {
				require.Equal(t, o, outputs.HeadOutputs[i], "mismatch at index %d", i)
			}
		})
	}
}

func TestLiveOutputs(t *testing.T) {
	if !doLive(t) {
		return
	}

	c := newClient()

	// Request all outputs and check that HeadOutputs is not empty
	// OutgoingOutputs and IncomingOutputs are variable and could be empty
	outputs, err := c.Outputs()
	require.NoError(t, err)
	require.NotEmpty(t, outputs.HeadOutputs)

	outputs, err = c.OutputsForAddresses(nil)
	require.NoError(t, err)
	require.NotEmpty(t, outputs.HeadOutputs)

	outputs, err = c.OutputsForHashes(nil)
	require.NoError(t, err)
	require.NotEmpty(t, outputs.HeadOutputs)
}

func TestStableBlock(t *testing.T) {
	if !doStable(t) {
		return
	}

	testKnownBlocks(t)
}

// These blocks were affected by the coinhour overflow issue or by coinhour fee calculation bugs,
// make sure that they can be queried
var knownBadBlockSeqs = []uint64{
	// coinhour fee calculation mistake, related to distribution addresses:
	297,
	741,
	743,
	749,
	796,
	4956,
	10125,
	// coinhour overflow related:
	11685,
	11707,
	11710,
	11709,
	11705,
	11708,
	11711,
	11706,
	11699,
	13277,
}

func TestLiveBlock(t *testing.T) {
	if !doLive(t) {
		return
	}

	if liveDisableNetworking(t) {
		t.Skip("Skipping slow block tests when networking disabled")
		return
	}

	testKnownBlocks(t)

	// Check the knownBadBlockSeqs
	c := newClient()
	for _, seq := range knownBadBlockSeqs {
		b, err := c.BlockBySeq(seq)
		require.NoError(t, err)
		require.Equal(t, seq, b.Head.BkSeq)
	}
}

func testKnownBlocks(t *testing.T) {
	c := newClient()

	cases := []struct {
		name    string
		golden  string
		hash    string
		seq     uint64
		errCode int
		errMsg  string
	}{
		{
			name:    "unknown hash",
			hash:    "80744ec25e6233f40074d35bf0bfdbddfac777869b954a96833cb89f44204444",
			errCode: http.StatusNotFound,
			errMsg:  "404 Not Found",
		},
		{
			name:   "valid hash",
			golden: "block-hash.golden",
			hash:   "70584db7fb8ab88b8dbcfed72ddc42a1aeb8c4882266dbb78439ba3efcd0458d",
		},
		{
			name:   "genesis hash",
			golden: "block-hash-genesis.golden",
			hash:   "0551a1e5af999fe8fff529f6f2ab341e1e33db95135eef1b2be44fe6981349f3",
		},
		{
			name:   "genesis seq",
			golden: "block-seq-0.golden",
			seq:    0,
		},
		{
			name:   "seq 1",
			golden: "block-seq-1.golden",
			seq:    1,
		},
		{
			name:   "seq 100",
			golden: "block-seq-100.golden",
			seq:    100,
		},
		{
			name:    "unknown seq",
			seq:     99999999999999,
			errCode: http.StatusNotFound,
			errMsg:  "404 Not Found",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var b *readable.Block
			var err error

			if tc.hash != "" {
				b, err = c.BlockByHash(tc.hash)
			} else {
				b, err = c.BlockBySeq(tc.seq)
			}

			if tc.errCode != 0 && tc.errCode != http.StatusOK {
				assertResponseError(t, err, tc.errCode, tc.errMsg)
				return
			}

			require.NotNil(t, b)

			var expected readable.Block
			checkGoldenFile(t, tc.golden, TestData{*b, &expected})
		})
	}

	t.Logf("Querying every block in the blockchain")

	// Scan every block by seq
	progress, err := c.BlockchainProgress()
	require.NoError(t, err)

	var prevBlock *readable.Block
	for i := uint64(0); i < progress.Current; i++ {
		t.Run(fmt.Sprintf("block-seq-%d", i), func(t *testing.T) {
			b, err := c.BlockBySeq(i)
			require.NoError(t, err)
			require.NotNil(t, b)
			require.Equal(t, i, b.Head.BkSeq)

			if prevBlock != nil {
				require.Equal(t, prevBlock.Head.Hash, b.Head.PreviousHash, "%s != %s", prevBlock.Head.Hash, b.Head.PreviousHash)
			}

			bHash, err := c.BlockByHash(b.Head.Hash)
			require.NoError(t, err)
			require.NotNil(t, bHash)
			require.Equal(t, b, bHash)

			prevBlock = b
		})
	}
}

func TestStableBlockVerbose(t *testing.T) {
	if !doStable(t) {
		return
	}

	testKnownBlocksVerbose(t)
}

func TestLiveBlockVerbose(t *testing.T) {
	if !doLive(t) {
		return
	}

	if liveDisableNetworking(t) {
		t.Skip("Skipping slow block tests when networking disabled")
		return
	}

	testKnownBlocksVerbose(t)

	// Check the knownBadBlockSeqs
	c := newClient()
	for _, seq := range knownBadBlockSeqs {
		b, err := c.BlockBySeqVerbose(seq)
		require.NoError(t, err)
		require.Equal(t, seq, b.Head.BkSeq)
		assertVerboseBlockFee(t, b)
	}
}

func testKnownBlocksVerbose(t *testing.T) {
	c := newClient()

	cases := []struct {
		name    string
		golden  string
		hash    string
		seq     uint64
		errCode int
		errMsg  string
	}{
		{
			name:    "unknown hash",
			hash:    "80744ec25e6233f40074d35bf0bfdbddfac777869b954a96833cb89f44204444",
			errCode: http.StatusNotFound,
			errMsg:  "404 Not Found",
		},
		{
			name:   "valid hash",
			golden: "block-hash-verbose.golden",
			hash:   "70584db7fb8ab88b8dbcfed72ddc42a1aeb8c4882266dbb78439ba3efcd0458d",
		},
		{
			name:   "genesis hash",
			golden: "block-hash-verbose-genesis.golden",
			hash:   "0551a1e5af999fe8fff529f6f2ab341e1e33db95135eef1b2be44fe6981349f3",
		},
		{
			name:   "genesis seq",
			golden: "block-seq-verbose-0.golden",
			seq:    0,
		},
		{
			name:   "seq 1",
			golden: "block-seq-verbose-1.golden",
			seq:    1,
		},
		{
			name:   "seq 100",
			golden: "block-seq-verbose-100.golden",
			seq:    100,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var b *readable.BlockVerbose
			var err error

			if tc.hash != "" {
				b, err = c.BlockByHashVerbose(tc.hash)
			} else {
				b, err = c.BlockBySeqVerbose(tc.seq)
			}

			if tc.errCode != 0 && tc.errCode != http.StatusOK {
				assertResponseError(t, err, tc.errCode, tc.errMsg)
				return
			}

			require.NotNil(t, b)
			assertVerboseBlockFee(t, b)

			var expected readable.BlockVerbose
			checkGoldenFile(t, tc.golden, TestData{*b, &expected})
		})
	}

	t.Logf("Querying every block in the blockchain")

	// Scan every block by seq
	progress, err := c.BlockchainProgress()
	require.NoError(t, err)

	var prevBlock *readable.BlockVerbose
	for i := uint64(0); i < progress.Current; i++ {
		t.Run(fmt.Sprintf("block-seq-verbose-%d", i), func(t *testing.T) {
			b, err := c.BlockBySeqVerbose(i)
			require.NoError(t, err)
			require.NotNil(t, b)
			require.Equal(t, i, b.Head.BkSeq)
			assertVerboseBlockFee(t, b)

			if prevBlock != nil {
				require.Equal(t, prevBlock.Head.Hash, b.Head.PreviousHash)
			}

			bHash, err := c.BlockByHashVerbose(b.Head.Hash)
			require.NoError(t, err)
			require.NotNil(t, bHash)
			require.Equal(t, b, bHash)

			prevBlock = b
		})
	}
}

// assertVerboseBlockFee checks that the block's fee matches the calculated fee of the block's transactions
func assertVerboseBlockFee(t *testing.T, b *readable.BlockVerbose) {
	fee := uint64(0)
	for _, txn := range b.Body.Transactions {
		var err error
		fee, err = mathutil.AddUint64(fee, txn.Fee)
		require.NoError(t, err)
	}

	// The estimated transaction fees should equal the block fee, but in a few cases
	// it doesn't due to older bugs in fee calculation
	if b.Head.Fee != fee {
		switch b.Head.BkSeq {
		case 297:
			require.Equal(t, b.Head.Fee, uint64(3477395194))
			require.Equal(t, fee, uint64(8601490771))
		case 741:
			require.Equal(t, b.Head.Fee, uint64(2093567995))
			require.Equal(t, fee, uint64(17465854723))
		case 743:
			require.Equal(t, b.Head.Fee, uint64(2093809661))
			require.Equal(t, fee, uint64(17466096389))
		case 749:
			require.Equal(t, b.Head.Fee, uint64(1572050737))
			require.Equal(t, fee, uint64(16944337465))
		case 796:
			require.Equal(t, b.Head.Fee, uint64(3197771253))
			require.Equal(t, fee, uint64(13445962405))
		case 4956:
			require.Equal(t, b.Head.Fee, uint64(2309386399))
			require.Equal(t, fee, uint64(22805768703))
		case 10125:
			require.Equal(t, b.Head.Fee, uint64(1938082460))
			require.Equal(t, fee, uint64(22434464764))
		case 13277:
			// In this case, the hours overflow, so the API reports calculated_hours and fee as 0
			require.Equal(t, b.Head.Fee, uint64(3))
			require.Equal(t, fee, uint64(0))
		default:
			require.Equal(t, b.Head.Fee, fee, "Block seq=%d fee does not match sum of transaction fees %d != %d", b.Head.BkSeq, b.Head.Fee, fee)
		}
	}
}

func TestStableBlockchainMetadata(t *testing.T) {
	if !doStable(t) {
		return
	}

	c := newClient()

	metadata, err := c.BlockchainMetadata()
	require.NoError(t, err)

	var expected readable.BlockchainMetadata

	goldenFile := "blockchain-metadata.golden"
	if dbNoUnconfirmed(t) {
		goldenFile = "blockchain-metadata-no-unconfirmed.golden"
	}

	checkGoldenFile(t, goldenFile, TestData{*metadata, &expected})
}

func TestLiveBlockchainMetadata(t *testing.T) {
	if !doLive(t) {
		return
	}

	c := newClient()

	metadata, err := c.BlockchainMetadata()
	require.NoError(t, err)

	require.NotEqual(t, uint64(0), metadata.Head.BkSeq)
}

func TestStableBlockchainProgress(t *testing.T) {
	if !doStable(t) {
		return
	}

	c := newClient()

	progress, err := c.BlockchainProgress()
	require.NoError(t, err)

	var expected readable.BlockchainProgress
	checkGoldenFile(t, "blockchain-progress.golden", TestData{*progress, &expected})
}

func TestLiveBlockchainProgress(t *testing.T) {
	if !doLive(t) {
		return
	}

	c := newClient()

	progress, err := c.BlockchainProgress()
	require.NoError(t, err)

	require.NotEqual(t, uint64(0), progress.Current)

	if liveDisableNetworking(t) {
		require.Empty(t, progress.Peers)
		require.Equal(t, progress.Current, progress.Highest)
	} else {
		require.NotEmpty(t, progress.Peers)
		require.True(t, progress.Current <= progress.Highest)
	}
}

func TestStableBalance(t *testing.T) {
	if !doStable(t) {
		return
	}

	type balanceTestCase struct {
		name   string
		golden string
		addrs  []string
	}

	c := newClient()

	cases := []balanceTestCase{
		{
			name:   "unknown address",
			addrs:  []string{"prRXwTcDK24hs6AFxj69UuWae3LzhrsPW9"},
			golden: "balance-noaddrs.golden",
		},
		{
			name:   "one address",
			addrs:  []string{"2THDupTBEo7UqB6dsVizkYUvkKq82Qn4gjf"},
			golden: "balance-2THDupTBEo7UqB6dsVizkYUvkKq82Qn4gjf.golden",
		},
		{
			name:   "duplicate addresses",
			addrs:  []string{"2THDupTBEo7UqB6dsVizkYUvkKq82Qn4gjf", "2THDupTBEo7UqB6dsVizkYUvkKq82Qn4gjf"},
			golden: "balance-2THDupTBEo7UqB6dsVizkYUvkKq82Qn4gjf.golden",
		},
		{
			name:   "two addresses",
			addrs:  []string{"2THDupTBEo7UqB6dsVizkYUvkKq82Qn4gjf", "qxmeHkwgAMfwXyaQrwv9jq3qt228xMuoT5"},
			golden: "balance-two-addrs.golden",
		},
	}

	if !dbNoUnconfirmed(t) {
		cases = append(cases, balanceTestCase{
			name:   "balance affected by unconfirmed transaction",
			addrs:  []string{"R6aHqKWSQfvpdo2fGSrq4F1RYXkBWR9HHJ", "212mwY3Dmey6vwnWpiph99zzCmopXTqeVEN"},
			golden: "balance-affected-by-unconfirmed-txns.golden",
		})
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			balance, err := c.Balance(tc.addrs)
			require.NoError(t, err)

			var expected api.BalanceResponse
			checkGoldenFile(t, tc.golden, TestData{*balance, &expected})
		})
	}
}

func TestLiveBalance(t *testing.T) {
	if !doLive(t) {
		return
	}

	c := newClient()

	// Genesis address check, should not have a balance
	b, err := c.Balance([]string{"2jBbGxZRGoQG1mqhPBnXnLTxK6oxsTf8os6"})
	require.NoError(t, err)
	require.Equal(t, api.BalanceResponse{
		Addresses: readable.AddressBalances{
			"2jBbGxZRGoQG1mqhPBnXnLTxK6oxsTf8os6": readable.BalancePair{},
		},
	}, *b)

	// Balance of final distribution address. Should have the same coins balance
	// for the next 15-20 years.
	b, err = c.Balance([]string{"ejJjiCwp86ykmFr5iTJ8LxQXJ2wJPTYmkm"})
	require.NoError(t, err)
	require.Equal(t, b.Confirmed, b.Predicted)
	require.NotEmpty(t, b.Confirmed.Hours)
	// Add 1e4 because someone sent 0.01 coins to it
	expectedBalance := uint64(1e6*1e6 + 1e4)
	require.Equal(t, expectedBalance, b.Confirmed.Coins)

	// Check that the balance is queryable for addresses known to be affected
	// by the coinhour overflow problem
	addrs := []string{
		"n7AR1VMW1pK7F9TxhYdnr3HoXEQ3g9iTNP",
		"2aTzmXi9jyiq45oTRFCP9Y7dcvnT6Rsp7u",
		"FjFLnus2ePxuaPTXFXfpw6cVAE5owT1t3P",
		"KT9vosieyWhn9yWdY8w7UZ6tk31KH4NAQK",
	}
	for _, a := range addrs {
		_, err := c.Balance([]string{a})
		require.NoError(t, err, "Failed to get balance of address %s", a)
	}
	_, err = c.Balance(addrs)
	require.NoError(t, err)
}

func TestStableUxOut(t *testing.T) {
	if !doStable(t) {
		return
	}

	c := newClient()

	cases := []struct {
		name   string
		golden string
		uxID   string
	}{
		{
			name:   "valid uxID - unspent",
			golden: "uxout.golden",
			uxID:   "fe6762d753d626115c8dd3a053b5fb75d6d419a8d0fb1478c5fffc1fe41c5f20",
		},
		{
			name:   "valid uxID - spent",
			golden: "uxout-spent-179.golden",
			uxID:   "8e55f10a0615a0737e6906132e09ac08a206971ba4b656f004acc7f4b7889bc8",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ux, err := c.UxOut(tc.uxID)
			require.NoError(t, err)

			var expected readable.SpentOutput
			checkGoldenFile(t, tc.golden, TestData{*ux, &expected})
		})
	}

	// Scan all uxouts from the result of /outputs
	scanUxOuts(t)
}

func TestLiveUxOut(t *testing.T) {
	if !doLive(t) {
		return
	}

	if liveDisableNetworking(t) {
		t.Skip("Skipping slow ux out tests when networking disabled")
		return
	}

	c := newClient()

	// A spent uxout should never change
	ux, err := c.UxOut("fe6762d753d626115c8dd3a053b5fb75d6d419a8d0fb1478c5fffc1fe41c5f20")
	require.NoError(t, err)

	var expected readable.SpentOutput
	checkGoldenFile(t, "uxout-spent.golden", TestData{*ux, &expected})
	require.NotEqual(t, uint64(0), ux.SpentBlockSeq)

	// Scan all uxouts from the result of /outputs
	scanUxOuts(t)
}

func scanUxOuts(t *testing.T) {
	c := newClient()

	outputs, err := c.Outputs()
	require.NoError(t, err)

	for _, ux := range outputs.HeadOutputs {
		t.Run(ux.Hash, func(t *testing.T) {
			foundUx, err := c.UxOut(ux.Hash)
			require.NoError(t, err)

			require.Equal(t, ux.Hash, foundUx.Uxid)
			require.Equal(t, ux.Time, foundUx.Time)
			require.Equal(t, ux.BkSeq, foundUx.SrcBkSeq)
			require.Equal(t, ux.SourceTransaction, foundUx.SrcTx)
			require.Equal(t, ux.Address, foundUx.OwnerAddress)
			require.Equal(t, ux.Hours, foundUx.Hours)
			coinsStr, err := droplet.ToString(foundUx.Coins)
			require.NoError(t, err)
			require.Equal(t, ux.Coins, coinsStr)

			if foundUx.SpentBlockSeq == 0 {
				require.Equal(t, "0000000000000000000000000000000000000000000000000000000000000000", foundUx.SpentTxnID)
			} else {
				require.NotEqual(t, "0000000000000000000000000000000000000000000000000000000000000000", foundUx.SpentTxnID)
			}
		})
	}
}

func TestStableAddressUxOuts(t *testing.T) {
	if !doStable(t) {
		return
	}

	c := newClient()

	cases := []struct {
		name    string
		errCode int
		errMsg  string
		golden  string
		addr    string
	}{
		{
			name:    "no addresses",
			errCode: http.StatusBadRequest,
			errMsg:  "400 Bad Request - address is empty",
		},
		{
			name:   "unknown address",
			addr:   "prRXwTcDK24hs6AFxj69UuWae3LzhrsPW9",
			golden: "uxout-noaddr.golden",
		},
		{
			name:   "one address",
			addr:   "2THDupTBEo7UqB6dsVizkYUvkKq82Qn4gjf",
			golden: "uxout-addr.golden",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ux, err := c.AddressUxOuts(tc.addr)
			if tc.errCode != 0 && tc.errCode != http.StatusOK {
				assertResponseError(t, err, tc.errCode, tc.errMsg)
				return
			}
			require.NoError(t, err)
			var expected []readable.SpentOutput
			checkGoldenFile(t, tc.golden, TestData{ux, &expected})
		})
	}
}

func TestLiveAddressUxOuts(t *testing.T) {
	if !doLive(t) {
		return
	}

	c := newClient()

	cases := []struct {
		name         string
		errCode      int
		errMsg       string
		addr         string
		moreThanZero bool
	}{
		{
			name:    "no addresses",
			errCode: http.StatusBadRequest,
			errMsg:  "400 Bad Request - address is empty",
		},
		{
			name:    "invalid address length",
			errCode: http.StatusBadRequest,
			errMsg:  "400 Bad Request - Invalid address length",
			addr:    "prRXwTcDK24hs6AFxj",
		},
		{
			name: "unknown address",
			addr: "prRXwTcDK24hs6AFxj69UuWae3LzhrsPW9",
		},
		{
			name: "one address",
			addr: "2THDupTBEo7UqB6dsVizkYUvkKq82Qn4gjf",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ux, err := c.AddressUxOuts(tc.addr)
			if tc.errCode != 0 && tc.errCode != http.StatusOK {
				assertResponseError(t, err, tc.errCode, tc.errMsg)
				return
			}
			require.NoError(t, err)
			if tc.moreThanZero {
				require.NotEqual(t, 0, len(ux))
			}
		})
	}
}

func TestStableBlocksInRange(t *testing.T) {
	if !doStable(t) {
		return
	}

	c := newClient()

	progress, err := c.BlockchainProgress()
	require.NoError(t, err)

	var lastNBlocks uint64 = 10
	require.True(t, progress.Current > lastNBlocks+1)

	cases := []struct {
		name    string
		golden  string
		start   uint64
		end     uint64
		errCode int
		errMsg  string
	}{
		{
			name:   "first 10",
			golden: "blocks-first-10.golden",
			start:  1,
			end:    10,
		},
		{
			name:   "last 10",
			golden: "blocks-last-10.golden",
			start:  progress.Current - lastNBlocks,
			end:    progress.Current,
		},
		{
			name:   "first block",
			golden: "blocks-first-1.golden",
			start:  1,
			end:    1,
		},
		{
			name:   "all blocks",
			golden: "blocks-all.golden",
			start:  0,
			end:    progress.Current,
		},
		{
			name:   "start > end",
			golden: "blocks-end-less-than-start.golden",
			start:  10,
			end:    9,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.errMsg == "" {
				resp := testBlocksInRange(t, tc.start, tc.end)

				var expected readable.Blocks
				checkGoldenFile(t, tc.golden, TestData{*resp, &expected})
			} else {
				_, err := c.BlocksInRange(tc.start, tc.end)
				assertResponseError(t, err, tc.errCode, tc.errMsg)
			}
		})
	}
}

func TestLiveBlocksInRange(t *testing.T) {
	if !doLive(t) {
		return
	}

	testBlocksInRange(t, 1, 10)
}

func testBlocksInRange(t *testing.T, start, end uint64) *readable.Blocks {
	c := newClient()

	blocks, err := c.BlocksInRange(start, end)
	require.NoError(t, err)

	if start > end {
		require.Empty(t, blocks.Blocks)
	} else {
		require.Len(t, blocks.Blocks, int(end-start+1))
	}

	var prevBlock *readable.Block
	for idx, b := range blocks.Blocks {
		if prevBlock != nil {
			require.Equal(t, prevBlock.Head.Hash, b.Head.PreviousHash)
		}

		bHash, err := c.BlockByHash(b.Head.Hash)
		require.Equal(t, uint64(idx)+start, b.Head.BkSeq)
		require.NoError(t, err)
		require.NotNil(t, bHash)
		require.Equal(t, b, *bHash)

		prevBlock = &blocks.Blocks[idx]
	}

	return blocks
}

func TestStableBlocksInRangeVerbose(t *testing.T) {
	if !doStable(t) {
		return
	}

	c := newClient()

	progress, err := c.BlockchainProgress()
	require.NoError(t, err)

	var lastNBlocks uint64 = 10
	require.True(t, progress.Current > lastNBlocks+1)

	cases := []struct {
		name    string
		golden  string
		start   uint64
		end     uint64
		errCode int
		errMsg  string
	}{
		{
			name:   "genesis",
			golden: "blocks-verbose-genesis.golden",
			start:  0,
			end:    0,
		},
		{
			name:   "first 10",
			golden: "blocks-verbose-first-10.golden",
			start:  1,
			end:    10,
		},
		{
			name:   "last 10",
			golden: "blocks-verbose-last-10.golden",
			start:  progress.Current - lastNBlocks,
			end:    progress.Current,
		},
		{
			name:   "first block",
			golden: "blocks-verbose-first-1.golden",
			start:  1,
			end:    1,
		},
		{
			name:   "all blocks",
			golden: "blocks-verbose-all.golden",
			start:  0,
			end:    progress.Current,
		},
		{
			name:   "start > end",
			golden: "blocks-verbose-end-less-than-start.golden",
			start:  10,
			end:    9,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.errMsg == "" {
				resp := testBlocksInRangeVerbose(t, tc.start, tc.end)

				var expected readable.BlocksVerbose
				checkGoldenFile(t, tc.golden, TestData{*resp, &expected})
			} else {
				blocks, err := c.BlocksInRangeVerbose(tc.start, tc.end)
				require.Nil(t, blocks)
				assertResponseError(t, err, tc.errCode, tc.errMsg)
			}
		})
	}
}

func TestLiveBlocksInRangeVerbose(t *testing.T) {
	if !doLive(t) {
		return
	}

	testBlocksInRangeVerbose(t, 1, 10)
}

func testBlocksInRangeVerbose(t *testing.T, start, end uint64) *readable.BlocksVerbose {
	c := newClient()

	blocks, err := c.BlocksInRangeVerbose(start, end)
	require.NoError(t, err)

	if start > end {
		require.Empty(t, blocks.Blocks)
	} else {
		require.Len(t, blocks.Blocks, int(end-start+1))
	}

	var prevBlock *readable.BlockVerbose
	for idx, b := range blocks.Blocks {
		assertVerboseBlockFee(t, &b)

		if prevBlock != nil {
			require.Equal(t, prevBlock.Head.Hash, b.Head.PreviousHash)
		}

		bHash, err := c.BlockByHashVerbose(b.Head.Hash)
		require.Equal(t, uint64(idx)+start, b.Head.BkSeq)
		require.NoError(t, err)
		require.NotNil(t, bHash)
		require.Equal(t, b, *bHash)

		prevBlock = &blocks.Blocks[idx]
	}

	return blocks
}

func TestStableBlocks(t *testing.T) {
	if !doStable(t) {
		return
	}

	c := newClient()

	cases := []struct {
		name    string
		golden  string
		seqs    []uint64
		errCode int
		errMsg  string
	}{
		{
			name:   "multiple sequences",
			golden: "blocks-3-5-7.golden",
			seqs:   []uint64{3, 5, 7},
		},
		{
			name:    "block seq not found",
			seqs:    []uint64{3, 5, 7, 99999},
			errCode: http.StatusNotFound,
			errMsg:  "404 Not Found - block does not exist seq=99999",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.errMsg == "" {
				resp := testBlocks(t, tc.seqs)

				var expected readable.Blocks
				checkGoldenFile(t, tc.golden, TestData{*resp, &expected})
			} else {
				_, err := c.Blocks(tc.seqs)
				assertResponseError(t, err, tc.errCode, tc.errMsg)
			}
		})
	}
}

func TestLiveBlocks(t *testing.T) {
	if !doLive(t) {
		return
	}

	testBlocks(t, []uint64{3, 5, 7})
}

func testBlocks(t *testing.T, seqs []uint64) *readable.Blocks {
	c := newClient()

	blocks, err := c.Blocks(seqs)
	require.NoError(t, err)

	require.Equal(t, len(seqs), len(blocks.Blocks))

	seqsMap := make(map[uint64]struct{}, len(seqs))
	for _, x := range seqs {
		seqsMap[x] = struct{}{}
	}

	for _, b := range blocks.Blocks {
		_, ok := seqsMap[b.Head.BkSeq]
		require.True(t, ok)
		delete(seqsMap, b.Head.BkSeq)

		bHash, err := c.BlockByHash(b.Head.Hash)
		require.NoError(t, err)
		require.NotNil(t, bHash)
		require.Equal(t, b, *bHash)
	}

	require.Empty(t, seqsMap)

	return blocks
}

func TestStableBlocksVerbose(t *testing.T) {
	if !doStable(t) {
		return
	}

	c := newClient()

	cases := []struct {
		name    string
		golden  string
		seqs    []uint64
		errCode int
		errMsg  string
	}{
		{
			name:   "multiple sequences",
			golden: "blocks-3-5-7-verbose.golden",
			seqs:   []uint64{3, 5, 7},
		},
		{
			name:    "block seq not found",
			seqs:    []uint64{3, 5, 7, 99999},
			errCode: http.StatusNotFound,
			errMsg:  "404 Not Found - block does not exist seq=99999",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.errMsg == "" {
				resp := testBlocksVerbose(t, tc.seqs)

				var expected readable.BlocksVerbose
				checkGoldenFile(t, tc.golden, TestData{*resp, &expected})
			} else {
				blocks, err := c.BlocksVerbose(tc.seqs)
				require.Nil(t, blocks)
				assertResponseError(t, err, tc.errCode, tc.errMsg)
			}
		})
	}
}

func TestLiveBlocksVerbose(t *testing.T) {
	if !doLive(t) {
		return
	}

	testBlocksVerbose(t, []uint64{3, 5, 7})
}

func testBlocksVerbose(t *testing.T, seqs []uint64) *readable.BlocksVerbose {
	c := newClient()

	blocks, err := c.BlocksVerbose(seqs)
	require.NoError(t, err)

	require.Equal(t, len(seqs), len(blocks.Blocks))

	seqsMap := make(map[uint64]struct{}, len(seqs))
	for _, x := range seqs {
		seqsMap[x] = struct{}{}
	}

	for _, b := range blocks.Blocks {
		_, ok := seqsMap[b.Head.BkSeq]
		require.True(t, ok)
		delete(seqsMap, b.Head.BkSeq)

		assertVerboseBlockFee(t, &b)

		bHash, err := c.BlockByHashVerbose(b.Head.Hash)
		require.NoError(t, err)
		require.NotNil(t, bHash)
		require.Equal(t, b, *bHash)
	}

	require.Empty(t, seqsMap)

	return blocks
}

func TestStableLastBlocks(t *testing.T) {
	if !doStable(t) {
		return
	}

	c := newClient()

	blocks, err := c.LastBlocks(1)
	require.NoError(t, err)

	var expected *readable.Blocks
	checkGoldenFile(t, "block-last.golden", TestData{blocks, &expected})

	var prevBlock *readable.Block
	blocks, err = c.LastBlocks(10)
	require.NoError(t, err)
	require.Equal(t, 10, len(blocks.Blocks))
	for idx, b := range blocks.Blocks {
		if prevBlock != nil {
			require.Equal(t, prevBlock.Head.Hash, b.Head.PreviousHash)
		}

		bHash, err := c.BlockByHash(b.Head.Hash)
		require.NoError(t, err)
		require.NotNil(t, bHash)
		require.Equal(t, b, *bHash)

		prevBlock = &blocks.Blocks[idx]
	}
}

func TestLiveLastBlocks(t *testing.T) {
	if !doLive(t) {
		return
	}
	c := newClient()
	var prevBlock *readable.Block
	blocks, err := c.LastBlocks(10)
	require.NoError(t, err)
	require.Equal(t, 10, len(blocks.Blocks))
	for idx, b := range blocks.Blocks {
		if prevBlock != nil {
			require.Equal(t, prevBlock.Head.Hash, b.Head.PreviousHash)
		}

		bHash, err := c.BlockByHash(b.Head.Hash)
		require.NoError(t, err)
		require.NotNil(t, bHash)
		require.Equal(t, b, *bHash)

		prevBlock = &blocks.Blocks[idx]
	}
}

func TestStableLastBlocksVerbose(t *testing.T) {
	if !doStable(t) {
		return
	}

	c := newClient()

	blocks, err := c.LastBlocksVerbose(1)
	require.NoError(t, err)

	var expected *readable.BlocksVerbose
	checkGoldenFile(t, "block-last-verbose.golden", TestData{blocks, &expected})

	blocks, err = c.LastBlocksVerbose(10)
	require.NoError(t, err)
	require.Equal(t, 10, len(blocks.Blocks))

	var prevBlock *readable.BlockVerbose
	for idx, b := range blocks.Blocks {
		assertVerboseBlockFee(t, &b)

		if prevBlock != nil {
			require.Equal(t, prevBlock.Head.Hash, b.Head.PreviousHash)
		}

		bHash, err := c.BlockByHashVerbose(b.Head.Hash)
		require.NoError(t, err)
		require.NotNil(t, bHash)
		require.Equal(t, b, *bHash)

		prevBlock = &blocks.Blocks[idx]
	}
}

func TestLiveLastBlocksVerbose(t *testing.T) {
	if !doLive(t) {
		return
	}
	c := newClient()

	blocks, err := c.LastBlocksVerbose(10)
	require.NoError(t, err)
	require.Equal(t, 10, len(blocks.Blocks))

	var prevBlock *readable.BlockVerbose
	for idx, b := range blocks.Blocks {
		assertVerboseBlockFee(t, &b)

		if prevBlock != nil {
			require.Equal(t, prevBlock.Head.Hash, b.Head.PreviousHash)
		}

		bHash, err := c.BlockByHashVerbose(b.Head.Hash)
		require.NoError(t, err)
		require.NotNil(t, bHash)
		require.Equal(t, b, *bHash)

		prevBlock = &blocks.Blocks[idx]
	}
}

func TestStableNetworkConnections(t *testing.T) {
	if !doStable(t) {
		return
	}

	c := newClient()
	connections, err := c.NetworkConnections(nil)
	require.NoError(t, err)
	require.Empty(t, connections.Connections)

	connection, err := c.NetworkConnection("127.0.0.1:4444")
	assertResponseError(t, err, http.StatusNotFound, "404 Not Found")
	require.Nil(t, connection)
}

func TestLiveNetworkConnections(t *testing.T) {
	if !doLive(t) {
		return
	}

	c := newClient()
	connections, err := c.NetworkConnections(nil)
	require.NoError(t, err)

	if liveDisableNetworking(t) {
		require.Empty(t, connections.Connections)
		return
	}

	require.NotEmpty(t, connections.Connections)

	checked := false

	for _, cc := range connections.Connections {
		connection, err := c.NetworkConnection(cc.Addr)

		// The connection may have disconnected by now
		if err != nil {
			assertResponseError(t, err, http.StatusNotFound, "404 Not Found")
			continue
		}

		require.NoError(t, err)
		require.NotEmpty(t, cc.Addr)
		require.Equal(t, cc.Addr, connection.Addr)
		require.Equal(t, cc.GnetID, connection.GnetID)
		require.Equal(t, cc.ListenPort, connection.ListenPort)
		require.Equal(t, cc.Mirror, connection.Mirror)

		switch cc.State {
		case daemon.ConnectionStateIntroduced:
			// If the connection was introduced it should stay introduced
			require.Equal(t, daemon.ConnectionStateIntroduced, connection.State)
		case daemon.ConnectionStateConnected:
			// If the connection was connected it should stay connected or have become introduced
			require.NotEqual(t, daemon.ConnectionStatePending, connection.State)
		}

		// The GnetID should be 0 if pending, otherwise it should not be 0
		if cc.State == daemon.ConnectionStatePending {
			require.Equal(t, uint64(0), cc.GnetID)
		} else {
			require.NotEmpty(t, cc.GnetID)
		}

		require.Equal(t, cc.Outgoing, connection.Outgoing)
		require.True(t, cc.LastReceived <= connection.LastReceived)
		require.True(t, cc.LastSent <= connection.LastSent)
		require.Equal(t, cc.ConnectedAt, connection.ConnectedAt)

		checked = true
	}

	// This could unfortunately occur if a connection disappeared in between the two calls,
	// which will require a test re-run.
	require.True(t, checked, "Was not able to find any connection by address, despite finding connections when querying all")

	connections, err = c.NetworkConnections(&api.NetworkConnectionsFilter{
		States: []daemon.ConnectionState{daemon.ConnectionStatePending},
	})
	require.NoError(t, err)

	for _, cc := range connections.Connections {
		require.Equal(t, daemon.ConnectionStatePending, cc.State)
	}

	connections, err = c.NetworkConnections(&api.NetworkConnectionsFilter{
		Direction: "incoming",
	})
	require.NoError(t, err)

	for _, cc := range connections.Connections {
		require.False(t, cc.Outgoing)
	}
}

func TestNetworkDefaultConnections(t *testing.T) {
	if !doLiveOrStable(t) {
		return
	}

	c := newClient()
	connections, err := c.NetworkDefaultPeers()
	require.NoError(t, err)
	require.NotEmpty(t, connections)
	sort.Strings(connections)

	var expected []string
	checkGoldenFile(t, "network-default-peers.golden", TestData{connections, &expected})
}

func TestNetworkTrustedConnections(t *testing.T) {
	if !doLiveOrStable(t) {
		return
	}

	c := newClient()
	connections, err := c.NetworkTrustedPeers()
	require.NoError(t, err)
	require.NotEmpty(t, connections)
	sort.Strings(connections)

	var expected []string
	checkGoldenFile(t, "network-trusted-peers.golden", TestData{connections, &expected})
}

func TestStableNetworkExchangeableConnections(t *testing.T) {
	if !doStable(t) {
		return
	}

	c := newClient()
	connections, err := c.NetworkExchangedPeers()
	require.NoError(t, err)

	var expected []string
	checkGoldenFile(t, "network-exchanged-peers.golden", TestData{connections, &expected})
}

func TestLiveNetworkExchangeableConnections(t *testing.T) {
	if !doLive(t) {
		return
	}

	c := newClient()
	_, err := c.NetworkExchangedPeers()
	require.NoError(t, err)
}

type transactionTestCase struct {
	name       string
	txID       string
	err        api.ClientError
	goldenFile string
}

func TestLiveTransaction(t *testing.T) {
	if !doLive(t) {
		return
	}

	cases := []transactionTestCase{
		{
			name: "invalid txID",
			txID: "abcd",
			err: api.ClientError{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Message:    "400 Bad Request - Invalid hex length",
			},
		},
		{
			name: "empty txID",
			txID: "",
			err: api.ClientError{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Message:    "400 Bad Request - txid is empty",
			},
		},
		{
			name:       "OK",
			txID:       "76ecbabc53ea2a3be46983058433dda6a3cf7ea0b86ba14d90b932fa97385de7",
			goldenFile: "transaction-block-517.golden",
		},
	}

	c := newClient()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tx, err := c.Transaction(tc.txID)
			if err != nil {
				require.Equal(t, tc.err, err)
				return
			}

			// tx.Status.Height is how many blocks are above this transaction,
			// make sure it is past some checkpoint height
			require.True(t, tx.Status.Height >= 50836)

			// readable.TransactionWithStatus.Status.Height is not stable
			tx.Status.Height = 0

			var expected readable.TransactionWithStatus
			loadGoldenFile(t, tc.goldenFile, TestData{tx, &expected})
			require.Equal(t, &expected, tx)
		})
	}
}

func TestStableTransaction(t *testing.T) {
	if !doStable(t) {
		return
	}

	cases := []transactionTestCase{
		{
			name: "invalid txId",
			txID: "abcd",
			err: api.ClientError{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Message:    "400 Bad Request - Invalid hex length",
			},
			goldenFile: "",
		},
		{
			name: "not exist",
			txID: "540582ee4128b733f810f149e908d984a5f403ad2865108e6c1c5423aeefc759",
			err: api.ClientError{
				Status:     "404 Not Found",
				StatusCode: http.StatusNotFound,
				Message:    "404 Not Found",
			},
			goldenFile: "",
		},
		{
			name: "empty txId",
			txID: "",
			err: api.ClientError{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Message:    "400 Bad Request - txid is empty",
			},
			goldenFile: "",
		},
		{
			name:       "genesis transaction",
			txID:       "d556c1c7abf1e86138316b8c17183665512dc67633c04cf236a8b7f332cb4add",
			goldenFile: "genesis-transaction.golden",
		},
		{
			name:       "transaction in block 101",
			txID:       "e8fe5290afba3933389fd5860dca2cbcc81821028be9c65d0bb7cf4e8d2c4c18",
			goldenFile: "transaction-block-101.golden",
		},
	}

	if !dbNoUnconfirmed(t) {
		cases = append(cases, transactionTestCase{
			name:       "unconfirmed",
			txID:       "701d23fd513bad325938ba56869f9faba19384a8ec3dd41833aff147eac53947",
			goldenFile: "transaction-unconfirmed.golden",
		})
	}

	c := newClient()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tx, err := c.Transaction(tc.txID)
			if err != nil {
				require.Equal(t, tc.err, err)
				return
			}

			var expected readable.TransactionWithStatus
			loadGoldenFile(t, tc.goldenFile, TestData{tx, &expected})
			require.Equal(t, &expected, tx)
		})
	}
}

func TestLiveTransactionVerbose(t *testing.T) {
	if !doLive(t) {
		return
	}

	cases := []transactionTestCase{
		{
			name: "invalid txID",
			txID: "abcd",
			err: api.ClientError{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Message:    "400 Bad Request - Invalid hex length",
			},
		},
		{
			name: "empty txID",
			txID: "",
			err: api.ClientError{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Message:    "400 Bad Request - txid is empty",
			},
		},
		{
			name:       "OK",
			txID:       "76ecbabc53ea2a3be46983058433dda6a3cf7ea0b86ba14d90b932fa97385de7",
			goldenFile: "transaction-verbose-block-517.golden",
		},
	}

	c := newClient()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tx, err := c.TransactionVerbose(tc.txID)
			if err != nil {
				require.Equal(t, tc.err, err)
				return
			}

			// tx.Status.Height is how many blocks are above this transaction,
			// make sure it is past some checkpoint height
			require.True(t, tx.Status.Height >= 50836)

			// readable.TransactionWithStatus.Status.Height is not stable
			tx.Status.Height = 0

			var expected readable.TransactionWithStatusVerbose
			loadGoldenFile(t, tc.goldenFile, TestData{tx, &expected})
			require.Equal(t, &expected, tx)
		})
	}
}

func TestStableTransactionVerbose(t *testing.T) {
	if !doStable(t) {
		return
	}

	cases := []transactionTestCase{
		{
			name: "invalid txId",
			txID: "abcd",
			err: api.ClientError{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Message:    "400 Bad Request - Invalid hex length",
			},
			goldenFile: "",
		},
		{
			name: "not exist",
			txID: "540582ee4128b733f810f149e908d984a5f403ad2865108e6c1c5423aeefc759",
			err: api.ClientError{
				Status:     "404 Not Found",
				StatusCode: http.StatusNotFound,
				Message:    "404 Not Found",
			},
			goldenFile: "",
		},
		{
			name: "empty txId",
			txID: "",
			err: api.ClientError{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Message:    "400 Bad Request - txid is empty",
			},
			goldenFile: "",
		},
		{
			name:       "genesis transaction",
			txID:       "d556c1c7abf1e86138316b8c17183665512dc67633c04cf236a8b7f332cb4add",
			goldenFile: "genesis-transaction-verbose.golden",
		},
		{
			name:       "transaction in block 101",
			txID:       "e8fe5290afba3933389fd5860dca2cbcc81821028be9c65d0bb7cf4e8d2c4c18",
			goldenFile: "transaction-verbose-block-101.golden",
		},
	}

	if !dbNoUnconfirmed(t) {
		cases = append(cases, transactionTestCase{
			name:       "unconfirmed",
			txID:       "701d23fd513bad325938ba56869f9faba19384a8ec3dd41833aff147eac53947",
			goldenFile: "transaction-unconfirmed-verbose.golden",
		})
	}

	c := newClient()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tx, err := c.TransactionVerbose(tc.txID)
			if err != nil {
				require.Equal(t, tc.err, err)
				return
			}

			var expected readable.TransactionWithStatusVerbose
			loadGoldenFile(t, tc.goldenFile, TestData{tx, &expected})
			require.Equal(t, &expected, tx)
		})
	}
}

func TestLiveTransactionEncoded(t *testing.T) {
	if !doLive(t) {
		return
	}

	cases := []transactionTestCase{
		{
			name: "invalid txID",
			txID: "abcd",
			err: api.ClientError{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Message:    "400 Bad Request - Invalid hex length",
			},
		},
		{
			name: "empty txID",
			txID: "",
			err: api.ClientError{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Message:    "400 Bad Request - txid is empty",
			},
		},
		{
			name:       "OK",
			txID:       "76ecbabc53ea2a3be46983058433dda6a3cf7ea0b86ba14d90b932fa97385de7",
			goldenFile: "transaction-encoded.golden",
		},
	}

	c := newClient()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			testTransactionEncoded(t, c, tc, false)
		})
	}
}

func TestStableTransactionEncoded(t *testing.T) {
	if !doStable(t) {
		return
	}

	cases := []transactionTestCase{
		{
			name: "invalid txId",
			txID: "abcd",
			err: api.ClientError{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Message:    "400 Bad Request - Invalid hex length",
			},
			goldenFile: "",
		},
		{
			name: "not exist",
			txID: "540582ee4128b733f810f149e908d984a5f403ad2865108e6c1c5423aeefc759",
			err: api.ClientError{
				Status:     "404 Not Found",
				StatusCode: http.StatusNotFound,
				Message:    "404 Not Found",
			},
			goldenFile: "",
		},
		{
			name: "empty txId",
			txID: "",
			err: api.ClientError{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Message:    "400 Bad Request - txid is empty",
			},
			goldenFile: "",
		},
		{
			name:       "genesis transaction",
			txID:       "d556c1c7abf1e86138316b8c17183665512dc67633c04cf236a8b7f332cb4add",
			goldenFile: "genesis-transaction-encoded.golden",
		},
		{
			name:       "transaction in block 101",
			txID:       "e8fe5290afba3933389fd5860dca2cbcc81821028be9c65d0bb7cf4e8d2c4c18",
			goldenFile: "transaction-encoded-block-101.golden",
		},
		{
			name:       "transaction in block 105",
			txID:       "41ec724bd40c852096379d1ae57d3f27606877fa95ac9c082fbf63900e6c5cb5",
			goldenFile: "transaction-encoded-block-105.golden",
		},
	}

	if !dbNoUnconfirmed(t) {
		cases = append(cases, transactionTestCase{
			name:       "unconfirmed",
			txID:       "701d23fd513bad325938ba56869f9faba19384a8ec3dd41833aff147eac53947",
			goldenFile: "transaction-unconfirmed-encoded.golden",
		})
	}

	c := newClient()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			testTransactionEncoded(t, c, tc, true)
		})
	}
}

func testTransactionEncoded(t *testing.T, c *api.Client, tc transactionTestCase, stable bool) {
	encodedTxn, err := c.TransactionEncoded(tc.txID)
	if err != nil {
		require.Equal(t, tc.err, err)
		return
	}

	if !stable {
		encodedTxn.Status.Height = 0
	}

	decodedTxn, err := coin.DeserializeTransactionHex(encodedTxn.EncodedTransaction)
	require.NoError(t, err)

	txnResult, err := readable.NewTransactionWithStatus(&visor.Transaction{
		Transaction: decodedTxn,
		Status: visor.TransactionStatus{
			Confirmed: encodedTxn.Status.Confirmed,
			Height:    encodedTxn.Status.Height,
			BlockSeq:  encodedTxn.Status.BlockSeq,
		},
		Time: encodedTxn.Time,
	})
	require.NoError(t, err)

	txn, err := c.Transaction(tc.txID)
	require.NoError(t, err)

	if !stable {
		txn.Status.Height = 0
	}

	require.Equal(t, txn, txnResult)

	var expected api.TransactionEncodedResponse
	loadGoldenFile(t, tc.goldenFile, TestData{encodedTxn, &expected})
	require.Equal(t, &expected, encodedTxn)
}

func TestLiveTransactions(t *testing.T) {
	if !doLive(t) {
		return
	}

	c := newClient()
	addrs := []string{
		"2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKt",
	}
	txns, err := c.Transactions(addrs)
	require.NoError(t, err)
	require.True(t, len(txns) > 0)
	assertNoTransactionsDupes(t, txns)

	// Two addresses with a mutual transaction between the two, to test deduplication
	addrs = []string{
		"7cpQ7t3PZZXvjTst8G7Uvs7XH4LeM8fBPD",
		"2K6NuLBBapWndAssUtkxKfCtyjDQDHrEhhT",
	}
	txns, err = c.Transactions(addrs)
	require.NoError(t, err)
	// There were 4 transactions amonst these two addresses at the time this was written
	require.True(t, len(txns) >= 4)
	assertNoTransactionsDupes(t, txns)
}

type transactionsTestCase struct {
	name       string
	addrs      []string
	err        api.ClientError
	goldenFile string
}

func TestStableTransactions(t *testing.T) {
	if !doStable(t) {
		return
	}

	cases := []transactionsTestCase{
		{
			name:  "invalid addr length",
			addrs: []string{"abcd"},
			err: api.ClientError{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Message:    "400 Bad Request - parse parameter: 'addrs' failed: address \"abcd\" is invalid: Invalid address length",
			},
		},
		{
			name:  "invalid addr character",
			addrs: []string{"701d23fd513bad325938ba56869f9faba19384a8ec3dd41833aff147eac53947"},
			err: api.ClientError{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Message:    "400 Bad Request - parse parameter: 'addrs' failed: address \"701d23fd513bad325938ba56869f9faba19384a8ec3dd41833aff147eac53947\" is invalid: Invalid base58 character",
			},
		},
		{
			name:  "invalid checksum",
			addrs: []string{"2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKk"},
			err: api.ClientError{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Message:    "400 Bad Request - parse parameter: 'addrs' failed: address \"2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKk\" is invalid: Invalid checksum",
			},
		},
		{
			name:  "empty addrs",
			addrs: []string{},
			err: api.ClientError{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Message:    "400 Bad Request - txId is empty",
			},
			goldenFile: "empty-addrs-transactions.golden",
		},
		{
			name:       "single addr",
			addrs:      []string{"2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKt"},
			goldenFile: "single-addr-transactions.golden",
		},
		{
			name:       "genesis",
			addrs:      []string{"2jBbGxZRGoQG1mqhPBnXnLTxK6oxsTf8os6"},
			goldenFile: "genesis-addr-transactions.golden",
		},
		{
			name:       "multiple addrs",
			addrs:      []string{"2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKt", "2JJ8pgq8EDAnrzf9xxBJapE2qkYLefW4uF8"},
			goldenFile: "multiple-addr-transactions.golden",
		},
	}

	if !dbNoUnconfirmed(t) {
		cases = append(cases, transactionsTestCase{
			name:       "confirmed and unconfirmed transactions",
			addrs:      []string{"212mwY3Dmey6vwnWpiph99zzCmopXTqeVEN"},
			goldenFile: "confirmed-and-unconfirmed-transactions.golden",
		})
	}

	c := newClient()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			txnResult, err := c.Transactions(tc.addrs)
			if err != nil {
				require.Equal(t, tc.err, err, "case: "+tc.name)
				return
			}

			assertNoTransactionsDupes(t, txnResult)

			var expected []readable.TransactionWithStatus
			checkGoldenFile(t, tc.goldenFile, TestData{txnResult, &expected})
		})
	}
}

func TestLiveConfirmedTransactions(t *testing.T) {
	if !doLive(t) {
		return
	}

	c := newClient()

	cTxsSingle, err := c.ConfirmedTransactions([]string{"2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKt"})
	require.NoError(t, err)
	require.True(t, len(cTxsSingle) > 0)
	assertNoTransactionsDupes(t, cTxsSingle)

	cTxsAll, err := c.ConfirmedTransactions([]string{})
	require.NoError(t, err)
	require.True(t, len(cTxsAll) > 0)
	require.True(t, len(cTxsAll) > len(cTxsSingle))
	assertNoTransactionsDupes(t, cTxsAll)
}

func TestStableConfirmedTransactions(t *testing.T) {
	if !doStable(t) {
		return
	}

	cases := []transactionsTestCase{
		{
			name:  "invalid addr length",
			addrs: []string{"abcd"},
			err: api.ClientError{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Message:    "400 Bad Request - parse parameter: 'addrs' failed: address \"abcd\" is invalid: Invalid address length",
			},
		},
		{
			name:  "invalid addr character",
			addrs: []string{"701d23fd513bad325938ba56869f9faba19384a8ec3dd41833aff147eac53947"},
			err: api.ClientError{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Message:    "400 Bad Request - parse parameter: 'addrs' failed: address \"701d23fd513bad325938ba56869f9faba19384a8ec3dd41833aff147eac53947\" is invalid: Invalid base58 character",
			},
		},
		{
			name:  "invalid checksum",
			addrs: []string{"2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKk"},
			err: api.ClientError{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Message:    "400 Bad Request - parse parameter: 'addrs' failed: address \"2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKk\" is invalid: Invalid checksum",
			},
		},
		{
			name:       "empty addrs",
			addrs:      []string{},
			goldenFile: "empty-addrs-transactions.golden",
		},
		{
			name:       "single addr",
			addrs:      []string{"2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKt"},
			goldenFile: "single-addr-transactions.golden",
		},
		{
			name:       "genesis",
			addrs:      []string{"2jBbGxZRGoQG1mqhPBnXnLTxK6oxsTf8os6"},
			goldenFile: "genesis-addr-transactions.golden",
		},
		{
			name:       "multiple addrs",
			addrs:      []string{"2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKt", "2JJ8pgq8EDAnrzf9xxBJapE2qkYLefW4uF8"},
			goldenFile: "multiple-addr-transactions.golden",
		},
		{
			name:       "unconfirmed should be excluded",
			addrs:      []string{"212mwY3Dmey6vwnWpiph99zzCmopXTqeVEN"},
			goldenFile: "unconfirmed-excluded-from-transactions.golden",
		},
	}

	c := newClient()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			txnResult, err := c.ConfirmedTransactions(tc.addrs)
			if err != nil {
				require.Equal(t, tc.err, err, "case: "+tc.name)
				return
			}

			assertNoTransactionsDupes(t, txnResult)

			var expected []readable.TransactionWithStatus
			checkGoldenFile(t, tc.goldenFile, TestData{txnResult, &expected})
		})
	}
}

func TestStableUnconfirmedTransactions(t *testing.T) {
	if !doStable(t) {
		return
	}

	cases := []transactionsTestCase{
		{
			name:  "invalid addr length",
			addrs: []string{"abcd"},
			err: api.ClientError{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Message:    "400 Bad Request - parse parameter: 'addrs' failed: address \"abcd\" is invalid: Invalid address length",
			},
		},
		{
			name:  "invalid addr character",
			addrs: []string{"701d23fd513bad325938ba56869f9faba19384a8ec3dd41833aff147eac53947"},
			err: api.ClientError{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Message:    "400 Bad Request - parse parameter: 'addrs' failed: address \"701d23fd513bad325938ba56869f9faba19384a8ec3dd41833aff147eac53947\" is invalid: Invalid base58 character",
			},
		},
		{
			name:  "invalid checksum",
			addrs: []string{"2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKk"},
			err: api.ClientError{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Message:    "400 Bad Request - parse parameter: 'addrs' failed: address \"2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKk\" is invalid: Invalid checksum",
			},
		},
	}

	if dbNoUnconfirmed(t) {
		cases = append(cases, transactionsTestCase{
			name:       "empty addrs",
			addrs:      []string{},
			goldenFile: "no-unconfirmed-txns.golden",
		})
	} else {
		cases = append(cases, transactionsTestCase{
			name:       "empty addrs (all unconfirmed txns)",
			addrs:      []string{},
			goldenFile: "all-unconfirmed-txns.golden",
		})
	}

	c := newClient()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			txnResult, err := c.UnconfirmedTransactions(tc.addrs)
			if err != nil {
				require.Equal(t, tc.err, err, "case: "+tc.name)
				return
			}

			assertNoTransactionsDupes(t, txnResult)

			var expected []readable.TransactionWithStatus
			checkGoldenFile(t, tc.goldenFile, TestData{txnResult, &expected})
		})
	}
}

func TestLiveUnconfirmedTransactions(t *testing.T) {
	if !doLive(t) {
		return
	}

	c := newClient()

	cTxsSingle, err := c.UnconfirmedTransactions([]string{"2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKt"})
	require.NoError(t, err)
	require.True(t, len(cTxsSingle) >= 0)
	assertNoTransactionsDupes(t, cTxsSingle)

	cTxsAll, err := c.UnconfirmedTransactions([]string{})
	require.NoError(t, err)
	require.True(t, len(cTxsAll) >= 0)
	require.True(t, len(cTxsAll) >= len(cTxsSingle))
	assertNoTransactionsDupes(t, cTxsAll)
}

func assertNoTransactionsDupes(t *testing.T, r []readable.TransactionWithStatus) {
	txids := make(map[string]struct{})

	for _, x := range r {
		_, ok := txids[x.Transaction.Hash]
		require.False(t, ok)
		txids[x.Transaction.Hash] = struct{}{}
	}
}

func TestLiveTransactionsVerbose(t *testing.T) {
	if !doLive(t) {
		return
	}

	c := newClient()
	addrs := []string{
		"2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKt",
	}
	txns, err := c.TransactionsVerbose(addrs)
	require.NoError(t, err)
	require.True(t, len(txns) > 0)
	assertNoTransactionsDupesVerbose(t, txns)
}

func TestStableTransactionsVerbose(t *testing.T) {
	if !doStable(t) {
		return
	}

	cases := []transactionsTestCase{
		{
			name:  "invalid addr length",
			addrs: []string{"abcd"},
			err: api.ClientError{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Message:    "400 Bad Request - parse parameter: 'addrs' failed: address \"abcd\" is invalid: Invalid address length",
			},
		},
		{
			name:  "invalid addr character",
			addrs: []string{"701d23fd513bad325938ba56869f9faba19384a8ec3dd41833aff147eac53947"},
			err: api.ClientError{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Message:    "400 Bad Request - parse parameter: 'addrs' failed: address \"701d23fd513bad325938ba56869f9faba19384a8ec3dd41833aff147eac53947\" is invalid: Invalid base58 character",
			},
		},
		{
			name:  "invalid checksum",
			addrs: []string{"2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKk"},
			err: api.ClientError{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Message:    "400 Bad Request - parse parameter: 'addrs' failed: address \"2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKk\" is invalid: Invalid checksum",
			},
		},
		{
			name:       "empty addrs",
			addrs:      []string{},
			goldenFile: "empty-addrs-transactions-verbose.golden",
		},
		{
			name:       "single addr",
			addrs:      []string{"2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKt"},
			goldenFile: "single-addr-transactions-verbose.golden",
		},
		{
			name:       "genesis",
			addrs:      []string{"2jBbGxZRGoQG1mqhPBnXnLTxK6oxsTf8os6"},
			goldenFile: "genesis-addr-transactions-verbose.golden",
		},
		{
			name:       "multiple addrs",
			addrs:      []string{"2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKt", "2JJ8pgq8EDAnrzf9xxBJapE2qkYLefW4uF8"},
			goldenFile: "multiple-addr-transactions-verbose.golden",
		},
	}

	if !dbNoUnconfirmed(t) {
		cases = append(cases, transactionsTestCase{
			name:       "confirmed and unconfirmed transactions",
			addrs:      []string{"212mwY3Dmey6vwnWpiph99zzCmopXTqeVEN"},
			goldenFile: "confirmed-and-unconfirmed-transactions-verbose.golden",
		})
	}

	c := newClient()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			txnResult, err := c.TransactionsVerbose(tc.addrs)
			if err != nil {
				require.Equal(t, tc.err, err, "case: "+tc.name)
				return
			}

			assertNoTransactionsDupesVerbose(t, txnResult)

			var expected []readable.TransactionWithStatusVerbose
			checkGoldenFile(t, tc.goldenFile, TestData{txnResult, &expected})
		})
	}
}

func TestLiveConfirmedTransactionsVerbose(t *testing.T) {
	if !doLive(t) {
		return
	}

	c := newClient()

	cTxsSingle, err := c.ConfirmedTransactionsVerbose([]string{"2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKt"})
	require.NoError(t, err)
	require.True(t, len(cTxsSingle) > 0)
	assertNoTransactionsDupesVerbose(t, cTxsSingle)

	cTxsAll, err := c.ConfirmedTransactionsVerbose([]string{})
	require.NoError(t, err)
	require.True(t, len(cTxsAll) > 0)
	require.True(t, len(cTxsAll) > len(cTxsSingle))
	assertNoTransactionsDupesVerbose(t, cTxsAll)
}

func TestStableConfirmedTransactionsVerbose(t *testing.T) {
	if !doStable(t) {
		return
	}

	cases := []transactionsTestCase{
		{
			name:  "invalid addr length",
			addrs: []string{"abcd"},
			err: api.ClientError{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Message:    "400 Bad Request - parse parameter: 'addrs' failed: address \"abcd\" is invalid: Invalid address length",
			},
		},
		{
			name:  "invalid addr character",
			addrs: []string{"701d23fd513bad325938ba56869f9faba19384a8ec3dd41833aff147eac53947"},
			err: api.ClientError{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Message:    "400 Bad Request - parse parameter: 'addrs' failed: address \"701d23fd513bad325938ba56869f9faba19384a8ec3dd41833aff147eac53947\" is invalid: Invalid base58 character",
			},
		},
		{
			name:  "invalid checksum",
			addrs: []string{"2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKk"},
			err: api.ClientError{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Message:    "400 Bad Request - parse parameter: 'addrs' failed: address \"2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKk\" is invalid: Invalid checksum",
			},
		},
		{
			name:       "empty addrs",
			addrs:      []string{},
			goldenFile: "empty-addrs-transactions-verbose.golden",
		},
		{
			name:       "single addr",
			addrs:      []string{"2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKt"},
			goldenFile: "single-addr-transactions-verbose.golden",
		},
		{
			name:       "genesis",
			addrs:      []string{"2jBbGxZRGoQG1mqhPBnXnLTxK6oxsTf8os6"},
			goldenFile: "genesis-addr-transactions-verbose.golden",
		},
		{
			name:       "multiple addrs",
			addrs:      []string{"2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKt", "2JJ8pgq8EDAnrzf9xxBJapE2qkYLefW4uF8"},
			goldenFile: "multiple-addr-transactions-verbose.golden",
		},
		{
			name:       "unconfirmed should be excluded",
			addrs:      []string{"212mwY3Dmey6vwnWpiph99zzCmopXTqeVEN"},
			goldenFile: "unconfirmed-excluded-from-transactions-verbose.golden",
		},
	}

	c := newClient()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			txnResult, err := c.ConfirmedTransactionsVerbose(tc.addrs)
			if err != nil {
				require.Equal(t, tc.err, err, "case: "+tc.name)
				return
			}

			assertNoTransactionsDupesVerbose(t, txnResult)

			var expected []readable.TransactionWithStatusVerbose
			checkGoldenFile(t, tc.goldenFile, TestData{txnResult, &expected})
		})
	}
}

func TestStableUnconfirmedTransactionsVerbose(t *testing.T) {
	if !doStable(t) {
		return
	}

	cases := []transactionsTestCase{
		{
			name:  "invalid addr length",
			addrs: []string{"abcd"},
			err: api.ClientError{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Message:    "400 Bad Request - parse parameter: 'addrs' failed: address \"abcd\" is invalid: Invalid address length",
			},
		},
		{
			name:  "invalid addr character",
			addrs: []string{"701d23fd513bad325938ba56869f9faba19384a8ec3dd41833aff147eac53947"},
			err: api.ClientError{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Message:    "400 Bad Request - parse parameter: 'addrs' failed: address \"701d23fd513bad325938ba56869f9faba19384a8ec3dd41833aff147eac53947\" is invalid: Invalid base58 character",
			},
		},
		{
			name:  "invalid checksum",
			addrs: []string{"2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKk"},
			err: api.ClientError{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Message:    "400 Bad Request - parse parameter: 'addrs' failed: address \"2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKk\" is invalid: Invalid checksum",
			},
		},
	}

	if dbNoUnconfirmed(t) {
		cases = append(cases, transactionsTestCase{
			name:       "empty addrs",
			addrs:      []string{},
			goldenFile: "no-unconfirmed-txns.golden",
		})
	} else {
		cases = append(cases, transactionsTestCase{
			name:       "empty addrs (all unconfirmed txns)",
			addrs:      []string{},
			goldenFile: "all-unconfirmed-txns-verbose.golden",
		})
	}

	c := newClient()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			txnResult, err := c.UnconfirmedTransactionsVerbose(tc.addrs)
			if err != nil {
				require.Equal(t, tc.err, err, "case: "+tc.name)
				return
			}

			assertNoTransactionsDupesVerbose(t, txnResult)

			var expected []readable.TransactionWithStatusVerbose
			checkGoldenFile(t, tc.goldenFile, TestData{txnResult, &expected})
		})
	}
}

func assertNoTransactionsDupesVerbose(t *testing.T, r []readable.TransactionWithStatusVerbose) {
	txids := make(map[string]struct{})

	for _, x := range r {
		_, ok := txids[x.Transaction.Hash]
		require.False(t, ok)
		txids[x.Transaction.Hash] = struct{}{}
	}
}

func TestLiveUnconfirmedTransactionsVerbose(t *testing.T) {
	if !doLive(t) {
		return
	}
	c := newClient()

	cTxsSingle, err := c.UnconfirmedTransactionsVerbose([]string{"2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKt"})
	require.NoError(t, err)
	require.True(t, len(cTxsSingle) >= 0)

	cTxsAll, err := c.UnconfirmedTransactionsVerbose([]string{})
	require.NoError(t, err)
	require.True(t, len(cTxsAll) >= 0)
	require.True(t, len(cTxsAll) >= len(cTxsSingle))
}

func TestStableResendUnconfirmedTransactions(t *testing.T) {
	if !doStable(t) {
		return
	}
	c := newClient()
	_, err := c.ResendUnconfirmedTransactions()
	require.NotNil(t, err)
	respErr, ok := err.(api.ClientError)
	require.True(t, ok)
	require.Equal(t, fmt.Sprintf("503 Service Unavailable - %s", daemon.ErrNetworkingDisabled), respErr.Message)
	require.Equal(t, http.StatusServiceUnavailable, respErr.StatusCode)
}

func TestLiveResendUnconfirmedTransactions(t *testing.T) {
	if !doLive(t) {
		return
	}
	c := newClient()
	_, err := c.ResendUnconfirmedTransactions()

	if liveDisableNetworking(t) {
		require.NotNil(t, err)
		respErr, ok := err.(api.ClientError)
		require.True(t, ok)
		require.Equal(t, fmt.Sprintf("503 Service Unavailable - %s", daemon.ErrNetworkingDisabled), respErr.Message)
		require.Equal(t, http.StatusServiceUnavailable, respErr.StatusCode)
	} else {
		require.NoError(t, err)
	}
}

type rawTransactionTestCase struct {
	name   string
	txID   string
	err    api.ClientError
	rawTxn string
}

func TestStableRawTransaction(t *testing.T) {
	if !doStable(t) {
		return
	}

	cases := []rawTransactionTestCase{
		{
			name: "invalid hex length",
			txID: "abcd",
			err: api.ClientError{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Message:    "400 Bad Request - Invalid hex length",
			},
		},
		{
			name: "not found",
			txID: "540582ee4128b733f810f149e908d984a5f403ad2865108e6c1c5423aeefc759",
			err: api.ClientError{
				Status:     "404 Not Found",
				StatusCode: http.StatusNotFound,
				Message:    "404 Not Found",
			},
		},
		{
			name: "odd length hex string",
			txID: "abcdeffedca",
			err: api.ClientError{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Message:    "400 Bad Request - encoding/hex: odd length hex string",
			},
		},
		{
			name:   "OK",
			txID:   "d556c1c7abf1e86138316b8c17183665512dc67633c04cf236a8b7f332cb4add",
			rawTxn: "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000100000000f8f9c644772dc5373d85e11094e438df707a42c900407a10f35a000000407a10f35a0000",
		},
	}

	if !dbNoUnconfirmed(t) {
		cases = append(cases, rawTransactionTestCase{
			name:   "unconfirmed",
			txID:   "701d23fd513bad325938ba56869f9faba19384a8ec3dd41833aff147eac53947",
			rawTxn: "dc00000000f8293dbfdddcc56a97664655ceee650715d35a0dda32a9f0ce0e2e99d4899124010000003981061c7275ae9cc936e902a5367fdd87ef779bbdb31e1e10d325d17a129abb34f6e597ceeaf67bb051774b41c58276004f6a63cb81de61d4693bc7a5536f320001000000fe6762d753d626115c8dd3a053b5fb75d6d419a8d0fb1478c5fffc1fe41c5f2002000000003be2537f8c0893fddcddc878518f38ea493d949e008988068d0000002739570000000000009037ff169fbec6db95e2537e4ff79396c050aeeb00e40b54020000002739570000000000",
		})
	}

	c := newClient()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			txnResult, err := c.RawTransaction(tc.txID)
			if err != nil {
				require.Equal(t, tc.err, err, "case: "+tc.name)
				return
			}
			require.Equal(t, tc.rawTxn, txnResult, "case: "+tc.name)
		})
	}
}

func TestLiveRawTransaction(t *testing.T) {
	if !doLive(t) {
		return
	}

	cases := []rawTransactionTestCase{
		{
			name: "invalid hex length",
			txID: "abcd",
			err: api.ClientError{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Message:    "400 Bad Request - Invalid hex length",
			},
		},
		{
			name: "odd length hex string",
			txID: "abcdeffedca",
			err: api.ClientError{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Message:    "400 Bad Request - encoding/hex: odd length hex string",
			},
		},
		{
			name:   "OK - genesis tx",
			txID:   "d556c1c7abf1e86138316b8c17183665512dc67633c04cf236a8b7f332cb4add",
			rawTxn: "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000100000000f8f9c644772dc5373d85e11094e438df707a42c900407a10f35a000000407a10f35a0000",
		},
		{
			name:   "OK",
			txID:   "540582ee4128b733f810f149e908d984a5f403ad2865108e6c1c5423aeefc759",
			rawTxn: "3d0100000088b4e967d77a8b7155c5378a85c199fabf94048aa84833ef5eab7818545bcda80200000071985c70041fe5a6408a2dfac2ea4963820bc603059521259debb114b2f6630b5658e7ff665b2db7878ce9b0d1d051ec66b5dea23274e52642bc7e451b273a90008afb06133958b03c4795d5a7acd001f3942cc6d3b19e93d357d2675fe9ba8bbf3db30b3cda779e441fced581aee88f48c8af017b30dc276b15be25d4bb44260c000200000050386f195b367f8261e66e3fdfbc942fbacfe25e117e554ca1c1caf8993454767afab03c823346ff8b00c29df6acc05841583d90dfd451ba09e66884a48e83f70200000000ef3b60779f014b3c7acf27c16c9acc3ff3bea61600a8b54b06000000c2ba2400000000000037274869aaa4c2e2e5c91595024c65f8f9458102404b4c0000000000c2ba240000000000",
		},
	}

	c := newClient()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			txnResult, err := c.RawTransaction(tc.txID)
			if err != nil {
				require.Equal(t, tc.err, err, "case: "+tc.name)
				return
			}

			require.Equal(t, tc.rawTxn, txnResult, "case: "+tc.name)
		})
	}
}

type addressTransactionsTestCase struct {
	name    string
	address string
	golden  string
	errCode int
	errMsg  string
}

func TestStableAddressTransactions(t *testing.T) {
	// Formerly tested /api/v1/explorer/address, now tests /api/v1/transactions?verbose=1
	if !doStable(t) {
		return
	}

	cases := []addressTransactionsTestCase{
		{
			name:    "genesis address",
			address: "2jBbGxZRGoQG1mqhPBnXnLTxK6oxsTf8os6",
			golden:  "address-transactions-2jBbGxZRGoQG1mqhPBnXnLTxK6oxsTf8os6.golden",
		},
		{
			name:    "address with transactions",
			address: "ALJVNKYL7WGxFBSriiZuwZKWD4b7fbV1od",
			golden:  "address-transactions-ALJVNKYL7WGxFBSriiZuwZKWD4b7fbV1od.golden",
		},
		{
			name:    "address without transactions",
			address: "2b8ourW8fbTkC1yQBSLseVt6srhXvNMHvn9",
			golden:  "address-transactions-2b8ourW8fbTkC1yQBSLseVt6srhXvNMHvn9.golden",
		},
		{
			name:    "invalid address",
			address: "prRXwTcDK24hs6AFxj",
			errCode: http.StatusBadRequest,
			errMsg:  "400 Bad Request - parse parameter: 'addrs' failed: address \"prRXwTcDK24hs6AFxj\" is invalid: Invalid address length",
		},
	}

	if !dbNoUnconfirmed(t) {
		cases = append(cases, []addressTransactionsTestCase{
			{
				name:    "address with outgoing transaction",
				address: "R6aHqKWSQfvpdo2fGSrq4F1RYXkBWR9HHJ",
				golden:  "address-transactions-outgoing-R6aHqKWSQfvpdo2fGSrq4F1RYXkBWR9HHJ.golden",
			},
			{
				name:    "address with incoming transaction",
				address: "212mwY3Dmey6vwnWpiph99zzCmopXTqeVEN",
				golden:  "address-transactions-incoming-212mwY3Dmey6vwnWpiph99zzCmopXTqeVEN.golden",
			},
		}...)
	}

	c := newClient()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			txns, err := c.TransactionsVerbose([]string{tc.address})
			if tc.errMsg != "" {
				assertResponseError(t, err, tc.errCode, tc.errMsg)
				return
			}

			require.NoError(t, err)

			var expected []readable.TransactionWithStatusVerbose
			checkGoldenFile(t, tc.golden, TestData{txns, &expected})
		})
	}
}

func TestLiveAddressTransactions(t *testing.T) {
	// Formerly tested /api/v1/explorer/address, now tests /api/v1/transactions?verbose=1
	if !doLive(t) {
		return
	}

	cases := []addressTransactionsTestCase{
		{
			name: "address with transactions",
			// This is the first distribution address which has spent all of its coins
			// Its transactions list should not change, unless someone sends coins to it
			address: "R6aHqKWSQfvpdo2fGSrq4F1RYXkBWR9HHJ",
			golden:  "address-transactions-R6aHqKWSQfvpdo2fGSrq4F1RYXkBWR9HHJ.golden",
		},
		{
			name: "address without transactions",
			// This is a randomly generated address, never used
			// It should never see new transactions
			// (if it ever does, somebody managed to generate this address for use and there is a serious bug)
			address: "2RRpfMDmPHEyG4LWmNYT6eWj5VcmUfCJY6D",
			golden:  "address-transactions-2RRpfMDmPHEyG4LWmNYT6eWj5VcmUfCJY6D.golden",
		},
		{
			name:    "invalid address",
			address: "prRXwTcDK24hs6AFxj",
			errCode: http.StatusBadRequest,
			errMsg:  "400 Bad Request - parse parameter: 'addrs' failed: address \"prRXwTcDK24hs6AFxj\" is invalid: Invalid address length",
		},
	}

	c := newClient()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			txns, err := c.TransactionsVerbose([]string{tc.address})
			if tc.errMsg != "" {
				assertResponseError(t, err, tc.errCode, tc.errMsg)
				return
			}

			require.NoError(t, err)

			// Unset height since it is not stable
			for i := range txns {
				txns[i].Status.Height = 0
			}

			var expected []readable.TransactionWithStatusVerbose
			checkGoldenFile(t, tc.golden, TestData{txns, &expected})
		})
	}
}

func TestStableRichlist(t *testing.T) {
	if !doStable(t) {
		return
	}

	c := newClient()

	richlist, err := c.Richlist(nil)
	require.NoError(t, err)

	var expected api.Richlist
	checkGoldenFile(t, "richlist-default.golden", TestData{*richlist, &expected})

	richlist, err = c.Richlist(&api.RichlistParams{
		N:                   0,
		IncludeDistribution: false,
	})
	require.NoError(t, err)

	expected = api.Richlist{}
	checkGoldenFile(t, "richlist-all.golden", TestData{*richlist, &expected})

	richlist, err = c.Richlist(&api.RichlistParams{
		N:                   0,
		IncludeDistribution: true,
	})
	require.NoError(t, err)

	expected = api.Richlist{}
	checkGoldenFile(t, "richlist-all-include-distribution.golden", TestData{*richlist, &expected})

	richlist, err = c.Richlist(&api.RichlistParams{
		N:                   8,
		IncludeDistribution: false,
	})
	require.NoError(t, err)

	expected = api.Richlist{}
	checkGoldenFile(t, "richlist-8.golden", TestData{*richlist, &expected})

	richlist, err = c.Richlist(&api.RichlistParams{
		N:                   150,
		IncludeDistribution: true,
	})
	require.NoError(t, err)

	expected = api.Richlist{}
	checkGoldenFile(t, "richlist-150-include-distribution.golden", TestData{*richlist, &expected})
}

func TestLiveRichlist(t *testing.T) {
	if !doLive(t) {
		return
	}

	c := newClient()

	richlist, err := c.Richlist(nil)
	require.NoError(t, err)

	require.NotEmpty(t, richlist.Richlist)
	require.Len(t, richlist.Richlist, 20)

	richlist, err = c.Richlist(&api.RichlistParams{
		N:                   150,
		IncludeDistribution: true,
	})
	require.NoError(t, err)

	require.Len(t, richlist.Richlist, 150)
}

func TestStableAddressCount(t *testing.T) {
	if !doStable(t) {
		return
	}

	c := newClient()

	count, err := c.AddressCount()
	require.NoError(t, err)

	require.Equal(t, uint64(155), count)
}

func TestLiveAddressCount(t *testing.T) {
	if !doLive(t) {
		return
	}

	c := newClient()

	count, err := c.AddressCount()
	require.NoError(t, err)

	// 5296 addresses as of 2018-03-06, the count could decrease but is unlikely to
	require.True(t, count > 5000)
}

func TestStableNoUnconfirmedPendingTransactions(t *testing.T) {
	if !doStable(t) || !dbNoUnconfirmed(t) {
		return
	}

	c := newClient()

	txns, err := c.PendingTransactions()
	require.NoError(t, err)
	require.Empty(t, txns)
}

func TestStablePendingTransactions(t *testing.T) {
	if !doStable(t) || dbNoUnconfirmed(t) {
		return
	}

	c := newClient()

	txns, err := c.PendingTransactions()
	require.NoError(t, err)

	// Convert Received and Checked times to UTC for stable comparison
	for i, txn := range txns {
		require.False(t, txn.Received.IsZero())
		require.False(t, txn.Checked.IsZero())

		txns[i].Received = txn.Received.UTC()
		txns[i].Checked = txn.Checked.UTC()
	}

	var expect []readable.UnconfirmedTransactions
	checkGoldenFile(t, "pending-transactions.golden", TestData{txns, &expect})
}

func TestLivePendingTransactions(t *testing.T) {
	if !doLive(t) {
		return
	}

	c := newClient()

	_, err := c.PendingTransactions()
	require.NoError(t, err)
}

func TestStableNoUnconfirmedPendingTransactionsVerbose(t *testing.T) {
	if !doStable(t) || !dbNoUnconfirmed(t) {
		return
	}

	c := newClient()

	txns, err := c.PendingTransactionsVerbose()
	require.NoError(t, err)
	require.Empty(t, txns)
}

func TestStablePendingTransactionsVerbose(t *testing.T) {
	if !doStable(t) || dbNoUnconfirmed(t) {
		return
	}

	c := newClient()

	txns, err := c.PendingTransactionsVerbose()
	require.NoError(t, err)

	// Convert Received and Checked times to UTC for stable comparison
	for i, txn := range txns {
		require.False(t, txn.Received.IsZero())
		require.False(t, txn.Checked.IsZero())

		txns[i].Received = txn.Received.UTC()
		txns[i].Checked = txn.Checked.UTC()
	}

	var expect []readable.UnconfirmedTransactionVerbose
	checkGoldenFile(t, "verbose-pending-transactions.golden", TestData{txns, &expect})
}

func TestLivePendingTransactionsVerbose(t *testing.T) {
	if !doLive(t) {
		return
	}

	c := newClient()

	_, err := c.PendingTransactionsVerbose()
	require.NoError(t, err)
}

func TestDisableWalletAPI(t *testing.T) {
	if !doDisableWalletAPI(t) {
		return
	}

	changeAddress := testutil.MakeAddress().String()

	type testCase struct {
		name        string
		method      string
		endpoint    string
		contentType string
		body        func() io.Reader
		json        func() interface{}
		expectErr   string
		code        int
	}

	tt := []testCase{
		{
			name:      "get wallet",
			method:    http.MethodGet,
			endpoint:  "/api/v1/wallet?id=test.wlt",
			expectErr: "403 Forbidden - Endpoint is disabled",
			code:      http.StatusForbidden,
		},
		{
			name:     "create wallet",
			method:   http.MethodPost,
			endpoint: "/api/v1/wallet/create",
			body: func() io.Reader {
				v := url.Values{}
				v.Add("seed", "seed")
				v.Add("label", "label")
				v.Add("scan", "1")
				return strings.NewReader(v.Encode())
			},
			expectErr: "403 Forbidden - Endpoint is disabled",
			code:      http.StatusForbidden,
		},
		{
			name:     "generate new address",
			method:   http.MethodPost,
			endpoint: "/api/v1/wallet/newAddress",
			body: func() io.Reader {
				v := url.Values{}
				v.Add("id", "test.wlt")
				return strings.NewReader(v.Encode())
			},
			expectErr: "403 Forbidden - Endpoint is disabled",
			code:      http.StatusForbidden,
		},
		{
			name:      "get wallet balance",
			method:    http.MethodGet,
			endpoint:  "/api/v1/wallet/balance?id=test.wlt",
			expectErr: "403 Forbidden - Endpoint is disabled",
			code:      http.StatusForbidden,
		},
		{
			name:      "get wallet unconfirmed transactions",
			method:    http.MethodGet,
			endpoint:  "/api/v1/wallet/transactions?id=test.wlt",
			expectErr: "403 Forbidden - Endpoint is disabled",
			code:      http.StatusForbidden,
		},
		{
			name:     "update wallet label",
			method:   http.MethodPost,
			endpoint: "/api/v1/wallet/update",
			body: func() io.Reader {
				v := url.Values{}
				v.Add("id", "test.wlt")
				v.Add("label", "label")
				return strings.NewReader(v.Encode())
			},
			expectErr: "403 Forbidden - Endpoint is disabled",
			code:      http.StatusForbidden,
		},
		{
			name:      "new seed",
			method:    http.MethodGet,
			endpoint:  "/api/v1/wallet/newSeed",
			expectErr: "403 Forbidden - Endpoint is disabled",
			code:      http.StatusForbidden,
		},
		{
			name:     "new seed",
			method:   http.MethodPost,
			endpoint: "/api/v1/wallet/seed",
			body: func() io.Reader {
				v := url.Values{}
				v.Add("id", "test.wlt")
				return strings.NewReader(v.Encode())
			},
			expectErr: "403 Forbidden - Endpoint is disabled",
			code:      http.StatusForbidden,
		},
		{
			name:      "get wallets",
			method:    http.MethodGet,
			endpoint:  "/api/v1/wallets",
			expectErr: "403 Forbidden - Endpoint is disabled",
			code:      http.StatusForbidden,
		},
		{
			name:      "get wallets folder name",
			method:    http.MethodGet,
			endpoint:  "/api/v1/wallets/folderName",
			expectErr: "403 Forbidden - Endpoint is disabled",
			code:      http.StatusForbidden,
		},
		{
			name:      "main index.html 404 not found",
			method:    http.MethodGet,
			endpoint:  "/api/v1/",
			expectErr: "404 Not Found",
			code:      http.StatusNotFound,
		},
		{
			name:     "encrypt wallet",
			method:   http.MethodPost,
			endpoint: "/api/v1/wallet/encrypt",
			body: func() io.Reader {
				v := url.Values{}
				v.Add("id", "test.wlt")
				v.Add("password", "pwd")
				return strings.NewReader(v.Encode())
			},
			expectErr: "403 Forbidden - Endpoint is disabled",
			code:      http.StatusForbidden,
		},
		{
			name:     "decrypt wallet",
			method:   http.MethodPost,
			endpoint: "/api/v1/wallet/decrypt",
			body: func() io.Reader {
				v := url.Values{}
				v.Add("id", "test.wlt")
				v.Add("password", "pwd")
				return strings.NewReader(v.Encode())
			},
			expectErr: "403 Forbidden - Endpoint is disabled",
			code:      http.StatusForbidden,
		},
		{
			name:     "get wallet seed",
			method:   http.MethodPost,
			endpoint: "/api/v1/wallet/seed",
			body: func() io.Reader {
				v := url.Values{}
				v.Add("id", "test.wlt")
				v.Add("password", "pwd")
				return strings.NewReader(v.Encode())
			},
			expectErr: "403 Forbidden - Endpoint is disabled",
			code:      http.StatusForbidden,
		},
		{
			name:        "create transaction",
			method:      http.MethodPost,
			endpoint:    "/api/v1/wallet/transaction",
			contentType: api.ContentTypeJSON,
			json: func() interface{} {
				return api.WalletCreateTransactionRequest{
					WalletID: "test.wlt",
					CreateTransactionRequest: api.CreateTransactionRequest{
						HoursSelection: api.HoursSelection{
							Type: transaction.HoursSelectionTypeManual,
						},
						ChangeAddress: &changeAddress,
						To: []api.Receiver{
							{
								Address: changeAddress,
								Coins:   "0.001",
								Hours:   "1",
							},
						},
					},
				}
			},
			expectErr: "403 Forbidden - Endpoint is disabled",
			code:      http.StatusForbidden,
		},
	}

	c := newClient()
	for _, tc := range tt {
		f := func(tc testCase) func(t *testing.T) {
			return func(t *testing.T) {
				var err error
				switch tc.method {
				case http.MethodGet:
					err = c.Get(tc.endpoint, nil)
				case http.MethodPost:
					switch tc.contentType {
					case api.ContentTypeJSON:
						err = c.PostJSON(tc.endpoint, tc.json(), nil)
					default:
						err = c.PostForm(tc.endpoint, tc.body(), nil)
					}
				}
				assertResponseError(t, err, tc.code, tc.expectErr)
			}
		}

		t.Run(tc.name, f(tc))
	}

	// Confirms that no new wallet is created
	// API_WALLET_DIR environment variable is set in ci-script/integration-test-disable-wallet-api.sh
	walletDir := os.Getenv("API_WALLET_DIR")
	if walletDir == "" {
		t.Fatal("API_WALLET_DIR is not set")
	}

	// Confirms that the wallet directory does not exist
	testutil.RequireFileNotExists(t, walletDir)
}

func checkHealthResponse(t *testing.T, r *api.HealthResponse) {
	require.NotEmpty(t, r.BlockchainMetadata.Unspents)
	require.NotEmpty(t, r.BlockchainMetadata.Head.BkSeq)
	require.NotEmpty(t, r.BlockchainMetadata.Head.Time)
	require.NotEmpty(t, r.Version.Version)
	require.True(t, r.Uptime.Duration > time.Duration(0))
	require.NotEmpty(t, r.CoinName)
	require.NotEmpty(t, r.DaemonUserAgent)

	_, err := useragent.Parse(r.DaemonUserAgent)
	require.NoError(t, err)
}

func TestStableHealth(t *testing.T) {
	if !doStable(t) {
		return
	}

	c := newClient()

	r, err := c.Health()
	require.NoError(t, err)

	checkHealthResponse(t, r)

	require.Equal(t, 0, r.OpenConnections)
	require.Equal(t, 0, r.IncomingConnections)
	require.Equal(t, 0, r.OutgoingConnections)

	require.True(t, r.BlockchainMetadata.TimeSinceLastBlock.Duration > time.Duration(0))

	// The stable node is always run with the commit and branch ldflags, so they should appear
	require.NotEmpty(t, r.Version.Commit)
	require.NotEmpty(t, r.Version.Branch)

	coinName := os.Getenv("COIN")
	require.Equal(t, coinName, r.CoinName)
	require.Equal(t, fmt.Sprintf("%s:%s", coinName, r.Version.Version), r.DaemonUserAgent)

	_, err = useragent.Parse(r.DaemonUserAgent)
	require.NoError(t, err)

	require.Equal(t, useCSRF(t), r.CSRFEnabled)
	require.Equal(t, doHeaderCheck(t), r.HeaderCheckEnabled)
	require.True(t, r.CSPEnabled)
	require.True(t, r.WalletAPIEnabled)
	require.False(t, r.GUIEnabled)
}

func TestLiveHealth(t *testing.T) {
	if !doLive(t) {
		return
	}

	c := newClient()

	r, err := c.Health()
	require.NoError(t, err)

	checkHealthResponse(t, r)

	if liveDisableNetworking(t) {
		require.Equal(t, 0, r.OpenConnections)
		require.Equal(t, 0, r.OutgoingConnections)
		require.Equal(t, 0, r.IncomingConnections)
	} else {
		require.NotEqual(t, 0, r.OpenConnections)
	}

	require.Equal(t, r.OutgoingConnections+r.IncomingConnections, r.OpenConnections)

	// The TimeSinceLastBlock can be any value, including negative values, due to clock skew
	// The live node is not necessarily run with the commit and branch ldflags, so don't check them
}

func TestDisableGUIAPI(t *testing.T) {
	if !doDisableGUI(t) {
		return
	}

	c := newClient()
	err := c.Get("/", nil)
	assertResponseError(t, err, http.StatusNotFound, "404 Not Found")
}

func TestInvalidAuth(t *testing.T) {
	if !doLiveOrStable(t) {
		return
	}

	user := nodeUsername()
	pass := nodePassword()

	c := newClient()

	require.Equal(t, user, c.Username)
	require.Equal(t, pass, c.Password)

	if user == "" && pass == "" {
		// If neither user nor pass are set, assume the node is run without auth,
		// and check that providing auth returns a 401 error
		testInvalidAuthNoAuthConfigured(t)
	} else {
		// If either user or pass are set, assume the node is run with auth,
		// and check that missing or invalid auth returns a 401 error
		testInvalidAuthAuthConfigured(t)
	}
}

func testInvalidAuthAuthConfigured(t *testing.T) {
	cases := []struct {
		user string
		pass string
	}{
		{}, // both missing
		{
			user: nodeUsername(), // user right, pass missing
		},
		{
			pass: nodePassword(), // pass right, user missing
		},
		{
			user: nodeUsername() + "x", // user wrong, pass missing
		},
		{
			pass: nodePassword() + "x", // pass wrong, user missing
		},
		{
			user: nodeUsername() + "x", // both wrong
			pass: nodePassword() + "x",
		},
		{
			user: nodeUsername(), // user right, pass wrong
			pass: nodePassword() + "x",
		},
		{
			user: nodeUsername() + "x", // user wrong, pass right
			pass: nodePassword(),
		},
	}

	for _, tc := range cases {
		name := fmt.Sprintf("u=%s p=%s", tc.user, tc.pass)
		t.Run(name, func(t *testing.T) {
			c := newClient()
			c.SetAuth(tc.user, tc.pass)
			require.Equal(t, tc.user, c.Username)
			require.Equal(t, tc.pass, c.Password)

			_, err := c.Health()
			assertResponseError(t, err, http.StatusUnauthorized, "401 Unauthorized")
		})
	}
}

func testInvalidAuthNoAuthConfigured(t *testing.T) {
	cases := []struct {
		user string
		pass string
	}{
		{
			user: "foo",
		},
		{
			pass: "bar",
		},
		{
			user: "foo",
			pass: "bar",
		},
	}

	for _, tc := range cases {
		name := fmt.Sprintf("u=%s p=%s", tc.user, tc.pass)
		t.Run(name, func(t *testing.T) {
			c := newClient()
			c.SetAuth(tc.user, tc.pass)
			require.Equal(t, tc.user, c.Username)
			require.Equal(t, tc.pass, c.Password)

			_, err := c.Health()
			assertResponseError(t, err, http.StatusUnauthorized, "401 Unauthorized")
		})
	}
}
