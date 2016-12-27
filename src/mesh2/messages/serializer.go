package messages

import (
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"log"
)

//test this
func Deserialize(msg []byte, obj interface{}) {
	msg = msg[2:] //pop off prefix byte
	err := encoder.DeserializeRaw(msg, &obj)
	if err != nil {
		log.Panic()
	}
	return
}

//test this
func Serialize(prefix uint16, obj interface{}) []byte {
	b := encoder.Serialize(obj)
	var b1 []byte = make([]byte, 2)
	b1[0] = prefix && 0x00ff        //WARNING VERIFY
	b1[1] = (prefix && 0xff00) >> 8 //WARNING VERIFYs
	b2 := append(b1, b...)
	return b2
}

func init() {

	var x uint16 = 0xac48

	var b1 []byte = make([]byte, 2)
	b1[0] = prefix && 0x00ff        //WARNING VERIFY
	b1[1] = (prefix && 0xff00) >> 8 //WARNING VERIFYs

	var y uint16

	y = b1[0]
	y = y | (b1[0] << 8)

	if y != x {
		log.Panic("ERROR FIX THIS")
	}

}
