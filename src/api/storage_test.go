package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/SkycoinProject/skycoin/src/kvstorage"
)

func TestGetAllStorageValuesHandler(t *testing.T) {
	tt := []struct {
		name                      string
		method                    string
		contentType               string
		query                     string
		status                    int
		storageType               kvstorage.Type
		getAllStorageValuesResult map[string]string
		getAllStorageValuesErr    error
		httpResponse              HTTPResponse
		csrfDisabled              bool
	}{
		{
			name:        "405",
			method:      http.MethodPut,
			contentType: ContentTypeForm,
			query: url.Values{
				"type": []string{string(kvstorage.TypeTxIDNotes)},
			}.Encode(),
			status:                    http.StatusMethodNotAllowed,
			storageType:               kvstorage.TypeTxIDNotes,
			getAllStorageValuesResult: make(map[string]string),
			getAllStorageValuesErr:    nil,
			httpResponse:              NewHTTPErrorResponse(http.StatusMethodNotAllowed, ""),
		},
		{
			name:        "403",
			method:      http.MethodGet,
			contentType: ContentTypeForm,
			query: url.Values{
				"type": []string{string(kvstorage.TypeTxIDNotes)},
			}.Encode(),
			status:                    http.StatusForbidden,
			storageType:               kvstorage.TypeTxIDNotes,
			getAllStorageValuesResult: nil,
			getAllStorageValuesErr:    kvstorage.ErrStorageAPIDisabled,
			httpResponse:              NewHTTPErrorResponse(http.StatusForbidden, ""),
		},
		{
			name:                      "400 - missing type",
			method:                    http.MethodGet,
			contentType:               ContentTypeForm,
			status:                    http.StatusBadRequest,
			storageType:               "",
			getAllStorageValuesResult: nil,
			getAllStorageValuesErr:    kvstorage.ErrUnknownKVStorageType,
			httpResponse:              NewHTTPErrorResponse(http.StatusBadRequest, "type is required"),
		},
		{
			name:        "400 - unknown type",
			method:      http.MethodGet,
			contentType: ContentTypeForm,
			query: url.Values{
				"type": []string{"unknown"},
			}.Encode(),
			status:                    http.StatusBadRequest,
			storageType:               "unknown",
			getAllStorageValuesResult: nil,
			getAllStorageValuesErr:    kvstorage.ErrUnknownKVStorageType,
			httpResponse:              NewHTTPErrorResponse(http.StatusBadRequest, "unknown storage"),
		},
		{
			name:        "404 - storage not loaded",
			method:      http.MethodGet,
			contentType: ContentTypeForm,
			query: url.Values{
				"type": []string{string(kvstorage.TypeTxIDNotes)},
			}.Encode(),
			status:                    http.StatusNotFound,
			storageType:               kvstorage.TypeTxIDNotes,
			getAllStorageValuesResult: nil,
			getAllStorageValuesErr:    kvstorage.ErrNoSuchStorage,
			httpResponse:              NewHTTPErrorResponse(http.StatusNotFound, "storage is not loaded"),
		},
		{
			name:        "200",
			method:      http.MethodGet,
			contentType: ContentTypeForm,
			query: url.Values{
				"type": []string{string(kvstorage.TypeTxIDNotes)},
			}.Encode(),
			status:      http.StatusOK,
			storageType: kvstorage.TypeTxIDNotes,
			getAllStorageValuesResult: map[string]string{
				"test1": "some value",
				"test2": "{\"key\":\"val\",\"key2\":2}",
			},
			getAllStorageValuesErr: nil,
			httpResponse: HTTPResponse{
				Data: map[string]string{
					"test1": "some value",
					"test2": "{\"key\":\"val\",\"key2\":2}",
				},
			},
		},
		{
			name:        "200 - csrf disabled",
			method:      http.MethodGet,
			contentType: ContentTypeForm,
			query: url.Values{
				"type": []string{string(kvstorage.TypeTxIDNotes)},
			}.Encode(),
			status:      http.StatusOK,
			storageType: kvstorage.TypeTxIDNotes,
			getAllStorageValuesResult: map[string]string{
				"test1": "some value",
				"test2": "{\"key\":\"val\",\"key2\":2}",
			},
			getAllStorageValuesErr: nil,
			httpResponse: HTTPResponse{
				Data: map[string]string{
					"test1": "some value",
					"test2": "{\"key\":\"val\",\"key2\":2}",
				},
			},
			csrfDisabled: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			gateway := &MockGatewayer{}
			gateway.On("GetAllStorageValues", tc.storageType).Return(tc.getAllStorageValuesResult,
				tc.getAllStorageValuesErr)

			endpoint := "/api/v2/data"

			if tc.query != "" {
				endpoint += "?" + tc.query
			}

			req, err := http.NewRequest(tc.method, endpoint, strings.NewReader(""))
			require.NoError(t, err)

			req.Header.Set("Content-Type", tc.contentType)

			if tc.csrfDisabled {
				setCSRFParameters(t, tokenInvalid, req)
			} else {
				setCSRFParameters(t, tokenValid, req)
			}

			rr := httptest.NewRecorder()

			cfg := defaultMuxConfig()
			cfg.disableCSRF = false

			handler := newServerMux(cfg, gateway)
			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "got `%v` want `%v`", status, tc.status)

			var rsp ReceivedHTTPResponse
			err = json.Unmarshal(rr.Body.Bytes(), &rsp)
			require.NoError(t, err)

			require.Equal(t, tc.httpResponse.Error, rsp.Error)

			if rsp.Data == nil {
				require.Nil(t, tc.httpResponse.Data)
			} else {
				require.NotNil(t, tc.httpResponse.Data)

				var getAllStorageValuesRsp map[string]string
				err := json.Unmarshal(rsp.Data, &getAllStorageValuesRsp)
				require.NoError(t, err)

				require.Equal(t, tc.httpResponse.Data, getAllStorageValuesRsp)
			}
		})
	}
}

