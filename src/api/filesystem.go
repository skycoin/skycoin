package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	wh "github.com/skycoin/skycoin/src/util/http"
)

// saveDataRequest
type saveDataRequest struct {
	Filename string                 `json:"filename"`
	Data     map[string]interface{} `json:"data"`
	Update   bool                   `json:"update"`
}

// Save arbitrary data to disk
// URI: /api/v2/data/save
// Method: POST, PATCH
// Args:
//     filename: filename [required]
//     data: arbitrary data to save [required]
//     update: update existing values [optional]
func dataSaveHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost && r.Method != http.MethodPatch {
			wh.Error405(w)
			return
		}

		if r.Header.Get("Content-Type") != ContentTypeJSON {
			wh.Error415(w)
			return
		}

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()

		var params saveDataRequest
		err := decoder.Decode(&params)
		if err != nil {
			logger.WithError(err).Error("invalid save data request")
			wh.Error400(w, err.Error())
			return
		}
		defer r.Body.Close()

		if params.Filename == "" {
			wh.Error400(w, "missing filename")
			return
		}

		if params.Data == nil {
			wh.Error400(w, "empty data")
			return
		}

		err = gateway.SaveData(params.Filename, params.Data, params.Update)
		if err != nil {
			switch {
			case os.IsNotExist(err):
				wh.Error404(w, fmt.Sprintf("file %s does not exist", params.Filename))
			case os.IsPermission(err):
				wh.Error403(w, fmt.Sprintf("cannot access %s - permission denied", params.Filename))
			default:
				resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
				writeHTTPResponse(w, resp)
			}
			return
		}

		wh.SendJSONOr500(logger, w, "success")
	}
}

// Get data from a file on disk
// URI: /api/v2/data/get
// Method: GET
// Args:
//     filename: filename [required]
//     keys: comma separated list of keys to retrieve [required]
func dataGetHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		filename := r.FormValue("filename")
		if filename == "" {
			wh.Error400(w, "missing filename")
			return
		}

		keys := r.FormValue("keys")
		if keys == "" {
			wh.Error400(w, "missing keys")
			return
		}

		keyArr := strings.Split(keys, ",")

		data, err := gateway.GetData(filename, keyArr)
		if err != nil {
			switch {
			case err.Error() == "empty file":
				wh.Error400(w, err.Error())
			case os.IsNotExist(err):
				wh.Error404(w, fmt.Sprintf("file %s does not exist", filename))
			case os.IsPermission(err):
				wh.Error403(w, fmt.Sprintf("cannot access %s: permission denied", filename))
			default:
				resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
				writeHTTPResponse(w, resp)
			}
			return
		}

		wh.SendJSONOr500(logger, w, data)
	}
}

// Delete data from a file on disk
// URI: /api/v2/data/delete
// Method: Delete
// Args:
//     filename: filename [required]
//     keys: list of keys to retrieve [required]
func dataDeleteHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			wh.Error405(w)
			return
		}

		filename := r.FormValue("filename")
		if filename == "" {
			wh.Error400(w, "missing filename")
			return
		}

		keys := r.FormValue("keys")
		if keys == "" {
			wh.Error400(w, "missing keys")
			return
		}

		keyArr := strings.Split(keys, ",")

		err := gateway.DeleteData(filename, keyArr)
		if err != nil {
			switch {
			case err.Error() == "empty file":
				wh.Error400(w, err.Error())
			case os.IsNotExist(err):
				wh.Error404(w, fmt.Sprintf("file %s does not exist", filename))
			case os.IsPermission(err):
				wh.Error403(w, fmt.Sprintf("cannot access %s: permission denied", filename))
			default:
				resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
				writeHTTPResponse(w, resp)
			}
			return
		}

		wh.SendJSONOr500(logger, w, "success")
	}
}
