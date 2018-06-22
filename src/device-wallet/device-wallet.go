
package deviceWallet

import (
	"net"
	"bytes"
	"encoding/binary"
	"fmt"
    "time"
    "io"
    "log"
	"github.com/wire"
	"github.com/skycoin/skycoin/src/device-wallet/usb"

    proto "github.com/proto"
	messages "github.com/skycoin/skycoin/protob"
)
// DeviceType type of device: emulated or usb
type DeviceType int32

const (
    // DeviceTypeEmulator use emulator
    DeviceTypeEmulator DeviceType = 1
    // DeviceTypeUsb use usb
	DeviceTypeUsb      DeviceType = 2 
)

func getEmulatorDevice() (net.Conn, error) {
	return net.Dial("udp", "127.0.0.1:21324")
}


func getUsbDevice() (usb.Device, error) {
	w, err := usb.InitWebUSB()
	if err != nil {
		log.Fatalf("webusb: %s", err)
	}
	h, err := usb.InitHIDAPI()
	if err != nil {
		log.Fatalf("hidapi: %s", err)
	}
	b := usb.Init(w, h)

	var infos []usb.Info
	infos, _ = b.Enumerate()

	tries := 0
	dev, err := b.Connect(infos[0].Path)
	if err != nil {
		fmt.Printf(err.Error())
		if tries < 3 {
			tries++
			time.Sleep(100 * time.Millisecond)
		}
	}
	return dev, err
}

func sendToDeviceNoAnswer(dev io.ReadWriteCloser, chunks [][64]byte) {
    for _, element := range chunks {
        _, _ = dev.Write(element[:])
    }
}
func sendToDevice(dev io.ReadWriteCloser, chunks [][64]byte) wire.Message {
    for _, element := range chunks {
        _, _ = dev.Write(element[:])
    }
    var msg wire.Message
    msg.ReadFrom(dev)
    return msg
}

func makeTrezorHeader(data []byte, msgID messages.MessageType) []byte {
	header := new(bytes.Buffer)
	binary.Write(header, binary.BigEndian, []byte("?##"))
	binary.Write(header, binary.BigEndian, uint16(msgID))
	binary.Write(header, binary.BigEndian, uint32(len(data)))
	binary.Write(header, binary.BigEndian, []byte("\n"))
	return header.Bytes()
}

func makeTrezorMessage(data []byte, msgID messages.MessageType) [][64]byte {
	message := new(bytes.Buffer)
	binary.Write(message, binary.BigEndian, []byte("##"))
	binary.Write(message, binary.BigEndian, uint16(msgID))
	binary.Write(message, binary.BigEndian, uint32(len(data)))
    binary.Write(message, binary.BigEndian, []byte("\n"))
    if (len(data) > 0) {
        binary.Write(message, binary.BigEndian, data[1:])
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

func getDevice(deviceType DeviceType) (io.ReadWriteCloser) {
    var dev io.ReadWriteCloser
    switch (deviceType) {
    case DeviceTypeEmulator:
        dev, _ = getEmulatorDevice()
        break
    case DeviceTypeUsb:
        dev, _ = getUsbDevice()
        break
    }
    return dev
}

// DeviceCheckMessageSignature Check a message signature matches the given address.
func DeviceCheckMessageSignature(deviceType DeviceType, message string, signature string, address string) {
    
    dev := getDevice(deviceType)
    
    // Send CheckMessageSignature
    
	skycoinCheckMessageSignature := &messages.SkycoinCheckMessageSignature{
		Address:   proto.String(address),
		Message:   proto.String(message),
		Signature: proto.String(signature),
	}

	data, _ := proto.Marshal(skycoinCheckMessageSignature)
	chunks := makeTrezorMessage(data, messages.MessageType_MessageType_SkycoinCheckMessageSignature)
    msg := sendToDevice(dev, chunks)
    fmt.Printf("Success %d! address that issued the signature is: %s\n", msg.Kind, msg.Data)
    dev.Close()
}

func MessageButtonAck() [][64]byte{
    buttonAck := &messages.ButtonAck{}
    data, _ := proto.Marshal(buttonAck)
    chunks := makeTrezorMessage(data, messages.MessageType_MessageType_ButtonAck)
    return chunks
}

// DeviceSetMnemonic Configure the device with a mnemonic.
func DeviceSetMnemonic(deviceType DeviceType, mnemonic string) {

    dev := getDevice(deviceType)
    
    // Send SetMnemonic
    
	skycoinSetMnemonic := &messages.SetMnemonic{
		Mnemonic:     proto.String(mnemonic),
	}
	
	data, _ := proto.Marshal(skycoinSetMnemonic)
	chunks := makeTrezorMessage(data, messages.MessageType_MessageType_SetMnemonic)
	
    msg := sendToDevice(dev, chunks)
	
    fmt.Printf("Success %d! Mnemonic %s\n", msg.Kind, msg.Data)
    
    // Send ButtonAck
    chunks = MessageButtonAck()
    sendToDeviceNoAnswer(dev, chunks)

    time.Sleep(1*time.Second)
    _, err := msg.ReadFrom(dev)
	if err != nil {
        fmt.Printf(err.Error())
		return
    }
    
    fmt.Printf("MessageButtonAck Answer is: %d / %s\n", msg.Kind, msg.Data)
    dev.Close()
}

// DeviceAddressGen Ask the device to generate an address
func DeviceAddressGen(deviceType DeviceType, coinType messages.SkycoinAddressType, addressN int) (uint16, string) {

	dev := getDevice(deviceType)
	defer dev.Close()
	skycoinAddress := &messages.SkycoinAddress{
		AddressN:    proto.Uint32(uint32(addressN)),
		AddressType: coinType.Enum(),
	}
	data, _ := proto.Marshal(skycoinAddress)

	chunks := makeTrezorMessage(data, messages.MessageType_MessageType_SkycoinAddress)

    msg := sendToDevice(dev, chunks)
	if (msg.Kind == uint16(messages.MessageType_MessageType_ResponseSkycoinAddress)) {
        responseSkycoinAddress := &messages.ResponseSkycoinAddress{}
        err := proto.Unmarshal(msg.Data, responseSkycoinAddress)
        if err != nil {
            return msg.Kind, ""
        }
        return msg.Kind, responseSkycoinAddress.GetAddress()
    } 
	return msg.Kind, ""
}

// DeviceSignMessage Ask the device to sign a message using the secret key at given index.
func DeviceSignMessage(deviceType DeviceType, addressN int, message string) (uint16, []byte) {
    
    dev := getDevice(deviceType)
    
	skycoinSignMessage := &messages.SkycoinSignMessage{
		AddressN:    proto.Uint32(uint32(addressN)),
		Message:     proto.String(message),
	}

	data, _ := proto.Marshal(skycoinSignMessage)

	chunks := makeTrezorMessage(data, messages.MessageType_MessageType_SkycoinSignMessage)

    msg := sendToDevice(dev, chunks)

    dev.Close()

	return msg.Kind, msg.Data
}

// DeviceConnected check if a device is connected
func DeviceConnected(deviceType DeviceType) bool {
    dev := getDevice(deviceType)
    if dev == nil {
        return false
    }
    msgRaw := &messages.Ping{}
    data, err := proto.Marshal(msgRaw)
    chunks := makeTrezorMessage(data, messages.MessageType_MessageType_Ping)
    for _, element := range chunks {
        _, err = dev.Write(element[:])
        if err != nil {
            return false
        }
    }
    var msg wire.Message
    _, err = msg.ReadFrom(dev)
    if err != nil {
        return false
    }
    return msg.Kind == uint16(messages.MessageType_MessageType_Success)
}
