package visor

import (
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/visor/historydb"
	"github.com/skycoin/skycoin/src/wallet"
)

//go:generate goautomock -template=testify Visorer
type Visorer interface {
	GetConfig() Config
	Wallets() *wallet.Service
	GetUnconfirmed() UnconfirmedTxnPooler
	GetBlockchain() Blockchainer
	Run() error
	Shutdown()
	maybeCreateGenesisBlock() error
	GenesisPreconditions()
	RefreshUnconfirmed() ([]cipher.SHA256, error)
	RemoveInvalidUnconfirmed() ([]cipher.SHA256, error)
	CreateBlock(when uint64) (coin.SignedBlock, error)
	CreateAndExecuteBlock() (coin.SignedBlock, error)
	ExecuteSignedBlock(b coin.SignedBlock) error
	SignBlock(b coin.Block) coin.SignedBlock
	GetUnspentOutputs() ([]coin.UxOut, error)
	UnconfirmedSpendingOutputs() (coin.UxArray, error)
	UnconfirmedIncomingOutputs() (coin.UxArray, error)
	GetSignedBlocksSince(seq, ct uint64) ([]coin.SignedBlock, error)
	HeadBkSeq() uint64
	GetBlockchainMetadata() BlockchainMetadata
	GetBlock(seq uint64) (*coin.SignedBlock, error)
	GetBlocks(start, end uint64) []coin.SignedBlock
	InjectTransaction(txn coin.Transaction) (bool, *ErrTxnViolatesSoftConstraint, error)
	InjectTransactionStrict(txn coin.Transaction) (bool, error)
	GetAddressTxns(a cipher.Address) ([]Transaction, error)
	GetTransaction(txHash cipher.SHA256) (*Transaction, error)
	GetTransactions(flts ...TxFilter) ([]Transaction, error)
	getTransactionsOfAddrs(addrs []cipher.Address) (map[cipher.Address][]Transaction, error)
	traverseTxns(flts ...TxFilter) ([]Transaction, error)
	AddressBalance(auxs coin.AddressUxOuts) (uint64, uint64, error)
	GetUnconfirmedTxns(filter func(UnconfirmedTxn) bool) []UnconfirmedTxn
	GetAllUnconfirmedTxns() []UnconfirmedTxn
	GetAllValidUnconfirmedTxHashes() []cipher.SHA256
	GetBlockByHash(hash cipher.SHA256) (*coin.SignedBlock, error)
	GetBlockBySeq(seq uint64) (*coin.SignedBlock, error)
	GetLastBlocks(num uint64) []coin.SignedBlock
	GetLastTxs() ([]*Transaction, error)
	GetHeadBlock() (*coin.SignedBlock, error)
	GetUxOutByID(id cipher.SHA256) (*historydb.UxOut, error)
	GetAddrUxOuts(address cipher.Address) ([]*historydb.UxOut, error)
	ScanAheadWalletAddresses(wltName string, scanN uint64) (wallet.Wallet, error)
	GetBalanceOfAddrs(addrs []cipher.Address) ([]wallet.BalancePair, error)
}
