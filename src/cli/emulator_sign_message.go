package cli

import (
	"fmt"

	gcli "github.com/urfave/cli"

	deviceWallet "github.com/skycoin/skycoin/src/device-wallet"
)

func emulatorSignMessageCmd() gcli.Command {
	name := "emulatorSignMessage"
	return gcli.Command{
		Name:        name,
		Usage:       "Ask the emulated device to sign a message using the secret key at given index.",
		Description: "",
		Flags: []gcli.Flag{
			gcli.IntFlag{
				Name:  "addressN",
				Value: 1,
				Usage: "Index of the address that will issue the signature. Assume 1 if not set.",
			},
			gcli.StringFlag{
				Name:  "message",
				Usage: "The message that the signature claims to be signing.",
			},
		},
		OnUsageError: onCommandUsageError(name),
		Action: func(c *gcli.Context) {
			addressN := c.Int("addressN")
			message := c.String("message")
			kind, data := deviceWallet.DeviceSignMessage(deviceWallet.DeviceTypeEmulator, addressN, message)
			fmt.Printf("Success %d! address that issued the signature is: %s\n", kind, data[2:])
		},
	}
}
