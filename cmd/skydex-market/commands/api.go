// Package commands cmd/apps/skydex-market/commands/api.go
//
// Operator API consumed by the market's local dashboard. It exposes the
// market's configuration (explorers, SKY fullnode, escrow wallet, ban rules)
// for reading and editing, plus read-only monitoring of products, orders and
// active bans. It is a localhost operator surface, not part of the client<->
// market dmsg protocol.
package commands

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/skycoin/skycoin/src/skydex/chain"
	"github.com/skycoin/skycoin/src/skydex/db"
	"github.com/skycoin/skycoin/src/skydex/protocol"
)

// editableConfigKeys is the whitelist of fixed market_config keys the operator
// UI may change. Per-coin explorer keys (explorer_<coin>_{provider,url,key}) are
// validated dynamically in validateConfigKey.
var editableConfigKeys = map[string]bool{
	"market_name":             true,
	"wallet_sky":              true,
	"sky_fullnode_url":        true,
	"sky_wallet_seed":         true,
	"commission_rate_percent": true,
	"commission_min_sky":      true,
	"commission_max_sky":      true,
	"confirmations_required":  true,
	"min_trade_sky":           true,
	"max_trade_sky":           true,
	"freeze_violations_limit": true,
	"ban_duration_days":       true,
	"buy_cancel_limit":        true,
	"listing_expiry_minutes":  true,
	"order_expiry_minutes":    true,
	"cleanup_days":            true,
}

// isSecretConfigKey reports whether a market_config key holds a credential that
// must never leave the process.
//
// sky_wallet_seed is the escrow hot-wallet signing key; explorer_<coin>_key is a
// third-party API key. Both are settable by the operator but neither is ever
// readable back — the same discipline db.SellCoin applies to its per-coin
// wallet_seed via `json:"-"`.
func isSecretConfigKey(key string) bool {
	if key == "sky_wallet_seed" {
		return true
	}
	_, field, ok := parseExplorerKey(key)

	return ok && field == "key"
}

// redactConfig copies cfg, replacing every secret value with the empty string
// and reporting separately which secrets are set. The UI renders those as
// "configured" without ever receiving the value, and submits blank to keep it.
func redactConfig(cfg map[string]string) map[string]any {
	values := make(map[string]string, len(cfg))
	set := make(map[string]bool)
	for k, v := range cfg {
		if isSecretConfigKey(k) {
			values[k] = ""
			set[k] = strings.TrimSpace(v) != ""
			continue
		}
		values[k] = v
	}

	return map[string]any{"config": values, "secrets_set": set}
}

// validateConfigKey allows a fixed whitelisted key, or a per-coin explorer key
// for ANY coin symbol (the operator may add arbitrary payment coins). For a
// provider assignment it checks the chosen provider can actually verify that coin.
func validateConfigKey(key, value string) error {
	if editableConfigKeys[key] {
		return nil
	}
	coin, field, ok := parseExplorerKey(key)
	if !ok {
		return fmt.Errorf("unknown or read-only config key: %s", key)
	}
	if field == "provider" && strings.TrimSpace(value) != "" && !chain.SupportsProvider(coin, value) {
		return fmt.Errorf("explorer provider %q does not support %s", value, coin)
	}
	return nil
}

// validateExplorerEnable rejects enabling a coin (setting a non-empty provider)
// when no explorer URL is available — neither a built-in default, a URL in this
// same update, nor one already stored. Such a coin could never verify a payment,
// so enabling it would be a silent misconfiguration.
func validateExplorerEnable(updates map[string]string, database *db.Database) error {
	for k, v := range updates {
		coin, field, ok := parseExplorerKey(k)
		if !ok || field != "provider" || strings.TrimSpace(v) == "" {
			continue
		}
		if chain.DefaultExplorerURL(coin) != "" {
			continue // has a built-in endpoint
		}
		if u, ok := updates["explorer_"+strings.ToLower(coin)+"_url"]; ok && strings.TrimSpace(u) != "" {
			continue // a URL is being set in this same request
		}
		_, storedURL, _, err := database.ExplorerConfig(coin)
		if err != nil {
			return err
		}
		if strings.TrimSpace(storedURL) != "" {
			continue // a URL was configured earlier
		}
		return fmt.Errorf("enable %s: set an explorer URL first (it has no built-in default)", coin)
	}
	return nil
}

