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
	passphrase     string
	label          string
	language       string
	inputs         []string
	inputHash      []string
	prevHash       []string
	inputIndex     []int
	outputs        []string
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
	DeviceModeUnknown DeviceMode = iota
	DeviceModeFirmware
	DeviceModeBootloader
)

// checkDeviceMode checks if the device is in the expected mode
// Returns the actual mode found and an error if the check fails
func checkDeviceMode(device *skyWallet.Device, expectedMode DeviceMode) (DeviceMode, error) {
	msg, err := device.GetFeatures()
	if err != nil {
		return DeviceModeUnknown, fmt.Errorf("failed to get device features: %w", err)
	}

	if msg.Kind != uint16(messages.MessageType_MessageType_Features) {
		return DeviceModeUnknown, fmt.Errorf("unexpected response type: %s", messages.MessageType(msg.Kind))
	}

	features := &messages.Features{}
	if err := proto.Unmarshal(msg.Data, features); err != nil {
		return DeviceModeUnknown, fmt.Errorf("failed to decode features: %w", err)
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
		return actualMode, fmt.Errorf("device is in %s mode, but %s mode is required", actualStr, expectedStr)
	}

	return actualMode, nil
}

// requireFirmwareMode checks that the device is in firmware mode (not bootloader)
// If skipModeCheck is true, the check is skipped
func requireFirmwareMode(device *skyWallet.Device) error {
	if skipModeCheck {
		return nil
	}
	_, err := checkDeviceMode(device, DeviceModeFirmware)
	return err
}

// requireBootloaderMode checks that the device is in bootloader mode
// If skipModeCheck is true, the check is skipped
func requireBootloaderMode(device *skyWallet.Device) error {
	if skipModeCheck {
		return nil
	}
	_, err := checkDeviceMode(device, DeviceModeBootloader)
	return err
}
