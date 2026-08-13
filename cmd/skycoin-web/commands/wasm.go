package commands

// The skycoin-lite cipher WebAssembly module and its JavaScript loader, served
// at /assets/scripts/skycoin-lite.wasm and /assets/scripts/wasm_exec.js. The
// wallet loads them to get its Cipher and CipherExtras globals, which is how
// key material stays in the browser.
//
// These are registered rather than embedded here, because whether this command
// should carry a cipher wasm at all depends on what is hosting it:
//
//   - skycoin's own binaries embed it. cmd/skycoin-web and
//     cmd/skycoin-wallet/commands both import cmd/skycoin-web/wasmassets for its
//     side effect, so `skycoin web` serves the blob exactly as before.
//   - a host that already ships a wasm providing the same globals registers
//     that one instead. skywire embeds the wasm visor, whose
//     skycoin-lite/wasmcipher.Register() publishes the identical Cipher and
//     CipherExtras API — so a second copy of the cipher would be ~1.8 MB of
//     duplicate in its binary.
//
// The point of the indirection is the linker. A package-level embed is reachable
// from the handlers below and so is linked into every importer whether or not it
// wants it; a registered one is only linked by an importer that asks for it.
//
// The blob is held gzipped, as committed, and served that way.
var (
	wasmFileGz []byte
	wasmExecJS []byte
)

// RegisterCipherWasm supplies the cipher wasm and the wasm_exec.js loader that
// matches the toolchain which compiled it. The two always travel together —
// TinyGo's loader provides the WASI shims a TinyGo module imports, and Go's does
// not — so they are set in one call and never mixed.
//
// Call it before serve() runs; cmd/skycoin-web/wasmassets does so from an init.
// Without it the two asset routes report 404, which is what a host serving its
// own cipher wants: the routes are absent rather than serving an empty body that
// would fail deep inside WebAssembly.instantiate.
func RegisterCipherWasm(gz, execJS []byte) {
	wasmFileGz = gz
	wasmExecJS = execJS
}

// cipherWasmRegistered reports whether a cipher wasm is available to serve.
func cipherWasmRegistered() bool { return len(wasmFileGz) > 0 && len(wasmExecJS) > 0 }
