// Package libusb provides a wrapper around github.com/google/gousb for compatibility.
// This package uses system-installed libusb-1.0.
package libusb

import (
	"fmt"

	"github.com/google/gousb"
)

// Type aliases for compatibility with existing code
type Context = *gousb.Context
type Device = *gousb.Device
type Device_Handle = *gousb.Device

// Device_Descriptor wraps gousb.DeviceDesc with compatible field names
type Device_Descriptor struct {
	*gousb.DeviceDesc
}

// Compatibility field accessors
func (d *Device_Descriptor) IdVendor() gousb.ID {
	return d.Vendor
}

func (d *Device_Descriptor) IdProduct() gousb.ID {
	return d.Product
}

func (d *Device_Descriptor) BcdDevice() gousb.BCD {
	return d.Device
}

type Config_Descriptor struct {
	BNumInterfaces uint8
	Interface      []Interface
}

type Interface struct {
	Num_altsetting int
	Altsetting     []Interface_Descriptor
}

type Interface_Descriptor struct {
	BInterfaceNumber   uint8
	BAlternateSetting  uint8
	BNumEndpoints      uint8
	BInterfaceClass    uint8
	Endpoint           []Endpoint_Descriptor
}

type Endpoint_Descriptor struct {
	BEndpointAddress uint8
}

// Constants
const (
	CLASS_HID         = uint8(gousb.ClassHID)
	CLASS_VENDOR_SPEC = uint8(gousb.ClassVendorSpec)
	
	ERROR_IO        = -1
	ERROR_NO_DEVICE = -4
	ERROR_OTHER     = -99
	ERROR_PIPE      = -9
)

// Init initializes a new libusb context
func Init(ctx *Context) error {
	newCtx := gousb.NewContext()
	*ctx = newCtx
	return nil
}

// Exit closes the libusb context
func Exit(ctx Context) {
	if ctx != nil {
		ctx.Close()
	}
}

// Get_Device_List_Filtered returns a list of USB devices matching VID/PID
// This avoids permission errors by only opening devices we actually need
func Get_Device_List_Filtered(ctx Context, vendorID, productID uint16) ([]*gousb.Device, error) {
	devices, err := ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		match := true
		if vendorID != 0 && uint16(desc.Vendor) != vendorID {
			match = false
		}
		if productID != 0 && uint16(desc.Product) != productID {
			match = false
		}
		return match
	})
	// OpenDevices returns (devices, error) where devices contains successfully opened devices
	// and error is only set if NO devices could be opened at all
	// Individual device open failures are logged but don't prevent returning other devices
	if len(devices) == 0 && err != nil {
		return nil, err
	}
	return devices, nil
}

// Get_Device_List returns a list of all USB devices
// Prefer Get_Device_List_Filtered when possible to avoid permission errors
func Get_Device_List(ctx Context) ([]*gousb.Device, error) {
	devices, err := ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		return true // Return all devices
	})
	if err != nil {
		return nil, err
	}
	return devices, nil
}

// Free_Device_List closes the device list
func Free_Device_List(list []*gousb.Device, unref int) {
	for _, dev := range list {
		if dev != nil {
			dev.Close()
		}
	}
}

// Get_Device_Descriptor returns the device descriptor
func Get_Device_Descriptor(dev Device) (*gousb.DeviceDesc, error) {
	if dev == nil {
		return nil, fmt.Errorf("device is nil")
	}
	desc := dev.Desc
	return desc, nil
}

// Get_Config_Descriptor returns a configuration descriptor by index
func Get_Config_Descriptor(dev Device, index uint8) (*Config_Descriptor, error) {
	if dev == nil {
		return nil, fmt.Errorf("device is nil")
	}
	
	desc := dev.Desc
	if desc == nil {
		return nil, fmt.Errorf("device descriptor is nil")
	}
	
	if len(desc.Configs) == 0 {
		return nil, fmt.Errorf("device has no configs in descriptor")
	}
	
	// gousb uses map[int]ConfigDesc where key is config number (usually 1-based)
	// For index 0, we want config number 1 (the default/first config)
	configNum := int(index) + 1
	cfg, ok := desc.Configs[configNum]
	if !ok {
		// Fallback: try index as direct key
		cfg, ok = desc.Configs[int(index)]
		if !ok {
			return nil, fmt.Errorf("config index %d (or %d) not found", index, configNum)
		}
	}
	
	result := &Config_Descriptor{
		BNumInterfaces: uint8(len(cfg.Interfaces)),
		Interface:      make([]Interface, len(cfg.Interfaces)),
	}
	
	// Build interface descriptors
	for i, iface := range cfg.Interfaces {
		result.Interface[i] = Interface{
			Num_altsetting: len(iface.AltSettings),
			Altsetting:     make([]Interface_Descriptor, len(iface.AltSettings)),
		}
		
		for j, alt := range iface.AltSettings {
			numEndpoints := len(alt.Endpoints)
			result.Interface[i].Altsetting[j] = Interface_Descriptor{
				BInterfaceNumber:  uint8(alt.Number),
				BAlternateSetting: uint8(alt.Alternate),
				BNumEndpoints:     uint8(numEndpoints),
				BInterfaceClass:   uint8(alt.Class),
				Endpoint:          make([]Endpoint_Descriptor, numEndpoints),
			}
			
			// Populate endpoint descriptors from map
			// Endpoints is map[EndpointAddress]EndpointDesc
			epIdx := 0
			for epAddr := range alt.Endpoints {
				if epIdx < numEndpoints {
					result.Interface[i].Altsetting[j].Endpoint[epIdx] = Endpoint_Descriptor{
						BEndpointAddress: uint8(epAddr), // epAddr is the EndpointAddress key
					}
					epIdx++
				}
			}
		}
	}
	
	return result, nil
}

