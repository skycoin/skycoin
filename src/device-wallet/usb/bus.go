package usb

import (
	"errors"
	"fmt"
	"io"
)

const (
	vendorT1            = 0x534c
	productT1Bootloader = 0x0000
	productT1Firmware   = 0x0001
	vendorT2            = 0x1209
	productT2Bootloader = 0x53C0
	productT2Firmware   = 0x53C1
)

var (
	// ErrNotFound device not found
	ErrNotFound = fmt.Errorf("device not found")
)

// Info driver information about the usb device
type Info struct {
	Path      string
	VendorID  int
	ProductID int
}

// Device interface for object that has methods Read Write and Close
type Device interface {
	io.ReadWriteCloser
}

// Bus interface for object that has Enumerate Connect and Has functions
type Bus interface {
	Enumerate() ([]Info, error)
	Connect(path string) (Device, error)
	Has(path string) bool
}

// USB list of Bus
type USB struct {
	buses []Bus
}

// Init Creates USB structure from Bus objects
func Init(buses ...Bus) *USB {
	return &USB{
		buses: buses,
	}
}

// Enumerate enumerates the devices connected to the usb
func (b *USB) Enumerate() ([]Info, error) {
	var infos []Info

	for _, b := range b.buses {
		l, err := b.Enumerate()
		if err != nil {
			return nil, err
		}
		infos = append(infos, l...)
	}
	return infos, nil
}

// Connect try to connect to a device at the given path
func (b *USB) Connect(path string) (Device, error) {
	for _, b := range b.buses {
		if b.Has(path) {
			return b.Connect(path)
		}
	}
	return nil, ErrNotFound
}

var errDisconnect = errors.New("Device disconnected during action")
var errClosedDeviceError = errors.New("Closed device")
