package visor

// import (
// 	"os"
// 	"path/filepath"
// 	"testing"
// 	"time"

// 	"github.com/skycoin/skycoin/src/cipher"
// 	"github.com/skycoin/skycoin/src/coin"
// 	"github.com/skycoin/skycoin/src/util"
// 	"github.com/skycoin/skycoin/src/wallet"
// 	"github.com/stretchr/testify/assert"
// )

// /* Helper functions */

// func setupVisorWriting(vc VisorConfig) *Visor {
// 	vc.WalletDirectory = testWalletDir
// 	v := NewVisor(vc)
// 	v.Wallets[0].SetFilename(testWalletFile)
// 	cleanupVisor() // delete the automatically saved wallet
// 	return v
// }

// func writeVisorFilesDirect(t *testing.T, v *Visor) {
// 	assert.True(t, len(v.Wallets) > 0)
// 	for _, w := range v.Wallets {
// 		assert.NotEqual(t, w.GetFilename(), "")
// 	}
// 	assert.Nil(t, v.SaveWallet(v.Wallets[0].GetFilename()))
// 	assert.Nil(t, v.SaveBlockSigs())
// 	assert.Nil(t, v.SaveBlockchain())
// 	assertDirExists(t, v.Config.WalletDirectory)
// 	walletFile := filepath.Join(v.Config.WalletDirectory, testWalletFile)
// 	assertFileExists(t, walletFile)
// 	assertFileExists(t, v.Config.BlockSigsFile)
// 	assertFileExists(t, v.Config.BlockchainFile)
// }

// func writeVisorFiles(t *testing.T, vc VisorConfig) *Visor {
// 	cleanupVisor()
// 	v := setupVisorWriting(vc)
// 	writeVisorFilesDirect(t, v)
// 	return v
// }

// func newWalletEntry(t *testing.T) wallet.WalletEntry {
// 	we := wallet.NewWalletEntry()
// 	assert.Nil(t, we.Verify())
// 	return we
// }

// func setupGenesis(t *testing.T) (wallet.WalletEntry, cipher.Sig, uint64) {
// 	we := newWalletEntry(t)
// 	vc := NewVisorConfig()
// 	vc.IsMaster = true
// 	vc.MasterKeys = we
// 	vc.GenesisSignature = createGenesisSignature(we)
// 	v := NewVisor(vc)
// 	we.Secret = cipher.SecKey{}
// 	return we, v.blockSigs.Sigs[0], v.blockchain.Blocks[0].Head.Time
// }

// func newGenesisConfig(t *testing.T) VisorConfig {
// 	refvc := NewVisorConfig()
// 	we, sig, ts := setupGenesis(t)
// 	refvc.MasterKeys = we
// 	refvc.GenesisSignature = sig
// 	refvc.GenesisTimestamp = ts
// 	refvc.IsMaster = false
// 	refvc.WalletDirectory = testWalletDir
// 	return refvc
// }

// func corruptFile(t *testing.T, filename string) {
// 	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC, 0600)
// 	assert.Nil(t, err)
// 	_, err = f.Write([]byte("xxxxx"))
// 	assert.Nil(t, err)
// 	f.Close()
// }

// func setupChildVisorConfig(refvc VisorConfig, master bool) VisorConfig {
// 	vc := NewVisorConfig()
// 	vc.IsMaster = master
// 	vc.MasterKeys = refvc.MasterKeys
// 	vc.WalletDirectory = testWalletDir
// 	vc.BlockchainFile = testBlockchainFile
// 	vc.BlockSigsFile = testBlocksigsFile
// 	return vc
// }

// func newMasterVisorConfig(t *testing.T) VisorConfig {
// 	vc := NewVisorConfig()
// 	vc.CoinHourBurnFactor = 0
// 	mw := newWalletEntry(t)
// 	vc.MasterKeys = mw
// 	vc.GenesisSignature = createGenesisSignature(mw)
// 	vc.IsMaster = true
// 	return vc
// }

// func addValidTxns(t *testing.T, v *Visor, n int) coin.Transactions {
// 	txns := make(coin.Transactions, n)
// 	for i := 0; i < len(txns); i++ {
// 		txn, err := makeValidTxn(v)
// 		assert.Nil(t, err)
// 		txns[i] = txn
// 	}
// 	for _, txn := range txns {
// 		err, known := v.InjectTxn(txn)
// 		assert.Nil(t, err)
// 		assert.False(t, known)
// 	}
// 	txns = coin.SortTransactions(txns, getFee)
// 	assert.Equal(t, len(v.Unconfirmed.Txns), n)
// 	return txns
// }

// func addSignedBlockAt(t *testing.T, v *Visor, when uint64) coin.SignedBlock {
// 	we := wallet.NewWalletEntry()
// 	tx, err := v.Spend(v.Wallets[0].GetFilename(), wallet.Balance{1e6, 0}, 0, we.Address)
// 	assert.Nil(t, err)
// 	err, known := v.InjectTxn(tx)
// 	assert.Nil(t, err)
// 	assert.False(t, known)
// 	sb, err := v.CreateBlock(when)
// 	assert.Nil(t, err)
// 	if err != nil {
// 		return sb
// 	}
// 	err = v.ExecuteSignedBlock(sb)
// 	assert.Nil(t, err)
// 	return sb
// }

// func addSignedBlock(t *testing.T, v *Visor) coin.SignedBlock {
// 	return addSignedBlockAt(t, v, uint64(utc.UnixNow()))
// }

// func addSignedBlocks(t *testing.T, v *Visor, n int) []coin.SignedBlock {
// 	sbs := make([]SignedBlock, n)
// 	now := uint64(utc.UnixNow())
// 	for i := 0; i < n; i++ {
// 		sbs[i] = addSignedBlockAt(t, v, now+1+uint64(i))
// 	}
// 	return sbs
// }

