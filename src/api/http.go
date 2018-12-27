/*
Package api implements the REST API interface
*/
package api

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"github.com/NYTimes/gziphandler"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"

	"github.com/skycoin/skycoin/src/api/webrpc"
	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/readable"
	"github.com/skycoin/skycoin/src/util/file"
	wh "github.com/skycoin/skycoin/src/util/http"
	"github.com/skycoin/skycoin/src/util/logging"
	"github.com/skycoin/skycoin/src/util/useragent"
)

var (
	logger = logging.MustGetLogger("api")
)

const (
	resourceDir = "dist/"
	devDir      = "dev/"
	indexPage   = "index.html"

	apiVersion1 = "v1"
	apiVersion2 = "v2"

	defaultReadTimeout  = time.Second * 10
	defaultWriteTimeout = time.Second * 60
	defaultIdleTimeout  = time.Second * 120

	// EndpointsRead endpoints with no side-effects and no changes in node state
	EndpointsRead = "READ"
	// EndpointsStatus endpoints offer (meta,runtime)data to dashboard and monitoring clients
	EndpointsStatus = "STATUS"
	// EndpointsTransaction endpoints export operations on transactions that modify node state
	EndpointsTransaction = "TXN"
	// EndpointsWallet endpoints implement wallet interface
	EndpointsWallet = "WALLET"
	// EndpointsInsecureWalletSeed endpoints implement wallet interface
	EndpointsInsecureWalletSeed = "INSECURE_WALLET_SEED"
	// EndpointsDeprecatedWalletSpend endpoints implement the deprecated /api/v1/wallet/spend method
	EndpointsDeprecatedWalletSpend = "DEPRECATED_WALLET_SPEND"
	// EndpointsPrometheus endpoints for Go application metrics
	EndpointsPrometheus = "PROMETHEUS"
	// EndpointsNetCtrl endpoints for managing network connections
	EndpointsNetCtrl = "NET_CTRL"
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
	Health               HealthConfig
	HostWhitelist        []string
	EnabledAPISets       map[string]struct{}
	Username             string
	Password             string
}

// HealthConfig configuration data exposed in /health
type HealthConfig struct {
	BuildInfo       readable.BuildInfo
	CoinName        string
	DaemonUserAgent useragent.Data
}

type muxConfig struct {
	host                 string
	appLoc               string
	enableGUI            bool
	enableJSON20RPC      bool
	enableUnversionedAPI bool
	disableCSRF          bool
	disableCSP           bool
	enabledAPISets       map[string]struct{}
	hostWhitelist        []string
	username             string
	password             string
	health               HealthConfig
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

func writeHTTPResponse(w http.ResponseWriter, resp HTTPResponse) {
	out, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		wh.Error500(w, "json.MarshalIndent failed")
		return
	}

	w.Header().Add("Content-Type", ContentTypeJSON)

	if resp.Error == nil {
		w.WriteHeader(http.StatusOK)
	} else {
		if resp.Error.Code < 400 || resp.Error.Code >= 600 {
			logger.Critical().Errorf("writeHTTPResponse invalid error status code: %d", resp.Error.Code)
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(resp.Error.Code)
		}
	}

	if _, err := w.Write(out); err != nil {
		logger.WithError(err).Error("http Write failed")
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
		disableCSRF:          c.DisableCSRF,
		disableCSP:           c.DisableCSP,
		health:               c.Health,
		enabledAPISets:       c.EnabledAPISets,
		hostWhitelist:        c.HostWhitelist,
		username:             c.Username,
		password:             c.Password,
	}

	srvMux := newServerMux(mc, gateway, rpc)
	srv := &http.Server{
		Handler:      srvMux,
		ReadTimeout:  c.ReadTimeout,
		WriteTimeout: c.WriteTimeout,
		IdleTimeout:  c.IdleTimeout,
		// MaxHeaderBytes: http.DefaultMaxHeaderBytes, // adjust this to allow longer GET queries
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
		if closeErr := s.listener.Close(); closeErr != nil {
			logger.WithError(err).Warning("s.listener.Close() error")
		}
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
		if closeErr := s.listener.Close(); closeErr != nil {
			logger.WithError(err).Warning("s.listener.Close() error")
		}
		return nil, err
	}

	s.listener = listener

	return s, nil
}

// Addr returns the listening address of the Server
func (s *Server) Addr() string {
	if s == nil || s.listener == nil {
		return ""
	}
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
	if err := s.listener.Close(); err != nil {
		logger.WithError(err).Warning("s.listener.Close() error")
	}
	<-s.done
}

