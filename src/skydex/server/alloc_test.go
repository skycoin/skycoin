package server

import (
	"errors"
	"math"
	"testing"
)

// TestRoundToDecimals normalizes to an asset's on-chain precision. SKY rounds to
// 3 decimals (the network's spendable precision); payment coins to 8.
func TestRoundToDecimals(t *testing.T) {
	cases := []struct {
		in   float64
		d    int
		want float64
	}{
		{10.123456789, skyOnChainDecimals, 10.123}, //nolint // SKY: 3 dp
		{2.123456789, paymentDecimals, 2.12345679}, // payment coin: 8 dp
		{0.0004, skyOnChainDecimals, 0},            // below half a milli-SKY -> 0
		{10.0, skyOnChainDecimals, 10.0},
	}
	for _, c := range cases {
		if got := roundToDecimals(c.in, c.d); math.Abs(got-c.want) > 1e-12 {
			t.Errorf("roundToDecimals(%v, %d) = %v, want %v", c.in, c.d, got, c.want)
		}
	}
}

// The non-round unique-amount mechanism below is used only for external payment
// coins (BTC/LTC, 8 dp). SKY deposits use exact round amounts identified by
// sender address + time window, so they don't go through nonRound.

// TestNonRoundDeltaAtCap verifies the delta survives even at the max allowed
// price, so uniqueness isn't silently lost (the draw is always above the base).
func TestNonRoundDeltaAtCap(t *testing.T) {
	for i := 0; i < 500; i++ {
		v := nonRound(maxPrice, paymentDecimals)
		if v <= maxPrice {
			t.Fatalf("delta was lost at the cap: %v <= %v", v, maxPrice)
		}
	}
}

// TestAllocUniqueAmountRetries verifies the allocator retries past collisions
// and returns a non-colliding amount.
func TestAllocUniqueAmountRetries(t *testing.T) {
	calls := 0
	// Report the first 3 generated amounts as taken, then free.
	exists := func(float64) (bool, error) {
		calls++
		return calls <= 3, nil
	}
	amt, err := allocUniqueAmount(10, paymentDecimals, exists)
	if err != nil {
		t.Fatalf("allocUniqueAmount: %v", err)
	}
	if amt <= 10 {
		t.Fatalf("amount %v should be a non-round value above the base", amt)
	}
	if calls != 4 {
		t.Fatalf("expected 4 attempts (3 collisions + 1 free), got %d", calls)
	}
}

// TestAllocUniqueAmountExhausts returns an error when every amount collides.
func TestAllocUniqueAmountExhausts(t *testing.T) {
	exists := func(float64) (bool, error) { return true, nil }
	if _, err := allocUniqueAmount(10, paymentDecimals, exists); !errors.Is(err, errNoUniqueAmount) {
		t.Fatalf("expected errNoUniqueAmount, got %v", err)
	}
}

// TestAllocUniqueAmountPropagatesError surfaces a DB error from the check.
func TestAllocUniqueAmountPropagatesError(t *testing.T) {
	boom := errors.New("db down")
	exists := func(float64) (bool, error) { return false, boom }
	if _, err := allocUniqueAmount(10, paymentDecimals, exists); !errors.Is(err, boom) {
		t.Fatalf("expected the db error, got %v", err)
	}
}

// TestNonRoundPrecision confirms payment amounts are aligned to 8 decimals and
// the delta stays in the intended small range.
func TestNonRoundPrecision(t *testing.T) {
	for i := 0; i < 1000; i++ {
		v := nonRound(10, paymentDecimals)
		// v must already be an 8-decimal value (rounding it changes nothing).
		if r := roundToDecimals(v, paymentDecimals); math.Abs(v-r) > 1e-12 {
			t.Fatalf("payment amount %v is not aligned to %d decimals", v, paymentDecimals)
		}
		if v <= 10 || v > 10.00001 {
			t.Fatalf("delta out of expected small range: %v", v)
		}
	}
}
