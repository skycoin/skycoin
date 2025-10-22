package wasmgo

import _ "embed"

//go:embed skycoin-lite.wasm
var WasmFile []byte

//go:embed wasm_exec.js
var WasmExecJS []byte
