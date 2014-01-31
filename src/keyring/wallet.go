package keyring

import (
    "encoding/hex"
    "errors"
    "fmt"
    "github.com/skycoin/skycoin/src/coin"
    "log"
    "math/rand"
)

func ToHex(b []byte) string {
    return hex.EncodeToString(b)
}

func MustFromHex(s string) []byte {
    b, err := hex.DecodeString(s)
    if err != nil {
        log.Panic(err)
    }
    return b
}

// An Address contains a public and secret key,
type Address struct {
    Address coin.Address
    PubKey  coin.PubKey
    SecKey  coin.SecKey // Keep secret
}

func NewAddress() Address {
    pub, sec := coin.GenerateKeyPair()
    return Address{
        Address: coin.AddressFromPubKey(pub),
        PubKey:  pub,
        SecKey:  sec,
    }
}

/*
   genesis address
   pub: 02fa939957e9fc52140e180264e621c2576a1bfe781f88792fb315ca3d1786afb8
   sec: f2de015e7096a8b2416ea6232ea3ee2f9b928bf4672bc9801dea6ef9a0aee7a0

*/

type Wallet struct {
    Addresses []Address    //address array
    Outputs   []coin.UxOut //unspent outputs, must be refreshed
}

// Returns a random Address
func (self *Wallet) GetRandomAddress() Address {
    index := rand.Int() % len(self.Addresses)
    return self.Addresses[index]
}

// Signs a hash for a given address. nil on failure.
// use error
func (self *Wallet) Sign(address coin.Address, hash coin.SHA256) (coin.Sig, error) {
    for _, a := range self.Addresses {
        if a.Address == address {
            sig, err := coin.SignHash(hash, a.SecKey)
            if err != nil {
                log.Panic("Error signing")
                return coin.Sig{}, err
            }
            return sig, err
        }
    }
    return coin.Sig{}, errors.New("Failed to sign: Unknown address")
}

// Refresh the unspent outputs for the wallet
func (self *Wallet) RefeshUnspentOutputs(bc *coin.Blockchain) {
    outputs := make([]coin.UxOut, 0)
    for _, a := range self.Addresses {
        unspentOutputs := bc.GetUnspentOutputs(a.Address)
        outputs = append(outputs, unspentOutputs...)
    }
    self.Outputs = outputs
}

// Returns the wallet's coins and coin hours balance
func (self *Wallet) Balance(bc *coin.Blockchain) (coins uint64, hours uint64) {
    self.RefeshUnspentOutputs(bc)
    t := bc.Head.Header.Time
    for _, ux := range self.Outputs {
        coins += ux.Body.Coins
        hours += ux.CoinHours(t)
    }
    return
}

/*
func (self *Wallet) NewTransaction(bc *coin.Blockchain, Address coin.Address, amt1 uint64, amt2 uint64) (coin.Transaction, error) {
	self.RefeshUnspentOutputs(bc)
	bal1, bal2 := self.Balance()

	if bal1 < amt1 {
		return coin.Transaction{}, errors.New("insufficient coin balance")
	}
	if bal2 < amt2 {
		return coin.Transaction{}, errors.New("insufficient coinhour balance")
	}

	//decide which outputs get spent

	var ti coin.TransactionInput
	ti.SigIdx = uint16(0)
	ti.UxOut = genesisWallet.Outputs[0].Hash()
	T.TI = append(T.TI, ti)

	var to coin.TransactionOutput
	to.DestinationAddress = genesisWallet.Addresses[0].Address
	to.Value1 = uint64(100*1000000 - wn*1000)
	T.TO = append(T.TO, to)

	for i := 0; i < wn; i++ {
		var to coin.TransactionOutput
		a := WA[i].GetRandom()
		to.DestinationAddress = a.Address
		to.Value1 = 1000
		to.Value2 = 1024
		T.TO = append(T.TO, to)
	}

	var sec coin.SecKey
	sec.Set(genesisAddress.SecKey)
	T.SetSig(0, sec)

}
*/

// Creates a new wallet with n addresses
func NewWallet(n int) Wallet {
    var w Wallet
    for i := 0; i < n; i++ {
        w.Addresses = append(w.Addresses, NewAddress())
    }
    return w
}

//func (self *Address) GetOutputs(bc coin.Blockchain) []coin.UxOut {
//	ux := bc.GetUnspentOutputs(*self.Address)
//	return ux
//}

//pub, sec := util.GenerateKeyPair()
//pub := MustF("02fa939957e9fc52140e180264e621c2576a1bfe781f88792fb315ca3d1786afb8")
//sec := MustF("f2de015e7096a8b2416ea6232ea3ee2f9b928bf4672bc9801dea6ef9a0aee7a0")
//fmt.Printf("sec: %s \n", ToHex(sec))
//fmt.Printf("pub: %s \n", ToHex(pub))

func uxStr(ux coin.UxOut) string {
    return fmt.Sprintf("%s, %d: %d %d", ux.Body.Address.String(), ux.Head.Time,
        ux.Body.Coins, ux.Body.Hours)
}

// Prints the balances for multiple wallets
func PrintWalletBalances(bc *coin.Blockchain, wallets []Wallet) {
    for i, w := range wallets {
        b1, b2 := w.Balance(bc)
        fmt.Printf("PWB: %v: %v %v \n", i, b1, b2)
    }

    for i, ux := range bc.Unspent {
        fmt.Printf("PWB,UX: %v: %v \n", i, uxStr(ux))
    }
}
