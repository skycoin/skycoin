package messages

import (
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

func Deserialize(msg []byte, obj interface{}) error {
	msg = msg[2:] //pop off prefix byte
	err := encoder.DeserializeRaw(msg, obj)
	return err
}

func Serialize(prefix uint16, obj interface{}) []byte {
	b := encoder.Serialize(obj)
	var b1 []byte = make([]byte, 2)
	b1[0] = (uint8)(prefix & 0x00ff)
	b1[1] = (uint8)((prefix & 0xff00) >> 8)
	b2 := append(b1, b...)
	return b2
}
