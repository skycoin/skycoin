package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSaveData(t *testing.T) {
	type httpBody struct {
		Data   map[string]interface{}
		Update bool
	}

	tt := []struct {
		name               string
		method             string
		body               *httpBody
		status             int
		err                string
		contentType        string
		gatewaySaveDataErr error
		responseBody       string
	}{
		{
			name:        "415",
			method:      http.MethodPost,
			status:      http.StatusUnsupportedMediaType,
			contentType: ContentTypeForm,
			err:         "415 Unsupported Media Type",
		},
		{
			name:   "400 - empty data",
			method: http.MethodPost,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - empty data",
			body:   &httpBody{},
		},
		{
			name:   "200",
			method: http.MethodPost,
			body: &httpBody{
				Data: map[string]interface{}{
					"key1": "value1",
					"key2": "value2",
				},
				Update: false,
			},
			status:             http.StatusOK,
			contentType:        ContentTypeJSON,
			err:                "",
			gatewaySaveDataErr: nil,
			responseBody:       "\"success\"",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			endpoint := "/api/v2/data"
			gateway := &MockGatewayer{}

			serializedBody, err := json.Marshal(tc.body)
			require.NoError(t, err)

			var body saveDataRequest
			err = json.Unmarshal(serializedBody, &body)

			if err == nil {
				gateway.On("SaveData", body.Data, body.Update).Return(tc.gatewaySaveDataErr)
			}

			requestJSON, err := json.Marshal(tc.body)
			require.NoError(t, err)

			req, err := http.NewRequest(tc.method, endpoint, bytes.NewBuffer(requestJSON))
			require.NoError(t, err)

			contentType := tc.contentType
			if contentType == "" {
				contentType = ContentTypeJSON
			}

			req.Header.Add("Content-Type", contentType)

			rr := httptest.NewRecorder()

			cfg := defaultMuxConfig()

			handler := newServerMux(cfg, gateway)
			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "got `%v` want `%v`", status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "got `%v`| %d, want `%v`",
					strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				require.Equal(t, tc.responseBody, rr.Body.String(), tc.name)
			}
		})
	}
}

func TestGetData(t *testing.T) {
	type httpBody struct {
		Keys string
	}

	tt := []struct {
		name                 string
		method               string
		body                 *httpBody
		keys                 []string
		status               int
		err                  string
		gatewayGetDataErr    error
		gatewatGetDataResult map[string]interface{}
		responseBody         map[string]string
	}{
		{
			name:   "400 - missing keys",
			method: http.MethodGet,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - missing keys",
		},
		{
			name:   "400 - empty file",
			method: http.MethodGet,
			status: http.StatusBadRequest,
			body: &httpBody{
				Keys: "key1,key2",
			},
			keys:              []string{"key1", "key2"},
			err:               "400 Bad Request - empty file",
			gatewayGetDataErr: errors.New("empty file"),
		},
		{
			name:   "200",
			method: http.MethodGet,
			body: &httpBody{
				Keys: "key1,key2",
			},
			keys:              []string{"key1", "key2"},
			status:            http.StatusOK,
			err:               "",
			gatewayGetDataErr: nil,
			gatewatGetDataResult: map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
			},
			responseBody: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			endpoint := "/api/v2/data"
			gateway := &MockGatewayer{}

			gateway.On("GetData", tc.keys).Return(tc.gatewatGetDataResult, tc.gatewayGetDataErr)

			v := url.Values{}
			if tc.body != nil {
				if tc.body.Keys != "" {
					v.Add("keys", tc.body.Keys)
				}
				if len(v) > 0 {
					endpoint += "?" + v.Encode()
				}
			}

			req, err := http.NewRequest(tc.method, endpoint, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()

			cfg := defaultMuxConfig()

			handler := newServerMux(cfg, gateway)
			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "got `%v` want `%v`", status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()),
					"got `%v`| %d, want `%v`",
					strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				var rlt map[string]string
				err = json.Unmarshal(rr.Body.Bytes(), &rlt)
				require.NoError(t, err)
				fmt.Println(tc.responseBody)
				require.Equal(t, tc.responseBody, rlt)
			}
		})
	}
}

func TestDeleteData(t *testing.T) {
	type httpBody struct {
		Keys string
	}

	tt := []struct {
		name                 string
		method               string
		body                 *httpBody
		keys                 []string
		status               int
		err                  string
		gatewayDeleteDataErr error
		responseBody         string
	}{
		{
			name:   "400 - missing keys",
			method: http.MethodDelete,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - missing keys",
		},
		{
			name:   "400 - empty file",
			method: http.MethodDelete,
			status: http.StatusBadRequest,
			body: &httpBody{
				Keys: "key1,key2",
			},
			keys:                 []string{"key1", "key2"},
			err:                  "400 Bad Request - empty file",
			gatewayDeleteDataErr: errors.New("empty file"),
		},
		{
			name:   "200",
			method: http.MethodDelete,
			body: &httpBody{
				Keys: "key1,key2",
			},
			keys:                 []string{"key1", "key2"},
			status:               http.StatusOK,
			err:                  "",
			gatewayDeleteDataErr: nil,
			responseBody:         "\"success\"",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			endpoint := "/api/v2/data"
			gateway := &MockGatewayer{}

			gateway.On("DeleteData", tc.keys).Return(tc.gatewayDeleteDataErr)

			v := url.Values{}
			if tc.body != nil {
				if tc.body.Keys != "" {
					v.Add("keys", tc.body.Keys)
				}
				if len(v) > 0 {
					endpoint += "?" + v.Encode()
				}
			}

			req, err := http.NewRequest(tc.method, endpoint, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()

			cfg := defaultMuxConfig()

			handler := newServerMux(cfg, gateway)
			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "got `%v` want `%v`", status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "got `%v`| %d, want `%v`",
					strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				require.Equal(t, tc.responseBody, rr.Body.String(), tc.name)
			}
		})
	}
}
