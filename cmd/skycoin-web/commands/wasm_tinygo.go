//go:build tinygo

package commands

import wasmtinygo "github.com/skycoin/skycoin/src/skycoin-lite/wasm-tinygo"

// See wasm_go.go. TinyGo builds embed the TinyGo-compiled skycoin-lite wasm and
// its matching wasm_exec.js.
var (
	wasmFile   = wasmtinygo.WasmFile
	wasmExecJS = wasmtinygo.WasmExecJS
)
