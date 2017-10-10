package visor

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
)

const (
	blockchainPubkeyStr = "0328c576d3f420e7682058a981173a4b374c7cc5ff55bf394d3cf57059bbe6456a"
)

func readAll(t *testing.T, f string) []byte {
	fi, err := os.Open(f)
	require.NoError(t, err)
	defer fi.Close()

	b, err := ioutil.ReadAll(fi)
	require.NoError(t, err)

	return b
}

func mustParsePubkey(t *testing.T) cipher.PubKey {
	// Parse the blockchain pubkey associated with this corrupted test db
	t.Helper()
	pubkey, err := cipher.PubKeyFromHex(blockchainPubkeyStr)
	require.NoError(t, err)
	return pubkey
}

func writeDBFile(t *testing.T, badDBFile string, badDBData []byte) {
	t.Logf("Writing the original bad db file back to %s", badDBFile)
	fi, err := os.OpenFile(badDBFile, os.O_WRONLY, 0600)
	require.NoError(t, err)
	defer fi.Close()

	_, err = io.Copy(fi, bytes.NewBuffer(badDBData))
	require.NoError(t, err)
}

func findCorruptDBFiles(t *testing.T, badDBFile string) []string {
	corruptFiles, err := filepath.Glob(badDBFile + ".corrupt.*")
	require.NoError(t, err)
	return corruptFiles
}

func removeCorruptDBFiles(t *testing.T, badDBFile string) {
	corruptFiles := findCorruptDBFiles(t, badDBFile)
	for _, m := range corruptFiles {
		err := os.Remove(m)
		require.NoError(t, err)
	}
}

func TestErrSignatureLostRecreateDB(t *testing.T) {
	badDBFile := "./testdata/data.db.nosig" // about 8MB size
	badDBData := readAll(t, badDBFile)

	pubkey := mustParsePubkey(t)

	// Remove any existing corrupt db files from testdata
	removeCorruptDBFiles(t, badDBFile)
	corruptFiles := findCorruptDBFiles(t, badDBFile)
	require.Len(t, corruptFiles, 0)

	// Cleanup
	defer func() {
		// Write the bad db data back to badDBFile
		writeDBFile(t, badDBFile, badDBData)
		// Remove leftover corrupt db copies
		removeCorruptDBFiles(t, badDBFile)
	}()

	// Make sure that the database file causes ErrSignatureLost error
	t.Logf("Checking that %s is a corrupted database", badDBFile)
	func() {
		db, err := openDB(badDBFile)
		require.NoError(t, err)
		defer func() {
			err := db.Close()
			assert.NoError(t, err)
		}()

		_, err = NewBlockchain(db, pubkey, Arbitrating(false))
		require.Error(t, err)
		require.Contains(t, err.Error(), "find no signature of block:")
	}()

	// Loading this invalid db should cause load() to recreate the db
	t.Logf("Loading the corrupted db")
	db, bc, err := load(badDBFile, pubkey, false)
	require.NoError(t, err)

	err = db.Close()
	require.NoError(t, err)

	require.NotNil(t, db)
	require.NotNil(t, bc)

	// A corrupted database file should exist
	corruptFiles = findCorruptDBFiles(t, badDBFile)
	require.Len(t, corruptFiles, 1)

	// A new db should be written in place of the old bad db, and not be corrupted
	t.Logf("Checking that the new db file is valid")
	func() {
		db, err := openDB(badDBFile)
		require.NoError(t, err)
		defer func() {
			err := db.Close()
			assert.NoError(t, err)
		}()

		// The new db is not corrupted and loads without error
		bc, err := NewBlockchain(db, pubkey, Arbitrating(false))
		require.NoError(t, err)
		require.NotNil(t, bc)
	}()
}

func TestNormalLoadDBErr(t *testing.T) {
	// If load() returns an error other than ErrSignatureLost,
	// it should not recreate a db
	badDBFile := "./testdata/data.db.garbage"
	badDBData := readAll(t, badDBFile)

	pubkey := mustParsePubkey(t)

	// Cleanup
	defer func() {
		// Write the bad db data back to badDBFile
		writeDBFile(t, badDBFile, badDBData)
		// Remove leftover corrupt db copies
		removeCorruptDBFiles(t, badDBFile)
	}()

	db, bc, err := load(badDBFile, pubkey, false)
	require.Error(t, err)
	require.NotEqual(t, ErrSignatureLost, err)
	require.Nil(t, db)
	require.Nil(t, bc)

	// Garbage db should not have been modified
	badDBDataCopy := readAll(t, badDBFile)
	require.True(t, bytes.Equal(badDBData, badDBDataCopy))
}
