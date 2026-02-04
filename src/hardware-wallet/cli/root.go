package cli

import (
	"fmt"

	"github.com/gogo/protobuf/proto"
	messages "github.com/skycoin/hardware-wallet-protob/go"

	skyWallet "github.com/skycoin/skycoin/src/hardware-wallet/skywallet"
)

var (
	deviceType     string
	addressN       int
	address        string
	message        string
	startIndex     int
	confirmAddress bool
	coinTypeStr    string
	usePassphrase  bool
	wordCount      int
	mnemonic       string
	label          string
	language       string
	inputHash      []string
	prevHash       []string
	inputIndex     []int
	outputAddress  []string
	coins          []int64
	hours          []int64
	addressIndex   []int
	entropyBytes   int
	signature      string
	skipModeCheck  bool // Global --skip flag
)

// DeviceMode represents the device's current mode
type DeviceMode int

const (
	// DeviceModeUnknown indicates the device mode could not be determined.
	DeviceModeUnknown DeviceMode = iota
	// DeviceModeFirmware indicates the device is running firmware.
	DeviceModeFirmware
	// DeviceModeBootloader indicates the device is in bootloader mode.
	DeviceModeBootloader
)

// checkDeviceMode checks if the device is in the expected mode
func checkDeviceMode(device *skyWallet.Device, expectedMode DeviceMode) error {
	msg, err := device.GetFeatures()
	if err != nil {
		return fmt.Errorf("failed to get device features: %w", err)
	}

	if msg.Kind != uint16(messages.MessageType_MessageType_Features) {
		return fmt.Errorf("unexpected response type: %s", messages.MessageType(msg.Kind))
	}

	features := &messages.Features{}
	if err := proto.Unmarshal(msg.Data, features); err != nil {
		return fmt.Errorf("failed to decode features: %w", err)
	}

	actualMode := DeviceModeFirmware
	if features.GetBootloaderMode() {
		actualMode = DeviceModeBootloader
	}

	if expectedMode != DeviceModeUnknown && actualMode != expectedMode {
		expectedStr := "firmware"
		actualStr := "firmware"
		if expectedMode == DeviceModeBootloader {
			expectedStr = "bootloader"
		}
		if actualMode == DeviceModeBootloader {
			actualStr = "bootloader"
		}
		return fmt.Errorf("device is in %s mode, but %s mode is required", actualStr, expectedStr)
	}

	return nil
}

// requireFirmwareMode checks that the device is in firmware mode (not bootloader)
// If skipModeCheck is true, the check is skipped
func requireFirmwareMode(device *skyWallet.Device) error {
	if skipModeCheck {
		return nil
	}
	return checkDeviceMode(device, DeviceModeFirmware)
}

// requireBootloaderMode checks that the device is in bootloader mode
// If skipModeCheck is true, the check is skipped
func requireBootloaderMode(device *skyWallet.Device) error {
	if skipModeCheck {
		return nil
	}
	return checkDeviceMode(device, DeviceModeBootloader)
}
