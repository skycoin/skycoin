package emulatorWallet

import (
	"fmt"
    "net"
	"github.com/golang/protobuf/proto"
	messages "github.com/skycoin/skycoin/protob"
	"github.com/wire"
)

func SendToDeviceNoAnswer(dev net.Conn, chunks [][64]byte) {
    for _, element := range chunks {
        _, _ = dev.Write(element[:])
    }
}
func SendToDevice(dev net.Conn, chunks [][64]byte) wire.Message {
    for _, element := range chunks {
        _, _ = dev.Write(element[:])
    }
    var msg wire.Message
    msg.ReadFrom(dev)
    return msg
}

// DeviceCheckMessageSignature Check a message signature matches the given address.
func DeviceCheckMessageSignature(message string, signature string, address string) {

	dev, _ := GetTrezorDevice()
    defer dev.Close();
    
    // Send CheckMessageSignature
    
	skycoinCheckMessageSignature := &messages.SkycoinCheckMessageSignature{
		Address:   proto.String(address),
		Message:   proto.String(message),
		Signature: proto.String(signature),
	}
	
	data, _ := proto.Marshal(skycoinCheckMessageSignature)
	chunks := MakeTrezorMessage(data, messages.MessageType_MessageType_SkycoinCheckMessageSignature)

	var msg wire.Message
	msg = SendToDevice(dev, chunks)

	fmt.Printf("Success %d! address that issued the signature is: %s\n", msg.Kind, msg.Data)

}

// DeviceSetMnemonic Configure the device with a mnemonic.
func DeviceSetMnemonic(mnemonic string) {

	dev, _ := GetTrezorDevice()
    defer dev.Close();
    
    // Send SetMnemonic
    
	skycoinSetMnemonic := &messages.SetMnemonic{
		Mnemonic:     proto.String(mnemonic),
	}
	
	data, _ := proto.Marshal(skycoinSetMnemonic)
	chunks := MakeTrezorMessage(data, messages.MessageType_MessageType_SetMnemonic)
	
	var msg wire.Message
	msg = SendToDevice(dev, chunks)
	
    fmt.Printf("Success %d! address that issued the signature is: %s\n", msg.Kind, msg.Data)
    
    // Send ButtonAck
    
    buttonAck := &messages.ButtonAck{}
    data, _ = proto.Marshal(buttonAck)
    chunks = MakeTrezorMessage(data, messages.MessageType_MessageType_ButtonAck)

    SendToDeviceNoAnswer(dev, chunks)
    
    _, err := msg.ReadFrom(dev)
	if err != nil {
        fmt.Printf(err.Error())
		return
    }
    
    fmt.Printf("MessageButtonAck Answer is: %d / %s\n", msg.Kind, msg.Data)
}

// DeviceAddressGen Ask the device to generate an address
func DeviceAddressGen(coinType messages.SkycoinAddressType, addressN int) {
	dev, _ := GetTrezorDevice()
    defer dev.Close();
	skycoinAddress := &messages.SkycoinAddress{
		AddressN:    proto.Uint32(uint32(addressN)),
		AddressType: coinType.Enum(),
	}
	data, _ := proto.Marshal(skycoinAddress)

	chunks := MakeTrezorMessage(data, messages.MessageType_MessageType_SkycoinAddress)

	var msg wire.Message
	msg = SendToDevice(dev, chunks)

	fmt.Printf("Success %d! Address is: %s\n", msg.Kind, msg.Data)
}

// DeviceSignMessage Ask the device to sign a message using the secret key at given index.
func DeviceSignMessage(addressN int, message string) (uint16, []byte) {
	dev, _ := GetTrezorDevice()
    defer dev.Close();
	skycoinSignMessage := &messages.SkycoinSignMessage{
		AddressN:    proto.Uint32(uint32(addressN)),
		Message:     proto.String(message),
	}

	data, _ := proto.Marshal(skycoinSignMessage)

	chunks := MakeTrezorMessage(data, messages.MessageType_MessageType_SkycoinSignMessage)

    var msg wire.Message
	msg = SendToDevice(dev, chunks)

	return msg.Kind, msg.Data
}
