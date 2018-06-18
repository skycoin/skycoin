// package integration_test implements API integration tests
package integration_test

import (
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
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

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/api"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/testutil"
	"github.com/skycoin/skycoin/src/util/droplet" //http,json helpers
	"github.com/skycoin/skycoin/src/util/fee"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/visor/historydb"
	"github.com/skycoin/skycoin/src/wallet"
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
	testModeDisableWalletApi = "disable-wallet-api"
	testModeDisableSeedApi   = "disable-seed-api"

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

func mode(t *testing.T) string {
	mode := os.Getenv("SKYCOIN_INTEGRATION_TEST_MODE")
	switch mode {
	case "":
		mode = testModeStable
	case testModeLive,
		testModeStable,
		testModeDisableWalletApi,
		testModeDisableSeedApi:
	default:
		t.Fatal("Invalid test mode, must be stable, live or disable-wallet-api")
	}
	return mode
}

func enabled() bool {
	return os.Getenv("SKYCOIN_INTEGRATION_TESTS") == "1"
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

func doDisableWalletApi(t *testing.T) bool {
	if enabled() && mode(t) == testModeDisableWalletApi {
		return true
	}

	t.Skip("DisableWalletApi tests disabled")
	return false
}

func doDisableSeedApi(t *testing.T) bool {
	if enabled() && mode(t) == testModeDisableSeedApi {
		return true
	}

	t.Skip("EnableSeedAPI tests disabled")
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

func loadJSON(t *testing.T, filename string, obj interface{}) {
	f, err := os.Open(filename)
	require.NoError(t, err, filename)
	defer f.Close()

	d := json.NewDecoder(f)
	d.DisallowUnknownFields()

	err = d.Decode(obj)
	require.NoError(t, err, filename)
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

	require.Equal(t, string(c), string(b)+"\n", "json struct output differs from golden file, was a field added to the struct?")
}

func assertResponseError(t *testing.T, err error, errCode int, errMsg string) {
	require.Error(t, err)
	require.IsType(t, api.ClientError{}, err)
	require.Equal(t, errCode, err.(api.ClientError).StatusCode)
	require.Equal(t, errMsg, err.(api.ClientError).Message)
}

func TestStableCoinSupply(t *testing.T) {
	if !doStable(t) {
		return
	}

	c := api.NewClient(nodeAddress())

	cs, err := c.CoinSupply()
	require.NoError(t, err)

	var expected api.CoinSupply
	checkGoldenFile(t, "coinsupply.golden", TestData{*cs, &expected})
}

func TestLiveCoinSupply(t *testing.T) {
	if !doLive(t) {
		return
	}

	c := api.NewClient(nodeAddress())

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

	c := api.NewClient(nodeAddress())

	v, err := c.Version()
	require.NoError(t, err)

	require.NotEmpty(t, v.Version)
}

func TestVerifyAddress(t *testing.T) {
	if !doLiveOrStable(t) {
		return
	}

	c := api.NewClient(nodeAddress())

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

	c := api.NewClient(nodeAddress())

	badSigStr := "71f2c01516fe696328e79bcf464eb0db374b63d494f7a307d1e77114f18581d7a81eed5275a9e04a336292dd2fd16977d9bef2a54ea3161d0876603d00c53bc9dd"
	badSigBytes, err := hex.DecodeString(badSigStr)
	require.NoError(t, err)
	badSig := cipher.NewSig(badSigBytes)

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
	badSignatureTxn.UpdateHeader()

	cases := []struct {
		name    string
		golden  string
		txn     coin.Transaction
		errCode int
		errMsg  string
	}{
		{
			name:    "invalid transaction empty",
			txn:     coin.Transaction{},
			golden:  "verify-transaction-invalid-empty.golden",
			errCode: http.StatusUnprocessableEntity,
			errMsg:  "Transaction violates soft constraint: Transaction has zero coinhour fee",
		},

		{
			name:    "invalid transaction bad signature",
			txn:     badSignatureTxn,
			golden:  "verify-transaction-invalid-bad-sig.golden",
			errCode: http.StatusUnprocessableEntity,
			errMsg:  "Transaction violates hard constraint: Signature invalid for hash",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			encodedTxn := hex.EncodeToString(tc.txn.Serialize())

			resp, err := c.VerifyTransaction(encodedTxn)

			if tc.errCode != 0 && tc.errCode != http.StatusOK {
				assertResponseError(t, err, tc.errCode, tc.errMsg)
				if tc.errCode != http.StatusUnprocessableEntity {
					return
				}
			}

			if tc.errCode != http.StatusUnprocessableEntity {
				require.NoError(t, err)
			}

			var expected api.VerifyTxnResponse
			checkGoldenFile(t, tc.golden, TestData{*resp, &expected})
		})
	}

}

func TestStableOutputs(t *testing.T) {
	if !doStable(t) {
		return
	}

	c := api.NewClient(nodeAddress())

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
			},
			golden: "outputs-hashes.golden",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			require.False(t, tc.addrs != nil && tc.hashes != nil)

			var outputs *visor.ReadableOutputSet
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

			var expected visor.ReadableOutputSet
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

	c := api.NewClient(nodeAddress())

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

func TestLiveBlock(t *testing.T) {
	if !doLive(t) {
		return
	}

	testKnownBlocks(t)

	// These blocks were affected by the coinhour overflow issue, make sure that they can be queried
	blockSeqs := []uint64{11685, 11707, 11710, 11709, 11705, 11708, 11711, 11706, 11699}

	c := api.NewClient(nodeAddress())
	for _, seq := range blockSeqs {
		b, err := c.BlockBySeq(seq)
		require.NoError(t, err)
		require.Equal(t, seq, b.Head.BkSeq)
	}
}

func testKnownBlocks(t *testing.T) {
	c := api.NewClient(nodeAddress())

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
			errMsg:  "404 Not Found\n",
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
			name:   "seq 100",
			golden: "block-seq-100.golden",
			seq:    100,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var b *visor.ReadableBlock
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

			var expected visor.ReadableBlock
			checkGoldenFile(t, tc.golden, TestData{*b, &expected})
		})
	}

	t.Logf("Querying every block in the blockchain")

	// Scan every block by seq
	progress, err := c.BlockchainProgress()
	require.NoError(t, err)

	var prevBlock *visor.ReadableBlock
	for i := uint64(0); i < progress.Current; i++ {
		t.Run(fmt.Sprintf("block-seq-%d", i), func(t *testing.T) {
			b, err := c.BlockBySeq(i)
			require.NoError(t, err)
			require.NotNil(t, b)
			require.Equal(t, i, b.Head.BkSeq)

			if prevBlock != nil {
				require.Equal(t, prevBlock.Head.BlockHash, b.Head.PreviousBlockHash)
			}

			bHash, err := c.BlockByHash(b.Head.BlockHash)
			require.NoError(t, err)
			require.NotNil(t, bHash)
			require.Equal(t, b, bHash)

			prevBlock = b
		})
	}
}

func TestStableBlockchainMetadata(t *testing.T) {
	if !doStable(t) {
		return
	}

	c := api.NewClient(nodeAddress())

	metadata, err := c.BlockchainMetadata()
	require.NoError(t, err)

	var expected visor.BlockchainMetadata
	checkGoldenFile(t, "blockchain-metadata.golden", TestData{*metadata, &expected})
}

func TestLiveBlockchainMetadata(t *testing.T) {
	if !doLive(t) {
		return
	}

	c := api.NewClient(nodeAddress())

	metadata, err := c.BlockchainMetadata()
	require.NoError(t, err)

	require.NotEqual(t, uint64(0), metadata.Head.BkSeq)
}

func TestStableBlockchainProgress(t *testing.T) {
	if !doStable(t) {
		return
	}

	c := api.NewClient(nodeAddress())

	progress, err := c.BlockchainProgress()
	require.NoError(t, err)

	var expected daemon.BlockchainProgress
	checkGoldenFile(t, "blockchain-progress.golden", TestData{*progress, &expected})
}

func TestLiveBlockchainProgress(t *testing.T) {
	if !doLive(t) {
		return
	}

	c := api.NewClient(nodeAddress())

	progress, err := c.BlockchainProgress()
	require.NoError(t, err)

	require.NotEqual(t, uint64(0), progress.Current)
	require.True(t, progress.Current <= progress.Highest)
	require.NotEmpty(t, progress.Peers)
}

func TestStableBalance(t *testing.T) {
	if !doStable(t) {
		return
	}

	c := api.NewClient(nodeAddress())

	cases := []struct {
		name   string
		golden string
		addrs  []string
	}{
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

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			balance, err := c.Balance(tc.addrs)
			require.NoError(t, err)

			var expected wallet.BalancePair
			checkGoldenFile(t, tc.golden, TestData{*balance, &expected})
		})
	}
}

func TestLiveBalance(t *testing.T) {
	if !doLive(t) {
		return
	}

	c := api.NewClient(nodeAddress())

	// Genesis address check, should not have a balance
	b, err := c.Balance([]string{"2jBbGxZRGoQG1mqhPBnXnLTxK6oxsTf8os6"})
	require.NoError(t, err)
	require.Equal(t, wallet.BalancePair{}, *b)

	// Balance of final distribution address. Should have the same coins balance
	// for the next 15-20 years.
	b, err = c.Balance([]string{"ejJjiCwp86ykmFr5iTJ8LxQXJ2wJPTYmkm"})
	require.NoError(t, err)
	require.Equal(t, b.Confirmed, b.Predicted)
	require.NotEmpty(t, b.Confirmed.Hours)
	require.Equal(t, uint64(1e6*1e6), b.Confirmed.Coins)

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

	c := api.NewClient(nodeAddress())

	cases := []struct {
		name   string
		golden string
		uxID   string
	}{
		{
			name:   "valid uxID",
			golden: "uxout.golden",
			uxID:   "fe6762d753d626115c8dd3a053b5fb75d6d419a8d0fb1478c5fffc1fe41c5f20",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ux, err := c.UxOut(tc.uxID)
			require.NoError(t, err)

			var expected historydb.UxOutJSON
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

	c := api.NewClient(nodeAddress())

	// A spent uxout should never change
	ux, err := c.UxOut("fe6762d753d626115c8dd3a053b5fb75d6d419a8d0fb1478c5fffc1fe41c5f20")
	require.NoError(t, err)

	var expected historydb.UxOutJSON
	checkGoldenFile(t, "uxout-spent.golden", TestData{*ux, &expected})
	require.NotEqual(t, uint64(0), ux.SpentBlockSeq)

	// Scan all uxouts from the result of /outputs
	scanUxOuts(t)
}

func scanUxOuts(t *testing.T) {
	c := api.NewClient(nodeAddress())

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
				require.Equal(t, "0000000000000000000000000000000000000000000000000000000000000000", foundUx.SpentTxID)
			} else {
				require.NotEqual(t, "0000000000000000000000000000000000000000000000000000000000000000", foundUx.SpentTxID)
			}
		})
	}
}

func TestStableAddressUxOuts(t *testing.T) {
	if !doStable(t) {
		return
	}

	c := api.NewClient(nodeAddress())

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
			errMsg:  "400 Bad Request - address is empty\n",
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
			var expected []*historydb.UxOutJSON
			checkGoldenFile(t, tc.golden, TestData{ux, &expected})
		})
	}
}

func TestLiveAddressUxOuts(t *testing.T) {
	if !doLive(t) {
		return
	}

	c := api.NewClient(nodeAddress())

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
			errMsg:  "400 Bad Request - address is empty\n",
		},
		{
			name:    "invalid address length",
			errCode: http.StatusBadRequest,
			errMsg:  "400 Bad Request - Invalid address length\n",
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

func TestStableBlocks(t *testing.T) {
	if !doStable(t) {
		return
	}

	c := api.NewClient(nodeAddress())

	progress, err := c.BlockchainProgress()
	require.NoError(t, err)

	lastNBlocks := 10
	require.True(t, int(progress.Current) > lastNBlocks+1)

	cases := []struct {
		name    string
		golden  string
		start   int
		end     int
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
			start:  int(progress.Current) - lastNBlocks,
			end:    int(progress.Current),
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
			end:    int(progress.Current),
		},
		{
			name:   "start > end",
			golden: "blocks-end-less-than-start.golden",
			start:  10,
			end:    9,
		},
		{
			name:    "start negative",
			start:   -10,
			end:     9,
			errCode: http.StatusBadRequest,
			errMsg:  "400 Bad Request - Invalid start value \"-10\"\n",
		},
		{
			name:    "end negative",
			start:   10,
			end:     -9,
			errCode: http.StatusBadRequest,
			errMsg:  "400 Bad Request - Invalid end value \"-9\"\n",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.errMsg == "" {
				resp := testBlocks(t, tc.start, tc.end)

				var expected visor.ReadableBlocks
				checkGoldenFile(t, tc.golden, TestData{*resp, &expected})
			} else {
				_, err := c.Blocks(tc.start, tc.end)
				assertResponseError(t, err, tc.errCode, tc.errMsg)
			}
		})
	}
}

