// Package main provides the GopherJS skycoin-lite browser implementation
//
// The functions published here mirror the wasm build in
// src/skycoin-lite/wasm/main_wasm.go: on failure they return {error: string}
// rather than throwing. liteclient returns errors now, and letting them surface
// as JavaScript exceptions instead would leave the two builds of the same
// library disagreeing about how a caller is meant to check for failure.
package main

import (
	"github.com/gopherjs/gopherjs/js"

	"github.com/skycoin/skycoin/src/skycoin-lite/liteclient"
)

// errorResult is what every entry point returns when it could not do its job.
func errorResult(err error) map[string]interface{} {
	return map[string]interface{}{"error": err.Error()}
}

// result returns value, or an error object if err is set.
func result(value interface{}, err error) interface{} {
	if err != nil {
		return errorResult(err)
	}

	return value
}

// checked returns null when err is nil, mirroring the verify* helpers of the
// wasm build.
func checked(err error) interface{} {
	if err != nil {
		return errorResult(err)
	}

	return nil
}

func main() {
	js.Global.Set("Cipher", map[string]interface{}{
		"GenerateAddresses": func(seed string, num int) interface{} {
			return result(liteclient.GenerateAddresses(seed, num))
		},
		"PrepareTransaction": func(inputsBody, outputsBody string) interface{} {
			return result(liteclient.PrepareTransaction(inputsBody, outputsBody))
		},
		"PrepareTransactionWithSignatures": func(inputsBody, outputsBody, signatureList string) interface{} {
			return result(liteclient.PrepareTransactionWithSignatures(inputsBody, outputsBody, signatureList))
		},
	})

	js.Global.Set("CipherExtras", map[string]interface{}{
		"VerifyPubKeySignedHash": func(pubkey, sig, hash string) interface{} {
			return checked(liteclient.VerifyPubKeySignedHash(pubkey, sig, hash))
		},
		"VerifyAddressSignedHash": func(address, sig, hash string) interface{} {
			return checked(liteclient.VerifyAddressSignedHash(address, sig, hash))
		},
		"VerifySignatureRecoverPubKey": func(sig, hash string) interface{} {
			return checked(liteclient.VerifySignatureRecoverPubKey(sig, hash))
		},
		"VerifySeckey": func(seckey string) interface{} {
			return checked(liteclient.VerifySeckey(seckey))
		},
		"VerifyPubkey": func(pubkey string) interface{} {
			return checked(liteclient.VerifyPubkey(pubkey))
		},
		"AddressFromPubKey": func(pubkey string) interface{} {
			return result(liteclient.AddressFromPubKey(pubkey))
		},
		"AddressFromSecKey": func(seckey string) interface{} {
			return result(liteclient.AddressFromSecKey(seckey))
		},
		"PubKeyFromSig": func(sig, hash string) interface{} {
			return result(liteclient.PubKeyFromSig(sig, hash))
		},
		"SignHash": func(hash, seckey string) interface{} {
			return result(liteclient.SignHash(hash, seckey))
		},
	})
}
