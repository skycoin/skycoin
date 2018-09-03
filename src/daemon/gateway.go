package daemon

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/daemon/strand"
	"github.com/skycoin/skycoin/src/readable"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/wallet"

	"github.com/skycoin/skycoin/src/visor/historydb"
)

// GatewayConfig configuration set of gateway.
type GatewayConfig struct {
	BufferSize      int
	EnableWalletAPI bool
}

// NewGatewayConfig create and init an GatewayConfig
func NewGatewayConfig() GatewayConfig {
	return GatewayConfig{
		BufferSize:      32,
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
	// Requests are queued on this channel
	requests chan strand.Request
	quit     chan struct{}
}

// NewGateway create and init an Gateway instance.
func NewGateway(c GatewayConfig, d *Daemon) *Gateway {
	return &Gateway{
		Config:   c,
		d:        d,
		v:        d.visor,
		requests: make(chan strand.Request, c.BufferSize),
		quit:     make(chan struct{}),
	}
}

// Shutdown closes the Gateway
func (gw *Gateway) Shutdown() {
	close(gw.quit)
	// wait for strand to complete
	gw.strand("wait-shutdown", func() {})
}

func (gw *Gateway) strand(name string, f func()) {
	name = fmt.Sprintf("daemon.Gateway.%s", name)
	if err := strand.Strand(logger, gw.requests, name, func() error {
		f()
		return nil
	}, gw.quit, nil); err != nil {
		logger.WithError(err).Error("Gateway.strand.Strand failed")
	}
}

// Connection a connection's state within the daemon
type Connection struct {
	ID           int    `json:"id"`
	Addr         string `json:"address"`
	LastSent     int64  `json:"last_sent"`
	LastReceived int64  `json:"last_received"`
	// Whether the connection is from us to them (true, outgoing),
	// or from them to us (false, incoming)
	Outgoing bool `json:"outgoing"`
	// Whether the client has identified their version, mirror etc
	Introduced bool   `json:"introduced"`
	Mirror     uint32 `json:"mirror"`
	ListenPort uint16 `json:"listen_port"`
}

// Connections an array of connections
// Arrays must be wrapped in structs to avoid certain javascript exploits
type Connections struct {
	Connections []*Connection `json:"connections"`
}

// GetConnections returns a *Connections
func (gw *Gateway) GetConnections() *Connections {
	var conns *Connections
	gw.strand("GetConnections", func() {
		conns = gw.getConnections()
	})
	return conns
}

func (gw *Gateway) getConnections() *Connections {
	if gw.d.pool.Pool == nil {
		return nil
	}

	n, err := gw.d.pool.Pool.Size()
	if err != nil {
		logger.Error(err)
		return nil
	}

	conns := make([]*Connection, 0, n)
	cs, err := gw.d.pool.Pool.GetConnections()
	if err != nil {
		logger.Error(err)
		return nil
	}

	for _, c := range cs {
		if c.Solicited {
			conn := gw.getConnection(c.Addr())
			if conn != nil {
				conns = append(conns, conn)
			}
		}
	}

	// Sort connnections by IP address
	sort.Slice(conns, func(i, j int) bool {
		return strings.Compare(conns[i].Addr, conns[j].Addr) < 0
	})

	return &Connections{Connections: conns}

}

// GetDefaultConnections returns default connections
func (gw *Gateway) GetDefaultConnections() []string {
	var conns []string
	gw.strand("GetDefaultConnections", func() {
		conns = make([]string, len(gw.d.Config.DefaultConnections))
		copy(conns[:], gw.d.Config.DefaultConnections[:])
	})
	return conns
}

// GetConnection returns a *Connection of specific address
func (gw *Gateway) GetConnection(addr string) *Connection {
	var conn *Connection
	gw.strand("GetConnection", func() {
		conn = gw.getConnection(addr)
	})
	return conn
}

func (gw *Gateway) getConnection(addr string) *Connection {
	if gw.d.pool.Pool == nil {
		return nil
	}

	c, err := gw.d.pool.Pool.GetConnection(addr)
	if err != nil {
		logger.Error(err)
		return nil
	}

	if c == nil {
		return nil
	}

	mirror, exist := gw.d.connectionMirrors.Get(addr)
	if !exist {
		return nil
	}

	return &Connection{
		ID:           c.ID,
		Addr:         addr,
		LastSent:     c.LastSent.Unix(),
		LastReceived: c.LastReceived.Unix(),
		Outgoing:     gw.d.outgoingConnections.Get(addr),
		Introduced:   !gw.d.needsIntro(addr),
		Mirror:       mirror,
		ListenPort:   gw.d.GetListenPort(addr),
	}
}

// GetTrustConnections returns all trusted connections,
// including private and public
func (gw *Gateway) GetTrustConnections() []string {
	var conn []string
	gw.strand("GetTrustConnections", func() {
		conn = gw.d.pex.Trusted().ToAddrs()
	})
	return conn
}

// GetExchgConnection returns all exchangeable connections,
// including private and public
func (gw *Gateway) GetExchgConnection() []string {
	var conn []string
	gw.strand("GetExchgConnection", func() {
		conn = gw.d.pex.RandomExchangeable(0).ToAddrs()
	})
	return conn
}

/* Blockchain & Transaction status */

// BlockchainProgress current sync blockchain status
type BlockchainProgress struct {
	// Our current blockchain length
	Current uint64 `json:"current"`
	// Our best guess at true blockchain length
	Highest uint64                 `json:"highest"`
	Peers   []PeerBlockchainHeight `json:"peers"`
}

// GetBlockchainProgress returns a *BlockchainProgress
func (gw *Gateway) GetBlockchainProgress() (*BlockchainProgress, error) {
	var bcp *BlockchainProgress
	var err error
	gw.strand("GetBlockchainProgress", func() {
		var headSeq uint64
		headSeq, _, err = gw.d.visor.HeadBkSeq()
		if err != nil {
			return
		}

		bcp = &BlockchainProgress{
			Current: headSeq,
			Highest: gw.d.Heights.Estimate(headSeq),
			Peers:   gw.d.Heights.All(),
		}
	})

	if err != nil {
		return nil, err
	}

	return bcp, nil
}

// ResendResult rebroadcast tx result
type ResendResult struct {
	Txids []string `json:"txids"` // transaction id
}

// ResendUnconfirmedTxns resents all unconfirmed transactions
func (gw *Gateway) ResendUnconfirmedTxns() (*ResendResult, error) {
	var hashes []cipher.SHA256
	var err error
	gw.strand("ResendUnconfirmedTxns", func() {
		hashes, err = gw.d.ResendUnconfirmedTxns()
	})

	if err != nil {
		return nil, err
	}

	var rlt ResendResult
	for _, txid := range hashes {
		rlt.Txids = append(rlt.Txids, txid.Hex())
	}
	return &rlt, nil
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
func (gw *Gateway) GetSignedBlockByHash(hash cipher.SHA256) (*coin.SignedBlock, error) {
	var b *coin.SignedBlock
	var err error
	gw.strand("GetSignedBlockByHash", func() {
		b, err = gw.v.GetSignedBlockByHash(hash)
	})
	return b, err
}

// GetBlockByHashVerbose returns the block by hash
func (gw *Gateway) GetBlockByHashVerbose(hash cipher.SHA256) (*readable.BlockVerbose, error) {
	var b *coin.SignedBlock
	var inputs [][]visor.TransactionInput
	var err error
	gw.strand("GetBlockByHashVerbose", func() {
		b, inputs, err = gw.v.GetBlockByHashVerbose(hash)
	})

	if err != nil {
		return nil, err
	}

	if b == nil {
		return nil, err
	}

	return readable.NewBlockVerbose(&b.Block, inputs)
}

// GetSignedBlockBySeq returns block by seq
func (gw *Gateway) GetSignedBlockBySeq(seq uint64) (*coin.SignedBlock, error) {
	var b *coin.SignedBlock
	var err error
	gw.strand("GetSignedBlockBySeq", func() {
		b, err = gw.v.GetSignedBlockBySeq(seq)
	})
	return b, err
}

// GetBlockBySeqVerbose returns the block by seq
func (gw *Gateway) GetBlockBySeqVerbose(seq uint64) (*readable.BlockVerbose, error) {
	var b *coin.SignedBlock
	var inputs [][]visor.TransactionInput
	var err error
	gw.strand("GetBlockBySeqVerbose", func() {
		b, inputs, err = gw.v.GetBlockBySeqVerbose(seq)
	})

	if err != nil {
		return nil, err
	}

	if b == nil {
		return nil, err
	}

	return readable.NewBlockVerbose(&b.Block, inputs)
}

// GetBlocks returns a *readable.Blocks
func (gw *Gateway) GetBlocks(start, end uint64) (*readable.Blocks, error) {
	var blocks []coin.SignedBlock
	var err error

	gw.strand("GetBlocks", func() {
		blocks, err = gw.v.GetBlocks(start, end)
	})
	if err != nil {
		return nil, err
	}

	return readable.NewBlocks(blocks)
}

// GetBlocksVerbose returns a *readable.BlocksVerbose
func (gw *Gateway) GetBlocksVerbose(start, end uint64) (*readable.BlocksVerbose, error) {
	var blocks []coin.SignedBlock
	var inputs [][][]visor.TransactionInput
	var err error

	gw.strand("GetBlocksVerbose", func() {
		blocks, inputs, err = gw.v.GetBlocksVerbose(start, end)
	})
	if err != nil {
		return nil, err
	}

	return readable.NewBlocksVerbose(blocks, inputs)
}

// GetBlocksInDepth returns blocks in different depth
func (gw *Gateway) GetBlocksInDepth(vs []uint64) (*readable.Blocks, error) {
	blocks := []coin.SignedBlock{}
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

	return readable.NewBlocks(blocks)
}

// GetLastBlocks get last N blocks
func (gw *Gateway) GetLastBlocks(num uint64) (*readable.Blocks, error) {
	var blocks []coin.SignedBlock
	var err error

	gw.strand("GetLastBlocks", func() {
		blocks, err = gw.v.GetLastBlocks(num)
	})
	if err != nil {
		return nil, err
	}

	return readable.NewBlocks(blocks)
}

// GetLastBlocksVerbose get last N blocks with verbose transaction input data
func (gw *Gateway) GetLastBlocksVerbose(num uint64) (*readable.BlocksVerbose, error) {
	var blocks []coin.SignedBlock
	var inputs [][][]visor.TransactionInput
	var err error

	gw.strand("GetLastBlocksVerbose", func() {
		blocks, inputs, err = gw.v.GetLastBlocksVerbose(num)
	})
	if err != nil {
		return nil, err
	}

	return readable.NewBlocksVerbose(blocks, inputs)
}

// OutputsFilter used as optional arguments in GetUnspentOutputs method
type OutputsFilter func(outputs coin.UxArray) coin.UxArray

// GetUnspentOutputs gets unspent outputs and returns the filtered results,
// Note: all filters will be executed as the pending sequence in 'AND' mode.
func (gw *Gateway) GetUnspentOutputs(filters ...OutputsFilter) (*readable.OutputSet, error) {
	// unspent outputs
	var unspentOutputs []coin.UxOut
	// unconfirmed spending outputs
	var uncfmSpendingOutputs coin.UxArray
	// unconfirmed incoming outputs
	var uncfmIncomingOutputs coin.UxArray
	var head *coin.SignedBlock
	var err error
	gw.strand("GetUnspentOutputs", func() {
		head, err = gw.v.GetHeadBlock()
		if err != nil {
			err = fmt.Errorf("v.GetHeadBlock failed: %v", err)
			return
		}

		unspentOutputs, err = gw.v.GetAllUnspentOutputs()
		if err != nil {
			err = fmt.Errorf("v.GetAllUnspentOutputs failed: %v", err)
			return
		}

		uncfmSpendingOutputs, err = gw.v.UnconfirmedSpendingOutputs()
		if err != nil {
			err = fmt.Errorf("v.UnconfirmedSpendingOutputs failed: %v", err)
			return
		}

		uncfmIncomingOutputs, err = gw.v.UnconfirmedIncomingOutputs()
		if err != nil {
			err = fmt.Errorf("v.UnconfirmedIncomingOutputs failed: %v", err)
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

	outputSet := readable.OutputSet{}
	outputSet.HeadOutputs, err = readable.NewOutputs(head.Time(), unspentOutputs)
	if err != nil {
		return nil, err
	}

	outputSet.OutgoingOutputs, err = readable.NewOutputs(head.Time(), uncfmSpendingOutputs)
	if err != nil {
		return nil, err
	}

	outputSet.IncomingOutputs, err = readable.NewOutputs(head.Time(), uncfmIncomingOutputs)
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
func (gw *Gateway) GetTransaction(txid cipher.SHA256) (*visor.Transaction, error) {
	var txn *visor.Transaction
	var err error

	gw.strand("GetTransaction", func() {
		txn, err = gw.v.GetTransaction(txid)
	})

	return txn, err
}

// GetTransactionWithStatus gets transaction result by txid.
func (gw *Gateway) GetTransactionWithStatus(txid cipher.SHA256) (*readable.TransactionWithStatus, error) {
	var txn *visor.Transaction
	var err error

	gw.strand("GetTransactionWithStatus", func() {
		txn, err = gw.v.GetTransaction(txid)
	})

	if err != nil {
		return nil, err
	}

	if txn == nil {
		return nil, nil
	}

	return readable.NewTransactionWithStatus(txn)
}

// GetTransactionWithStatusVerbose gets verbose transaction result by txid.
func (gw *Gateway) GetTransactionWithStatusVerbose(txid cipher.SHA256) (*readable.TransactionWithStatusVerbose, error) {
	var txn *visor.Transaction
	var inputs []visor.TransactionInput
	var err error

	gw.strand("GetTransactionWithStatusVerbose", func() {
		txn, inputs, err = gw.v.GetTransactionWithInputs(txid)
	})

	if err != nil {
		return nil, err
	}

	if txn == nil {
		return nil, nil
	}

	return readable.NewTransactionWithStatusVerbose(txn, inputs)
}

// InjectBroadcastTransaction injects and broadcasts a transaction
func (gw *Gateway) InjectBroadcastTransaction(txn coin.Transaction) error {
	var err error
	gw.strand("InjectBroadcastTransaction", func() {
		err = gw.d.InjectBroadcastTransaction(txn)
	})
	return err
}

// GetVerboseTransactionsForAddress returns []readable.TransactionVerbose for a given address.
// These transactions include confirmed and unconfirmed transactions
func (gw *Gateway) GetVerboseTransactionsForAddress(a cipher.Address) ([]readable.TransactionVerbose, error) {
	var err error
	var txns []visor.Transaction
	var inputs [][]visor.TransactionInput

	gw.strand("GetVerboseTransactionsForAddress", func() {
		txns, inputs, err = gw.v.GetVerboseTransactionsForAddress(a)
	})

	if err != nil {
		return nil, err
	}

	vb := make([]readable.TransactionVerbose, len(txns))
	for i, txn := range txns {
		v, err := readable.NewTransactionVerbose(txn, inputs[i])
		if err != nil {
			return nil, err
		}

		vb[i] = v
	}

	return vb, nil
}

// GetTransactions returns transactions filtered by zero or more visor.TxFilter
func (gw *Gateway) GetTransactions(flts []visor.TxFilter) ([]visor.Transaction, error) {
	var txns []visor.Transaction
	var err error
	gw.strand("GetTransactions", func() {
		txns, err = gw.v.GetTransactions(flts)
	})
	return txns, err
}

// GetTransactionsWithStatus returns transactions filtered by zero or more visor.TxFilter
func (gw *Gateway) GetTransactionsWithStatus(flts []visor.TxFilter) (*readable.TransactionsWithStatus, error) {
	var txns []visor.Transaction
	var err error
	gw.strand("GetTransactionsWithStatus", func() {
		txns, err = gw.v.GetTransactions(flts)
	})
	if err != nil {
		return nil, err
	}

	return readable.NewTransactionsWithStatus(txns)
}

// GetTransactionsWithStatusVerbose returns transactions filtered by zero or more visor.TxFilter
func (gw *Gateway) GetTransactionsWithStatusVerbose(flts []visor.TxFilter) (*readable.TransactionsWithStatusVerbose, error) {
	var txns []visor.Transaction
	var inputs [][]visor.TransactionInput
	var err error
	gw.strand("GetTransactionsWithStatusVerbose", func() {
		txns, inputs, err = gw.v.GetTransactionsWithInputs(flts)
	})
	if err != nil {
		return nil, err
	}

	return readable.NewTransactionsWithStatusVerbose(txns, inputs)
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
func (gw *Gateway) GetAddrUxOuts(addresses []cipher.Address) ([]*historydb.UxOut, error) {
	var uxOuts []*historydb.UxOut
	var err error

	gw.strand("GetAddrUxOuts", func() {
		for _, addr := range addresses {
			var result []*historydb.UxOut
			result, err = gw.v.GetAddrUxOuts(addr)
			if err != nil {
				return
			}

			uxOuts = append(uxOuts, result...)
		}
	})

	if err != nil {
		return nil, err
	}

	return uxOuts, nil
}

// GetTimeNow returns the current Unix time
func (gw *Gateway) GetTimeNow() uint64 {
	return uint64(time.Now().UTC().Unix())
}

// GetAllUnconfirmedTxns returns all unconfirmed transactions
func (gw *Gateway) GetAllUnconfirmedTxns() ([]visor.UnconfirmedTransaction, error) {
	var txns []visor.UnconfirmedTransaction
	var err error
	gw.strand("GetAllUnconfirmedTxns", func() {
		txns, err = gw.v.GetAllUnconfirmedTxns()
	})
	return txns, err
}

// GetAllUnconfirmedTxnsVerbose returns all unconfirmed transactions with verbose transaction inputs
func (gw *Gateway) GetAllUnconfirmedTxnsVerbose() ([]readable.UnconfirmedTransactionVerbose, error) {
	var txns []visor.UnconfirmedTransaction
	var inputs [][]visor.TransactionInput
	var err error
	gw.strand("GetAllUnconfirmedTxnsVerbose", func() {
		txns, inputs, err = gw.v.GetAllUnconfirmedTxnsVerbose()
	})
	if err != nil {
		return nil, err
	}

	vb := make([]readable.UnconfirmedTransactionVerbose, len(txns))
	for i, txn := range txns {
		v, err := readable.NewUnconfirmedTransactionVerbose(&txn, inputs[i])
		if err != nil {
			return nil, err
		}

		vb[i] = *v
	}

	return vb, nil
}

// GetUnconfirmedTxns returns addresses related unconfirmed transactions
func (gw *Gateway) GetUnconfirmedTxns(addrs []cipher.Address) ([]visor.UnconfirmedTransaction, error) {
	var txns []visor.UnconfirmedTransaction
	var err error
	gw.strand("GetUnconfirmedTxns", func() {
		txns, err = gw.v.GetUnconfirmedTxns(visor.SendsToAddresses(addrs))
	})
	return txns, err
}

// Spend spends coins from given wallet and broadcast it,
// set password as nil if wallet is not encrypted, otherwise the password must be provied.
// return transaction or error.
func (gw *Gateway) Spend(wltID string, password []byte, coins uint64, dest cipher.Address) (*coin.Transaction, error) {
	if !gw.Config.EnableWalletAPI {
		return nil, wallet.ErrWalletAPIDisabled
	}

	var txn *coin.Transaction
	var err error
	gw.strand("Spend", func() {
		txn, err = gw.v.CreateTransactionDeprecated(wltID, password, coins, dest)
		if err != nil {
			return
		}

		// Inject transaction
		err = gw.d.InjectBroadcastTransaction(*txn)
		if err != nil {
			logger.Errorf("Inject transaction failed: %v", err)
			return
		}
	})

	if err != nil {
		return nil, err
	}

	return txn, nil
}

// CreateTransaction creates a transaction based upon parameters in wallet.CreateTransactionParams
func (gw *Gateway) CreateTransaction(params wallet.CreateTransactionParams) (*coin.Transaction, []wallet.UxBalance, error) {
	if !gw.Config.EnableWalletAPI {
		return nil, nil, wallet.ErrWalletAPIDisabled
	}

	var txn *coin.Transaction
	var inputs []wallet.UxBalance
	var err error

	gw.strand("CreateTransaction", func() {
		txn, inputs, err = gw.v.CreateTransaction(params)
	})

	if err != nil {
		return nil, nil, err
	}

	return txn, inputs, err
}

// CreateWallet creates wallet
func (gw *Gateway) CreateWallet(wltName string, options wallet.Options) (*wallet.Wallet, error) {
	if !gw.Config.EnableWalletAPI {
		return nil, wallet.ErrWalletAPIDisabled
	}

	var wlt *wallet.Wallet
	var err error
	gw.strand("CreateWallet", func() {
		wlt, err = gw.v.Wallets.CreateWallet(wltName, options, gw.v)
	})
	return wlt, err
}

// EncryptWallet encrypts the wallet
func (gw *Gateway) EncryptWallet(wltName string, password []byte) (*wallet.Wallet, error) {
	if !gw.Config.EnableWalletAPI {
		return nil, wallet.ErrWalletAPIDisabled
	}

	var err error
	var w *wallet.Wallet
	gw.strand("EncryptWallet", func() {
		w, err = gw.v.Wallets.EncryptWallet(wltName, password)
	})
	return w, err
}

// DecryptWallet decrypts wallet
func (gw *Gateway) DecryptWallet(wltID string, password []byte) (*wallet.Wallet, error) {
	if !gw.Config.EnableWalletAPI {
		return nil, wallet.ErrWalletAPIDisabled
	}

	var err error
	var w *wallet.Wallet
	gw.strand("DecryptWallet", func() {
		w, err = gw.v.Wallets.DecryptWallet(wltID, password)
	})
	return w, err
}

// GetWalletBalance returns balance pairs of specific wallet
func (gw *Gateway) GetWalletBalance(wltID string) (wallet.BalancePair, wallet.AddressBalance, error) {
	var addressBalances wallet.AddressBalance
	var walletBalance wallet.BalancePair
	if !gw.Config.EnableWalletAPI {
		return walletBalance, addressBalances, wallet.ErrWalletAPIDisabled
	}

	var err error
	gw.strand("GetWalletBalance", func() {
		var addrs []cipher.Address
		addrs, err = gw.v.Wallets.GetAddresses(wltID)
		if err != nil {
			return
		}

		// get list of address balances
		addrsBalanceList, err := gw.v.GetBalanceOfAddrs(addrs)
		if err != nil {
			return
		}

		// create map of address to balance
		addressBalances = make(wallet.AddressBalance, len(addrs))
		for idx, addr := range addrs {
			addressBalances[addr.String()] = addrsBalanceList[idx]
		}

		// compute the sum of all addresses
		for _, addrBalance := range addressBalances {
			// compute confirmed balance
			walletBalance.Confirmed.Coins, err = coin.AddUint64(walletBalance.Confirmed.Coins, addrBalance.Confirmed.Coins)
			if err != nil {
				return
			}
			walletBalance.Confirmed.Hours, err = coin.AddUint64(walletBalance.Confirmed.Hours, addrBalance.Confirmed.Hours)
			if err != nil {
				return
			}

			// compute predicted balance
			walletBalance.Predicted.Coins, err = coin.AddUint64(walletBalance.Predicted.Coins, addrBalance.Predicted.Coins)
			if err != nil {
				return
			}
			walletBalance.Predicted.Hours, err = coin.AddUint64(walletBalance.Predicted.Hours, addrBalance.Predicted.Hours)
			if err != nil {
				return
			}
		}
	})

	return walletBalance, addressBalances, err
}

// GetBalanceOfAddrs gets balance of given addresses
func (gw *Gateway) GetBalanceOfAddrs(addrs []cipher.Address) ([]wallet.BalancePair, error) {
	var balance []wallet.BalancePair
	var err error

	gw.strand("GetBalanceOfAddrs", func() {
		balance, err = gw.v.GetBalanceOfAddrs(addrs)
	})

	if err != nil {
		return nil, err
	}

	return balance, nil
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

	var addrs []cipher.Address
	var err error
	gw.strand("NewAddresses", func() {
		addrs, err = gw.v.Wallets.NewAddresses(wltID, password, n)
	})
	return addrs, err
}

// UpdateWalletLabel updates the label of wallet
func (gw *Gateway) UpdateWalletLabel(wltID, label string) error {
	if !gw.Config.EnableWalletAPI {
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
	if !gw.Config.EnableWalletAPI {
		return nil, wallet.ErrWalletAPIDisabled
	}

	var w *wallet.Wallet
	var err error
	gw.strand("GetWallet", func() {
		w, err = gw.v.Wallets.GetWallet(wltID)
	})
	return w, err
}

// GetWallets returns wallets
func (gw *Gateway) GetWallets() (wallet.Wallets, error) {
	if !gw.Config.EnableWalletAPI {
		return nil, wallet.ErrWalletAPIDisabled
	}

	var w wallet.Wallets
	var err error
	gw.strand("GetWallets", func() {
		w, err = gw.v.Wallets.GetWallets()
	})
	return w, err
}

// GetWalletUnconfirmedTxns returns all unconfirmed transactions in given wallet
func (gw *Gateway) GetWalletUnconfirmedTxns(wltID string) ([]visor.UnconfirmedTransaction, error) {
	if !gw.Config.EnableWalletAPI {
		return nil, wallet.ErrWalletAPIDisabled
	}

	var txns []visor.UnconfirmedTransaction
	var err error
	gw.strand("GetWalletUnconfirmedTxns", func() {
		var addrs []cipher.Address
		addrs, err = gw.v.Wallets.GetAddresses(wltID)
		if err != nil {
			return
		}

		txns, err = gw.v.GetUnconfirmedTxns(visor.SendsToAddresses(addrs))
	})

	return txns, err
}

// GetWalletUnconfirmedTxnsVerbose returns all unconfirmed transactions in given wallet
func (gw *Gateway) GetWalletUnconfirmedTxnsVerbose(wltID string) ([]readable.UnconfirmedTransactionVerbose, error) {
	if !gw.Config.EnableWalletAPI {
		return nil, wallet.ErrWalletAPIDisabled
	}

	var txns []visor.UnconfirmedTransaction
	var inputs [][]visor.TransactionInput
	var err error
	gw.strand("GetWalletUnconfirmedTxnsVerbose", func() {
		var addrs []cipher.Address
		addrs, err = gw.v.Wallets.GetAddresses(wltID)
		if err != nil {
			return
		}

		txns, inputs, err = gw.v.GetUnconfirmedTxnsVerbose(visor.SendsToAddresses(addrs))
	})

	if err != nil {
		return nil, err
	}

	vb := make([]readable.UnconfirmedTransactionVerbose, len(txns))
	for i, txn := range txns {
		v, err := readable.NewUnconfirmedTransactionVerbose(&txn, inputs[i])
		if err != nil {
			return nil, err
		}
		vb[i] = *v
	}

	return vb, nil
}

// ReloadWallets reloads all wallets
func (gw *Gateway) ReloadWallets() error {
	if !gw.Config.EnableWalletAPI {
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
	if !gw.Config.EnableWalletAPI {
		return wallet.ErrWalletAPIDisabled
	}

	var err error

	gw.strand("UnloadWallet", func() {
		err = gw.v.Wallets.Remove(id)
	})

	return err
}

// GetWalletSeed returns seed of wallet of given id,
// returns wallet.ErrWalletNotEncrypted if the wallet is not encrypted.
func (gw *Gateway) GetWalletSeed(id string, password []byte) (string, error) {
	if !gw.Config.EnableWalletAPI {
		return "", wallet.ErrWalletAPIDisabled
	}

	var seed string
	var err error
	gw.strand("GetWalletSeed", func() {
		seed, err = gw.v.Wallets.GetWalletSeed(id, password)
	})
	return seed, err
}

// IsWalletAPIEnabled returns if all wallet related apis are disabled
func (gw *Gateway) IsWalletAPIEnabled() bool {
	return gw.Config.EnableWalletAPI
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
	var count uint64
	var err error

	gw.strand("GetAddressCount", func() {
		count, err = gw.v.AddressCount()
	})

	return count, err
}

// Health is returned by the /health endpoint
type Health struct {
	BlockchainMetadata visor.BlockchainMetadata
	Version            visor.BuildInfo
	OpenConnections    int
	Uptime             time.Duration
}

// GetHealth returns statistics about the running node
func (gw *Gateway) GetHealth() (*Health, error) {
	var health *Health
	var err error
	gw.strand("GetHealth", func() {
		var metadata *visor.BlockchainMetadata
		metadata, err = gw.v.GetBlockchainMetadata()
		if err != nil {
			return
		}

		conns := gw.getConnections()

		health = &Health{
			BlockchainMetadata: *metadata,
			Version:            gw.v.Config.BuildInfo,
			OpenConnections:    len(conns.Connections),
			Uptime:             time.Since(gw.v.StartedAt),
		}
	})

	return health, err
}

// VerifyTxnVerbose verifies an isolated transaction and returns []wallet.UxBalance of
// transaction inputs, whether the transaction is confirmed and error if any
func (gw *Gateway) VerifyTxnVerbose(txn *coin.Transaction) ([]wallet.UxBalance, bool, error) {
	var uxs []wallet.UxBalance
	var isTxnConfirmed bool
	var err error
	gw.strand("VerifyTxnVerbose", func() {
		uxs, isTxnConfirmed, err = gw.v.VerifyTxnVerbose(txn)
	})
	return uxs, isTxnConfirmed, err
}
