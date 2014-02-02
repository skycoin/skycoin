// File and Filesystem related utilities
package util

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os"
    "os/user"
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

func LoadJSON(filename string, thing interface{}) error {
    file, err := ioutil.ReadFile(filename)
    if err != nil {
        return err
    }
    return json.Unmarshal(file, thing)
}

func SaveJSON(filename string, thing interface{}, mode os.FileMode) error {
    data, err := json.Marshal(thing)
    if err != nil {
        return err
    }
    // Write the new file to a temporary
    tmpname := filename + ".tmp"
    err = ioutil.WriteFile(tmpname, data, mode)
    if err != nil {
        return err
    }
    // Backup the previous file
    err := os.Rename(filename, filename+".bak")
    if err != nil {
        return err
    }
    // Move the temporary to the new file
    return os.Rename(tmpname, filename)
}
