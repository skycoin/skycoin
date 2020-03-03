package skycoin

import (
	"errors"
	"testing"

	"github.com/SkycoinProject/skycoin/src/testutil"
	"github.com/SkycoinProject/skycoin/src/visor/dbutil"
	"github.com/blang/semver"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCheckDB(t *testing.T) {
	type expect struct {
		action dbAction
		err    error
	}

	v241, err := semver.New("0.24.1")
	require.NoError(t, err)

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
			params:    dbCheckConfig{AppVersion: v26, DBCheckpointVersion: v25},
			dbVersion: v26,
			exp:       expect{action: doNothing},
		},
		{
			name:   "db version nil, check db",
			params: dbCheckConfig{AppVersion: v26},
			exp:    expect{action: doCheckDB},
		},
		{
			name:   "db version nil, reset corrupt db",
			params: dbCheckConfig{AppVersion: v26, ResetCorruptDB: true, DBCheckpointVersion: v25},
			exp:    expect{action: doResetCorruptDB},
		},
		{
			name:      "db version > app version get err",
			params:    dbCheckConfig{AppVersion: v241, DBCheckpointVersion: v25},
			dbVersion: v25,
			exp:       expect{action: doNothing, err: errors.New("Cannot use newer DB version=0.25.0 with older software version=0.24.1")},
		},
		{
			name:      "db version < check point < app version, check DB",
			params:    dbCheckConfig{AppVersion: v26, DBCheckpointVersion: v25},
			dbVersion: v241,
			exp:       expect{action: doCheckDB},
		},
		{
			name:      "db version < check point < app version, reset corrupt db",
			params:    dbCheckConfig{AppVersion: v26, ResetCorruptDB: true, DBCheckpointVersion: v25},
			dbVersion: v241,
			exp:       expect{action: doResetCorruptDB},
		},
		{
			name:      "db version < check point == app version, check DB",
			params:    dbCheckConfig{AppVersion: v25, DBCheckpointVersion: v25},
			dbVersion: v241,
			exp:       expect{action: doCheckDB},
		},
		{
			name:      "db version < check point == app version, reset corrupt DB",
			params:    dbCheckConfig{AppVersion: v25, ResetCorruptDB: true, DBCheckpointVersion: v25},
			dbVersion: v241,
			exp:       expect{action: doResetCorruptDB},
		},
		{
			name:      "db version == app version, force verify, check DB",
			params:    dbCheckConfig{AppVersion: v26, ForceVerify: true, DBCheckpointVersion: v25},
			dbVersion: v26,
			exp:       expect{action: doCheckDB},
		},
		{
			name:      "db version == app version, force verify, reset corrupt DB",
			params:    dbCheckConfig{AppVersion: v26, ForceVerify: true, ResetCorruptDB: true, DBCheckpointVersion: v25},
			dbVersion: v26,
			exp:       expect{action: doResetCorruptDB},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			action, err := checkDB(tc.params, tc.dbVersion)
			require.Equal(t, tc.exp.action, action)
			require.Equal(t, tc.exp.err, err)
		})
	}
}

