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
	"github.com/skycoin/skycoin/src/params"
)

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
				if encoder.TagOmitempty(f.Tag.Get("enc")) {
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

/**************************************
 *
 * Test cases
 *
 *************************************/

var hashes = []cipher.SHA256{
	// buffer => 0x00, 0x01, ... 0xff
	GetSHAFromHex("40aff2e9d2d8922e47afd4648e6967497158785fbd1da870e7110266bf944880"),
	// buffer => 0x01, 0x01, ... 0xff
	GetSHAFromHex("7bb462c3bd371dd81c06ad1d2b635971cb56eb22233dfc9febe83e44c840b8d7"),
	// buffer => 0x02, 0x01, ... 0xff
	GetSHAFromHex("e75ac801c13f3da9c7a124ca313be2a373f64ad97c58a1b6febc0e0ca5c5c873"),
	// buffer => 0x03, 0x01, ... 0xff
	GetSHAFromHex("f4457de9f5a5942e076a7f2b28e1842ab61f1bfc39e4ca557536600fd64209f6"),
	// buffer => 0x04, 0x01, ... 0xff
	GetSHAFromHex("c4220982f988b625d0afc12c7ffd06a7fe89bbe6602c1f20d908913fe9381047"),
	// buffer => 0x05, 0x01, ... 0xff
	GetSHAFromHex("3862f193a1564e5e260f827da8e169cad811d81d6a7c4ffd661c00b199941781"),
	// buffer => 0x06, 0x01, ... 0xff
	GetSHAFromHex("05640e4480739e879757b0a2d1bd59dea7dfccfef3df75a1830a502001106721"),
	// buffer => 0x07, 0x01, ... 0xff
	GetSHAFromHex("8a5dbfbb7e6466495e30781c1540b5e398e0844f60c91ec6789d4bbb367e33a6"),
	// buffer => 0x08, 0x01, ... 0xff
	GetSHAFromHex("1c1d7dbfd7ba2bb1aa9b56edae26ea565cbf72f98cc6a62c729723cbc0750d3b"),
	// buffer => 0x09, 0x01, ... 0xff
	GetSHAFromHex("66dd3fc45be9b4fbb1fbed2be5de4e8a479f6638adfe4675b8544ae84eca3f75"),
}

func initSecKey(secKeyHex string) cipher.SecKey {
	sk, err := cipher.SecKeyFromHex(secKeyHex)
	if err != nil {
		sk = cipher.SecKey{}
	}
	return sk
}

func initAddress(sk cipher.SecKey) cipher.Address {
	addr, err := cipher.AddressFromSecKey(sk)
	if err != nil {
		addr = cipher.Address{}
	}
	return addr
}

var (
	// seed = 'w'
	secKey1 = initSecKey("2c08fb99b0ba9b8c32072ae52719f25c064735d03618354656fdb31a88150f7f")
	// seed = 'x'
	secKey2 = initSecKey("7a0f56ee1c49ef669235065f0b1d7a03c252b7f7bdcf39704c7cc1feff41cf53")
	// seed = 'y'
	secKey3 = initSecKey("def165004ebc06530105cc82496816ef754933d6668db5780c16c00c853a7991")
	// seed = 'z'
	secKey4 = initSecKey("061035ba600bcc4442389011d310ffcb19177626eaf582b344804d395805b5e1")
)

var addresses = []cipher.Address{
	initAddress(secKey1),
	initAddress(secKey2),
	initAddress(secKey3),
	initAddress(secKey4),
}

func GetSHAFromHex(hex string) cipher.SHA256 {
	var sha, _ = cipher.SHA256FromHex(hex)
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
	err := NewFromIterator(gnet.EncodeMessage(&message), &mai, w)
	if err != nil {
		fmt.Println(err)
	}
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

	pk := cipher.MustPubKeyFromHex("0328c576d3f420e7682058a981173a4b374c7cc5ff55bf394d3cf57059bbe6456a")

	var message = NewIntroductionMessage(1234, 5, 7890, pk, "skycoin:0.24.1", params.VerifyTxn{
		BurnFactor:          2,
		MaxTransactionSize:  32768,
		MaxDropletPrecision: 3,
	})
	fmt.Println("IntroductionMessage:")
	var mai = NewMessagesAnnotationsIterator(message)
	w := bufio.NewWriter(os.Stdout)
	err := NewFromIterator(gnet.EncodeMessage(message), &mai, w)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
	// IntroductionMessage:
	// 0x0000 | 4e 00 00 00 ....................................... Length
	// 0x0004 | 49 4e 54 52 ....................................... Prefix
	// 0x0008 | d2 04 00 00 ....................................... Mirror
	// 0x000c | d2 1e ............................................. ListenPort
	// 0x000e | 05 00 00 00 ....................................... ProtocolVersion
	// 0x0012 | 3c 00 00 00 ....................................... Extra length
	// 0x0016 | 03 ................................................ Extra[0]
	// 0x0017 | 28 ................................................ Extra[1]
	// 0x0018 | c5 ................................................ Extra[2]
	// 0x0019 | 76 ................................................ Extra[3]
	// 0x001a | d3 ................................................ Extra[4]
	// 0x001b | f4 ................................................ Extra[5]
	// 0x001c | 20 ................................................ Extra[6]
	// 0x001d | e7 ................................................ Extra[7]
	// 0x001e | 68 ................................................ Extra[8]
	// 0x001f | 20 ................................................ Extra[9]
	// 0x0020 | 58 ................................................ Extra[10]
	// 0x0021 | a9 ................................................ Extra[11]
	// 0x0022 | 81 ................................................ Extra[12]
	// 0x0023 | 17 ................................................ Extra[13]
	// 0x0024 | 3a ................................................ Extra[14]
	// 0x0025 | 4b ................................................ Extra[15]
	// 0x0026 | 37 ................................................ Extra[16]
	// 0x0027 | 4c ................................................ Extra[17]
	// 0x0028 | 7c ................................................ Extra[18]
	// 0x0029 | c5 ................................................ Extra[19]
	// 0x002a | ff ................................................ Extra[20]
	// 0x002b | 55 ................................................ Extra[21]
	// 0x002c | bf ................................................ Extra[22]
	// 0x002d | 39 ................................................ Extra[23]
	// 0x002e | 4d ................................................ Extra[24]
	// 0x002f | 3c ................................................ Extra[25]
	// 0x0030 | f5 ................................................ Extra[26]
	// 0x0031 | 70 ................................................ Extra[27]
	// 0x0032 | 59 ................................................ Extra[28]
	// 0x0033 | bb ................................................ Extra[29]
	// 0x0034 | e6 ................................................ Extra[30]
	// 0x0035 | 45 ................................................ Extra[31]
	// 0x0036 | 6a ................................................ Extra[32]
	// 0x0037 | 02 ................................................ Extra[33]
	// 0x0038 | 00 ................................................ Extra[34]
	// 0x0039 | 00 ................................................ Extra[35]
	// 0x003a | 00 ................................................ Extra[36]
	// 0x003b | 00 ................................................ Extra[37]
	// 0x003c | 80 ................................................ Extra[38]
	// 0x003d | 00 ................................................ Extra[39]
	// 0x003e | 00 ................................................ Extra[40]
	// 0x003f | 03 ................................................ Extra[41]
	// 0x0040 | 0e ................................................ Extra[42]
	// 0x0041 | 00 ................................................ Extra[43]
	// 0x0042 | 00 ................................................ Extra[44]
	// 0x0043 | 00 ................................................ Extra[45]
	// 0x0044 | 73 ................................................ Extra[46]
	// 0x0045 | 6b ................................................ Extra[47]
	// 0x0046 | 79 ................................................ Extra[48]
	// 0x0047 | 63 ................................................ Extra[49]
	// 0x0048 | 6f ................................................ Extra[50]
	// 0x0049 | 69 ................................................ Extra[51]
	// 0x004a | 6e ................................................ Extra[52]
	// 0x004b | 3a ................................................ Extra[53]
	// 0x004c | 30 ................................................ Extra[54]
	// 0x004d | 2e ................................................ Extra[55]
	// 0x004e | 32 ................................................ Extra[56]
	// 0x004f | 34 ................................................ Extra[57]
	// 0x0050 | 2e ................................................ Extra[58]
	// 0x0051 | 31 ................................................ Extra[59]
	// 0x0052 |
}

func ExampleGetPeersMessage() {
	defer gnet.EraseMessages()
	setupMsgEncoding()
	var message = NewGetPeersMessage()
	fmt.Println("GetPeersMessage:")
	var mai = NewMessagesAnnotationsIterator(message)
	w := bufio.NewWriter(os.Stdout)
	err := NewFromIterator(gnet.EncodeMessage(message), &mai, w)
	if err != nil {
		fmt.Println(err)
	}
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
	var mai = NewMessagesAnnotationsIterator(message)
	w := bufio.NewWriter(os.Stdout)
	err := NewFromIterator(gnet.EncodeMessage(message), &mai, w)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
	// GivePeersMessage:
	// 0x0000 | 1a 00 00 00 ....................................... Length
	// 0x0004 | 47 49 56 50 ....................................... Prefix
	// 0x0008 | 03 00 00 00 ....................................... Peers length
	// 0x000c | 5d 87 b2 76 70 17 ................................. Peers[0]
	// 0x0012 | 9c 21 58 2f 70 17 ................................. Peers[1]
	// 0x0018 | 94 67 29 79 70 17 ................................. Peers[2]
	// 0x001e |
}

func ExampleGetBlocksMessage() {
	defer gnet.EraseMessages()
	setupMsgEncoding()
	var message = NewGetBlocksMessage(1234, 5678)
	fmt.Println("GetBlocksMessage:")
	var mai = NewMessagesAnnotationsIterator(message)
	w := bufio.NewWriter(os.Stdout)
	err := NewFromIterator(gnet.EncodeMessage(message), &mai, w)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
	// GetBlocksMessage:
	// 0x0000 | 14 00 00 00 ....................................... Length
	// 0x0004 | 47 45 54 42 ....................................... Prefix
	// 0x0008 | d2 04 00 00 00 00 00 00 ........................... LastBlock
	// 0x0010 | 2e 16 00 00 00 00 00 00 ........................... RequestedBlocks
	// 0x0018 |
}

func ExampleGiveBlocksMessage() {
	defer gnet.EraseMessages()
	setupMsgEncoding()
	var blocks = make([]coin.SignedBlock, 2)
	var body1 = coin.BlockBody{
		Transactions: make([]coin.Transaction, 0),
	}
	var body2 = coin.BlockBody{
		Transactions: make([]coin.Transaction, 0),
	}
	var block1 = coin.Block{
		Body: body1,
		Head: coin.BlockHeader{
			Version:  0x02,
			Time:     100,
			BkSeq:    0,
			Fee:      10,
			PrevHash: hashes[0],
			BodyHash: body1.Hash(),
		}}
	var block2 = coin.Block{
		Body: body2,
		Head: coin.BlockHeader{
			Version:  0x02,
			Time:     100,
			BkSeq:    0,
			Fee:      250,
			PrevHash: hashes[1],
			BodyHash: body1.Hash(),
		}}
	var sig, _ = cipher.SigFromHex(sig1hex)
	var signedBlock = coin.SignedBlock{
		Sig:   sig,
		Block: block1,
	}
	blocks[0] = signedBlock
	sig, _ = cipher.SigFromHex(sig2hex) // nolint: errcheck
	signedBlock = coin.SignedBlock{
		Sig:   sig,
		Block: block2,
	}
	blocks[1] = signedBlock
	var message = NewGiveBlocksMessage(blocks)
	fmt.Println("GiveBlocksMessage:")
	var mai = NewMessagesAnnotationsIterator(message)
	w := bufio.NewWriter(os.Stdout)
	err := NewFromIterator(gnet.EncodeMessage(message), &mai, w)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
	// GiveBlocksMessage:
	// 0x0000 | 8a 01 00 00 ....................................... Length
	// 0x0004 | 47 49 56 42 ....................................... Prefix
	// 0x0008 | 02 00 00 00 ....................................... Blocks length
	// 0x000c | 02 00 00 00 64 00 00 00 00 00 00 00 00 00 00 00
	// 0x001c | 00 00 00 00 0a 00 00 00 00 00 00 00 40 af f2 e9
	// 0x002c | d2 d8 92 2e 47 af d4 64 8e 69 67 49 71 58 78 5f
	// 0x003c | bd 1d a8 70 e7 11 02 66 bf 94 48 80 00 00 00 00
	// 0x004c | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x005c | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x006c | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x007c | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x008c | 03 21 3f dd 6d df 86 0e 40 53 e1 a9 7e 42 76 d6
	// 0x009c | 34 54 f5 19 5a 83 21 35 70 04 d5 2c db bf d3 88
	// 0x00ac | 6f c7 ad 3f 3f 63 b6 5d 4a 87 9c e3 08 6d ae b3
	// 0x00bc | e5 4a 93 d3 c2 f9 6a 50 61 f9 bc 49 36 83 ca 8e
	// 0x00cc | 01 ................................................ Blocks[0]
	// 0x00cd | 02 00 00 00 64 00 00 00 00 00 00 00 00 00 00 00
	// 0x00dd | 00 00 00 00 fa 00 00 00 00 00 00 00 7b b4 62 c3
	// 0x00ed | bd 37 1d d8 1c 06 ad 1d 2b 63 59 71 cb 56 eb 22
	// 0x00fd | 23 3d fc 9f eb e8 3e 44 c8 40 b8 d7 00 00 00 00
	// 0x010d | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x011d | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x012d | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x013d | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x014d | ae ab ca 1d ec b2 ab 5b a8 af 93 c6 f0 33 de 6c
	// 0x015d | bf 1d 50 31 4d f2 75 f9 40 77 8b c7 20 43 3c bc
	// 0x016d | 21 94 aa cf a2 43 cc 2a 21 f8 5f 2f ff 71 d3 16
	// 0x017d | 6d 18 75 e1 98 1a 0d a5 a2 3d 28 96 81 fc 1f a3
	// 0x018d | 00 ................................................ Blocks[1]
	// 0x018e |
}

func ExampleAnnounceBlocksMessage() {
	defer gnet.EraseMessages()
	setupMsgEncoding()
	var message = NewAnnounceBlocksMessage(123456)
	fmt.Println("AnnounceBlocksMessage:")
	var mai = NewMessagesAnnotationsIterator(message)
	w := bufio.NewWriter(os.Stdout)
	err := NewFromIterator(gnet.EncodeMessage(message), &mai, w)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
	// AnnounceBlocksMessage:
	// 0x0000 | 0c 00 00 00 ....................................... Length
	// 0x0004 | 41 4e 4e 42 ....................................... Prefix
	// 0x0008 | 40 e2 01 00 00 00 00 00 ........................... MaxBkSeq
	// 0x0010 |
}

func ExampleGetTxnsMessage() {
	defer gnet.EraseMessages()
	setupMsgEncoding()
	var shas = make([]cipher.SHA256, 0)

	shas = append(shas, hashes[1], hashes[2])
	var message = NewGetTxnsMessage(shas)
	fmt.Println("GetTxnsMessage:")
	var mai = NewMessagesAnnotationsIterator(message)
	w := bufio.NewWriter(os.Stdout)
	err := NewFromIterator(gnet.EncodeMessage(message), &mai, w)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
	// GetTxnsMessage:
	// 0x0000 | 48 00 00 00 ....................................... Length
	// 0x0004 | 47 45 54 54 ....................................... Prefix
	// 0x0008 | 02 00 00 00 ....................................... Transactions length
	// 0x000c | 7b b4 62 c3 bd 37 1d d8 1c 06 ad 1d 2b 63 59 71
	// 0x001c | cb 56 eb 22 23 3d fc 9f eb e8 3e 44 c8 40 b8 d7 ... Transactions[0]
	// 0x002c | e7 5a c8 01 c1 3f 3d a9 c7 a1 24 ca 31 3b e2 a3
	// 0x003c | 73 f6 4a d9 7c 58 a1 b6 fe bc 0e 0c a5 c5 c8 73 ... Transactions[1]
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

	var sig0, sig1, sig2, sig3 cipher.Sig
	sig0, _ = cipher.SigFromHex(sig1hex) // nolint: errcheck
	sig1, _ = cipher.SigFromHex(sig2hex) // nolint: errcheck
	sig2, _ = cipher.SigFromHex(sig3hex) // nolint: errcheck
	sig3, _ = cipher.SigFromHex(sig4hex) // nolint: errcheck
	var transaction0 = coin.Transaction{
		Type:      123,
		In:        []cipher.SHA256{hashes[3], hashes[4]},
		InnerHash: hashes[5],
		Length:    5000,
		Out:       transactionOutputs0,
		Sigs:      []cipher.Sig{sig0, sig1},
	}
	var transaction1 = coin.Transaction{
		Type:      123,
		In:        []cipher.SHA256{hashes[5], hashes[6]},
		InnerHash: hashes[6],
		Length:    5000,
		Out:       transactionOutputs1,
		Sigs:      []cipher.Sig{sig2, sig3},
	}
	transactions = append(transactions, transaction0, transaction1)
	var message = NewGiveTxnsMessage(transactions)
	fmt.Println("GiveTxnsMessage:")
	var mai = NewMessagesAnnotationsIterator(message)
	w := bufio.NewWriter(os.Stdout)
	err := NewFromIterator(gnet.EncodeMessage(message), &mai, w)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
	// GiveTxnsMessage:
	// 0x0000 | 82 02 00 00 ....................................... Length
	// 0x0004 | 47 49 56 54 ....................................... Prefix
	// 0x0008 | 02 00 00 00 ....................................... Transactions length
	// 0x000c | 88 13 00 00 7b 38 62 f1 93 a1 56 4e 5e 26 0f 82
	// 0x001c | 7d a8 e1 69 ca d8 11 d8 1d 6a 7c 4f fd 66 1c 00
	// 0x002c | b1 99 94 17 81 02 00 00 00 03 21 3f dd 6d df 86
	// 0x003c | 0e 40 53 e1 a9 7e 42 76 d6 34 54 f5 19 5a 83 21
	// 0x004c | 35 70 04 d5 2c db bf d3 88 6f c7 ad 3f 3f 63 b6
	// 0x005c | 5d 4a 87 9c e3 08 6d ae b3 e5 4a 93 d3 c2 f9 6a
	// 0x006c | 50 61 f9 bc 49 36 83 ca 8e 01 ae ab ca 1d ec b2
	// 0x007c | ab 5b a8 af 93 c6 f0 33 de 6c bf 1d 50 31 4d f2
	// 0x008c | 75 f9 40 77 8b c7 20 43 3c bc 21 94 aa cf a2 43
	// 0x009c | cc 2a 21 f8 5f 2f ff 71 d3 16 6d 18 75 e1 98 1a
	// 0x00ac | 0d a5 a2 3d 28 96 81 fc 1f a3 00 02 00 00 00 f4
	// 0x00bc | 45 7d e9 f5 a5 94 2e 07 6a 7f 2b 28 e1 84 2a b6
	// 0x00cc | 1f 1b fc 39 e4 ca 55 75 36 60 0f d6 42 09 f6 c4
	// 0x00dc | 22 09 82 f9 88 b6 25 d0 af c1 2c 7f fd 06 a7 fe
	// 0x00ec | 89 bb e6 60 2c 1f 20 d9 08 91 3f e9 38 10 47 02
	// 0x00fc | 00 00 00 00 ad dc d4 a7 19 6a 8c a8 6b 9b 3d 74
	// 0x010c | 16 95 f3 69 ef 1b 3d ba 0c 00 00 00 00 00 00 00
	// 0x011c | 22 00 00 00 00 00 00 00 00 31 7c 95 cb 68 79 ac
	// 0x012c | 7e e5 3a f9 98 c0 b6 70 37 6c 7c 51 38 38 00 00
	// 0x013c | 00 00 00 00 00 4e 00 00 00 00 00 00 00 ............ Transactions[0]
	// 0x0149 | 88 13 00 00 7b 05 64 0e 44 80 73 9e 87 97 57 b0
	// 0x0159 | a2 d1 bd 59 de a7 df cc fe f3 df 75 a1 83 0a 50
	// 0x0169 | 20 01 10 67 21 02 00 00 00 67 65 27 8a fc 9f 3e
	// 0x0179 | 0f fb 95 cb b3 f0 18 72 e9 2e d1 d5 1e 7a 83 d1
	// 0x0189 | 6d 49 9e 95 97 e2 4f a6 f3 08 d8 85 c8 31 c4 6b
	// 0x0199 | 69 9a d6 7b 2f dd 2f 76 2f d6 5f 4f bc f6 6f 98
	// 0x01a9 | 1d 76 a9 ad fe 42 0d 16 14 01 5e db 9b bd 4a 78
	// 0x01b9 | 20 f4 53 91 f0 8e 75 72 c9 1f b6 1b e7 10 6d db
	// 0x01c9 | 53 46 ac 00 59 6c 89 92 cf a2 3a ac 06 46 9a ec
	// 0x01d9 | 55 33 71 02 41 8c eb 11 52 9b 3b 30 57 84 02 e5
	// 0x01e9 | 8c 2f f6 6d cd 8b bf 35 21 e9 01 02 00 00 00 38
	// 0x01f9 | 62 f1 93 a1 56 4e 5e 26 0f 82 7d a8 e1 69 ca d8
	// 0x0209 | 11 d8 1d 6a 7c 4f fd 66 1c 00 b1 99 94 17 81 05
	// 0x0219 | 64 0e 44 80 73 9e 87 97 57 b0 a2 d1 bd 59 de a7
	// 0x0229 | df cc fe f3 df 75 a1 83 0a 50 20 01 10 67 21 02
	// 0x0239 | 00 00 00 00 89 e6 17 1b d0 e4 8b fc 31 68 a5 f0
	// 0x0249 | 3c 5e 16 10 8a ed e7 6f 09 00 00 00 00 00 00 00
	// 0x0259 | 0c 00 00 00 00 00 00 00 00 76 60 d2 e5 38 f9 69
	// 0x0269 | 46 5f 22 a2 bb ae 9b a8 4e c5 f3 77 a1 22 00 00
	// 0x0279 | 00 00 00 00 00 38 00 00 00 00 00 00 00 ............ Transactions[1]
	// 0x0286 |
}

func ExampleAnnounceTxnsMessage() {
	defer gnet.EraseMessages()
	setupMsgEncoding()
	var message = NewAnnounceTxnsMessage([]cipher.SHA256{hashes[7], hashes[8]})
	fmt.Println("AnnounceTxnsMessage:")
	var mai = NewMessagesAnnotationsIterator(message)
	w := bufio.NewWriter(os.Stdout)
	err := NewFromIterator(gnet.EncodeMessage(message), &mai, w)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
	// AnnounceTxnsMessage:
	// 0x0000 | 48 00 00 00 ....................................... Length
	// 0x0004 | 41 4e 4e 54 ....................................... Prefix
	// 0x0008 | 02 00 00 00 ....................................... Transactions length
	// 0x000c | 8a 5d bf bb 7e 64 66 49 5e 30 78 1c 15 40 b5 e3
	// 0x001c | 98 e0 84 4f 60 c9 1e c6 78 9d 4b bb 36 7e 33 a6 ... Transactions[0]
	// 0x002c | 1c 1d 7d bf d7 ba 2b b1 aa 9b 56 ed ae 26 ea 56
	// 0x003c | 5c bf 72 f9 8c c6 a6 2c 72 97 23 cb c0 75 0d 3b ... Transactions[1]
	// 0x004c |
}

func ExampleDisconnectMessage() {
	defer gnet.EraseMessages()
	setupMsgEncoding()

	message := NewDisconnectMessage(ErrDisconnectIdle)
	fmt.Println("DisconnectMessage:")
	var mai = NewMessagesAnnotationsIterator(message)
	w := bufio.NewWriter(os.Stdout)
	err := NewFromIterator(gnet.EncodeMessage(message), &mai, w)
	if err != nil {
		fmt.Println(err)
	}
	// DisconnectMessage:
	// 0x0000 | 31 00 00 00 ....................................... Length
	// 0x0004 | 52 4a 43 54 ....................................... Prefix
	// 0x0008 | 49 4e 54 52 ....................................... TargetPrefix
	// 0x000c | 13 00 00 00 ....................................... ErrorCode
	// 0x0010 | 1d 00 00 00 45 78 61 6d 70 6c 65 52 65 6a 65 63
	// 0x0020 | 74 57 69 74 68 50 65 65 72 73 4d 65 73 73 61 67
	// 0x0030 | 65 ................................................ Reason
	// 0x0031 | 00 00 00 00 ....................................... Reserved length
	// 0x0035 |

}
