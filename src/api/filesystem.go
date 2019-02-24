package api

import (
	"encoding/json"
	"net/http"
	"strings"

	wh "github.com/skycoin/skycoin/src/util/http"
)

// saveDataRequest
type saveDataRequest struct {
	Data   map[string]interface{} `json:"data"`
	Update bool                   `json:"update"`
}

// Save arbitrary data to disk
// URI: /api/v2/data
// Method: POST
// Args:
//     data: arbitrary data to save [required]
//     update: update existing values [optional]
func dataSaveHandler(gateway Gatewayer, w http.ResponseWriter, r *http.Request) {
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

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var params saveDataRequest
	err := decoder.Decode(&params)
	if err != nil {
		resp := NewHTTPErrorResponse(http.StatusBadRequest, err.Error())
		writeHTTPResponse(w, resp)
		return
	}
	defer r.Body.Close()

	if params.Data == nil {
		resp := NewHTTPErrorResponse(http.StatusBadRequest, "empty data")
		writeHTTPResponse(w, resp)
		return
	}

	err = gateway.SaveData(params.Data, params.Update)
	if err != nil {
		resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
		writeHTTPResponse(w, resp)
		return
	}

	writeHTTPResponse(w, HTTPResponse{Data: struct{}{}})
}

// Get data from a file on disk
// URI: /api/v2/data
// Method: GET
// Args:
//     keys: comma separated list of keys to retrieve [required]
func dataGetHandler(gateway Gatewayer, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		resp := NewHTTPErrorResponse(http.StatusMethodNotAllowed, "")
		writeHTTPResponse(w, resp)
		return
	}

	keys := r.FormValue("keys")
	if keys == "" {
		resp := NewHTTPErrorResponse(http.StatusBadRequest, "missing keys")
		writeHTTPResponse(w, resp)
		return
	}

	keyArr := strings.Split(keys, ",")

	data, err := gateway.GetData(keyArr)
	if err != nil {
		switch {
		case err.Error() == "empty file":
			resp := NewHTTPErrorResponse(http.StatusBadRequest, err.Error())
			writeHTTPResponse(w, resp)
		default:
			resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
			writeHTTPResponse(w, resp)
		}
		return
	}

	writeHTTPResponse(w, HTTPResponse{Data: data})
}

// Delete data from a file on disk
// URI: /api/v2/data
// Method: Delete
// Args:
//     keys: list of keys to retrieve [required]
func dataDeleteHandler(gateway Gatewayer, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		resp := NewHTTPErrorResponse(http.StatusMethodNotAllowed, "")
		writeHTTPResponse(w, resp)
		return
	}

	keys := r.FormValue("keys")
	if keys == "" {
		resp := NewHTTPErrorResponse(http.StatusBadRequest, "missing keys")
		writeHTTPResponse(w, resp)
		return
	}

	keyArr := strings.Split(keys, ",")

	err := gateway.DeleteData(keyArr)
	if err != nil {
		switch {
		case err.Error() == "empty file":
			resp := NewHTTPErrorResponse(http.StatusBadRequest, err.Error())
			writeHTTPResponse(w, resp)
		default:
			resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
			writeHTTPResponse(w, resp)
		}
		return
	}

	writeHTTPResponse(w, HTTPResponse{Data: struct{}{}})
}

func dataHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			dataGetHandler(gateway, w, r)
		case http.MethodPost:
			dataSaveHandler(gateway, w, r)
		case http.MethodDelete:
			dataDeleteHandler(gateway, w, r)
		default:
			wh.Error405(w)
			return
		}
	}
}
