package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
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
			method: http.MethodPost,
			status: http.StatusMethodNotAllowed,
			err:    "405 Method Not Allowed",
		},
		{
			name:   "200 - OK",
			method: http.MethodGet,
			body:   nil,
			status: http.StatusOK,
			gatewayGetAllNotesResult: []notes.Note{
				{
					TxIDHex: "62b1e205aa2895b7094f708d853a64709e14d467ef3e3eee54ef79bcefdbd4c8",
					Notes:   "A Note... ",
				},
				{
					TxIDHex: "9c8995afd843372636ae66991797c824e2fd8dfffa77c901c7f9e8d4f5e87113",
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
					TxIDHex: "9c8995afd843372636ae66991797c824e2fd8dfffa77c901c7f9e8d4f5e87113",
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
			endpoint := "/api/v2/notes"
			gateway := &MockGatewayer{}
			gateway.On("GetAllNotes").Return(tc.gatewayGetAllNotesResult, tc.gatewayGetAllNotesErr)

			req, err := http.NewRequest(tc.method, endpoint, nil)
			req.Header.Add("Content-Type", "application/json")
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
		txID         string
		status       int
		err          string
		responseBody notes.Note
	}{
		{
			name:   "405",
			method: http.MethodPut,
			status: http.StatusMethodNotAllowed,
			err:    "405 Method Not Allowed",
		},
		{
			name:   "400",
			method: http.MethodGet,
			status: http.StatusBadRequest,
			txID:   "tooShortTxID",
			err:    "400 Bad Request - " + ErrorWrongTxID,
		},
		{
			name:         "200 - OK",
			method:       http.MethodGet,
			status:       http.StatusOK,
			txID:         "9c8995afd843372636ae66991797c824e2fd8dfffa77c901c7f9e8d4f5e87114",
			responseBody: notes.Note{TxIDHex: "9c8995afd843372636ae66991797c824e2fd8dfffa77c901c7f9e8d4f5e87114", Notes: ""},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			gateway := &MockGatewayer{}

			endpoint := "/api/v2/note"
			gateway.On("GetNoteByTxID", tc.txID).Return(tc.responseBody, tc.err)

			v := url.Values{}
			v.Add("txid", tc.txID)

			endpoint += "?" + v.Encode()

			req, err := http.NewRequest(tc.method, endpoint, bytes.NewBufferString(v.Encode()))
			req.Header.Add("Content-Type", "application/json")
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
			method: http.MethodPatch,
			status: http.StatusMethodNotAllowed,
			err:    "405 Method Not Allowed",
		},
		{
			name:              "400",
			method:            http.MethodPost,
			status:            http.StatusBadRequest,
			err:               "400 Bad Request - " + ErrorBadParams,
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
			err:               "400 Bad Request - " + ErrorBadParams,
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
			endpoint := "/api/v2/note"
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
			req.Header.Add("Content-Type", "application/json")
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
		txID         string
		status       int
		err          string
		responseBody notes.Note
	}{
		{
			name:   "405",
			method: http.MethodPut,
			status: http.StatusMethodNotAllowed,
			err:    "405 Method Not Allowed",
			txID:   "62b1e205aa2895b7094f708d853a64709e14d467ef3e3eee54ef79bcefdbd4c8",
		},
		{
			name:         "400",
			method:       http.MethodDelete,
			status:       http.StatusBadRequest,
			err:          "400 Bad Request - " + ErrorWrongTxID,
			txID:         "wrongtxid",
			responseBody: notes.Note{},
		},
		{
			name:         "200 - OK",
			method:       http.MethodDelete,
			status:       http.StatusOK,
			txID:         "62b1e205aa2895b7094f708d853a64709e14d467ef3e3eee54ef79bcefdbd4c8",
			responseBody: notes.Note{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var err error

			gateway := &MockGatewayer{}
			gateway.On("RemoveNote", tc.txID).Return(err, nil)

			endpoint := "/api/v2/note"

			v := url.Values{}
			v.Add("txid", tc.txID)

			endpoint += "?" + v.Encode()

			req, err := http.NewRequest(tc.method, endpoint, bytes.NewBufferString(v.Encode()))
			req.Header.Add("Content-Type", "application/json")
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
