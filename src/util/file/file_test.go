package file

import (
	"bytes"
	"crypto/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"encoding/json"

	"github.com/stretchr/testify/require"
)

func requireFileMode(t *testing.T, filename string, mode os.FileMode) {
	stat, err := os.Stat(filename)
	require.Nil(t, err)
	require.Equal(t, stat.Mode(), mode)
}

func requireFileContentsBinary(t *testing.T, filename string, contents []byte) {
	f, err := os.Open(filename)
	defer f.Close()
	require.Nil(t, err)
	b := make([]byte, len(contents)*16)
	n, err := f.Read(b)
	require.Nil(t, err)

	require.Equal(t, n, len(contents))
	require.True(t, bytes.Equal(b[:n], contents))
}

func requireFileContents(t *testing.T, filename, contents string) {
	requireFileContentsBinary(t, filename, []byte(contents))
}

func requireFileExists(t *testing.T, filename string) {
	stat, err := os.Stat(filename)
	require.Nil(t, err)
	require.True(t, stat.Mode().IsRegular())
}

func requireFileNotExists(t *testing.T, filename string) {
	_, err := os.Stat(filename)
	require.NotNil(t, err)
	require.True(t, os.IsNotExist(err))
}

func cleanup(fn string) {
	os.Remove(fn)
	os.Remove(fn + ".tmp")
	os.Remove(fn + ".bak")
}

func TestBuildDataDir(t *testing.T) {
	dir := "./.test-skycoin/test"
	builtDir, err := buildDataDir(dir)
	require.NoError(t, err)

	cleanDir := filepath.Clean(dir)
	require.True(t, strings.HasSuffix(builtDir, cleanDir))

	home := filepath.Clean(UserHome())
	if home == "" {
		require.Equal(t, cleanDir, builtDir)
	} else {
		require.True(t, strings.HasPrefix(builtDir, home))
		require.NotEqual(t, builtDir, filepath.Clean(home))
	}
}

func TestBuildDataDirEmptyError(t *testing.T) {
	dir, err := buildDataDir("")
	require.Empty(t, dir)
	require.Error(t, err)
	require.Equal(t, ErrEmptyDirectoryName, err)
}

func TestBuildDataDirDotError(t *testing.T) {
	bad := []string{".", "./", "./.", "././", "./../"}
	for _, b := range bad {
		dir, err := buildDataDir(b)
		require.Empty(t, dir)
		require.Error(t, err)
		require.Equal(t, ErrDotDirectoryName, err)
	}
}

func TestUserHome(t *testing.T) {
	home := UserHome()
	require.NotEqual(t, home, "")
}

func TestLoadJSON(t *testing.T) {
	obj := struct{ Key string }{}
	fn := "test.json"
	defer cleanup(fn)

	// Loading nonexistant file
	requireFileNotExists(t, fn)
	err := LoadJSON(fn, &obj)
	require.NotNil(t, err)
	require.True(t, os.IsNotExist(err))

	f, err := os.Create(fn)
	require.Nil(t, err)
	_, err = f.WriteString("{\"key\":\"value\"}")
	require.Nil(t, err)
	f.Close()

	err = LoadJSON(fn, &obj)
	require.Nil(t, err)
	require.Equal(t, obj.Key, "value")
}

func TestSaveJSON(t *testing.T) {
	fn := "test.json"
	defer cleanup(fn)
	obj := struct {
		Key string `json:"key"`
	}{Key: "value"}

	b, err := json.MarshalIndent(obj, "", "    ")
	require.Nil(t, err)

	err = SaveJSON(fn, obj, 0644)
	require.Nil(t, err)

	requireFileExists(t, fn)
	requireFileNotExists(t, fn+".bak")
	requireFileMode(t, fn, 0644)
	requireFileContents(t, fn, string(b))

	// Saving again should result in a .bak file same as original
	obj.Key = "value2"
	err = SaveJSON(fn, obj, 0644)
	require.Nil(t, err)
	b2, err := json.MarshalIndent(obj, "", "    ")
	require.Nil(t, err)

	requireFileMode(t, fn, 0644)
	requireFileExists(t, fn)
	requireFileExists(t, fn+".bak")
	requireFileContents(t, fn, string(b2))
	requireFileContents(t, fn+".bak", string(b))
	requireFileNotExists(t, fn+".tmp")
}

func TestSaveJSONSafe(t *testing.T) {
	fn := "test.json"
	defer cleanup(fn)
	obj := struct {
		Key string `json:"key"`
	}{Key: "value"}
	err := SaveJSONSafe(fn, obj, 0600)
	require.Nil(t, err)
	b, err := json.MarshalIndent(obj, "", "    ")
	require.Nil(t, err)

	requireFileExists(t, fn)
	requireFileMode(t, fn, 0600)
	requireFileContents(t, fn, string(b))

	// Saving again should result in error, and original file not changed
	obj.Key = "value2"
	err = SaveJSONSafe(fn, obj, 0600)
	require.NotNil(t, err)

	requireFileExists(t, fn)
	requireFileContents(t, fn, string(b))
	requireFileNotExists(t, fn+".bak")
	requireFileNotExists(t, fn+".tmp")
}

func TestSaveBinary(t *testing.T) {
	fn := "test.bin"
	defer cleanup(fn)
	b := make([]byte, 128)
	rand.Read(b)
	err := SaveBinary(fn, b, 0644)
	require.Nil(t, err)
	requireFileNotExists(t, fn+".tmp")
	requireFileNotExists(t, fn+".bak")
	requireFileExists(t, fn)
	requireFileContentsBinary(t, fn, b)
	requireFileMode(t, fn, 0644)

	b2 := make([]byte, 128)
	rand.Read(b2)
	require.False(t, bytes.Equal(b, b2))

	err = SaveBinary(fn, b2, 0644)
	requireFileExists(t, fn)
	requireFileExists(t, fn+".bak")
	requireFileNotExists(t, fn+".tmp")
	requireFileContentsBinary(t, fn, b2)
	requireFileContentsBinary(t, fn+".bak", b)
	requireFileMode(t, fn, 0644)
	requireFileMode(t, fn+".bak", 0644)
}
