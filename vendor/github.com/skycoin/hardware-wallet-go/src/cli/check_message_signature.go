package cli

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
	skyWallet "github.com/skycoin/hardware-wallet-go/src/skywallet"
)

func init() {
	checkMessageSignatureCmd.Flags().StringVar(&message, "message", "", "The message that the signature claims to be signing.")
	checkMessageSignatureCmd.Flags().StringVar(&signature, "signature", "", "Signature of the message.")
	checkMessageSignatureCmd.Flags().StringVar(&address, "address", "", "Address to verify against the signature.")
	checkMessageSignatureCmd.Flags().StringVar(&deviceType, "deviceType", "USB", "Device type to send instructions to, hardware wallet (USB) or emulator.")
}

var checkMessageSignatureCmd = &cobra.Command{
		Use:   "checkMessageSignature",
		Short: "Check a message signature matches the given address.",
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

			msg, err := device.CheckMessageSignature(message, signature, address)
			if err != nil {
				return err
			}

			responseMsg, err := skyWallet.DecodeSuccessOrFailMsg(msg)
			if err != nil {
				return err
			}

			fmt.Println(responseMsg)
			return nil
		},
	}
