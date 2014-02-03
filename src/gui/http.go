package gui

import (
    "fmt"
    "github.com/lonnc/golang-nw"
    "github.com/op/go-logging"
    "github.com/skycoin/skycoin/src/daemon"
    "log"
    "net/http"
    "path/filepath"
)

var (
    logger      = logging.MustGetLogger("skycoin.gui")
    resources   = []string{"js", "css", "lib", "partials", "img", "assets"}
    resourceDir = "app/"
    indexPage   = "index.html"
)

type HTTPHandler func(w http.ResponseWriter, r *http.Request)

// Begins listening on the node-webkit localhost
func LaunchGUI(daemon *daemon.Daemon) {
    // Create a link back to node-webkit using the environment variable
    // populated by golang-nw's node-webkit code
    nodeWebkit, err := nw.New()
    if err != nil {
        log.Panic(err)
    }

    // Pick a random localhost port, start listening for http requests using
    // default handler and send a message back to node-webkit to redirect
    logger.Info("Launching GUI server")
    mux := NewGUIMux(resourceDir, daemon)
    if err := nodeWebkit.ListenAndServe(mux); err != nil {
        log.Panic(err)
    }
}

// Begins listening on http://$host, for enabling remote web access
// Does NOT use HTTPS
func LaunchWebInterface(host, staticDir string, daemon *daemon.Daemon) {
    logger.Warning("Starting web interface on http://%s", host)
    logger.Warning("HTTPS not in use!")
    appLoc := filepath.Join(staticDir, resourceDir)
    mux := NewGUIMux(appLoc, daemon)
    if err := http.ListenAndServe(host, mux); err != nil {
        log.Panic(err)
    }
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
    for _, s := range resources {
        route := fmt.Sprintf("/%s/", s)
        mux.HandleFunc(route, newStaticHandler(appLoc))
    }
    // Wallet interface
    RegisterWalletHandlers(mux, daemon.RPC)
    // Network stats interface
    RegisterNetworkHandlers(mux, daemon.RPC)
    return mux
}

// Returns a func(http.ResponseWriter, *http.Request) for index.html,
// where index.html is in appLoc
func newIndexHandler(appLoc string) func(http.ResponseWriter, *http.Request) {
    // Serves the main page
    return func(w http.ResponseWriter, r *http.Request) {
        page := filepath.Join(appLoc, indexPage)
        logger.Debug("Serving %s", page)
        http.ServeFile(w, r, page)
    }
}

// Returns a func(http.ResponseWriter, *http.Request) for files in
// appLoc/static/
func newStaticHandler(appLoc string) func(http.ResponseWriter, *http.Request) {
    // Serves files out of ./static/
    return func(w http.ResponseWriter, r *http.Request) {
        fp := filepath.Join(appLoc, r.URL.Path[1:])
        logger.Debug("Serving %s", fp)
        http.ServeFile(w, r, fp)
    }
}