func TestLiveBlocks(t *testing.T) {
	if !doLive(t) {
		return
	}

	testBlocks(t, 1, 10)
}

func testBlocks(t *testing.T, start, end int) *visor.ReadableBlocks {
	c := api.NewClient(nodeAddress())

	blocks, err := c.Blocks(start, end)
	require.NoError(t, err)

	if start > end {
		require.Empty(t, blocks.Blocks)
	} else {
		require.Len(t, blocks.Blocks, end-start+1)
	}

	var prevBlock *visor.ReadableBlock
	for idx, b := range blocks.Blocks {
		if prevBlock != nil {
			require.Equal(t, prevBlock.Head.BlockHash, b.Head.PreviousBlockHash)
		}

		bHash, err := c.BlockByHash(b.Head.BlockHash)
		require.Equal(t, uint64(idx+start), b.Head.BkSeq)
		require.NoError(t, err)
		require.NotNil(t, bHash)
		require.Equal(t, b, *bHash)

		prevBlock = &blocks.Blocks[idx]
	}

	return blocks
}

func TestStableLastBlocks(t *testing.T) {
	if !doStable(t) {
		return
	}

	c := api.NewClient(nodeAddress())

	blocks, err := c.LastBlocks(1)
	require.NoError(t, err)

	var expected *visor.ReadableBlocks
	checkGoldenFile(t, "block-last.golden", TestData{blocks, &expected})

	var prevBlock *visor.ReadableBlock
	blocks, err = c.LastBlocks(10)
	require.NoError(t, err)
	require.Equal(t, 10, len(blocks.Blocks))
	for idx, b := range blocks.Blocks {
		if prevBlock != nil {
			require.Equal(t, prevBlock.Head.BlockHash, b.Head.PreviousBlockHash)
		}

		bHash, err := c.BlockByHash(b.Head.BlockHash)
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
	c := api.NewClient(nodeAddress())
	var prevBlock *visor.ReadableBlock
	blocks, err := c.LastBlocks(10)
	require.NoError(t, err)
	require.Equal(t, 10, len(blocks.Blocks))
	for idx, b := range blocks.Blocks {
		if prevBlock != nil {
			require.Equal(t, prevBlock.Head.BlockHash, b.Head.PreviousBlockHash)
		}

		bHash, err := c.BlockByHash(b.Head.BlockHash)
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

	c := api.NewClient(nodeAddress())
	connections, err := c.NetworkConnections()
	require.NoError(t, err)
	require.Empty(t, connections.Connections)

	connection, err := c.NetworkConnection("127.0.0.1:4444")
	assertResponseError(t, err, http.StatusNotFound, "404 Not Found\n")
	require.Nil(t, connection)
}

func TestLiveNetworkConnections(t *testing.T) {
	if !doLive(t) {
		return
	}

	c := api.NewClient(nodeAddress())
	connections, err := c.NetworkConnections()
	require.NoError(t, err)
	require.NotEmpty(t, connections.Connections)

	for _, cc := range connections.Connections {
		connection, err := c.NetworkConnection(cc.Addr)
		require.NoError(t, err)
		require.NotEmpty(t, cc.Addr)
		require.Equal(t, cc.Addr, connection.Addr)
		require.Equal(t, cc.ID, connection.ID)
		require.Equal(t, cc.ListenPort, connection.ListenPort)
		require.Equal(t, cc.Mirror, connection.Mirror)
		require.Equal(t, cc.Introduced, connection.Introduced)
		require.Equal(t, cc.Outgoing, connection.Outgoing)
		require.True(t, cc.LastReceived <= connection.LastReceived)
		require.True(t, cc.LastSent <= connection.LastSent)
		require.True(t, cc.Height >= 0)
	}
}

func TestNetworkDefaultConnections(t *testing.T) {
	if !doLiveOrStable(t) {
		return
	}

	c := api.NewClient(nodeAddress())
	connections, err := c.NetworkDefaultConnections()
	require.NoError(t, err)
	require.NotEmpty(t, connections)
	sort.Strings(connections)

	var expected []string
	checkGoldenFile(t, "network-default-connections.golden", TestData{connections, &expected})
}

func TestNetworkTrustedConnections(t *testing.T) {
	if !doLiveOrStable(t) {
		return
	}

	c := api.NewClient(nodeAddress())
	connections, err := c.NetworkTrustedConnections()
	require.NoError(t, err)
	require.NotEmpty(t, connections)
	sort.Strings(connections)

	var expected []string
	checkGoldenFile(t, "network-trusted-connections.golden", TestData{connections, &expected})
}

func TestStableNetworkExchangeableConnections(t *testing.T) {
	if !doStable(t) {
		return
	}

	c := api.NewClient(nodeAddress())
	connections, err := c.NetworkExchangeableConnections()
	require.NoError(t, err)

	var expected []string
	checkGoldenFile(t, "network-exchangeable-connections.golden", TestData{connections, &expected})
}

func TestLiveNetworkExchangeableConnections(t *testing.T) {
	if !doLive(t) {
		return
	}

	c := api.NewClient(nodeAddress())
	_, err := c.NetworkExchangeableConnections()
	require.NoError(t, err)
}

func TestLiveTransaction(t *testing.T) {
	if !doLive(t) {
		return
	}

	cases := []struct {
		name       string
		txId       string
		err        api.ClientError
		goldenFile string
	}{
		{
			name: "invalid txId",
			txId: "abcd",
			err: api.ClientError{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Message:    "400 Bad Request - Invalid hex length\n",
			},
		},
		{
			name: "empty txId",
			txId: "",
			err: api.ClientError{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Message:    "400 Bad Request - txid is empty\n",
			},
		},
		{
			name:       "OK",
			txId:       "76ecbabc53ea2a3be46983058433dda6a3cf7ea0b86ba14d90b932fa97385de7",
			goldenFile: "./transaction.golden",
		},
	}

	c := api.NewClient(nodeAddress())
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tx, err := c.Transaction(tc.txId)
			if err != nil {
				require.Equal(t, tc.err, err)
				return
			}
			var expected *visor.ReadableTransaction
			loadGoldenFile(t, tc.goldenFile, TestData{tx, &expected})
			require.Equal(t, expected, &tx.Transaction)
		})
	}
}

func TestStableTransaction(t *testing.T) {
	if !doStable(t) {
		return
	}

	cases := []struct {
		name       string
		txId       string
		err        api.ClientError
		goldenFile string
	}{
		{
			name: "invalid txId",
			txId: "abcd",
			err: api.ClientError{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Message:    "400 Bad Request - Invalid hex length\n",
			},
			goldenFile: "",
		},
		{
			name: "not exist",
			txId: "701d23fd513bad325938ba56869f9faba19384a8ec3dd41833aff147eac53947",
			err: api.ClientError{
				Status:     "404 Not Found",
				StatusCode: http.StatusNotFound,
				Message:    "404 Not Found\n",
			},
			goldenFile: "",
		},
		{
			name: "empty txId",
			txId: "",
			err: api.ClientError{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Message:    "400 Bad Request - txid is empty\n",
			},
			goldenFile: "",
		},
		{
			name:       "genesis transaction",
			txId:       "d556c1c7abf1e86138316b8c17183665512dc67633c04cf236a8b7f332cb4add",
			goldenFile: "genesis-transaction.golden",
		},
	}

	c := api.NewClient(nodeAddress())
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tx, err := c.Transaction(tc.txId)
			if err != nil {
				require.Equal(t, tc.err, err)
				return
			}

			var expected *visor.ReadableTransaction
			loadGoldenFile(t, tc.goldenFile, TestData{tx, &expected})
			require.Equal(t, expected, &tx.Transaction)
		})
	}
}

func TestLiveTransactions(t *testing.T) {
	if !doLive(t) {
		return
	}

	c := api.NewClient(nodeAddress())
	addrs := []string{
		"2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKt",
	}
	txns, err := c.Transactions(addrs)
	require.NoError(t, err)
	require.True(t, len(*txns) > 0)
}

func TestStableTransactions(t *testing.T) {
	if !doStable(t) {
		return
	}

	cases := []struct {
		name       string
		addrs      []string
		err        api.ClientError
		goldenFile string
	}{
		{
			name:  "invalid addr length",
			addrs: []string{"abcd"},
			err: api.ClientError{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Message:    "400 Bad Request - parse parameter: 'addrs' failed: Invalid address length\n",
			},
		},
		{
			name:  "invalid addr character",
			addrs: []string{"701d23fd513bad325938ba56869f9faba19384a8ec3dd41833aff147eac53947"},
			err: api.ClientError{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Message:    "400 Bad Request - parse parameter: 'addrs' failed: Invalid base58 character\n",
			},
		},
		{
			name:  "invalid checksum",
			addrs: []string{"2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKk"},
			err: api.ClientError{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Message:    "400 Bad Request - parse parameter: 'addrs' failed: Invalid checksum\n",
			},
		},
		{
			name:  "empty addrs",
			addrs: []string{},
			err: api.ClientError{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Message:    "400 Bad Request - txId is empty\n",
			},
			goldenFile: "./empty-addrs.golden",
		},
		{
			name:       "single addr",
			addrs:      []string{"2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKt"},
			goldenFile: "./single-addr.golden",
		},
	}

	c := api.NewClient(nodeAddress())
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			txResult, err := c.Transactions(tc.addrs)
			if err != nil {
				require.Equal(t, tc.err, err, "case: "+tc.name)
				return
			}

			var expected *[]daemon.TransactionResult
			checkGoldenFile(t, tc.goldenFile, TestData{txResult, &expected})
		})
	}
}

func TestLiveConfirmedTransactions(t *testing.T) {
	if !doLive(t) {
		return
	}
	c := api.NewClient(nodeAddress())

	ctxsSingle, err := c.ConfirmedTransactions([]string{"2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKt"})
	require.NoError(t, err)
	require.True(t, len(*ctxsSingle) > 0)

	ctxsAll, err := c.ConfirmedTransactions([]string{})
	require.NoError(t, err)
	require.True(t, len(*ctxsAll) > 0)
	require.True(t, len(*ctxsAll) > len(*ctxsSingle))
}

func TestStableConfirmedTransactions(t *testing.T) {
	if !doStable(t) {
		return
	}
	cases := []struct {
		name       string
		addrs      []string
		err        api.ClientError
		goldenFile string
	}{
		{
			name:  "invalid addr length",
			addrs: []string{"abcd"},
			err: api.ClientError{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Message:    "400 Bad Request - parse parameter: 'addrs' failed: Invalid address length\n",
			},
		},
		{
			name:  "invalid addr character",
			addrs: []string{"701d23fd513bad325938ba56869f9faba19384a8ec3dd41833aff147eac53947"},
			err: api.ClientError{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Message:    "400 Bad Request - parse parameter: 'addrs' failed: Invalid base58 character\n",
			},
		},
		{
			name:  "invalid checksum",
			addrs: []string{"2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKk"},
			err: api.ClientError{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Message:    "400 Bad Request - parse parameter: 'addrs' failed: Invalid checksum\n",
			},
		},
		{
			name:       "empty addrs",
			addrs:      []string{},
			goldenFile: "./empty-addrs.golden",
		},
		{
			name:       "single addr",
			addrs:      []string{"2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKt"},
			goldenFile: "./single-addr.golden",
		},
	}

	c := api.NewClient(nodeAddress())
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			txResult, err := c.ConfirmedTransactions(tc.addrs)
			if err != nil {
				require.Equal(t, tc.err, err, "case: "+tc.name)
				return
			}

			var expected *[]daemon.TransactionResult
			checkGoldenFile(t, tc.goldenFile, TestData{txResult, &expected})
		})
	}
}

