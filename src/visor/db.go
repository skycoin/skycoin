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
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/visor/blockdb"
	"github.com/skycoin/skycoin/src/visor/dbutil"
	"github.com/skycoin/skycoin/src/visor/historydb"
)

var (
	// BlockchainVerifyTheadNum number of goroutines to use for signature and historydb verification
	BlockchainVerifyTheadNum = 4
)

// ErrCorruptDB is returned if the database is corrupted
// The original corruption error is embedded
type ErrCorruptDB struct {
	error
}

// CheckDatabase checks the database for corruption
func CheckDatabase(db *dbutil.DB, pubkey cipher.PubKey, quit chan struct{}) error {
	var blocksBktExist bool
	db.View("CheckDatabase", func(tx *dbutil.Tx) error {
		blocksBktExist = dbutil.Exists(tx, blockdb.BlocksBkt)
		return nil
	})

	// Don't verify the db if the blocks bucket does not exist
	if !blocksBktExist {
		return nil
	}

	bc, err := NewBlockchain(db, BlockchainConfig{Pubkey: pubkey})
	if err != nil {
		return err
	}

	history := historydb.New()
	indexesMap := historydb.NewIndexesMap()
	verifyFunc := func(tx *dbutil.Tx, b *coin.SignedBlock) error {
		// Verify signature
		if err := bc.VerifySignature(b); err != nil {
			return err
		}

		// Verify historydb
		return history.Verify(tx, b, indexesMap)
	}

	err = bc.WalkChain(BlockchainVerifyTheadNum, verifyFunc, quit)
	switch err.(type) {
	case nil:
		return nil
	case blockdb.ErrMissingSignature, historydb.ErrHistoryDBCorrupted:
		return ErrCorruptDB{err}
	default:
		return err
	}
}

// ResetCorruptDB checks the database for corruption and if corrupted, then it erases the db and starts over.
// A copy of the corrupted database is saved.
func ResetCorruptDB(db *dbutil.DB, pubkey cipher.PubKey, quit chan struct{}) (*dbutil.DB, error) {
	err := CheckDatabase(db, pubkey, quit)

	switch err.(type) {
	case nil:
		return db, nil
	case ErrCorruptDB:
		logger.Critical().Errorf("Database is corrupted, recreating db: %v", err)
		return handleCorruptDB(db)
	default:
		return nil, err
	}
}

// handleCorruptDB recreates the DB, making a backup copy marked as corrupted
func handleCorruptDB(db *dbutil.DB) (*dbutil.DB, error) {
	dbReadOnly := db.IsReadOnly()
	dbPath := db.Path()

	if err := db.Close(); err != nil {
		return nil, fmt.Errorf("Failed to close db: %v", err)
	}

	corruptDBPath, err := moveCorruptDB(dbPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to copy corrupted db: %v", err)
	}

	logger.Critical().Infof("Moved corrupted db to %s", corruptDBPath)

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
		logger.Errorf("os.Rename(%s, %s) failed: %v", dbPath, newDBPath, err)
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