// func assertSignedBlocks(t *testing.T, v *Visor, sbs []coin.SignedBlock,
// 	start, ct uint64) {
// 	have := v.HeadBkSeq()
// 	if have <= start {
// 		assert.Equal(t, len(sbs), 0)
// 	} else if have-start < ct {
// 		assert.Equal(t, len(sbs), int(have-start))
// 	} else {
// 		assert.Equal(t, len(sbs), int(ct))
// 	}
// 	for i, sb := range sbs {
// 		assert.Nil(t, v.verifySignedBlock(&sb))
// 		assert.Equal(t, sb.Sig, v.blockSigs.Sigs[sb.Block.Head.BkSeq])
// 		assert.Equal(t, sb.Block.Head.BkSeq, start+uint64(i)+1)
// 		assert.Equal(t, v.blockchain.Blocks[start+uint64(i)+1], sb.Block)
// 	}
// }

// func assertReadableBlocks(t *testing.T, v *Visor, rbs []ReadableBlock,
// 	sbs []coin.SignedBlock) {
// 	assert.Equal(t, len(rbs), len(sbs))
// 	for i, rb := range rbs {
// 		assertReadableBlock(t, rb, sbs[i].Block)
// 	}
// }

// func assertBlocks(t *testing.T, v *Visor, bs []coin.Block, sbs []coin.SignedBlock) {
// 	assert.Equal(t, len(bs), len(sbs))
// 	for i, b := range bs {
// 		assert.Equal(t, b, sbs[i].Block)
// 	}
// }

// /* Actual tests */

// func TestNewVisorConfig(t *testing.T) {
// 	vc := NewVisorConfig()
// 	assert.False(t, vc.IsMaster)
// 	assert.Equal(t, vc.WalletDirectory, "")
// 	assert.Equal(t, vc.BlockchainFile, "")
// 	assert.Equal(t, vc.BlockSigsFile, "")
// 	assert.Panics(t, func() { vc.MasterKeys.Verify() })
// 	assert.NotNil(t, vc.MasterKeys.VerifyPublic())
// 	assert.Equal(t, vc.GenesisSignature, cipher.Sig{})
// }

// func TestNewVisor(t *testing.T) {
// 	defer cleanupVisor()

// 	// Not master, Invalid master keys
// 	cleanupVisor()
// 	we := wallet.NewWalletEntry()
// 	we.Public = cipher.PubKey{}
// 	vc := NewVisorConfig()
// 	vc.IsMaster = false
// 	assert.Panics(t, func() { NewVisor(vc) })
// 	vc.MasterKeys = we
// 	assert.Panics(t, func() { NewVisor(vc) })

// 	// Master, invalid master keys
// 	cleanupVisor()
// 	vc.IsMaster = true
// 	vc.MasterKeys = wallet.WalletEntry{}
// 	assert.Panics(t, func() { NewVisor(vc) })
// 	vc.MasterKeys = we
// 	assert.Panics(t, func() { NewVisor(vc) })

// 	// Not master, no wallet, blockchain, blocksigs file
// 	cleanupVisor()
// 	vc = NewVisorConfig()
// 	vc.IsMaster = false
// 	we, sig, ts := setupGenesis(t)
// 	vc.MasterKeys = we
// 	vc.GenesisSignature = sig
// 	vc.GenesisTimestamp = ts
// 	v := NewVisor(vc)
// 	assert.Equal(t, len(v.blockchain.Blocks), 1)
// 	assert.Equal(t, len(v.blockSigs.Sigs), 1)

// 	// Master, no wallet, blockchain, blocksigs file
// 	cleanupVisor()
// 	vc = NewVisorConfig()
// 	we = newWalletEntry(t)
// 	vc.MasterKeys = we
// 	vc.GenesisSignature = createGenesisSignature(we)
// 	vc.IsMaster = true
// 	v = NewVisor(vc)
// 	assert.Equal(t, len(v.Wallets), 1)
// 	// Wallet should only have 1 entry if master
// 	assert.Equal(t, v.Wallets[0].NumEntries(), 1)
// 	assert.Equal(t, len(v.blockchain.Blocks), 1)
// 	assert.Equal(t, len(v.blockSigs.Sigs), 1)

// 	// Not master, has all files
// 	cleanupVisor()
// 	refvc := newGenesisConfig(t)
// 	refv := writeVisorFiles(t, refvc)
// 	vc = setupChildVisorConfig(refvc, false)
// 	v = NewVisor(vc)
// 	assert.Equal(t, v.Wallets, refv.Wallets)
// 	assert.Equal(t, len(v.Wallets), 1)
// 	assert.Equal(t, v.blockchain, refv.blockchain)
// 	assert.Equal(t, v.blockSigs, refv.blockSigs)

// 	// Master, has all files
// 	cleanupVisor()
// 	refvc = newMasterVisorConfig(t)
// 	refv = writeVisorFiles(t, refvc)
// 	vc = setupChildVisorConfig(refvc, true)
// 	v = NewVisor(vc)
// 	assert.Equal(t, v.Wallets[0].GetEntries(), refv.Wallets[0].GetEntries())
// 	assert.Equal(t, v.blockchain, refv.blockchain)

// 	// Not master, wallet is corrupt
// 	cleanupVisor()
// 	refvc = newGenesisConfig(t)
// 	refv = writeVisorFiles(t, refvc)
// 	walletFile := filepath.Join(testWalletDir, testWalletFile)
// 	assertFileExists(t, walletFile)
// 	corruptFile(t, walletFile)
// 	vc = setupChildVisorConfig(refvc, false)
// 	assert.Panics(t, func() { NewVisor(vc) })

// 	// Master, wallet is corrupt.  Nothing happens because master ignores
// 	// wallet
// 	cleanupVisor()
// 	refvc = newMasterVisorConfig(t)
// 	refv = writeVisorFiles(t, refvc)
// 	assertFileExists(t, walletFile)
// 	corruptFile(t, walletFile)
// 	vc = setupChildVisorConfig(refvc, true)
// 	assert.NotPanics(t, func() { NewVisor(vc) })

// 	// Not master, blocksigs is corrupt
// 	cleanupVisor()
// 	refvc = newGenesisConfig(t)
// 	assertFileNotExists(t, testWalletFile)
// 	refv = writeVisorFiles(t, refvc)
// 	corruptFile(t, testBlocksigsFile)
// 	vc = setupChildVisorConfig(refvc, false)
// 	assert.Panics(t, func() { NewVisor(vc) })

