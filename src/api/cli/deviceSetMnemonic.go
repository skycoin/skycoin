package cli

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	messages "github.com/skycoin/skycoin/protob"
	hardwareWallet "github.com/skycoin/skycoin/src/hardware-wallet"
	"github.com/wire"
	gcli "github.com/urfave/cli"
)

func deviceSetMnemonicCmd() gcli.Command {
	name := "deviceSetMnemonic"
	return gcli.Command{
		Name:        name,
		Usage:       "Configure the device with a mnemonic.",
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
			dev, _ := hardwareWallet.GetTrezorDevice()

			skycoinSetMnemonic := &messages.SetMnemonic{
				Mnemonic:     proto.String(mnemonic),
			}

			data, _ := proto.Marshal(skycoinSetMnemonic)

			chunks := hardwareWallet.MakeTrezorMessage(data, messages.MessageType_MessageType_SetMnemonic)


			for _, element := range chunks {
				_, _ = dev.Write(element[:])
			}

			var msg wire.Message
			msg.ReadFrom(dev)

			fmt.Printf("Success %d! address that issued the signature is: %s\n", msg.Kind, msg.Data)
		},
	}
}
