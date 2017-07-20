package wallet

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/util/file"
)

// NotesExtension file extension of notes
const NotesExtension = "nts"

// Notes array of notes
type Notes []Note

// Note note struct
type Note struct {
	TxID  string
	Value string
}

// ReadableNotes readable notes
type ReadableNotes []ReadableNote

// ReadableNote readable note struct
type ReadableNote struct {
	TransactionID string `json:"transaction_id"`
	ActualNote    string `json:"note_val"`
}

// NewNotesFilename check for collisions and retry if failure
func NewNotesFilename() string {
	timestamp := time.Now().Format(WalletTimestampFormat)
	//should read in wallet files and make sure does not exist
	padding := hex.EncodeToString((cipher.RandByte(2)))
	return fmt.Sprintf("%s_%s.%s", timestamp, padding, NotesExtension)
}

// LoadNotes loads notes from given dir
func LoadNotes(dir string) (Notes, error) {
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	bkpath := dir + "/backup/"
	if _, err := os.Stat(bkpath); os.IsNotExist(err) {
		// create the backup dir
		logger.Critical("create wallet backup dir, %v", bkpath)
		if err := os.Mkdir(bkpath, 0777); err != nil {
			return nil, err
		}
	}

	//have := make(map[WalletID]Wallet, len(entries))
	wallets := make(Notes, 0)
	for _, e := range entries {
		if e.Mode().IsRegular() {
			name := e.Name()
			if !strings.HasSuffix(name, NotesExtension) {
				continue
			}
			fullpath := filepath.Join(dir, name)
			rw, err := LoadReadableNotes(fullpath)
			if err != nil {
				return nil, err
			}
			w, err := rw.ToNotes()
			if err != nil {
				return nil, err
			}
			return w, nil
		}
	}
	return wallets, nil
}

// LoadReadableNotes loads readable notes from given file
func LoadReadableNotes(filename string) (*ReadableNotes, error) {
	w := &ReadableNotes{}
	err := w.Load(filename)
	return w, err
}

// Load loads readable notes from given file
func (rns *ReadableNotes) Load(filename string) error {
	return file.LoadJSON(filename, rns)
}

// ToNotes converts from readable notes to Notes
func (rns ReadableNotes) ToNotes() ([]Note, error) {
	notes := make([]Note, len(rns))
	for i, e := range rns {
		notes[i] = Note{
			TxID:  e.TransactionID,
			Value: e.ActualNote,
		}
	}
	return notes, nil
}

// Save persists readable notes to disk
func (rns *ReadableNotes) Save(filename string) error {
	return file.SaveJSON(filename, rns, 0600)
}

// NewReadableNote creates readable note
func NewReadableNote(note Note) ReadableNote {
	return ReadableNote{
		TransactionID: note.TxID,
		ActualNote:    note.Value,
	}
}

// NewReadableNotesFromNotes creates readable notes from notes
func NewReadableNotesFromNotes(w Notes) ReadableNotes {
	readable := make(ReadableNotes, len(w))
	i := 0
	for _, e := range w {
		readable[i] = NewReadableNote(e)
		i++
	}
	return readable
}

// Save persists notes to disk
func (notes *Notes) Save(dir string, fileName string) error {
	r := notes.ToReadable()
	return r.Save(filepath.Join(dir, fileName))
}

// SaveNote save new note
func (notes *Notes) SaveNote(dir string, note Note) error {
	newNotes := make([]Note, len(*notes)+1)
	for i, e := range *notes {
		newNotes[i] = e
		i++
	}
	newNotes[len(*notes)] = note

	*notes = newNotes

	readableNotesToBeSaved := NewReadableNotesFromNotes(newNotes)
	fileName, error := getNoteFileName(dir)
	if error != nil {
		return error
	}
	readableNotesToBeSaved.Save(fileName)
	return nil
}

func getNoteFileName(dir string) (string, error) {
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		return "", err
	}
	for _, e := range entries {
		if e.Mode().IsRegular() {
			name := e.Name()
			if !strings.HasSuffix(name, NotesExtension) {
				continue
			}
			fullPath := filepath.Join(dir, name)
			return fullPath, nil
		}
	}
	return "", nil
}

// ToReadable converts Notes to readable notes
func (notes Notes) ToReadable() ReadableNotes {
	return NewReadableNotesFromNotes(notes)
}

// NotesFileExist checks if there're notes exist
func NotesFileExist(dir string) (bool, error) {
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		return false, err
	}
	for _, e := range entries {
		if e.Mode().IsRegular() {
			name := e.Name()
			if !strings.HasSuffix(name, NotesExtension) {
				continue
			}
			return true, nil
		}
	}
	return false, nil
}

// CreateNoteFileIfNotExist creates note file if not exist
func CreateNoteFileIfNotExist(dir string) {
	exist, err := NotesFileExist(dir)
	if err != nil {
		return
	}
	if exist == false {
		noteFileName := NewNotesFilename()
		dummyNotes := ReadableNotes{}
		fullpath := filepath.Join(dir, noteFileName)
		dummyNotes.Save(fullpath)
	}

}
