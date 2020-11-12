package skycoin

import (
	"errors"
	"testing"

	"github.com/blang/semver"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/testutil"
	"github.com/skycoin/skycoin/src/visor/dbutil"
)

func TestCheckDB(t *testing.T) {
	type expect struct {
		action dbAction
		err    error
	}

	v25, err := semver.New("0.25.0")
	require.NoError(t, err)

	v26, err := semver.New("0.26.0")
	require.NoError(t, err)

	tt := []struct {
		name      string
		params    dbCheckConfig
		dbVersion *semver.Version
		exp       expect
	}{
		{
			name:      "do nothing",
			params:    dbCheckConfig{DBCheckpointVersion: v25},
			dbVersion: v25,
			exp:       expect{action: doNothing},
		},
		{
			name:   "db version nil, check db",
			params: dbCheckConfig{DBCheckpointVersion: v25},
			exp:    expect{action: doCheckDB},
		},
		{
			name:   "db version nil, reset corrupt db",
			params: dbCheckConfig{ResetCorruptDB: true, DBCheckpointVersion: v25},
			exp:    expect{action: doResetCorruptDB},
		},
		{
			name:      "db version > check point, get err",
			params:    dbCheckConfig{DBCheckpointVersion: v25},
			dbVersion: v26,
			exp:       expect{action: doNothing, err: errors.New("Cannot use newer DB version=0.26.0 with older check point version=0.25.0")},
		},
		{
			name:      "db version < check point, check DB",
			params:    dbCheckConfig{DBCheckpointVersion: v26},
			dbVersion: v25,
			exp:       expect{action: doCheckDB},
		},
		{
			name:      "db version < check point, reset corrupt db",
			params:    dbCheckConfig{ResetCorruptDB: true, DBCheckpointVersion: v26},
			dbVersion: v25,
			exp:       expect{action: doResetCorruptDB},
		},
		{
			name:      "db version == check point, force verify, check DB",
			params:    dbCheckConfig{ForceVerify: true, DBCheckpointVersion: v25},
			dbVersion: v25,
			exp:       expect{action: doCheckDB},
		},
		{
			name:      "db version == app version, force verify, reset corrupt DB",
			params:    dbCheckConfig{ForceVerify: true, ResetCorruptDB: true, DBCheckpointVersion: v25},
			dbVersion: v25,
			exp:       expect{action: doResetCorruptDB},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			action, err := checkDBVersion(tc.params, tc.dbVersion)
			require.Equal(t, tc.exp.action, action)
			require.Equal(t, tc.exp.err, err)
		})
	}
}

