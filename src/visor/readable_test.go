// +build ignore
package visor

import (
	"crypto/rand"
	"time"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/util/utc"
)

const (
	testMasterKeysFile  = "testmaster.keys"
	testWalletFile      = "testwallet.wlt"
	testBlocksigsFile   = "testblockchain.sigs"
	testBlockchainFile  = "testblockchain.bin"
	testWalletEntryFile = "testwalletentry.json"
	testWalletDir       = "./"
)

// func createGenesisSignature(master wallet.WalletEntry) cipher.Sig {
// 	c := NewVisorConfig()
// 	bc := coin.NewBlockchain()
// 	gb := bc.CreateGenesisBlock(master.Address, c.GenesisTimestamp,
// 		c.GenesisCoinVolume)
// 	return cipher.SignHash(gb.HashHeader(), master.Secret)
// }

// Returns an appropriate VisorConfig and a master visor
// func setupVisorConfig() (VisorConfig, *Visor) {
// 	// Make a new master visor + blockchain
// 	// Get the signed genesis block,
// 	mw := wallet.NewWalletEntry()
// 	mvc := NewVisorConfig()
// 	mvc.CoinHourBurnFactor = 0
// 	mvc.IsMaster = true
// 	mvc.MasterKeys = mw
// 	mvc.GenesisSignature = createGenesisSignature(mw)
// 	mv := NewVisor(mvc)

// 	// Use the master values for a client configuration
// 	c := NewVisorConfig()
// 	c.IsMaster = false
// 	c.GenesisSignature = mvc.GenesisSignature
// 	c.GenesisTimestamp = mvc.GenesisTimestamp
// 	c.MasterKeys = mw
// 	c.MasterKeys.Secret = cipher.SecKey{}
// 	c.WalletDirectory = testWalletDir
// 	return c, mv
// }

// func setupVisor() (v *Visor, mv *Visor) {
// 	vc, mv := setupVisorConfig()
// 	v = NewVisor(vc)
// 	return
// }

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

// func cleanupVisor() {
// 	filenames := []string{
// 		testMasterKeysFile,
// 		testBlockchainFile,
// 		testBlocksigsFile,
// 		testWalletFile,
// 		testWalletEntryFile,
// 	}
// 	for _, fn := range filenames {
// 		os.Remove(fn)
// 		os.Remove(fn + ".bak")
// 		os.Remove(fn + ".tmp")
// 	}
// 	wallets, err := filepath.Glob("*." + wallet.WalletExt)
// 	if err != nil {
// 		logger.Critical("Failed to glob wallet files: %v", err)
// 	} else {
// 		for _, w := range wallets {
// 			os.Remove(w)
// 			os.Remove(w + ".bak")
// 			os.Remove(w + ".tmp")
// 		}
// 	}
// }

func randSHA256() cipher.SHA256 {
	b := make([]byte, 128)
	rand.Read(b)
	return cipher.SumSHA256(b)
}

func createUnconfirmedTxn() UnconfirmedTxn {
	ut := UnconfirmedTxn{}
	ut.Txn = coin.Transaction{}
	ut.Txn.InnerHash = randSHA256()
	ut.Received = utc.Now().UnixNano()
	ut.Checked = ut.Received
	ut.Announced = time.Time{}.UnixNano()
	return ut
}

// func addUnconfirmedTxn(v *Visor) UnconfirmedTxn {
// 	ut := createUnconfirmedTxn()
// 	ut.Hash()
// 	v.Unconfirmed.txns.put(&ut)
// 	return ut
// }

// func addUnconfirmedTxnToPool(utp *UnconfirmedTxnPool) UnconfirmedTxn {
// 	ut := createUnconfirmedTxn()
// 	utp.txns.put(&ut)
// 	return ut
// }

// func transferCoinsToSelf(v *Visor, addr cipher.Address) error {
// 	tx, err := v.Spend(v.Wallets[0].GetFilename(), wallet.Balance{1e6, 0}, 0, addr)
// 	if err != nil {
// 		return err
// 	}
// 	v.InjectTxn(tx)
// 	_, err = v.CreateAndExecuteBlock()
// 	return err
// }

