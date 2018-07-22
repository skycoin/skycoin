// hid - Gopher Interface Devices (USB HID)
// Copyright (c) 2017 Péter Szilágyi. All rights reserved.
//
// This file is released under the 3-clause BSD license. Note however that Linux
// support depends on libusb, released under GNU LGPL 2.1 or later.

// Package hid provides an interface for USB HID devices.

// +build !linux,!darwin,!windows ios !cgo

package usbhid

import (
	"errors"
)

// ErrDeviceClosed is returned for operations where the device closed before or
// during the execution.
var ErrDeviceClosed = errors.New("hid: device closed")

// ErrUnsupportedPlatform is returned for all operations where the underlying
// operating system is not supported by the library.
var ErrUnsupportedPlatform = errors.New("hid: unsupported platform")

type Context struct {}
type Device_Handle struct {}
type Device struct {}
type Endpoint_Descriptor struct {}
// A structure representing the standard USB device descriptor.
// This descriptor is documented in section 9.6.1 of the USB 3.0 specification.
// All multiple-byte fields are represented in host-endian format.
type Device_Descriptor struct {
	ptr                uintptr
	BLength            uint8
	BDescriptorType    uint8
	BcdUSB             uint16
	BDeviceClass       uint8
	BDeviceSubClass    uint8
	BDeviceProtocol    uint8
	BMaxPacketSize0    uint8
	IdVendor           uint16
	IdProduct          uint16
	BcdDevice          uint16
	IManufacturer      uint8
	IProduct           uint8
	ISerialNumber      uint8
	BNumConfigurations uint8
}

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

// A structure representing the standard USB configuration descriptor.
// This descriptor is documented in section 9.6.3 of the USB 3.0 specification.
// All multiple-byte fields are represented in host-endian format.
type Config_Descriptor struct {
	ptr                 uintptr
	BLength             uint8
	BDescriptorType     uint8
	WTotalLength        uint16
	BNumInterfaces      uint8
	BConfigurationValue uint8
	IConfiguration      uint8
	BmAttributes        uint8
	MaxPower            uint8
	Interface           []*Interface
	Extra               []byte
}

// A collection of alternate settings for a particular USB interface.
type Interface struct {
	ptr            uintptr
	Num_altsetting int
	Altsetting     []*Interface_Descriptor
}

// A structure representing the standard USB interface descriptor.
// This descriptor is documented in section 9.6.5 of the USB 3.0 specification.
// All multiple-byte fields are represented in host-endian format.
type Interface_Descriptor struct {
	ptr                uintptr
	BLength            uint8
	BDescriptorType    uint8
	BInterfaceNumber   uint8
	BAlternateSetting  uint8
	BNumEndpoints      uint8
	BInterfaceClass    uint8
	BInterfaceSubClass uint8
	BInterfaceProtocol uint8
	IInterface         uint8
	Endpoint           []*Endpoint_Descriptor
	Extra              []byte
}

// Log message levels.
const (
    LOG_LEVEL_NONE    = 0
    LOG_LEVEL_ERROR   = 1
    LOG_LEVEL_WARNING = 2
    LOG_LEVEL_INFO    = 3
    LOG_LEVEL_DEBUG   = 4
)

// Device and/or Interface Class codes.
const (
    CLASS_PER_INTERFACE       = 0
    CLASS_AUDIO               = 1
    CLASS_COMM                = 2
    CLASS_HID                 = 3
    CLASS_PHYSICAL            = 4
    CLASS_PRINTER             = 5
    CLASS_PTP                 = 6
    CLASS_IMAGE               = 7
    CLASS_MASS_STORAGE        = 8
    CLASS_HUB                 = 9
    CLASS_DATA                = 10
    CLASS_SMART_CARD          = 11
    CLASS_CONTENT_SECURITY    = 12
    CLASS_VIDEO               = 13
    CLASS_PERSONAL_HEALTHCARE = 14
    CLASS_DIAGNOSTIC_DEVICE   = 15
    CLASS_WIRELESS            = 16
    CLASS_APPLICATION         = 17
    CLASS_VENDOR_SPEC         = 18
)


// Error codes.
const (
	SUCCESS             = 0
	ERROR_IO            = 1
	ERROR_INVALID_PARAM = 2
	ERROR_ACCESS        = 3
	ERROR_NO_DEVICE     = 4
	ERROR_NOT_FOUND     = 5
	ERROR_BUSY          = 6
	ERROR_TIMEOUT       = 7
	ERROR_OVERFLOW      = 8
	ERROR_PIPE          = 9
	ERROR_INTERRUPTED   = 10
	ERROR_NO_MEM        = 11
	ERROR_NOT_SUPPORTED = 12
	ERROR_OTHER         = 13
)

// Enumerate returns a list of all the HID devices attached to the system which
// match the vendor and product id:
//  - If the vendor id is set to 0 then any vendor matches.
//  - If the product id is set to 0 then any product matches.
//  - If the vendor and product id are both 0, all HID devices are returned.
func HidEnumerate(vendorID uint16, productID uint16) []HidDeviceInfo {
	return nil
}

// Open connects to an HID device by its path name.
func (info HidDeviceInfo) Open() (*HidDevice, error) {
	return nil, ErrUnsupportedPlatform
}

// Device is a live HID USB connected device handle.
type HidDevice struct {
}

// Close releases the HID USB device handle.
func (dev *HidDevice) Close() error {
	return nil
}

// Write sends an output report to a HID device.
//
// Write will send the data on the first OUT endpoint, if one exists. If it does
// not, it will send the data through the Control Endpoint (Endpoint 0).
func (dev *HidDevice) Write(b []byte, prepend bool) (int, error) {
	return 0, ErrUnsupportedPlatform
}

// Read retrieves an input report from a HID device.
func (dev *HidDevice) Read(b []byte) (int, error) {
	
	return 0, ErrUnsupportedPlatform
}


func Set_Debug(ctx Context, level int) {}

func Init(ctx *Context) error {
    return ErrUnsupportedPlatform
}
func Exit(ctx Context) {}

func Get_Device_List(ctx Context) ([]Device, error) {
    return nil, ErrUnsupportedPlatform
}

func Free_Device_List(list []Device, unref_devices int) {
    
}
func Get_Device_Descriptor(dev Device) (*Device_Descriptor, error) {
    return nil, ErrUnsupportedPlatform
}

func Open(dev Device) (Device_Handle, error) {	var handle Device_Handle
	return handle, ErrUnsupportedPlatform
}

func Open_Device_With_VID_PID(ctx Context, vendor_id uint16, product_id uint16) Device_Handle {
	var handle Device_Handle
	return handle
}

func Close(hdl Device_Handle) {}

func Set_Configuration(hdl Device_Handle, configuration int) error {
	return ErrUnsupportedPlatform
}

func Claim_Interface(hdl Device_Handle, interface_number int) error {
	return ErrUnsupportedPlatform
}

func Reset_Device(hdl Device_Handle) error {
	return ErrUnsupportedPlatform
}

func Get_Port_Numbers(dev Device, ports []byte) ([]byte, error) {
	return nil, ErrUnsupportedPlatform
}
func Error_Name(code int) string {
	return ""
}
func Get_Active_Config_Descriptor(dev Device) (*Config_Descriptor, error) {
	return nil, ErrUnsupportedPlatform
}
func Interrupt_Transfer(hdl Device_Handle, endpoint uint8, data []byte, timeout uint) ([]byte, error) {
	return nil, ErrUnsupportedPlatform
}