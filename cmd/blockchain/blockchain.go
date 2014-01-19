package main

import (
    "encoding/hex"
    "fmt"
    "github.com/skycoin/skycoin/src/coin"
    "log"
    "math/rand"
)

func Run() {
    tests()
}

func ToHex(b []byte) string {
    return hex.EncodeToString(b)
}

func FromHex(s string) []byte {
    b, err := hex.DecodeString(s)
    if err != nil {
        log.Panic(err)
    }
    return b
}

/*
   genesis address
   pub: 02fa939957e9fc52140e180264e621c2576a1bfe781f88792fb315ca3d1786afb8
   sec: f2de015e7096a8b2416ea6232ea3ee2f9b928bf4672bc9801dea6ef9a0aee7a0

*/

type Wallet struct {
    AA      []Address    //address array
    Outputs []coin.UxOut //unspent outputs
}

func (self *Wallet) GetRandom() Address {
    index := rand.Int() % len(self.AA)
    return self.AA[index]
}

//signs a hash for a given address. nil on failure
func (self *Wallet) Sign(address coin.Address, hash []byte) []byte {
    for _, a := range self.AA {
        if a.address == address {
            //func GenerateSignature(seckey []byte, msg []byte) []byte
            sig := coin.GenerateSignature(a.seckey, hash)
            return sig
        }

    }
    return nil
}

//refresh the unspent outputs for the wallet
func (self *Wallet) RefeshUnspentOutputs(bc *coin.BlockChain) {
    var outputs []coin.UxOut
    for _, a := range self.AA {
        //get unspent outputs for the address
        unspentOutputs := bc.GetUnspentOutputs(a.address)
        outputs = append(outputs, unspentOutputs...)
    }
    self.Outputs = outputs
}

func (self *Wallet) Balance(bc *coin.BlockChain) (uint64, uint64) {
    self.RefeshUnspentOutputs(bc)

    var balance1 uint64 = 0
    var balance2 uint64 = 0

    var time uint64 = bc.Head.Header.Time

    for _, ux := range self.Outputs {
        balance1 += ux.Body.Value1
        balance2 += ux.CoinHours(time)
    }
    return balance1, balance2
}

func NewWallet(n int) Wallet {
    var w Wallet
    for i := 0; i < n; i++ {
        w.AA = append(w.AA, GenerateAddress())
    }
    return w
}

type Address struct {
    pubkey  []byte
    seckey  []byte
    address coin.Address
}

func GenerateAddress() Address {
    var A Address
    A.pubkey, A.seckey = coin.GenerateKeyPair()
    A.address = coin.AddressFromRawPubkey(A.pubkey)
    return A
}

//func (self *Address) GetOutputs(bc coin.BlockChain) []coin.UxOut {
//	ux := bc.GetUnspentOutputs(*self.address)
//	return ux
//}

//pub, sec := util.GenerateKeyPair()
//pub := FromHex("02fa939957e9fc52140e180264e621c2576a1bfe781f88792fb315ca3d1786afb8")
//sec := FromHex("f2de015e7096a8b2416ea6232ea3ee2f9b928bf4672bc9801dea6ef9a0aee7a0")
//fmt.Printf("sec: %s \n", ToHex(sec))
//fmt.Printf("pub: %s \n", ToHex(pub))

func WalletBalances(bc *coin.BlockChain, wallets []Wallet) {
    for i, w := range wallets {

        b1, b2 := w.Balance(bc)
        fmt.Printf("%v: %v %v \n", i, b1, b2)
    }

    for i, ux := range bc.Unspent {
        fmt.Printf("%v: %v \n", i, ux.String())
    }
}

func tests() {

    genesisWallet := NewWallet(1)
    genesisAddress := genesisWallet.AA[0]
    var BC *coin.BlockChain = coin.NewBlockChain(genesisAddress.address)

    genesisWallet.RefeshUnspentOutputs(BC)
    //create 16 wallets with 64 addresses
    wn := 16 //number of wallets to create
    var WA []Wallet
    for i := 0; i < wn; i++ {
        WA = append(WA, NewWallet(64))
    }

    B := BC.NewBlock()

    var T coin.Transaction

    if true {
        //create transaction by hand
        var ti coin.TransactionInput
        ti.SigIdx = uint16(0)
        ti.UxOut = genesisWallet.Outputs[0].Hash()
        T.TI = append(T.TI, ti)

        var to coin.TransactionOutput
        to.DestinationAddress = genesisWallet.AA[0].address
        to.Value1 = uint64(100*1000000 - wn*1000)
        T.TO = append(T.TO, to)

        for i := 0; i < wn; i++ {
            var to coin.TransactionOutput
            a := WA[i].GetRandom()
            to.DestinationAddress = a.address
            to.Value1 = 1000
            to.Value2 = 1024
            T.TO = append(T.TO, to)
        }

        var sec coin.SecKey
        sec.Set(genesisAddress.seckey)
        T.SetSig(0, sec)

    } else {
        T.PushInput(genesisWallet.Outputs[0].Hash())
        T.PushOutput(genesisWallet.AA[0].address, uint64(100*1000000-wn*1000), 0)

        for i := 0; i < wn; i++ {
            a := WA[i].GetRandom()
            T.PushOutput(a.address, uint64(1000), 1024*1024)
        }

        var sec coin.SecKey
        sec.Set(genesisAddress.seckey)
        T.SetSig(0, sec)
    }
    T.UpdateHeader() //sets hash

    fmt.Printf("genesis transaction: \n")
    /*
       Need to add input
    */
    err := BC.AppendTransaction(B, &T)

    if err != nil {
        log.Panic(err)
    }

    WalletBalances(BC, WA)
    err = BC.ExecuteBlock(B)
    if err != nil {
        log.Panic(err)
    }

    WalletBalances(BC, WA)

}

func main() {
    tests()
}
