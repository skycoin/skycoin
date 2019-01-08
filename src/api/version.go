package api

import (
	"net/http"

	"github.com/skycoin/skycoin/src/readable"
	wh "github.com/skycoin/skycoin/src/util/http"
)

// versionHandler returns the application version info
// URI: /api/v1/version
// Method: GET
func versionHandler(bi readable.BuildInfo) http.HandlerFunc {

	// swagger:route GET /api/v1/version version
	//
	// versionHandler returns the application version info
	//
	//     Produces:
	//     - application/json
	//
	//     Schemes: http, https
	//
	//     Security:
	//       api_key:
	//       oauth: read, write
	//
	//     Responses:
	//       default: genericError
	//       200: OK

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		wh.SendJSONOr500(logger, w, bi)
	}
}
