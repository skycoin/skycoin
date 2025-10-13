package api

import (
	skyWallet "github.com/skycoin/hardware-wallet-go/src/skywallet"
	messages "github.com/skycoin/hardware-wallet-protob/go"

	"net/http"
)

// URI: /api/v1/cancel
// Method: PUT
func cancel(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			resp := NewHTTPErrorResponse(http.StatusMethodNotAllowed, "")
			writeHTTPResponse(w, resp)
			return
		}

		// for integration tests
		if autoPressEmulatorButtons {
			err := gateway.SetAutoPressButton(true, skyWallet.ButtonRight)
			if err != nil {
				logger.Error("cancel failed: %s", err.Error())
				resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
				writeHTTPResponse(w, resp)
				return
			}
		}

		msg, err := gateway.Cancel()
		if err != nil {
			logger.Errorf("cancel failed: %s", err.Error())
			resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
			writeHTTPResponse(w, resp)
			return
		}

		if msg.Kind == uint16(messages.MessageType_MessageType_Failure) {
			failureMsg, err := skyWallet.DecodeFailMsg(msg)
			if err != nil {
				logger.Errorf("cancel failed: %s", err.Error())
				resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
				writeHTTPResponse(w, resp)
				return
			}

			writeHTTPResponse(w, HTTPResponse{
				Data: []string{failureMsg},
			})
		} else {
			HandleFirmwareResponseMessages(w, msg)
		}
	}
}
