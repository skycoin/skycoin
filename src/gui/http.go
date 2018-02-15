package gui

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
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
)

// Server exposes an HTTP API
type Server struct {
	mux      *http.ServeMux
	listener net.Listener
	done     chan struct{}
}

type ServerConfig struct {
	StaticDir   string
	DisableCSRF bool
}

func create(host string, serverConfig ServerConfig, daemon *daemon.Daemon) (*Server, error) {
	appLoc, err := file.DetermineResourcePath(serverConfig.StaticDir, resourceDir, devDir)
	if err != nil {
		return nil, err
	}
	logger.Info("Web resources directory: %s", appLoc)

	csrfStore := &CSRFStore{
		Enabled: !serverConfig.DisableCSRF,
	}
	if serverConfig.DisableCSRF {
		logger.Warning("CSRF check disabled")
	}

	return &Server{
		mux:  NewServerMux(host, appLoc, daemon.Gateway, csrfStore),
		done: make(chan struct{}),
	}, nil
}

// Create creates a new Server instance that listens on HTTP
func Create(host string, serverConfig ServerConfig, daemon *daemon.Daemon) (*Server, error) {
	s, err := create(host, serverConfig, daemon)
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
func CreateHTTPS(host string, serverConfig ServerConfig, daemon *daemon.Daemon, certFile, keyFile string) (*Server, error) {
	s, err := create(host, serverConfig, daemon)
	if err != nil {
		return nil, err
	}

	logger.Info("Using %s for the certificate", certFile)
	logger.Info("Using %s for the key", keyFile)

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
	logger.Info("Starting web interface on %s", s.listener.Addr())
	defer logger.Info("Web interface closed")
	defer close(s.done)

	if err := http.Serve(s.listener, s.mux); err != nil {
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

// NewServerMux creates an http.ServeMux with handlers registered
func NewServerMux(host, appLoc string, gateway Gatewayer, csrfStore *CSRFStore) *http.ServeMux {
	mux := http.NewServeMux()

	headerCheck := func(host string, handler http.Handler) http.Handler {
		handler = OriginRefererCheck(host, handler)
		handler = wh.HostCheck(logger, host, handler)
		return handler
	}

	webHandler := func(endpoint string, handler http.Handler) {
		handler = CSRFCheck(csrfStore, handler)
		handler = headerCheck(host, handler)
		mux.Handle(endpoint, handler)
	}

	webHandler("/", newIndexHandler(appLoc))

	fileInfos, _ := ioutil.ReadDir(appLoc)
	for _, fileInfo := range fileInfos {
		route := fmt.Sprintf("/%s", fileInfo.Name())
		if fileInfo.IsDir() {
			route = route + "/"
		}
		webHandler(route, http.FileServer(http.Dir(appLoc)))
	}

	// get the current CSRF token
	mux.Handle("/csrf", headerCheck(host, getCSRFToken(gateway, csrfStore)))

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
	//  hours: Number of hours to spends
	//  fee: Number of hours to use as fee, on top of the default fee.
	//  Returns total amount spent if successful, otherwise error describing
	//  failure status.
	webHandler("/wallet/spend", walletSpendHandler(gateway))

	// GET Arguments:
	//      id: Wallet ID
	// Returns all pending transanction for all addresses by selected Wallet
	webHandler("/wallet/transactions", walletTransactionsHandler(gateway))

	// Update wallet label
	//      POST Arguments:
	//          id: wallet id
	//          label: wallet label
	webHandler("/wallet/update", walletUpdateHandler(gateway))

	// Returns all loaded wallets
	// returns sensitive information
	webHandler("/wallets", walletsHandler(gateway))

	webHandler("/wallets/folderName", getWalletFolder(gateway))

	// generate wallet seed
	webHandler("/wallet/newSeed", newWalletSeed(gateway))

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
	webHandler("/network/connections/trust", trustConnectionsHandler(gateway))
	webHandler("/network/connections/exchange", exchgConnectionsHandler(gateway))

	// Transaction handler

	// get set of pending transactions
	webHandler("/pendingTxs", getPendingTxs(gateway))
	// get latest confirmed transactions
	webHandler("/lastTxs", getLastTxs(gateway))
	// get txn by txid
	webHandler("/transaction", getTransactionByID(gateway))

	// Returns transactions that match the filters.
	// Method: GET
	// Args:
	//     addrs: Comma seperated addresses [optional, returns all transactions if no address is provided]
	//     confirmed: Whether the transactions should be confirmed [optional, must be 0 or 1; if not provided, returns all]
	webHandler("/transactions", getTransactions(gateway))
	//inject a transaction into network
	webHandler("/injectTransaction", injectTransaction(gateway))
	webHandler("/resendUnconfirmedTxns", resendUnconfirmedTxns(gateway))
	// get raw tx by txid.
	webHandler("/rawtx", getRawTx(gateway))

	// UxOUt api handler

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
func newIndexHandler(appLoc string) http.HandlerFunc {
	// Serves the main page
	return func(w http.ResponseWriter, r *http.Request) {
		page := filepath.Join(appLoc, indexPage)
		logger.Debug("Serving index page: %s", page)
		if r.URL.Path == "/" {
			http.ServeFile(w, r, page)
		} else {
			wh.Error404(w)
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
			logger.Error("get unspent outputs failed: %v", err)
			wh.Error500(w)
			return
		}

		wh.SendOr404(w, outs)
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
			errMsg := fmt.Sprintf("Get balance failed: %v", err)
			logger.Error("%s", errMsg)
			wh.Error500Msg(w, errMsg)
			return
		}

		var balance wallet.BalancePair
		for _, bal := range bals {
			var err error
			balance.Confirmed, err = balance.Confirmed.Add(bal.Confirmed)
			if err != nil {
				wh.Error500Msg(w, err.Error())
				return
			}

			balance.Predicted, err = balance.Predicted.Add(bal.Predicted)
			if err != nil {
				wh.Error500Msg(w, err.Error())
				return
			}
		}

		wh.SendOr404(w, balance)
	}
}

func versionHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		wh.SendOr404(w, gateway.GetBuildInfo())
	}
}

/*
attrActualLog remove color char in log
origin: "\u001b[36m[skycoin.daemon:DEBUG] Trying to connect to 47.88.33.156:6000\u001b[0m",
*/
func attrActualLog(logInfo string) string {
	//return logInfo
	var actualLog string
	actualLog = logInfo
	if strings.HasPrefix(logInfo, "[skycoin") {
		if strings.Contains(logInfo, "\u001b") {
			actualLog = logInfo[0 : len(logInfo)-4]
		}
	} else {
		if len(logInfo) > 5 {
			if strings.Contains(logInfo, "\u001b") {
				actualLog = logInfo[5 : len(logInfo)-4]
			}
		}
	}
	return actualLog
}
func getLogsHandler(logbuf *bytes.Buffer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		var err error
		defaultLineNum := 1000 // default line numbers
		linenum := defaultLineNum
		if lines := r.FormValue("lines"); lines != "" {
			linenum, err = strconv.Atoi(lines)
			if err != nil {
				linenum = defaultLineNum
			}
		}
		keyword := r.FormValue("include")
		excludeKeyword := r.FormValue("exclude")
		logs := []string{}
		logList := strings.Split(logbuf.String(), "\n")
		for _, logInfo := range logList {
			if excludeKeyword != "" && strings.Contains(logInfo, excludeKeyword) {
				continue
			}
			if keyword != "" && !strings.Contains(logInfo, keyword) {
				continue
			}

			if len(logs) >= linenum {
				logger.Debug("logs size %d,total size:%d", len(logs), len(logList))
				break
			}
			log := attrActualLog(logInfo)
			if "" != log {
				logs = append(logs, log)
			}

		}

		wh.SendOr404(w, logs)
	}
}
