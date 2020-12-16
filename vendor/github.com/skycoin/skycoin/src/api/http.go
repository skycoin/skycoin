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
	"github.com/rs/cors"

	"github.com/skycoin/skycoin/src/cipher"
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
	// EndpointsPrometheus endpoints for Go application metrics
	EndpointsPrometheus = "PROMETHEUS"
	// EndpointsNetCtrl endpoints for managing network connections
	EndpointsNetCtrl = "NET_CTRL"
	// EndpointsStorage endpoints implement interface for key-value storage for arbitrary data
	EndpointsStorage = "STORAGE"
)

// Server exposes an HTTP API
type Server struct {
	server   *http.Server
	listener net.Listener
	done     chan struct{}
}

// Config configures Server
type Config struct {
	StaticDir          string
	DisableCSRF        bool
	DisableHeaderCheck bool
	DisableCSP         bool
	EnableGUI          bool
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
	IdleTimeout        time.Duration
	Health             HealthConfig
	HostWhitelist      []string
	EnabledAPISets     map[string]struct{}
	Username           string
	Password           string
}

// HealthConfig configuration data exposed in /health
type HealthConfig struct {
	BuildInfo       readable.BuildInfo
	Fiber           readable.FiberConfig
	DaemonUserAgent useragent.Data
	BlockPublisher  bool
}

