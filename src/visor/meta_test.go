package visor

import (
	"testing"

	"github.com/blang/semver"
	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/testutil"
)

func TestGetSetDBVersion(t *testing.T) {
	db, shutdown := testutil.PrepareDB(t)
	defer shutdown()

	// No version yet
	v, err := GetDBVersion(db)
	require.NoError(t, err)
	require.Nil(t, v)

	// Bucket exists, but still no version
	v, err = GetDBVersion(db)
	require.NoError(t, err)
	require.Nil(t, v)

	// Set the version
	x := semver.MustParse("0.25.0")
	err = SetDBVersion(db, x)
	require.NoError(t, err)

	// Get the version
	v, err = GetDBVersion(db)
	require.NoError(t, err)
	require.NotNil(t, v)
	require.True(t, x.EQ(*v))
	require.Equal(t, "0.25.0", v.String())

	// Set to a new version succeeds
	err = SetDBVersion(db, semver.MustParse("0.26.0"))
	require.NoError(t, err)

	// Set to an older version fails
	err = SetDBVersion(db, x)
	testutil.RequireError(t, err, "SetDBVersion cannot regress version from 0.26.0 to 0.25.0")
}
