package liteclient

import "testing"

// One consistent set of key material, the same vector the GopherJS layer's test
// has always used.
const (
	testPubKey    = "037c4cff096a7219b17f8502b9ed643c947d5d4929c1a141b3240f70b60a15a7b8"
	testSecKey    = "697c7cfba3c6d13dc6bd3f063c60ef4d25de903e50fe8c5e123e5efb08e21e29"
	testAddress   = "XZ9S3QKN5tSRVswDNE6GLCtTfm8DqRthyA"
	testSignature = "ba85034f675b8a284537a0c23e4d55d5f8b4cc1620eee077a478482b8f7bf9304d17c16dd5aaf70d24df262c531f37e58a57469d1161d6134075ed8c203f7cbc01"
	testHash      = "72cd6e8422c407fb6d098690f1130b7ded7ec2f7f5e1d30bd9d521f015363793"
)

func TestVerifyHelpersAcceptValidInput(t *testing.T) {
	if err := VerifyPubKeySignedHash(testPubKey, testSignature, testHash); err != nil {
		t.Errorf("VerifyPubKeySignedHash: %v", err)
	}

	if err := VerifyAddressSignedHash(testAddress, testSignature, testHash); err != nil {
		t.Errorf("VerifyAddressSignedHash: %v", err)
	}

	if err := VerifySignatureRecoverPubKey(testSignature, testHash); err != nil {
		t.Errorf("VerifySignatureRecoverPubKey: %v", err)
	}

	if err := VerifySeckey(testSecKey); err != nil {
		t.Errorf("VerifySeckey: %v", err)
	}

	if err := VerifyPubkey(testPubKey); err != nil {
		t.Errorf("VerifyPubkey: %v", err)
	}
}

func TestAddressDerivationAgrees(t *testing.T) {
	fromPubKey, err := AddressFromPubKey(testPubKey)
	if err != nil {
		t.Fatalf("AddressFromPubKey: %v", err)
	}

	fromSecKey, err := AddressFromSecKey(testSecKey)
	if err != nil {
		t.Fatalf("AddressFromSecKey: %v", err)
	}

	if fromPubKey != testAddress || fromSecKey != testAddress {
		t.Errorf("derived %q and %q, want %q", fromPubKey, fromSecKey, testAddress)
	}
}

func TestSignAndRecover(t *testing.T) {
	recovered, err := PubKeyFromSig(testSignature, testHash)
	if err != nil {
		t.Fatalf("PubKeyFromSig: %v", err)
	}

	if recovered != testPubKey {
		t.Errorf("recovered %q, want %q", recovered, testPubKey)
	}

	signature, err := SignHash(testHash, testSecKey)
	if err != nil {
		t.Fatalf("SignHash: %v", err)
	}

	if signature == "" {
		t.Error("SignHash returned an empty signature")
	}
}

// Every one of these used to panic — through an explicit panic(err) or through
// a cipher Must* constructor — which the wasm build could only turn back into an
// error on the standard toolchain.
func TestVerifyHelpersRejectMalformedInput(t *testing.T) {
	cases := []struct {
		name string
		err  error
	}{
		{"VerifyPubKeySignedHash/pubkey", VerifyPubKeySignedHash("nonsense", testSignature, testHash)},
		{"VerifyPubKeySignedHash/sig", VerifyPubKeySignedHash(testPubKey, "nonsense", testHash)},
		{"VerifyPubKeySignedHash/hash", VerifyPubKeySignedHash(testPubKey, testSignature, "nonsense")},
		{"VerifyAddressSignedHash/address", VerifyAddressSignedHash("nonsense", testSignature, testHash)},
		{"VerifySignatureRecoverPubKey/sig", VerifySignatureRecoverPubKey("nonsense", testHash)},
		{"VerifySeckey", VerifySeckey("nonsense")},
		{"VerifyPubkey", VerifyPubkey("nonsense")},
	}

	for _, c := range cases {
		if c.err == nil {
			t.Errorf("%s: expected an error", c.name)
		}
	}
}

func TestDerivationRejectsMalformedInput(t *testing.T) {
	if _, err := AddressFromPubKey("nonsense"); err == nil {
		t.Error("AddressFromPubKey: expected an error")
	}

	if _, err := AddressFromSecKey("nonsense"); err == nil {
		t.Error("AddressFromSecKey: expected an error")
	}

	if _, err := PubKeyFromSig("nonsense", testHash); err == nil {
		t.Error("PubKeyFromSig: expected an error")
	}

	if _, err := SignHash(testHash, "nonsense"); err == nil {
		t.Error("SignHash: expected an error")
	}
}

// A secret key that is well-formed hex of the right length but not a usable
// key. cipher.MustSignHash panics on it.
func TestSignHashRejectsUnusableSecretKey(t *testing.T) {
	zero := "0000000000000000000000000000000000000000000000000000000000000000"

	if err := VerifySeckey(zero); err == nil {
		t.Error("VerifySeckey: expected an error for the zero key")
	}

	if _, err := SignHash(testHash, zero); err == nil {
		t.Error("SignHash: expected an error for the zero key")
	}
}
