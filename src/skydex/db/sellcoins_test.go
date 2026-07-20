package db

import (
	"path/filepath"
	"testing"
)

func newSellCoinDB(t *testing.T) *Database {
	t.Helper()
	d, err := New(filepath.Join(t.TempDir(), "sc.db"), "")
	if err != nil {
		t.Fatalf("db.New: %v", err)
	}
	t.Cleanup(func() { _ = d.Close() }) //nolint:errcheck
	if err := d.Migrate(); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	if err := d.InitDefaultConfig(); err != nil {
		t.Fatalf("init config: %v", err)
	}
	return d
}

// TestSellCoinDefaultSKY: InitDefaultConfig seeds a DISABLED SKY row pre-filled
// with the public fullnode; it is unavailable until the operator adds the escrow
// seed + address and enables it.
func TestSellCoinDefaultSKY(t *testing.T) {
	d := newSellCoinDB(t)
	sky, err := d.GetSellCoin("sky") // case-insensitive
	if err != nil {
		t.Fatalf("GetSellCoin: %v", err)
	}
	if sky == nil || sky.Symbol != "SKY" {
		t.Fatalf("default SKY sell coin missing: %+v", sky)
	}
	if sky.Enabled {
		t.Fatalf("default SKY row should be disabled until configured: %+v", sky)
	}
	// The effective node defaults to the public fullnode even though the row is
	// empty (so a pre-fibercoin sky_fullnode_url can still override it).
	if node, _, _, _, _, _ := d.SellCoinConfig("SKY"); node != "https://node.skycoin.com" { //nolint:errcheck
		t.Fatalf("default SKY node = %q, want the public fullnode", node)
	}
	// Not available yet: disabled and no escrow seed/address.
	if avail, err := d.AvailableSellCoins(); err != nil || len(avail) != 0 {
		t.Fatalf("AvailableSellCoins = %v, %v; want none until SKY is configured", avail, err)
	}
	// Configure escrow + enable -> now available.
	if err := d.SetConfig("sky_wallet_seed", "seed"); err != nil {
		t.Fatal(err)
	}
	if err := d.SetConfig("wallet_sky", "SKYaddr"); err != nil {
		t.Fatal(err)
	}
	if err := d.SetSellCoinEnabled("SKY", true); err != nil {
		t.Fatal(err)
	}
	if avail, _ := d.AvailableSellCoins(); len(avail) != 1 || avail[0] != "SKY" { //nolint:errcheck
		t.Fatalf("AvailableSellCoins = %v, want [SKY]", avail)
	}
}

// TestSellCoinLegacyBridge: SKY's escrow config falls back to the legacy
// sky_fullnode_url/sky_wallet_seed/wallet_sky keys when its row fields are empty.
func TestSellCoinLegacyBridge(t *testing.T) {
	d := newSellCoinDB(t)
	for k, v := range map[string]string{
		"sky_fullnode_url": "http://node:6420",
		"sky_wallet_seed":  "seed words",
		"wallet_sky":       "SKYescrowaddr",
	} {
		if err := d.SetConfig(k, v); err != nil {
			t.Fatalf("set %s: %v", k, err)
		}
	}
	// The seeded SKY row is disabled by default; enable it (now complete via the
	// legacy keys) so the bridge reports it enabled.
	if err := d.SetSellCoinEnabled("SKY", true); err != nil {
		t.Fatal(err)
	}
	url, seed, wallet, confs, enabled, err := d.SellCoinConfig("SKY")
	if err != nil {
		t.Fatalf("SellCoinConfig: %v", err)
	}
	if url != "http://node:6420" || seed != "seed words" || wallet != "SKYescrowaddr" || !enabled || confs < 1 {
		t.Fatalf("bridge = (%q,%q,%q,%d,%v)", url, seed, wallet, confs, enabled)
	}
	if w, _ := d.SellCoinWallet("SKY"); w != "SKYescrowaddr" { //nolint:errcheck
		t.Fatalf("SellCoinWallet(SKY) = %q, want SKYescrowaddr", w)
	}
}

