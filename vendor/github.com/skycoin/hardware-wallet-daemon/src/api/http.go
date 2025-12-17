package api

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"regexp"
	"runtime"

	"github.com/NYTimes/gziphandler"
	"github.com/rs/cors"
	skyWallet "github.com/skycoin/hardware-wallet-go/src/skywallet"
	wh "github.com/skycoin/skycoin/src/util/http"
	"github.com/skycoin/skycoin/src/util/logging"
)

const (
	// ContentTypeJSON json content type header
	ContentTypeJSON = "application/json"
	// ContentTypeForm form data content type header
	ContentTypeForm = "application/x-www-form-urlencoded"

	apiVersion1 = "v1"
)

var (
	logger = logging.MustGetLogger("daemon-api")

	// custom lock to help with serializing requests
	ongoingOperation chan struct{}
)

// corsRegex matches all localhost origin headers
var corsRegex *regexp.Regexp

// autoPressEmulatorButtons is used to automatically press emulator buttons
// Used in integration testing
var autoPressEmulatorButtons bool

func init() {
	var err error
	corsRegex, err = regexp.Compile(`(^https?:\/\/)?^?(localhost|127\.0\.0\.1):\d+$`)
	if err != nil {
		logger.Panic(err)
	}

	// set size to 1 to allow only 1 request at a time
	ongoingOperation = make(chan struct{}, 1)

	apb := os.Getenv("AUTO_PRESS_BUTTONS")
	if apb == "1" && runtime.GOOS == "linux" {
		autoPressEmulatorButtons = true
	}
}

// Config configures Server
type Config struct {
	EnableCSRF         bool
	DisableHeaderCheck bool
	HostWhitelist      []string
	Mode               skyWallet.DeviceType
	Build              BuildInfo
}

type muxConfig struct {
	host               string
	enableCSRF         bool
	disableHeaderCheck bool
	hostWhitelist      []string
	mode               skyWallet.DeviceType
	build              BuildInfo
}

// Server exposes an HTTP API
type Server struct {
	server   *http.Server
	listener net.Listener
	done     chan struct{}
}

// Serve serves the web interface on the configured host
func (s *Server) Serve() error {
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

func create(host string, c Config, gateway *Gateway) *Server {
	mc := muxConfig{
		host:               host,
		enableCSRF:         c.EnableCSRF,
		disableHeaderCheck: c.DisableHeaderCheck,
		hostWhitelist:      c.HostWhitelist,
		mode:               c.Mode,
		build:              c.Build,
	}

	srvMux := newServerMux(mc, gateway.Device)

	srv := &http.Server{
		Handler: srvMux,
	}

	return &Server{
		server: srv,
		done:   make(chan struct{}),
	}
}

// Create create a new http server
func Create(host string, c Config, gateway *Gateway) (*Server, error) {
	listener, err := net.Listen("tcp", host)
	if err != nil {
		return nil, err
	}

	// If the host did not specify a port, allowing the kernel to assign one,
	// we need to get the assigned address to know the full hostname
	host = listener.Addr().String()

	s := create(host, c, gateway)

	s.listener = listener

	return s, nil
}

func newServerMux(c muxConfig, gateway Gatewayer) *http.ServeMux {
	mux := http.NewServeMux()

	allowedOrigins := []string{
		fmt.Sprintf("http://%s", c.host),
		"https://staging.wallet.skycoin.net",
		"https://wallet.skycoin.net",
	}

	for _, s := range c.hostWhitelist {
		allowedOrigins = append(allowedOrigins, fmt.Sprintf("http://%s", s))
	}

	corsValidator := func(origin string) bool {
		if corsRegex.MatchString(origin) {
			return true
		}

		for _, allowedOrigin := range allowedOrigins {
			if allowedOrigin == origin {
				return true
			}
		}

		return false
	}

	corsHandler := cors.New(cors.Options{
		AllowOriginFunc:    corsValidator,
		Debug:              false,
		AllowedMethods:     []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut},
		AllowedHeaders:     []string{"Origin", "Accept", "Content-Type", "X-Requested-With", CSRFHeaderName},
		AllowCredentials:   false, // credentials are not used, but it would be safe to enable if necessary
		OptionsPassthrough: false,
	})

	headerCheck := func(host string, hostWhitelist []string, handler http.Handler) http.Handler {
		handler = originRefererCheck(host, hostWhitelist, handler)
		handler = hostCheck(host, hostWhitelist, handler)
		return handler
	}

	webHandlerWithOptionals := func(endpoint string, handlerFunc http.Handler, checkCSRF, checkHeaders bool) {
		handler := wh.ElapsedHandler(logger, handlerFunc)

		handler = corsHandler.Handler(handler)

		if checkCSRF {
			handler = CSRFCheck(c.enableCSRF, handler)
		}

		if checkHeaders {
			handler = headerCheck(c.host, c.hostWhitelist, handler)
		}

		handler = gziphandler.GzipHandler(handler)

		mux.Handle(endpoint, handler)
	}

	webHandler := func(endpoint string, handler http.Handler) {
		webHandlerWithOptionals(endpoint, handler, c.enableCSRF, !c.disableHeaderCheck)
	}

	webHandlerV1 := func(endpoint string, handler http.Handler) {
		webHandler("/api/"+apiVersion1+endpoint, handler)
	}

	if autoPressEmulatorButtons && c.mode != skyWallet.DeviceTypeEmulator {
		logger.Panic("auto press buttons enabled but device mode is not emulator")
	}

	// get the current CSRF token
	csrfHandlerV1 := func(endpoint string, handler http.Handler) {
		webHandlerWithOptionals("/api/"+apiVersion1+endpoint, handler, false, !c.disableHeaderCheck)
	}
	csrfHandlerV1("/csrf", getCSRFToken(c.enableCSRF)) // csrf is always available, regardless of the API set

	// hw daemon endpoints
	webHandlerV1("/generate_addresses", generateAddresses(gateway))
	webHandlerV1("/apply_settings", applySettings(gateway))
	webHandlerV1("/backup", backup(gateway))
	webHandlerV1("/cancel", cancel(gateway))
	webHandlerV1("/check_message_signature", checkMessageSignature(gateway))
	webHandlerV1("/features", features(gateway))
	// enable firmware update endpoint only for hw wallet
	if c.mode == skyWallet.DeviceTypeUSB {
		webHandlerV1("/firmware_update", firmwareUpdate(gateway))
		webHandlerV1("/available", available(gateway))
	}
	webHandlerV1("/generate_mnemonic", generateMnemonic(gateway))
	webHandlerV1("/recovery", recovery(gateway))
	webHandlerV1("/set_mnemonic", setMnemonic(gateway))
	webHandlerV1("/configure_pin_code", configurePinCode(gateway))
	webHandlerV1("/sign_message", signMessage(gateway))
	webHandlerV1("/transaction_sign", transactionSign(gateway))
	webHandlerV1("/wipe", wipe(gateway))

	webHandlerV1("/intermediate/pin_matrix", pinMatrixRequestHandler(gateway))
	webHandlerV1("/intermediate/passphrase", passphraseRequestHandler(gateway))
	webHandlerV1("/intermediate/word", wordRequestHandler(gateway))
	webHandlerV1("/intermediate/button", buttonRequestHandler(gateway))

	webHandlerV1("/version", versionHandler(c))
	return mux
}