// Get_Active_Config_Descriptor returns the active configuration descriptor
func Get_Active_Config_Descriptor(dev Device) (*Config_Descriptor, error) {
	return Get_Config_Descriptor(dev, 0)
}

// Open opens a device and returns a device handle
func Open(dev Device) (Device_Handle, error) {
	// In gousb, the device is already opened when returned from OpenDevices
	// We just return it as-is since Device and Device_Handle are aliases
	return dev, nil
}

// Close closes a device handle
func Close(handle Device_Handle) {
	if handle != nil {
		handle.Close()
	}
}

// Get_Configuration returns the current configuration value
func Get_Configuration(handle Device_Handle) (int, error) {
	cfg, err := handle.ActiveConfigNum()
	return cfg, err
}

// Set_Configuration sets the configuration value and returns the config object
// IMPORTANT: The returned *gousb.Config must be kept alive to maintain the interface claim
func Set_Configuration(handle Device_Handle, config int) (*gousb.Config, error) {
	return handle.Config(config)
}

// Claim_Interface claims an interface
func Claim_Interface(handle Device_Handle, iface int) error {
	// In gousb, interfaces are claimed when you call Config.Interface()
	// This is done in Interrupt_Transfer, so this is essentially a no-op
	return nil
}

// Release_Interface releases an interface
func Release_Interface(handle Device_Handle, iface int) error {
	// gousb automatically releases interfaces when the device is closed
	return nil
}

// Kernel_Driver_Active checks if a kernel driver is active
func Kernel_Driver_Active(handle Device_Handle, iface int) (bool, error) {
	// gousb handles this internally
	return false, nil
}

// Detach_Kernel_Driver detaches the kernel driver
func Detach_Kernel_Driver(handle Device_Handle, iface int) error {
	// gousb handles kernel driver detachment automatically
	// when claiming interfaces
	return nil
}

// Attach_Kernel_Driver attaches the kernel driver
func Attach_Kernel_Driver(handle Device_Handle, iface int) error {
	// gousb handles this internally
	return nil
}

// Interrupt_Transfer performs an interrupt transfer
// intf must be a valid Interface obtained from Config.Interface() and kept alive by caller
func Interrupt_Transfer(intf *gousb.Interface, endpoint uint8, data []byte, timeout int) ([]byte, error) {
	if intf == nil {
		return nil, fmt.Errorf("interface is nil")
	}
	
	// Determine if this is an IN or OUT endpoint
	isIn := (endpoint & 0x80) != 0
	
	if isIn {
		epIn, err := intf.InEndpoint(int(endpoint & 0x7F))
		if err != nil {
			return nil, err
		}
		n, err := epIn.Read(data)
		if err != nil {
			return nil, err
		}
		return data[:n], nil
	} else {
		epOut, err := intf.OutEndpoint(int(endpoint))
		if err != nil {
			return nil, err
		}
		n, err := epOut.Write(data)
		if err != nil {
			return nil, err
		}
		return data[:n], nil
	}
}

// Get_Port_Numbers returns the port numbers for a device
func Get_Port_Numbers(dev Device, ports []byte) ([]byte, error) {
	if dev == nil {
		return nil, fmt.Errorf("device is nil")
	}
	
	// Create a unique identifier from bus and address
	bus, addr := dev.Desc.Bus, dev.Desc.Address
	result := []byte{byte(bus), byte(addr)}
	
	if len(result) > len(ports) {
		result = result[:len(ports)]
	}
	
	copy(ports, result)
	return result, nil
}

// Cancel_Sync_Transfers_On_Device cancels synchronous transfers
func Cancel_Sync_Transfers_On_Device(handle Device_Handle) {
	// gousb handles cancellation differently
	// This is a no-op for compatibility
}

// Error_Name returns the error name for an error code
func Error_Name(code int) string {
	switch code {
	case ERROR_IO:
		return "LIBUSB_ERROR_IO"
	case ERROR_NO_DEVICE:
		return "LIBUSB_ERROR_NO_DEVICE"
	case ERROR_OTHER:
		return "LIBUSB_ERROR_OTHER"
	case ERROR_PIPE:
		return "LIBUSB_ERROR_PIPE"
	default:
		return fmt.Sprintf("LIBUSB_ERROR_%d", code)
	}
}
