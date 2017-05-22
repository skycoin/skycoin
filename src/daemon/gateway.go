package daemon

import (
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/wallet"
	//"github.com/skycoin/skycoin/src/wallet"
	"time"

	"github.com/skycoin/skycoin/src/visor/historydb"
)

// Exposes a read-only api for use by the gui rpc interface

// GatewayConfig configuration set of gateway.
type GatewayConfig struct {
	BufferSize int
}

// NewGatewayConfig create and init an GatewayConfig
func NewGatewayConfig() GatewayConfig {
	return GatewayConfig{
		BufferSize: 32,
	}
}

// Gateway RPC interface wrapper for daemon state
type Gateway struct {
	Config GatewayConfig
	drpc   RPC
	vrpc   visor.RPC

	// Backref to Daemon
	d *Daemon
	// Backref to Visor
	v *visor.Visor
	// Requests are queued on this channel
	requests chan func()
	// When a request is done processing, it is placed on this channel
	// Responses chan interface{}
}

// NewGateway create and init an Gateway instance.
func NewGateway(c GatewayConfig, D *Daemon) *Gateway {
	return &Gateway{
		Config:   c,
		drpc:     RPC{},
		vrpc:     visor.RPC{},
		d:        D,
		v:        D.Visor.v,
		requests: make(chan func(), c.BufferSize),
	}
}

func (gw *Gateway) strand(f func()) {
	done := make(chan struct{})
	gw.requests <- func() {
		defer close(done)
		f()
	}
	<-done
}

// GetConnections returns a *Connections
func (gw *Gateway) GetConnections() interface{} {
	var conns interface{}
	gw.strand(func() {
		conns = gw.drpc.GetConnections(gw.d)
	})
	return conns
}

// GetDefaultConnections returns default connections
func (gw *Gateway) GetDefaultConnections() interface{} {
	var conns interface{}
	gw.strand(func() {
		conns = gw.drpc.GetDefaultConnections(gw.d)
	})
	return conns
}

// GetConnection returns a *Connection of specific address
func (gw *Gateway) GetConnection(addr string) interface{} {
	logger.Critical("here")
	var conn interface{}
	gw.strand(func() {
		conn = gw.drpc.GetConnection(gw.d, addr)
	})
	return conn
}

// GetTrustConnections returns all trusted connections,
// including private and public
func (gw *Gateway) GetTrustConnections() interface{} {
	var conn interface{}
	gw.strand(func() {
		conn = gw.drpc.GetTrustConnections(gw.d)
	})
	return conn
}

// GetExchgConnection returns all exchangeable connections,
// including private and public
func (gw *Gateway) GetExchgConnection() interface{} {
	var conn interface{}
	gw.strand(func() {
		conn = gw.drpc.GetAllExchgConnections(gw.d)
	})
	return conn
}

/* Blockchain & Transaction status */
//DEPRECATE

// GetBlockchainProgress returns a *BlockchainProgress
func (gw *Gateway) GetBlockchainProgress() interface{} {
	var bcp interface{}
	gw.strand(func() {
		bcp = gw.drpc.GetBlockchainProgress(gw.d.Visor)
	})
	return bcp
}

// ResendTransaction resent the transaction and return a *ResendResult
func (gw *Gateway) ResendTransaction(txn cipher.SHA256) interface{} {
	var result interface{}
	gw.strand(func() {
		result = gw.drpc.ResendTransaction(gw.d.Visor, gw.d.Pool, txn)
	})
	return result
}

// ResendUnconfirmedTxns resents all unconfirmed transactions
func (gw *Gateway) ResendUnconfirmedTxns() (rlt *ResendResult) {
	gw.strand(func() {
		rlt = gw.drpc.ResendUnconfirmedTxns(gw.d.Visor, gw.d.Pool)
	})
	return
}

// GetBlockchainMetadata returns a *visor.BlockchainMetadata
func (gw *Gateway) GetBlockchainMetadata() interface{} {
	var bcm interface{}
	gw.strand(func() {
		bcm = gw.vrpc.GetBlockchainMetadata(gw.v)
	})
	return bcm
}

// GetBlockByHash returns the block by hash
func (gw *Gateway) GetBlockByHash(hash cipher.SHA256) (block coin.Block, ok bool) {
	gw.strand(func() {
		b := gw.v.GetBlockByHash(hash)
		if b == nil {
			return
		}
		block = *b
		ok = true
	})
	return
}

// GetBlockBySeq returns blcok by seq
func (gw *Gateway) GetBlockBySeq(seq uint64) (block coin.Block, ok bool) {
	gw.strand(func() {
		b := gw.v.GetBlockBySeq(seq)
		if b == nil {
			return
		}
		block = *b
		ok = true
	})
	return
}

// GetBlocks returns a *visor.ReadableBlocks
func (gw *Gateway) GetBlocks(start, end uint64) *visor.ReadableBlocks {
	var blocks *visor.ReadableBlocks
	gw.strand(func() {
		blocks = gw.vrpc.GetBlocks(gw.v, start, end)
	})
	return blocks
}

