package cipher

import (
	"bytes"
	"testing"
)

// Regression test: a compressed pubkey whose X coordinate is >= the field
// prime must be rejected as invalid without panicking (audit HIGH finding).
func TestPubKeyXCoordinateAboveFieldPrimeDoesNotPanic(t *testing.T) {
	b := append([]byte{0x02}, bytes.Repeat([]byte{0xFF}, 32)...)
	pk, err := NewPubKey(b)
	if err != nil {
		return // rejected at construction is also acceptable
	}
	if err := pk.Verify(); err == nil {
		t.Fatal("expected Verify() to reject X>=p pubkey, got nil error")
	}
}
