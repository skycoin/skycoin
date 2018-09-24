package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/notes"
)

func TestGetAllNotes(t *testing.T) {
	tt := []struct {
		name                     string
		method                   string
		body                     *notes.Note
		status                   int
		err                      string
		gatewayGetAllNotesResult []notes.Note
		responseBody             []notes.Note
		gatewayGetAllNotesErr    error
	}{
		{
			name:   "405",
			method: http.MethodGet,
			status: http.StatusMethodNotAllowed,
			err:    "405 Method Not Allowed",
		},
		{
			name:   "200 - OK",
			method: http.MethodPost,
			body:   nil,
			status: http.StatusOK,
			gatewayGetAllNotesResult: []notes.Note{
				{
					TxIDHex: "62b1e205aa2895b7094f708d853a64709e14d467ef3e3eee54ef79bcefdbd4c8",
					Notes:   "A Note... ",
				},
				{
					TxIDHex: "9c8995afd843372636ae66991797c824e2fd8dfffa77c901c7f9e8d4f5e8711",
					Notes:   "Another note...",
				},
				{
					TxIDHex: "9c8995afd843372636ae66991797c824e2fd8dfffa77c901c7f9e8d4f5e87114",
					Notes:   "Last note",
				},
			},
			responseBody: []notes.Note{
				{
					TxIDHex: "62b1e205aa2895b7094f708d853a64709e14d467ef3e3eee54ef79bcefdbd4c8",
					Notes:   "A Note... ",
				},
				{
					TxIDHex: "9c8995afd843372636ae66991797c824e2fd8dfffa77c901c7f9e8d4f5e8711",
					Notes:   "Another note...",
				},
				{
					TxIDHex: "9c8995afd843372636ae66991797c824e2fd8dfffa77c901c7f9e8d4f5e87114",
					Notes:   "Last note",
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			endpoint := "/api/v2/notes/notes"
			gateway := &MockGatewayer{}
			gateway.On("GetAllNotes").Return(tc.gatewayGetAllNotesResult, tc.gatewayGetAllNotesErr)

			req, err := http.NewRequest(tc.method, endpoint, nil)
			require.NoError(t, err)

			csrfStore := &CSRFStore{
				Enabled: true,
			}
			setCSRFParameters(csrfStore, tokenValid, req)

			rr := httptest.NewRecorder()
			handler := newServerMux(defaultMuxConfig(), gateway, csrfStore, nil)

			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "got `%v` want `%v`", status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()),
					"got `%v`| %d, want `%v`",
					strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				var rlt []notes.Note
				err = json.Unmarshal(rr.Body.Bytes(), &rlt)
				require.NoError(t, err)
				require.Equal(t, tc.responseBody, rlt)
			}
		})
	}
}

func TestGetNoteByTxID(t *testing.T) {
	tt := []struct {
		name         string
		method       string
		body         *notes.Note
		status       int
		err          string
		responseBody notes.Note
	}{
		{
			name:   "405",
			method: http.MethodGet,
			status: http.StatusMethodNotAllowed,
			err:    "405 Method Not Allowed",
		},
		{
			name:         "400",
			method:       http.MethodPost,
			status:       http.StatusBadRequest,
			body:         &notes.Note{TxIDHex: "tooShortTxID"},
			responseBody: notes.Note{TxIDHex: "9c8995afd843372636ae66991797c824e2fd8dfffa77c901c7f9e8d4f5e87114", Notes: ""},
			err:          "400 Bad Request - Wrong txid",
		},
		{
			name:         "200 - OK",
			method:       http.MethodPost,
			status:       http.StatusOK,
			body:         &notes.Note{TxIDHex: "9c8995afd843372636ae66991797c824e2fd8dfffa77c901c7f9e8d4f5e87114", Notes: ""},
			responseBody: notes.Note{TxIDHex: "9c8995afd843372636ae66991797c824e2fd8dfffa77c901c7f9e8d4f5e87114", Notes: ""},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var jsonStr []byte
			var err error
			gateway := &MockGatewayer{}

			endpoint := "/api/v2/notes/noteByTxid"

			if tc.body != nil {
				gateway.On("GetNoteByTxID", tc.body.TxIDHex).Return(tc.responseBody, tc.err)

				jsonStr, err = json.Marshal(tc.body)
				if err != nil {
					t.Error(err)
				}
			}

			req, err := http.NewRequest(tc.method, endpoint, bytes.NewBuffer(jsonStr))
			require.NoError(t, err)

			csrfStore := &CSRFStore{
				Enabled: true,
			}
			setCSRFParameters(csrfStore, tokenValid, req)

			rr := httptest.NewRecorder()
			handler := newServerMux(defaultMuxConfig(), gateway, csrfStore, nil)

			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "got `%v` want `%v`", status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()),
					"got `%v`| %d, want `%v`",
					strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				var rlt notes.Note
				err = json.Unmarshal(rr.Body.Bytes(), &rlt)
				require.NoError(t, err)
				require.Equal(t, tc.responseBody, rlt)
			}
		})
	}
}

