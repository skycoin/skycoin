package api

import (
	"encoding/json"
	"net/http"

	wh "github.com/skycoin/skycoin/src/util/http"
)

// saveDataRequest
type saveDataRequest struct {
	Filename string                 `json:"filename"`
	Data     map[string]interface{} `json:"data"`
	Update   bool                   `json:"update"`
}

// getDataRequest
type getDataRequest struct {
	Filename string   `json:"filename"`
	Keys     []string `json:"keys"`
}

// deleteDataRequest
type deleteDataRequest struct {
	Filename string   `json:"filename"`
	Keys     []string `json:"keys"`
}

// Save arbitrary data to disk
// URI: /api/v2/data/save
// Method: POST
// Args:
//     filename: filename [required]
//     data: arbitrary data to save [required]
//     update: update existing values [optional]
func dataSaveHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		err = gateway.SaveData(params.Filename, params.Data, params.Update)
		if err != nil {
			wh.Error400(w, err.Error())
			return
		}

		wh.SendJSONOr500(logger, w, "success")
	}
}

// Get data from a file on disk
// URI: /api/v2/data/get
// Method: POST
// Args:
//     filename: filename [required]
//     keys: list of keys to retrieve [required]
func dataGetHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		var params getDataRequest
		err := decoder.Decode(&params)
		if err != nil {
			logger.WithError(err).Error("invalid get data request")
			wh.Error400(w, err.Error())
			return
		}
		defer r.Body.Close()

		data, err := gateway.GetData(params.Filename, params.Keys)
		if err != nil {
			wh.Error400(w, err.Error())
			return
		}

		wh.SendJSONOr500(logger, w, data)
	}
}

// Delete data from a file on disk
// URI: /api/v2/data/delete
// Method: POST
// Args:
//     filename: filename [required]
//     keys: list of keys to retrieve [required]
func dataDeleteHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		var params deleteDataRequest
		err := decoder.Decode(&params)
		if err != nil {
			logger.WithError(err).Error("invalid get data request")
			wh.Error400(w, err.Error())
			return
		}
		defer r.Body.Close()

		err = gateway.DeleteData(params.Filename, params.Keys)
		if err != nil {
			wh.Error400(w, err.Error())
			return
		}

		wh.SendJSONOr500(logger, w, "success")
	}
}
