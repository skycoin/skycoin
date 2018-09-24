package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/notes"
	wh "github.com/skycoin/skycoin/src/util/http"
)

// URI: /api/v2/notes/notes
// Method: POST
// Content-Type: application/json
// Body: -
// Response:
//      200 - ok, returns all notes
func getAllNotesHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			wh.Error405(w)
			return
		}

		savedNotes := gateway.GetAllNotes()

		wh.SendJSONOr500(logger, w, savedNotes)
	}
}

// URI: /api/v2/notes/noteByTxid
// Method: POST
// Content-Type: application/json
// Body: { "txid": "<Transaction ID>", "notes": "<Notes>" }
// Response:
//      422 - wrong parameters
//      400 - internal server error
//      200 - ok, returns note by TxId
func getNoteByIDHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			wh.Error405(w)
			return
		}

		var note notes.Note
		if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
			wh.Error422(w, err.Error())
			return
		}

		if _, err := cipher.SHA256FromHex(note.TxIDHex); err == nil {
			savedNotes := gateway.GetNoteByTxID(note.TxIDHex)
			wh.SendJSONOr500(logger, w, savedNotes)
		} else {
			wh.Error422(w, fmt.Errorf("Wrong txid").Error())
			return
		}
	}
}

// URI: /api/v2/notes/addNote
// Method: POST
// Content-Type: application/json
// Body: { "txid": "<Transaction ID>", "notes": "<Notes>" }
// Response:
//      400 - wrong parameters
//      200 - ok, note added
func addNoteHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			wh.Error405(w)
			return
		}

		var note notes.Note
		var retNote notes.Note
		if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
			wh.Error400(w, fmt.Errorf("bad parameters").Error())
			return
		}

		if _, err := cipher.SHA256FromHex(note.TxIDHex); err == nil && len(note.Notes) > 0 {
			retNote, err = gateway.AddNote(note)
			if err != nil {
				wh.Error400(w, err.Error())
				return
			}
		} else {
			wh.Error400(w, fmt.Errorf("bad parameters").Error())
			return
		}

		wh.SendJSONOr500(logger, w, retNote)
	}
}

// URI: /api/v2/notes/notes
// Method: POST
// Content-Type: application/json
// Body: { "txid": "<Transaction ID>", "notes": "<Notes>" }
// Response:
//      422 - wrong parameters
//      400 - internal server error
//      200 - ok, note removed by TxId
func removeNoteHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			wh.Error405(w)
			return
		}

		var note notes.Note
		if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
			wh.Error422(w, err.Error())
			return
		}

		if _, err := cipher.SHA256FromHex(note.TxIDHex); err == nil {
			if err := gateway.RemoveNote(note.TxIDHex); err != nil {
				wh.Error400(w, err.Error())
				return
			}
		} else {
			wh.Error422(w, fmt.Errorf("Wrong 'txid'").Error())
			return
		}

		wh.SendJSONOr500(logger, w, gateway.GetNoteByTxID(note.TxIDHex))
	}
}
