package util

import (
	"bytes"
	"crypto/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertFileMode(t *testing.T, filename string, mode os.FileMode) {
	stat, err := os.Stat(filename)
	assert.Nil(t, err)
	assert.Equal(t, stat.Mode(), mode)
}

func assertFileContentsBinary(t *testing.T, filename string, contents []byte) {
	f, err := os.Open(filename)
	defer f.Close()
	assert.Nil(t, err)
	b := make([]byte, len(contents)*16)
	n, err := f.Read(b)
	assert.Nil(t, err)
	assert.Equal(t, n, len(contents))
	assert.True(t, bytes.Equal(b[:n], contents))
}

func assertFileContents(t *testing.T, filename, contents string) {
	assertFileContentsBinary(t, filename, []byte(contents))
}

func assertFileExists(t *testing.T, filename string) {
	stat, err := os.Stat(filename)
	assert.Nil(t, err)
	assert.True(t, stat.Mode().IsRegular())
}

func assertFileNotExists(t *testing.T, filename string) {
	_, err := os.Stat(filename)
	assert.NotNil(t, err)
	assert.True(t, os.IsNotExist(err))
}

func assertDirExists(t *testing.T, dirname string) {
	stat, err := os.Stat(dirname)
	assert.Nil(t, err)
	assert.True(t, stat.IsDir())
}

func assertDirNotExists(t *testing.T, dirname string) {
	_, err := os.Stat(dirname)
	assert.NotNil(t, err)
	assert.True(t, os.IsNotExist(err))
}

func cleanup(fn string) {
	os.Remove(fn)
	os.Remove(fn + ".tmp")
	os.Remove(fn + ".bak")
}

func TestInitDataDir(t *testing.T) {
	defer os.RemoveAll("./.test")
	d := "./.test/test"
	assertDirNotExists(t, d)
	dir := InitDataDir(d)
	assertDirExists(t, dir)
	_, err := os.Stat(dir)
	assert.Nil(t, err)
	os.RemoveAll(dir)
}

func TestInitDataDirDefault(t *testing.T) {
	defaultDataDir := ".skycointestXCAWDAWD232232"
	home := UserHome()
	assertDirNotExists(t, filepath.Join(home, defaultDataDir))
	dir := InitDataDir(defaultDataDir)
	assert.NotEqual(t, dir, "")
	assert.True(t, strings.HasSuffix(dir, defaultDataDir))
	assertDirExists(t, dir)
	os.RemoveAll(dir)

}

func TestUserHome(t *testing.T) {
	home := UserHome()
	assert.NotEqual(t, home, "")
}

func TestLoadJSON(t *testing.T) {
	obj := struct{ Key string }{}
	fn := "test.json"
	defer cleanup(fn)

	// Loading nonexistant file
	assertFileNotExists(t, fn)
	err := LoadJSON(fn, &obj)
	assert.NotNil(t, err)
	assert.True(t, os.IsNotExist(err))

	f, err := os.Create(fn)
	assert.Nil(t, err)
	_, err = f.WriteString("{\"key\":\"value\"}")
	assert.Nil(t, err)
	f.Close()

	err = LoadJSON(fn, &obj)
	assert.Nil(t, err)
	assert.Equal(t, obj.Key, "value")
}

func TestSaveJSON(t *testing.T) {
	fn := "test.json"
	defer cleanup(fn)
	obj := struct {
		Key string `json:"key"`
	}{Key: "value"}
	err := SaveJSON(fn, obj, 0644)
	assert.Nil(t, err)
	assertFileExists(t, fn)
	assertFileNotExists(t, fn+".bak")
	assertFileMode(t, fn, 0644)
	assertFileContents(t, fn, "{\"key\":\"value\"}")

	// Saving again should result in a .bak file same as original
	obj.Key = "value2"
	err = SaveJSON(fn, obj, 0644)
	assert.Nil(t, err)
	assertFileMode(t, fn, 0644)
	assertFileExists(t, fn)
	assertFileExists(t, fn+".bak")
	assertFileContents(t, fn, "{\"key\":\"value2\"}")
	assertFileContents(t, fn+".bak", "{\"key\":\"value\"}")
	assertFileNotExists(t, fn+".tmp")
}

func TestSaveJSONSafe(t *testing.T) {
	fn := "test.json"
	defer cleanup(fn)
	obj := struct {
		Key string `json:"key"`
	}{Key: "value"}
	err := SaveJSONSafe(fn, obj, 0600)
	assert.Nil(t, err)
	assertFileExists(t, fn)
	assertFileMode(t, fn, 0600)
	assertFileContents(t, fn, "{\"key\":\"value\"}")

	// Saving again should result in error, and original file not changed
	obj.Key = "value2"
	err = SaveJSONSafe(fn, obj, 0600)
	assert.NotNil(t, err)
	assertFileExists(t, fn)
	assertFileContents(t, fn, "{\"key\":\"value\"}")
	assertFileNotExists(t, fn+".bak")
	assertFileNotExists(t, fn+".tmp")
}

func TestSaveBinary(t *testing.T) {
	fn := "test.bin"
	defer cleanup(fn)
	b := make([]byte, 128)
	rand.Read(b)
	err := SaveBinary(fn, b, 0644)
	assert.Nil(t, err)
	assertFileNotExists(t, fn+".tmp")
	assertFileNotExists(t, fn+".bak")
	assertFileExists(t, fn)
	assertFileContentsBinary(t, fn, b)
	assertFileMode(t, fn, 0644)

	b2 := make([]byte, 128)
	rand.Read(b2)
	assert.False(t, bytes.Equal(b, b2))

	err = SaveBinary(fn, b2, 0644)
	assertFileExists(t, fn)
	assertFileExists(t, fn+".bak")
	assertFileNotExists(t, fn+".tmp")
	assertFileContentsBinary(t, fn, b2)
	assertFileContentsBinary(t, fn+".bak", b)
	assertFileMode(t, fn, 0644)
	assertFileMode(t, fn+".bak", 0644)
}
