//go:build !tinygo

package commands

// net/http/pprof registers debug profiling handlers on the default
// http.ServeMux as a side effect. TinyGo's net/http does not provide the
// full server APIs that pprof depends on, so it is only registered for
// standard builds.
import _ "net/http/pprof"
