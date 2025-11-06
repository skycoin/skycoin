// Package gui provides embedded web UI assets.
package gui

import (
	"embed"
)

//go:embed all:dist
// DistFS contains the embedded web UI distribution files
var DistFS embed.FS
