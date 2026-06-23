// Package gui pkg/gui/static.go
package gui

import (
	"embed"
)

// GuiFiles includes the embedded gui sources.
//
// Only the built output under static/dist is embedded, since that is the only
// subtree served (see fs.Sub(GuiFiles, "static/dist") in src/api/http.go).
// Embedding static/* instead pulls in the entire front-end working tree (src,
// e2e, configs and, when present locally, node_modules) — tens of thousands of
// files that bloat the binary and, under TinyGo, blow past the linker's
// argument limit because TinyGo emits one object file per embedded file.
//
//go:embed static/dist
var GuiFiles embed.FS
