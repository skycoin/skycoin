package daemon

import (
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/wallet"
)

//go:generate goautomock -template=testify Visorer
type Visorer interface {
	GetVisor() visor.Visorer
	GetConfig() VisorConfig
	Run() error
	Shutdown()
	RefreshUnconfirmed() ([]cipher.SHA256, error)
	RemoveInvalidUnconfirmed() ([]cipher.SHA256, error)
	RequestBlocks(pool *Pool) error
	AnnounceBlocks(pool *Pool) error
	AnnounceAllTxns(pool *Pool) error
	AnnounceTxns(pool *Pool, txns []cipher.SHA256) error
	RequestBlocksFromAddr(pool *Pool, addr string) error
	SetTxnsAnnounced(txns []cipher.SHA256)
	InjectBroadcastTransaction(txn coin.Transaction, pool *Pool) error
	InjectTransaction(tx coin.Transaction) (bool, *visor.ErrTxnViolatesSoftConstraint, error)
	broadcastBlock(sb coin.SignedBlock, pool *Pool) error
	broadcastTransaction(t coin.Transaction, pool *Pool) error
	ResendTransaction(h cipher.SHA256, pool *Pool) error
	ResendUnconfirmedTxns(pool *Pool) []cipher.SHA256
	CreateAndPublishBlock(pool *Pool) (coin.SignedBlock, error)
	RemoveConnection(addr string)
	RecordBlockchainHeight(addr string, bkLen uint64)
	EstimateBlockchainHeight() uint64
	ScanAheadWalletAddresses(wltName string, scanN uint64) (wallet.Wallet, error)
	GetPeerBlockchainHeights() []PeerBlockchainHeight
	HeadBkSeq() uint64
	ExecuteSignedBlock(b coin.SignedBlock) error
	GetSignedBlock(seq uint64) (*coin.SignedBlock, error)
	GetSignedBlocksSince(seq uint64, ct uint64) ([]coin.SignedBlock, error)
	UnConfirmFilterKnown(txns []cipher.SHA256) []cipher.SHA256
	UnConfirmKnow(hashes []cipher.SHA256) coin.Transactions
}
