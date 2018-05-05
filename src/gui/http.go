package gui

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/util/file"
	wh "github.com/skycoin/skycoin/src/util/http"
	"github.com/skycoin/skycoin/src/util/logging"
	"github.com/skycoin/skycoin/src/wallet"
)

var (
	logger = logging.MustGetLogger("gui")
)

const (
	resourceDir = "dist/"
	devDir      = "dev/"
	indexPage   = "index.html"

	defaultReadTimeout  = time.Second * 10
	defaultWriteTimeout = time.Second * 60
	defaultIdleTimeout  = time.Second * 120
)

// Server exposes an HTTP API
type Server struct {
	server   *http.Server
	listener net.Listener
	done     chan struct{}
}

// Config configures Server
type Config struct {
	StaticDir       string
	DisableCSRF     bool
	EnableWalletAPI bool
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
}

type muxConfig struct {
	host            string
	appLoc          string
	enableWalletAPI bool
}

func create(host string, c Config, daemon *daemon.Daemon) (*Server, error) {
	var appLoc string
	if c.EnableWalletAPI {
		var err error
		appLoc, err = file.DetermineResourcePath(c.StaticDir, resourceDir, devDir)
		if err != nil {
			return nil, err
		}
		logger.Infof("Web resources directory: %s", appLoc)
	}

	csrfStore := &CSRFStore{
		Enabled: !c.DisableCSRF,
	}
	if c.DisableCSRF {
		logger.Warning("CSRF check disabled")
	}

	if c.ReadTimeout == 0 {
		c.ReadTimeout = defaultReadTimeout
	}
	if c.WriteTimeout == 0 {
		c.WriteTimeout = defaultWriteTimeout
	}
	if c.IdleTimeout == 0 {
		c.IdleTimeout = defaultIdleTimeout
	}

	mc := muxConfig{
		host:            host,
		appLoc:          appLoc,
		enableWalletAPI: c.EnableWalletAPI,
	}

	srvMux := newServerMux(mc, daemon.Gateway, csrfStore)
	srv := &http.Server{
		Handler:      srvMux,
		ReadTimeout:  c.ReadTimeout,
		WriteTimeout: c.WriteTimeout,
		IdleTimeout:  c.IdleTimeout,
	}

	return &Server{
		server: srv,
		done:   make(chan struct{}),
	}, nil
}

// Create creates a new Server instance that listens on HTTP
func Create(host string, c Config, daemon *daemon.Daemon) (*Server, error) {
	s, err := create(host, c, daemon)
	if err != nil {
		return nil, err
	}

	logger.Warning("HTTPS not in use!")

	s.listener, err = net.Listen("tcp", host)
	if err != nil {
		return nil, err
	}

	return s, nil
}

// CreateHTTPS creates a new Server instance that listens on HTTPS
func CreateHTTPS(host string, c Config, daemon *daemon.Daemon, certFile, keyFile string) (*Server, error) {
	s, err := create(host, c, daemon)
	if err != nil {
		return nil, err
	}

	logger.Infof("Using %s for the certificate", certFile)
	logger.Infof("Using %s for the key", keyFile)

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	s.listener, err = tls.Listen("tcp", host, &tls.Config{
		Certificates: []tls.Certificate{cert},
	})
	if err != nil {
		return nil, err
	}

	return s, nil
}

// Serve serves the web interface on the configured host
func (s *Server) Serve() error {
	logger.Infof("Starting web interface on %s", s.listener.Addr())
	defer logger.Info("Web interface closed")
	defer close(s.done)

	if err := s.server.Serve(s.listener); err != nil {
		if err != http.ErrServerClosed {
			return err
		}
	}
	return nil
}

// Shutdown closes the HTTP service. This can only be called after Serve or ServeHTTPS has been called.
func (s *Server) Shutdown() {
	logger.Info("Shutting down web interface")
	defer logger.Info("Web interface shut down")
	s.listener.Close()
	<-s.done
}

