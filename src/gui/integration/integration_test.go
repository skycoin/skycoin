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

// DefaultConnections the default trust node addresses
var defaultConnections = map[string]bool{
	"118.178.135.93:6000":  true,
	"47.88.33.156:6000":    true,
	"121.41.103.148:6000":  true,
	"120.77.69.188:6000":   true,
	"104.237.142.206:6000": true,
	"176.58.126.224:6000":  true,
	"172.104.85.6:6000":    true,
	"139.162.7.132:6000":   true,
}

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

func loadJSON(t *testing.T, filename string, testData *TestData) {
	goldenFile := filepath.Join("test-fixtures", filename)
	f, err := os.Open(goldenFile)

	if *update {
		updateGoldenFile(t, goldenFile, testData.actual)
	}

	require.NoError(t, err)
	defer f.Close()

	err = json.NewDecoder(f).Decode(testData.expected)
	require.NoError(t, err, filename)
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
	if !doLiveOrStable(t) {
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

func TestStableAddressUxOuts(t *testing.T) {
	if !doStable(t) {
		return
	}

	c := gui.NewClient(nodeAddress())

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
			loadJSON(t, tc.golden, &TestData{ux, &expected})
			require.Equal(t, expected, ux, tc.name)
		})
	}
}

func TestLiveAddressUxOuts(t *testing.T) {
	if !doLive(t) {
		return
	}

	c := gui.NewClient(nodeAddress())

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

	c := gui.NewClient(nodeAddress())

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
				loadJSON(t, tc.golden, &TestData{resp, &expected})

				require.Equal(t, expected, *resp)
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
	c := gui.NewClient(nodeAddress())

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

	c := gui.NewClient(nodeAddress())

	blocks, err := c.LastBlocks(1)
	require.NoError(t, err)

	var expected *visor.ReadableBlocks
	loadJSON(t, "block-last.golden", &TestData{blocks, &expected})
	require.Equal(t, expected, blocks)

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
	c := gui.NewClient(nodeAddress())
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

func TestLiveConnections(t *testing.T) {
	if !doLive(t) {
		return
	}

	c := gui.NewClient(nodeAddress())
	connections, err := c.NetworkConnections()
	require.NoError(t, err)
	require.True(t, len(connections.Connections) > 0)

}

func TestLiveConnection(t *testing.T) {
	if !doLive(t) {
		return
	}
	c := gui.NewClient(nodeAddress())
	connections, err := c.NetworkConnections()
	require.NoError(t, err)
	for _, cc := range connections.Connections {
		connection, err := c.NetworkConnection(cc.Addr)
		require.Equal(t, cc.Addr, connection.Addr)
		require.Equal(t, cc.ID, connection.ID)
		require.Equal(t, cc.ListenPort, connection.ListenPort)
		require.Equal(t, cc.Mirror, connection.Mirror)
		require.Equal(t, cc.Introduced, connection.Introduced)
		require.Equal(t, cc.Outgoing, connection.Outgoing)
		require.True(t, cc.LastReceived <= connection.LastReceived)
		require.True(t, cc.LastSent <= connection.LastReceived)
		require.NoError(t, err)
	}
}

func TestLiveNetworkDefaultConnections(t *testing.T) {
	if !doLive(t) {
		return
	}

	c := gui.NewClient(nodeAddress())
	connections, err := c.NetworkDefaultConnections()
	require.NoError(t, err)
	require.True(t, len(connections) == len(defaultConnections))
	for _, connection := range connections {
		require.True(t, defaultConnections[connection])
	}
}

func TestNetworkTrustedConnections(t *testing.T) {
	if !doLive(t) {
		return
	}
	c := gui.NewClient(nodeAddress())
	connections, err := c.NetworkTrustedConnections()
	require.NoError(t, err)
	require.True(t, len(connections) > 0)
	for _, connection := range connections {
		require.True(t, defaultConnections[connection])
	}
}

func TestNetworkExchangeableConnections(t *testing.T) {
	if !doLive(t) {
		return
	}
	c := gui.NewClient(nodeAddress())
	connections, err := c.NetworkExchangeableConnections()
	require.NoError(t, err)
	require.True(t, len(connections) > 0)
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

func TestLiveTransaction(t *testing.T) {
	if !doLive(t) {
		return
	}

	tt := []struct {
		name       string
		txId       string
		err        gui.APIError
		goldenFile string
	}{
		{
			name: "invalid txId",
			txId: "abcd",
			err: gui.APIError{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Message:    "400 Bad Request - Invalid hex length\n",
			},
		},
		{
			name: "empty txId",
			txId: "",
			err: gui.APIError{
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

	c := gui.NewClient(nodeAddress())
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tx, err := c.Transaction(tc.txId)
			if err != nil {
				require.Equal(t, tc.err, err)
				return
			}
			var expected *visor.ReadableTransaction
			loadJSON(t, tc.goldenFile, &TestData{tx, &expected})
			require.Equal(t, expected, tx)
		})
	}
}

func TestStableTransaction(t *testing.T) {
	if !doStable(t) {
		return
	}

	tt := []struct {
		name       string
		txId       string
		err        gui.APIError
		goldenFile string
	}{
		{
			name: "invalid txId",
			txId: "abcd",
			err: gui.APIError{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Message:    "400 Bad Request - Invalid hex length\n",
			},
			goldenFile: "",
		},
		{
			name: "not exist",
			txId: "701d23fd513bad325938ba56869f9faba19384a8ec3dd41833aff147eac53947",
			err: gui.APIError{
				Status:     "404 Not Found",
				StatusCode: http.StatusNotFound,
				Message:    "404 Not Found\n",
			},
			goldenFile: "",
		},
		{
			name: "empty txId",
			txId: "",
			err: gui.APIError{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Message:    "400 Bad Request - txid is empty\n",
			},
			goldenFile: "",
		},
		{
			name:       "genesis transaction",
			txId:       "d556c1c7abf1e86138316b8c17183665512dc67633c04cf236a8b7f332cb4add",
			goldenFile: "./genesis-transaction.golden",
		},
	}

	c := gui.NewClient(nodeAddress())
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tx, err := c.Transaction(tc.txId)
			if err != nil {
				require.Equal(t, tc.err, err)
				return
			}
			var expected *visor.ReadableTransaction
			loadJSON(t, tc.goldenFile, &TestData{tx, &expected})
			require.Equal(t, expected, tx)
		})
	}
}