func TestStableUnconfirmedTransactions(t *testing.T) {
	if !doStable(t) {
		return
	}
	cases := []struct {
		name       string
		addrs      []string
		err        api.ClientError
		goldenFile string
	}{
		{
			name:  "invalid addr length",
			addrs: []string{"abcd"},
			err: api.ClientError{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Message:    "400 Bad Request - parse parameter: 'addrs' failed: Invalid address length\n",
			},
		},
		{
			name:  "invalid addr character",
			addrs: []string{"701d23fd513bad325938ba56869f9faba19384a8ec3dd41833aff147eac53947"},
			err: api.ClientError{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Message:    "400 Bad Request - parse parameter: 'addrs' failed: Invalid base58 character\n",
			},
		},
		{
			name:  "invalid checksum",
			addrs: []string{"2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKk"},
			err: api.ClientError{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Message:    "400 Bad Request - parse parameter: 'addrs' failed: Invalid checksum\n",
			},
		},
		{
			name:       "empty addrs",
			addrs:      []string{},
			goldenFile: "./empty-addrs-unconfirmed-txs.golden",
		},
	}

	c := api.NewClient(nodeAddress())
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			txResult, err := c.UnconfirmedTransactions(tc.addrs)
			if err != nil {
				require.Equal(t, tc.err, err, "case: "+tc.name)
				return
			}

			var expected *[]daemon.TransactionResult
			checkGoldenFile(t, tc.goldenFile, TestData{txResult, &expected})
		})
	}
}

func TestLiveUnconfirmedTransactions(t *testing.T) {
	if !doLive(t) {
		return
	}
	c := api.NewClient(nodeAddress())

	cTxsSingle, err := c.UnconfirmedTransactions([]string{"2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKt"})
	require.NoError(t, err)
	require.True(t, len(*cTxsSingle) >= 0)

	cTxsAll, err := c.UnconfirmedTransactions([]string{})
	require.NoError(t, err)
	require.True(t, len(*cTxsAll) >= 0)
	require.True(t, len(*cTxsAll) >= len(*cTxsSingle))
}

func TestStableResendUnconfirmedTransactions(t *testing.T) {
	if !doStable(t) {
		return
	}
	c := api.NewClient(nodeAddress())
	res, err := c.ResendUnconfirmedTransactions()
	require.NoError(t, err)
	require.True(t, len(res.Txids) == 0)
}

func TestLiveResendUnconfirmedTransactions(t *testing.T) {
	if !doLive(t) {
		return
	}
	c := api.NewClient(nodeAddress())
	_, err := c.ResendUnconfirmedTransactions()
	require.NoError(t, err)
}

func TestStableRawTransaction(t *testing.T) {
	if !doStable(t) {
		return
	}

	cases := []struct {
		name  string
		txId  string
		err   api.ClientError
		rawTx string
	}{
		{
			name: "invalid hex length",
			txId: "abcd",
			err: api.ClientError{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Message:    "400 Bad Request - Invalid hex length\n",
			},
		},
		{
			name: "not found",
			txId: "701d23fd513bad325938ba56869f9faba19384a8ec3dd41833aff147eac53947",
			err: api.ClientError{
				Status:     "404 Not Found",
				StatusCode: http.StatusNotFound,
				Message:    "404 Not Found\n",
			},
		},
		{
			name: "odd length hex string",
			txId: "abcdeffedca",
			err: api.ClientError{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Message:    "400 Bad Request - encoding/hex: odd length hex string\n",
			},
		},
		{
			name:  "OK",
			txId:  "d556c1c7abf1e86138316b8c17183665512dc67633c04cf236a8b7f332cb4add",
			rawTx: "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000100000000f8f9c644772dc5373d85e11094e438df707a42c900407a10f35a000000407a10f35a0000",
		},
	}

	c := api.NewClient(nodeAddress())
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			txResult, err := c.RawTransaction(tc.txId)
			if err != nil {
				require.Equal(t, tc.err, err, "case: "+tc.name)
				return
			}
			require.Equal(t, tc.rawTx, txResult, "case: "+tc.name)
		})
	}
}

func TestLiveRawTransaction(t *testing.T) {
	if !doLive(t) {
		return
	}

	cases := []struct {
		name  string
		txId  string
		err   api.ClientError
		rawTx string
	}{
		{
			name: "invalid hex length",
			txId: "abcd",
			err: api.ClientError{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Message:    "400 Bad Request - Invalid hex length\n",
			},
		},
		{
			name: "odd length hex string",
			txId: "abcdeffedca",
			err: api.ClientError{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Message:    "400 Bad Request - encoding/hex: odd length hex string\n",
			},
		},
		{
			name:  "OK - genesis tx",
			txId:  "d556c1c7abf1e86138316b8c17183665512dc67633c04cf236a8b7f332cb4add",
			rawTx: "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000100000000f8f9c644772dc5373d85e11094e438df707a42c900407a10f35a000000407a10f35a0000",
		},
		{
			name:  "OK",
			txId:  "701d23fd513bad325938ba56869f9faba19384a8ec3dd41833aff147eac53947",
			rawTx: "dc00000000f8293dbfdddcc56a97664655ceee650715d35a0dda32a9f0ce0e2e99d4899124010000003981061c7275ae9cc936e902a5367fdd87ef779bbdb31e1e10d325d17a129abb34f6e597ceeaf67bb051774b41c58276004f6a63cb81de61d4693bc7a5536f320001000000fe6762d753d626115c8dd3a053b5fb75d6d419a8d0fb1478c5fffc1fe41c5f2002000000003be2537f8c0893fddcddc878518f38ea493d949e008988068d0000002739570000000000009037ff169fbec6db95e2537e4ff79396c050aeeb00e40b54020000002739570000000000",
		},
	}

	c := api.NewClient(nodeAddress())
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			txResult, err := c.RawTransaction(tc.txId)
			if err != nil {
				require.Equal(t, tc.err, err, "case: "+tc.name)
				return
			}
			require.Equal(t, tc.rawTx, txResult, "case: "+tc.name)
		})
	}
}

func TestWalletNewSeed(t *testing.T) {
	if !doLiveOrStable(t) {
		return
	}

	cases := []struct {
		name     string
		entropy  int
		numWords int
		errCode  int
		errMsg   string
	}{
		{
			name:     "entropy 128",
			entropy:  128,
			numWords: 12,
		},
		{
			name:     "entropy 256",
			entropy:  256,
			numWords: 24,
		},
		{
			name:    "entropy 100",
			entropy: 100,
			errCode: http.StatusBadRequest,
			errMsg:  "400 Bad Request - entropy length must be 128 or 256\n",
		},
	}

	c := api.NewClient(nodeAddress())
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			seed, err := c.NewSeed(tc.entropy)
			if tc.errMsg != "" {
				assertResponseError(t, err, tc.errCode, tc.errMsg)
				return
			}

			require.NoError(t, err)
			words := strings.Split(seed, " ")
			require.Len(t, words, tc.numWords)

			// no extra whitespace on the seed
			require.Equal(t, seed, strings.TrimSpace(seed))

			// should generate a different seed each time
			seed2, err := c.NewSeed(tc.entropy)
			require.NoError(t, err)
			require.NotEqual(t, seed, seed2)
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
	if !doStable(t) {
		return
	}

	cases := []addressTransactionsTestCase{
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
			errMsg:  "400 Bad Request - invalid address\n",
		},
	}

	c := api.NewClient(nodeAddress())
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			txns, err := c.AddressTransactions(tc.address)
			if tc.errMsg != "" {
				assertResponseError(t, err, tc.errCode, tc.errMsg)
				return
			}

			require.NoError(t, err)

			var expected []daemon.ReadableTransaction
			checkGoldenFile(t, tc.golden, TestData{txns, &expected})
		})
	}
}

func TestLiveAddressTransactions(t *testing.T) {
	if !doLive(t) {
		return
	}

	cases := []addressTransactionsTestCase{
		{
			name: "address with transactions",
			// This is the first distribution address which has spent all of its coins
			// It's transactions list should not change, unless someone sends coins to it
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
			errMsg:  "400 Bad Request - invalid address\n",
		},
	}

	c := api.NewClient(nodeAddress())
	// Get current blockchain height
	bp, err := c.BlockchainProgress()
	require.NoError(t, err)
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			txns, err := c.AddressTransactions(tc.address)
			if tc.errMsg != "" {
				assertResponseError(t, err, tc.errCode, tc.errMsg)
				return
			}

			require.NoError(t, err)

			var expected []daemon.ReadableTransaction
			loadGoldenFile(t, tc.golden, TestData{txns, &expected})

			// Recaculate the height if it's live test
			for i := range expected {
				expected[i].Status.Height = bp.Current - expected[i].Status.BlockSeq + 1
			}

			require.Equal(t, expected, txns)
		})
	}
}

func TestStableRichlist(t *testing.T) {
	if !doStable(t) {
		return
	}

	c := api.NewClient(nodeAddress())

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

	c := api.NewClient(nodeAddress())

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

	c := api.NewClient(nodeAddress())

	count, err := c.AddressCount()
	require.NoError(t, err)

	require.Equal(t, uint64(155), count)
}

func TestLiveAddressCount(t *testing.T) {
	if !doLive(t) {
		return
	}

	c := api.NewClient(nodeAddress())

	count, err := c.AddressCount()
	require.NoError(t, err)

	// 5296 addresses as of 2018-03-06, the count could decrease but is unlikely to
	require.True(t, count > 5000)
}

func TestStablePendingTransactions(t *testing.T) {
	if !doStable(t) {
		return
	}

	c := api.NewClient(nodeAddress())

	txns, err := c.PendingTransactions()
	require.NoError(t, err)
	require.Empty(t, txns)
}

func TestLivePendingTransactions(t *testing.T) {
	if !doLive(t) {
		return
	}

	c := api.NewClient(nodeAddress())

	_, err := c.PendingTransactions()
	require.NoError(t, err)
}

