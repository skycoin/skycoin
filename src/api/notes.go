package api

import (
	"encoding/json"
	"net/http"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/notes"

	wh "github.com/skycoin/skycoin/src/util/http"
)

var (
	// ErrorWrongTxID appears when TransactionID has wrong format
	ErrorWrongTxID = "wrong 'txid'"
	// ErrorBadParams appears when note obj isn't complete
	ErrorBadParams = "bad parameters"
)

// URI: /api/v2/notes
// Method: GET
// Content-Type: application/json
// Body: -
// Response:
//      200 - ok, returns all notes
func getAllNotesHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			resp := NewHTTPErrorResponse(http.StatusUnsupportedMediaType, "")
			writeHTTPResponse(w, resp)
			return
		}

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
//			- GET: return note by Transaction ID
//			- DELETE: removes note by Transaction ID
func noteHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			resp := NewHTTPErrorResponse(http.StatusUnsupportedMediaType, "")
			writeHTTPResponse(w, resp)
			return
		}

		if r.Method == http.MethodPost {
			var note notes.Note
			if err := json.NewDecoder(r.Body).Decode(&note); err != nil || len(note.Notes) == 0 {
				wh.Error400(w, ErrorBadParams)
				return
			}

			if err := validateTxID(note.TxIDHex); err == nil {
				note, err := gateway.AddNote(note)
				if err != nil {
					wh.Error400(w, err.Error())
					return
				}
				wh.SendJSONOr500(logger, w, note)
				return
			}
			wh.Error400(w, ErrorWrongTxID)
			return
		}

		txID := r.FormValue("txid")

		switch r.Method {
		case http.MethodGet:
			if err := validateTxID(txID); err != nil {
				wh.Error400(w, ErrorWrongTxID)
				return
			}

			noteByTxID := gateway.GetNoteByTxID(txID)
			wh.SendJSONOr500(logger, w, noteByTxID)
			return

		case http.MethodDelete:
			if err := validateTxID(txID); err != nil {
				wh.Error400(w, ErrorWrongTxID)
				return
			}

			if err := gateway.RemoveNote(txID); err != nil {
				wh.Error400(w, err.Error())
				return
			}
			wh.SendJSONOr500(logger, w, notes.Note{})
			return

		default:
			// Bad Method
			wh.Error405(w)
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
