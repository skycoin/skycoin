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

func byteSliceToHexStrSlice(bts []byte) []string {

	var hex= make([]string, len(bts))
	for i := 0; i < len(bts); i++ {
		var n int64
		_ = binary.Read(bytes.NewReader(bts), binary.BigEndian, &n)
		hex[i] = strconv.FormatInt(n,16)


	}
	return hex
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
	//if offset == -1{
		//fmt.Println(buffer,name)
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
	//} else{
	//	fmt.Println(offset,buffer,name)
	//}

}

func HexDump(message *dm.AnnounceTxnsMessage){


	var serializedMsg = serializeMessage(message)

	printLHexDumpWithFormat(-1,"Full message",serializedMsg)

	fmt.Println("------------------------------------------------------------------------")

	printLHexDumpWithFormat(0,"Prefix",serializedMsg[0:8])
	//fmt.Println(serializedMsg[0:8],"Prefix")

	var v = reflect.ValueOf(*message)

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
	// see comment for corresponding code in decoder.value()
	v_f := v.Field(i)
	f := t.Field(i)
	if f.Tag.Get("enc") != "-" {
		if v_f.CanSet() || f.Name != "_" {
			if v.Field(i).Kind() == reflect.Slice {
				printLHexDumpWithFormat(0,f.Name + " header",encoder.Serialize(v.Field(i).Slice(0, v.Field(i).Len()).Interface())[0:4])
				//fmt.Println(encoder.Serialize(v.Field(i).Slice(0, v.Field(i).Len()).Interface())[0:4],f.Name,"header")
				for j := 0; j < v.Field(i).Len(); j++ {
					printLHexDumpWithFormat(0,f.Name + "#" + strconv.Itoa(j),encoder.Serialize(v.Field(i).Slice(j, j+1).Interface()))
					//fmt.Println(encoder.Serialize(v.Field(i).Slice(j, j+1).Interface()),f.Name + "#" + strconv.Itoa(j))
				}
			} else {
				printLHexDumpWithFormat(0,f.Name,encoder.Serialize(v.Field(i).Interface()))
				//fmt.Println(encoder.Serialize(v.Field(i).Interface()),f.Name)
			}
		} else {
				//dont write anything
				//e.skip(v)
		}
	}
}


}

func main() {

	h := sha256.New()
	h.Write([]byte("hello world\n"))
	h.Sum(nil)
	var sha, _ = cipher.SHA256FromHex("ffgd")
	var message = dm.NewAnnounceTxnsMessage([]cipher.SHA256 {sha,sha,sha})

	HexDump(message)

}