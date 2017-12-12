package logging

import (
	"io"
	"io/ioutil"
	"log"
	"os"

	logging "github.com/op/go-logging"
)

const (
	defaultLogFormat = "[%{module}:%{level}] %{message}"
)

// Logger wraps op/go-logging.Logger
type Logger struct {
	*logging.Logger
}

// Level embedes the logging's level
type Level int

// Log levels.
const (
	CRITICAL Level = iota
	ERROR
	WARNING
	NOTICE
	INFO
	DEBUG
)

var levelNames = []string{
	"CRITICAL",
	"ERROR",
	"WARNING",
	"NOTICE",
	"INFO",
	"DEBUG",
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

// LogLevel parse the log level string
func LogLevel(level string) (logging.Level, error) {
	return logging.LogLevel(level)
}

// TODO:
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
	level, err := logging.LogLevel(l.Level)
	if err != nil {
		log.Panicf("Invalid -log-level %s: %v", l.Level, err)
	}
	l.level = Level(level)
}

// InitLogger initialize logging using this LogConfig;
// it panics if l.Format is invalid or l.Level is invalid
func (l *LogConfig) InitLogger() {
	l.initLevel()

	format := logging.MustStringFormatter(l.Format)
	logging.SetFormatter(format)
	for _, s := range l.Modules {
		logging.SetLevel(logging.Level(l.level), s)
	}
	stdout := logging.NewLogBackend(l.Output, "", 0)
	stdout.Color = l.Colors
	logging.SetBackend(stdout)
}

// MustGetLogger safe initialize global logger
func MustGetLogger(module string) *Logger {
	return &Logger{logging.MustGetLogger(module)}
}

// Disable disables the logger completely
func Disable() {
	logging.SetBackend(logging.NewLogBackend(ioutil.Discard, "", 0))
}
