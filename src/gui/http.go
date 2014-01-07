package gui

import (
    "fmt"
    "github.com/op/go-logging"
    "github.com/lonnc/golang-nw"
    "log"
    "net/http"
)

var (
    logger = logging.MustGetLogger("skycoin.gui")
)

// Begins listening on the node-webkit localhost
func LaunchGUI() {
    // Create a link back to node-webkit using the environment variable
    // populated by golang-nw's node-webkit code
    nodeWebkit, err := nw.New()
    if err != nil {
        log.Panic(err)
    }

    // Pick a random localhost port, start listening for http requests using default handler
    // and send a message back to node-webkit to redirect
    logger.Info("Launching GUI server")
    if err := nodeWebkit.ListenAndServe(NewGUIMux()); err != nil {
        log.Panic(err)
    }
}

// Begins listening on addr:port, for enabling remote web access
func LaunchWebInterface(addr string, port int) {
    log.Panic("Web interface is not supported yet, needs TLS support")
    a := fmt.Sprintf("%s:%d", addr, port)
    // TODO -- use ListenAndServeTLS. Will need to generate a pem file
    // and allow the user to override it with their own
    if err := http.ListenAndServe(a, NewGUIMux()); err != nil {
        log.Panic(err)
    }
}

// Creates an http.ServeMux with handlers registered
func NewGUIMux() *http.ServeMux {
    mux := http.NewServeMux()
    mux.HandleFunc("/", indexHandler)
    mux.HandleFunc("/static/", staticHandler)
    // Wallet interface
    RegisterWalletHandlers(mux)
    // Network stats interface
    RegisterNetworkHandlers(mux)
    return mux
}

// Serves the main page
func indexHandler(w http.ResponseWriter, r *http.Request) {
    logger.Debug("Serving index.html\n")
    http.ServeFile(w, r, "./static/index.html")
}

// Serves files out of ./static/
func staticHandler(w http.ResponseWriter, r *http.Request) {
    fp := r.URL.Path[1:]
    logger.Debug("Serving %s\n", fp)
    http.ServeFile(w, r, fp)
}
