// Package templates provides embedded template files for fibercoin generation
package templates

import (
	"embed"
)

// FS embeds all template files for fibercoin generation
//
//go:embed *.template
var FS embed.FS

// Template file names
const (
	CoinTemplate     = "coin.template"
	CoinTestTemplate = "coin_test.template"
	CommandTemplate  = "command.template"
	ParamsTemplate   = "params.template"
)
