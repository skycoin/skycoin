// Package usb provides hardware wallet USB device communication.
package usb

import (
	"errors"
	"io"

	"github.com/skycoin/skycoin/src/util/logging"
)

var (
	log = logging.MustGetLogger("walletusb")
)

const (
	// VendorT1 is the USB vendor ID for Trezor One devices.
	VendorT1 = 0x313a
	// ProductT1Bootloader is the USB product ID for Trezor One in bootloader mode.
	ProductT1Bootloader = 0x0000
	// ProductT1Firmware is the USB product ID for Trezor One in firmware mode.
	ProductT1Firmware = 0x0001
	// VendorT2 is the USB vendor ID for Trezor T devices.
	VendorT2 = 0x1209
	// ProductT2Bootloader is the USB product ID for Trezor T in bootloader mode.
	ProductT2Bootloader = 0x53C0
	// ProductT2Firmware is the USB product ID for Trezor T in firmware mode.
	ProductT2Firmware = 0x53C1
)

var (
	// ErrNotFound is returned when the device is not found.
	ErrNotFound = errors.New("device not found")
	// ErrDisconnect is returned when the device disconnects during an action.
	ErrDisconnect = errors.New("device disconnected during action")
	// ErrClosedDevice is returned when attempting to use a closed device.
	ErrClosedDevice = errors.New("closed device")
)

// DeviceType represents the type of hardware wallet device.
type DeviceType int

const (
	// TypeT1Hid is a Trezor One device using HID protocol.
	TypeT1Hid DeviceType = 0
	// TypeT1Webusb is a Trezor One device using WebUSB protocol.
	TypeT1Webusb DeviceType = 1
	// TypeT1WebusbBoot is a Trezor One device in bootloader mode using WebUSB.
	TypeT1WebusbBoot DeviceType = 2
	// TypeT2 is a Trezor T device.
	TypeT2 DeviceType = 3
	// TypeT2Boot is a Trezor T device in bootloader mode.
	TypeT2Boot DeviceType = 4
	// TypeEmulator is an emulated device.
	TypeEmulator DeviceType = 5
)

// Info contains USB device identification information.
type Info struct {
	Path      string
	VendorID  int
	ProductID int
	Type      DeviceType
}

// Device represents a connected USB hardware wallet device.
type Device interface {
	io.ReadWriter
	Close(disconnected bool) error
}

// Bus represents a USB bus that can enumerate and connect to devices.
type Bus interface {
	// Enumerate returns a list of all the devices accessible in the the system
	// - If the vendor id is set to 0 then any vendor matches.
	// - If the product id is set to 0 then any product matches.
	// - If the vendor and product id are both 0, all devices are returned.
	Enumerate(vendorID, productID uint16) ([]Info, error)
	Connect(path string) (Device, error)
	Has(path string) bool
	Close() // called on program exit
}

// USB manages multiple USB buses for device communication.
type USB struct {
	buses []Bus
}

// Init creates a new USB manager with the given buses.
func Init(buses ...Bus) *USB {
	return &USB{
		buses: buses,
	}
}

// Has returns true if any bus has a device at the given path.
func (b *USB) Has(path string) bool {
	for _, b := range b.buses {
		if b.Has(path) {
			return true
		}
	}
	return false
}

// Enumerate returns device info for all matching devices across all buses.
func (b *USB) Enumerate(vendorID, productID uint16) ([]Info, error) {
	var infos []Info

	for _, b := range b.buses {
		l, err := b.Enumerate(vendorID, productID)
		if err != nil {
			return nil, err
		}
		infos = append(infos, l...)
	}
	return infos, nil
}

// Connect opens a connection to the device at the given path.
func (b *USB) Connect(path string) (Device, error) {
	for _, b := range b.buses {
		if b.Has(path) {
			return b.Connect(path)
		}
	}
	return nil, ErrNotFound
}

// Close shuts down all USB buses.
func (b *USB) Close() {
	for _, b := range b.buses {
		b.Close()
	}
}
