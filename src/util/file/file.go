// Package file provides filesystem related utilities
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

	"github.com/SkycoinProject/skycoin/src/util/logging"
)

var (
	// ErrEmptyDirectoryName is returned by constructing the full path
	// of data directory if the passed argument is empty
	ErrEmptyDirectoryName = errors.New("data directory must not be empty")
	// ErrDotDirectoryName is returned by constructing the full path of
	// data directory if the passed argument is "."
	ErrDotDirectoryName = errors.New("data directory must not be equal to \".\"")

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

	// check if dir already exists
	st, err := os.Stat(dir)
	if !os.IsNotExist(err) {
		if !st.IsDir() {
			return "", fmt.Errorf("%s is not a directory", dir)
		}
		// dir already exists
		return dir, nil
	}

	if err := os.MkdirAll(dir, os.FileMode(0700)); err != nil {
		logger.Errorf("Failed to create directory %s: %v", dir, err)
		return "", err
	}

	logger.Infof("Created data directory %s", dir)
	return dir, nil
}

// Construct the full data directory by adding to $HOME or ./
func buildDataDir(dir string) (string, error) {
	if dir == "" {
		logger.Error("data directory is empty")
		return "", ErrEmptyDirectoryName
	}

	home := filepath.Clean(UserHome())
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	wd = filepath.Clean(wd)

	fullDir, err := filepath.Abs(dir)

	if err != nil {
		return "", err
	}

	// The joined directory must not be equal to $HOME or a parent path of $HOME
	// The joined directory must not be equal to `pwd` or a parent path of `pwd`
	if strings.HasPrefix(home, fullDir) || strings.HasPrefix(wd, fullDir) {
		logger.Errorf("join(%[1]s, %[2]s) == %[1]s", home, dir)
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
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	dec := json.NewDecoder(file)
	dec.UseNumber()
	return dec.Decode(thing)
}

// SaveJSON write value into json file
func SaveJSON(filename string, thing interface{}, mode os.FileMode) error {
	data, err := json.MarshalIndent(thing, "", "    ")
	if err != nil {
		return err
	}
	return SaveBinary(filename, data, mode)
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
		if removeErr := os.Remove(filename); removeErr != nil {
			logger.WithError(removeErr).Warningf("os.Remove(%s) failed", filename)
		}
	}
	return err
}

// SaveBinary persists data into given file in binary,
// backup the previous file, if there was one
func SaveBinary(filename string, data []byte, mode os.FileMode) error {
	// Write the new file to a temporary
	tmpname := filename + ".tmp"
	if err := ioutil.WriteFile(tmpname, data, mode); err != nil {
		return err
	}

	// Write the new file to the target wallet file
	if err := ioutil.WriteFile(filename, data, mode); err != nil {
		return err
	}

	// Remove the tmp file
	return os.Remove(tmpname)
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
		absPath, err := filepath.Abs(dirs[i])
		if err != nil {
			logger.WithError(err).Errorf("filepath.Abs(%s) failed", dirs[i])
			continue
		}
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

// Copy copies file. Will overwrite dst if dst exists.
func Copy(dst, src string) (err error) {
	f, err := os.Open(src)
	if err != nil {
		return err
	}

	defer func() {
		cerr := f.Close()
		if err == nil {
			err = cerr
		}
	}()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}

	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()

	_, err = io.Copy(out, f)
	return
}

// Exists checks whether the file exists in the file system
func Exists(fn string) (bool, error) {
	_, err := os.Stat(fn)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// IsWritable checks if the file is writable
func IsWritable(fn string) bool {
	st, _ := os.Stat(fn)
	return (st.Mode().Perm()&(1<<uint(7)) == 0)
}