// 	// Master, blocksigs is corrupt
// 	cleanupVisor()
// 	refvc = newMasterVisorConfig(t)
// 	refv = writeVisorFiles(t, refvc)
// 	corruptFile(t, testBlocksigsFile)
// 	assertFileExists(t, testBlocksigsFile)
// 	vc = setupChildVisorConfig(refvc, true)
// 	assert.Panics(t, func() { NewVisor(vc) })

// 	// Not master, blockchain is corrupt
// 	cleanupVisor()
// 	refvc = newGenesisConfig(t)
// 	refv = writeVisorFiles(t, refvc)
// 	corruptFile(t, testBlockchainFile)
// 	vc = setupChildVisorConfig(refvc, false)
// 	assert.Panics(t, func() { NewVisor(vc) })

// 	// Master, blockchain is corrupt
// 	cleanupVisor()
// 	refvc = newMasterVisorConfig(t)
// 	refv = writeVisorFiles(t, refvc)
// 	corruptFile(t, testBlockchainFile)
// 	vc = setupChildVisorConfig(refvc, true)
// 	assert.Panics(t, func() { NewVisor(vc) })

// 	// Not master, blocksigs is not valid for blockchain
// 	cleanupVisor()
// 	refvc = newGenesisConfig(t)
// 	refv = setupVisorWriting(refvc)
// 	// Corrupt the signature
// 	refv.blockSigs.Sigs[uint64(0)] = cipher.Sig{}
// 	writeVisorFilesDirect(t, refv)
// 	vc = setupChildVisorConfig(refvc, false)
// 	assert.Panics(t, func() { NewVisor(vc) })

// 	// Master, blocksigs is not valid for blockchain
// 	cleanupVisor()
// 	refvc = newMasterVisorConfig(t)
// 	refv = setupVisorWriting(refvc)
// 	// Corrupt the signature
// 	refv.blockSigs.Sigs[uint64(0)] = cipher.Sig{}
// 	writeVisorFilesDirect(t, refv)
// 	vc = setupChildVisorConfig(refvc, true)
// 	assert.Panics(t, func() { NewVisor(vc) })
// }

// func TestNewMinimalVisor(t *testing.T) {
// 	defer cleanupVisor()
// 	vc := newMasterVisorConfig(t)
// 	v := NewMinimalVisor(vc)
// 	assert.Equal(t, v.Config, vc)
// 	assert.NotNil(t, v.Unconfirmed)
// 	assert.Nil(t, v.Wallets)
// 	assert.Equal(t, len(v.blockchain.Blocks), 0)
// 	assert.Equal(t, len(v.blockSigs.Sigs), 0)
// }

// func TestCreateGenesisBlockVisor(t *testing.T) {
// 	defer cleanupVisor()
// 	// Test as master, successful
// 	vc := newMasterVisorConfig(t)
// 	v := NewMinimalVisor(vc)
// 	assert.True(t, v.Config.IsMaster)
// 	assert.Equal(t, len(v.blockchain.Blocks), 0)
// 	assert.Equal(t, len(v.blockSigs.Sigs), 0)
// 	sb := v.CreateGenesisBlock()
// 	assert.NotEqual(t, sb.Block, coin.Block{})
// 	assert.NotEqual(t, sb.Sig, cipher.Sig{})
// 	assert.Equal(t, len(v.blockchain.Blocks), 1)
// 	assert.Equal(t, len(v.blockSigs.Sigs), 1)
// 	assert.Nil(t, v.blockSigs.Verify(vc.MasterKeys.Public, v.blockchain))

// 	// Test as not master, successful
// 	vc = newGenesisConfig(t)
// 	v = NewMinimalVisor(vc)
// 	assert.False(t, v.Config.IsMaster)
// 	assert.Equal(t, len(v.blockchain.Blocks), 0)
// 	assert.Equal(t, len(v.blockSigs.Sigs), 0)
// 	sb = v.CreateGenesisBlock()
// 	assert.NotEqual(t, sb.Block, coin.Block{})
// 	assert.NotEqual(t, sb.Sig, cipher.Sig{})
// 	assert.Equal(t, len(v.blockchain.Blocks), 1)
// 	assert.Equal(t, len(v.blockSigs.Sigs), 1)
// 	assert.Nil(t, v.blockSigs.Verify(vc.MasterKeys.Public, v.blockchain))
// 	assert.Equal(t, v.Config.GenesisSignature, sb.Sig)
// 	assert.Equal(t, v.blockchain.Blocks[0].Head.Time, v.Config.GenesisTimestamp)

// 	// Test as master, blockSigs invalid for pubkey
// 	vc = newMasterVisorConfig(t)
// 	vc.MasterKeys.Public = cipher.PubKey{}
// 	v = NewMinimalVisor(vc)
// 	assert.True(t, v.Config.IsMaster)
// 	assert.Equal(t, len(v.blockchain.Blocks), 0)
// 	assert.Equal(t, len(v.blockSigs.Sigs), 0)
// 	assert.Panics(t, func() { v.CreateGenesisBlock() })

// 	// Test as not master, blockSigs invalid for pubkey
// 	vc = newGenesisConfig(t)
// 	vc.MasterKeys.Public = cipher.PubKey{}
// 	v = NewMinimalVisor(vc)
// 	assert.False(t, v.Config.IsMaster)
// 	assert.Equal(t, len(v.blockchain.Blocks), 0)
// 	assert.Equal(t, len(v.blockSigs.Sigs), 0)
// 	assert.Panics(t, func() { v.CreateGenesisBlock() })

// 	// Test as master, signing failed
// 	vc = newMasterVisorConfig(t)
// 	vc.MasterKeys.Secret = cipher.SecKey{}
// 	vc.GenesisSignature = cipher.Sig{}
// 	assert.Equal(t, vc.MasterKeys.Secret, cipher.SecKey{})
// 	v = NewMinimalVisor(vc)
// 	assert.True(t, v.Config.IsMaster)
// 	assert.Equal(t, v.Config, vc)
// 	assert.Equal(t, v.Config.MasterKeys.Secret, cipher.SecKey{})
// 	assert.Equal(t, len(v.blockchain.Blocks), 0)
// 	assert.Equal(t, len(v.blockSigs.Sigs), 0)
// 	assert.Panics(t, func() { v.CreateGenesisBlock() })
// }

