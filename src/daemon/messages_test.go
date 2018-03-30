package daemon

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/daemon/gnet"
	"github.com/skycoin/skycoin/src/daemon/pex"
)

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

var addresses = []cipher.Address{
	cipher.Address{
		Version: 1,
		Key:     cipher.HashRipemd160([]byte("123")),
	},
	cipher.Address{
		Version: 1,
		Key:     cipher.HashRipemd160([]byte("456")),
	},
	cipher.Address{
		Version: 1,
		Key:     cipher.HashRipemd160([]byte("789")),
	},
	cipher.Address{
		Version: 1,
		Key:     cipher.HashRipemd160([]byte("abc")),
	},
	cipher.Address{
		Version: 1,
		Key:     cipher.HashRipemd160([]byte("def")),
	},
	cipher.Address{
		Version: 1,
		Key:     cipher.HashRipemd160([]byte("101")),
	},
}

func GetSHAFromHex(hex string) cipher.SHA256 {
	var sha, _ = cipher.SHA256FromHex(hex)
	return sha
}

func setupMsgEncoding() {
	gnet.EraseMessages()
	var messagesConfig = NewMessagesConfig()
	messagesConfig.Register()
}

func MessageHexDump(message gnet.Message, printFull bool) {
	var serializedMsg = gnet.EncodeMessage(message)

	PrintLHexDumpWithFormat(-1, "Full message", serializedMsg)

	fmt.Println("------------------------------------------------------------------------")
	var offset int = 0
	PrintLHexDumpWithFormat(0, "Length", serializedMsg[0:4])
	PrintLHexDumpWithFormat(4, "Prefix", serializedMsg[4:8])
	offset += len(serializedMsg[0:8])
	var v = reflect.Indirect(reflect.ValueOf(message))

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		v_f := v.Field(i)
		f := t.Field(i)
		if f.Tag.Get("enc") != "-" {
			if v_f.CanSet() || f.Name != "_" {
				if v.Field(i).Kind() == reflect.Slice {
					PrintLHexDumpWithFormat(offset, f.Name+" length", encoder.Serialize(v.Field(i).Slice(0, v.Field(i).Len()).Interface())[0:4])
					offset += len(encoder.Serialize(v.Field(i).Slice(0, v.Field(i).Len()).Interface())[0:4])

					for j := 0; j < v.Field(i).Len(); j++ {
						PrintLHexDumpWithFormat(offset, f.Name+"#"+strconv.Itoa(j), encoder.Serialize(v.Field(i).Slice(j, j+1).Interface()))
						offset += len(encoder.Serialize(encoder.Serialize(v.Field(i).Slice(j, j+1).Interface())))
					}
				} else {
					PrintLHexDumpWithFormat(offset, f.Name, encoder.Serialize(v.Field(i).Interface()))
					offset += len(encoder.Serialize(v.Field(i).Interface()))
				}
			} else {
				//don't write anything
			}
		}
	}

	PrintFinalHex(len(serializedMsg))
}

func ExampleNewIntroductionMessage() {
	defer gnet.EraseMessages()
	setupMsgEncoding()

	var message = NewIntroductionMessage(1234, 5, 7890)
	fmt.Println("IntroductionMessage:")
	MessageHexDump(message, true)
	// Output:
	// IntroductionMessage:
	// 0e 00 00 00 49 4e 54 52 d2 04 00 00 d2 1e 05 00
	// 00 00 ............................................. Full message
	// ------------------------------------------------------------------------
	// 0x0000 | 0e 00 00 00 ....................................... Length
	// 0x0004 | 49 4e 54 52 ....................................... Prefix
	// 0x0008 | d2 04 00 00 ....................................... Mirror
	// 0x000c | d2 1e ............................................. Port
	// 0x000e | 05 00 00 00 ....................................... Version
	// 0x0012 |
}

func ExampleNewGetPeersMessage() {
	defer gnet.EraseMessages()
	setupMsgEncoding()

	var message = NewGetPeersMessage()
	fmt.Println("GetPeersMessage:")
	MessageHexDump(message, true)
	// Output:
	// GetPeersMessage:
	// 04 00 00 00 47 45 54 50 ........................... Full message
	// ------------------------------------------------------------------------
	// 0x0000 | 04 00 00 00 ....................................... Length
	// 0x0004 | 47 45 54 50 ....................................... Prefix
	// 0x0008 |
}

