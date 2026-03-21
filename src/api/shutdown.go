package api

import (
	"net/http"

	wh "github.com/skycoin/skycoin/src/util/http"
)

// Shuts down the node
// URI: /api/v1/shutdown
// Method: POST
func shutdownHandler(quit chan struct{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			wh.Error405(w)
			return
		}

		logger.Info("Shutdown requested via API")
		wh.SendJSONOr500(logger, w, struct {
			Message string `json:"message"`
		}{
			Message: "shutting down",
		})

		// Close the quit channel to trigger graceful shutdown
		go func() {
			close(quit)
		}()
	}
}
