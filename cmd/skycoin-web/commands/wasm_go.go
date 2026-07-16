//go:build !tinygo

package commands

import wasmgo "github.com/skycoin/skycoin/src/skycoin-lite/wasm-go"

// wasmFile and wasmExecJS are the skycoin-lite WebAssembly module and its
// JavaScript loader served by the web wallet. Standard-toolchain builds embed
// the Go-compiled wasm; TinyGo builds embed the TinyGo-compiled wasm (see
// wasm_tinygo.go). Selecting one via build tags keeps a single wasm in each
// binary and pairs it with the matching wasm_exec.js (which differs between the
// two toolchains).
var (
	wasmFile   = wasmgo.WasmFile
	wasmExecJS = wasmgo.WasmExecJS
)
