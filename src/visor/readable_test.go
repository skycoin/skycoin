package visor

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/boltdb/bolt"
	"github.com/skycoin/skycoin/src/util/droplet"
	"github.com/skycoin/skycoin/src/visor/blockdb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/testutil"
	"github.com/skycoin/skycoin/src/util/utc"
	"github.com/skycoin/skycoin/src/wallet"
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
func setupVisorConfig(t *testing.T) (Config, *Visor) {
	wltDir := prepareWltDir()
	mvc := NewVisorConfig()
	mvc.WalletDirectory = wltDir
	db, close := testutil.PrepareDB(t)
	defer close()
	mvc.GenesisSignature = createGenesisSignature(t)
	mv, err := NewVisor(mvc, db)
	require.NoError(t, err)
	// Use the master values for a client configuration
	c := NewVisorConfig()
	c.IsMaster = false
	c.GenesisSignature = mvc.GenesisSignature
	c.GenesisTimestamp = mvc.GenesisTimestamp
	c.WalletDirectory = wltDir
	return c, mv
}

func setupVisor(t *testing.T) (v *Visor, mv *Visor) {
	//TODO: close file
	db, _ := testutil.PrepareDB(t)
	vc, mv := setupVisorConfig(t)
	v, err := NewVisor(vc, db)
	require.NoError(t, err)
	return
}

// func setupVisorFromMaster(mv *Visor) *Visor {
// 	vc := NewVisorConfig()
// 	vc.IsMaster = false
// 	vc.MasterKeys = mv.Config.MasterKeys
// 	vc.MasterKeys.Secret = cipher.SecKey{}
// 	vc.GenesisSignature = mv.blockSigs.Sigs[0]
// 	vc.GenesisTimestamp = mv.blockchain.Blocks[0].Head.Time
// 	return NewVisor(vc)
// }

// func setupMasterVisorConfig() VisorConfig {
// 	// Create testmaster.keys file
// 	c := NewVisorConfig()
// 	c.CoinHourBurnFactor = 0
// 	c.IsMaster = true
// 	mw := wallet.NewWalletEntry()
// 	c.MasterKeys = mw
// 	c.GenesisSignature = createGenesisSignature(mw)
// 	return c
// }

// func setupMasterVisor() *Visor {
// 	return NewVisor(setupMasterVisorConfig())
// }

func cleanupVisor() {
	filenames := []string{
	//testMasterKeysFile,
	//testBlockchainFile,
	//testBlocksigsFile,
	//testWalletFile,
	//testWalletEntryFile,
	}
	for _, fn := range filenames {
		os.Remove(fn)
		os.Remove(fn + ".bak")
		os.Remove(fn + ".tmp")
	}
	wallets, err := filepath.Glob("*." + wallet.WalletExt)
	if err != nil {
		logger.Critical("Failed to glob wallet files: %v", err)
	} else {
		for _, w := range wallets {
			os.Remove(w)
			os.Remove(w + ".bak")
			os.Remove(w + ".tmp")
		}
	}
}

func createUnconfirmedTxn(t *testing.T) UnconfirmedTxn {
	ut := UnconfirmedTxn{}
	ut.Txn = coin.Transaction{}
	ut.Txn.InnerHash = testutil.RandSHA256(t)
	ut.Received = utc.Now().UnixNano()
	ut.Checked = ut.Received
	ut.Announced = time.Time{}.UnixNano()
	return ut
}

func addUnconfirmedTxn(t *testing.T, v *Visor) error {
	//ut := createUnconfirmedTxn(t)
	// Create enough unspent outputs to create all of these transactions
	//sb, err := v.Blockchain.Head()
	bc, ok := v.Blockchain.(*Blockchain)
	require.True(t, ok)
	uxOuts, err := v.Blockchain.Unspent().GetAll()
	require.NoError(t, err)
	//require.Len(t, uxOuts, 1)
	//uxs := coin.CreateUnspents(sb.Head, sb.Body.Transactions[0])
	_, _, toAddr := MakeAddress()
	txn := MakeTransactionForChain(t, bc, uxOuts[0], GenesisSecret, toAddr, 1, 1, 1)
	require.Equal(t, txn.Out[0].Address.String(), toAddr.String())
	//nUnspents := 100
	//txn := makeUnspentsTx(t, uxs, []cipher.SecKey{genSecret}, genAddress, nUnspents, maxDropletDivisor)
	up := NewUnconfirmedTxnPool(v.db)
	res, softErr, err := up.InjectTransaction(v.Blockchain, txn, 1)
	require.Equal(t, res, true)
	require.NoError(t, err)
	require.NoError(t, softErr)
	err = v.Blockchain.UpdateDB(func(tx *bolt.Tx) error {
		// update unconfirmed unspent
		head, err := v.Blockchain.Head()
		if err != nil {
			return err
		}
		b, err := v.Blockchain.NewBlock(coin.Transactions{txn}, head.Time()+uint64(100))
		require.NoError(t, err)

		sb := &coin.SignedBlock{
			Block: *b,
			Sig:   cipher.SignHash(b.HashHeader(), genSecret),
		}

		bcc, ok := v.Blockchain.(*Blockchain)
		require.True(t, ok)
		//bcc.store.UnspentPool().
		//bcc.store.AddBlockWithTx(tx, sb)
		up, err := blockdb.NewUnspentPool(bcc.db)
		//oldUxHash := up.GetUxHash()
		txHandler := up.ProcessBlock(sb)
		rb, err := txHandler(tx)
		rb()
		return nil
	})
	return err
}

