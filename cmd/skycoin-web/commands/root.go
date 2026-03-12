// Package commands provides commands for the skycoin web interface.
package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/skycoin/skywire/pkg/skywire-utilities/pkg/calvin"
	"github.com/spf13/cobra"

	"github.com/skycoin/skycoin/src/cipher/bip44"
	"github.com/skycoin/skycoin/src/cipher/crypto"
	"github.com/skycoin/skycoin/src/readable"
	wasmtinygo "github.com/skycoin/skycoin/src/skycoin-lite/wasm-tinygo"
	"github.com/skycoin/skycoin/src/skycoin-web/src/gui"
	"github.com/skycoin/skycoin/src/wallet"
)

var (
	port          int
	host          string
	nodeURLs      []string
	walletDirs    []string
	enableSeedAPI bool
)

// proxyCache caches responses for slow read-only endpoints (e.g. transactions)
type proxyCache struct {
	mu      sync.RWMutex
	entries map[string]proxyCacheEntry
}

type proxyCacheEntry struct {
	body        []byte
	contentType string
	statusCode  int
	cachedAt    time.Time
}

var queryCache = &proxyCache{entries: make(map[string]proxyCacheEntry)}

func (pc *proxyCache) get(key string, maxAge time.Duration) (proxyCacheEntry, bool) {
	pc.mu.RLock()
	defer pc.mu.RUnlock()
	entry, ok := pc.entries[key]
	if !ok || time.Since(entry.cachedAt) > maxAge {
		return proxyCacheEntry{}, false
	}
	return entry, true
}

func (pc *proxyCache) set(key string, entry proxyCacheEntry) {
	pc.mu.Lock()
	defer pc.mu.Unlock()
	pc.entries[key] = entry
}

// discoveredCoin represents a fibercoin discovered from a node's health endpoint
type discoveredCoin struct {
	ID                int    `json:"id"`
	NodeURL           string `json:"nodeUrl"`
	CoinName          string `json:"coinName"`
	CoinSymbol        string `json:"coinSymbol"`
	HoursName         string `json:"hoursName"`
	PriceTickerID     string `json:"priceTickerId"`
	PriceTickerSource string `json:"priceTickerSource"`
	CoinExplorer      string `json:"coinExplorer"`
	// internal: the actual remote node URL (not exposed to frontend)
	remoteNodeURL string
}

// RootCmd is the root cil command
var RootCmd = &cobra.Command{
	Use:   "skycoin-web",
	Short: "Skycoin Web Wallet",
	Long: func() (ret string) {
		ret = calvin.AsciiFont("skycoin-web")
		ret += "\nThin client web wallet for Skycoin and fibercoins."
		return ret
	}(),
	Run: func(_ *cobra.Command, _ []string) {
		serve()
	},
}

func init() {
	RootCmd.Flags().IntVarP(&port, "port", "p", 8001, "Port to serve on")
	RootCmd.Flags().StringVarP(&host, "host", "H", "127.0.0.1", "Host to bind to")
	RootCmd.Flags().StringArrayVarP(&nodeURLs, "node-url", "n", []string{"https://node.skycoin.com"}, "Node URL (can be specified multiple times)")
	RootCmd.Flags().StringArrayVarP(&walletDirs, "wallet-dir", "w", nil, "Local wallet directory (can be specified multiple times)")
	RootCmd.Flags().BoolVar(&enableSeedAPI, "enable-seed-api", false, "Enable the wallet seed API (requires --wallet-dir)")
}

