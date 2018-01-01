package visor

import (
	"sort"
	"strings"

	"github.com/skycoin/skycoin/src/util/droplet"
)

type Account struct {
	Addr   string
	Coins  uint64
	Locked bool
}

//AccountJson topn elements
type AccountJSON struct {
	Addr   string `json:"address"`
	Coins  string `json:"coins"`
	Locked bool   `json:"locked"`
}

//AccountMgr manager all unspent outputs
type AccountMgr struct {
	Accounts []Account
}

//NewAccountMgr create AccountMgr via unspent outputs map
func NewAccountMgr(allAccount map[string]uint64, distributionMap map[string]struct{}) *AccountMgr {
	am := &AccountMgr{}
	am.Accounts = make([]Account, 0, len(allAccount))
	var islocked bool
	for acc, value := range allAccount {
		if _, ok := distributionMap[acc]; ok {
			islocked = true
		} else {
			islocked = false
		}
		am.Accounts = append(am.Accounts, Account{Addr: acc, Coins: value, Locked: islocked})
	}
	return am
}

func (am AccountMgr) Len() int {
	return len(am.Accounts)
}

func (am AccountMgr) Less(i, j int) bool {
	if am.Accounts[i].Locked == am.Accounts[j].Locked {
		if am.Accounts[i].Coins == am.Accounts[j].Coins {
			//sort alphabetically
			cp := strings.Compare(am.Accounts[i].Addr, am.Accounts[j].Addr)
			if cp > 0 {
				return true
			}
			return false
		}
		return am.Accounts[i].Coins > am.Accounts[j].Coins
	}
	if am.Accounts[i].Locked {
		return true
	}
	return false
}

func (am AccountMgr) Swap(i, j int) {
	am.Accounts[i], am.Accounts[j] = am.Accounts[j], am.Accounts[i]
}

//Sort sort coin owner desc order
func (am *AccountMgr) Sort() {
	sort.Sort(am)
}

//GetTopn returns topn rich owner, returns all if topn <= 0, exclude distribution if includeDistribution = false
func (am *AccountMgr) GetTopn(topn int, includeDistribution bool) ([]AccountJSON, error) {
	topnAccount := []AccountJSON{}
	for _, acc := range am.Accounts {
		//skip special address
		if !includeDistribution {
			if acc.Locked {
				continue
			}
		}
		coinsStr, err := droplet.ToString(acc.Coins)
		if err != nil {
			return topnAccount, err
		}
		topnAccount = append(topnAccount, AccountJSON{Addr: acc.Addr, Coins: coinsStr, Locked: acc.Locked})
		//return all if topn is -1
		if topn > 0 && len(topnAccount) >= topn {
			break
		}
	}
	return topnAccount, nil
}
