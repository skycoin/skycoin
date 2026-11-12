//go:build tinygo

// Package wasmassets cmd/skycoin-web/wasmassets/wasmassets_tinygo.go c4-ui-wallet
package wasmassets

import (
	wasmtinygo "github.com/skycoin/skycoin/src/skycoin-lite/wasm-tinygo"
	"github.com/skycoin/skycoin/src/skycoin-web/src/gui"
)

func init() {
	gui.RegisterCipherWasm(wasmtinygo.WasmFileGz, wasmtinygo.WasmExecJS)
}