// func transferCoinsAdvanced(mv *Visor, v *Visor, b wallet.Balance, fee uint64,
// 	addr cipher.Address) error {
// 	tx, err := mv.Spend(mv.Wallets[0].GetFilename(), b, fee, addr)
// 	if err != nil {
// 		return err
// 	}
// 	mv.InjectTxn(tx)
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

// func transferCoins(mv *Visor, v *Visor) error {
// 	// Give the nonmaster some money to spend
// 	addr := v.Wallets[0].GetAddresses()[0]
// 	return transferCoinsAdvanced(mv, v, wallet.Balance{10e6, 0}, 0, addr)
// }

// func assertJSONSerializability(t *testing.T, thing interface{}) {
// 	b, err := json.Marshal(thing)
// 	assert.Nil(t, err)
// 	rt := reflect.TypeOf(reflect.Indirect(reflect.ValueOf(thing)).Interface())
// 	newThing := reflect.New(rt).Interface()
// 	err = json.Unmarshal(b, newThing)
// 	assert.Nil(t, err)
// 	assert.True(t, reflect.DeepEqual(thing, newThing))
// }

// func TestNewBlockchainMetadata(t *testing.T) {
// 	defer cleanupVisor()
// 	v, mv := setupVisor()
// 	assert.Nil(t, transferCoins(mv, v))
// 	addUnconfirmedTxn(v)
// 	addUnconfirmedTxn(v)

// 	bcm := NewBlockchainMetadata(v)
// 	assert.Equal(t, bcm.Unspents, uint64(2))
// 	assert.Equal(t, bcm.Unconfirmed, uint64(2))
// 	assertReadableBlockHeader(t, bcm.Head, v.blockchain.Head().Head)
// 	assertJSONSerializability(t, &bcm)
// }

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

// func assertReadableTransactionHeader(t *testing.T,
// 	rth ReadableTransactionHeader, th coin.TransactionHeader) {
// 	assert.Equal(t, len(rth.Sigs), len(th.Sigs))
// 	assert.NotPanics(t, func() {
// 		for i, s := range rth.Sigs {
// 			assert.Equal(t, cipher.MustSigFromHex(s), th.Sigs[i])
// 		}
// 		assert.Equal(t, cipher.MustSHA256FromHex(rth.Hash), th.Hash)
// 	})
// 	assertJSONSerializability(t, &rth)
// }

// func TestReadableTransactionHeader(t *testing.T) {
// 	defer cleanupVisor()
// 	v, mv := setupVisor()
// 	assert.Nil(t, transferCoins(mv, v))
// 	b := mv.blockchain.Head()
// 	th := b.Body.Transactions[0].Head
// 	rth := NewReadableTransactionHeader(&th)
// 	assertReadableTransactionHeader(t, rth, th)
// }

// func assertReadableTransactionOutput(t *testing.T,
// 	rto ReadableTransactionOutput, to coin.TransactionOutput) {
// 	assert.NotPanics(t, func() {
// 		assert.Equal(t, cipher.MustDecodeBase58Address(rto.Address),
// 			to.Address)
// 	})
// 	assert.Equal(t, rto.Coins, to.Coins)
// 	assert.Equal(t, rto.Hours, to.Hours)
// 	assertJSONSerializability(t, &rto)
// }

// func TestReadableTransactionOutput(t *testing.T) {
// 	defer cleanupVisor()
// 	v, mv := setupVisor()
// 	assert.Nil(t, transferCoins(mv, v))
// 	b := mv.blockchain.Head()
// 	to := b.Body.Transactions[0].Out[0]

// 	rto := NewReadableTransactionOutput(&to)
// 	assertReadableTransactionOutput(t, rto, to)
// }

