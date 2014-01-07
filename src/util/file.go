// File and Filesystem related utilities
package util

import (
    "encoding/gob"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "os"
    "os/user"
    "path"
    "path/filepath"
)

// If dir is "", uses the default directory of ~/.skycoin.  The path to dir
// is created, and the dir used is returned
func InitDataDir(dir string) string {
    if dir == "" {
        dir = ".skycoin/"
        home, err := UserHome()
        if err == nil {
            dir = filepath.Join(home, dir)
        } else {
            fmt.Printf("Warning, failed to get home directory: %v\n", err)
        }
    }
    os.MkdirAll(dir, os.FileMode(0755))
    return dir
}

func UserHome() (string, error) {
    usr, err := user.Current()
    if err != nil {
        return "", err
    }
    return usr.HomeDir, nil
}

//should read by hand in 50 meg chunks
func ReadFile(fp string) []byte {
    b, err := ioutil.ReadFile(fp)
    if err != nil {
        log.Panic(err)
    }
    return b
}

func FileSizeString(n int64) string {
    const KB = 1024.0
    const MB = 1024.0 * 1024.0
    if n == 0 {
        return "0 KB"
    }
    if n < 1024 {
        return "<1 KB"
    }
    if n < 1024*1024 {
        return fmt.Sprintf("%0.0f KB", float64(n)/KB)
    }
    return fmt.Sprintf("%0.0f MB", float64(n)/MB)
}

func FileExists(dpath string, fpath string) bool {
    filename := path.Join(dpath, fpath)
    _, err := os.Stat(filename)
    if os.IsNotExist(err) {
        return false
    }
    if err != nil {
        return false
    }

    return true //can fail for other reason
}

func LoadGob(filename string, thing interface{}) error {
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer file.Close()
    g := gob.NewDecoder(file)
    err = g.Decode(thing)
    if err != nil {
        return err
    }
    return nil
}

func SaveGob(filename string, thing interface{}) error {
    tmpname := filename + ".tmp"
    file, err := os.Create(tmpname)
    if err != nil {
        return err
    }
    defer file.Close()
    g := gob.NewEncoder(file)
    err = g.Encode(thing)
    if err != nil {
        return err
    }
    return os.Rename(tmpname, filename)
}

func LoadJSON(filename string, thing interface{}) error {
    file, err := ioutil.ReadFile(filename)
    if err != nil {
        return err
    }
    return json.Unmarshal(file, thing)
}

func SaveJSON(filename string, thing interface{}) error {
    data, err := json.Marshal(thing)
    if err != nil {
        return err
    }
    tmpname := filename + ".tmp"
    err = ioutil.WriteFile(tmpname, data, os.FileMode(0644))
    if err != nil {
        return err
    }
    return os.Rename(tmpname, filename)
}
