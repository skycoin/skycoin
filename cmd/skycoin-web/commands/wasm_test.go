package commands

import (
	"os/exec"
	"strings"
	"testing"
)

const cipherWasmBlob = "github.com/skycoin/skycoin/src/skycoin-lite/wasm-go"

// deps returns the dependency graph of a package pattern, as a set.
func deps(t *testing.T, pkg string) map[string]bool {
	t.Helper()

	// pkg is never attacker controlled: both call sites pass a literal, either
	// "." or one of the four in-repo package paths listed below.
	out, err := exec.Command("go", "list", "-deps", pkg).Output() //nolint:gosec // G204: package pattern is a literal, not input
	if err != nil {
		t.Skipf("go list unavailable: %v", err)
	}

	set := make(map[string]bool)
	for _, line := range strings.Split(string(out), "\n") {
		set[strings.TrimSpace(line)] = true
	}

	return set
}

// TestCipherWasmIsNotLinkedByTheCommand pins why wasm.go registers the cipher
// wasm rather than embedding it.
//
// Embedding it here made the blob reachable from the handlers, so every importer
// of this command linked ~1.8 MB of gzipped wasm whether or not it wanted one.
// skywire imports this command and ships the wasm visor, which publishes the
// same Cipher and CipherExtras globals — so it was carrying the cipher twice.
//
// Re-adding the embed would compile, serve identically for skycoin, and pass
// every other test in this package. Only the dependency graph shows it.
func TestCipherWasmIsNotLinkedByTheCommand(t *testing.T) {
	if deps(t, ".")[cipherWasmBlob] {
		t.Errorf("%s is reachable from this command; it is registered by "+
			"cmd/skycoin-web/wasmassets precisely so importers that do not want "+
			"it do not link it", cipherWasmBlob)
	}
}

// TestCipherWasmIsLinkedBySkycoinsOwnAssembly guards the opposite regression.
//
// Making the wasm optional means it can go missing. If cmd/skycoin-wallet/commands
// were to drop its wasmassets import, `skycoin web` would still build and still
// serve the wallet — the two asset routes would just 404, and the wallet would
// fail in the browser with no cipher. Nothing in a Go build catches that, so it
// is asserted here.
func TestCipherWasmIsLinkedBySkycoinsOwnAssembly(t *testing.T) {
	// Every skycoin build that serves the wallet. cmd/skycoin-wallet/commands
	// covers the skycoin binary and the duplicate root skycoin.go; the other
	// three are the standalone thin client, its src/ twin, and the combined
	// release build.
	for _, pkg := range []string{
		"github.com/skycoin/skycoin/cmd/skycoin-wallet/commands",
		"github.com/skycoin/skycoin/cmd/skycoin-web",
		"github.com/skycoin/skycoin/src/skycoin-web",
		"github.com/skycoin/skycoin/cmd/release/commands",
	} {
		if !deps(t, pkg)[cipherWasmBlob] {
			t.Errorf("%s does not link %s; `skycoin web` would 404 both "+
				"/assets/scripts routes and the wallet would have no cipher",
				pkg, cipherWasmBlob)
		}
	}
}

func TestRegisterCipherWasm(t *testing.T) {
	origGz, origJS := wasmFileGz, wasmExecJS
	defer func() { wasmFileGz, wasmExecJS = origGz, origJS }()

	wasmFileGz, wasmExecJS = nil, nil
	if cipherWasmRegistered() {
		t.Error("no cipher wasm registered, but cipherWasmRegistered() is true")
	}

	// Both halves are required: serving a wasm without the matching loader, or a
	// loader with no wasm, fails in the browser rather than at the route.
	RegisterCipherWasm([]byte("gz"), nil)
	if cipherWasmRegistered() {
		t.Error("only the wasm was registered, but cipherWasmRegistered() is true")
	}

	RegisterCipherWasm(nil, []byte("js"))
	if cipherWasmRegistered() {
		t.Error("only the loader was registered, but cipherWasmRegistered() is true")
	}

	RegisterCipherWasm([]byte("gz"), []byte("js"))
	if !cipherWasmRegistered() {
		t.Error("both halves registered, but cipherWasmRegistered() is false")
	}
}
