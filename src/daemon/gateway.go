package daemon

import (
	"sort"
	"strings"
	"time"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/daemon/gnet"
	"github.com/skycoin/skycoin/src/daemon/pex"
	"github.com/skycoin/skycoin/src/params"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/visor/dbutil"
	"github.com/skycoin/skycoin/src/visor/historydb"
	"github.com/skycoin/skycoin/src/wallet"
)

// GatewayConfig configuration set of gateway.
type GatewayConfig struct {
	EnableWalletAPI bool
}

// NewGatewayConfig create and init an GatewayConfig
func NewGatewayConfig() GatewayConfig {
	return GatewayConfig{
		EnableWalletAPI: false,
	}
}

// Gateway RPC interface wrapper for daemon state
type Gateway struct {
	Config GatewayConfig

	// Backref to Daemon
	d *Daemon
	// Backref to Visor
	v *visor.Visor
}

// NewGateway create and init an Gateway instance.
func NewGateway(c GatewayConfig, d *Daemon) *Gateway {
	return &Gateway{
		Config: c,
		d:      d,
		v:      d.visor,
	}
}

// Connection a connection's state within the daemon
type Connection struct {
	Addr string
	Pex  pex.Peer
	Gnet GnetConnectionDetails
	ConnectionDetails
}

// GnetConnectionDetails connection data from gnet
type GnetConnectionDetails struct {
	ID           uint64
	LastSent     time.Time
	LastReceived time.Time
}

func newConnection(dc *connection, gc *gnet.Connection, pp *pex.Peer) Connection {
	c := Connection{}

	if dc != nil {
		c.Addr = dc.Addr
		c.ConnectionDetails = dc.ConnectionDetails
	}

	if gc != nil {
		c.Gnet = GnetConnectionDetails{
			ID:           gc.ID,
			LastSent:     gc.LastSent,
			LastReceived: gc.LastReceived,
		}
	}

	if pp != nil {
		c.Pex = *pp
	}

	return c
}

// newConnection creates a Connection from daemon.connection, gnet.Connection and pex.Peer
func (gw *Gateway) newConnection(c *connection) (*Connection, error) {
	if c == nil {
		return nil, nil
	}

	gc, err := gw.d.pool.Pool.GetConnection(c.Addr)
	if err != nil {
		return nil, err
	}

	var pp *pex.Peer
	listenAddr := c.ListenAddr()
	if listenAddr != "" {
		p, ok := gw.d.pex.GetPeer(listenAddr)
		if ok {
			pp = &p
		}
	}

	cc := newConnection(c, gc, pp)
	return &cc, nil
}

// GetConnections returns solicited (outgoing) connections
func (gw *Gateway) GetConnections(f func(c Connection) bool) ([]Connection, error) {
	return gw.getConnections(f)
}

func (gw *Gateway) getConnections(f func(c Connection) bool) ([]Connection, error) {
	if gw.d.pool.Pool == nil {
		return nil, nil
	}

	cs := gw.d.connections.all()

	conns := make([]Connection, 0)

	for _, c := range cs {
		cc, err := gw.newConnection(&c)
		if err != nil {
			return nil, err
		}

		ccc := *cc

		if !f(ccc) {
			continue
		}

		conns = append(conns, ccc)
	}

	// Sort connnections by IP address
	sort.Slice(conns, func(i, j int) bool {
		return strings.Compare(conns[i].Addr, conns[j].Addr) < 0
	})

	return conns, nil
}

// GetDefaultConnections returns the default hardcoded connection addresses
func (gw *Gateway) GetDefaultConnections() []string {
	conns := make([]string, len(gw.d.Config.DefaultConnections))
	copy(conns[:], gw.d.Config.DefaultConnections[:])
	return conns
}

// GetConnection returns a *Connection of specific address
func (gw *Gateway) GetConnection(addr string) (*Connection, error) {
	c := gw.d.connections.get(addr)
	if c == nil {
		return nil, nil
	}

	return gw.newConnection(c)
}