func ExampleNewGivePeersMessage() {
	defer gnet.EraseMessages()
	setupMsgEncoding()

	var peers = make([]pex.Peer, 3)
	var peer0 pex.Peer = *pex.NewPeer("118.178.135.93:6000")
	var peer1 pex.Peer = *pex.NewPeer("47.88.33.156:6000")
	var peer2 pex.Peer = *pex.NewPeer("121.41.103.148:6000")
	peers = append(peers, peer0, peer1, peer2)
	var message = NewGivePeersMessage(peers)
	fmt.Println("GivePeersMessage:")
	MessageHexDump(message, true)
	// Output:
	// GivePeersMessage:
	// 1a 00 00 00 47 49 56 50 03 00 00 00 5d 87 b2 76
	// 70 17 9c 21 58 2f 70 17 94 67 29 79 70 17 ......... Full message
	// ------------------------------------------------------------------------
	// 0x0000 | 1a 00 00 00 ....................................... Length
	// 0x0004 | 47 49 56 50 ....................................... Prefix
	// 0x0008 | 03 00 00 00 ....................................... Peers length
	// 0x000c | 01 00 00 00 5d 87 b2 76 70 17 ..................... Peers#0
	// 0x001a | 01 00 00 00 9c 21 58 2f 70 17 ..................... Peers#1
	// 0x0028 | 01 00 00 00 94 67 29 79 70 17 ..................... Peers#2
	// 0x001e |
}

/*
func TestNewPingMessage(t *testing.T) {
	var message gnet.Message = &PingMessage{}
	MessageHexDump(message, true)
}


func ExampleNewPongMessage() {
	//var message = PongMessage{
	//}
	//MessageHexDump(message, true)
}

*/

func ExampleNewGetBlocksMessage() {
	defer gnet.EraseMessages()
	setupMsgEncoding()

	var message = NewGetBlocksMessage(1234, 5678)
	fmt.Println("GetBlocksMessage:")
	MessageHexDump(message, true)
	// Output:
	// GetBlocksMessage:
	// 14 00 00 00 47 45 54 42 d2 04 00 00 00 00 00 00
	// 2e 16 00 00 00 00 00 00 ........................... Full message
	// ------------------------------------------------------------------------
	// 0x0000 | 14 00 00 00 ....................................... Length
	// 0x0004 | 47 45 54 42 ....................................... Prefix
	// 0x0008 | d2 04 00 00 00 00 00 00 ........................... LastBlock
	// 0x0010 | 2e 16 00 00 00 00 00 00 ........................... RequestedBlocks
	// 0x0018 |
}

