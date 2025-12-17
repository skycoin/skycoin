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
	generateMnemonicCmd.Flags().BoolVar(&usePassphrase, "usePassphrase", false, "Configure a passphrase")
	generateMnemonicCmd.Flags().IntVar(&wordCount, "wordCount", 12, "Use a specific (12 | 24) number of words for the Mnemonic")
	generateMnemonicCmd.Flags().StringVar(&deviceType, "deviceType", "USB", "Device type to send instructions to, hardware wallet (USB) or emulator.")
}

var generateMnemonicCmd = &cobra.Command{
		Use:   "generateMnemonic",
		Short: "Ask the device to generate a mnemonic and configure itself with it.",
		RunE: func(_ *cobra.Command, _ []string) error {
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

			msg, err := device.GenerateMnemonic(uint32(wordCount), usePassphrase)
			if err != nil {
				return err
			}

			if msg.Kind == uint16(messages.MessageType_MessageType_ButtonRequest) {
				msg, err = device.ButtonAck()
				if err != nil {
					return err
				}
			}

			responseMsg, err := skyWallet.DecodeSuccessOrFailMsg(msg)
			if err != nil {
				return err
			}

			fmt.Println(responseMsg)
			return nil
		},
	}
