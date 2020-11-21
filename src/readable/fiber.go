package readable

import "github.com/skycoin/skycoin/src/cipher/bip44"

// FiberConfig is fiber configuration parameters
type FiberConfig struct {
	Name                  string         `json:"name"`
	DisplayName           string         `json:"display_name"`
	Ticker                string         `json:"ticker"`
	CoinHoursName         string         `json:"coin_hours_display_name"`
	CoinHoursNameSingular string         `json:"coin_hours_display_name_singular"`
	CoinHoursTicker       string         `json:"coin_hours_ticker"`
	ExplorerURL           string         `json:"explorer_url"`
	VersionURL            string         `json:"version_url"`
	Bip44Coin             bip44.CoinType `json:"bip44_coin"`
}
