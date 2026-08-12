//go:build wasm

// Package wasmcipher publishes the skycoin browser cipher on the JavaScript
// global object.
//
// It lives in a package rather than in a main so that any wasm program can carry
// the cipher, not just the one built from src/skycoin-lite/wasm. The wallet
// front end reaches the cipher only through window.SkycoinCipher — see
// cipher.provider.ts in src/skycoin-web — so a larger wasm that calls Register
// is a drop-in for skycoin-lite.wasm, and the page cannot tell the difference.
// That is what would let a single wasm serve both a visor and the wallet.
//
// Register does not block. The caller owns the program's lifetime, because a
// host that carries the cipher alongside other things has its own reason to
// stay alive.
package wasmcipher

import (
	"fmt"
	"runtime/debug"
	"syscall/js"

	"github.com/skycoin/skycoin/src/skycoin-lite/liteclient"
)

// buildVersion reports what this wasm was built from.
//
// The toolchain already records it: a `go build` inside a git work tree stamps
// the module version and the vcs.* settings into the binary, so nothing needs
// injecting through -ldflags. Reading it back is the only way to tell, from a
// running wallet, which cipher it is actually holding — the wasm is a committed
// artifact, so the commit it was built at is necessarily an earlier one than the
// commit that carries it.
//
// modified is the interesting field. It means the working tree had uncommitted
// changes when this was compiled, so no commit describes what is in it.
//
// version describes the wasm program that called Register, which is what
// identifies the module the browser actually loaded. When that program is not
// skycoin — a visor carrying the cipher — cipherVersion additionally reports the
// skycoin version it was built against, read from the dependency list the way
// `skywire -d | grep skycoin` shows it. Both matter: one says which binary is
// running, the other which cipher is inside it.
//
// The upstream TinyGo does not record any of this and reports empty strings;
// github.com/0magnet/tinygo does.
// skycoinModulePath is this module, used to find the cipher's own version when
// the wasm was built from a different one.
const skycoinModulePath = "github.com/skycoin/skycoin"

func buildVersion() map[string]interface{} {
	version := map[string]interface{}{
		"version":       "",
		"commit":        "",
		"date":          "",
		"modified":      false,
		"cipherVersion": "",
	}

	info, ok := debug.ReadBuildInfo()
	if !ok {
		return version
	}

	version["version"] = info.Main.Version

	// The cipher's own version, when the host is something else.
	if info.Main.Path == skycoinModulePath || info.Main.Path == "" {
		version["cipherVersion"] = info.Main.Version
	} else {
		for _, dep := range info.Deps {
			if dep.Path == skycoinModulePath {
				version["cipherVersion"] = dep.Version

				break
			}
		}
	}

	for _, setting := range info.Settings {
		switch setting.Key {
		case "vcs.revision":
			version["commit"] = setting.Value
		case "vcs.time":
			version["date"] = setting.Value
		case "vcs.modified":
			version["modified"] = setting.Value == "true"
		}
	}

	return version
}

// errorResult is what every entry point returns when it could not do its job.
// Callers check for the "error" key, so the shape has to be a map even when the
// success value is a plain string.
func errorResult(err interface{}) map[string]interface{} {
	return map[string]interface{}{"error": fmt.Sprint(err)}
}

// guard runs fn and turns a panic into an {"error": ...} result.
//
// liteclient returns errors now, so nothing below is expected to panic and this
// should never fire. It stays as a backstop for the cipher package underneath,
// which still reports some failures that way.
//
// It cannot be relied on. Under TinyGo the recover does not fire at all and a
// panic traps the whole module, leaving the cipher dead for the rest of the
// page — which is why the errors are returned rather than recovered.
func guard(fn func() (interface{}, error)) (result interface{}) {
	defer func() {
		if r := recover(); r != nil {
			result = errorResult(r)
		}
	}()

	value, err := fn()
	if err != nil {
		return errorResult(err)
	}

	return value
}

// arity wraps a function of n string arguments, rejecting short calls before
// they reach the cipher.
func arity(n int, name string, fn func(args []js.Value) (interface{}, error)) js.Func {
	return js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		if len(args) < n {
			return errorResult(fmt.Sprintf("%s requires %d argument(s)", name, n))
		}

		return guard(func() (interface{}, error) { return fn(args) })
	})
}

