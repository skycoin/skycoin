//go:build !tinygo

package coloredcobra

// installManualHelp is a no-op on toolchains where cobra's templated help works.
// See manualhelp_tinygo.go for the TinyGo implementation.
func installManualHelp(_ *Config) {}
