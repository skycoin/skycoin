package liteclient

import (
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/testutil"
	"github.com/skycoin/skycoin/src/wallet"
	"github.com/skycoin/skycoin/src/wallet/pqwallet"
)

// The browser derives the same post-quantum addresses as a node's pq wallet made
// from the same seed.
func TestGeneratePQAddressesMatchNodeWallet(t *testing.T) {
	w, err := pqwallet.NewWallet("t.wlt", "t", "abcdefg", wallet.OptionGenerateN(3))
	if err != nil {
		t.Fatal(err)
	}
	want, err := w.GetAddresses()
	if err != nil {
		t.Fatal(err)
	}

	got, err := GeneratePQAddresses(testSeed(), 0, 3)
	if err != nil {
		t.Fatal(err)
	}
	for i := range want {
		if got[i].Address != want[i].String() {
			t.Errorf("address %d: got %s, want %s", i, got[i].Address, want[i])
		}
	}

	one, err := GeneratePQAddresses(testSeed(), 2, 1)
	if err != nil {
		t.Fatal(err)
	}
	if one[0] != got[2] {
		t.Errorf("index 2 alone gave %v, want %v", one[0], got[2])
	}
}

func TestGeneratePQAddressesRejects(t *testing.T) {
	if _, err := GeneratePQAddresses("zz", 0, 1); err == nil {
		t.Error("non-hex seed was accepted")
	}
	if _, err := GeneratePQAddresses(testSeed(), -1, 1); err != ErrPQIndex {
		t.Errorf("negative index: got %v", err)
	}
}

// A transaction that spends a secp256k1 output and a post-quantum one verifies
// after the round trip through its raw form.
func TestPrepareTransactionPQ(t *testing.T) {
	secp := mustGenerate(t, testSeed())
	pq, err := GeneratePQAddresses(testSeed(), 0, 1)
	if err != nil {
		t.Fatal(err)
	}
	pqAddr := cipher.MustDecodeBase58Address(pq[0].Address)

	uxs := coin.UxArray{
		{Body: coin.UxBody{SrcTransaction: testutil.RandSHA256(t), Address: cipher.MustDecodeBase58Address(secp.Address), Coins: 2e6, Hours: 10}},
		{Body: coin.UxBody{SrcTransaction: testutil.RandSHA256(t), Address: pqAddr, Coins: 2e6, Hours: 10}},
	}
	inputs, err := json.Marshal([]TransactionInput{
		{Hash: uxs[0].Hash().Hex(), Secret: secp.Secret},
		{Hash: uxs[1].Hash().Hex(), Secret: pq[0].Secret, Address: pq[0].Address},
	})
	if err != nil {
		t.Fatal(err)
	}
	outputs, err := json.Marshal([]TransactionOutput{{Address: secp.Address, Coins: 4e6, Hours: 5}})
	if err != nil {
		t.Fatal(err)
	}

	raw, err := PrepareTransaction(string(inputs), string(outputs))
	if err != nil {
		t.Fatalf("PrepareTransaction: %v", err)
	}
	b, err := hex.DecodeString(raw)
	if err != nil {
		t.Fatal(err)
	}
	txn, err := coin.DeserializeTransactionWithPQ(b)
	if err != nil {
		t.Fatal(err)
	}
	if txn.Type != coin.TxnTypePQ {
		t.Errorf("type %d, want %d", txn.Type, coin.TxnTypePQ)
	}
	if err := txn.Verify(); err != nil {
		t.Errorf("Verify: %v", err)
	}
	if err := txn.VerifyInputSignatures(uxs); err != nil {
		t.Errorf("VerifyInputSignatures: %v", err)
	}
}

func TestPrepareTransactionRejectsBadPQSecret(t *testing.T) {
	pq, err := GeneratePQAddresses(testSeed(), 0, 1)
	if err != nil {
		t.Fatal(err)
	}
	inputs, err := json.Marshal([]TransactionInput{{Hash: testutil.RandSHA256(t).Hex(), Secret: "abcd", Address: pq[0].Address}})
	if err != nil {
		t.Fatal(err)
	}
	outputs, err := json.Marshal([]TransactionOutput{{Address: pq[0].Address, Coins: 1e6}})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := PrepareTransaction(string(inputs), string(outputs)); err != ErrPQSecretKey {
		t.Errorf("got %v, want %v", err, ErrPQSecretKey)
	}
}
