// +build ignore

package coin

import (
	"testing"

	"github.com/ShanghaiKuaibei/mzcoin/src/coin"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/stretchr/testify/assert"
)

func TestNewBlock(t *testing.T) {
	// TODO -- update this test for newBlock changes
	prev := coin.Block{Head: coin.BlockHeader{Version: 0x02, Time: 100, BkSeq: 98}}
	unsp := coin.NewUnspentPool()
	unsp.XorHash = randSHA256()
	txns := coin.Transactions{coin.Transaction{}}
	// invalid txn fees panics
	assert.Panics(t, func() { coin.NewBlock(prev, 133, unsp, txns, _badFeeCalc) })
	// no txns panics
	assert.Panics(t, func() {
		coin.NewBlock(prev, 133, unsp, nil, _feeCalc)
	})
	assert.Panics(t, func() {
		coin.NewBlock(prev, 133, unsp, coin.Transactions{}, _feeCalc)
	})
	// valid block is fine
	fee := uint64(121)
	currentTime := uint64(133)
	b := coin.NewBlock(prev, currentTime, unsp, txns, _makeFeeCalc(fee))
	assert.Equal(t, b.Body.Transactions, txns)
	assert.Equal(t, b.Head.Fee, fee*uint64(len(txns)))
	assert.Equal(t, b.Body, coin.BlockBody{Transactions: txns})
	assert.Equal(t, b.Head.PrevHash, prev.HashHeader())
	assert.Equal(t, b.Head.Time, currentTime)
	assert.Equal(t, b.Head.BkSeq, prev.Head.BkSeq+1)
	assert.Equal(t, b.Head.UxHash,
		unsp.GetUxHash())
}

func TestBlockHashHeader(t *testing.T) {
	b := makeNewBlock()
	assert.Equal(t, b.HashHeader(), b.Head.Hash())
	assert.NotEqual(t, b.HashHeader(), cipher.SHA256{})
}

func TestBlockHashBody(t *testing.T) {
	b := makeNewBlock()
	assert.Equal(t, b.HashBody(), b.Body.Hash())
	hb := b.HashBody()
	hashes := b.Body.Transactions.Hashes()
	tx := addTransactionToBlock(t, &b)
	assert.NotEqual(t, b.HashBody(), hb)
	hashes = append(hashes, tx.Hash())
	assert.Equal(t, b.HashBody(), cipher.Merkle(hashes))
	assert.Equal(t, b.HashBody(), b.Body.Hash())
}
