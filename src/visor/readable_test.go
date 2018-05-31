package visor

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/util/droplet"
	"github.com/skycoin/skycoin/src/visor/dbutil"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/testutil"
)

func prepareWltDir(t *testing.T) string {
	dir, err := ioutil.TempDir("", "wallets")
	require.NoError(t, err)
	return dir
}

// Returns an appropriate VisorConfig and a master visor
func setupVisorConfig(t *testing.T) Config {
	wltDir := prepareWltDir(t)
	c := NewVisorConfig()
	c.WalletDirectory = wltDir
	c.BlockchainSeckey = genSecret
	c.BlockchainPubkey = genPublic
	c.GenesisAddress = genAddress
	return c
}

func setupVisor(t *testing.T) (*Visor, func()) {
	db, shutdown := prepareDB(t)
	vc := setupVisorConfig(t)
	v, err := NewVisor(vc, db)
	require.NoError(t, err)
	return v, shutdown
}

func transferCoins(t *testing.T, v *Visor) {
	head := addGenesisBlockToVisor(t, v)
	toAddrs := make([]cipher.Address, 10)
	keys := make([]cipher.SecKey, 10)
	for i := 0; i < 10; i++ {
		p, s := cipher.GenerateKeyPair()
		toAddrs[i] = cipher.AddressFromPubKey(p)
		keys[i] = s
	}

	var spend = spending{
		TxIndex: 0,
		UxIndex: 0,
		Keys:    []cipher.SecKey{genSecret},
		ToAddr:  toAddrs[0],
		Coins:   10e6,
	}
	// create normal spending tx
	uxs := coin.CreateUnspents(head.Head, head.Body.Transactions[0])
	txn := makeSpendTx(t, coin.UxArray{uxs[spend.UxIndex]}, spend.Keys, spend.ToAddr, spend.Coins)

	var b *coin.Block
	err := v.DB.View("", func(tx *dbutil.Tx) error {
		var err error
		b, err = v.Blockchain.NewBlock(tx, coin.Transactions{txn}, head.Time()+uint64(100))
		require.NoError(t, err)
		return nil
	})
	require.NoError(t, err)

	sb := &coin.SignedBlock{
		Block: *b,
		Sig:   cipher.SignHash(b.HashHeader(), genSecret),
	}
	v.DB.Update("", func(tx *dbutil.Tx) error {
		bcc, ok := v.Blockchain.(*Blockchain)
		require.True(t, ok)
		return bcc.store.AddBlock(tx, sb)
	})
	head = sb

}

func assertJSONSerializability(t *testing.T, thing interface{}) {
	b, err := json.Marshal(thing)
	require.NoError(t, err)
	rt := reflect.TypeOf(reflect.Indirect(reflect.ValueOf(thing)).Interface())
	newThing := reflect.New(rt).Interface()
	err = json.Unmarshal(b, newThing)
	require.NoError(t, err)
	require.True(t, reflect.DeepEqual(thing, newThing))
}

func TestNewBlockchainMetadata(t *testing.T) {
	v, shutdown := setupVisor(t)
	defer shutdown()

	transferCoins(t, v)

	bcm, err := v.GetBlockchainMetadata()
	require.NoError(t, err)
	require.Equal(t, uint64(2), bcm.Unspents)
	require.Equal(t, uint64(0), bcm.Unconfirmed)
	b, err := v.GetHeadBlock()
	require.NoError(t, err)
	assertReadableBlockHeader(t, bcm.Head, b.Block.Head)
	assertJSONSerializability(t, &bcm)
}

func TestNewTransactionStatus(t *testing.T) {
	ts := NewUnconfirmedTransactionStatus()
	require.True(t, ts.Unconfirmed)
	require.False(t, ts.Unknown)
	require.False(t, ts.Confirmed)
	require.Equal(t, ts.Height, uint64(0))
	assertJSONSerializability(t, &ts)

	ts = NewUnknownTransactionStatus()
	require.False(t, ts.Unconfirmed)
	require.True(t, ts.Unknown)
	require.False(t, ts.Confirmed)
	require.Equal(t, ts.Height, uint64(0))
	assertJSONSerializability(t, &ts)

	ts = NewConfirmedTransactionStatus(uint64(7), uint64(7))
	require.False(t, ts.Unconfirmed)
	require.False(t, ts.Unknown)
	require.True(t, ts.Confirmed)
	require.Equal(t, ts.Height, uint64(7))
	assertJSONSerializability(t, &ts)

	require.Panics(t, func() {
		NewConfirmedTransactionStatus(uint64(0), uint64(0))
	})
}

