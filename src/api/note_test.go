package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/skycoin/skycoin/src/note"

	"github.com/stretchr/testify/require"
)

func TestGetNotesHandler(t *testing.T) {
	tt := []struct {
		name           string
		method         string
		status         int
		httpBody       string
		getNotesResult map[string]string
		getNotesErr    error
		httpResponse   HTTPResponse
		csrfDisabled   bool
	}{
		{
			name:     "405",
			method:   http.MethodPost,
			status:   http.StatusMethodNotAllowed,
			httpBody: "",
			getNotesResult: map[string]string{
				"a5cf149da9cab9fdff681cec9fe83983aada218a46e26292a2c977ceff5bb1a5": "note1",
				"db6fec68266296fcf6bf98a26cf25d86c83bfc31b8248575724977d90426addd": "note2",
				"fef07801a566c3eafd680c9d29ccc18657c600e8b9d8f2c0eb89e3c98f5019c4": "note3",
			},
			getNotesErr:  nil,
			httpResponse: NewHTTPErrorResponse(405, ""),
		},
		{
			name:           "403",
			method:         http.MethodGet,
			status:         http.StatusForbidden,
			httpBody:       "",
			getNotesResult: nil,
			getNotesErr:    note.ErrNoteAPIDisabled,
			httpResponse:   NewHTTPErrorResponse(403, ""),
		},
		{
			name:     "200",
			method:   http.MethodGet,
			status:   http.StatusOK,
			httpBody: "",
			getNotesResult: map[string]string{
				"a5cf149da9cab9fdff681cec9fe83983aada218a46e26292a2c977ceff5bb1a5": "note1",
				"db6fec68266296fcf6bf98a26cf25d86c83bfc31b8248575724977d90426addd": "note2",
				"fef07801a566c3eafd680c9d29ccc18657c600e8b9d8f2c0eb89e3c98f5019c4": "note3",
			},
			getNotesErr: nil,
			httpResponse: HTTPResponse{
				Data: map[string]string{
					"a5cf149da9cab9fdff681cec9fe83983aada218a46e26292a2c977ceff5bb1a5": "note1",
					"db6fec68266296fcf6bf98a26cf25d86c83bfc31b8248575724977d90426addd": "note2",
					"fef07801a566c3eafd680c9d29ccc18657c600e8b9d8f2c0eb89e3c98f5019c4": "note3",
				},
			},
		},
		{
			name:     "200 - csrf disabled",
			method:   http.MethodGet,
			status:   http.StatusOK,
			httpBody: "",
			getNotesResult: map[string]string{
				"a5cf149da9cab9fdff681cec9fe83983aada218a46e26292a2c977ceff5bb1a5": "note1",
				"db6fec68266296fcf6bf98a26cf25d86c83bfc31b8248575724977d90426addd": "note2",
				"fef07801a566c3eafd680c9d29ccc18657c600e8b9d8f2c0eb89e3c98f5019c4": "note3",
			},
			getNotesErr: nil,
			httpResponse: HTTPResponse{
				Data: map[string]string{
					"a5cf149da9cab9fdff681cec9fe83983aada218a46e26292a2c977ceff5bb1a5": "note1",
					"db6fec68266296fcf6bf98a26cf25d86c83bfc31b8248575724977d90426addd": "note2",
					"fef07801a566c3eafd680c9d29ccc18657c600e8b9d8f2c0eb89e3c98f5019c4": "note3",
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			gateway := &MockGatewayer{}
			gateway.On("GetNotes").Return(tc.getNotesResult, tc.getNotesErr)

			endpoint := "/api/v2/notes"

			req, err := http.NewRequest(tc.method, endpoint, strings.NewReader(tc.httpBody))
			require.NoError(t, err)

			req.Header.Set("Content-Type", ContentTypeForm)

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
			err = json.NewDecoder(rr.Body).Decode(&rsp)
			require.NoError(t, err)

			require.Equal(t, tc.httpResponse.Error, rsp.Error)

			if rsp.Data == nil {
				require.Nil(t, tc.httpResponse.Data)
			} else {
				require.NotNil(t, tc.httpResponse.Data)

				var getNotesRsp map[string]string
				err := json.Unmarshal(rsp.Data, &getNotesRsp)
				require.NoError(t, err)

				require.Equal(t, tc.httpResponse.Data, getNotesRsp)
			}
		})
	}
}

func TestGetNoteHandler(t *testing.T) {
	tt := []struct {
		name          string
		method        string
		status        int
		contentType   string
		query         string
		txID          string
		getNoteResult string
		getNoteErr    error
		httpResponse  HTTPResponse
		csrfDisabled  bool
	}{
		{
			name:        "405",
			method:      http.MethodPut,
			status:      http.StatusMethodNotAllowed,
			contentType: ContentTypeForm,
			query: url.Values{
				"txid": []string{"a5cf149da9cab9fdff681cec9fe83983aada218a46e26292a2c977ceff5bb1a5"},
			}.Encode(),
			txID:          "a5cf149da9cab9fdff681cec9fe83983aada218a46e26292a2c977ceff5bb1a5",
			getNoteResult: "",
			getNoteErr:    note.ErrNoteAPIDisabled,
			httpResponse:  NewHTTPErrorResponse(405, ""),
		},
		{
			name:        "403",
			method:      http.MethodGet,
			status:      http.StatusForbidden,
			contentType: ContentTypeForm,
			query: url.Values{
				"txid": []string{"a5cf149da9cab9fdff681cec9fe83983aada218a46e26292a2c977ceff5bb1a5"},
			}.Encode(),
			txID:          "a5cf149da9cab9fdff681cec9fe83983aada218a46e26292a2c977ceff5bb1a5",
			getNoteResult: "",
			getNoteErr:    note.ErrNoteAPIDisabled,
			httpResponse:  NewHTTPErrorResponse(403, ""),
		},
		{
			name:        "415",
			method:      http.MethodGet,
			status:      http.StatusUnsupportedMediaType,
			contentType: ContentTypeJSON,
			query: url.Values{
				"txid": []string{"a5cf149da9cab9fdff681cec9fe83983aada218a46e26292a2c977ceff5bb1a5"},
			}.Encode(),
			txID:          "a5cf149da9cab9fdff681cec9fe83983aada218a46e26292a2c977ceff5bb1a5",
			getNoteResult: "note1",
			getNoteErr:    nil,
			httpResponse:  NewHTTPErrorResponse(415, ""),
		},
		{
			name:          "400 - Missing txid",
			method:        http.MethodGet,
			status:        http.StatusBadRequest,
			contentType:   ContentTypeForm,
			query:         "",
			txID:          "",
			getNoteResult: "",
			getNoteErr:    note.ErrInvalidTxID,
			httpResponse:  NewHTTPErrorResponse(400, "txid is required"),
		},
		{
			name:        "400 - Invalid txid",
			method:      http.MethodGet,
			status:      http.StatusBadRequest,
			contentType: ContentTypeForm,
			query: url.Values{
				"txid": []string{"txid1"},
			}.Encode(),
			txID:          "txid1",
			getNoteResult: "",
			getNoteErr:    note.ErrInvalidTxID,
			httpResponse:  NewHTTPErrorResponse(400, "txid is invalid"),
		},
		{
			name:        "404",
			method:      http.MethodGet,
			status:      http.StatusNotFound,
			contentType: ContentTypeForm,
			query: url.Values{
				"txid": []string{"a5cf149da9cab9fdff681cec9fe83983aada218a46e26292a2c977ceff5bb1a5"},
			}.Encode(),
			txID:          "a5cf149da9cab9fdff681cec9fe83983aada218a46e26292a2c977ceff5bb1a5",
			getNoteResult: "",
			getNoteErr:    note.ErrNoteNotExist,
			httpResponse:  NewHTTPErrorResponse(404, ""),
		},
		{
			name:        "200",
			method:      http.MethodGet,
			status:      http.StatusOK,
			contentType: ContentTypeForm,
			query: url.Values{
				"txid": []string{"a5cf149da9cab9fdff681cec9fe83983aada218a46e26292a2c977ceff5bb1a5"},
			}.Encode(),
			txID:          "a5cf149da9cab9fdff681cec9fe83983aada218a46e26292a2c977ceff5bb1a5",
			getNoteResult: "note1",
			getNoteErr:    nil,
			httpResponse: HTTPResponse{
				Data: "note1",
			},
		},
		{
			name:        "200 - csrf disabled",
			method:      http.MethodGet,
			status:      http.StatusOK,
			contentType: ContentTypeForm,
			query: url.Values{
				"txid": []string{"a5cf149da9cab9fdff681cec9fe83983aada218a46e26292a2c977ceff5bb1a5"},
			}.Encode(),
			txID:          "a5cf149da9cab9fdff681cec9fe83983aada218a46e26292a2c977ceff5bb1a5",
			getNoteResult: "note1",
			getNoteErr:    nil,
			httpResponse: HTTPResponse{
				Data: "note1",
			},
			csrfDisabled: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			gateway := &MockGatewayer{}
			gateway.On("GetNote", tc.txID).Return(tc.getNoteResult, tc.getNoteErr)

			endpoint := "/api/v2/note"

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
			err = json.NewDecoder(rr.Body).Decode(&rsp)
			require.NoError(t, err)

			require.Equal(t, tc.httpResponse.Error, rsp.Error)

			if rsp.Data == nil {
				require.Nil(t, tc.httpResponse.Data)
			} else {
				require.NotNil(t, tc.httpResponse.Data)

				var getNoteRsp string
				err := json.Unmarshal(rsp.Data, &getNoteRsp)
				require.NoError(t, err)

				require.Equal(t, tc.httpResponse.Data, getNoteRsp)
			}
		})
	}
}

