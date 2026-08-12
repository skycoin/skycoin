package liteclient

import (
	"encoding/hex"
	"encoding/json"
	"strings"
	"testing"

	"github.com/skycoin/skycoin/src/cipher"
)

// The seed the browser wallet derives from, hex-encoded the way
// convertAsciiToHexa does it in src/skycoin-web.
func testSeed() string {
	return hex.EncodeToString([]byte("abcdefg"))
}

func mustGenerate(t *testing.T, seed string) Address {
	t.Helper()

	address, err := GenerateAddress(seed)
	if err != nil {
		t.Fatalf("GenerateAddress(%q): %v", seed, err)
	}

	return address
}

func TestGenerateAddress(t *testing.T) {
	address := mustGenerate(t, testSeed())

	if address.Address != "gyJNKKj95bCn6o5mUCQgz8SCem7av2W3CG" {
		t.Errorf("address = %q", address.Address)
	}

	if address.Public != "02e1f33d00576ef4d89adcd1c8cb732810de862e8c5c6ecdc0208ccc6faefd1e09" {
		t.Errorf("public key = %q", address.Public)
	}

	next := mustGenerate(t, address.NextSeed)

	if next.Address != "2YzQDZFpS64u5ydnE9iKK9WJPXU2EbmSzmW" {
		t.Errorf("second address = %q", next.Address)
	}
}

// A seed that is not hex used to panic. Every caller in the browser wallet
// passes hex, but the entry point is reachable from JavaScript with anything.
func TestGenerateAddressRejectsNonHexSeed(t *testing.T) {
	if _, err := GenerateAddress("not a hex seed"); err == nil {
		t.Fatal("expected an error for a non-hex seed")
	}
}

func TestGenerateAddressesCount(t *testing.T) {
	addresses, err := GenerateAddresses(testSeed(), 3)
	if err != nil {
		t.Fatalf("GenerateAddresses: %v", err)
	}

	if len(addresses) != 3 {
		t.Fatalf("got %d addresses, want 3", len(addresses))
	}

	// Each is derived from the seed the previous one returned.
	seen := map[string]bool{}
	for _, address := range addresses {
		if seen[address.Address] {
			t.Errorf("address %q repeated", address.Address)
		}
		seen[address.Address] = true
	}
}

// buildInputs returns a spendable input owned by the given address entry.
func buildInputs(t *testing.T, address Address) string {
	t.Helper()

	body, err := json.Marshal([]TransactionInput{{
		Hash:   "0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f",
		Secret: address.Secret,
	}})
	if err != nil {
		t.Fatalf("marshalling inputs: %v", err)
	}

	return string(body)
}

func buildOutputs(t *testing.T, address string, coins, hours uint64) string {
	t.Helper()

	body, err := json.Marshal([]TransactionOutput{{Address: address, Coins: coins, Hours: hours}})
	if err != nil {
		t.Fatalf("marshalling outputs: %v", err)
	}

	return string(body)
}

func TestPrepareTransaction(t *testing.T) {
	address := mustGenerate(t, testSeed())

	txn, err := PrepareTransaction(buildInputs(t, address), buildOutputs(t, address.Address, 1000000, 100))
	if err != nil {
		t.Fatalf("PrepareTransaction: %v", err)
	}

	if txn == "" {
		t.Fatal("PrepareTransaction returned an empty transaction")
	}

	if _, err := hex.DecodeString(txn); err != nil {
		t.Errorf("transaction is not hex: %v", err)
	}
}

// The reachable one: the destination address is typed by the user, and a typo
// used to panic. On the standard build that surfaced as an error through
// recover; under TinyGo it trapped the whole wasm module.
func TestPrepareTransactionRejectsBadDestinationAddress(t *testing.T) {
	address := mustGenerate(t, testSeed())

	_, err := PrepareTransaction(buildInputs(t, address), buildOutputs(t, "not an address", 1000000, 100))
	if err == nil {
		t.Fatal("expected an error for an invalid destination address")
	}
}

func TestPrepareTransactionRejectsNullDestinationAddress(t *testing.T) {
	address := mustGenerate(t, testSeed())

	// A null address would burn the coins, so it is rejected before signing
	// rather than producing a spendable transaction. Asserting the specific
	// error matters: any invalid string fails the base58 decode first, which
	// would pass this test without ever reaching the check it is named for.
	_, err := PrepareTransaction(
		buildInputs(t, address), buildOutputs(t, cipher.Address{}.String(), 1000000, 100))
	if err != ErrNullOutputAddress {
		t.Fatalf("err = %v, want %v", err, ErrNullOutputAddress)
	}
}

func TestPrepareTransactionRejectsMalformedJSON(t *testing.T) {
	cases := []struct {
		name    string
		inputs  string
		outputs string
	}{
		{"inputs", "not json", "[]"},
		{"outputs", "[]", "not json"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if _, err := PrepareTransaction(c.inputs, c.outputs); err == nil {
				t.Fatalf("expected an error for malformed %s", c.name)
			}
		})
	}
}

// coin.Transaction.SignInputs reports this through log.Panic, which would take
// the process — or the wasm module — down rather than returning.
func TestPrepareTransactionRejectsNoInputs(t *testing.T) {
	address := mustGenerate(t, testSeed())

	_, err := PrepareTransaction("[]", buildOutputs(t, address.Address, 1000000, 100))
	if err != ErrNoInputs {
		t.Fatalf("err = %v, want %v", err, ErrNoInputs)
	}
}

// Well-formed hex of the right length, but not a valid secret key. This reaches
// cipher.MustSignHash, which panics.
func TestPrepareTransactionRejectsUnusableSecretKey(t *testing.T) {
	address := mustGenerate(t, testSeed())

	inputs, err := json.Marshal([]TransactionInput{{
		Hash:   "0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f",
		Secret: strings.Repeat("0", 64),
	}})
	if err != nil {
		t.Fatalf("marshalling inputs: %v", err)
	}

	if _, err := PrepareTransaction(string(inputs), buildOutputs(t, address.Address, 1000000, 100)); err == nil {
		t.Fatal("expected an error for an unusable secret key")
	}
}

func TestPrepareTransactionWithSignaturesRejectsMalformedSignatures(t *testing.T) {
	address := mustGenerate(t, testSeed())

	_, err := PrepareTransactionWithSignatures(
		buildInputs(t, address), buildOutputs(t, address.Address, 1000000, 100), `["not a signature"]`)
	if err == nil {
		t.Fatal("expected an error for a malformed signature")
	}

	_, err = PrepareTransactionWithSignatures(
		buildInputs(t, address), buildOutputs(t, address.Address, 1000000, 100), "not json")
	if err == nil {
		t.Fatal("expected an error for a malformed signature list")
	}
}
