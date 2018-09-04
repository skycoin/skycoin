package readable

import "github.com/skycoin/skycoin/src/visor"

// RichlistBalance holds info an address balance holder
type RichlistBalance struct {
	Address string `json:"address"`
	Coins   string `json:"coins"`
	Locked  bool   `json:"locked"`
}

// NewRichlistBalances copies from visor.Richlist
func NewRichlistBalances(visorRichlist visor.Richlist) []RichlistBalance {
	richlist := make([]RichlistBalance, len(visorRichlist))
	for i, v := range visorRichlist {
		richlist[i] = RichlistBalance{
			Address: v.Address,
			Coins:   v.Coins,
			Locked:  v.Locked,
		}
	}
	return richlist
}
