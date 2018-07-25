package cli

import (
	"fmt"

	gcli "github.com/urfave/cli"

	proto "github.com/golang/protobuf/proto"
	messages "github.com/skycoin/skycoin/protob"
	deviceWallet "github.com/skycoin/skycoin/src/device-wallet"
)

func emulatorAddressGenCmd() gcli.Command {
	name := "emulatorAddressGen"
	return gcli.Command{
		Name:        name,
		Usage:       "Generate skycoin addresses using an emulated device.",
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
			var data []byte
			var pinEnc string
			kind, addresses := deviceWallet.DeviceAddressGen(deviceWallet.DeviceTypeEmulator, addressN, startIndex)
			if kind == uint16(messages.MessageType_MessageType_PinMatrixRequest) {
				fmt.Printf("PinMatrixRequest response: ")
				fmt.Scanln(&pinEnc)
				kind, data = deviceWallet.DevicePinMatrixAck(deviceWallet.DeviceTypeEmulator, pinEnc)

				if kind == uint16(messages.MessageType_MessageType_ResponseSkycoinAddress) {
					responseSkycoinAddress := &messages.ResponseSkycoinAddress{}
					err := proto.Unmarshal(data, responseSkycoinAddress)
					if err != nil {
						fmt.Printf("unmarshaling error: %s\n", err.Error())
						return
					}
					fmt.Print("Successfully got address")
					fmt.Print(responseSkycoinAddress.GetAddresses())
				}
			} else {
				if kind == uint16(messages.MessageType_MessageType_ResponseSkycoinAddress) {
					fmt.Println("Got addresses without pin code")
				}
				fmt.Print(addresses)
				fmt.Print("\n")
			}
		},
	}
}
