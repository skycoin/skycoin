package visor

import (
	"bytes"
	"sort"

	"github.com/SkycoinProject/skycoin/src/cipher"
)

// RichlistBalance holds info an address balance holder
type RichlistBalance struct {
	Address cipher.Address
	Coins   uint64
	Locked  bool
}

// Richlist contains RichlistBalances
type Richlist []RichlistBalance

// NewRichlist create Richlist via unspent outputs map
func NewRichlist(allAccounts map[cipher.Address]uint64, lockedAddrs map[cipher.Address]struct{}) (Richlist, error) {
	richlist := make(Richlist, 0, len(allAccounts))

	for addr, coins := range allAccounts {
		var islocked bool
		if _, ok := lockedAddrs[addr]; ok {
			islocked = true
		}

		richlist = append(richlist, RichlistBalance{
			Address: addr,
			Coins:   coins,
			Locked:  islocked,
		})
	}

	// Sort order:
	// Higher coins
	// Locked > unlocked
	// Address bytes
	sort.Slice(richlist, func(i, j int) bool {
		if richlist[i].Coins == richlist[j].Coins {
			if richlist[i].Locked == richlist[j].Locked {
				return bytes.Compare(richlist[i].Address.Bytes(), richlist[j].Address.Bytes()) < 0
			}
			return richlist[i].Locked
		}

		return richlist[i].Coins > richlist[j].Coins
	})

	return richlist, nil
}

// FilterAddresses returns the richlist without addresses from the map
func (r Richlist) FilterAddresses(addrs map[cipher.Address]struct{}) Richlist {
	var s Richlist
	for _, b := range r {
		if _, ok := addrs[b.Address]; !ok {
			s = append(s, b)
		}
	}
	return s
}
