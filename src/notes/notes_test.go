package notes

import (
	"crypto/rand"
	"encoding/base64"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/skycoin/skycoin/src/cipher"
)

var (
	noteServ *Service
	noteCFG  = Config{
		NotesPath: "transactionnotes_temp.json",
	}
)

const (
	totalNotes     = 60
	txIDByteLength = 32
)

func TestMain(m *testing.M) {
	os.Exit(teardown(m.Run()))
}

func TestNewService(t *testing.T) {
	var err error

	noteServ, err = NewService(noteCFG)
	if err != nil {
		t.Error(err)
	}

	assert.NotNil(t, noteServ)
}

func teardown(i int) int {
	fi, err := os.Stat(noteCFG.NotesPath)

	if err != nil {
		panic(err)
	}

	if fi.Size() > 0 {

		if fi != nil {

			if err := os.Remove(noteCFG.NotesPath); err != nil {
				panic(err)
			}
		}
	}

	return i
}

func TestAddNotes(t *testing.T) {
	beforeAddCount := len(noteServ.GetAll())

	for i := 0; i < totalNotes; i++ {
		key, err := generateRandomBytes(txIDByteLength)
		if err != nil {
			t.Error(err)
			return
		}

		sha, err := cipher.SHA256FromBytes(key)
		if err != nil {
			t.Error(err)
			return
		}

		rndmTxIDHex := sha.Hex()
		rndmNotes, err := getRndmID(txIDByteLength)

		assert.NoError(t, err)

		note := Note{}
		note.TxIDHex = rndmTxIDHex
		note.Notes = rndmNotes

		note, err = noteServ.Add(note)

		if err != nil {
			t.Error(err)
		}
	}

	allNotesCount := len(noteServ.GetAll())

	assert.True(t, allNotesCount == (totalNotes+beforeAddCount))
}

func getRndmID(s int) (string, error) {
	b, err := generateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func TestGetNoteByID(t *testing.T) {
	var testNotes []Note
	allNotes := noteServ.GetAll()

	// Get every 2nd Note for testing
	for i := 0; i < len(allNotes); i += 2 {
		testNotes = append(testNotes, allNotes[i])
	}

	// check
	for i := 0; i < len(testNotes); i++ {
		txID := testNotes[i].TxIDHex
		note := noteServ.GetByTxID(txID)

		assert.NotNil(t, note)

		assert.NotEmpty(t, note.Notes)
		assert.True(t, len(note.Notes) > 0)

		assert.NotEmpty(t, note.TxIDHex)
		assert.True(t, len(note.TxIDHex) > 0)
	}

	// Test not existent note
	noteByTransID := noteServ.GetByTxID("TRANSAKTIONSID")

	assert.Empty(t, noteByTransID.TxIDHex)
	assert.Empty(t, noteByTransID.Notes)
}

func TestRemoveNote(t *testing.T) {
	allNotes := noteServ.GetAll()

	var notesToRem []Note

	// Get every 2nd Note for testing
	for i := 0; i < len(allNotes); i += 4 {
		notesToRem = append(notesToRem, allNotes[i])
	}

	// Remove Notes
	for i := 0; i < len(notesToRem); i++ {
		txID := notesToRem[i].TxIDHex

		err := noteServ.RemoveByTxID(txID)

		assert.Nil(t, err)
	}

	// Test remove non existent Note
	err := noteServ.RemoveByTxID("NotExistentTxID")

	assert.Error(t, err)
}

func TestOverwriteNotes(t *testing.T) {
	allNotes := noteServ.GetAll()

	var testNotes []Note

	// Get every 2nd Note for testing
	for i := 0; i < len(allNotes); i += 10 {
		testNotes = append(testNotes, allNotes[i])
	}

	for i := 0; i < len(testNotes); i++ {
		// Cache Note
		noteToOverwrite := testNotes[i]

		// Modify Note
		noteToOverwrite.Notes = "New Note"
		_, err := noteServ.Add(noteToOverwrite)

		assert.Nil(t, err)

		// Check if Note with TransId has changed
		checkNote := noteServ.GetByTxID(noteToOverwrite.TxIDHex)

		assert.NotNil(t, checkNote)
		assert.True(t, checkNote.Notes == noteToOverwrite.Notes)
	}

	assert.NotNil(t, testNotes)
}