// Execute runs the root command
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// discoverCoin queries a node's /api/v1/health endpoint and builds a discoveredCoin
func discoverCoin(index int, nodeURL string) (*discoveredCoin, error) {
	nodeURL = strings.TrimRight(nodeURL, "/")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(nodeURL + "/api/v1/health")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to node %s: %v", nodeURL, err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			log.Printf("Error closing health response body: %v", cerr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("node %s returned status %d", nodeURL, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read health response from %s: %v", nodeURL, err)
	}

	var healthResp struct {
		Fiber readable.FiberConfig `json:"fiber"`
	}
	if err := json.Unmarshal(body, &healthResp); err != nil {
		return nil, fmt.Errorf("failed to parse health response from %s: %v", nodeURL, err)
	}

	f := healthResp.Fiber
	coin := &discoveredCoin{
		ID:                index,
		NodeURL:           fmt.Sprintf("/coin/%d", index),
		CoinName:          f.DisplayName,
		CoinSymbol:        f.Ticker,
		HoursName:         f.CoinHoursName,
		PriceTickerID:     f.PriceTickerID,
		PriceTickerSource: f.PriceTickerSource,
		CoinExplorer:      f.ExplorerURL,
		remoteNodeURL:     nodeURL,
	}

	if coin.CoinName == "" {
		coin.CoinName = f.Name
	}
	if coin.CoinName == "" {
		coin.CoinName = fmt.Sprintf("Coin %d", index)
	}
	if coin.CoinSymbol == "" {
		coin.CoinSymbol = strings.ToUpper(f.Name)
	}

	// Apply default price ticker for Skycoin nodes that don't include it in health response
	if coin.PriceTickerID == "" && strings.EqualFold(f.Name, "skycoin") {
		coin.PriceTickerID = "sky-skycoin"
		coin.PriceTickerSource = "coinpaprika"
	}

	return coin, nil
}

