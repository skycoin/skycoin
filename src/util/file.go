// File and Filesystem related utilities
package util

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/op/go-logging.v1"
)

var (
	DefaultDataDir = ".skycoin/"

	logger = logging.MustGetLogger("skycoin.util")
)

// Disables the logger completely
func DisableLogging() {
	logging.SetBackend(logging.NewLogBackend(ioutil.Discard, "", 0))
}

// If dir is "", uses the default directory of ~/.skycoin.  The path to dir
// is created, and the dir used is returned
func InitDataDir(dir string) string {
	if dir == "" {
		home := UserHome()
		if home == "" {
			logger.Warning("Failed to get home directory")
			dir = filepath.Join("./", DefaultDataDir)
		} else {
			dir = filepath.Join(home, DefaultDataDir)
		}
	}
	if err := os.MkdirAll(dir, os.FileMode(0700)); err != nil {
		logger.Error("Failed to create directory %s: %v", dir, err)
	}
	return dir
}

func UserHome() string {
	// os/user relies on cgo which is disabled when cross compiling
	// use fallbacks for various OSes instead
	// usr, err := user.Current()
	// if err == nil {
	// 	return usr.HomeDir
	// }
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}

	return os.Getenv("HOME")
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
	if err != nil {
		return err
	}
	err = SaveBinary(filename, data, mode)
	return err
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
