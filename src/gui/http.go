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

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/wallet"

	"github.com/skycoin/skycoin/src/util/file"
	wh "github.com/skycoin/skycoin/src/util/http" //http,json helpers

	"github.com/skycoin/skycoin/src/util/logging"
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

func create(host, staticDir string, daemon *daemon.Daemon) (*Server, error) {
	appLoc, err := file.DetermineResourcePath(staticDir, resourceDir, devDir)
	if err != nil {
		return nil, err
	}
	logger.Info("Web resources directory: %s", appLoc)

	return &Server{
		mux:  NewServerMux(appLoc, daemon),
		done: make(chan struct{}),
	}, nil
}

// Create creates a new Server instance that listens on HTTP
func Create(host, staticDir string, daemon *daemon.Daemon) (*Server, error) {
	s, err := create(host, staticDir, daemon)
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
func CreateHTTPS(host, staticDir string, daemon *daemon.Daemon, certFile, keyFile string) (*Server, error) {
	s, err := create(host, staticDir, daemon)
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
func NewServerMux(appLoc string, daemon *daemon.Daemon) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", newIndexHandler(appLoc))

	fileInfos, _ := ioutil.ReadDir(appLoc)
	for _, fileInfo := range fileInfos {
		route := fmt.Sprintf("/%s", fileInfo.Name())
		if fileInfo.IsDir() {
			route = route + "/"
		}
		mux.Handle(route, http.FileServer(http.Dir(appLoc)))
	}

	mux.HandleFunc("/logs", getLogsHandler(&daemon.LogBuff))

	mux.HandleFunc("/version", versionHandler(daemon.Gateway))

	//get set of unspent outputs
	mux.HandleFunc("/outputs", getOutputsHandler(daemon.Gateway))

	// get balance of addresses
	mux.HandleFunc("/balance", getBalanceHandler(daemon.Gateway))

	// Wallet interface
	RegisterWalletHandlers(mux, daemon.Gateway)
	// Blockchain interface
	RegisterBlockchainHandlers(mux, daemon.Gateway)
	// Network stats interface
	RegisterNetworkHandlers(mux, daemon.Gateway)
	// Transaction handler
	RegisterTxHandlers(mux, daemon.Gateway)
	// UxOUt api handler
	RegisterUxOutHandlers(mux, daemon.Gateway)
	// expplorer handler
	RegisterExplorerHandlers(mux, daemon.Gateway)
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

// getOutputsHandler get utxos base on the filters in url params.
// mode: GET
// url: /outputs?addrs=[:addrs]&hashes=[:hashes]
// if addrs and hashes are not specificed, return all unspent outputs.
// if both addrs and hashes are specificed, then both those filters are need to be matched.
// if only specify one filter, then return outputs match the filter.
func getOutputsHandler(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		var addrs []string
		var hashes []string

		trimSpace := func(vs []string) []string {
			for i := range vs {
				vs[i] = strings.TrimSpace(vs[i])
			}
			return vs
		}

		addrStr := r.FormValue("addrs")
		if addrStr != "" {
			addrs = trimSpace(strings.Split(addrStr, ","))
		}

		hashStr := r.FormValue("hashes")
		if hashStr != "" {
			hashes = trimSpace(strings.Split(hashStr, ","))
		}

		filters := []daemon.OutputsFilter{}
		if len(addrs) > 0 {
			filters = append(filters, daemon.FbyAddresses(addrs))
		}

		if len(hashes) > 0 {
			filters = append(filters, daemon.FbyHashes(hashes))
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

func getBalanceHandler(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		addrsParam := r.FormValue("addrs")
		addrsStr := strings.Split(addrsParam, ",")
		addrs := make([]cipher.Address, 0, len(addrsStr))
		for _, addr := range addrsStr {
			// trim space
			addr = strings.Trim(addr, " ")
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

func versionHandler(gateway *daemon.Gateway) http.HandlerFunc {
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
