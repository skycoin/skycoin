package cli

import (
	"fmt"
	"os"
	"runtime"

	messages "github.com/skycoin/hardware-wallet-protob/go"

	"github.com/spf13/cobra"

	skyWallet "github.com/skycoin/hardware-wallet-go/src/skywallet"
)

func init() {
		addressGenCmd.Flags().IntVar(&addressN, "addressN", 1, "Number of addresses to generate. Assume 1 if not set.")
		addressGenCmd.Flags().IntVar(&startIndex, "startIndex", 0, "Index where deterministic key generation will start from. Assume 0 if not set.")
		addressGenCmd.Flags().BoolVar(&confirmAddress, "confirmAddress", false, "If requesting one address it will be sent only if user confirms operation by pressing device's button.")
		addressGenCmd.Flags().StringVar(&deviceType, "deviceType", "USB", "Device type to send instructions to, hardware wallet (USB) or emulator.")
		addressGenCmd.Flags().StringVar(&coinTypeStr, "coinTypeStr", "SKY", "Coin type to use on hardware-wallet.")

}
var addressGenCmd = &cobra.Command{
		Use:   "addressGen",
		Short: "Generate skycoin addresses using the firmware",
		RunE: func(_ *cobra.Command, _ []string) error {
			coinType, err := skyWallet.CoinTypeFromString(coinTypeStr)
			if err != nil {
				return err
			}

			device := skyWallet.NewDevice(skyWallet.DeviceTypeFromString(deviceType))
			if device == nil {
				return fmt.Errorf("failed to create device")
			}
			defer device.Close()

			if os.Getenv("AUTO_PRESS_BUTTONS") == "1" && device.Driver.DeviceType() == skyWallet.DeviceTypeEmulator && runtime.GOOS == "linux" {
				err := device.SetAutoPressButton(true, skyWallet.ButtonRight)
				if err != nil {
					return err
				}
			}

			var pinEnc string
			msg, err := device.AddressGen(uint32(addressN), uint32(startIndex), confirmAddress, coinType)
			if err != nil {
				return err
			}

			for msg.Kind != uint16(messages.MessageType_MessageType_ResponseSkycoinAddress) && msg.Kind != uint16(messages.MessageType_MessageType_Failure) {
				if msg.Kind == uint16(messages.MessageType_MessageType_PinMatrixRequest) {
					fmt.Printf("PinMatrixRequest response: ")
					fmt.Scanln(&pinEnc)
					pinAckResponse, err := device.PinMatrixAck(pinEnc)
					if err != nil {
						return err
					}
					log.Infof("PinMatrixAck response: %s", pinAckResponse)
					continue
				}

				if msg.Kind == uint16(messages.MessageType_MessageType_PassphraseRequest) {
					var passphrase string
					fmt.Printf("Input passphrase: ")
					fmt.Scanln(&passphrase)
					passphraseAckResponse, err := device.PassphraseAck(passphrase)
					if err != nil {
						return err
					}
					log.Infof("PinMatrixAck response: %s", passphraseAckResponse)
					continue
				}

				if msg.Kind == uint16(messages.MessageType_MessageType_ButtonRequest) {
					msg, err = device.ButtonAck()
					if err != nil {
						return err
					}
					continue
				}
			}

			if msg.Kind == uint16(messages.MessageType_MessageType_ResponseSkycoinAddress) {
				addresses, err := skyWallet.DecodeResponseSkycoinAddress(msg)
				if err != nil {
					return err
				}
				fmt.Println(addresses)
			} else {
				failMsg, err := skyWallet.DecodeFailMsg(msg)
				if err != nil {
					return err
				}
				return fmt.Errorf("failed with code: %s", failMsg)
			}
			return nil
		},
	}
