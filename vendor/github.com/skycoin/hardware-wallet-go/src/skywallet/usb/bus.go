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
	VendorT1            = 0x313a
	ProductT1Bootloader = 0x0000
	ProductT1Firmware   = 0x0001
	VendorT2            = 0x1209
	ProductT2Bootloader = 0x53C0
	ProductT2Firmware   = 0x53C1
)

var (
	ErrNotFound     = errors.New("device not found")
	ErrDisconnect   = errors.New("device disconnected during action")
	ErrClosedDevice = errors.New("closed device")
)

type DeviceType int

const (
	TypeT1Hid        DeviceType = 0
	TypeT1Webusb     DeviceType = 1
	TypeT1WebusbBoot DeviceType = 2
	TypeT2           DeviceType = 3
	TypeT2Boot       DeviceType = 4
	TypeEmulator     DeviceType = 5
)

type Info struct {
	Path      string
	VendorID  int
	ProductID int
	Type      DeviceType
}

type Device interface {
	io.ReadWriter
	Close(disconnected bool) error
}

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

type USB struct {
	buses []Bus
}

func Init(buses ...Bus) *USB {
	return &USB{
		buses: buses,
	}
}

func (b *USB) Has(path string) bool {
	for _, b := range b.buses {
		if b.Has(path) {
			return true
		}
	}
	return false
}

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

func (b *USB) Connect(path string) (Device, error) {
	for _, b := range b.buses {
		if b.Has(path) {
			return b.Connect(path)
		}
	}
	return nil, ErrNotFound
}

func (b *USB) Close() {
	for _, b := range b.buses {
		b.Close()
	}
}