func serve() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Discover coins from configured nodes
	var coins []*discoveredCoin
	for i, rawURL := range nodeURLs {
		nodeURL := strings.TrimRight(rawURL, "/")
		coin, err := discoverCoin(i, nodeURL)
		if err != nil {
			log.Printf("[WARN] Could not discover coin from %s: %v (will use as unconfigured node)", nodeURL, err)
			// Create a placeholder coin for nodes that don't respond to health
			coin = &discoveredCoin{
				ID:            i,
				NodeURL:       fmt.Sprintf("/coin/%d", i),
				CoinName:      fmt.Sprintf("Node %d", i),
				CoinSymbol:    fmt.Sprintf("N%d", i),
				HoursName:     "Coin Hours",
				remoteNodeURL: nodeURL,
			}
		}
		coins = append(coins, coin)
		log.Printf("[COIN] Discovered %s (%s) at %s → proxy /coin/%d", coin.CoinName, coin.CoinSymbol, coin.remoteNodeURL, i)
	}

	// Initialize wallet services for each --wallet-dir
	var wltServices []*wallet.Service
	for _, dir := range walletDirs {
		if dir == "" {
			continue
		}
		bc := bip44.CoinTypeSkycoin
		cfg := wallet.Config{
			WalletDir:       dir,
			CryptoType:      crypto.DefaultCryptoType,
			EnableWalletAPI: true,
			EnableSeedAPI:   enableSeedAPI,
			Bip44Coin:       &bc,
		}

		svc, err := wallet.NewService(cfg)
		if err != nil {
			log.Fatalf("Failed to initialize wallet service for %s: %v", dir, err)
		}
		wltServices = append(wltServices, svc)
		log.Printf("[WALLET] Wallet service initialized: %s", dir)
	}

	// Map wallet services to coins by index.
	// If counts match, wallet dir i is used for coin i.
	// Otherwise, all wallet services are shared across all coins.
	coinWltServices := make(map[int][]*wallet.Service)
	if len(wltServices) == len(coins) {
		for i, svc := range wltServices {
			coinWltServices[i] = []*wallet.Service{svc}
		}
	} else if len(wltServices) > 0 {
		for i := range coins {
			coinWltServices[i] = wltServices
		}
	}

	// Get the embedded dist directory
	distSub, err := fs.Sub(gui.DistFS, "dist")
	if err != nil {
		log.Fatalf("Failed to get dist subdirectory: %v", err)
	}

	// Serve embedded WASM files from skycoin-lite
	router.GET("/assets/scripts/skycoin-lite.wasm", func(c *gin.Context) {
		c.Header("Content-Type", "application/wasm")
		c.Data(http.StatusOK, "application/wasm", wasmtinygo.WasmFile)
	})

	router.GET("/assets/scripts/wasm_exec.js", func(c *gin.Context) {
		c.Header("Content-Type", "application/javascript")
		c.Data(http.StatusOK, "application/javascript", wasmtinygo.WasmExecJS)
	})

	// Per-coin proxy routes: /coin/{index}/api/*
	router.Any("/coin/:coinIndex/api/*path", func(c *gin.Context) {
		setCORSHeaders(c)
		if c.Request.Method == "OPTIONS" {
			c.Status(http.StatusOK)
			return
		}

		coinIndex := 0
		if _, err := fmt.Sscanf(c.Param("coinIndex"), "%d", &coinIndex); err != nil || coinIndex < 0 || coinIndex >= len(coins) {
			c.String(http.StatusBadRequest, "invalid coin index")
			return
		}

		coin := coins[coinIndex]
		apiPath := c.Param("path")

		// Intercept read-only POST requests — convert to GET to avoid CSRF issues
		if c.Request.Method == http.MethodPost {
			trimmed := strings.TrimSuffix(apiPath, "/")
			if handled := handleReadOnlyPost(c, trimmed, coin.remoteNodeURL); handled {
				return
			}
		}

		// Try local wallet handling first
		if services, ok := coinWltServices[coinIndex]; ok && len(services) > 0 {
			if handleMultiWalletAPI(c, apiPath, services, coin.remoteNodeURL) {
				return
			}
		}

		// Proxy to the coin's node
		proxyToNodeWithBase(c, coin.remoteNodeURL, "/api"+apiPath)
	})

	// Legacy /api/* route — proxies to the first node for backwards compatibility
	// Also serves the /api/v1/coins discovery endpoint
	router.Any("/api/*path", func(c *gin.Context) {
		setCORSHeaders(c)
		if c.Request.Method == "OPTIONS" {
			c.Status(http.StatusOK)
			return
		}

		apiPath := c.Param("path")

		// Coins discovery endpoint
		if apiPath == "/v1/coins" && c.Request.Method == http.MethodGet {
			c.JSON(http.StatusOK, coins)
			return
		}

		// Intercept read-only POST requests — convert to GET to avoid CSRF issues
		if c.Request.Method == http.MethodPost {
			trimmed := strings.TrimSuffix(apiPath, "/")
			if len(coins) > 0 {
				if handled := handleReadOnlyPost(c, trimmed, coins[0].remoteNodeURL); handled {
					return
				}
			}
		}

		// Wallet endpoints served locally (legacy route uses coin 0's wallet services)
		if services, ok := coinWltServices[0]; ok && len(services) > 0 {
			defaultNodeURL := ""
			if len(coins) > 0 {
				defaultNodeURL = coins[0].remoteNodeURL
			}
			if handleMultiWalletAPI(c, apiPath, services, defaultNodeURL) {
				return
			}
		}

		// Proxy to first node
		if len(coins) > 0 {
			proxyToNode(c, coins[0].remoteNodeURL)
		} else {
			c.String(http.StatusBadGateway, "no nodes configured")
		}
	})

	// Serve static files from embedded dist directory
	router.NoRoute(func(c *gin.Context) {
		if c.Request.URL.Path != "/" && c.Request.URL.Path != "/favicon.ico" {
			log.Printf("[STATIC] %s %s", c.Request.Method, c.Request.URL.Path)
		}
		fileServer := http.FileServer(http.FS(distSub))
		fileServer.ServeHTTP(c.Writer, c.Request)
	})

	addr := fmt.Sprintf("%s:%d", host, port)
	fmt.Printf("Skycoin Web Wallet starting...\n")
	fmt.Printf("Server listening on http://%s\n", addr)
	if len(coins) > 0 {
		fmt.Printf("Configured coins:\n")
		for _, coin := range coins {
			fmt.Printf("  [%d] %s (%s) → %s\n", coin.ID, coin.CoinName, coin.CoinSymbol, coin.remoteNodeURL)
		}
	}
	for _, svc := range wltServices {
		dir, err := svc.WalletDir()
		if err != nil {
			log.Printf("[WARN] Could not get wallet dir: %v", err)
			continue
		}
		fmt.Printf("Local wallet dir: %s\n", dir)
	}
	fmt.Printf("Open your browser and navigate to the address above\n")
	fmt.Printf("Press Ctrl+C to stop the server\n\n")

	if err := router.Run(addr); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func setCORSHeaders(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Content-Type, X-CSRF-Token")
}

