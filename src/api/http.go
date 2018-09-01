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
	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/util/collections"
	"github.com/skycoin/skycoin/src/util/file"
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
	APIDefault = "READ_ONLY"
	// APIStatus endpoints offer (meta,runtime)data to dashboard and monitoring clients
	APIStatus = "STATUS"
	// APIWallet endpoints implement wallet interface
	APIWallet = "WALLET"
	// APISeed endpoints implement wallet interface
	APISeed = "WALLET_SEED"
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
	DisableCSP           bool
	EnableJSON20RPC      bool
	EnableGUI            bool
	EnableUnversionedAPI bool
	ReadTimeout          time.Duration
	WriteTimeout         time.Duration
	IdleTimeout          time.Duration
	EnabledAPISets       collections.StringSet
}

type muxConfig struct {
	host                 string
	appLoc               string
	enableGUI            bool
	enableJSON20RPC      bool
	enableUnversionedAPI bool
	disableCSP           bool
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
		disableCSP:           c.DisableCSP,
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
	if !c.disableCSP {
		indexHandler = wh.CSPHandler(indexHandler)
	}
	webHandler("/", indexHandler)

	if c.enableGUI {
		fileInfos, _ := ioutil.ReadDir(c.appLoc)

		fs := http.FileServer(http.Dir(c.appLoc))
		if !c.disableCSP {
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
	webHandlerV1("/outputs", forAPISet(getOutputsHandler(gateway), APIDefault))

	// get balance of addresses
	webHandlerV1("/balance", forAPISet(getBalanceHandler(gateway), APIDefault))

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

	webHandlerV1("/blockchain/metadata", forAPISet(blockchainHandler(gateway), APIStatus, APIDefault))
	webHandlerV1("/blockchain/progress", forAPISet(blockchainProgressHandler(gateway), APIStatus, APIDefault))

	// get block by hash or seq
	webHandlerV1("/block", forAPISet(blockHandler(gateway), APIDefault))
	// get blocks in specific range
	webHandlerV1("/blocks", forAPISet(blocksHandler(gateway), APIDefault))
	// get last N blocks
	webHandlerV1("/last_blocks", forAPISet(lastBlocksHandler(gateway), APIDefault))

	// Network stats interface
	webHandlerV1("/network/connection", forAPISet(connectionHandler(gateway), APIStatus, APIDefault))
	webHandlerV1("/network/connections", forAPISet(connectionsHandler(gateway), APIStatus, APIDefault))
	webHandlerV1("/network/defaultConnections", forAPISet(defaultConnectionsHandler(gateway), APIStatus, APIDefault))
	webHandlerV1("/network/connections/trust", forAPISet(trustConnectionsHandler(gateway), APIStatus, APIDefault))
	webHandlerV1("/network/connections/exchange", forAPISet(exchgConnectionsHandler(gateway), APIStatus, APIDefault))

	// Transaction handler

	// get set of pending transactions
	webHandlerV1("/pendingTxs", forAPISet(pendingTxnsHandler(gateway), APIDefault))
	// get txn by txid
	webHandlerV1("/transaction", forAPISet(transactionHandler(gateway), APIDefault))

	// parse and verify transaction
	webHandlerV2("/transaction/verify", verifyTxnHandler(gateway))

	// Health check handler
	webHandlerV1("/health", forAPISet(healthHandler(c, csrfStore, gateway), APIStatus, APIDefault))

	// Returns transactions that match the filters.
	// Method: GET
	// Args:
	//     addrs: Comma separated addresses [optional, returns all transactions if no address is provided]
	//     confirmed: Whether the transactions should be confirmed [optional, must be 0 or 1; if not provided, returns all]
	webHandlerV1("/transactions", forAPISet(getTransactions(gateway), APIDefault))
	// inject a transaction into network
	webHandlerV1("/injectTransaction", forAPISet(injectTransaction(gateway), APIDefault))
	webHandlerV1("/resendUnconfirmedTxns", forAPISet(resendUnconfirmedTxns(gateway), APIDefault))
	// get raw tx by txid.
	webHandlerV1("/rawtx", forAPISet(getRawTxn(gateway), APIDefault))

	// UxOut api handler

	// get uxout by id.
	webHandlerV1("/uxout", forAPISet(getUxOutByID(gateway), APIDefault))
	// get all the address affected uxouts.
	webHandlerV1("/address_uxouts", forAPISet(getAddrUxOuts(gateway), APIDefault))

	webHandlerV2("/address/verify", http.HandlerFunc(forAPISet(addressVerify, APIDefault)))

	// Explorer handler

	// get set of pending transactions
	webHandlerV1("/explorer/address", forAPISet(getTransactionsForAddress(gateway), APIDefault))

	webHandlerV1("/coinSupply", forAPISet(getCoinSupply(gateway), APIStatus, APIDefault))

	webHandlerV1("/richlist", forAPISet(getRichlist(gateway), APIStatus, APIDefault))

	webHandlerV1("/addresscount", forAPISet(getAddressCount(gateway), APIStatus, APIDefault))

	return mux
}

// newIndexHandler returns a http.HandlerFunc for index.html, where index.html is in appLoc
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

// splitCommaString splits a string separated by commas or whitespace into tokens
// and returns an array of unique tokens split from that string
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
