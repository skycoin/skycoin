package harwareWallet

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	messages "github.com/skycoin/skycoin/protob"
	"github.com/wire"
)

// DeviceCheckMessageSignature Check a message signature matches the given address.
func DeviceCheckMessageSignature(message string, signature string, address string) {

	dev, _ := GetTrezorDevice()

	skycoinCheckMessageSignature := &messages.SkycoinCheckMessageSignature{
		Address:   proto.String(address),
		Message:   proto.String(message),
		Signature: proto.String(signature),
	}

	data, _ := proto.Marshal(skycoinCheckMessageSignature)

	chunks := MakeTrezorMessage(data, messages.MessageType_MessageType_SkycoinCheckMessageSignature)

	for _, element := range chunks {
		_, _ = dev.Write(element[:])
	}

	var msg wire.Message
	msg.ReadFrom(dev)

	fmt.Printf("Success %d! address that issued the signature is: %s\n", msg.Kind, msg.Data)

}

// DeviceSetMnemonic Configure the device with a mnemonic.
func DeviceSetMnemonic(mnemonic string) {

	dev, _ := GetTrezorDevice()

	skycoinSetMnemonic := &messages.SetMnemonic{
		Mnemonic:     proto.String(mnemonic),
	}

	data, _ := proto.Marshal(skycoinSetMnemonic)

	chunks := MakeTrezorMessage(data, messages.MessageType_MessageType_SetMnemonic)


	for _, element := range chunks {
		_, _ = dev.Write(element[:])
	}

	var msg wire.Message
	msg.ReadFrom(dev)

	fmt.Printf("Success %d! address that issued the signature is: %s\n", msg.Kind, msg.Data)

}

// DeviceAddressGen Ask the device to generate an address
func DeviceAddressGen(coinType messages.SkycoinAddressType, addressN int) {
	dev, _ := GetTrezorDevice()

	skycoinAddress := &messages.SkycoinAddress{
		AddressN:    proto.Uint32(uint32(addressN)),
		AddressType: coinType.Enum(),
	}
	data, _ := proto.Marshal(skycoinAddress)

	chunks := MakeTrezorMessage(data, messages.MessageType_MessageType_SkycoinAddress)

	for _, element := range chunks {
		_, _ = dev.Write(element[:])
	}

	var msg wire.Message
	msg.ReadFrom(dev)

	fmt.Printf("Success %d! Address is: %s\n", msg.Kind, msg.Data)
}

// DeviceSignMessage Ask the device to sign a message using the secret key at given index.
func DeviceSignMessage(addressN int, message string) (uint16, []byte) {
	dev, _ := GetTrezorDevice()

	skycoinSignMessage := &messages.SkycoinSignMessage{
		AddressN:    proto.Uint32(uint32(addressN)),
		Message:     proto.String(message),
	}

	data, _ := proto.Marshal(skycoinSignMessage)

	chunks := MakeTrezorMessage(data, messages.MessageType_MessageType_SkycoinSignMessage)


	for _, element := range chunks {
		_, _ = dev.Write(element[:])
	}

	var msg wire.Message
	msg.ReadFrom(dev)

	return msg.Kind, msg.Data
}
