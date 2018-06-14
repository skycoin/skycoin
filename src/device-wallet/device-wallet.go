
package deviceWallet

import (
	"github.com/skycoin/skycoin/src/device-wallet/hardware-wallet/usb"
	"net"
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/wire"

    proto "github.com/proto"
	messages "github.com/skycoin/skycoin/protob"
	hardwareWallet "github.com/skycoin/skycoin/src/device-wallet/hardware-wallet"
	emulatorWallet "github.com/skycoin/skycoin/src/device-wallet/emulator-wallet"
)
// DeviceType type of device: emulated or usb
type DeviceType int32

const (
    // DeviceTypeEmulator use emulator
    DeviceTypeEmulator DeviceType = 1
    // DeviceTypeUsb use usb
	DeviceTypeUsb      DeviceType = 2 
)

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
	binary.Write(message, binary.BigEndian, data[1:])

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

func getDevice(deviceType DeviceType) (net.Conn, usb.Device) {
        
    var emulator net.Conn
    var dev usb.Device
    switch (deviceType) {
    case DeviceTypeEmulator:
        emulator, _ = emulatorWallet.GetTrezorDevice()
        break
    case DeviceTypeUsb:
        dev, _ = hardwareWallet.GetTrezorDevice()
        break
    }
    return emulator, dev
}

// DeviceCheckMessageSignature Check a message signature matches the given address.
func DeviceCheckMessageSignature(deviceType DeviceType, message string, signature string, address string) {
    
    emulator, dev := getDevice(deviceType)
    switch (deviceType) {
    case DeviceTypeEmulator:
        defer emulator.Close();
        break
    case DeviceTypeUsb:
        defer dev.Close();
        break
    }
    
    // Send CheckMessageSignature
    
	skycoinCheckMessageSignature := &messages.SkycoinCheckMessageSignature{
		Address:   proto.String(address),
		Message:   proto.String(message),
		Signature: proto.String(signature),
	}

	data, _ := proto.Marshal(skycoinCheckMessageSignature)
	chunks := makeTrezorMessage(data, messages.MessageType_MessageType_SkycoinCheckMessageSignature)
    var msg wire.Message
    switch (deviceType) {
    case DeviceTypeEmulator:
        msg = emulatorWallet.SendToDevice(emulator, chunks)
        break
    case DeviceTypeUsb:
        msg = hardwareWallet.SendToDevice(dev, chunks)
        break
    }

	fmt.Printf("Success %d! address that issued the signature is: %s\n", msg.Kind, msg.Data)

}

// DeviceSetMnemonic Configure the device with a mnemonic.
func DeviceSetMnemonic(deviceType DeviceType, mnemonic string) {

    emulator, dev := getDevice(deviceType)
    switch (deviceType) {
    case DeviceTypeEmulator:
        defer emulator.Close();
        break
    case DeviceTypeUsb:
        defer dev.Close();
        break
    }
    
    // Send SetMnemonic
    
	skycoinSetMnemonic := &messages.SetMnemonic{
		Mnemonic:     proto.String(mnemonic),
	}
	
	data, _ := proto.Marshal(skycoinSetMnemonic)
	chunks := makeTrezorMessage(data, messages.MessageType_MessageType_SetMnemonic)
	
    var msg wire.Message
    switch (deviceType) {
    case DeviceTypeEmulator:
        msg = emulatorWallet.SendToDevice(emulator, chunks)
        break
    case DeviceTypeUsb:
        msg = hardwareWallet.SendToDevice(dev, chunks)
        break
    }
	
    fmt.Printf("Success %d! address that issued the signature is: %s\n", msg.Kind, msg.Data)
    
    // Send ButtonAck
    
    buttonAck := &messages.ButtonAck{}
    data, _ = proto.Marshal(buttonAck)
    chunks = makeTrezorMessage(data, messages.MessageType_MessageType_ButtonAck)

    switch (deviceType) {
    case DeviceTypeEmulator:
        emulatorWallet.SendToDeviceNoAnswer(emulator, chunks)
        break
    case DeviceTypeUsb:
        hardwareWallet.SendToDeviceNoAnswer(dev, chunks)
        break
    }
    
    _, err := msg.ReadFrom(dev)
	if err != nil {
        fmt.Printf(err.Error())
		return
    }
    
    fmt.Printf("MessageButtonAck Answer is: %d / %s\n", msg.Kind, msg.Data)
}

// DeviceAddressGen Ask the device to generate an address
func DeviceAddressGen(deviceType DeviceType, coinType messages.SkycoinAddressType, addressN int) {

    emulator, dev := getDevice(deviceType)
    switch (deviceType) {
    case DeviceTypeEmulator:
        defer emulator.Close();
        break
    case DeviceTypeUsb:
        defer dev.Close();
        break
    }
	skycoinAddress := &messages.SkycoinAddress{
		AddressN:    proto.Uint32(uint32(addressN)),
		AddressType: coinType.Enum(),
	}
	data, _ := proto.Marshal(skycoinAddress)

	chunks := makeTrezorMessage(data, messages.MessageType_MessageType_SkycoinAddress)

    var msg wire.Message
    switch (deviceType) {
    case DeviceTypeEmulator:
        msg = emulatorWallet.SendToDevice(emulator, chunks)
        break
    case DeviceTypeUsb:
        msg = hardwareWallet.SendToDevice(dev, chunks)
        break
    }

	fmt.Printf("Success %d! Address is: %s\n", msg.Kind, msg.Data)
}

// DeviceSignMessage Ask the device to sign a message using the secret key at given index.
func DeviceSignMessage(deviceType DeviceType, addressN int, message string) (uint16, []byte) {
    
    emulator, dev := getDevice(deviceType)
    
	skycoinSignMessage := &messages.SkycoinSignMessage{
		AddressN:    proto.Uint32(uint32(addressN)),
		Message:     proto.String(message),
	}

	data, _ := proto.Marshal(skycoinSignMessage)

	chunks := makeTrezorMessage(data, messages.MessageType_MessageType_SkycoinSignMessage)

    var msg wire.Message
    switch (deviceType) {
    case DeviceTypeEmulator:
        msg = emulatorWallet.SendToDevice(emulator, chunks)
        break
    case DeviceTypeUsb:
        msg = hardwareWallet.SendToDevice(dev, chunks)
        break
    }

    switch (deviceType) {
    case DeviceTypeEmulator:
        emulator.Close();
        break
    case DeviceTypeUsb:
        dev.Close();
        break
    }

	return msg.Kind, msg.Data
}
