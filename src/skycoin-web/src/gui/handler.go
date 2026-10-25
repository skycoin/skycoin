package gui

import (
	"io/fs"
	"net/http"
	"path"
	"strings"

	"github.com/skycoin/skycoin/src/skycoin-lite/wasmgz"
)

// The skycoin-lite cipher wasm and the wasm_exec.js loader of the toolchain
// that compiled it. They are registered by the binary that carries them rather
// than embedded here, so a host that ships its own cipher links only that one.
var (
	cipherWasmGz []byte
	cipherExecJS []byte
)

// RegisterCipherWasm supplies the gzipped cipher wasm and its wasm_exec.js.
// Without it the two cipher routes answer 404.
func RegisterCipherWasm(gz, execJS []byte) {
	cipherWasmGz = gz
	cipherExecJS = execJS
}

// CipherWasmRegistered reports whether a cipher wasm is available to serve.
func CipherWasmRegistered() bool { return len(cipherWasmGz) > 0 && len(cipherExecJS) > 0 }

// ServeCipherWasm serves the cipher wasm at /assets/scripts/skycoin-lite.wasm.
// It is handed over gzipped, as committed, to any client that accepts that.
func ServeCipherWasm(w http.ResponseWriter, r *http.Request) {
	if !CipherWasmRegistered() {
		http.Error(w, "no cipher wasm is registered in this build", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/wasm")
	if acceptsGzip(r) {
		w.Header().Set("Content-Encoding", "gzip")
		_, _ = w.Write(cipherWasmGz) //nolint:errcheck
		return
	}

	wasm, err := wasmgz.Decompress(cipherWasmGz)
	if err != nil {
		http.Error(w, "could not decompress the wasm cipher: "+err.Error(), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(wasm) //nolint:errcheck
}

// ServeWasmExecJS serves the cipher's loader at /assets/scripts/wasm_exec.js.
func ServeWasmExecJS(w http.ResponseWriter, _ *http.Request) {
	if !CipherWasmRegistered() {
		http.Error(w, "no cipher wasm is registered in this build", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/javascript")
	_, _ = w.Write(cipherExecJS) //nolint:errcheck
}

// Handler serves the embedded web wallet and its cipher wasm from the site
// root. The wallet talks to the node API on the same origin, so it needs
// nothing else from the server.
func Handler() (http.Handler, error) {
	dist, err := fs.Sub(DistFS, "dist")
	if err != nil {
		return nil, err
	}
	return HandlerFS(dist), nil
}

// HandlerFS serves the web wallet built into dist, such as a local build in
// place of the embedded one.
func HandlerFS(dist fs.FS) http.Handler {
	files := http.FileServer(http.FS(dist))

	mux := http.NewServeMux()
	mux.HandleFunc("/assets/scripts/skycoin-lite.wasm", ServeCipherWasm)
	mux.HandleFunc("/assets/scripts/wasm_exec.js", ServeWasmExecJS)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Go does not know this extension, and browsers want the manifest type.
		if path.Ext(r.URL.Path) == ".webmanifest" {
			w.Header().Set("Content-Type", "application/manifest+json")
		}
		files.ServeHTTP(w, r)
	})

	return mux
}

// acceptsGzip reports whether the client will decompress a gzipped response.
func acceptsGzip(r *http.Request) bool {
	for _, encoding := range strings.Split(r.Header.Get("Accept-Encoding"), ",") {
		if strings.EqualFold(strings.TrimSpace(strings.SplitN(encoding, ";", 2)[0]), "gzip") {
			return true
		}
	}

	return false
}
