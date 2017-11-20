package daemon

import (
	"errors"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/daemon/strand"
	"github.com/skycoin/skycoin/src/util/utc"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/wallet"

	"fmt"

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

	// Backref to Daemon
	d *Daemon
	// Backref to Visor
	v *visor.Visor
	// Requests are queued on this channel
	requests chan strand.Request

	quit chan struct{}
}

// NewGateway create and init an Gateway instance.
func NewGateway(c GatewayConfig, D *Daemon) *Gateway {
	return &Gateway{
		Config:   c,
		drpc:     RPC{},
		d:        D,
		v:        D.Visor.v,
		requests: make(chan strand.Request, c.BufferSize),
		quit:     make(chan struct{}),
	}
}

// Shutdown shuts down the gateway
func (gw *Gateway) Shutdown() {
	close(gw.quit)
}

func (gw *Gateway) strand(name string, f func()) {
	name = fmt.Sprintf("daemon.Gateway.%s", name)
	strand.Strand(logger, gw.requests, name, func() error {
		f()
		return nil
	}, gw.quit, nil)
}

// GetConnections returns a *Connections
func (gw *Gateway) GetConnections() interface{} {
	var conns interface{}
	gw.strand("GetConnections", func() {
		conns = gw.drpc.GetConnections(gw.d)
	})
	return conns
}

// GetDefaultConnections returns default connections
func (gw *Gateway) GetDefaultConnections() interface{} {
	var conns interface{}
	gw.strand("GetDefaultConnections", func() {
		conns = gw.drpc.GetDefaultConnections(gw.d)
	})
	return conns
}

// GetConnection returns a *Connection of specific address
func (gw *Gateway) GetConnection(addr string) interface{} {
	var conn interface{}
	gw.strand("GetConnection", func() {
		conn = gw.drpc.GetConnection(gw.d, addr)
	})
	return conn
}

// GetTrustConnections returns all trusted connections,
// including private and public
func (gw *Gateway) GetTrustConnections() interface{} {
	var conn interface{}
	gw.strand("GetTrustConnections", func() {
		conn = gw.drpc.GetTrustConnections(gw.d)
	})
	return conn
}

// GetExchgConnection returns all exchangeable connections,
// including private and public
func (gw *Gateway) GetExchgConnection() interface{} {
	var conn interface{}
	gw.strand("GetExchgConnection", func() {
		conn = gw.drpc.GetAllExchgConnections(gw.d)
	})
	return conn
}

/* Blockchain & Transaction status */

// GetBlockchainProgress gets the blockchain progress
func (gw *Gateway) GetBlockchainProgress() (*BlockchainProgress, error) {
	var bcp *BlockchainProgress
	var err error

	gw.strand("GetBlockchainProgress", func() {
		if gw.d.Visor.v == nil {
			return
		}

		var headSeq uint64
		headSeq, _, err = gw.d.Visor.HeadBkSeq()
		if err != nil {
			return
		}

		var height uint64
		height, err = gw.d.Visor.EstimateBlockchainHeight()
		if err != nil {
			return
		}

		bcp = &BlockchainProgress{
			Current: headSeq,
			Highest: height,
		}

		peerHeights := gw.d.Visor.GetPeerBlockchainHeights()
		for _, ph := range peerHeights {
			bcp.Peers = append(bcp.Peers, BlockchainPeer{
				Address: ph.Address,
				Height:  ph.Height,
			})
		}
	})

	return bcp, err
}

// ResendTransaction resent the transaction and return a *ResendResult
func (gw *Gateway) ResendTransaction(txn cipher.SHA256) interface{} {
	var result interface{}
	gw.strand("ResendTransaction", func() {
		result = gw.drpc.ResendTransaction(gw.d.Visor, gw.d.Pool, txn)
	})
	return result
}

