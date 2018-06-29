package usb

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"

	"github.com/usbhid"
)

const (
	hidapiPrefix = "hid"
	hidIfaceNum  = 0
	hidUsagePage = 0xFF00
)

// HIDAPI TODO documentation
type HIDAPI struct {
}

// InitHIDAPI TODO documentation
func InitHIDAPI() (*HIDAPI, error) {
	return &HIDAPI{}, nil
}

// Enumerate TODO documentation
func (b *HIDAPI) Enumerate() ([]Info, error) {
	var infos []Info

	for _, dev := range usbhid.HidEnumerate(0, 0) { // enumerate all devices
		if b.match(&dev) {
			infos = append(infos, Info{
				Path:      b.identify(&dev),
				VendorID:  int(dev.VendorID),
				ProductID: int(dev.ProductID),
			})
		}
	}
	return infos, nil
}

// Has TODO documentation
func (b *HIDAPI) Has(path string) bool {
	return strings.HasPrefix(path, hidapiPrefix)
}

// Connect TODO documentation
func (b *HIDAPI) Connect(path string) (Device, error) {
	for _, dev := range usbhid.HidEnumerate(0, 0) { // enumerate all devices
		if b.match(&dev) && b.identify(&dev) == path {
			d, err := dev.Open()
			if err != nil {
				return nil, err
			}
			prepend, err := detectPrepend(d)
			if err != nil {
				return nil, err
			}
			return &HID{
				dev:     d,
				prepend: prepend,
			}, nil
		}
	}
	return nil, ErrNotFound
}

func (b *HIDAPI) match(d *usbhid.HidDeviceInfo) bool {
	vid := d.VendorID
	pid := d.ProductID
	trezor1 := vid == vendorT1 && (pid == productT1Firmware || pid == productT1Bootloader)
	trezor2 := vid == vendorT2 && (pid == productT2Firmware || pid == productT2Bootloader)
	return (trezor1 || trezor2) && (d.Interface == hidIfaceNum || d.UsagePage == hidUsagePage)
}

func (b *HIDAPI) identify(dev *usbhid.HidDeviceInfo) string {
	path := []byte(dev.Path)
	digest := sha256.Sum256(path)
	return hidapiPrefix + hex.EncodeToString(digest[:])
}

// HID TODO documentation
type HID struct {
	dev     *usbhid.HidDevice
	prepend bool // on windows, see detectPrepend
}

// Close TODO documentation
func (d *HID) Close() error {
	return d.dev.Close()
}

var unknownErrorMessage = "hidapi: unknown failure"

// This will write a useless buffer to trezor
// to test whether it is an older HID version on reportid 63
// or a newer one that is on id 0.
// The older one does not need prepending, the newer one does
// This makes difference only on windows
func detectPrepend(dev *usbhid.HidDevice) (bool, error) {
	buf := []byte{63}
	for i := 0; i < 63; i++ {
		buf = append(buf, 0xff)
	}

	// first test newer version
	w, _ := dev.Write(buf, true)
	if w == 65 {
		return true, nil
	}

	// then test older version
	w, err := dev.Write(buf, false)
	if err != nil {
		return false, err
	}
	if w == 64 {
		return false, nil
	}

	return false, errors.New("Unknown HID version")
}

func (d *HID) readWrite(buf []byte, read bool) (int, error) {
	var w int
	var err error

	if read {
		w, err = d.dev.Read(buf)
	} else {
		w, err = d.dev.Write(buf, d.prepend)
	}

	if err != nil && err.Error() == unknownErrorMessage {
		return 0, errDisconnect
	}
	return w, err
}

// Write TODO documentation
func (d *HID) Write(buf []byte) (int, error) {
	return d.readWrite(buf, false)
}

// Read TODO documentation
func (d *HID) Read(buf []byte) (int, error) {
	return d.readWrite(buf, true)
}