// Disconnect disconnects a connection by gnet ID
func (gw *Gateway) Disconnect(gnetID uint64) error {
	c := gw.d.connections.getByGnetID(gnetID)
	if c == nil {
		return ErrConnectionNotExist
	}

	return gw.d.Disconnect(c.Addr, ErrDisconnectRequestedByOperator)
}

// GetTrustConnections returns all trusted connections
func (gw *Gateway) GetTrustConnections() []string {
	return gw.d.pex.Trusted().ToAddrs()
}

// GetExchgConnection returns all connections to peers found through peer exchange
func (gw *Gateway) GetExchgConnection() []string {
	return gw.d.pex.RandomExchangeable(0).ToAddrs()
}

/* Blockchain & Transaction status */

// BlockchainProgress is the current blockchain syncing status
type BlockchainProgress struct {
	// Our current blockchain length
	Current uint64
	// Our best guess at true blockchain length
	Highest uint64
	// Individual blockchain length reports from peers
	Peers []PeerBlockchainHeight
}

// newBlockchainProgress creates BlockchainProgress from the local head blockchain sequence number
// and a list of remote peers
func newBlockchainProgress(headSeq uint64, conns []connection) *BlockchainProgress {
	peers := newPeerBlockchainHeights(conns)

	return &BlockchainProgress{
		Current: headSeq,
		Highest: EstimateBlockchainHeight(headSeq, peers),
		Peers:   peers,
	}
}

// PeerBlockchainHeight records blockchain height for an address
type PeerBlockchainHeight struct {
	Address string
	Height  uint64
}

func newPeerBlockchainHeights(conns []connection) []PeerBlockchainHeight {
	peers := make([]PeerBlockchainHeight, 0, len(conns))
	for _, c := range conns {
		if c.State != ConnectionStatePending {
			peers = append(peers, PeerBlockchainHeight{
				Address: c.Addr,
				Height:  c.Height,
			})
		}
	}
	return peers
}

// EstimateBlockchainHeight estimates the blockchain sync height.
// The highest height reported amongst all peers, and including the node itself, is returned.
func EstimateBlockchainHeight(headSeq uint64, peers []PeerBlockchainHeight) uint64 {
	for _, c := range peers {
		if c.Height > headSeq {
			headSeq = c.Height
		}
	}
	return headSeq
}

// GetBlockchainProgress returns a *BlockchainProgress
func (gw *Gateway) GetBlockchainProgress() (*BlockchainProgress, error) {
	headSeq, _, err := gw.v.HeadBkSeq()
	if err != nil {
		return nil, err
	}

	conns := gw.d.connections.all()

	return newBlockchainProgress(headSeq, conns), nil
}

// ResendUnconfirmedTxns resents all unconfirmed transactions, returning the txids
// of the transactions that were resent
func (gw *Gateway) ResendUnconfirmedTxns() ([]cipher.SHA256, error) {
	return gw.d.ResendUnconfirmedTxns()
}

// GetBlockchainMetadata returns a *visor.BlockchainMetadata
func (gw *Gateway) GetBlockchainMetadata() (*visor.BlockchainMetadata, error) {
	return gw.v.GetBlockchainMetadata()
}

// GetSignedBlockByHash returns the block by hash
func (gw *Gateway) GetSignedBlockByHash(hash cipher.SHA256) (*coin.SignedBlock, error) {
	return gw.v.GetSignedBlockByHash(hash)
}

// GetSignedBlockByHashVerbose returns the block by hash with verbose transaction inputs
func (gw *Gateway) GetSignedBlockByHashVerbose(hash cipher.SHA256) (*coin.SignedBlock, [][]visor.TransactionInput, error) {
	return gw.v.GetSignedBlockByHashVerbose(hash)
}

