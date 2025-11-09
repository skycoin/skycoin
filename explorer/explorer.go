// Package explorer provides explorer functionality.
package explorer

import (
	"embed"
)

// EmbeddedFiles is the embedded gui source
//
//go:embed dist/*
var EmbeddedFiles embed.FS
