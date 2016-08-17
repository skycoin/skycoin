package gui

import (
	"fmt"
	"os"
	"time"

	"github.com/skycoin/skycoin/src/util"
)

// Returns true if both exist and are files, else false.
// Returns a slice of errors, indicating whether certFile and/or keyFile
// exist and are a file.  However, if neither exist, no error is returned.
func certKeyXor(certFile, keyFile string) (bool, []error) {
	certInfo, err := os.Stat(certFile)
	certExists := !os.IsNotExist(err)
	certIsFile := certExists && certInfo.Mode().IsRegular()

	keyInfo, err := os.Stat(keyFile)
	keyExists := !os.IsNotExist(err)
	keyIsFile := keyExists && keyInfo.Mode().IsRegular()

	errors := make([]error, 0)
	if certExists && certIsFile && keyExists && keyIsFile {
		return true, errors
	}
	if !certExists && !keyExists {
		return false, errors
	}
	if !certExists {
		errors = append(errors, fmt.Errorf("Cert %s does not exist", certFile))
	} else if !certIsFile {
		errors = append(errors, fmt.Errorf("Cert %s is not a file", certFile))
	}
	if !keyExists {
		errors = append(errors, fmt.Errorf("Key %s does not exist", keyFile))
	} else if !keyIsFile {
		errors = append(errors, fmt.Errorf("Key %s is not a file", keyFile))
	}
	return false, errors
}

// Checks that certFile and keyFile exist and are files, and if not,
// returns a slice of errors indicating status.
// If neither certFile nor keyFile exist, they are automatically created
// for host
func CreateCertIfNotExists(host, certFile, keyFile string) []error {
	// check that cert/key both exist, or dont
	exist, errs := certKeyXor(certFile, keyFile)
	// Automatically create a new cert if neither files exist
	if !exist && len(errs) == 0 {
		logger.Info("Creating certificate %s", certFile)
		logger.Info("Creating key %s", keyFile)
		err := util.GenerateCert(certFile, keyFile, host, "Skycoind", 2048,
			false, util.Now(), 365*24*time.Hour)
		if err == nil {
			logger.Info("Created certificate %s for host %s", certFile, host)
			logger.Info("Created key %s for host %s", keyFile, host)
		} else {
			errs = append(errs, err)
		}
	}
	return errs
}
