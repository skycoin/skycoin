package skywallet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"runtime"
	"sync"
	"time"

	"github.com/gogo/protobuf/proto"

	messages "github.com/skycoin/hardware-wallet-protob/go"

	"github.com/skycoin/hardware-wallet-go/src/skywallet/usb"
	"github.com/skycoin/hardware-wallet-go/src/skywallet/wire"
)

// DeviceType type of device: emulator or usb
type DeviceType int32

func (dt DeviceType) String() string {
	switch dt {
	case DeviceTypeEmulator:
		return "EMULATOR"
	case DeviceTypeUSB:
		return "USB"
	default:
		return "Invalid"
	}
}

const (
	// DeviceTypeEmulator use emulator
	DeviceTypeEmulator DeviceType = iota + 1
	// DeviceTypeUSB use usb
	DeviceTypeUSB
	// DeviceTypeInvalid not valid value
	DeviceTypeInvalid
)

// CoinType type of coin, that will be used: Skycoin, Bitcoin, etc
type CoinType int32

const (
	// InvalidCoinType not valid value
	InvalidCoinType CoinType = iota + 1
	// SkycoinCoinType Skycoin coin type
	SkycoinCoinType
	// BitcoinCoinType Bitcoin coin type
	BitcoinCoinType
)

// CoinTypeFromString returns CoinType from String (i.e. SkycoinCoinType from 'SKY')
func CoinTypeFromString(ct string) (CoinType, error) {
	switch ct {
	case "SKY":
		return SkycoinCoinType, nil
	case "BTC":
		return BitcoinCoinType, nil
	default:
		return InvalidCoinType, fmt.Errorf("invalid coin type: %s", ct)
	}
}

const (
	// SkycoinVendorID from https://github.com/SkycoinProject/hardware-wallet/blob/50000f674c56c0cc18eec30d55978b73ed279b2e/tiny-firmware/bootloader/usb.c#L57
	SkycoinVendorID = 0x313A

	// SkycoinHwProductID from https://github.com/SkycoinProject/hardware-wallet/blob/50000f674c56c0cc18eec30d55978b73ed279b2e/tiny-firmware/bootloader/usb.c#L58
	SkycoinHwProductID = 0x0001

	// EmulatorPort is the emulator udp port
	EmulatorPort = 21324
)

//go:generate mockery -name DeviceDriver -case underscore -inpkg -testonly

// DeviceDriver is the api for hardware wallet communication
type DeviceDriver interface {
	SendToDevice(dev usb.Device, chunks [][64]byte) (wire.Message, error)
	SendToDeviceNoAnswer(dev usb.Device, chunks [][64]byte) error
	GetDevice() (usb.Device, error)
	GetDeviceInfos() ([]usb.Info, error)
	DeviceType() DeviceType
	Close()
}

// Driver represents a particular device (USB / Emulator)
type Driver struct {
	deviceType DeviceType
	bus        usb.Bus
}

func initUsb() []usb.Bus {
	w, err := usb.InitLibUSB(!usb.HIDUse, allowCancel(), detachKernelDriver())
	if err != nil {
		log.Fatalf("libusb: %s", err)
	}

	if !usb.HIDUse {
		return []usb.Bus{w}
	}

	h, err := usb.InitHIDAPI()
	if err != nil {
		log.Fatalf("hidapi: %s", err)
	}
	return []usb.Bus{w, h}
}

// NewDriver create a new device driver
func NewDriver(deviceType DeviceType) (*Driver, error) {
	switch deviceType {
	case DeviceTypeUSB:
		return &Driver{
			deviceType: deviceType,
			bus:        usb.Init(initUsb()...),
		}, nil
	case DeviceTypeEmulator:
		udpBus, err := usb.InitUDP([]int{EmulatorPort})
		if err != nil {
			return nil, err
		}

		return &Driver{
			deviceType: deviceType,
			bus:        usb.Init(udpBus),
		}, nil
	}

	return nil, fmt.Errorf("invalid device %s", deviceType)
}

// Close closes the bus
func (drv *Driver) Close() {
	drv.bus.Close()
}

// DeviceType return driver device type
func (drv *Driver) DeviceType() DeviceType {
	return drv.deviceType
}

// SendToDeviceNoAnswer sends msg to device and doesnt return response
func (drv *Driver) SendToDeviceNoAnswer(dev usb.Device, chunks [][64]byte) error {
	return sendToDeviceNoAnswer(dev, chunks)
}

// SendToDevice sends msg to device and returns response
func (drv *Driver) SendToDevice(dev usb.Device, chunks [][64]byte) (wire.Message, error) {
	return sendToDevice(dev, chunks)
}

