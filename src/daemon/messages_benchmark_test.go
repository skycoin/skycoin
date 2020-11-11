package daemon

import (
	"testing"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"github.com/skycoin/skycoin/src/coin"
)

var getPeersMessageObj = GetPeersMessage{}

func BenchmarkDeserializeRawGetPeersMessage(b *testing.B) {
	byt := encoder.Serialize(getPeersMessageObj)
	result := &GetPeersMessage{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		encoder.DeserializeRaw(byt, result) //nolint:errcheck
	}
}

func BenchmarkSerializeGetPeersMessage(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		encoder.Serialize(&getPeersMessageObj)
	}
}

var givePeersMessageObj = GivePeersMessage{
	Peers: []IPAddr{
		{
			IP:   1234,
			Port: 1234,
		},
		{
			IP:   5678,
			Port: 5678,
		},
		{
			IP:   9876,
			Port: 9876,
		},
	},
}

func BenchmarkDeserializeRawGivePeersMessage(b *testing.B) {
	byt := encoder.Serialize(givePeersMessageObj)
	result := &GivePeersMessage{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		encoder.DeserializeRaw(byt, result) //nolint:errcheck
	}
}

func BenchmarkSerializeGivePeersMessage(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		encoder.Serialize(&givePeersMessageObj)
	}
}

var introPubKey = cipher.MustPubKeyFromHex("03cd7dfcd8c3452d1bb5d9d9e34dd95d6848cb9f66c2aad127b60578f4be7498f2")

var introMessageObj = IntroductionMessage{
	Mirror:          1234,
	ListenPort:      5678,
	ProtocolVersion: 1,
	Extra:           introPubKey[:],
}

func BenchmarkDeserializeRawIntroductionMessage(b *testing.B) {
	byt := encoder.Serialize(introMessageObj)
	result := &IntroductionMessage{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		encoder.DeserializeRaw(byt, result) //nolint:errcheck
	}
}

func BenchmarkSerializeIntroductionMessage(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		encoder.Serialize(&introMessageObj)
	}
}

var giveBlocksMessageObj = GiveBlocksMessage{
	Blocks: []coin.SignedBlock{
		{
			Block: coin.Block{
				Body: coin.BlockBody{
					Transactions: []coin.Transaction{
						{
							Sigs: make([]cipher.Sig, 3),
							In:   make([]cipher.SHA256, 3),
							Out:  make([]coin.TransactionOutput, 3),
						},
						{
							Sigs: make([]cipher.Sig, 3),
							In:   make([]cipher.SHA256, 3),
							Out:  make([]coin.TransactionOutput, 3),
						},
						{
							Sigs: make([]cipher.Sig, 3),
							In:   make([]cipher.SHA256, 3),
							Out:  make([]coin.TransactionOutput, 3),
						},
					},
				},
			},
		},
		{
			Block: coin.Block{
				Body: coin.BlockBody{
					Transactions: []coin.Transaction{
						{
							Sigs: make([]cipher.Sig, 3),
							In:   make([]cipher.SHA256, 3),
							Out:  make([]coin.TransactionOutput, 3),
						},
						{
							Sigs: make([]cipher.Sig, 3),
							In:   make([]cipher.SHA256, 3),
							Out:  make([]coin.TransactionOutput, 3),
						},
						{
							Sigs: make([]cipher.Sig, 3),
							In:   make([]cipher.SHA256, 3),
							Out:  make([]coin.TransactionOutput, 3),
						},
					},
				},
			},
		},
		{
			Block: coin.Block{
				Body: coin.BlockBody{
					Transactions: []coin.Transaction{
						{
							Sigs: make([]cipher.Sig, 3),
							In:   make([]cipher.SHA256, 3),
							Out:  make([]coin.TransactionOutput, 3),
						},
						{
							Sigs: make([]cipher.Sig, 3),
							In:   make([]cipher.SHA256, 3),
							Out:  make([]coin.TransactionOutput, 3),
						},
						{
							Sigs: make([]cipher.Sig, 3),
							In:   make([]cipher.SHA256, 3),
							Out:  make([]coin.TransactionOutput, 3),
						},
					},
				},
			},
		},
	},
}

func BenchmarkDeserializeRawGiveBlocksMessage(b *testing.B) {
	byt := encoder.Serialize(giveBlocksMessageObj)
	result := &GiveBlocksMessage{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		encoder.DeserializeRaw(byt, result) //nolint:errcheck
	}
}

func BenchmarkSerializeGiveBlocksMessage(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		encoder.Serialize(&giveBlocksMessageObj)
	}
}

var announceTxnsMessageObj = AnnounceTxnsMessage{
	Transactions: make([]cipher.SHA256, 3),
}

func BenchmarkDeserializeRawAnnounceTxnsMessage(b *testing.B) {
	byt := encoder.Serialize(announceTxnsMessageObj)
	result := &AnnounceTxnsMessage{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		encoder.DeserializeRaw(byt, result) //nolint:errcheck
	}
}

func BenchmarkSerializeAnnounceTxnsMessage(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		encoder.Serialize(&announceTxnsMessageObj)
	}
}

var giveTxnsMessageObj = GiveTxnsMessage{
	Transactions: []coin.Transaction{
		{
			Sigs: make([]cipher.Sig, 3),
			In:   make([]cipher.SHA256, 3),
			Out:  make([]coin.TransactionOutput, 3),
		},
		{
			Sigs: make([]cipher.Sig, 3),
			In:   make([]cipher.SHA256, 3),
			Out:  make([]coin.TransactionOutput, 3),
		},
		{
			Sigs: make([]cipher.Sig, 3),
			In:   make([]cipher.SHA256, 3),
			Out:  make([]coin.TransactionOutput, 3),
		},
	},
}

func BenchmarkDeserializeRawGiveTxnsMessage(b *testing.B) {
	byt := encoder.Serialize(giveTxnsMessageObj)
	result := &GiveTxnsMessage{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		encoder.DeserializeRaw(byt, result) //nolint:errcheck
	}
}

func BenchmarkSerializeGiveTxnsMessage(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		encoder.Serialize(&giveTxnsMessageObj)
	}
}
