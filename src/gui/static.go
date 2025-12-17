// Package gui pkg/gui/static.go
package gui

import (
	"embed"
)

// GuiFiles includes the embedded gui sources
//
//go:embed static/*
var GuiFiles embed.FS
