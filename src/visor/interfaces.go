package visor

import (
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/params"
	"github.com/skycoin/skycoin/src/visor/blockdb"
	"github.com/skycoin/skycoin/src/visor/dbutil"
	"github.com/skycoin/skycoin/src/visor/historydb"
)

//go:generate mockery -name Historyer -case underscore -inpkg -testonly
//go:generate mockery -name Blockchainer -case underscore -inpkg -testonly
//go:generate mockery -name UnconfirmedTransactionPooler -case underscore -inpkg -testonly

// Historyer is the interface that provides methods for accessing history data that are parsed from blockchain.
type Historyer interface {
	GetUxOuts(tx *dbutil.Tx, uxids []cipher.SHA256) ([]historydb.UxOut, error)
	ParseBlock(tx *dbutil.Tx, b coin.Block) error
	GetTransaction(tx *dbutil.Tx, hash cipher.SHA256) (*historydb.Transaction, error)
	GetOutputsForAddress(tx *dbutil.Tx, address cipher.Address) ([]historydb.UxOut, error)
	GetTransactionHashesForAddresses(tx *dbutil.Tx, addresses []cipher.Address) ([]cipher.SHA256, error)
	AddressSeen(tx *dbutil.Tx, address cipher.Address) (bool, error)
	NeedsReset(tx *dbutil.Tx) (bool, error)
	Erase(tx *dbutil.Tx) error
	ParsedBlockSeq(tx *dbutil.Tx) (uint64, bool, error)
	ForEachTxn(tx *dbutil.Tx, f func(cipher.SHA256, *historydb.Transaction) error) error
}

// Blockchainer is the interface that provides methods for accessing the blockchain data
type Blockchainer interface {
	GetGenesisBlock(tx *dbutil.Tx) (*coin.SignedBlock, error)
	GetBlocks(tx *dbutil.Tx, seqs []uint64) ([]coin.SignedBlock, error)
	GetBlocksInRange(tx *dbutil.Tx, start, end uint64) ([]coin.SignedBlock, error)
	GetLastBlocks(tx *dbutil.Tx, n uint64) ([]coin.SignedBlock, error)
	GetSignedBlockByHash(tx *dbutil.Tx, hash cipher.SHA256) (*coin.SignedBlock, error)
	GetSignedBlockBySeq(tx *dbutil.Tx, seq uint64) (*coin.SignedBlock, error)
	Unspent() blockdb.UnspentPooler
	Len(tx *dbutil.Tx) (uint64, error)
	Head(tx *dbutil.Tx) (*coin.SignedBlock, error)
	HeadSeq(tx *dbutil.Tx) (uint64, bool, error)
	Time(tx *dbutil.Tx) (uint64, error)
	NewBlock(tx *dbutil.Tx, txns coin.Transactions, currentTime uint64) (*coin.Block, error)
	ExecuteBlock(tx *dbutil.Tx, sb *coin.SignedBlock) error
	VerifyBlock(tx *dbutil.Tx, sb *coin.SignedBlock) error
	VerifyBlockTxnConstraints(tx *dbutil.Tx, txn coin.Transaction) error
	VerifySingleTxnHardConstraints(tx *dbutil.Tx, txn coin.Transaction, signed TxnSignedFlag) error
	VerifySingleTxnSoftHardConstraints(tx *dbutil.Tx, txn coin.Transaction, distParams params.Distribution, verifyParams params.VerifyTxn, signed TxnSignedFlag) (*coin.SignedBlock, coin.UxArray, error)
	TransactionFee(tx *dbutil.Tx, hours uint64) coin.FeeCalculator
}

// UnconfirmedTransactionPooler is the interface that provides methods for
// accessing the unconfirmed transaction pool
type UnconfirmedTransactionPooler interface {
	SetTransactionsAnnounced(tx *dbutil.Tx, hashes map[cipher.SHA256]int64) error
	InjectTransaction(tx *dbutil.Tx, bc Blockchainer, t coin.Transaction, distParams params.Distribution, verifyParams params.VerifyTxn) (bool, *ErrTxnViolatesSoftConstraint, error)
	AllRawTransactions(tx *dbutil.Tx) (coin.Transactions, error)
	RemoveTransactions(tx *dbutil.Tx, txns []cipher.SHA256) error
	Refresh(tx *dbutil.Tx, bc Blockchainer, distParams params.Distribution, verifyParams params.VerifyTxn) ([]cipher.SHA256, error)
	RemoveInvalid(tx *dbutil.Tx, bc Blockchainer) ([]cipher.SHA256, error)
	FilterKnown(tx *dbutil.Tx, txns []cipher.SHA256) ([]cipher.SHA256, error)
	GetKnown(tx *dbutil.Tx, txns []cipher.SHA256) (coin.Transactions, error)
	RecvOfAddresses(tx *dbutil.Tx, bh coin.BlockHeader, addrs []cipher.Address) (coin.AddressUxOuts, error)
	GetIncomingOutputs(tx *dbutil.Tx, bh coin.BlockHeader) (coin.UxArray, error)
	Get(tx *dbutil.Tx, hash cipher.SHA256) (*UnconfirmedTransaction, error)
	GetFiltered(tx *dbutil.Tx, filter func(tx UnconfirmedTransaction) bool) ([]UnconfirmedTransaction, error)
	GetHashes(tx *dbutil.Tx, filter func(tx UnconfirmedTransaction) bool) ([]cipher.SHA256, error)
	ForEach(tx *dbutil.Tx, f func(cipher.SHA256, UnconfirmedTransaction) error) error
	GetUnspentsOfAddr(tx *dbutil.Tx, addr cipher.Address) (coin.UxArray, error)
	Len(tx *dbutil.Tx) (uint64, error)
}
