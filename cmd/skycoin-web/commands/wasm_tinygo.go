//go:build tinygo

package commands

import wasmtinygo "github.com/skycoin/skycoin/src/skycoin-lite/wasm-tinygo"

// wasmFileGz and wasmExecJS are the skycoin-lite WebAssembly module and its
// JavaScript loader served by the web wallet. Standard-toolchain builds embed
// the Go-compiled wasm; TinyGo builds embed the TinyGo-compiled wasm (see
// wasm_go.go). Selecting one via build tags keeps a single wasm in each
// binary and pairs it with the matching wasm_exec.js (which differs between the
// two toolchains).
//
// The wasm is held gzipped, as committed, and served that way.
var (
	wasmFileGz = wasmtinygo.WasmFileGz
	wasmExecJS = wasmtinygo.WasmExecJS
)
