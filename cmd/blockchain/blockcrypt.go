package main


import (
    //"encoding/hex"
    //"errors"
    "fmt"
    "github.com/skycoin/skycoin/src/coin"
    //"github.com/skycoin/skycoin/src/keyring"

    "log"
    //"math/rand"
    "encoding/hex"
)


var master_pubkey coin.PubKey
var master_seckey coin.SecKey

func init() {

	a := "5a42c0643bdb465d90bf673b99c14f5fa02db71513249d904573d2b8b63d353d"
	b, err := hex.DecodeString(a)
	if err != nil || len(b) != 32 {
		log.Panic(err)
	}

    master_seckey := NewSecKey(b)
    master_pubkey := PubKeyFromSecKey(seckey)

	master_seckey = seckey
	master_pubkey = pubkey
}

//func SignBlock(block coin.Block, seckey coin.SecKey) (coin.Sig, error) {
	//var block_hash coin.SHA256 = block.HashHeader()
	//sig, err := SignHash(block_hash, seckey)
	//return sig,err
//}

//sign a block with a private key
func SignBlock(block coin.Block, seckey coin.SecKey) (coin.Sig, error) {
	return coin.SignHash(block.HashHeader(), seckey)
}

func VerifyBlockSignature(block coin.Block, sig Sig) error {

	block_hash := block.HashHeader()
	pubkey = 
}