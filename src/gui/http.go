package gui

import (
    "fmt"
    "github.com/lonnc/golang-nw"
    "github.com/op/go-logging"
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
    if err := nodeWebkit.ListenAndServe(NewGUIMux(resourceDir)); err != nil {
        log.Panic(err)
    }
}

// Begins listening on addr:port, for enabling remote web access
func LaunchWebInterface(addr string, port int, staticDir string) {
    logger.Info("Starting web interface on http://%s:%d", addr, port)
    if addr != "localhost" && addr != "127.0.0.1" {
        log.Panic("Remote web interface is not supported yet, " +
            "needs TLS support")
    }
    a := fmt.Sprintf("%s:%d", addr, port)
    // TODO -- use ListenAndServeTLS. Will need to generate a pem file
    // and allow the user to override it with their own
    appLoc := filepath.Join(staticDir, resourceDir)
    if err := http.ListenAndServe(a, NewGUIMux(appLoc)); err != nil {
        log.Panic(err)
    }
}

// Creates an http.ServeMux with handlers registered
func NewGUIMux(appLoc string) *http.ServeMux {
    mux := http.NewServeMux()
    mux.HandleFunc("/", newIndexHandler(appLoc))
    for _, s := range resources {
        route := fmt.Sprintf("/%s/", s)
        mux.HandleFunc(route, newStaticHandler(appLoc))
    }
    // Wallet interface
    RegisterWalletHandlers(mux)
    // Network stats interface
    RegisterNetworkHandlers(mux)
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