func TestLiveWalletSpend(t *testing.T) {
	if !doLive(t) {
		return
	}

	requireWalletEnv(t)

	c := api.NewClient(nodeAddress())
	w, totalCoins, _, password := prepareAndCheckWallet(t, c, 2e6, 2)

	tt := []struct {
		name    string
		to      string
		coins   uint64
		errMsg  []byte
		checkTx func(t *testing.T, tx *daemon.TransactionResult)
	}{
		{
			name:  "send all coins to the first address",
			to:    w.Entries[0].Address.String(),
			coins: totalCoins,
			checkTx: func(t *testing.T, tx *daemon.TransactionResult) {
				// Confirms the total output coins are equal to the totalCoins
				var coins uint64
				for _, o := range tx.Transaction.Out {
					c, err := droplet.FromString(o.Coins)
					require.NoError(t, err)
					coins, err = coin.AddUint64(coins, c)
					require.NoError(t, err)
				}

				// Confirms the address balance are equal to the totalCoins
				coins, _ = getAddressBalance(t, c, w.Entries[0].Address.String())
				require.Equal(t, totalCoins, coins)
			},
		},
		{
			// send 0.003 coin to the second address,
			// this amount is chosen to not interfere with TestLiveWalletCreateTransaction
			name:  "send 0.003 coin to second address",
			to:    w.Entries[1].Address.String(),
			coins: 3e3,
			checkTx: func(t *testing.T, tx *daemon.TransactionResult) {
				// Confirms there're two outputs, one to the second address, one as change output to the first address.
				require.Len(t, tx.Transaction.Out, 2)

				// Gets the output of the second address in the transaction
				getAddrOutputInTx := func(t *testing.T, tx *daemon.TransactionResult, addr string) *visor.ReadableTransactionOutput {
					for _, output := range tx.Transaction.Out {
						if output.Address == addr {
							return &output
						}
					}
					t.Fatalf("transaction doesn't have output to address: %v", addr)
					return nil
				}

				out := getAddrOutputInTx(t, tx, w.Entries[1].Address.String())

				// Confirms the second address has 0.003 coin
				require.Equal(t, out.Coins, "0.003000")
				require.Equal(t, out.Address, w.Entries[1].Address.String())

				coin, err := droplet.FromString(out.Coins)
				require.NoError(t, err)

				// Gets the expected change coins
				expectChangeCoins := totalCoins - coin

				// Gets the real change coins
				changeOut := getAddrOutputInTx(t, tx, w.Entries[0].Address.String())
				changeCoins, err := droplet.FromString(changeOut.Coins)
				require.NoError(t, err)
				// Confirms the change coins are matched.
				require.Equal(t, expectChangeCoins, changeCoins)
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result, err := c.Spend(w.Filename(), tc.to, tc.coins, password)
			if err != nil {
				t.Fatalf("spend failed: %v", err)
			}

			tk := time.NewTicker(time.Second)
			var tx *daemon.TransactionResult
		loop:
			for {
				select {
				case <-time.After(30 * time.Second):
					t.Fatal("Waiting for transaction to be confirmed timeout")
				case <-tk.C:
					tx = getTransaction(t, c, result.Transaction.Hash)
					if tx.Status.Confirmed {
						break loop
					}
				}
			}
			tc.checkTx(t, tx)
		})
	}

	// Return if wallet is encrypted, cause the rest of the tests will spend a lot of time.
	if w.IsEncrypted() {
		return
	}

	// Confirms sending coins less than 0.001 is not allowed
	errMsg := "500 Internal Server Error - Transaction violates soft constraint: invalid amount, too many decimal places\n"
	for i := uint64(1); i < uint64(1000); i++ {
		cs, err := droplet.ToString(i)
		require.NoError(t, err)
		name := fmt.Sprintf("send invalid coin %v", cs)
		t.Run(name, func(t *testing.T) {
			result, err := c.Spend(w.Filename(), w.Entries[0].Address.String(), i, password)
			if w.IsEncrypted() && len(password) == 0 {
				assertResponseError(t, err, http.StatusBadRequest, "400 Bad Request - missing password\n")
				return
			}
			assertResponseError(t, err, http.StatusInternalServerError, errMsg)
			require.Nil(t, result)
		})
	}
}

