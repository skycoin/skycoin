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

type UDP struct {
	ports []int
}

func InitUDP(ports []int) (*UDP, error) {
	udp := UDP{
		ports: ports,
	}

	return &udp, nil
}

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

func (udp *UDP) Has(path string) bool {
	return strings.HasPrefix(path, emulatorPrefix)
}

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

func (udp *UDP) Close() {
	// nothing
}

type UDPDevice struct {
	dev io.ReadWriteCloser

	closed int32 // atomic
}

func (d *UDPDevice) Close(disconnected bool) error {
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
