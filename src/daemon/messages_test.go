package daemon

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strconv"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/daemon/gnet"
	"github.com/skycoin/skycoin/src/daemon/pex"
	"github.com/skycoin/skycoin/src/util"
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
func (mai *MessagesAnnotationsIterator) Next() (util.Annotation, bool) {
	if !mai.LengthCalled {
		mai.LengthCalled = true
		return util.Annotation{Size: 4, Name: "Length"}, true
	}
	if !mai.PrefixCalled {
		mai.PrefixCalled = true
		return util.Annotation{Size: 4, Name: "Prefix"}, true

	}
	if mai.CurrentField >= mai.MaxField {
		return util.Annotation{}, false
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
							return util.Annotation{}, false
						}
					} else {
						panic(encoder.ErrInvalidOmitEmpty)
					}
				}
			}
		} else {
			return util.Annotation{}, false
		}
	}
	if f.Tag.Get("enc") != "-" {
		if vF.CanSet() || f.Name != "_" {
			if v.Field(i).Kind() == reflect.Slice {
				if mai.CurrentIndex == -1 {
					mai.CurrentIndex = 0
					return util.Annotation{Size: 4, Name: f.Name + " length"}, true
				}
				sliceLen := v.Field(i).Len()
				mai.CurrentIndex++
				if mai.CurrentIndex < sliceLen {
					// Emit annotation for slice item
					return util.Annotation{Size: len(encoder.Serialize(v.Field(i).Slice(j, j+1).Interface())[4:]), Name: f.Name + "[" + strconv.Itoa(j) + "]"}, true
				}
				// No more annotation tokens for current slice field
				mai.CurrentIndex = -1
				mai.CurrentField++
				if sliceLen > 0 {
					// Emit annotation for last item
					return util.Annotation{Size: len(encoder.Serialize(v.Field(i).Slice(j, j+1).Interface())[4:]), Name: f.Name + "[" + strconv.Itoa(j) + "]"}, true
				}
				// Zero length slice. Start over
				return mai.Next()
			}

			mai.CurrentField++
			return util.Annotation{Size: len(encoder.Serialize(v.Field(i).Interface())), Name: f.Name}, true

		}
	}

	return util.Annotation{}, false
}

/**************************************
 *
 * Test cases
 *
 *************************************/