func TestLiveWalletCreateTransactionSpecific(t *testing.T) {
	if !doLive(t) {
		return
	}

	requireWalletEnv(t)

	c := api.NewClient(nodeAddress())

	w, totalCoins, totalHours, password := prepareAndCheckWallet(t, c, 2e6, 20)

	remainingHours := fee.RemainingHours(totalHours)
	require.True(t, remainingHours > 1)

	addresses := make([]string, len(w.Entries))
	addressMap := make(map[string]struct{}, len(w.Entries))
	for i, e := range w.Entries {
		addresses[i] = e.Address.String()
		addressMap[e.Address.String()] = struct{}{}
	}

	// Get all outputs
	outputs, err := c.Outputs()
	require.NoError(t, err)

	// Split outputs into those held by the wallet and those not
	var walletOutputHashes []string
	var walletOutputs visor.ReadableOutputs
	walletAuxs := make(map[string][]string)
	var nonWalletOutputs visor.ReadableOutputs
	for _, o := range outputs.HeadOutputs {
		if _, ok := addressMap[o.Address]; ok {
			walletOutputs = append(walletOutputs, o)
			walletOutputHashes = append(walletOutputHashes, o.Hash)
			walletAuxs[o.Address] = append(walletAuxs[o.Address], o.Hash)
		} else {
			nonWalletOutputs = append(nonWalletOutputs, o)
		}
	}

	require.NotEmpty(t, walletOutputs)
	require.NotEmpty(t, nonWalletOutputs)

	unknownOutput := testutil.RandSHA256(t)

	toDropletString := func(i uint64) string {
		x, err := droplet.ToString(i)
		require.NoError(t, err)
		return x
	}

	defaultChangeAddress := w.Entries[0].Address.String()

	type testCase struct {
		name                 string
		req                  api.CreateTransactionRequest
		outputs              []coin.TransactionOutput
		outputsSubset        []coin.TransactionOutput
		err                  string
		code                 int
		ignoreHours          bool
		additionalRespVerify func(t *testing.T, r *api.CreateTransactionResponse)
	}

	cases := []testCase{
		{
			name: "invalid decimals",
			req: api.CreateTransactionRequest{
				HoursSelection: api.HoursSelection{
					Type: wallet.HoursSelectionTypeManual,
				},
				Wallet: api.CreateTransactionRequestWallet{
					ID:       w.Filename(),
					Password: password,
				},
				ChangeAddress: &defaultChangeAddress,
				To: []api.Receiver{
					{
						Address: w.Entries[0].Address.String(),
						Coins:   "0.0001",
						Hours:   "1",
					},
				},
			},
			err:  "400 Bad Request - to[0].coins has too many decimal places\n",
			code: http.StatusBadRequest,
		},

		{
			name: "overflowing hours",
			req: api.CreateTransactionRequest{
				HoursSelection: api.HoursSelection{
					Type: wallet.HoursSelectionTypeManual,
				},
				Wallet: api.CreateTransactionRequestWallet{
					ID:       w.Filename(),
					Password: password,
				},
				ChangeAddress: &defaultChangeAddress,
				To: []api.Receiver{
					{
						Address: w.Entries[0].Address.String(),
						Coins:   "0.001",
						Hours:   "1",
					},
					{
						Address: w.Entries[0].Address.String(),
						Coins:   "0.001",
						Hours:   fmt.Sprint(uint64(math.MaxUint64)),
					},
					{
						Address: w.Entries[0].Address.String(),
						Coins:   "0.001",
						Hours:   fmt.Sprint(uint64(math.MaxUint64) - 1),
					},
				},
			},
			err:  "400 Bad Request - total output hours error: uint64 addition overflow\n",
			code: http.StatusBadRequest,
		},

		{
			name: "insufficient coins",
			req: api.CreateTransactionRequest{
				HoursSelection: api.HoursSelection{
					Type: wallet.HoursSelectionTypeManual,
				},
				Wallet: api.CreateTransactionRequestWallet{
					ID:       w.Filename(),
					Password: password,
				},
				ChangeAddress: &defaultChangeAddress,
				To: []api.Receiver{
					{
						Address: w.Entries[0].Address.String(),
						Coins:   fmt.Sprint(totalCoins + 1),
						Hours:   "1",
					},
				},
			},
			err:  "400 Bad Request - balance is not sufficient\n",
			code: http.StatusBadRequest,
		},

		{
			name: "insufficient hours",
			req: api.CreateTransactionRequest{
				HoursSelection: api.HoursSelection{
					Type: wallet.HoursSelectionTypeManual,
				},
				Wallet: api.CreateTransactionRequestWallet{
					ID:       w.Filename(),
					Password: password,
				},
				ChangeAddress: &defaultChangeAddress,
				To: []api.Receiver{
					{
						Address: w.Entries[0].Address.String(),
						Coins:   toDropletString(totalCoins),
						Hours:   fmt.Sprint(totalHours + 1),
					},
				},
			},
			err:  "400 Bad Request - hours are not sufficient\n",
			code: http.StatusBadRequest,
		},

		{
			// NOTE: this test will fail if "totalCoins - 1e3" does not require
			// all of the outputs to be spent, e.g. if there is an output with
			// "totalCoins - 1e3" coins in it.
			// TODO -- Check that the wallet does not have an output of 0.001,
			// because then this test cannot be performed, since there is no
			// way to use all outputs and produce change in that case.
			name: "valid request, manual one output with change, spend all",
			req: api.CreateTransactionRequest{
				HoursSelection: api.HoursSelection{
					Type: wallet.HoursSelectionTypeManual,
				},
				Wallet: api.CreateTransactionRequestWallet{
					ID:       w.Filename(),
					Password: password,
				},
				ChangeAddress: &defaultChangeAddress,
				To: []api.Receiver{
					{
						Address: w.Entries[1].Address.String(),
						Coins:   toDropletString(totalCoins - 1e3),
						Hours:   "1",
					},
				},
			},
			outputs: []coin.TransactionOutput{
				{
					Address: w.Entries[1].Address,
					Coins:   totalCoins - 1e3,
					Hours:   1,
				},
				{
					Address: w.Entries[0].Address,
					Coins:   1e3,
					Hours:   remainingHours - 1,
				},
			},
		},

		{
			// NOTE: this test will fail if "totalCoins - 1e3" does not require
			// all of the outputs to be spent, e.g. if there is an output with
			// "totalCoins - 1e3" coins in it.
			// TODO -- Check that the wallet does not have an output of 0.001,
			// because then this test cannot be performed, since there is no
			// way to use all outputs and produce change in that case.
			name: "valid request, manual one output with change, spend all, unspecified change address",
			req: api.CreateTransactionRequest{
				HoursSelection: api.HoursSelection{
					Type: wallet.HoursSelectionTypeManual,
				},
				Wallet: api.CreateTransactionRequestWallet{
					ID:       w.Filename(),
					Password: password,
				},
				To: []api.Receiver{
					{
						Address: w.Entries[1].Address.String(),
						Coins:   toDropletString(totalCoins - 1e3),
						Hours:   "1",
					},
				},
			},
			outputs: []coin.TransactionOutput{
				{
					Address: w.Entries[1].Address,
					Coins:   totalCoins - 1e3,
					Hours:   1,
				},
				{
					// Address omitted -- will be check later in the test body
					Coins: 1e3,
					Hours: remainingHours - 1,
				},
			},
		},

		{
			name: "valid request, manual one output with change, don't spend all",
			req: api.CreateTransactionRequest{
				HoursSelection: api.HoursSelection{
					Type: wallet.HoursSelectionTypeManual,
				},
				Wallet: api.CreateTransactionRequestWallet{
					ID:       w.Filename(),
					Password: password,
				},
				ChangeAddress: &defaultChangeAddress,
				To: []api.Receiver{
					{
						Address: w.Entries[1].Address.String(),
						Coins:   toDropletString(1e3),
						Hours:   "1",
					},
				},
			},
			outputsSubset: []coin.TransactionOutput{
				{
					Address: w.Entries[1].Address,
					Coins:   1e3,
					Hours:   1,
				},
				// NOTE: change omitted,
				// change is too difficult to predict in this case, we are
				// just checking that not all uxouts get spent in the transaction
			},
		},

		{
			name: "valid request, manual one output no change",
			req: api.CreateTransactionRequest{
				HoursSelection: api.HoursSelection{
					Type: wallet.HoursSelectionTypeManual,
				},
				Wallet: api.CreateTransactionRequestWallet{
					ID:       w.Filename(),
					Password: password,
				},
				ChangeAddress: &defaultChangeAddress,
				To: []api.Receiver{
					{
						Address: w.Entries[1].Address.String(),
						Coins:   toDropletString(totalCoins),
						Hours:   "1",
					},
				},
			},
			outputs: []coin.TransactionOutput{
				{
					Address: w.Entries[1].Address,
					Coins:   totalCoins,
					Hours:   1,
				},
			},
		},

		{
			// NOTE: no reliable way to test the ignore unconfirmed behavior,
			// this test only checks that if IgnoreUnconfirmed is specified,
			// the API doesn't throw up some parsing error
			name: "valid request, manual one output no change, ignore unconfirmed",
			req: api.CreateTransactionRequest{
				IgnoreUnconfirmed: true,
				HoursSelection: api.HoursSelection{
					Type: wallet.HoursSelectionTypeManual,
				},
				Wallet: api.CreateTransactionRequestWallet{
					ID:       w.Filename(),
					Password: password,
				},
				ChangeAddress: &defaultChangeAddress,
				To: []api.Receiver{
					{
						Address: w.Entries[1].Address.String(),
						Coins:   toDropletString(totalCoins),
						Hours:   "1",
					},
				},
			},
			outputs: []coin.TransactionOutput{
				{
					Address: w.Entries[1].Address,
					Coins:   totalCoins,
					Hours:   1,
				},
			},
		},

		{
			name: "valid request, auto one output no change, share factor recalculates to 1.0",
			req: api.CreateTransactionRequest{
				HoursSelection: api.HoursSelection{
					Type:        wallet.HoursSelectionTypeAuto,
					Mode:        wallet.HoursSelectionModeShare,
					ShareFactor: "0.5",
				},
				Wallet: api.CreateTransactionRequestWallet{
					ID:       w.Filename(),
					Password: password,
				},
				ChangeAddress: &defaultChangeAddress,
				To: []api.Receiver{
					{
						Address: w.Entries[1].Address.String(),
						Coins:   toDropletString(totalCoins),
					},
				},
			},
			outputs: []coin.TransactionOutput{
				{
					Address: w.Entries[1].Address,
					Coins:   totalCoins,
					Hours:   remainingHours,
				},
			},
		},

		{
			name: "valid request, auto two outputs with change",
			req: api.CreateTransactionRequest{
				HoursSelection: api.HoursSelection{
					Type:        wallet.HoursSelectionTypeAuto,
					Mode:        wallet.HoursSelectionModeShare,
					ShareFactor: "0.5",
				},
				Wallet: api.CreateTransactionRequestWallet{
					ID:       w.Filename(),
					Password: password,
				},
				ChangeAddress: &defaultChangeAddress,
				To: []api.Receiver{
					{
						Address: w.Entries[1].Address.String(),
						Coins:   toDropletString(1e3),
					},
					{
						Address: w.Entries[1].Address.String(),
						Coins:   toDropletString(totalCoins - 2e3),
					},
				},
			},
			outputs: []coin.TransactionOutput{
				{
					Address: w.Entries[1].Address,
					Coins:   1e3,
				},
				{
					Address: w.Entries[1].Address,
					Coins:   totalCoins - 2e3,
				},
				{
					Address: w.Entries[0].Address,
					Coins:   1e3,
				},
			},
			ignoreHours: true, // the hours are too unpredictable
		},

		{
			name: "uxout does not exist",
			req: api.CreateTransactionRequest{
				HoursSelection: api.HoursSelection{
					Type: wallet.HoursSelectionTypeManual,
				},
				Wallet: api.CreateTransactionRequestWallet{
					ID:       w.Filename(),
					Password: password,
					UxOuts:   []string{unknownOutput.Hex()},
				},
				ChangeAddress: &defaultChangeAddress,
				To: []api.Receiver{
					{
						Address: w.Entries[1].Address.String(),
						Coins:   toDropletString(totalCoins),
						Hours:   "1",
					},
				},
			},
			err:  fmt.Sprintf("400 Bad Request - unspent output of %s does not exist\n", unknownOutput.Hex()),
			code: http.StatusBadRequest,
		},

		{
			name: "uxout not held by the wallet",
			req: api.CreateTransactionRequest{
				HoursSelection: api.HoursSelection{
					Type: wallet.HoursSelectionTypeManual,
				},
				Wallet: api.CreateTransactionRequestWallet{
					ID:       w.Filename(),
					Password: password,
					UxOuts:   []string{nonWalletOutputs[0].Hash},
				},
				ChangeAddress: &defaultChangeAddress,
				To: []api.Receiver{
					{
						Address: w.Entries[1].Address.String(),
						Coins:   toDropletString(totalCoins),
						Hours:   "1",
					},
				},
			},
			err:  "400 Bad Request - uxout is not owned by any address in the wallet\n",
			code: http.StatusBadRequest,
		},

		{
			name: "insufficient balance with uxouts",
			req: api.CreateTransactionRequest{
				HoursSelection: api.HoursSelection{
					Type: wallet.HoursSelectionTypeManual,
				},
				Wallet: api.CreateTransactionRequestWallet{
					ID:       w.Filename(),
					Password: password,
					UxOuts:   []string{walletOutputs[0].Hash},
				},
				ChangeAddress: &defaultChangeAddress,
				To: []api.Receiver{
					{
						Address: w.Entries[1].Address.String(),
						Coins:   toDropletString(totalCoins + 1e3),
						Hours:   "1",
					},
				},
			},
			err:  "400 Bad Request - balance is not sufficient\n",
			code: http.StatusBadRequest,
		},

		{
			// NOTE: expects wallet to have multiple outputs with non-zero coins
			name: "insufficient hours with uxouts",
			req: api.CreateTransactionRequest{
				HoursSelection: api.HoursSelection{
					Type: wallet.HoursSelectionTypeManual,
				},
				Wallet: api.CreateTransactionRequestWallet{
					ID:       w.Filename(),
					Password: password,
					UxOuts:   []string{walletOutputs[0].Hash},
				},
				ChangeAddress: &defaultChangeAddress,
				To: []api.Receiver{
					{
						Address: w.Entries[1].Address.String(),
						Coins:   toDropletString(1e3),
						Hours:   fmt.Sprint(totalHours + 1),
					},
				},
			},
			err:  "400 Bad Request - hours are not sufficient\n",
			code: http.StatusBadRequest,
		},

		{
			name: "valid request, uxouts specified",
			req: api.CreateTransactionRequest{
				HoursSelection: api.HoursSelection{
					Type: wallet.HoursSelectionTypeManual,
				},
				Wallet: api.CreateTransactionRequestWallet{
					ID:       w.Filename(),
					Password: password,
					// NOTE: all uxouts are provided, which has the same behavior as
					// not providing any uxouts or addresses.
					// Using a subset of uxouts makes the wallet setup very
					// difficult, especially to make deterministic, in the live test
					// More complex cases should be covered by unit tests
					UxOuts: walletOutputHashes,
				},
				ChangeAddress: &defaultChangeAddress,
				To: []api.Receiver{
					{
						Address: w.Entries[1].Address.String(),
						Coins:   toDropletString(totalCoins - 1e3),
						Hours:   "1",
					},
				},
			},
			outputs: []coin.TransactionOutput{
				{
					Address: w.Entries[1].Address,
					Coins:   totalCoins - 1e3,
					Hours:   1,
				},
				{
					Address: w.Entries[0].Address,
					Coins:   1e3,
					Hours:   remainingHours - 1,
				},
			},
			additionalRespVerify: func(t *testing.T, r *api.CreateTransactionResponse) {
				require.Equal(t, len(walletOutputHashes), len(r.Transaction.In))
			},
		},

		{
			name: "specified addresses not in wallet",
			req: api.CreateTransactionRequest{
				HoursSelection: api.HoursSelection{
					Type: wallet.HoursSelectionTypeManual,
				},
				Wallet: api.CreateTransactionRequestWallet{
					ID:        w.Filename(),
					Password:  password,
					Addresses: []string{testutil.MakeAddress().String()},
				},
				ChangeAddress: &defaultChangeAddress,
				To: []api.Receiver{
					{
						Address: w.Entries[1].Address.String(),
						Coins:   toDropletString(totalCoins),
						Hours:   "1",
					},
				},
			},
			err:  "400 Bad Request - address not found in wallet\n",
			code: http.StatusBadRequest,
		},

		{
			name: "valid request, addresses specified",
			req: api.CreateTransactionRequest{
				HoursSelection: api.HoursSelection{
					Type: wallet.HoursSelectionTypeManual,
				},
				Wallet: api.CreateTransactionRequestWallet{
					ID:       w.Filename(),
					Password: password,
					// NOTE: all addresses are provided, which has the same behavior as
					// not providing any addresses.
					// Using a subset of addresses makes the wallet setup very
					// difficult, especially to make deterministic, in the live test
					// More complex cases should be covered by unit tests
					Addresses: addresses,
				},
				ChangeAddress: &defaultChangeAddress,
				To: []api.Receiver{
					{
						Address: w.Entries[1].Address.String(),
						Coins:   toDropletString(totalCoins - 1e3),
						Hours:   "1",
					},
				},
			},
			outputs: []coin.TransactionOutput{
				{
					Address: w.Entries[1].Address,
					Coins:   totalCoins - 1e3,
					Hours:   1,
				},
				{
					Address: w.Entries[0].Address,
					Coins:   1e3,
					Hours:   remainingHours - 1,
				},
			},
		},
	}

	if w.IsEncrypted() {
		cases = append(cases, testCase{
			name: "invalid password",
			req: api.CreateTransactionRequest{
				HoursSelection: api.HoursSelection{
					Type: wallet.HoursSelectionTypeManual,
				},
				Wallet: api.CreateTransactionRequestWallet{
					ID:       w.Filename(),
					Password: password + "foo",
				},
				ChangeAddress: &defaultChangeAddress,
				To: []api.Receiver{
					{
						Address: w.Entries[0].Address.String(),
						Coins:   "1000",
						Hours:   "1",
					},
				},
			},
			err:  "401 Unauthorized - invalid password\n",
			code: http.StatusUnauthorized,
		})

		cases = append(cases, testCase{
			name: "password not provided",
			req: api.CreateTransactionRequest{
				HoursSelection: api.HoursSelection{
					Type: wallet.HoursSelectionTypeManual,
				},
				Wallet: api.CreateTransactionRequestWallet{
					ID:       w.Filename(),
					Password: "",
				},
				ChangeAddress: &defaultChangeAddress,
				To: []api.Receiver{
					{
						Address: w.Entries[0].Address.String(),
						Coins:   "1000",
						Hours:   "1",
					},
				},
			},
			err:  "400 Bad Request - missing password\n",
			code: http.StatusBadRequest,
		})

	} else {
		cases = append(cases, testCase{
			name: "password provided for unencrypted wallet",
			req: api.CreateTransactionRequest{
				HoursSelection: api.HoursSelection{
					Type: wallet.HoursSelectionTypeManual,
				},
				Wallet: api.CreateTransactionRequestWallet{
					ID:       w.Filename(),
					Password: password + "foo",
				},
				ChangeAddress: &defaultChangeAddress,
				To: []api.Receiver{
					{
						Address: w.Entries[0].Address.String(),
						Coins:   "1000",
						Hours:   "1",
					},
				},
			},
			err:  "400 Bad Request - wallet is not encrypted\n",
			code: http.StatusBadRequest,
		})
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			require.False(t, len(tc.outputs) != 0 && len(tc.outputsSubset) != 0, "outputs and outputsSubset can't both be set")

			result, err := c.CreateTransaction(tc.req)
			if tc.err != "" {
				assertResponseError(t, err, tc.code, tc.err)
				return
			}

			require.NoError(t, err)

			d, err := json.MarshalIndent(result, "", "    ")
			require.NoError(t, err)
			fmt.Println(string(d))

			if len(tc.outputsSubset) == 0 {
				require.Equal(t, len(tc.outputs), len(result.Transaction.Out))
			}

			for i, o := range tc.outputs {
				// The final change output may not have the address specified,
				// if the ChangeAddress was not specified in the wallet params.
				// Calculate it automatically based upon the transaction inputs
				if o.Address.Null() {
					require.Equal(t, i, len(tc.outputs)-1)
					require.Nil(t, tc.req.ChangeAddress)

					changeAddr := result.Transaction.Out[i].Address
					// The changeAddr must be associated with one of the transaction inputs
					changeAddrFound := false
					for _, x := range result.Transaction.In {
						require.NotNil(t, x.Address)
						if changeAddr == x.Address {
							changeAddrFound = true
							break
						}
					}

					require.True(t, changeAddrFound)
				} else {
					require.Equal(t, o.Address.String(), result.Transaction.Out[i].Address)
				}

				coins, err := droplet.FromString(result.Transaction.Out[i].Coins)
				require.NoError(t, err)
				require.Equal(t, o.Coins, coins, "[%d] %d != %d", i, o.Coins, coins)

				if !tc.ignoreHours {
					hours, err := strconv.ParseUint(result.Transaction.Out[i].Hours, 10, 64)
					require.NoError(t, err)
					require.Equal(t, o.Hours, hours, "[%d] %d != %d", i, o.Hours, hours)
				}
			}

			assertEncodeTxnMatchesTxn(t, result)
			assertRequestedCoins(t, tc.req.To, result.Transaction.Out)
			assertCreatedTransactionValid(t, result.Transaction)

			if tc.req.HoursSelection.Type == wallet.HoursSelectionTypeManual {
				assertRequestedHours(t, tc.req.To, result.Transaction.Out)
			}

			if tc.additionalRespVerify != nil {
				tc.additionalRespVerify(t, result)
			}
		})
	}
}

