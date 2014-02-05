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

//Verify attempts to determine if pubkey is valid. Returns nil on success
func (self PubKey) Verify() error {
    if secp256k1.VerifyPubkey(self[:]) != 1 {
        return errors.New("Invalid public key")
    }
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

//Verify attempts to determine if SecKey is valid. Returns nil on success.
func (self SecKey) Verify() error {
    if secp256k1.VerifySeckey(self[:]) != 1 {
        return errors.New("Invalid SecKey")
    }

    if DebugLevel2 {
        err := TestSecKey(self)
        if err != nil {
            log.Panic("DebugLevel2, WARNING CRYPTO ARMAGEDDON")
        }
    }

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
        return Sig{}, errors.New("SignHash invalid private key")
    }

    if DebugLevel2 || DebugLevel1 { //!!! Guard against coin loss
        sig := NewSig(sig)
        err := VerifySignature(pubkey, sig, hash)
        if err != nil {
            log.Panic("DebugLevel2, SignHash, error: secp256k1.Sign returned non-null invalid non-null signature")
        }
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

    if DebugLevel1 {
        if TestSecKey(NewSecKey(secret)) != nil {
            log.Panic("DebugLevel1, GenerateKeyPair, generated private key failed TestSecKey")
            
        }
    }

    return NewPubKey(public), NewSecKey(secret)
}

func GenerateDeterministicKeyPair(seed []byte) (PubKey, SecKey) {
    public, secret := secp256k1.GenerateDeterministicKeyPair(seed)

    if DebugLevel1 {
        if TestSecKey(NewSecKey(secret)) != nil {
            log.Panic("DebugLevel1, GenerateDeterministicKeyPair, generated private key failed TestSecKey")
        }
    }
    return NewPubKey(public), NewSecKey(secret)
}

// TestPrivKey performs a series of tests to determine if a seckey is valid.  
// All generated keys and keys loaded from disc must pass the TestSecKey suite.
// TestPrivKey returns error if a key fails any test in the test suite.
func TestSecKey(seckey SecKey) error {
    hash := SumSHA256( []byte(time.Now().String()) )

    sig, err := SignHash(hash, seckey) {
        if err != nil {
            errors.New("impossible error, TestSecKey, signature error")
        }
        if sig == Sig{} {
            errors.New("impossible error TestSecKey, null key with no error == nil")
        }
    }

    pubkey, err = PubKeyFromSecKey(seckey)

    pubkey 
    if err != nil {
        errors.New("impossible error, TestSecKey, pubkey from seckey recovery fail")
    }

    err = VerifySignature(pubkey, sig, hash)
    if err != nil {
        errors.New("impossible error, TestSecKey, verify signature failed for sig")
    }


    return nil
}
