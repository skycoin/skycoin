package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"

	"github.com/skycoin/skycoin/src/eth"
)

// ethEndpointsFile keeps the RPC and indexer addresses chosen in the wallet,
// beside the wallets it serves.
const ethEndpointsFile = "ethereum-endpoints.json"

// ethEndpoints are the RPC and indexer addresses every Ethereum coin shares.
// The wallet can change them, and they are kept in the wallet directory.
type ethEndpoints struct {
	rpc            *eth.RPC
	indexer        *eth.Indexer
	defaultRPC     []string
	defaultIndexer string
	file           string

	mu sync.Mutex
}

type ethEndpointsBody struct {
	RPCURLs    []string `json:"rpc_urls"`
	IndexerURL string   `json:"indexer_url"`
}

func newEthEndpoints(rpc *eth.RPC, indexer *eth.Indexer, dir string) *ethEndpoints {
	return &ethEndpoints{
		rpc:            rpc,
		indexer:        indexer,
		defaultRPC:     rpc.URLs(),
		defaultIndexer: indexer.Base(),
		file:           filepath.Join(dir, ethEndpointsFile),
	}
}

// load applies the saved addresses, if there are any.
func (e *ethEndpoints) load() error {
	b, err := os.ReadFile(e.file)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	var saved ethEndpointsBody
	if err := json.Unmarshal(b, &saved); err != nil {
		return fmt.Errorf("%s: %w", e.file, err)
	}
	return e.apply(saved)
}

// apply checks body and uses it, an empty list or address meaning the default.
func (e *ethEndpoints) apply(body ethEndpointsBody) error {
	given := body.RPCURLs
	if len(given) == 0 {
		given = e.defaultRPC
	}
	urls := make([]string, len(given))
	for i, u := range given {
		u, err := checkedHTTPURL(u)
		if err != nil {
			return err
		}
		urls[i] = u
	}
	indexer := e.defaultIndexer
	if body.IndexerURL != "" {
		u, err := checkedHTTPURL(body.IndexerURL)
		if err != nil {
			return err
		}
		indexer = u
	}
	if err := e.rpc.SetURLs(urls); err != nil {
		return err
	}
	e.indexer.SetBase(indexer)
	return nil
}

func checkedHTTPURL(raw string) (string, error) {
	s := strings.TrimSuffix(strings.TrimSpace(raw), "/")
	u, err := url.Parse(s)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return "", fmt.Errorf("not an http or https address: %q", raw)
	}
	return s, nil
}

func (e *ethEndpoints) body() H {
	return H{
		"rpc_urls":            e.rpc.URLs(),
		"indexer_url":         e.indexer.Base(),
		"default_rpc_urls":    e.defaultRPC,
		"default_indexer_url": e.defaultIndexer,
	}
}

func (e *ethEndpoints) get(c *webCtx) {
	c.JSON(http.StatusOK, e.body())
}

// put replaces the addresses and saves them. Sending the defaults, or empty
// values, removes the saved file so later default changes apply again.
func (e *ethEndpoints) put(c *webCtx) {
	var body ethEndpointsBody
	if err := json.NewDecoder(c.Request.Body).Decode(&body); err != nil {
		c.JSON(http.StatusBadRequest, H{"error": "invalid JSON body"})
		return
	}

	e.mu.Lock()
	defer e.mu.Unlock()
	if err := e.apply(body); err != nil {
		c.JSON(http.StatusBadRequest, H{"error": err.Error()})
		return
	}

	saved := ethEndpointsBody{RPCURLs: e.rpc.URLs(), IndexerURL: e.indexer.Base()}
	var err error
	if slices.Equal(saved.RPCURLs, e.defaultRPC) && saved.IndexerURL == e.defaultIndexer {
		err = os.Remove(e.file)
		if errors.Is(err, os.ErrNotExist) {
			err = nil
		}
	} else {
		var b []byte
		b, err = json.MarshalIndent(saved, "", "  ")
		if err == nil {
			err = os.WriteFile(e.file, b, 0600)
		}
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, H{"error": "saving the addresses: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, e.body())
}