func TestLiveTransactions(t *testing.T) {
	if !doLive(t) {
		return
	}

	c := gui.NewClient(nodeAddress())
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

	tt := []struct {
		name       string
		addrs      []string
		err        gui.APIError
		goldenFile string
	}{
		{
			name:  "invalid addr length",
			addrs: []string{"abcd"},
			err: gui.APIError{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Message:    "400 Bad Request - parse parameter: 'addrs' failed: Invalid address length\n",
			},
		},
		{
			name:  "invalid addr character",
			addrs: []string{"701d23fd513bad325938ba56869f9faba19384a8ec3dd41833aff147eac53947"},
			err: gui.APIError{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Message:    "400 Bad Request - parse parameter: 'addrs' failed: Invalid base58 character\n",
			},
		},
		{
			name:  "invalid checksum",
			addrs: []string{"2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKk"},
			err: gui.APIError{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Message:    "400 Bad Request - parse parameter: 'addrs' failed: Invalid checksum\n",
			},
		},
		{
			name:  "empty addrs",
			addrs: []string{},
			err: gui.APIError{
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

	c := gui.NewClient(nodeAddress())
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			txResult, err := c.Transactions(tc.addrs)
			if err != nil {
				require.Equal(t, tc.err, err, "case: "+tc.name)
				return
			}

			var expected *[]visor.TransactionResult
			loadJSON(t, tc.goldenFile, &TestData{txResult, &expected})
			require.Equal(t, expected, txResult, "case: "+tc.name)
		})
	}
}

