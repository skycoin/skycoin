package logging

import (
	"fmt"
	"io"
	"os"
	"strings"

	prefixed "github.com/gz-c/logrus-prefixed-formatter"
	"github.com/sirupsen/logrus"
)

const (
	defaultLogFormat = "[%{module}:%{level}] %{message}"
)

// Level embedes the logging's level
type Level uint32

// Log levels. Based (approximately) on Python defaults
const (
	QUIET Level = 10 * iota // No logging
	PANIC
	FATAL
	CRITICAL
	ERROR
	WARNING
	NOTICE
	INFO
	DEBUG
)

func (lvl Level) toImplLevel() (logrus.Level, error) {
	switch lvl {
	case PANIC:
		return logrus.PanicLevel, nil
	case FATAL:
		return logrus.FatalLevel, nil
	case ERROR, CRITICAL:
		return logrus.ErrorLevel, nil
	case WARNING:
		return logrus.WarnLevel, nil
	case INFO, NOTICE:
		return logrus.InfoLevel, nil
	case DEBUG:
		return logrus.DebugLevel, nil
	}
	var l logrus.Level
	return l, fmt.Errorf("logrus.ing implementation does not support level: %q", lvl)
}

// LogLevel parse the log level string
func LogLevel(levelStr string) (Level, error) {
	switch strings.ToLower(levelStr) {
	case "panic":
		return PANIC, nil
	case "fatal":
		return FATAL, nil
	case "critical", "error":
		return ERROR, nil
	case "warn", "warning":
		return WARNING, nil
	case "notice", "info":
		return INFO, nil
	case "debug":
		return DEBUG, nil
	}

	var l Level
	return l, fmt.Errorf("not a valid logging Level: %q", levelStr)
}

// Return the
func (l Level) String() string {
	switch l {
	case PANIC:
		return "PANIC"
	case FATAL:
		return "FATAL"
	case CRITICAL:
		return "CRITICAL"
	case ERROR:
		return "ERROR"
	case WARNING:
		return "WARN"
	case INFO:
		return "INFO"
	case NOTICE:
		return "NOTICE"
	case DEBUG:
		return "DEBUG"
	}

	return fmt.Sprintf("LOGLEVEL%d", l)
}

// LogConfig logger configurations
type LogConfig struct {
	// for internal usage
	level Level
	// Level convertes to level during initialization
	Level string
	// list of all modules
	Modules []string
	// format
	Format string
	// enable colors
	Colors bool
	// output
	Output io.Writer
}

// DefaultLogConfig vs (DevLogConfig + ProdLogConfig) ?

// DevLogConfig default development config for logging
func DevLogConfig(modules []string) *LogConfig {
	return &LogConfig{
		level:   DEBUG,   // int
		Level:   "debug", // string
		Modules: modules,
		Format:  defaultLogFormat,
		Colors:  true,
		Output:  os.Stdout,
	}
}

// ProdLogConfig Default production config for logging
func ProdLogConfig(modules []string) *LogConfig {
	return &LogConfig{
		level:   ERROR,
		Level:   "error",
		Modules: modules,
		Format:  defaultLogFormat,
		Colors:  false,
		Output:  os.Stdout,
	}
}

// convertes l.Level (string) to l.level (int)
// or panics if l.Level is invalid
func (l *LogConfig) initLevel() {
	level, err := LogLevel(l.Level)
	if err != nil {
		log.Panicf("Invalid -log-level %s: %v", l.Level, err)
	}
	l.level = Level(level)
}

var log = NewLogger()

// InitLogger initialize logging using this LogConfig;
// it panics if l.Format is invalid or l.Level is invalid
func (l *LogConfig) InitLogger() {
	l.initLevel()

	formatter := prefixed.TextFormatter{
		FullTimestamp:      true,
		AlwaysQuoteStrings: true,
		QuoteEmptyFields:   true,
		ForceFormatting:    true,
	}
	formatter.ForceColors = l.Colors
	formatter.DisableColors = !l.Colors
	log.Formatter = &formatter

	log.Out = l.Output
	if level, err := l.level.toImplLevel(); err == nil {
		log.SetLevel(level)
	}

	log.DisableAllModules()
	log.EnableModules(l.Modules...)
}

// MustGetLogger safe initialize global logger
func MustGetLogger(module string) *Logger {
	return log.MustGetLogger(module)
}

// Disable disables the logger completely
func Disable() {
	for k := range log.moduleLoggers {
		log.moduleLoggers[k].Disable()
	}
}

// Disable disables the logger completely
func RedirectTo(w io.Writer) {
	for k := range log.moduleLoggers {
		log.moduleLoggers[k].Out = w
	}
}
