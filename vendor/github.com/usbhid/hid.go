// hid - Gopher Interface Devices (USB HID)
// Copyright (c) 2017 Péter Szilágyi. All rights reserved.
//
// This file is released under the 3-clause BSD license. Note however that Linux
// support depends on libusb, released under GNU LGPL 2.1 or later.

// Package hid provides an interface for USB HID devices.

// +build linux,cgo darwin,!ios,cgo windows,cgo

package usbhid

/*
#cgo CFLAGS: -I./c/hidapi -I./c/libusb

#cgo linux CFLAGS: -DDEFAULT_VISIBILITY="" -DOS_LINUX -D_GNU_SOURCE -DPOLL_NFDS_TYPE=int
#cgo linux,!android LDFLAGS: -lrt
#cgo darwin CFLAGS: -DOS_DARWIN -DDEFAULT_VISIBILITY="" -DPOLL_NFDS_TYPE="unsigned int"
#cgo darwin LDFLAGS: -framework CoreFoundation -framework IOKit -lobjc
#cgo windows CFLAGS: -DOS_WINDOWS -DDEFAULT_VISIBILITY="" -DPOLL_NFDS_TYPE="unsigned int"
#cgo windows LDFLAGS: -lsetupapi


#ifdef OS_LINUX
	#include <sys/poll.h>

	#include "os/threads_posix.c"
	#include "os/poll_posix.c"
	#include "os/linux_usbfs.c"
	#include "os/linux_netlink.c"
#elif OS_DARWIN
	#include <sys/poll.h>

	#include "os/threads_posix.c"
	#include "os/poll_posix.c"
	#include "os/darwin_usb.c"
#elif OS_WINDOWS
	#include <oledlg.h>

  #include "os/poll_windows.c"
	#include "os/threads_windows.c"
#endif

#include "core.c"
#include "descriptor.c"
#include "hotplug.c"
#include "io.c"
#include "strerror.c"
#include "sync.c"

#ifdef OS_LINUX
	#include "linux/hid.c"
#elif OS_DARWIN
	#include "mac/hid.c"
#elif OS_WINDOWS
	#include "windows/hid.c"

	#include "os/windows_nt_common.c"
	#include "os/windows_winusb.c"
#endif

*/
import "C"

import (
	"errors"
	"runtime"
	"sync"
	"unsafe"
)

// ErrDeviceClosed is returned for operations where the device closed before or
// during the execution.
var ErrDeviceClosed = errors.New("hid: device closed")

// ErrUnsupportedPlatform is returned for all operations where the underlying
// operating system is not supported by the library.
var ErrUnsupportedPlatform = errors.New("hid: unsupported platform")

// HidDeviceInfo is a hidapi info structure.
type HidDeviceInfo struct {
	Path         string // Platform-specific device path
	VendorID     uint16 // Device Vendor ID
	ProductID    uint16 // Device Product ID
	Release      uint16 // Device Release Number in binary-coded decimal, also known as Device Version Number
	Serial       string // Serial Number
	Manufacturer string // Manufacturer String
	Product      string // Product string
	UsagePage    uint16 // Usage Page for this Device/Interface (Windows/Mac only)
	Usage        uint16 // Usage for this Device/Interface (Windows/Mac only)

	// The USB interface which this logical device
	// represents. Valid on both Linux implementations
	// in all cases, and valid on the Windows implementation
	// only if the device contains more than one interface.
	Interface int
}

// enumerateLock is a mutex serializing access to USB device enumeration needed
// by the macOS USB HID system calls, which require 2 consecutive method calls
// for enumeration, causing crashes if called concurrently.
//
// For more details, see:
//   https://developer.apple.com/documentation/iokit/1438371-iohidmanagersetdevicematching
//   > "subsequent calls will cause the hid manager to release previously enumerated devices"
var enumerateLock sync.Mutex

func init() {
	// Initialize the HIDAPI library
	C.hid_init()
}

