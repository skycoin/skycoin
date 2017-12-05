package daemon

import (
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/daemon/strand"
	"github.com/skycoin/skycoin/src/util/utc"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/wallet"

	"fmt"

	"github.com/skycoin/skycoin/src/visor/blockdb"
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
	requests chan strand.Request
}

// NewGateway create and init an Gateway instance.
func NewGateway(c GatewayConfig, D *Daemon) *Gateway {
	return &Gateway{
		Config:   c,
		drpc:     RPC{},
		vrpc:     visor.MakeRPC(D.Visor.v),
		d:        D,
		v:        D.Visor.v,
		requests: make(chan strand.Request, c.BufferSize),
	}
}

func (gw *Gateway) strand(name string, f func()) {
	name = fmt.Sprintf("daemon.Gateway.%s", name)
	strand.Strand(logger, gw.requests, name, func() error {
		f()
		return nil
	})
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

// GetBlockchainProgress returns a *BlockchainProgress
func (gw *Gateway) GetBlockchainProgress() interface{} {
	var bcp interface{}
	gw.strand("GetBlockchainProgress", func() {
		bcp = gw.drpc.GetBlockchainProgress(gw.d.Visor)
	})
	return bcp
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
func (gw *Gateway) ResendUnconfirmedTxns() (rlt *ResendResult) {
	gw.strand("ResendUnconfirmedTxns", func() {
		rlt = gw.drpc.ResendUnconfirmedTxns(gw.d.Visor, gw.d.Pool)
	})
	return
}

// GetBlockchainMetadata returns a *visor.BlockchainMetadata
func (gw *Gateway) GetBlockchainMetadata() interface{} {
	var bcm interface{}
	gw.strand("GetBlockchainMetadata", func() {
		bcm = gw.vrpc.GetBlockchainMetadata(gw.v)
	})
	return bcm
}

// GetBlockByHash returns the block by hash
func (gw *Gateway) GetBlockByHash(hash cipher.SHA256) (block coin.SignedBlock, ok bool) {
	gw.strand("GetBlockByHash", func() {
		b, err := gw.v.GetBlockByHash(hash)
		if err != nil {
			logger.Error("gateway.GetBlockByHash failed: %v", err)
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

// GetBlockBySeq returns blcok by seq
func (gw *Gateway) GetBlockBySeq(seq uint64) (block coin.SignedBlock, ok bool) {
	gw.strand("GetBlockBySeq", func() {
		b, err := gw.v.GetBlockBySeq(seq)
		if err != nil {
			logger.Error("gateway.GetBlockBySeq failed: %v", err)
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
	gw.strand("GetBlocks", func() {
		blocks = gw.vrpc.GetBlocks(gw.v, start, end)
	})

	return visor.NewReadableBlocks(blocks)
}

// GetBlocksInDepth returns blocks in different depth
func (gw *Gateway) GetBlocksInDepth(vs []uint64) (*visor.ReadableBlocks, error) {
	blocks := []coin.SignedBlock{}
	var err error

	gw.strand("GetBlocksInDepth", func() {
		for _, n := range vs {
			var b *coin.SignedBlock
			b, err = gw.vrpc.GetBlockBySeq(gw.v, n)
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
	gw.strand("GetLastBlocks", func() {
		blocks = gw.vrpc.GetLastBlocks(gw.v, num)
	})

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
func (gw *Gateway) GetTransaction(txid cipher.SHA256) (tx *visor.Transaction, err error) {
	gw.strand("GetTransaction", func() {
		tx, err = gw.v.GetTransaction(txid)
	})
	return
}

// GetTransactionResult gets transaction result by txid.
func (gw *Gateway) GetTransactionResult(txid cipher.SHA256) (*visor.TransactionResult, error) {
	var tx *visor.Transaction
	var err error
	gw.strand("GetTransactionResult", func() {
		tx, err = gw.vrpc.GetTransaction(gw.v, txid)
	})

	if err != nil {
		return nil, err
	}

	return visor.NewTransactionResult(tx)
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
func (gw *Gateway) GetAddressTxns(a cipher.Address) (*visor.TransactionResults, error) {
	var txs []visor.Transaction
	var err error

	gw.strand("GetAddressesTxns", func() {
		txs, err = gw.vrpc.GetAddressTxns(gw.v, a)
	})

	if err != nil {
		return nil, err
	}

	return visor.NewTransactionResults(txs)
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
func (gw *Gateway) GetAllUnconfirmedTxns() []visor.UnconfirmedTxn {
	var txns []visor.UnconfirmedTxn
	gw.strand("GetAllUnconfirmedTxns", func() {
		txns = gw.v.GetAllUnconfirmedTxns()
	})
	return txns
}

// GetUnconfirmedTxns returns addresses related unconfirmed transactions
func (gw *Gateway) GetUnconfirmedTxns(addrs []cipher.Address) []visor.UnconfirmedTxn {
	var txns []visor.UnconfirmedTxn
	gw.strand("GetUnconfirmedTxns", func() {
		txns = gw.v.GetUnconfirmedTxns(visor.ToAddresses(addrs))
	})
	return txns
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

// GetUnspent returns the unspent pool
func (gw *Gateway) GetUnspent() blockdb.UnspentPool {
	var unspent blockdb.UnspentPool
	gw.strand("GetUnspent", func() {
		unspent = gw.v.Blockchain.Unspent()
	})
	return unspent
}

// impelemts the wallet.Validator interface
type spendValidator struct {
	uncfm   *visor.UnconfirmedTxnPool
	unspent blockdb.UnspentPool
}

func newSpendValidator(uncfm *visor.UnconfirmedTxnPool, unspent blockdb.UnspentPool) *spendValidator {
	return &spendValidator{
		uncfm:   uncfm,
		unspent: unspent,
	}
}

func (sv spendValidator) HasUnconfirmedSpendTx(addr []cipher.Address) (bool, error) {
	aux, err := sv.uncfm.SpendsOfAddresses(addr, sv.unspent)
	if err != nil {
		return false, err
	}

	return len(aux) > 0, nil
}

// Spend spends coins from given wallet and broadcast it,
// return transaction or error.
func (gw *Gateway) Spend(wltID string, amt wallet.Balance, dest cipher.Address) (*coin.Transaction, error) {
	var tx *coin.Transaction
	var err error
	gw.strand("Spend", func() {
		// create spend validator
		unspent := gw.v.Blockchain.Unspent()
		sv := newSpendValidator(gw.v.Unconfirmed, unspent)
		// create and sign transaction
		tx, err = gw.vrpc.CreateAndSignTransaction(wltID,
			sv,
			unspent,
			gw.v.Blockchain.Time(),
			amt,
			dest)
		if err != nil {
			err = fmt.Errorf("Create transaction failed: %v", err)
			return
		}

		// inject transaction
		if err = gw.d.Visor.InjectTransaction(*tx, gw.d.Pool); err != nil {
			err = fmt.Errorf("Inject transaction failed: %v", err)
		}
	})

	return tx, err
}

// NewWallet creates wallet
func (gw *Gateway) NewWallet(wltName string, options ...wallet.Option) (wallet.Wallet, error) {
	var wlt wallet.Wallet
	var err error
	gw.strand("NewWallet", func() {
		wlt, err = gw.vrpc.NewWallet(wltName, options...)
	})
	return wlt, err
}

// LoadAndScanWallet loads wallet from given seed and scan ahead N addresses
func (gw *Gateway) LoadAndScanWallet(wltName string, seed string, scanN uint64, options ...wallet.Option) (wallet.Wallet, error) {
	var wlt wallet.Wallet
	var err error
	gw.strand("LoadAndScanWallet", func() {
		wlt, err = gw.v.LoadAndScanWallet(wltName, seed, scanN, options...)
	})
	return wlt, err
}

// CreateSpendingTransaction creates spending transactions
func (gw *Gateway) CreateSpendingTransaction(wlt wallet.Wallet, amt wallet.Balance, dest cipher.Address) (*coin.Transaction, error) {
	var tx *coin.Transaction
	var err error
	gw.strand("CreateSpendingTransaction", func() {
		// generate spend validator
		unspent := gw.v.Blockchain.Unspent()
		sv := newSpendValidator(gw.v.Unconfirmed, unspent)

		// create and sign transaction
		headTime := gw.v.Blockchain.Time()
		tx, err = wlt.CreateAndSignTransaction(sv, unspent, headTime, amt, dest)
	})
	return tx, err
}

// GetWalletBalance returns balance pair of specific wallet
func (gw *Gateway) GetWalletBalance(wltID string) (wallet.BalancePair, error) {
	var balance wallet.BalancePair
	var err error
	gw.strand("GetWalletBalance", func() {
		var addrs []cipher.Address
		addrs, err = gw.vrpc.GetWalletAddresses(wltID)
		if err != nil {
			return
		}
		auxs := gw.vrpc.GetUnspent(gw.v).GetUnspentsOfAddrs(addrs)

		var spendUxs coin.AddressUxOuts
		spendUxs, err = gw.vrpc.GetUnconfirmedSpends(gw.v, addrs)
		if err != nil {
			err = fmt.Errorf("get unconfimed spending failed when checking wallet balance: %v", err)
			return
		}

		var recvUxs coin.AddressUxOuts
		recvUxs, err = gw.vrpc.GetUnconfirmedReceiving(gw.v, addrs)
		if err != nil {
			err = fmt.Errorf("get unconfirmed receiving failed when when checking wallet balance: %v", err)
			return
		}

		coins1, hours1 := gw.v.AddressBalance(auxs)
		coins2, hours2 := gw.v.AddressBalance(auxs.Sub(spendUxs).Add(recvUxs))
		balance = wallet.BalancePair{
			Confirmed: wallet.Balance{Coins: coins1, Hours: hours1},
			Predicted: wallet.Balance{Coins: coins2, Hours: hours2},
		}
	})

	return balance, err
}

// GetBalanceOfAddrs gets balance of given addresses
func (gw *Gateway) GetBalanceOfAddrs(addrs []cipher.Address) ([]wallet.BalancePair, error) {
	var bps []wallet.BalancePair
	var err error
	gw.strand("GetBalanceOfAddrs", func() {
		bps, err = gw.v.GetBalanceOfAddrs(addrs)
	})

	return bps, err
}

// GetWalletDir returns path for storing wallet files
func (gw *Gateway) GetWalletDir() string {
	return gw.v.Config.WalletDirectory
}

// NewAddresses generate addresses in given wallet
func (gw *Gateway) NewAddresses(wltID string, n uint64) ([]cipher.Address, error) {
	var addrs []cipher.Address
	var err error
	gw.strand("NewAddresses", func() {
		addrs, err = gw.vrpc.NewAddresses(wltID, n)
	})
	return addrs, err
}

// UpdateWalletLabel updates the label of wallet
func (gw *Gateway) UpdateWalletLabel(wltID, label string) error {
	var err error
	gw.strand("UpdateWalletLabel", func() {
		err = gw.vrpc.UpdateWalletLabel(wltID, label)
	})
	return err
}

// GetWallet returns wallet by id
func (gw *Gateway) GetWallet(wltID string) (wallet.Wallet, bool) {
	var w wallet.Wallet
	var ok bool
	gw.strand("GetWallet", func() {
		w, ok = gw.vrpc.GetWallet(wltID)
	})
	return w, ok
}

// GetWallets returns wallets
func (gw *Gateway) GetWallets() wallet.Wallets {
	var w wallet.Wallets
	gw.strand("GetWallets", func() {
		w = gw.vrpc.GetWallets()
	})
	return w
}

// GetWalletUnconfirmedTxns returns all unconfirmed transactions in given wallet
func (gw *Gateway) GetWalletUnconfirmedTxns(wltID string) ([]visor.UnconfirmedTxn, error) {
	var txns []visor.UnconfirmedTxn
	var err error
	gw.strand("GetWalletUnconfirmedTxns", func() {
		var addrs []cipher.Address
		addrs, err = gw.vrpc.GetWalletAddresses(wltID)
		if err != nil {
			return
		}

		txns = gw.v.GetUnconfirmedTxns(visor.ToAddresses(addrs))
	})

	return txns, err
}

// ReloadWallets reloads all wallets
func (gw *Gateway) ReloadWallets() error {
	var err error
	gw.strand("ReloadWallets", func() {
		err = gw.vrpc.ReloadWallets()
	})
	return err
}

// GetBuildInfo returns node build info.
func (gw *Gateway) GetBuildInfo() visor.BuildInfo {
	var bi visor.BuildInfo
	gw.strand("GetBuildInfo", func() {
		bi = gw.vrpc.GetBuildInfo()
	})
	return bi
}
