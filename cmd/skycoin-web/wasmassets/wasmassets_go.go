//go:build !tinygo

// Package wasmassets registers the skycoin-lite cipher wasm with the
// skycoin-web command, for the side effect of importing it:
//
//	import _ "github.com/skycoin/skycoin/cmd/skycoin-web/wasmassets"
//
// It exists so the blob is linked only by binaries that want to serve it. The
// command package holds the assets in registered variables rather than embedding
// them (see cmd/skycoin-web/commands/wasm.go), which means importing the command
// no longer drags ~1.8 MB of gzipped wasm along with it. skycoin's own binaries
// import this; a host that already ships a wasm publishing the same Cipher and
// CipherExtras globals — skywire, via the wasm visor — registers that one
// instead and never links this package.
//
// Standard-toolchain builds provide the Go-compiled wasm; TinyGo builds provide
// the TinyGo-compiled one (see wasmassets_tinygo.go). Selecting by build tag
// keeps a single wasm in each binary, paired with the wasm_exec.js of the
// toolchain that produced it.
package wasmassets

import (
	"github.com/skycoin/skycoin/cmd/skycoin-web/commands"
	wasmgo "github.com/skycoin/skycoin/src/skycoin-lite/wasm-go"
)

func init() {
	commands.RegisterCipherWasm(wasmgo.WasmFileGz, wasmgo.WasmExecJS)
}
