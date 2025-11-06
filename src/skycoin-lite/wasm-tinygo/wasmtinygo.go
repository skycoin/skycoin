// Package wasmtinygo provides WebAssembly binaries compiled with TinyGo.
package wasmtinygo

import _ "embed"

// WasmFile contains the compiled WebAssembly binary
//
//go:embed skycoin-lite.wasm
var WasmFile []byte

// WasmExecJS contains the JavaScript WebAssembly executor
//
//go:embed wasm_exec.js
var WasmExecJS []byte
