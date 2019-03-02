package note

import (
	"errors"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/util/file"
)

// Error wraps note-related errors.
// It wraps errors caused by user input, but not errors caused by programmer input or internal issues.
type Error struct {
	error
}

// NewError creates an Error
func NewError(err error) error {
	if err == nil {
		return nil
	}
	return Error{err}
}

var (
	// ErrNoteAPIDisabled is returned when trying to do note actions while the EnableNoteAPI option is false
	ErrNoteAPIDisabled = NewError(errors.New("wallet api is disabled"))
	// ErrNoteNotExist is returned if a note does not exist
	ErrNoteNotExist = NewError(errors.New("note doesn't exist"))
	// ErrInvalidTxID is returned if a transaction ID is not valud
	ErrInvalidTxID = NewError(errors.New("invalid transaction ID"))
)

const (
	// NoteExt file extension
	NoteExt = ".txnote"
)

// LoadNotes Loads all notes contained in notes dir. If any regular file in notes
// dir fails to load, loading is aborted and error returned. Only files with
// extension NoteExt are considered.
func LoadNotes(dir string) (map[string]string, error) {
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	notes := make(map[string]string)
	for _, e := range entries {
		if e.Mode().IsRegular() {
			fileName := e.Name()
			if !strings.HasSuffix(fileName, NoteExt) {
				continue
			}

			path := filepath.Join(dir, fileName)
			n, err := loadNote(path)
			if err != nil {
				return nil, err
			}

			txID := strings.Replace(fileName, NoteExt, "", -1)
			notes[txID] = n
		}
	}
	return notes, nil
}

func loadNote(path string) (string, error) {
	data, err := file.LoadBinary(path)
	return string(data), err
}

// Save saves the note to the given dir
func Save(dir, txID, note string) error {
	fileName := txID + NoteExt
	path := filepath.Join(dir, fileName)
	return file.RewriteBinary(path, []byte(note), 0600)
}

// Remove removes the note from the given dir
func Remove(dir, txID string) error {
	fileName := txID + NoteExt
	path := filepath.Join(dir, fileName)
	return file.RemoveFile(path)
}

func validateTxID(txID string) error {
	_, err := cipher.SHA256FromHex(txID)
	return err
}
