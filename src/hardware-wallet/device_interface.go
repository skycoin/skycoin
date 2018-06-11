package harwareWallet

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	messages "github.com/skycoin/skycoin/protob"
	"github.com/wire"
)

// DeviceSignMessageCmd Ask the device to sign a message using the secret key at given index.
func DeviceSignMessageCmd(addressN int, message string) {
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

	fmt.Printf("Success %d! address that issued the signature is: %s\n", msg.Kind, msg.Data)
}
