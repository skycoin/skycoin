package api

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"time"
	"unicode"

	"github.com/NYTimes/gziphandler"

	"github.com/skycoin/skycoin/src/api/webrpc"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/util/file"
	"github.com/skycoin/skycoin/src/util/flagutils"
	wh "github.com/skycoin/skycoin/src/util/http"
	"github.com/skycoin/skycoin/src/util/logging"
)

var (
	logger = logging.MustGetLogger("api")
)

const (
	resourceDir = "dist/"
	devDir      = "dev/"
	indexPage   = "index.html"

	defaultReadTimeout  = time.Second * 10
	defaultWriteTimeout = time.Second * 60
	defaultIdleTimeout  = time.Second * 120

	// APIDefault endpoints available when nodes executed with no CLI args
	APIDefault = "DEFAULT"
	// APIBlockchain endpoints expose blockchain to clients
	APIBlockchain = "BLOCKCHAIN"
	// APIStatus endpoints offer (meta,runtime)data to dashboard and monitoring clients
	APIStatus = "STATUS"
	// APIWallet endpoints implement wallet interface
	APIWallet = "WALLET"
	// APISeed endpoints implement wallet interface
	APISeed = "SEED"
	// APIPex endpoints expose peer exchange data to clients
	APIPex = "PEX"
	// APITxn endpoints for transaction data and related operations
	APITxn = "TX"
	// APIUxOut endpoints expose UxOut data to clients
	APIUxOut = "UX"
	// APIExplorer endpoints consumed by or related to Skycoin blockchain explorer
	APIExplorer = "EXPLORER"
)

// Server exposes an HTTP API
type Server struct {
	server   *http.Server
	listener net.Listener
	done     chan struct{}
}

// Config configures Server
type Config struct {
	StaticDir            string
	DisableCSRF          bool
	EnableWalletAPI      bool
	EnableJSON20RPC      bool
	EnableGUI            bool
	EnableUnversionedAPI bool
	ReadTimeout          time.Duration
	WriteTimeout         time.Duration
	IdleTimeout          time.Duration
	EnabledAPISets       flagutils.StringSet
}

type muxConfig struct {
	host                 string
	appLoc               string
	enableGUI            bool
	enableJSON20RPC      bool
	enableUnversionedAPI bool
}

