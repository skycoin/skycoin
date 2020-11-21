package file

import (
	"bytes"
	"crypto/rand"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"encoding/json"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/testutil"
)

func requireFileMode(t *testing.T, filename string, mode os.FileMode) {
	stat, err := os.Stat(filename)
	require.NoError(t, err)
	require.Equal(t, stat.Mode(), mode)
}

func requireFileContentsBinary(t *testing.T, filename string, contents []byte) {
	f, err := os.Open(filename)
	require.NoError(t, err)
	defer f.Close()
	b := make([]byte, len(contents)*16)
	n, err := f.Read(b)
	require.NoError(t, err)

	require.Equal(t, n, len(contents))
	require.True(t, bytes.Equal(b[:n], contents))
}

func requireFileContents(t *testing.T, filename, contents string) { //nolint:unparam
	requireFileContentsBinary(t, filename, []byte(contents))
}

func requireIsRegularFile(t *testing.T, filename string) {
	stat := testutil.RequireFileExists(t, filename)
	require.True(t, stat.Mode().IsRegular())
}

func cleanup(t *testing.T, fn string) {
	paths := []string{fn}
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			panic(err)
		}
		if strings.Contains(info.Name(), fn+".tmp") {
			paths = append(paths, info.Name())
			return nil
		}

		if strings.Contains(info.Name(), fn+".bak") {
			paths = append(paths, info.Name())
			return nil
		}

		return nil
	})
	require.NoError(t, err)

	for _, f := range paths {
		os.Remove(f)
	}
}

func TestCleanup(t *testing.T) {
	fn := "test.bin"
	defer cleanup(t, fn)
	b := make([]byte, 128)
	_, err := rand.Read(b)
	require.NoError(t, err)
	require.NoError(t, ioutil.WriteFile(fn, b, 0600))
	require.NoError(t, ioutil.WriteFile(fn+".tmp.abc", b, 0600))
	require.NoError(t, ioutil.WriteFile(fn+".bak.abc", b, 0600))
}

