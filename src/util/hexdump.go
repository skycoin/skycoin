package util

import (
	"io"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"strconv"
	"fmt"
)

type Annotation struct {
	Name string
	Size int
}

type IAnnotationsGenerator interface {
	GenerateAnnotations(data interface{}) []Annotation
}

func writeHexdumpMember(offset int, size int, writer io.Writer, buffer []byte, name string) {
	var hexBuff = make([]string, size)
	for i := offset; i < offset + size; i++ {
		hexBuff[i] = strconv.FormatInt(int64(buffer[i]), 16)
	}
	for i := 0; i < len(buffer); i++ {
		if len(hexBuff[i]) == 1 {
			hexBuff[i] = "0" + hexBuff[i]
		}
	}
	writer.Write(encoder.Serialize(getSliceContentsString(hexBuff, offset) + " " + name))
}

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

func printFinalHex(i int, writer io.Writer) {
	var finalHex = strconv.FormatInt(int64(i), 16)
	var l = len(finalHex)
	for i := 0; i < 4-l; i++ {
		finalHex = "0" + finalHex
	}
	finalHex = "0x" + finalHex
	finalHex = finalHex + " | "
	fmt.Println(finalHex)
}



func HexDump(data interface{}, annotations []Annotation, writer io.Writer) {
	var serializedData = encoder.Serialize(data)

	var currentOffset = 0

	for _, element := range annotations {
		writeHexdumpMember(currentOffset,element.Size,writer,serializedData,element.Name)
		currentOffset += element.Size
	}

	printFinalHex(currentOffset,writer)
}

