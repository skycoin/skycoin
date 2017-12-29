package webrpc

import (
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/visor/historydb"
)

//go:generate goautomock -template=testify Gatewayer

// Gatewayer provides interfaces for getting skycoin related info.
type Gatewayer interface {
	GetLastBlocks(num uint64) (*visor.ReadableBlocks, error)
	GetBlocks(start, end uint64) (*visor.ReadableBlocks, error)
	GetBlocksInDepth(vs []uint64) (*visor.ReadableBlocks, error)
	GetUnspentOutputs(filters ...daemon.OutputsFilter) (visor.ReadableOutputSet, error)
	GetTransaction(txid cipher.SHA256) (*visor.Transaction, error)
	InjectTransaction(tx coin.Transaction) error
	GetAddrUxOuts(addr cipher.Address) ([]*historydb.UxOutJSON, error)
	GetTimeNow() uint64
}
