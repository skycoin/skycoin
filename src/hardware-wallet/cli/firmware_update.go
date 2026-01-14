package cli

import (
	"crypto/sha256"
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
	skyWallet "github.com/skycoin/skycoin/src/hardware-wallet/skywallet"
)

func init() {
	firmwareUpdate.Flags().StringVar(&deviceType, "deviceType", "USB", "Device type to send instructions to, hardware wallet (USB) or emulator.")
}

var firmwareUpdate = &cobra.Command{
	Use:   "firmwareUpdate [firmware-file]",
	Short: "Update device's firmware.",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		firmwarePath := args[0]

		// Read firmware file
		fileBytes, err := os.ReadFile(firmwarePath)
		if err != nil {
			fmt.Printf("Error: failed to read firmware file: %v\n", err)
			return err
		}

		// Compute hash of firmware data (skip first 256 bytes which is the header)
		if len(fileBytes) <= 0x100 {
			fmt.Println("Error: firmware file too small (must be > 256 bytes)")
			return fmt.Errorf("firmware file too small")
		}
		hash := sha256.Sum256(fileBytes[0x100:])

		device := skyWallet.NewDevice(skyWallet.DeviceTypeFromString(deviceType))
		if device == nil {
			fmt.Println("Error: failed to create device (is device connected in bootloader mode?)")
			return fmt.Errorf("failed to create device")
		}
		defer device.Close()

		if os.Getenv("AUTO_PRESS_BUTTONS") == "1" && device.Driver.DeviceType() == skyWallet.DeviceTypeEmulator && runtime.GOOS == "linux" {
			err := device.SetAutoPressButton(true, skyWallet.ButtonRight)
			if err != nil {
				fmt.Printf("Error: failed to set auto press button: %v\n", err)
				return err
			}
		}

		fmt.Printf("Uploading firmware: %s (%d bytes)\n", firmwarePath, len(fileBytes))
		err = device.FirmwareUpload(fileBytes, hash)
		if err != nil {
			fmt.Printf("Error: firmware upload failed: %v\n", err)
			return err
		}

		fmt.Println("Firmware uploaded successfully")
		return nil
	},
}
