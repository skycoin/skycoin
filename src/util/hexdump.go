package util

import (
	"io"
	"strconv"
	"bufio"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"strings"
)

type Annotation struct {
	Name string
	Size int
}

type IAnnotationsGenerator interface {
	GenerateAnnotations() []Annotation
}

type IAnnotationsIterator interface {

	Next() (Annotation,bool)
}

func writeHexdumpMember(offset int, size int, writer io.Writer, buffer []byte, name string) {
	var hexBuff = make([]string, size)
	var j = 0
	for i := offset; i < offset + size; i++ {
		hexBuff[j] = strconv.FormatInt(int64(buffer[i]), 16)
		j++
	}
	for i := 0; i < len(hexBuff); i++ {
		if len(hexBuff[i]) == 1 {
			hexBuff[i] = "0" + hexBuff[i]
		}
	}

	var sliceContents = getSliceContentsString(hexBuff, offset)
	var serialized = encoder.Serialize(sliceContents + " " + name + "\n")

	f := bufio.NewWriter(writer)
	defer f.Flush()
	f.Write(serialized[4:])

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
			if (i != len(sl) - 1) {
				res = strings.TrimRight(res," ")
				res += "\n"
				currentOff += 16
				if offset != -1 {
					//res += "         " //9 spaces
					var hex= strconv.FormatInt(int64(currentOff), 16)
					var l= len(hex)
					for i := 0; i < 4-l; i++ {
						hex = "0" + hex
					}
					hex = "0x" + hex
					res += hex + " | "
				}
				counter = 0
			} else{
				res += "..."
				return res
			}
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

	var serialized = encoder.Serialize(finalHex)

	f := bufio.NewWriter(writer)
	defer f.Flush()
	f.Write(serialized[4:])
}



func HexDump(buffer []byte, annotations []Annotation, writer io.Writer) {
	//var serializedData = encoder.Serialize(data)

	var currentOffset = 0

	for _, element := range annotations {
		writeHexdumpMember(currentOffset,element.Size,writer,buffer,element.Name)
		currentOffset += element.Size
	}

	printFinalHex(currentOffset,writer)
}

func HexDumpFromIterator(buffer []byte, annotationsIterator IAnnotationsIterator, writer io.Writer) {
	//var serializedData = encoder.Serialize(data)

	var currentOffset = 0

	var current, valid = annotationsIterator.Next()

	for ; ;  {
		if !valid {
			break
		}
		writeHexdumpMember(currentOffset,current.Size,writer,buffer,current.Name)
		currentOffset += current.Size
		current, valid = annotationsIterator.Next()
	}



	printFinalHex(currentOffset,writer)
}