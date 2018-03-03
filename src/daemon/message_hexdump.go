package daemon

import (
	"fmt"
	"reflect"
	"strconv"

	"io/ioutil"
	"os"

	"github.com/skycoin/skycoin/src/cipher/encoder"
	"github.com/skycoin/skycoin/src/daemon/gnet"
)

var registered = false

func getSliceContentsString(sl []string, offset int) string {
	var res string = ""
	var counter int = 0
	var currentOff = offset
	if offset != -1 {
		var hex = strconv.FormatInt(int64(offset), 16)
		var l = len(hex)
		for i := 0; i < 4-l; i++ {
			hex = "0" + hex
		}
		hex = "0x" + hex
		res += hex + " | "
	}
	for i := 0; i < len(sl); i++ {
		counter++
		res += sl[i] + " "
		if counter == 16 {
			res += "\n"
			currentOff += 16
			if offset != -1 {
				//res += "         " //9 spaces
				var hex = strconv.FormatInt(int64(currentOff), 16)
				var l = len(hex)
				for i := 0; i < 4-l; i++ {
					hex = "0" + hex
				}
				hex = "0x" + hex
				res += hex + " | "
			}
			counter = 0
		}
	}
	for i := 0; i < (16 - counter); i++ {
		res += "..."
	}
	res += "..."
	return res
}

func printLHexDumpWithFormat(offset int, name string, buffer []byte) {
	var hexBuff = make([]string, len(buffer))
	for i := 0; i < len(buffer); i++ {
		hexBuff[i] = strconv.FormatInt(int64(buffer[i]), 16)
	}
	for i := 0; i < len(buffer); i++ {
		if len(hexBuff[i]) == 1 {
			hexBuff[i] = "0" + hexBuff[i]
		}
	}
	fmt.Println(getSliceContentsString(hexBuff, offset), name)
}

func HexDump(message gnet.Message) string {

	//setting stdout to a temp file
	defaultStdOut := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	var messagesConfig = NewMessagesConfig()
	if registered == false {
		messagesConfig.Register()
		registered = true
	}

	var serializedMsg = gnet.EncodeMessage(message)

	printLHexDumpWithFormat(-1, "Full message", serializedMsg)

	fmt.Println("------------------------------------------------------------------------")
	var offset int = 0
	printLHexDumpWithFormat(0, "Length", serializedMsg[0:4])
	printLHexDumpWithFormat(4, "Prefix", serializedMsg[4:8])
	offset += len(serializedMsg[0:8])
	var v = reflect.Indirect(reflect.ValueOf(message))

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		v_f := v.Field(i)
		f := t.Field(i)
		if f.Tag.Get("enc") != "-" {
			if v_f.CanSet() || f.Name != "_" {
				if v.Field(i).Kind() == reflect.Slice {
					printLHexDumpWithFormat(offset, f.Name+" length", encoder.Serialize(v.Field(i).Slice(0, v.Field(i).Len()).Interface())[0:4])
					offset += len(encoder.Serialize(v.Field(i).Slice(0, v.Field(i).Len()).Interface())[0:4])

					for j := 0; j < v.Field(i).Len(); j++ {
						printLHexDumpWithFormat(offset, f.Name+"#"+strconv.Itoa(j), encoder.Serialize(v.Field(i).Slice(j, j+1).Interface()))
						offset += len(encoder.Serialize(encoder.Serialize(v.Field(i).Slice(j, j+1).Interface())))
					}
				} else {
					printLHexDumpWithFormat(offset, f.Name, encoder.Serialize(v.Field(i).Interface()))
					offset += len(encoder.Serialize(v.Field(i).Interface()))
				}
			} else {
				//don't write anything
			}
		}
	}

	printFinalHex(len(serializedMsg))

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = defaultStdOut

	var strOut string = ""

	for i := 0; i < len(out); i++ {
		strOut += string(out[i])
	}

	return strOut
}
func printFinalHex(i int) {
	var finalHex = strconv.FormatInt(int64(i), 16)
	var l = len(finalHex)
	for i := 0; i < 4-l; i++ {
		finalHex = "0" + finalHex
	}
	finalHex = "0x" + finalHex
	finalHex = finalHex + " | "
	fmt.Println(finalHex)
}
