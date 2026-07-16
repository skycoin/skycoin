//go:build tinygo && !linux

package coloredcobra

// isTerminal conservatively reports false on non-Linux TinyGo targets, so colors
// are only emitted where terminal detection is known to work. The colored help
// renderer is intended for the native (host) build.
func isTerminal(_ uintptr) bool { return false }
