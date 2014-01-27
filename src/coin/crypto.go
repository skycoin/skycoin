package coin

import (
    "errors"
    "github.com/skycoin/skycoin/src/lib/secp256k1-go"
    "log"
)

type PubKey [33]byte

func NewPubKey(b []byte) *PubKey {
    var p PubKey
    p.Set(b)
    return &p
}

func (g *PubKey) Set(b []byte) {
    if len(b) != 33 {
        log.Panic("Invalid public key length")
    }
    copy(g[:], b[:])
}

type SecKey [32]byte

func NewSecKey(b []byte) *SecKey {
    var p SecKey
    p.Set(b)
    return &p
}

func (g *SecKey) Set(b []byte) {
    if len(b) != 32 {
        log.Panic("Invalid secret key length")
    }
    copy(g[:], b[:])
}

type Sig [64 + 1]byte

func NewSig(b []byte) *Sig {
    var p Sig
    p.Set(b)
    return &p
}

func (g *Sig) Set(b []byte) {
    if len(b) != 65 {
        log.Panic("Invalid signature length")
    }
    copy(g[:], b[:])
}

/*
	Checks whether PubKey corresponding to address hash signed hash
	- recovers the PubKey from sig and hash
	- fail if PubKey cannot be be recovered
	- computes the address from the PubKey
	- fail if recovered address does not match PubKey hash
	- verify that signature is valid for hash for PubKey
*/
func ChkSig(address *Address, hash *SHA256, sig *Sig) error {
    rawPubKey := secp256.RecoverPubkey(hash[:], sig[:])
    if rawPubKey == nil {
        return errors.New("ChkSig Error: signature invalid, PubKey recovery failed")
    }
    if !address.Equals(AddressFromRawPubkey(rawPubKey)) {
        return errors.New("ChkSig Error: signature invalid, address does not match output address")
    }
    if secp256.VerifySignature(hash[:], sig[:], rawPubKey) != 1 {
        return errors.New("ChkSig Error: signature invalid, signature invalid for hash")
    }
    return nil
}

func SignHash(hash *SHA256, sec *SecKey) (*Sig, error) {
    sig := secp256.Sign(hash[:], sec[:])
    if sig == nil {
        log.Panic("SignHash invalid private key")
        return &Sig{}, errors.New("SignHash invalid private key")
    }
    return NewSig(sig), nil
}

// TODO -- implement
func PubKeyFromSec(sec *SecKey) *PubKey {
    return &PubKey{}
}

func GenerateSignature(seckey *SecKey, msg []byte) *Sig {
    if secp256.VerifySeckey(seckey[:]) != 1 {
        log.Panic("Invalid secret key")
    }
    sig := secp256.Sign(msg, seckey[:])
    return NewSig(sig)
}

func VerifySignature(pubkey *PubKey, sig *Sig, msg []byte) error {
    if secp256.VerifyPubkey(pubkey[:]) != 1 {
        log.Panic("Invalid public key")
        return errors.New("Invalid public key")
    }
    if secp256.VerifySignatureValidity(sig[:]) != 1 {
        log.Panic("Invalid signature")
        return errors.New("Invalid signature")
    }
    if secp256.VerifySignature(msg, sig[:], pubkey[:]) != 1 {
        return errors.New("Invalid signature for this message")
    }
    return nil
}

func GenerateKeyPair() (*PubKey, *SecKey) {
    public, secret := secp256.GenerateKeyPair()
    return NewPubKey(public), NewSecKey(secret)
}
