package gui

import (
	"bytes"
	"compress/gzip"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRegisterCipherWasm(t *testing.T) {
	origGz, origJS := cipherWasmGz, cipherExecJS
	defer RegisterCipherWasm(origGz, origJS)

	RegisterCipherWasm(nil, nil)
	if CipherWasmRegistered() {
		t.Error("no cipher wasm registered, but CipherWasmRegistered() is true")
	}

	// Both halves are required: serving a wasm without the matching loader, or a
	// loader with no wasm, fails in the browser rather than at the route.
	RegisterCipherWasm([]byte("gz"), nil)
	if CipherWasmRegistered() {
		t.Error("only the wasm was registered, but CipherWasmRegistered() is true")
	}

	RegisterCipherWasm(nil, []byte("js"))
	if CipherWasmRegistered() {
		t.Error("only the loader was registered, but CipherWasmRegistered() is true")
	}

	RegisterCipherWasm([]byte("gz"), []byte("js"))
	if !CipherWasmRegistered() {
		t.Error("both halves registered, but CipherWasmRegistered() is false")
	}
}

func get(t *testing.T, h http.Handler, path string, gzipOK bool) *httptest.ResponseRecorder {
	t.Helper()
	r := httptest.NewRequest(http.MethodGet, path, nil)
	if gzipOK {
		r.Header.Set("Accept-Encoding", "gzip")
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w
}

func TestHandler(t *testing.T) {
	origGz, origJS := cipherWasmGz, cipherExecJS
	defer RegisterCipherWasm(origGz, origJS)

	h, err := Handler()
	if err != nil {
		t.Fatal(err)
	}

	w := get(t, h, "/", false)
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), "<app-root>") {
		t.Fatalf("/ returned %d, want the wallet's index.html", w.Code)
	}

	w = get(t, h, "/manifest.webmanifest", false)
	if w.Code != http.StatusOK || w.Header().Get("Content-Type") != "application/manifest+json" {
		t.Errorf("manifest returned %d with type %q", w.Code, w.Header().Get("Content-Type"))
	}

	RegisterCipherWasm(nil, nil)
	for _, p := range []string{"/assets/scripts/skycoin-lite.wasm", "/assets/scripts/wasm_exec.js"} {
		if w := get(t, h, p, true); w.Code != http.StatusNotFound {
			t.Errorf("%s returned %d with no cipher registered, want 404", p, w.Code)
		}
	}

	var gz bytes.Buffer
	zw := gzip.NewWriter(&gz)
	if _, err := zw.Write([]byte("\x00asm")); err != nil {
		t.Fatal(err)
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	RegisterCipherWasm(gz.Bytes(), []byte("// loader"))

	w = get(t, h, "/assets/scripts/skycoin-lite.wasm", true)
	if w.Header().Get("Content-Encoding") != "gzip" || !bytes.Equal(w.Body.Bytes(), gz.Bytes()) {
		t.Error("a gzip client did not get the wasm as committed")
	}

	w = get(t, h, "/assets/scripts/skycoin-lite.wasm", false)
	if w.Header().Get("Content-Encoding") != "" || w.Body.String() != "\x00asm" {
		t.Error("a client without gzip did not get the decompressed wasm")
	}

	w = get(t, h, "/assets/scripts/wasm_exec.js", false)
	if w.Body.String() != "// loader" {
		t.Errorf("wasm_exec.js returned %q", w.Body.String())
	}
}