// GetBlocksInDepth returns blocks in different depth
func (gw *Gateway) GetBlocksInDepth(vs []uint64) *visor.ReadableBlocks {
	var blocks *visor.ReadableBlocks
	gw.strand(func() {
		blks := visor.ReadableBlocks{}
		for _, n := range vs {
			if b := gw.vrpc.GetBlockInDepth(gw.v, n); b != nil {
				blks.Blocks = append(blks.Blocks, *b)
			}
		}
		blocks = &blks
	})
	return blocks
}

// GetLastBlocks get last N blocks
func (gw *Gateway) GetLastBlocks(num uint64) *visor.ReadableBlocks {
	var blocks *visor.ReadableBlocks
	gw.strand(func() {
		headSeq := gw.v.HeadBkSeq()
		var start uint64
		if (headSeq + 1) > num {
			start = headSeq - num + 1
		}

		blocks = gw.vrpc.GetBlocks(gw.v, start, headSeq)
	})
	return blocks
}

// OutputsFilter used as optional arguments in GetUnspentOutputs method
type OutputsFilter func(outputs []visor.ReadableOutput) []visor.ReadableOutput

// GetUnspentOutputs gets unspent outputs and returns the filtered results,
// Note: all filters will be executed as the pending sequence in 'AND' mode.
func (gw *Gateway) GetUnspentOutputs(filters ...OutputsFilter) visor.ReadableOutputSet {
	var allOutputs []visor.ReadableOutput
	var spendingOutputs []visor.ReadableOutput
	var inOutputs []visor.ReadableOutput
	gw.strand(func() {
		allOutputs = gw.v.GetUnspentOutputReadables()
		spendingOutputs = gw.v.AllSpendsOutputs()
		inOutputs = gw.v.AllIncommingOutputs()
	})

	for _, flt := range filters {
		allOutputs = flt(allOutputs)
		spendingOutputs = flt(spendingOutputs)
		inOutputs = flt(inOutputs)
	}

	return visor.ReadableOutputSet{
		HeadOutputs:      allOutputs,
		OutgoingOutputs:  spendingOutputs,
		IncommingOutputs: inOutputs,
	}
}

// FbyAddressesNotIncluded filters the unspent outputs that are not owned by the addresses
func FbyAddressesNotIncluded(addrs []string) OutputsFilter {
	return func(outputs []visor.ReadableOutput) []visor.ReadableOutput {
		addrMatch := []visor.ReadableOutput{}
		addrMap := make(map[string]bool)
		for _, addr := range addrs {
			addrMap[addr] = false
		}

		for _, u := range outputs {
			_, ok := addrMap[u.Address]
			if !ok {
				addrMatch = append(addrMatch, u)
			}
		}
		return addrMatch
	}
}

// FbyAddresses filters the unspent outputs that owned by the addresses
func FbyAddresses(addrs []string) OutputsFilter {
	return func(outputs []visor.ReadableOutput) []visor.ReadableOutput {
		addrMatch := []visor.ReadableOutput{}
		addrMap := make(map[string]bool)
		for _, addr := range addrs {
			addrMap[addr] = true
		}

		for _, u := range outputs {
			if _, ok := addrMap[u.Address]; ok {
				addrMatch = append(addrMatch, u)
			}
		}
		return addrMatch
	}
}

// FbyHashes filters the unspent outputs that have hashes matched.
func FbyHashes(hashes []string) OutputsFilter {
	return func(outputs []visor.ReadableOutput) []visor.ReadableOutput {
		hsMatch := []visor.ReadableOutput{}
		hsMap := make(map[string]bool)
		for _, h := range hashes {
			hsMap[h] = true
		}

		for _, u := range outputs {
			if _, ok := hsMap[u.Hash]; ok {
				hsMatch = append(hsMatch, u)
			}
		}
		return hsMatch
	}
}

// GetTransaction returns transaction by txid
func (gw *Gateway) GetTransaction(txid cipher.SHA256) (tx *visor.Transaction, err error) {
	gw.strand(func() {
		tx, err = gw.v.GetTransaction(txid)
	})
	return
}

// GetTransactionResult gets transaction result by txid.
func (gw *Gateway) GetTransactionResult(txid cipher.SHA256) (*visor.TransactionResult, error) {
	var tx *visor.TransactionResult
	var err error
	gw.strand(func() {
		tx, err = gw.vrpc.GetTransaction(gw.v, txid)
	})
	return tx, err
}

// InjectTransaction injects transaction
func (gw *Gateway) InjectTransaction(txn coin.Transaction) (tx coin.Transaction, err error) {
	gw.strand(func() {
		tx, err = gw.d.Visor.InjectTransaction(txn, gw.d.Pool)
	})
	return
}

