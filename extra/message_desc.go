package main

import (

	"fmt"
	"reflect"
	dm "github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/daemon/gnet"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"crypto/sha256"
	"github.com/skycoin/skycoin/src/cipher"
	//"crypto/sha1"
	//"encoding/base64"
	//"math"
	//"log"
)

func serializeMessage(msg gnet.Message) []byte {
	t := reflect.ValueOf(msg).Elem().Type()
	msgID := gnet.MessageIDMap[t]
	//if !succ {
		//txt := "Attempted to serialize message struct not in MessageIdMap: %v"
		//logger.Panicf(txt, msg)
	//	panic(msg)//TODO: Log
	//}
	bMsg := encoder.Serialize(msg)

	// message length
	bLen := encoder.SerializeAtomic(uint32(len(bMsg) + len(msgID)))
	m := make([]byte, 0)
	m = append(m, bLen...)     // length prefix
	m = append(m, msgID[:]...) // message id
	m = append(m, bMsg...)     // message bytes
	return m
}



func HexDump(message interface{}){
	var v = reflect.ValueOf(message)

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
	// see comment for corresponding code in decoder.value()
	v_f := v.Field(i)
	f := t.Field(i)
	if f.Tag.Get("enc") != "-" {
		if v_f.CanSet() || f.Name != "_" {
			if v.Field(i).Kind() == reflect.Slice {
				fmt.Println(encoder.Serialize(v.Field(i).Slice(0, v.Field(i).Len()).Interface()),f.Name)
			} else {
				fmt.Println(encoder.Serialize(v.Field(i).Interface()),f.Name)
			}
		} else {
				//dont write anything
				//e.skip(v)
		}
	}
}

fmt.Println("_________--------_______")
var bts = serializeMessage((gnet.Message)message)
fmt.Println(bts)
}

func main() {

	h := sha256.New()
	h.Write([]byte("hello world\n"))
	h.Sum(nil)
	var sha, _ = cipher.SHA256FromHex("ffgd")
	var message = dm.NewAnnounceTxnsMessage([]cipher.SHA256 {sha})

	HexDump(message)

}