// HTTPResponse represents the http response struct
type HTTPResponse struct {
	Error *HTTPError  `json:"error,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}

// HTTPError is included in an HTTPResponse
type HTTPError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// NewHTTPErrorResponse returns an HTTPResponse with the Error field populated
func NewHTTPErrorResponse(code int, msg string) HTTPResponse {
	if msg == "" {
		msg = http.StatusText(code)
	}

	return HTTPResponse{
		Error: &HTTPError{
			Code:    code,
			Message: msg,
		},
	}
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
		host:                 host,
		appLoc:               appLoc,
		enableGUI:            c.EnableGUI,
		enableJSON20RPC:      c.EnableJSON20RPC,
		enableUnversionedAPI: c.EnableUnversionedAPI,
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
	logger.Warning("HTTPS not in use!")

	listener, err := net.Listen("tcp", host)
	if err != nil {
		return nil, err
	}

	// If the host did not specify a port, allowing the kernel to assign one,
	// we need to get the assigned address to know the full hostname
	host = listener.Addr().String()

	s, err := create(host, c, gateway)
	if err != nil {
		s.listener.Close()
		return nil, err
	}

	s.listener = listener

	return s, nil
}

// CreateHTTPS creates a new Server instance that listens on HTTPS
func CreateHTTPS(host string, c Config, gateway Gatewayer, certFile, keyFile string) (*Server, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	logger.Infof("Using %s for the certificate", certFile)
	logger.Infof("Using %s for the key", keyFile)

	listener, err := tls.Listen("tcp", host, &tls.Config{
		Certificates: []tls.Certificate{cert},
	})
	if err != nil {
		return nil, err
	}

	// If the host did not specify a port, allowing the kernel to assign one,
	// we need to get the assigned address to know the full hostname
	host = listener.Addr().String()

	s, err := create(host, c, gateway)
	if err != nil {
		s.listener.Close()
		return nil, err
	}

	s.listener = listener

	return s, nil
}

// Addr returns the listening address of the Server
func (s *Server) Addr() string {
	return s.listener.Addr().String()
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

	forAPISet := func(hf http.HandlerFunc, mainAPIName string, otherAPINames ...string) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if !gateway.IsAPISetEnabled(mainAPIName, otherAPINames...) {
				funcName := runtime.FuncForPC(reflect.ValueOf(hf).Pointer()).Name()
				// FIXME: Debugf ?
				logger.Infof("Handler %s not executed because API set %v not enabled", funcName, otherAPINames)
				wh.Error403(w, "")
			} else {
				hf(w, r)
			}
		}
	}

	webHandler := func(endpoint string, handler http.Handler) {
		handler = wh.ElapsedHandler(logger, handler)
		handler = CSRFCheck(csrfStore, handler)
		handler = headerCheck(c.host, handler)
		handler = gziphandler.GzipHandler(handler)
		mux.Handle(endpoint, handler)
	}

	webHandlerV1 := func(endpoint string, handler http.Handler) {
		if c.enableUnversionedAPI {
			webHandler(endpoint, handler)
		}
		webHandler("/api/v1"+endpoint, handler)
	}

	webHandlerV2 := func(endpoint string, handler http.Handler) {
		webHandler("/api/v2"+endpoint, handler)
	}

	indexHandler := newIndexHandler(c.appLoc, c.enableGUI)
	if gateway.IsCSPEnabled() {
		indexHandler = wh.CSPHandler(indexHandler)
	}
	webHandler("/", indexHandler)

	if c.enableGUI {
		fileInfos, _ := ioutil.ReadDir(c.appLoc)

		fs := http.FileServer(http.Dir(c.appLoc))
		if gateway.IsCSPEnabled() {
			fs = wh.CSPHandler(fs)
		}

		for _, fileInfo := range fileInfos {
			route := fmt.Sprintf("/%s", fileInfo.Name())
			if fileInfo.IsDir() {
				route = route + "/"
			}

			webHandler(route, fs)
		}
	}

	if c.enableJSON20RPC {
		webHandlerV1("/webrpc", http.HandlerFunc(rpc.Handler))
	}

	// get the current CSRF token
	csrfHandler := headerCheck(c.host, getCSRFToken(csrfStore))
	mux.Handle("/csrf", csrfHandler)
	mux.Handle("/api/v1/csrf", csrfHandler)

	webHandlerV1("/version", versionHandler(gateway))

	// get set of unspent outputs
	webHandlerV1("/outputs", forAPISet(getOutputsHandler(gateway), APIUxOut, APIBlockchain, APIDefault))

	// get balance of addresses
	webHandlerV1("/balance", forAPISet(getBalanceHandler(gateway), APIBlockchain, APIDefault))

	// Wallet interface

	// Returns wallet info
	// Method: GET
	// Args:
	//      id - Wallet ID [required]
	webHandlerV1("/wallet", forAPISet(walletGet(gateway), APIWallet))

	// Loads wallet from seed, will scan ahead N address and
	// load addresses till the last one that have coins.
	// Method: POST
	// Args:
	//     seed: wallet seed [required]
	//     label: wallet label [required]
	//     scan: the number of addresses to scan ahead for balances [optional, must be > 0]
	webHandlerV1("/wallet/create", forAPISet(walletCreate(gateway), APIWallet))

	webHandlerV1("/wallet/newAddress", forAPISet(walletNewAddresses(gateway), APIWallet))

	// Returns the confirmed and predicted balance for a specific wallet.
	// The predicted balance is the confirmed balance minus any pending
	// spent amount.
	// GET arguments:
	//      id: Wallet ID
	webHandlerV1("/wallet/balance", forAPISet(walletBalanceHandler(gateway), APIWallet))

	// Sends coins&hours to another address.
	// POST arguments:
	//  id: Wallet ID
	//  coins: Number of coins to spend
	//  dst: Destination address
	//  Returns total amount spent if successful, otherwise error describing
	//  failure status.
	webHandlerV1("/wallet/spend", forAPISet(walletSpendHandler(gateway), APIWallet))

	// Creates a transaction from a wallet
	webHandlerV1("/wallet/transaction", forAPISet(createTransactionHandler(gateway), APIWallet))

	// GET Arguments:
	//      id: Wallet ID
	// Returns all pending transanction for all addresses by selected Wallet
	webHandlerV1("/wallet/transactions", forAPISet(walletTransactionsHandler(gateway), APIWallet))

	// Update wallet label
	// POST Arguments:
	//     id: wallet id
	//     label: wallet label
	webHandlerV1("/wallet/update", forAPISet(walletUpdateHandler(gateway), APIWallet))

	// Returns all loaded wallets
	// returns sensitive information
	webHandlerV1("/wallets", forAPISet(walletsHandler(gateway), APIWallet))

	// Returns wallets directory path
	webHandlerV1("/wallets/folderName", forAPISet(getWalletFolder(gateway), APIWallet))

	// Generate wallet seed
	// GET Arguments:
	//     entropy: entropy bitsize.
	webHandlerV1("/wallet/newSeed", forAPISet(newWalletSeed(gateway), APIWallet))

	// Gets seed of wallet of given id
	// GET Arguments:
	//     id: wallet id
	//     password: wallet password
	webHandlerV1("/wallet/seed", forAPISet(walletSeedHandler(gateway), APISeed)) // FIXME: APIWallet?

	// unload wallet
	// POST Argument:
	//         id: wallet id
	webHandlerV1("/wallet/unload", forAPISet(walletUnloadHandler(gateway), APIWallet))

	// Encrypts wallet
	// POST arguments:
	//     id: wallet id
	//     password: wallet password
	// Returns an encrypted wallet json without sensitive data
	webHandlerV1("/wallet/encrypt", forAPISet(walletEncryptHandler(gateway), APIWallet))

	// Decrypts wallet
	// POST arguments:
	//     id: wallet id
	//     password: wallet password
	webHandlerV1("/wallet/decrypt", forAPISet(walletDecryptHandler(gateway), APIWallet))

	// Blockchain interface

	webHandlerV1("/blockchain/metadata", forAPISet(blockchainHandler(gateway), APIBlockchain, APIStatus, APIDefault))
	webHandlerV1("/blockchain/progress", forAPISet(blockchainProgressHandler(gateway), APIBlockchain, APIStatus, APIDefault))

	// get block by hash or seq
	webHandlerV1("/block", forAPISet(getBlock(gateway), APIBlockchain, APIDefault))
	// get blocks in specific range
	webHandlerV1("/blocks", forAPISet(getBlocks(gateway), APIBlockchain, APIDefault))
	// get last N blocks
	webHandlerV1("/last_blocks", forAPISet(getLastBlocks(gateway), APIBlockchain, APIDefault))

	// Network stats interface
	webHandlerV1("/network/connection", forAPISet(connectionHandler(gateway), APIPex, ApiStatus, APIDefault))
	webHandlerV1("/network/connections", forAPISet(connectionsHandler(gateway), APIPex, ApiStatus, APIDefault))
	webHandlerV1("/network/defaultConnections", forAPISet(defaultConnectionsHandler(gateway), APIPex, ApiStatus, APIDefault))
	webHandlerV1("/network/connections/trust", forAPISet(trustConnectionsHandler(gateway), APIPex, ApiStatus, APIDefault))
	webHandlerV1("/network/connections/exchange", forAPISet(exchgConnectionsHandler(gateway), APIPex, ApiStatus, APIDefault))

	// Transaction handler

	// get set of pending transactions
	webHandlerV1("/pendingTxs", forAPISet(getPendingTxns(gateway), APITxn, APIDefault))
	// get txn by txid
	webHandlerV1("/transaction", forAPISet(getTransactionByID(gateway), APITxn, APIDefault))

	// parse and verify transaction
	webHandlerV2("/transaction/verify", verifyTxnHandler(gateway))

	// Health check handler
	webHandlerV1("/health", forAPISet(healthCheck(gateway), APIStatus, APIBlockchain, APIPex, APITxn, APIDefault))

	// Returns transactions that match the filters.
	// Method: GET
	// Args:
	//     addrs: Comma seperated addresses [optional, returns all transactions if no address is provided]
	//     confirmed: Whether the transactions should be confirmed [optional, must be 0 or 1; if not provided, returns all]
	webHandlerV1("/transactions", forAPISet(getTransactions(gateway), ApiBlockchain, APITxn, APIDefault))
	// inject a transaction into network
	webHandlerV1("/injectTransaction", forAPISet(injectTransaction(gateway), APITxn, APIDefault))
	webHandlerV1("/resendUnconfirmedTxns", forAPISet(resendUnconfirmedTxns(gateway), APITxn, APIDefault))
	// get raw tx by txid.
	webHandlerV1("/rawtx", forAPISet(getRawTxn(gateway), APITxn, APIDefault))

	// UxOut api handler

	// get uxout by id.
	webHandlerV1("/uxout", forAPISet(getUxOutByID(gateway), APIUxOut, APIBlockchain, APIDefault))
	// get all the address affected uxouts.
	webHandlerV1("/address_uxouts", forAPISet(getAddrUxOuts(gateway), APIUxOut, APIBlockchain, APIDefault))

	webHandlerV2("/address/verify", http.HandlerFunc(forAPISet(addressVerify, APIBlockchain, APITxn, APIExplorer, APIUxOut, APIDefault)))

	// Explorer handler

	// get set of pending transactions
	webHandlerV1("/explorer/address", forAPISet(getTransactionsForAddress(gateway), APIExplorer, APIBlockchain, APIDefault))

	webHandlerV1("/coinSupply", forAPISet(getCoinSupply(gateway), APIBlockchain, APIStatus, APIExplorer, APIDefault))

	webHandlerV1("/richlist", forAPISet(getRichlist(gateway), APIBlockchain, APIStatus, APIExplorer, APIDefault))

	webHandlerV1("/addresscount", forAPISet(getAddressCount(gateway), APIBlockchain, APIStatus, APIExplorer, APIDefault))

	return mux
}

// Returns a http.HandlerFunc for index.html, where index.html is in appLoc
func newIndexHandler(appLoc string, enableGUI bool) http.Handler {
	// Serves the main page
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	})
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
// URI: /api/v1/outputs
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
