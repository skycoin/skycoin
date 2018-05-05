package logging

import (
	"bytes"
	"fmt"
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
	return append([]byte(text), '\n'), nil
}

// Independent log levels for packages. Output redirects
func TestPkgLevelIo(t *testing.T) {
	var buff1, buff2 bytes.Buffer

	// Empty buffers
	require.Equal(t, "", buff1.String())
	require.Equal(t, "", buff2.String())

	var formatter TestMsgCollector

	// A master logger logging to a memory buffer at WARN level
	log := NewMasterLogger()
	log.Formatter = &formatter
	log.Out = &OutputRealm{Writer: &buff1}
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

	checkAllLevels := func() {
		// Clear previous log history
		formatter.Messages = nil

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

	checkAllLevels()
	logOutput := buff1.String()
	require.NotEqual(t, "", logOutput)
	require.Equal(t, "", buff2.String())

	setLoggerOutput(log, &buff2)
	checkAllLevels()
	require.Equal(t, logOutput, buff1.String())
	require.Equal(t, logOutput, buff2.String())

	var buff3 bytes.Buffer
	setLoggerOutput(log, &buff3)
	disableLoggerOutput(log)
	checkAllLevels()
	require.Equal(t, "", buff3.String())
}
