package gui_integration_test

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/gui"
	"github.com/skycoin/skycoin/src/util/droplet"
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
TODO: When go1.10 is released, use the new DisallowUnknownFields property of the JSON decoder, to detect when
an API adds a new field to the response. See: https://tip.golang.org/doc/go1.10#encoding/json

When update flag is set to true all tests pass
*/

const (
	testModeStable = "stable"
	testModeLive   = "live"
)

type TestData struct {
	actual   interface{}
	expected interface{}
}

var update = flag.Bool("update", false, "update golden files")

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
	case testModeLive, testModeStable:
	default:
		t.Fatal("Invalid test mode, must be stable or live")
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

func loadJSON(t *testing.T, filename string, testData *TestData) {
	goldenFile := filepath.Join("test-fixtures", filename)
	f, err := os.Open(goldenFile)

	if *update {
		updateGoldenFile(t, goldenFile, testData.actual)
	}

	require.NoError(t, err)
	defer f.Close()

	err = json.NewDecoder(f).Decode(testData.expected)
	require.NoError(t, err)
}

func updateGoldenFile(t *testing.T, filename string, content interface{}) {
	contentJson, err := json.MarshalIndent(content, "", "\t")
	require.NoError(t, err)
	err = ioutil.WriteFile(filename, contentJson, 0644)
	require.NoError(t, err)
}

func assertResponseError(t *testing.T, err error, errCode int, errMsg string) {
	require.Error(t, err)
	require.IsType(t, gui.APIError{}, err)
	require.Equal(t, errCode, err.(gui.APIError).StatusCode)
	require.Equal(t, errMsg, err.(gui.APIError).Message)
}

func TestStableCoinSupply(t *testing.T) {
	if !doStable(t) {
		return
	}

	c := gui.NewClient(nodeAddress())

	cs, err := c.CoinSupply()
	require.NoError(t, err)

	var expected gui.CoinSupply
	loadJSON(t, "coinsupply.golden", &TestData{cs, &expected})

	require.Equal(t, expected, *cs)
}

func TestLiveCoinSupply(t *testing.T) {
	if !doLive(t) {
		return
	}

	c := gui.NewClient(nodeAddress())

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
	if !doStable(t) && !doLive(t) {
		return
	}

	c := gui.NewClient(nodeAddress())

	v, err := c.Version()
	require.NoError(t, err)

	require.NotEmpty(t, v.Version)
}

func TestStableOutputs(t *testing.T) {
	if !doStable(t) {
		return
	}

	c := gui.NewClient(nodeAddress())

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
		{
			name: "addrs and hashes",
			addrs: []string{
				"ALJVNKYL7WGxFBSriiZuwZKWD4b7fbV1od",
				"2THDupTBEo7UqB6dsVizkYUvkKq82Qn4gjf",
				"qxmeHkwgAMfwXyaQrwv9jq3qt228xMuoT5",
			},
			hashes: []string{
				"9e53268a18f8d32a44b4fb183033b49bebfe9d0da3bf3ef2ad1d560500aa54c6",
				"d91e07318227651129b715d2db448ae245b442acd08c8b4525a934f0e87efce9",
				"01f9c1d6c83dbc1c993357436cdf7f214acd0bfa107ff7f1466d1b18ec03563e",
				"fe6762d753d626115c8dd3a053b5fb75d6d419a8d0fb1478c5fffc1fe41c5f20",
			},
			errCode: http.StatusBadRequest,
			errMsg:  "400 Bad Request - addrs and hashes cannot be specified together\n",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			outputs, err := c.Outputs(tc.addrs, tc.hashes)

			if tc.errCode != 0 && tc.errCode != http.StatusOK {
				assertResponseError(t, err, tc.errCode, tc.errMsg)
				return
			}

			require.NoError(t, err)

			var expected visor.ReadableOutputSet
			loadJSON(t, tc.golden, &TestData{outputs, &expected})

			require.Equal(t, len(expected.HeadOutputs), len(outputs.HeadOutputs))
			require.Equal(t, len(expected.OutgoingOutputs), len(outputs.OutgoingOutputs))
			require.Equal(t, len(expected.IncomingOutputs), len(outputs.IncomingOutputs))

			for i, o := range expected.HeadOutputs {
				require.Equal(t, o, outputs.HeadOutputs[i], "mismatch at index %d", i)
			}

			require.Equal(t, expected, *outputs)
		})
	}
}

