package visor

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/boltdb/bolt"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/visor/blockdb"
	"github.com/skycoin/skycoin/src/visor/dbutil"
)

var (
	// SigVerifyTheadNum number of goroutines to use for signature verification
	SigVerifyTheadNum = 4
)

// CheckAndRepairDatabase loads blockchain from DB and if any error occurs then delete
// the db and create an empty blockchain.
func CheckAndRepairDatabase(db *dbutil.DB, pubkey cipher.PubKey) (*dbutil.DB, error) {
	logger.Info("Loading blockchain")

	bc, err := NewBlockchain(db, BlockchainConfig{
		Pubkey: pubkey,
	})
	if err != nil {
		return nil, err
	}

	err = db.View("VerifySignatures", func(tx *dbutil.Tx) error {
		return bc.VerifySignatures(tx, SigVerifyTheadNum)
	})
	if err != nil {
		switch err.(type) {
		case blockdb.ErrMissingSignature:
			// Recreate the block database if ErrMissingSignature occurs
			logger.Critical().Errorf("Block database signature missing, recreating db: %v", err)
			return handleCorruptedDB(db)
		default:
			return nil, err
		}
	}

	return db, nil
}

// handleCorruptedDB recreates the DB, making a backup copy marked as corrupted
func handleCorruptedDB(db *dbutil.DB) (*dbutil.DB, error) {
	dbReadOnly := db.IsReadOnly()
	dbPath := db.Path()

	if err := db.Close(); err != nil {
		return nil, fmt.Errorf("Failed to close db: %v", err)
	}

	corruptDBPath, err := moveCorruptDB(dbPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to copy corrupted db: %v", err)
	}

	logger.Critical().Errorf("Moved corrupted db to %s", corruptDBPath)

	return OpenDB(dbPath, dbReadOnly)
}

// OpenDB opens the blockdb
func OpenDB(dbFile string, readOnly bool) (*dbutil.DB, error) {
	db, err := bolt.Open(dbFile, 0600, &bolt.Options{
		Timeout:  500 * time.Millisecond,
		ReadOnly: readOnly,
	})
	if err != nil {
		return nil, fmt.Errorf("Open boltdb failed, %v", err)
	}

	return dbutil.WrapDB(db), nil
}

// moveCorruptDB moves a file to makeCorruptDBPath(dbPath)
func moveCorruptDB(dbPath string) (string, error) {
	newDBPath, err := makeCorruptDBPath(dbPath)
	if err != nil {
		return "", err
	}

	if err := os.Rename(dbPath, newDBPath); err != nil {
		logger.Infof("os.Rename(%s, %s) failed: %v", dbPath, newDBPath, err)
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
