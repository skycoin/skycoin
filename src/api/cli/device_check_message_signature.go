package cli

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	messages "github.com/skycoin/skycoin/protob"
	hardwareWallet "github.com/skycoin/skycoin/src/hardware-wallet"
	"github.com/wire"
	gcli "github.com/urfave/cli"
)

func deviceCheckMessageSignatureCmd() gcli.Command {
	name := "deviceCheckMessageSignature"
	return gcli.Command{
		Name:        name,
		Usage:       "Check a message signature matches the given address.",
		Description: "",
		Flags: []gcli.Flag{
			gcli.StringFlag{
				Name:  "message",
				Usage: "The message that the signature claims to be signing.",
			},
			gcli.StringFlag{
				Name:  "signature",
				Usage: "Signature of the message.",
			},
			gcli.StringFlag{
				Name:  "address",
				Usage: "Address that issued the signature.",
			},
		},
		OnUsageError: onCommandUsageError(name),
		Action: func(c *gcli.Context) {
			message := c.String("message")
			signature := c.String("signature")
			address := c.String("address")
			dev, _ := hardwareWallet.GetTrezorDevice()

			skycoinCheckMessageSignature := &messages.SkycoinCheckMessageSignature{
				Address:   proto.String(address),
				Message:   proto.String(message),
				Signature: proto.String(signature),
			}

			data, _ := proto.Marshal(skycoinCheckMessageSignature)

			chunks := hardwareWallet.MakeTrezorMessage(data, messages.MessageType_MessageType_SkycoinCheckMessageSignature)

			for _, element := range chunks {
				_, _ = dev.Write(element[:])
			}

			var msg wire.Message
			msg.ReadFrom(dev)

			fmt.Printf("Success %d! address that issued the signature is: %s\n", msg.Kind, msg.Data)
		},
	}
}
