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

	// swagger:operation POST /api/v2/address/verify verifyAddress
	//
	// Verifies a Skycoin address.
	//
	// ---
	//
	// produces:
	// - application/json
	// parameters:
	// - name: address
	//   in: query
	//   description: Address id.
	//   required: true
	//   type: string
	//
	// security:
	// - csrfAuth: []
	//
	//
	// responses:
	//   200:
	//     description: Response verifies a Skycoin address
	//     schema:
	//       type: object
	//       properties:
	//         error:
	//           type: object
	//         data:
	//           type: object
	//           properties:
	//             version:
	//               type: integer
	//               format: int64
	//   default:
	//     $ref: '#/responses/genericError'

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