// func TestVisorRefreshUnconfirmed(t *testing.T) {
// 	defer cleanupVisor()
// 	mv := setupMasterVisor()
// 	testRefresh(t, mv, func(checkPeriod, maxAge time.Duration) {
// 		mv.Config.UnconfirmedCheckInterval = checkPeriod
// 		mv.Config.UnconfirmedMaxAge = maxAge
// 		mv.RefreshUnconfirmed()
// 	})
// }

// func TestVisorSaveBlockchain(t *testing.T) {
// 	cleanupVisor()
// 	defer cleanupVisor()
// 	vc := newMasterVisorConfig(t)
// 	vc.BlockchainFile = ""

// 	// Test with no blockchain file set
// 	v := NewVisor(vc)
// 	assertFileNotExists(t, testBlockchainFile)
// 	err := v.SaveBlockchain()
// 	assert.NotNil(t, err)
// 	assert.Equal(t, err.Error(), "No BlockchainFile location set")
// 	assertFileNotExists(t, testBlockchainFile)

// 	// Test with blockchain file set
// 	vc.BlockchainFile = testBlockchainFile
// 	v = NewVisor(vc)
// 	assert.Nil(t, v.SaveBlockchain())
// 	assertFileExists(t, testBlockchainFile)
// 	assert.NotPanics(t, func() {
// 		loadBlockchain(testBlockchainFile, vc.MasterKeys.Address)
// 	})
// 	bc := loadBlockchain(testBlockchainFile, vc.MasterKeys.Address)
// 	assert.Equal(t, v.blockchain, bc)
// }

// func TestVisorSaveWallets(t *testing.T) {
// 	cleanupVisor()
// 	defer cleanupVisor()
// 	vc := newGenesisConfig(t)
// 	assert.False(t, vc.IsMaster)
// 	vc.WalletDirectory = testWalletDir
// 	v := NewVisor(vc)
// 	assertFileNotExists(t, filepath.Join(testWalletDir, testWalletFile))
// 	v.Wallets[0].SetFilename(testWalletFile)
// 	assert.Equal(t, len(v.SaveWallets()), 0)
// 	assertFileExists(t, filepath.Join(testWalletDir, testWalletFile))
// 	w, err := wallet.LoadSimpleWallet(testWalletDir, testWalletFile)
// 	assert.Nil(t, err)
// 	assert.Equal(t, v.Wallets[0], w)
// }

// func TestVisorSaveBlockSigs(t *testing.T) {
// 	cleanupVisor()
// 	defer cleanupVisor()
// 	vc := newMasterVisorConfig(t)
// 	vc.BlockSigsFile = ""

// 	// Test with no blocksigs file set
// 	v := NewVisor(vc)
// 	assertFileNotExists(t, testBlocksigsFile)
// 	err := v.SaveBlockSigs()
// 	assert.NotNil(t, err)
// 	assert.Equal(t, err.Error(), "No BlockSigsFile location set")
// 	assertFileNotExists(t, testBlocksigsFile)

// 	vc.BlockSigsFile = testBlocksigsFile
// 	v = NewVisor(vc)
// 	assert.Nil(t, v.SaveBlockSigs())
// 	assertFileExists(t, testBlocksigsFile)

// 	bs, err := LoadBlockSigs(testBlocksigsFile)
// 	assert.Nil(t, err)
// 	assert.Equal(t, v.blockSigs, bs)
// }

// func TestCreateAndExecuteBlock(t *testing.T) {
// 	defer cleanupVisor()

// 	// Test as not master, should fail
// 	vc := newGenesisConfig(t)
// 	v := NewVisor(vc)
// 	assert.Panics(t, func() { v.CreateAndExecuteBlock() })

// 	// Test as master, no txns
// 	vc = newMasterVisorConfig(t)
// 	v = NewVisor(vc)
// 	_, err := v.CreateAndExecuteBlock()
// 	assert.NotNil(t, err)
// 	assert.Equal(t, err.Error(), "No transactions")

// 	// Test as master, more txns than allowed
// 	vc.BlockCreationInterval = uint64(101)
// 	v = NewVisor(vc)
// 	txns := addValidTxns(t, v, 3)
// 	txns = coin.SortTransactions(txns, v.blockchain.TransactionFee)
// 	v.Config.MaxBlockSize = txns[0].Size()
// 	assert.Equal(t, len(v.blockchain.Blocks), 1)
// 	assert.Equal(t, len(v.blockSigs.Sigs), 1)
// 	sb, err := v.CreateAndExecuteBlock()
// 	assert.Nil(t, err)

// 	assert.Equal(t, len(sb.Block.Body.Transactions), 1)
// 	assert.Equal(t, sb.Block.Body.Transactions[0], txns[0])
// 	assert.Equal(t, len(v.blockchain.Blocks), 2)
// 	assert.Equal(t, len(v.blockSigs.Sigs), 2)
// 	assert.Equal(t, v.blockchain.Blocks[1], sb.Block)
// 	assert.Equal(t, v.blockSigs.Sigs[1], sb.Sig)
// 	assert.Equal(t, len(v.Unconfirmed.Txns), 2)
// 	assert.True(t, sb.Block.Head.Time > v.blockchain.Blocks[0].Head.Time)
// 	rawTxns := v.Unconfirmed.RawTxns()
// 	assert.Equal(t, len(rawTxns), 2)
// 	for _, tx := range sb.Block.Body.Transactions {
// 		assert.NotEqual(t, tx.Hash(), rawTxns[0].Hash())
// 		assert.NotEqual(t, tx.Hash(), rawTxns[1].Hash())
// 	}
// 	if txns[1].Hash() == rawTxns[0].Hash() {
// 		assert.Equal(t, txns[2].Hash(), rawTxns[1].Hash())
// 	} else {
// 		assert.Equal(t, txns[2].Hash(), rawTxns[0].Hash())
// 	}
// 	assert.Nil(t, v.blockSigs.Verify(v.Config.MasterKeys.Public, v.blockchain))

