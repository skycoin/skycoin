package integration_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/kvstorage"
)

func TestGetAllStorageValues(t *testing.T) {
	if !doStable(t) {
		return
	}

	wantVals := map[string]string{
		"key": "val",
	}

	c := newClient()

	err := c.AddStorageValue(kvstorage.TypeNotes, "key", "val")
	require.NoError(t, err)

	vals, err := c.GetAllStorageValues(kvstorage.TypeNotes)
	require.NoError(t, err)
	require.Equal(t, wantVals, vals)
}

func TestGetStorageValue(t *testing.T) {
	if !doStable(t) {
		return
	}

	wantVal := "val"

	c := newClient()

	err := c.AddStorageValue(kvstorage.TypeNotes, "key", "val")
	require.NoError(t, err)

	val, err := c.GetStorageValue(kvstorage.TypeNotes, "key")
	require.NoError(t, err)
	require.Equal(t, wantVal, val)
}

func TestAddStorageValue(t *testing.T) {
	if !doStable(t) {
		return
	}

	c := newClient()

	err := c.AddStorageValue(kvstorage.TypeNotes, "key", "val")
	require.NoError(t, err)
}

func TestRemoveStorageValue(t *testing.T) {
	if !doStable(t) {
		return
	}

	wantVals := map[string]string{
		"key": "val",
	}

	c := newClient()

	err := c.AddStorageValue(kvstorage.TypeNotes, "key", "val")
	require.NoError(t, err)
	err = c.AddStorageValue(kvstorage.TypeNotes, "key2", "val2")
	require.NoError(t, err)

	err = c.RemoveStorageValue(kvstorage.TypeNotes, "key2")
	require.NoError(t, err)

	vals, err := c.GetAllStorageValues(kvstorage.TypeNotes)
	require.NoError(t, err)
	require.Equal(t, wantVals, vals)
}
