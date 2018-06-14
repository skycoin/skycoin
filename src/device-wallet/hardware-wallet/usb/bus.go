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
	ErrNotFound = fmt.Errorf("device not found")
)

type Info struct {
	Path      string
	VendorID  int
	ProductID int
}

type Device interface {
	io.ReadWriteCloser
}

type Bus interface {
	Enumerate() ([]Info, error)
	Connect(path string) (Device, error)
	Has(path string) bool
}

type USB struct {
	buses []Bus
}

func Init(buses ...Bus) *USB {
	return &USB{
		buses: buses,
	}
}

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

func (b *USB) Connect(path string) (Device, error) {
	for _, b := range b.buses {
		if b.Has(path) {
			return b.Connect(path)
		}
	}
	return nil, ErrNotFound
}

var disconnectError = errors.New("Device disconnected during action")
var closedDeviceError = errors.New("Closed device")
