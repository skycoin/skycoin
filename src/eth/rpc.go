package eth

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"
	"sync"
)

// DefaultRPCURLs are public Ethereum mainnet endpoints that need no key.
var DefaultRPCURLs = []string{
	"https://ethereum-rpc.publicnode.com",
	"https://cloudflare-eth.com",
	"https://eth.drpc.org",
	"https://1rpc.io/eth",
	"https://rpc.mevblocker.io",
}

// RPC is a JSON-RPC client that tries its URLs in order, starting from the
// last one that answered.
// https://ethereum.org/en/developers/docs/apis/json-rpc/
type RPC struct {
	urls   []string
	client *http.Client

	mu   sync.Mutex
	next int
	gen  int // counts SetURLs, so a call that began on older URLs does not move next
}

// RPCError is an error the node returned, as opposed to a failure to reach it.
type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *RPCError) Error() string { return fmt.Sprintf("rpc error %d: %s", e.Code, e.Message) }

// NewRPC makes a client for urls. A nil client is http.DefaultClient.
func NewRPC(urls []string, client *http.Client) (*RPC, error) {
	if len(urls) == 0 {
		return nil, errors.New("no Ethereum RPC URLs")
	}
	if client == nil {
		client = http.DefaultClient
	}
	return &RPC{urls: urls, client: client}, nil
}

// URLs are the endpoints in the order the client tries them.
func (r *RPC) URLs() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]string{}, r.urls...)
}

// SetURLs replaces the endpoints, starting again from the first.
func (r *RPC) SetURLs(urls []string) error {
	if len(urls) == 0 {
		return errors.New("no Ethereum RPC URLs")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.urls = append([]string{}, urls...)
	r.next = 0
	r.gen++
	return nil
}

// Call runs method, moving to the next URL when one cannot be reached or
// answers badly. An error from the node itself is returned at once.
func (r *RPC) Call(ctx context.Context, result interface{}, method string, params ...interface{}) error {
	if params == nil {
		params = []interface{}{}
	}
	body, err := json.Marshal(map[string]interface{}{"jsonrpc": "2.0", "id": 1, "method": method, "params": params})
	if err != nil {
		return err
	}

	r.mu.Lock()
	urls, start, gen := r.urls, r.next, r.gen
	r.mu.Unlock()

	var errs []string
	for i := 0; i < len(urls); i++ {
		n := (start + i) % len(urls)
		raw, err := r.post(ctx, urls[n], body)
		if err == nil {
			r.mu.Lock()
			if r.gen == gen {
				r.next = n
			}
			r.mu.Unlock()
			if result == nil {
				return nil
			}
			return json.Unmarshal(raw, result)
		}
		var rpcErr *RPCError
		if errors.As(err, &rpcErr) {
			return err
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
		errs = append(errs, fmt.Sprintf("%s: %v", urls[n], err))
	}
	return fmt.Errorf("no Ethereum RPC answered %s: %s", method, strings.Join(errs, "; "))
}

func (r *RPC) post(ctx context.Context, url string, body []byte) (json.RawMessage, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body)) //nolint:gosec // the RPC URLs are the operator's configuration
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := r.client.Do(req) //nolint:gosec // the RPC URLs are the operator's configuration
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() //nolint:errcheck
	b, err := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}
	var out struct {
		Result json.RawMessage `json:"result"`
		Error  *RPCError       `json:"error"`
	}
	if err := json.Unmarshal(b, &out); err != nil {
		return nil, fmt.Errorf("bad response: %w", err)
	}
	if out.Error != nil {
		return nil, out.Error
	}
	return out.Result, nil
}

func (r *RPC) quantity(ctx context.Context, method string, params ...interface{}) (*big.Int, error) {
	var s string
	if err := r.Call(ctx, &s, method, params...); err != nil {
		return nil, err
	}
	return parseQuantity(s)
}

func parseQuantity(s string) (*big.Int, error) {
	v, ok := new(big.Int).SetString(strings.TrimPrefix(s, "0x"), 16)
	if !ok || !strings.HasPrefix(s, "0x") {
		return nil, fmt.Errorf("bad quantity %q", s)
	}
	return v, nil
}

func quantityHex(v *big.Int) string { return "0x" + v.Text(16) }

// ChainID is the chain the node follows, 1 for mainnet.
func (r *RPC) ChainID(ctx context.Context) (*big.Int, error) {
	return r.quantity(ctx, "eth_chainId")
}

// Balance is an address's balance in wei.
func (r *RPC) Balance(ctx context.Context, a Address) (*big.Int, error) {
	return r.quantity(ctx, "eth_getBalance", a.String(), "latest")
}

// Nonce is the next nonce for a, counting pending transactions.
func (r *RPC) Nonce(ctx context.Context, a Address) (uint64, error) {
	n, err := r.quantity(ctx, "eth_getTransactionCount", a.String(), "pending")
	if err != nil {
		return 0, err
	}
	return n.Uint64(), nil
}

// BaseFee is the base fee per gas of the latest block.
func (r *RPC) BaseFee(ctx context.Context) (*big.Int, error) {
	var block struct {
		BaseFeePerGas string `json:"baseFeePerGas"`
	}
	if err := r.Call(ctx, &block, "eth_getBlockByNumber", "latest", false); err != nil {
		return nil, err
	}
	return parseQuantity(block.BaseFeePerGas)
}

// MaxPriorityFee is the node's suggested tip per gas.
func (r *RPC) MaxPriorityFee(ctx context.Context) (*big.Int, error) {
	return r.quantity(ctx, "eth_maxPriorityFeePerGas")
}

// EstimateGas is the gas a call from from to to with value and data needs.
func (r *RPC) EstimateGas(ctx context.Context, from, to Address, value *big.Int, data []byte) (uint64, error) {
	tx := map[string]string{"from": from.String(), "to": to.String(), "value": quantityHex(value)}
	if len(data) > 0 {
		tx["data"] = "0x" + hex.EncodeToString(data)
	}
	g, err := r.quantity(ctx, "eth_estimateGas", tx)
	if err != nil {
		return 0, err
	}
	return g.Uint64(), nil
}

// TokenBalance is owner's balance of an ERC-20 token, in its base units.
func (r *RPC) TokenBalance(ctx context.Context, token, owner Address) (*big.Int, error) {
	var s string
	call := map[string]string{"to": token.String(), "data": "0x" + hex.EncodeToString(ERC20BalanceOf(owner))}
	if err := r.Call(ctx, &s, "eth_call", call, "latest"); err != nil {
		return nil, err
	}
	b, err := hex.DecodeString(strings.TrimPrefix(s, "0x"))
	if err != nil {
		return nil, fmt.Errorf("bad balanceOf result %q", s)
	}
	return new(big.Int).SetBytes(b), nil
}

// SendRaw broadcasts a signed transaction and returns its hash.
func (r *RPC) SendRaw(ctx context.Context, raw []byte) (string, error) {
	var h string
	if err := r.Call(ctx, &h, "eth_sendRawTransaction", "0x"+hex.EncodeToString(raw)); err != nil {
		return "", err
	}
	return h, nil
}

// BlockNumber is the number of the latest block.
func (r *RPC) BlockNumber(ctx context.Context) (uint64, error) {
	n, err := r.quantity(ctx, "eth_blockNumber")
	if err != nil {
		return 0, err
	}
	return n.Uint64(), nil
}
