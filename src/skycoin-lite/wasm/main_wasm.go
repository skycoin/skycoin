//go:build wasm

// Package main provides the wasm build of the browser cipher: the entry points
// src/skycoin-web calls to derive addresses and sign transactions, plus the
// verification helpers its cipher spec checks them with.
//
// The registration itself lives in src/skycoin-lite/wasmcipher, so a wasm
// program that does more than this one — a visor serving the wallet, say — can
// carry the same cipher by calling Register. This build is the cipher on its
// own, which is what src/skycoin-web serves today.
package main

import (
	"github.com/skycoin/skycoin/src/skycoin-lite/wasmcipher"
)

func main() {
	wasmcipher.Register()

	// Keep the Go program running: JavaScript calls into the exported functions
	// long after main would otherwise return, and exiting would take them with
	// it. A host with its own reason to stay alive does not need this, which is
	// why Register does not block.
	<-make(chan struct{})
}
