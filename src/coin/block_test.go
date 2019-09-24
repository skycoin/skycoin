package coin

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/SkycoinProject/skycoin/src/cipher"
	"github.com/SkycoinProject/skycoin/src/testutil"
)

var (
	genPublic, genSecret        = cipher.GenerateKeyPair()
	genAddress                  = cipher.AddressFromPubKey(genPublic)
	_genTime             uint64 = 1000
	_genCoins            uint64 = 1000e6
	_genCoinHours        uint64 = 1000 * 1000
)

func tNow() uint64 {
	return uint64(time.Now().UTC().Unix())
}

func feeCalc(t *Transaction) (uint64, error) {
	return 0, nil
}

func badFeeCalc(t *Transaction) (uint64, error) {
	return 0, errors.New("Bad")
}

func makeNewBlock(t *testing.T, uxHash cipher.SHA256) *Block {
	body := BlockBody{
		Transactions: Transactions{Transaction{}},
	}

	prev := Block{
		Body: body,
		Head: BlockHeader{
			Version:  0x02,
			Time:     100,
			BkSeq:    0,
			Fee:      10,
			PrevHash: cipher.SHA256{},
			BodyHash: body.Hash(),
		}}
	b, err := NewBlock(prev, 100+20, uxHash, Transactions{Transaction{}}, feeCalc)
	require.NoError(t, err)
	return b
}

func addTransactionToBlock(t *testing.T, b *Block) Transaction {
	txn := makeTransaction(t)
	b.Body.Transactions = append(b.Body.Transactions, txn)
	return txn
}

func TestNewBlock(t *testing.T) {
	// TODO -- update this test for newBlock changes
	prev := Block{Head: BlockHeader{Version: 0x02, Time: 100, BkSeq: 98}}
	uxHash := testutil.RandSHA256(t)
	txns := Transactions{Transaction{}}
	// invalid txn fees panics
	_, err := NewBlock(prev, 133, uxHash, txns, badFeeCalc)
	require.EqualError(t, err, fmt.Sprintf("Invalid transaction fees: Bad"))

	// no txns panics
	_, err = NewBlock(prev, 133, uxHash, nil, feeCalc)
	require.EqualError(t, err, "Refusing to create block with no transactions")

	_, err = NewBlock(prev, 133, uxHash, Transactions{}, feeCalc)
	require.EqualError(t, err, "Refusing to create block with no transactions")

	// valid block is fine
	fee := uint64(121)
	currentTime := uint64(133)
	b, err := NewBlock(prev, currentTime, uxHash, txns, func(t *Transaction) (uint64, error) {
		return fee, nil
	})
	require.NoError(t, err)
	require.Equal(t, b.Body.Transactions, txns)
	require.Equal(t, b.Head.Fee, fee*uint64(len(txns)))
	require.Equal(t, b.Body, BlockBody{Transactions: txns})
	require.Equal(t, b.Head.PrevHash, prev.HashHeader())
	require.Equal(t, b.Head.Time, currentTime)
	require.Equal(t, b.Head.BkSeq, prev.Head.BkSeq+1)
	require.Equal(t, b.Head.UxHash, uxHash)
}

func TestBlockHashHeader(t *testing.T) {
	uxHash := testutil.RandSHA256(t)
	b := makeNewBlock(t, uxHash)
	require.Equal(t, b.HashHeader(), b.Head.Hash())
	require.NotEqual(t, b.HashHeader(), cipher.SHA256{})
}

func TestBlockBodyHash(t *testing.T) {
	uxHash := testutil.RandSHA256(t)
	b := makeNewBlock(t, uxHash)
	hb := b.Body.Hash()
	hashes := b.Body.Transactions.Hashes()
	txn := addTransactionToBlock(t, b)
	require.NotEqual(t, hb, b.Body.Hash())
	hashes = append(hashes, txn.Hash())
	require.Equal(t, b.Body.Hash(), cipher.Merkle(hashes))
}

func TestNewGenesisBlock(t *testing.T) {
	gb, err := NewGenesisBlock(genAddress, _genCoins, _genTime)
	require.NoError(t, err)

	require.Equal(t, cipher.SHA256{}, gb.Head.PrevHash)
	require.Equal(t, _genTime, gb.Head.Time)
	require.Equal(t, uint64(0), gb.Head.BkSeq)
	require.Equal(t, uint32(0), gb.Head.Version)
	require.Equal(t, uint64(0), gb.Head.Fee)
	require.Equal(t, cipher.SHA256{}, gb.Head.UxHash)

	require.Equal(t, 1, len(gb.Body.Transactions))
	txn := gb.Body.Transactions[0]
	require.Len(t, txn.In, 0)
	require.Len(t, txn.Sigs, 0)
	require.Len(t, txn.Out, 1)

	require.Equal(t, genAddress, txn.Out[0].Address)
	require.Equal(t, _genCoins, txn.Out[0].Coins)
	require.Equal(t, _genCoins, txn.Out[0].Hours)
}

func TestCreateUnspent(t *testing.T) {
	txn := Transaction{}
	err := txn.PushOutput(genAddress, 11e6, 255)
	require.NoError(t, err)
	bh := BlockHeader{
		Time:  tNow(),
		BkSeq: uint64(1),
	}

	tt := []struct {
		name    string
		txIndex int
		err     error
	}{
		{
			"ok",
			0,
			nil,
		},
		{
			"index overflow",
			10,
			errors.New("Transaction out index overflows transaction outputs"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			uxout, err := CreateUnspent(bh, txn, tc.txIndex)
			require.Equal(t, tc.err, err)
			if err != nil {
				return
			}
			requireUnspent(t, bh, txn, tc.txIndex, uxout)
		})
	}
}

func TestCreateUnspents(t *testing.T) {
	txn := Transaction{}
	err := txn.PushOutput(genAddress, 11e6, 255)
	require.NoError(t, err)
	bh := BlockHeader{
		Time:  tNow(),
		BkSeq: uint64(1),
	}
	uxouts := CreateUnspents(bh, txn)
	require.Equal(t, len(uxouts), 1)
	requireValidUnspents(t, bh, txn, uxouts)
}

func requireUnspent(t *testing.T, bh BlockHeader, txn Transaction, txIndex int, ux UxOut) {
	require.Equal(t, bh.Time, ux.Head.Time)
	require.Equal(t, bh.BkSeq, ux.Head.BkSeq)
	require.Equal(t, txn.Hash(), ux.Body.SrcTransaction)
	require.Equal(t, txn.Out[txIndex].Address, ux.Body.Address)
	require.Equal(t, txn.Out[txIndex].Coins, ux.Body.Coins)
	require.Equal(t, txn.Out[txIndex].Hours, ux.Body.Hours)
}

func requireValidUnspents(t *testing.T, bh BlockHeader, txn Transaction,
	uxo UxArray) {
	require.Equal(t, len(txn.Out), len(uxo))
	for i, ux := range uxo {
		require.Equal(t, bh.Time, ux.Head.Time)
		require.Equal(t, bh.BkSeq, ux.Head.BkSeq)
		require.Equal(t, txn.Hash(), ux.Body.SrcTransaction)
		require.Equal(t, txn.Out[i].Address, ux.Body.Address)
		require.Equal(t, txn.Out[i].Coins, ux.Body.Coins)
		require.Equal(t, txn.Out[i].Hours, ux.Body.Hours)
	}
}
