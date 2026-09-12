//go:build js && !tinygo

package coloredcobra

import (
	"os"

	"github.com/fatih/color"
)

// Color is on in a browser terminal. fatih/color's isatty probe sees a pipe —
// under js/wasm stdout always is one — but the terminal rendering the other
// end draws ANSI just fine. Only the host knows, and it says so by exporting
// TERM (a browser shell sets TERM=xterm-256color for the instances it runs).
// NO_COLOR and TERM=dumb still win, and an environment that sets no TERM (a
// test harness capturing output) stays plain.
func init() {
	if os.Getenv("NO_COLOR") != "" {
		return
	}
	if t := os.Getenv("TERM"); t != "" && t != "dumb" {
		color.NoColor = false
	}
}
