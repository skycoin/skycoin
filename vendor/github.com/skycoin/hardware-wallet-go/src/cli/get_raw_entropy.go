package cli

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"

	"github.com/spf13/cobra"
	skyWallet "github.com/skycoin/hardware-wallet-go/src/skywallet"
)

func init() {
	getRawEntropyCmd.Flags().IntVar(&entropyBytes, "entropyBytes", 1048576, "Number of how many bytes of entropy to read.")
	getRawEntropyCmd.Flags().StringVar(&deviceType, "deviceType", "USB", "Device type to send instructions to, hardware wallet (USB) or emulator.")
}

var getRawEntropyCmd = &cobra.Command{
		Use:   "getRawEntropy",
		Short: "Get device raw internal entropy and write it down to a file",
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

			entropy, err := skyWallet.MessageDeviceGetRawEntropy(uint32(entropyBytes))
			if err != nil {
				return err
			}

			var entropyData []byte
			for _, chunk := range entropy {
				entropyData = append(entropyData, chunk[:]...)
			}

			err = ioutil.WriteFile("/tmp/entropy.dump", entropyData, 0644)
			if err != nil {
				return err
			}

			fmt.Println("Raw entropy dumped to: /tmp/entropy.dump")
			return nil
		},
	}