func TestBuildDataDirDotOk(t *testing.T) {
	dir := "./.test-skycoin/test"
	builtDir, err := buildDataDir(dir)
	require.NoError(t, err)

	cleanDir := filepath.Clean(dir)
	require.True(t, strings.HasSuffix(builtDir, cleanDir))

	gopath := os.Getenv("GOPATH")
	// by default go uses GOPATH=$HOME/go if it is not set
	if gopath == "" {
		home := filepath.Clean(UserHome())
		gopath = filepath.Join(home, "go")
	}

	require.True(t, strings.HasPrefix(builtDir, gopath))
	require.NotEqual(t, builtDir, filepath.Clean(gopath))
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

func TestBuildDataDirDefault(t *testing.T) {
	home := UserHome()
	defaultDir := filepath.Join(home, ".skycoin")
	dir, err := buildDataDir(defaultDir)
	require.NoError(t, err)
	expectedPath := filepath.Join(home, ".skycoin")
	require.Equal(t, dir, expectedPath)
}

func TestBuildDataDirAbsolute(t *testing.T) {
	abspath := "/opt/.skycoin"
	dir, err := buildDataDir(abspath)
	require.NoError(t, err)
	require.Equal(t, abspath, dir)
}

func TestLoadJSON(t *testing.T) {
	obj := struct{ Key string }{}
	fn := "test.json"
	defer cleanup(t, fn)

	// Loading nonexistant file
	testutil.RequireFileNotExists(t, fn)
	err := LoadJSON(fn, &obj)
	require.Error(t, err)
	require.True(t, os.IsNotExist(err))

	f, err := os.Create(fn)
	require.NoError(t, err)
	_, err = f.WriteString("{\"key\":\"value\"}")
	require.NoError(t, err)
	f.Close()

	err = LoadJSON(fn, &obj)
	require.NoError(t, err)
	require.Equal(t, obj.Key, "value")
}

func TestSaveJSON(t *testing.T) {
	fn := "test.json"
	defer cleanup(t, fn)
	obj := struct {
		Key string `json:"key"`
	}{Key: "value"}

	b, err := json.MarshalIndent(obj, "", "    ")
	require.NoError(t, err)

	err = SaveJSON(fn, obj, 0644)
	require.NoError(t, err)

	objHash := cipher.SumSHA256(b).Hex()[:8]

	requireIsRegularFile(t, fn)
	testutil.RequireFileNotExists(t, fn+".bak."+objHash)
	requireFileMode(t, fn, 0644)
	requireFileContents(t, fn, string(b))

	// Saving again should result in a .bak file same as original
	obj.Key = "value2"
	err = SaveJSON(fn, obj, 0644)
	require.NoError(t, err)
	b2, err := json.MarshalIndent(obj, "", "    ")
	require.NoError(t, err)

	requireFileMode(t, fn, 0644)
	requireIsRegularFile(t, fn)
	requireFileContents(t, fn, string(b2))
	objHash = cipher.SumSHA256(b2).Hex()[:8]
	testutil.RequireFileNotExists(t, fn+".tmp."+objHash)
}

func TestSaveJSONSafe(t *testing.T) {
	fn := "test.json"
	defer cleanup(t, fn)
	obj := struct {
		Key string `json:"key"`
	}{Key: "value"}
	err := SaveJSONSafe(fn, obj, 0600)
	require.NoError(t, err)
	b, err := json.MarshalIndent(obj, "", "    ")
	require.NoError(t, err)

	requireIsRegularFile(t, fn)
	requireFileMode(t, fn, 0600)
	requireFileContents(t, fn, string(b))

	// Saving again should result in error, and original file not changed
	obj.Key = "value2"
	err = SaveJSONSafe(fn, obj, 0600)
	require.Error(t, err)

	requireIsRegularFile(t, fn)
	requireFileContents(t, fn, string(b))
	objHash := cipher.SumSHA256(b).Hex()[:8]
	testutil.RequireFileNotExists(t, fn+".bak."+objHash)
	testutil.RequireFileNotExists(t, fn+".tmp."+objHash)
}

func TestSaveBinary(t *testing.T) {
	fn := "test.bin"
	defer cleanup(t, fn)
	b := make([]byte, 128)
	_, err := rand.Read(b)
	require.NoError(t, err)
	err = SaveBinary(fn, b, 0644)
	require.NoError(t, err)
	objHash := cipher.SumSHA256(b).Hex()[:8]
	testutil.RequireFileNotExists(t, fn+".tmp."+objHash)
	testutil.RequireFileNotExists(t, fn+".bak."+objHash)
	requireIsRegularFile(t, fn)
	requireFileContentsBinary(t, fn, b)
	requireFileMode(t, fn, 0644)

	b2 := make([]byte, 128)
	_, err = rand.Read(b2)
	require.NoError(t, err)
	require.False(t, bytes.Equal(b, b2))

	err = SaveBinary(fn, b2, 0644)
	require.NoError(t, err)
	requireIsRegularFile(t, fn)
	objHash = cipher.SumSHA256(b2).Hex()[:8]
	testutil.RequireFileNotExists(t, fn+".tmp."+objHash)
	requireFileContentsBinary(t, fn, b2)
	// requireFileContentsBinary(t, fn+".bak", b)
	requireFileMode(t, fn, 0644)
	// requireFileMode(t, fn+".bak", 0644)
}

func TestIsWritable(t *testing.T) {
	fn := "test.bin"
	defer cleanup(t, fn)
	b := make([]byte, 32)
	_, err := rand.Read(b)
	require.NoError(t, err)
	err = SaveBinary(fn, b, 0400)
	require.NoError(t, err)
	require.False(t, IsWritable(fn))

	err = os.Chmod(fn, 0600)
	require.NoError(t, err)
	require.True(t, IsWritable(fn))
}