func TestGetStorageValueHandler(t *testing.T) {
	tt := []struct {
		name                  string
		method                string
		contentType           string
		query                 string
		status                int
		storageType           kvstorage.Type
		key                   string
		getStorageValueResult string
		getStorageValueErr    error
		httpResponse          HTTPResponse
		csrfDisabled          bool
	}{
		{
			name:        "405",
			method:      http.MethodPut,
			contentType: ContentTypeForm,
			query: url.Values{
				"type": []string{string(kvstorage.TypeTxIDNotes)},
				"key":  []string{"test1"},
			}.Encode(),
			status:                http.StatusMethodNotAllowed,
			storageType:           kvstorage.TypeTxIDNotes,
			key:                   "test1",
			getStorageValueResult: "some value",
			getStorageValueErr:    nil,
			httpResponse:          NewHTTPErrorResponse(http.StatusMethodNotAllowed, ""),
		},
		{
			name:        "403",
			method:      http.MethodGet,
			contentType: ContentTypeForm,
			query: url.Values{
				"type": []string{string(kvstorage.TypeTxIDNotes)},
				"key":  []string{"test1"},
			}.Encode(),
			status:                http.StatusForbidden,
			storageType:           kvstorage.TypeTxIDNotes,
			key:                   "test1",
			getStorageValueResult: "",
			getStorageValueErr:    kvstorage.ErrStorageAPIDisabled,
			httpResponse:          NewHTTPErrorResponse(http.StatusForbidden, ""),
		},
		{
			name:        "400 - missing type",
			method:      http.MethodGet,
			contentType: ContentTypeForm,
			query: url.Values{
				"key": []string{"test1"},
			}.Encode(),
			status:                http.StatusBadRequest,
			storageType:           "",
			key:                   "test1",
			getStorageValueResult: "",
			getStorageValueErr:    kvstorage.ErrUnknownKVStorageType,
			httpResponse:          NewHTTPErrorResponse(http.StatusBadRequest, "type is required"),
		},
		{
			name:        "400 - unknown type",
			method:      http.MethodGet,
			contentType: ContentTypeForm,
			query: url.Values{
				"type": []string{"unknown"},
				"key":  []string{"test1"},
			}.Encode(),
			status:                http.StatusBadRequest,
			storageType:           "unknown",
			key:                   "test1",
			getStorageValueResult: "",
			getStorageValueErr:    kvstorage.ErrUnknownKVStorageType,
			httpResponse:          NewHTTPErrorResponse(http.StatusBadRequest, "unknown storage"),
		},
		{
			name:        "404 - storage not loaded",
			method:      http.MethodGet,
			contentType: ContentTypeForm,
			query: url.Values{
				"type": []string{string(kvstorage.TypeTxIDNotes)},
				"key":  []string{"test1"},
			}.Encode(),
			status:                http.StatusNotFound,
			storageType:           kvstorage.TypeTxIDNotes,
			key:                   "test1",
			getStorageValueResult: "",
			getStorageValueErr:    kvstorage.ErrNoSuchStorage,
			httpResponse:          NewHTTPErrorResponse(http.StatusNotFound, "storage is not loaded"),
		},
		{
			name:        "400 - not found",
			method:      http.MethodGet,
			contentType: ContentTypeForm,
			query: url.Values{
				"type": []string{string(kvstorage.TypeTxIDNotes)},
				"key":  []string{"test1"},
			}.Encode(),
			status:                http.StatusNotFound,
			storageType:           kvstorage.TypeTxIDNotes,
			key:                   "test1",
			getStorageValueResult: "",
			getStorageValueErr:    kvstorage.ErrNoSuchKey,
			httpResponse:          NewHTTPErrorResponse(http.StatusNotFound, ""),
		},
		{
			name:        "200",
			method:      http.MethodGet,
			contentType: ContentTypeForm,
			query: url.Values{
				"type": []string{string(kvstorage.TypeTxIDNotes)},
				"key":  []string{"test1"},
			}.Encode(),
			status:                http.StatusOK,
			storageType:           kvstorage.TypeTxIDNotes,
			key:                   "test1",
			getStorageValueResult: "some value",
			getStorageValueErr:    nil,
			httpResponse: HTTPResponse{
				Data: "some value",
			},
		},
		{
			name:        "200 - csrf disabled",
			method:      http.MethodGet,
			contentType: ContentTypeForm,
			query: url.Values{
				"type": []string{string(kvstorage.TypeTxIDNotes)},
				"key":  []string{"test1"},
			}.Encode(),
			status:                http.StatusOK,
			storageType:           kvstorage.TypeTxIDNotes,
			key:                   "test1",
			getStorageValueResult: "some value",
			getStorageValueErr:    nil,
			httpResponse: HTTPResponse{
				Data: "some value",
			},
			csrfDisabled: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			gateway := &MockGatewayer{}
			gateway.On("GetStorageValue", tc.storageType, tc.key).Return(tc.getStorageValueResult,
				tc.getStorageValueErr)

			endpoint := "/api/v2/data"

			if tc.query != "" {
				endpoint += "?" + tc.query
			}

			req, err := http.NewRequest(tc.method, endpoint, strings.NewReader(""))
			require.NoError(t, err)

			req.Header.Set("Content-Type", tc.contentType)

			if tc.csrfDisabled {
				setCSRFParameters(t, tokenInvalid, req)
			} else {
				setCSRFParameters(t, tokenValid, req)
			}

			rr := httptest.NewRecorder()

			cfg := defaultMuxConfig()
			cfg.disableCSRF = false

			handler := newServerMux(cfg, gateway)
			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "got `%v` want `%v`", status, tc.status)

			var rsp ReceivedHTTPResponse
			err = json.Unmarshal(rr.Body.Bytes(), &rsp)
			require.NoError(t, err)

			require.Equal(t, tc.httpResponse.Error, rsp.Error)

			if rsp.Data == nil {
				require.Nil(t, tc.httpResponse.Data)
			} else {
				require.NotNil(t, tc.httpResponse.Data)

				var getStorageValueRsp string
				err := json.Unmarshal(rsp.Data, &getStorageValueRsp)
				require.NoError(t, err)

				require.Equal(t, tc.httpResponse.Data, getStorageValueRsp)
			}
		})
	}
}

