package skydexui

import (
	"testing"
	"time"

	"github.com/skycoin/skycoin/src/skydex/protocol"
)

func TestMergeTrades(t *testing.T) {
	now := time.UnixMilli(1_000_000)
	saved := []Trade{
		{"id": "old", "status": "completed", "created_at": "2026-01-01T00:00:00Z", "market_pk": "m0", "saved_at": float64(5), "extra": "kept"},
		{"id": "both", "status": "expired", "created_at": "2026-02-01T00:00:00Z", "market_name": "Old name", "saved_at": float64(7)},
		{"status": "completed"},
	}
	orders := []protocol.OrderView{
		{ID: "both", Status: "completed", CreatedAt: "2026-02-01T00:00:00Z", Amount: 3},
		{ID: "new", Status: "canceled", CreatedAt: "2026-03-01T00:00:00Z"},
		{ID: "live", Status: "pending_payment", CreatedAt: "2026-04-01T00:00:00Z"},
	}
	got := mergeTrades(saved, orders, "m1", "", now)
	if len(got) != 3 {
		t.Fatalf("got %d trades, want 3: %v", len(got), got)
	}
	if got[0].str("id") != "new" || got[1].str("id") != "both" || got[2].str("id") != "old" {
		t.Fatalf("order = %s %s %s, want new both old", got[0].str("id"), got[1].str("id"), got[2].str("id"))
	}
	if got[0].str("market_pk") != "m1" || got[0].num("saved_at") != 0 && got[0]["saved_at"] != now.UnixMilli() {
		t.Errorf("new trade = %v", got[0])
	}
	b := got[1]
	if b.str("status") != "completed" || b.num("amount") != 3 || b.str("market_name") != "Old name" || b.num("saved_at") != 7 {
		t.Errorf("refreshed trade = %v", b)
	}
	if got[2].str("extra") != "kept" || got[2].str("market_pk") != "m0" {
		t.Errorf("purged trade = %v", got[2])
	}
}
