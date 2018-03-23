package visor

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/boltdb/bolt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/util/droplet"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/testutil"
)

func prepareWltDir() string {
	dir, err := ioutil.TempDir("", "wallets")
	if err != nil {
		panic(err)
	}

	return dir
}

func createGenesisSignature(t *testing.T) cipher.Sig {

	_, s := cipher.GenerateKeyPair()
	gb, err := coin.NewGenesisBlock(GenesisAddress, GenesisCoins, GenesisTime)
	if err != nil {
		panic(fmt.Errorf("create genesis block failed: %v", err))
	}

	sig := cipher.SignHash(gb.HashHeader(), s)
	return sig
}

// Returns an appropriate VisorConfig and a master visor
func setupVisorConfig(t *testing.T) Config {
	wltDir := prepareWltDir()
	c := NewVisorConfig()
	c.WalletDirectory = wltDir
	c.GenesisSignature = createGenesisSignature(t)
	return c
}

func setupVisor(t *testing.T) (v *Visor, close func()) {
	db, close := testutil.PrepareDB(t)
	vc := setupVisorConfig(t)
	v, err := NewVisor(vc, db)
	require.NoError(t, err)
	return
}

func transferCoins(t *testing.T, v *Visor) error {
	head := addGenesisBlock(t, v.Blockchain)
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
	tx := makeSpendTx(t, coin.UxArray{uxs[spend.UxIndex]}, spend.Keys, spend.ToAddr, spend.Coins)
	b, err := v.Blockchain.NewBlock(coin.Transactions{tx}, head.Time()+uint64(100))
	require.NoError(t, err)

	sb := &coin.SignedBlock{
		Block: *b,
		Sig:   cipher.SignHash(b.HashHeader(), genSecret),
	}
	v.db.Update(func(tx *bolt.Tx) error {
		bcc, ok := v.Blockchain.(*Blockchain)
		require.True(t, ok)
		return bcc.store.AddBlockWithTx(tx, sb)
	})
	head = sb
	return nil
}

func assertJSONSerializability(t *testing.T, thing interface{}) {
	b, err := json.Marshal(thing)
	assert.Nil(t, err)
	rt := reflect.TypeOf(reflect.Indirect(reflect.ValueOf(thing)).Interface())
	newThing := reflect.New(rt).Interface()
	err = json.Unmarshal(b, newThing)
	assert.Nil(t, err)
	assert.True(t, reflect.DeepEqual(thing, newThing))
}

func TestNewBlockchainMetadata(t *testing.T) {
	v, close := setupVisor(t)
	defer close()
	assert.Nil(t, transferCoins(t, v))

	bcm := NewBlockchainMetadata(v)
	assert.Equal(t, uint64(2), bcm.Unspents)
	assert.Equal(t, uint64(0), bcm.Unconfirmed)
	b, err := v.Blockchain.Head()
	require.NoError(t, err)
	assertReadableBlockHeader(t, bcm.Head, b.Block.Head)
	assertJSONSerializability(t, &bcm)
}

func TestNewTransactionStatus(t *testing.T) {
	ts := NewUnconfirmedTransactionStatus()
	assert.True(t, ts.Unconfirmed)
	assert.False(t, ts.Unknown)
	assert.False(t, ts.Confirmed)
	assert.Equal(t, ts.Height, uint64(0))
	assertJSONSerializability(t, &ts)

	ts = NewUnknownTransactionStatus()
	assert.False(t, ts.Unconfirmed)
	assert.True(t, ts.Unknown)
	assert.False(t, ts.Confirmed)
	assert.Equal(t, ts.Height, uint64(0))
	assertJSONSerializability(t, &ts)

	ts = NewConfirmedTransactionStatus(uint64(7), uint64(7))
	assert.False(t, ts.Unconfirmed)
	assert.False(t, ts.Unknown)
	assert.True(t, ts.Confirmed)
	assert.Equal(t, ts.Height, uint64(7))
	assertJSONSerializability(t, &ts)

	assert.Panics(t, func() { NewConfirmedTransactionStatus(uint64(0), uint64(0)) })
}

