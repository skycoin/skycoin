//go:build linux || freebsd
// +build linux freebsd

// shim for linux and freebsd so that daemon.go builds

// Package usb src/hardware-wallet/skywallet/usb/hidapi_shim.go c1-app-hw
package usb

import "errors"

// HIDUse indicates whether HID is available on this platform.
const HIDUse = false

var errHIDNotImplemented = errors.New("HID is not implemented for linux and freebsd")

// HIDAPI is a stub HID API for platforms where HID is not supported.
type HIDAPI struct {
}

// InitHIDAPI returns a stub HIDAPI instance.
func InitHIDAPI() (*HIDAPI, error) {
	return &HIDAPI{}, nil
}

// Enumerate is not implemented on this platform.
func (b *HIDAPI) Enumerate(_, _ uint16) ([]Info, error) {
	return nil, errHIDNotImplemented
}

// Has is not implemented on this platform.
func (b *HIDAPI) Has(_ string) bool {
	return false
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
	return errHIDNotImplemented
}

// Write is not implemented on this platform.
func (d *HID) Write(_ []byte) (int, error) {
	return 0, errHIDNotImplemented
}

// Read is not implemented on this platform.
func (d *HID) Read(_ []byte) (int, error) {
	return 0, errHIDNotImplemented
}

// Close does nothing on this platform.
func (b *HIDAPI) Close() {}