// GetAddressTransactions returns a *visor.TransactionResults
func (gw *Gateway) GetAddressTransactions(a cipher.Address) interface{} {
	var tx interface{}
	gw.strand(func() {
		tx = gw.vrpc.GetAddressTransactions(gw.v, a)
	})
	return tx
}

// GetUxOutByID gets UxOut by hash id.
func (gw *Gateway) GetUxOutByID(id cipher.SHA256) (*historydb.UxOut, error) {
	var uxout *historydb.UxOut
	var err error
	gw.strand(func() {
		uxout, err = gw.v.GetUxOutByID(id)
	})
	return uxout, err
}

// GetAddrUxOuts gets all the address affected UxOuts.
func (gw *Gateway) GetAddrUxOuts(addr cipher.Address) ([]*historydb.UxOutJSON, error) {
	var (
		uxouts []*historydb.UxOut
		err    error
	)
	gw.strand(func() {
		uxouts, err = gw.v.GetAddrUxOuts(addr)
	})
	uxs := make([]*historydb.UxOutJSON, len(uxouts))
	for i, ux := range uxouts {
		uxs[i] = historydb.NewUxOutJSON(ux)
	}
	return uxs, err
}

// GetAddressUxOuts gets all the address affected UxOuts.
func (gw *Gateway) GetAddressUxOuts(addr cipher.Address) ([]*historydb.UxOut, error) {
	var (
		uxouts []*historydb.UxOut
		err    error
	)
	gw.strand(func() {
		uxouts, err = gw.v.GetAddrUxOuts(addr)
	})
	return uxouts, err
}

// GetTimeNow returns the current Unix time
func (gw *Gateway) GetTimeNow() uint64 {
	return uint64(time.Now().Unix())
}

// GetAllUnconfirmedTxns returns all unconfirmed transactions
func (gw *Gateway) GetAllUnconfirmedTxns() (txns []visor.UnconfirmedTxn) {
	gw.strand(func() {
		txns = gw.v.GetAllUnconfirmedTxns()
	})
	return
}

// GetUnconfirmedTxns returns addresses related unconfirmed transactions
func (gw *Gateway) GetUnconfirmedTxns(addrs []cipher.Address) (txns []visor.UnconfirmedTxn) {
	gw.strand(func() {
		txns = gw.v.GetUnconfirmedTxns(visor.ToAddresses(addrs))
	})
	return
}

// GetLastTxs returns last confirmed transactions, return nil if empty
func (gw *Gateway) GetLastTxs() (txns []*visor.Transaction, err error) {
	gw.strand(func() {
		txns, err = gw.v.GetLastTxs()
	})
	return
}

// GetUnspent returns the unspent pool
func (gw *Gateway) GetUnspent() (unspent coin.UnspentPool) {
	gw.strand(func() {
		unspent = gw.vrpc.GetUnspent(gw.v)
	})
	return
}

// CreateSpendingTransaction creates spending transactions
func (gw *Gateway) CreateSpendingTransaction(wlt wallet.Wallet,
	amt wallet.Balance,
	dest cipher.Address) (tx coin.Transaction, err error) {
	gw.strand(func() {
		tx, err = gw.vrpc.CreateSpendingTransaction(gw.v, wlt, amt, dest)
	})
	return
}

// WalletBalance returns balance pair of specific wallet
func (gw *Gateway) WalletBalance(wlt wallet.Wallet) (balance wallet.BalancePair) {
	gw.strand(func() {
		auxs := gw.vrpc.GetUnspent(gw.v).AllForAddresses(wlt.GetAddresses())
		puxs := gw.vrpc.GetUnconfirmedSpends(gw.v, wlt.GetAddressSet())

		coins1, hours1 := gw.v.AddressBalance(auxs)
		coins2, hours2 := gw.v.AddressBalance(auxs.Sub(puxs))
		balance = wallet.BalancePair{
			Confirmed: wallet.Balance{Coins: coins1, Hours: hours1},
			Predicted: wallet.Balance{Coins: coins2, Hours: hours2},
		}
	})
	return
}

// AddressesBalance gets balance of given addresses
func (gw *Gateway) AddressesBalance(addrs []cipher.Address) (balance wallet.BalancePair) {
	addrMap := make(map[cipher.Address]byte)
	gw.strand(func() {
		for _, a := range addrs {
			addrMap[a] = byte(1)
		}

		auxs := gw.vrpc.GetUnspent(gw.v).AllForAddresses(addrs)
		puxs := gw.vrpc.GetUnconfirmedSpends(gw.v, addrMap)

		coins1, hours1 := gw.v.AddressBalance(auxs)
		coins2, hours2 := gw.v.AddressBalance(auxs.Sub(puxs))
		balance = wallet.BalancePair{
			Confirmed: wallet.Balance{Coins: coins1, Hours: hours1},
			Predicted: wallet.Balance{Coins: coins2, Hours: hours2},
		}
	})
	return
}
