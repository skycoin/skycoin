package usb

import (
	"io"
	"net"
	"strconv"
	"strings"
	"sync/atomic"
)

const (
	emulatorPrefix  = "emulator"
	emulatorAddress = "127.0.0.1"
)

// UDP implements Bus for emulator communication over UDP.
type UDP struct {
	ports []int
}

// InitUDP creates a new UDP bus for the given ports.
func InitUDP(ports []int) (*UDP, error) {
	udp := UDP{
		ports: ports,
	}

	return &udp, nil
}

// Enumerate returns emulator device info for all configured ports.
func (udp *UDP) Enumerate(_, _ uint16) ([]Info, error) {
	var infos []Info

	for _, port := range udp.ports {
		info := Info{
			Path:      emulatorPrefix + strconv.Itoa(port),
			VendorID:  0,
			ProductID: 0,
			Type:      TypeEmulator,
		}

		infos = append(infos, info)
	}
	return infos, nil
}

// Has returns true if the path matches an emulator device.
func (udp *UDP) Has(path string) bool {
	return strings.HasPrefix(path, emulatorPrefix)
}

// Connect establishes a UDP connection to an emulator device.
func (udp *UDP) Connect(path string) (Device, error) {
	port, err := strconv.Atoi(strings.TrimPrefix(path, emulatorPrefix))
	if err != nil {
		return nil, err
	}

	address := emulatorAddress + ":" + strconv.Itoa(port)
	dev, err := net.Dial("udp", address)
	if err != nil {
		return nil, err
	}

	d := &UDPDevice{
		dev: dev,
	}
	return d, nil
}

// Close is a no-op for UDP connections.
func (udp *UDP) Close() {
	// nothing
}

// UDPDevice represents a UDP-connected emulated device.
type UDPDevice struct {
	dev io.ReadWriteCloser

	closed int32 // atomic
}

// Close closes the UDP device connection.
func (d *UDPDevice) Close(_ bool) error {
	atomic.StoreInt32(&d.closed, 1)
	return d.dev.Close()
}

func (d *UDPDevice) Write(buf []byte) (int, error) {
	closed := (atomic.LoadInt32(&d.closed)) == 1
	if closed {
		return 0, ErrClosedDevice
	}
	return d.dev.Write(buf)
}

func (d *UDPDevice) Read(buf []byte) (int, error) {
	closed := (atomic.LoadInt32(&d.closed)) == 1
	if closed {
		return 0, ErrClosedDevice
	}

	return d.dev.Read(buf)
}