// 	// No txns, forcing NewBlockFromTransactions to fail
// 	v = NewVisor(vc)
// 	assert.Equal(t, len(v.Unconfirmed.Txns), 0)
// 	txns = addValidTxns(t, v, 3)
// 	v.Config.MaxBlockSize = 0
// 	sb, err = v.CreateAndExecuteBlock()
// 	assert.NotNil(t, err)
// 	assert.Equal(t, len(v.blockchain.Blocks), 1)
// 	assert.Equal(t, len(v.blockSigs.Sigs), 1)
// 	assert.Equal(t, len(v.Unconfirmed.Txns), 3)
// }

// func TestVisorSpend(t *testing.T) {
// 	defer cleanupVisor()
// 	we := wallet.NewWalletEntry()
// 	addr := we.Address
// 	vc := newMasterVisorConfig(t)
// 	assert.Equal(t, vc.CoinHourBurnFactor, uint64(0))
// 	v := NewVisor(vc)
// 	wid := v.Wallets[0].GetFilename()
// 	ogb := v.WalletBalance(wid).Confirmed

// 	// Test spend 0 amount
// 	v = NewVisor(vc)
// 	b = wallet.Balance{0, 0}
// 	_, err = v.Spend(v.Wallets[0].GetFilename(), b, 0, addr)
// 	assert.NotNil(t, err)
// 	assert.Equal(t, err.Error(), "Zero spend amount")

// 	// Test lacking funds
// 	v = NewVisor(vc)
// 	b = wallet.Balance{10e16, 10e16}
// 	_, err = v.Spend(v.Wallets[0].GetFilename(), b, 10e16, addr)
// 	assert.NotNil(t, err)
// 	assert.Equal(t, err.Error(), "Not enough coins")

// 	// Test created txn too large
// 	v = NewVisor(vc)
// 	v.Config.MaxBlockSize = 0
// 	b = wallet.Balance{10e6, 10}
// 	assert.Panics(t, func() { v.Spend(v.Wallets[0].GetFilename(), b, 0, addr) })

// 	// Test simple spend (we have only 1 address to spend from, no fee)
// 	v = NewVisor(vc)
// 	assert.Equal(t, v.Config.CoinHourBurnFactor, uint64(0))
// 	b = wallet.Balance{10e6, 10}
// 	tx, err := v.Spend(v.Wallets[0].GetFilename(), b, 0, addr)
// 	assert.Nil(t, err)
// 	assert.Equal(t, len(tx.In), 1)
// 	assert.Equal(t, len(tx.Out), 2)
// 	// Hash should be updated
// 	assert.NotEqual(t, tx.Head.Hash, cipher.SHA256{})
// 	// Should be 1 signature for the single input
// 	assert.Equal(t, len(tx.Head.Sigs), 1)
// 	// Spent amount should be correct
// 	assert.Equal(t, tx.Out[1].Address, addr)
// 	assert.Equal(t, tx.Out[1].Coins, b.Coins)
// 	assert.Equal(t, tx.Out[1].Hours, b.Hours)
// 	// Change amount should be correct
// 	ourAddr := v.Wallets[0].GetAddresses()[0]
// 	assert.Equal(t, tx.Out[0].Address, ourAddr)
// 	assert.Equal(t, tx.Out[0].Coins, ogb.Coins-b.Coins)
// 	assert.Equal(t, tx.Out[0].Hours, ogb.Hours-b.Hours)
// 	assert.Nil(t, tx.Verify())
// }

// func TestExecuteSignedBlock(t *testing.T) {
// 	defer cleanupVisor()
// 	cleanupVisor()
// 	we := wallet.NewWalletEntry()
// 	vc := newMasterVisorConfig(t)
// 	v := NewVisor(vc)
// 	wid := v.Wallets[0].GetFilename()
// 	assert.Equal(t, len(v.Unconfirmed.Txns), 0)
// 	tx, err := v.Spend(wid, wallet.Balance{1e6, 0}, 0, we.Address)
// 	assert.Nil(t, err)
// 	err, known := v.InjectTxn(tx)
// 	assert.Nil(t, err)
// 	assert.False(t, known)
// 	assert.Equal(t, len(v.Unconfirmed.Txns), 1)
// 	assert.Equal(t, len(v.blockSigs.Sigs), 1)
// 	now := uint64(utc.UnixNow())

// 	// Invalid signed block
// 	sb, err := v.CreateBlock(now)
// 	assert.Equal(t, len(v.blockSigs.Sigs), 1)
// 	assert.Nil(t, err)
// 	sb.Sig = cipher.Sig{}
// 	err = v.ExecuteSignedBlock(sb)
// 	assert.NotNil(t, err)
// 	assert.Equal(t, len(v.Unconfirmed.Txns), 1)
// 	assert.Equal(t, len(v.blockSigs.Sigs), 1)

// 	// Invalid block
// 	sb, err = v.CreateBlock(now)
// 	assert.Nil(t, err)
// 	// TODO -- empty BodyHash is being accepted, fix blockchain verification
// 	sb.Block.Head.BodyHash = cipher.SHA256{}
// 	sb.Block.Body.Transactions = make(coin.Transactions, 0)
// 	sb = v.SignBlock(sb.Block)
// 	err = v.ExecuteSignedBlock(sb)
// 	assert.NotNil(t, err)
// 	assert.Equal(t, len(v.Unconfirmed.Txns), 1)
// 	assert.Equal(t, len(v.blockSigs.Sigs), 1)

// 	// Valid block
// 	sb, err = v.CreateBlock(now)
// 	assert.Nil(t, err)
// 	err = v.ExecuteSignedBlock(sb)
// 	assert.Nil(t, err)
// 	assert.Equal(t, len(v.blockSigs.Sigs), 2)
// 	assert.Equal(t, v.blockSigs.Sigs[uint64(1)], sb.Sig)
// 	assert.Equal(t, v.blockchain.Blocks[1], sb.Block)
// 	assert.Equal(t, len(v.Unconfirmed.Txns), 0)

