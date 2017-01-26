package wallet

import (
	"io/ioutil"
	"os"
	"time"
	"encoding/hex"
	"path/filepath"
	"strings"
	"fmt"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/util"
)

const NotesExtension = "nts"

type Notes []Note

type Note struct {
	TransactionId string
	Value string
}

type ReadableNotes []ReadableNote


type ReadableNote struct {
	TransactionId string `json:"transaction_id"`
	ActualNote string `json:"note_val"`
}


//check for collisions and retry if failure
func NewNotesFilename() string {
	timestamp := time.Now().Format(WalletTimestampFormat)
	//should read in wallet files and make sure does not exist
	padding := hex.EncodeToString((cipher.RandByte(2)))
	return fmt.Sprintf("%s_%s.%s", timestamp, padding, NotesExtension)
}

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
			return w,nil
		}
	}
	return wallets, nil
}

func LoadReadableNotes(filename string) (*ReadableNotes, error) {
	w := &ReadableNotes{}
	err := w.Load(filename)
	return w, err
}

func (self *ReadableNotes) Load(filename string) error {
	return util.LoadJSON(filename, self)
}

// Loads from filename
func (self ReadableNotes) ToNotes() ([]Note, error)  {
	notes := make([]Note, len(self))
	for i, e := range self {
		notes[i] = Note {
			TransactionId:e.TransactionId,
			Value:e.ActualNote,
		}
	}
	return notes,nil
}

func NewReadableNote(note Note) ReadableNote{
	return ReadableNote{
		TransactionId:note.TransactionId,
		ActualNote:note.Value,
	}
}

func NewReadableNotesFromNotes(w Notes) ReadableNotes {
	readable := make(ReadableNotes, len(w))
	i := 0
	for _, e := range w {
		readable[i] = NewReadableNote(e)
		i++
	}
	return readable
}

func (notes *Notes) Save(dir string, fileName string) error {
	r := notes.ToReadable()
	return r.Save(filepath.Join(dir, fileName))
}



func (notes *Notes) SaveNote(dir string, note Note) error {
	newNotes :=make([]Note,len(*notes)+1)
	for i,e :=range *notes{
		newNotes[i] = e
		i++
	}
	newNotes[len(*notes)] = note

	*notes = newNotes

	readableNotesToBeSaved:=NewReadableNotesFromNotes(newNotes)
	fileName, error := getNoteFileName(dir)
	if error !=nil{
		return error
	}
	readableNotesToBeSaved.Save(fileName)
	return nil
}

func getNoteFileName(dir string) (string, error){
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
			fullPath :=filepath.Join(dir, name)
			return fullPath,nil
		}
	}
	return "",nil
}

func (notes Notes) ToReadable() ReadableNotes {
	return NewReadableNotesFromNotes(notes)
}

// Saves to filename
func (self *ReadableNotes) Save(filename string) error {
	return util.SaveJSON(filename, self, 0600)
}

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
			return true,nil
		}
	}
	return false,nil
}

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