// handleMultiWalletAPI dispatches wallet API requests across all configured wallet services.
// For read operations, it aggregates results from all services.
// For write operations, it uses the first (primary) service.
func handleMultiWalletAPI(c *gin.Context, apiPath string, services []*wallet.Service, nodeURL string) bool {
	path := strings.TrimSuffix(apiPath, "/")
	method := c.Request.Method

	// Aggregated read endpoints
	switch {
	case path == "/v1/wallets" && method == http.MethodGet:
		handleGetWalletsMulti(c, services)
		return true
	case path == "/v1/wallets/folderName" && method == http.MethodGet:
		handleWalletFolderMulti(c, services)
		return true
	}

	// For wallet operations that need to find a wallet by ID, search all services
	if needsWalletLookup(path, method) {
		wltID := c.Request.FormValue("id")
		if wltID != "" {
			for _, svc := range services {
				if _, err := svc.GetWallet(wltID); err == nil {
					return handleWalletAPI(c, apiPath, svc, nodeURL)
				}
			}
		}
	}

	// Write operations and new wallet creation use the primary service
	return handleWalletAPI(c, apiPath, services[0], nodeURL)
}

// needsWalletLookup returns true if the endpoint operates on an existing wallet by ID
func needsWalletLookup(path, method string) bool {
	switch {
	case path == "/v1/wallet" && method == http.MethodGet:
		return true
	case path == "/v1/wallet/newAddress" && method == http.MethodPost:
		return true
	case path == "/v1/wallet/update" && method == http.MethodPost:
		return true
	case path == "/v1/wallet/unload" && method == http.MethodPost:
		return true
	case path == "/v1/wallet/encrypt" && method == http.MethodPost:
		return true
	case path == "/v1/wallet/decrypt" && method == http.MethodPost:
		return true
	case path == "/v1/wallet/seed" && method == http.MethodPost:
		return true
	case path == "/v1/wallet/xpub" && method == http.MethodGet:
		return true
	case path == "/v1/wallet/balance" && method == http.MethodGet:
		return true
	case path == "/v2/wallet/recover" && method == http.MethodPost:
		return true
	}
	return false
}

// handleGetWalletsMulti aggregates wallets from all services
func handleGetWalletsMulti(c *gin.Context, services []*wallet.Service) {
	var allWallets []*walletResponse
	for _, svc := range services {
		wlts, err := svc.GetWallets()
		if err != nil {
			continue
		}
		for _, wlt := range wlts {
			wr, err := newWalletResponse(wlt)
			if err != nil {
				continue
			}
			allWallets = append(allWallets, wr)
		}
	}
	if allWallets == nil {
		allWallets = make([]*walletResponse, 0)
	}
	c.JSON(http.StatusOK, allWallets)
}

// handleWalletFolderMulti returns the primary wallet directory
func handleWalletFolderMulti(c *gin.Context, services []*wallet.Service) {
	addr, err := services[0].WalletDir()
	if err != nil {
		handleWalletError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"address": addr})
}

