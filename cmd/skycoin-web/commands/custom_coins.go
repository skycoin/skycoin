package commands

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/skycoin/skycoin/src/eth"
	"github.com/skycoin/skycoin/src/wallet"
)

// coinTable is every coin the server offers, by index, with what serves it.
// It grows when the wallet adds a coin, so it is read under coinsMu.
type coinTable struct {
	coins []*discoveredCoin
	eth   map[int]*ethHandler
	btc   map[int]*btcHandler
	wlt   map[int][]*wallet.Service
}

// get is coin i and what serves it, or a nil coin when there is none.
func (t *coinTable) get(i int) (*discoveredCoin, *ethHandler, *btcHandler, []*wallet.Service) {
	coinsMu.RLock()
	defer coinsMu.RUnlock()
	if i < 0 || i >= len(t.coins) || t.coins[i].hidden {
		return nil, nil, nil, nil
	}
	return t.coins[i], t.eth[i], t.btc[i], t.wlt[i]
}

// visible are the coins the wallet lists. Removed ones keep their index, so
// the other coins' addresses do not move. Call it under coinsMu.
func (t *coinTable) visible() []*discoveredCoin {
	out := make([]*discoveredCoin, 0, len(t.coins))
	for _, c := range t.coins {
		if !c.hidden {
			out = append(out, c)
		}
	}
	return out
}

// customCoinsFile keeps the coins added in the wallet, beside its wallets.
const customCoinsFile = "custom-coins.json"

// maxCoinIcon bounds an icon, which is kept and served as a data URL.
const maxCoinIcon = 96 << 10

// customCoin is a fiber coin or an ERC-20 token added in the wallet.
type customCoin struct {
	ID       string `json:"id"`
	Kind     string `json:"kind"`
	Name     string `json:"name"`
	Ticker   string `json:"ticker"`
	NodeURL  string `json:"node_url,omitempty"`
	Contract string `json:"contract,omitempty"`
	Decimals int    `json:"decimals,omitempty"`
	Icon     string `json:"icon,omitempty"`
}

type customCoins struct {
	table *coinTable
	file  string

	mu   sync.Mutex // serializes changes and saves
	list []customCoin
}

func newCustomCoins(table *coinTable) *customCoins {
	cc := &customCoins{table: table}
	for _, dir := range walletDirs {
		if dir != "" {
			cc.file = filepath.Join(dir, customCoinsFile)
			break
		}
	}
	return cc
}

// load adds the saved coins. A fiber coin whose node does not answer is still
// added under its saved name and ticker.
func (cc *customCoins) load() error {
	if cc.file == "" {
		return nil
	}
	b, err := os.ReadFile(cc.file)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	var saved []customCoin
	if err := json.Unmarshal(b, &saved); err != nil {
		return fmt.Errorf("%s: %w", cc.file, err)
	}
	cc.mu.Lock()
	defer cc.mu.Unlock()
	for _, c := range saved {
		if _, err := cc.add(c, false); err != nil {
			log.Printf("[COIN] Skipping saved coin %s: %v", c.Name, err)
			continue
		}
		cc.list = append(cc.list, c)
	}
	return nil
}

// handle serves /v1/coins/custom: GET lists, POST adds and DELETE ?id= removes.
func (cc *customCoins) handle(c *webCtx, apiPath string) bool {
	if strings.TrimSuffix(apiPath, "/") != "/v1/coins/custom" {
		return false
	}
	if cc.file == "" {
		c.JSON(http.StatusNotFound, H{"error": "adding coins needs --wallet-dir"})
		return true
	}
	cc.mu.Lock()
	defer cc.mu.Unlock()
	switch c.Request.Method {
	case http.MethodGet:
		c.JSON(http.StatusOK, append([]customCoin{}, cc.list...))
	case http.MethodPost:
		var req customCoin
		if err := json.NewDecoder(http.MaxBytesReader(c.Writer, c.Request.Body, 2*maxCoinIcon)).Decode(&req); err != nil {
			c.JSON(http.StatusBadRequest, H{"error": "invalid JSON body"})
			return true
		}
		req.ID = newCustomCoinID()
		coin, err := cc.add(req, true)
		if err != nil {
			c.JSON(http.StatusBadRequest, H{"error": err.Error()})
			return true
		}
		cc.list = append(cc.list, cc.stored(req, coin))
		if err := cc.save(); err != nil {
			c.JSON(http.StatusInternalServerError, H{"error": "saving the coin: " + err.Error()})
			return true
		}
		c.JSON(http.StatusOK, coin)
	case http.MethodDelete:
		id := c.Query("id")
		for i, saved := range cc.list {
			if saved.ID != id {
				continue
			}
			coinsMu.Lock()
			for _, coin := range cc.table.coins {
				if coin.CustomID == id {
					coin.hidden = true
				}
			}
			coinsMu.Unlock()
			cc.list = append(cc.list[:i], cc.list[i+1:]...)
			if err := cc.save(); err != nil {
				c.JSON(http.StatusInternalServerError, H{"error": "saving the coins: " + err.Error()})
				return true
			}
			c.JSON(http.StatusOK, H{"removed": id})
			return true
		}
		c.JSON(http.StatusNotFound, H{"error": "no added coin " + id})
	default:
		c.JSON(http.StatusMethodNotAllowed, H{"error": "method not allowed"})
	}
	return true
}

