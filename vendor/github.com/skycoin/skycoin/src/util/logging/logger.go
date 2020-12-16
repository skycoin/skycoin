package logging

import (
	"os"
	"time"

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

// WithTime overrides time, used by logger.
func (logger *Logger) WithTime(t time.Time) *logrus.Entry {
	return logger.WithFields(logrus.Fields{}).WithTime(t)
}

// MasterLogger wraps logrus.Logger and is able to create new package-aware loggers
type MasterLogger struct {
	*logrus.Logger
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
	}
}

// PackageLogger instantiates a package-aware logger
func (logger *MasterLogger) PackageLogger(moduleName string) *Logger {
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