func TestLiveOutputs(t *testing.T) {
	if !doLive(t) {
		return
	}

	c := gui.NewClient(nodeAddress())

	// Request all outputs and check that HeadOutputs is not empty
	// OutgoingOutputs and IncomingOutputs are variable and could be empty
	outputs, err := c.Outputs(nil, nil)
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

	c := gui.NewClient(nodeAddress())
	for _, seq := range blockSeqs {
		b, err := c.BlockBySeq(seq)
		require.NoError(t, err)
		require.Equal(t, seq, b.Head.BkSeq)
	}
}

func testKnownBlocks(t *testing.T) {
	c := gui.NewClient(nodeAddress())

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
			loadJSON(t, tc.golden, &TestData{b, &expected})

			require.Equal(t, expected, *b)
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

	c := gui.NewClient(nodeAddress())

	metadata, err := c.BlockchainMetadata()
	require.NoError(t, err)

	var expected visor.BlockchainMetadata
	loadJSON(t, "blockchain-metadata.golden", &TestData{metadata, &expected})

	require.Equal(t, expected, *metadata)
}

func TestLiveBlockchainMetadata(t *testing.T) {
	if !doLive(t) {
		return
	}

	c := gui.NewClient(nodeAddress())

	metadata, err := c.BlockchainMetadata()
	require.NoError(t, err)

	require.NotEqual(t, uint64(0), metadata.Head.BkSeq)
}

func TestStableBlockchainProgress(t *testing.T) {
	if !doStable(t) {
		return
	}

	c := gui.NewClient(nodeAddress())

	progress, err := c.BlockchainProgress()
	require.NoError(t, err)

	var expected daemon.BlockchainProgress
	loadJSON(t, "blockchain-progress.golden", &TestData{progress, &expected})

	require.Equal(t, expected, *progress)
}

func TestLiveBlockchainProgress(t *testing.T) {
	if !doLive(t) {
		return
	}

	c := gui.NewClient(nodeAddress())

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

	c := gui.NewClient(nodeAddress())

	cases := []struct {
		name   string
		golden string
		addrs  []string
	}{
		{
			name:   "no addresses",
			golden: "balance-noaddrs.golden",
		},
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
			loadJSON(t, tc.golden, &TestData{balance, &expected})

			require.Equal(t, expected, *balance)
		})
	}
}

func TestLiveBalance(t *testing.T) {
	if !doLive(t) {
		return
	}

	c := gui.NewClient(nodeAddress())

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

	c := gui.NewClient(nodeAddress())

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
			loadJSON(t, tc.golden, &TestData{ux, &expected})

			require.Equal(t, expected, *ux)
		})
	}

	// Scan all uxouts from the result of /outputs
	scanUxOuts(t)
}

func TestLiveUxOut(t *testing.T) {
	if !doLive(t) {
		return
	}

	c := gui.NewClient(nodeAddress())

	// A spent uxout should never change
	ux, err := c.UxOut("fe6762d753d626115c8dd3a053b5fb75d6d419a8d0fb1478c5fffc1fe41c5f20")
	require.NoError(t, err)

	var expected historydb.UxOutJSON
	loadJSON(t, "uxout-spent.golden", &TestData{ux, &expected})
	require.Equal(t, expected, *ux)
	require.NotEqual(t, uint64(0), ux.SpentBlockSeq)

	// Scan all uxouts from the result of /outputs
	scanUxOuts(t)
}

func scanUxOuts(t *testing.T) {
	c := gui.NewClient(nodeAddress())

	outputs, err := c.Outputs(nil, nil)
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