// 	// Test a valid block created by a master but executing in non master
// 	vc2, mv := setupVisorConfig()
// 	v2 := NewVisor(vc2)
// 	w := v2.Wallets[0]
// 	addr := w.GetAddresses()[0]
// 	tx, err = mv.Spend(mv.Wallets[0].GetFilename(), wallet.Balance{1e6, 0}, 0, addr)
// 	assert.Nil(t, err)
// 	err, known = mv.InjectTxn(tx)
// 	assert.Nil(t, err)
// 	assert.False(t, known)
// 	sb, err = mv.CreateAndExecuteBlock()
// 	assert.Nil(t, err)
// 	err = v2.ExecuteSignedBlock(sb)
// 	assert.Nil(t, err)
// 	assert.Equal(t, len(v2.blockSigs.Sigs), 2)
// 	assert.Equal(t, v2.blockSigs.Sigs[uint64(1)], sb.Sig)
// 	assert.Equal(t, v2.blockchain.Blocks[1], sb.Block)
// 	assert.Equal(t, len(v2.Unconfirmed.Txns), 0)
// }

// func TestGetSignedBlocksSince(t *testing.T) {
// 	defer cleanupVisor()
// 	cleanupVisor()
// 	vc := newMasterVisorConfig(t)
// 	v := NewVisor(vc)

// 	// No blocks
// 	sbs := v.GetSignedBlocksSince(0, 10)
// 	assert.Equal(t, len(sbs), 0)

// 	// All available blocks
// 	addSignedBlocks(t, v, 10)
// 	sbs = v.GetSignedBlocksSince(2, 4)
// 	assertSignedBlocks(t, v, sbs, 2, 4)

// 	// No available blocks
// 	sbs = v.GetSignedBlocksSince(100, 20)
// 	assert.Equal(t, len(sbs), 0)

// 	// Some, but not all
// 	sbs = v.GetSignedBlocksSince(7, 5)
// 	assertSignedBlocks(t, v, sbs, 7, 5)
// }

// func TestGetGenesisBlock(t *testing.T) {
// 	defer cleanupVisor()
// 	vc := newMasterVisorConfig(t)

// 	// Panics with no signed genesis block
// 	v := NewMinimalVisor(vc)
// 	assert.Panics(t, func() { v.GetGenesisBlock() })

// 	// Panics with no blocks
// 	v = NewMinimalVisor(vc)
// 	v.blockSigs.Sigs[0] = cipher.Sig{}
// 	assert.Panics(t, func() { v.GetGenesisBlock() })

// 	// Correct result
// 	v = NewVisor(vc)
// 	gb := v.GetGenesisBlock()
// 	assert.Equal(t, v.blockSigs.Sigs[0], gb.Sig)
// 	assert.Equal(t, v.blockchain.Blocks[0], gb.Block)
// }

// func TestHeadBkSeq(t *testing.T) {
// 	defer cleanupVisor()
// 	vc := newMasterVisorConfig(t)
// 	v := NewVisor(vc)
// 	assert.Equal(t, v.HeadBkSeq(), uint64(0))
// 	addSignedBlocks(t, v, 10)
// 	assert.Equal(t, v.HeadBkSeq(), uint64(10))
// 	v = NewMinimalVisor(vc)
// 	assert.Panics(t, func() { v.HeadBkSeq() })
// }

// func TestGetBlockchainMetadata(t *testing.T) {
// 	defer cleanupVisor()
// 	vc := newMasterVisorConfig(t)
// 	v := NewVisor(vc)
// 	addSignedBlocks(t, v, 8)
// 	addUnconfirmedTxn(v)
// 	addUnconfirmedTxn(v)
// 	bcm := v.GetBlockchainMetadata()
// 	assert.Equal(t, bcm.Unspents, uint64(9))
// 	assert.Equal(t, bcm.Unconfirmed, uint64(2))
// 	assertReadableBlockHeader(t, bcm.Head, v.blockchain.Head().Head)
// }

// func TestGetReadableBlock(t *testing.T) {
// 	defer cleanupVisor()
// 	vc := newMasterVisorConfig(t)
// 	v := NewVisor(vc)

// 	rb, err := v.GetReadableBlock(1)
// 	assert.NotNil(t, err)
// 	sb := addSignedBlock(t, v)
// 	rb, err = v.GetReadableBlock(1)
// 	assert.Nil(t, err)
// 	assertReadableBlock(t, rb, sb.Block)
// }

// func TestGetReadableBlocks(t *testing.T) {
// 	defer cleanupVisor()
// 	vc := newMasterVisorConfig(t)
// 	v := NewVisor(vc)

// 	rbs := v.GetReadableBlocks(1, 10)
// 	assert.Equal(t, len(rbs), 0)
// 	rbs = v.GetReadableBlocks(0, 10)
// 	sbs := []SignedBlock{SignedBlock{
// 		Sig:   v.blockSigs.Sigs[0],
// 		Block: v.blockchain.Blocks[0],
// 	}}
// 	assertReadableBlocks(t, v, rbs, sbs)
// 	sbs = append(sbs, addSignedBlocks(t, v, 5)...)
// 	rbs = v.GetReadableBlocks(0, 10)
// 	assertReadableBlocks(t, v, rbs, sbs)
// 	rbs = v.GetReadableBlocks(2, 4)
// 	sbs = sbs[2:4]
// 	assertReadableBlocks(t, v, rbs, sbs)
// }

// func TestGetBlock(t *testing.T) {
// 	defer cleanupVisor()
// 	vc := newMasterVisorConfig(t)
// 	v := NewVisor(vc)

// 	b, err := v.GetBlock(1)
// 	assert.NotNil(t, err)
// 	sb := addSignedBlock(t, v)
// 	b, err = v.GetBlock(1)
// 	assert.Nil(t, err)
// 	assert.Equal(t, b, sb.Block)
// }

// func TestGetBlocks(t *testing.T) {
// 	defer cleanupVisor()
// 	vc := newMasterVisorConfig(t)
// 	v := NewVisor(vc)