// func addUnconfirmedTxnToPool(utp *UnconfirmedTxnPool) UnconfirmedTxn {
// 	ut := createUnconfirmedTxn()
// 	utp.txns.put(&ut)
// 	return ut
// }

//func transferCoinsToSelf(v *Visor, addr cipher.Address) error {
//	tx, err := v.Spend(v.Wallets[0].GetFilename(), wallet.Balance{1e6, 0}, 0, addr)
//	if err != nil {
//		return err
//	}
//	v.InjectTransaction(tx)
//	_, err = v.CreateAndExecuteBlock()
//	return err
//}

// func transferCoinsAdvanced(mv *Visor, v *Visor, b wallet.Balance, fee uint64,
// 	addr cipher.Address) error {
// 	tx, err := mv.Spend(mv.Wallets[0].GetFilename(), b, fee, addr)
// 	if err != nil {
// 		return err
// 	}
// 	mv.InjectTransaction(tx)
// 	now := uint64(utc.UnixNow())
// 	if len(mv.blockchain.Blocks) > 0 {
// 		now = mv.blockchain.Time() + 1
// 	}
// 	sb, err := mv.CreateBlock(now)
// 	if err != nil {
// 		return err
// 	}
// 	err = mv.ExecuteSignedBlock(sb)
// 	if err != nil {
// 		return err
// 	}
// 	return v.ExecuteSignedBlock(sb)
// }

func transferCoins(t *testing.T, mv *Visor, v *Visor) error {
	head := addGenesisBlock(t, v.Blockchain)
	toAddrs := make([]cipher.Address, 10)
	keys := make([]cipher.SecKey, 10)
	for i := 0; i < 10; i++ {
		p, s := cipher.GenerateKeyPair()
		toAddrs[i] = cipher.AddressFromPubKey(p)
		keys[i] = s
	}
	//toAddr := testutil.MakeAddress()
	//coins := uint64(10e6)
	var spend = spending{
		TxIndex: 0,
		UxIndex: 0,
		Keys:    []cipher.SecKey{genSecret},
		ToAddr:  toAddrs[0],
		Coins:   10e6,
	}
	// create normal spending tx
	uxs := coin.CreateUnspents(head.Head, head.Body.Transactions[0])
	//tx := makeSpendTx(t, uxs, []cipher.SecKey{genSecret}, toAddr, coins)
	//err := mv.Blockchain.VerifySingleTxnAllConstraints(tx, DefaultMaxBlockSize)
	//require.NoError(t, err)
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
	//Give the nonmaster some money to spend
	//addr := v.Wallets[0].GetAddresses()[0]
	//return transferCoinsAdvanced(mv, v, wallet.Balance{10e6, 0}, 0, addr)
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
	defer cleanupVisor()
	v, mv := setupVisor(t)
	assert.Nil(t, transferCoins(t, mv, v))
	//addUnconfirmedTxn(t, v)
	//addUnconfirmedTxn(v)

	bcm := NewBlockchainMetadata(v)
	assert.Equal(t, uint64(2), bcm.Unspents)
	assert.Equal(t, uint64(0), bcm.Unconfirmed)
	b, err := v.Blockchain.Head()
	require.NoError(t, err)
	//require.Equal(t, err, errors.New("found no head block: 0"))
	//fmt.Printf("%s", b)
	assertReadableBlockHeader(t, bcm.Head, b.Block.Head)
	//assertJSONSerializability(t, &bcm)
}

