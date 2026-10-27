package daemon

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher/encoder"
	"github.com/skycoin/skycoin/src/cipher/pqsig"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/daemon/gnet"
	"github.com/skycoin/skycoin/src/testutil"
)

func pqTestSig(b byte) []byte {
	sig := make([]byte, pqsig.SigSize)
	for i := range sig {
		sig[i] = b
	}
	return sig
}

func TestGivePQSigsMessageEncoding(t *testing.T) {
	m := &GivePQSigsMessage{Sigs: []BlockPQSig{
		{Hash: testutil.RandSHA256(t), Sig: pqTestSig(1)},
		{Hash: testutil.RandSHA256(t), Sig: pqTestSig(2)},
	}}
	buf := make([]byte, m.EncodeSize())
	require.NoError(t, m.Encode(buf))

	var got GivePQSigsMessage
	n, err := got.Decode(buf)
	require.NoError(t, err)
	require.Equal(t, uint64(len(buf)), n)
	require.Equal(t, m.Sigs, got.Sigs)

	_, err = got.Decode(buf[:len(buf)-1])
	require.Equal(t, encoder.ErrBufferUnderflow, err)

	tooMany := &GivePQSigsMessage{Sigs: make([]BlockPQSig, maxPQSigsPerMessage+1)}
	require.Equal(t, encoder.ErrMaxLenExceeded, tooMany.Encode(make([]byte, tooMany.EncodeSize())))

	tooLong := &GivePQSigsMessage{Sigs: []BlockPQSig{{Sig: make([]byte, pqsig.SigSize+1)}}}
	require.Equal(t, encoder.ErrMaxLenExceeded, tooLong.Encode(make([]byte, tooLong.EncodeSize())))
}

func TestNewGivePQSigsMessages(t *testing.T) {
	blocks := make([]coin.SignedBlock, maxPQSigsPerMessage+5)
	for i := range blocks {
		blocks[i].Head.BkSeq = uint64(i) //nolint:gosec
		if i != 3 {
			blocks[i].PQSig = pqTestSig(byte(i))
		}
	}
	msgs := newGivePQSigsMessages(blocks)
	require.Len(t, msgs, 2)
	require.Len(t, msgs[0].Sigs, maxPQSigsPerMessage)
	require.Len(t, msgs[1].Sigs, 4, "the block with no signature is skipped")
	require.Equal(t, blocks[0].HashHeader(), msgs[0].Sigs[0].Hash)

	for _, m := range msgs {
		require.LessOrEqual(t, m.EncodeSize()+4, uint64(256*1024), "fits the default MaxOutgoingMessageLength")
	}
}

func TestPQSigCache(t *testing.T) {
	c := newPQSigCache()
	first := testutil.RandSHA256(t)
	c.add([]BlockPQSig{{Hash: first, Sig: pqTestSig(1)}})
	require.Equal(t, pqTestSig(1), c.get(first))

	for i := 0; i < pqSigCacheSize; i++ {
		c.add([]BlockPQSig{{Hash: testutil.RandSHA256(t), Sig: pqTestSig(2)}})
	}
	require.Nil(t, c.get(first), "the oldest entry goes once the cache is full")
	require.Len(t, c.sigs, pqSigCacheSize)
}

func TestGiveBlocksMessageAttachesPQSig(t *testing.T) {
	var b coin.SignedBlock
	b.Head.BkSeq = 8
	sig := pqTestSig(7)

	d := &mockDaemoner{}
	d.On("DaemonConfig").Return(DaemonConfig{})
	d.On("headBkSeq").Return(uint64(7), true, nil)
	d.On("pqSigFor", b.HashHeader()).Return(sig)
	d.On("cosigsFor", b.HashHeader()).Return(nil)
	d.On("executeSignedBlock", mock.MatchedBy(func(got coin.SignedBlock) bool {
		return string(got.PQSig) == string(sig)
	})).Return(nil)
	d.On("requestBlocksFromAddr", mock.Anything).Return(nil).Maybe()
	d.On("broadcastMessage", mock.Anything).Return(nil, nil).Maybe()
	d.On("announceAllValidTxns").Return(nil).Maybe()

	m := &GiveBlocksMessage{Blocks: []coin.SignedBlock{b}, c: &gnet.MessageContext{Addr: "127.0.0.1:1234"}}
	m.process(d)
	d.AssertCalled(t, "executeSignedBlock", mock.Anything)
}