// ResendUnconfirmedTxns resents all unconfirmed transactions
func (gw *Gateway) ResendUnconfirmedTxns() (*ResendResult, error) {
	var rlt *ResendResult
	var err error

	gw.strand("ResendUnconfirmedTxns", func() {
		rlt, err = gw.drpc.ResendUnconfirmedTxns(gw.d.Visor, gw.d.Pool)
	})

	return rlt, err
}

// GetBlockchainMetadata returns a *visor.BlockchainMetadata
func (gw *Gateway) GetBlockchainMetadata() (*visor.BlockchainMetadata, error) {
	var bcm *visor.BlockchainMetadata
	var err error

	gw.strand("GetBlockchainMetadata", func() {
		bcm, err = gw.v.GetBlockchainMetadata()
	})

	return bcm, err
}

// GetSignedBlockByHash returns the block by hash
func (gw *Gateway) GetSignedBlockByHash(hash cipher.SHA256) (block coin.SignedBlock, ok bool) {
	gw.strand("GetSignedBlockByHash", func() {
		b, err := gw.v.GetSignedBlockByHash(hash)
		if err != nil {
			logger.Error("gateway.GetSignedBlockByHash failed: %v", err)
			return
		}
		if b == nil {
			return
		}

		block = *b
		ok = true
	})
	return
}

// GetSignedBlockBySeq returns blcok by seq
func (gw *Gateway) GetSignedBlockBySeq(seq uint64) (block coin.SignedBlock, ok bool) {
	gw.strand("GetSignedBlockBySeq", func() {
		b, err := gw.v.GetSignedBlockBySeq(seq)
		if err != nil {
			logger.Error("gateway.GetSignedBlockBySeq failed: %v", err)
			return
		}
		if b == nil {
			return
		}
		block = *b
		ok = true
	})
	return
}

// GetBlocks returns a *visor.ReadableBlocks
func (gw *Gateway) GetBlocks(start, end uint64) (*visor.ReadableBlocks, error) {
	var blocks []coin.SignedBlock
	var err error

	gw.strand("GetBlocks", func() {
		blocks, err = gw.v.GetBlocks(start, end)
	})
	if err != nil {
		return nil, err
	}

	return visor.NewReadableBlocks(blocks)
}

// GetBlocksInDepth returns blocks in different depth
func (gw *Gateway) GetBlocksInDepth(vs []uint64) (*visor.ReadableBlocks, error) {
	var blocks []coin.SignedBlock
	var err error

	gw.strand("GetBlocksInDepth", func() {
		for _, n := range vs {
			var b *coin.SignedBlock
			b, err = gw.v.GetSignedBlockBySeq(n)
			if err != nil {
				err = fmt.Errorf("get block %v failed: %v", n, err)
				return
			}

			if b == nil {
				return
			}

			blocks = append(blocks, *b)
		}
	})

	if err != nil {
		return nil, err
	}

	return visor.NewReadableBlocks(blocks)
}

// GetLastBlocks get last N blocks
func (gw *Gateway) GetLastBlocks(num uint64) (*visor.ReadableBlocks, error) {
	var blocks []coin.SignedBlock
	var err error

	gw.strand("GetLastBlocks", func() {
		blocks, err = gw.v.GetLastBlocks(num)
	})
	if err != nil {
		return nil, err
	}

	return visor.NewReadableBlocks(blocks)
}

// OutputsFilter used as optional arguments in GetUnspentOutputs method
type OutputsFilter func(outputs coin.UxArray) coin.UxArray

