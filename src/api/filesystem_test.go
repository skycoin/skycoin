package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
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
		contentType        string
		gatewaySaveDataErr error
		httpResponse       HTTPResponse
	}{
		{
			name:         "415 - Unsupported Media Type",
			method:       http.MethodPost,
			contentType:  ContentTypeForm,
			status:       http.StatusUnsupportedMediaType,
			httpResponse: NewHTTPErrorResponse(http.StatusUnsupportedMediaType, ""),
		},
		{
			name:         "400 - empty data",
			method:       http.MethodPost,
			status:       http.StatusBadRequest,
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "empty data"),
			body:         &httpBody{},
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
			gatewaySaveDataErr: nil,
			httpResponse:       HTTPResponse{Data: struct{}{}},
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

			var rsp ReceivedHTTPResponse
			err = json.NewDecoder(rr.Body).Decode(&rsp)
			require.NoError(t, err)

			require.Equal(t, tc.httpResponse.Error, rsp.Error)

			if rsp.Data == nil {
				require.Nil(t, tc.httpResponse.Data)
			} else {
				require.NotNil(t, tc.httpResponse.Data)

				var saveRsp struct{}
				err := json.Unmarshal(rsp.Data, &saveRsp)
				require.NoError(t, err)

				require.Equal(t, tc.httpResponse.Data, saveRsp)
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
		gatewayGetDataErr    error
		gatewatGetDataResult map[string]interface{}
		httpResponse         HTTPResponse
	}{
		{
			name:         "400 - missing keys",
			method:       http.MethodGet,
			status:       http.StatusBadRequest,
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "missing keys"),
		},
		{
			name:   "400 - empty file",
			method: http.MethodGet,
			status: http.StatusBadRequest,
			body: &httpBody{
				Keys: "key1,key2",
			},
			keys:              []string{"key1", "key2"},
			gatewayGetDataErr: errors.New("empty file"),
			httpResponse:      NewHTTPErrorResponse(http.StatusBadRequest, "empty file"),
		},
		{
			name:   "200",
			method: http.MethodGet,
			body: &httpBody{
				Keys: "key1,key2",
			},
			keys:              []string{"key1", "key2"},
			status:            http.StatusOK,
			gatewayGetDataErr: nil,
			gatewatGetDataResult: map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
			},
			httpResponse: HTTPResponse{
				Data: map[string]string{
					"key1": "value1",
					"key2": "value2",
				},
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

			var rsp ReceivedHTTPResponse
			err = json.NewDecoder(rr.Body).Decode(&rsp)
			require.NoError(t, err)

			require.Equal(t, tc.httpResponse.Error, rsp.Error)

			if rsp.Data == nil {
				require.Nil(t, tc.httpResponse.Data)
			} else {
				require.NotNil(t, tc.httpResponse.Data)

				var getRsp map[string]string
				err := json.Unmarshal(rsp.Data, &getRsp)
				require.NoError(t, err)

				require.Equal(t, tc.httpResponse.Data, getRsp)
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
		gatewayDeleteDataErr error
		httpResponse         HTTPResponse
	}{
		{
			name:         "400 - missing keys",
			method:       http.MethodDelete,
			status:       http.StatusBadRequest,
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "missing keys"),
		},
		{
			name:   "400 - empty file",
			method: http.MethodDelete,
			status: http.StatusBadRequest,
			body: &httpBody{
				Keys: "key1,key2",
			},
			keys:                 []string{"key1", "key2"},
			gatewayDeleteDataErr: errors.New("empty file"),
			httpResponse:         NewHTTPErrorResponse(http.StatusBadRequest, "empty file"),
		},
		{
			name:   "200",
			method: http.MethodDelete,
			body: &httpBody{
				Keys: "key1,key2",
			},
			keys:                 []string{"key1", "key2"},
			status:               http.StatusOK,
			gatewayDeleteDataErr: nil,
			httpResponse:         HTTPResponse{Data: struct{}{}},
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

			var rsp ReceivedHTTPResponse
			err = json.NewDecoder(rr.Body).Decode(&rsp)
			require.NoError(t, err)

			require.Equal(t, tc.httpResponse.Error, rsp.Error)

			if rsp.Data == nil {
				require.Nil(t, tc.httpResponse.Data)
			} else {
				require.NotNil(t, tc.httpResponse.Data)

				var deleteRsp struct{}
				err := json.Unmarshal(rsp.Data, &deleteRsp)
				require.NoError(t, err)

				require.Equal(t, tc.httpResponse.Data, deleteRsp)
			}
		})
	}
}