func assertReadableTransactionOutput(t *testing.T,
	rto ReadableTransactionOutput, to coin.TransactionOutput) {
	require.NotPanics(t, func() {
		require.Equal(t, cipher.MustDecodeBase58Address(rto.Address), to.Address)
	})
	coins, err := droplet.ToString(to.Coins)
	require.NoError(t, err)
	require.Equal(t, rto.Coins, coins)
	require.Equal(t, rto.Hours, to.Hours)
	assertJSONSerializability(t, &rto)
}

func TestReadableTransactionOutput(t *testing.T) {
	v, shutdown := setupVisor(t)
	defer shutdown()

	transferCoins(t, v)
	b, err := v.GetHeadBlock()
	require.NoError(t, err)
	to := b.Body.Transactions[0].Out[0]

	rto, err := NewReadableTransactionOutput(&to, testutil.RandSHA256(t))
	assertReadableTransactionOutput(t, *rto, to)
}

func assertReadableTransactionInput(t *testing.T, rti string, ti cipher.SHA256) {
	require.NotPanics(t, func() {
		require.Equal(t, cipher.MustSHA256FromHex(rti), ti)
	})
	assertJSONSerializability(t, &rti)
}

func TestReadableTransactionInput(t *testing.T) {
	v, shutdown := setupVisor(t)
	defer shutdown()

	transferCoins(t, v)
	b, err := v.GetHeadBlock()
	require.NoError(t, err)
	ti := b.Body.Transactions[0].In[0]
	rti := ti.Hex()
	assertReadableTransactionInput(t, rti, ti)
}

func assertReadableTransaction(t *testing.T, rtx ReadableTransaction, tx coin.Transaction) {
	require.Equal(t, len(tx.In), len(rtx.In))
	require.Equal(t, len(tx.Out), len(rtx.Out))
	for i, ti := range rtx.In {
		assertReadableTransactionInput(t, ti, tx.In[i])
	}
	for i, to := range rtx.Out {
		assertReadableTransactionOutput(t, to, tx.Out[i])
	}
	assertJSONSerializability(t, &rtx)
}

func TestReadableTransaction(t *testing.T) {
	v, shutdown := setupVisor(t)
	defer shutdown()

	transferCoins(t, v)
	b, err := v.GetHeadBlock()
	require.NoError(t, err)
	tx := b.Body.Transactions[0]

	rtx, err := NewReadableTransaction(&Transaction{
		Txn: tx,
	})
	assertReadableTransaction(t, *rtx, tx)
}

func assertReadableBlockHeader(t *testing.T, rb ReadableBlockHeader, bh coin.BlockHeader) {
	require.Equal(t, rb.Version, bh.Version)
	require.Equal(t, rb.Time, bh.Time)
	require.Equal(t, rb.BkSeq, bh.BkSeq)
	require.Equal(t, rb.Fee, bh.Fee)
	require.NotPanics(t, func() {
		require.Equal(t, cipher.MustSHA256FromHex(rb.PreviousBlockHash), bh.PrevHash)
		require.Equal(t, cipher.MustSHA256FromHex(rb.BodyHash), bh.BodyHash)
	})
	assertJSONSerializability(t, &rb)
}

func TestNewReadableBlockHeader(t *testing.T) {
	v, shutdown := setupVisor(t)
	defer shutdown()

	transferCoins(t, v)

	bh, err := v.GetHeadBlock()
	require.NoError(t, err)
	require.Equal(t, bh.Head.BkSeq, uint64(1))
	rb := NewReadableBlockHeader(&bh.Head)
	assertReadableBlockHeader(t, rb, bh.Head)
}

func assertReadableBlockBody(t *testing.T, rbb ReadableBlockBody, bb coin.BlockBody) {
	require.Equal(t, len(rbb.Transactions), len(bb.Transactions))
	for i, rt := range rbb.Transactions {
		assertReadableTransaction(t, rt, bb.Transactions[i])
	}
	assertJSONSerializability(t, &rbb)
}

func assertReadableBlock(t *testing.T, rb ReadableBlock, b coin.Block) {
	assertReadableBlockHeader(t, rb.Head, b.Head)
	assertReadableBlockBody(t, rb.Body, b.Body)
	assertJSONSerializability(t, &rb)
}

func TestNewReadableBlock(t *testing.T) {
	v, shutdown := setupVisor(t)
	defer shutdown()

	transferCoins(t, v)
	sb, err := v.GetHeadBlock()
	require.NoError(t, err)
	require.Equal(t, sb.Head.BkSeq, uint64(1))
	rb, err := NewReadableBlock(&sb.Block)
	assertReadableBlock(t, *rb, sb.Block)
}
