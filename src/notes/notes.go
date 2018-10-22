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

// Store for caching notes
type Store struct {
	notes     []Note
	notesPath string
}

var log = logging.MustGetLogger("notes")

// NewStore returns an instance of a Notes Store
func NewStore() *Store {
	return &Store{
		notes: make([]Note, 0),
	}
}

// GetAll returns all saved Notes
func (c *Store) GetAll() []Note {
	return c.notes
}

// GetByTxID If note wasn't found by ID -> return empty Note
func (c *Store) GetByTxID(txID string) Note {
	for i := 0; i < len(c.notes); i++ {

		if note := c.notes[i]; note.TxIDHex == txID {
			return note
		}
	}

	return Note{}
}

// Add Note, if Note already exists, the old one will be overwritten
func (c *Store) Add(note Note) (Note, error) {
	if !c.isNoteExist(note.TxIDHex) {
		log.Infof("Adding Note with txID=&v", note.TxIDHex)
		c.notes = append(c.notes, note)
	} else {
		log.Infof("Overwriting Note with txID=%v", note.TxIDHex)

		for i := 0; i < len(c.notes); i++ {
			if c.notes[i].TxIDHex == note.TxIDHex {
				c.notes[i] = note
			}
		}
	}

	if err := c.writeJSON(); err != nil {
		return Note{}, err
	}

	return note, nil
}

// Remove Note by txID
func (c *Store) Remove(txID string) error {
	log.Infof("Removing note with txID=%v", txID)

	for i := 0; i < len(c.notes); i++ {

		if c.notes[i].TxIDHex == txID {
			c.notes = append(c.notes[:i], c.notes[i+1:]...)

			if err := c.writeJSON(); err != nil {
				return err
			}
			return nil
		}
	}

	return fmt.Errorf("note with txID='%v' has not been removed: Note doesn't exist", txID)
}

// Check if Note with given txID exists
func (c *Store) isNoteExist(txID string) bool {
	for i := 0; i < len(c.notes); i++ {

		if c.notes[i].TxIDHex == txID {
			return true
		}
	}
	return false
}

// Write Notes to configured notes path
func (c *Store) writeJSON() error {
	notesJSON, err := json.Marshal(c.notes)

	if err != nil {
		log.Error(err)
		return err
	}
	return ioutil.WriteFile(c.notesPath, notesJSON, 0644)
}

// LoadJSON loads Notes from configured path
func (c *Store) loadJSON() error {
	var notes []Note

	// Open jsonFile
	jsonFile, err := os.Open(c.notesPath)
	if err != nil {

		if os.IsExist(err) {
			log.Error(err)
			return err
		}

		log.Infof("File does not exist: %v; Creating empty File...", c.notesPath)

		err = ioutil.WriteFile(c.notesPath, []byte{}, 0644)
		if err != nil {
			log.Error(err)
			return err
		}

		jsonFile, err = os.Open(c.notesPath)
		if err != nil {
			log.Error(err)
			return err
		}
	}

	// When given json path doesn't exist, create an empty json file
	if jsonFile != nil {
		var fi os.FileInfo
		var byteValue []byte

		fi, err = jsonFile.Stat()
		if err != nil {
			log.Error(err)
			return err
		}

		if fi.Size() > 0 {

			byteValue, err = ioutil.ReadAll(jsonFile)
			if err != nil {
				log.Error(err)
				return err
			}

		} else {
			log.Info("Failed to load Notes: File is empty")
			return err
		}

		if len(byteValue) > 0 {

			err = json.Unmarshal(byteValue, &notes)
			if err != nil {
				log.Error(err)
				return err
			}

			c.notes = notes

			log.Infof("Loaded Notes from %v", fi.Name())
		}
	}
	return nil
}
