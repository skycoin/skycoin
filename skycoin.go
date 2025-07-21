// Package skycoin github.com/skycoin/skycoin/skycoin.go
//
package skycoin

import (
	"embed"
)

/*
Embedded Files
*/

// GuiFiles is the embedded gui sources
//
//go:embed src/gui/static/*
var GuiFiles embed.FS


/*
//TODO: embed files for use with newcoin

// FiberToml is the embedded fiber.toml default node configuraion file
//go:embed fiber.toml
//var FiberToml []byte

// CoinTemplate is embedded template/coin.template
//go:embed template/coin.template
//var CoinTemplate []byte

// CoinTemplate is embedded template/coin_test.template
//go:embed template/coin_test.template
//var CoinTestTemplate []byte

// CommandTemplate is embedded template/command.template
//go:embed template/command.template
//var CommandTemplate []byte

// ParamsTemplate is embedded template/params.template
//go:embed template/params.template
//var ParamsTemplate []byte

//TODO: embed default peers

// PeersTxt is the embedded fiber.toml default node configuraion file
//go:embed peers.txt
//var PeersTxt []byte
*/