func TestAddStorageValueHandler(t *testing.T) {
	tt := []struct {
		name               string
		method             string
		contentType        string
		httpBody           string
		status             int
		storageType        kvstorage.Type
		key                string
		val                string
		addStorageValueErr error
		httpResponse       HTTPResponse
		csrfDisabled       bool
	}{
		{
			name:        "405",
			method:      http.MethodPut,
			contentType: ContentTypeJSON,
			httpBody: toJSON(t, StorageRequest{
				StorageType: kvstorage.TypeTxIDNotes,
				Key:         "test",
				Val:         "qwe",
			}),
			status:             http.StatusMethodNotAllowed,
			storageType:        kvstorage.TypeTxIDNotes,
			key:                "test",
			val:                "qwe",
			addStorageValueErr: nil,
			httpResponse:       NewHTTPErrorResponse(http.StatusMethodNotAllowed, ""),
		},
		{
			name:        "403",
			method:      http.MethodPost,
			contentType: ContentTypeJSON,
			httpBody: toJSON(t, StorageRequest{
				StorageType: kvstorage.TypeTxIDNotes,
				Key:         "test",
				Val:         "qwe",
			}),
			status:             http.StatusForbidden,
			storageType:        kvstorage.TypeTxIDNotes,
			key:                "test",
			val:                "qwe",
			addStorageValueErr: kvstorage.ErrStorageAPIDisabled,
			httpResponse:       NewHTTPErrorResponse(http.StatusForbidden, ""),
		},
		{
			name:        "415",
			method:      http.MethodPost,
			contentType: ContentTypeForm,
			httpBody: toJSON(t, StorageRequest{
				StorageType: kvstorage.TypeTxIDNotes,
				Key:         "test",
				Val:         "qwe",
			}),
			status:             http.StatusUnsupportedMediaType,
			storageType:        kvstorage.TypeTxIDNotes,
			key:                "test",
			val:                "qwe",
			addStorageValueErr: nil,
			httpResponse:       NewHTTPErrorResponse(http.StatusUnsupportedMediaType, ""),
		},
		{
			name:               "400 - EOF",
			method:             http.MethodPost,
			contentType:        ContentTypeJSON,
			httpBody:           "",
			status:             http.StatusBadRequest,
			storageType:        "",
			key:                "",
			val:                "",
			addStorageValueErr: kvstorage.ErrUnknownKVStorageType,
			httpResponse:       NewHTTPErrorResponse(http.StatusBadRequest, "EOF"),
		},
		{
			name:        "400 - missing type",
			method:      http.MethodPost,
			contentType: ContentTypeJSON,
			httpBody: toJSON(t, StorageRequest{
				Key: "test",
				Val: "qwe",
			}),
			status:             http.StatusBadRequest,
			storageType:        "",
			key:                "test",
			val:                "qwe",
			addStorageValueErr: kvstorage.ErrUnknownKVStorageType,
			httpResponse:       NewHTTPErrorResponse(http.StatusBadRequest, "type is required"),
		},
		{
			name:        "400 - unknown type",
			method:      http.MethodPost,
			contentType: ContentTypeJSON,
			httpBody: toJSON(t, StorageRequest{
				StorageType: "unknown",
				Key:         "test",
				Val:         "qwe",
			}),
			status:             http.StatusBadRequest,
			storageType:        "unknown",
			key:                "test",
			val:                "qwe",
			addStorageValueErr: kvstorage.ErrUnknownKVStorageType,
			httpResponse:       NewHTTPErrorResponse(http.StatusBadRequest, "unknown storage"),
		},
		{
			name:        "404 - storage not loaded",
			method:      http.MethodPost,
			contentType: ContentTypeJSON,
			httpBody: toJSON(t, StorageRequest{
				StorageType: kvstorage.TypeTxIDNotes,
				Key:         "test",
				Val:         "qwe",
			}),
			status:             http.StatusNotFound,
			storageType:        kvstorage.TypeTxIDNotes,
			key:                "test",
			val:                "qwe",
			addStorageValueErr: kvstorage.ErrNoSuchStorage,
			httpResponse:       NewHTTPErrorResponse(http.StatusNotFound, "storage is not loaded"),
		},
		{
			name:        "400 - missing key",
			method:      http.MethodPost,
			contentType: ContentTypeJSON,
			httpBody: toJSON(t, StorageRequest{
				StorageType: kvstorage.TypeTxIDNotes,
				Val:         "qwe",
			}),
			status:             http.StatusBadRequest,
			storageType:        kvstorage.TypeTxIDNotes,
			val:                "qwe",
			addStorageValueErr: nil,
			httpResponse:       NewHTTPErrorResponse(http.StatusBadRequest, "key is required"),
		},
		{
			name:        "200",
			method:      http.MethodPost,
			contentType: ContentTypeJSON,
			httpBody: toJSON(t, StorageRequest{
				StorageType: kvstorage.TypeTxIDNotes,
				Key:         "test",
				Val:         "qwe",
			}),
			status:             http.StatusOK,
			storageType:        kvstorage.TypeTxIDNotes,
			key:                "test",
			val:                "qwe",
			addStorageValueErr: nil,
			httpResponse:       HTTPResponse{},
		},
		{
			name:        "403 - csrf disabled",
			method:      http.MethodPost,
			contentType: ContentTypeJSON,
			httpBody: toJSON(t, StorageRequest{
				StorageType: kvstorage.TypeTxIDNotes,
				Key:         "test",
				Val:         "qwe",
			}),
			status:             http.StatusForbidden,
			storageType:        kvstorage.TypeTxIDNotes,
			key:                "test",
			val:                "qwe",
			addStorageValueErr: nil,
			httpResponse:       NewHTTPErrorResponse(http.StatusForbidden, "invalid CSRF token"),
			csrfDisabled:       true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			gateway := &MockGatewayer{}
			gateway.On("AddStorageValue", tc.storageType, tc.key, tc.val).Return(tc.addStorageValueErr)

			endpoint := "/api/v2/data"

			req, err := http.NewRequest(tc.method, endpoint, strings.NewReader(tc.httpBody))
			require.NoError(t, err)

			req.Header.Set("Content-Type", tc.contentType)

			if tc.csrfDisabled {
				setCSRFParameters(t, tokenInvalid, req)
			} else {
				setCSRFParameters(t, tokenValid, req)
			}

			rr := httptest.NewRecorder()

			cfg := defaultMuxConfig()
			cfg.disableCSRF = false

			handler := newServerMux(cfg, gateway)
			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "got `%v` want `%v`", status, tc.status)

			var rsp ReceivedHTTPResponse
			err = json.Unmarshal(rr.Body.Bytes(), &rsp)
			require.NoError(t, err)

			require.Equal(t, tc.httpResponse.Error, rsp.Error)

			if rsp.Data == nil {
				require.Nil(t, tc.httpResponse.Data)
			} else {
				require.NotNil(t, tc.httpResponse.Data)
			}
		})
	}
}

