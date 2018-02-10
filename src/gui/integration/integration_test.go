package gui_integration_test

import (
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/skycoin/skycoin/src/gui"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/stretchr/testify/require"
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

*/

const (
	testModeStable = "stable"
	testModeLive   = "live"
)

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
	return enabled() && mode(t) == testModeStable
}

func doLive(t *testing.T) bool {
	return enabled() && mode(t) == testModeLive
}

func loadJSON(t *testing.T, filename string, obj interface{}) {
	f, err := os.Open(filename)
	require.NoError(t, err)
	defer f.Close()

	err = json.NewDecoder(f).Decode(obj)
	require.NoError(t, err)
}

func TestStableCoinSupply(t *testing.T) {
	if !doStable(t) {
		t.Skip()
		return
	}

	c := gui.NewClient(nodeAddress())

	cs, err := c.CoinSupply()
	require.NoError(t, err)

	var expected gui.CoinSupply
	loadJSON(t, "coinsupply.golden", &expected)

	require.Equal(t, expected, *cs)
}

func TestLiveCoinSupply(t *testing.T) {
	if !doLive(t) {
		t.Skip()
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
		t.Skip()
		return
	}

	c := gui.NewClient(nodeAddress())

	v, err := c.Version()
	require.NoError(t, err)

	require.NotEmpty(t, v.Version)
}

func TestStableOutputs(t *testing.T) {
	if !doStable(t) {
		t.Skip()
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
				require.Error(t, err)
				require.IsType(t, gui.APIError{}, err)
				require.Equal(t, tc.errCode, err.(gui.APIError).StatusCode)
				require.Equal(t, tc.errMsg, err.(gui.APIError).Message)
				return
			}

			require.NoError(t, err)

			var expected visor.ReadableOutputSet
			loadJSON(t, tc.golden, &expected)

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