// newServerMux creates an http.ServeMux with handlers registered
func newServerMux(c muxConfig, gateway Gatewayer, csrfStore *CSRFStore) *http.ServeMux {
	mux := http.NewServeMux()

	headerCheck := func(host string, handler http.Handler) http.Handler {
		handler = OriginRefererCheck(host, handler)
		handler = wh.HostCheck(logger, host, handler)
		return handler
	}

	webHandler := func(endpoint string, handler http.Handler) {

		handler = wh.ElapsedHandler(logger, handler)
		handler = CSRFCheck(csrfStore, handler)
		handler = headerCheck(c.host, handler)
		mux.Handle(endpoint, handler)
	}

	webHandlerAPI := func(endpoint string, handler http.Handler, stable bool) {
		version := "/v2" + endpoint
		if stable {
			version = "/v1" + endpoint
		}
		handler = wh.ElapsedHandler(logger, handler)
		handler = CSRFCheck(csrfStore, handler)
		handler = headerCheck(c.host, handler)
		mux.Handle(version, handler)
	}

	if c.enableWalletAPI {
		webHandler("/", newIndexHandler(c.appLoc))

		fileInfos, _ := ioutil.ReadDir(c.appLoc)
		for _, fileInfo := range fileInfos {
			route := fmt.Sprintf("/%s", fileInfo.Name())
			if fileInfo.IsDir() {
				route = route + "/"
			}
			webHandler(route, http.FileServer(http.Dir(c.appLoc)))
		}
	}

	// get the current CSRF token
	mux.Handle("/v1/csrf", headerCheck(c.host, getCSRFToken(gateway, csrfStore)))

	webHandlerAPI("/version", versionHandler(gateway), true)

	// get set of unspent outputs
	webHandlerAPI("/outputs", getOutputsHandler(gateway), true)

	// get balance of addresses
	webHandlerAPI("/balance", getBalanceHandler(gateway), true)

	// Wallet interface

	// Returns wallet info
	// Method: GET
	// Args:
	//      id - Wallet ID [required]
	webHandlerAPI("/wallet", walletGet(gateway), true)

	// Loads wallet from seed, will scan ahead N address and
	// load addresses till the last one that have coins.
	// Method: POST
	// Args:
	//     seed: wallet seed [required]
	//     label: wallet label [required]
	//     scan: the number of addresses to scan ahead for balances [optional, must be > 0]
	webHandlerAPI("/wallet/create", walletCreate(gateway), true)

	webHandlerAPI("/wallet/newAddress", walletNewAddresses(gateway), true)

	// Returns the confirmed and predicted balance for a specific wallet.
	// The predicted balance is the confirmed balance minus any pending
	// spent amount.
	// GET arguments:
	//      id: Wallet ID
	webHandlerAPI("/wallet/balance", walletBalanceHandler(gateway), true)

	// Sends coins&hours to another address.
	// POST arguments:
	//  id: Wallet ID
	//  coins: Number of coins to spend
	//  dst: Destination address
	//  Returns total amount spent if successful, otherwise error describing
	//  failure status.
	webHandlerAPI("/wallet/spend", walletSpendHandler(gateway), true)

	// Creates a transaction from a wallet
	webHandlerAPI("/wallet/transaction", createTransactionHandler(gateway), true)

	// GET Arguments:
	//      id: Wallet ID
	// Returns all pending transanction for all addresses by selected Wallet
	webHandlerAPI("/wallet/transactions", walletTransactionsHandler(gateway), true)

	// Update wallet label
	// POST Arguments:
	//     id: wallet id
	//     label: wallet label
	webHandlerAPI("/wallet/update", walletUpdateHandler(gateway), true)

	// Returns all loaded wallets
	// returns sensitive information
	webHandlerAPI("/wallets", walletsHandler(gateway), true)

	// Returns wallets directory path
	webHandlerAPI("/wallets/folderName", getWalletFolder(gateway), true)

	// Generate wallet seed
	// GET Arguments:
	//     entropy: entropy bitsize.
	webHandlerAPI("/wallet/newSeed", newWalletSeed(gateway), true)

	// Gets seed of wallet of given id
	// GET Arguments:
	//     id: wallet id
	//     password: wallet password
	webHandlerAPI("/wallet/seed", walletSeedHandler(gateway), true)

	// unload wallet
	// POST Argument:
	//         id: wallet id
	webHandlerAPI("/wallet/unload", walletUnloadHandler(gateway), true)

	// Encrypts wallet
	// POST arguments:
	//     id: wallet id
	//     password: wallet password
	// Returns an encrypted wallet json without sensitive data
	webHandlerAPI("/wallet/encrypt", walletEncryptHandler(gateway), true)

	// Decrypts wallet
	// POST arguments:
	//     id: wallet id
	//     password: wallet password
	webHandlerAPI("/wallet/decrypt", walletDecryptHandler(gateway), true)

	// Blockchain interface

	webHandlerAPI("/blockchain/metadata", blockchainHandler(gateway), true)
	webHandlerAPI("/blockchain/progress", blockchainProgressHandler(gateway), true)

	// get block by hash or seq
	webHandlerAPI("/block", getBlock(gateway), true)
	// get blocks in specific range
	webHandlerAPI("/blocks", getBlocks(gateway), true)
	// get last N blocks
	webHandlerAPI("/last_blocks", getLastBlocks(gateway), true)

	// Network stats interface
	webHandlerAPI("/network/connection", connectionHandler(gateway), true)
	webHandlerAPI("/network/connections", connectionsHandler(gateway), true)
	webHandlerAPI("/network/connections", connectionsHandler(gateway), false)
	webHandlerAPI("/network/defaultConnections", defaultConnectionsHandler(gateway), true)
	webHandlerAPI("/network/connections/trust", trustConnectionsHandler(gateway), true)
	webHandlerAPI("/network/connections/exchange", exchgConnectionsHandler(gateway), true)

	// Transaction handler

	// get set of pending transactions
	webHandlerAPI("/pendingTxs", getPendingTxs(gateway), true)
	// get txn by txid
	webHandlerAPI("/transaction", getTransactionByID(gateway), true)

	// Health check handler
	webHandlerAPI("/health", healthCheck(gateway), true)

	// Returns transactions that match the filters.
	// Method: GET
	// Args:
	//     addrs: Comma seperated addresses [optional, returns all transactions if no address is provided]
	//     confirmed: Whether the transactions should be confirmed [optional, must be 0 or 1; if not provided, returns all]
	webHandlerAPI("/transactions", getTransactions(gateway), true)
	// inject a transaction into network
	webHandlerAPI("/injectTransaction", injectTransaction(gateway), true)
	webHandlerAPI("/resendUnconfirmedTxns", resendUnconfirmedTxns(gateway), true)
	// get raw tx by txid.
	webHandlerAPI("/rawtx", getRawTx(gateway), true)

	// UxOut api handler

	// get uxout by id.
	webHandlerAPI("/uxout", getUxOutByID(gateway), true)
	// get all the address affected uxouts.
	webHandlerAPI("/address_uxouts", getAddrUxOuts(gateway), true)

	// Explorer handler

	// get set of pending transactions
	webHandlerAPI("/explorer/address", getTransactionsForAddress(gateway), true)

	webHandlerAPI("/coinSupply", getCoinSupply(gateway), true)

	webHandlerAPI("/richlist", getRichlist(gateway), true)

	webHandlerAPI("/addresscount", getAddressCount(gateway), true)

	return mux
}

