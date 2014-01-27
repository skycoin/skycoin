package util

import (
    "github.com/stretchr/testify/assert"
    "os"
    "testing"
)

func TestInitDataDir(t *testing.T) {
    defer os.RemoveAll("./.test")
    d := "./.test/test"
    dir := InitDataDir(d)
    assert.Equal(t, dir, d)
    _, err := os.Stat(d)
    assert.Nil(t, err)
}

func TestUserHome(t *testing.T) {
    home, err := UserHome()
    assert.NotEqual(t, home, "")
    assert.Nil(t, err)
}

func TestLoadJSON(t *testing.T) {
    fn := "test.json"
    defer os.Remove(fn)
    f, err := os.Create(fn)
    assert.Nil(t, err)
    _, err = f.WriteString("{\"key\":\"value\"}")
    assert.Nil(t, err)
    f.Close()

    obj := struct{ Key string }{}
    err = LoadJSON(fn, &obj)
    assert.Nil(t, err)
    assert.Equal(t, obj.Key, "value")
}

func TestSaveJSON(t *testing.T) {
    fn := "test.json"
    defer os.Remove(fn)
    obj := struct {
        Key string `json:"key"`
    }{Key: "value"}
    err := SaveJSON(fn, obj)
    assert.Nil(t, err)

    f, err := os.Open(fn)
    assert.Nil(t, err)
    b := make([]byte, 128)
    n, err := f.Read(b)
    assert.Nil(t, err)
    assert.Equal(t, string(b[:n]), "{\"key\":\"value\"}")
}