func ExampleNewGiveBlocksMessage() {
	defer gnet.EraseMessages()
	setupMsgEncoding()

	var blocks = make([]coin.SignedBlock, 1)
	var body1 = coin.BlockBody{
		Transactions: make([]coin.Transaction, 0),
	}
	var block1 coin.Block = coin.Block{
		Body: body1,
		Head: coin.BlockHeader{
			Version:  0x02,
			Time:     100,
			BkSeq:    0,
			Fee:      10,
			PrevHash: hashes[8],
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
	MessageHexDump(message, true)
	// Output:
	// 8a 01 00 00 47 49 56 42 02 00 00 00 00 00 00 00
	// 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 00 00 00 00 00 00 00 00 00 00 00 00 00 02 00 00
	// 00 64 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 00 0a 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 00 00 00 00 00 00 00 00 00 00 00 00 00 00 ......... Full message
	// ------------------------------------------------------------------------
	// 0x0000 | 8a 01 00 00 ....................................... Length
	// 0x0004 | 47 49 56 42 ....................................... Prefix
	// 0x0008 | 02 00 00 00 ....................................... Blocks length
	// 0x000c | 01 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
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
	// 0x00cc | 00 00 00 00 00 .................................... Blocks#0
	// 0x00d5 | 01 00 00 00 02 00 00 00 64 00 00 00 00 00 00 00
	// 0x00e5 | 00 00 00 00 00 00 00 00 0a 00 00 00 00 00 00 00
	// 0x00f5 | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x0105 | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x0115 | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x0125 | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x0135 | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x0145 | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x0155 | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x0165 | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x0175 | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x0185 | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x0195 | 00 00 00 00 00 .................................... Blocks#1
	// 0x018e |
}

func ExampleNewAnnounceBlocksMessage() {
	defer gnet.EraseMessages()
	setupMsgEncoding()

	var message = NewAnnounceBlocksMessage(123456)
	fmt.Println("AnnounceBlocksMessage:")
	MessageHexDump(message, true)
	// Output:
	// AnnounceBlocksMessage:
	// 0c 00 00 00 41 4e 4e 42 40 e2 01 00 00 00 00 00
	// ................................................... Full message
	// ------------------------------------------------------------------------
	// 0x0000 | 0c 00 00 00 ....................................... Length
	// 0x0004 | 41 4e 4e 42 ....................................... Prefix
	// 0x0008 | 40 e2 01 00 00 00 00 00 ........................... MaxBkSeq
	// 0x0010 |
}

func ExampleNewGetTxnsMessage() {
	defer gnet.EraseMessages()
	setupMsgEncoding()

	var shas = make([]cipher.SHA256, 0)

	shas = append(shas, hashes[0], hashes[1])
	var message = NewGetTxnsMessage(shas)
	fmt.Println("GetTxns:")
	MessageHexDump(message, true)
	// Output:
	// GetTxns:
	// 48 00 00 00 47 45 54 54 02 00 00 00 00 00 00 00
	// 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 00 00 00 00 00 00 00 00 00 00 00 00 ............... Full message
	// ------------------------------------------------------------------------
	// 0x0000 | 48 00 00 00 ....................................... Length
	// 0x0004 | 47 45 54 54 ....................................... Prefix
	// 0x0008 | 02 00 00 00 ....................................... Txns length
	// 0x000c | 01 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x001c | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x002c | 00 00 00 00 ....................................... Txns#0
	// 0x0034 | 01 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x0044 | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x0054 | 00 00 00 00 ....................................... Txns#1
	// 0x004c |
}

func ExampleNewGiveTxnsMessage() {
	defer gnet.EraseMessages()
	setupMsgEncoding()

	var transactions coin.Transactions = make([]coin.Transaction, 0)
	var transactionOutputs0 []coin.TransactionOutput = make([]coin.TransactionOutput, 0)
	var transactionOutputs1 []coin.TransactionOutput = make([]coin.TransactionOutput, 0)
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
		Address: addresses[2],
		Coins:   9,
		Hours:   12,
	}
	var txOutput3 = coin.TransactionOutput{
		Address: addresses[3],
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
		In:        []cipher.SHA256{hashes[2], hashes[3]},
		InnerHash: hashes[4],
		Length:    5000,
		Out:       transactionOutputs0,
		Sigs:      []cipher.Sig{sig0, sig1},
	}
	var transaction1 = coin.Transaction{
		Type:   123,
		In:     []cipher.SHA256{hashes[5], hashes[6], hashes[7]},
		Length: 5000,
		Out:    transactionOutputs1,
		Sigs:   []cipher.Sig{sig2, sig3},
	}
	transactions = append(transactions, transaction0, transaction1)
	var message = NewGiveTxnsMessage(transactions)
	fmt.Println("GiveTxnsMessage:")
	MessageHexDump(message, true)
	// Output:
	// GiveTxnsMessage:
	// a2 02 00 00 47 49 56 54 02 00 00 00 88 13 00 00
	// 7b 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 00 02 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 00 00 00 00 00 00 00 02 00 00 00 00 00 00 00 00
	// 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 00 00 00 00 00 00 00 00 00 00 00 02 00 00 00 01
	// e3 43 1a 8e 0a db f9 6f d1 40 10 3d c6 f6 3a 3f
	// 8f a3 43 ab 0c 00 00 00 00 00 00 00 22 00 00 00
	// 00 00 00 00 01 08 4f 32 85 28 6c a7 84 bb ff d8
	// 23 de 73 62 76 03 64 c9 a8 38 00 00 00 00 00 00
	// 00 4e 00 00 00 00 00 00 00 88 13 00 00 7b 00 00
	// 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 00 00 00 00 00 00 00 00 00 00 00 00 00 00 02 00
	// 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 00 00 00 00 03 00 00 00 00 00 00 00 00 00 00 00
	// 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 00 00 00 00 00 00 00 00 02 00 00 00 01 dd c7 8d
	// 8b 29 8a 9c 51 9e 36 5c 9f c6 8c 40 e3 17 d0 98
	// e9 09 00 00 00 00 00 00 00 0c 00 00 00 00 00 00
	// 00 01 8e b2 08 f7 e0 5d 98 7a 9b 04 4a 8e 98 c6
	// b0 87 f1 5a 0b fc 22 00 00 00 00 00 00 00 38 00
	// 00 00 00 00 00 00 ................................. Full message
	// ------------------------------------------------------------------------
	// 0x0000 | a2 02 00 00 ....................................... Length
	// 0x0004 | 47 49 56 54 ....................................... Prefix
	// 0x0008 | 02 00 00 00 ....................................... Txns length
	// 0x000c | 01 00 00 00 88 13 00 00 7b 00 00 00 00 00 00 00
	// 0x001c | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x002c | 00 00 00 00 00 00 00 00 00 02 00 00 00 00 00 00
	// 0x003c | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x004c | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x005c | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x006c | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x007c | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x008c | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x009c | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x00ac | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 02
	// 0x00bc | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x00cc | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x00dc | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x00ec | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x00fc | 00 00 00 02 00 00 00 01 e3 43 1a 8e 0a db f9 6f
	// 0x010c | d1 40 10 3d c6 f6 3a 3f 8f a3 43 ab 0c 00 00 00
	// 0x011c | 00 00 00 00 22 00 00 00 00 00 00 00 01 08 4f 32
	// 0x012c | 85 28 6c a7 84 bb ff d8 23 de 73 62 76 03 64 c9
	// 0x013c | a8 38 00 00 00 00 00 00 00 4e 00 00 00 00 00 00
	// 0x014c | 00 ................................................ Txns#0
	// 0x0151 | 01 00 00 00 88 13 00 00 7b 00 00 00 00 00 00 00
	// 0x0161 | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x0171 | 00 00 00 00 00 00 00 00 00 02 00 00 00 00 00 00
	// 0x0181 | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x0191 | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x01a1 | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x01b1 | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x01c1 | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x01d1 | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x01e1 | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x01f1 | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 03
	// 0x0201 | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x0211 | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x0221 | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x0231 | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x0241 | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x0251 | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x0261 | 00 00 00 02 00 00 00 01 dd c7 8d 8b 29 8a 9c 51
	// 0x0271 | 9e 36 5c 9f c6 8c 40 e3 17 d0 98 e9 09 00 00 00
	// 0x0281 | 00 00 00 00 0c 00 00 00 00 00 00 00 01 8e b2 08
	// 0x0291 | f7 e0 5d 98 7a 9b 04 4a 8e 98 c6 b0 87 f1 5a 0b
	// 0x02a1 | fc 22 00 00 00 00 00 00 00 38 00 00 00 00 00 00
	// 0x02b1 | 00 ................................................ Txns#1
	// 0x02a6 |
}

func ExampleNewAnnounceTxnsMessage() {
	defer gnet.EraseMessages()
	setupMsgEncoding()

	var message = NewAnnounceTxnsMessage([]cipher.SHA256{hashes[8], hashes[9]})
	fmt.Println("AnnounceTxnsMessage:")
	MessageHexDump(message, true)
	// Output:
	// AnnounceTxnsMessage:
	// 48 00 00 00 41 4e 4e 54 02 00 00 00 83 d5 44 cc
	// c2 23 c0 57 d2 bf 80 d3 f2 a3 29 82 c3 2c 3c 0d
	// b8 e2 67 48 20 da 50 64 78 3f b0 97 83 d5 44 cc
	// c2 23 c0 57 d2 bf 80 d3 f2 a3 29 82 c3 2c 3c 0d
	// b8 e2 67 48 20 da 50 64 78 3f b0 97 ............... Full message
	// ------------------------------------------------------------------------
	// 0x0000 | 48 00 00 00 ....................................... Length
	// 0x0004 | 41 4e 4e 54 ....................................... Prefix
	// 0x0008 | 02 00 00 00 ....................................... Txns length
	// 0x000c | 01 00 00 00 83 d5 44 cc c2 23 c0 57 d2 bf 80 d3
	// 0x001c | f2 a3 29 82 c3 2c 3c 0d b8 e2 67 48 20 da 50 64
	// 0x002c | 78 3f b0 97 ....................................... Txns#0
	// 0x0034 | 01 00 00 00 83 d5 44 cc c2 23 c0 57 d2 bf 80 d3
	// 0x0044 | f2 a3 29 82 c3 2c 3c 0d b8 e2 67 48 20 da 50 64
	// 0x0054 | 78 3f b0 97 ....................................... Txns#1
	// 0x004c |
}

func TestSucceed(t *testing.T) {
	// Succeed
}
