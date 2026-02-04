package cli

import (
	"crypto/sha256"
	"fmt"
	"os"
	"runtime"

	"github.com/skycoin/hardware-wallet/firmware"
	"github.com/spf13/cobra"

	skyWallet "github.com/skycoin/skycoin/src/hardware-wallet/skywallet"
)

var (
	firmwareFile     string
	firmwareEmbedded string
	firmwareList     bool
)

func init() {
	firmwareUpdate.Flags().StringVar(&deviceType, "deviceType", "USB", "Device type to send instructions to, hardware wallet (USB) or emulator.")
	firmwareUpdate.Flags().StringVar(&firmwareFile, "file", "", "Path to firmware file to upload")
	firmwareUpdate.Flags().StringVar(&firmwareEmbedded, "embedded", "", "Name of embedded firmware to flash (use --list to see options)")
	firmwareUpdate.Flags().BoolVar(&firmwareList, "list", false, "List available embedded firmwares")
}

var firmwareUpdate = &cobra.Command{
	Use:   "firmwareUpdate [firmware-file]",
	Short: "Update device's firmware.",
	Long: `Update the device's firmware from a file or from embedded firmware binaries.

Use --list to see available embedded firmwares.
Use --embedded <name> to flash an embedded firmware.
Use --file <path> to flash a firmware from a file.

The device must be in bootloader mode (hold buttons while plugging in USB).`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		// Handle --list flag
		if firmwareList {
			fmt.Println("Available embedded firmwares:")
			fmt.Println()
			for _, fw := range firmware.Available() {
				fmt.Printf("  %-12s  v%-8s  %s (%d bytes)\n", fw.Name, fw.Version, fw.Description, len(fw.Data))
			}
			fmt.Println()
			fmt.Println("Use --embedded <name> to flash one of these firmwares")
			return nil
		}

		var fileBytes []byte
		var sourceName string

		// Determine firmware source
		if firmwareEmbedded != "" {
			// Use embedded firmware
			fileBytes = firmware.GetByName(firmwareEmbedded)
			if fileBytes == nil {
				fmt.Printf("Error: unknown embedded firmware '%s'\n", firmwareEmbedded)
				fmt.Println("Use --list to see available firmwares")
				return fmt.Errorf("unknown embedded firmware")
			}
			sourceName = fmt.Sprintf("embedded:%s", firmwareEmbedded)
		} else if firmwareFile != "" {
			// Use firmware from file
			var err error
			fileBytes, err = os.ReadFile(firmwareFile) //nolint:gosec // user-specified firmware file path is intentional
			if err != nil {
				fmt.Printf("Error: failed to read firmware file: %v\n", err)
				return err
			}
			sourceName = firmwareFile
		} else if len(args) == 1 {
			// Use firmware from positional argument
			var err error
			fileBytes, err = os.ReadFile(args[0])
			if err != nil {
				fmt.Printf("Error: failed to read firmware file: %v\n", err)
				return err
			}
			sourceName = args[0]
		} else {
			return fmt.Errorf("firmware source required: use --file <path>, --embedded <name>, or --list")
		}

		// Compute hash of firmware data (skip first 256 bytes which is the header)
		if len(fileBytes) <= 0x100 {
			fmt.Println("Error: firmware file too small (must be > 256 bytes)")
			return fmt.Errorf("firmware file too small")
		}
		hash := sha256.Sum256(fileBytes[0x100:])

		device := skyWallet.NewDevice(skyWallet.DeviceTypeFromString(deviceType))
		if device == nil {
			fmt.Println("Error: failed to create device (is device connected?)")
			return fmt.Errorf("failed to create device")
		}
		defer device.Close()

		// Check that device is in bootloader mode (required for firmware upload)
		if err := requireBootloaderMode(device); err != nil {
			fmt.Printf("Error: %v\n", err)
			fmt.Println("Hint: Put the device in bootloader mode by holding buttons while plugging in USB, or use --skip to bypass this check")
			return err
		}

		if os.Getenv("AUTO_PRESS_BUTTONS") == "1" && device.Driver.DeviceType() == skyWallet.DeviceTypeEmulator && runtime.GOOS == "linux" {
			err := device.SetAutoPressButton(true, skyWallet.ButtonRight)
			if err != nil {
				fmt.Printf("Error: failed to set auto press button: %v\n", err)
				return err
			}
		}

		fmt.Printf("Uploading firmware: %s (%d bytes)\n", sourceName, len(fileBytes))
		if err := device.FirmwareUpload(fileBytes, hash); err != nil {
			fmt.Printf("Error: firmware upload failed: %v\n", err)
			return err
		}

		fmt.Println("Firmware uploaded successfully")
		return nil
	},
}
