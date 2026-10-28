package daemon

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/daemon/gnet"
	"github.com/skycoin/skycoin/src/testutil"
)

func testCosig(t *testing.T, hash cipher.SHA256, i int) coin.BlockCosig {
	_, sk := cipher.GenerateKeyPair()
	c, err := coin.SignBlockCosig(hash, i, sk)
	require.NoError(t, err)
	return c
}

func TestGiveBlockCosigsMessageEncoding(t *testing.T) {
	h1 := testutil.RandSHA256(t)
	h2 := testutil.RandSHA256(t)
	m := &GiveBlockCosigsMessage{Cosigs: []BlockCosigEntry{
		{Hash: h1, Cosig: testCosig(t, h1, 0)},
		{Hash: h1, Cosig: testCosig(t, h1, 3)},
		{Hash: h2, Cosig: testCosig(t, h2, 1)},
	}}
	buf := make([]byte, m.EncodeSize())
	require.NoError(t, m.Encode(buf))

	var got GiveBlockCosigsMessage
	n, err := got.Decode(buf)
	require.NoError(t, err)
	require.Equal(t, uint64(len(buf)), n)
	require.Equal(t, m.Cosigs, got.Cosigs)

	_, err = got.Decode(buf[:len(buf)-1])
	require.Equal(t, encoder.ErrBufferUnderflow, err)

	tooMany := &GiveBlockCosigsMessage{Cosigs: make([]BlockCosigEntry, maxBlockCosigsPerMessage+1)}
	require.Equal(t, encoder.ErrMaxLenExceeded, tooMany.Encode(make([]byte, tooMany.EncodeSize())))

	big := make([]byte, 4)
	big[0], big[1] = 0xff, 0xff
	_, err = got.Decode(big)
	require.Equal(t, encoder.ErrMaxLenExceeded, err)
}

func TestProposeBlockMessageEncoding(t *testing.T) {
	pk, sk := cipher.GenerateKeyPair()
	addr := cipher.AddressFromPubKey(pk)
	var txn coin.Transaction
	require.NoError(t, txn.PushOutput(addr, 1e6, 100))
	txn.InnerHash = txn.HashInner()

	b := coin.SignedBlock{Block: coin.Block{
		Head: coin.BlockHeader{BkSeq: 12, Time: 1000, Fee: 3},
		Body: coin.BlockBody{Transactions: coin.Transactions{txn}},
	}}
	b.Sig = cipher.MustSignHash(b.HashHeader(), sk)

	m := &ProposeBlockMessage{Round: 7, Block: b}
	buf := make([]byte, m.EncodeSize())
	require.NoError(t, m.Encode(buf))

	var got ProposeBlockMessage
	n, err := got.Decode(buf)
	require.NoError(t, err)
	require.Equal(t, uint64(len(buf)), n)
	require.Equal(t, uint64(7), got.Round)
	require.Equal(t, b.HashHeader(), got.Block.HashHeader())
	require.Equal(t, b.Sig, got.Block.Sig)
	require.Equal(t, b.Body.Transactions, got.Block.Body.Transactions)

	_, err = got.Decode(buf[:7])
	require.Error(t, err)
	_, err = got.Decode(buf[:len(buf)-1])
	require.Error(t, err)
}

func TestNewGiveBlockCosigsMessages(t *testing.T) {
	require.Empty(t, newGiveBlockCosigsMessages(nil))

	blocks := make([]coin.SignedBlock, 100)
	for i := range blocks {
		blocks[i].Head.BkSeq = uint64(i) //nolint:gosec
		h := blocks[i].HashHeader()
		blocks[i].Cosigs = []coin.BlockCosig{testCosig(t, h, 0), testCosig(t, h, 1), testCosig(t, h, 2)}
	}
	entries := blockCosigEntries(blocks)
	require.Len(t, entries, 300)
	require.Equal(t, blocks[5].HashHeader(), entries[15].Hash)

	msgs := newGiveBlockCosigsMessages(entries)
	require.Len(t, msgs, 2)
	require.Len(t, msgs[0].Cosigs, maxBlockCosigsPerMessage)
	require.Len(t, msgs[1].Cosigs, 300-maxBlockCosigsPerMessage)
	for _, m := range msgs {
		require.LessOrEqual(t, m.EncodeSize()+4, uint64(256*1024), "fits the default MaxOutgoingMessageLength")
	}
}

