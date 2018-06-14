package harwareWallet

import (
	"fmt"
	"log"
	"time"

	"github.com/wire"
	"github.com/skycoin/skycoin/src/device-wallet/hardware-wallet/usb"
)

func GetTrezorDevice() (usb.Device, error) {
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

func SendToDeviceNoAnswer(dev usb.Device, chunks [][64]byte) {
    for _, element := range chunks {
        _, _ = dev.Write(element[:])
    }
}
func SendToDevice(dev usb.Device, chunks [][64]byte) wire.Message {
    for _, element := range chunks {
        _, _ = dev.Write(element[:])
    }
    var msg wire.Message
    msg.ReadFrom(dev)
    return msg
}
