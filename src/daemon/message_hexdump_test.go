package daemon

import (
	"testing"

	"fmt"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/daemon/pex"
	"github.com/skycoin/skycoin/src/testutil"
)

func TestIntroductionMessage(t *testing.T) {
	var message = NewIntroductionMessage(1234, 5, 7890)
	fmt.Println("IntroductionMessage:")
	fmt.Println(HexDump(message))
}

func TestGetPeersMessage(t *testing.T) {
	var message = NewGetPeersMessage()
	fmt.Println("GetPeersMessage:")
	fmt.Println(HexDump(message))
}

func TestGivePeersMessage(t *testing.T) {
	var peers = make([]pex.Peer, 3)
	var peer0 pex.Peer = *pex.NewPeer("118.178.135.93:6000")
	var peer1 pex.Peer = *pex.NewPeer("47.88.33.156:6000")
	var peer2 pex.Peer = *pex.NewPeer("121.41.103.148:6000")
	peers = append(peers, peer0, peer1, peer2)
	var message = NewGivePeersMessage(peers)
	fmt.Println("GivePeersMessage:")
	fmt.Println(HexDump(message))
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

func TestGetBlocksMessage(t *testing.T) {
	var message = NewGetBlocksMessage(1234, 5678)
	fmt.Println("GetBlocksMessage:")
	fmt.Println(HexDump(message))
}

func TestGiveBlocksMessage(t *testing.T) {
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
}

func TestAnnounceBlocksMessage(t *testing.T) {
	var message = NewAnnounceBlocksMessage(123456)
	fmt.Println("AnnounceBlocksMessage:")
	fmt.Println(HexDump(message))
}

func TestGetTxnsMessage(t *testing.T) {
	var shas = make([]cipher.SHA256, 0)

	shas = append(shas, testutil.RandSHA256(t), testutil.RandSHA256(t))
	var message = NewGetTxnsMessage(shas)
	fmt.Println("GetTxns:")
	fmt.Println(HexDump(message))
}

func TestGiveTxnsMessage(t *testing.T) {
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
		In:        []cipher.SHA256{testutil.RandSHA256(t), testutil.RandSHA256(t)},
		InnerHash: testutil.RandSHA256(t),
		Length:    5000,
		Out:       transactionOutputs0,
		Sigs:      []cipher.Sig{sig0, sig1},
	}
	var transaction1 = coin.Transaction{
		Type:      123,
		In:        []cipher.SHA256{testutil.RandSHA256(t), testutil.RandSHA256(t)},
		InnerHash: testutil.RandSHA256(t),
		Length:    5000,
		Out:       transactionOutputs1,
		Sigs:      []cipher.Sig{sig2, sig3},
	}
	transactions = append(transactions, transaction0, transaction1)
	var message = NewGiveTxnsMessage(transactions)
	fmt.Println("GiveTxnsMessage:")
	fmt.Println(HexDump(message))
}

func TestAnnounceTxnsMessage(t *testing.T) {
	var message = NewAnnounceTxnsMessage([]cipher.SHA256{testutil.RandSHA256(t), testutil.RandSHA256(t)})
	fmt.Println("AnnounceTxnsMessage:")
	fmt.Println(HexDump(message))
}
