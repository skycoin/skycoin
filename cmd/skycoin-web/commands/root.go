// Package commands provides commands for the skycoin web interface.
package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	nethttppprof "net/http/pprof" //nolint:gosec
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/skycoin/skywire/pkg/skywire-utilities/pkg/calvin"
	"github.com/spf13/cobra"

	"github.com/skycoin/skycoin/src/btc"
	"github.com/skycoin/skycoin/src/cipher/bip44"
	"github.com/skycoin/skycoin/src/cipher/crypto"
	"github.com/skycoin/skycoin/src/fiber"
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

	guiDir string // custom GUI directory, overrides embedded GUI

	// Profiling flags
	pprofMode string
	pprofAddr string

	// Bitcoin flags
	btcNodeURL     string
	btcElectrumURL string
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
	CoinType          string `json:"coinType"`
	ServerWallets     bool   `json:"serverWallets"`
	// internal: the actual remote node URL (not exposed to frontend)
	remoteNodeURL string
}

// RootCmd is the root cil command
var RootCmd = &cobra.Command{
	Use:   "skycoin-web",
	Short: "Skycoin Web Wallet",
	Long: func() (ret string) {
		coinName := "skycoin"
		if fiberTomlPath := os.Getenv("FIBER_TOML"); fiberTomlPath != "" {
			if absPath, err := filepath.Abs(fiberTomlPath); err == nil {
				fiberTomlPath = absPath
			}
			if fiberCfg, err := fiber.NewConfig(filepath.Base(fiberTomlPath), filepath.Dir(fiberTomlPath)); err == nil {
				if fiberCfg.Node.DisplayName != "" {
					coinName = fiberCfg.Node.DisplayName
				}
			}
		}
		ret = calvin.AsciiFont(strings.ToLower(coinName) + "-web")
		ret += fmt.Sprintf("\nThin client web wallet for %s and fibercoins.", coinName)
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
	RootCmd.Flags().StringArrayVarP(&walletDirs, "wallet-dir", "w", nil, "Local wallet directory (e.g. ~/.skycoin/wallets)")
	RootCmd.Flags().BoolVar(&enableSeedAPI, "enable-seed-api", false, "Enable the wallet seed API (requires --wallet-dir)")
	RootCmd.Flags().StringVarP(&guiDir, "gui-dir", "g", "", "Custom GUI directory (overrides embedded GUI)")

	// Profiling flags
	RootCmd.Flags().StringVarP(&pprofMode, "pprofmode", "q", "", "[ cpu | mem | mutex | block | trace | http ]")
	RootCmd.Flags().StringVarP(&pprofAddr, "pprofaddr", "r", "localhost:6060", "pprof http port")

	// Bitcoin flags (mutually exclusive)
	RootCmd.Flags().StringVar(&btcNodeURL, "btc-node-url", "", "Bitcoin Core RPC URL (e.g. http://user:pass@127.0.0.1:8332)")
	RootCmd.Flags().StringVar(&btcElectrumURL, "btc-electrum-url", "", "Electrum server URL (e.g. ssl://electrum.blockstream.info:50002)")
	RootCmd.MarkFlagsMutuallyExclusive("btc-node-url", "btc-electrum-url")
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
		CoinType:          "skycoin",
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

// discoverCoins queries each configured node URL and returns discovered coins.
func discoverCoins() []*discoveredCoin {
	var coins []*discoveredCoin
	for i, rawURL := range nodeURLs {
		nodeURL := strings.TrimRight(rawURL, "/")
		if nodeURL == "" {
			continue
		}
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
				CoinType:      "skycoin",
				remoteNodeURL: nodeURL,
			}
		}
		coins = append(coins, coin)
		log.Printf("[COIN] Discovered %s (%s) at %s → proxy /coin/%d", coin.CoinName, coin.CoinSymbol, coin.remoteNodeURL, i)
	}
	return coins
}

// initWalletServices initializes Skycoin wallet services for each configured wallet directory.
func initWalletServices() []*wallet.Service {
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
	return wltServices
}

