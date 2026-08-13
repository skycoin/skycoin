//go:build tinygo

package wasmassets

import (
	"github.com/skycoin/skycoin/cmd/skycoin-web/commands"
	wasmtinygo "github.com/skycoin/skycoin/src/skycoin-lite/wasm-tinygo"
)

func init() {
	commands.RegisterCipherWasm(wasmtinygo.WasmFileGz, wasmtinygo.WasmExecJS)
}
