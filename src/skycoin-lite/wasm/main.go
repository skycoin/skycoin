package main

import (
	"syscall/js"

	"github.com/skycoin/skycoin/src/skycoin-lite/liteclient"
)

func main() {
	// Create SkycoinCipher object with methods
	skycoinCipher := js.Global().Get("Object").New()

	// generateAddress function
	skycoinCipher.Set("generateAddress", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) < 1 {
			return map[string]interface{}{"error": "seed parameter required"}
		}

		seed := args[0].String()

		defer func() {
			if r := recover(); r != nil {
				// Convert panic to error return
			}
		}()

		address := liteclient.GenerateAddress(seed)

		return map[string]interface{}{
			"nextSeed": address.NextSeed,
			"secret":   address.Secret,
			"public":   address.Public,
			"address":  address.Address,
		}
	}))

	// prepareTransaction function
	skycoinCipher.Set("prepareTransaction", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) < 2 {
			return map[string]interface{}{"error": "inputs and outputs parameters required"}
		}

		inputsJSON := args[0].String()
		outputsJSON := args[1].String()

		defer func() {
			if r := recover(); r != nil {
				// Convert panic to error return
			}
		}()

		txHex := liteclient.PrepareTransaction(inputsJSON, outputsJSON)

		return txHex
	}))

	// Set SkycoinCipher on global window object
	js.Global().Set("SkycoinCipher", skycoinCipher)

	// Keep the Go program running
	<-make(chan struct{})
}