// initBitcoinBackend initializes the Bitcoin backend and wallet services if configured.
// It also appends a Bitcoin coin entry to the provided coins slice.
func initBitcoinBackend(coins []*discoveredCoin) (btc.Backend, []*wallet.Service, []*discoveredCoin) {
	var btcBackend btc.Backend
	var btcWltServices []*wallet.Service
	if btcNodeURL == "" && btcElectrumURL == "" {
		return nil, nil, coins
	}

	// Initialize Bitcoin backend
	var berr error
	if btcElectrumURL != "" {
		btcBackend, berr = btc.NewElectrumBackend(btcElectrumURL)
		if berr != nil {
			log.Printf("[WARN] Failed to connect to Electrum server %s: %v", btcElectrumURL, berr)
		} else {
			log.Printf("[BTC] Connected to Electrum server: %s", btcElectrumURL)
		}
	} else {
		btcBackend, berr = btc.NewCoreBackend(btcNodeURL)
		if berr != nil {
			log.Printf("[WARN] Failed to connect to Bitcoin Core %s: %v", btcNodeURL, berr)
		} else {
			log.Printf("[BTC] Connected to Bitcoin Core: %s", btcNodeURL)
		}
	}

	if btcBackend != nil {
		// Create Bitcoin wallet services using the same --wallet-dir directories.
		// The wallet service filters by coin type, so BTC and SKY wallets coexist
		// in the same directory without conflict.
		for _, dir := range walletDirs {
			if dir == "" {
				continue
			}
			btcBip44Coin := bip44.CoinTypeBitcoin
			btcCfg := wallet.Config{
				WalletDir:       dir,
				CryptoType:      crypto.DefaultCryptoType,
				EnableWalletAPI: true,
				EnableSeedAPI:   enableSeedAPI,
				Bip44Coin:       &btcBip44Coin,
			}
			btcSvc, btcErr := wallet.NewService(btcCfg)
			if btcErr != nil {
				log.Fatalf("Failed to initialize Bitcoin wallet service for %s: %v", dir, btcErr)
			}
			btcWltServices = append(btcWltServices, btcSvc)
			log.Printf("[BTC] Bitcoin wallet service initialized: %s", dir)
		}
		if len(walletDirs) == 0 {
			log.Printf("[BTC] No --wallet-dir specified, Bitcoin wallet management disabled (web-only mode)")
		}

		// Add Bitcoin as a discovered coin
		btcCoinIndex := len(coins)
		btcCoin := &discoveredCoin{
			ID:                btcCoinIndex,
			NodeURL:           fmt.Sprintf("/coin/%d", btcCoinIndex),
			CoinName:          "Bitcoin",
			CoinSymbol:        "BTC",
			HoursName:         "",
			PriceTickerID:     "btc-bitcoin",
			PriceTickerSource: "coinpaprika",
			CoinExplorer:      "https://blockchair.com/bitcoin",
			CoinType:          "bitcoin",
		}
		coins = append(coins, btcCoin)
		log.Printf("[COIN] Added Bitcoin at index %d", btcCoinIndex)
	}

	return btcBackend, btcWltServices, coins
}

// mapWalletsToCoin maps wallet services to coin indices and sets up Bitcoin handlers.
// It also marks coins that have server-side wallet management.
func mapWalletsToCoin(coins []*discoveredCoin, wltServices []*wallet.Service, btcBackend btc.Backend, btcWltServices []*wallet.Service) (map[int][]*wallet.Service, map[int]*btcHandler) {
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

	// Map Bitcoin wallet services and handler to the Bitcoin coin
	btcHandlers := make(map[int]*btcHandler)
	if btcBackend != nil {
		for i, coin := range coins {
			if coin.CoinType == "bitcoin" {
				if len(btcWltServices) > 0 {
					coinWltServices[i] = btcWltServices
					btcHandlers[i] = &btcHandler{
						backend:    btcBackend,
						wltService: btcWltServices[0],
					}
				}
			}
		}
	}

	// Mark coins that have server-side wallet management
	for i := range coins {
		if _, ok := coinWltServices[i]; ok {
			coins[i].ServerWallets = true
		}
	}

	return coinWltServices, btcHandlers
}

// initGUIFS returns the GUI filesystem, either from a custom directory or embedded.
func initGUIFS() fs.FS {
	if guiDir != "" {
		log.Printf("Serving GUI from local folder: %s", guiDir)
		return os.DirFS(guiDir)
	}
	log.Println("Serving embedded GUI")
	guiFS, err := fs.Sub(gui.DistFS, "dist")
	if err != nil {
		log.Fatalf("Failed to get dist subdirectory: %v", err)
	}
	return guiFS
}

