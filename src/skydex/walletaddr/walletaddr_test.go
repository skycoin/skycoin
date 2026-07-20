package walletaddr

import (
	"testing"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/base58"
	"github.com/skycoin/skycoin/src/cipher/bech32"
)

// b58check builds a valid base58check address (version + 20-byte payload +
// double-SHA256 checksum), so tests use authoritative, self-consistent vectors
// instead of hard-coded strings.
func b58check(version byte, payload20 []byte) string {
	b := append([]byte{version}, payload20...)
	sum := cipher.DoubleSHA256(b)
	b = append(b, sum[:4]...)
	return base58.Encode(b)
}

// segwit builds a valid v0 P2WPKH address for the given hrp.
func segwit(t *testing.T, hrp string, program20 []byte) string {
	t.Helper()
	a, err := bech32.SegwitEncode(hrp, 0, program20)
	if err != nil {
		t.Fatalf("SegwitEncode(%s): %v", hrp, err)
	}
	return a
}

func TestValidateValid(t *testing.T) {
	h20 := make([]byte, 20)
	for i := range h20 {
		h20[i] = byte(i + 1)
	}
	pk, _ := cipher.GenerateKeyPair()
	skyAddr := cipher.AddressFromPubKey(pk).String()

	cases := []struct{ currency, addr string }{
		{"SKY", skyAddr},
		{"BTC", b58check(0x00, h20)},               // P2PKH "1..."
		{"BTC", b58check(0x05, h20)},               // P2SH  "3..."
		{"BTC", segwit(t, "bc", h20)},              // bech32 "bc1..."
		{"LTC", b58check(0x30, h20)},               // P2PKH "L..."
		{"LTC", b58check(0x32, h20)},               // P2SH  "M..."
		{"LTC", segwit(t, "ltc", h20)},             // bech32 "ltc1..."
		{"DOGE", "anything-goes-no-validator-yet"}, // unsupported => accept
	}
	for _, c := range cases {
		if err := Validate(c.currency, c.addr); err != nil {
			t.Errorf("Validate(%s, %q) = %v, want nil", c.currency, c.addr, err)
		}
	}
}

func TestValidateInvalid(t *testing.T) {
	h20 := make([]byte, 20)

	// A base58check BTC address with a tampered final char (broken checksum).
	goodBTC := b58check(0x00, h20)
	badBTC := goodBTC[:len(goodBTC)-1] + flip(goodBTC[len(goodBTC)-1])

	// A valid LTC segwit address — its hrp is "ltc", so decoding under BTC fails.
	ltcSegwit := segwit(t, "ltc", h20)

	cases := []struct{ currency, addr, why string }{
		{"SKY", "", "empty"},
		{"SKY", "not-an-address", "garbage"},
		{"SKY", b58check(0x00, h20), "BTC address is not a SKY address"},
		{"BTC", badBTC, "broken checksum"},
		{"BTC", b58check(0x30, h20), "LTC version byte for BTC"},
		{"BTC", ltcSegwit, "ltc hrp under BTC"},
		{"BTC", "bc1qinvalid!!!", "malformed bech32"},
		{"LTC", b58check(0x00, h20), "BTC P2PKH version for LTC"},
	}
	for _, c := range cases {
		if err := Validate(c.currency, c.addr); err == nil {
			t.Errorf("Validate(%s, %q) = nil, want error (%s)", c.currency, c.addr, c.why)
		}
	}
}

// flip returns a different character than c (to break a checksum deterministically).
func flip(c byte) string {
	if c == 'a' {
		return "b"
	}
	return "a"
}
