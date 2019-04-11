package integration_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/kvstorage"
)

func TestStableStorageGetAllValues(t *testing.T) {
	if !doStable(t) {
		return
	}

	wantVals := map[string]string{
		"key": "val",
	}

	c := newClient()

	err := c.AddStorageValue(kvstorage.TypeTxIDNotes, "key", "val")
	require.NoError(t, err)

	vals, err := c.GetAllStorageValues(kvstorage.TypeTxIDNotes)
	require.NoError(t, err)
	require.Equal(t, wantVals, vals)
}

func TestStableStorageGetValue(t *testing.T) {
	if !doStable(t) {
		return
	}

	wantVal := "val"

	c := newClient()

	err := c.AddStorageValue(kvstorage.TypeTxIDNotes, "key", "val")
	require.NoError(t, err)

	val, err := c.GetStorageValue(kvstorage.TypeTxIDNotes, "key")
	require.NoError(t, err)
	require.Equal(t, wantVal, val)
}

func TestStableStorageAddValue(t *testing.T) {
	if !doStable(t) {
		return
	}

	c := newClient()

	err := c.AddStorageValue(kvstorage.TypeTxIDNotes, "key", "val")
	require.NoError(t, err)
}

func TestStableStorageRemoveValue(t *testing.T) {
	if !doStable(t) {
		return
	}

	wantVals := map[string]string{
		"key": "val",
	}

	c := newClient()

	err := c.AddStorageValue(kvstorage.TypeTxIDNotes, "key", "val")
	require.NoError(t, err)
	err = c.AddStorageValue(kvstorage.TypeTxIDNotes, "key2", "val2")
	require.NoError(t, err)

	err = c.RemoveStorageValue(kvstorage.TypeTxIDNotes, "key2")
	require.NoError(t, err)

	vals, err := c.GetAllStorageValues(kvstorage.TypeTxIDNotes)
	require.NoError(t, err)
	require.Equal(t, wantVals, vals)
}
