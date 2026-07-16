//go:build !tinygo

package commands

// net/http/pprof registers debug profiling handlers on the default
// http.ServeMux as a side effect. TinyGo's net/http does not provide the
// full server APIs that pprof depends on, so it is only registered for
// standard builds.
import _ "net/http/pprof" //nolint:gosec // G108: pprof endpoints only reachable when --pprofmode=http is set
