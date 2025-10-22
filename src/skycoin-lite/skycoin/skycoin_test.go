package main

import (
	"encoding/hex"
	"testing"

	"github.com/skycoin/skycoin/src/skycoin-lite/liteclient"
)

func TestGenerateAddress(t *testing.T) {

	//client.go

	stringSeed := "abcdefg"
	seed := hex.EncodeToString([]byte(stringSeed))

	addr := liteclient.GenerateAddress(seed)

	if addr.Address != "gyJNKKj95bCn6o5mUCQgz8SCem7av2W3CG" {
		t.Fatalf("GenerateAddress address is invalid")
	}

	if addr.Public != "02e1f33d00576ef4d89adcd1c8cb732810de862e8c5c6ecdc0208ccc6faefd1e09" {
		t.Fatalf("GenerateAddress pubkey is invalid")
	}

	nextAddr := liteclient.GenerateAddress(addr.NextSeed)

	if nextAddr.Address != "2YzQDZFpS64u5ydnE9iKK9WJPXU2EbmSzmW" {
		t.Fatalf("GenerateAddress address is invalid 2")
	}

	if nextAddr.Public != "023ea091a649aef9d8fe5b956bd5cf2a9e47b8a7ba41f8b83997de0d1b1a851e73" {
		t.Fatalf("GenerateAddress pubkey is invalid 2")
	}

	//extras.go

	pubkey := "037c4cff096a7219b17f8502b9ed643c947d5d4929c1a141b3240f70b60a15a7b8"
	seckey := "697c7cfba3c6d13dc6bd3f063c60ef4d25de903e50fe8c5e123e5efb08e21e29"
	address := "XZ9S3QKN5tSRVswDNE6GLCtTfm8DqRthyA"
	signature := "ba85034f675b8a284537a0c23e4d55d5f8b4cc1620eee077a478482b8f7bf9304d17c16dd5aaf70d24df262c531f37e58a57469d1161d6134075ed8c203f7cbc01"
	hash := "72cd6e8422c407fb6d098690f1130b7ded7ec2f7f5e1d30bd9d521f015363793"

	liteclient.VerifyPubKeySignedHash(pubkey, signature, hash)

	liteclient.VerifyAddressSignedHash(address, signature, hash)

	liteclient.VerifySignatureRecoverPubKey(signature, hash)

	if liteclient.VerifySeckey(seckey) != 1 {
		t.Fatalf("secp256k1.VerifySeckey failed")
	}

	if liteclient.VerifyPubkey(pubkey) != 1 {
		t.Fatalf("secp256k1.VerifyPubkey failed")
	}

	if liteclient.AddressFromPubKey(pubkey) != address {
		t.Fatalf("cipher.AddressFromPubKey failed")
	}

	if liteclient.AddressFromSecKey(seckey) != address {
		t.Fatalf("cipher.AddressFromPubKey failed")
	}

	pubkey2 := liteclient.PubKeyFromSig(signature, hash)
	if pubkey2 != pubkey {
		t.Fatalf("cipher.AddressFromPubKey failed")
	}

	signHash := liteclient.SignHash(hash, seckey)
	if signHash == "" {
		t.Fatalf("created signature is null")
	}

}
