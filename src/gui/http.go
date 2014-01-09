package gui

import (
    "fmt"
    "log"
    "net/http"
    "path/filepath"

    "github.com/lonnc/golang-nw"
    "github.com/op/go-logging"
)

var (
    logger      = logging.MustGetLogger("skycoin.gui")
    resources   = []string{"js", "css", "lib", "partials", "img", "assets"}
    resourceDir = "app/"
    indexPage   = filepath.Join(resourceDir, "index.html")
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
    for _, s := range resources {
        route := fmt.Sprintf("/%s/", s)
        mux.HandleFunc(route, staticHandler)
    }
    // Wallet interface
    RegisterWalletHandlers(mux)
    // Network stats interface
    RegisterNetworkHandlers(mux)
    return mux
}

// Serves the main page
func indexHandler(w http.ResponseWriter, r *http.Request) {
    logger.Debug("Serving %s\n", indexPage)
    http.ServeFile(w, r, indexPage)
}

// Serves files out of ./static/
func staticHandler(w http.ResponseWriter, r *http.Request) {
    fp := filepath.Join(resourceDir, r.URL.Path[1:])
    logger.Debug("Serving %s\n", fp)
    http.ServeFile(w, r, fp)
}
