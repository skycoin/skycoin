package usb

import (
	"encoding/hex"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/usbhid"
)

const (
	webusbPrefix  = "web"
	webConfigNum  = 1
	webIfaceNum   = 0
	webAltSetting = 0
	webEpIn       = 0x81
	webEpOut      = 0x01
	usbTimeout    = 5000
)

type WebUSB struct {
	usb usbhid.Context
}

func InitWebUSB() (*WebUSB, error) {
	var usb usbhid.Context
	err := usbhid.Init(&usb)
	if err != nil {
		return nil, err
	}
	usbhid.Set_Debug(usb, usbhid.LOG_LEVEL_NONE)

	return &WebUSB{
		usb: usb,
	}, nil
}

func (b *WebUSB) Close() {
	usbhid.Exit(b.usb)
}

func (b *WebUSB) Enumerate() ([]Info, error) {
	list, err := usbhid.Get_Device_List(b.usb)
	if err != nil {
		return nil, err
	}
	defer usbhid.Free_Device_List(list, 1) // unlink devices

	var infos []Info

	// There is a bug in either Trezor T or libusb that makes
	// device appear twice with the same path
	paths := make(map[string]bool)

	for _, dev := range list {
		if b.match(dev) {
			dd, err := usbhid.Get_Device_Descriptor(dev)
			if err != nil {
				continue
			}
			path := b.identify(dev)
			inset := paths[path]
			if !inset {
				infos = append(infos, Info{
					Path:      path,
					VendorID:  int(dd.IdVendor),
					ProductID: int(dd.IdProduct),
				})
				paths[path] = true
			}
		}
	}
	return infos, nil
}

func (b *WebUSB) Has(path string) bool {
	return strings.HasPrefix(path, webusbPrefix)
}

func (b *WebUSB) Connect(path string) (Device, error) {
	list, err := usbhid.Get_Device_List(b.usb)
	if err != nil {
		return nil, err
	}
	defer usbhid.Free_Device_List(list, 1) // unlink devices

	// There is a bug in either Trezor T or libusb that makes
	// device appear twice with the same path

	// We try both and return the first that works

	mydevs := make([]usbhid.Device, 0)
	for _, dev := range list {
		if b.match(dev) && b.identify(dev) == path {
			mydevs = append(mydevs, dev)
		}
	}

	err = ErrNotFound
	for _, dev := range mydevs {
		res, err := b.connect(dev)
		if err == nil {
			return res, nil
		}
	}
	return nil, err
}

func (b *WebUSB) connect(dev usbhid.Device) (*WUD, error) {
	d, err := usbhid.Open(dev)
	if err != nil {
		return nil, err
	}
	err = usbhid.Reset_Device(d)
	if err != nil {
		// don't abort if reset fails
		// usbhid.Close(d)
		// return nil, err
	}
	err = usbhid.Set_Configuration(d, webConfigNum)
	if err != nil {
		// don't abort if set configuration fails
		// usbhid.Close(d)
		// return nil, err
	}
	err = usbhid.Claim_Interface(d, webIfaceNum)
	if err != nil {
		usbhid.Close(d)
		return nil, err
	}
	return &WUD{
		dev:    d,
		closed: 0,
	}, nil
}

func (b *WebUSB) match(dev usbhid.Device) bool {
	dd, err := usbhid.Get_Device_Descriptor(dev)
	if err != nil {
		return false
	}
	vid := dd.IdVendor
	pid := dd.IdProduct
	trezor1 := vid == vendorT1 && (pid == productT1Firmware || pid == productT1Bootloader)
	trezor2 := vid == vendorT2 && (pid == productT2Firmware || pid == productT2Bootloader)
	if !trezor1 && !trezor2 {
		return false
	}
	c, err := usbhid.Get_Active_Config_Descriptor(dev)
	if err != nil {
		return false
	}
	return (c.BNumInterfaces > webIfaceNum &&
		c.Interface[webIfaceNum].Num_altsetting > webAltSetting &&
		c.Interface[webIfaceNum].Altsetting[webAltSetting].BInterfaceClass == usbhid.CLASS_VENDOR_SPEC)
}

func (b *WebUSB) identify(dev usbhid.Device) string {
	var ports [8]byte
	p, err := usbhid.Get_Port_Numbers(dev, ports[:])
	if err != nil {
		return ""
	}
	return webusbPrefix + hex.EncodeToString(p)
}

type WUD struct {
	dev usbhid.Device_Handle

	closed int32 // atomic

	transferMutex sync.Mutex
	// closing cannot happen while interrupt_transfer is hapenning,
	// otherwise interrupt_transfer hangs forever
}

func (d *WUD) Close() error {
	atomic.StoreInt32(&d.closed, 1)

	d.transferMutex.Lock()
	usbhid.Close(d.dev)
	d.transferMutex.Unlock()

	return nil
}

func (d *WUD) readWrite(buf []byte, endpoint uint8) (int, error) {
	for {
		closed := (atomic.LoadInt32(&d.closed)) == 1
		if closed {
			return 0, closedDeviceError
		}

		d.transferMutex.Lock()
		p, err := usbhid.Interrupt_Transfer(d.dev, endpoint, buf, usbTimeout)
		d.transferMutex.Unlock()

		if err == nil {
			// sometimes, empty report is read, skip it
			if len(p) > 0 {
				return len(p), err
			}
		}

		if err != nil {
			if err.Error() == usbhid.Error_Name(usbhid.ERROR_IO) ||
				err.Error() == usbhid.Error_Name(usbhid.ERROR_NO_DEVICE) {
				return 0, disconnectError
			}

			if err.Error() != usbhid.Error_Name(usbhid.ERROR_TIMEOUT) {
				return 0, err
			}
		}

		// continue the for cycle
	}
}

func (d *WUD) Write(buf []byte) (int, error) {
	return d.readWrite(buf, webEpOut)
}

func (d *WUD) Read(buf []byte) (int, error) {
	return d.readWrite(buf, webEpIn)
}