type muxConfig struct {
	host               string
	appLoc             string
	enableGUI          bool
	disableCSRF        bool
	disableHeaderCheck bool
	disableCSP         bool
	enabledAPISets     map[string]struct{}
	hostWhitelist      []string
	username           string
	password           string
	health             HealthConfig
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

	if c.DisableHeaderCheck {
		logger.Warning("Header check disabled")
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
		host:               host,
		appLoc:             appLoc,
		enableGUI:          c.EnableGUI,
		disableCSRF:        c.DisableCSRF,
		disableHeaderCheck: c.DisableHeaderCheck,
		disableCSP:         c.DisableCSP,
		health:             c.Health,
		enabledAPISets:     c.EnabledAPISets,
		hostWhitelist:      c.HostWhitelist,
		username:           c.Username,
		password:           c.Password,
	}

	srvMux := newServerMux(mc, gateway)
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
func newServerMux(c muxConfig, gateway Gatewayer) *http.ServeMux {
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

	forMethodAPISets := func(apiVersion string, f http.Handler, methodsAPISets map[string][]string) http.Handler {
		if len(methodsAPISets) == 0 {
			logger.Panic("methodsAPISets should not be empty")
		}

		switch apiVersion {
		case apiVersion1, apiVersion2:
		default:
			logger.Panicf("Invalid API version %q", apiVersion)
		}

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiSets := methodsAPISets[r.Method]

			// If no API sets are specified for a given method, return 405 Method Not Allowed
			if len(apiSets) == 0 {
				switch apiVersion {
				case apiVersion1:
					wh.Error405(w)
				case apiVersion2:
					resp := NewHTTPErrorResponse(http.StatusMethodNotAllowed, "")
					writeHTTPResponse(w, resp)
				}
				return
			}

			for _, k := range apiSets {
				if _, ok := c.enabledAPISets[k]; ok {
					f.ServeHTTP(w, r)
					return
				}
			}

			switch apiVersion {
			case apiVersion1:
				wh.Error403(w, "Endpoint is disabled")
			case apiVersion2:
				resp := NewHTTPErrorResponse(http.StatusForbidden, "Endpoint is disabled")
				writeHTTPResponse(w, resp)
			}
		})
	}

	webHandlerWithOptionals := func(apiVersion, endpoint string, handlerFunc http.Handler, checkCSRF, checkHeaders bool) {
		handler := wh.ElapsedHandler(logger, handlerFunc)

		handler = corsHandler.Handler(handler)

		if checkCSRF {
			handler = CSRFCheck(apiVersion, c.disableCSRF, handler)
		}

		if checkHeaders {
			handler = headerCheck(apiVersion, c.host, c.hostWhitelist, handler)
		}

		if apiVersion == apiVersion2 {
			handler = ContentTypeJSONRequired(handler)
		}

		handler = basicAuth(apiVersion, c.username, c.password, "skycoin daemon", handler)
		handler = gziphandler.GzipHandler(handler)
		mux.Handle(endpoint, handler)
	}

	webHandler := func(apiVersion, endpoint string, handler http.Handler, methodAPISets map[string][]string) {
		// methodAPISets can be nil to ignore the concept of API sets for an endpoint. It will always be enabled.
		// Explicitly check nil, caller should not pass empty initialized map
		if methodAPISets != nil {
			handler = forMethodAPISets(apiVersion, handler, methodAPISets)
		}

		webHandlerWithOptionals(apiVersion, endpoint, handler, true, !c.disableHeaderCheck)
	}

	webHandlerV1 := func(endpoint string, handler http.Handler, methodAPISets map[string][]string) {
		webHandler(apiVersion1, "/api/v1"+endpoint, handler, methodAPISets)
	}

	webHandlerV2 := func(endpoint string, handler http.Handler, methodAPISets map[string][]string) {
		webHandler(apiVersion2, "/api/v2"+endpoint, handler, methodAPISets)
	}

	indexHandler := newIndexHandler(c.appLoc, c.enableGUI)
	if !c.disableCSP {
		indexHandler = CSPHandler(indexHandler, ContentSecurityPolicy)
	}
	webHandler(apiVersion1, "/", indexHandler, nil)

	if c.enableGUI {
		fileInfos, err := ioutil.ReadDir(c.appLoc)
		if err != nil {
			logger.WithError(err).Panicf("ioutil.ReadDir(%s) failed", c.appLoc)
		}

		fs := http.FileServer(http.Dir(c.appLoc))
		if !c.disableCSP {
			fs = CSPHandler(fs, ContentSecurityPolicy)
		}

		for _, fileInfo := range fileInfos {
			route := fmt.Sprintf("/%s", fileInfo.Name())
			if fileInfo.IsDir() {
				route = route + "/"
			}

			webHandler(apiVersion1, route, fs, nil)
		}
	}

	// get the current CSRF token
	csrfHandlerV1 := func(endpoint string, handler http.Handler) {
		webHandlerWithOptionals(apiVersion1, "/api/v1"+endpoint, handler, false, !c.disableHeaderCheck)
	}
	csrfHandlerV1("/csrf", getCSRFToken(c.disableCSRF)) // csrf is always available, regardless of the API set

	// Status endpoints
	webHandlerV1("/version", versionHandler(c.health.BuildInfo), nil) // version is always available, regardless of the API set
	webHandlerV1("/health", healthHandler(c, gateway), map[string][]string{
		http.MethodGet: []string{EndpointsRead, EndpointsStatus},
	})

	// Wallet endpoints
	webHandlerV1("/wallet", walletHandler(gateway), map[string][]string{
		http.MethodGet: []string{EndpointsWallet},
	})
	webHandlerV1("/wallet/create", walletCreateHandler(gateway), map[string][]string{
		http.MethodPost: []string{EndpointsWallet},
	})
	webHandlerV1("/wallet/newAddress", walletNewAddressesHandler(gateway), map[string][]string{
		http.MethodPost: []string{EndpointsWallet},
	})
	webHandlerV1("/wallet/balance", walletBalanceHandler(gateway), map[string][]string{
		http.MethodGet: []string{EndpointsWallet},
	})
	webHandlerV1("/wallet/transaction", walletCreateTransactionHandler(gateway), map[string][]string{
		http.MethodPost: []string{EndpointsWallet},
	})
	webHandlerV2("/wallet/transaction/sign", walletSignTransactionHandler(gateway), map[string][]string{
		http.MethodPost: []string{EndpointsWallet},
	})
	webHandlerV1("/wallet/transactions", walletTransactionsHandler(gateway), map[string][]string{
		http.MethodGet: []string{EndpointsWallet},
	})
	webHandlerV1("/wallet/update", walletUpdateHandler(gateway), map[string][]string{
		http.MethodPost: []string{EndpointsWallet},
	})
	webHandlerV1("/wallets", walletsHandler(gateway), map[string][]string{
		http.MethodGet: []string{EndpointsWallet},
	})
	webHandlerV1("/wallets/folderName", walletFolderHandler(gateway), map[string][]string{
		http.MethodGet: []string{EndpointsWallet},
	})
	webHandlerV1("/wallet/newSeed", newSeedHandler(), map[string][]string{
		http.MethodGet: []string{EndpointsWallet},
	})
	webHandlerV1("/wallet/seed", walletSeedHandler(gateway), map[string][]string{
		http.MethodPost: []string{EndpointsInsecureWalletSeed},
	})
	webHandlerV2("/wallet/seed/verify", http.HandlerFunc(walletVerifySeedHandler), map[string][]string{
		http.MethodPost: []string{EndpointsWallet},
	})

	webHandlerV1("/wallet/unload", walletUnloadHandler(gateway), map[string][]string{
		http.MethodPost: []string{EndpointsWallet},
	})
	webHandlerV1("/wallet/encrypt", walletEncryptHandler(gateway), map[string][]string{
		http.MethodPost: []string{EndpointsWallet},
	})
	webHandlerV1("/wallet/decrypt", walletDecryptHandler(gateway), map[string][]string{
		http.MethodPost: []string{EndpointsWallet},
	})
	webHandlerV2("/wallet/recover", walletRecoverHandler(gateway), map[string][]string{
		http.MethodPost: []string{EndpointsWallet},
	})

	// Blockchain interface
	webHandlerV1("/blockchain/metadata", blockchainMetadataHandler(gateway), map[string][]string{
		http.MethodGet: []string{EndpointsRead, EndpointsStatus},
	})
	webHandlerV1("/blockchain/progress", blockchainProgressHandler(gateway), map[string][]string{
		http.MethodGet: []string{EndpointsRead, EndpointsStatus},
	})
	webHandlerV1("/block", blockHandler(gateway), map[string][]string{
		http.MethodGet: []string{EndpointsRead},
	})
	webHandlerV1("/blocks", blocksHandler(gateway), map[string][]string{
		http.MethodGet:  []string{EndpointsRead},
		http.MethodPost: []string{EndpointsRead},
	})
	webHandlerV1("/last_blocks", lastBlocksHandler(gateway), map[string][]string{
		http.MethodGet: []string{EndpointsRead},
	})

	// Network stats endpoints
	webHandlerV1("/network/connection", connectionHandler(gateway), map[string][]string{
		http.MethodGet: []string{EndpointsRead, EndpointsStatus},
	})
	webHandlerV1("/network/connections", connectionsHandler(gateway), map[string][]string{
		http.MethodGet: []string{EndpointsRead, EndpointsStatus},
	})
	webHandlerV1("/network/defaultConnections", defaultConnectionsHandler(gateway), map[string][]string{
		http.MethodGet: []string{EndpointsRead, EndpointsStatus},
	})
	webHandlerV1("/network/connections/trust", trustConnectionsHandler(gateway), map[string][]string{
		http.MethodGet: []string{EndpointsRead, EndpointsStatus},
	})
	webHandlerV1("/network/connections/exchange", exchgConnectionsHandler(gateway), map[string][]string{
		http.MethodGet: []string{EndpointsRead, EndpointsStatus},
	})

	// Network admin endpoints
	webHandlerV1("/network/connection/disconnect", disconnectHandler(gateway), map[string][]string{
		http.MethodPost: []string{EndpointsNetCtrl},
	})

	// Transaction related endpoints
	webHandlerV1("/pendingTxs", pendingTxnsHandler(gateway), map[string][]string{
		http.MethodGet: []string{EndpointsRead},
	})
	webHandlerV1("/transaction", transactionHandler(gateway), map[string][]string{
		http.MethodGet: []string{EndpointsRead},
	})
	webHandlerV2("/transaction", transactionHandlerV2(gateway), map[string][]string{
		// http.MethodGet:  []string{EndpointsRead},
		http.MethodPost: []string{EndpointsTransaction},
	})
	webHandlerV2("/transaction/verify", verifyTxnHandler(gateway), map[string][]string{
		http.MethodPost: []string{EndpointsRead},
	})
	webHandlerV1("/transactions", transactionsHandler(gateway), map[string][]string{
		http.MethodGet:  []string{EndpointsRead},
		http.MethodPost: []string{EndpointsRead},
	})
	webHandlerV1("/injectTransaction", injectTransactionHandler(gateway), map[string][]string{
		http.MethodPost: []string{EndpointsTransaction, EndpointsWallet},
	})
	webHandlerV1("/resendUnconfirmedTxns", resendUnconfirmedTxnsHandler(gateway), map[string][]string{
		http.MethodPost: []string{EndpointsTransaction, EndpointsWallet},
	})
	webHandlerV1("/rawtx", rawTxnHandler(gateway), map[string][]string{
		http.MethodGet: []string{EndpointsRead},
	})

	// Unspent output related endpoints
	webHandlerV1("/outputs", outputsHandler(gateway), map[string][]string{
		http.MethodGet:  []string{EndpointsRead},
		http.MethodPost: []string{EndpointsRead},
	})
	webHandlerV1("/balance", balanceHandler(gateway), map[string][]string{
		http.MethodGet:  []string{EndpointsRead},
		http.MethodPost: []string{EndpointsRead},
	})
	webHandlerV1("/uxout", uxOutHandler(gateway), map[string][]string{
		http.MethodGet: []string{EndpointsRead},
	})
	webHandlerV1("/address_uxouts", addrUxOutsHandler(gateway), map[string][]string{
		http.MethodGet: []string{EndpointsRead},
	})

	// golang process internal metrics for Prometheus
	webHandlerV2("/metrics", metricsHandler(c, gateway), map[string][]string{
		http.MethodGet: []string{EndpointsPrometheus},
	})

	// Address related endpoints
	webHandlerV2("/address/verify", http.HandlerFunc(addressVerifyHandler), map[string][]string{
		http.MethodPost: []string{EndpointsRead},
	})

	// Explorer endpoints
	webHandlerV1("/coinSupply", coinSupplyHandler(gateway), map[string][]string{
		http.MethodGet: []string{EndpointsRead},
	})
	webHandlerV1("/richlist", richlistHandler(gateway), map[string][]string{
		http.MethodGet: []string{EndpointsRead},
	})
	webHandlerV1("/addresscount", addressCountHandler(gateway), map[string][]string{
		http.MethodGet: []string{EndpointsRead},
	})

	// Storage endpoint
	webHandlerV2("/data", storageHandler(gateway), map[string][]string{
		http.MethodGet:    []string{EndpointsStorage},
		http.MethodPost:   []string{EndpointsStorage},
		http.MethodDelete: []string{EndpointsStorage},
	})

	return mux
}

// newIndexHandler returns a http.Handler for index.html, where index.html is in appLoc
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

// parseAddressesFromStr parses comma-separated addresses string into []cipher.Address
func parseAddressesFromStr(s string) ([]cipher.Address, error) {
	addrsStr := splitCommaString(s)

	addrs := make([]cipher.Address, len(addrsStr))
	for i, s := range addrsStr {
		a, err := cipher.DecodeBase58Address(s)
		if err != nil {
			return nil, fmt.Errorf("address %q is invalid: %v", s, err)
		}

		addrs[i] = a
	}

	return addrs, nil
}

// parseAddressesFromStr parses comma-separated hashes string into []cipher.SHA256
func parseHashesFromStr(s string) ([]cipher.SHA256, error) {
	hashesStr := splitCommaString(s)

	hashes := make([]cipher.SHA256, len(hashesStr))
	for i, s := range hashesStr {
		h, err := cipher.SHA256FromHex(s)
		if err != nil {
			return nil, fmt.Errorf("SHA256 hash %q is invalid: %v", s, err)
		}

		hashes[i] = h
	}

	return hashes, nil
}