// Register publishes SkycoinCipher and SkycoinCipherExtras on the JavaScript
// global object. It returns once they are in place.
func Register() {
	// The wallet's own entry points.
	skycoinCipher := js.Global().Get("Object").New()

	skycoinCipher.Set("generateAddress", arity(1, "generateAddress", func(args []js.Value) (interface{}, error) {
		address, err := liteclient.GenerateAddress(args[0].String())
		if err != nil {
			return nil, err
		}

		return map[string]interface{}{
			"nextSeed": address.NextSeed,
			"secret":   address.Secret,
			"public":   address.Public,
			"address":  address.Address,
		}, nil
	}))

	skycoinCipher.Set("prepareTransaction", arity(2, "prepareTransaction", func(args []js.Value) (interface{}, error) {
		return liteclient.PrepareTransaction(args[0].String(), args[1].String())
	}))

	// prepareTransactionWithSignatures is used for hardware wallet signing.
	skycoinCipher.Set("prepareTransactionWithSignatures", arity(3, "prepareTransactionWithSignatures", func(args []js.Value) (interface{}, error) {
		return liteclient.PrepareTransactionWithSignatures(args[0].String(), args[1].String(), args[2].String())
	}))

	skycoinCipher.Set("version", buildVersion())

	js.Global().Set("SkycoinCipher", skycoinCipher)

	// The verification helpers. These exist so the browser cipher can be checked
	// against src/cipher/testsuite/testdata — the same golden vectors the Go
	// implementation is checked against — by
	// src/skycoin-web/src/app/services/cipher.provider.lib.spec.ts. The GopherJS
	// build publishes the same set as CipherExtras.
	//
	// Each verify* returns null when the check passes and an error string when it
	// does not, rather than the int the underlying secp256k1 helpers use.
	skycoinCipherExtras := js.Global().Get("Object").New()

	skycoinCipherExtras.Set("verifyPubKeySignedHash", arity(3, "verifyPubKeySignedHash", func(args []js.Value) (interface{}, error) {
		return nil, liteclient.VerifyPubKeySignedHash(args[0].String(), args[1].String(), args[2].String())
	}))

	skycoinCipherExtras.Set("verifyAddressSignedHash", arity(3, "verifyAddressSignedHash", func(args []js.Value) (interface{}, error) {
		return nil, liteclient.VerifyAddressSignedHash(args[0].String(), args[1].String(), args[2].String())
	}))

	skycoinCipherExtras.Set("verifySignatureRecoverPubKey", arity(2, "verifySignatureRecoverPubKey", func(args []js.Value) (interface{}, error) {
		return nil, liteclient.VerifySignatureRecoverPubKey(args[0].String(), args[1].String())
	}))

	skycoinCipherExtras.Set("verifySeckey", arity(1, "verifySeckey", func(args []js.Value) (interface{}, error) {
		return nil, liteclient.VerifySeckey(args[0].String())
	}))

	skycoinCipherExtras.Set("verifyPubkey", arity(1, "verifyPubkey", func(args []js.Value) (interface{}, error) {
		return nil, liteclient.VerifyPubkey(args[0].String())
	}))

	skycoinCipherExtras.Set("addressFromPubKey", arity(1, "addressFromPubKey", func(args []js.Value) (interface{}, error) {
		return liteclient.AddressFromPubKey(args[0].String())
	}))

	skycoinCipherExtras.Set("addressFromSecKey", arity(1, "addressFromSecKey", func(args []js.Value) (interface{}, error) {
		return liteclient.AddressFromSecKey(args[0].String())
	}))

	skycoinCipherExtras.Set("pubKeyFromSig", arity(2, "pubKeyFromSig", func(args []js.Value) (interface{}, error) {
		return liteclient.PubKeyFromSig(args[0].String(), args[1].String())
	}))

	skycoinCipherExtras.Set("signHash", arity(2, "signHash", func(args []js.Value) (interface{}, error) {
		return liteclient.SignHash(args[0].String(), args[1].String())
	}))

	js.Global().Set("SkycoinCipherExtras", skycoinCipherExtras)
}
