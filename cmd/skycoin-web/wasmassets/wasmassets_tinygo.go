//go:build tinygo

package wasmassets

import (
	wasmtinygo "github.com/skycoin/skycoin/src/skycoin-lite/wasm-tinygo"
	"github.com/skycoin/skycoin/src/skycoin-web/src/gui"
)

func init() {
	gui.RegisterCipherWasm(wasmtinygo.WasmFileGz, wasmtinygo.WasmExecJS)
}
