package logging

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Logger wraps logrus.FieldLogger
type Logger struct {
	logrus.FieldLogger
}

// Critical adds special critical-level fields for specially highlighted logging,
// since logrus lacks a distinct critical field and does not have configurable log levels
func (logger *Logger) Critical() logrus.FieldLogger {
	return logger.WithField(logPriorityKey, logPriorityCritical)
}

// MasterLogger wraps logrus.Logger and is able to create new package-aware loggers
// sharing configuration realm.
type MasterLogger struct {
	*logrus.Logger
	PkgConfig map[string]PkgLogConfig
}

type PkgLogConfig struct {
	PkgName string
	Level   logrus.Level
}

// NewMasterLogger creates a new package-aware logger with formatting string
func NewMasterLogger() *MasterLogger {
	hooks := make(logrus.LevelHooks)

	return &MasterLogger{
		Logger: &logrus.Logger{
			Out: os.Stdout,
			Formatter: &TextFormatter{
				FullTimestamp:      true,
				AlwaysQuoteStrings: true,
				QuoteEmptyFields:   true,
				ForceFormatting:    true,
				DisableColors:      false,
				ForceColors:        false,
			},
			Hooks: hooks,
			Level: logrus.DebugLevel,
		},
		PkgConfig: make(map[string]PkgLogConfig, 5),
	}
}

func copyMasterLogger(logger *MasterLogger) *MasterLogger {
	return &MasterLogger{
		Logger: &logrus.Logger{
			Out:       logger.Out,
			Formatter: logger.Formatter,
			Hooks:     logger.Hooks,
			Level:     logger.Level,
		},
	}
}

// PackageLogger instantiates a package-aware logger
func (logger *MasterLogger) PackageLogger(moduleName string) *Logger {
	if pkgcfg, ok := logger.PkgConfig[moduleName]; ok {
		logger = copyMasterLogger(logger)
		logger.SetLevel(pkgcfg.Level)
	}
	return &Logger{
		FieldLogger: logger.WithField(logModuleKey, moduleName),
	}
}

// AddHook adds a logrus.Hook to the logger and its module loggers
func (logger *MasterLogger) AddHook(hook logrus.Hook) {
	logger.Hooks.Add(hook)
}

// SetLevel sets the log level for the logger and its module loggers
func (logger *MasterLogger) SetLevel(level logrus.Level) {
	logger.Level = level
}

// EnableColors enables colored logging
func (logger *MasterLogger) EnableColors() {
	logger.Formatter.(*TextFormatter).DisableColors = false
}

// DisableColors disables colored logging
func (logger *MasterLogger) DisableColors() {
	logger.Formatter.(*TextFormatter).DisableColors = true
}