func serve() {
	stopPProf := initPProf(pprofMode, pprofAddr)
	defer stopPProf()

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	coins := discoverCoins()
	wltServices := initWalletServices()
	btcBackend, btcWltServices, coins := initBitcoinBackend(coins)
	coinWltServices, btcHandlers := mapWalletsToCoin(coins, wltServices, btcBackend, btcWltServices)
	guiFS := initGUIFS()

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

		// Bitcoin coins are handled by the btcHandler, not proxied to a node
		if coin.CoinType == "bitcoin" {
			// Return stub responses for Skycoin-specific endpoints the frontend calls on all coins
			if handleBtcStubEndpoints(c, apiPath) {
				return
			}

			if handler, ok := btcHandlers[coinIndex]; ok {
				// Try BTC-specific API endpoints first
				if handler.handleBtcAPI(c, apiPath) {
					return
				}
			}
			// Try wallet endpoints (create, list, etc.)
			if services, ok := coinWltServices[coinIndex]; ok && len(services) > 0 {
				if handleMultiWalletAPI(c, apiPath, services, "") {
					return
				}
			}
			c.String(http.StatusNotFound, "endpoint not available for bitcoin")
			return
		}

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

	// Serve static files from GUI filesystem
	router.NoRoute(func(c *gin.Context) {
		if c.Request.URL.Path != "/" && c.Request.URL.Path != "/favicon.ico" {
			log.Printf("[STATIC] %s %s", c.Request.Method, c.Request.URL.Path)
		}
		fileServer := http.FileServer(http.FS(guiFS))
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

// handleBtcStubEndpoints returns stub responses for Skycoin-specific endpoints
// that the frontend calls on all coins (network/connections, health, blockchain/progress, etc.)
func handleBtcStubEndpoints(c *gin.Context, apiPath string) bool {
	path := strings.TrimSuffix(apiPath, "/")
	switch path {
	case "/v1/network/connections":
		c.JSON(http.StatusOK, gin.H{"connections": []any{}})
	case "/v1/health":
		c.JSON(http.StatusOK, gin.H{
			"blockchain":       gin.H{"head": gin.H{"seq": 0, "timestamp": 0}},
			"version":          gin.H{"version": "0.27.0", "commit": "bitcoin"},
			"open_connections": 0,
			"uptime":           "0s",
		})
	case "/v1/blockchain/progress":
		c.JSON(http.StatusOK, gin.H{"current": 1, "highest": 1, "peers": []any{}})
	case "/v1/blockchain/metadata":
		c.JSON(http.StatusOK, gin.H{"head": gin.H{"seq": 0, "fee": 0}})
	case "/v1/csrf":
		c.JSON(http.StatusOK, gin.H{"csrf_token": ""})
	default:
		return false
	}
	return true
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
	var allWallets []*readable.WalletResponse
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
		allWallets = make([]*readable.WalletResponse, 0)
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

// handleReadOnlyPost forwards read-only POST requests to the node as POST with CSRF token.
// This preserves form body data (e.g. long address lists) that would exceed URI length
// limits if converted to GET query parameters.
// Returns true if the request was handled, false otherwise.
func handleReadOnlyPost(c *gin.Context, trimmedPath string, nodeURL string) bool {
	readOnlyEndpoints := map[string]bool{
		"/v1/balance":      true,
		"/v1/transactions": true,
		"/v1/outputs":      true,
	}

	if !readOnlyEndpoints[trimmedPath] {
		return false
	}

	targetURL := fmt.Sprintf("%s/api%s", nodeURL, trimmedPath)

	// Cache transaction queries for 30 seconds to avoid slow repeated lookups
	if trimmedPath == "/v1/transactions" {
		cacheKey := targetURL + "?" + c.Request.FormValue("addrs")
		if entry, ok := queryCache.get(cacheKey, 30*time.Second); ok {
			log.Printf("[PROXY] POST %s -> cached (%s ago)", c.Request.URL.Path, time.Since(entry.cachedAt).Round(time.Second))
			c.Data(entry.statusCode, entry.contentType, entry.body)
			return true
		}
	}

	log.Printf("[PROXY] POST %s -> %s", c.Request.URL.Path, targetURL)

	// Build form body from the original request
	if err := c.Request.ParseForm(); err != nil {
		errInternal(c, fmt.Sprintf("failed to parse form: %v", err))
		return true
	}
	formData := c.Request.PostForm.Encode()

	req, err := http.NewRequest(http.MethodPost, targetURL, strings.NewReader(formData))
	if err != nil {
		errInternal(c, fmt.Sprintf("failed to create request: %v", err))
		return true
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Fetch and attach CSRF token
	csrfToken, err := fetchCSRFToken(nodeURL)
	if err != nil {
		log.Printf("[PROXY] Warning: could not fetch CSRF token: %v", err)
	} else if csrfToken != "" {
		req.Header.Set("X-CSRF-Token", csrfToken)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
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

	// Cache successful transaction responses
	if trimmedPath == "/v1/transactions" && resp.StatusCode == http.StatusOK {
		cacheKey := targetURL + "?" + c.Request.FormValue("addrs")
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

func initPProf(profMode string, profAddr string) (stop func()) {
	stop = func() {}
	switch profMode {
	case "http":
		go func() {
			mux := http.NewServeMux()
			mux.HandleFunc("/debug/pprof/", nethttppprof.Index)
			mux.HandleFunc("/debug/pprof/cmdline", nethttppprof.Cmdline)
			mux.HandleFunc("/debug/pprof/profile", nethttppprof.Profile)
			mux.HandleFunc("/debug/pprof/symbol", nethttppprof.Symbol)
			mux.HandleFunc("/debug/pprof/trace", nethttppprof.Trace)
			srv := &http.Server{ //nolint:gosec
				Addr:              profAddr,
				Handler:           mux,
				ReadHeaderTimeout: 5 * time.Second,
				WriteTimeout:      30 * time.Second,
			}
			log.Printf("Serving pprof on http://%s/debug/pprof/", profAddr)
			if err := srv.ListenAndServe(); err != nil {
				log.Printf("pprof http server stopped: %v", err)
			}
		}()
	}
	return stop
}
