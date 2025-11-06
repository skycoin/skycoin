// Package main provides the GopherJS skycoin-lite browser implementation
package main

import (
	"github.com/gopherjs/gopherjs/js"

	"github.com/skycoin/skycoin/src/skycoin-lite/liteclient"
)

func main() {
	js.Global.Set("Cipher", map[string]interface{}{
		"GenerateAddresses":                liteclient.GenerateAddress,
		"PrepareTransaction":               liteclient.PrepareTransaction,
		"PrepareTransactionWithSignatures": liteclient.PrepareTransactionWithSignatures,
	})

	js.Global.Set("CipherExtras", map[string]interface{}{
		"VerifyPubKeySignedHash":       liteclient.VerifyPubKeySignedHash,
		"VerifyAddressSignedHash":      liteclient.VerifyAddressSignedHash,
		"VerifySignatureRecoverPubKey": liteclient.VerifySignatureRecoverPubKey,
		"VerifySeckey":                 liteclient.VerifySeckey,
		"VerifyPubkey":                 liteclient.VerifyPubkey,
		"AddressFromPubKey":            liteclient.AddressFromPubKey,
		"AddressFromSecKey":            liteclient.AddressFromSecKey,
		"PubKeyFromSig":                liteclient.PubKeyFromSig,
		"SignHash":                     liteclient.SignHash,
	})
}
