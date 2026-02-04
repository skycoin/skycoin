// Package firmware provides embedded firmware binaries for the Skycoin hardware wallet.
// This package embeds multiple firmware versions that can be flashed to the device.
package firmware

import (
	_ "embed"
)

// Legacy firmware versions (C firmware from official releases)

//go:embed skywallet-firmware-v1.0.0.bin
var firmwareV100 []byte

//go:embed skywallet-firmware-v1.8.0.bin
var firmwareV180 []byte

// Current firmware versions (built from source)

//go:embed skywallet-firmware-c.bin
var firmwareC []byte

//go:embed skywallet-firmware-tinygo.bin
var firmwareTinyGo []byte

// FirmwareType identifies which firmware variant to use
type FirmwareType int

const (
	FirmwareV100    FirmwareType = iota // Legacy v1.0.0 release
	FirmwareV180                        // Legacy v1.8.0 release
	FirmwareC                           // Current C firmware (v1.1.0)
	FirmwareTinyGo                      // TinyGo firmware
)

// FirmwareInfo contains metadata about a firmware variant
type FirmwareInfo struct {
	Type        FirmwareType
	Name        string
	Version     string
	Description string
	Data        []byte
}

// Available returns all available firmware variants
func Available() []FirmwareInfo {
	return []FirmwareInfo{
		{
			Type:        FirmwareV100,
			Name:        "v1.0.0",
			Version:     "1.0.0",
			Description: "Original release firmware",
			Data:        firmwareV100,
		},
		{
			Type:        FirmwareV180,
			Name:        "v1.8.0",
			Version:     "1.8.0",
			Description: "Latest stable C firmware release",
			Data:        firmwareV180,
		},
		{
			Type:        FirmwareC,
			Name:        "C (dev)",
			Version:     "1.1.0",
			Description: "Current C firmware built from source",
			Data:        firmwareC,
		},
		{
			Type:        FirmwareTinyGo,
			Name:        "TinyGo",
			Version:     "0.1.0",
			Description: "Experimental TinyGo firmware",
			Data:        firmwareTinyGo,
		},
	}
}

// Get returns the firmware binary for the specified type
func Get(t FirmwareType) []byte {
	switch t {
	case FirmwareV100:
		return firmwareV100
	case FirmwareV180:
		return firmwareV180
	case FirmwareC:
		return firmwareC
	case FirmwareTinyGo:
		return firmwareTinyGo
	default:
		return nil
	}
}

// GetByName returns firmware by its name string
func GetByName(name string) []byte {
	for _, fw := range Available() {
		if fw.Name == name {
			return fw.Data
		}
	}
	return nil
}
