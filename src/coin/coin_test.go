package coin

import (
    //"crypto/sha256"
    //"hash"
    "encoding/hex"
    //"fmt"
    //"log"
    "math/rand"
    "testing"
    "errors"
    "github.com/skycoin/skycoin/src/lib/secp256k1-go"
)

func TestAddress1(t *testing.T) {
    a := "02fa939957e9fc52140e180264e621c2576a1bfe781f88792fb315ca3d1786afb8"
    b, err := hex.DecodeString(a)
    if err != nil {
        t.Fatal(err)
    }
    addr := AddressFromPubKey(NewPubKey(b))
    _ = addr

    ///func SignHash(hash SHA256, sec SecKey) (Sig, error) {

}

func TestAddress2(t *testing.T) {
    a := "5a42c0643bdb465d90bf673b99c14f5fa02db71513249d904573d2b8b63d353d"
    b, err := hex.DecodeString(a)
    if err != nil {
        t.Fail()
    }

    if len(b) != 32 {
        t.Fail()
    }

    seckey := NewSecKey(b)
    pubkey := PubKeyFromSecKey(seckey)
    addr := AddressFromPubKey(pubkey)
    _ = addr

    ///func SignHash(hash SHA256, sec SecKey) (Sig, error) {

}

//TODO: 100% coverage over cryptographic functions

//Crypto Functions to Test
//func ChkSig(address Address, hash SHA256, sig Sig) error {
//func SignHash(hash SHA256, sec SecKey) (Sig, error) {
//func PubKeyFromSecKey(seckey SecKey) (PubKey) {
//func PubKeyFromSig(sig Sig, hash SHA256) (PubKey, error) {
//func VerifySignature(pubkey PubKey, sig Sig, hash SHA256) error {
//func GenerateKeyPair() (PubKey, SecKey) {
//func GenerateDeterministicKeyPair(seed []byte) (PubKey, SecKey) {
//func TestSecKey(seckey SecKey) error {


func TestCrypto1(t *testing.T) {
    for i:=0; i<10; i++ {
        _,seckey := GenerateKeyPair()
        if TestSecKey(seckey) != nil {
            t.Fatal("CRYPTOGRAPHIC INTEGRITY CHECK FAILED")
        }
    }
}

//test signatures
func TestCrypto2(t *testing.T) {
    a := "5a42c0643bdb465d90bf673b99c14f5fa02db71513249d904573d2b8b63d353d"
    b, err := hex.DecodeString(a)
    if err != nil {
        t.Fatal(err)
    }

    if len(b) != 32 {
        t.Fatal()
    }

    seckey := NewSecKey(b)
    pubkey := PubKeyFromSecKey(seckey)

    addr := AddressFromPubKey(pubkey)
    _ = addr

    test := []byte("test message")
    hash := SumSHA256(test)

    err = testSecKey(seckey, hash)
    if err != nil {
        t.Fatal()
    }

}

//testSecKey tests a hash and seckey in all known cryptographic contexts
func testSecKey(seckey SecKey, hash SHA256) error {
    //check seckey with verify
    if secp256k1.VerifySeckey(seckey[:]) != 1 {
        return errors.New("Seckey verification failed")
    }

    //check pubkey recovery
    pubkey := PubKeyFromSecKey(seckey)
    if pubkey == (PubKey{}) {
        errors.New("impossible error, TestSecKey, nil pubkey recovered")
    }
    //verify recovered pubkey
    if secp256k1.VerifyPubkey(pubkey[:]) != 1 {
        return errors.New("impossible error, TestSecKey, Derived Pubkey verification failed")
    }

    //check signature production
    sig, err := SignHash(hash, seckey)
    if err != nil {
        errors.New("impossible error, TestSecKey, signature error")
    }
    if sig == (Sig{}) {
        errors.New("impossible error TestSecKey, nil sig with no error == nil")
    }

    //check pubkey recovered from sig
    recovered_pubkey, err := PubKeyFromSig(sig, hash)
    if err != nil {
        return errors.New("impossible error, TestSecKey, pubkey recovery from signature failed")
    }
    if pubkey != recovered_pubkey {
        return errors.New("impossible error TestSecKey, pubkey does not match recovered pubkey")
    }

    //verify produced signature
    err = VerifySignature(pubkey, sig, hash)
    if err != nil {
        errors.New("impossible error, TestSecKey, verify signature failed for sig")
    }

    //verify ChkSig
    addr := AddressFromPubKey(pubkey)
    err = ChkSig(addr,hash, sig)
    if err != nil {
        return errors.New("impossible error TestSecKey, ChkSig Failed, should not get this far")
    }

    return nil
}

func _gensec() SecKey {
    _, s := GenerateKeyPair()
    return s
}

func _gpub(s SecKey) PubKey {
    return PubKeyFromSecKey(s)
}

func _gaddr(s SecKey) Address {
    return AddressFromPubKey(PubKeyFromSecKey(s))
}

