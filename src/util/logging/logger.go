package logging

import (
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
)

// ExtendedFieldLogger is an enhanced logger supporting critical and important levels
type ExtendedFieldLogger interface {
	logrus.FieldLogger

	// Critical level
	Criticalf(format string, args ...interface{})
	Critical(args ...interface{})
	Criticalln(args ...interface{})

	// Notice level
	Noticef(format string, args ...interface{})
	Notice(args ...interface{})
	Noticeln(args ...interface{})
}

// Logger wraps sirupsen/logrus.Logger to implement ExtendendFieldLogger
type Logger struct {
	*logrus.Logger
	formatter         *TextFormatter
	module            string
	allModulesEnabled bool
	moduleLoggers     map[string]*Logger
	PriorityKey       string
	CriticalPriority  string
}

var (
	// QuietLogger disables all log output
	QuietLogger = logrus.Logger{
		Out:       ioutil.Discard,
		Formatter: new(logrus.TextFormatter), // FIXME: Performance?
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.FatalLevel,
	}
)

// NewLogger creates a new modules-aware logger with formatting string
func NewLogger(priorityKey, criticalPriority string) (logger *Logger) {
	formatter := &TextFormatter{
		FullTimestamp:          true,
		AlwaysQuoteStrings:     true,
		QuoteEmptyFields:       true,
		ForceFormatting:        true,
		PriorityKey:            priorityKey,
		HighlightPriorityValue: criticalPriority,
	}

	logger = &Logger{
		Logger: &logrus.Logger{
			Out:       os.Stdout,
			Formatter: formatter,
			Hooks:     make(logrus.LevelHooks),
			Level:     logrus.InfoLevel,
		},
		formatter:         formatter,
		allModulesEnabled: true,
		moduleLoggers:     make(map[string]*Logger),
		PriorityKey:       priorityKey,
		CriticalPriority:  criticalPriority,
	}
	logger.Hooks.Add(NewModuleLogHook(""))
	logger.moduleLoggers[""] = logger
	return
}

func (logger *Logger) cloneForModule(moduleName string) *Logger {
	newLogger := &Logger{
		Logger: &logrus.Logger{
			Out:       logger.Out,
			Formatter: logger.Formatter,
			Hooks:     make(logrus.LevelHooks),
			Level:     logger.Level,
		},
		allModulesEnabled: logger.allModulesEnabled,
		moduleLoggers:     logger.moduleLoggers,
		PriorityKey:       logger.PriorityKey,
		CriticalPriority:  logger.CriticalPriority,
	}
	newLogger.Hooks.Add(NewReplayHook(logger.moduleLoggers[""]))
	newLogger.Hooks.Add(NewModuleLogHook(moduleName))
	return newLogger
}

// MustGetLogger returns an existing logger for a given module or creates a new one
func (logger *Logger) MustGetLogger(moduleName string) *Logger {
	newLogger, isInCache := logger.moduleLoggers[moduleName]
	if !(isInCache && newLogger != nil) {
		if isInCache || logger.allModulesEnabled {
			newLogger = logger.cloneForModule(moduleName)
			logger.moduleLoggers[moduleName] = newLogger
		} else {
			newLogger = &Logger{
				Logger:            &QuietLogger,
				allModulesEnabled: logger.allModulesEnabled,
				moduleLoggers:     logger.moduleLoggers,
			}
		}
	}
	return newLogger
}

// AddHook adds a logrus.Hook to the logger and its module loggers
func (logger *Logger) AddHook(hook logrus.Hook) {
	logger.Logger.Hooks.Add(hook)
	for _, module := range logger.moduleLoggers {
		module.Logger.Hooks.Add(hook)
	}
}

// SetLevel sets the log level for the logger and its module loggers
func (logger *Logger) SetLevel(level logrus.Level) {
	logger.Logger.Level = level
	for _, module := range logger.moduleLoggers {
		module.Logger.Level = level
	}
}

// EnableColors enables colored logging
func (logger *Logger) EnableColors() {
	logger.formatter.DisableColors = false
}

// DisableColors disables colored logging
func (logger *Logger) DisableColors() {
	logger.formatter.DisableColors = true
}

// Criticalf formatted log with Critical priority
func (logger *Logger) Criticalf(format string, args ...interface{}) {
	logger.WithField(logger.PriorityKey, logger.CriticalPriority).Error(args...)
}

// Critical log with Critical priority
func (logger *Logger) Critical(args ...interface{}) {
	logger.WithField(logger.PriorityKey, logger.CriticalPriority).Error(args...)
}

// Criticalln log line with Critical priority
func (logger *Logger) Criticalln(args ...interface{}) {
	logger.WithField(logger.PriorityKey, logger.CriticalPriority).Error(args...)
}

// Noticef formatted log with Critical priority
func (logger *Logger) Noticef(format string, args ...interface{}) {
	logger.WithField(logger.PriorityKey, logger.CriticalPriority).Info(args...)
}

// Notice log with Critical priority
func (logger *Logger) Notice(args ...interface{}) {
	logger.WithField(logger.PriorityKey, logger.CriticalPriority).Info(args...)
}

// Noticeln log line with Critical priority
func (logger *Logger) Noticeln(args ...interface{}) {
	logger.WithField(logger.PriorityKey, logger.CriticalPriority).Info(args...)
}

// Disable discards all log output
func (logger *Logger) Disable() {
	logger.Out = ioutil.Discard
}
