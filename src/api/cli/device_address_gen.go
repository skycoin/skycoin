package cli

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	messages "github.com/skycoin/skycoin/protob"
	hardwareWallet "github.com/skycoin/skycoin/src/hardware-wallet"
	"github.com/wire"
	gcli "github.com/urfave/cli"
)

func deviceAddressGenCmd() gcli.Command {
	name := "deviceAddressGen"
	return gcli.Command{
		Name:        name,
		Usage:       "Generate skycoin or bitcoin addresses using the firmware",
		Description: "",
		Flags: []gcli.Flag{
			gcli.IntFlag{
				Name:  "addressN",
				Value: 1,
				Usage: "Index for deterministic key generation. Assume 1 if not set.",
			},
			gcli.BoolFlag{
				Name:  "bitcoin,b",
				Usage: "Output the addresses as bitcoin addresses instead of skycoin addresses",
			},
		},
		OnUsageError: onCommandUsageError(name),
		Action: func(c *gcli.Context) {
			var coinType messages.SkycoinAddressType
			if c.Bool("bitcoin") {
				coinType = messages.SkycoinAddressType_AddressTypeBitcoin
			} else {
				coinType = messages.SkycoinAddressType_AddressTypeSkycoin
			}

			addressN := c.Int("addressN")
			dev, _ := hardwareWallet.GetTrezorDevice()

			skycoinAddress := &messages.SkycoinAddress{
				AddressN:    proto.Uint32(uint32(addressN)),
				AddressType: coinType.Enum(),
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
