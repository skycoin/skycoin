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

//var (
//	logger   = logging.MustGetLogger("gui")
//	listener net.Listener
//	quit     chan struct{}
//)

type server struct {
	daemon *daemon.Daemon
	logger *logging.Logger
	listener net.Listener
	quit chan struct{}
}

const (
	resourceDir = "dist/"
	devDir      = "dev/"
	indexPage   = "index.html"
)

// LaunchWebInterface begins listening on http://$host, for enabling remote web access
// Does NOT use HTTPS
func LaunchWebInterface(host, staticDir string, daemon *daemon.Daemon) *server {
	s := server{
		daemon: daemon,
		quit: make(chan struct{}),
		logger: logging.MustGetLogger("gui"),
	}
	s.logger.Info("Starting web interface on http://%s", host)
	s.logger.Warning("HTTPS not in use!")
	appLoc, err := file.DetermineResourcePath(staticDir, resourceDir, devDir)
	if err != nil {
		return nil
	}
	s.logger.Info("Web resources directory: %s", appLoc)

	s.listener, err = net.Listen("tcp", host)
	if err != nil {
		return nil
	}

	// Runs http.Serve() in a goroutine
	s.serve(NewGUIMux(appLoc, daemon, &s))
	return &s
}

// LaunchWebInterfaceHTTPS begins listening on https://$host, for enabling remote web access
// Uses HTTPS
func LaunchWebInterfaceHTTPS(host, staticDir string, daemon *daemon.Daemon, certFile, keyFile string) *server {
	s := server{
		quit: make(chan struct{}),
		logger: logging.MustGetLogger("gui"),
	}
	s.logger.Info("Starting web interface on https://%s", host)
	s.logger.Info("Using %s for the certificate", certFile)
	s.logger.Info("Using %s for the key", keyFile)
	s.logger.Info("Web resources directory: %s", staticDir)

	appLoc, err := file.DetermineResourcePath(staticDir, devDir, resourceDir)
	if err != nil {
		return nil
	}

	certs := make([]tls.Certificate, 1)
	if certs[0], err = tls.LoadX509KeyPair(certFile, keyFile); err != nil {
		return nil
	}

	s.listener, err = tls.Listen("tcp", host, &tls.Config{Certificates: certs})
	if err != nil {
		return nil
	}

	// Runs http.Serve() in a goroutine
	s.serve(NewGUIMux(appLoc, daemon, &s))
	return &s
}

func (s *server) serve(mux *http.ServeMux) {
	go func() {
		for {
			if err := http.Serve(s.listener, mux); err != nil {
				select {
				case <-s.quit:
					return
				default:
				}
				continue
			}
		}
	}()
}

// Shutdown close http service
func (s *server)Shutdown() {
	if s.quit != nil {
		// must close quit first
		close(s.quit)
		s.listener.Close()
		s.listener = nil
	}
}

// NewGUIMux creates an http.ServeMux with handlers registered
func (s *server)NewGUIMux(appLoc string, daemon *daemon.Daemon) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.newIndexHandler(appLoc))

	fileInfos, _ := ioutil.ReadDir(appLoc)
	for _, fileInfo := range fileInfos {
		route := fmt.Sprintf("/%s", fileInfo.Name())
		if fileInfo.IsDir() {
			route = route + "/"
		}
		mux.Handle(route, http.FileServer(http.Dir(appLoc)))
	}

	mux.HandleFunc("/logs", s.getLogsHandler(&daemon.LogBuff))

	mux.HandleFunc("/version", s.versionHandler(daemon.Gateway))

	//get set of unspent outputs
	mux.HandleFunc("/outputs", s.getOutputsHandler(daemon.Gateway))

	// get balance of addresses
	mux.HandleFunc("/balance", s.getBalanceHandler(daemon.Gateway))

	// Wallet interface
	s.RegisterWalletHandlers(mux, s.daemon.Gateway)
	// Blockchain interface
	RegisterBlockchainHandlers(mux, s.daemon.Gateway)
	// Network stats interface
	RegisterNetworkHandlers(mux, s.daemon.Gateway)
	// Transaction handler
	RegisterTxHandlers(mux, s.daemon.Gateway)
	// UxOUt api handler
	RegisterUxOutHandlers(mux, s.daemon.Gateway)
	// expplorer handler
	RegisterExplorerHandlers(mux, s.daemon.Gateway)
	return mux
}

// Returns a http.HandlerFunc for index.html, where index.html is in appLoc
func (s *server) newIndexHandler(appLoc string) http.HandlerFunc {
	// Serves the main page
	return func(w http.ResponseWriter, r *http.Request) {
		page := filepath.Join(appLoc, indexPage)
		s.logger.Debug("Serving index page: %s", page)
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
func (s *server) getOutputsHandler(gateway *daemon.Gateway) http.HandlerFunc {
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
			s.logger.Error("get unspent outputs failed: %v", err)
			wh.Error500(w)
			return
		}

		wh.SendOr404(w, outs)
	}
}

func (s *server) getBalanceHandler(gateway *daemon.Gateway) http.HandlerFunc {
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
			s.logger.Error("Get balance failed: %v", err)
			wh.Error500(w)
			return
		}

		var balance wallet.BalancePair
		for _, bal := range bals {
			balance.Confirmed = balance.Confirmed.Add(bal.Confirmed)
			balance.Predicted = balance.Predicted.Add(bal.Predicted)
		}

		wh.SendOr404(w, balance)
	}
}

func (s *server) versionHandler(gateway *daemon.Gateway) http.HandlerFunc {
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
func (s *server) getLogsHandler(logbuf *bytes.Buffer) http.HandlerFunc {
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
				s.logger.Debug("logs size %d,total size:%d", len(logs), len(logList))
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
