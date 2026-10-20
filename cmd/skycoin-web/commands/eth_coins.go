package commands

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/skycoin/skycoin/src/cipher/bip44"
	"github.com/skycoin/skycoin/src/cipher/crypto"
	"github.com/skycoin/skycoin/src/eth"
	"github.com/skycoin/skycoin/src/wallet"
)

// Ethereum flags
var (
	ethRPCURL     string
	ethIndexerURL string
	erc20Tokens   string
)

// defaultERC20Tokens is USDT, as SYMBOL:contract:decimals.
const defaultERC20Tokens = "USDT:0xdAC17F958D2ee523a2206206994597C13D831ec7:6" //nolint:gosec // a public contract address, not a credential

func init() {
	RootCmd.Flags().StringVar(&ethRPCURL, "eth-rpc-url", "default", "Ethereum JSON-RPC URLs, comma separated and tried in order, default for the built-in list, or none for no ETH")
	RootCmd.Flags().StringVar(&ethIndexerURL, "eth-indexer-url", eth.DefaultIndexerURL, "Etherscan-style API for Ethereum history, such as a Blockscout instance")
	RootCmd.Flags().StringVar(&erc20Tokens, "erc20-tokens", defaultERC20Tokens, "ERC-20 tokens to offer, comma separated SYMBOL:contract:decimals, or none")
}

// ethRPCURLs is the URL list --eth-rpc-url names, or nil for none.
func ethRPCURLs(flag string) []string {
	switch strings.TrimSpace(flag) {
	case "", "none":
		return nil
	case "default":
		return eth.DefaultRPCURLs
	}
	var urls []string
	for _, u := range strings.Split(flag, ",") {
		if u = strings.TrimSpace(u); u != "" {
			urls = append(urls, u)
		}
	}
	return urls
}

// parseERC20Tokens reads --erc20-tokens.
func parseERC20Tokens(flag string) ([]ethToken, error) {
	flag = strings.TrimSpace(flag)
	if flag == "" || flag == "none" {
		return nil, nil
	}
	var tokens []ethToken
	for _, spec := range strings.Split(flag, ",") {
		parts := strings.Split(strings.TrimSpace(spec), ":")
		if len(parts) != 3 || parts[0] == "" {
			return nil, fmt.Errorf("--erc20-tokens: %q is not SYMBOL:contract:decimals", spec)
		}
		contract, err := eth.ParseAddress(parts[1])
		if err != nil {
			return nil, fmt.Errorf("--erc20-tokens: %s contract: %w", parts[0], err)
		}
		decimals, err := strconv.Atoi(parts[2])
		if err != nil || decimals < 0 || decimals > 36 {
			return nil, fmt.Errorf("--erc20-tokens: %s decimals %q", parts[0], parts[2])
		}
		tokens = append(tokens, ethToken{Symbol: parts[0], Contract: contract, Decimals: decimals})
	}
	return tokens, nil
}

// knownTokenTickers are coinpaprika ids for tokens the wallet can price.
var knownTokenTickers = map[string]string{"USDT": "usdt-tether", "USDC": "usdc-usd-coin"}

// initEthereumBackend appends Ether and the configured tokens to coins, with
// handlers keyed by coin index, when RPC URLs are configured.
func initEthereumBackend(coins []*discoveredCoin) ([]*discoveredCoin, map[int]*ethHandler, []*wallet.Service, error) {
	urls := ethRPCURLs(ethRPCURL)
	if len(urls) == 0 {
		return coins, nil, nil, nil
	}
	tokens, err := parseERC20Tokens(erc20Tokens)
	if err != nil {
		return coins, nil, nil, err
	}
	client := &http.Client{Transport: nodeTransport, Timeout: 30 * time.Second}
	rpc, err := eth.NewRPC(urls, client)
	if err != nil {
		return coins, nil, nil, err
	}
	indexer := eth.NewIndexer(ethIndexerURL, client)

	var services []*wallet.Service
	var endpointsDir string
	for _, dir := range walletDirs {
		if dir == "" {
			continue
		}
		ethCoin := bip44.CoinTypeEthereum
		svc, err := wallet.NewService(wallet.Config{
			WalletDir:       dir,
			CryptoType:      crypto.DefaultCryptoType,
			EnableWalletAPI: true,
			EnableSeedAPI:   enableSeedAPI,
			Bip44Coin:       &ethCoin,
		})
		if err != nil {
			return coins, nil, nil, fmt.Errorf("ethereum wallet service for %s: %w", dir, err)
		}
		if endpointsDir == "" {
			endpointsDir = dir
		}
		services = append(services, svc)
	}
	if len(services) == 0 {
		log.Printf("[ETH] No --wallet-dir specified, Ethereum wallets disabled")
		return coins, nil, nil, nil
	}

	endpoints := newEthEndpoints(rpc, indexer, endpointsDir)
	if err := endpoints.load(); err != nil {
		log.Printf("[ETH] Ignoring the saved RPC and indexer addresses: %v", err)
	}

	handlers := make(map[int]*ethHandler)
	add := func(c *discoveredCoin, token *ethToken) {
		c.ID = len(coins)
		c.NodeURL = fmt.Sprintf("/coin/%d", c.ID)
		c.NodeOptions = urls
		c.CoinExplorer = strings.TrimSuffix(ethIndexerURL, "/")
		c.ServerWallets = true
		coins = append(coins, c)
		handlers[c.ID] = &ethHandler{rpc: rpc, indexer: indexer, wltService: services[0], token: token, endpoints: endpoints}
		log.Printf("[COIN] Added %s at index %d", c.CoinSymbol, c.ID)
	}
	add(&discoveredCoin{
		CoinName: "Ethereum", CoinSymbol: "ETH", CoinType: "ethereum", Decimals: 18,
		PriceTickerID: "eth-ethereum", PriceTickerSource: "coinpaprika",
	}, nil)
	for i := range tokens {
		t := tokens[i]
		add(&discoveredCoin{
			CoinName: t.Symbol, CoinSymbol: t.Symbol, CoinType: "erc20", Decimals: t.Decimals,
			Contract: t.Contract.String(), PriceTickerID: knownTokenTickers[t.Symbol], PriceTickerSource: "coinpaprika",
		}, &t)
	}
	log.Printf("[ETH] RPC endpoints, in order: %s", strings.Join(urls, ", "))
	return coins, handlers, services, nil
}

// isEthereumCoin reports whether a coin is Ether or an ERC-20 token.
func isEthereumCoin(c *discoveredCoin) bool {
	return c.CoinType == "ethereum" || c.CoinType == "erc20"
}