// GetSignedBlockBySeq returns block by seq
func (gw *Gateway) GetSignedBlockBySeq(seq uint64) (*coin.SignedBlock, error) {
	return gw.v.GetSignedBlockBySeq(seq)
}

// GetSignedBlockBySeqVerbose returns the block by seq with verbose transaction inputs
func (gw *Gateway) GetSignedBlockBySeqVerbose(seq uint64) (*coin.SignedBlock, [][]visor.TransactionInput, error) {
	return gw.v.GetSignedBlockBySeqVerbose(seq)
}

// GetBlocks returns blocks matching given block sequences
func (gw *Gateway) GetBlocks(seqs []uint64) ([]coin.SignedBlock, error) {
	return gw.v.GetBlocks(seqs)
}

// GetBlocksVerbose returns blocks matching given block sequences, with verbose transaction input data
func (gw *Gateway) GetBlocksVerbose(seqs []uint64) ([]coin.SignedBlock, [][][]visor.TransactionInput, error) {
	return gw.v.GetBlocksVerbose(seqs)
}

// GetBlocksInRange returns blocks between start and end, including start and end
func (gw *Gateway) GetBlocksInRange(start, end uint64) ([]coin.SignedBlock, error) {
	return gw.v.GetBlocksInRange(start, end)
}

// GetBlocksInRangeVerbose returns blocks between start and end, including start and end,
// and returns the blocks' verbose transaction input data
func (gw *Gateway) GetBlocksInRangeVerbose(start, end uint64) ([]coin.SignedBlock, [][][]visor.TransactionInput, error) {
	return gw.v.GetBlocksInRangeVerbose(start, end)
}

// GetLastBlocks get last N blocks
func (gw *Gateway) GetLastBlocks(num uint64) ([]coin.SignedBlock, error) {
	return gw.v.GetLastBlocks(num)
}

// GetLastBlocksVerbose get last N blocks with verbose transaction input data
func (gw *Gateway) GetLastBlocksVerbose(num uint64) ([]coin.SignedBlock, [][][]visor.TransactionInput, error) {
	return gw.v.GetLastBlocksVerbose(num)
}

// GetUnspentOutputsSummary gets unspent outputs and returns the filtered results,
// Note: all filters will be executed as the pending sequence in 'AND' mode.
func (gw *Gateway) GetUnspentOutputsSummary(filters []visor.OutputsFilter) (*visor.UnspentOutputsSummary, error) {
	return gw.v.GetUnspentOutputsSummary(filters)
}

// GetTransaction returns transaction by txid
func (gw *Gateway) GetTransaction(txid cipher.SHA256) (*visor.Transaction, error) {
	return gw.v.GetTransaction(txid)
}

// GetTransactionVerbose gets verbose transaction result by txid.
func (gw *Gateway) GetTransactionVerbose(txid cipher.SHA256) (*visor.Transaction, []visor.TransactionInput, error) {
	return gw.v.GetTransactionWithInputs(txid)
}

// InjectBroadcastTransaction injects transaction to the unconfirmed pool and broadcasts it.
// If the transaction violates either hard or soft constraints, it is not broadcast.
// This method is to be used by user-initiated transaction injections.
// For transactions received over the network, use daemon.injectTransaction and check the result to
// decide on repropagation.
func (gw *Gateway) InjectBroadcastTransaction(txn coin.Transaction) error {
	return gw.v.WithUpdateTx("gateway.InjectBroadcastTransaction", func(tx *dbutil.Tx) error {
		_, head, inputs, err := gw.v.InjectUserTransactionTx(tx, txn)
		if err != nil {
			logger.WithError(err).Error("InjectUserTransactionTx failed")
			return err
		}

		if err := gw.d.BroadcastUserTransaction(txn, head, inputs); err != nil {
			logger.WithError(err).Error("BroadcastUserTransaction failed")
			return err
		}

		return nil
	})
}

