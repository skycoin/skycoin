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
		encoder.DeserializeRaw(byt, result) // nolint: errcheck
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
		encoder.DeserializeRaw(byt, result) // nolint: errcheck
	}
}

func BenchmarkSerializeGivePeersMessage(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		encoder.Serialize(&givePeersMessageObj)
	}
}

var introMessageObj = IntroductionMessage{
	Mirror:  1234,
	Port:    5678,
	Version: 1,
	Extra:   []byte("abcdefghijklmnoqrstuvwxyz1234567890"),
}

func BenchmarkDeserializeRawIntroductionMessage(b *testing.B) {
	byt := encoder.Serialize(introMessageObj)
	result := &IntroductionMessage{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		encoder.DeserializeRaw(byt, result) // nolint: errcheck
	}
}

func BenchmarkSerializeIntroductionMessage(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		encoder.Serialize(&introMessageObj)
	}
}

var giveBlocksMessageObj = GiveBlocksMessage{
	Blocks: make([]coin.SignedBlock, 3),
}

func BenchmarkDeserializeRawGiveBlocksMessage(b *testing.B) {
	byt := encoder.Serialize(giveBlocksMessageObj)
	result := &GiveBlocksMessage{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		encoder.DeserializeRaw(byt, result) // nolint: errcheck
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
		encoder.DeserializeRaw(byt, result) // nolint: errcheck
	}
}

func BenchmarkSerializeAnnounceTxnsMessage(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		encoder.Serialize(&announceTxnsMessageObj)
	}
}

var giveTxnsMessageObj = GiveTxnsMessage{
	Transactions: make(coin.Transactions, 3),
}

func BenchmarkDeserializeRawGiveTxnsMessage(b *testing.B) {
	byt := encoder.Serialize(giveTxnsMessageObj)
	result := &GiveTxnsMessage{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		encoder.DeserializeRaw(byt, result) // nolint: errcheck
	}
}

func BenchmarkSerializeGiveTxnsMessage(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		encoder.Serialize(&giveTxnsMessageObj)
	}
}
