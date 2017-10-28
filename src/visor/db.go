package visor

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	"github.com/skycoin/skycoin/src/cipher"
)

// loadBlockchain loads blockchain from DB and if any error occurs then delete
// the db and create an empty blockchain.
func loadBlockchain(dbPath string, pubkey cipher.PubKey, arbitrating bool) (*bolt.DB, *Blockchain, error) {
	// creates blockchain instance
	db, err := openDB(dbPath)
	if err != nil {
		return nil, nil, err
	}

	bc, err := NewBlockchain(db, pubkey, Arbitrating(arbitrating))

	if err == nil {
		return db, bc, nil
	}

	if !strings.Contains(err.Error(), "find no signature of block") {
		return nil, nil, err
	}

	// Recreate the block database if ErrSignatureLost occurs
	logger.Critical("Block database signature missing, recreating db: %v", err)
	if err := db.Close(); err != nil {
		return nil, nil, fmt.Errorf("failed to close db: %v", err)
	}

	corruptDBPath, err := moveCorruptDB(dbPath)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to copy corrupted db: %v", err)
	}

	logger.Critical("Moved corrupted db to %s", corruptDBPath)

	db, err = openDB(dbPath)
	if err != nil {
		return nil, nil, err
	}

	bc, err = NewBlockchain(db, pubkey, Arbitrating(arbitrating))
	if err != nil {
		return nil, nil, err
	}

	return db, bc, nil
}

// open the blockdb.
func openDB(dbFile string) (*bolt.DB, error) {
	db, err := bolt.Open(dbFile, 0600, &bolt.Options{
		Timeout: 500 * time.Millisecond,
	})
	if err != nil {
		return nil, fmt.Errorf("Open boltdb failed, %v", err)
	}

	return db, nil
}

// moveCorruptDB moves a file to makeCorruptDBPath(dbPath)
func moveCorruptDB(dbPath string) (string, error) {
	newDBPath, err := makeCorruptDBPath(dbPath)
	if err != nil {
		return "", err
	}

	if err := os.Rename(dbPath, newDBPath); err != nil {
		return "", err
	}

	return newDBPath, nil
}

// makeCorruptDBPath creates a $FILE.corrupt.$HASH string based on dbPath,
// where $HASH is truncated SHA1 of $FILE.
func makeCorruptDBPath(dbPath string) (string, error) {
	dbFileHash, err := shaFileID(dbPath)
	if err != nil {
		return "", err
	}

	dbDir, dbFile := filepath.Split(dbPath)
	newDBFile := fmt.Sprintf("%s.corrupt.%s", dbFile, dbFileHash)
	newDBPath := filepath.Join(dbDir, newDBFile)

	return newDBPath, nil
}

// shaFileID return the first 8 bytes of the SHA1 hash of the file,
// base64-encoded
func shaFileID(dbPath string) (string, error) {
	fi, err := os.Open(dbPath)
	if err != nil {
		return "", err
	}
	defer fi.Close()

	h := sha1.New()
	if _, err := io.Copy(h, fi); err != nil {
		return "", err
	}

	sum := h.Sum(nil)
	encodedSum := base64.RawStdEncoding.EncodeToString(sum[:8])

	return encodedSum, nil
}
