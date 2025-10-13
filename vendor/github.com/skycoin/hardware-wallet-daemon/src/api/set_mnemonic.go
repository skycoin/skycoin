package api

import (
	"encoding/json"
	"net/http"

	"github.com/skycoin/hardware-wallet-go/src/skywallet/wire"

	skyWallet "github.com/skycoin/hardware-wallet-go/src/skywallet"
	"github.com/skycoin/skycoin/src/cipher/bip39"
)

// SetMnemonicRequest is request data for /api/v1/set_mnemonic
type SetMnemonicRequest struct {
	Mnemonic string `json:"mnemonic"`
}

// URI: /api/v1/set_mnemonic
// Method: POST
// Args: JSON Body
func setMnemonic(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		var req SetMnemonicRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			resp := NewHTTPErrorResponse(http.StatusBadRequest, err.Error())
			writeHTTPResponse(w, resp)
			return
		}
		defer r.Body.Close()

		if err := bip39.ValidateMnemonic(req.Mnemonic); err != nil {
			resp := NewHTTPErrorResponse(http.StatusUnprocessableEntity, "seed is not a valid bip39 seed")
			writeHTTPResponse(w, resp)
			return
		}

		// for integration tests
		if autoPressEmulatorButtons {
			err := gateway.SetAutoPressButton(true, skyWallet.ButtonRight)
			if err != nil {
				logger.Error("setMnemonic failed: %s", err.Error())
				resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
				writeHTTPResponse(w, resp)
				return
			}
		}

		var msg wire.Message
		var err error
		retCH := make(chan int)
		errCH := make(chan int)
		ctx := r.Context()

		go func() {
			msg, err = gateway.SetMnemonic(req.Mnemonic)
			if err != nil {
				errCH <- 1
				return
			}
			retCH <- 1
		}()

		select {
		case <-retCH:
			HandleFirmwareResponseMessages(w, msg)
		case <-errCH:
			logger.Errorf("setMnemonic failed: %s", err.Error())
			resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
			writeHTTPResponse(w, resp)
		case <-ctx.Done():
			disConnErr := gateway.Disconnect()
			if disConnErr != nil {
				resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
				writeHTTPResponse(w, resp)
			} else {
				resp := NewHTTPErrorResponse(499, "Client Closed Request")
				writeHTTPResponse(w, resp)
			}
		}
	}
}