// Enumerate returns a list of all the HID devices attached to the system which
// match the vendor and product id:
//  - If the vendor id is set to 0 then any vendor matches.
//  - If the product id is set to 0 then any product matches.
//  - If the vendor and product id are both 0, all HID devices are returned.
func HidEnumerate(vendorID uint16, productID uint16) []HidDeviceInfo {
	enumerateLock.Lock()
	defer enumerateLock.Unlock()

	// Gather all device infos and ensure they are freed before returning
	head := C.hid_enumerate(C.ushort(vendorID), C.ushort(productID))
	if head == nil {
		return nil
	}
	defer C.hid_free_enumeration(head)

	// Iterate the list and retrieve the device details
	var infos []HidDeviceInfo
	for ; head != nil; head = head.next {
		info := HidDeviceInfo{
			Path:      C.GoString(head.path),
			VendorID:  uint16(head.vendor_id),
			ProductID: uint16(head.product_id),
			Release:   uint16(head.release_number),
			UsagePage: uint16(head.usage_page),
			Usage:     uint16(head.usage),
			Interface: int(head.interface_number),
		}
		if head.serial_number != nil {
			info.Serial, _ = wcharTToString(head.serial_number)
		}
		if head.product_string != nil {
			info.Product, _ = wcharTToString(head.product_string)
		}
		if head.manufacturer_string != nil {
			info.Manufacturer, _ = wcharTToString(head.manufacturer_string)
		}
		infos = append(infos, info)
	}
	return infos
}

// Open connects to an HID device by its path name.
func (info HidDeviceInfo) Open() (*HidDevice, error) {
	path := C.CString(info.Path)
	defer C.free(unsafe.Pointer(path))

	device := C.hid_open_path(path)
	if device == nil {
		return nil, errors.New("hidapi: failed to open device")
	}
	return &HidDevice{
		HidDeviceInfo: info,
		device:        device,
	}, nil
}

// Device is a live HID USB connected device handle.
type HidDevice struct {
	HidDeviceInfo // Embed the infos for easier access

	device *C.hid_device // Low level HID device to communicate through
	lock   sync.Mutex
}

// Close releases the HID USB device handle.
func (dev *HidDevice) Close() error {
	dev.lock.Lock()
	defer dev.lock.Unlock()

	if dev.device != nil {
		C.hid_close(dev.device)
		dev.device = nil
	}
	return nil
}

// Write sends an output report to a HID device.
//
// Write will send the data on the first OUT endpoint, if one exists. If it does
// not, it will send the data through the Control Endpoint (Endpoint 0).
func (dev *HidDevice) Write(b []byte, prepend bool) (int, error) {
	// Abort if nothing to write
	if len(b) == 0 {
		return 0, nil
	}
	// Abort if device closed in between
	dev.lock.Lock()
	device := dev.device
	dev.lock.Unlock()

	if device == nil {
		return 0, ErrDeviceClosed
	}
	// Prepend a HID report ID on Windows, other OSes don't need it
	var report []byte
	if prepend && runtime.GOOS == "windows" {
		report = append([]byte{0x00}, b...)
	} else {
		report = b
	}
	// Execute the write operation
	written := int(C.hid_write(device, (*C.uchar)(&report[0]), C.size_t(len(report))))
	if written == -1 {
		// If the write failed, verify if closed or other error
		dev.lock.Lock()
		device = dev.device
		dev.lock.Unlock()

		if device == nil {
			return 0, ErrDeviceClosed
		}
		// Device not closed, some other error occurred
		message := C.hid_error(device)
		if message == nil {
			return 0, errors.New("hidapi: unknown failure")
		}
		failure, _ := wcharTToString(message)
		return 0, errors.New("hidapi: " + failure)
	}
	return written, nil
}

// Read retrieves an input report from a HID device.
func (dev *HidDevice) Read(b []byte) (int, error) {
	// Aborth if nothing to read
	if len(b) == 0 {
		return 0, nil
	}
	// Abort if device closed in between
	dev.lock.Lock()
	device := dev.device
	dev.lock.Unlock()

	if device == nil {
		return 0, ErrDeviceClosed
	}
	// Execute the read operation
	read := int(C.hid_read(device, (*C.uchar)(&b[0]), C.size_t(len(b))))
	if read == -1 {
		// If the read failed, verify if closed or other error
		dev.lock.Lock()
		device = dev.device
		dev.lock.Unlock()

		if device == nil {
			return 0, ErrDeviceClosed
		}
		// Device not closed, some other error occurred
		message := C.hid_error(device)
		if message == nil {
			return 0, errors.New("hidapi: unknown failure")
		}
		failure, _ := wcharTToString(message)
		return 0, errors.New("hidapi: " + failure)
	}
	return read, nil
}
