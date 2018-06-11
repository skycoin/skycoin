package cli

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	messages "github.com/skycoin/skycoin/protob"
	hardwareWallet "github.com/skycoin/skycoin/src/hardware-wallet"
	"github.com/wire"
	gcli "github.com/urfave/cli"
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
			dev, _ := hardwareWallet.GetTrezorDevice()

			skycoinSignMessage := &messages.SkycoinSignMessage{
				AddressN:    proto.Uint32(uint32(addressN)),
				Message:     proto.String(message),
			}

			data, _ := proto.Marshal(skycoinSignMessage)

			chunks := hardwareWallet.MakeTrezorMessage(data, messages.MessageType_MessageType_SkycoinSignMessage)


			for _, element := range chunks {
				_, _ = dev.Write(element[:])
			}

			var msg wire.Message
			msg.ReadFrom(dev)

			fmt.Printf("Success %d! address that issued the signature is: %s\n", msg.Kind, msg.Data)
		},
	}
}
