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

	if params.Data == nil {
		wh.Error400(w, "empty data")
		return
	}

	err = gateway.SaveData(params.Data, params.Update)
	if err != nil {
		resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
		writeHTTPResponse(w, resp)
		return
	}

	wh.SendJSONOr500(logger, w, "success")
}

// Get data from a file on disk
// URI: /api/v2/data
// Method: GET
// Args:
//     keys: comma separated list of keys to retrieve [required]
func dataGetHandler(gateway Gatewayer, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		wh.Error405(w)
		return
	}

	keys := r.FormValue("keys")
	if keys == "" {
		wh.Error400(w, "missing keys")
		return
	}

	keyArr := strings.Split(keys, ",")

	data, err := gateway.GetData(keyArr)
	if err != nil {
		switch {
		case err.Error() == "empty file":
			wh.Error400(w, err.Error())
		default:
			resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
			writeHTTPResponse(w, resp)
		}
		return
	}

	wh.SendJSONOr500(logger, w, data)
}

// Delete data from a file on disk
// URI: /api/v2/data
// Method: Delete
// Args:
//     keys: list of keys to retrieve [required]
func dataDeleteHandler(gateway Gatewayer, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		wh.Error405(w)
		return
	}

	keys := r.FormValue("keys")
	if keys == "" {
		wh.Error400(w, "missing keys")
		return
	}

	keyArr := strings.Split(keys, ",")

	err := gateway.DeleteData(keyArr)
	if err != nil {
		switch {
		case err.Error() == "empty file":
			wh.Error400(w, err.Error())
		default:
			resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
			writeHTTPResponse(w, resp)
		}
		return
	}

	wh.SendJSONOr500(logger, w, "success")
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