// Returns a http.HandlerFunc for index.html, where index.html is in appLoc
func newIndexHandler(appLoc string) http.HandlerFunc {
	// Serves the main page
	return func(w http.ResponseWriter, r *http.Request) {
		page := filepath.Join(appLoc, indexPage)
		logger.Debugf("Serving index page: %s", page)
		if r.URL.Path == "/" {
			http.ServeFile(w, r, page)
		} else {
			wh.Error404(w, "")
		}
	}
}

func splitCommaString(s string) []string {
	words := strings.FieldsFunc(s, func(r rune) bool {
		return r == ',' || unicode.IsSpace(r)
	})

	// Deduplicate
	var dedupWords []string
	wordsMap := make(map[string]struct{})
	for _, w := range words {
		if _, ok := wordsMap[w]; !ok {
			dedupWords = append(dedupWords, w)
		}
		wordsMap[w] = struct{}{}
	}

	return dedupWords
}

// getOutputsHandler returns UxOuts filtered by a set of addresses or a set of hashes
// URI: /outputs
// Method: GET
// Args:
//    addrs: comma-separated list of addresses
//    hashes: comma-separated list of uxout hashes
// If neither addrs nor hashes are specificed, return all unspent outputs.
// If only one filter is specified, then return outputs match the filter.
// Both filters cannot be specified.
func getOutputsHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		var addrs []string
		var hashes []string

		addrStr := r.FormValue("addrs")
		hashStr := r.FormValue("hashes")

		if addrStr != "" && hashStr != "" {
			wh.Error400(w, "addrs and hashes cannot be specified together")
			return
		}

		filters := []daemon.OutputsFilter{}

		if addrStr != "" {
			addrs = splitCommaString(addrStr)

			for _, a := range addrs {
				if _, err := cipher.DecodeBase58Address(a); err != nil {
					wh.Error400(w, "addrs contains invalid address")
					return
				}
			}

			if len(addrs) > 0 {
				filters = append(filters, daemon.FbyAddresses(addrs))
			}
		}

		if hashStr != "" {
			hashes = splitCommaString(hashStr)
			if len(hashes) > 0 {
				filters = append(filters, daemon.FbyHashes(hashes))
			}
		}

		outs, err := gateway.GetUnspentOutputs(filters...)
		if err != nil {
			err = fmt.Errorf("get unspent outputs failed: %v", err)
			wh.Error500(w, err.Error())
			return
		}

		wh.SendJSONOr500(logger, w, outs)
	}
}

func getBalanceHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		addrsParam := r.FormValue("addrs")
		addrsStr := splitCommaString(addrsParam)

		addrs := make([]cipher.Address, 0, len(addrsStr))
		for _, addr := range addrsStr {
			a, err := cipher.DecodeBase58Address(addr)
			if err != nil {
				wh.Error400(w, fmt.Sprintf("address %s is invalid: %v", addr, err))
				return
			}
			addrs = append(addrs, a)
		}

		bals, err := gateway.GetBalanceOfAddrs(addrs)
		if err != nil {
			err = fmt.Errorf("gateway.GetBalanceOfAddrs failed: %v", err)
			wh.Error500(w, err.Error())
			return
		}

		var balance wallet.BalancePair
		for _, bal := range bals {
			var err error
			balance.Confirmed, err = balance.Confirmed.Add(bal.Confirmed)
			if err != nil {
				wh.Error500(w, err.Error())
				return
			}

			balance.Predicted, err = balance.Predicted.Add(bal.Predicted)
			if err != nil {
				wh.Error500(w, err.Error())
				return
			}
		}

		wh.SendJSONOr500(logger, w, balance)
	}
}

func versionHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		wh.SendJSONOr500(logger, w, gateway.GetBuildInfo())
	}
}
