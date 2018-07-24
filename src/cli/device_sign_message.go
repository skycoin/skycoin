package cli

import (
	"fmt"

	gcli "github.com/urfave/cli"

	deviceWallet "github.com/skycoin/skycoin/src/device-wallet"
)

func deviceSignMessageCmd() gcli.Command {
	name := "deviceSignMessage"
	return gcli.Command{
		Name:        name,
		Usage:       "Ask the device to sign a message using the secret key at given index.",
		Description: "",
		Flags: []gcli.Flag{
			gcli.IntFlag{
				Name:  "addressN",
				Value: 0,
				Usage: "Index of the address that will issue the signature. Assume 0 if not set.",
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
			kind, data := deviceWallet.DeviceSignMessage(deviceWallet.DeviceTypeUsb, addressN, message)
			fmt.Printf("Success %d! address that issued the signature is: %s\n", kind, data[2:])
		},
	}
}