func TestCheckAndUpdateDB(t *testing.T) {
	v241, err := semver.New("0.24.1")
	require.NoError(t, err)

	v25, err := semver.New("0.25.0")
	require.NoError(t, err)

	v26, err := semver.New("0.26.0")
	require.NoError(t, err)

	matchFunc := mock.MatchedBy(func(db *dbutil.DB) bool {
		return true
	})

	resetedDB, closeDB := testutil.PrepareDB(t)
	defer closeDB()

	tt := []struct {
		name            string
		config          dbCheckConfig
		dbVersion       *semver.Version
		dbVersionErr    error
		checkDBErr      error
		resetedDB       *dbutil.DB
		resetDBErr      error
		setDBVersion    *semver.Version
		setDBVersionErr error
		retDB           *dbutil.DB
		retErr          error
		assertCalled    func(t *testing.T, db *dbutil.DB, m *mockDbCheckCorruptResetter)
	}{
		{
			name:         "do nothing",
			config:       dbCheckConfig{AppVersion: v26, DBCheckpointVersion: v25},
			dbVersion:    v26,
			setDBVersion: v26,
			assertCalled: func(t *testing.T, db *dbutil.DB, m *mockDbCheckCorruptResetter) {
				require.True(t, m.AssertCalled(t, "GetDBVersion", matchFunc))
				require.True(t, m.AssertCalled(t, "SetDBVersion", matchFunc, v26))
				require.True(t, m.AssertNotCalled(t, "CheckDatabase", matchFunc))
				require.True(t, m.AssertNotCalled(t, "ResetCorruptDB", matchFunc))
			},
		},
		{
			name:         "db version nil - check db",
			config:       dbCheckConfig{AppVersion: v26, DBCheckpointVersion: v25},
			dbVersion:    nil,
			setDBVersion: v26,
			assertCalled: func(t *testing.T, db *dbutil.DB, m *mockDbCheckCorruptResetter) {
				require.True(t, m.AssertCalled(t, "GetDBVersion", matchFunc))
				require.True(t, m.AssertCalled(t, "CheckDatabase", matchFunc))
				require.True(t, m.AssertCalled(t, "SetDBVersion", matchFunc, v26))
				require.True(t, m.AssertNotCalled(t, "ResetCorruptDB", matchFunc))
			},
		},
		{
			name:         "db version nil - reset corrupt db",
			config:       dbCheckConfig{AppVersion: v26, ResetCorruptDB: true, DBCheckpointVersion: v25},
			dbVersion:    nil,
			setDBVersion: v26,
			resetedDB:    resetedDB,
			assertCalled: func(t *testing.T, db *dbutil.DB, m *mockDbCheckCorruptResetter) {
				require.True(t, m.AssertCalled(t, "GetDBVersion", matchFunc))
				require.True(t, m.AssertNotCalled(t, "CheckDatabase", matchFunc))
				require.True(t, m.AssertCalled(t, "SetDBVersion", matchFunc, v26))
				require.True(t, m.AssertCalled(t, "ResetCorruptDB", matchFunc))
			},
		},
		{
			name: "db version > app version get err",
			config: dbCheckConfig{
				AppVersion:          v241,
				DBCheckpointVersion: v25,
			},
			dbVersion: v25,
			retErr:    errors.New("Cannot use newer DB version=0.25.0 with older software version=0.24.1"),
			assertCalled: func(t *testing.T, db *dbutil.DB, m *mockDbCheckCorruptResetter) {
				require.True(t, m.AssertCalled(t, "GetDBVersion", matchFunc))
				require.True(t, m.AssertNotCalled(t, "CheckDatabase", matchFunc))
				require.True(t, m.AssertNotCalled(t, "SetDBVersion", matchFunc, v26))
				require.True(t, m.AssertNotCalled(t, "ResetCorruptDB", matchFunc))
			},
		},
		{
			name: "db version < check point < app version - check db",
			config: dbCheckConfig{
				AppVersion:          v26,
				DBCheckpointVersion: v25,
			},
			dbVersion:    v241,
			checkDBErr:   nil,
			setDBVersion: v26,
			assertCalled: func(t *testing.T, db *dbutil.DB, m *mockDbCheckCorruptResetter) {
				require.True(t, m.AssertCalled(t, "GetDBVersion", matchFunc))
				require.True(t, m.AssertCalled(t, "CheckDatabase", matchFunc))
				require.True(t, m.AssertCalled(t, "SetDBVersion", matchFunc, v26))
				require.True(t, m.AssertNotCalled(t, "ResetCorruptDB", matchFunc))
			},
		},
		{
			name: "db version < check point < app version - check db - got err",
			config: dbCheckConfig{
				AppVersion:          v26,
				DBCheckpointVersion: v25,
			},
			dbVersion:  v241,
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
			name: "db version < check point < app version - reset corrupt db",
			config: dbCheckConfig{
				AppVersion:          v26,
				DBCheckpointVersion: v25,
				ResetCorruptDB:      true,
			},
			dbVersion:    v241,
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
			name: "db version < check point < app version - reset corrupt db - got err",
			config: dbCheckConfig{
				AppVersion:          v26,
				DBCheckpointVersion: v25,
				ResetCorruptDB:      true,
			},
			dbVersion:  v241,
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
			name: "db version < check point == app version - check DB",
			config: dbCheckConfig{
				AppVersion:          v25,
				DBCheckpointVersion: v25,
			},
			dbVersion:    v241,
			setDBVersion: v25,
			assertCalled: func(t *testing.T, db *dbutil.DB, m *mockDbCheckCorruptResetter) {
				require.True(t, m.AssertCalled(t, "GetDBVersion", matchFunc))
				require.True(t, m.AssertCalled(t, "CheckDatabase", matchFunc))
				require.True(t, m.AssertNotCalled(t, "ResetCorruptDB", matchFunc))
				require.True(t, m.AssertCalled(t, "SetDBVersion", matchFunc, v25))
			},
		},
		{
			name: "db version < check point == app version - reset DB",
			config: dbCheckConfig{
				AppVersion:          v25,
				DBCheckpointVersion: v25,
				ResetCorruptDB:      true,
			},
			dbVersion:    v241,
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
			name: "db version == app version - force verify - check DB",
			config: dbCheckConfig{
				AppVersion:          v26,
				DBCheckpointVersion: v25,
				ForceVerify:         true,
			},
			dbVersion:    v26,
			setDBVersion: v26,
			assertCalled: func(t *testing.T, db *dbutil.DB, m *mockDbCheckCorruptResetter) {
				require.True(t, m.AssertCalled(t, "GetDBVersion", matchFunc))
				require.True(t, m.AssertCalled(t, "CheckDatabase", matchFunc))
				require.True(t, m.AssertNotCalled(t, "ResetCorruptDB", matchFunc))
				require.True(t, m.AssertCalled(t, "SetDBVersion", matchFunc, v26))
			},
		},
		{
			name: "db version == app version - force verify - reset corrupt DB",
			config: dbCheckConfig{
				AppVersion:          v26,
				DBCheckpointVersion: v25,
				ForceVerify:         true,
				ResetCorruptDB:      true,
			},
			dbVersion:    v26,
			setDBVersion: v26,
			resetedDB:    resetedDB,
			assertCalled: func(t *testing.T, db *dbutil.DB, m *mockDbCheckCorruptResetter) {
				require.True(t, m.AssertCalled(t, "GetDBVersion", matchFunc))
				require.True(t, m.AssertNotCalled(t, "CheckDatabase", matchFunc))
				require.True(t, m.AssertCalled(t, "ResetCorruptDB", matchFunc))
				require.True(t, m.AssertCalled(t, "SetDBVersion", matchFunc, v26))
			},
		},
	}

	db, closeDBV2 := testutil.PrepareDB(t)
	defer closeDBV2()

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			m := &mockDbCheckCorruptResetter{}
			m.On("GetDBVersion", matchFunc).Return(tc.dbVersion, tc.dbVersionErr)
			m.On("SetDBVersion", matchFunc, tc.setDBVersion).Return(tc.setDBVersionErr)
			m.On("CheckDatabase", matchFunc).Return(tc.checkDBErr)
			m.On("ResetCorruptDB", matchFunc).Return(tc.resetedDB, tc.resetDBErr)

			dbAfter, err := checkAndUpdateDB(db, tc.config, m)
			require.Equal(t, tc.retErr, err)
			tc.assertCalled(t, db, m)
			if err != nil {
				return
			}

			if tc.resetedDB != nil {
				require.Equal(t, tc.resetedDB, dbAfter)
			} else {
				require.Equal(t, db, dbAfter)
			}
		})
	}
}
