// Package wasmtinygo provides WebAssembly binaries compiled with TinyGo.
package wasmtinygo

import _ "embed"

//go:embed skycoin-lite.wasm
// WasmFile contains the compiled WebAssembly binary
var WasmFile []byte

//go:embed wasm_exec.js
// WasmExecJS contains the JavaScript WebAssembly executor
var WasmExecJS []byte