func assertReadableTransactionOutput(t *testing.T,
	rto ReadableTransactionOutput, to coin.TransactionOutput) {
	assert.NotPanics(t, func() {
		assert.Equal(t, cipher.MustDecodeBase58Address(rto.Address),
			to.Address)
	})
	coins, err := droplet.ToString(to.Coins)
	require.NoError(t, err)
	assert.Equal(t, rto.Coins, coins)
	assert.Equal(t, rto.Hours, to.Hours)
	assertJSONSerializability(t, &rto)
}

func TestReadableTransactionOutput(t *testing.T) {
	v, close := setupVisor(t)
	defer close()
	assert.Nil(t, transferCoins(t, v))
	b, err := v.Blockchain.Head()
	require.NoError(t, err)
	to := b.Body.Transactions[0].Out[0]

	rto, err := NewReadableTransactionOutput(&to, testutil.RandSHA256(t))
	assertReadableTransactionOutput(t, *rto, to)
}

func assertReadableTransactionInput(t *testing.T, rti string, ti cipher.SHA256) {
	assert.NotPanics(t, func() {
		assert.Equal(t, cipher.MustSHA256FromHex(rti), ti)
	})
	assertJSONSerializability(t, &rti)
}

func TestReadableTransactionInput(t *testing.T) {
	v, close := setupVisor(t)
	defer close()
	assert.Nil(t, transferCoins(t, v))
	b, err := v.Blockchain.Head()
	require.NoError(t, err)
	ti := b.Body.Transactions[0].In[0]
	rti := ti.Hex()
	assertReadableTransactionInput(t, rti, ti)
}

func assertReadableTransaction(t *testing.T, rtx ReadableTransaction,
	tx coin.Transaction) {
	assert.Equal(t, len(tx.In), len(rtx.In))
	assert.Equal(t, len(tx.Out), len(rtx.Out))
	for i, ti := range rtx.In {
		assertReadableTransactionInput(t, ti, tx.In[i])
	}
	for i, to := range rtx.Out {
		assertReadableTransactionOutput(t, to, tx.Out[i])
	}
	assertJSONSerializability(t, &rtx)
}

func TestReadableTransaction(t *testing.T) {
	v, close := setupVisor(t)
	defer close()
	assert.Nil(t, transferCoins(t, v))
	b, err := v.Blockchain.Head()
	require.NoError(t, err)
	tx := b.Body.Transactions[0]

	rtx, err := NewReadableTransaction(&Transaction{Txn: tx})
	assertReadableTransaction(t, *rtx, tx)
}

func assertReadableBlockHeader(t *testing.T, rb ReadableBlockHeader,
	bh coin.BlockHeader) {
	assert.Equal(t, rb.Version, bh.Version)
	assert.Equal(t, rb.Time, bh.Time)
	assert.Equal(t, rb.BkSeq, bh.BkSeq)
	assert.Equal(t, rb.Fee, bh.Fee)
	assert.NotPanics(t, func() {
		assert.Equal(t, cipher.MustSHA256FromHex(rb.PreviousBlockHash), bh.PrevHash)
		assert.Equal(t, cipher.MustSHA256FromHex(rb.BodyHash), bh.BodyHash)
	})
	assertJSONSerializability(t, &rb)
}

func TestNewReadableBlockHeader(t *testing.T) {
	v, close := setupVisor(t)
	defer close()
	assert.Nil(t, transferCoins(t, v))
	bh, err := v.Blockchain.Head()
	require.NoError(t, err)
	assert.Equal(t, bh.Head.BkSeq, uint64(1))
	rb := NewReadableBlockHeader(&bh.Head)
	assertReadableBlockHeader(t, rb, bh.Head)
}

func assertReadableBlockBody(t *testing.T, rbb ReadableBlockBody,
	bb coin.BlockBody) {
	assert.Equal(t, len(rbb.Transactions), len(bb.Transactions))
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
	v, close := setupVisor(t)
	defer close()
	assert.Nil(t, transferCoins(t, v))
	sb, err := v.Blockchain.Head()
	require.NoError(t, err)
	assert.Equal(t, sb.Head.BkSeq, uint64(1))
	rb, err := NewReadableBlock(&sb.Block)
	assertReadableBlock(t, *rb, sb.Block)
}
