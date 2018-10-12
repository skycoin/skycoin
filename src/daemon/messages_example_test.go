package daemon

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/daemon/gnet"
	"github.com/skycoin/skycoin/src/daemon/pex"
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

func ExampleAnnounceTxnsMessage() {
	defer gnet.EraseMessages()
	setupMsgEncoding()

	var message = NewAnnounceTxnsMessage([]cipher.SHA256{hashes[7], hashes[8]})
	fmt.Println("AnnounceTxnsMessage:")
	var mai = NewDeepMessagesAnnotationsIterator(message, 3)
	w := bufio.NewWriter(os.Stdout)
	err := NewFromIterator(gnet.EncodeMessage(message), &mai, w)
	if err != nil {
		fmt.Println(err)
	}

	// Output:
	// AnnounceTxnsMessage:
	// 0x0000 | 48 00 00 00 ....................................... Length
	// 0x0004 | 41 4e 4e 54 ....................................... Prefix
	// 0x0008 | 02 00 00 00 ....................................... .Transactions length
	// 0x000c | 8a 5d bf bb 7e 64 66 49 5e 30 78 1c 15 40 b5 e3
	// 0x001c | 98 e0 84 4f 60 c9 1e c6 78 9d 4b bb 36 7e 33 a6 ... .Transactions[0]
	// 0x002c | 1c 1d 7d bf d7 ba 2b b1 aa 9b 56 ed ae 26 ea 56
	// 0x003c | 5c bf 72 f9 8c c6 a6 2c 72 97 23 cb c0 75 0d 3b ... .Transactions[1]
	// 0x004c |
}