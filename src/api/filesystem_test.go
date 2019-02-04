package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSaveData(t *testing.T) {
	type httpBody struct {
		filename string
		data     map[string]interface{}
		update   bool
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
			name:   "405",
			method: http.MethodGet,
			status: http.StatusMethodNotAllowed,
			err:    "405 Method Not Allowed",
		},

		{
			name:        "415",
			method:      http.MethodPost,
			status:      http.StatusUnsupportedMediaType,
			contentType: ContentTypeForm,
			err:         "415 Unsupported Media Type",
		},
		{
			name:   "200",
			method: http.MethodPost,
			body: &httpBody{
				filename: "foo.json",
				data: map[string]interface{}{
					"key1": "value1",
					"key2": "value2",
				},
				update: false,
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
			endpoint := "/api/v2/data/save"
			gateway := &MockGatewayer{}

			serializedBody, err := json.Marshal(tc.body)
			require.NoError(t, err)

			var body saveDataRequest
			err = json.Unmarshal(serializedBody, &body)

			if err == nil {
				gateway.On("SaveData", body.Filename, body.Data, body.Update).Return(tc.gatewaySaveDataErr)
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

			handler := newServerMux(cfg, gateway, nil)
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
		filename string
		keys     []string
	}

	tt := []struct {
		name                 string
		method               string
		body                 *httpBody
		status               int
		err                  string
		contentType          string
		gatewayGetDataErr    error
		gatewatGetDataResult map[string]interface{}
		responseBody         map[string]string
	}{
		{
			name:   "405",
			method: http.MethodGet,
			status: http.StatusMethodNotAllowed,
			err:    "405 Method Not Allowed",
		},

		{
			name:        "415",
			method:      http.MethodPost,
			status:      http.StatusUnsupportedMediaType,
			contentType: ContentTypeForm,
			err:         "415 Unsupported Media Type",
		},
		{
			name:   "200",
			method: http.MethodPost,
			body: &httpBody{
				filename: "foo.json",
				keys: []string{
					"key1",
					"key2",
				},
			},
			status:            http.StatusOK,
			contentType:       ContentTypeJSON,
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
			endpoint := "/api/v2/data/get"
			gateway := &MockGatewayer{}

			serializedBody, err := json.Marshal(tc.body)
			require.NoError(t, err)

			var body getDataRequest
			err = json.Unmarshal(serializedBody, &body)

			if err == nil {
				gateway.On("GetData", body.Filename, body.Keys).Return(tc.gatewatGetDataResult, tc.gatewayGetDataErr)
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

			handler := newServerMux(cfg, gateway, nil)
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
		filename string
		keys     []string
	}

	tt := []struct {
		name                 string
		method               string
		body                 *httpBody
		status               int
		err                  string
		contentType          string
		gatewayDeleteDataErr error
		responseBody         string
	}{
		{
			name:   "405",
			method: http.MethodGet,
			status: http.StatusMethodNotAllowed,
			err:    "405 Method Not Allowed",
		},

		{
			name:        "415",
			method:      http.MethodPost,
			status:      http.StatusUnsupportedMediaType,
			contentType: ContentTypeForm,
			err:         "415 Unsupported Media Type",
		},
		{
			name:   "200",
			method: http.MethodPost,
			body: &httpBody{
				filename: "foo.json",
				keys: []string{
					"key1",
					"key2",
				},
			},
			status:               http.StatusOK,
			contentType:          ContentTypeJSON,
			err:                  "",
			gatewayDeleteDataErr: nil,
			responseBody:         "\"success\"",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			endpoint := "/api/v2/data/delete"
			gateway := &MockGatewayer{}

			serializedBody, err := json.Marshal(tc.body)
			require.NoError(t, err)

			var body getDataRequest
			err = json.Unmarshal(serializedBody, &body)

			if err == nil {
				gateway.On("DeleteData", body.Filename, body.Keys).Return(tc.gatewayDeleteDataErr)
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

			handler := newServerMux(cfg, gateway, nil)
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
