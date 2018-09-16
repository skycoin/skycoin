package daemon

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/daemon/gnet"
	"github.com/skycoin/skycoin/src/daemon/pex"
	"github.com/skycoin/skycoin/src/util/hexdump"
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
func (mai *MessagesAnnotationsIterator) Next() (hexdump.Annotation, bool) {
	if !mai.LengthCalled {
		mai.LengthCalled = true
		return hexdump.Annotation{Size: 4, Name: "Length"}, true
	}
	if !mai.PrefixCalled {
		mai.PrefixCalled = true
		return hexdump.Annotation{Size: 4, Name: "Prefix"}, true

	}
	if mai.CurrentField >= mai.MaxField {
		return hexdump.Annotation{}, false
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
							return hexdump.Annotation{}, false
						}
					} else {
						panic(encoder.ErrInvalidOmitEmpty)
					}
				}
			}
		} else {
			return hexdump.Annotation{}, false
		}
	}
	if f.Tag.Get("enc") != "-" {
		if vF.CanSet() || f.Name != "_" {
			if v.Field(i).Kind() == reflect.Slice {
				if mai.CurrentIndex == -1 {
					mai.CurrentIndex = 0
					return hexdump.Annotation{Size: 4, Name: f.Name + " length"}, true
				}
				sliceLen := v.Field(i).Len()
				mai.CurrentIndex++
				if mai.CurrentIndex < sliceLen {
					// Emit annotation for slice item
					return hexdump.Annotation{Size: len(encoder.Serialize(v.Field(i).Slice(j, j+1).Interface())[4:]), Name: f.Name + "[" + strconv.Itoa(j) + "]"}, true
				}
				// No more annotation tokens for current slice field
				mai.CurrentIndex = -1
				mai.CurrentField++
				if sliceLen > 0 {
					// Emit annotation for last item
					return hexdump.Annotation{Size: len(encoder.Serialize(v.Field(i).Slice(j, j+1).Interface())[4:]), Name: f.Name + "[" + strconv.Itoa(j) + "]"}, true
				}
				// Zero length slice. Start over
				return mai.Next()
			}

			mai.CurrentField++
			return hexdump.Annotation{Size: len(encoder.Serialize(v.Field(i).Interface())), Name: f.Name}, true

		}
	}

	return hexdump.Annotation{}, false
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

var secKey1 = (cipher.NewSecKey([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32}))
var secKey2 = cipher.NewSecKey([]byte{33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64})
var secKey3 = cipher.NewSecKey([]byte{65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96})
var secKey4 = cipher.NewSecKey([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96})

var addresses = []cipher.Address{
	cipher.AddressFromSecKey(secKey1),
	cipher.AddressFromSecKey(secKey2),
	cipher.AddressFromSecKey(secKey3),
	cipher.AddressFromSecKey(secKey4),
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
	err := hexdump.NewFromIterator(gnet.EncodeMessage(&message), &mai, w)
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
	err := hexdump.NewFromIterator(gnet.EncodeMessage(&message), &mai, w)
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
	var mai = NewMessagesAnnotationsIterator(message)
	w := bufio.NewWriter(os.Stdout)
	err := hexdump.NewFromIterator(gnet.EncodeMessage(message), &mai, w)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
	// IntroductionMessage:
	// 0x0000 | 0e 00 00 00 ....................................... Length
	// 0x0004 | 49 4e 54 52 ....................................... Prefix
	// 0x0008 | d2 04 00 00 ....................................... Mirror
	// 0x000c | d2 1e ............................................. Port
	// 0x000e | 05 00 00 00 ....................................... Version
	// 0x0012 |
}

func ExampleGetPeersMessage() {
	defer gnet.EraseMessages()
	setupMsgEncoding()
	var message = NewGetPeersMessage()
	fmt.Println("GetPeersMessage:")
	var mai = NewMessagesAnnotationsIterator(message)
	w := bufio.NewWriter(os.Stdout)
	err := hexdump.NewFromIterator(gnet.EncodeMessage(message), &mai, w)
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
	err := hexdump.NewFromIterator(gnet.EncodeMessage(message), &mai, w)
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
	err := hexdump.NewFromIterator(gnet.EncodeMessage(message), &mai, w)
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
	err := hexdump.NewFromIterator(gnet.EncodeMessage(message), &mai, w)
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
	err := hexdump.NewFromIterator(gnet.EncodeMessage(message), &mai, w)
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
	err := hexdump.NewFromIterator(gnet.EncodeMessage(message), &mai, w)
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
	err := hexdump.NewFromIterator(gnet.EncodeMessage(message), &mai, w)
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
	// 0x00fc | 00 00 00 00 07 6d ca 32 de 03 4e 48 67 fa 7a 2a
	// 0x010c | a9 ee fe 91 f2 0b a0 74 0c 00 00 00 00 00 00 00
	// 0x011c | 22 00 00 00 00 00 00 00 00 e9 cb 47 35 e3 95 cf
	// 0x012c | 36 b0 d1 a6 f2 21 bb 23 b3 f7 bf b1 f9 38 00 00
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
	// 0x0239 | 00 00 00 00 7e f9 b1 b9 40 6f 8d b3 99 b2 5f d0
	// 0x0249 | e9 f4 f0 88 7b 08 4b 43 09 00 00 00 00 00 00 00
	// 0x0259 | 0c 00 00 00 00 00 00 00 00 83 f1 96 59 16 14 99
	// 0x0269 | 2f a6 03 13 38 6f 72 88 ac 40 14 c8 bc 22 00 00
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
	err := hexdump.NewFromIterator(gnet.EncodeMessage(message), &mai, w)
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
