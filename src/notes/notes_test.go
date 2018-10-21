package notes

import (
	"crypto/rand"
	"encoding/base64"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
)

var (
	noteServ *Service
	noteCFG = Config{
		NotesPath:  "./transactionnotes_temp.json",
		NotesStore: NewStore(),
	}
)

const (
	totalNotes     = 60
	txIDByteLength = 32
)

func TestMain(m *testing.M) {
	os.Exit(teardown(m.Run()))
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

func TestNewService(t *testing.T) {
	var err error

	noteServ, err = NewService(noteCFG)
	require.NoError(t, err)
	require.NotNil(t, noteServ)
}

func TestAddNotes(t *testing.T) {
	beforeAddCount := len(noteServ.GetAll())

	for i := 0; i < totalNotes; i++ {
		key, err := generateRandomBytes(txIDByteLength)
		require.NoError(t, err)

		sha, err := cipher.SHA256FromBytes(key)
		require.NoError(t, err)

		rndmTxIDHex := sha.Hex()
		rndmNotes, err := getRndmID(txIDByteLength)
		require.NoError(t, err)

		note := Note{}
		note.TxIDHex = rndmTxIDHex
		note.Notes = rndmNotes

		note, err = noteServ.Add(note)
		require.NoError(t, err)
	}

	allNotesCount := len(noteServ.GetAll())
	require.Equal(t, totalNotes+beforeAddCount, allNotesCount)
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

	return b, err
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

		require.NotNil(t, note)

		require.NotEmpty(t, note.Notes)
		require.True(t, len(note.Notes) > 0)

		require.NotEmpty(t, note.TxIDHex)
		require.True(t, len(note.TxIDHex) > 0)
	}

	// Test not existent note
	noteByTransID := noteServ.GetByTxID("TRANSAKTIONSID")

	require.Empty(t, noteByTransID.TxIDHex)
	require.Empty(t, noteByTransID.Notes)
}

func TestRemoveNotes(t *testing.T) {
	allNotes := noteServ.GetAll()

	var notesToRem []Note

	// Get every 2nd Note to test remove mechanism
	for i := 0; i < len(allNotes); i += 4 {
		notesToRem = append(notesToRem, allNotes[i])
	}

	for i := 0; i < len(notesToRem); i++ {
		txID := notesToRem[i].TxIDHex

		err := noteServ.RemoveByTxID(txID)
		require.NoError(t, err)
	}

	// Test remove non existent Note
	err := noteServ.RemoveByTxID("NotExistentTxID")
	require.Error(t, err)
}

func TestOverwriteNotes(t *testing.T) {
	allNotes := noteServ.GetAll()

	log.Info(len(allNotes))

	var testNotes []Note

	// Get every 2nd Note for testing
	for i := 0; i < len(allNotes); i += 10 {
		log.Info(allNotes[i].Notes)
		testNotes = append(testNotes, allNotes[i])
	}

	for i := 0; i < len(testNotes); i++ {
		// Cache Note
		noteToOverwrite := testNotes[i]

		// Modify Note
		noteToOverwrite.Notes = "New Note"
		_, err := noteServ.Add(noteToOverwrite)

		require.Nil(t, err)

		// Check if Note with TxID has changed
		checkNote := noteServ.GetByTxID(noteToOverwrite.TxIDHex)
		require.NotNil(t, checkNote)
		require.Equal(t, checkNote.Notes, noteToOverwrite.Notes)
		require.NotEqual(t, checkNote.Notes, testNotes[i].Notes)
	}

	require.NotNil(t, testNotes)
}
