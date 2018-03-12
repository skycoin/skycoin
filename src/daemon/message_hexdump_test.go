package daemon

import (
	"testing"

	"fmt"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/daemon/pex"
	"github.com/skycoin/skycoin/src/testutil"
)

func ExampleIntroductionMessage() {
	var message = NewIntroductionMessage(1234, 5, 7890)
	fmt.Println("IntroductionMessage:")
	fmt.Println(HexDump(message))
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

func ExampleGetPeersMessage() {
	var message = NewGetPeersMessage()
	fmt.Println("GetPeersMessage:")
	fmt.Println(HexDump(message))
	// Output:
	// GetPeersMessage:
	// 04 00 00 00 47 45 54 50 ........................... Full message
	// ------------------------------------------------------------------------
	// 0x0000 | 04 00 00 00 ....................................... Length
	// 0x0004 | 47 45 54 50 ....................................... Prefix
	// 0x0008 |
}

func ExampleGivePeersMessage() {
	var peers = make([]pex.Peer, 3)
	var peer0 pex.Peer = *pex.NewPeer("118.178.135.93:6000")
	var peer1 pex.Peer = *pex.NewPeer("47.88.33.156:6000")
	var peer2 pex.Peer = *pex.NewPeer("121.41.103.148:6000")
	peers = append(peers, peer0, peer1, peer2)
	var message = NewGivePeersMessage(peers)
	fmt.Println("GivePeersMessage:")
	fmt.Println(HexDump(message))
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

func TestPingMessage(t *testing.T) {
	//var message gnet.Message = daemon.PingMessage{
	//}
	//HexDump(message)
}

func TestPongMessage(t *testing.T) {
	//var message = daemon.PongMessage{
	//}
	//HexDump(message)
}

func ExampleGetBlocksMessage() {
	var message = NewGetBlocksMessage(1234, 5678)
	fmt.Println("GetBlocksMessage:")
	fmt.Println(HexDump(message))
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

func ExampleGiveBlocksMessage() {
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
	fmt.Println(HexDump(message))
	// Output:
	// GiveBlocksMessage:
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

func ExampleAnnounceBlocksMessage() {
	var message = NewAnnounceBlocksMessage(123456)
	fmt.Println("AnnounceBlocksMessage:")
	fmt.Println(HexDump(message))
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

func ExampleGetTxnsMessage() { //TODO: Conflict Here
	var shas = make([]cipher.SHA256, 0)

	shas = append(shas, GenerateRandomSha256(), GenerateRandomSha256())
	var message = NewGetTxnsMessage(shas)
	fmt.Println("GetTxnsMessage:")
	fmt.Println(HexDump(message))
	// Output:
	// GetTxnsMessage:
	// 48 00 00 00 47 45 54 54 02 00 00 00 81 53 df 72
	// 1b 46 7a a3 94 2d 8c a4 e9 bc e5 13 45 ef ca 14
	// 12 bd cd 65 ab d8 a3 7c a7 9e d6 3a 07 ac 56 ca
	// 77 a9 4b cf 88 54 c2 90 d1 8d 66 ef 4e a3 62 cd
	// 90 fb e4 5b d9 ef f4 e5 38 1f db 73 ............... Full message
	// ------------------------------------------------------------------------
	// 0x0000 | 48 00 00 00 ....................................... Length
	// 0x0004 | 47 45 54 54 ....................................... Prefix
	// 0x0008 | 02 00 00 00 ....................................... Txns length
	// 0x000c | 01 00 00 00 81 53 df 72 1b 46 7a a3 94 2d 8c a4
	// 0x001c | e9 bc e5 13 45 ef ca 14 12 bd cd 65 ab d8 a3 7c
	// 0x002c | a7 9e d6 3a ....................................... Txns#0
	// 0x0034 | 01 00 00 00 07 ac 56 ca 77 a9 4b cf 88 54 c2 90
	// 0x0044 | d1 8d 66 ef 4e a3 62 cd 90 fb e4 5b d9 ef f4 e5
	// 0x0054 | 38 1f db 73 ....................................... Txns#1
	// 0x004c |
}

func ExampleGiveTxnsMessage() {
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
		In:        []cipher.SHA256{GenerateRandomSha256(), GenerateRandomSha256()},
		InnerHash: GenerateRandomSha256(),
		Length:    5000,
		Out:       transactionOutputs0,
		Sigs:      []cipher.Sig{sig0, sig1},
	}
	var transaction1 = coin.Transaction{
		Type:      123,
		In:        []cipher.SHA256{GenerateRandomSha256(), GenerateRandomSha256()},
		InnerHash: GenerateRandomSha256(),
		Length:    5000,
		Out:       transactionOutputs1,
		Sigs:      []cipher.Sig{sig2, sig3},
	}
	transactions = append(transactions, transaction0, transaction1)
	var message = NewGiveTxnsMessage(transactions)
	fmt.Println("GiveTxnsMessage:")
	fmt.Println(HexDump(message))
	// Output:
	// GiveTxnsMessage:
	// 82 02 00 00 47 49 56 54 02 00 00 00 88 13 00 00
	// 7b 91 fe 83 9d 92 d1 e4 d1 77 a0 33 1a 90 c9 1e
	// ea 0e 3c 09 0c 73 c6 c4 b3 6e 38 f0 08 be ae 8c
	// 26 02 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 00 00 00 00 00 00 00 02 00 00 00 31 79 64 90 a8
	// d2 0d 56 5e 4b 33 bd 0f e8 66 d6 db c0 4a 6f ce
	// 83 e7 1b de a9 2a 27 6b e7 04 eb 6b 1e e4 fb 9c
	// 0f d8 69 50 de 98 c3 fd 04 4c f3 26 8e 88 a6 f8
	// 25 d7 29 1d 25 ee a7 5d f8 bb 01 02 00 00 00 00
	// d5 9b 21 22 d8 f0 69 99 eb 27 86 59 9f b1 12 08
	// 02 84 cc 66 0c 00 00 00 00 00 00 00 22 00 00 00
	// 00 00 00 00 00 56 04 df 45 c3 0f f0 29 bb bd dd
	// b2 53 76 a2 26 f9 1f ca d4 38 00 00 00 00 00 00
	// 00 4e 00 00 00 00 00 00 00 88 13 00 00 7b e5 d2
	// a9 03 5a 69 d2 34 23 33 28 47 f3 66 3c b8 6e 52
	// c8 7a 20 3f 65 1a c3 25 c9 1a 00 40 cc cc 02 00
	// 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 00 00 00 00 02 00 00 00 01 9b e9 28 40 5e fa 49
	// 42 f8 ca 89 52 31 f4 0d 5a 73 62 d7 c7 c1 29 1c
	// e8 44 ea 2d 59 04 c5 7f a2 13 2e f3 5e f3 c0 00
	// 27 e8 73 e5 d2 dc f8 a0 8b 07 8e d3 68 8e fb da
	// 54 de 73 51 4b 8f 9a 68 02 00 00 00 00 4e 55 f7
	// 8c 22 d3 8a f3 73 a9 93 14 21 1b 55 96 99 c6 2f
	// e6 09 00 00 00 00 00 00 00 0c 00 00 00 00 00 00
	// 00 00 20 28 01 fe f3 4a de 97 6b f2 ea 18 c3 17
	// f9 f9 0a f1 6f 19 22 00 00 00 00 00 00 00 38 00
	// 00 00 00 00 00 00 ................................. Full message
	// ------------------------------------------------------------------------
	// 0x0000 | 82 02 00 00 ....................................... Length
	// 0x0004 | 47 49 56 54 ....................................... Prefix
	// 0x0008 | 02 00 00 00 ....................................... Txns length
	// 0x000c | 01 00 00 00 88 13 00 00 7b 91 fe 83 9d 92 d1 e4
	// 0x001c | d1 77 a0 33 1a 90 c9 1e ea 0e 3c 09 0c 73 c6 c4
	// 0x002c | b3 6e 38 f0 08 be ae 8c 26 02 00 00 00 00 00 00
	// 0x003c | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x004c | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x005c | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x006c | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x007c | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x008c | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x009c | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x00ac | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 02
	// 0x00bc | 00 00 00 31 79 64 90 a8 d2 0d 56 5e 4b 33 bd 0f
	// 0x00cc | e8 66 d6 db c0 4a 6f ce 83 e7 1b de a9 2a 27 6b
	// 0x00dc | e7 04 eb 6b 1e e4 fb 9c 0f d8 69 50 de 98 c3 fd
	// 0x00ec | 04 4c f3 26 8e 88 a6 f8 25 d7 29 1d 25 ee a7 5d
	// 0x00fc | f8 bb 01 02 00 00 00 00 d5 9b 21 22 d8 f0 69 99
	// 0x010c | eb 27 86 59 9f b1 12 08 02 84 cc 66 0c 00 00 00
	// 0x011c | 00 00 00 00 22 00 00 00 00 00 00 00 00 56 04 df
	// 0x012c | 45 c3 0f f0 29 bb bd dd b2 53 76 a2 26 f9 1f ca
	// 0x013c | d4 38 00 00 00 00 00 00 00 4e 00 00 00 00 00 00
	// 0x014c | 00 ................................................ Txns#0
	// 0x0151 | 01 00 00 00 88 13 00 00 7b e5 d2 a9 03 5a 69 d2
	// 0x0161 | 34 23 33 28 47 f3 66 3c b8 6e 52 c8 7a 20 3f 65
	// 0x0171 | 1a c3 25 c9 1a 00 40 cc cc 02 00 00 00 00 00 00
	// 0x0181 | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x0191 | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x01a1 | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x01b1 | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x01c1 | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x01d1 | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x01e1 | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
	// 0x01f1 | 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 02
	// 0x0201 | 00 00 00 01 9b e9 28 40 5e fa 49 42 f8 ca 89 52
	// 0x0211 | 31 f4 0d 5a 73 62 d7 c7 c1 29 1c e8 44 ea 2d 59
	// 0x0221 | 04 c5 7f a2 13 2e f3 5e f3 c0 00 27 e8 73 e5 d2
	// 0x0231 | dc f8 a0 8b 07 8e d3 68 8e fb da 54 de 73 51 4b
	// 0x0241 | 8f 9a 68 02 00 00 00 00 4e 55 f7 8c 22 d3 8a f3
	// 0x0251 | 73 a9 93 14 21 1b 55 96 99 c6 2f e6 09 00 00 00
	// 0x0261 | 00 00 00 00 0c 00 00 00 00 00 00 00 00 20 28 01
	// 0x0271 | fe f3 4a de 97 6b f2 ea 18 c3 17 f9 f9 0a f1 6f
	// 0x0281 | 19 22 00 00 00 00 00 00 00 38 00 00 00 00 00 00
	// 0x0291 | 00 ................................................ Txns#1
	// 0x0286 |
}

func TestAnnounceTxnsMessage(t *testing.T) {
	var message = NewAnnounceTxnsMessage([]cipher.SHA256{testutil.RandSHA256(t), testutil.RandSHA256(t)})
	fmt.Println("AnnounceTxnsMessage:")
	fmt.Println(HexDump(message))
	// Output:
	// AnnounceTxnsMessage:
	// 48 00 00 00 41 4e 4e 54 02 00 00 00 37 07 e8 23
	// 50 93 3d f7 e6 09 dc 41 3d c1 77 3c ae 9d ba af
	// fd fc 56 7a 43 5d 2f e9 e0 fb 63 20 50 4c 6a 6e
	// 6d 17 a4 91 29 42 8e 14 41 df 1f 50 b2 2b 03 75
	// f8 34 99 ce 47 24 59 e5 86 f2 9c 8d ............... Full message
	// ------------------------------------------------------------------------
	// 0x0000 | 48 00 00 00 ....................................... Length
	// 0x0004 | 41 4e 4e 54 ....................................... Prefix
	// 0x0008 | 02 00 00 00 ....................................... Txns length
	// 0x000c | 01 00 00 00 37 07 e8 23 50 93 3d f7 e6 09 dc 41
	// 0x001c | 3d c1 77 3c ae 9d ba af fd fc 56 7a 43 5d 2f e9
	// 0x002c | e0 fb 63 20 ....................................... Txns#0
	// 0x0034 | 01 00 00 00 50 4c 6a 6e 6d 17 a4 91 29 42 8e 14
	// 0x0044 | 41 df 1f 50 b2 2b 03 75 f8 34 99 ce 47 24 59 e5
	// 0x0054 | 86 f2 9c 8d ....................................... Txns#1
	// 0x004c |
}