func TestAddNote(t *testing.T) {
	tt := []struct {
		name              string
		method            string
		body              *notes.Note
		status            int
		err               string
		responseBody      notes.Note
		gatewayAddNoteErr *error
	}{
		{
			name:   "405",
			method: http.MethodGet,
			status: http.StatusMethodNotAllowed,
			err:    "405 Method Not Allowed",
		},
		{
			name:              "400",
			method:            http.MethodPost,
			status:            http.StatusBadRequest,
			err:               "400 Bad Request - bad parameters",
			gatewayAddNoteErr: nil,
			body: &notes.Note{
				TxIDHex: "wrongtxid",
				Notes:   "",
			},
		},
		{
			name:              "400",
			method:            http.MethodPost,
			status:            http.StatusBadRequest,
			err:               "400 Bad Request - bad parameters",
			gatewayAddNoteErr: nil,
			body: &notes.Note{
				TxIDHex: "62b1e205aa2895b7094f708d853a64709e14d467ef3e3eee54ef79bcefdbd4c8",
				Notes:   "",
			},
			responseBody: notes.Note{
				TxIDHex: "62b1e205aa2895b7094f708d853a64709e14d467ef3e3eee54ef79bcefdbd4c8",
				Notes:   "A Note that is not empty",
			},
		},
		{
			name:              "200 - OK",
			method:            http.MethodPost,
			status:            http.StatusOK,
			gatewayAddNoteErr: nil,
			body: &notes.Note{
				TxIDHex: "62b1e205aa2895b7094f708d853a64709e14d467ef3e3eee54ef79bcefdbd4c8",
				Notes:   "A Note that is not empty",
			},
			responseBody: notes.Note{
				TxIDHex: "62b1e205aa2895b7094f708d853a64709e14d467ef3e3eee54ef79bcefdbd4c8",
				Notes:   "A Note that is not empty",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var jsonStr []byte
			var err error
			endpoint := "/api/v2/notes/addNote"
			gateway := &MockGatewayer{}

			if tc.body != nil {
				gateway.On("AddNote", *tc.body).Return(tc.responseBody, nil)

				note := notes.Note{TxIDHex: tc.body.TxIDHex, Notes: tc.body.Notes}
				jsonStr, err = json.Marshal(note)

				if err != nil {
					t.Error(err)
				}
			}

			req, err := http.NewRequest(tc.method, endpoint, bytes.NewBuffer(jsonStr))
			require.NoError(t, err)

			csrfStore := &CSRFStore{
				Enabled: true,
			}
			setCSRFParameters(csrfStore, tokenValid, req)

			rr := httptest.NewRecorder()
			handler := newServerMux(defaultMuxConfig(), gateway, csrfStore, nil)
			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "got `%v` want `%v`", status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()),
					"got `%v`| %d, want `%v`",
					strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				var rlt notes.Note
				err = json.Unmarshal(rr.Body.Bytes(), &rlt)
				require.NoError(t, err)
				require.Equal(t, tc.responseBody, rlt)
			}
		})
	}
}

func TestRemoveNote(t *testing.T) {
	tt := []struct {
		name         string
		method       string
		body         *notes.Note
		status       int
		err          string
		responseBody notes.Note
	}{
		{
			name:   "405",
			method: http.MethodGet,
			status: http.StatusMethodNotAllowed,
			err:    "405 Method Not Allowed",
		},
		{
			name:   "400",
			method: http.MethodPost,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - wrong 'txid'",
			body: &notes.Note{
				TxIDHex: "wrongtxid",
				Notes:   "",
			},
			responseBody: notes.Note{},
		},
		{
			name:   "200 - OK",
			method: http.MethodPost,
			status: http.StatusOK,
			body: &notes.Note{
				TxIDHex: "62b1e205aa2895b7094f708d853a64709e14d467ef3e3eee54ef79bcefdbd4c8",
				Notes:   "",
			},
			responseBody: notes.Note{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var jsonStr []byte
			var err error
			endpoint := "/api/v2/notes/removeNote"
			gateway := &MockGatewayer{}

			if tc.body != nil {
				gateway.On("RemoveNote", tc.body.TxIDHex).Return(err, nil)

				jsonStr, err = json.Marshal(tc.body)

				if err != nil {
					t.Error(err)
				}
			}

			req, err := http.NewRequest(tc.method, endpoint, bytes.NewBuffer(jsonStr))
			require.NoError(t, err)

			csrfStore := &CSRFStore{
				Enabled: true,
			}
			setCSRFParameters(csrfStore, tokenValid, req)

			rr := httptest.NewRecorder()
			handler := newServerMux(defaultMuxConfig(), gateway, csrfStore, nil)
			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "got `%v` want `%v`", status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()),
					"got `%v`| %d, want `%v`",
					strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				var rlt notes.Note
				err = json.Unmarshal(rr.Body.Bytes(), &rlt)
				require.NoError(t, err)
				require.Equal(t, tc.responseBody, rlt)
			}
		})
	}
}
