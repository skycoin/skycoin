package visor

import (
	"fmt"

	"github.com/blang/semver"

	"github.com/SkycoinProject/skycoin/src/visor/dbutil"
)

var (
	// MetaBkt stores data about the application DB
	MetaBkt = []byte("db_meta")

	versionKey = []byte("version")
)

// GetDBVersion returns the saved DB version
func GetDBVersion(db *dbutil.DB) (*semver.Version, error) {
	var v *semver.Version
	if err := db.View("GetDBVersion", func(tx *dbutil.Tx) error {
		var err error
		v, err = getDBVersion(tx)
		return err
	}); err != nil {
		return nil, err
	}

	return v, nil
}

func getDBVersion(tx *dbutil.Tx) (*semver.Version, error) {
	v, err := dbutil.GetBucketValue(tx, MetaBkt, versionKey)
	if err != nil {
		switch err.(type) {
		case dbutil.ErrBucketNotExist:
			return nil, nil
		default:
			return nil, err
		}
	} else if v == nil {
		return nil, nil
	}

	sv, err := semver.Make(string(v))
	if err != nil {
		return nil, err
	}

	return &sv, nil
}

// SetDBVersion sets the DB version
func SetDBVersion(db *dbutil.DB, version semver.Version) error {
	return db.Update("SetDBVersion", func(tx *dbutil.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(MetaBkt); err != nil {
			return err
		}

		oldVersion, err := getDBVersion(tx)
		if err != nil {
			return err
		}

		if oldVersion != nil && oldVersion.GT(version) {
			return fmt.Errorf("SetDBVersion cannot regress version from %v to %v", oldVersion, version)
		}

		return dbutil.PutBucketValue(tx, MetaBkt, versionKey, []byte(version.String()))
	})
}