func TestLiveWalletCreateTransactionRandom(t *testing.T) {
	if !doLive(t) {
		return
	}

	requireWalletEnv(t)

	c := api.NewClient(nodeAddress())

	w, totalCoins, totalHours, password := prepareAndCheckWallet(t, c, 2e6, 20)

	if w.IsEncrypted() {
		t.Skip("Skipping TestLiveWalletCreateTransactionRandom tests with encrypted wallet")
		return
	}

	remainingHours := fee.RemainingHours(totalHours)
	require.True(t, remainingHours > 1)

	assertTxnOutputCount := func(t *testing.T, changeAddress string, nOutputs int, result *api.CreateTransactionResponse) {
		nResultOutputs := len(result.Transaction.Out)
		require.True(t, nResultOutputs == nOutputs || nResultOutputs == nOutputs+1)
		hasChange := nResultOutputs == nOutputs+1
		changeOutput := result.Transaction.Out[nResultOutputs-1]
		if hasChange {
			require.Equal(t, changeOutput.Address, changeAddress)
		}

		t.Log("hasChange", hasChange)
		if hasChange {
			t.Log("changeCoins", changeOutput.Coins)
			t.Log("changeHours", changeOutput.Hours)
		}
	}

	iterations := 250
	maxOutputs := 10
	destAddrs := make([]cipher.Address, maxOutputs)
	for i := range destAddrs {
		destAddrs[i] = testutil.MakeAddress()
	}

	for i := 0; i < iterations; i++ {
		t.Log("iteration", i)
		t.Log("totalCoins", totalCoins)
		t.Log("totalHours", totalHours)

		spendableHours := fee.RemainingHours(totalHours)
		t.Log("spendableHours", spendableHours)

		coins := rand.Intn(int(totalCoins)) + 1
		coins -= coins % int(visor.MaxDropletDivisor())
		if coins == 0 {
			coins = int(visor.MaxDropletDivisor())
		}
		hours := rand.Intn(int(spendableHours + 1))
		nOutputs := rand.Intn(maxOutputs) + 1

		t.Log("sendCoins", coins)
		t.Log("sendHours", hours)

		changeAddress := w.Entries[0].Address.String()

		shareFactor := strconv.FormatFloat(rand.Float64(), 'f', 8, 64)

		t.Log("shareFactor", shareFactor)

		to := make([]api.Receiver, 0, nOutputs)
		remainingHours := hours
		remainingCoins := coins
		for i := 0; i < nOutputs; i++ {
			if remainingCoins == 0 {
				break
			}

			receiver := api.Receiver{}
			receiver.Address = destAddrs[rand.Intn(len(destAddrs))].String()

			if i == nOutputs-1 {
				var err error
				receiver.Coins, err = droplet.ToString(uint64(remainingCoins))
				require.NoError(t, err)
				receiver.Hours = fmt.Sprint(remainingHours)

				remainingCoins = 0
				remainingHours = 0
			} else {
				receiverCoins := rand.Intn(remainingCoins) + 1
				receiverCoins -= receiverCoins % int(visor.MaxDropletDivisor())
				if receiverCoins == 0 {
					receiverCoins = int(visor.MaxDropletDivisor())
				}

				var err error
				receiver.Coins, err = droplet.ToString(uint64(receiverCoins))
				require.NoError(t, err)
				remainingCoins -= receiverCoins

				receiverHours := rand.Intn(remainingHours + 1)
				receiver.Hours = fmt.Sprint(receiverHours)
				remainingHours -= receiverHours
			}

			to = append(to, receiver)
		}

		// Remove duplicate outputs
		dup := make(map[api.Receiver]struct{}, len(to))
		newTo := make([]api.Receiver, 0, len(dup))
		for _, o := range to {
			if _, ok := dup[o]; !ok {
				dup[o] = struct{}{}
				newTo = append(newTo, o)
			}
		}
		to = newTo

		nOutputs = len(to)
		t.Log("nOutputs", nOutputs)

		rand.Shuffle(len(to), func(i, j int) {
			to[i], to[j] = to[j], to[i]
		})

		for i, o := range to {
			t.Logf("to[%d].Hours %s\n", i, o.Hours)
		}

		autoTo := make([]api.Receiver, len(to))
		for i, o := range to {
			autoTo[i] = api.Receiver{
				Address: o.Address,
				Coins:   o.Coins,
				Hours:   "",
			}
		}

		// Remove duplicate outputs
		dup = make(map[api.Receiver]struct{}, len(autoTo))
		newAutoTo := make([]api.Receiver, 0, len(dup))
		for _, o := range autoTo {
			if _, ok := dup[o]; !ok {
				dup[o] = struct{}{}
				newAutoTo = append(newAutoTo, o)
			}
		}
		autoTo = newAutoTo

		nAutoOutputs := len(autoTo)
		t.Log("nAutoOutputs", nAutoOutputs)

		for i, o := range autoTo {
			t.Logf("autoTo[%d].Coins %s\n", i, o.Coins)
		}

		// Auto, random share factor

		result, err := c.CreateTransaction(api.CreateTransactionRequest{
			HoursSelection: api.HoursSelection{
				Type:        wallet.HoursSelectionTypeAuto,
				Mode:        wallet.HoursSelectionModeShare,
				ShareFactor: shareFactor,
			},
			ChangeAddress: &changeAddress,
			Wallet: api.CreateTransactionRequestWallet{
				ID:       w.Filename(),
				Password: password,
			},
			To: autoTo,
		})
		require.NoError(t, err)

		assertEncodeTxnMatchesTxn(t, result)
		assertTxnOutputCount(t, changeAddress, nAutoOutputs, result)
		assertRequestedCoins(t, autoTo, result.Transaction.Out)
		assertCreatedTransactionValid(t, result.Transaction)

		// Auto, share factor 0

		result, err = c.CreateTransaction(api.CreateTransactionRequest{
			HoursSelection: api.HoursSelection{
				Type:        wallet.HoursSelectionTypeAuto,
				Mode:        wallet.HoursSelectionModeShare,
				ShareFactor: "0",
			},
			ChangeAddress: &changeAddress,
			Wallet: api.CreateTransactionRequestWallet{
				ID:       w.Filename(),
				Password: password,
			},
			To: autoTo,
		})
		require.NoError(t, err)

		assertEncodeTxnMatchesTxn(t, result)
		assertTxnOutputCount(t, changeAddress, nAutoOutputs, result)
		assertRequestedCoins(t, autoTo, result.Transaction.Out)
		assertCreatedTransactionValid(t, result.Transaction)

		// Check that the non-change outputs have 0 hours
		for _, o := range result.Transaction.Out[:nAutoOutputs] {
			require.Equal(t, "0", o.Hours)
		}

		// Auto, share factor 1

		result, err = c.CreateTransaction(api.CreateTransactionRequest{
			HoursSelection: api.HoursSelection{
				Type:        wallet.HoursSelectionTypeAuto,
				Mode:        wallet.HoursSelectionModeShare,
				ShareFactor: "1",
			},
			ChangeAddress: &changeAddress,
			Wallet: api.CreateTransactionRequestWallet{
				ID:       w.Filename(),
				Password: password,
			},
			To: autoTo,
		})
		require.NoError(t, err)

		assertEncodeTxnMatchesTxn(t, result)
		assertTxnOutputCount(t, changeAddress, nAutoOutputs, result)
		assertRequestedCoins(t, autoTo, result.Transaction.Out)
		assertCreatedTransactionValid(t, result.Transaction)

		// Check that the change output has 0 hours
		if len(result.Transaction.Out) > nAutoOutputs {
			require.Equal(t, "0", result.Transaction.Out[nAutoOutputs].Hours)
		}

		// Manual

		result, err = c.CreateTransaction(api.CreateTransactionRequest{
			HoursSelection: api.HoursSelection{
				Type: wallet.HoursSelectionTypeManual,
			},
			ChangeAddress: &changeAddress,
			Wallet: api.CreateTransactionRequestWallet{
				ID:       w.Filename(),
				Password: password,
			},
			To: to,
		})
		require.NoError(t, err)

		assertEncodeTxnMatchesTxn(t, result)
		assertTxnOutputCount(t, changeAddress, nOutputs, result)
		assertRequestedCoins(t, to, result.Transaction.Out)
		assertRequestedHours(t, to, result.Transaction.Out)
		assertCreatedTransactionValid(t, result.Transaction)
	}
}

func assertEncodeTxnMatchesTxn(t *testing.T, result *api.CreateTransactionResponse) {
	require.NotEmpty(t, result.EncodedTransaction)
	emptyTxn := &coin.Transaction{}
	require.NotEqual(t, hex.EncodeToString(emptyTxn.Serialize()), result.EncodedTransaction)
	txn, err := result.Transaction.ToTransaction()
	require.NoError(t, err)

	serializedTxn := txn.Serialize()
	require.Equal(t, hex.EncodeToString(serializedTxn), result.EncodedTransaction)

	require.Equal(t, int(txn.Length), len(serializedTxn))
}

func assertRequestedCoins(t *testing.T, to []api.Receiver, out []api.CreatedTransactionOutput) {
	var requestedCoins uint64
	for _, o := range to {
		c, err := droplet.FromString(o.Coins)
		require.NoError(t, err)
		requestedCoins += c
	}

	var sentCoins uint64
	for _, o := range out[:len(to)] { // exclude change output
		c, err := droplet.FromString(o.Coins)
		require.NoError(t, err)
		sentCoins += c
	}

	require.Equal(t, requestedCoins, sentCoins)
}

func assertRequestedHours(t *testing.T, to []api.Receiver, out []api.CreatedTransactionOutput) {
	for i, o := range out[:len(to)] { // exclude change output
		toHours, err := strconv.ParseUint(to[i].Hours, 10, 64)
		require.NoError(t, err)

		outHours, err := strconv.ParseUint(o.Hours, 10, 64)

		require.Equal(t, toHours, outHours)
	}
}

func assertCreatedTransactionValid(t *testing.T, r api.CreatedTransaction) {
	require.NotEmpty(t, r.In)
	require.NotEmpty(t, r.Out)

	fee, err := strconv.ParseUint(r.Fee, 10, 64)
	require.NoError(t, err)

	require.NotEqual(t, uint64(0), fee)

	var inputHours uint64
	var inputCoins uint64
	for _, in := range r.In {
		require.NotNil(t, in.CalculatedHours)
		calculatedHours, err := strconv.ParseUint(in.CalculatedHours, 10, 64)
		require.NoError(t, err)
		inputHours, err = coin.AddUint64(inputHours, calculatedHours)
		require.NoError(t, err)

		require.NotNil(t, in.Hours)
		hours, err := strconv.ParseUint(in.Hours, 10, 64)
		require.NoError(t, err)

		require.True(t, hours <= calculatedHours)

		require.NotNil(t, in.Coins)
		coins, err := droplet.FromString(in.Coins)
		require.NoError(t, err)
		inputCoins, err = coin.AddUint64(inputCoins, coins)
		require.NoError(t, err)
	}

	var outputHours uint64
	var outputCoins uint64
	for _, out := range r.Out {
		hours, err := strconv.ParseUint(out.Hours, 10, 64)
		require.NoError(t, err)
		outputHours, err = coin.AddUint64(outputHours, hours)
		require.NoError(t, err)

		coins, err := droplet.FromString(out.Coins)
		require.NoError(t, err)
		outputCoins, err = coin.AddUint64(outputCoins, coins)
		require.NoError(t, err)
	}

	require.True(t, inputHours > outputHours)
	require.Equal(t, inputHours-outputHours, fee)

	require.Equal(t, inputCoins, outputCoins)

	require.Equal(t, uint8(0), r.Type)
	require.NotEmpty(t, r.Length)
}

