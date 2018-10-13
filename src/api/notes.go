package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/notes"

	wh "github.com/skycoin/skycoin/src/util/http"
)

const (
	// ErrorWrongTxID appears when TransactionID has wrong format
	ErrorWrongTxID = "wrong 'txid'"
	// ErrorBadParams appears when note obj isn't complete
	ErrorBadParams = "bad parameters"
)

// URI: /api/v2/notes
// Method: GET
// Response:
//      200 - ok, returns all notes
func getAllNotesHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}
		savedNotes := gateway.GetAllNotes()

		wh.SendJSONOr500(logger, w, savedNotes)
	}
}

// URI: /api/v2/note
// Method: POST, GET, DELETE
// Content-Type: application/json
// Body: { "txid": "<Transaction ID>" }
// Response:
//      400 - bad parameters
//      200 - POST: returns added Note
//			- GET: returns note by Transaction ID
//			- DELETE: removes note by Transaction ID
func noteHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" && r.Method == http.MethodPost {
			resp := NewHTTPErrorResponse(http.StatusUnsupportedMediaType, "")
			writeHTTPResponse(w, resp)
			return
		}

		// Add Note
		if r.Method == http.MethodPost {
			var note notes.Note

			if err := json.NewDecoder(r.Body).Decode(&note); err != nil || len(note.Notes) == 0 {
				resp := NewHTTPErrorResponse(http.StatusBadRequest, fmt.Sprint(http.StatusText(http.StatusBadRequest), " - ", ErrorBadParams))
				writeHTTPResponse(w, resp)
				return
			}

			if err := validateTxID(note.TxIDHex); err != nil {
				resp := NewHTTPErrorResponse(http.StatusBadRequest, fmt.Sprint(http.StatusText(http.StatusBadRequest), " - ", ErrorWrongTxID))
				writeHTTPResponse(w, resp)
				return
			}

			note, err := gateway.AddNote(note)
			if err != nil {
				resp := NewHTTPErrorResponse(http.StatusUnprocessableEntity, fmt.Sprint(http.StatusText(http.StatusUnprocessableEntity), " - ", err.Error()))
				writeHTTPResponse(w, resp)
				return
			}

			resp := HTTPResponse{Data: note}
			writeHTTPResponse(w, resp)
			return
		}

		txID := r.FormValue("txid")
		txID = strings.Replace(txID, "\"", "", -1)

		switch r.Method {
		case http.MethodGet:
			// Get Note by TxID
			if err := validateTxID(txID); err != nil {
				resp := NewHTTPErrorResponse(http.StatusBadRequest, fmt.Sprint(http.StatusText(http.StatusBadRequest), " - ", ErrorWrongTxID))
				writeHTTPResponse(w, resp)
				return
			}

			noteByTxID := gateway.GetNoteByTxID(txID)

			resp := HTTPResponse{Data: noteByTxID}
			writeHTTPResponse(w, resp)
			return

		case http.MethodDelete:
			// Remove Note by TxID
			if err := validateTxID(txID); err != nil {
				resp := NewHTTPErrorResponse(http.StatusBadRequest, fmt.Sprint(http.StatusText(http.StatusBadRequest), " - ", ErrorWrongTxID))
				writeHTTPResponse(w, resp)
				return
			}

			if err := gateway.RemoveNote(txID); err != nil {
				resp := NewHTTPErrorResponse(http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity))
				writeHTTPResponse(w, resp)
				return
			}

			resp := HTTPResponse{Data: struct{}{}}
			writeHTTPResponse(w, resp)
			return

		default:
			// Bad request method
			resp := NewHTTPErrorResponse(http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
			writeHTTPResponse(w, resp)
			return
		}
	}
}

func validateTxID(txID string) error {
	if _, err := cipher.SHA256FromHex(txID); err != nil {
		return err
	}
	return nil
}
