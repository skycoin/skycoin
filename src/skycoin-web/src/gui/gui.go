// Package gui provides embedded web UI assets.
package gui

import (
	"embed"
)

// DistFS contains the embedded web UI distribution files
//
//go:embed all:dist
var DistFS embed.FS