// handleReadOnlyPost converts read-only POST requests to GET to avoid CSRF issues.
// Returns true if the request was handled, false otherwise.
func handleReadOnlyPost(c *gin.Context, trimmedPath string, nodeURL string) bool {
	// Map of read-only POST endpoints and the form fields to forward as query params
	type readOnlyEndpoint struct {
		fields []string
	}
	endpoints := map[string]readOnlyEndpoint{
		"/v1/balance":      {fields: []string{"addrs"}},
		"/v1/transactions": {fields: []string{"addrs", "verbose"}},
		"/v1/outputs":      {fields: []string{"addrs", "hashes"}},
	}

	ep, ok := endpoints[trimmedPath]
	if !ok {
		return false
	}

	apiPath := trimmedPath
	query := ""
	for _, field := range ep.fields {
		val := c.Request.FormValue(field)
		if val != "" {
			if query != "" {
				query += "&"
			}
			query += field + "=" + val
		}
	}

	getURL := fmt.Sprintf("%s/api%s", nodeURL, apiPath)
	if query != "" {
		getURL += "?" + query
	}

	// Cache transaction queries for 30 seconds to avoid slow repeated lookups
	cacheKey := getURL
	if trimmedPath == "/v1/transactions" {
		if entry, ok := queryCache.get(cacheKey, 30*time.Second); ok {
			log.Printf("[PROXY] POST→GET %s -> cached (%s ago)", c.Request.URL.Path, time.Since(entry.cachedAt).Round(time.Second))
			c.Data(entry.statusCode, entry.contentType, entry.body)
			return true
		}
	}

	log.Printf("[PROXY] POST→GET %s -> %s", c.Request.URL.Path, getURL)

	resp, err := http.Get(getURL) //nolint:gosec
	if err != nil {
		errInternal(c, fmt.Sprintf("failed to query node: %v", err))
		return true
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			log.Printf("Error closing response body: %v", cerr)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		errInternal(c, fmt.Sprintf("failed to read response: %v", err))
		return true
	}

	// Cache successful responses
	if trimmedPath == "/v1/transactions" && resp.StatusCode == http.StatusOK {
		queryCache.set(cacheKey, proxyCacheEntry{
			body:        body,
			contentType: resp.Header.Get("Content-Type"),
			statusCode:  resp.StatusCode,
			cachedAt:    time.Now(),
		})
	}

	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
	return true
}

// fetchCSRFToken fetches a CSRF token from the remote node
func fetchCSRFToken(nodeURL string) (string, error) {
	resp, err := http.Get(nodeURL + "/api/v1/csrf") //nolint:gosec
	if err != nil {
		return "", fmt.Errorf("failed to fetch CSRF token: %v", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			log.Printf("Error closing CSRF response body: %v", cerr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("CSRF endpoint returned status %d", resp.StatusCode)
	}

	var result struct {
		Token string `json:"csrf_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode CSRF response: %v", err)
	}
	return result.Token, nil
}

// proxyToNode forwards an API request to the remote node
func proxyToNode(c *gin.Context, remoteNodeURL string) {
	proxyToNodeWithBase(c, remoteNodeURL, c.Request.URL.Path)
}

// proxyToNodeWithBase forwards a request with a custom target path
func proxyToNodeWithBase(c *gin.Context, remoteNodeURL string, targetPath string) {
	targetURL := remoteNodeURL + targetPath
	if c.Request.URL.RawQuery != "" {
		targetURL += "?" + c.Request.URL.RawQuery
	}

	log.Printf("[PROXY] %s %s -> %s", c.Request.Method, c.Request.URL.Path, targetURL)

	proxyReq, err := http.NewRequest(c.Request.Method, targetURL, c.Request.Body)
	if err != nil {
		log.Printf("[PROXY] Failed to create request: %v", err)
		c.String(http.StatusInternalServerError, "Failed to create proxy request")
		return
	}

	for name, values := range c.Request.Header {
		if name == "Referer" || name == "Origin" || name == "Host" {
			continue
		}
		for _, value := range values {
			proxyReq.Header.Add(name, value)
		}
	}

	// For POST/PUT/DELETE requests, fetch a CSRF token from the node
	if c.Request.Method != http.MethodGet && c.Request.Method != http.MethodHead {
		csrfToken, err := fetchCSRFToken(remoteNodeURL)
		if err != nil {
			log.Printf("[PROXY] Warning: could not fetch CSRF token: %v", err)
		} else if csrfToken != "" {
			proxyReq.Header.Set("X-CSRF-Token", csrfToken)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(proxyReq)
	if err != nil {
		log.Printf("[PROXY] Request failed: %v", err)
		c.String(http.StatusBadGateway, "Failed to proxy request to node: %v", err)
		return
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("Error closing response body: %v", err)
		}
	}()

	log.Printf("[PROXY] Response: %d", resp.StatusCode)

	for name, values := range resp.Header {
		if name != "Access-Control-Allow-Origin" && name != "Access-Control-Allow-Methods" && name != "Access-Control-Allow-Headers" {
			for _, value := range values {
				c.Header(name, value)
			}
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[PROXY] Failed to read response: %v", err)
		c.String(http.StatusInternalServerError, "Failed to read proxy response")
		return
	}

	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
}
