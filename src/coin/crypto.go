package coin

import (
    "errors"
    "github.com/skycoin/skycoin/src/lib/secp256k1-go"
    "log"
)

type PubKey [33]byte

func NewPubKey(b []byte) PubKey {
    if len(b) != 33 {
        log.Panic("Invalid public key length")
    }
    var p PubKey
    copy(p[:], b[:])
    return p
}
type SecKey [32]byte

func NewSecKey(b []byte) SecKey {
    if len(b) != 32 {
        log.Panic("Invalid secret key length")
    }
    var p SecKey
    copy(p[:], b[:])
    return p
}

type Sig [64 + 1]byte

func NewSig(b []byte) Sig {
    if len(b) != 65 {
        log.Panic("Invalid secret key length")
    }
    var p Sig
    copy(p[:],b[:])
    return p
}

/*
	Checks whether PubKey corresponding to address hash signed hash
	- recovers the PubKey from sig and hash
	- fail if PubKey cannot be be recovered
	- computes the address from the PubKey
	- fail if recovered address does not match PubKey hash
	- verify that signature is valid for hash for PubKey
*/
func ChkSig(address Address, hash SHA256, sig Sig) error {
    rawPubKey := secp256k1.RecoverPubkey(hash[:], sig[:])
    if rawPubKey == nil {
        return errors.New("Invalig sig: PubKey recovery failed")
    }
    if address != AddressFromPubKey(NewPubKey(rawPubKey)) {
        return errors.New("Invalid sig: address does not match output address")
    }
    if secp256k1.VerifySignature(hash[:], sig[:], rawPubKey[:]) != 1 {
        return errors.New("Invalid sig: invalid for hash")
    }
    return nil
}

func SignHash(hash SHA256, sec SecKey) (Sig, error) {
    sig := secp256k1.Sign(hash[:], sec[:])
    if sig == nil {
        log.Panic("SignHash invalid private key")
        return Sig{}, errors.New("SignHash invalid private key")
    }
    return NewSig(sig), nil
}

/*
func SignMessage(seckey SecKey, msg []byte) Sig {
    if secp256k1.VerifySeckey(seckey[:]) != 1 {
        log.Panic("Invalid secret key")
    }
    sig := secp256k1.Sign(msg, seckey[:])
    return NewSig(sig)
}
*/

func PubKeyFromSecKey(seckey SecKey) PubKey {
    b := secp256k1.PubkeyFromSeckey(seckey[:])
    if b == nil {
        log.Panic("could not recover pubkey form sec key \n")
        return PubKey{}
    }
    return NewPubKey(b)
}

func PubKeyFromSig(sig Sig, hash SHA256) (PubKey, error) {
    rawPubKey := secp256k1.RecoverPubkey(hash[:], sig[:])
    if rawPubKey == nil {
        return PubKey{}, errors.New("Invalig sig: PubKey recovery failed")
    }
    return NewPubKey(rawPubKey), nil
}

//verifies that mesh hash was signed by pubkey
func VerifySignature(pubkey PubKey, sig Sig, hash SHA256) error {
    pubkey_rec, err := PubKeyFromSig(sig, hash) //recovered pubkey
    if err != nil {
        return errors.New("Invalig sig: PubKey recovery failed")
    }
    if pubkey_rec != pubkey {
        return errors.New("Recovered pubkey does not match pubkey")
    }
    if secp256k1.VerifyPubkey(pubkey[:]) != 1 {
        log.Panic("Invalid public key")
        return errors.New("Invalid public key")
    }
    if secp256k1.VerifySignatureValidity(sig[:]) != 1 {
        log.Panic("Invalid signature")
        return errors.New("Invalid signature")
    }
    if secp256k1.VerifySignature(hash[:], sig[:], pubkey[:]) != 1 {
        return errors.New("Invalid signature for this message")
    }
    return nil
}

func GenerateKeyPair() (PubKey, SecKey) {
    public, secret := secp256k1.GenerateKeyPair()
    return NewPubKey(public), NewSecKey(secret)
}