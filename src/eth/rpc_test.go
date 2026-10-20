package eth

import (
	"context"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func rpcServer(t *testing.T, handle func(method string, params []json.RawMessage) (interface{}, *RPCError)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Method string            `json:"method"`
			Params []json.RawMessage `json:"params"`
		}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&req))
		result, rpcErr := handle(req.Method, req.Params)
		resp := map[string]interface{}{"jsonrpc": "2.0", "id": 1}
		if rpcErr != nil {
			resp["error"] = rpcErr
		} else {
			resp["result"] = result
		}
		require.NoError(t, json.NewEncoder(w).Encode(resp))
	}))
}

func TestRPCFailover(t *testing.T) {
	down := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
	}))
	defer down.Close()
	calls := 0
	up := rpcServer(t, func(_ string, _ []json.RawMessage) (interface{}, *RPCError) {
		calls++
		return "0x1", nil
	})
	defer up.Close()

	r, err := NewRPC([]string{down.URL, up.URL}, nil)
	require.NoError(t, err)
	id, err := r.ChainID(context.Background())
	require.NoError(t, err)
	require.Equal(t, int64(1), id.Int64())

	// The next call starts at the server that answered.
	_, err = r.ChainID(context.Background())
	require.NoError(t, err)
	require.Equal(t, 2, calls)
}

func TestRPCErrorIsNotRetried(t *testing.T) {
	calls := 0
	srv := func() *httptest.Server {
		return rpcServer(t, func(string, []json.RawMessage) (interface{}, *RPCError) {
			calls++
			return nil, &RPCError{Code: -32000, Message: "insufficient funds"}
		})
	}
	a, b := srv(), srv()
	defer a.Close()
	defer b.Close()

	r, err := NewRPC([]string{a.URL, b.URL}, nil)
	require.NoError(t, err)
	_, err = r.SendRaw(context.Background(), []byte{1})
	var rpcErr *RPCError
	require.ErrorAs(t, err, &rpcErr)
	require.Equal(t, 1, calls)
}

func TestIndexer(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api", r.URL.Path)
		if r.URL.Query().Get("address") == "0xempty" {
			_, _ = w.Write([]byte(`{"message":"No transactions found","result":[]}`)) //nolint:errcheck
			return
		}
		require.Equal(t, "txlist", r.URL.Query().Get("action"))
		_, _ = w.Write([]byte(`{"message":"OK","result":[{"hash":"0xab","blockNumber":"12","timeStamp":"1700000000",` + //nolint:errcheck
			`"from":"0x1","to":"0x2","value":"1000","gasUsed":"21000","gasPrice":"2","isError":"0"}]}`))
	}))
	defer srv.Close()

	ix := NewIndexer(srv.URL+"/", nil)
	txs, err := ix.Transactions(context.Background(), "0xabc")
	require.NoError(t, err)
	require.Len(t, txs, 1)
	require.Equal(t, "0xab", txs[0].Hash)
	require.Equal(t, uint64(12), txs[0].Block)
	require.Equal(t, int64(1700000000), txs[0].Timestamp)
	require.Equal(t, "1000", txs[0].Value.String())
	require.Equal(t, "42000", txs[0].Fee.String())

	none, err := ix.Transactions(context.Background(), "0xempty")
	require.NoError(t, err)
	require.Empty(t, none)
}

// TestTransferLive sends Ether and a token on a development chain such as
// `npx hardhat node`, whose funded accounts come from testMnemonic.
// ETH_TEST_RPC is its URL and ETH_TEST_TOKEN an ERC-20 the first account holds.
func TestTransferLive(t *testing.T) {
	url := os.Getenv("ETH_TEST_RPC")
	if url == "" {
		t.Skip("ETH_TEST_RPC is not set")
	}
	ctx := context.Background()
	r, err := NewRPC([]string{url}, nil)
	require.NoError(t, err)

	sec, err := Key(testMnemonic, "", 0)
	require.NoError(t, err)
	toSec, err := Key(testMnemonic, "", 1)
	require.NoError(t, err)
	to, err := AddressOf(toSec)
	require.NoError(t, err)

	before, err := r.Balance(ctx, to)
	require.NoError(t, err)
	amount := new(big.Int).Mul(big.NewInt(3), Gwei)
	_, signed, err := Transfer(ctx, r, TransferRequest{From: sec, To: to, Amount: amount})
	require.NoError(t, err)
	hash, err := r.SendRaw(ctx, signed.Raw)
	require.NoError(t, err)
	require.Equal(t, signed.Hash, mustHash(t, hash))
	waitFor(t, func() bool {
		after, err := r.Balance(ctx, to)
		return err == nil && new(big.Int).Sub(after, before).Cmp(amount) == 0
	})

	tokenHex := os.Getenv("ETH_TEST_TOKEN")
	if tokenHex == "" {
		return
	}
	token, err := ParseAddress(tokenHex)
	require.NoError(t, err)
	tokBefore, err := r.TokenBalance(ctx, token, to)
	require.NoError(t, err)
	units := big.NewInt(1234567)
	tx, signed, err := Transfer(ctx, r, TransferRequest{From: sec, To: to, Amount: units, Token: &token})
	require.NoError(t, err)
	require.Equal(t, token, tx.To)
	require.Zero(t, tx.Value.Sign())
	_, err = r.SendRaw(ctx, signed.Raw)
	require.NoError(t, err)
	waitFor(t, func() bool {
		after, err := r.TokenBalance(ctx, token, to)
		return err == nil && new(big.Int).Sub(after, tokBefore).Cmp(units) == 0
	})
}

func mustHash(t *testing.T, s string) [32]byte {
	var h [32]byte
	copy(h[:], mustHex(t, s[2:]))
	return h
}

func waitFor(t *testing.T, cond func() bool) {
	deadline := time.Now().Add(20 * time.Second)
	for !cond() {
		if time.Now().After(deadline) {
			t.Fatal("condition not met in time")
		}
		time.Sleep(200 * time.Millisecond)
	}
}

func TestSuggestFees(t *testing.T) {
	srv := rpcServer(t, func(method string, _ []json.RawMessage) (interface{}, *RPCError) {
		switch method {
		case "eth_getBlockByNumber":
			return map[string]string{"baseFeePerGas": "0x3b9aca00"}, nil
		case "eth_maxPriorityFeePerGas":
			return "0x0", nil
		}
		return nil, &RPCError{Code: -32601, Message: "no method"}
	})
	defer srv.Close()
	r, err := NewRPC([]string{srv.URL}, nil)
	require.NoError(t, err)
	fees, base, err := SuggestFees(context.Background(), r)
	require.NoError(t, err)
	require.Equal(t, "1000000000", base.String())
	require.Equal(t, minTip, fees.MaxPriorityFeePerGas)
	require.Equal(t, "2100000000", fees.MaxFeePerGas.String())
}
