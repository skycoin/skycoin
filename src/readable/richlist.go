package readable

import (
	"github.com/skycoin/skycoin/src/util/droplet"
	"github.com/skycoin/skycoin/src/visor"
)

// RichlistBalance holds info an address balance holder
type RichlistBalance struct {
	Address string `json:"address"`
	Coins   string `json:"coins"`
	Locked  bool   `json:"locked"`
}

// NewRichlistBalances copies from visor.Richlist
func NewRichlistBalances(visorRichlist visor.Richlist) ([]RichlistBalance, error) {
	richlist := make([]RichlistBalance, len(visorRichlist))
	for i, v := range visorRichlist {
		coins, err := droplet.ToString(v.Coins)
		if err != nil {
			return nil, err
		}

		richlist[i] = RichlistBalance{
			Address: v.Address.String(),
			Coins:   coins,
			Locked:  v.Locked,
		}
	}

	return richlist, nil
}
