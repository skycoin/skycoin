// +build darwin,!ios,cgo windows,cgo

package usb

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"sync"
	"sync/atomic"

	lowlevel "github.com/skycoin/hardware-wallet-go/src/usb/lowlevel/hidapi"
)

const (
	hidapiPrefix = "hid"
	hidUsagePage = 0xFF00
	hidTimeout   = 50
	HIDUse       = true
)

type HIDAPI struct {
}

func InitHIDAPI() (*HIDAPI, error) {
	return &HIDAPI{}, nil
}

func (b *HIDAPI) Enumerate(vendorID, productID uint16) ([]Info, error) {
	var infos []Info

	devs := lowlevel.HidEnumerate(vendorID, productID)

	for _, dev := range devs { // enumerate all devices
		if b.match(&dev) {
			infos = append(infos, Info{
				Path:      b.identify(&dev),
				VendorID:  int(dev.VendorID),
				ProductID: int(dev.ProductID),
				Type:      TypeT1Hid,
			})
		}
	}
	return infos, nil
}

func (b *HIDAPI) Has(path string) bool {
	return strings.HasPrefix(path, hidapiPrefix)
}

func (b *HIDAPI) Connect(path string) (Device, error) {
	devs := lowlevel.HidEnumerate(0, 0)

	for _, dev := range devs { // enumerate all devices
		if b.match(&dev) && b.identify(&dev) == path {
			d, err := dev.Open()
			if err != nil {
				return nil, err
			}
			prepend, err := b.detectPrepend(d)
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

func (b *HIDAPI) match(d *lowlevel.HidDeviceInfo) bool {
	vid := d.VendorID
	pid := d.ProductID
	wallet1 := vid == VendorT1 && (pid == ProductT1Firmware)
	wallet2 := vid == VendorT2 && (pid == ProductT2Firmware || pid == ProductT2Bootloader)
	return (wallet1 || wallet2) && (d.Interface == int(normalIface.number) || d.UsagePage == hidUsagePage)
}

func (b *HIDAPI) identify(dev *lowlevel.HidDeviceInfo) string {
	path := []byte(dev.Path)
	digest := sha256.Sum256(path)
	return hidapiPrefix + hex.EncodeToString(digest[:])
}

func (b *HIDAPI) Close() {
	// nothing
}

type HID struct {
	dev     *lowlevel.HidDevice
	prepend bool // on windows, see detectPrepend

	closed        int32 // atomic
	transferMutex sync.Mutex
	// closing cannot happen while read/write is hapenning,
	// otherwise it segfaults on windows
}

func (d *HID) Close(disconnected bool) error {
	atomic.StoreInt32(&d.closed, 1)

	d.transferMutex.Lock()
	err := d.dev.Close()
	d.transferMutex.Unlock()

	return err
}

var unknownErrorMessage = "hidapi: unknown failure"

// This will write a useless buffer to trezor
// to test whether it is an older HID version on reportid 63
// or a newer one that is on id 0.
// The older one does not need prepending, the newer one does
// This makes difference only on windows
func (b *HIDAPI) detectPrepend(dev *lowlevel.HidDevice) (bool, error) {
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

	return false, errors.New("unknown HID version")
}

func (d *HID) readWrite(buf []byte, read bool) (int, error) {
	for {
		closed := (atomic.LoadInt32(&d.closed)) == 1
		if closed {
			return 0, ErrClosedDevice
		}

		d.transferMutex.Lock()

		var w int
		var err error

		if read {
			w, err = d.dev.Read(buf, hidTimeout)
		} else {
			w, err = d.dev.Write(buf, d.prepend)
		}

		d.transferMutex.Unlock()

		if err == nil {
			// sometimes, empty report is read, skip it
			if w > 0 {
				return w, err
			}
			if !read {
				return 0, errors.New("HID - empty write")
			}

		} else {
			if err.Error() == unknownErrorMessage {
				return 0, ErrDisconnect
			}
			return 0, err
		}

		// continue the for cycle
	}
}

func (d *HID) Write(buf []byte) (int, error) {
	return d.readWrite(buf, false)
}

func (d *HID) Read(buf []byte) (int, error) {
	return d.readWrite(buf, true)
}
