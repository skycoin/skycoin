package api

import (
	"encoding/json"
	"net/http"

	"github.com/SkycoinProject/skycoin/src/kvstorage"
)

// Dispatches /data endpoint.
// Method: GET, POST, DELETE
// URI: /api/v2/data
func storageHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getStorageValuesHandler(w, r, gateway)
		case http.MethodPost:
			addStorageValueHandler(w, r, gateway)
		case http.MethodDelete:
			removeStorageValueHandler(w, r, gateway)
		default:
			resp := NewHTTPErrorResponse(http.StatusMethodNotAllowed, "")
			writeHTTPResponse(w, resp)
		}
	}
}

// serves GET requests for /data enpdoint
func getStorageValuesHandler(w http.ResponseWriter, r *http.Request, gateway Gatewayer) {
	storageType := r.FormValue("type")
	if storageType == "" {
		resp := NewHTTPErrorResponse(http.StatusBadRequest, "type is required")
		writeHTTPResponse(w, resp)
		return
	}

	key := r.FormValue("key")

	if key == "" {
		getAllStorageValuesHandler(w, gateway, kvstorage.Type(storageType))
	} else {
		getStorageValueHandler(w, gateway, kvstorage.Type(storageType), key)
	}
}

// Returns all existing storage values of a given storage type.
// Args:
//     type: storage type to get values from
func getAllStorageValuesHandler(w http.ResponseWriter, gateway Gatewayer, storageType kvstorage.Type) {
	data, err := gateway.GetAllStorageValues(kvstorage.Type(storageType))
	if err != nil {
		var resp HTTPResponse
		switch err {
		case kvstorage.ErrStorageAPIDisabled:
			resp = NewHTTPErrorResponse(http.StatusForbidden, "")
		case kvstorage.ErrNoSuchStorage:
			resp = NewHTTPErrorResponse(http.StatusNotFound, "storage is not loaded")
		case kvstorage.ErrUnknownKVStorageType:
			resp = NewHTTPErrorResponse(http.StatusBadRequest, "unknown storage")
		default:
			resp = NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
		}
		writeHTTPResponse(w, resp)
		return
	}

	writeHTTPResponse(w, HTTPResponse{
		Data: data,
	})
}

// Returns value from storage of a given type by key.
// Args:
//     key: key for a value to be retrieved
func getStorageValueHandler(w http.ResponseWriter, gateway Gatewayer, storageType kvstorage.Type, key string) {
	val, err := gateway.GetStorageValue(storageType, key)
	if err != nil {
		var resp HTTPResponse
		switch err {
		case kvstorage.ErrStorageAPIDisabled:
			resp = NewHTTPErrorResponse(http.StatusForbidden, "")
		case kvstorage.ErrNoSuchStorage:
			resp = NewHTTPErrorResponse(http.StatusNotFound, "storage is not loaded")
		case kvstorage.ErrUnknownKVStorageType:
			resp = NewHTTPErrorResponse(http.StatusBadRequest, "unknown storage")
		case kvstorage.ErrNoSuchKey:
			resp = NewHTTPErrorResponse(http.StatusNotFound, "")
		default:
			resp = NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
		}
		writeHTTPResponse(w, resp)
		return
	}

	writeHTTPResponse(w, HTTPResponse{
		Data: val,
	})
}

// StorageRequest is the request data for POST /api/v2/data
type StorageRequest struct {
	StorageType kvstorage.Type `json:"type"`
	Key         string         `json:"key"`
	Val         string         `json:"val"`
}

// Adds the value to the storage of a given type
// Args:
//     type: storage type
//     key: key
//     val: value
func addStorageValueHandler(w http.ResponseWriter, r *http.Request, gateway Gatewayer) {
	var req StorageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		resp := NewHTTPErrorResponse(http.StatusBadRequest, err.Error())
		writeHTTPResponse(w, resp)
		return
	}

	if req.StorageType == "" {
		resp := NewHTTPErrorResponse(http.StatusBadRequest, "type is required")
		writeHTTPResponse(w, resp)
		return
	}

	if req.Key == "" {
		resp := NewHTTPErrorResponse(http.StatusBadRequest, "key is required")
		writeHTTPResponse(w, resp)
		return
	}

	if err := gateway.AddStorageValue(req.StorageType, req.Key, req.Val); err != nil {
		var resp HTTPResponse
		switch err {
		case kvstorage.ErrStorageAPIDisabled:
			resp = NewHTTPErrorResponse(http.StatusForbidden, "")
		case kvstorage.ErrNoSuchStorage:
			resp = NewHTTPErrorResponse(http.StatusNotFound, "storage is not loaded")
		case kvstorage.ErrUnknownKVStorageType:
			resp = NewHTTPErrorResponse(http.StatusBadRequest, "unknown storage")
		default:
			resp = NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
		}
		writeHTTPResponse(w, resp)
		return
	}

	writeHTTPResponse(w, HTTPResponse{})
}

// Removes the value by key from the storage of a given type
// Args:
//     type: storage type
//     key: key
func removeStorageValueHandler(w http.ResponseWriter, r *http.Request, gateway Gatewayer) {
	storageType := r.FormValue("type")
	if storageType == "" {
		resp := NewHTTPErrorResponse(http.StatusBadRequest, "type is required")
		writeHTTPResponse(w, resp)
		return
	}

	key := r.FormValue("key")
	if key == "" {
		resp := NewHTTPErrorResponse(http.StatusBadRequest, "key is required")
		writeHTTPResponse(w, resp)
		return
	}

	if err := gateway.RemoveStorageValue(kvstorage.Type(storageType), key); err != nil {
		var resp HTTPResponse
		switch err {
		case kvstorage.ErrStorageAPIDisabled:
			resp = NewHTTPErrorResponse(http.StatusForbidden, "")
		case kvstorage.ErrNoSuchStorage:
			resp = NewHTTPErrorResponse(http.StatusNotFound, "storage is not loaded")
		case kvstorage.ErrUnknownKVStorageType:
			resp = NewHTTPErrorResponse(http.StatusBadRequest, "unknown storage")
		case kvstorage.ErrNoSuchKey:
			resp = NewHTTPErrorResponse(http.StatusNotFound, "")
		default:
			resp = NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
		}
		writeHTTPResponse(w, resp)
		return
	}

	writeHTTPResponse(w, HTTPResponse{})
}
