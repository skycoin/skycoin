// Package walletaddr validates user-supplied cryptocurrency payout addresses by
// currency. The exchange never learns an address is wrong until it tries to pay
// out to it; validating at registration turns a silent typo into an immediate,
// actionable error at the boundary.
//
// Phase 1 trades SKY (sold) for BTC or LTC (payment), so those three are
// validated. Any other currency returns nil (accepted) until it becomes a
// supported payment coin and gets its own validator here.
//
// Validation is format/checksum-level (base58check version+checksum, bech32/
// bech32m segwit, SKY base58check) — it proves an address is well-formed and
// self-consistent, not that it exists or is controlled by the user.
package walletaddr

import (
	"fmt"
	"slices"
	"strings"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/base58"
	"github.com/skycoin/skycoin/src/cipher/bech32"
)

// Validate checks addr against the address format(s) valid for currency (an
// uppercase code: SKY, BTC, LTC). It returns a descriptive error for a malformed
// address and nil for a well-formed one, or nil for a currency without a
// validator yet.
func Validate(currency, addr string) error {
	addr = strings.TrimSpace(addr)
	if addr == "" {
		return fmt.Errorf("empty address")
	}
	switch currency {
	case "SKY":
		if _, err := cipher.DecodeBase58Address(addr); err != nil {
			return fmt.Errorf("invalid SKY address: %w", err)
		}
		return nil
	case "BTC":
		// P2PKH (0x00, "1"), P2SH (0x05, "3"), or segwit "bc1".
		return validateBase58OrSegwit(addr, []byte{0x00, 0x05}, "bc")
	case "LTC":
		// P2PKH (0x30, "L"), P2SH (0x32, "M"; legacy 0x05, "3"), or segwit "ltc1".
		return validateBase58OrSegwit(addr, []byte{0x30, 0x32, 0x05}, "ltc")
	default:
		return nil
	}
}

// validateBase58OrSegwit accepts a legacy base58check address with one of the
// given version bytes, or a mainnet bech32/bech32m segwit address under hrp.
func validateBase58OrSegwit(addr string, versions []byte, hrp string) error {
	if strings.HasPrefix(strings.ToLower(addr), hrp+"1") {
		if _, _, err := bech32.SegwitDecode(hrp, addr); err != nil {
			return fmt.Errorf("invalid %s segwit address: %w", strings.ToUpper(hrp), err)
		}
		return nil
	}
	if err := base58Check(addr, versions); err != nil {
		return fmt.Errorf("invalid address: %w", err)
	}
	return nil
}

// base58Check decodes a base58check address (1 version byte + 20-byte payload +
// 4-byte double-SHA256 checksum) and verifies its length, checksum, and that the
// version byte is one of the allowed set.
func base58Check(addr string, versions []byte) error {
	b, err := base58.Decode(addr)
	if err != nil {
		return fmt.Errorf("not valid base58: %w", err)
	}
	if len(b) != 25 {
		return fmt.Errorf("wrong decoded length %d (want 25)", len(b))
	}
	sum := cipher.DoubleSHA256(b[:21])
	if !equal4(b[21:25], sum[:4]) {
		return fmt.Errorf("bad checksum")
	}
	if !slices.Contains(versions, b[0]) {
		return fmt.Errorf("unexpected version byte 0x%02x", b[0])
	}
	return nil
}

func equal4(a, b []byte) bool {
	return len(a) == 4 && len(b) == 4 &&
		a[0] == b[0] && a[1] == b[1] && a[2] == b[2] && a[3] == b[3]
}
