package gui

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/skycoin/skycoin/src/daemon"
)

var (
	logger      = logging.MustGetLogger("skycoin.gui")
	resourceDir = "dist/"
	indexPage   = "index.html"
)

// Begins listening on http://$host, for enabling remote web access
// Does NOT use HTTPS
func LaunchWebInterface(host, staticDir string, daemon *daemon.Daemon) error {
	logger.Info("Starting web interface on http://%s", host)
	logger.Warning("HTTPS not in use!")
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Panic(err)
	}
	appLoc := filepath.Join(dir, staticDir, resourceDir)
	mux := NewGUIMux(appLoc, daemon)
	//if err := http.ListenAndServe(host, mux); err != nil {
	//	log.Panic(err)
	//}
	webInterfaceActive := make(chan bool, 1) //do not return until webserver is running
	listener, err := net.Listen("tcp", host)
	if err != nil {
		log.Panic(err)
	}
	go func() {
		webInterfaceActive <- true
		err = http.Serve(listener, mux) //blocks
		if err != nil {
			log.Panic()
		}
	}()
	value := <-webInterfaceActive
	if value == true {
		log.Printf("webservice should be running: RUN POPUP")
	}
	return nil
}

// Begins listening on https://$host, for enabling remote web access
// Uses HTTPS
func LaunchWebInterfaceHTTPS(host, staticDir string, daemon *daemon.Daemon,
	certFile, keyFile string) error {
	logger.Info("Starting web interface on https://%s", host)
	logger.Info("Using %s for the certificate", certFile)
	logger.Info("Using %s for the key", keyFile)

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Panic(err)
	}
	appLoc := filepath.Join(dir, staticDir, resourceDir)
	mux := NewGUIMux(appLoc, daemon)

	//err := http.ListenAndServeTLS(host, certFile, keyFile, mux)
	//if err != nil {
	//	log.Panic(err)
	//}
	config := new(tls.Config)
	config.Certificates = make([]tls.Certificate, 1)
	config.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Panic(err)
	}

	// do not return until webserver is running
	webInterfaceActive := make(chan bool, 1)
	listener, err := tls.Listen("tcp", host, config)
	if err != nil {
		log.Panic(err)
	}

	go func() {
		webInterfaceActive <- true
		err = http.Serve(listener, mux) //blocks
		if err != nil {
			log.Panic()
		}
	}()

	value := <-webInterfaceActive
	if value {
		log.Printf("webservice should be running: RUN POPUP")
	}
	return nil
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