// GetVerboseTransactionsForAddress returns transactions and their verbose input data for a given address.
// These transactions include confirmed and unconfirmed transactions
func (gw *Gateway) GetVerboseTransactionsForAddress(a cipher.Address) ([]visor.Transaction, [][]visor.TransactionInput, error) {
	return gw.v.GetVerboseTransactionsForAddress(a)
}

// GetTransactions returns transactions filtered by zero or more visor.TxFilter
func (gw *Gateway) GetTransactions(flts []visor.TxFilter) ([]visor.Transaction, error) {
	return gw.v.GetTransactions(flts)
}

// GetTransactionsVerbose returns transactions filtered by zero or more visor.TxFilter
func (gw *Gateway) GetTransactionsVerbose(flts []visor.TxFilter) ([]visor.Transaction, [][]visor.TransactionInput, error) {
	return gw.v.GetTransactionsWithInputs(flts)
}

// GetUxOutByID gets UxOut by hash id.
func (gw *Gateway) GetUxOutByID(id cipher.SHA256) (*historydb.UxOut, error) {
	return gw.v.GetUxOutByID(id)
}

// GetSpentOutputsForAddresses gets all the spent outputs of a set of addresses
func (gw *Gateway) GetSpentOutputsForAddresses(addresses []cipher.Address) ([][]historydb.UxOut, error) {
	return gw.v.GetSpentOutputsForAddresses(addresses)
}

// GetAllUnconfirmedTransactions returns all unconfirmed transactions
func (gw *Gateway) GetAllUnconfirmedTransactions() ([]visor.UnconfirmedTransaction, error) {
	return gw.v.GetAllUnconfirmedTransactions()
}

// GetAllUnconfirmedTransactionsVerbose returns all unconfirmed transactions with verbose transaction inputs
func (gw *Gateway) GetAllUnconfirmedTransactionsVerbose() ([]visor.UnconfirmedTransaction, [][]visor.TransactionInput, error) {
	return gw.v.GetAllUnconfirmedTransactionsVerbose()
}

// GetUnconfirmedTransactions returns addresses related unconfirmed transactions
func (gw *Gateway) GetUnconfirmedTransactions(addrs []cipher.Address) ([]visor.UnconfirmedTransaction, error) {
	return gw.v.GetUnconfirmedTransactions(visor.SendsToAddresses(addrs))
}

// CreateTransaction creates a transaction based upon parameters in wallet.CreateTransactionParams
func (gw *Gateway) CreateTransaction(params wallet.CreateTransactionParams) (*coin.Transaction, []wallet.UxBalance, error) {
	if !gw.Config.EnableWalletAPI {
		return nil, nil, wallet.ErrWalletAPIDisabled
	}

	return gw.v.CreateTransaction(params)
}

// SignTransaction signs an unsigned transaction using a wallet. Specific inputs may be signed by specifying signIndexes.
// If signIndexes is empty, all inputs will be signed.
func (gw *Gateway) SignTransaction(wltName string, password []byte, txn *coin.Transaction, signIndexes []int) ([]wallet.UxBalance, error) {
	if !gw.Config.EnableWalletAPI {
		return nil, nil, wallet.ErrWalletAPIDisabled
	}

	return gw.v.SignTransaction(wltName, password, txn, signIndexes)
}

// CreateWallet creates wallet
func (gw *Gateway) CreateWallet(wltName string, options wallet.Options) (*wallet.Wallet, error) {
	if !gw.Config.EnableWalletAPI {
		return nil, wallet.ErrWalletAPIDisabled
	}

	return gw.v.Wallets.CreateWallet(wltName, options, gw.v)
}

// RecoverWallet recovers an encrypted wallet from seed
func (gw *Gateway) RecoverWallet(wltName, seed string, password []byte) (*wallet.Wallet, error) {
	if !gw.Config.EnableWalletAPI {
		return nil, wallet.ErrWalletAPIDisabled
	}

	return gw.v.Wallets.RecoverWallet(wltName, seed, password)
}

