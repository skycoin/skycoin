//go:build linux || freebsd
// +build linux freebsd

// shim for linux and freebsd so that daemon.go builds

package usb

// HIDUse indicates whether HID is available on this platform.
const HIDUse = false

// HIDAPI is a stub HID API for platforms where HID is not supported.
type HIDAPI struct {
}

// InitHIDAPI returns a stub HIDAPI instance.
func InitHIDAPI() (*HIDAPI, error) {
	return &HIDAPI{}, nil
}

// Enumerate is not implemented on this platform.
func (b *HIDAPI) Enumerate(_, _ uint16) ([]Info, error) {
	panic("not implemented for linux and freebsd")
}

// Has is not implemented on this platform.
func (b *HIDAPI) Has(_ string) bool {
	panic("not implemented for linux and freebsd")
}

// Connect is not implemented on this platform.
func (b *HIDAPI) Connect(_ string) (Device, error) {
	return &HID{}, nil
}

// HID is a stub HID device for platforms where HID is not supported.
type HID struct {
}

// Close is not implemented on this platform.
func (d *HID) Close(_ bool) error {
	panic("not implemented for linux and freebsd")
}

// Write is not implemented on this platform.
func (d *HID) Write(_ []byte) (int, error) {
	panic("not implemented for linux and freebsd")
}

// Read is not implemented on this platform.
func (d *HID) Read(_ []byte) (int, error) {
	panic("not implemented for linux and freebsd")
}

// Close is not implemented on this platform.
func (b *HIDAPI) Close() {
	panic("not implemented for linux and freebsd")
}
