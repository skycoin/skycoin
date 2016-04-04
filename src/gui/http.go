package gui

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"path/filepath"

	"github.com/skycoin/skycoin/src/daemon"
	"gopkg.in/op/go-logging.v1"
)

var (
	logger      = logging.MustGetLogger("skycoin.gui")
	resourceDir = "dist/"
	indexPage   = "index.html"
)

// Begins listening on http://$host, for enabling remote web access
// Does NOT use HTTPS
func LaunchWebInterface(host, staticDir string, daemon *daemon.Daemon) error {
	logger.Warning("Starting web interface on http://%s", host)
	logger.Warning("HTTPS not in use!")
	appLoc := filepath.Join(staticDir, resourceDir)
	mux := NewGUIMux(appLoc, daemon)
	//if err := http.ListenAndServe(host, mux); err != nil {
	//	log.Panic(err)
	//}
	web_interface_active := make(chan bool, 1) //do not return until webserver is running
	listener, err := net.Listen("tcp", host)
	if err != nil {
		log.Panic(err)
	}
	go func() {
		web_interface_active <- true
		err = http.Serve(listener, mux) //blocks
		if err != nil {
			log.Panic()
		}
	}()
	value := <-web_interface_active
	if value == true {
		log.Printf("webservice should be running: RUN POPUP")
	}
	return nil
}

// Begins listening on https://$host, for enabling remote web access
// Uses HTTPS
func LaunchWebInterfaceHTTPS(host, staticDir string, daemon *daemon.Daemon,
	certFile, keyFile string) {
	logger.Info("Starting web interface on https://%s", host)
	logger.Info("Using %s for the certificate", certFile)
	logger.Info("Using %s for the key", keyFile)
	appLoc := filepath.Join(staticDir, resourceDir)
	mux := NewGUIMux(appLoc, daemon)
	err := http.ListenAndServeTLS(host, certFile, keyFile, mux)
	if err != nil {
		log.Panic(err)
	}
}

// Creates an http.ServeMux with handlers registered
func NewGUIMux(appLoc string, daemon *daemon.Daemon) *http.ServeMux {
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

	// Wallet interface
	RegisterWalletHandlers(mux, daemon.Gateway)
	// Blockchain interface
	RegisterBlockchainHandlers(mux, daemon.Gateway)
	// Network stats interface
	RegisterNetworkHandlers(mux, daemon.Gateway)
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
			Error404(w)
		}
	}
}