func _gaddr_a1(S []SecKey) []Address {
    A := make([]Address, len(S))
    for i := 0; i < len(S); i++ {
        A[i] = AddressFromPubKey(PubKeyFromSecKey(S[i]))
    }
    return A
}

func _gaddr_a2(S []SecKey, O []UxOut) []int {
    A := _gaddr_a1(S)
    var M map[Address]int //address to int
    for i, a := range A {
        M[a] = i
    }

    I := make([]int, len(O)) //output to seckey/address index
    for i, o := range O {
        I[i] = M[o.Body.Address]
    }

    return I
}

func _gaddr_a3(S []SecKey, O []UxOut) map[Address]int {
    A := _gaddr_a1(S)
    M := make(map[Address]int) //address to int
    for i, a := range A {
        M[a] = i
    }
    return M
}

//assign amt to n bins in randomized manner
func _rand_bins(amt uint64, n int) []uint64 {
    var bins []uint64 = make([]uint64, n)
    var max uint64 = amt / (4 * uint64(n))
    for i := 0; amt > 0; i++ {
        //amount going into this bin
        var b uint64 = 1 + (uint64(rand.Int63()) % max)
        if b > amt {
            b = amt
        }
        bins[i%n] += b
        amt -= b
    }
    return bins
}

/*
TODO: check block header of new block
TODO: check that coins are not created or destroyed
TODO:
*/

//create 4096 addresses
//send addreses randomly between each other over 1024 blocks

/*
func TestBlockchain1(t *testing.T) {

    var S []SecKey
    for i := 0; i < 4096; i++ {
        S = append(S, _gensec())
    }

    A := _gaddr_a1(S)

    var bc *Blockchain = NewBlockchain(A[0])

    for i := 0; i < 1024; i++ {
        b := bc.NewBlock()

        //unspent outputs
        U := make([]UxOut, len(bc.Unspent))
        copy(U, bc.Unspent)

        //for _,Ux := range U {
        //    if Ux.Hours() < Ux.Body.
        //}
        //I := _gaddr_a2(S,U)
        M := _gaddr_a3(S, U)
        var num_in int = 1 + rand.Intn(len(U))%15
        var num_out int = 1 + rand.Int()%30

        var t Transaction

        SigIdx := make([]int, num_in)

        var v1 uint64 = 0
        var v2 uint64 = 0
        for i := 0; i < num_in; i++ {
            idx := rand.Intn(len(U))
            var Ux UxOut = U[idx]                 //unspent output to spend
            U[idx], U = U[len(U)-1], U[:len(U)-1] //remove output idx

            v1 += Ux.Body.Coins
            v2 += Ux.Body.Hours

            //index of signature that must sign input
            SigIdx[i] = M[Ux.Body.Address] //signature index

            var ti TransactionInput
            ti.SigIdx = uint16(i)
            ti.UxOut = Ux.Hash()
            t.TxIn = append(t.TxIn, ti) //append input to transaction
        }

        //assign coins to output addresses in random manner

        //check that inputs/outputs sum
        v1_ := v1
        v2_ := v2

        vo1 := _rand_bins(v1, num_out)
        vo2 := _rand_bins(v2, num_out)

        var v1_t uint64
        var v2_t uint64
        for i, _ := range vo1 {
            v1_t += vo1[i]
            v2_t += vo2[i]
        }

        if v1_t != v1_ {
            log.Panic()
        }
        if v2_t != v2_ {
            log.Panic()
        }
        //log.Printf("%v %v, %v %v \n", v1_,v2_, v1_t, v2_t)

        for i := 0; i < num_out; i++ {
            var to TransactionOutput
            to.DestinationAddress = A[rand.Intn(len(A))]
            to.Coins = vo1[i]
            to.Hours = vo2[i]
            t.TxOut = append(t.TxOut, to)
        }

        //transaction complete, now set signatures
        for i := 0; i < num_in; i++ {
            t.SetSig(uint16(i), S[SigIdx[i]])
        }
        t.UpdateHeader() //sets hash

        err := bc.AppendTransaction(&b, t)
        if err != nil {
            log.Panic(err)
        }

        fmt.Printf("Block %v \n", i)
        err = bc.ExecuteBlock(b)
        if err != nil {
            log.Panic(err)
        }

    }
}
*/

/*
func TestGetListenPort(t *testing.T) {
    // No connectionMirror found
    assert.Equal(t, getListenPort(addr), uint16(0))
    // No mirrorConnection map exists
    connectionMirrors[addr] = uint32(4)
    assert.Panics(t, func() { getListenPort(addr) })
    // Everything is good
    m := make(map[string]uint16)
    mirrorConnections[uint32(4)] = m
    m[addrIP] = uint16(6667)
    assert.Equal(t, getListenPort(addr), uint16(6667))

    // cleanup
    delete(mirrorConnections, uint32(4))
    delete(connectionMirrors, addr)
}
*/