func TestRemoveStorageValueHandler(t *testing.T) {
	tt := []struct {
		name                  string
		method                string
		query                 string
		status                int
		storageType           kvstorage.Type
		key                   string
		removeStorageValueErr error
		httpResponse          HTTPResponse
		csrfDisabled          bool
	}{
		{
			name:   "405",
			method: http.MethodPut,
			query: url.Values{
				"type": []string{string(kvstorage.TypeTxIDNotes)},
				"key":  []string{"test"},
			}.Encode(),
			status:                http.StatusMethodNotAllowed,
			storageType:           kvstorage.TypeTxIDNotes,
			key:                   "test",
			removeStorageValueErr: nil,
			httpResponse:          NewHTTPErrorResponse(http.StatusMethodNotAllowed, ""),
		},
		{
			name:   "403",
			method: http.MethodDelete,
			query: url.Values{
				"type": []string{string(kvstorage.TypeTxIDNotes)},
				"key":  []string{"test"},
			}.Encode(),
			status:                http.StatusForbidden,
			storageType:           kvstorage.TypeTxIDNotes,
			key:                   "test",
			removeStorageValueErr: kvstorage.ErrStorageAPIDisabled,
			httpResponse:          NewHTTPErrorResponse(http.StatusForbidden, ""),
		},
		{
			name:   "400 - missing type",
			method: http.MethodDelete,
			query: url.Values{
				"key": []string{"test"},
			}.Encode(),
			status:                http.StatusBadRequest,
			storageType:           kvstorage.TypeTxIDNotes,
			key:                   "test",
			removeStorageValueErr: kvstorage.ErrUnknownKVStorageType,
			httpResponse:          NewHTTPErrorResponse(http.StatusBadRequest, "type is required"),
		},
		{
			name:   "400 - unknown type",
			method: http.MethodDelete,
			query: url.Values{
				"type": []string{"unknown"},
				"key":  []string{"test"},
			}.Encode(),
			status:                http.StatusBadRequest,
			storageType:           "unknown",
			key:                   "test",
			removeStorageValueErr: kvstorage.ErrUnknownKVStorageType,
			httpResponse:          NewHTTPErrorResponse(http.StatusBadRequest, "unknown storage"),
		},
		{
			name:   "404 - storage not loaded",
			method: http.MethodDelete,
			query: url.Values{
				"type": []string{string(kvstorage.TypeTxIDNotes)},
				"key":  []string{"test"},
			}.Encode(),
			status:                http.StatusNotFound,
			storageType:           kvstorage.TypeTxIDNotes,
			key:                   "test",
			removeStorageValueErr: kvstorage.ErrNoSuchStorage,
			httpResponse:          NewHTTPErrorResponse(http.StatusNotFound, "storage is not loaded"),
		},
		{
			name:   "400 - missing key",
			method: http.MethodDelete,
			query: url.Values{
				"type": []string{string(kvstorage.TypeTxIDNotes)},
			}.Encode(),
			status:                http.StatusBadRequest,
			storageType:           kvstorage.TypeTxIDNotes,
			removeStorageValueErr: nil,
			httpResponse:          NewHTTPErrorResponse(http.StatusBadRequest, "key is required"),
		},
		{
			name:   "404 - not found",
			method: http.MethodDelete,
			query: url.Values{
				"type": []string{string(kvstorage.TypeTxIDNotes)},
				"key":  []string{"test"},
			}.Encode(),
			status:                http.StatusNotFound,
			storageType:           kvstorage.TypeTxIDNotes,
			key:                   "test",
			removeStorageValueErr: kvstorage.ErrNoSuchKey,
			httpResponse:          NewHTTPErrorResponse(http.StatusNotFound, ""),
		},
		{
			name:   "200",
			method: http.MethodDelete,
			query: url.Values{
				"type": []string{string(kvstorage.TypeTxIDNotes)},
				"key":  []string{"test"},
			}.Encode(),
			status:                http.StatusOK,
			storageType:           kvstorage.TypeTxIDNotes,
			key:                   "test",
			removeStorageValueErr: nil,
			httpResponse:          HTTPResponse{},
		},
		{
			name:   "403 - csrf disabled",
			method: http.MethodDelete,
			query: url.Values{
				"type": []string{string(kvstorage.TypeTxIDNotes)},
				"key":  []string{"test"},
			}.Encode(),
			status:                http.StatusForbidden,
			storageType:           kvstorage.TypeTxIDNotes,
			key:                   "test",
			removeStorageValueErr: nil,
			httpResponse:          NewHTTPErrorResponse(http.StatusForbidden, "invalid CSRF token"),
			csrfDisabled:          true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			gateway := &MockGatewayer{}
			gateway.On("RemoveStorageValue", tc.storageType, tc.key).Return(tc.removeStorageValueErr)

			endpoint := "/api/v2/data"

			if tc.query != "" {
				endpoint += "?" + tc.query
			}

			req, err := http.NewRequest(tc.method, endpoint, strings.NewReader(""))
			require.NoError(t, err)

			if tc.csrfDisabled {
				setCSRFParameters(t, tokenInvalid, req)
			} else {
				setCSRFParameters(t, tokenValid, req)
			}

			rr := httptest.NewRecorder()

			cfg := defaultMuxConfig()
			cfg.disableCSRF = false

			handler := newServerMux(cfg, gateway)
			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "got `%v` want `%v`", status, tc.status)

			var rsp ReceivedHTTPResponse
			err = json.Unmarshal(rr.Body.Bytes(), &rsp)
			require.NoError(t, err)

			require.Equal(t, tc.httpResponse.Error, rsp.Error)

			if rsp.Data == nil {
				require.Nil(t, tc.httpResponse.Data)
			} else {
				require.NotNil(t, tc.httpResponse.Data)
			}
		})
	}
}
