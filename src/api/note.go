package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/skycoin/skycoin/src/note"
)

// Returns all existing notes.
// Method: GET
// URI: /api/v2/notes
func notesHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			resp := NewHTTPErrorResponse(http.StatusMethodNotAllowed, "")
			writeHTTPResponse(w, resp)
			return
		}

		notes, err := gateway.GetNotes()
		if err != nil {
			var resp HTTPResponse
			switch err {
			case note.ErrNoteAPIDisabled:
				resp = NewHTTPErrorResponse(http.StatusForbidden, "")
			default:
				resp = NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
			}
			writeHTTPResponse(w, resp)
			return
		}

		writeHTTPResponse(w, HTTPResponse{
			Data: notes,
		})
	}
}

// NoteRequest is the request data for POST /api/v2/note
type NoteRequest struct {
	TxID string `json:"txid"`
	Note string `json:"note"`
}

// Dispatches /note endpoint.
// Method: GET, POST, DELETE
// URI: /api/v2/note
func noteHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getNoteHandler(w, r, gateway)
		case http.MethodPost:
			addNoteHandler(w, r, gateway)
		case http.MethodDelete:
			removeNoteHandler(w, r, gateway)
		default:
			resp := NewHTTPErrorResponse(http.StatusMethodNotAllowed, "")
			writeHTTPResponse(w, resp)
			return
		}
	}
}

// Returns a note by txid
// Args:
//	txid: transaction id
func getNoteHandler(w http.ResponseWriter, r *http.Request, gateway Gatewayer) {
	if r.Header.Get("Content-Type") != ContentTypeForm {
		resp := NewHTTPErrorResponse(http.StatusUnsupportedMediaType, "")
		writeHTTPResponse(w, resp)
		return
	}

	txID := r.FormValue("txid")
	if txID == "" {
		resp := NewHTTPErrorResponse(http.StatusBadRequest, "txid is required")
		writeHTTPResponse(w, resp)
		return
	}

	n, err := gateway.GetNote(txID)
	if err != nil {
		var resp HTTPResponse
		switch err {
		case note.ErrNoteAPIDisabled:
			resp = NewHTTPErrorResponse(http.StatusForbidden, "")
		case note.ErrInvalidTxID:
			resp = NewHTTPErrorResponse(http.StatusBadRequest, "txid is invalid")
		case note.ErrNoteNotExist:
			resp = NewHTTPErrorResponse(http.StatusNotFound, "")
		default:
			resp = NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
		}
		writeHTTPResponse(w, resp)
		return
	}

	writeHTTPResponse(w, HTTPResponse{
		Data: n,
	})
}

// Adds a note to txid. Note leading and trailing white spaces are removed.
// Args:
//	txid: transaction id
//	note: transaction note
func addNoteHandler(w http.ResponseWriter, r *http.Request, gateway Gatewayer) {
	if r.Header.Get("Content-Type") != ContentTypeJSON {
		resp := NewHTTPErrorResponse(http.StatusUnsupportedMediaType, "")
		writeHTTPResponse(w, resp)
		return
	}

	var req NoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		resp := NewHTTPErrorResponse(http.StatusBadRequest, err.Error())
		writeHTTPResponse(w, resp)
		return
	}

	if req.TxID == "" {
		resp := NewHTTPErrorResponse(http.StatusBadRequest, "txid is required")
		writeHTTPResponse(w, resp)
		return
	}

	trimmedNote := strings.TrimSpace(req.Note)

	if err := gateway.AddNote(req.TxID, trimmedNote); err != nil {
		var resp HTTPResponse
		switch err {
		case note.ErrNoteAPIDisabled:
			resp = NewHTTPErrorResponse(http.StatusForbidden, "")
		case note.ErrInvalidTxID:
			resp = NewHTTPErrorResponse(http.StatusBadRequest, "txid is invalid")
		default:
			resp = NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
		}
		writeHTTPResponse(w, resp)
		return
	}

	writeHTTPResponse(w, HTTPResponse{})
}

// Removes a note by txid
// Args:
//	txid: transaction id
func removeNoteHandler(w http.ResponseWriter, r *http.Request, gateway Gatewayer) {
	if r.Header.Get("Content-Type") != ContentTypeForm {
		resp := NewHTTPErrorResponse(http.StatusUnsupportedMediaType, "")
		writeHTTPResponse(w, resp)
		return
	}

	txID := r.FormValue("txid")
	if txID == "" {
		resp := NewHTTPErrorResponse(http.StatusBadRequest, "txid is required")
		writeHTTPResponse(w, resp)
		return
	}

	if err := gateway.RemoveNote(txID); err != nil {
		var resp HTTPResponse
		switch err {
		case note.ErrNoteAPIDisabled:
			resp = NewHTTPErrorResponse(http.StatusForbidden, "")
		case note.ErrNoteNotExist:
			resp = NewHTTPErrorResponse(http.StatusNotFound, "")
		case note.ErrInvalidTxID:
			resp = NewHTTPErrorResponse(http.StatusBadRequest, "txid is invalid")
		default:
			resp = NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
		}
		writeHTTPResponse(w, resp)
		return
	}

	writeHTTPResponse(w, HTTPResponse{})
}
