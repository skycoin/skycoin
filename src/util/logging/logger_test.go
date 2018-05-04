package logging

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

type TestMsgCollector struct {
	Messages []string
}

func (w *TestMsgCollector) Format(entry *logrus.Entry) ([]byte, error) {
	moduleName, ok := entry.Data[logModuleKey]
	if ok {
		moduleName = fmt.Sprintf("[%s]", moduleName)
	} else {
		moduleName = ""
	}
	text := fmt.Sprintf("%s %s: %s", strings.ToUpper(entry.Level.String()), moduleName, entry.Message)
	w.Messages = append(w.Messages, text)
	return []byte{}, nil
}

// Independent log levels for packages
func TestPackageLevel(t *testing.T) {
	// A master logger logging to a memory buffer at WARN level
	var formatter TestMsgCollector

	log := NewMasterLogger()
	log.Formatter = &formatter
	log.Out = ioutil.Discard
	log.Level = logrus.WarnLevel

	// Configure package logging levels
	ConfigPkgLogging(log, []PkgLogConfig{
		PkgLogConfig{"pkgdebug", logrus.DebugLevel},
		PkgLogConfig{"pkginfo", logrus.InfoLevel},
	})

	loggers := []*Logger{
		log.PackageLogger("pkgdebug"),
		log.PackageLogger("pkginfo"),
		log.PackageLogger("pkgnocfg"),
	}

	for _, logger := range loggers {
		logger.Error("Error")
		logger.Warn("Warn")
		logger.Info("Info")
		logger.Debug("Debug")
	}
	require.Equal(t, []string{
		"ERROR [pkgdebug]: Error",
		"WARNING [pkgdebug]: Warn",
		"INFO [pkgdebug]: Info",
		"DEBUG [pkgdebug]: Debug",
		"ERROR [pkginfo]: Error",
		"WARNING [pkginfo]: Warn",
		"INFO [pkginfo]: Info",
		"ERROR [pkgnocfg]: Error",
		"WARNING [pkgnocfg]: Warn",
	}, formatter.Messages)
}
