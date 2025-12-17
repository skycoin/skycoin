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
	setPinCode.Flags().StringVar(&deviceType, "deviceType", "USB", "Device type to send instructions to, hardware wallet (USB) or emulator.")
}

var setPinCode = &cobra.Command{
		Use:   "setPinCode",
		Short: "Configure a PIN code on a device.",
		RunE: func(cmd *cobra.Command, args []string) error {
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

			removePin := false
			msg, err := device.ChangePin(&removePin)
			if err != nil {
				return err
			}

			for {
				switch msg.Kind {
				case uint16(messages.MessageType_MessageType_Success):
					responseMsg, err := skyWallet.DecodeSuccessOrFailMsg(msg)
					if err != nil {
						return err
					}
					fmt.Println(responseMsg)
					return nil
				case uint16(messages.MessageType_MessageType_Failure):
					failMsg, err := skyWallet.DecodeFailMsg(msg)
					if err != nil {
						return err
					}
					fmt.Println(failMsg)
					return nil
				case uint16(messages.MessageType_MessageType_PinMatrixRequest):
					var pinEnc string
					fmt.Printf("PinMatrixRequest response: ")
					fmt.Scanln(&pinEnc)
					msg, err = device.PinMatrixAck(pinEnc)
					if err != nil {
						return err
					}
				case uint16(messages.MessageType_MessageType_ButtonRequest):
					msg, err = device.ButtonAck()
					if err != nil {
						return err
					}
				}
			}
		},
	}
