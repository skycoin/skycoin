package api

import (
	"encoding/json"
	"net/http"

	"github.com/skycoin/skycoin/src/cipher"
)

// VerifyAddressRequest is the request data for POST /api/v2/address/verify
type VerifyAddressRequest struct {
	Address string `json:"address"`
}

// VerifyAddressResponse is returned by POST /api/v2/address/verify
type VerifyAddressResponse struct {
	Version byte `json:"version"`
}

// addressVerifyHandler verifies a Skycoin address
// Method: POST
// URI: /api/v2/address/verify
func addressVerifyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		resp := NewHTTPErrorResponse(http.StatusMethodNotAllowed, "")
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
