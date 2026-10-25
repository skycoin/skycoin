package eth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
)

// DefaultIndexerURL is Blockscout's Ethereum mainnet instance, whose
// Etherscan-style API needs no key.
// https://docs.blockscout.com/devs/apis/rpc/account
const DefaultIndexerURL = "https://eth.blockscout.com"

// Indexer reads an address's history from an Etherscan-style API.
type Indexer struct {
	client *http.Client

	mu   sync.Mutex
	base string
}

// IndexedTx is one transaction or token transfer the indexer lists.
type IndexedTx struct {
	Hash      string
	Block     uint64
	Timestamp int64
	From      string
	To        string
	Value     *big.Int
	Fee       *big.Int
	Failed    bool
}

// NewIndexer makes a client for the API at base. A nil client is http.DefaultClient.
func NewIndexer(base string, client *http.Client) *Indexer {
	if client == nil {
		client = http.DefaultClient
	}
	return &Indexer{base: strings.TrimSuffix(base, "/"), client: client}
}

// Base is the API address the indexer is queried at.
func (ix *Indexer) Base() string {
	ix.mu.Lock()
	defer ix.mu.Unlock()
	return ix.base
}

// SetBase moves the indexer to another API address.
func (ix *Indexer) SetBase(base string) {
	ix.mu.Lock()
	defer ix.mu.Unlock()
	ix.base = strings.TrimSuffix(base, "/")
}

// Transactions are the Ether transactions of address, newest first.
func (ix *Indexer) Transactions(ctx context.Context, address string) ([]IndexedTx, error) {
	return ix.list(ctx, url.Values{"action": {"txlist"}, "address": {address}})
}

// TokenTransfers are the transfers of one token to or from address, newest first.
func (ix *Indexer) TokenTransfers(ctx context.Context, address, token string) ([]IndexedTx, error) {
	return ix.list(ctx, url.Values{"action": {"tokentx"}, "address": {address}, "contractaddress": {token}})
}

func (ix *Indexer) list(ctx context.Context, q url.Values) ([]IndexedTx, error) {
	q.Set("module", "account")
	q.Set("sort", "desc")
	q.Set("page", "1")
	q.Set("offset", "100")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ix.Base()+"/api?"+q.Encode(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := ix.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() //nolint:errcheck
	b, err := io.ReadAll(io.LimitReader(resp.Body, 16<<20))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("indexer status %d", resp.StatusCode)
	}
	var out struct {
		Message string          `json:"message"`
		Result  json.RawMessage `json:"result"`
	}
	if err := json.Unmarshal(b, &out); err != nil {
		return nil, fmt.Errorf("bad indexer response: %w", err)
	}
	var rows []map[string]string
	if err := json.Unmarshal(out.Result, &rows); err != nil {
		// No transactions comes back as a message and a non-list result.
		if strings.HasPrefix(out.Message, "No ") {
			return nil, nil
		}
		return nil, fmt.Errorf("indexer: %s", out.Message)
	}
	txs := make([]IndexedTx, 0, len(rows))
	for _, row := range rows {
		txs = append(txs, indexedTx(row))
	}
	return txs, nil
}

func indexedTx(row map[string]string) IndexedTx {
	num := func(k string) *big.Int {
		v, ok := new(big.Int).SetString(row[k], 10)
		if !ok {
			return new(big.Int)
		}
		return v
	}
	block, _ := strconv.ParseUint(row["blockNumber"], 10, 64) //nolint:errcheck // a bad number reads as 0
	ts, _ := strconv.ParseInt(row["timeStamp"], 10, 64)       //nolint:errcheck // a bad number reads as 0
	return IndexedTx{
		Hash:      row["hash"],
		Block:     block,
		Timestamp: ts,
		From:      row["from"],
		To:        row["to"],
		Value:     num("value"),
		Fee:       new(big.Int).Mul(num("gasUsed"), num("gasPrice")),
		Failed:    row["isError"] == "1",
	}
}
