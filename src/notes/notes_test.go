package notes

import (
	"testing"
	"math/rand"
	"time"
	"github.com/stretchr/testify/assert"
	"os"
	"github.com/skycoin/skycoin/src/cipher"
)

var (
	specialChars = []string{"É", "G", "É", "ì", "É", "R", "Å", "[", "É", "f", "É", "B", "É", "ì", "É", "O", "Ç", "Õ", "ì", "Ô", "Ç", "µ", "Ç", "≠", "Ç", "»", "Ç", "¢", "*", "'", "Ü", "Ä", "Ö", "!", "§", "$", "%", "&", "/", ">", "<", "(", ")", "=", "?", "}", "{"}
	noteServ     *Service
	noteCFG      = Config{
		NotesPath: "transactionnotes_temp.json",
	}
)

const (
	totalNotes = 60
	noteLength = 20
	txIdByteLength = 32
)

func TestMain(m *testing.M) {
	defer teardown(1)

	m.Run()
}

func TestNewService(t *testing.T) {
	var err error

	noteServ, err = NewService(noteCFG)

	if err != nil {
		t.Error(err)
	}

	assert.NotNil(t, noteServ)
}

func teardown(i int) {
	defer os.Exit(i)

	err := os.Remove(noteCFG.NotesPath)

	if err != nil {
		panic(err)
		return
	}
}

func TestAddNotes(t *testing.T) {
	for i := 0; i < totalNotes; i++ {
		key := make([]byte, txIdByteLength)

		_, err := rand.Read(key)
		if err != nil {
			t.Error(err)
		}

		sha, err := cipher.SHA256FromBytes(key)

		if err != nil {
			t.Error(err)
			return
		}

		rndmTxIdHex := sha.Hex()
		rndmNotes := getRndmId()

		note := Note{}
		note.TxIdHex = rndmTxIdHex
		note.Notes = rndmNotes

		noteServ.Add(note)
	}

	assert.True(t, len(noteServ.GetAll()) == totalNotes)
}

func getRndmId() string {
	var rndTxId = ""
	for i := 0; i < noteLength; i++ {
		rand.Seed(time.Now().UTC().UnixNano())
		rndTxId += specialChars[rand.Intn(len(specialChars))]
	}

	return rndTxId
}

func TestGetNoteById(t *testing.T) {
	var testNotes []Note
	allNotes := noteServ.GetAll()

	// Get every 2nd Note for testing
	for i := 0; i < len(allNotes); i += 2 {
		testNotes = append(testNotes, allNotes[i])
	}

	// check
	for i := 0; i < len(testNotes); i++ {
		txId := testNotes[i].TxIdHex
		note := noteServ.GetByTransId(txId)

		assert.NotNil(t, note)

		assert.NotEmpty(t, note.Notes)
		assert.True(t, len(note.Notes) > 0)

		assert.NotEmpty(t, note.TxIdHex)
		assert.True(t, len(note.TxIdHex) > 0)
	}

	// Test not existent note
	noteByTransId := noteServ.GetByTransId("TRANSAKTIONSID")

	assert.Empty(t, noteByTransId.TxIdHex)
	assert.Empty(t, noteByTransId.Notes)
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
		txId := notesToRem[i].TxIdHex

		err := noteServ.RemoveByTxId(txId)

		assert.Nil(t, err)
	}

	// Test remove non existent Note
	err := noteServ.RemoveByTxId("NotExistentTxid")

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
		err := noteServ.Add(noteToOverwrite)

		assert.Nil(t, err)

		// Check if Note with TransId has changed
		checkNote := noteServ.GetByTransId(noteToOverwrite.TxIdHex)

		assert.NotNil(t, checkNote)
		assert.True(t, checkNote.Notes == noteToOverwrite.Notes)
	}

	assert.NotNil(t, testNotes)
}