// EncryptWallet encrypts the wallet
func (gw *Gateway) EncryptWallet(wltName string, password []byte) (*wallet.Wallet, error) {
	if !gw.Config.EnableWalletAPI {
		return nil, wallet.ErrWalletAPIDisabled
	}

	return gw.v.Wallets.EncryptWallet(wltName, password)
}

// DecryptWallet decrypts wallet
func (gw *Gateway) DecryptWallet(wltID string, password []byte) (*wallet.Wallet, error) {
	if !gw.Config.EnableWalletAPI {
		return nil, wallet.ErrWalletAPIDisabled
	}

	return gw.v.Wallets.DecryptWallet(wltID, password)
}

// GetWalletBalance returns balance pairs of specific wallet
func (gw *Gateway) GetWalletBalance(wltID string) (wallet.BalancePair, wallet.AddressBalances, error) {
	if !gw.Config.EnableWalletAPI {
		return wallet.BalancePair{}, wallet.AddressBalances{}, wallet.ErrWalletAPIDisabled
	}

	return gw.v.GetWalletBalance(wltID)
}

// GetBalanceOfAddrs gets balance of given addresses
func (gw *Gateway) GetBalanceOfAddrs(addrs []cipher.Address) ([]wallet.BalancePair, error) {
	return gw.v.GetBalanceOfAddrs(addrs)
}

// GetWalletDir returns path for storing wallet files
func (gw *Gateway) GetWalletDir() (string, error) {
	if !gw.Config.EnableWalletAPI {
		return "", wallet.ErrWalletAPIDisabled
	}
	return gw.v.Config.WalletDirectory, nil
}

// NewAddresses generate addresses in given wallet
func (gw *Gateway) NewAddresses(wltID string, password []byte, n uint64) ([]cipher.Address, error) {
	if !gw.Config.EnableWalletAPI {
		return nil, wallet.ErrWalletAPIDisabled
	}

	return gw.v.Wallets.NewAddresses(wltID, password, n)
}

// UpdateWalletLabel updates the label of wallet
func (gw *Gateway) UpdateWalletLabel(wltID, label string) error {
	if !gw.Config.EnableWalletAPI {
		return wallet.ErrWalletAPIDisabled
	}

	return gw.v.Wallets.UpdateWalletLabel(wltID, label)
}

// GetWallet returns wallet by id
func (gw *Gateway) GetWallet(wltID string) (*wallet.Wallet, error) {
	if !gw.Config.EnableWalletAPI {
		return nil, wallet.ErrWalletAPIDisabled
	}

	return gw.v.Wallets.GetWallet(wltID)
}

// GetWallets returns wallets
func (gw *Gateway) GetWallets() (wallet.Wallets, error) {
	if !gw.Config.EnableWalletAPI {
		return nil, wallet.ErrWalletAPIDisabled
	}

	return gw.v.Wallets.GetWallets()
}

// GetWalletUnconfirmedTransactions returns all unconfirmed transactions in given wallet
func (gw *Gateway) GetWalletUnconfirmedTransactions(wltID string) ([]visor.UnconfirmedTransaction, error) {
	if !gw.Config.EnableWalletAPI {
		return nil, wallet.ErrWalletAPIDisabled
	}

	return gw.v.GetWalletUnconfirmedTransactions(wltID)
}

// GetWalletUnconfirmedTransactionsVerbose returns all unconfirmed transactions in given wallet
func (gw *Gateway) GetWalletUnconfirmedTransactionsVerbose(wltID string) ([]visor.UnconfirmedTransaction, [][]visor.TransactionInput, error) {
	if !gw.Config.EnableWalletAPI {
		return nil, nil, wallet.ErrWalletAPIDisabled
	}

	return gw.v.GetWalletUnconfirmedTransactionsVerbose(wltID)
}

