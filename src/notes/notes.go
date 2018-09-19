package notes

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"github.com/skycoin/skycoin/src/util/logging"
	"fmt"
)

type Note struct {
	TxIdHex string `json:"txid"`
	Notes   string `json:"notes"`
}

var (
	gNotes []Note
	log    = logging.MustGetLogger("notes")
	gNotesPath string
)

// Get all notes
func GetAll() []Note {
	return gNotes
}

// If note wasn't found by Id -> return empty Note
func GetByTransId(txId string) Note {
	for i := 0; i < len(gNotes); i++ {

		if note := gNotes[i]; note.TxIdHex == txId {
			return note
		}
	}

	return Note{}
}

// Add Note, if Note already exists, the old one will be overwritten
func Add(note Note) error {
	if !isNoteExist(note.TxIdHex) {
		log.Info("Adding Note with txid=" + note.TxIdHex)

		gNotes = append(gNotes, note)
	} else {
		log.Info("Overwriting Note with txid=" + note.TxIdHex)

		for i := 0; i < len(gNotes); i++ {
			if gNotes[i].TxIdHex == note.TxIdHex {
				gNotes[i] = note
			}
		}
	}

	return writeJson()
}

// Remove Note by txId
func Remove(txId string) error {
	for i := 0; i < len(gNotes); i++ {
		if gNotes[i].TxIdHex == txId {
			gNotes = append(gNotes[:i], gNotes[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Note with txid='" + txId + "' has not been removed")
}

// Check if Note with txId exists
func isNoteExist(txId string) bool {
	for i := 0; i < len(gNotes); i++ {
		if gNotes[i].TxIdHex == txId {
			return true
		}
	}
	return false
}

// Write Notes to configured gNotesPath
func writeJson() error {
	notesJson, _ := json.Marshal(gNotes)

	return ioutil.WriteFile(gNotesPath, notesJson, 0644)
}

// Read Notes from configured gNotesPath
func loadJson(notesPath string) {
	var notes []Note

	// Set Path for transactionnotes file
	gNotesPath = notesPath

	// Open jsonFile
	jsonFile, err := os.Open(notesPath)
	if err != nil {

		if os.IsExist(err) {
			log.Error(err)
			return
		} else {
			fmt.Print("File does not Exist: " + notesPath + "; Creating empty File...")
			err = ioutil.WriteFile(notesPath, nil, 0644)

			if err != nil {
				log.Error(err)
				return
			}
		}
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)

	if len(byteValue) > 0 {
		err = json.Unmarshal(byteValue, &notes)

		if err != nil {
			log.Error(err)
			return
		}

		gNotes = notes

		var fi os.FileInfo
		fi, err = jsonFile.Stat()

		if err != nil {
			log.Error(err)
			return
		}

		log.Info("Loaded Notes from " + fi.Name())
	} else {
		log.Info("Failed to load Notes: File is empty")
	}
}