// GetUnspentOutputs gets unspent outputs and returns the filtered results,
// Note: all filters will be executed as the pending sequence in 'AND' mode.
func (gw *Gateway) GetUnspentOutputs(filters ...OutputsFilter) (visor.ReadableOutputSet, error) {
	// unspent outputs
	var unspentOutputs []coin.UxOut
	// unconfirmed spending outputs
	var uncfmSpendingOutputs coin.UxArray
	// unconfirmed incoming outputs
	var uncfmIncomingOutputs coin.UxArray
	var err error
	gw.strand("GetUnspentOutputs", func() {
		unspentOutputs, err = gw.v.GetUnspentOutputs()
		if err != nil {
			err = fmt.Errorf("get unspent output readables failed: %v", err)
			return
		}

		uncfmSpendingOutputs, err = gw.v.UnconfirmedSpendingOutputs()
		if err != nil {
			err = fmt.Errorf("get unconfirmed spending outputs failed: %v", err)
			return
		}

		uncfmIncomingOutputs, err = gw.v.UnconfirmedIncomingOutputs()
		if err != nil {
			err = fmt.Errorf("get all incomming outputs failed: %v", err)
			return
		}
	})

	if err != nil {
		return visor.ReadableOutputSet{}, err
	}

	for _, flt := range filters {
		unspentOutputs = flt(unspentOutputs)
		uncfmSpendingOutputs = flt(uncfmSpendingOutputs)
		uncfmIncomingOutputs = flt(uncfmIncomingOutputs)
	}

	outputSet := visor.ReadableOutputSet{}
	outputSet.HeadOutputs, err = visor.NewReadableOutputs(unspentOutputs)
	if err != nil {
		return visor.ReadableOutputSet{}, err
	}

	outputSet.OutgoingOutputs, err = visor.NewReadableOutputs(uncfmSpendingOutputs)
	if err != nil {
		return visor.ReadableOutputSet{}, err
	}

	outputSet.IncomingOutputs, err = visor.NewReadableOutputs(uncfmIncomingOutputs)
	if err != nil {
		return visor.ReadableOutputSet{}, err
	}

	return outputSet, nil
}

// FbyAddressesNotIncluded filters the unspent outputs that are not owned by the addresses
func FbyAddressesNotIncluded(addrs []string) OutputsFilter {
	return func(outputs coin.UxArray) coin.UxArray {
		addrMatch := coin.UxArray{}
		addrMap := make(map[string]bool)
		for _, addr := range addrs {
			addrMap[addr] = false
		}

		for _, u := range outputs {
			_, ok := addrMap[u.Body.Address.String()]
			if !ok {
				addrMatch = append(addrMatch, u)
			}
		}
		return addrMatch
	}
}

// FbyAddresses filters the unspent outputs that owned by the addresses
func FbyAddresses(addrs []string) OutputsFilter {
	return func(outputs coin.UxArray) coin.UxArray {
		addrMatch := coin.UxArray{}
		addrMap := make(map[string]bool)
		for _, addr := range addrs {
			addrMap[addr] = true
		}

		for _, u := range outputs {
			if _, ok := addrMap[u.Body.Address.String()]; ok {
				addrMatch = append(addrMatch, u)
			}
		}
		return addrMatch
	}
}

// FbyHashes filters the unspent outputs that have hashes matched.
func FbyHashes(hashes []string) OutputsFilter {
	return func(outputs coin.UxArray) coin.UxArray {
		hsMatch := coin.UxArray{}
		hsMap := make(map[string]bool)
		for _, h := range hashes {
			hsMap[h] = true
		}

		for _, u := range outputs {
			if _, ok := hsMap[u.Hash().Hex()]; ok {
				hsMatch = append(hsMatch, u)
			}
		}
		return hsMatch
	}
}

// GetTransaction returns transaction by txid
func (gw *Gateway) GetTransaction(txid cipher.SHA256) (*visor.Transaction, error) {
	var tx *visor.Transaction
	var err error
	gw.strand("GetTransaction", func() {
		tx, err = gw.v.GetTransaction(txid)
	})
	return tx, err
}

// GetTransactionResult gets transaction result by txid.
func (gw *Gateway) GetTransactionResult(txid cipher.SHA256) (*TransactionResult, error) {
	tx, err := gw.GetTransaction(txid)
	if err != nil {
		return nil, err
	}

	return NewTransactionResult(tx)
}

// InjectTransaction injects transaction
func (gw *Gateway) InjectTransaction(txn coin.Transaction) error {
	var err error
	gw.strand("InjectTransaction", func() {
		err = gw.d.Visor.InjectTransaction(txn, gw.d.Pool)
	})
	return err
}