// GetDevice returns a device instance
func (drv *Driver) GetDevice() (usb.Device, error) {
	var vendorID, productID uint16

	if drv.deviceType == DeviceTypeUSB {
		vendorID = SkycoinVendorID
		productID = SkycoinHwProductID
	} else if drv.deviceType != DeviceTypeEmulator {
		return nil, fmt.Errorf("invalid device type: %s", drv.deviceType)
	}

	infos, err := drv.bus.Enumerate(vendorID, productID)
	if len(infos) <= 0 {
		return nil, ErrNoDeviceConnected
	}

	if err != nil {
		return nil, err
	}

	tries := 0
	for tries < 3 {
		dev, err := drv.bus.Connect(infos[0].Path)
		if err != nil {
			log.Print(err.Error())
			tries++
			time.Sleep(100 * time.Millisecond)
		} else {
			return dev, err
		}
	}
	return nil, err
}

// GetDeviceInfos returns information from the attached usb
func (drv *Driver) GetDeviceInfos() ([]usb.Info, error) {
	if drv.DeviceType() == DeviceTypeUSB {
		infos, err := drv.bus.Enumerate(SkycoinVendorID, SkycoinHwProductID)
		if err != nil {
			return nil, err
		}

		return infos, nil
	}
	return nil, errors.New("reading device info make sense for physical devices only")
}

func sendToDeviceNoAnswer(dev usb.Device, chunks [][64]byte) error {
	for _, element := range chunks {
		_, err := dev.Write(element[:])
		if err != nil {
			return err
		}
	}
	return nil
}

func sendToDevice(dev usb.Device, chunks [][64]byte) (wire.Message, error) {
	var msg *wire.Message
	var err error
	for _, element := range chunks {
		_, err := dev.Write(element[:])
		if err != nil {
			return wire.Message{}, err
		}
	}

	msg, err = wire.ReadFrom(dev)
	if err != nil {
		return wire.Message{}, err
	}

	for msg.Kind == uint16(messages.MessageType_MessageType_EntropyRequest) {
		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			defer wg.Done()
			entropyChunks, err := MessageEntropyAck(entropyBufferSize)
			if err != nil {
				log.Errorf("failed to create entropy ack msg: %v", err)
				return
			}

			for _, element := range entropyChunks {
				_, err := dev.Write(element[:])
				if err != nil {
					log.Errorf("entropy ack error: %v", err)
					return
				}
			}
		}()

		msg, err = wire.ReadFrom(dev)
		if err != nil {
			return wire.Message{}, err
		}
		wg.Wait()
	}
	for msg.Kind == uint16(messages.MessageType_MessageType_Success) {
		success, err := decodeSuccessMsgStruct(*msg)
		if err != nil {
			return wire.Message{}, err
		}
		if success.MsgType != nil && *success.MsgType == messages.MessageType(messages.MessageType_MessageType_EntropyAck) {
			msg, err = wire.ReadFrom(dev)
			if err != nil {
				return wire.Message{}, err
			}
		} else {
			break
		}
	}

	return *msg, err
}

func binaryWrite(message io.Writer, data interface{}) {
	err := binary.Write(message, binary.BigEndian, data)
	if err != nil {
		log.Panic(err)
	}
}

func makeSkyWalletMessage(data []byte, msgID messages.MessageType) [][64]byte {
	message := new(bytes.Buffer)
	binaryWrite(message, []byte("##"))
	binaryWrite(message, uint16(msgID))
	binaryWrite(message, uint32(len(data)))
	binaryWrite(message, []byte("\n"))
	if len(data) > 0 {
		binaryWrite(message, data[1:])
	}

	messageLen := message.Len()
	var chunks [][64]byte
	i := 0
	for messageLen > 0 {
		var chunk [64]byte
		chunk[0] = '?'
		copy(chunk[1:], message.Bytes()[63*i:63*(i+1)])
		chunks = append(chunks, chunk)
		messageLen -= 63
		i = i + 1
	}
	return chunks
}

// Initialize send an init request to the device
func Initialize(dev usb.Device) error {
	var chunks [][64]byte

	chunks, err := MessageInitialize()
	if err != nil {
		return err
	}
	_, err = sendToDevice(dev, chunks)
	return err
}

// DecodeSuccessOrFailMsg parses a success or failure msg
func DecodeSuccessOrFailMsg(msg wire.Message) (string, error) {
	if msg.Kind == uint16(messages.MessageType_MessageType_Success) {
		return DecodeSuccessMsg(msg)
	}
	if msg.Kind == uint16(messages.MessageType_MessageType_Failure) {
		return DecodeFailMsg(msg)
	}

	return "", fmt.Errorf("calling DecodeSuccessOrFailMsg on message kind %s", messages.MessageType(msg.Kind))
}