func TestAddNoteHandler(t *testing.T) {
	tt := []struct {
		name         string
		method       string
		status       int
		contentType  string
		httpBody     string
		txID         string
		note         string
		addNoteErr   error
		httpResponse HTTPResponse
		csrfDisabled bool
	}{
		{
			name:        "405",
			method:      http.MethodPut,
			status:      http.StatusMethodNotAllowed,
			contentType: ContentTypeJSON,
			httpBody: toJSON(t, NoteRequest{
				TxID: "a5cf149da9cab9fdff681cec9fe83983aada218a46e26292a2c977ceff5bb1a5",
				Note: "note4",
			}),
			txID:         "a5cf149da9cab9fdff681cec9fe83983aada218a46e26292a2c977ceff5bb1a5",
			note:         "note4",
			addNoteErr:   note.ErrNoteAPIDisabled,
			httpResponse: NewHTTPErrorResponse(405, ""),
		},
		{
			name:        "403",
			method:      http.MethodPost,
			status:      http.StatusForbidden,
			contentType: ContentTypeJSON,
			httpBody: toJSON(t, NoteRequest{
				TxID: "a5cf149da9cab9fdff681cec9fe83983aada218a46e26292a2c977ceff5bb1a5",
				Note: "note4",
			}),
			txID:         "a5cf149da9cab9fdff681cec9fe83983aada218a46e26292a2c977ceff5bb1a5",
			note:         "note4",
			addNoteErr:   note.ErrNoteAPIDisabled,
			httpResponse: NewHTTPErrorResponse(403, ""),
		},
		{
			name:        "415",
			method:      http.MethodPost,
			status:      http.StatusUnsupportedMediaType,
			contentType: ContentTypeForm,
			httpBody: toJSON(t, NoteRequest{
				TxID: "a5cf149da9cab9fdff681cec9fe83983aada218a46e26292a2c977ceff5bb1a5",
				Note: "note4",
			}),
			txID:         "a5cf149da9cab9fdff681cec9fe83983aada218a46e26292a2c977ceff5bb1a5",
			note:         "note4",
			addNoteErr:   nil,
			httpResponse: NewHTTPErrorResponse(415, ""),
		},
		{
			name:         "400 - EOF",
			method:       http.MethodPost,
			status:       http.StatusBadRequest,
			contentType:  ContentTypeJSON,
			txID:         "",
			note:         "",
			addNoteErr:   note.ErrInvalidTxID,
			httpResponse: NewHTTPErrorResponse(400, "EOF"),
		},
		{
			name:         "400 - Missing txid",
			method:       http.MethodPost,
			status:       http.StatusBadRequest,
			contentType:  ContentTypeJSON,
			httpBody:     "{}",
			txID:         "",
			note:         "",
			addNoteErr:   note.ErrInvalidTxID,
			httpResponse: NewHTTPErrorResponse(400, "txid is required"),
		},
		{
			name:        "400 - Invalid txid",
			method:      http.MethodPost,
			status:      http.StatusBadRequest,
			contentType: ContentTypeJSON,
			httpBody: toJSON(t, NoteRequest{
				TxID: "txid1",
				Note: "note4",
			}),
			txID:         "txid1",
			note:         "note4",
			addNoteErr:   note.ErrInvalidTxID,
			httpResponse: NewHTTPErrorResponse(400, "txid is invalid"),
		},
		{
			name:        "200",
			method:      http.MethodPost,
			status:      http.StatusOK,
			contentType: ContentTypeJSON,
			httpBody: toJSON(t, NoteRequest{
				TxID: "a5cf149da9cab9fdff681cec9fe83983aada218a46e26292a2c977ceff5bb1a5",
				Note: "note4",
			}),
			txID:         "a5cf149da9cab9fdff681cec9fe83983aada218a46e26292a2c977ceff5bb1a5",
			note:         "note4",
			addNoteErr:   nil,
			httpResponse: HTTPResponse{},
		},
		{
			name:        "403 - csrf disabled",
			method:      http.MethodPost,
			status:      http.StatusForbidden,
			contentType: ContentTypeJSON,
			httpBody: toJSON(t, NoteRequest{
				TxID: "a5cf149da9cab9fdff681cec9fe83983aada218a46e26292a2c977ceff5bb1a5",
				Note: "note4",
			}),
			txID:         "a5cf149da9cab9fdff681cec9fe83983aada218a46e26292a2c977ceff5bb1a5",
			note:         "note4",
			addNoteErr:   nil,
			httpResponse: NewHTTPErrorResponse(403, "invalid CSRF token"),
			csrfDisabled: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			gateway := &MockGatewayer{}
			gateway.On("AddNote", tc.txID, tc.note).Return(tc.addNoteErr)

			endpoint := "/api/v2/note"

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
			err = json.NewDecoder(rr.Body).Decode(&rsp)
			require.NoError(t, err)

			require.Equal(t, tc.httpResponse.Error, rsp.Error)

			require.Nil(t, tc.httpResponse.Data)
		})
	}
}

