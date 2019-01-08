package api

import (
	"encoding/json"
	"net/http"

	"github.com/skycoin/skycoin/src/cipher"
)

// VerifyAddressRequest is the request data for POST /api/v2/address/verify
// swagger:parameters verifyAddressRequest
type VerifyAddressRequest struct {
	Address string `json:"address"`
}

// VerifyAddressResponse is returned by POST /api/v2/address/verify
// swagger:response verifyAddressResponse
type VerifyAddressResponse struct {
	Version byte `json:"version"`
}

// addressVerifyHandler verifies a Skycoin address
// Method: POST
// URI: /api/v2/address/verify
func addressVerifyHandler(w http.ResponseWriter, r *http.Request) {

	// swagger:route POST /api/v2/address/verify addressVerify
	//
	// addressVerifyHandler verifies a Skycoin address
	//
	//     Consumes:
	//     - application/json
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

	if r.Method != http.MethodPost {
		resp := NewHTTPErrorResponse(http.StatusMethodNotAllowed, "")
		writeHTTPResponse(w, resp)
		return
	}

	if r.Header.Get("Content-Type") != ContentTypeJSON {
		resp := NewHTTPErrorResponse(http.StatusUnsupportedMediaType, "")
		writeHTTPResponse(w, resp)
		return
	}

	var req VerifyAddressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		resp := NewHTTPErrorResponse(http.StatusBadRequest, err.Error())
		writeHTTPResponse(w, resp)
		return
	}

	if req.Address == "" {
		resp := NewHTTPErrorResponse(http.StatusBadRequest, "address is required")
		writeHTTPResponse(w, resp)
		return
	}

	addr, err := cipher.DecodeBase58Address(req.Address)

	if err != nil {
		resp := NewHTTPErrorResponse(http.StatusUnprocessableEntity, err.Error())
		writeHTTPResponse(w, resp)
		return
	}

	writeHTTPResponse(w, HTTPResponse{
		Data: VerifyAddressResponse{
			Version: addr.Version,
		},
	})
}