// GetAddressTxns returns a *visor.TransactionResults
func (gw *Gateway) GetAddressTxns(a cipher.Address) (*TransactionResults, error) {
	var txs []visor.Transaction
	var err error

	gw.strand("GetAddressesTxns", func() {
		txs, err = gw.v.GetAddressTxns(a)
	})

	if err != nil {
		return nil, err
	}

	return NewTransactionResults(txs)
}

// GetUxOutByID gets UxOut by hash id.
func (gw *Gateway) GetUxOutByID(id cipher.SHA256) (*historydb.UxOut, error) {
	var uxout *historydb.UxOut
	var err error
	gw.strand("GetUxOutByID", func() {
		uxout, err = gw.v.GetUxOutByID(id)
	})
	return uxout, err
}

// GetAddrUxOuts gets all the address affected UxOuts.
func (gw *Gateway) GetAddrUxOuts(addr cipher.Address) ([]*historydb.UxOutJSON, error) {
	var uxouts []*historydb.UxOut
	var err error
	gw.strand("GetAddrUxOuts", func() {
		uxouts, err = gw.v.GetAddrUxOuts(addr)
	})

	uxs := make([]*historydb.UxOutJSON, len(uxouts))
	for i, ux := range uxouts {
		uxs[i] = historydb.NewUxOutJSON(ux)
	}

	return uxs, err
}

// GetTimeNow returns the current Unix time
func (gw *Gateway) GetTimeNow() uint64 {
	return uint64(utc.UnixNow())
}

// GetAllUnconfirmedTxns returns all unconfirmed transactions
func (gw *Gateway) GetAllUnconfirmedTxns() ([]visor.UnconfirmedTxn, error) {
	var txns []visor.UnconfirmedTxn
	var err error
	gw.strand("GetAllUnconfirmedTxns", func() {
		txns, err = gw.v.GetAllUnconfirmedTxns()
	})
	return txns, err
}

// GetUnconfirmedTxns returns addresses related unconfirmed transactions
func (gw *Gateway) GetUnconfirmedTxns(addrs []cipher.Address) ([]visor.UnconfirmedTxn, error) {
	var txns []visor.UnconfirmedTxn
	var err error
	gw.strand("GetUnconfirmedTxns", func() {
		txns, err = gw.v.GetUnconfirmedTxns(visor.ToAddresses(addrs))
	})
	return txns, err
}

// GetLastTxs returns last confirmed transactions, return nil if empty
func (gw *Gateway) GetLastTxs() ([]*visor.Transaction, error) {
	var txns []*visor.Transaction
	var err error
	gw.strand("GetLastTxs", func() {
		txns, err = gw.v.GetLastTxs()
	})
	return txns, err
}

// // spendValidator implements the wallet.Validator interface
// type spendValidator struct {
// 	v       *visor.Visor
// 	unspent blockdb.UnspentPool
// }

// func newSpendValidator(v *visor.Visor, unspent blockdb.UnspentPool) *spendValidator {
// 	return &spendValidator{
// 		v:       v,
// 		unspent: unspent,
// 	}
// }

// func (sv spendValidator) HasUnconfirmedSpendTx(addr []cipher.Address) (bool, error) {
// 	aux, err := sv.v.SpendsOfAddresses(addr, sv.unspent)
// 	if err != nil {
// 		return false, err
// 	}

// 	return len(aux) > 0, nil
// }