func decodeSuccessMsgStruct(msg wire.Message) (messages.Success, error) {
	if msg.Kind == uint16(messages.MessageType_MessageType_Success) {
		success := messages.Success{}
		err := proto.Unmarshal(msg.Data, &success)
		if err != nil {
			return messages.Success{}, err
		}
		return success, nil
	}
	return messages.Success{}, fmt.Errorf("calling DecodeSuccessMsg with wrong message type: %s", messages.MessageType(msg.Kind))
}

// DecodeSuccessMsg convert byte data into string containing the success message returned by the device
func DecodeSuccessMsg(msg wire.Message) (string, error) {
	success, err := decodeSuccessMsgStruct(msg)
	if err != nil {
		return "", err
	}
	return success.GetMessage(), nil
}

// DecodeFailMsg convert byte data into string containing the failure returned by the device
func DecodeFailMsg(msg wire.Message) (string, error) {
	if msg.Kind == uint16(messages.MessageType_MessageType_Failure) {
		failure := &messages.Failure{}
		err := proto.Unmarshal(msg.Data, failure)
		if err != nil {
			return "", err
		}
		return failure.GetMessage(), nil
	}
	return "", fmt.Errorf("calling DecodeFailMsg with wrong message type: %s", messages.MessageType(msg.Kind))
}

// DecodeResponseSkycoinAddress convert byte data into list of addresses, meant to be used after DevicePinMatrixAck
func DecodeResponseSkycoinAddress(msg wire.Message) ([]string, error) {
	if msg.Kind == uint16(messages.MessageType_MessageType_ResponseSkycoinAddress) {
		responseSkycoinAddress := &messages.ResponseSkycoinAddress{}
		err := proto.Unmarshal(msg.Data, responseSkycoinAddress)
		if err != nil {
			return []string{}, err
		}
		return responseSkycoinAddress.GetAddresses(), nil
	}

	return []string{}, fmt.Errorf("calling DecodeResponseSkycoinAddress with wrong message type: %s", messages.MessageType(msg.Kind))
}

// DecodeResponseTransactionSign convert byte data into list of signatures
func DecodeResponseTransactionSign(msg wire.Message) ([]string, error) {
	if msg.Kind == uint16(messages.MessageType_MessageType_ResponseTransactionSign) {
		responseSkycoinTransactionSign := &messages.ResponseTransactionSign{}
		err := proto.Unmarshal(msg.Data, responseSkycoinTransactionSign)
		if err != nil {
			return make([]string, 0), err
		}
		return responseSkycoinTransactionSign.GetSignatures(), nil
	}

	return []string{}, fmt.Errorf("calling DecodeResponseeSkycoinSignMessage with wrong message type: %s", messages.MessageType(msg.Kind))
}

// DecodeResponseSkycoinSignMessage convert byte data into signed message, meant to be used after DevicePinMatrixAck
func DecodeResponseSkycoinSignMessage(msg wire.Message) (string, error) {
	if msg.Kind == uint16(messages.MessageType_MessageType_ResponseSkycoinSignMessage) {
		responseSkycoinSignMessage := &messages.ResponseSkycoinSignMessage{}
		err := proto.Unmarshal(msg.Data, responseSkycoinSignMessage)
		if err != nil {
			return "", err
		}
		return responseSkycoinSignMessage.GetSignedMessage(), nil
	}
	return "", fmt.Errorf("calling DecodeResponseeSkycoinSignMessage with wrong message type: %s", messages.MessageType(msg.Kind))
}

// DecodeResponseEntropyMessage convert byte data into entropy message, meant to be used after GetEntropy
func DecodeResponseEntropyMessage(msg wire.Message) (*messages.Entropy, error) {
	if msg.Kind == uint16(messages.MessageType_MessageType_Entropy) {
		responseEntropyMessage := &messages.Entropy{}
		err := proto.Unmarshal(msg.Data, responseEntropyMessage)
		if err != nil {
			return nil, err
		}
		return responseEntropyMessage, nil
	}
	return nil, fmt.Errorf("calling DecodeResponseEntropyMessage with wrong message type: %s", messages.MessageType(msg.Kind))
}

// Does OS allow sync canceling via our custom libusb patches?
func allowCancel() bool {
	return runtime.GOOS != "freebsd"
}

// Does OS detach kernel driver in libusb?
func detachKernelDriver() bool {
	return runtime.GOOS == "linux"
}
