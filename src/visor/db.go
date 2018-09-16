package visor

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
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

// CheckDatabase checks the database for corruption, rebuild history if corrupted
func CheckDatabase(db *dbutil.DB, pubkey cipher.PubKey, quit chan struct{}) error {
	var blocksBktExist bool
	if err := db.View("CheckDatabase", func(tx *dbutil.Tx) error {
		blocksBktExist = dbutil.Exists(tx, blockdb.BlocksBkt)
		return nil
	}); err != nil {
		return err
	}

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

	var historyVerifyErr error
	var lock sync.Mutex
	verifyFunc := func(tx *dbutil.Tx, b *coin.SignedBlock) error {
		// Verify signature
		if err := bc.VerifySignature(b); err != nil {
			return err
		}

		// Verify historydb, we don't return the error of history.Verify here,
		// as we have to check all signature, if we return error early here, the
		// potential bad signature won't be detected.
		lock.Lock()
		defer lock.Unlock()
		if historyVerifyErr == nil {
			historyVerifyErr = history.Verify(tx, b, indexesMap)
		}
		return nil
	}

	err = bc.WalkChain(BlockchainVerifyTheadNum, verifyFunc, quit)
	switch err.(type) {
	case nil:
		lock.Lock()
		err = historyVerifyErr
		lock.Unlock()
		return err
	default:
		return err
	}
}

// backup the corrypted db first, then rebuild the history DB.
func rebuildHistoryDB(db *dbutil.DB, history *historydb.HistoryDB, bc *Blockchain, quit chan struct{}) (*dbutil.DB, error) { // nolint: unused,megacheck
	db, err := backupDB(db)
	if err != nil {
		return nil, err
	}

	if err := db.Update("Rebuild history db", func(tx *dbutil.Tx) error {
		if err := history.Erase(tx); err != nil {
			return err
		}

		headSeq, ok, err := bc.HeadSeq(tx)
		if err != nil {
			return err
		}

		if !ok {
			return errors.New("head block does not exist")
		}

		for i := uint64(0); i <= headSeq; i++ {
			select {
			case <-quit:
				return nil
			default:
				b, err := bc.GetSignedBlockBySeq(tx, i)
				if err != nil {
					return err
				}

				if err := history.ParseBlock(tx, b.Block); err != nil {
					return err
				}

				if i%1000 == 0 {
					logger.Critical().Infof("Parse block: %d", i)
				}
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return db, nil
}

// backupDB makes a backup copy of the DB
func backupDB(db *dbutil.DB) (*dbutil.DB, error) { // nolint: unused,megacheck
	// backup the corrupted database
	dbReadOnly := db.IsReadOnly()

	dbPath := db.Path()

	if err := db.Close(); err != nil {
		return nil, fmt.Errorf("Failed to close db: %v", err)
	}

	corruptDBPath, err := copyCorruptDB(dbPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to copy corrupted db: %v", err)
	}

	logger.Critical().Infof("Copy corrupted db to %s", corruptDBPath)

	// Open the database again
	return OpenDB(dbPath, dbReadOnly)
}

// RepairCorruptDB checks the database for corruption and if corrupted and
// is ErrMissingSignature, then then it erases the db and starts over.
// If it's ErrHistoryDBCorrupted, then rebuild historydb from scratch.
// A copy of the corrupted database is saved.
func RepairCorruptDB(db *dbutil.DB, pubkey cipher.PubKey, quit chan struct{}) (*dbutil.DB, error) {
	err := CheckDatabase(db, pubkey, quit)
	switch err.(type) {
	case nil:
		return db, nil
	case blockdb.ErrMissingSignature,
		historydb.ErrHistoryDBCorrupted:
		logger.Critical().Errorf("Database is corrupted, recreating db: %v", err)
		return resetCorruptDB(db)
	default:
		return nil, err
	}
}

func rebuildCorruptDB(db *dbutil.DB, pubkey cipher.PubKey, quit chan struct{}) (*dbutil.DB, error) { //nolint: deadcode,unused,megacheck
	history := historydb.New()
	bc, err := NewBlockchain(db, BlockchainConfig{Pubkey: pubkey})
	if err != nil {
		return nil, err
	}

	return rebuildHistoryDB(db, history, bc, quit)
}

// resetCorruptDB recreates the DB, making a backup copy marked as corrupted
func resetCorruptDB(db *dbutil.DB) (*dbutil.DB, error) {
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
		Timeout:  5000 * time.Millisecond,
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

// copyCorruptDB copy a file to makeCorruptDBPath(dbPath)
func copyCorruptDB(dbPath string) (string, error) { // nolint: unused,megacheck
	newDBPath, err := makeCorruptDBPath(dbPath)
	if err != nil {
		return "", err
	}

	in, err := os.Open(dbPath)
	if err != nil {
		return "", err
	}
	defer in.Close()

	out, err := os.Create(newDBPath)
	if err != nil {
		return "", err
	}
	defer out.Close()
	logger.Critical().Info(out.Name())

	_, err = io.Copy(in, out)
	if err != nil {
		return "", err
	}

	if err := out.Close(); err != nil {
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
// hex-encoded
func shaFileID(dbPath string) (string, error) {
	fi, err := os.Open(dbPath)
	if err != nil {
		return "", err
	}
	defer fi.Close()

	h := sha256.New()
	if _, err := io.Copy(h, fi); err != nil {
		return "", err
	}

	sum := h.Sum(nil)
	encodedSum := base64.RawURLEncoding.EncodeToString(sum[:8])
	return encodedSum, nil
}
