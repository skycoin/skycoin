package skydexui

import (
	"testing"
	"time"
)

func TestBuildURI(t *testing.T) {
	for _, c := range []struct {
		amount float64
		coin   string
		want   string
	}{
		{0.00020038, "BTC", "bitcoin:A?amount=0.00020038"},
		{1e-7, "btc", "bitcoin:A?amount=0.0000001"},
		{10, "SKY", "skycoin:A?amount=10"},
		{2.5, "LTC", "litecoin:A?amount=2.5"},
		{1, "ETH", ""},
		{0, "BTC", ""},
		{-1, "BTC", ""},
	} {
		if got := buildURI("A", c.amount, c.coin); got != c.want {
			t.Errorf("buildURI(%v %s) = %q, want %q", c.amount, c.coin, got, c.want)
		}
	}
	if got := buildURI("", 1, "BTC"); got != "" {
		t.Errorf("no address gives %q", got)
	}
}

func TestTimeAgo(t *testing.T) {
	now := time.UnixMilli(10_000_000_000)
	for ago, want := range map[time.Duration]string{
		10 * time.Second: "just now",
		5 * time.Minute:  "5m ago",
		3 * time.Hour:    "3h ago",
		50 * time.Hour:   "2d ago",
	} {
		if got := timeAgo(now.Add(-ago).UnixMilli(), now); got != want {
			t.Errorf("timeAgo(%v) = %q, want %q", ago, got, want)
		}
	}
	if timeAgo(0, now) != "" {
		t.Error("no time renders")
	}
}

func TestShortKey(t *testing.T) {
	if got := shortKey("0123456789abcdefghij", 8, 4); got != "01234567…ghij" {
		t.Errorf("shortKey = %q", got)
	}
	if got := shortKey("short", 8, 4); got != "short" {
		t.Errorf("shortKey of a short key = %q", got)
	}
}
