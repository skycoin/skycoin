package visor

import (
	"bytes"
	"errors"
	//"fmt"
	"log"
	"sort"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/wallet"
)

// Deprecate dependency on wallet
// DEPRECATED. CAN BE DELETED

/*

Sort unspents oldest to newest

Keep adding until either exact amount (coins+hours), or if hours exceeded,
also exceed coins by at least 1e6


*/

// OldestUxOut sorts a UxArray oldest to newest.
type OldestUxOut coin.UxArray

func (ouo OldestUxOut) Len() int      { return len(ouo) }
func (ouo OldestUxOut) Swap(i, j int) { ouo[i], ouo[j] = ouo[j], ouo[i] }
func (ouo OldestUxOut) Less(i, j int) bool {
	a := ouo[i].Head.BkSeq
	b := ouo[j].Head.BkSeq
	// Use hash to break ties
	if a == b {
		ih := ouo[i].Hash()
		jh := ouo[j].Hash()
		cmp := bytes.Compare(ih[:], jh[:])
		if cmp == 0 {
			log.Panic("Duplicate UxOut when sorting")
		}
		return cmp < 0
	}
	return a < b
}

//delete this function
func createSpends(headTime uint64, uxa coin.UxArray,
	amt wallet.Balance) (coin.UxArray, error) {
	if amt.Coins == 0 {
		return nil, errors.New("Zero spend amount")
	}
	if amt.Coins%1e6 != 0 {
		return nil, errors.New("Coins must be multiple of 1e6")
	}

	uxs := OldestUxOut(uxa)
	sort.Sort(uxs)

	have := wallet.Balance{Coins: 0, Hours: 0}
	spending := make(coin.UxArray, 0)
	for i := range uxs {
		b := wallet.NewBalanceFromUxOut(headTime, &uxs[i]) //this is bullshit
		if b.Coins == 0 || b.Coins%1e6 != 0 {
			logger.Error("UxOut coins are 0 or 1e6, can't spend")
			continue
		}
		have = have.Add(b)
		spending = append(spending, uxs[i])
	}

	if amt.Coins > have.Coins {
		return nil, errors.New("Not enough coins")
	}

	return spending, nil
}

// CreateSpendingTransaction DEPRECATE
// deprecate dependency on wallet
// Creates a Transaction spending coins and hours from our coins
// MOVE SOMEWHERE ELSE
// Move to wallet or move to ???
func CreateSpendingTransaction(wlt wallet.Wallet,
	unconfirmed *UnconfirmedTxnPool, unspent *coin.UnspentPool,
	headTime uint64, amt wallet.Balance,
	dest cipher.Address) (coin.Transaction, error) {
	txn := coin.Transaction{}
	auxs := unspent.AllForAddresses(wlt.GetAddresses())
	// Subtract pending spends from available
	puxs := unconfirmed.SpendsForAddresses(unspent, wlt.GetAddressSet())
	auxs = auxs.Sub(puxs)

	// Determine which unspents to spend
	spends, err := createSpends(headTime, auxs.Flatten(), amt)
	if err != nil {
		return txn, err
	}

	// Add these unspents as tx inputs
	toSign := make([]cipher.SecKey, len(spends))
	spending := wallet.Balance{Coins: 0, Hours: 0}
	for i, au := range spends {
		entry, exists := wlt.GetEntry(au.Body.Address)
		if !exists {
			log.Panic("On second thought, the wallet entry does not exist")
		}
		txn.PushInput(au.Hash())
		toSign[i] = entry.Secret
		spending.Coins += au.Body.Coins
		spending.Hours += au.CoinHours(headTime)
	}

	//keep 1/4th of hours as change
	//send half to each address
	var changeHours = uint64(spending.Hours / 4)

	if amt.Coins == spending.Coins {
		txn.PushOutput(dest, amt.Coins, changeHours/2)
		txn.SignInputs(toSign)
		txn.UpdateHeader()
		return txn, nil
	}

	change := wallet.NewBalance(spending.Coins-amt.Coins, changeHours/2)
	// TODO -- send change to a new address
	changeAddr := spends[0].Body.Address

	//create transaction
	txn.PushOutput(changeAddr, change.Coins, change.Hours)
	txn.PushOutput(dest, amt.Coins, changeHours/2)
	txn.SignInputs(toSign)
	txn.UpdateHeader()
	return txn, nil
}