// func assertReadableTransactionInput(t *testing.T, rti string, ti cipher.SHA256) {
// 	assert.NotPanics(t, func() {
// 		assert.Equal(t, cipher.MustSHA256FromHex(rti), ti)
// 	})
// 	assertJSONSerializability(t, &rti)
// }

// func TestReadableTransactionInput(t *testing.T) {
// 	defer cleanupVisor()
// 	v, mv := setupVisor()
// 	assert.Nil(t, transferCoins(mv, v))
// 	b := mv.blockchain.Head()
// 	ti := b.Body.Transactions[0].In[0]
// 	rti := ti.Hex()
// 	assertReadableTransactionInput(t, rti, ti)
// }

// func assertReadableTransaction(t *testing.T, rtx ReadableTransaction,
// 	tx coin.Transaction) {
// 	assert.Equal(t, len(tx.In), len(rtx.In))
// 	assert.Equal(t, len(tx.Out), len(rtx.Out))
// 	assertReadableTransactionHeader(t, rtx.Head, tx.Head)
// 	for i, ti := range rtx.In {
// 		assertReadableTransactionInput(t, ti, tx.In[i])
// 	}
// 	for i, to := range rtx.Out {
// 		assertReadableTransactionOutput(t, to, tx.Out[i])
// 	}
// 	assertJSONSerializability(t, &rtx)
// }

// func TestReadableTransaction(t *testing.T) {
// 	defer cleanupVisor()
// 	v, mv := setupVisor()
// 	assert.Nil(t, transferCoins(mv, v))
// 	b := mv.blockchain.Head()
// 	tx := b.Body.Transactions[0]

// 	rtx := NewReadableTransaction(&tx)
// 	assertReadableTransaction(t, rtx, tx)
// }

// func assertReadableBlockHeader(t *testing.T, rb ReadableBlockHeader,
// 	bh coin.BlockHeader) {
// 	assert.Equal(t, rb.Version, bh.Version)
// 	assert.Equal(t, rb.Time, bh.Time)
// 	assert.Equal(t, rb.BkSeq, bh.BkSeq)
// 	assert.Equal(t, rb.Fee, bh.Fee)
// 	assert.NotPanics(t, func() {
// 		assert.Equal(t, cipher.MustSHA256FromHex(rb.PrevHash), bh.PrevHash)
// 		assert.Equal(t, cipher.MustSHA256FromHex(rb.BodyHash), bh.BodyHash)
// 	})
// 	assertJSONSerializability(t, &rb)
// }

// func TestNewReadableBlockHeader(t *testing.T) {
// 	defer cleanupVisor()
// 	v, mv := setupVisor()
// 	assert.Nil(t, transferCoins(mv, v))
// 	bh := mv.blockchain.Head().Head
// 	assert.Equal(t, bh.BkSeq, uint64(1))
// 	rb := NewReadableBlockHeader(&bh)
// 	assertReadableBlockHeader(t, rb, bh)
// }

// func assertReadableBlockBody(t *testing.T, rbb ReadableBlockBody,
// 	bb coin.BlockBody) {
// 	assert.Equal(t, len(rbb.Transactions), len(bb.Transactions))
// 	for i, rt := range rbb.Transactions {
// 		assertReadableTransaction(t, rt, bb.Transactions[i])
// 	}
// 	assertJSONSerializability(t, &rbb)
// }

// func assertReadableBlock(t *testing.T, rb ReadableBlock, b coin.Block) {
// 	assertReadableBlockHeader(t, rb.Head, b.Head)
// 	assertReadableBlockBody(t, rb.Body, b.Body)
// 	assertJSONSerializability(t, &rb)
// }

// func TestNewReadableBlock(t *testing.T) {
// 	defer cleanupVisor()
// 	v, mv := setupVisor()
// 	assert.Nil(t, transferCoins(mv, v))
// 	b := mv.blockchain.Head()
// 	assert.Equal(t, b.Head.BkSeq, uint64(1))
// 	rb := NewReadableBlock(&b)
// 	assertReadableBlock(t, rb, b)
// }