func TestCheckAndUpdateDB(t *testing.T) {
	v25, err := semver.New("0.25.0")
	require.NoError(t, err)

	v26, err := semver.New("0.26.0")
	require.NoError(t, err)

	matchFunc := mock.MatchedBy(func(db *dbutil.DB) bool {
		return true
	})

	resetedDB, closeRSDB := testutil.PrepareDB(t)
	defer closeRSDB()

	readOnlyDB, closeRODB := testutil.PrepareDBReadOnly(t)
	defer closeRODB()

	db, closeDB := testutil.PrepareDB(t)
	defer closeDB()

	tt := []struct {
		name            string
		config          dbCheckConfig
		db              *dbutil.DB
		dbVersion       *semver.Version
		dbVersionErr    error
		checkDBErr      error
		resetedDB       *dbutil.DB
		resetDBErr      error
		setDBVersion    *semver.Version
		setDBVersionErr error
		retErr          error
		assertCalled    func(t *testing.T, db *dbutil.DB, m *mockDbCheckCorruptResetter)
	}{
		{
			name:         "do nothing",
			config:       dbCheckConfig{DBCheckpointVersion: v26},
			db:           db,
			dbVersion:    v26,
			setDBVersion: v26,
			assertCalled: func(t *testing.T, db *dbutil.DB, m *mockDbCheckCorruptResetter) {
				require.True(t, m.AssertCalled(t, "GetDBVersion", matchFunc))
				require.True(t, m.AssertNotCalled(t, "CheckDatabase", matchFunc))
				require.True(t, m.AssertNotCalled(t, "ResetCorruptDB", matchFunc))
				require.True(t, m.AssertNotCalled(t, "SetDBVersion", matchFunc, v26))
			},
		},
		{
			name:         "db version nil - check db",
			config:       dbCheckConfig{DBCheckpointVersion: v25},
			db:           db,
			dbVersion:    nil,
			setDBVersion: v25,
			assertCalled: func(t *testing.T, db *dbutil.DB, m *mockDbCheckCorruptResetter) {
				require.True(t, m.AssertCalled(t, "GetDBVersion", matchFunc))
				require.True(t, m.AssertCalled(t, "CheckDatabase", matchFunc))
				require.True(t, m.AssertNotCalled(t, "ResetCorruptDB", matchFunc))
				require.True(t, m.AssertCalled(t, "SetDBVersion", matchFunc, v25))
			},
		},
		{
			name:         "db version nil - reset corrupt db",
			config:       dbCheckConfig{ResetCorruptDB: true, DBCheckpointVersion: v25},
			db:           db,
			dbVersion:    nil,
			setDBVersion: v25,
			resetedDB:    resetedDB,
			assertCalled: func(t *testing.T, db *dbutil.DB, m *mockDbCheckCorruptResetter) {
				require.True(t, m.AssertCalled(t, "GetDBVersion", matchFunc))
				require.True(t, m.AssertNotCalled(t, "CheckDatabase", matchFunc))
				require.True(t, m.AssertCalled(t, "ResetCorruptDB", matchFunc))
				require.True(t, m.AssertCalled(t, "SetDBVersion", matchFunc, v25))
			},
		},
		{
			name:      "db version > check point get err",
			config:    dbCheckConfig{DBCheckpointVersion: v25},
			db:        db,
			dbVersion: v26,
			retErr:    errors.New("Cannot use newer DB version=0.26.0 with older check point version=0.25.0"),
			assertCalled: func(t *testing.T, db *dbutil.DB, m *mockDbCheckCorruptResetter) {
				require.True(t, m.AssertCalled(t, "GetDBVersion", matchFunc))
				require.True(t, m.AssertNotCalled(t, "CheckDatabase", matchFunc))
				require.True(t, m.AssertNotCalled(t, "SetDBVersion", matchFunc, v25))
				require.True(t, m.AssertNotCalled(t, "ResetCorruptDB", matchFunc))
			},
		},
		{
			name:         "db version < check point - check db",
			config:       dbCheckConfig{DBCheckpointVersion: v26},
			db:           db,
			dbVersion:    v25,
			setDBVersion: v26,
			assertCalled: func(t *testing.T, db *dbutil.DB, m *mockDbCheckCorruptResetter) {
				require.True(t, m.AssertCalled(t, "GetDBVersion", matchFunc))
				require.True(t, m.AssertCalled(t, "CheckDatabase", matchFunc))
				require.True(t, m.AssertCalled(t, "SetDBVersion", matchFunc, v26))
				require.True(t, m.AssertNotCalled(t, "ResetCorruptDB", matchFunc))
			},
		},
		{
			name:      "db version < check point - check db - read only db",
			config:    dbCheckConfig{DBCheckpointVersion: v26},
			db:        readOnlyDB,
			dbVersion: v25,
			assertCalled: func(t *testing.T, db *dbutil.DB, m *mockDbCheckCorruptResetter) {
				require.True(t, m.AssertCalled(t, "GetDBVersion", matchFunc))
				require.True(t, m.AssertCalled(t, "CheckDatabase", matchFunc))
				require.True(t, m.AssertNotCalled(t, "SetDBVersion", matchFunc, v26))
				require.True(t, m.AssertNotCalled(t, "ResetCorruptDB", matchFunc))
			},
		},
		{
			name:       "db version < check point - check db - got err",
			config:     dbCheckConfig{DBCheckpointVersion: v26},
			db:         db,
			dbVersion:  v25,
			checkDBErr: errors.New("check db error"),
			retErr:     errors.New("check db error"),
			assertCalled: func(t *testing.T, db *dbutil.DB, m *mockDbCheckCorruptResetter) {
				require.True(t, m.AssertCalled(t, "GetDBVersion", matchFunc))
				require.True(t, m.AssertCalled(t, "CheckDatabase", matchFunc))
				require.True(t, m.AssertNotCalled(t, "SetDBVersion", matchFunc, v26))
				require.True(t, m.AssertNotCalled(t, "ResetCorruptDB", matchFunc))
			},
		},
		{
			name:         "db version < check point - reset corrupt db",
			config:       dbCheckConfig{DBCheckpointVersion: v26, ResetCorruptDB: true},
			db:           db,
			dbVersion:    v25,
			resetedDB:    resetedDB,
			setDBVersion: v26,
			assertCalled: func(t *testing.T, db *dbutil.DB, m *mockDbCheckCorruptResetter) {
				require.True(t, m.AssertCalled(t, "GetDBVersion", matchFunc))
				require.True(t, m.AssertNotCalled(t, "CheckDatabase", matchFunc))
				require.True(t, m.AssertCalled(t, "ResetCorruptDB", matchFunc))
				require.True(t, m.AssertCalled(t, "SetDBVersion", matchFunc, v26))
			},
		},
		{
			name:       "db version < check point - reset corrupt db - got err",
			config:     dbCheckConfig{DBCheckpointVersion: v26, ResetCorruptDB: true},
			db:         db,
			dbVersion:  v25,
			resetDBErr: errors.New("reset corrupt db failed"),
			retErr:     errors.New("reset corrupt db failed"),
			assertCalled: func(t *testing.T, db *dbutil.DB, m *mockDbCheckCorruptResetter) {
				require.True(t, m.AssertCalled(t, "GetDBVersion", matchFunc))
				require.True(t, m.AssertNotCalled(t, "CheckDatabase", matchFunc))
				require.True(t, m.AssertCalled(t, "ResetCorruptDB", matchFunc))
				require.True(t, m.AssertNotCalled(t, "SetDBVersion", matchFunc, v26))
			},
		},
		{
			name: "db version == check point - force verify - check DB",
			config: dbCheckConfig{
				DBCheckpointVersion: v25,
				ForceVerify:         true,
			},
			db:           db,
			dbVersion:    v25,
			setDBVersion: v25,
			assertCalled: func(t *testing.T, db *dbutil.DB, m *mockDbCheckCorruptResetter) {
				require.True(t, m.AssertCalled(t, "GetDBVersion", matchFunc))
				require.True(t, m.AssertCalled(t, "CheckDatabase", matchFunc))
				require.True(t, m.AssertNotCalled(t, "ResetCorruptDB", matchFunc))
				require.True(t, m.AssertCalled(t, "SetDBVersion", matchFunc, v25))
			},
		},
		{
			name: "db version == app version - force verify - reset corrupt DB",
			config: dbCheckConfig{
				DBCheckpointVersion: v25,
				ForceVerify:         true,
				ResetCorruptDB:      true,
			},
			db:           db,
			dbVersion:    v25,
			setDBVersion: v25,
			resetedDB:    resetedDB,
			assertCalled: func(t *testing.T, db *dbutil.DB, m *mockDbCheckCorruptResetter) {
				require.True(t, m.AssertCalled(t, "GetDBVersion", matchFunc))
				require.True(t, m.AssertNotCalled(t, "CheckDatabase", matchFunc))
				require.True(t, m.AssertCalled(t, "ResetCorruptDB", matchFunc))
				require.True(t, m.AssertCalled(t, "SetDBVersion", matchFunc, v25))
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			m := &mockDbCheckCorruptResetter{}
			m.On("GetDBVersion", matchFunc).Return(tc.dbVersion, tc.dbVersionErr)
			m.On("SetDBVersion", matchFunc, tc.setDBVersion).Return(tc.setDBVersionErr)
			m.On("CheckDatabase", matchFunc).Return(tc.checkDBErr)
			m.On("ResetCorruptDB", matchFunc).Return(tc.resetedDB, tc.resetDBErr)

			dbAfter, err := checkAndUpdateDB(tc.db, tc.config, m)
			require.Equal(t, tc.retErr, err)
			tc.assertCalled(t, tc.db, m)
			if err != nil {
				return
			}

			if tc.resetedDB != nil {
				require.Equal(t, tc.resetedDB, dbAfter)
			} else {
				require.Equal(t, tc.db, dbAfter)
			}
		})
	}
}