func TestMultiSignerStateCosigs(t *testing.T) {
	s := newMultiSignerState()
	first := testutil.RandSHA256(t)
	c := testCosig(t, first, 2)
	require.True(t, s.addCosig(first, c))
	require.False(t, s.addCosig(first, c), "a known cosignature is not new")
	require.False(t, s.addCosig(first, coin.BlockCosig{Signer: 2}), "one per signer")
	require.Equal(t, []coin.BlockCosig{c}, s.cosigsFor(first))

	for i := 0; i < blockCosigCacheSize; i++ {
		s.addCosig(testutil.RandSHA256(t), coin.BlockCosig{})
	}
	require.Nil(t, s.cosigsFor(first), "the oldest block goes once the cache is full")
	require.Len(t, s.cosigs, blockCosigCacheSize)
}

func TestMultiSignerStateProposals(t *testing.T) {
	s := newMultiSignerState()
	s.setHeight(5)

	var b coin.SignedBlock
	b.Head.BkSeq = 5
	require.True(t, s.addProposal(b, 0))
	require.False(t, s.addProposal(b, 1), "already held")

	var other coin.SignedBlock
	other.Head.BkSeq = 6
	require.False(t, s.addProposal(other, 0), "not the next block")

	for i := 1; i < maxProposalsPerHeight+5; i++ {
		b.Head.Time = uint64(i) //nolint:gosec
		s.addProposal(b, 0)
	}
	require.Len(t, s.proposals, maxProposalsPerHeight)

	s.setHeight(6)
	require.Empty(t, s.proposals, "a new head drops the old proposals")
	require.True(t, s.addProposal(other, 0))
}

// A block that arrives with too few signatures is held for them, not dropped.
func TestGiveBlocksMessageHoldsUnderSignedBlock(t *testing.T) {
	var b coin.SignedBlock
	b.Head.BkSeq = 8
	h := b.HashHeader()
	cosigs := []coin.BlockCosig{testCosig(t, h, 1)}

	d := &mockDaemoner{}
	d.On("DaemonConfig").Return(DaemonConfig{})
	d.On("headBkSeq").Return(uint64(7), true, nil)
	d.On("pqSigFor", h).Return(nil)
	d.On("cosigsFor", h).Return(cosigs)
	d.On("executeSignedBlock", mock.MatchedBy(func(got coin.SignedBlock) bool {
		return len(got.Cosigs) == 1
	})).Return(coin.ErrNotEnoughBlockSigners{Have: 2, Need: 3})
	d.On("holdBlock", mock.Anything).Return()

	m := &GiveBlocksMessage{Blocks: []coin.SignedBlock{b}, c: &gnet.MessageContext{Addr: "127.0.0.1:1234"}}
	m.process(d)
	d.AssertCalled(t, "holdBlock", mock.Anything)
}

func TestMultiSignerMessagesRegistered(t *testing.T) {
	prefixes := make(map[gnet.MessagePrefix]bool)
	for _, mc := range getMessageConfigs() {
		require.False(t, prefixes[mc.Prefix], "prefix %s is used twice", mc.Prefix)
		prefixes[mc.Prefix] = true
	}
	require.True(t, prefixes[gnet.MessagePrefixFromString("PRPB")])
	require.True(t, prefixes[gnet.MessagePrefixFromString("GIVC")])
	require.Equal(t, multiSignerProtocolVersion, NewDaemonConfig().ProtocolVersion)
}
