package webrpc

import (
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/visor/historydb"
)

//go:generate goautomock -template=testify Gatewayer

// Gatewayer provides interfaces for getting skycoin related info.
type Gatewayer interface {
	GetLastBlocks(num uint64) *visor.ReadableBlocks
	GetBlocks(start, end uint64) *visor.ReadableBlocks
	GetUnspentByAddrs(addrs []string) []visor.ReadableOutput
	GetUnspentByHashes(hashes []string) []visor.ReadableOutput
	GetTransaction(txid cipher.SHA256) (*visor.TransactionResult, error)
	InjectTransaction(tx coin.Transaction) (coin.Transaction, error)
	GetRecvUxOutOfAddr(addr cipher.Address) ([]*historydb.UxOut, error)
	GetSpentUxOutOfAddr(addr cipher.Address) ([]*historydb.UxOut, error)
}
