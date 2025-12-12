// Package templates provides embedded template files for fibercoin generation
package templates

import (
	"embed"
)

//go:embed *.template
// FS embeds all template files for fibercoin generation
var FS embed.FS

// Template file names
const (
	CoinTemplate     = "coin.template"
	CoinTestTemplate = "coin_test.template"
	CommandTemplate  = "command.template"
	ParamsTemplate   = "params.template"
)
