package coin

import (
	"encoding/hex"
	"math/rand"
	"testing"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/util/utc"
)

var (
	genPublic, genSecret        = cipher.GenerateKeyPair()
	genAddress                  = cipher.AddressFromPubKey(genPublic)
	testMaxSize                 = 1024 * 1024
	_genTime             uint64 = 1000
	_incTime             uint64 = 3600 * 1000
	_genCoins            uint64 = 1000e6
	_genCoinHours        uint64 = 1000 * 1000
)

func tNow() uint64 {
	return uint64(utc.UnixNow())
}

func _feeCalc(t *Transaction) (uint64, error) {
	return 0, nil
}

func _makeFeeCalc(fee uint64) FeeCalculator {
	return func(t *Transaction) (uint64, error) {
		return fee, nil
	}
}

func TestAddress1(t *testing.T) {
	a := "02fa939957e9fc52140e180264e621c2576a1bfe781f88792fb315ca3d1786afb8"
	b, err := hex.DecodeString(a)
	if err != nil {
		t.Fatal(err)
	}
	addr := cipher.AddressFromPubKey(cipher.NewPubKey(b))
	_ = addr

	///func SignHash(hash cipher.SHA256, sec SecKey) (Sig, error) {

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

	seckey := cipher.NewSecKey(b)
	pubkey := cipher.PubKeyFromSecKey(seckey)
	addr := cipher.AddressFromPubKey(pubkey)
	_ = addr

	///func SignHash(hash cipher.SHA256, sec SecKey) (Sig, error) {

}

//TODO: 100% coverage over cryptographic functions

//Crypto Functions to Test
//func ChkSig(address Address, hash cipher.SHA256, sig Sig) error {
//func SignHash(hash cipher.SHA256, sec SecKey) (Sig, error) {
//func cipher.PubKeyFromSecKey(seckey SecKey) (PubKey) {
//func PubKeyFromSig(sig Sig, hash cipher.SHA256) (PubKey, error) {
//func VerifySignature(pubkey PubKey, sig Sig, hash cipher.SHA256) error {
//func GenerateKeyPair() (PubKey, SecKey) {
//func GenerateDeterministicKeyPair(seed []byte) (PubKey, SecKey) {
//func testSecKey(seckey SecKey) error {

func TestCrypto1(t *testing.T) {
	for i := 0; i < 10; i++ {
		_, seckey := cipher.GenerateKeyPair()
		if cipher.TestSecKey(seckey) != nil {
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

	seckey := cipher.NewSecKey(b)
	pubkey := cipher.PubKeyFromSecKey(seckey)

	addr := cipher.AddressFromPubKey(pubkey)
	_ = addr

	test := []byte("test message")
	hash := cipher.SumSHA256(test)
	err = cipher.TestSecKeyHash(seckey, hash)
	if err != nil {
		t.Fatal()
	}

}

func _gensec() cipher.SecKey {
	_, s := cipher.GenerateKeyPair()
	return s
}

func _gpub(s cipher.SecKey) cipher.PubKey {
	return cipher.PubKeyFromSecKey(s)
}

func _gaddr(s cipher.SecKey) cipher.Address {
	return cipher.AddressFromSecKey(s)
}

func _gaddrA1(S []cipher.SecKey) []cipher.Address {
	A := make([]cipher.Address, len(S))
	for i := 0; i < len(S); i++ {
		A[i] = cipher.AddressFromSecKey(S[i])
	}
	return A
}

func _gaddrA2(S []cipher.SecKey, O []UxOut) []int {
	A := _gaddrA1(S)
	var M map[cipher.Address]int //address to int
	for i, a := range A {
		M[a] = i
	}

	I := make([]int, len(O)) //output to seckey/address index
	for i, o := range O {
		I[i] = M[o.Body.Address]
	}

	return I
}

func _gaddrA3(S []cipher.SecKey, O []UxOut) map[cipher.Address]int {
	A := _gaddrA1(S)
	M := make(map[cipher.Address]int) //address to int
	for i, a := range A {
		M[a] = i
	}
	return M
}

//assign amt to n bins in randomized manner
func _randBins(amt uint64, n int) []uint64 {
	bins := make([]uint64, n)
	max := amt / (4 * uint64(n))
	for i := 0; amt > 0; i++ {
		//amount going into this bin
		b := 1 + (uint64(rand.Int63()) % max)
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
            to.Address = A[rand.Intn(len(A))]
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
    ConnectionMirrors[addr] = uint32(4)
    assert.Panics(t, func() { getListenPort(addr) })
    // Everything is good
    m := make(map[string]uint16)
    mirrorConnections[uint32(4)] = m
    m[addrIP] = uint16(6667)
    assert.Equal(t, getListenPort(addr), uint16(6667))

    // cleanup
    delete(mirrorConnections, uint32(4))
    delete(ConnectionMirrors, addr)
}
*/