// 	bs := v.GetBlocks(1, 10)
// 	assert.Equal(t, len(bs), 0)
// 	bs = v.GetBlocks(0, 10)
// 	sbs := []SignedBlock{SignedBlock{
// 		Sig:   v.blockSigs.Sigs[0],
// 		Block: v.blockchain.Blocks[0],
// 	}}
// 	assertBlocks(t, v, bs, sbs)
// 	sbs = append(sbs, addSignedBlocks(t, v, 5)...)
// 	bs = v.GetBlocks(0, 10)
// 	assertBlocks(t, v, bs, sbs)
// 	bs = v.GetBlocks(2, 4)
// 	sbs = sbs[2:4]
// 	assertBlocks(t, v, bs, sbs)
// }

// /*
// func TestVisorSetAnnounced(t *testing.T) {
// 	defer cleanupVisor()
// 	vc := newMasterVisorConfig(t)
// 	v := NewVisor(vc)

// 	now := utc.Now()
// 	utx := addUnconfirmedTxn(v)
// 	assert.True(t, utx.Announced.IsZero())
// 	assert.True(t, v.Unconfirmed.Txns[utx.Hash()].Announced.IsZero())
// 	v.SetAnnounced(utx.Hash(), now)
// 	assert.False(t, v.Unconfirmed.Txns[utx.Hash()].Announced.IsZero())
// 	assert.Equal(t, v.Unconfirmed.Txns[utx.Hash()].Announced, now)
// }
// */

// func TestVisorInjectTxn(t *testing.T) {
// 	defer cleanupVisor()
// 	vc := newMasterVisorConfig(t)
// 	v := NewVisor(vc)

// 	// Setup txns
// 	tx, err := makeValidTxn(v)
// 	assert.Nil(t, err)
// 	we := v.Wallets[0].CreateEntry()
// 	tx2, err := v.Spend(v.Wallets[0].GetFilename(), wallet.Balance{1e6, 0}, 0, we.Address)
// 	assert.Nil(t, err)

// 	// Valid record, did not announce
// 	assert.Equal(t, len(v.Unconfirmed.Txns), 0)
// 	err, known := v.InjectTxn(tx)
// 	assert.Nil(t, err)
// 	assert.False(t, known)
// 	assert.Equal(t, len(v.Unconfirmed.Txns), 1)
// 	assert.True(t, v.Unconfirmed.Txns[tx.Hash()].Announced.IsZero())

// 	// Invalid txn
// 	tx.Out = make([]coin.TransactionOutput, 0)
// 	err, known = v.InjectTxn(tx)
// 	assert.NotNil(t, err)
// 	assert.False(t, known)
// 	assert.Equal(t, len(v.Unconfirmed.Txns), 1)
// 	assert.True(t, v.Unconfirmed.Txns[tx.Hash()].Announced.IsZero())

// 	// Make sure isOurSpend and isOurReceive is correct
// 	tx = tx2
// 	err, known = v.InjectTxn(tx)
// 	assert.Nil(t, err)
// 	assert.False(t, known)
// 	assert.Equal(t, len(v.Unconfirmed.Txns), 2)
// 	assert.True(t, v.Unconfirmed.Txns[tx.Hash()].Announced.IsZero())
// }

// func TestGetAddressTransactions(t *testing.T) {
// 	defer cleanupVisor()
// 	vc := newMasterVisorConfig(t)
// 	v := NewVisor(vc)

// 	// An confirmed txn
// 	w := v.Wallets[0]
// 	we := w.CreateEntry()
// 	tx, err := v.Spend(w.GetFilename(), wallet.Balance{1e6, 0}, 0, we.Address)
// 	assert.Nil(t, err)
// 	err, known := v.InjectTxn(tx)
// 	assert.Nil(t, err)
// 	assert.False(t, known)
// 	assert.Equal(t, len(v.Unconfirmed.Txns), 1)
// 	_, err = v.CreateAndExecuteBlock()
// 	assert.Nil(t, err)
// 	assert.Equal(t, len(v.Unconfirmed.Txns), 0)
// 	txns := v.GetAddressTransactions(we.Address)
// 	assert.Equal(t, len(txns), 1)
// 	assert.Equal(t, txns[0].Txn, tx)
// 	assert.True(t, txns[0].Status.Confirmed)
// 	assert.Equal(t, txns[0].Status.Height, uint64(1))

// 	// An unconfirmed txn
// 	assert.Equal(t, len(v.Unconfirmed.Txns), 0)
// 	assert.Equal(t, len(v.Unconfirmed.Unspent), 0)
// 	we = w.CreateEntry()
// 	tx, err = v.Spend(w.GetFilename(), wallet.Balance{2e6, 0}, 0, we.Address)
// 	err, known = v.InjectTxn(tx)
// 	assert.Nil(t, err)
// 	assert.False(t, known)
// 	assert.Equal(t, len(v.Unconfirmed.Txns), 1)
// 	assert.Equal(t, len(v.Unconfirmed.Unspent), 1)
// 	assert.Equal(t, len(v.Unconfirmed.Unspent[tx.Hash()]), 2)
// 	found := false
// 	for _, uxs := range v.Unconfirmed.Unspent {
// 		if found {
// 			break
// 		}
// 		for _, ux := range uxs {
// 			if ux.Body.Address == we.Address {
// 				found = true
// 				break
// 			}
// 		}
// 	}
// 	auxs := v.Unconfirmed.Unspent.AllForAddress(we.Address)
// 	assert.Equal(t, len(auxs), 1)
// 	assert.True(t, found)
// 	txns = v.GetAddressTransactions(we.Address)
// 	assert.Equal(t, len(txns), 1)
// 	assert.Equal(t, txns[0].Txn, tx)
// 	assert.True(t, txns[0].Status.Unconfirmed)

// 	// An unconfirmed txn, but pool is corrupted
// 	assert.True(t, len(v.Unconfirmed.Unspent) > 0)
// 	ux := coin.UxOut{}
// 	found = false
// 	for _, uxs := range v.Unconfirmed.Unspent {
// 		if len(uxs) > 0 {
// 			ux = uxs[0]
// 			found = true
// 			break
// 		}
// 	}
// 	assert.True(t, found)
// 	srcTxn := ux.Body.SrcTransaction
// 	delete(v.Unconfirmed.Txns, srcTxn)
// 	txns = v.GetAddressTransactions(we.Address)
// 	assert.Equal(t, len(txns), 0)
// }

