// Package file Filesystem related utilities
package file

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/skycoin/skycoin/src/util/logging"
)

var (
	// ErrEmptyDirectoryName is returned by constructing the full path
	// of data directory if the passed argument is empty
	ErrEmptyDirectoryName = errors.New("data directory must not be empty")
	// ErrDotDirectoryName is returned by constructing the full path of
	// data directory if the passed argument is "."
	ErrDotDirectoryName = errors.New("data directory must not be equivalent to .")

	logger = logging.MustGetLogger("file")
)

// InitDataDir Joins dir with the user's $HOME directory.
// If $HOME cannot be determined, uses the current working directory.
// dir must not be the empty string
func InitDataDir(dir string) (string, error) {
	dir, err := buildDataDir(dir)
	if err != nil {
		return "", err
	}

	// check if dir already exist
	st, err := os.Stat(dir)
	if !os.IsNotExist(err) {
		if !st.IsDir() {
			return "", fmt.Errorf("%s is not a directory", dir)
		}
		// dir already exist
		return dir, nil
	}

	if err := os.MkdirAll(dir, os.FileMode(0700)); err != nil {
		logger.Error("Failed to create directory %s: %v", dir, err)
		return "", err
	}

	logger.Info("Created data directory %s", dir)
	return dir, nil
}

// Construct the full data directory by adding to $HOME or ./
func buildDataDir(dir string) (string, error) {
	if dir == "" {
		logger.Error("data directory is empty")
		return "", ErrEmptyDirectoryName
	}

	home := UserHome()
	if home == "" {
		logger.Warning("Failed to get home directory, using ./")
		home = "./"
	} else {
		home = filepath.Clean(home)
	}

	fullDir := filepath.Join(home, dir)
	fullDir = filepath.Clean(fullDir)

	// The joined directory must not be equal to $HOME or a parent path of $HOME
	if strings.HasPrefix(home, fullDir) {
		logger.Error("join(%[1]s, %[2]s) == %[1]s", home, dir)
		return "", ErrDotDirectoryName
	}

	return fullDir, nil
}

// UserHome returns the current user home path
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

// LoadJSON load json file
func LoadJSON(filename string, thing interface{}) error {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(file, thing)
}

// SaveJSON write value into json file
func SaveJSON(filename string, thing interface{}, mode os.FileMode) error {
	data, err := json.MarshalIndent(thing, "", "    ")
	if err != nil {
		return err
	}
	err = SaveBinary(filename, data, mode)
	return err
}

// SaveJSONSafe saves json to disk, but refuses if file already exists
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

// SaveBinary persists data into given file in binary
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

//TODO: require file named after application and then hashcode, in static directory

// ResolveResourceDirectory searches locations for a research directory and returns absolute path
func ResolveResourceDirectory(path string) string {
	workDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logger.Panic(err)
	}

	_, rtFilename, _, _ := runtime.Caller(1)
	rtDirectory := filepath.Dir(rtFilename)

	pathAbs, err := filepath.Abs(path)
	if err != nil {
		logger.Panic(err)
	}
	fmt.Println("abs path:", pathAbs)

	fmt.Printf("runtime.Caller= %s \n", rtFilename)
	//fmt.Printf("Filepath Raw= %s \n")
	fmt.Printf("Filepath Directory= %s \n", filepath.Dir(path))
	fmt.Printf("Filepath Absolute Directory= %s \n", pathAbs)

	fmt.Printf("Working Directory= %s \n", workDir)
	fmt.Printf("Runtime Filename= %s \n", rtFilename)
	fmt.Printf("Runtime Directory= %s \n", rtDirectory)

	//dir1 := filepath.Join(workDir, filepath.Dir(path))
	//fmt.Printf("Dir1= %s \n", dir1)

	dirs := []string{
		pathAbs, //try direct path first
		filepath.Join(workDir, filepath.Dir(path)), //default
		//filepath.Join(rt_directory, "./", filepath.Dir(path)),
		filepath.Join(rtDirectory, "./", filepath.Dir(path)),
		filepath.Join(rtDirectory, "../", filepath.Dir(path)),
		filepath.Join(rtDirectory, "../../", filepath.Dir(path)),
		filepath.Join(rtDirectory, "../../../", filepath.Dir(path)),
	}

	//for i, dir := range dirs {
	//	fmt.Printf("Dir[%d]= %s \n", i, dir)
	//}

	//must be an absolute path
	//error and problem and crash if not absolute path
	for i := range dirs {
		absPath, _ := filepath.Abs(dirs[i])
		dirs[i] = absPath
	}

	for _, dir := range dirs {
		if _, err := os.Stat(dir); !os.IsNotExist(err) {
			fmt.Printf("ResolveResourceDirectory: static resource dir= %s \n", dir)
			return dir
		}
	}
	logger.Panic("GUI directory not found")
	return ""
}

// DetermineResourcePath DEPRECATE
// From src/gui/http.go and src/mesh/gui/http.go
func DetermineResourcePath(staticDir string, resourceDir string, devDir string) (string, error) {
	//check "dev" directory first
	appLoc := filepath.Join(staticDir, devDir)
	// if !strings.HasPrefix(appLoc, "/") {
	// 	// Prepend the binary's directory path if appLoc is relative
	// 	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	// 	if err != nil {
	// 		return "", err
	// 	}

	// 	appLoc = filepath.Join(dir, appLoc)
	// }
	if _, err := os.Stat(appLoc); os.IsNotExist(err) {
		//check dist directory
		appLoc = filepath.Join(staticDir, resourceDir)
		// if !strings.HasPrefix(appLoc, "/") {
		// 	// Prepend the binary's directory path if appLoc is relative
		// 	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		// 	if err != nil {
		// 		return "", err
		// 	}

		// 	appLoc = filepath.Join(dir, appLoc)
		// }

		if _, err := os.Stat(appLoc); os.IsNotExist(err) {
			return "", err
		}
	}

	return appLoc, nil
}

// CopyFile copy file
func CopyFile(dst string, src io.Reader) (n int64, err error) {
	// check the existence of dst file.
	if _, err := os.Stat(dst); err == nil {
		return 0, nil
	}
	err = nil

	out, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()

	n, err = io.Copy(out, src)
	return
}
