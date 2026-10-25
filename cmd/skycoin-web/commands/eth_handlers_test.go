package commands

import (
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher/bip44"
	"github.com/skycoin/skycoin/src/cipher/crypto"
	"github.com/skycoin/skycoin/src/eth"
	"github.com/skycoin/skycoin/src/wallet"
)

const hardhatMnemonic = "test test test test test test test test test test test junk"

// fakeNode answers the JSON-RPC calls the handler makes, holding balance wei
// at every address, and records raw transactions sent to it.
func fakeNode(t *testing.T, balance *big.Int, sent *[]string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Method string            `json:"method"`
			Params []json.RawMessage `json:"params"`
		}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&req))
		var result interface{}
		switch req.Method {
		case "eth_chainId":
			result = "0x1"
		case "eth_blockNumber":
			result = "0x10"
		case "eth_getBalance":
			result = "0x" + balance.Text(16)
		case "eth_getTransactionCount":
			result = "0x5"
		case "eth_getBlockByNumber":
			result = map[string]string{"baseFeePerGas": "0x3b9aca00"}
		case "eth_maxPriorityFeePerGas":
			result = "0x3b9aca00"
		case "eth_sendRawTransaction":
			var raw string
			require.NoError(t, json.Unmarshal(req.Params[0], &raw))
			*sent = append(*sent, raw)
			result = "0xfeed"
		default:
			t.Fatalf("unexpected method %s", req.Method)
		}
		require.NoError(t, json.NewEncoder(w).Encode(map[string]interface{}{"jsonrpc": "2.0", "id": 1, "result": result}))
	}))
}

func ethTestHandler(t *testing.T, balance *big.Int, sent *[]string) (*ethHandler, string) {
	node := fakeNode(t, balance, sent)
	t.Cleanup(node.Close)
	rpc, err := eth.NewRPC([]string{node.URL}, nil)
	require.NoError(t, err)

	dir, err := os.MkdirTemp("", "ethwallets")
	require.NoError(t, err)
	t.Cleanup(func() { _ = os.RemoveAll(dir) }) //nolint:errcheck
	coin := bip44.CoinTypeEthereum
	svc, err := wallet.NewService(wallet.Config{
		WalletDir: dir, CryptoType: crypto.CryptoTypeScryptChacha20poly1305Insecure, EnableWalletAPI: true, Bip44Coin: &coin,
	})
	require.NoError(t, err)
	w, err := svc.CreateWallet("", wallet.Options{
		Label: "eth", Seed: hardhatMnemonic, Type: wallet.WalletTypeBip44,
		Encrypt: true, Password: []byte("pw"), CryptoType: crypto.CryptoTypeScryptChacha20poly1305Insecure,
	})
	require.NoError(t, err)
	return &ethHandler{rpc: rpc, wltService: svc}, w.Filename()
}

