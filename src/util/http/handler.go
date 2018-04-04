package httphelper

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/skycoin/skycoin/src/util/iputil"
	"github.com/skycoin/skycoin/src/util/logging"
)

// HostCheck checks that the request's Host header is 127.0.0.1:$port or localhost:$port
// if the HTTP interface host is also a localhost address.
// This prevents DNS rebinding attacks, where an attacker uses a DNS rebinding service
// to bypass CORS checks.
// If the HTTP interface host is not a localhost address,
// the Host header is not checked. This is considered a public interface.
// If the Host header is not set, it is not checked.
// All major browsers send the Host header as required by the HTTP spec.
// TODO: move this back into gui/ library after webrpc interface is removed
func HostCheck(logger *logging.Logger, host string, handler http.Handler) http.Handler {
	addr := host
	var port uint16
	if strings.Contains(host, ":") {
		var err error
		addr, port, err = iputil.SplitAddr(host)
		if err != nil {
			log.Panic(err)
		}
	}

	isLocalhost := iputil.IsLocalhost(addr)

	if isLocalhost && port == 0 {
		log.Panic("localhost with no port specified is unsupported")
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// NOTE: The "Host" header is not in http.Request.Header, it's put in the http.Request.Host field
		if r.Host != "" && isLocalhost && r.Host != fmt.Sprintf("127.0.0.1:%d", port) && r.Host != fmt.Sprintf("localhost:%d", port) {
			logger.Criticalf("Detected DNS rebind attempt - configured-host=%s header-host=%s", host, r.Host)
			Error403(w)
			return
		}

		handler.ServeHTTP(w, r)
	})
}