func TestLiveConfirmedTransactions(t *testing.T) {
	if !doLive(t) {
		return
	}
	c := gui.NewClient(nodeAddress())

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
	tt := []struct {
		name       string
		addrs      []string
		err        gui.APIError
		goldenFile string
	}{
		{
			name:  "invalid addr length",
			addrs: []string{"abcd"},
			err: gui.APIError{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Message:    "400 Bad Request - parse parameter: 'addrs' failed: Invalid address length\n",
			},
		},
		{
			name:  "invalid addr character",
			addrs: []string{"701d23fd513bad325938ba56869f9faba19384a8ec3dd41833aff147eac53947"},
			err: gui.APIError{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Message:    "400 Bad Request - parse parameter: 'addrs' failed: Invalid base58 character\n",
			},
		},
		{
			name:  "invalid checksum",
			addrs: []string{"2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKk"},
			err: gui.APIError{
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

	c := gui.NewClient(nodeAddress())
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			txResult, err := c.ConfirmedTransactions(tc.addrs)
			if err != nil {
				require.Equal(t, tc.err, err, "case: "+tc.name)
				return
			}

			var expected *[]visor.TransactionResult
			loadJSON(t, tc.goldenFile, &TestData{txResult, &expected})
			require.Equal(t, expected, txResult, "case: "+tc.name)
		})
	}
}

func TestStableUnconfirmedTransactions(t *testing.T) {
	if !doStable(t) {
		return
	}
	tt := []struct {
		name       string
		addrs      []string
		err        gui.APIError
		goldenFile string
	}{
		{
			name:  "invalid addr length",
			addrs: []string{"abcd"},
			err: gui.APIError{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Message:    "400 Bad Request - parse parameter: 'addrs' failed: Invalid address length\n",
			},
		},
		{
			name:  "invalid addr character",
			addrs: []string{"701d23fd513bad325938ba56869f9faba19384a8ec3dd41833aff147eac53947"},
			err: gui.APIError{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Message:    "400 Bad Request - parse parameter: 'addrs' failed: Invalid base58 character\n",
			},
		},
		{
			name:  "invalid checksum",
			addrs: []string{"2kvLEyXwAYvHfJuFCkjnYNRTUfHPyWgVwKk"},
			err: gui.APIError{
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

	c := gui.NewClient(nodeAddress())
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			txResult, err := c.UnconfirmedTransactions(tc.addrs)
			if err != nil {
				require.Equal(t, tc.err, err, "case: "+tc.name)
				return
			}

			var expected *[]visor.TransactionResult
			loadJSON(t, tc.goldenFile, &TestData{txResult, &expected})
			require.Equal(t, expected, txResult, "case: "+tc.name)
		})
	}
}

func TestLiveUnconfirmedTransactions(t *testing.T) {
	if !doLive(t) {
		return
	}
	c := gui.NewClient(nodeAddress())

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
	c := gui.NewClient(nodeAddress())
	res, err := c.ResendUnconfirmedTransactions()
	require.NoError(t, err)
	require.True(t, len(res.Txids) == 0)
}

func TestLiveResendUnconfirmedTransactions(t *testing.T) {
	if !doLive(t) {
		return
	}
	c := gui.NewClient(nodeAddress())
	_, err := c.ResendUnconfirmedTransactions()
	require.NoError(t, err)
}

func TestStableRawTransaction(t *testing.T) {
	if !doStable(t) {
		return
	}

	tt := []struct {
		name  string
		txId  string
		err   gui.APIError
		rawTx string
	}{
		{
			name: "invalid hex length",
			txId: "abcd",
			err: gui.APIError{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Message:    "400 Bad Request - Invalid hex length\n",
			},
		},
		{
			name: "not found",
			txId: "701d23fd513bad325938ba56869f9faba19384a8ec3dd41833aff147eac53947",
			err: gui.APIError{
				Status:     "404 Not Found",
				StatusCode: http.StatusNotFound,
				Message:    "404 Not Found\n",
			},
		},
		{
			name: "odd length hex string",
			txId: "abcdeffedca",
			err: gui.APIError{
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

	c := gui.NewClient(nodeAddress())
	for _, tc := range tt {
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

	tt := []struct {
		name  string
		txId  string
		err   gui.APIError
		rawTx string
	}{
		{
			name: "invalid hex length",
			txId: "abcd",
			err: gui.APIError{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Message:    "400 Bad Request - Invalid hex length\n",
			},
		},
		{
			name: "odd length hex string",
			txId: "abcdeffedca",
			err: gui.APIError{
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

	c := gui.NewClient(nodeAddress())
	for _, tc := range tt {
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