func TestCreateWallet(t *testing.T) {
	if !doLiveOrStable(t) {
		return
	}

	c := api.NewClient(nodeAddress())

	w, seed, clean := createWallet(t, c, false, "", "")
	defer clean()
	require.False(t, w.Meta.Encrypted)

	walletDir := getWalletDir(t, c)

	// Confirms the wallet does exist
	walletPath := filepath.Join(walletDir, w.Meta.Filename)
	_, err := os.Stat(walletPath)
	require.NoError(t, err)

	// Loads the wallet and confirms that the wallet has the same seed
	lw, err := wallet.Load(walletPath)
	require.NoError(t, err)
	require.False(t, lw.IsEncrypted())
	require.Equal(t, seed, lw.Meta["seed"])
	require.Equal(t, len(w.Entries), len(lw.Entries))

	for i := range w.Entries {
		require.Equal(t, w.Entries[i].Address, lw.Entries[i].Address.String())
		require.Equal(t, w.Entries[i].Public, lw.Entries[i].Public.Hex())
	}

	// Creates wallet with encryption
	encW, _, encWClean := createWallet(t, c, true, "pwd", "")
	defer encWClean()
	require.True(t, encW.Meta.Encrypted)

	walletPath = filepath.Join(walletDir, encW.Meta.Filename)
	encLW, err := wallet.Load(walletPath)
	require.NoError(t, err)

	// Confirms the loaded wallet is encrypted and has the same address entries
	require.True(t, encLW.IsEncrypted())
	require.Equal(t, len(encW.Entries), len(encLW.Entries))

	for i := range encW.Entries {
		require.Equal(t, encW.Entries[i].Address, encLW.Entries[i].Address.String())
		require.Equal(t, encW.Entries[i].Public, encLW.Entries[i].Public.Hex())
	}
}

func TestGetWallet(t *testing.T) {
	if !doLiveOrStable(t) {
		return
	}

	c := api.NewClient(nodeAddress())

	// Create a wallet
	w, _, clean := createWallet(t, c, false, "", "")
	defer clean()

	// Confirms the wallet can be acquired
	w1, err := c.Wallet(w.Meta.Filename)
	require.NoError(t, err)
	require.Equal(t, *w, *w1)
}

func TestGetWallets(t *testing.T) {
	if !doLiveOrStable(t) {
		return
	}

	c := api.NewClient(nodeAddress())

	// Creates 2 new wallets
	var ws []api.WalletResponse
	for i := 0; i < 2; i++ {
		w, _, clean := createWallet(t, c, false, "", "")
		defer clean()
		// cleaners = append(cleaners, clean)
		ws = append(ws, *w)
	}

	// Gets wallet from node
	wlts, err := c.Wallets()
	require.NoError(t, err)

	// Create the wallet map
	walletMap := make(map[string]api.WalletResponse)
	for _, w := range wlts {
		walletMap[w.Meta.Filename] = *w
	}

	// Confirms the returned wallets contains the wallet we created.
	for _, w := range ws {
		retW, ok := walletMap[w.Meta.Filename]
		require.True(t, ok)
		require.Equal(t, w, retW)
	}
}

// TestWalletNewAddress will generate 30 wallets for testing, and they will
// be removed automatically after testing.
func TestWalletNewAddress(t *testing.T) {
	if !doLiveOrStable(t) {
		return
	}

	// We only test 30 cases, cause the more addresses we generate, the longer
	// it takes, we don't want to spend much time here.
	for i := 1; i <= 30; i++ {
		name := fmt.Sprintf("generate %v addresses", i)
		t.Run(name, func(t *testing.T) {
			c := api.NewClient(nodeAddress())
			var encrypt bool
			var password string
			// Test wallet with encryption only when i == 2, so that
			// the tests won't time out.
			if i == 2 {
				encrypt = true
				password = "pwd"
			}

			w, seed, clean := createWallet(t, c, encrypt, password, "")
			defer clean()

			addrs, err := c.NewWalletAddress(w.Meta.Filename, i, password)
			if err != nil {
				t.Fatalf("%v", err)
				return
			}
			require.NoError(t, err)

			seckeys := cipher.GenerateDeterministicKeyPairs([]byte(seed), i+1)
			var as []string
			for _, k := range seckeys {
				as = append(as, cipher.AddressFromSecKey(k).String())
			}

			// Confirms thoses new generated addresses are the same.
			require.Equal(t, len(addrs), len(as)-1)
			for i := range addrs {
				require.Equal(t, as[i+1], addrs[i])
			}
		})
	}
}

func TestStableWalletBalance(t *testing.T) {
	if !doStable(t) {
		return
	}

	c := api.NewClient(nodeAddress())
	w, _, clean := createWallet(t, c, false, "", "casino away claim road artist where blossom warrior demise royal still palm")
	defer clean()

	bp, err := c.WalletBalance(w.Meta.Filename)
	require.NoError(t, err)

	var expect api.BalanceResponse
	checkGoldenFile(t, "wallet-balance.golden", TestData{*bp, &expect})
}

func TestLiveWalletBalance(t *testing.T) {
	if !doLive(t) {
		return
	}

	requireWalletEnv(t)

	c := api.NewClient(nodeAddress())
	_, walletName, _ := getWalletFromEnv(t, c)
	bp, err := c.WalletBalance(walletName)
	require.NoError(t, err)
	require.NotNil(t, bp)
	require.NotNil(t, bp.Addresses)
}

func TestWalletUpdate(t *testing.T) {
	if !doLiveOrStable(t) {
		return
	}

	c := api.NewClient(nodeAddress())
	w, _, clean := createWallet(t, c, false, "", "")
	defer clean()

	err := c.UpdateWallet(w.Meta.Filename, "new wallet")
	require.NoError(t, err)

	// Confirms the wallet has label of "new wallet"
	w1, err := c.Wallet(w.Meta.Filename)
	require.NoError(t, err)
	require.Equal(t, w1.Meta.Label, "new wallet")
}

func TestStableWalletTransactions(t *testing.T) {
	if !doStable(t) {
		return
	}

	c := api.NewClient(nodeAddress())
	w, _, clean := createWallet(t, c, false, "", "")
	defer clean()

	txns, err := c.WalletTransactions(w.Meta.Filename)
	require.NoError(t, err)

	var expect api.UnconfirmedTxnsResponse
	checkGoldenFile(t, "wallet-transactions.golden", TestData{*txns, &expect})
}

func TestLiveWalletTransactions(t *testing.T) {
	if !doLive(t) {
		return
	}

	requireWalletEnv(t)

	c := api.NewClient(nodeAddress())
	w, _, _, _ := prepareAndCheckWallet(t, c, 1e6, 1)
	txns, err := c.WalletTransactions(w.Filename())
	require.NoError(t, err)

	bp, err := c.WalletBalance(w.Filename())
	require.NoError(t, err)
	// There's pending transactions if predicted coins are not the same as confirmed coins
	if bp.Predicted.Coins != bp.Confirmed.Coins {
		require.NotEmpty(t, txns.Transactions)
		return
	}

	require.Empty(t, txns.Transactions)
}

func TestWalletFolderName(t *testing.T) {
	if !doLiveOrStable(t) {
		return
	}

	c := api.NewClient(nodeAddress())
	folderName, err := c.WalletFolderName()
	require.NoError(t, err)

	require.NotNil(t, folderName)
	require.NotEmpty(t, folderName.Address)
}

func TestEncryptWallet(t *testing.T) {
	if !doLiveOrStable(t) {
		return
	}

	c := api.NewClient(nodeAddress())

	// Create a unencrypted wallet
	w, _, clean := createWallet(t, c, false, "", "")
	defer clean()

	// Encrypts the wallet
	rlt, err := c.EncryptWallet(w.Meta.Filename, "pwd")
	require.NoError(t, err)
	require.NotEmpty(t, rlt.Meta.CryptoType)
	require.True(t, rlt.Meta.Encrypted)

	//  Encrypt the wallet again, should returns error
	_, err = c.EncryptWallet(w.Meta.Filename, "pwd")
	assertResponseError(t, err, http.StatusBadRequest, "400 Bad Request - wallet is encrypted\n")

	// Confirms that no sensitive data do exist in wallet file
	wf, err := c.WalletFolderName()
	require.NoError(t, err)
	wltPath := filepath.Join(wf.Address, w.Meta.Filename)
	lw, err := wallet.Load(wltPath)
	require.NoError(t, err)
	require.Empty(t, lw.Meta["seed"])
	require.Empty(t, lw.Meta["lastSeed"])
	require.NotEmpty(t, lw.Meta["secrets"])

	// Decrypts the wallet, and confirms that the
	// seed and address entries are the same as it was before being encrypted.
	dw, err := c.DecryptWallet(w.Meta.Filename, "pwd")
	require.NoError(t, err)
	require.Equal(t, w, dw)
}

func TestDecryptWallet(t *testing.T) {
	if !doLiveOrStable(t) {
		return
	}

	c := api.NewClient(nodeAddress())
	w, seed, clean := createWallet(t, c, true, "pwd", "")
	defer clean()

	// Decrypt wallet with different password, must fail
	_, err := c.DecryptWallet(w.Meta.Filename, "pwd1")
	assertResponseError(t, err, http.StatusUnauthorized, "401 Unauthorized - invalid password\n")

	// Decrypt wallet with no password, must fail
	_, err = c.DecryptWallet(w.Meta.Filename, "")
	assertResponseError(t, err, http.StatusBadRequest, "400 Bad Request - missing password\n")

	// Decrypts wallet with correct password
	dw, err := c.DecryptWallet(w.Meta.Filename, "pwd")
	require.NoError(t, err)

	// Confirms that no sensitive data are returned
	require.Empty(t, dw.Meta.CryptoType)
	require.False(t, dw.Meta.Encrypted)

	// Loads wallet from file
	wf, err := c.WalletFolderName()
	require.NoError(t, err)
	wltPath := filepath.Join(wf.Address, w.Meta.Filename)
	lw, err := wallet.Load(wltPath)
	require.NoError(t, err)

	require.Equal(t, lw.Meta["seed"], seed)
	require.Len(t, lw.Entries, 1)

	// Confirms the last seed is matched
	lseed, seckeys := cipher.GenerateDeterministicKeyPairsSeed([]byte(seed), 1)
	require.Equal(t, hex.EncodeToString(lseed), lw.Meta["lastSeed"])

	// Confirms that the first address is derivied from the private key
	pubkey := cipher.PubKeyFromSecKey(seckeys[0])
	require.Equal(t, w.Entries[0].Address, cipher.AddressFromPubKey(pubkey).String())
	require.Equal(t, lw.Entries[0].Address.String(), w.Entries[0].Address)
}

func TestGetWalletSeedDisabledAPI(t *testing.T) {
	if !doDisableSeedApi(t) {
		return
	}

	c := api.NewClient(nodeAddress())

	// Create an encrypted wallet
	w, _, clean := createWallet(t, c, true, "pwd", "")
	defer clean()

	_, err := c.GetWalletSeed(w.Meta.Filename, "pwd")
	assertResponseError(t, err, http.StatusForbidden, "403 Forbidden\n")
}

func TestGetWalletSeedEnabledAPI(t *testing.T) {
	if !doLiveOrStable(t) {
		return
	}

	c := api.NewClient(nodeAddress())

	// Create an encrypted wallet
	w, seed, clean := createWallet(t, c, true, "pwd", "")
	defer clean()

	require.NotEmpty(t, seed)

	sd, err := c.GetWalletSeed(w.Meta.Filename, "pwd")
	require.NoError(t, err)

	// Confirms the seed are matched
	require.Equal(t, seed, sd)

	// Get seed of wrong wallet id
	_, err = c.GetWalletSeed("w.wlt", "pwd")
	assertResponseError(t, err, http.StatusNotFound, "404 Not Found\n")

	// Check with invalid password
	_, err = c.GetWalletSeed(w.Meta.Filename, "wrong password")
	assertResponseError(t, err, http.StatusUnauthorized, "401 Unauthorized - invalid password\n")

	// Check with missing password
	_, err = c.GetWalletSeed(w.Meta.Filename, "")
	assertResponseError(t, err, http.StatusBadRequest, "400 Bad Request - missing password\n")

	// Create unencrypted wallet to check against
	nw, _, nclean := createWallet(t, c, false, "", "")
	defer nclean()
	_, err = c.GetWalletSeed(nw.Meta.Filename, "pwd")
	assertResponseError(t, err, http.StatusBadRequest, "400 Bad Request - wallet is not encrypted\n")
}

