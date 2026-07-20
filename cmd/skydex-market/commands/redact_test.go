// Package commands cmd/apps/skydex-market/commands/redact_test.go
package commands

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestIsSecretConfigKey(t *testing.T) {
	secret := []string{
		"sky_wallet_seed",
		"explorer_btc_key",
		"explorer_ltc_key",
		"explorer_somefibercoin_key",
	}
	for _, k := range secret {
		if !isSecretConfigKey(k) {
			t.Fatalf("%q should be treated as secret", k)
		}
	}

	notSecret := []string{
		"market_name",
		"wallet_sky", // an address, not a key — safe and useful to display
		"sky_fullnode_url",
		"explorer_btc_url",
		"explorer_btc_provider",
		"commission_rate_percent",
	}
	for _, k := range notSecret {
		if isSecretConfigKey(k) {
			t.Fatalf("%q should not be treated as secret", k)
		}
	}
}

func TestRedactConfig_NeverEmitsSecretValues(t *testing.T) {
	const seed = "correct horse battery staple escrow seed"
	const apiKey = "sk_live_supersecretkey" //nolint

	cfg := map[string]string{
		"market_name":        "Test Market",
		"wallet_sky":         "2abcdefgHIJK",
		"sky_wallet_seed":    seed,
		"explorer_btc_key":   apiKey,
		"explorer_btc_url":   "https://mempool.space",
		"explorer_ltc_key":   "", // configured but empty
		"commission_min_sky": "0.001",
	}

	out := redactConfig(cfg)

	// The single most important assertion in this file: serialize exactly what
	// the HTTP handler would send and prove the secrets are not in the bytes.
	raw, err := json.Marshal(out)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	body := string(raw)
	if strings.Contains(body, seed) {
		t.Fatal("escrow seed leaked in the /api/config response body")
	}
	if strings.Contains(body, apiKey) {
		t.Fatal("explorer API key leaked in the /api/config response body")
	}

	values, ok := out["config"].(map[string]string)
	if !ok {
		t.Fatal("config block missing from redacted output")
	}
	if values["sky_wallet_seed"] != "" || values["explorer_btc_key"] != "" {
		t.Fatal("secret values should be blanked")
	}

	// Non-secret config must survive untouched, or the operator UI breaks.
	if values["market_name"] != "Test Market" {
		t.Fatalf("market_name = %q, want %q", values["market_name"], "Test Market")
	}
	if values["wallet_sky"] != "2abcdefgHIJK" {
		t.Fatal("wallet_sky (an address, not a secret) should be preserved")
	}
	if values["explorer_btc_url"] != "https://mempool.space" {
		t.Fatal("explorer url should be preserved")
	}

	// The UI still needs to know whether a secret is configured.
	set, ok := out["secrets_set"].(map[string]bool)
	if !ok {
		t.Fatal("secrets_set block missing")
	}
	if !set["sky_wallet_seed"] {
		t.Fatal("sky_wallet_seed is set and should report true")
	}
	if !set["explorer_btc_key"] {
		t.Fatal("explorer_btc_key is set and should report true")
	}
	if set["explorer_ltc_key"] {
		t.Fatal("empty explorer_ltc_key should report false")
	}
	if _, present := set["market_name"]; present {
		t.Fatal("non-secret keys should not appear in secrets_set")
	}
}

func TestRedactConfig_EmptyInput(t *testing.T) {
	out := redactConfig(map[string]string{})

	values, ok := out["config"].(map[string]string)
	if !ok || len(values) != 0 {
		t.Fatal("empty config should redact to an empty config block")
	}
	if _, ok := out["secrets_set"].(map[string]bool); !ok {
		t.Fatal("secrets_set should always be present")
	}
}
