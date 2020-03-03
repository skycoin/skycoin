package skycoin

import (
	"errors"
	"testing"

	"github.com/blang/semver"
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
			params:    dbCheckConfig{AppVersion: &dbVerifyCheckpointVersionParsed, DBCheckpointVersion: v25},
			dbVersion: v241,
			exp:       expect{action: doCheckDB},
		},
		{
			name:      "db version < check point == app version, reset corrupt DB",
			params:    dbCheckConfig{AppVersion: &dbVerifyCheckpointVersionParsed, ResetCorruptDB: true, DBCheckpointVersion: v25},
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

}
