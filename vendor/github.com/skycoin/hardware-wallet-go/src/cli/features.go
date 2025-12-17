package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"

	"github.com/gogo/protobuf/proto"
	"github.com/spf13/cobra"

	messages "github.com/skycoin/hardware-wallet-protob/go"

	skyWallet "github.com/skycoin/hardware-wallet-go/src/skywallet"
)

func init() {
	featuresCmd.Flags().StringVar(&deviceType, "deviceType", "USB", "Device type to send instructions to, hardware wallet (USB) or emulator.")
}

var featuresCmd = &cobra.Command{
		Use:   "features",
		Short: "Ask the device Features.",
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

			msg, err := device.GetFeatures()
			if err != nil {
				return err
			}

			switch msg.Kind {
			case uint16(messages.MessageType_MessageType_Features):
				features := &messages.Features{}
				err = proto.Unmarshal(msg.Data, features)
				if err != nil {
					return err
				}

				enc := json.NewEncoder(os.Stdout)
				if err = enc.Encode(features); err != nil {
					return err
				}
				ff := skyWallet.NewFirmwareFeatures(uint64(*features.FirmwareFeatures))
				if err := ff.Unmarshal(); err != nil {
					return err
				}
				log.Printf("\n\nFirmware features:\n%s", ff)
			case uint16(messages.MessageType_MessageType_Failure), uint16(messages.MessageType_MessageType_Success):
				msgData, err := skyWallet.DecodeSuccessOrFailMsg(msg)
				if err != nil {
					return err
				}

				fmt.Println(msgData)
			default:
				return fmt.Errorf("received unexpected message type: %s", messages.MessageType(msg.Kind))
			}
			return nil
		},
}