// TestSellCoinCRUD covers add/edit/enable/disable/delete of a fibercoin.
func TestSellCoinCRUD(t *testing.T) {
	d := newSellCoinDB(t)

	// Configure + enable the seeded SKY row so it counts as available alongside
	// the fibercoin added below.
	if err := d.SetConfig("sky_wallet_seed", "sky seed"); err != nil {
		t.Fatal(err)
	}
	if err := d.SetConfig("wallet_sky", "SKYaddr"); err != nil {
		t.Fatal(err)
	}
	if err := d.SetSellCoinEnabled("SKY", true); err != nil {
		t.Fatal(err)
	}

	// Add a fibercoin (lower-case symbol is normalized to upper-case).
	if err := d.UpsertSellCoin(&SellCoin{
		Symbol: "mdl", Name: "Mobile", NodeURL: "http://mdl:6420",
		WalletSeed: "mdl seed", WalletAddr: "MDLaddr", Confirmations: 2, Enabled: true,
	}); err != nil {
		t.Fatalf("upsert: %v", err)
	}
	mdl, err := d.GetSellCoin("MDL")
	if err != nil || mdl == nil {
		t.Fatalf("GetSellCoin(MDL): %v / %v", mdl, err)
	}
	if mdl.NodeURL != "http://mdl:6420" || mdl.Confirmations != 2 || !mdl.Enabled {
		t.Fatalf("unexpected MDL row: %+v", mdl)
	}
	// Its config resolves purely from the row (no SKY legacy bridge).
	url, seed, wallet, confs, enabled, err := d.SellCoinConfig("MDL")
	if err != nil || url != "http://mdl:6420" || seed != "mdl seed" || wallet != "MDLaddr" || confs != 2 || !enabled {
		t.Fatalf("MDL config = (%q,%q,%q,%d,%v) err=%v", url, seed, wallet, confs, enabled, err)
	}

	// Now SKY + MDL are both available.
	avail, _ := d.AvailableSellCoins() //nolint:errcheck
	if len(avail) != 2 || avail[0] != "SKY" || avail[1] != "MDL" {
		t.Fatalf("AvailableSellCoins = %v, want [SKY MDL]", avail)
	}

	// Disable MDL: still configured, no longer available.
	if err := d.SetSellCoinEnabled("MDL", false); err != nil {
		t.Fatalf("disable: %v", err)
	}
	if ok, _ := d.IsSellCoinAvailable("MDL"); ok { //nolint:errcheck
		t.Fatal("MDL should be unavailable while disabled")
	}

	// Edit with an empty seed via UpsertSellCoin does NOT preserve the seed at the
	// DB layer (that merge lives in the operator API); assert the raw behavior.
	if err := d.UpsertSellCoin(&SellCoin{Symbol: "MDL", Name: "Mobile2", NodeURL: "http://mdl2:6420", WalletAddr: "MDLaddr", Confirmations: 3, Enabled: true}); err != nil {
		t.Fatalf("edit: %v", err)
	}
	mdl2, _ := d.GetSellCoin("MDL") //nolint:errcheck
	if mdl2.Name != "Mobile2" || mdl2.NodeURL != "http://mdl2:6420" || mdl2.Confirmations != 3 {
		t.Fatalf("edit not applied: %+v", mdl2)
	}

	// Delete it.
	if err := d.DeleteSellCoin("MDL"); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if got, _ := d.GetSellCoin("MDL"); got != nil { //nolint:errcheck
		t.Fatalf("MDL still present after delete: %+v", got)
	}
	// SKY remains.
	if ok, _ := d.IsSellCoinAvailable("SKY"); !ok { //nolint:errcheck
		t.Fatal("SKY should still be available after deleting MDL")
	}
}

// TestUnknownSellCoinUnavailable: an unconfigured symbol is never available.
func TestUnknownSellCoinUnavailable(t *testing.T) {
	d := newSellCoinDB(t)
	if ok, err := d.IsSellCoinAvailable("NOPE"); err != nil || ok {
		t.Fatalf("IsSellCoinAvailable(NOPE) = %v, %v; want false, nil", ok, err)
	}
	if _, _, _, _, enabled, _ := d.SellCoinConfig("NOPE"); enabled { //nolint:errcheck
		t.Fatal("unknown coin should report disabled")
	}
}