// prepareAndCheckWallet gets wallet from environment, and confirms:
// 1. The minimal coins and coin hours requirements are met.
// 2. The wallet has at least two address entry.
// Returns the loaded wallet, total coins, total coin hours and password of the wallet.
func prepareAndCheckWallet(t *testing.T, c *api.Client, miniCoins, miniCoinHours uint64) (*wallet.Wallet, uint64, uint64, string) {
	walletDir, walletName, password := getWalletFromEnv(t, c)
	walletPath := filepath.Join(walletDir, walletName)

	// Checks if the wallet does exist
	if _, err := os.Stat(walletPath); os.IsNotExist(err) {
		t.Fatalf("Wallet %v doesn't exist", walletPath)
	}

	w, err := wallet.Load(walletPath)
	if err != nil {
		t.Fatalf("Load wallet %v failed: %v", walletPath, err)
	}

	if w.IsEncrypted() && password == "" {
		t.Fatalf("Wallet is encrypted, must set WALLET_PASSWORD env var")
	}

	// Generate more addresses if address entries less than 2.
	if len(w.Entries) < 2 {
		_, err := c.NewWalletAddress(w.Filename(), 2-len(w.Entries), password)
		if err != nil {
			t.Fatalf("New wallet address failed: %v", err)
		}

		w, err = wallet.Load(walletPath)
		if err != nil {
			t.Fatalf("Reload wallet %v failed: %v", walletPath, err)
		}
	}

	coins, hours := getWalletBalance(t, c, walletName)
	if coins < miniCoins {
		t.Fatalf("Wallet must have at least %d coins", miniCoins)
	}

	if hours < miniCoinHours {
		t.Fatalf("Wallet must have at least %d coin hours", miniCoinHours)
	}

	if err := w.Save(walletDir); err != nil {
		t.Fatalf("%v", err)
	}

	return w, coins, hours, password
}

// getWalletFromEnv loads wallet from envrionment variables.
// Returns wallet dir, wallet name and wallet password is any.
func getWalletFromEnv(t *testing.T, c *api.Client) (string, string, string) {
	walletDir := getWalletDir(t, c)

	walletName := os.Getenv("WALLET_NAME")
	if walletName == "" {
		t.Fatal("Missing WALLET_NAME environment value")
	}

	walletPassword := os.Getenv("WALLET_PASSWORD")
	return walletDir, walletName, walletPassword
}

func requireWalletEnv(t *testing.T) {
	if !doLiveWallet(t) {
		return
	}

	walletName := os.Getenv("WALLET_NAME")
	if walletName == "" {
		t.Fatal("missing WALLET_NAME environment value")
	}
}

// getWalletBalance gets wallet balance.
// Returns coins and hours
func getWalletBalance(t *testing.T, c *api.Client, walletName string) (uint64, uint64) {
	wp, err := c.WalletBalance(walletName)
	if err != nil {
		t.Fatalf("Get wallet balance of %v failed: %v", walletName, err)
	}

	return wp.Confirmed.Coins, wp.Confirmed.Hours
}

func getTransaction(t *testing.T, c *api.Client, txid string) *daemon.TransactionResult {
	tx, err := c.Transaction(txid)
	if err != nil {
		t.Fatalf("%v", err)
	}

	return tx
}

// getAddressBalance gets balance of given address.
// Returns coins and coin hours.
func getAddressBalance(t *testing.T, c *api.Client, addr string) (uint64, uint64) { // nolint: unparam
	bp, err := c.Balance([]string{addr})
	if err != nil {
		t.Fatalf("%v", err)
	}
	return bp.Confirmed.Coins, bp.Confirmed.Hours
}

func checkNoSensitiveData(t *testing.T, w *wallet.Wallet) {
	require.Empty(t, w.Meta["seed"])
	require.Empty(t, w.Meta["lastSeed"])
	require.Empty(t, w.Meta["secrets"])
	for _, e := range w.Entries {
		require.Equal(t, cipher.SecKey{}, e.Secret)
	}
}

// checkWalletEntriesAndLastSeed confirms the wallet entries and lastSeed are derivied
// from the seed.
func checkWalletEntriesAndLastSeed(t *testing.T, w *wallet.Wallet) {
	seed, ok := w.Meta["seed"]
	require.True(t, ok)
	newSeed, seckeys := cipher.GenerateDeterministicKeyPairsSeed([]byte(seed), len(w.Entries))
	require.Len(t, seckeys, len(w.Entries))
	for i, sk := range seckeys {
		require.Equal(t, w.Entries[i].Secret, sk)
		pk := cipher.PubKeyFromSecKey(sk)
		require.Equal(t, w.Entries[i].Public, pk)
	}
	lastSeed, ok := w.Meta["lastSeed"]
	require.True(t, ok)
	require.Equal(t, lastSeed, hex.EncodeToString(newSeed))
}

// createWallet creates a wallet with rand seed.
// Returns the generated wallet, seed and clean up function.
func createWallet(t *testing.T, c *api.Client, encrypt bool, password string, seed string) (*api.WalletResponse, string, func()) {
	if seed == "" {
		seed = hex.EncodeToString(cipher.RandByte(32))
	}
	// Use the first 6 letter of the seed as label.
	var w *api.WalletResponse
	var err error
	if encrypt {
		w, err = c.CreateEncryptedWallet(seed, seed[:6], password, 0)
	} else {
		w, err = c.CreateUnencryptedWallet(seed, seed[:6], 0)
	}

	require.NoError(t, err)

	walletDir := getWalletDir(t, c)

	return w, seed, func() {
		// Cleaner function to delete the wallet and bak wallet
		walletPath := filepath.Join(walletDir, w.Meta.Filename)
		err = os.Remove(walletPath)
		require.NoError(t, err)

		bakWalletPath := walletPath + ".bak"
		if _, err := os.Stat(bakWalletPath); !os.IsNotExist(err) {
			// Return directly if no .bak file does exist
			err = os.Remove(bakWalletPath)
		}

		require.NoError(t, err)

		// Removes the wallet from memory
		c.UnloadWallet(w.Meta.Filename)
	}
}

func getWalletDir(t *testing.T, c *api.Client) string {
	wf, err := c.WalletFolderName()
	if err != nil {
		t.Fatalf("%v", err)
	}
	return wf.Address
}

func TestDisableWalletApi(t *testing.T) {
	if !doDisableWalletApi(t) {
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
			expectErr: "403 Forbidden\n",
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
			expectErr: "403 Forbidden\n",
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
			expectErr: "403 Forbidden\n",
			code:      http.StatusForbidden,
		},
		{
			name:      "get wallet balance",
			method:    http.MethodGet,
			endpoint:  "/api/v1/wallet/balance?id=test.wlt",
			expectErr: "403 Forbidden\n",
			code:      http.StatusForbidden,
		},
		{
			name:     "wallet spending",
			method:   http.MethodPost,
			endpoint: "/api/v1/wallet/spend",
			body: func() io.Reader {
				v := url.Values{}
				v.Add("id", "test.wlt")
				v.Add("coins", "100000") // 1e5
				v.Add("dst", "2jBbGxZRGoQG1mqhPBnXnLTxK6oxsTf8os6")
				return strings.NewReader(v.Encode())
			},
			expectErr: "403 Forbidden\n",
			code:      http.StatusForbidden,
		},
		{
			name:      "get wallet unconfirmed transactions",
			method:    http.MethodGet,
			endpoint:  "/api/v1/wallet/transactions?id=test.wlt",
			expectErr: "403 Forbidden\n",
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
			expectErr: "403 Forbidden\n",
			code:      http.StatusForbidden,
		},
		{
			name:      "new seed",
			method:    http.MethodGet,
			endpoint:  "/api/v1/wallet/newSeed",
			expectErr: "403 Forbidden\n",
			code:      http.StatusForbidden,
		},
		{
			name:      "get wallets",
			method:    http.MethodGet,
			endpoint:  "/api/v1/wallets",
			expectErr: "403 Forbidden\n",
			code:      http.StatusForbidden,
		},
		{
			name:      "get wallets folder name",
			method:    http.MethodGet,
			endpoint:  "/api/v1/wallets/folderName",
			expectErr: "403 Forbidden\n",
			code:      http.StatusForbidden,
		},
		{
			name:      "main index.html 404 not found",
			method:    http.MethodGet,
			endpoint:  "/api/v1/",
			expectErr: "404 Not Found\n",
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
			expectErr: "403 Forbidden\n",
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
			expectErr: "403 Forbidden\n",
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
			expectErr: "403 Forbidden\n",
			code:      http.StatusForbidden,
		},
		{
			name:        "create transaction",
			method:      http.MethodPost,
			endpoint:    "/api/v1/wallet/transaction",
			contentType: "application/json",
			json: func() interface{} {
				return api.CreateTransactionRequest{
					HoursSelection: api.HoursSelection{
						Type: wallet.HoursSelectionTypeManual,
					},
					Wallet: api.CreateTransactionRequestWallet{
						ID: "test.wlt",
					},
					ChangeAddress: &changeAddress,
					To: []api.Receiver{
						{
							Address: changeAddress,
							Coins:   "0.001",
							Hours:   "1",
						},
					},
				}
			},
			expectErr: "403 Forbidden\n",
			code:      http.StatusForbidden,
		},
	}

	c := api.NewClient(nodeAddress())
	for _, tc := range tt {
		f := func(tc testCase) func(t *testing.T) {
			return func(t *testing.T) {
				var err error
				switch tc.method {
				case http.MethodGet:
					err = c.Get(tc.endpoint, nil)
				case http.MethodPost:
					switch tc.contentType {
					case "application/json":
						err = c.PostJSON(tc.endpoint, tc.json(), nil)
					default:
						err = c.PostForm(tc.endpoint, tc.body(), nil)
					}
				}
				assertResponseError(t, err, tc.code, tc.expectErr)
			}
		}

		t.Run(tc.name, f(tc))

		if strings.HasPrefix(tc.endpoint, "/api/v1") {
			tc.endpoint = strings.TrimPrefix(tc.endpoint, "/api/v1")
			t.Run(tc.name, f(tc))
		}
	}

	// Confirms that no new wallet is created
	// WALLET_DIR environment variable is set in ci-script/integration-test-disable-wallet-api.sh
	walletDir := os.Getenv("WALLET_DIR")
	if walletDir == "" {
		t.Fatal("WALLET_DIR is not set")
	}

	// Confirms that the wallet directory does not exist
	_, err := os.Stat(walletDir)
	require.True(t, os.IsNotExist(err))
}

func checkHealthResponse(t *testing.T, r *api.HealthResponse) {
	require.NotEmpty(t, r.BlockchainMetadata.Unspents)
	require.NotEmpty(t, r.BlockchainMetadata.Head.BkSeq)
	require.NotEmpty(t, r.BlockchainMetadata.Head.Time)
	require.NotEmpty(t, r.Version.Version)
	require.True(t, r.Uptime.Duration > time.Duration(0))
}

func TestStableHealth(t *testing.T) {
	if !doStable(t) {
		return
	}

	c := api.NewClient(nodeAddress())

	r, err := c.Health()
	require.NoError(t, err)

	checkHealthResponse(t, r)

	require.Equal(t, 0, r.OpenConnections)

	require.True(t, r.BlockchainMetadata.TimeSinceLastBlock.Duration > time.Duration(0))

	// The stable node is always run with the commit and branch ldflags, so they should appear
	require.NotEmpty(t, r.Version.Commit)
	require.NotEmpty(t, r.Version.Branch)
}

func TestLiveHealth(t *testing.T) {
	if !doLive(t) {
		return
	}

	c := api.NewClient(nodeAddress())

	r, err := c.Health()
	require.NoError(t, err)

	checkHealthResponse(t, r)

	require.NotEqual(t, 0, r.OpenConnections)

	// The TimeSinceLastBlock can be any value, including negative values, due to clock skew
	// The live node is not necessarily run with the commit and branch ldflags, so don't check them
}

func TestDisableGUIAPI(t *testing.T) {
	if !doLiveOrStable(t) {
		return
	}

	c := api.NewClient(nodeAddress())
	err := c.Get("/", nil)
	assertResponseError(t, err, http.StatusNotFound, "404 Not Found\n")
}