// Spend spends coins from given wallet and broadcast it,
// return transaction or error.
func (gw *Gateway) Spend(wltID string, amt wallet.Balance, dest cipher.Address) (*coin.Transaction, error) {
	var tx *coin.Transaction
	var err error

	gw.strand("Spend", func() {
		// NOTE: visor.Visor makes multiple bolt.Tx calls
		// gw.d.Visor.InjectTransaction eventually calls visor.Visor methods,
		// these would all need to be gathered into one bolt.Tx,
		// this is too difficult.
		// This should be safe anyway as long as we use strand()
		var addrs []cipher.Address
		addrs, err = gw.v.Wallets.GetAddresses(wltID)
		if err != nil {
			return
		}

		// Check that this is not trying to spend unconfirmed outputs
		var auxs coin.AddressUxOuts
		auxs, err = gw.getUnspentsForSpending(addrs)
		if err != nil {
			return
		}

		var headTime uint64
		headTime, err = gw.v.GetHeadBlockTime()
		if err != nil {
			err = fmt.Errorf("GetHeadBlockTime failed: %v", err)
			return
		}

		// Create and sign transaction
		tx, err = gw.v.Wallets.CreateAndSignTransaction(wltID, auxs, headTime, amt, dest)
		if err != nil {
			err = fmt.Errorf("CreateAndSignTransaction failed: %v", err)
			return
		}

		// Inject transaction to the network
		if err = gw.d.Visor.InjectTransaction(*tx, gw.d.Pool); err != nil {
			err = fmt.Errorf("InjectTransaction failed: %v", err)
			return
		}
	})

	return tx, err
}

// CreateSpendingTransaction creates spending transactions
func (gw *Gateway) CreateSpendingTransaction(wlt wallet.Wallet, amt wallet.Balance, dest cipher.Address) (*coin.Transaction, error) {
	var tx *coin.Transaction
	var err error

	gw.strand("CreateSpendingTransaction", func() {
		// NOTE: visor.Visor makes multiple bolt.Tx calls
		// gw.d.Visor.InjectTransaction eventually calls visor.Visor methods,
		// these would all need to be gathered into one bolt.Tx,
		// this is too difficult.
		// This should be safe anyway as long as we use strand()
		addrs := wlt.GetAddresses()

		// Check that this is not trying to spend unconfirmed outputs
		var auxs coin.AddressUxOuts
		auxs, err = gw.getUnspentsForSpending(addrs)
		if err != nil {
			return
		}

		var headTime uint64
		headTime, err = gw.v.GetHeadBlockTime()
		if err != nil {
			return
		}

		// Create and sign transaction
		tx, err = wlt.CreateAndSignTransaction(auxs, headTime, amt, dest)
	})
	return tx, err
}

func (gw *Gateway) getUnspentsForSpending(addrs []cipher.Address) (coin.AddressUxOuts, error) {
	auxs, err := gw.v.UnconfirmedSpendsOfAddresses(addrs)
	if err != nil {
		err = fmt.Errorf("UnconfirmedSpendsOfAddresses failed: %v", err)
		return nil, err
	}

	// Check that this is not trying to spend unconfirmed outputs
	if len(auxs) > 0 {
		return nil, errors.New("please spend after your pending transaction is confirmed")
	}

	auxs, err = gw.v.GetUnspentsOfAddrs(addrs)
	if err != nil {
		err = fmt.Errorf("GetUnspentsOfAddrs failed: %v", err)
		return nil, err
	}

	return auxs, nil
}

// NewWallet creates wallet
func (gw *Gateway) NewWallet(wltName string, options ...wallet.Option) (wallet.Wallet, error) {
	var wlt wallet.Wallet
	var err error
	gw.strand("NewWallet", func() {
		wlt, err = gw.v.Wallets.CreateWallet(wltName, options...)
	})
	return wlt, err
}

// GetWalletBalance returns balance pair of specific wallet
func (gw *Gateway) GetWalletBalance(wltID string) (wallet.BalancePair, error) {
	var balance wallet.BalancePair
	var err error

	gw.strand("GetWalletBalance", func() {
		var addrs []cipher.Address
		addrs, err = gw.v.Wallets.GetAddresses(wltID)
		if err != nil {
			return
		}

		balance, err = gw.v.GetBalanceOfAddrs(addrs)
	})

	return balance, err
}

// GetBalanceOfAddrs gets balance of given addresses
func (gw *Gateway) GetBalanceOfAddrs(addrs []cipher.Address) (wallet.BalancePair, error) {
	var balance wallet.BalancePair
	var err error

	gw.strand("GetBalanceOfAddrs", func() {
		balance, err = gw.v.GetBalanceOfAddrs(addrs)
	})

	return balance, err
}

