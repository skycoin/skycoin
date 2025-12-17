package cli

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
	messages "github.com/skycoin/hardware-wallet-protob/go"
	skyWallet "github.com/skycoin/hardware-wallet-go/src/skywallet"
)

func init() {
	applySettingsCmd.Flags().BoolVar(&usePassphrase, "usePassphrase", false, "Configure a passphrase (true or false)")
	applySettingsCmd.Flags().StringVar(&label, "label", "", "Configure a device label")
	applySettingsCmd.Flags().StringVar(&deviceType, "deviceType", "USB", "Device type to send instructions to, hardware wallet (USB) or emulator.")
	applySettingsCmd.Flags().StringVar(&language, "language", "", "Configure a device language")
}

var applySettingsCmd = &cobra.Command{
		Use:   "applySettings",
		Short: "Apply settings.",
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

			msg, err := device.ApplySettings(&usePassphrase, label, language)
			if err != nil {
				return err
			}

			for msg.Kind != uint16(messages.MessageType_MessageType_Failure) && msg.Kind != uint16(messages.MessageType_MessageType_Success) {
				if msg.Kind == uint16(messages.MessageType_MessageType_ButtonRequest) {
					msg, err = device.ButtonAck()
					if err != nil {
						return err
					}
					continue
				}

				if msg.Kind == uint16(messages.MessageType_MessageType_PinMatrixRequest) {
					var pinEnc string
					fmt.Printf("PinMatrixRequest response: ")
					fmt.Scanln(&pinEnc)
					pinAckResponse, err := device.PinMatrixAck(pinEnc)
					if err != nil {
						return err
					}
					log.Infof("PinMatrixAck response: %s", pinAckResponse)
					continue
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
