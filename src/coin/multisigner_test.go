package coin

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

func fourSigners() ([]cipher.PubKey, []cipher.SecKey) {
	const n = 4
	pks := make([]cipher.PubKey, n)
	sks := make([]cipher.SecKey, n)
	for i := range pks {
		pks[i], sks[i] = cipher.GenerateKeyPair()
	}
	return pks, sks
}

func multiSignedBlock(t *testing.T, seq uint64, sks []cipher.SecKey, proposer int, cosigners ...int) SignedBlock {
	b := Block{Head: BlockHeader{BkSeq: seq, Time: 100 + seq}}
	h := b.HashHeader()
	sb := SignedBlock{Block: b, Sig: cipher.MustSignHash(h, sks[proposer])}
	for _, i := range cosigners {
		c, err := SignBlockCosig(h, i, sks[i])
		require.NoError(t, err)
		sb.Cosigs = append(sb.Cosigs, c)
	}
	return sb
}

func TestMultiSignerRuleApplies(t *testing.T) {
	require.False(t, MultiSignerRule{ForkSeq: 1}.Applies(5), "no signers, no rule")

	pks, _ := fourSigners()
	r := MultiSignerRule{Signers: pks, Threshold: 3, ForkSeq: 10}
	require.False(t, r.Applies(0))
	require.False(t, r.Applies(9))
	require.True(t, r.Applies(10))
	require.True(t, r.Applies(11))

	r.ForkSeq = 0
	require.False(t, r.Applies(0), "genesis is exempt")
	require.True(t, r.Applies(1))
}

func TestMultiSignerRuleValidate(t *testing.T) {
	require.NoError(t, MultiSignerRule{}.Validate())

	pks, _ := fourSigners()
	require.Equal(t, 3, DefaultBlockSignerThreshold(4))
	require.Equal(t, 1, DefaultBlockSignerThreshold(1))
	require.Equal(t, 3, DefaultBlockSignerThreshold(3))
	require.Equal(t, 5, DefaultBlockSignerThreshold(7))

	require.NoError(t, MultiSignerRule{Signers: pks, Threshold: 3}.Validate())
	require.NoError(t, MultiSignerRule{Signers: pks, Threshold: 4}.Validate())
	require.Error(t, MultiSignerRule{Signers: pks, Threshold: 2}.Validate(), "2 of 4 is not more than 2/3")
	require.Error(t, MultiSignerRule{Signers: pks, Threshold: 5}.Validate(), "more than N")
	require.Error(t, MultiSignerRule{Signers: pks, Threshold: 0}.Validate())

	dup := append([]cipher.PubKey{}, pks...)
	dup[3] = dup[0]
	require.Error(t, MultiSignerRule{Signers: dup, Threshold: 3}.Validate(), "duplicate key")

	bad := append([]cipher.PubKey{}, pks...)
	bad[1] = cipher.PubKey{}
	require.Error(t, MultiSignerRule{Signers: bad, Threshold: 3}.Validate(), "invalid key")
}

func TestMultiSignerRuleProposer(t *testing.T) {
	pks, _ := fourSigners()
	r := MultiSignerRule{Signers: pks, Threshold: 3}
	require.Equal(t, 1, r.Proposer(1, 0))
	require.Equal(t, 2, r.Proposer(1, 1))
	require.Equal(t, 0, r.Proposer(3, 1))
	require.Equal(t, 3, r.Proposer(7, 4))
}

