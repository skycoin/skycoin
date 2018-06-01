package api

import (
	"encoding/json"
	"net/http"

	"github.com/skycoin/skycoin/src/cipher"
	//http,json helpers
)

// VerifyAddressRequest is the request data for POST /api/v1/address/verify
type VerifyAddressRequest struct {
	Address string `json:"address"`
}

// VerifyAddressResponse is returned by POST /api/v1/address/verify
type VerifyAddressResponse struct {
	Address string `json:"address"`
	Version *byte  `json:"version,omitempty"`
	Valid   bool   `json:"valid"`
}

// addressVerify verifies a Skycoin address
// Method: POST
// URI: /api/v1/address/verify
func addressVerify(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		resp := NewHTTPErrorResponse(http.StatusMethodNotAllowed, "")
		writeHTTPResponse(w, resp)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
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
		writeHTTPResponse(w, HTTPResponse{
			Error: &HTTPError{
				Message: err.Error(),
				Code:    http.StatusUnprocessableEntity,
			},
			Data: VerifyAddressResponse{
				Address: req.Address,
				Valid:   false,
			},
		})
		return
	}

	writeHTTPResponse(w, HTTPResponse{
		Data: VerifyAddressResponse{
			Address: req.Address,
			Version: &addr.Version,
			Valid:   true,
		},
	})
}
