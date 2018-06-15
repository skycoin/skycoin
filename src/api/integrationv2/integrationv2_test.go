// package integrationv2_test implements API integration tests
package integrationv2_test

import (
	"net/http"
	"os"
	"testing"

	"github.com/skycoin/skycoin/src/api"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/stretchr/testify/require"
)

const (
	testModeStable           = "stable"
	testModeLive             = "live"
	testModeDisableWalletApi = "disable-wallet-api"
	testModeDisableSeedApi   = "disable-seed-api"
)

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

func nodeAddress() string {
	addr := os.Getenv("SKYCOIN_NODE_HOST")
	if addr == "" {
		return "http://127.0.0.1:6420"
	}
	return addr
}

func assertResponseError(t *testing.T, err error, errCode int, errMsg string) {
	require.Error(t, err)
	require.IsType(t, api.ClientError{}, err)
	require.Equal(t, errCode, err.(api.ClientError).StatusCode)
	require.Equal(t, errMsg, err.(api.ClientError).Message)
}

func testBlock(t *testing.T, block *visor.ReadableBlockV2) {
	for _, trans := range block.Body.Transactions {
		for _, in := range trans.In {
			found := false
			for _, input := range trans.InData {
				if input.Hash == in {
					found = true
				}
			}
			require.Equal(t, found, true)
		}
	}
}

func testKnownBlocks(t *testing.T) {
	if !doStable(t) {
		return
	}

	c := api.NewClientV2(nodeAddress())

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
			var b *visor.ReadableBlockV2
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
			testBlock(t, b)
		})
	}
}

func TestBlocks(t *testing.T) {
	if !doStable(t) {
		return
	}
	c := api.NewClientV2(nodeAddress())
	blocks, err := c.Blocks(0, 1)
	require.NoError(t, err)
	for _, block := range blocks.Blocks {
		testBlock(t, &block)
	}
}

func TestLastBlocks(t *testing.T) {
	if !doStable(t) {
		return
	}
	c := api.NewClientV2(nodeAddress())
	blocks, err := c.LastBlocks(10)
	require.NoError(t, err)
	for _, block := range blocks.Blocks {
		testBlock(t, &block)
	}
}