func TestRemoveNoteHandler(t *testing.T) {
	tt := []struct {
		name          string
		method        string
		status        int
		contentType   string
		query         string
		txID          string
		removeNoteErr error
		httpResponse  HTTPResponse
		csrfDisabled  bool
	}{
		{
			name:        "405",
			method:      http.MethodPut,
			status:      http.StatusMethodNotAllowed,
			contentType: ContentTypeForm,
			query: url.Values{
				"txid": []string{"a5cf149da9cab9fdff681cec9fe83983aada218a46e26292a2c977ceff5bb1a5"},
			}.Encode(),
			txID:          "a5cf149da9cab9fdff681cec9fe83983aada218a46e26292a2c977ceff5bb1a5",
			removeNoteErr: note.ErrNoteAPIDisabled,
			httpResponse:  NewHTTPErrorResponse(405, ""),
		},
		{
			name:        "403",
			method:      http.MethodDelete,
			status:      http.StatusForbidden,
			contentType: ContentTypeForm,
			query: url.Values{
				"txid": []string{"a5cf149da9cab9fdff681cec9fe83983aada218a46e26292a2c977ceff5bb1a5"},
			}.Encode(),
			txID:          "a5cf149da9cab9fdff681cec9fe83983aada218a46e26292a2c977ceff5bb1a5",
			removeNoteErr: note.ErrNoteAPIDisabled,
			httpResponse:  NewHTTPErrorResponse(403, ""),
		},
		{
			name:          "415",
			method:        http.MethodDelete,
			status:        http.StatusUnsupportedMediaType,
			contentType:   ContentTypeJSON,
			txID:          "",
			removeNoteErr: note.ErrInvalidTxID,
			httpResponse:  NewHTTPErrorResponse(415, ""),
		},
		{
			name:          "400 - Missing txid",
			method:        http.MethodDelete,
			status:        http.StatusBadRequest,
			contentType:   ContentTypeForm,
			txID:          "",
			removeNoteErr: note.ErrInvalidTxID,
			httpResponse:  NewHTTPErrorResponse(400, "txid is required"),
		},
		{
			name:        "400 - Invalid txid",
			method:      http.MethodDelete,
			status:      http.StatusBadRequest,
			contentType: ContentTypeForm,
			query: url.Values{
				"txid": []string{"txid1"},
			}.Encode(),
			txID:          "txid1",
			removeNoteErr: note.ErrInvalidTxID,
			httpResponse:  NewHTTPErrorResponse(400, "txid is invalid"),
		},
		{
			name:        "404",
			method:      http.MethodDelete,
			status:      http.StatusNotFound,
			contentType: ContentTypeForm,
			query: url.Values{
				"txid": []string{"a5cf149da9cab9fdff681cec9fe83983aada218a46e26292a2c977ceff5bb1a5"},
			}.Encode(),
			txID:          "a5cf149da9cab9fdff681cec9fe83983aada218a46e26292a2c977ceff5bb1a5",
			removeNoteErr: note.ErrNoteNotExist,
			httpResponse:  NewHTTPErrorResponse(404, ""),
		},
		{
			name:        "200",
			method:      http.MethodDelete,
			status:      http.StatusOK,
			contentType: ContentTypeForm,
			query: url.Values{
				"txid": []string{"a5cf149da9cab9fdff681cec9fe83983aada218a46e26292a2c977ceff5bb1a5"},
			}.Encode(),
			txID:          "a5cf149da9cab9fdff681cec9fe83983aada218a46e26292a2c977ceff5bb1a5",
			removeNoteErr: nil,
			httpResponse:  HTTPResponse{},
		},
		{
			name:        "403 - csrf disabled",
			method:      http.MethodDelete,
			status:      http.StatusForbidden,
			contentType: ContentTypeForm,
			query: url.Values{
				"txid": []string{"a5cf149da9cab9fdff681cec9fe83983aada218a46e26292a2c977ceff5bb1a5"},
			}.Encode(),
			txID:          "a5cf149da9cab9fdff681cec9fe83983aada218a46e26292a2c977ceff5bb1a5",
			removeNoteErr: nil,
			httpResponse:  NewHTTPErrorResponse(403, "invalid CSRF token"),
			csrfDisabled:  true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			gateway := &MockGatewayer{}
			gateway.On("RemoveNote", tc.txID).Return(tc.removeNoteErr)

			endpoint := "/api/v2/note"

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
			err = json.NewDecoder(rr.Body).Decode(&rsp)
			require.NoError(t, err)

			require.Equal(t, tc.httpResponse.Error, rsp.Error)

			require.Nil(t, tc.httpResponse.Data)
		})
	}
}
