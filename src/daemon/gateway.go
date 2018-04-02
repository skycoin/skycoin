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
	BufferSize       int
	DisableWalletAPI bool
}

// NewGatewayConfig create and init an GatewayConfig
func NewGatewayConfig() GatewayConfig {
	return GatewayConfig{
		BufferSize:       32,
		DisableWalletAPI: false,
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
	quit     chan struct{}
}

// NewGateway create and init an Gateway instance.
func NewGateway(c GatewayConfig, d *Daemon) *Gateway {
	return &Gateway{
		Config:   c,
		drpc:     RPC{},
		vrpc:     visor.MakeRPC(d.Visor.v),
		d:        d,
		v:        d.Visor.v,
		requests: make(chan strand.Request, c.BufferSize),
		quit:     make(chan struct{}),
	}
}

// Shutdown closes the Gateway
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
func (gw *Gateway) GetConnections() *Connections {
	var conns *Connections
	gw.strand("GetConnections", func() {
		conns = gw.drpc.GetConnections(gw.d)
	})
	return conns
}

// GetDefaultConnections returns default connections
func (gw *Gateway) GetDefaultConnections() []string {
	var conns []string
	gw.strand("GetDefaultConnections", func() {
		conns = gw.drpc.GetDefaultConnections(gw.d)
	})
	return conns
}

// GetConnection returns a *Connection of specific address
func (gw *Gateway) GetConnection(addr string) *Connection {
	var conn *Connection
	gw.strand("GetConnection", func() {
		conn = gw.drpc.GetConnection(gw.d, addr)
	})
	return conn
}

// GetTrustConnections returns all trusted connections,
// including private and public
func (gw *Gateway) GetTrustConnections() []string {
	var conn []string
	gw.strand("GetTrustConnections", func() {
		conn = gw.drpc.GetTrustConnections(gw.d)
	})
	return conn
}

// GetExchgConnection returns all exchangeable connections,
// including private and public
func (gw *Gateway) GetExchgConnection() []string {
	var conn []string
	gw.strand("GetExchgConnection", func() {
		conn = gw.drpc.GetAllExchgConnections(gw.d)
	})
	return conn
}

/* Blockchain & Transaction status */

// GetBlockchainProgress returns a *BlockchainProgress
func (gw *Gateway) GetBlockchainProgress() *BlockchainProgress {
	var bcp *BlockchainProgress
	gw.strand("GetBlockchainProgress", func() {
		bcp = gw.drpc.GetBlockchainProgress(gw.d.Visor)
	})
	return bcp
}

// ResendTransaction resent the transaction and return a *ResendResult
func (gw *Gateway) ResendTransaction(txn cipher.SHA256) *ResendResult {
	var result *ResendResult
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
func (gw *Gateway) GetBlockchainMetadata() *visor.BlockchainMetadata {
	var bcm *visor.BlockchainMetadata
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
			logger.Errorf("gateway.GetBlockByHash failed: %v", err)
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
			logger.Errorf("gateway.GetBlockBySeq failed: %v", err)
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
func (gw *Gateway) GetUnspentOutputs(filters ...OutputsFilter) (*visor.ReadableOutputSet, error) {
	// unspent outputs
	var unspentOutputs []coin.UxOut
	// unconfirmed spending outputs
	var uncfmSpendingOutputs coin.UxArray
	// unconfirmed incoming outputs
	var uncfmIncomingOutputs coin.UxArray
	var headTime uint64
	var err error
	gw.strand("GetUnspentOutputs", func() {
		headTime = gw.v.Blockchain.Time()

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
			err = fmt.Errorf("get all incoming outputs failed: %v", err)
			return
		}
	})

	if err != nil {
		return nil, err
	}

	for _, flt := range filters {
		unspentOutputs = flt(unspentOutputs)
		uncfmSpendingOutputs = flt(uncfmSpendingOutputs)
		uncfmIncomingOutputs = flt(uncfmIncomingOutputs)
	}

	outputSet := visor.ReadableOutputSet{}
	outputSet.HeadOutputs, err = visor.NewReadableOutputs(headTime, unspentOutputs)
	if err != nil {
		return nil, err
	}

	outputSet.OutgoingOutputs, err = visor.NewReadableOutputs(headTime, uncfmSpendingOutputs)
	if err != nil {
		return nil, err
	}

	outputSet.IncomingOutputs, err = visor.NewReadableOutputs(headTime, uncfmIncomingOutputs)
	if err != nil {
		return nil, err
	}

	return &outputSet, nil
}

// FbyAddressesNotIncluded filters the unspent outputs that are not owned by the addresses
func FbyAddressesNotIncluded(addrs []string) OutputsFilter {
	return func(outputs coin.UxArray) coin.UxArray {
		addrMatch := coin.UxArray{}
		addrMap := MakeSearchMap(addrs)

		for _, u := range outputs {
			if _, ok := addrMap[u.Body.Address.String()]; !ok {
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
		addrMap := MakeSearchMap(addrs)

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
		hsMap := MakeSearchMap(hashes)

		for _, u := range outputs {
			if _, ok := hsMap[u.Hash().Hex()]; ok {
				hsMatch = append(hsMatch, u)
			}
		}
		return hsMatch
	}
}

// MakeSearchMap returns a search indexed map for use in filters
func MakeSearchMap(addrs []string) map[string]struct{} {
	addrMap := make(map[string]struct{})
	for _, addr := range addrs {
		addrMap[addr] = struct{}{}
	}

	return addrMap
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

// InjectBroadcastTransaction injects and broadcasts a transaction
func (gw *Gateway) InjectBroadcastTransaction(txn coin.Transaction) error {
	var err error
	gw.strand("InjectBroadcastTransaction", func() {
		err = gw.d.Visor.InjectBroadcastTransaction(txn, gw.d.Pool)
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

// GetTransactions returns transactions filtered by zero or more visor.TxFilter
func (gw *Gateway) GetTransactions(flts ...visor.TxFilter) ([]visor.Transaction, error) {
	var txns []visor.Transaction
	var err error
	gw.strand("GetTransactions", func() {
		txns, err = gw.v.GetTransactions(flts...)
	})
	return txns, err
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

// GetUnspent returns the unspent pool
func (gw *Gateway) GetUnspent() blockdb.UnspentPool {
	var unspent blockdb.UnspentPool
	gw.strand("GetUnspent", func() {
		unspent = gw.v.Blockchain.Unspent()
	})
	return unspent
}

// implements the wallet.Validator interface
type spendValidator struct {
	uncfm   visor.UnconfirmedTxnPooler
	unspent blockdb.UnspentPool
}

func newSpendValidator(uncfm visor.UnconfirmedTxnPooler, unspent blockdb.UnspentPool) *spendValidator {
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
// set password as nil if wallet is not encrypted, otherwise the password must be provied.
// return transaction or error.
func (gw *Gateway) Spend(wltID string, password []byte, coins uint64, dest cipher.Address) (*coin.Transaction, error) {
	if gw.Config.DisableWalletAPI {
		return nil, wallet.ErrWalletAPIDisabled
	}

	var tx *coin.Transaction
	var err error
	gw.strand("Spend", func() {
		// create spend validator
		unspent := gw.v.Blockchain.Unspent()
		sv := newSpendValidator(gw.v.Unconfirmed, unspent)
		// create and sign transaction
		tx, err = gw.vrpc.CreateAndSignTransaction(wltID, password, sv, unspent, gw.v.Blockchain.Time(), coins, dest)
		if err != nil {
			logger.Errorf("Create transaction failed: %v", err)
			return
		}

		// Inject transaction
		err = gw.d.Visor.InjectBroadcastTransaction(*tx, gw.d.Pool)
		if err != nil {
			logger.Errorf("Inject transaction failed: %v", err)
			return
		}
	})

	return tx, err
}

// CreateWallet creates wallet
func (gw *Gateway) CreateWallet(wltName string, options wallet.Options) (*wallet.Wallet, error) {
	if gw.Config.DisableWalletAPI {
		return nil, wallet.ErrWalletAPIDisabled
	}

	var wlt *wallet.Wallet
	var err error
	gw.strand("CreateWallet", func() {
		wlt, err = gw.vrpc.CreateWallet(wltName, options)
	})
	return wlt, err
}

// ScanAheadWalletAddresses loads wallet from given seed and scan ahead N addresses
// Set password as nil if the wallet is not encrypted, otherwise the password must be provided
func (gw *Gateway) ScanAheadWalletAddresses(wltName string, password []byte, scanN uint64) (*wallet.Wallet, error) {
	if gw.Config.DisableWalletAPI {
		return nil, wallet.ErrWalletAPIDisabled
	}

	var wlt *wallet.Wallet
	var err error
	gw.strand("ScanAheadWalletAddresses", func() {
		wlt, err = gw.v.ScanAheadWalletAddresses(wltName, password, scanN)
	})
	return wlt, err
}

// EncryptWallet encrypts the wallet
func (gw *Gateway) EncryptWallet(wltName string, password []byte) error {
	if gw.Config.DisableWalletAPI {
		return wallet.ErrWalletAPIDisabled
	}

	var err error
	gw.strand("EncryptWallet", func() {
		err = gw.v.Wallets.EncryptWallet(wltName, password)
	})
	return err
}

// GetWalletBalance returns balance pair of specific wallet
func (gw *Gateway) GetWalletBalance(wltID string) (wallet.BalancePair, error) {
	var balance wallet.BalancePair
	if gw.Config.DisableWalletAPI {
		return balance, wallet.ErrWalletAPIDisabled
	}

	var err error
	gw.strand("GetWalletBalance", func() {
		var addrs []cipher.Address
		addrs, err = gw.v.Wallets.GetAddresses(wltID)
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

		coins1, hours1, err := gw.v.AddressBalance(auxs)
		if err != nil {
			err = fmt.Errorf("Computing confirmed address balance failed: %v", err)
			return
		}

		coins2, hours2, err := gw.v.AddressBalance(auxs.Sub(spendUxs).Add(recvUxs))
		if err != nil {
			err = fmt.Errorf("Computing predicted address balance failed: %v", err)
			return
		}

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
func (gw *Gateway) GetWalletDir() (string, error) {
	if gw.Config.DisableWalletAPI {
		return "", wallet.ErrWalletAPIDisabled
	}
	return gw.v.Config.WalletDirectory, nil
}

// NewAddresses generate addresses in given wallet
func (gw *Gateway) NewAddresses(wltID string, password []byte, n uint64) ([]cipher.Address, error) {
	if gw.Config.DisableWalletAPI {
		return nil, wallet.ErrWalletAPIDisabled
	}

	var addrs []cipher.Address
	var err error
	gw.strand("NewAddresses", func() {
		addrs, err = gw.v.Wallets.NewAddresses(wltID, password, n)
	})
	return addrs, err
}

// UpdateWalletLabel updates the label of wallet
func (gw *Gateway) UpdateWalletLabel(wltID, label string) error {
	if gw.Config.DisableWalletAPI {
		return wallet.ErrWalletAPIDisabled
	}

	var err error
	gw.strand("UpdateWalletLabel", func() {
		err = gw.v.Wallets.UpdateWalletLabel(wltID, label)
	})
	return err
}

// GetWallet returns wallet by id
func (gw *Gateway) GetWallet(wltID string) (*wallet.Wallet, error) {
	if gw.Config.DisableWalletAPI {
		return nil, wallet.ErrWalletAPIDisabled
	}

	var w *wallet.Wallet
	var err error
	gw.strand("GetWallet", func() {
		w, err = gw.vrpc.GetWallet(wltID)
	})
	return w, err
}

// GetWallets returns wallets
func (gw *Gateway) GetWallets() (wallet.Wallets, error) {
	if gw.Config.DisableWalletAPI {
		return nil, wallet.ErrWalletAPIDisabled
	}

	var w wallet.Wallets
	gw.strand("GetWallets", func() {
		w = gw.v.Wallets.GetWallets()
	})
	return w, nil
}

// GetWalletUnconfirmedTxns returns all unconfirmed transactions in given wallet
func (gw *Gateway) GetWalletUnconfirmedTxns(wltID string) ([]visor.UnconfirmedTxn, error) {
	if gw.Config.DisableWalletAPI {
		return nil, wallet.ErrWalletAPIDisabled
	}

	var txns []visor.UnconfirmedTxn
	var err error
	gw.strand("GetWalletUnconfirmedTxns", func() {
		var addrs []cipher.Address
		addrs, err = gw.v.Wallets.GetAddresses(wltID)
		if err != nil {
			return
		}

		txns = gw.v.GetUnconfirmedTxns(visor.ToAddresses(addrs))
	})

	return txns, err
}

// ReloadWallets reloads all wallets
func (gw *Gateway) ReloadWallets() error {
	if gw.Config.DisableWalletAPI {
		return wallet.ErrWalletAPIDisabled
	}

	var err error
	gw.strand("ReloadWallets", func() {
		err = gw.v.Wallets.ReloadWallets()
	})
	return err
}

// UnloadWallet removes wallet of given id from memory.
func (gw *Gateway) UnloadWallet(id string) error {
	if gw.Config.DisableWalletAPI {
		return wallet.ErrWalletAPIDisabled
	}

	gw.strand("UnloadWallet", func() {
		gw.vrpc.UnloadWallet(id)
	})

	return nil
}

// IsWalletAPIDisabled returns if all wallet related apis are disabled
func (gw *Gateway) IsWalletAPIDisabled() bool {
	return gw.Config.DisableWalletAPI
}

// GetBuildInfo returns node build info.
func (gw *Gateway) GetBuildInfo() visor.BuildInfo {
	var bi visor.BuildInfo
	gw.strand("GetBuildInfo", func() {
		bi = gw.v.Config.BuildInfo
	})
	return bi
}

// GetRichlist returns rich list as desc order.
func (gw *Gateway) GetRichlist(includeDistribution bool) (visor.Richlist, error) {
	rbOuts, err := gw.GetUnspentOutputs()
	if err != nil {
		return nil, err
	}

	allAccounts, err := rbOuts.AggregateUnspentOutputs()
	if err != nil {
		return nil, err
	}

	lockedAddrs := visor.GetLockedDistributionAddresses()
	addrsMap := make(map[string]struct{}, len(lockedAddrs))
	for _, a := range lockedAddrs {
		addrsMap[a] = struct{}{}
	}

	richlist, err := visor.NewRichlist(allAccounts, addrsMap)
	if err != nil {
		return nil, err
	}

	if !includeDistribution {
		unlockedAddrs := visor.GetUnlockedDistributionAddresses()
		for _, a := range unlockedAddrs {
			addrsMap[a] = struct{}{}
		}
		richlist = richlist.FilterAddresses(addrsMap)
	}

	return richlist, nil
}

// GetAddressCount returns count number of unique address with uxouts > 0.
func (gw *Gateway) GetAddressCount() (uint64, error) {
	rbOuts, err := gw.GetUnspentOutputs()
	if err != nil {
		return 0, err
	}

	allAccounts, err := rbOuts.AggregateUnspentOutputs()
	if err != nil {
		return 0, err
	}

	return uint64(len(allAccounts)), nil
}