// UnloadWallet removes wallet of given id from memory.
func (gw *Gateway) UnloadWallet(id string) error {
	if !gw.Config.EnableWalletAPI {
		return wallet.ErrWalletAPIDisabled
	}

	return gw.v.Wallets.Remove(id)
}

// GetWalletSeed returns seed of wallet of given id,
// returns wallet.ErrWalletNotEncrypted if the wallet is not encrypted.
func (gw *Gateway) GetWalletSeed(id string, password []byte) (string, error) {
	if !gw.Config.EnableWalletAPI {
		return "", wallet.ErrWalletAPIDisabled
	}

	return gw.v.Wallets.GetWalletSeed(id, password)
}

// GetRichlist returns rich list as desc order.
func (gw *Gateway) GetRichlist(includeDistribution bool) (visor.Richlist, error) {
	rbOuts, err := gw.GetUnspentOutputsSummary(nil)
	if err != nil {
		return nil, err
	}

	// Build a map from addresses to total coins held
	allAccounts := map[cipher.Address]uint64{}
	for _, out := range rbOuts.Confirmed {
		if _, ok := allAccounts[out.Body.Address]; ok {
			var err error
			allAccounts[out.Body.Address], err = coin.AddUint64(allAccounts[out.Body.Address], out.Body.Coins)
			if err != nil {
				return nil, err
			}
		} else {
			allAccounts[out.Body.Address] = out.Body.Coins
		}
	}

	lockedAddrs := params.GetLockedDistributionAddressesDecoded()
	addrsMap := make(map[cipher.Address]struct{}, len(lockedAddrs))
	for _, a := range lockedAddrs {
		addrsMap[a] = struct{}{}
	}

	richlist, err := visor.NewRichlist(allAccounts, addrsMap)
	if err != nil {
		return nil, err
	}

	if !includeDistribution {
		unlockedAddrs := params.GetUnlockedDistributionAddressesDecoded()
		for _, a := range unlockedAddrs {
			addrsMap[a] = struct{}{}
		}
		richlist = richlist.FilterAddresses(addrsMap)
	}

	return richlist, nil
}

// GetAddressCount returns count number of unique address with uxouts > 0.
func (gw *Gateway) GetAddressCount() (uint64, error) {
	return gw.v.AddressCount()
}

// Health is returned by the /health endpoint
type Health struct {
	BlockchainMetadata   visor.BlockchainMetadata
	OutgoingConnections  int
	IncomingConnections  int
	Uptime               time.Duration
	UnconfirmedVerifyTxn params.VerifyTxn
	StartedAt            time.Time
}

// GetHealth returns statistics about the running node
func (gw *Gateway) GetHealth() (*Health, error) {
	metadata, err := gw.v.GetBlockchainMetadata()
	if err != nil {
		return nil, err
	}

	conns, err := gw.getConnections(func(c Connection) bool {
		return c.State != ConnectionStatePending
	})
	if err != nil {
		return nil, err
	}

	outgoingConns := 0
	incomingConns := 0
	for _, c := range conns {
		if c.Outgoing {
			outgoingConns++
		} else {
			incomingConns++
		}
	}

	return &Health{
		BlockchainMetadata:   *metadata,
		OutgoingConnections:  outgoingConns,
		IncomingConnections:  incomingConns,
		Uptime:               time.Since(gw.v.StartedAt),
		UnconfirmedVerifyTxn: gw.d.Config.UnconfirmedVerifyTxn,
		StartedAt:            gw.v.StartedAt,
	}, nil
}

// VerifyTxnVerbose verifies an isolated transaction and returns []wallet.UxBalance of
// transaction inputs, whether the transaction is confirmed and error if any
func (gw *Gateway) VerifyTxnVerbose(txn *coin.Transaction) ([]wallet.UxBalance, bool, error) {
	return gw.v.VerifyTxnVerbose(txn)
}