func TestMultiSignerRuleVerify(t *testing.T) {
	pks, sks := fourSigners()
	_, outSK := cipher.GenerateKeyPair()
	r := MultiSignerRule{Signers: pks, Threshold: 3, ForkSeq: 1}

	// Proposer plus two cosigners make three
	b := multiSignedBlock(t, 5, sks, 1, 2, 3)
	valid, err := r.Verify(b)
	require.NoError(t, err)
	require.Len(t, valid, 2)
	require.Equal(t, []cipher.PubKey{pks[1], pks[2], pks[3]}, r.BlockSigners(b))

	// All four
	b = multiSignedBlock(t, 5, sks, 0, 1, 2, 3)
	_, err = r.Verify(b)
	require.NoError(t, err)

	// k-1 signatures
	b = multiSignedBlock(t, 5, sks, 1, 2)
	_, err = r.Verify(b)
	require.Equal(t, ErrNotEnoughBlockSigners{Have: 2, Need: 3}, err)

	// The proposer's cosignature does not count twice
	b = multiSignedBlock(t, 5, sks, 1, 1, 2)
	_, err = r.Verify(b)
	require.Equal(t, ErrNotEnoughBlockSigners{Have: 2, Need: 3}, err)

	// A duplicate cosigner counts once
	b = multiSignedBlock(t, 5, sks, 1, 2, 2)
	valid, err = r.Verify(b)
	require.Equal(t, ErrNotEnoughBlockSigners{Have: 2, Need: 3}, err)
	require.Nil(t, valid)

	// A proposer outside the list is rejected however many cosigners it has
	sks5 := append(append([]cipher.SecKey{}, sks...), outSK)
	b = multiSignedBlock(t, 5, sks5, 4, 0, 1, 2, 3)
	_, err = r.Verify(b)
	require.Equal(t, ErrBlockProposerNotSigner, err)

	// A cosignature by a key outside the list does not count, under any index
	b = multiSignedBlock(t, 5, sks, 1, 2)
	h := b.HashHeader()
	for _, i := range []uint32{3, 4, 1000} {
		b.Cosigs = append(b.Cosigs, BlockCosig{Signer: i, Sig: cipher.MustSignHash(h, outSK)})
	}
	_, err = r.Verify(b)
	require.Equal(t, ErrNotEnoughBlockSigners{Have: 2, Need: 3}, err)

	// A cosignature over another block does not count
	b = multiSignedBlock(t, 5, sks, 1, 2)
	other := multiSignedBlock(t, 6, sks, 3, 3)
	b.Cosigs = append(b.Cosigs, other.Cosigs...)
	_, err = r.Verify(b)
	require.Equal(t, ErrNotEnoughBlockSigners{Have: 2, Need: 3}, err)

	// Bad cosignatures are dropped from the result, not counted against the block
	b = multiSignedBlock(t, 5, sks, 1, 2, 3)
	b.Cosigs = append(b.Cosigs, BlockCosig{Signer: 0, Sig: cipher.MustSignHash(h, outSK)})
	valid, err = r.Verify(b)
	require.NoError(t, err)
	require.Len(t, valid, 2)
}

func TestMultiSignerRuleForkSeq(t *testing.T) {
	pks, sks := fourSigners()
	r := MultiSignerRule{Signers: pks, Threshold: 3, ForkSeq: 10}

	// Before the fork the rule does not apply and the block is not checked here
	b := multiSignedBlock(t, 9, sks, 1)
	_, err := r.Verify(b)
	require.NoError(t, err)
	require.Nil(t, r.BlockSigners(b))

	b = multiSignedBlock(t, 10, sks, 1)
	_, err = r.Verify(b)
	require.Error(t, err)

	// Off when there are no signers
	_, err = MultiSignerRule{}.Verify(multiSignedBlock(t, 10, sks, 1))
	require.NoError(t, err)
}

func TestBlockCosigsEncoding(t *testing.T) {
	_, sks := fourSigners()
	b := multiSignedBlock(t, 5, sks, 0, 1, 2, 3)

	buf := EncodeBlockCosigs(b.Cosigs)
	require.Len(t, buf, 4+3*(4+65))
	cs, n, err := DecodeBlockCosigs(buf)
	require.NoError(t, err)
	require.Equal(t, len(buf), n)
	require.Equal(t, b.Cosigs, cs)

	empty := EncodeBlockCosigs(nil)
	cs, n, err = DecodeBlockCosigs(empty)
	require.NoError(t, err)
	require.Equal(t, 4, n)
	require.Empty(t, cs)

	// Trailing data is left unread
	cs, n, err = DecodeBlockCosigs(append(buf, 1, 2, 3))
	require.NoError(t, err)
	require.Equal(t, len(buf), n)
	require.Len(t, cs, 3)

	for _, bad := range [][]byte{nil, {1, 0}, buf[:len(buf)-1], {0xff, 0xff, 0xff, 0xff}} {
		_, _, err := DecodeBlockCosigs(bad)
		require.Equal(t, ErrBlockCosigsMalformed, err)
	}
}

func TestSignedBlockEncodingIgnoresCosigs(t *testing.T) {
	_, sks := fourSigners()
	b := multiSignedBlock(t, 5, sks, 0, 1, 2)
	plain := b
	plain.Cosigs = nil
	require.Equal(t, encoder.Serialize(plain), encoder.Serialize(b))
}
