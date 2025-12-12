// Package config provides embedded configuration files
package config

import (
	"embed"
)

// FS embeds the fiber.toml configuration file
//
//go:embed fiber.toml
var FS embed.FS

// Config file names
const (
	FiberToml = "fiber.toml"
)
