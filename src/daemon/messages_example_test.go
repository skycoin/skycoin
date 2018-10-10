package daemon

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"

	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/daemon/gnet"
	"github.com/skycoin/skycoin/src/daemon/pex"
	"github.com/stretchr/testify/require"
)

func setupMsgEncoding() {
	gnet.EraseMessages()
	var messagesConfig = NewMessagesConfig()
	messagesConfig.Register()
}

/**************************************
 *
 * Test helpers
 *
 *************************************/

// Annotation : Denotes a chunk of data to be dumped
type Annotation struct {
	Name string
	Size int
}

// IAnnotationsIterator : Interface to implement by types to use HexDumpFromIterator
type IAnnotationsIterator interface {
	Next() (Annotation, bool)
}

func writeHexdumpMember(offset int, size int, writer io.Writer, buffer []byte, name string) error {
	var hexBuff = make([]string, size)
	var j = 0
	if offset+size > len(buffer) {
		panic(encoder.ErrBufferUnderflow)
	}
	for i := offset; i < offset+size; i++ {
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
	_, err := f.Write(serialized[4:])
	return err
}

func getSliceContentsString(sl []string, offset int) string {
	var res = ""
	var counter = 0
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
			if i != len(sl)-1 {
				res = strings.TrimRight(res, " ")
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
			} else {
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

func printFinalHex(i int, writer io.Writer) error {
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

	_, err := f.Write(serialized[4:])
	return err
}

// NewFromIterator : Returns hexdump of buffer according to annotationsIterator, via writer
func NewFromIterator(buffer []byte, annotationsIterator IAnnotationsIterator, writer io.Writer) error {
	var currentOffset = 0

	var current, valid = annotationsIterator.Next()

	for {
		if !valid {
			break
		}
		if err := writeHexdumpMember(currentOffset, current.Size, writer, buffer, current.Name); err != nil {
			return err
		}
		currentOffset += current.Size
		current, valid = annotationsIterator.Next()
	}

	return printFinalHex(currentOffset, writer)
}

// MessagesAnnotationsIterator : Implementation of IAnnotationsIterator for type gnet.Message
type MessagesAnnotationsIterator struct {
	Message      gnet.Message
	LengthCalled bool
	PrefixCalled bool
	CurrentField int
	MaxField     int
	CurrentIndex int
}

// NewMessagesAnnotationsIterator : Initializes struct MessagesAnnotationsIterator
func NewMessagesAnnotationsIterator(message gnet.Message) MessagesAnnotationsIterator {
	var mai = MessagesAnnotationsIterator{}
	mai.Message = message
	mai.LengthCalled = false
	mai.PrefixCalled = false
	mai.CurrentField = 0
	mai.CurrentIndex = -1

	mai.MaxField = reflect.Indirect(reflect.ValueOf(mai.Message)).NumField()

	return mai
}

// Next : Yields next element of MessagesAnnotationsIterator
func (mai *MessagesAnnotationsIterator) Next() (Annotation, bool) {
	if !mai.LengthCalled {
		mai.LengthCalled = true
		return Annotation{Size: 4, Name: "Length"}, true
	}
	if !mai.PrefixCalled {
		mai.PrefixCalled = true
		return Annotation{Size: 4, Name: "Prefix"}, true

	}
	if mai.CurrentField >= mai.MaxField {
		return Annotation{}, false
	}

	var i = mai.CurrentField
	var j = mai.CurrentIndex

	var v = reflect.Indirect(reflect.ValueOf(mai.Message))
	t := v.Type()
	vF := v.Field(i)
	f := t.Field(i)
	for f.PkgPath != "" && i < mai.MaxField {
		i++
		mai.CurrentField++
		mai.CurrentIndex = -1
		j = -1
		if i < mai.MaxField {
			f = t.Field(i)
			if f.Type.Kind() == reflect.Slice {
				if _, omitempty := encoder.ParseTag(f.Tag.Get("enc")); omitempty {
					if i == mai.MaxField-1 {
						vF = v.Field(i)
						if vF.Len() == 0 {
							// Last field is empty slice. Nothing further tokens
							return Annotation{}, false
						}
					} else {
						panic(encoder.ErrInvalidOmitEmpty)
					}
				}
			}
		} else {
			return Annotation{}, false
		}
	}
	if f.Tag.Get("enc") != "-" {
		if vF.CanSet() || f.Name != "_" {
			if v.Field(i).Kind() == reflect.Slice {
				if mai.CurrentIndex == -1 {
					mai.CurrentIndex = 0
					return Annotation{Size: 4, Name: f.Name + " length"}, true
				}
				sliceLen := v.Field(i).Len()
				mai.CurrentIndex++
				if mai.CurrentIndex < sliceLen {
					// Emit annotation for slice item
					return Annotation{Size: len(encoder.Serialize(v.Field(i).Slice(j, j+1).Interface())[4:]), Name: f.Name + "[" + strconv.Itoa(j) + "]"}, true
				}
				// No more annotation tokens for current slice field
				mai.CurrentIndex = -1
				mai.CurrentField++
				if sliceLen > 0 {
					// Emit annotation for last item
					return Annotation{Size: len(encoder.Serialize(v.Field(i).Slice(j, j+1).Interface())[4:]), Name: f.Name + "[" + strconv.Itoa(j) + "]"}, true
				}
				// Zero length slice. Start over
				return mai.Next()
			}

			mai.CurrentField++
			return Annotation{Size: len(encoder.Serialize(v.Field(i).Interface())), Name: f.Name}, true

		}
	}

	return Annotation{}, false
}

//***Lazy iterator for DeepMessages

// DeepMessagesAnnotationsIterator : Implementation of deep IAnnotationsIterator for type gnet.Message
type DeepMessagesAnnotationsIterator struct {
	Message         gnet.Message
	LengthCalled    bool
	PrefixCalled    bool
	CurrentField    int
	CurrentIndex    int
	CurrentDepth    int
	MaxDepth        int
	CurrentTypology []reflect.Kind
	CurrentPosition []int
	CurrentMax      []int
	oldObj          reflect.Value
	obj             reflect.Value
}

// NewDeepMessagesAnnotationsIterator : Initializes struct DeepMessagesAnnotationsIterator
func NewDeepMessagesAnnotationsIterator(message gnet.Message, depth int) DeepMessagesAnnotationsIterator {
	var dmai = DeepMessagesAnnotationsIterator{}
	dmai.Message = message
	dmai.LengthCalled = false
	dmai.PrefixCalled = false
	dmai.CurrentField = 0
	dmai.CurrentDepth = 1
	dmai.MaxDepth = depth
	dmai.CurrentTypology = make([]reflect.Kind, 1)
	dmai.CurrentTypology[0] = reflect.Struct
	dmai.CurrentPosition = make([]int, 1)
	dmai.CurrentPosition[0] = 0
	dmai.CurrentMax = make([]int, 1)
	dmai.CurrentMax[0] = reflect.Indirect(reflect.ValueOf(dmai.Message)).NumField()
	dmai.oldObj = reflect.Indirect(reflect.ValueOf(dmai.Message))
	dmai.obj = reflect.Indirect(reflect.ValueOf(dmai.Message)).Field(0)
	return dmai
}

// Next : Yields next element of DeepMessagesAnnotationsIterator
func (dmai *DeepMessagesAnnotationsIterator) Next() (Annotation, bool) {
	if !dmai.LengthCalled {
		dmai.LengthCalled = true
		return Annotation{Size: 4, Name: "Length"}, true
	}
	if !dmai.PrefixCalled {
		dmai.PrefixCalled = true
		return Annotation{Size: 4, Name: "Prefix"}, true

	}
	if dmai.CurrentPosition[0] >= dmai.CurrentMax[0] {
		return Annotation{}, false
	}

	dmai.obj = reflect.Indirect(reflect.ValueOf(dmai.Message))
	depth := 1
	for (len(dmai.CurrentTypology) <= dmai.MaxDepth) && (dmai.obj.Kind() == reflect.Struct || dmai.obj.Kind() == reflect.Slice) && (len(dmai.CurrentPosition) != 1 || dmai.oldObj.Kind() == reflect.Struct && dmai.oldObj.Type().Field(dmai.CurrentPosition[len(dmai.CurrentPosition)-1]).Tag.Get("enc") != "-" && (dmai.oldObj.Field(dmai.CurrentPosition[len(dmai.CurrentPosition)-1]).CanSet() || dmai.oldObj.Type().Field(dmai.CurrentPosition[len(dmai.CurrentPosition)-1]).Name != "_") && !strings.Contains(dmai.oldObj.Type().Field(dmai.CurrentPosition[len(dmai.CurrentPosition)-1]).Tag.Get("enc"), "omitempty")) {

		if dmai.obj.Kind() == reflect.Struct {
			if len(dmai.CurrentPosition) >= depth {
				dmai.oldObj = dmai.obj
				dmai.obj = dmai.obj.Field(dmai.CurrentPosition[depth-1])
			} else {
				dmai.CurrentTypology = append(dmai.CurrentTypology, reflect.Struct)
				dmai.CurrentPosition = append(dmai.CurrentPosition, 0)
				dmai.CurrentMax = append(dmai.CurrentMax, dmai.obj.Type().NumField())
				dmai.oldObj = dmai.obj
				dmai.obj = dmai.obj.Field(0)
				dmai.CurrentDepth++
			}
		} else if dmai.obj.Kind() == reflect.Slice {
			if len(dmai.CurrentPosition) >= depth {
				dmai.oldObj = dmai.obj
				dmai.obj = dmai.obj.Index(dmai.CurrentPosition[depth-1])
			} else {
				dmai.CurrentTypology = append(dmai.CurrentTypology, reflect.Slice)
				dmai.CurrentPosition = append(dmai.CurrentPosition, 0)
				dmai.CurrentMax = append(dmai.CurrentMax, dmai.obj.Len())
				var _, fieldName = getCurrentObj(dmai.Message, dmai.CurrentDepth, dmai.CurrentTypology[0:len(dmai.CurrentTypology)-1], dmai.CurrentPosition[0:len(dmai.CurrentPosition)-1])
				dmai.oldObj = dmai.obj
				dmai.obj = dmai.obj.Index(0)
				dmai.CurrentDepth++
				return Annotation{Size: 4, Name: fieldName + " length"}, true
			}
		}
		depth++
	}

	var objName string
	dmai.obj, objName = getCurrentObj(dmai.Message, dmai.CurrentDepth, dmai.CurrentTypology, dmai.CurrentPosition)

	for len(dmai.CurrentPosition) != 1 && (dmai.CurrentPosition[len(dmai.CurrentPosition)-1] == dmai.CurrentMax[len(dmai.CurrentPosition)-1]-1) {
		dmai.CurrentPosition = dmai.CurrentPosition[0 : len(dmai.CurrentPosition)-1]
		dmai.CurrentTypology = dmai.CurrentTypology[0 : len(dmai.CurrentTypology)-1]
		dmai.CurrentMax = dmai.CurrentMax[0 : len(dmai.CurrentMax)-1]
		dmai.CurrentDepth--
	}

	var encTagIsDash bool
	var fieldIsSettable bool
	var fieldNameIsntUnderscore bool
	var encTagContainsOmitempty bool

	if len(dmai.CurrentPosition) == 1 {
		msg := reflect.Indirect(reflect.ValueOf(dmai.Message))
		encTagIsDash = msg.Type().Field(dmai.CurrentPosition[len(dmai.CurrentPosition)-1]).Tag.Get("enc") == "-"
		fieldIsSettable = msg.Field(dmai.CurrentPosition[len(dmai.CurrentPosition)-1]).CanSet()
		fieldNameIsntUnderscore = msg.Type().Field(dmai.CurrentPosition[len(dmai.CurrentPosition)-1]).Name != "_"
		encTagContainsOmitempty = strings.Contains(msg.Type().Field(dmai.CurrentPosition[len(dmai.CurrentPosition)-1]).Tag.Get("enc"), "omitempty")
	}

	dmai.CurrentPosition[len(dmai.CurrentPosition)-1]++

	if len(dmai.CurrentPosition) != 1 || (!encTagIsDash && (fieldIsSettable || fieldNameIsntUnderscore) && !encTagContainsOmitempty) {
		return Annotation{Name: objName, Size: len(encoder.Serialize(dmai.obj.Interface()))}, true
	}
	return dmai.Next()

}

func getCurrentObj(message gnet.Message, currentDepth int, currentTypology []reflect.Kind, currentPosition []int) (reflect.Value, string) {
	var obj = reflect.Indirect(reflect.ValueOf(message))
	var name string
	for i := 0; i < currentDepth; i++ {
		if currentTypology[i] == reflect.Slice {
			obj = obj.Index(currentPosition[i])
			name = name + "[" + strconv.Itoa(currentPosition[i]) + "]"
		}
		if currentTypology[i] == reflect.Struct {
			name = name + "." + obj.Type().Field(currentPosition[i]).Name
			obj = obj.Field(currentPosition[i])
		}
	}
	return obj, name
}

/**************************************
 *
 * Test cases
 *
 *************************************/
const (
	// seed = "x"
	// skhex = "2a2b56833e4fd6e30f4a653e62f4eb7761d7d383baa52496df08c237a2dfff98"
	// pkhex = "0268bc0885c2b9dc58199403bc5ac529e2962b326188cd09caabeb9d5e61f26c39"
	//
	// buffer => 0x00, 0x01, ... 0xff
	sig1hex = "03213fdd6ddf860e4053e1a97e4276d63454f5195a8321357004d52cdbbfd3886fc7ad3f3f63b65d4a879ce3086daeb3e54a93d3c2f96a5061f9bc493683ca8e01"
	// buffer => 0x01, 0x01, ... 0xff
	sig2hex = "aeabca1decb2ab5ba8af93c6f033de6cbf1d50314df275f940778bc720433cbc2194aacfa243cc2a21f85f2fff71d3166d1875e1981a0da5a23d289681fc1fa300"
	// buffer => 0x02, 0x01, ... 0xff
	sig3hex = "6765278afc9f3e0ffb95cbb3f01872e92ed1d51e7a83d16d499e9597e24fa6f308d885c831c46b699ad67b2fdd2f762fd65f4fbcf66f981d76a9adfe420d161401"
	// buffer => 0x03, 0x01, ... 0xff
	sig4hex = "5edb9bbd4a7820f45391f08e7572c91fb61be7106ddb5346ac00596c8992cfa23aac06469aec55337102418ceb11529b3b30578402e58c2ff66dcd8bbf3521e901"
)

var hashes = []cipher.SHA256{
	// buffer => 0x00, 0x01, ... 0xff
	getSHAFromHex("40aff2e9d2d8922e47afd4648e6967497158785fbd1da870e7110266bf944880"),
	// buffer => 0x01, 0x01, ... 0xff
	getSHAFromHex("7bb462c3bd371dd81c06ad1d2b635971cb56eb22233dfc9febe83e44c840b8d7"),
	// buffer => 0x02, 0x01, ... 0xff
	getSHAFromHex("e75ac801c13f3da9c7a124ca313be2a373f64ad97c58a1b6febc0e0ca5c5c873"),
	// buffer => 0x03, 0x01, ... 0xff
	getSHAFromHex("f4457de9f5a5942e076a7f2b28e1842ab61f1bfc39e4ca557536600fd64209f6"),
	// buffer => 0x04, 0x01, ... 0xff
	getSHAFromHex("c4220982f988b625d0afc12c7ffd06a7fe89bbe6602c1f20d908913fe9381047"),
	// buffer => 0x05, 0x01, ... 0xff
	getSHAFromHex("3862f193a1564e5e260f827da8e169cad811d81d6a7c4ffd661c00b199941781"),
	// buffer => 0x06, 0x01, ... 0xff
	getSHAFromHex("05640e4480739e879757b0a2d1bd59dea7dfccfef3df75a1830a502001106721"),
	// buffer => 0x07, 0x01, ... 0xff
	getSHAFromHex("8a5dbfbb7e6466495e30781c1540b5e398e0844f60c91ec6789d4bbb367e33a6"),
	// buffer => 0x08, 0x01, ... 0xff
	getSHAFromHex("1c1d7dbfd7ba2bb1aa9b56edae26ea565cbf72f98cc6a62c729723cbc0750d3b"),
	// buffer => 0x09, 0x01, ... 0xff
	getSHAFromHex("66dd3fc45be9b4fbb1fbed2be5de4e8a479f6638adfe4675b8544ae84eca3f75"),
}

var secKey1 = cipher.MustNewSecKey([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32})
var secKey2 = cipher.MustNewSecKey([]byte{33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64})
var secKey3 = cipher.MustNewSecKey([]byte{65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96})
var secKey4 = cipher.MustNewSecKey([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96})

var addresses = []cipher.Address{
	cipher.MustAddressFromSecKey(secKey1),
	cipher.MustAddressFromSecKey(secKey2),
	cipher.MustAddressFromSecKey(secKey3),
	cipher.MustAddressFromSecKey(secKey4),
}

var sig1, _ = cipher.SigFromHex(sig1hex) // nolint: errcheck
var sig2, _ = cipher.SigFromHex(sig2hex) // nolint: errcheck
var sig3, _ = cipher.SigFromHex(sig3hex) // nolint: errcheck
var sig4, _ = cipher.SigFromHex(sig4hex) // nolint: errcheck

var sigs = []cipher.Sig{
	sig1,
	sig2,
	sig3,
	sig4,
}

func getSHAFromHex(hex string) cipher.SHA256 {
	var sha, _ = cipher.SHA256FromHex(hex) // nolint: errcheck
	return sha
}

type EmptySliceStruct struct {
	A uint8
	e int16
	B string
	C int32
	D []byte
	f rune
}

func (m *EmptySliceStruct) Handle(mc *gnet.MessageContext, daemon interface{}) error {
	// Do nothing
	return nil
}

func ExampleEmptySliceStruct() {
	defer gnet.EraseMessages()
	setupMsgEncoding()
	gnet.RegisterMessage(gnet.MessagePrefixFromString("TEST"), EmptySliceStruct{})
	gnet.VerifyMessages()
	var message = EmptySliceStruct{
		0x01,
		0x2345,
		"",
		0x6789ABCD,
		nil,
		'a',
	}
	var mai = NewMessagesAnnotationsIterator(&message)
	w := bufio.NewWriter(os.Stdout)
	NewFromIterator(gnet.EncodeMessage(&message), &mai, w) // nolint: errcheck
	// Output:
	// 0x0000 | 11 00 00 00 ....................................... Length
	// 0x0004 | 54 45 53 54 ....................................... Prefix
	// 0x0008 | 01 ................................................ A
	// 0x0009 | 00 00 00 00 ....................................... B
	// 0x000d | cd ab 89 67 ....................................... C
	// 0x0011 | 00 00 00 00 ....................................... D length
	// 0x0015 |
}

type OmitEmptySliceTestStruct struct {
	A uint8
	B []byte
	c rune
	D []byte `enc:",omitempty"`
}

func (m *OmitEmptySliceTestStruct) Handle(mc *gnet.MessageContext, daemon interface{}) error {
	// Do nothing
	return nil
}

func ExampleOmitEmptySliceTestStruct() {
	defer gnet.EraseMessages()
	setupMsgEncoding()
	gnet.RegisterMessage(gnet.MessagePrefixFromString("TEST"), OmitEmptySliceTestStruct{})
	gnet.VerifyMessages()
	var message = OmitEmptySliceTestStruct{
		0x01,
		nil,
		'a',
		nil,
	}
	var mai = NewMessagesAnnotationsIterator(&message)
	w := bufio.NewWriter(os.Stdout)
	err := NewFromIterator(gnet.EncodeMessage(&message), &mai, w)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
	// 0x0000 | 09 00 00 00 ....................................... Length
	// 0x0004 | 54 45 53 54 ....................................... Prefix
	// 0x0008 | 01 ................................................ A
	// 0x0009 | 00 00 00 00 ....................................... B length
	// 0x000d |
}

func ExampleIntroductionMessage() {
	defer gnet.EraseMessages()
	setupMsgEncoding()
	var message = NewIntroductionMessage(1234, 5, 7890, nil)
	fmt.Println("IntroductionMessage:")
	var dmai = NewDeepMessagesAnnotationsIterator(message, 1)
	w := bufio.NewWriter(os.Stdout)
	NewFromIterator(gnet.EncodeMessage(message), &dmai, w) // nolint: errcheck
	// Output:
	// IntroductionMessage:
	// 0x0000 | 0e 00 00 00 ....................................... Length
	// 0x0004 | 49 4e 54 52 ....................................... Prefix
	// 0x0008 | d2 04 00 00 ....................................... .Mirror
	// 0x000c | d2 1e ............................................. .Port
	// 0x000e | 05 00 00 00 ....................................... .Version
	// 0x0012 |
}

func ExampleGetPeersMessage() {
	defer gnet.EraseMessages()
	setupMsgEncoding()
	var message = NewGetPeersMessage()
	fmt.Println("GetPeersMessage:")
	var dmai = NewDeepMessagesAnnotationsIterator(message, 1)
	w := bufio.NewWriter(os.Stdout)
	NewFromIterator(gnet.EncodeMessage(message), &dmai, w) // nolint: errcheck
	// Output:
	// GetPeersMessage:
	// 0x0000 | 04 00 00 00 ....................................... Length
	// 0x0004 | 47 45 54 50 ....................................... Prefix
	// 0x0008 |
}

func ExampleGivePeersMessage() {
	defer gnet.EraseMessages()
	setupMsgEncoding()
	var peers = make([]pex.Peer, 3)
	var peer0 = *pex.NewPeer("118.178.135.93:6000")
	var peer1 = *pex.NewPeer("47.88.33.156:6000")
	var peer2 = *pex.NewPeer("121.41.103.148:6000")
	peers = append(peers, peer0, peer1, peer2)
	var message = NewGivePeersMessage(peers)
	fmt.Println("GivePeersMessage:")
	var dmai = NewDeepMessagesAnnotationsIterator(message, 3)
	w := bufio.NewWriter(os.Stdout)
	NewFromIterator(gnet.EncodeMessage(message), &dmai, w) // nolint: errcheck
	// Output:
	// GivePeersMessage:
	// 0x0000 | 1a 00 00 00 ....................................... Length
	// 0x0004 | 47 49 56 50 ....................................... Prefix
	// 0x0008 | 03 00 00 00 ....................................... .Peers length
	// 0x000c | 5d 87 b2 76 ....................................... .Peers[0].IP
	// 0x0010 | 70 17 ............................................. .Peers[0].Port
	// 0x0012 | 9c 21 58 2f ....................................... .Peers[1].IP
	// 0x0016 | 70 17 ............................................. .Peers[1].Port
	// 0x0018 | 94 67 29 79 ....................................... .Peers[2].IP
	// 0x001c | 70 17 ............................................. .Peers[2].Port
	// 0x001e |
}

func ExampleGetBlocksMessage() {
	defer gnet.EraseMessages()
	setupMsgEncoding()
	var message = NewGetBlocksMessage(1234, 5678)
	fmt.Println("GetBlocksMessage:")
	var dmai = NewDeepMessagesAnnotationsIterator(message, 6)
	w := bufio.NewWriter(os.Stdout)
	NewFromIterator(gnet.EncodeMessage(message), &dmai, w) // nolint: errcheck
	// Output:
	// GetBlocksMessage:
	// 0x0000 | 14 00 00 00 ....................................... Length
	// 0x0004 | 47 45 54 42 ....................................... Prefix
	// 0x0008 | d2 04 00 00 00 00 00 00 ........................... .LastBlock
	// 0x0010 | 2e 16 00 00 00 00 00 00 ........................... .RequestedBlocks
	// 0x0018 |
}

func ExampleGiveBlocksMessage() {
	defer gnet.EraseMessages()
	setupMsgEncoding()
	var blocks = make([]coin.SignedBlock, 0)

	var transactions1 = make([]coin.Transaction, 0)
	var transactionOutput1_1 = coin.TransactionOutput{Address: addresses[0], Coins: 2, Hours: 4}
	var transactionOutput1_2 = coin.TransactionOutput{Address: addresses[1], Coins: 1, Hours: 2}
	var transactionOutput2_1 = coin.TransactionOutput{Address: addresses[2], Coins: 5, Hours: 5}
	var transactionOutput2_2 = coin.TransactionOutput{Address: addresses[3], Coins: 3, Hours: 3}
	var transaction1_1 = coin.Transaction{Type: 2, In: hashes[0:2], Out: []coin.TransactionOutput{transactionOutput1_1, transactionOutput1_2}, InnerHash: hashes[3], Length: 2, Sigs: sigs[0:2]}
	var transaction1_2 = coin.Transaction{Type: 1, In: hashes[3:4], Out: []coin.TransactionOutput{transactionOutput2_1, transactionOutput2_2}, InnerHash: hashes[4], Length: 2, Sigs: sigs[1:3]}
	transactions1 = append(transactions1, transaction1_1, transaction1_2)

	var body1 = coin.BlockBody{
		Transactions: transactions1,
	}
	var block1 = coin.Block{
		Body: body1,
		Head: coin.BlockHeader{
			Version:  0x02,
			Time:     100,
			BkSeq:    1,
			Fee:      10,
			PrevHash: hashes[0],
			BodyHash: body1.Hash(),
		}}
	var sig = sigs[3]
	var signedBlock = coin.SignedBlock{
		Sig:   sig,
		Block: block1,
	}
	blocks = append(blocks, signedBlock)
	var message = NewGiveBlocksMessage(blocks)
	fmt.Println("GiveBlocksMessage:")
	var dmai = NewDeepMessagesAnnotationsIterator(message, 8)
	w := bufio.NewWriter(os.Stdout)
	NewFromIterator(gnet.EncodeMessage(message), &dmai, w) // nolint: errcheck
	// Output:
	// GiveBlocksMessage:
	// 0x0000 | 23 03 00 00 ....................................... Length
	// 0x0004 | 47 49 56 42 ....................................... Prefix
	// 0x0008 | 01 00 00 00 ....................................... .Blocks length
	// 0x000c | 02 00 00 00 ....................................... .Blocks[0].Block.Head.Version
	// 0x0010 | 64 00 00 00 00 00 00 00 ........................... .Blocks[0].Block.Head.Time
	// 0x0018 | 01 00 00 00 00 00 00 00 ........................... .Blocks[0].Block.Head.BkSeq
	// 0x0020 | 0a 00 00 00 00 00 00 00 ........................... .Blocks[0].Block.Head.Fee
	// 0x0028 | 40 af f2 e9 d2 d8 92 2e 47 af d4 64 8e 69 67 49
	// 0x0038 | 71 58 78 5f bd 1d a8 70 e7 11 02 66 bf 94 48 80 ... .Blocks[0].Block.Head.PrevHash
	// 0x0048 | 6c 7b 5a 15 4b 6c 97 ea 8d 89 c1 51 79 a0 d0 03
	// 0x0058 | 1f 18 9e 3b 93 2a 9a 8e ac c7 70 57 36 09 64 85 ... .Blocks[0].Block.Head.BodyHash
	// 0x0068 | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x0078 | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 ... .Blocks[0].Block.Head.UxHash
	// 0x0088 | 02 00 00 00 ....................................... .Blocks[0].Block.Body.Transactions length
	// 0x008c | 02 00 00 00 ....................................... .Blocks[0].Block.Body.Transactions[0].Length
	// 0x0090 | 02 ................................................ .Blocks[0].Block.Body.Transactions[0].Type
	// 0x0091 | f4 45 7d e9 f5 a5 94 2e 07 6a 7f 2b 28 e1 84 2a
	// 0x00a1 | b6 1f 1b fc 39 e4 ca 55 75 36 60 0f d6 42 09 f6 ... .Blocks[0].Block.Body.Transactions[0].InnerHash
	// 0x00b1 | 02 00 00 00 ....................................... .Blocks[0].Block.Body.Transactions[0].Sigs length
	// 0x00b5 | 03 21 3f dd 6d df 86 0e 40 53 e1 a9 7e 42 76 d6
	// 0x00c5 | 34 54 f5 19 5a 83 21 35 70 04 d5 2c db bf d3 88
	// 0x00d5 | 6f c7 ad 3f 3f 63 b6 5d 4a 87 9c e3 08 6d ae b3
	// 0x00e5 | e5 4a 93 d3 c2 f9 6a 50 61 f9 bc 49 36 83 ca 8e
	// 0x00f5 | 01 ................................................ .Blocks[0].Block.Body.Transactions[0].Sigs[0]
	// 0x00f6 | ae ab ca 1d ec b2 ab 5b a8 af 93 c6 f0 33 de 6c
	// 0x0106 | bf 1d 50 31 4d f2 75 f9 40 77 8b c7 20 43 3c bc
	// 0x0116 | 21 94 aa cf a2 43 cc 2a 21 f8 5f 2f ff 71 d3 16
	// 0x0126 | 6d 18 75 e1 98 1a 0d a5 a2 3d 28 96 81 fc 1f a3
	// 0x0136 | 00 ................................................ .Blocks[0].Block.Body.Transactions[0].Sigs[1]
	// 0x0137 | 02 00 00 00 ....................................... .Blocks[0].Block.Body.Transactions[0].In length
	// 0x013b | 40 af f2 e9 d2 d8 92 2e 47 af d4 64 8e 69 67 49
	// 0x014b | 71 58 78 5f bd 1d a8 70 e7 11 02 66 bf 94 48 80 ... .Blocks[0].Block.Body.Transactions[0].In[0]
	// 0x015b | 7b b4 62 c3 bd 37 1d d8 1c 06 ad 1d 2b 63 59 71
	// 0x016b | cb 56 eb 22 23 3d fc 9f eb e8 3e 44 c8 40 b8 d7 ... .Blocks[0].Block.Body.Transactions[0].In[1]
	// 0x017b | 02 00 00 00 ....................................... .Blocks[0].Block.Body.Transactions[0].Out length
	// 0x017f | 00 07 6d ca 32 de 03 4e 48 67 fa 7a 2a a9 ee fe
	// 0x018f | 91 f2 0b a0 74 .................................... .Blocks[0].Block.Body.Transactions[0].Out[0].Address
	// 0x0194 | 02 00 00 00 00 00 00 00 ........................... .Blocks[0].Block.Body.Transactions[0].Out[0].Coins
	// 0x019c | 04 00 00 00 00 00 00 00 ........................... .Blocks[0].Block.Body.Transactions[0].Out[0].Hours
	// 0x01a4 | 00 e9 cb 47 35 e3 95 cf 36 b0 d1 a6 f2 21 bb 23
	// 0x01b4 | b3 f7 bf b1 f9 .................................... .Blocks[0].Block.Body.Transactions[0].Out[1].Address
	// 0x01b9 | 01 00 00 00 00 00 00 00 ........................... .Blocks[0].Block.Body.Transactions[0].Out[1].Coins
	// 0x01c1 | 02 00 00 00 00 00 00 00 ........................... .Blocks[0].Block.Body.Transactions[0].Out[1].Hours
	// 0x01c9 | 02 00 00 00 ....................................... .Blocks[0].Block.Body.Transactions[1].Length
	// 0x01cd | 01 ................................................ .Blocks[0].Block.Body.Transactions[1].Type
	// 0x01ce | c4 22 09 82 f9 88 b6 25 d0 af c1 2c 7f fd 06 a7
	// 0x01de | fe 89 bb e6 60 2c 1f 20 d9 08 91 3f e9 38 10 47 ... .Blocks[0].Block.Body.Transactions[1].InnerHash
	// 0x01ee | 02 00 00 00 ....................................... .Blocks[0].Block.Body.Transactions[1].Sigs length
	// 0x01f2 | ae ab ca 1d ec b2 ab 5b a8 af 93 c6 f0 33 de 6c
	// 0x0202 | bf 1d 50 31 4d f2 75 f9 40 77 8b c7 20 43 3c bc
	// 0x0212 | 21 94 aa cf a2 43 cc 2a 21 f8 5f 2f ff 71 d3 16
	// 0x0222 | 6d 18 75 e1 98 1a 0d a5 a2 3d 28 96 81 fc 1f a3
	// 0x0232 | 00 ................................................ .Blocks[0].Block.Body.Transactions[1].Sigs[0]
	// 0x0233 | 67 65 27 8a fc 9f 3e 0f fb 95 cb b3 f0 18 72 e9
	// 0x0243 | 2e d1 d5 1e 7a 83 d1 6d 49 9e 95 97 e2 4f a6 f3
	// 0x0253 | 08 d8 85 c8 31 c4 6b 69 9a d6 7b 2f dd 2f 76 2f
	// 0x0263 | d6 5f 4f bc f6 6f 98 1d 76 a9 ad fe 42 0d 16 14
	// 0x0273 | 01 ................................................ .Blocks[0].Block.Body.Transactions[1].Sigs[1]
	// 0x0274 | 01 00 00 00 ....................................... .Blocks[0].Block.Body.Transactions[1].In length
	// 0x0278 | f4 45 7d e9 f5 a5 94 2e 07 6a 7f 2b 28 e1 84 2a
	// 0x0288 | b6 1f 1b fc 39 e4 ca 55 75 36 60 0f d6 42 09 f6 ... .Blocks[0].Block.Body.Transactions[1].In[0]
	// 0x0298 | 02 00 00 00 ....................................... .Blocks[0].Block.Body.Transactions[1].Out length
	// 0x029c | 00 83 f1 96 59 16 14 99 2f a6 03 13 38 6f 72 88
	// 0x02ac | ac 40 14 c8 bc .................................... .Blocks[0].Block.Body.Transactions[1].Out[0].Address
	// 0x02b1 | 05 00 00 00 00 00 00 00 ........................... .Blocks[0].Block.Body.Transactions[1].Out[0].Coins
	// 0x02b9 | 05 00 00 00 00 00 00 00 ........................... .Blocks[0].Block.Body.Transactions[1].Out[0].Hours
	// 0x02c1 | 00 7e f9 b1 b9 40 6f 8d b3 99 b2 5f d0 e9 f4 f0
	// 0x02d1 | 88 7b 08 4b 43 .................................... .Blocks[0].Block.Body.Transactions[1].Out[1].Address
	// 0x02d6 | 03 00 00 00 00 00 00 00 ........................... .Blocks[0].Block.Body.Transactions[1].Out[1].Coins
	// 0x02de | 03 00 00 00 00 00 00 00 ........................... .Blocks[0].Block.Body.Transactions[1].Out[1].Hours
	// 0x02e6 | 5e db 9b bd 4a 78 20 f4 53 91 f0 8e 75 72 c9 1f
	// 0x02f6 | b6 1b e7 10 6d db 53 46 ac 00 59 6c 89 92 cf a2
	// 0x0306 | 3a ac 06 46 9a ec 55 33 71 02 41 8c eb 11 52 9b
	// 0x0316 | 3b 30 57 84 02 e5 8c 2f f6 6d cd 8b bf 35 21 e9
	// 0x0326 | 01 ................................................ .Blocks[0].Sig
	// 0x0327 |
}

func ExampleAnnounceBlocksMessage() {
	defer gnet.EraseMessages()
	setupMsgEncoding()
	var message = NewAnnounceBlocksMessage(123456)
	fmt.Println("AnnounceBlocksMessage:")
	var dmai = NewDeepMessagesAnnotationsIterator(message, 3)
	w := bufio.NewWriter(os.Stdout)
	NewFromIterator(gnet.EncodeMessage(message), &dmai, w) // nolint: errcheck
	// Output:
	// AnnounceBlocksMessage:
	// 0x0000 | 0c 00 00 00 ....................................... Length
	// 0x0004 | 41 4e 4e 42 ....................................... Prefix
	// 0x0008 | 40 e2 01 00 00 00 00 00 ........................... .MaxBkSeq
	// 0x0010 |
}

func ExampleGetTxnsMessage() {
	defer gnet.EraseMessages()
	setupMsgEncoding()
	var shas = make([]cipher.SHA256, 0)

	shas = append(shas, hashes[1], hashes[2])
	var message = NewGetTxnsMessage(shas)
	fmt.Println("GetTxnsMessage:")
	var dmai = NewDeepMessagesAnnotationsIterator(message, 3)
	w := bufio.NewWriter(os.Stdout)
	NewFromIterator(gnet.EncodeMessage(message), &dmai, w) // nolint: errcheck
	// Output:
	// GetTxnsMessage:
	// 0x0000 | 48 00 00 00 ....................................... Length
	// 0x0004 | 47 45 54 54 ....................................... Prefix
	// 0x0008 | 02 00 00 00 ....................................... .Transactions length
	// 0x000c | 7b b4 62 c3 bd 37 1d d8 1c 06 ad 1d 2b 63 59 71
	// 0x001c | cb 56 eb 22 23 3d fc 9f eb e8 3e 44 c8 40 b8 d7 ... .Transactions[0]
	// 0x002c | e7 5a c8 01 c1 3f 3d a9 c7 a1 24 ca 31 3b e2 a3
	// 0x003c | 73 f6 4a d9 7c 58 a1 b6 fe bc 0e 0c a5 c5 c8 73 ... .Transactions[1]
	// 0x004c |
}

func ExampleGiveTxnsMessage() {
	defer gnet.EraseMessages()
	setupMsgEncoding()
	var transactions coin.Transactions = make([]coin.Transaction, 0)
	var transactionOutputs0 = make([]coin.TransactionOutput, 0)
	var transactionOutputs1 = make([]coin.TransactionOutput, 0)
	var txOutput0 = coin.TransactionOutput{
		Address: addresses[0],
		Coins:   12,
		Hours:   34,
	}
	var txOutput1 = coin.TransactionOutput{
		Address: addresses[1],
		Coins:   56,
		Hours:   78,
	}
	var txOutput2 = coin.TransactionOutput{
		Address: addresses[3],
		Coins:   9,
		Hours:   12,
	}
	var txOutput3 = coin.TransactionOutput{
		Address: addresses[2],
		Coins:   34,
		Hours:   56,
	}
	transactionOutputs0 = append(transactionOutputs0, txOutput0, txOutput1)
	transactionOutputs1 = append(transactionOutputs1, txOutput2, txOutput3)

	var transaction0 = coin.Transaction{
		Type:      123,
		In:        []cipher.SHA256{hashes[3], hashes[4]},
		InnerHash: hashes[5],
		Length:    5000,
		Out:       transactionOutputs0,
		Sigs:      []cipher.Sig{sigs[0], sigs[1]},
	}
	var transaction1 = coin.Transaction{
		Type:      123,
		In:        []cipher.SHA256{hashes[5], hashes[6]},
		InnerHash: hashes[6],
		Length:    5000,
		Out:       transactionOutputs1,
		Sigs:      []cipher.Sig{sigs[2], sigs[3]},
	}
	transactions = append(transactions, transaction0, transaction1)
	var message = NewGiveTxnsMessage(transactions)
	fmt.Println("GiveTxnsMessage:")
	var dmai = NewDeepMessagesAnnotationsIterator(message, 5)
	w := bufio.NewWriter(os.Stdout)
	NewFromIterator(gnet.EncodeMessage(message), &dmai, w) // nolint: errcheck
	// Output:
	// GiveTxnsMessage:
	// 0x0000 | 82 02 00 00 ....................................... Length
	// 0x0004 | 47 49 56 54 ....................................... Prefix
	// 0x0008 | 02 00 00 00 ....................................... .Transactions length
	// 0x000c | 88 13 00 00 ....................................... .Transactions[0].Length
	// 0x0010 | 7b ................................................ .Transactions[0].Type
	// 0x0011 | 38 62 f1 93 a1 56 4e 5e 26 0f 82 7d a8 e1 69 ca
	// 0x0021 | d8 11 d8 1d 6a 7c 4f fd 66 1c 00 b1 99 94 17 81 ... .Transactions[0].InnerHash
	// 0x0031 | 02 00 00 00 ....................................... .Transactions[0].Sigs length
	// 0x0035 | 03 21 3f dd 6d df 86 0e 40 53 e1 a9 7e 42 76 d6
	// 0x0045 | 34 54 f5 19 5a 83 21 35 70 04 d5 2c db bf d3 88
	// 0x0055 | 6f c7 ad 3f 3f 63 b6 5d 4a 87 9c e3 08 6d ae b3
	// 0x0065 | e5 4a 93 d3 c2 f9 6a 50 61 f9 bc 49 36 83 ca 8e
	// 0x0075 | 01 ................................................ .Transactions[0].Sigs[0]
	// 0x0076 | ae ab ca 1d ec b2 ab 5b a8 af 93 c6 f0 33 de 6c
	// 0x0086 | bf 1d 50 31 4d f2 75 f9 40 77 8b c7 20 43 3c bc
	// 0x0096 | 21 94 aa cf a2 43 cc 2a 21 f8 5f 2f ff 71 d3 16
	// 0x00a6 | 6d 18 75 e1 98 1a 0d a5 a2 3d 28 96 81 fc 1f a3
	// 0x00b6 | 00 ................................................ .Transactions[0].Sigs[1]
	// 0x00b7 | 02 00 00 00 ....................................... .Transactions[0].In length
	// 0x00bb | f4 45 7d e9 f5 a5 94 2e 07 6a 7f 2b 28 e1 84 2a
	// 0x00cb | b6 1f 1b fc 39 e4 ca 55 75 36 60 0f d6 42 09 f6 ... .Transactions[0].In[0]
	// 0x00db | c4 22 09 82 f9 88 b6 25 d0 af c1 2c 7f fd 06 a7
	// 0x00eb | fe 89 bb e6 60 2c 1f 20 d9 08 91 3f e9 38 10 47 ... .Transactions[0].In[1]
	// 0x00fb | 02 00 00 00 ....................................... .Transactions[0].Out length
	// 0x00ff | 00 ................................................ .Transactions[0].Out[0].Address.Version
	// 0x0100 | 07 6d ca 32 de 03 4e 48 67 fa 7a 2a a9 ee fe 91
	// 0x0110 | f2 0b a0 74 ....................................... .Transactions[0].Out[0].Address.Key
	// 0x0114 | 0c 00 00 00 00 00 00 00 ........................... .Transactions[0].Out[0].Coins
	// 0x011c | 22 00 00 00 00 00 00 00 ........................... .Transactions[0].Out[0].Hours
	// 0x0124 | 00 ................................................ .Transactions[0].Out[1].Address.Version
	// 0x0125 | e9 cb 47 35 e3 95 cf 36 b0 d1 a6 f2 21 bb 23 b3
	// 0x0135 | f7 bf b1 f9 ....................................... .Transactions[0].Out[1].Address.Key
	// 0x0139 | 38 00 00 00 00 00 00 00 ........................... .Transactions[0].Out[1].Coins
	// 0x0141 | 4e 00 00 00 00 00 00 00 ........................... .Transactions[0].Out[1].Hours
	// 0x0149 | 88 13 00 00 ....................................... .Transactions[1].Length
	// 0x014d | 7b ................................................ .Transactions[1].Type
	// 0x014e | 05 64 0e 44 80 73 9e 87 97 57 b0 a2 d1 bd 59 de
	// 0x015e | a7 df cc fe f3 df 75 a1 83 0a 50 20 01 10 67 21 ... .Transactions[1].InnerHash
	// 0x016e | 02 00 00 00 ....................................... .Transactions[1].Sigs length
	// 0x0172 | 67 65 27 8a fc 9f 3e 0f fb 95 cb b3 f0 18 72 e9
	// 0x0182 | 2e d1 d5 1e 7a 83 d1 6d 49 9e 95 97 e2 4f a6 f3
	// 0x0192 | 08 d8 85 c8 31 c4 6b 69 9a d6 7b 2f dd 2f 76 2f
	// 0x01a2 | d6 5f 4f bc f6 6f 98 1d 76 a9 ad fe 42 0d 16 14
	// 0x01b2 | 01 ................................................ .Transactions[1].Sigs[0]
	// 0x01b3 | 5e db 9b bd 4a 78 20 f4 53 91 f0 8e 75 72 c9 1f
	// 0x01c3 | b6 1b e7 10 6d db 53 46 ac 00 59 6c 89 92 cf a2
	// 0x01d3 | 3a ac 06 46 9a ec 55 33 71 02 41 8c eb 11 52 9b
	// 0x01e3 | 3b 30 57 84 02 e5 8c 2f f6 6d cd 8b bf 35 21 e9
	// 0x01f3 | 01 ................................................ .Transactions[1].Sigs[1]
	// 0x01f4 | 02 00 00 00 ....................................... .Transactions[1].In length
	// 0x01f8 | 38 62 f1 93 a1 56 4e 5e 26 0f 82 7d a8 e1 69 ca
	// 0x0208 | d8 11 d8 1d 6a 7c 4f fd 66 1c 00 b1 99 94 17 81 ... .Transactions[1].In[0]
	// 0x0218 | 05 64 0e 44 80 73 9e 87 97 57 b0 a2 d1 bd 59 de
	// 0x0228 | a7 df cc fe f3 df 75 a1 83 0a 50 20 01 10 67 21 ... .Transactions[1].In[1]
	// 0x0238 | 02 00 00 00 ....................................... .Transactions[1].Out length
	// 0x023c | 00 ................................................ .Transactions[1].Out[0].Address.Version
	// 0x023d | 7e f9 b1 b9 40 6f 8d b3 99 b2 5f d0 e9 f4 f0 88
	// 0x024d | 7b 08 4b 43 ....................................... .Transactions[1].Out[0].Address.Key
	// 0x0251 | 09 00 00 00 00 00 00 00 ........................... .Transactions[1].Out[0].Coins
	// 0x0259 | 0c 00 00 00 00 00 00 00 ........................... .Transactions[1].Out[0].Hours
	// 0x0261 | 00 ................................................ .Transactions[1].Out[1].Address.Version
	// 0x0262 | 83 f1 96 59 16 14 99 2f a6 03 13 38 6f 72 88 ac
	// 0x0272 | 40 14 c8 bc ....................................... .Transactions[1].Out[1].Address.Key
	// 0x0276 | 22 00 00 00 00 00 00 00 ........................... .Transactions[1].Out[1].Coins
	// 0x027e | 38 00 00 00 00 00 00 00 ........................... .Transactions[1].Out[1].Hours
	// 0x0286 |
}

func TestIntroductionMessage(t *testing.T) {
	defer gnet.EraseMessages()
	setupMsgEncoding()

	pubkey, _ := cipher.GenerateKeyPair()
	pubkey2, _ := cipher.GenerateKeyPair()

	type mirrorPortResult struct {
		port  uint16
		exist bool
	}

	type daemonMockValue struct {
		version                    uint32
		mirror                     uint32
		isDefaultConnection        bool
		isMaxConnectionsReached    bool
		isMaxConnectionsReachedErr error
		setHasIncomingPortErr      error
		getMirrorPortResult        mirrorPortResult
		recordMessageEventErr      error
		pubkey                     cipher.PubKey
		disconnectReason           gnet.DisconnectReason
		disconnectErr              error
		addPeerArg                 string
		addPeerErr                 error
	}

	tt := []struct {
		name      string
		addr      string
		mockValue daemonMockValue
		intro     *IntroductionMessage
		err       error
	}{
		{
			name: "INTR message without extra bytes",
			addr: "121.121.121.121:6000",
			mockValue: daemonMockValue{
				mirror:  10000,
				version: 1,
				getMirrorPortResult: mirrorPortResult{
					exist: false,
				},
			},
			intro: &IntroductionMessage{
				Mirror:  10001,
				Port:    6000,
				Version: 1,
				valid:   true,
			},
			err: nil,
		},
		{
			name: "INTR message with pubkey",
			addr: "121.121.121.121:6000",
			mockValue: daemonMockValue{
				mirror:  10000,
				version: 1,
				getMirrorPortResult: mirrorPortResult{
					exist: false,
				},
				pubkey: pubkey,
			},
			intro: &IntroductionMessage{
				Mirror:  10001,
				Port:    6000,
				Version: 1,
				valid:   true,
				Extra:   pubkey[:],
			},
			err: nil,
		},
		{
			name: "INTR message with pubkey",
			addr: "121.121.121.121:6000",
			mockValue: daemonMockValue{
				mirror:  10000,
				version: 1,
				getMirrorPortResult: mirrorPortResult{
					exist: false,
				},
				pubkey: pubkey,
			},
			intro: &IntroductionMessage{
				Mirror:  10001,
				Port:    6000,
				Version: 1,
				valid:   true,
				Extra:   pubkey[:],
			},
			err: nil,
		},
		{
			name: "INTR message with pubkey and additional data",
			addr: "121.121.121.121:6000",
			mockValue: daemonMockValue{
				mirror:  10000,
				version: 1,
				getMirrorPortResult: mirrorPortResult{
					exist: false,
				},
				pubkey: pubkey,
			},
			intro: &IntroductionMessage{
				Mirror:  10001,
				Port:    6000,
				Version: 1,
				valid:   true,
				Extra:   append(pubkey[:], []byte("additional data")...),
			},
			err: nil,
		},
		{
			name: "INTR message with different pubkey",
			addr: "121.121.121.121:6000",
			mockValue: daemonMockValue{
				mirror:  10000,
				version: 1,
				getMirrorPortResult: mirrorPortResult{
					exist: false,
				},
				pubkey:           pubkey,
				disconnectReason: ErrDisconnectBlockchainPubkeyNotMatched,
			},
			intro: &IntroductionMessage{
				Mirror:  10001,
				Port:    6000,
				Version: 1,
				valid:   true,
				Extra:   pubkey2[:],
			},
			err: ErrDisconnectBlockchainPubkeyNotMatched,
		},
		{
			name: "INTR message with invalid pubkey",
			addr: "121.121.121.121:6000",
			mockValue: daemonMockValue{
				mirror:  10000,
				version: 1,
				getMirrorPortResult: mirrorPortResult{
					exist: false,
				},
				pubkey:           pubkey,
				disconnectReason: ErrDisconnectInvalidExtraData,
			},
			intro: &IntroductionMessage{
				Mirror:  10001,
				Port:    6000,
				Version: 1,
				valid:   true,
				Extra:   []byte("invalid extra data"),
			},
			err: ErrDisconnectInvalidExtraData,
		},
		{
			name: "Disconnect self connection",
			mockValue: daemonMockValue{
				mirror:           10000,
				disconnectReason: ErrDisconnectSelf,
			},
			intro: &IntroductionMessage{
				Mirror: 10000,
			},
			err: ErrDisconnectSelf,
		},
		{
			name: "Invalid version",
			mockValue: daemonMockValue{
				mirror:           10000,
				version:          1,
				disconnectReason: ErrDisconnectInvalidVersion,
			},
			intro: &IntroductionMessage{
				Mirror:  10001,
				Version: 0,
			},
			err: ErrDisconnectInvalidVersion,
		},
		{
			name: "Invalid address",
			addr: "121.121.121.121",
			mockValue: daemonMockValue{
				mirror:           10000,
				version:          1,
				disconnectReason: ErrDisconnectOtherError,
				pubkey:           pubkey,
			},
			intro: &IntroductionMessage{
				Mirror:  10001,
				Version: 1,
				Port:    6000,
			},
			err: ErrDisconnectOtherError,
		},
		{
			name: "incomming connection",
			addr: "121.121.121.121:12345",
			mockValue: daemonMockValue{
				mirror:                  10000,
				version:                 1,
				isDefaultConnection:     true,
				isMaxConnectionsReached: true,
				getMirrorPortResult: mirrorPortResult{
					exist: false,
				},
				pubkey:     pubkey,
				addPeerArg: "121.121.121.121:6000",
				addPeerErr: nil,
			},
			intro: &IntroductionMessage{
				Mirror:  10001,
				Version: 1,
				Port:    6000,
				valid:   true,
			},
		},
		{
			name: "Connect twice",
			addr: "121.121.121.121:6000",
			mockValue: daemonMockValue{
				mirror:              10000,
				version:             1,
				isDefaultConnection: true,
				getMirrorPortResult: mirrorPortResult{
					exist: true,
				},
				pubkey:           pubkey,
				addPeerArg:       "121.121.121.121:6000",
				addPeerErr:       nil,
				disconnectReason: ErrDisconnectConnectedTwice,
			},
			intro: &IntroductionMessage{
				Mirror:  10001,
				Version: 1,
				Port:    6000,
			},
			err: ErrDisconnectConnectedTwice,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			mc := &gnet.MessageContext{Addr: tc.addr}
			tc.intro.c = mc

			d := &MockDaemoner{}
			d.On("DaemonConfig").Return(DaemonConfig{Version: int32(tc.mockValue.version)})
			d.On("Mirror").Return(tc.mockValue.mirror)
			d.On("IsDefaultConnection", tc.addr).Return(tc.mockValue.isDefaultConnection)
			d.On("SetHasIncomingPort", tc.addr).Return(tc.mockValue.setHasIncomingPortErr)
			d.On("GetMirrorPort", tc.addr, tc.intro.Mirror).Return(tc.mockValue.getMirrorPortResult.port, tc.mockValue.getMirrorPortResult.exist)
			d.On("RecordMessageEvent", tc.intro, mc).Return(tc.mockValue.recordMessageEventErr)
			d.On("ResetRetryTimes", tc.addr)
			d.On("BlockchainPubkey").Return(tc.mockValue.pubkey)
			d.On("Disconnect", tc.addr, tc.mockValue.disconnectReason).Return(tc.mockValue.disconnectErr)
			d.On("IncreaseRetryTimes", tc.addr)
			d.On("RemoveFromExpectingIntroductions", tc.addr)
			d.On("IsMaxDefaultConnectionsReached").Return(tc.mockValue.isMaxConnectionsReached, tc.mockValue.isMaxConnectionsReachedErr)
			d.On("AddPeer", tc.mockValue.addPeerArg).Return(tc.mockValue.addPeerErr)

			err := tc.intro.Handle(mc, d)
			require.Equal(t, tc.err, err)
		})
	}
}

func TestMessageEncodeDecode(t *testing.T) {
	update := false

	cases := []struct {
		goldenFile string
		obj        interface{}
		msg        interface{}
	}{
		{
			goldenFile: "intro-msg.golden",
			obj:        &IntroductionMessage{},
			msg: &IntroductionMessage{
				Mirror:  99998888,
				Port:    8888,
				Version: 12341234,
			},
		},
		{
			goldenFile: "intro-msg-extra.golden",
			obj:        &IntroductionMessage{},
			msg: &IntroductionMessage{
				Mirror:  99998888,
				Port:    8888,
				Version: 12341234,
				Extra:   []byte("abcdef"),
			},
		},
		{
			goldenFile: "get-peers-msg.golden",
			obj:        &GetPeersMessage{},
			msg:        &GetPeersMessage{},
		},
		{
			goldenFile: "give-peers-msg.golden",
			obj:        &GivePeersMessage{},
			msg: &GivePeersMessage{
				Peers: []IPAddr{
					{
						IP:   12345678,
						Port: 1234,
					},
					{
						IP:   87654321,
						Port: 4321,
					},
				},
			},
		},
		{
			goldenFile: "ping-msg.golden",
			obj:        &PingMessage{},
			msg:        &PingMessage{},
		},
		{
			goldenFile: "pong-msg.golden",
			obj:        &PongMessage{},
			msg:        &PongMessage{},
		},
		{
			goldenFile: "get-blocks-msg.golden",
			obj:        &GetBlocksMessage{},
			msg: &GetBlocksMessage{
				LastBlock:       999988887777,
				RequestedBlocks: 888899997777,
			},
		},
		{
			goldenFile: "give-blocks-msg.golden",
			obj:        &GiveBlocksMessage{},
			msg: &GiveBlocksMessage{
				Blocks: []coin.SignedBlock{
					{
						Sig: cipher.MustSigFromHex("8cf145e9ef4a4a5254bc57798a7a61dfed238768f94edc5635175c6b91bccd8ec1555da603c5e31b018e135b82b1525be8a92973c468a74b5b40b8da189cb465eb"),
						Block: coin.Block{
							Head: coin.BlockHeader{
								Version:  1,
								Time:     1538036613,
								BkSeq:    9999999999,
								Fee:      1234123412341234,
								PrevHash: cipher.MustSHA256FromHex("59cb7d0e2ce8a03d1054afcc28a22fe864a8813460d241db38c59d10e7c29132"),
								BodyHash: cipher.MustSHA256FromHex("6d421469409591f0c3112884c8cf10f8bca5d8ab87c9c30dea2ea73b6751bbf9"),
								UxHash:   cipher.MustSHA256FromHex("6ea6a972cf06d25908b29953aeddb68c3b6f3a9903e8f964dc89b0abc0645dea"),
							},
							Body: coin.BlockBody{
								Transactions: coin.Transactions{
									{
										Length:    43214321,
										Type:      1,
										InnerHash: cipher.MustSHA256FromHex("cbedf8ef0bda91afc6a180eea0dddf8e3a986b6b6f87f70e8bffc63c6fbaa4e6"),
										Sigs: []cipher.Sig{
											cipher.MustSigFromHex("1cfd7a4db3a52a85d2a86708695112b6520acc8dc83c86e8da67915199fdf04964c168543598ab07c2b99c292899890891950364c2bf66f1aaa6d6a66a5c9a73ff"),
											cipher.MustSigFromHex("442167c6b3d13957bc32f83182c7f4fda0bb6bde893a41a6a04cdd8eecee0048d03a57eb2af04ea6050e1f418769c94c7f12fad9287dc650e6b307fdfce6b42a59"),
										},
										In: []cipher.SHA256{
											cipher.MustSHA256FromHex("536f0a1a915fadfa3a2720a0615641827ff67394d2b2149d6db63b8c619e14af"),
											cipher.MustSHA256FromHex("64ba5f01f90f97f84999f13aeaa75fed8d5b3e4a3a4a093dedf4795969e8bd27"),
										},
										Out: []coin.TransactionOutput{
											{
												Address: cipher.MustDecodeBase58Address("23FF4fshzD8tZk2d88P22WATfzUpNQF1x85"),
												Coins:   987987987,
												Hours:   789789789,
											},
											{
												Address: cipher.MustDecodeBase58Address("29V2iRpZAqHiFZHHRqaZLArZZuTcZM5owqT"),
												Coins:   123123,
												Hours:   321321,
											},
										},
									},
									{
										Length:    98769876,
										Type:      0,
										InnerHash: cipher.MustSHA256FromHex("46856af925fde9a1652d39eea479dd92589a741451a0228402e399fae02f8f3d"),
										Sigs: []cipher.Sig{
											cipher.MustSigFromHex("92e289792200518df9a82cf9dddd1f334bf0d47fb0ed4ff70c25403f39577af5ab24ef2d02a11cf6b76e6bd017457ad60d6ca85c0567c21f5c62599c93ee98e18c"),
											cipher.MustSigFromHex("e995da86ed87640ecb44e624074ba606b781aa0cbeb24e8c27ff30becf7181175479c0d74d93fe1e8692bba628b5cf532ca80fed4135148d84e6ecc2a762a10b19"),
										},
										In: []cipher.SHA256{
											cipher.MustSHA256FromHex("69b14a7ee184f24b95659d6887101ef7c921fa7977d95c73fbc0c4d0d22671bc"),
											cipher.MustSHA256FromHex("3a050b4ec33ec9ad2c789f24655ab1c8f7691d3a1c3d0e05cc14b022b4c360ea"),
										},
										Out: []coin.TransactionOutput{
											{
												Address: cipher.MustDecodeBase58Address("XvvjeyGcTBVXDXmfJoTUseFiqHvm12C6oQ"),
												Coins:   15,
												Hours:   1237882,
											},
											{
												Address: cipher.MustDecodeBase58Address("fQXVLq9fbCC9XVxDKLAGDLXmQYjPy59j22"),
												Coins:   2102123,
												Hours:   1003,
											},
										},
									},
								},
							},
						},
					},
					{
						Sig: cipher.MustSigFromHex("8015c8776de577d89c29d1cbd1d558ba4855dec94ba58f6c67d55ece5c85708b9906bd0b72b451e27008f3938fcec42c1a28ddac336ae8206d8e6443b95dde966c"),
						Block: coin.Block{
							Head: coin.BlockHeader{
								Version:  0,
								Time:     1427248825,
								BkSeq:    100,
								Fee:      120939323123,
								PrevHash: cipher.MustSHA256FromHex("04d40b5d27c539ab9d98934628604baef7dbfb1c35ddf9c0f96a67f6b061fa26"),
								BodyHash: cipher.MustSHA256FromHex("9a67fbb00216ae99f334d4efa2c9c42a25aac5d1a5bbb2058fe5705cfe0e30ea"),
								UxHash:   cipher.MustSHA256FromHex("58981d30da11be3c8e9dd8fdb7b51b48ba13dc0214cf211251308985bf089f76"),
							},
							Body: coin.BlockBody{
								Transactions: coin.Transactions{
									{
										Length:    128,
										Type:      99,
										InnerHash: cipher.MustSHA256FromHex("e943fd54a8071bb0ae92800c23c5a26443b5e5bf9b9321cefcdd9e80f518c37e"),
										Sigs: []cipher.Sig{
											cipher.MustSigFromHex("cff49d1d450db812d42748d4f7001e03a1dd2b98afcbb62eca1b3b1fa137e5095a0368250aabd3976008afe61471ecd31ed99185c3df49269d9aada4ca1dd2eecb"),
											cipher.MustSigFromHex("1313e5a80d6d9386fe2dffa13afba7277402f029d411e60f99b3806fee547d6157ca2d8d6407df3e858d6f3f58902f460412611282a0dec2468e41a2c5a39cc93e"),
										},
										In: []cipher.SHA256{
											cipher.MustSHA256FromHex("6a76c83b7b75075e2e34405e21d5e8d37adb69e4e6487a6179944ea7e04bc7db"),
											cipher.MustSHA256FromHex("a7555179a255e6a7dddb6121bd4c2259f75ebc321345be26b690f34094012f95"),
										},
										Out: []coin.TransactionOutput{
											{
												Address: cipher.MustDecodeBase58Address("2RmSTGbj5qaFT1WvKGz4SobaT4xSb9GvaCh"),
												Coins:   12301923233,
												Hours:   39932,
											},
											{
												Address: cipher.MustDecodeBase58Address("uA8XQnNzS4kit9DFzybyVSpWDEDy62MXur"),
												Coins:   9945924,
												Hours:   9030300895893902,
											},
										},
									},
									{
										Length:    1304,
										Type:      255,
										InnerHash: cipher.MustSHA256FromHex("d92057e9a4874aa876b7fd20074d78a4d890c2d3af483a10206f243308586763"),
										Sigs: []cipher.Sig{
											cipher.MustSigFromHex("394d53cc0bfeef11cc94bf39316d555549cf1a1afd14920be7d065e7940cc60752b8ade8c37991307a5681b06e0445c1c19ceb0e6611fd4593dcc65d18975c87be"),
											cipher.MustSigFromHex("50ad670bc672558c235653b6396135352bfbc8eec575de3cffce65d5a07076082f9694880eb6b1e708eb8fb39d21a96dd99615b5759fc917c3fdd4d9845489119b"),
										},
										In: []cipher.SHA256{
											cipher.MustSHA256FromHex("f37057407a6b5b103218abdfc5b5527f8abcc229256c912ec81ac6d72b68454e"),
											cipher.MustSHA256FromHex("9cd1fccddb5895ab77cd419802430e16a1e05f0f796d026fc69961c5c308b766"),
										},
										Out: []coin.TransactionOutput{
											{
												Address: cipher.MustDecodeBase58Address("MNf67cWXYmSizin4XUtGnFfQQzxkvNqCEH"),
												Coins:   1,
												Hours:   1,
											},
											{
												Address: cipher.MustDecodeBase58Address("HEkH8R1Uc58mAjZqGM15cqF4QMqG4mu4ry"),
												Coins:   1,
												Hours:   0,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			goldenFile: "announce-blocks-msg.golden",
			obj:        &AnnounceBlocksMessage{},
			msg: &AnnounceBlocksMessage{
				MaxBkSeq: 50000,
			},
		},
		{
			goldenFile: "announce-txns-msg.golden",
			obj:        &AnnounceTxnsMessage{},
			msg: &AnnounceTxnsMessage{
				Transactions: []cipher.SHA256{
					cipher.MustSHA256FromHex("23dc4b68c0fc790989bb82f04b9d5174baab6f0f6808ed35be9b93cb73c69108"),
					cipher.MustSHA256FromHex("2be4b0155c1ab9613007fe522e3b12bac4be79800a19bc8cd8ca343868caa583"),
				},
			},
		},
		{
			goldenFile: "get-txns-msg.golden",
			obj:        &GetTxnsMessage{},
			msg: &GetTxnsMessage{
				Transactions: []cipher.SHA256{
					cipher.MustSHA256FromHex("335b63b0f335c6aee5e7e1b3c62dd09bb6074e38b48e2469e294a019d5ae5aa1"),
					cipher.MustSHA256FromHex("619a367f4e5dee741348366899237ddc920335fc847ccafdf2d32ed57bb7b385"),
				},
			},
		},
		{
			goldenFile: "give-txns-msg.golden",
			obj:        &GiveTxnsMessage{},
			msg: &GiveTxnsMessage{
				Transactions: coin.Transactions{
					{
						Length:    256,
						Type:      0,
						InnerHash: cipher.MustSHA256FromHex("1773d8901df96bba4c6d65499e11e6ec73a9978c611d1463898ffbc2b49773fc"),
						Sigs: []cipher.Sig{
							cipher.MustSigFromHex("a711880ae54d1b6b9adade2ef1e743d6d539a78b0cecf1af08107e467956de80ef1d49fb5e896c9d0870ef8bf8a4d328ca0ecf7c1956866867ec56064e68f8a374"),
							cipher.MustSigFromHex("f9890ddd93f9479e364261ebc647326d2fd57e50b7728795adbf507c956f9eb44f77207b528700c4cef338290cdfc17f814dc3d94e3d711e92492ecc7b8abef808"),
						},
						In: []cipher.SHA256{
							cipher.MustSHA256FromHex("703f84ee0702b44fc89ce573a239d5fbf185bf5d4e7fc8f4930262bcda1e8fb0"),
							cipher.MustSHA256FromHex("c9e904862da01f2d7676c12c4342dde36d9a9a9d25be5351e2b57fae6f426bb9"),
						},
						Out: []coin.TransactionOutput{
							{
								Address: cipher.MustDecodeBase58Address("29VEn56iRr2TpVVpPoPxUJPfFWuhbLSBRdU"),
								Coins:   1111111111111111111,
								Hours:   9999999999999999999,
							},
							{
								Address: cipher.MustDecodeBase58Address("2bqs99tysFtfs8QPT81kpZWnzTT1rWd8xtQ"),
								Coins:   9922581002,
								Hours:   9932900022223334,
							},
						},
					},
					{
						Length:    13043,
						Type:      128,
						InnerHash: cipher.MustSHA256FromHex("a9da3e4acb1892a000c1b658a64d4e420d0c381862928ab820fb3f3a534a9674"),
						Sigs: []cipher.Sig{
							cipher.MustSigFromHex("7bbbdfd58c0533aed95f18d9413e0e0517892350eaf132eadf7a9a03d4a974ca0bc074abc001f86a34cf66c10f832dbcca20c2c67b5e8517f4ff0e1d0123fecb21"),
							cipher.MustSigFromHex("68732b78ac3a4e2fe146b8819c8b1c0b126a0188008c9c7c98fee965beba039778010ff7b0379dadeeadbbc42f9541ce4ad3c8cec12108d3aa58aca583bddd0df0"),
						},
						In: []cipher.SHA256{
							cipher.MustSHA256FromHex("766d6f6ed56599a91759c75466e3f09b9d6d5995b58dd5bbfba5af10b1a8cdea"),
							cipher.MustSHA256FromHex("2c7989f47524721bb2c7a7f967208c9b1c01829c9a55addf22d066e5c55ab3ac"),
						},
						Out: []coin.TransactionOutput{
							{
								Address: cipher.MustDecodeBase58Address("24iFsYHzVfYXo8cvWg1jhetpTMNvHH7j6AX"),
								Coins:   1123103123,
								Hours:   123000,
							},
							{
								Address: cipher.MustDecodeBase58Address("JV5xJ33po1Bj3dXZT3SYA3ZmnTibREFxxd"),
								Coins:   999999,
								Hours:   9043285343,
							},
						},
					},
				},
			},
		},
	}

	if update {
		for _, tc := range cases {
			t.Run(tc.goldenFile, func(t *testing.T) {
				fn := filepath.Join("testdata/", tc.goldenFile)

				f, err := os.Create(fn)
				require.NoError(t, err)
				defer f.Close()

				b := encoder.Serialize(tc.msg)
				_, err = f.Write(b)
				require.NoError(t, err)
			})
		}
	}

	for _, tc := range cases {
		t.Run(tc.goldenFile, func(t *testing.T) {
			fn := filepath.Join("testdata/", tc.goldenFile)

			f, err := os.Open(fn)
			require.NoError(t, err)
			defer f.Close()

			d, err := ioutil.ReadAll(f)
			require.NoError(t, err)

			err = encoder.DeserializeRaw(d, tc.obj)
			require.NoError(t, err)

			require.Equal(t, tc.msg, tc.obj)

			d2 := encoder.Serialize(tc.msg)
			require.Equal(t, d, d2)
		})
	}
}
