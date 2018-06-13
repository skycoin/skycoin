package cli

import (
	emulatorWallet "github.com/skycoin/skycoin/src/emulator-wallet"
	gcli "github.com/urfave/cli"
)

func emulatorSetMnemonicCmd() gcli.Command {
	name := "emulatorSetMnemonic"
	return gcli.Command{
		Name:        name,
		Usage:       "Configure an emulated device with a mnemonic.",
		Description: "",
		Flags: []gcli.Flag{
			gcli.StringFlag{
				Name:  "mnemonic",
				Usage: "Mnemonic that will be stored in the device to generate addresses.",
			},
		},
		OnUsageError: onCommandUsageError(name),
		Action: func(c *gcli.Context) {
			mnemonic := c.String("mnemonic")
			emulatorWallet.DeviceSetMnemonic(mnemonic)
		},
	}
}
