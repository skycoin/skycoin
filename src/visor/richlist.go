package visor

import (
	"sort"
	"strings"

	"github.com/skycoin/skycoin/src/util/droplet"
)

// RichlistBalance holds info an address balance holder
type RichlistBalance struct {
	Address string `json:"address"`
	Coins   string `json:"coins"`
	Locked  bool   `json:"locked"`
	coins   uint64
}

// Richlist contains RichlistBalances
type Richlist []RichlistBalance

// NewRichlist create Richlist via unspent outputs map
func NewRichlist(allAccounts map[string]uint64, lockedAddrs map[string]struct{}) (Richlist, error) {
	richlist := make(Richlist, 0, len(allAccounts))

	for addr, coins := range allAccounts {
		var islocked bool
		if _, ok := lockedAddrs[addr]; ok {
			islocked = true
		}

		coinsStr, err := droplet.ToString(coins)
		if err != nil {
			return nil, err
		}

		richlist = append(richlist, RichlistBalance{
			Address: addr,
			Coins:   coinsStr,
			coins:   coins,
			Locked:  islocked,
		})
	}

	// Sort order:
	// Higher coins
	// Locked > unlocked
	// Address alphabetical
	sort.Slice(richlist, func(i, j int) bool {
		if richlist[i].coins == richlist[j].coins {
			if richlist[i].Locked == richlist[j].Locked {
				return strings.Compare(richlist[i].Address, richlist[j].Address) < 0
			}
			return richlist[i].Locked
		}

		return richlist[i].coins > richlist[j].coins
	})

	return richlist, nil
}

// FilterAddresses returns the richlist without addresses from the map
func (r Richlist) FilterAddresses(addrs map[string]struct{}) Richlist {
	var s Richlist
	for _, b := range r {
		if _, ok := addrs[b.Address]; !ok {
			s = append(s, b)
		}
	}
	return s
}
