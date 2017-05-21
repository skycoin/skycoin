//
// logging helpers for cmd/mesh and others
//

// TODO: move other packages to use this

package util

import (
	"log"
	"os"

	logging "github.com/op/go-logging"
)

const (
	defaultLogFormat = "[%{module}:%{level}] %{message}"
)

// MustGetLogger wrapper for logging.MustGetLogger to avoid import
func MustGetLogger(moduleName string) *logging.Logger {
	// may be some stuff here (or may be not)
	return logging.MustGetLogger(moduleName)
}

// LogConfig logger configurations
type LogConfig struct {
	// for internal usage
	level logging.Level
	// Level convertes to level during initialization
	Level string
	// list of all modules
	Modules []string
	// format
	Format string
	// enable colors
	Colors bool
}

// TODO:
// DefaultLogConfig vs (DevLogConfig + ProdLogConfig) ?

// DevLogConfig default development config for logging
func DevLogConfig(modules []string) *LogConfig {
	return &LogConfig{
		level:   logging.DEBUG, // int
		Level:   "debug",       // string
		Modules: modules,
		Format:  defaultLogFormat,
		Colors:  true,
	}
}

// ProdLogConfig Default production config for logging
func ProdLogConfig(modules []string) *LogConfig {
	return &LogConfig{
		level:   logging.ERROR,
		Level:   "error",
		Modules: modules,
		Format:  defaultLogFormat,
		Colors:  false,
	}
}

// convertes l.Level (string) to l.level (int)
// or panics if l.Level is invalid
func (l *LogConfig) initLevel() {
	level, err := logging.LogLevel(l.Level)
	if err != nil {
		log.Panicf("Invalid -log-level %s: %v", l.Level, err)
	}
	l.level = level
}

// InitLogger initialize logging using this LogConfig;
// it panics if l.Format is invalid or l.Level is invalid
func (l *LogConfig) InitLogger() {
	l.initLevel()

	format := logging.MustStringFormatter(l.Format)
	logging.SetFormatter(format)
	for _, s := range l.Modules {
		logging.SetLevel(l.level, s)
	}
	stdout := logging.NewLogBackend(os.Stdout, "", 0)
	stdout.Color = l.Colors
	logging.SetBackend(stdout)
}
