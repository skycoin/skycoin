package notes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/skycoin/skycoin/src/util/logging"
)

// Note struct
type Note struct {
	TxIDHex string `json:"txid"`
	Notes   string `json:"notes"`
}

var (
	gNotes     []Note
	log        = logging.MustGetLogger("notes")
	gNotesPath string
)

// GetAll returns all saved Notes
func GetAll() []Note {
	return gNotes
}

// GetByTxID If note wasn't found by ID -> return empty Note
func GetByTxID(txID string) Note {
	for i := 0; i < len(gNotes); i++ {

		if note := gNotes[i]; note.TxIDHex == txID {
			return note
		}
	}

	return Note{}
}

// Add Note, if Note already exists, the old one will be overwritten
func Add(note Note) (Note, error) {
	if !isNoteExist(note.TxIDHex) {
		log.Infof("Adding Note with txID=&v", note.TxIDHex)

		gNotes = append(gNotes, note)
	} else {
		log.Infof("Overwriting Note with txID=%v", note.TxIDHex)

		for i := 0; i < len(gNotes); i++ {
			if gNotes[i].TxIDHex == note.TxIDHex {
				gNotes[i] = note
			}
		}
	}

	if err := writeJSON(); err != nil {
		return Note{}, err
	}

	return note, nil
}

// Remove Note by txID
func Remove(txID string) error {
	log.Infof("Removing note with txID=%v", txID)

	for i := 0; i < len(gNotes); i++ {

		if gNotes[i].TxIDHex == txID {
			gNotes = append(gNotes[:i], gNotes[i+1:]...)
			return nil
		}
	}

	if err := writeJSON(); err != nil {
		return err
	}

	return fmt.Errorf("note with txID='%v' has not been removed: Note doesn't exist", txID)
}

// Check if Note with given txID exists
func isNoteExist(txID string) bool {
	for i := 0; i < len(gNotes); i++ {
		if gNotes[i].TxIDHex == txID {
			return true
		}
	}
	return false
}

// Write Notes to configured gNotesPath
func writeJSON() error {
	notesJSON, err := json.Marshal(gNotes)

	if err != nil {
		log.Error(err)
		return err
	}

	return ioutil.WriteFile(gNotesPath, notesJSON, 0644)
}

// Read Notes from configured gNotesPath
func loadJSON(notesPath string) {
	var notes []Note

	// Set Path for notes file
	gNotesPath = notesPath

	// Open jsonFile
	jsonFile, err := os.Open(notesPath)
	if err != nil {

		if os.IsExist(err) {
			log.Error(err)
			return
		}

		log.Infof("File does not exist: %v; Creating empty File...", notesPath)

		err = ioutil.WriteFile(notesPath, []byte{}, 0644)
		if err != nil {
			log.Error(err)
			return
		}

		jsonFile, err = os.Open(notesPath)
		if err != nil {
			log.Error(err)
			return
		}
	}

	// When give json file doesn't exist, create an empty one
	if jsonFile != nil {
		var fi os.FileInfo
		var byteValue []byte

		fi, err = jsonFile.Stat()
		if err != nil {
			log.Error(err)
			return
		}

		if fi.Size() > 0 {

			byteValue, err = ioutil.ReadAll(jsonFile)
			if err != nil {
				log.Error(err)
				return
			}

		} else {
			log.Info("Failed to load Notes: File is empty")
			return
		}

		if len(byteValue) > 0 {

			err = json.Unmarshal(byteValue, &notes)
			if err != nil {
				log.Error(err)
				return
			}

			gNotes = notes

			log.Infof("Loaded Notes from %v", fi.Name())
		}
	}
}
