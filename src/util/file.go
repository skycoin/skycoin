// File and Filesystem related utilities
package util

import (
    "encoding/json"
    "errors"
    "fmt"
    "github.com/op/go-logging"
    "io/ioutil"
    "os"
    "os/user"
    "path/filepath"
)

var (
    defaultDataDir = ".skycoin/"
)

// Disable the logger completely
func DisableLogging() {
    logging.SetBackend(logging.NewLogBackend(ioutil.Discard, "", 0))
}

// If dir is "", uses the default directory of ~/.skycoin.  The path to dir
// is created, and the dir used is returned
func InitDataDir(dir string) string {
    if dir == "" {
        home, err := UserHome()
        if err == nil {
            dir = filepath.Join(home, defaultDataDir)
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

func LoadJSON(filename string, thing interface{}) error {
    file, err := ioutil.ReadFile(filename)
    if err != nil {
        return err
    }
    return json.Unmarshal(file, thing)
}

func SaveJSON(filename string, thing interface{}, mode os.FileMode) error {
    data, err := json.MarshalIndent(thing, "", "    ")
    if err == nil {
        return SaveBinary(filename, data, mode)
    } else {
        return err
    }
}

// Saves json to disk, but refuses if file already exists
func SaveJSONSafe(filename string, thing interface{}, mode os.FileMode) error {
    b, err := json.MarshalIndent(thing, "", "    ")
    if err != nil {
        return err
    }
    flags := os.O_WRONLY | os.O_CREATE | os.O_EXCL
    f, err := os.OpenFile(filename, flags, mode)
    if err != nil {
        return err
    }
    defer f.Close()
    n, err := f.Write(b)
    if n != len(b) && err != nil {
        err = errors.New("Failed to save complete file")
    }
    if err != nil {
        os.Remove(filename)
    }
    return err
}

func SaveBinary(filename string, data []byte, mode os.FileMode) error {
    // Write the new file to a temporary
    tmpname := filename + ".tmp"
    if err := ioutil.WriteFile(tmpname, data, mode); err != nil {
        return err
    }
    // Backup the previous file, if there was one
    _, err := os.Stat(filename)
    if !os.IsNotExist(err) {
        if err := os.Rename(filename, filename+".bak"); err != nil {
            return err
        }
    }
    // Move the temporary to the new file
    return os.Rename(tmpname, filename)

}