// newServerMux creates an http.ServeMux with handlers registered
func newServerMux(c muxConfig, gateway Gatewayer, rpc *webrpc.WebRPC) *http.ServeMux {
	mux := http.NewServeMux()

	allowedOrigins := []string{fmt.Sprintf("http://%s", c.host)}
	for _, s := range c.hostWhitelist {
		allowedOrigins = append(allowedOrigins, fmt.Sprintf("http://%s", s))
	}

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:     allowedOrigins,
		Debug:              false,
		AllowedMethods:     []string{http.MethodGet, http.MethodPost},
		AllowedHeaders:     []string{"Origin", "Accept", "Content-Type", "X-Requested-With", CSRFHeaderName},
		AllowCredentials:   false, // credentials are not used, but it would be safe to enable if necessary
		OptionsPassthrough: false,
	})

	headerCheck := func(apiVersion, host string, hostWhitelist []string, handler http.Handler) http.Handler {
		handler = originRefererCheck(apiVersion, host, hostWhitelist, handler)
		handler = hostCheck(apiVersion, host, hostWhitelist, handler)
		return handler
	}

	forAPISet := func(f http.HandlerFunc, apiNames []string) http.HandlerFunc {
		if len(apiNames) == 0 {
			logger.Panic("apiNames should not be empty")
		}

		isEnabled := false

		for _, k := range apiNames {
			if _, ok := c.enabledAPISets[k]; ok {
				isEnabled = true
				break
			}
		}

		return func(w http.ResponseWriter, r *http.Request) {
			if isEnabled {
				f(w, r)
			} else {
				wh.Error403(w, "Endpoint is disabled")
			}
		}
	}

	webHandlerCSRFOptional := func(apiVersion, endpoint string, handler http.Handler, checkCSRF bool) {
		handler = wh.ElapsedHandler(logger, handler)
		handler = corsHandler.Handler(handler)
		if checkCSRF {
			handler = CSRFCheck(apiVersion, c.disableCSRF, handler)
		}
		handler = headerCheck(apiVersion, c.host, c.hostWhitelist, handler)
		handler = basicAuth(apiVersion, c.username, c.password, "skycoin daemon", handler)
		handler = gziphandler.GzipHandler(handler)
		mux.Handle(endpoint, handler)
	}

	webHandler := func(apiVersion, endpoint string, handler http.Handler) {
		webHandlerCSRFOptional(apiVersion, endpoint, handler, true)
	}

	webHandlerV1 := func(endpoint string, handler http.Handler) {
		if c.enableUnversionedAPI {
			webHandler(apiVersion1, endpoint, handler)
		}
		webHandler(apiVersion1, "/api/v1"+endpoint, handler)
	}

	webHandlerV2 := func(endpoint string, handler http.Handler) {
		webHandler(apiVersion2, "/api/v2"+endpoint, handler)
	}

	indexHandler := newIndexHandler(c.appLoc, c.enableGUI)
	if !c.disableCSP {
		indexHandler = CSPHandler(indexHandler)
	}
	webHandler(apiVersion1, "/", indexHandler)

	if c.enableGUI {
		fileInfos, err := ioutil.ReadDir(c.appLoc)
		if err != nil {
			logger.WithError(err).Panicf("ioutil.ReadDir(%s) failed", c.appLoc)
		}

		fs := http.FileServer(http.Dir(c.appLoc))
		if !c.disableCSP {
			fs = CSPHandler(fs)
		}

		for _, fileInfo := range fileInfos {
			route := fmt.Sprintf("/%s", fileInfo.Name())
			if fileInfo.IsDir() {
				route = route + "/"
			}

			webHandler(apiVersion1, route, fs)
		}
	}

	if c.enableJSON20RPC {
		webHandlerV1("/webrpc", http.HandlerFunc(rpc.Handler))
	}

	// get the current CSRF token
	csrfHandlerV1 := func(endpoint string, handler http.Handler) {
		if c.enableUnversionedAPI {
			webHandlerCSRFOptional(apiVersion1, endpoint, handler, false)
		}
		webHandlerCSRFOptional(apiVersion1, "/api/v1"+endpoint, handler, false)
	}
	csrfHandlerV1("/csrf", getCSRFToken(c.disableCSRF)) // csrf is always available, regardless of the API set

	// Status endpoints
	webHandlerV1("/version", versionHandler(c.health.BuildInfo)) // version is always available, regardless of the API set
	webHandlerV1("/health", forAPISet(healthHandler(c, gateway), []string{EndpointsRead, EndpointsStatus}))

	// Wallet endpoints
	webHandlerV1("/wallet", forAPISet(walletHandler(gateway), []string{EndpointsWallet}))
	webHandlerV1("/wallet/create", forAPISet(walletCreateHandler(gateway), []string{EndpointsWallet}))
	webHandlerV1("/wallet/newAddress", forAPISet(walletNewAddressesHandler(gateway), []string{EndpointsWallet}))
	webHandlerV1("/wallet/balance", forAPISet(walletBalanceHandler(gateway), []string{EndpointsWallet}))
	webHandlerV1("/wallet/spend", forAPISet(walletSpendHandler(gateway), []string{EndpointsDeprecatedWalletSpend}))
	webHandlerV1("/wallet/transaction", forAPISet(createTransactionHandler(gateway), []string{EndpointsWallet}))
	webHandlerV1("/wallet/transactions", forAPISet(walletTransactionsHandler(gateway), []string{EndpointsWallet}))
	webHandlerV1("/wallet/update", forAPISet(walletUpdateHandler(gateway), []string{EndpointsWallet}))
	webHandlerV1("/wallets", forAPISet(walletsHandler(gateway), []string{EndpointsWallet}))
	webHandlerV1("/wallets/folderName", forAPISet(walletFolderHandler(gateway), []string{EndpointsWallet}))
	webHandlerV1("/wallet/newSeed", forAPISet(newSeedHandler(), []string{EndpointsWallet}))
	webHandlerV1("/wallet/seed", forAPISet(walletSeedHandler(gateway), []string{EndpointsInsecureWalletSeed}))

	webHandlerV1("/wallet/unload", forAPISet(walletUnloadHandler(gateway), []string{EndpointsWallet}))
	webHandlerV1("/wallet/encrypt", forAPISet(walletEncryptHandler(gateway), []string{EndpointsWallet}))
	webHandlerV1("/wallet/decrypt", forAPISet(walletDecryptHandler(gateway), []string{EndpointsWallet}))
	webHandlerV2("/wallet/recover", forAPISet(walletRecoverHandler(gateway), []string{EndpointsWallet}))

	// Blockchain interface
	webHandlerV1("/blockchain/metadata", forAPISet(blockchainMetadataHandler(gateway), []string{EndpointsRead, EndpointsStatus}))
	webHandlerV1("/blockchain/progress", forAPISet(blockchainProgressHandler(gateway), []string{EndpointsRead, EndpointsStatus}))
	webHandlerV1("/block", forAPISet(blockHandler(gateway), []string{EndpointsRead}))
	webHandlerV1("/blocks", forAPISet(blocksHandler(gateway), []string{EndpointsRead}))
	webHandlerV1("/last_blocks", forAPISet(lastBlocksHandler(gateway), []string{EndpointsRead}))

	// Network stats endpoints
	webHandlerV1("/network/connection", forAPISet(connectionHandler(gateway), []string{EndpointsRead, EndpointsStatus}))
	webHandlerV1("/network/connections", forAPISet(connectionsHandler(gateway), []string{EndpointsRead, EndpointsStatus}))
	webHandlerV1("/network/defaultConnections", forAPISet(defaultConnectionsHandler(gateway), []string{EndpointsRead, EndpointsStatus}))
	webHandlerV1("/network/connections/trust", forAPISet(trustConnectionsHandler(gateway), []string{EndpointsRead, EndpointsStatus}))
	webHandlerV1("/network/connections/exchange", forAPISet(exchgConnectionsHandler(gateway), []string{EndpointsRead, EndpointsStatus}))

	// Network admin endpoints
	webHandlerV1("/network/connection/disconnect", forAPISet(disconnectHandler(gateway), []string{EndpointsNetCtrl}))

	// Transaction related endpoints
	webHandlerV1("/pendingTxs", forAPISet(pendingTxnsHandler(gateway), []string{EndpointsRead}))
	webHandlerV1("/transaction", forAPISet(transactionHandler(gateway), []string{EndpointsRead}))
	webHandlerV2("/transaction/verify", forAPISet(verifyTxnHandler(gateway), []string{EndpointsRead}))
	webHandlerV1("/transactions", forAPISet(transactionsHandler(gateway), []string{EndpointsRead}))
	webHandlerV1("/injectTransaction", forAPISet(injectTransactionHandler(gateway), []string{EndpointsTransaction, EndpointsWallet}))
	webHandlerV1("/resendUnconfirmedTxns", forAPISet(resendUnconfirmedTxnsHandler(gateway), []string{EndpointsTransaction}))
	webHandlerV1("/rawtx", forAPISet(rawTxnHandler(gateway), []string{EndpointsRead}))

	// Unspent output related endpoints
	webHandlerV1("/outputs", forAPISet(outputsHandler(gateway), []string{EndpointsRead}))
	webHandlerV1("/balance", forAPISet(balanceHandler(gateway), []string{EndpointsRead}))
	webHandlerV1("/uxout", forAPISet(uxOutHandler(gateway), []string{EndpointsRead}))
	webHandlerV1("/address_uxouts", forAPISet(addrUxOutsHandler(gateway), []string{EndpointsRead}))

	// golang process internal metrics for Prometheus
	webHandlerV2("/metrics", forAPISet(promhttp.Handler().(http.HandlerFunc), []string{EndpointsPrometheus}))

	// Address related endpoints
	webHandlerV2("/address/verify", forAPISet(addressVerifyHandler, []string{EndpointsRead}))

	// Explorer endpoints
	webHandlerV1("/explorer/address", forAPISet(transactionsForAddressHandler(gateway), []string{EndpointsRead}))
	webHandlerV1("/coinSupply", forAPISet(coinSupplyHandler(gateway), []string{EndpointsRead}))
	webHandlerV1("/richlist", forAPISet(richlistHandler(gateway), []string{EndpointsRead}))
	webHandlerV1("/addresscount", forAPISet(addressCountHandler(gateway), []string{EndpointsRead}))

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