// stored is what is kept for req: the name and ticker the node reported when
// none were given, so a restart shows the same coin.
func (cc *customCoins) stored(req customCoin, coin *discoveredCoin) customCoin {
	req.Name = coin.CoinName
	req.Ticker = coin.CoinSymbol
	return req
}

func (cc *customCoins) save() error {
	b, err := json.MarshalIndent(cc.list, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(cc.file, b, 0600)
}

// add checks req and appends its coin to the table. With verify a fiber
// coin's node must answer.
func (cc *customCoins) add(req customCoin, verify bool) (*discoveredCoin, error) {
	req.Name = strings.TrimSpace(req.Name)
	req.Ticker = strings.ToUpper(strings.TrimSpace(req.Ticker))
	if req.Icon != "" {
		if len(req.Icon) > maxCoinIcon || !iconDataURL(req.Icon) {
			return nil, errors.New("the icon must be a PNG, JPEG, WebP or GIF image under 96 KB")
		}
	}

	switch req.Kind {
	case "fiber":
		return cc.addFiber(req, verify)
	case "erc20":
		return cc.addToken(req)
	default:
		return nil, errors.New(`kind must be "fiber" or "erc20"`)
	}
}

func (cc *customCoins) addFiber(req customCoin, verify bool) (*discoveredCoin, error) {
	nodeURL, err := checkedHTTPURL(req.NodeURL)
	if err != nil {
		return nil, err
	}

	coinsMu.Lock()
	index := len(cc.table.coins)
	coinsMu.Unlock()
	coin, err := discoverCoin(index, nodeURL)
	if err != nil {
		if verify {
			return nil, fmt.Errorf("the node did not answer: %w", err)
		}
		coin = &discoveredCoin{CoinType: "skycoin", HoursName: "Coin Hours", placeholder: true, remoteNodeURL: nodeURL}
	}
	if req.Name != "" {
		coin.CoinName = req.Name
	}
	if req.Ticker != "" {
		coin.CoinSymbol = req.Ticker
	}
	if coin.CoinName == "" || coin.CoinSymbol == "" {
		return nil, errors.New("give the coin a name and a ticker")
	}

	coinsMu.Lock()
	defer coinsMu.Unlock()
	cc.place(coin, req)
	// Fiber coins share the wallet directories of the first configured coin.
	if services := cc.table.wlt[0]; len(services) > 0 && !isEthereumCoin(cc.table.coins[0]) {
		cc.table.wlt[coin.ID] = services
		coin.ServerWallets = true
	}
	return coin, nil
}

func (cc *customCoins) addToken(req customCoin) (*discoveredCoin, error) {
	contract, err := eth.ParseAddress(strings.TrimSpace(req.Contract))
	if err != nil {
		return nil, fmt.Errorf("contract: %w", err)
	}
	if req.Decimals < 0 || req.Decimals > 36 {
		return nil, errors.New("decimals must be from 0 to 36")
	}
	if req.Ticker == "" {
		return nil, errors.New("give the token a ticker")
	}
	if req.Name == "" {
		req.Name = req.Ticker
	}

	coinsMu.Lock()
	defer coinsMu.Unlock()
	var ether *ethHandler
	for _, h := range cc.table.eth {
		ether = h
		break
	}
	if ether == nil {
		return nil, errors.New("Ethereum is not enabled on this server")
	}
	for i, h := range cc.table.eth {
		if h.token != nil && h.token.Contract == contract && !cc.table.coins[i].hidden {
			return nil, fmt.Errorf("%s is already listed", cc.table.coins[i].CoinSymbol)
		}
	}
	coin := &discoveredCoin{
		CoinName: req.Name, CoinSymbol: req.Ticker, CoinType: "erc20", Decimals: req.Decimals,
		Contract: contract.String(), PriceTickerID: knownTokenTickers[req.Ticker], PriceTickerSource: "coinpaprika",
		CoinExplorer: ether.indexer.Base(), ServerWallets: true,
	}
	cc.place(coin, req)
	token := &ethToken{Symbol: req.Ticker, Contract: contract, Decimals: req.Decimals}
	cc.table.eth[coin.ID] = &ethHandler{rpc: ether.rpc, indexer: ether.indexer, wltService: ether.wltService, token: token, endpoints: ether.endpoints}
	for i, h := range cc.table.eth {
		if h == ether {
			cc.table.wlt[coin.ID] = cc.table.wlt[i]
			break
		}
	}
	return coin, nil
}

// place appends coin to the table. Call it under coinsMu.
func (cc *customCoins) place(coin *discoveredCoin, req customCoin) {
	coin.ID = len(cc.table.coins)
	coin.NodeURL = fmt.Sprintf("/coin/%d", coin.ID)
	coin.CustomID = req.ID
	coin.IconURL = req.Icon
	cc.table.coins = append(cc.table.coins, coin)
	log.Printf("[COIN] Added %s at index %d", coin.CoinSymbol, coin.ID)
}

func iconDataURL(s string) bool {
	for _, t := range []string{"png", "jpeg", "webp", "gif"} {
		if strings.HasPrefix(s, "data:image/"+t+";base64,") {
			return true
		}
	}
	return false
}

func newCustomCoinID() string {
	b := make([]byte, 4)
	_, _ = rand.Read(b) //nolint:errcheck // crypto/rand does not fail
	return hex.EncodeToString(b)
}
