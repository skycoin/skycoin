package visor

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSkyencoderDBSafe(t *testing.T) {
	dbFile := "../api/integration/testdata/blockchain-180.db"

	db, err := OpenDB(dbFile, true)
	require.NoError(t, err)

	err = VerifyDBSkyencoderSafe(db, nil)
	require.NoError(t, err)
}