// func TestNewTransactionStatus(t *testing.T) {
// 	ts := NewUnconfirmedTransactionStatus()
// 	assert.True(t, ts.Unconfirmed)
// 	assert.False(t, ts.Unknown)
// 	assert.False(t, ts.Confirmed)
// 	assert.Equal(t, ts.Height, uint64(0))
// 	assertJSONSerializability(t, &ts)

// 	ts = NewUnknownTransactionStatus()
// 	assert.False(t, ts.Unconfirmed)
// 	assert.True(t, ts.Unknown)
// 	assert.False(t, ts.Confirmed)
// 	assert.Equal(t, ts.Height, uint64(0))
// 	assertJSONSerializability(t, &ts)

// 	ts = NewConfirmedTransactionStatus(uint64(7))
// 	assert.False(t, ts.Unconfirmed)
// 	assert.False(t, ts.Unknown)
// 	assert.True(t, ts.Confirmed)
// 	assert.Equal(t, ts.Height, uint64(7))
// 	assertJSONSerializability(t, &ts)

// 	assert.Panics(t, func() { NewConfirmedTransactionStatus(uint64(0)) })
// }

//func assertReadableTransactionHeader(t *testing.T,
//	rth ReadableTransactionHeader, th coin.TransactionHeader) {
//	assert.Equal(t, len(rth.Sigs), len(th.Sigs))
//	assert.NotPanics(t, func() {
//		for i, s := range rth.Sigs {
//			assert.Equal(t, cipher.MustSigFromHex(s), th.Sigs[i])
//		}
//		assert.Equal(t, cipher.MustSHA256FromHex(rth.Hash), th.Hash)
//	})
//	assertJSONSerializability(t, &rth)
//}

// func TestReadableTransactionHeader(t *testing.T) {
// 	defer cleanupVisor()
// 	v, mv := setupVisor()
// 	assert.Nil(t, transferCoins(mv, v))
// 	b := mv.blockchain.Head()
// 	th := b.Body.Transactions[0].Head
// 	rth := NewReadableTransactionHeader(&th)
// 	assertReadableTransactionHeader(t, rth, th)
// }

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

// func TestReadableTransactionOutput(t *testing.T) {
// 	defer cleanupVisor()
// 	v, mv := setupVisor()
// 	assert.Nil(t, transferCoins(mv, v))
// 	b := mv.blockchain.Head()
// 	to := b.Body.Transactions[0].Out[0]

// 	rto := NewReadableTransactionOutput(&to)
// 	assertReadableTransactionOutput(t, rto, to)
// }

func assertReadableTransactionInput(t *testing.T, rti string, ti cipher.SHA256) {
	assert.NotPanics(t, func() {
		assert.Equal(t, cipher.MustSHA256FromHex(rti), ti)
	})
	assertJSONSerializability(t, &rti)
}

// func TestReadableTransactionInput(t *testing.T) {
// 	defer cleanupVisor()
// 	v, mv := setupVisor()
// 	assert.Nil(t, transferCoins(mv, v))
// 	b := mv.blockchain.Head()
// 	ti := b.Body.Transactions[0].In[0]
// 	rti := ti.Hex()
// 	assertReadableTransactionInput(t, rti, ti)
// }

func assertReadableTransaction(t *testing.T, rtx ReadableTransaction,
	tx coin.Transaction) {
	assert.Equal(t, len(tx.In), len(rtx.In))
	assert.Equal(t, len(tx.Out), len(rtx.Out))
	//assertReadableTransactionHeader(t, rtx.Head, tx.Head)
	for i, ti := range rtx.In {
		assertReadableTransactionInput(t, ti, tx.In[i])
	}
	for i, to := range rtx.Out {
		assertReadableTransactionOutput(t, to, tx.Out[i])
	}
	assertJSONSerializability(t, &rtx)
}

// func TestReadableTransaction(t *testing.T) {
// 	defer cleanupVisor()
// 	v, mv := setupVisor()
// 	assert.Nil(t, transferCoins(mv, v))
// 	b := mv.blockchain.Head()
// 	tx := b.Body.Transactions[0]

// 	rtx := NewReadableTransaction(&tx)
// 	assertReadableTransaction(t, rtx, tx)
// }

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
	defer cleanupVisor()
	v, mv := setupVisor(t)
	assert.Nil(t, transferCoins(t, mv, v))
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
	defer cleanupVisor()
	v, mv := setupVisor(t)
	assert.Nil(t, transferCoins(t, mv, v))
	sb, err := v.Blockchain.Head()
	require.NoError(t, err)
	assert.Equal(t, sb.Head.BkSeq, uint64(1))
	rb, err := NewReadableBlock(&sb.Block)
	assertReadableBlock(t, *rb, sb.Block)
}
