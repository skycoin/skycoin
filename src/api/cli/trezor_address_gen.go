package cli

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	messages "github.com/skycoin/skycoin/protob"
	hardwareWallet "github.com/skycoin/skycoin/src/hardware-wallet"
	"github.com/trezor/trezord-go/wire"
	gcli "github.com/urfave/cli"
)

func trezorAddressGenCmd() gcli.Command {
	name := "trezorAddressGen"
	return gcli.Command{
		Name:        name,
		Usage:       "Generate skycoin or bitcoin addresses using the trezor",
		Description: "",
		Flags: []gcli.Flag{
			gcli.StringFlag{
				Name:  "seed",
				Usage: "Seed for deterministic key generation. Will use bip39 as the seed if not provided.",
			},
		},
		OnUsageError: onCommandUsageError(name),
		Action: func(c *gcli.Context) {
			seed := c.String("seed")
			fmt.Printf("Trezor! %s\n", seed)
			dev, _ := hardwareWallet.GetTrezorDevice()

			skycoinAddress := &messages.SkycoinAddress{
				Seed:        proto.String("seed"),
				AddressType: messages.SkycoinAddressType_AddressTypeSkycoin.Enum(),
			}
			data, _ := proto.Marshal(skycoinAddress)

			chunks := hardwareWallet.MakeTrezorMessage(data, messages.MessageType_MessageType_SkycoinAddress)

			for _, element := range chunks {
				_, _ = dev.Write(element[:])
			}

			var msg wire.Message
			msg.ReadFrom(dev)

			fmt.Printf("Success %d! Address is: %s\n", msg.Kind, msg.Data)
		},
	}
}
