package emulatorWallet

import (
    "net"
	"github.com/wire"
)

func GetTrezorDevice() (net.Conn, error) {
	return net.Dial("udp", "127.0.0.1:21324")
}

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
