package usb

import (
	"bytes"
	"io"
	"net"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

var emulatorPing = []byte("PINGPING")
var emulatorPong = []byte("PONGPONG")

const (
	emulatorPrefix      = "emulator"
	emulatorAddress     = "127.0.0.1"
	emulatorPingTimeout = 700 * time.Millisecond
)

// UDP TODO documentation
type UDP struct {
	ports []int

	pings   map[int](chan []byte)
	datas   map[int](chan []byte)
	writers map[int](io.Writer)
}

func listen(conn net.Conn) (chan []byte, chan []byte) {
	ping := make(chan []byte, 1)
	data := make(chan []byte, 100)
	go func() {
		for {
			buffer := make([]byte, 64)
			_, err := conn.Read(buffer)
			if err == nil {
				first := buffer[0]
				if first == '?' {
					data <- buffer
				}
				if first == 'P' {
					copied := make([]byte, 8)
					copy(copied, buffer)
					ping <- copied
				}
			}
		}
	}()
	return ping, data
}

// InitUDP TODO documentation
func InitUDP(ports []int) (*UDP, error) {
	udp := UDP{
		ports: ports,

		pings:   make(map[int](chan []byte)),
		datas:   make(map[int](chan []byte)),
		writers: make(map[int](io.Writer)),
	}
	for _, port := range ports {
		address := emulatorAddress + ":" + strconv.Itoa(port)

		connection, err := net.Dial("udp", address)
		if err != nil {
			return nil, err
		}

		ping, data := listen(connection)
		udp.pings[port] = ping
		udp.datas[port] = data
		udp.writers[port] = connection
	}
	return &udp, nil
}

func checkPort(ping chan []byte, w io.Writer) (bool, error) {
	_, err := w.Write(emulatorPing)
	if err != nil {
		return false, err
	}
	select {
	case response := <-ping:
		return bytes.Equal(response, emulatorPong), nil
	case <-time.After(emulatorPingTimeout):
		return false, nil
	}
}

// Enumerate TODO documentation
func (u *UDP) Enumerate() ([]Info, error) {
	var infos []Info

	for _, port := range u.ports {
		ping := u.pings[port]
		w := u.writers[port]
		present, err := checkPort(ping, w)
		if err != nil {
			return nil, err
		}
		if present {
			infos = append(infos, Info{
				Path:      emulatorPrefix + strconv.Itoa(port),
				VendorID:  0,
				ProductID: 0,
			})
		}
	}
	return infos, nil
}

// Has TODO documentation
func (u *UDP) Has(path string) bool {
	return strings.HasPrefix(path, emulatorPrefix)
}

// Connect TODO documentation
func (u *UDP) Connect(path string) (Device, error) {
	i, err := strconv.Atoi(strings.TrimPrefix(path, emulatorPrefix))
	if err != nil {
		return nil, err
	}
	return &UDPDevice{
		ping:   u.pings[i],
		data:   u.datas[i],
		writer: u.writers[i],
		closed: 0,
	}, nil
}

// UDPDevice TODO documentation
type UDPDevice struct {
	ping   chan []byte
	data   chan []byte
	writer io.Writer

	closed int32 // atomic
}

// Close TODO documentation
func (d *UDPDevice) Close() error {
	atomic.StoreInt32(&d.closed, 1)
	return nil
}

func (d *UDPDevice) readWrite(buf []byte, read bool) (int, error) {
	for {
		closed := (atomic.LoadInt32(&d.closed)) == 1
		if closed {
			return 0, errClosedDeviceError
		}
		check, err := checkPort(d.ping, d.writer)
		if err != nil {
			return 0, err
		}
		if !check {
			return 0, errDisconnect
		}
		if !read {
			return d.writer.Write(buf)
		}
		select {
		case response := <-d.data:
			copy(buf, response)
			return len(response), nil
		case <-time.After(emulatorPingTimeout):
			// timeout, continue for cycle
		}
	}
}

// Write TODO documentation
func (d *UDPDevice) Write(buf []byte) (int, error) {
	return d.readWrite(buf, false)
}

// Read TODO documentation
func (d *UDPDevice) Read(buf []byte) (int, error) {
	return d.readWrite(buf, true)
}
