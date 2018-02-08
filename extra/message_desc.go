package extra

import (

	"fmt"
	"reflect"
	//dm "github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/daemon/gnet"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	//"crypto/sha256"
	//"github.com/skycoin/skycoin/src/cipher"
	//"crypto/sha1"
	//"encoding/base64"
	//"math"
	//"log"
	"strconv"
	"encoding/binary"
	"bytes"
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

func getSliceContentsString(sl []string,offset int) string {
	var res string = ""
	var counter int = 0
	if offset != -1{
		var hex = strconv.FormatInt(int64(offset),16)
		var l = len(hex)
		for i := 0; i < 4 - l; i++ {
			hex = "0" + hex
		}
		hex = "0x" + hex
		res += hex +  " "
	}
	for i := 0; i < len(sl); i++ {
		counter++
		res += sl[i] + " "
		if counter == 16{
			res += "\n"
			if offset != -1 {
				res += "       "//7 spaces
			}
			counter = 0
		}
	}
	for i := 0; i < (16-counter); i++ {
		res += "..."
	}
	res += "..."
	return res
}

func printLHexDumpWithFormat(offset int, name string, buffer []byte){
	var hexBuff = make([]string,len(buffer))
	for  i := 0; i < len(buffer) ;i++  {
		hexBuff[i] = strconv.FormatInt(int64(buffer[i]),16)
	}
	for  i := 0; i < len(buffer) ;i++  {
		if len(hexBuff[i]) == 1{
			hexBuff[i] = "0" + hexBuff[i]
		}
	}
	fmt.Println(getSliceContentsString(hexBuff,offset),name)
}

func HexDump(message gnet.Message){
	var serializedMsg = serializeMessage(message)

	printLHexDumpWithFormat(-1,"Full message",serializedMsg)

	fmt.Println("------------------------------------------------------------------------")
	var offset int = 0
	printLHexDumpWithFormat(0,"Prefix",serializedMsg[0:8])
	offset += len(serializedMsg[0:8])
	var v = reflect.Indirect(reflect.ValueOf(message))

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
	v_f := v.Field(i)
	f := t.Field(i)
	if f.Tag.Get("enc") != "-" {
		if v_f.CanSet() || f.Name != "_" {
			if v.Field(i).Kind() == reflect.Slice {
				printLHexDumpWithFormat(offset,f.Name + " header",encoder.Serialize(v.Field(i).Slice(0, v.Field(i).Len()).Interface())[0:4])
				offset += len(encoder.Serialize(v.Field(i).Slice(0, v.Field(i).Len()).Interface())[0:4])

				for j := 0; j < v.Field(i).Len(); j++ {
					printLHexDumpWithFormat(offset,f.Name + "#" + strconv.Itoa(j),encoder.Serialize(v.Field(i).Slice(j, j+1).Interface()))
					offset += len(encoder.Serialize(encoder.Serialize(v.Field(i).Slice(j, j+1).Interface())))
					}
			} else {
				printLHexDumpWithFormat(offset,f.Name,encoder.Serialize(v.Field(i).Interface()))
				offset += len(encoder.Serialize(v.Field(i).Interface()))
				}
		} else {
				//dont write anything
				//e.skip(v)
		}
	}
}


}