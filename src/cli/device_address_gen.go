package cli

import (
	"fmt"

	gcli "github.com/urfave/cli"

	deviceWallet "github.com/skycoin/skycoin/src/device-wallet"
)

func deviceAddressGenCmd() gcli.Command {
	name := "deviceAddressGen"
	return gcli.Command{
		Name:        name,
		Usage:       "Generate skycoin addresses using the firmware",
		Description: "",
		Flags: []gcli.Flag{
			gcli.IntFlag{
				Name:  "addressN",
				Value: 1,
				Usage: "Number of addresses to generate. Assume 1 if not set.",
			},
			gcli.IntFlag{
				Name:  "startIndex",
				Value: 0,
				Usage: "Index where deterministic key generation will start from. Assume 0 if not set.",
			},
		},
		OnUsageError: onCommandUsageError(name),
		Action: func(c *gcli.Context) {
			addressN := c.Int("addressN")
			startIndex := c.Int("startIndex")
			kind, responseSkycoinAddress := deviceWallet.DeviceAddressGen(deviceWallet.DeviceTypeUsb, addressN, startIndex)
			fmt.Printf("MessageSkycoinAddress %d! array size is %d\n", kind, len(responseSkycoinAddress))
			for i := 0; i < len(responseSkycoinAddress); i++ {
				fmt.Printf("MessageSkycoinAddress %d! Answer is: %s\n", kind, responseSkycoinAddress[i])
			}
		},
	}
}