// parseExplorerKey splits "explorer_<coin>_<field>" into (COIN, field). field is
// the last "_"-segment; the coin is everything between.
func parseExplorerKey(key string) (coin, field string, ok bool) {
	rest, found := strings.CutPrefix(key, "explorer_")
	if !found {
		return "", "", false
	}
	i := strings.LastIndex(rest, "_")
	if i <= 0 {
		return "", "", false
	}
	field = rest[i+1:]
	if field != "provider" && field != "url" && field != "key" {
		return "", "", false
	}
	return strings.ToUpper(rest[:i]), field, true
}

// registerOperatorAPI mounts the operator API routes on mux.
func registerOperatorAPI(mux *http.ServeMux, database *db.Database, host Host) {
	// GET /api/info — market identity (its public key) for clients to use.
	mux.HandleFunc("/api/info", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			mWriteError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		mWriteJSON(w, http.StatusOK, map[string]any{"public_key": host.PubKey()})
	})

	// GET  /api/config — all market config; POST — update a subset.
	mux.HandleFunc("/api/config", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			cfg, err := database.GetAllConfig()
			if err != nil {
				mWriteError(w, http.StatusInternalServerError, err.Error())
				return
			}
			mWriteJSON(w, http.StatusOK, redactConfig(cfg))
		case http.MethodPost:
			var updates map[string]string
			if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
				mWriteError(w, http.StatusBadRequest, "invalid config payload")
				return
			}
			for k, v := range updates {
				if err := validateConfigKey(k, v); err != nil {
					mWriteError(w, http.StatusBadRequest, err.Error())
					return
				}
			}
			// Secrets are write-only: the UI can never read them back, so it
			// submits blank to mean "keep what's stored". Without this, every
			// config save would wipe the escrow seed.
			for k, v := range updates {
				if isSecretConfigKey(k) && strings.TrimSpace(v) == "" {
					delete(updates, k)
				}
			}
			if err := validateExplorerEnable(updates, database); err != nil {
				mWriteError(w, http.StatusBadRequest, err.Error())
				return
			}
			for k, v := range updates {
				if err := database.SetConfig(k, v); err != nil {
					mWriteError(w, http.StatusInternalServerError, err.Error())
					return
				}
			}
			cfg, err := database.GetAllConfig()
			if err != nil {
				mWriteError(w, http.StatusInternalServerError, err.Error())
				return
			}
			// Redacted here too — the POST response echoes the whole config.
			mWriteJSON(w, http.StatusOK, redactConfig(cfg))
		default:
			mWriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	})

	// GET /api/products — all products with status.
	mux.HandleFunc("/api/products", func(w http.ResponseWriter, _ *http.Request) {
		products, err := database.GetAllProducts()
		if err != nil {
			mWriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
		mWriteJSON(w, http.StatusOK, map[string]any{"products": products})
	})

	// GET /api/orders — all orders.
	mux.HandleFunc("/api/orders", func(w http.ResponseWriter, _ *http.Request) {
		orders, err := database.GetAllOrders()
		if err != nil {
			mWriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
		mWriteJSON(w, http.StatusOK, map[string]any{"orders": orders})
	})

	// GET /api/currencies — per supported coin: which explorer providers can
	// verify it, and the operator's current explorer config. Drives the UI rows.
	mux.HandleFunc("/api/currencies", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			mWriteError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		type coinCfg struct {
			Code      string   `json:"code"`
			Providers []string `json:"providers"`
			Provider  string   `json:"provider"`
			URL       string   `json:"url"`
			// HasKey replaces the explorer API key itself: the operator sets it
			// but never reads it back, mirroring has_seed on /api/sellcoins.
			HasKey     bool   `json:"has_key"`
			DefaultURL string `json:"default_url"` // built-in explorer endpoint, if any
			Available  bool   `json:"available"`   // enabled (a provider is configured)
		}
		// The canonical coins (so their built-in defaults are always offered),
		// plus any custom coin the operator has already configured.
		coins := append([]string{}, protocol.PaymentCurrencies...)
		seen := map[string]bool{}
		for _, c := range coins {
			seen[c] = true
		}
		if configured, err := database.ConfiguredPaymentCurrencies(); err == nil {
			for _, c := range configured {
				if !seen[c] {
					coins = append(coins, c)
					seen[c] = true
				}
			}
		}
		out := make([]coinCfg, 0, len(coins))
		for _, c := range coins {
			provider, url, key, err := database.ExplorerConfig(c)
			if err != nil {
				mWriteError(w, http.StatusInternalServerError, err.Error())
				return
			}
			providers := chain.ProvidersFor(c)
			if providers == nil {
				providers = []string{}
			}
			out = append(out, coinCfg{
				Code:       c,
				Providers:  providers,
				Provider:   provider,
				URL:        url,
				HasKey:     strings.TrimSpace(key) != "",
				DefaultURL: chain.DefaultExplorerURL(c),
				Available:  strings.TrimSpace(provider) != "",
			})
		}
		mWriteJSON(w, http.StatusOK, map[string]any{"currencies": out})
	})

	// GET  /api/sellcoins — every configured sell coin (SKY + fibercoins), with
	//                       escrow config (seed masked). POST — add/edit one.
	mux.HandleFunc("/api/sellcoins", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			coins, err := database.ListSellCoins()
			if err != nil {
				mWriteError(w, http.StatusInternalServerError, err.Error())
				return
			}
			// Report whether a seed is set without ever exposing it.
			type coinView struct {
				*db.SellCoin
				HasSeed bool `json:"has_seed"`
			}
			out := make([]coinView, 0, len(coins))
			for _, c := range coins {
				out = append(out, coinView{SellCoin: c, HasSeed: strings.TrimSpace(c.WalletSeed) != ""})
			}
			mWriteJSON(w, http.StatusOK, map[string]any{"sellcoins": out})
		case http.MethodPost:
			var req struct {
				Symbol        string `json:"symbol"`
				Name          string `json:"name"`
				NodeURL       string `json:"node_url"`
				WalletSeed    string `json:"wallet_seed"`
				WalletAddr    string `json:"wallet_addr"`
				Confirmations int    `json:"confirmations"`
				Enabled       bool   `json:"enabled"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				mWriteError(w, http.StatusBadRequest, "invalid sell-coin payload")
				return
			}
			if strings.TrimSpace(req.Symbol) == "" {
				mWriteError(w, http.StatusBadRequest, "symbol is required")
				return
			}
			sc := &db.SellCoin{
				Symbol:        req.Symbol,
				Name:          req.Name,
				NodeURL:       req.NodeURL,
				WalletSeed:    req.WalletSeed,
				WalletAddr:    req.WalletAddr,
				Confirmations: req.Confirmations,
				Enabled:       req.Enabled,
			}
			// A blank seed on edit preserves the stored one (the UI never re-shows it),
			// so saving other fields can't accidentally wipe the escrow seed.
			if strings.TrimSpace(req.WalletSeed) == "" {
				if existing, err := database.GetSellCoin(req.Symbol); err == nil && existing != nil {
					sc.WalletSeed = existing.WalletSeed
				}
			}
			// A coin can be saved as a draft, but it can only be ENABLED once it has
			// everything needed to settle: fullnode URL, escrow seed and escrow
			// address (SKY defaults to the public fullnode).
			if sc.Enabled {
				node := strings.TrimSpace(sc.NodeURL)
				if node == "" && strings.EqualFold(sc.Symbol, "SKY") {
					node = "https://node.skycoin.com"
				}
				if !db.SellCoinComplete(node, sc.WalletSeed, sc.WalletAddr) {
					mWriteError(w, http.StatusBadRequest,
						"to enable a sell coin, set its fullnode URL, escrow wallet address and escrow seed (save it disabled to keep a draft)")
					return
				}
			}
			if err := database.UpsertSellCoin(sc); err != nil {
				mWriteError(w, http.StatusInternalServerError, err.Error())
				return
			}
			mWriteJSON(w, http.StatusOK, map[string]any{"ok": true})
		default:
			mWriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	})

	// POST /api/sellcoins/delete — remove a sell coin by {symbol}.
	mux.HandleFunc("/api/sellcoins/delete", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			mWriteError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		var req struct {
			Symbol string `json:"symbol"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.Symbol) == "" {
			mWriteError(w, http.StatusBadRequest, "symbol is required")
			return
		}
		if err := database.DeleteSellCoin(req.Symbol); err != nil {
			mWriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
		mWriteJSON(w, http.StatusOK, map[string]any{"ok": true})
	})

	// POST /api/sellcoins/enable — toggle availability via {symbol, enabled}.
	mux.HandleFunc("/api/sellcoins/enable", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			mWriteError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		var req struct {
			Symbol  string `json:"symbol"`
			Enabled bool   `json:"enabled"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.Symbol) == "" {
			mWriteError(w, http.StatusBadRequest, "symbol is required")
			return
		}
		// Enabling requires a complete escrow config (node + seed + address).
		if req.Enabled {
			node, seed, wallet, _, _, err := database.SellCoinConfig(req.Symbol)
			if err != nil {
				mWriteError(w, http.StatusInternalServerError, err.Error())
				return
			}
			if !db.SellCoinComplete(node, seed, wallet) {
				mWriteError(w, http.StatusBadRequest,
					"cannot enable "+strings.ToUpper(strings.TrimSpace(req.Symbol))+": set its fullnode URL, escrow wallet address and escrow seed first")
				return
			}
		}
		if err := database.SetSellCoinEnabled(req.Symbol, req.Enabled); err != nil {
			mWriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
		mWriteJSON(w, http.StatusOK, map[string]any{"ok": true})
	})

	// GET /api/bans — currently active (temporary) bans.
	mux.HandleFunc("/api/bans", func(w http.ResponseWriter, _ *http.Request) {
		bans, err := database.GetActiveBans()
		if err != nil {
			mWriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
		mWriteJSON(w, http.StatusOK, map[string]any{"bans": bans})
	})

	// POST /api/bans/unban — lift a temporary (freeze-violation) ban via {pubkey}.
	mux.HandleFunc("/api/bans/unban", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			mWriteError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		var req struct {
			PubKey string `json:"pubkey"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.PubKey) == "" {
			mWriteError(w, http.StatusBadRequest, "pubkey is required")
			return
		}
		if err := database.DeleteBan(strings.TrimSpace(req.PubKey)); err != nil {
			mWriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
		mWriteJSON(w, http.StatusOK, map[string]any{"ok": true})
	})

	// GET /api/blocks — buyers blocked from re-buying a product (buy-cancel limit
	// reached). Includes the current limit so the UI can explain the rule.
	mux.HandleFunc("/api/blocks", func(w http.ResponseWriter, _ *http.Request) {
		limit := database.GetBuyCancelLimit()
		blocks, err := database.ListBuyerProductBlocks(limit)
		if err != nil {
			mWriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
		mWriteJSON(w, http.StatusOK, map[string]any{"blocks": blocks, "limit": limit})
	})

	// POST /api/blocks/clear — clear a buyer's block for one product so they can
	// buy it again, via {pubkey, product_id}.
	mux.HandleFunc("/api/blocks/clear", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			mWriteError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		var req struct {
			PubKey    string `json:"pubkey"`
			ProductID string `json:"product_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil ||
			strings.TrimSpace(req.PubKey) == "" || strings.TrimSpace(req.ProductID) == "" {
			mWriteError(w, http.StatusBadRequest, "pubkey and product_id are required")
			return
		}
		cleared, err := database.ClearBuyerProductBlock(strings.TrimSpace(req.PubKey), strings.TrimSpace(req.ProductID))
		if err != nil {
			mWriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
		mWriteJSON(w, http.StatusOK, map[string]any{"ok": true, "cleared": cleared})
	})
}

func mWriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v) //nolint:errcheck
}

func mWriteError(w http.ResponseWriter, status int, msg string) {
	mWriteJSON(w, status, map[string]any{"error": msg})
}