// GetWalletDir returns path for storing wallet files
func (gw *Gateway) GetWalletDir() string {
	return gw.v.Config.WalletDirectory
}

// NewAddresses generate addresses in given wallet
func (gw *Gateway) NewAddresses(wltID string, n int) ([]cipher.Address, error) {
	var addrs []cipher.Address
	var err error
	gw.strand("NewAddresses", func() {
		addrs, err = gw.v.Wallets.NewAddresses(wltID, n)
	})
	return addrs, err
}

// UpdateWalletLabel updates the label of wallet
func (gw *Gateway) UpdateWalletLabel(wltID, label string) error {
	var err error
	gw.strand("UpdateWalletLabel", func() {
		err = gw.v.Wallets.UpdateWalletLabel(wltID, label)
	})
	return err
}

// GetWallet returns wallet by id
func (gw *Gateway) GetWallet(wltID string) (wallet.Wallet, bool) {
	var w wallet.Wallet
	var ok bool
	gw.strand("GetWallet", func() {
		w, ok = gw.v.Wallets.GetWallet(wltID)
	})
	return w, ok
}

// GetWallets returns wallets
func (gw *Gateway) GetWallets() wallet.Wallets {
	var w wallet.Wallets
	gw.strand("GetWallets", func() {
		w = gw.v.Wallets.GetWallets()
	})
	return w
}

// GetWalletUnconfirmedTxns returns all unconfirmed transactions in given wallet
func (gw *Gateway) GetWalletUnconfirmedTxns(wltID string) ([]visor.UnconfirmedTxn, error) {
	var txns []visor.UnconfirmedTxn
	var err error
	gw.strand("GetWalletUnconfirmedTxns", func() {
		var addrs []cipher.Address
		addrs, err = gw.v.Wallets.GetAddresses(wltID)
		if err != nil {
			return
		}

		txns, err = gw.v.GetUnconfirmedTxns(visor.ToAddresses(addrs))
	})

	return txns, err
}

// ReloadWallets reloads all wallets
func (gw *Gateway) ReloadWallets() error {
	var err error
	gw.strand("ReloadWallets", func() {
		err = gw.v.Wallets.ReloadWallets()
	})
	return err
}

// GetBuildInfo returns node build info.
func (gw *Gateway) GetBuildInfo() visor.BuildInfo {
	var bi visor.BuildInfo
	gw.strand("GetBuildInfo", func() {
		bi = gw.v.Config.BuildInfo
	})
	return bi
}

// TransactionResult represents transaction result
type TransactionResult struct {
	Status      visor.TransactionStatus   `json:"status"`
	Time        uint64                    `json:"time"`
	Transaction visor.ReadableTransaction `json:"txn"`
}

// NewTransactionResult converts Transaction to TransactionResult
func NewTransactionResult(tx *visor.Transaction) (*TransactionResult, error) {
	if tx == nil {
		return nil, nil
	}

	rbTx, err := visor.NewReadableTransaction(tx)
	if err != nil {
		return nil, err
	}

	return &TransactionResult{
		Transaction: *rbTx,
		Status:      tx.Status,
		Time:        tx.Time,
	}, nil
}

// TransactionResults array of transaction results
type TransactionResults struct {
	Txns []TransactionResult `json:"txns"`
}

// NewTransactionResults converts []Transaction to []TransactionResults
func NewTransactionResults(txs []visor.Transaction) (*TransactionResults, error) {
	txRlts := make([]TransactionResult, 0, len(txs))
	for _, tx := range txs {
		rbTx, err := visor.NewReadableTransaction(&tx)
		if err != nil {
			return nil, err
		}

		txRlts = append(txRlts, TransactionResult{
			Transaction: *rbTx,
			Status:      tx.Status,
			Time:        tx.Time,
		})
	}

	return &TransactionResults{
		Txns: txRlts,
	}, nil
}
