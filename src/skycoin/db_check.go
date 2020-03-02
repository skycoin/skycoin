package skycoin

import (
	"fmt"

	"github.com/blang/semver"

	"github.com/SkycoinProject/skycoin/src/cipher"
	"github.com/SkycoinProject/skycoin/src/util/logging"
	"github.com/SkycoinProject/skycoin/src/visor"
	"github.com/SkycoinProject/skycoin/src/visor/dbutil"
)

// dbCheckParams contains the parameters for verifying db
type dbCheckConfig struct {
	// ForceVerify force verify DB
	ForceVerify bool
	// ResetCorruptDB reset the DB if it is corrupted
	ResetCorruptDB bool
	// AppVersion is the current wallet version
	AppVersion *semver.Version
	// DBVersion is the db version
	DBVersion *semver.Version
	// DBCheckpointVersion is the check point db version
	DBCheckpointVersion *semver.Version
}

type dbAction uint

const (
	doNothing dbAction = iota
	doCheckDB
	doResetCorruptDB
)

func checkAndUpdateDB(c dbCheckConfig, db *dbutil.DB, blockchainPubkey cipher.PubKey, logger *logging.Logger, quit chan struct{}) (*dbutil.DB, error) {
	action, err := checkDB(c)
	if err != nil {
		return nil, err
	}

	switch action {
	case doCheckDB:
		logger.Info("Checking database")
		if err := visor.CheckDatabase(db, blockchainPubkey, quit); err != nil {
			if err != visor.ErrVerifyStopped {
				logger.WithError(err).Error("visor.CheckDatabase failed")
			}
			return nil, err
		}
	case doResetCorruptDB:
		// Check the database integrity and recreate it if necessary
		logger.Info("Checking database and resetting if corrupted")
		newDB, err := visor.ResetCorruptDB(db, blockchainPubkey, quit)
		if err != nil {
			if err != visor.ErrVerifyStopped {
				logger.WithError(err).Error("visor.ResetCorruptDB failed")
			}
			return nil, err
		}
		db = newDB
	}

	if !db.IsReadOnly() {
		if err := visor.SetDBVersion(db, *c.AppVersion); err != nil {
			logger.WithError(err).Error("visor.SetDBVersion failed")
			return nil, err
		}
	}
	return db, nil
}

// checkDB checks whether need to verify or reset the DB version
func checkDB(c dbCheckConfig) (dbAction, error) {
	// If the saved DB version is higher than the app version, abort.
	// Otherwise DB corruption could occur.
	if c.DBVersion != nil && c.DBVersion.GT(*c.AppVersion) {
		return doNothing, fmt.Errorf("Cannot use newer DB version=%v with older software version=%v", c.DBVersion, c.AppVersion)
	}

	// Verify the DB if the version detection says to, or if it was requested on the command line
	if shouldVerifyDB(c.AppVersion, c.DBVersion, c.DBCheckpointVersion) || c.ForceVerify {
		if c.ResetCorruptDB {
			return doResetCorruptDB, nil
		}

		return doCheckDB, nil
	}

	return doNothing, nil
}