// func TestGetTransaction(t *testing.T) {
// 	defer cleanupVisor()
// 	vc := newMasterVisorConfig(t)
// 	v := NewVisor(vc)

// 	// Unknown
// 	tx, err := makeValidTxn(v)
// 	assert.Nil(t, err)
// 	tx2 := v.GetTransaction(tx.Hash())
// 	assert.True(t, tx2.Status.Unknown)

// 	// Unconfirmed
// 	err, known := v.InjectTxn(tx)
// 	assert.Nil(t, err)
// 	assert.False(t, known)
// 	tx2 = v.GetTransaction(tx.Hash())
// 	assert.True(t, tx2.Status.Unconfirmed)
// 	assert.Equal(t, tx, tx2.Txn)

// 	// Confirmed
// 	_, err = v.CreateAndExecuteBlock()
// 	assert.Nil(t, err)
// 	assert.Equal(t, len(v.Unconfirmed.Txns), 0)
// 	tx2 = v.GetTransaction(tx.Hash())
// 	assert.True(t, tx2.Status.Confirmed)
// 	assert.Equal(t, tx2.Status.Height, uint64(1))
// 	assert.Equal(t, tx, tx2.Txn)
// }

// func TestBalances(t *testing.T) {
// 	defer cleanupVisor()
// 	v, mv := setupVisor()
// 	w := v.Wallets[0]
// 	assert.Equal(t, len(w.GetEntries()), 1)
// 	we := w.CreateEntry()
// 	we2 := w.CreateEntry()
// 	assert.Equal(t, len(w.GetEntries()), 3)
// 	assert.Equal(t, v.TotalBalance().Confirmed, wallet.Balance{0, 0})
// 	startCoins := mv.Config.GenesisCoinVolume

// 	// Without predicted outputs
// 	assert.Nil(t,
// 		transferCoinsAdvanced(mv, v, wallet.Balance{10e6, 10}, 0, we.Address))
// 	assert.Nil(t,
// 		transferCoinsAdvanced(mv, v, wallet.Balance{10e6, 10}, 0, we.Address))
// 	assert.Nil(t,
// 		transferCoinsAdvanced(mv, v, wallet.Balance{5e6, 5}, 0, we2.Address))
// 	assert.Equal(t, v.WalletBalance(w.GetFilename()).Confirmed, wallet.Balance{25e6, 25})
// 	assert.Equal(t, v.AddressBalance(we.Address).Confirmed, wallet.Balance{20e6, 20})
// 	assert.Equal(t, v.AddressBalance(we2.Address).Confirmed, wallet.Balance{5e6, 5})
// 	assert.Equal(t, v.TotalBalance().Confirmed, wallet.Balance{25e6, 25})
// 	mvBalance := wallet.Balance{startCoins - 25e6, startCoins - 25}
// 	assert.Equal(t, mv.TotalBalance().Confirmed, mvBalance)
// 	assert.Equal(t, v.AddressBalance(we.Address).Confirmed, wallet.Balance{20e6, 20})
// 	assert.Equal(t, v.AddressBalance(we2.Address).Confirmed, wallet.Balance{5e6, 5})

// 	// TODO -- test the predicted balances
// }

// func TestVisorVerifySignedBlock(t *testing.T) {
// 	defer cleanupVisor()
// 	vc := newMasterVisorConfig(t)
// 	v := NewVisor(vc)
// 	w := v.Wallets[0]
// 	we := w.CreateEntry()

// 	// Master should verify its own blocks correctly
// 	txn, err := v.Spend(w.GetFilename(), wallet.Balance{1e6, 0}, 0, we.Address)
// 	assert.Nil(t, err)
// 	err, known := v.InjectTxn(txn)
// 	assert.Nil(t, err)
// 	assert.False(t, known)
// 	b, err := v.CreateBlock(uint64(utc.UnixNow()))
// 	assert.Nil(t, err)
// 	assert.Nil(t, v.verifySignedBlock(&b))
// 	badb := b
// 	badb.Sig = cipher.Sig{}
// 	assert.NotNil(t, v.verifySignedBlock(&badb))

// 	// Non master should verify signed blocks generated by master
// 	mv := v
// 	v = setupVisorFromMaster(mv)
// 	assert.Nil(t, v.verifySignedBlock(&b))
// 	assert.NotNil(t, v.verifySignedBlock(&badb))
// }

// func TestVisorSignBlock(t *testing.T) {
// 	defer cleanupVisor()
// 	vc := newMasterVisorConfig(t)
// 	v := NewVisor(vc)

// 	// Non master should panic
// 	b := v.blockchain.Blocks[0]
// 	v.Config.IsMaster = false
// 	assert.Panics(t, func() { v.SignBlock(b) })

// 	// Master should generate valid signed block
// 	v.Config.IsMaster = true
// 	sb := v.SignBlock(b)
// 	assert.Nil(t, v.verifySignedBlock(&sb))
// }

// func TestCreateMasterWallet(t *testing.T) {
// 	defer cleanupVisor()
// 	cleanupVisor()
// 	we := wallet.NewWalletEntry()
// 	w := CreateMasterWallet(we)
// 	assert.Equal(t, w.NumEntries(), 1)
// 	assert.Equal(t, w.GetAddresses()[0], we.Address)

// 	// Having a wallet file present should not affect loading master wallet
// 	w.Save(testWalletFile)
// 	we = wallet.NewWalletEntry()
// 	w = CreateMasterWallet(we)
// 	assert.Equal(t, w.NumEntries(), 1)
// 	assert.Equal(t, w.GetAddresses()[0], we.Address)

// 	// Creating with an invalid wallet entry should panic
// 	we = wallet.NewWalletEntry()
// 	we.Secret = cipher.SecKey{}
// 	assert.Panics(t, func() { CreateMasterWallet(we) })
// 	we = wallet.NewWalletEntry()
// 	we.Public = cipher.PubKey{}
// 	assert.Panics(t, func() { CreateMasterWallet(we) })
// }
