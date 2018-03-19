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
	"github.com/skycoin/skycoin/src/testutil"
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
}

func ExampleNewGetPeersMessage() {
	defer gnet.EraseMessages()
	setupMsgEncoding()

	var message = NewGetPeersMessage()
	fmt.Println("GetPeersMessage:")
	MessageHexDump(message, true)
	// Output:
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
}

func ExampleNewPingMessage() {
	//var message gnet.Message = PingMessage{
	//}
	//MessageHexDump(message, true)
}

func ExampleNewPongMessage() {
	//var message = PongMessage{
	//}
	//MessageHexDump(message, true)
}

func ExampleNewGetBlocksMessage() {
	defer gnet.EraseMessages()
	setupMsgEncoding()

	var message = NewGetBlocksMessage(1234, 5678)
	fmt.Println("GetBlocksMessage:")
	MessageHexDump(message, true)
	// Output:
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
			PrevHash: cipher.SHA256{},
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
}

func ExampleNewAnnounceBlocksMessage() {
	defer gnet.EraseMessages()
	setupMsgEncoding()

	var message = NewAnnounceBlocksMessage(123456)
	fmt.Println("AnnounceBlocksMessage:")
	MessageHexDump(message, true)
	// Output:
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
}

func ExampleNewGiveTxnsMessage() {
	defer gnet.EraseMessages()
	setupMsgEncoding()

	var transactions coin.Transactions = make([]coin.Transaction, 0)
	var transactionOutputs0 []coin.TransactionOutput = make([]coin.TransactionOutput, 0)
	var transactionOutputs1 []coin.TransactionOutput = make([]coin.TransactionOutput, 0)
	var txOutput0 = coin.TransactionOutput{
		Address: testutil.MakeAddress(),
		Coins:   12,
		Hours:   34,
	}
	var txOutput1 = coin.TransactionOutput{
		Address: testutil.MakeAddress(),
		Coins:   56,
		Hours:   78,
	}
	var txOutput2 = coin.TransactionOutput{
		Address: testutil.MakeAddress(),
		Coins:   9,
		Hours:   12,
	}
	var txOutput3 = coin.TransactionOutput{
		Address: testutil.MakeAddress(),
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
