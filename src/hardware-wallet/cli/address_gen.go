// Package cli provides the hardware wallet CLI commands.
package cli

import (
	"fmt"
	"os"
	"runtime"

	messages "github.com/skycoin/hardware-wallet-protob/go"

	"github.com/spf13/cobra"

	skyWallet "github.com/skycoin/skycoin/src/hardware-wallet/skywallet"
)

var addressNewline bool

func init() {
	addressGenCmd.Flags().IntVar(&addressN, "addressN", 1, "Number of addresses to generate. Assume 1 if not set.")
	addressGenCmd.Flags().IntVar(&startIndex, "startIndex", 0, "Index where deterministic key generation will start from. Assume 0 if not set.")
	addressGenCmd.Flags().BoolVar(&confirmAddress, "confirmAddress", false, "If requesting one address it will be sent only if user confirms operation by pressing device's button.")
	addressGenCmd.Flags().StringVar(&deviceType, "deviceType", "USB", "Device type to send instructions to, hardware wallet (USB) or emulator.")
	addressGenCmd.Flags().StringVar(&coinTypeStr, "coinTypeStr", "SKY", "Coin type to use on hardware-wallet.")
	addressGenCmd.Flags().BoolVar(&addressNewline, "newline", false, "Print each address on its own line instead of space-separated.")

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

		// Check that device is in firmware mode (not bootloader)
		if err := requireFirmwareMode(device); err != nil {
			fmt.Printf("Error: %v\n", err)
			fmt.Println("Hint: The device may be in bootloader mode. Flash firmware first, or use --skip to bypass this check")
			return err
		}

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
				_, _ = fmt.Scanln(&pinEnc) //nolint:errcheck // interactive user input
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
				_, _ = fmt.Scanln(&passphrase) //nolint:errcheck // interactive user input
				passphraseAckResponse, err := device.PassphraseAck(passphrase)
				if err != nil {
					return err
				}
				log.Infof("PassphraseAck response: %s", passphraseAckResponse)
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
			// Print addresses - with newline flag, print each on its own line
			if addressNewline {
				for _, addr := range addresses {
					fmt.Println(addr)
				}
			} else {
				fmt.Println(addresses)
			}
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
