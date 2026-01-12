package api

import (
	"crypto/sha256"
	"io/ioutil"
	"net/http"
)

const (
	// maxUploadSize is max firmware file size
	maxUploadSize = 1024 * 1024 // 1 MB
)

// URI: /api/v1/firmware_update
// Method: PUT
// Args:
//  file: firmware file
func firmwareUpdate(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			resp := NewHTTPErrorResponse(http.StatusMethodNotAllowed, "")
			writeHTTPResponse(w, resp)
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
		if err := r.ParseMultipartForm(maxUploadSize); err != nil {
			resp := NewHTTPErrorResponse(http.StatusBadRequest, err.Error())
			writeHTTPResponse(w, resp)
			return
		}

		file, _, err := r.FormFile("file")
		if err != nil {
			resp := NewHTTPErrorResponse(http.StatusUnprocessableEntity, err.Error())
			writeHTTPResponse(w, resp)
			return
		}
		defer file.Close()

		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			resp := NewHTTPErrorResponse(http.StatusUnprocessableEntity, err.Error())
			writeHTTPResponse(w, resp)
			return
		}

		retCH := make(chan int)
		errCH := make(chan int)
		ctx := r.Context()

		go func() {
			err = gateway.FirmwareUpload(fileBytes, sha256.Sum256(fileBytes[0x100:]))
			if err != nil {
				errCH <- 1
				return
			}
			retCH <- 1
		}()

		select {
		case <-retCH:
			writeHTTPResponse(w, HTTPResponse{})
		case <-errCH:
			logger.Errorf("firmwareUpdate failed: %s", err.Error())
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