var hashes = []cipher.SHA256{
	GetSHAFromHex("123"),
	GetSHAFromHex("456"),
	GetSHAFromHex("789"),
	GetSHAFromHex("abc"),
	GetSHAFromHex("def"),
	GetSHAFromHex("101"),
	GetSHAFromHex("111"),
	GetSHAFromHex("121"),
	GetSHAFromHex("314"),
	GetSHAFromHex("151"),
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
	util.HexDumpFromIterator(gnet.EncodeMessage(&message), &mai, w)
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
	util.HexDumpFromIterator(gnet.EncodeMessage(&message), &mai, w)
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
	var message = NewIntroductionMessage(1234, 5, 7890)
	fmt.Println("IntroductionMessage:")
	var mai = NewMessagesAnnotationsIterator(message)
	w := bufio.NewWriter(os.Stdout)
	util.HexDumpFromIterator(gnet.EncodeMessage(message), &mai, w)
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
	util.HexDumpFromIterator(gnet.EncodeMessage(message), &mai, w)
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
	util.HexDumpFromIterator(gnet.EncodeMessage(message), &mai, w)
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
	util.HexDumpFromIterator(gnet.EncodeMessage(message), &mai, w)
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
	var blocks = make([]coin.SignedBlock, 1)
	var body1 = coin.BlockBody{
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
	var sig, _ = cipher.SigFromHex("123")
	var signedBlock = coin.SignedBlock{
		Sig:   sig,
		Block: block1,
	}
	blocks = append(blocks, signedBlock)
	var message = NewGiveBlocksMessage(blocks)
	fmt.Println("GiveBlocksMessage:")
	var mai = NewMessagesAnnotationsIterator(message)
	w := bufio.NewWriter(os.Stdout)
	util.HexDumpFromIterator(gnet.EncodeMessage(message), &mai, w)
	// Output:
	// GiveBlocksMessage:
	// 0x0000 | 8a 01 00 00 ....................................... Length
	// 0x0004 | 47 49 56 42 ....................................... Prefix
	// 0x0008 | 02 00 00 00 ....................................... Blocks length
	// 0x000c | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x001c | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x002c | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x003c | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x004c | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x005c | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x006c | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x007c | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x008c | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x009c | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x00ac | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x00bc | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x00cc | 00 ................................................ Blocks[0]
	// 0x00cd | 02 00 00 00 64 00 00 00 00 00 00 00 00 00 00 00
	// 0x00dd | 00 00 00 00 0a 00 00 00 00 00 00 00 00 00 00 00
	// 0x00ed | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x00fd | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x010d | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x011d | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x012d | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x013d | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x014d | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x015d | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x016d | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x017d | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
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
	util.HexDumpFromIterator(gnet.EncodeMessage(message), &mai, w)
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
	util.HexDumpFromIterator(gnet.EncodeMessage(message), &mai, w)
	// Output:
	// GetTxnsMessage:
	// 0x0000 | 48 00 00 00 ....................................... Length
	// 0x0004 | 47 45 54 54 ....................................... Prefix
	// 0x0008 | 02 00 00 00 ....................................... Txns length
	// 0x000c | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x001c | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 ... Txns[0]
	// 0x002c | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x003c | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 ... Txns[1]
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
	sig0, _ = cipher.SigFromHex("sig0")
	sig1, _ = cipher.SigFromHex("sig1")
	sig2, _ = cipher.SigFromHex("sig2")
	sig3, _ = cipher.SigFromHex("sig3")
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
	util.HexDumpFromIterator(gnet.EncodeMessage(message), &mai, w)
	// Output:
	// GiveTxnsMessage:
	// 0x0000 | 82 02 00 00 ....................................... Length
	// 0x0004 | 47 49 56 54 ....................................... Prefix
	// 0x0008 | 02 00 00 00 ....................................... Txns length
	// 0x000c | 88 13 00 00 7b 00 00 00 00 00 00 00 00 00 00 00
	// 0x001c | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x002c | 00 00 00 00 00 02 00 00 00 00 00 00 00 00 00 00
	// 0x003c | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x004c | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x005c | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x006c | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x007c | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x008c | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x009c | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x00ac | 00 00 00 00 00 00 00 00 00 00 00 02 00 00 00 00
	// 0x00bc | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x00cc | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x00dc | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x00ec | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 02
	// 0x00fc | 00 00 00 00 07 6d ca 32 de 03 4e 48 67 fa 7a 2a
	// 0x010c | a9 ee fe 91 f2 0b a0 74 0c 00 00 00 00 00 00 00
	// 0x011c | 22 00 00 00 00 00 00 00 00 e9 cb 47 35 e3 95 cf
	// 0x012c | 36 b0 d1 a6 f2 21 bb 23 b3 f7 bf b1 f9 38 00 00
	// 0x013c | 00 00 00 00 00 4e 00 00 00 00 00 00 00 ............ Txns[0]
	// 0x0149 | 88 13 00 00 7b 00 00 00 00 00 00 00 00 00 00 00
	// 0x0159 | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x0169 | 00 00 00 00 00 02 00 00 00 00 00 00 00 00 00 00
	// 0x0179 | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x0189 | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x0199 | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x01a9 | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x01b9 | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x01c9 | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x01d9 | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x01e9 | 00 00 00 00 00 00 00 00 00 00 00 02 00 00 00 00
	// 0x01f9 | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x0209 | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x0219 | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x0229 | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 02
	// 0x0239 | 00 00 00 00 7e f9 b1 b9 40 6f 8d b3 99 b2 5f d0
	// 0x0249 | e9 f4 f0 88 7b 08 4b 43 09 00 00 00 00 00 00 00
	// 0x0259 | 0c 00 00 00 00 00 00 00 00 83 f1 96 59 16 14 99
	// 0x0269 | 2f a6 03 13 38 6f 72 88 ac 40 14 c8 bc 22 00 00
	// 0x0279 | 00 00 00 00 00 38 00 00 00 00 00 00 00 ............ Txns[1]
	// 0x0286 |
}

func ExampleAnnounceTxnsMessage() {
	defer gnet.EraseMessages()
	setupMsgEncoding()
	var message = NewAnnounceTxnsMessage([]cipher.SHA256{hashes[7], hashes[8]})
	fmt.Println("AnnounceTxnsMessage:")
	var mai = NewMessagesAnnotationsIterator(message)
	w := bufio.NewWriter(os.Stdout)
	util.HexDumpFromIterator(gnet.EncodeMessage(message), &mai, w)
	// Output:
	// AnnounceTxnsMessage:
	// 0x0000 | 48 00 00 00 ....................................... Length
	// 0x0004 | 41 4e 4e 54 ....................................... Prefix
	// 0x0008 | 02 00 00 00 ....................................... Txns length
	// 0x000c | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x001c | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 ... Txns[0]
	// 0x002c | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x003c | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 ... Txns[1]
	// 0x004c |
}
