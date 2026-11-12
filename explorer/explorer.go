// Package explorer explorer/explorer.go c0-ui-explorer
//
// Package explorer provides explorer functionality.
package explorer

import (
	"embed"
)

// EmbeddedFiles is the embedded gui source
//
//go:embed dist/*
var EmbeddedFiles embed.FS
