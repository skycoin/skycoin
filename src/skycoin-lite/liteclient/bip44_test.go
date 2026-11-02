package liteclient

import (
	"encoding/json"
	"testing"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/wallet"
	"github.com/skycoin/skycoin/src/wallet/bip44wallet"
)

// A throwaway mnemonic, the bip39 all-zero entropy test vector.
const testMnemonic = "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"

// nodeWallet is the node's own bip44 wallet for a mnemonic, with n addresses
// on each chain. It is the oracle the browser derivation has to match.
func nodeWallet(t *testing.T, passphrase string, n int) (xpub string, external, change wallet.Entries) {
	t.Helper()

	w, err := bip44wallet.NewWallet("test.wlt", "test", testMnemonic, passphrase,
		wallet.OptionCoinType(wallet.CoinTypeSkycoin))
	if err != nil {
		t.Fatalf("NewWallet: %v", err)
	}

	for _, chain := range []wallet.Option{wallet.OptionExternal(), wallet.OptionChange()} {
		entries, err := w.GetEntries(chain)
		if err != nil {
			t.Fatalf("GetEntries: %v", err)
		}
		if missing := n - len(entries); missing > 0 {
			if _, err := w.GenerateAddresses(chain, wallet.OptionGenerateN(uint64(missing))); err != nil {
				t.Fatalf("GenerateAddresses: %v", err)
			}
		}
	}

	if xpub, err = w.GetXPubKey("0"); err != nil {
		t.Fatalf("GetXPubKey: %v", err)
	}
	if external, err = w.GetEntries(wallet.OptionExternal()); err != nil {
		t.Fatalf("GetEntries: %v", err)
	}
	if change, err = w.GetEntries(wallet.OptionChange()); err != nil {
		t.Fatalf("GetEntries: %v", err)
	}

	return xpub, external[:n], change[:n]
}

func checkAgainstEntries(t *testing.T, name string, got []Bip44Address, want wallet.Entries, withSecrets bool) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("%s: got %d addresses, want %d", name, len(got), len(want))
	}
	for i, a := range got {
		e := want[i]
		if a.Index != e.ChildNumber || a.Address != e.Address.String() || a.Public != e.Public.Hex() {
			t.Errorf("%s[%d] = %+v, want index %d address %s", name, i, a, e.ChildNumber, e.Address)
		}
		if withSecrets && a.Secret != e.Secret.Hex() {
			t.Errorf("%s[%d] secret does not match the node wallet", name, i)
		}
		if !withSecrets && a.Secret != "" {
			t.Errorf("%s[%d] has a secret from an xpub", name, i)
		}
	}
}

func TestBip44MatchesNodeWallet(t *testing.T) {
	const n = 5
	for _, passphrase := range []string{"", "a seed passphrase"} {
		xpub, external, change := nodeWallet(t, passphrase, n)

		got, err := Bip44AccountXPub(testMnemonic, passphrase)
		if err != nil {
			t.Fatalf("Bip44AccountXPub: %v", err)
		}
		if got != xpub {
			t.Errorf("passphrase %q: xpub = %s, want %s", passphrase, got, xpub)
		}

		for chain, want := range []wallet.Entries{external, change} {
			addrs, err := Bip44Addresses(testMnemonic, passphrase, chain, 0, n)
			if err != nil {
				t.Fatalf("Bip44Addresses: %v", err)
			}
			checkAgainstEntries(t, "Bip44Addresses", addrs, want, true)

			addrs, err = Bip44XPubAddresses(xpub, chain, 0, n)
			if err != nil {
				t.Fatalf("Bip44XPubAddresses: %v", err)
			}
			checkAgainstEntries(t, "Bip44XPubAddresses", addrs, want, false)

			// A window that does not start at zero.
			addrs, err = Bip44XPubAddresses(xpub, chain, 2, 2)
			if err != nil {
				t.Fatalf("Bip44XPubAddresses: %v", err)
			}
			checkAgainstEntries(t, "Bip44XPubAddresses window", addrs, want[2:4], false)
		}
	}
}

// The passphrase is part of the bip39 seed, so it changes every key.
func TestBip44PassphraseChangesXPub(t *testing.T) {
	a, err := Bip44AccountXPub(testMnemonic, "")
	if err != nil {
		t.Fatal(err)
	}
	b, err := Bip44AccountXPub(testMnemonic, "x")
	if err != nil {
		t.Fatal(err)
	}
	if a == b {
		t.Fatal("the passphrase did not change the xpub")
	}
}

func TestBip44RejectsBadInput(t *testing.T) {
	if _, err := Bip44AccountXPub("abandon abandon abandon", ""); err == nil {
		t.Error("expected an error for an invalid mnemonic")
	}
	if _, err := Bip44Addresses(testMnemonic, "", 2, 0, 1); err != ErrBip44Range {
		t.Errorf("chain 2: err = %v", err)
	}
	if _, err := Bip44Addresses(testMnemonic, "", 0, -1, 1); err != ErrBip44Range {
		t.Errorf("negative start: err = %v", err)
	}
	if _, err := Bip44XPubAddresses("xpub-not", 0, 0, 1); err == nil {
		t.Error("expected an error for an invalid xpub")
	}
	if _, err := Bip44XPubAddresses("", 0, 1<<30, 1<<30+1); err != ErrBip44Range {
		t.Errorf("hardened index: err = %v", err)
	}
}

// A send from bip44 addresses on both chains, with change back to a change
// address, signs with the derived secrets and every signature verifies against
// the address that owns the input.
func TestBip44PrepareTransaction(t *testing.T) {
	external, err := Bip44Addresses(testMnemonic, "pw", 0, 0, 2)
	if err != nil {
		t.Fatal(err)
	}
	change, err := Bip44Addresses(testMnemonic, "pw", 1, 0, 2)
	if err != nil {
		t.Fatal(err)
	}

	owners := []Bip44Address{external[1], change[0]}
	inputs := make([]TransactionInput, len(owners))
	for i, a := range owners {
		h := cipher.SumSHA256([]byte(a.Address))
		inputs[i] = TransactionInput{Hash: h.Hex(), Secret: a.Secret, Address: a.Address}
	}
	outputs := []TransactionOutput{
		{Address: "2YzQDZFpS64u5ydnE9iKK9WJPXU2EbmSzmW", Coins: 1000000, Hours: 1},
		{Address: change[1].Address, Coins: 2000000, Hours: 1},
	}

	inBody, err := json.Marshal(inputs)
	if err != nil {
		t.Fatal(err)
	}
	outBody, err := json.Marshal(outputs)
	if err != nil {
		t.Fatal(err)
	}

	raw, err := PrepareTransaction(string(inBody), string(outBody))
	if err != nil {
		t.Fatalf("PrepareTransaction: %v", err)
	}

	txn, err := coin.DeserializeTransactionHex(raw)
	if err != nil {
		t.Fatalf("DeserializeTransactionHex: %v", err)
	}
	if len(txn.Sigs) != len(owners) {
		t.Fatalf("got %d signatures, want %d", len(txn.Sigs), len(owners))
	}
	for i, a := range owners {
		addr := cipher.MustDecodeBase58Address(a.Address)
		if err := cipher.VerifyAddressSignedHash(addr, txn.Sigs[i], cipher.AddSHA256(txn.InnerHash, txn.In[i])); err != nil {
			t.Errorf("input %d signature does not verify for %s: %v", i, a.Address, err)
		}
	}
	if txn.Out[1].Address.String() != change[1].Address {
		t.Errorf("change output went to %s", txn.Out[1].Address)
	}
}
