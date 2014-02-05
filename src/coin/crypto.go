package coin

import (
    "encoding/hex"
    "errors"
    "github.com/skycoin/skycoin/src/lib/secp256k1-go"
    "log"
)

type PubKey [33]byte

func NewPubKey(b []byte) PubKey {
    var p PubKey
    if len(b) != len(p) {
        log.Panic("Invalid public key length")
    }
    copy(p[:], b[:])
    return p
}

func PubKeyFromHex(s string) PubKey {
    b, err := hex.DecodeString(s)
    if err != nil {
        log.Panic(err)
    }
    return NewPubKey(b)
}

func (self *PubKey) Hex() string {
    return hex.EncodeToString(self[:])
}

// Returns the public key as ripemd160(sha256(sha256(key)))
func (self *PubKey) ToAddressHash() Ripemd160 {
    r1 := SumSHA256(self[:])
    r2 := SumSHA256(r1[:])
    return HashRipemd160(r2[:])
}

type SecKey [32]byte

func NewSecKey(b []byte) SecKey {
    var p SecKey
    if len(b) != len(p) {
        log.Panic("Invalid secret key length")
    }
    copy(p[:], b[:])
    return p
}

func SecKeyFromHex(s string) SecKey {
    b, err := hex.DecodeString(s)
    if err != nil {
        log.Panic(err)
    }
    return NewSecKey(b)
}

func (s SecKey) Hex() string {
    return hex.EncodeToString(s[:])
}

type Sig [64 + 1]byte //64 byte signature with 1 byte for key recovery

func NewSig(b []byte) Sig {
    var s Sig
    if len(b) != len(s) {
        log.Panic("Invalid secret key length")
    }
    copy(s[:], b[:])
    return s
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
        if Paranoid {
            log.Panic("Paranoid: SignHash invalid private key")
        }
        return Sig{}, errors.New("SignHash invalid private key")
    }
    return NewSig(sig), nil
}

func PubKeyFromSecKey(seckey SecKey) PubKey {
    b := secp256k1.PubkeyFromSeckey(seckey[:])
    if b == nil {
        log.Panic("could not recover pubkey from seckey")
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
        log.Panic("Invalid public key: check public key before signature verification")
        //return errors.New("Invalid public key")
    }
    if secp256k1.VerifySignatureValidity(sig[:]) != 1 {
        log.Panic("Invalid signature: check signature validity before verification")
        //return errors.New("Invalid signature")
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
func GenerateDeterministicKeyPair(seed []byte) (PubKey, SecKey) {
    public, secret := secp256k1.GenerateDeterministicKeyPair(seed)
    return NewPubKey(public), NewSecKey(secret)
}
