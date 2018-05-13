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

	"github.com/NYTimes/gziphandler"

	"github.com/skycoin/skycoin/src/api/webrpc"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/util/file"
	wh "github.com/skycoin/skycoin/src/util/http"
	"github.com/skycoin/skycoin/src/util/logging"
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
	EnableJSON20RPC bool
	EnableGUI       bool
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
}

type muxConfig struct {
	host            string
	appLoc          string
	enableGUI       bool
	enableJSON20RPC bool
}

func create(host string, c Config, gateway Gatewayer) (*Server, error) {
	var appLoc string
	if c.EnableGUI {
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

	var rpc *webrpc.WebRPC
	if c.EnableJSON20RPC {
		logger.Info("JSON 2.0 RPC enabled")
		var err error
		// TODO: change webprc to use http.Gatewayer
		rpc, err = webrpc.New(gateway.(*daemon.Gateway))
		if err != nil {
			return nil, err
		}
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
		enableGUI:       c.EnableGUI,
		enableJSON20RPC: c.EnableJSON20RPC,
	}

	srvMux := newServerMux(mc, gateway, csrfStore, rpc)
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
func Create(host string, c Config, gateway Gatewayer) (*Server, error) {
	s, err := create(host, c, gateway)
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
func CreateHTTPS(host string, c Config, gateway Gatewayer, certFile, keyFile string) (*Server, error) {
	s, err := create(host, c, gateway)
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
	if s == nil {
		return
	}

	logger.Info("Shutting down web interface")
	defer logger.Info("Web interface shut down")
	s.listener.Close()
	<-s.done
}

// newServerMux creates an http.ServeMux with handlers registered
func newServerMux(c muxConfig, gateway Gatewayer, csrfStore *CSRFStore, rpc *webrpc.WebRPC) *http.ServeMux {
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
		handler = gziphandler.GzipHandler(handler)
		mux.Handle(endpoint, handler)
	}

	webHandler("/", newIndexHandler(c.appLoc, c.enableGUI))

	if c.enableGUI {
		fileInfos, _ := ioutil.ReadDir(c.appLoc)
		for _, fileInfo := range fileInfos {
			route := fmt.Sprintf("/%s", fileInfo.Name())
			if fileInfo.IsDir() {
				route = route + "/"
			}
			webHandler(route, http.FileServer(http.Dir(c.appLoc)))
		}
	}

	if c.enableJSON20RPC {
		webHandler("/webrpc", http.HandlerFunc(rpc.Handler))
	}

	// get the current CSRF token
	mux.Handle("/csrf", headerCheck(c.host, getCSRFToken(csrfStore)))

	webHandler("/version", versionHandler(gateway))

	// get set of unspent outputs
	webHandler("/outputs", getOutputsHandler(gateway))

	// get balance of addresses
	webHandler("/balance", getBalanceHandler(gateway))

	// Wallet interface

	// Returns wallet info
	// Method: GET
	// Args:
	//      id - Wallet ID [required]
	webHandler("/wallet", walletGet(gateway))

	// Loads wallet from seed, will scan ahead N address and
	// load addresses till the last one that have coins.
	// Method: POST
	// Args:
	//     seed: wallet seed [required]
	//     label: wallet label [required]
	//     scan: the number of addresses to scan ahead for balances [optional, must be > 0]
	webHandler("/wallet/create", walletCreate(gateway))

	webHandler("/wallet/newAddress", walletNewAddresses(gateway))

	// Returns the confirmed and predicted balance for a specific wallet.
	// The predicted balance is the confirmed balance minus any pending
	// spent amount.
	// GET arguments:
	//      id: Wallet ID
	webHandler("/wallet/balance", walletBalanceHandler(gateway))

	// Sends coins&hours to another address.
	// POST arguments:
	//  id: Wallet ID
	//  coins: Number of coins to spend
	//  dst: Destination address
	//  Returns total amount spent if successful, otherwise error describing
	//  failure status.
	webHandler("/wallet/spend", walletSpendHandler(gateway))

	// Creates a transaction from a wallet
	webHandler("/wallet/transaction", createTransactionHandler(gateway))

	// GET Arguments:
	//      id: Wallet ID
	// Returns all pending transanction for all addresses by selected Wallet
	webHandler("/wallet/transactions", walletTransactionsHandler(gateway))

	// Update wallet label
	// POST Arguments:
	//     id: wallet id
	//     label: wallet label
	webHandler("/wallet/update", walletUpdateHandler(gateway))

	// Returns all loaded wallets
	// returns sensitive information
	webHandler("/wallets", walletsHandler(gateway))

	// Returns wallets directory path
	webHandler("/wallets/folderName", getWalletFolder(gateway))

	// Generate wallet seed
	// GET Arguments:
	//     entropy: entropy bitsize.
	webHandler("/wallet/newSeed", newWalletSeed(gateway))

	// Gets seed of wallet of given id
	// GET Arguments:
	//     id: wallet id
	//     password: wallet password
	webHandler("/wallet/seed", walletSeedHandler(gateway))

	// unload wallet
	// POST Argument:
	//         id: wallet id
	webHandler("/wallet/unload", walletUnloadHandler(gateway))

	// Encrypts wallet
	// POST arguments:
	//     id: wallet id
	//     password: wallet password
	// Returns an encrypted wallet json without sensitive data
	webHandler("/wallet/encrypt", walletEncryptHandler(gateway))

	// Decrypts wallet
	// POST arguments:
	//     id: wallet id
	//     password: wallet password
	webHandler("/wallet/decrypt", walletDecryptHandler(gateway))

	// Blockchain interface

	webHandler("/blockchain/metadata", blockchainHandler(gateway))
	webHandler("/blockchain/progress", blockchainProgressHandler(gateway))

	// get block by hash or seq
	webHandler("/block", getBlock(gateway))
	// get blocks in specific range
	webHandler("/blocks", getBlocks(gateway))
	// get last N blocks
	webHandler("/last_blocks", getLastBlocks(gateway))

	// Network stats interface
	webHandler("/network/connection", connectionHandler(gateway))
	webHandler("/network/connections", connectionsHandler(gateway))
	webHandler("/network/defaultConnections", defaultConnectionsHandler(gateway))
	webHandler("/network/status", networkStatusHandler(gateway))
	webHandler("/network/connections/trust", trustConnectionsHandler(gateway))
	webHandler("/network/connections/exchange", exchgConnectionsHandler(gateway))

	// Transaction handler

	// get set of pending transactions
	webHandler("/pendingTxs", getPendingTxs(gateway))
	// get txn by txid
	webHandler("/transaction", getTransactionByID(gateway))

	// Health check handler
	webHandler("/health", healthCheck(gateway))

	// Returns transactions that match the filters.
	// Method: GET
	// Args:
	//     addrs: Comma seperated addresses [optional, returns all transactions if no address is provided]
	//     confirmed: Whether the transactions should be confirmed [optional, must be 0 or 1; if not provided, returns all]
	webHandler("/transactions", getTransactions(gateway))
	// inject a transaction into network
	webHandler("/injectTransaction", injectTransaction(gateway))
	webHandler("/resendUnconfirmedTxns", resendUnconfirmedTxns(gateway))
	// get raw tx by txid.
	webHandler("/rawtx", getRawTx(gateway))

	// UxOut api handler

	// get uxout by id.
	webHandler("/uxout", getUxOutByID(gateway))
	// get all the address affected uxouts.
	webHandler("/address_uxouts", getAddrUxOuts(gateway))

	// Explorer handler

	// get set of pending transactions
	webHandler("/explorer/address", getTransactionsForAddress(gateway))

	webHandler("/coinSupply", getCoinSupply(gateway))

	webHandler("/richlist", getRichlist(gateway))

	webHandler("/addresscount", getAddressCount(gateway))

	return mux
}

// Returns a http.HandlerFunc for index.html, where index.html is in appLoc
func newIndexHandler(appLoc string, enableGUI bool) http.HandlerFunc {
	// Serves the main page
	return func(w http.ResponseWriter, r *http.Request) {
		if !enableGUI {
			wh.Error404(w, "")
			return
		}

		if r.URL.Path != "/" {
			wh.Error404(w, "")
			return
		}

		if r.URL.Path == "/" {
			page := filepath.Join(appLoc, indexPage)
			logger.Debugf("Serving index page: %s", page)
			http.ServeFile(w, r, page)
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

func versionHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		wh.SendJSONOr500(logger, w, gateway.GetBuildInfo())
	}
}
