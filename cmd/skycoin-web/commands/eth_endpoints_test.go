package commands

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/eth"
)

func TestEthEndpoints(t *testing.T) {
	dir := t.TempDir()
	rpc, err := eth.NewRPC([]string{"https://a.example"}, nil)
	require.NoError(t, err)
	ix := eth.NewIndexer("https://idx.example", nil)
	h := &ethHandler{rpc: rpc, indexer: ix, endpoints: newEthEndpoints(rpc, ix, dir)}
	file := filepath.Join(dir, ethEndpointsFile)

	rec := serveEth(h, http.MethodPut, "/v1/eth/endpoints", `{"rpc_urls":["ftp://x"]}`)
	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Equal(t, []string{"https://a.example"}, rpc.URLs())

	rec = serveEth(h, http.MethodPut, "/v1/eth/endpoints", `{"rpc_urls":[" https://b.example/ ","http://127.0.0.1:8545"],"indexer_url":"https://other.example/"}`)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, []string{"https://b.example", "http://127.0.0.1:8545"}, rpc.URLs())
	require.Equal(t, "https://other.example", ix.Base())

	var saved ethEndpointsBody
	b, err := os.ReadFile(file) //nolint:gosec // a file in the test directory
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(b, &saved))
	require.Equal(t, ethEndpointsBody{RPCURLs: rpc.URLs(), IndexerURL: "https://other.example"}, saved)

	// A restart picks the saved addresses up again.
	rpc2, err := eth.NewRPC([]string{"https://a.example"}, nil)
	require.NoError(t, err)
	ix2 := eth.NewIndexer("https://idx.example", nil)
	require.NoError(t, newEthEndpoints(rpc2, ix2, dir).load())
	require.Equal(t, rpc.URLs(), rpc2.URLs())
	require.Equal(t, "https://other.example", ix2.Base())

	rec = serveEth(h, http.MethodGet, "/v1/eth/endpoints", "")
	require.JSONEq(t, `{"rpc_urls":["https://b.example","http://127.0.0.1:8545"],"indexer_url":"https://other.example",
		"default_rpc_urls":["https://a.example"],"default_indexer_url":"https://idx.example"}`, rec.Body.String())

	// Empty values go back to the defaults and remove the file.
	rec = serveEth(h, http.MethodPut, "/v1/eth/endpoints", `{}`)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, []string{"https://a.example"}, rpc.URLs())
	require.Equal(t, "https://idx.example", ix.Base())
	_, err = os.Stat(file)
	require.True(t, os.IsNotExist(err))
}