func serveEth(h *ethHandler, method, target, body string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(method, target, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if !h.handleEthAPI(newCtx(rec, req), req.URL.Path) {
		rec.WriteHeader(http.StatusNotFound)
	}
	return rec
}

func TestEthBalance(t *testing.T) {
	h, _ := ethTestHandler(t, big.NewInt(7), nil)
	rec := serveEth(h, http.MethodGet, "/v1/eth/balance?addrs=0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266,0x70997970C51812dc3A010C7d01b50e0d17dc79C8", "")
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
	var got struct {
		Confirmed struct{ Coins string }
		Addresses map[string]struct{ Confirmed struct{ Coins string } }
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	require.Equal(t, "14", got.Confirmed.Coins)
	require.Equal(t, "7", got.Addresses["0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"].Confirmed.Coins)

	rec = serveEth(h, http.MethodGet, "/v1/eth/balance?addrs=nope", "")
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestEthHistory(t *testing.T) {
	ix := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		// The same transfer between the wallet's two addresses, and a failed one.
		_, _ = w.Write([]byte(`{"message":"OK","result":[` + //nolint:errcheck
			`{"hash":"0x1","blockNumber":"5","timeStamp":"200","from":"0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266","to":"0x70997970c51812dc3a010c7d01b50e0d17dc79c8","value":"9","gasUsed":"21000","gasPrice":"2","isError":"0"},` +
			`{"hash":"0x2","blockNumber":"4","timeStamp":"100","from":"0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266","to":"0x70997970c51812dc3a010c7d01b50e0d17dc79c8","value":"5","gasUsed":"21000","gasPrice":"2","isError":"1"}]}`))
	}))
	defer ix.Close()
	h, _ := ethTestHandler(t, big.NewInt(0), nil)
	h.indexer = eth.NewIndexer(ix.URL, nil)

	rec := serveEth(h, http.MethodGet, "/v1/eth/history?addrs=0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266,0x70997970C51812dc3A010C7d01b50e0d17dc79C8", "")
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
	var txs []struct {
		TxID    string `json:"txid"`
		Fee     string `json:"fee"`
		Failed  bool   `json:"failed"`
		Inputs  []struct{ Address string }
		Outputs []struct{ Address, Value string }
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &txs))
	require.Len(t, txs, 2)
	require.Equal(t, "0x1", txs[0].TxID)
	require.Equal(t, "42000", txs[0].Fee)
	require.Equal(t, "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266", txs[0].Inputs[0].Address)
	require.Equal(t, "9", txs[0].Outputs[0].Value)
	require.True(t, txs[1].Failed)
	require.Equal(t, "0", txs[1].Outputs[0].Value)
}

func TestEthSend(t *testing.T) {
	var sent []string
	h, id := ethTestHandler(t, new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil), &sent)
	to := "0x70997970C51812dc3A010C7d01b50e0d17dc79C8"

	rec := serveEth(h, http.MethodPost, "/v1/eth/send", `{"wallet_id":"`+id+`","password":"bad","to":"`+to+`","amount":"1000"}`)
	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Empty(t, sent)

	rec = serveEth(h, http.MethodPost, "/v1/eth/send", `{"wallet_id":"`+id+`","password":"pw","to":"`+to+`","amount":"1000"}`)
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
	require.Contains(t, rec.Body.String(), `"from":"0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"`)
	require.Len(t, sent, 1)
	require.True(t, strings.HasPrefix(sent[0], "0x02"))

	// More than any address holds.
	rec = serveEth(h, http.MethodPost, "/v1/eth/send", `{"wallet_id":"`+id+`","password":"pw","to":"`+to+`","amount":"1000000000000000000"}`)
	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Contains(t, rec.Body.String(), "no address in this wallet holds enough")
	require.Len(t, sent, 1)

	rec = serveEth(h, http.MethodPost, "/v1/eth/send", `{"wallet_id":"`+id+`","password":"pw","to":"0xnope","amount":"1"}`)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestParseERC20Tokens(t *testing.T) {
	tokens, err := parseERC20Tokens(defaultERC20Tokens)
	require.NoError(t, err)
	require.Len(t, tokens, 1)
	require.Equal(t, "USDT", tokens[0].Symbol)
	require.Equal(t, 6, tokens[0].Decimals)
	require.Equal(t, "0xdAC17F958D2ee523a2206206994597C13D831ec7", tokens[0].Contract.String())

	none, err := parseERC20Tokens("none")
	require.NoError(t, err)
	require.Empty(t, none)

	for _, bad := range []string{"USDT", "USDT:0x12:6", "USDT:0xdAC17F958D2ee523a2206206994597C13D831ec7:x"} {
		_, err := parseERC20Tokens(bad)
		require.Error(t, err, bad)
	}

	require.Equal(t, eth.DefaultRPCURLs, ethRPCURLs("default"))
	require.Nil(t, ethRPCURLs("none"))
	require.Equal(t, []string{"http://a", "http://b"}, ethRPCURLs("http://a, http://b"))
}
