package logging

import (
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
)

// An enhanced logger supporting critical and important levels
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
	module            string
	allModulesEnabled bool
	moduleLoggers     map[string]*Logger
}

var (
	QuietLogger = logrus.Logger{
		Out:       ioutil.Discard,
		Formatter: new(logrus.TextFormatter), // FIXME: Performance?
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.FatalLevel,
	}
)

// New modules-aware logger with formatting string
func NewLogger() (logger *Logger) {
	logger = &Logger{
		Logger: &logrus.Logger{
			Out: os.Stderr,
			Formatter: &TextFormatter{
				FullTimestamp:      true,
				AlwaysQuoteStrings: true,
				QuoteEmptyFields:   true,
				ForceFormatting:    true,
				PriorityKey:        "priority",
			},
			Hooks: make(logrus.LevelHooks),
			Level: logrus.InfoLevel,
		},
		allModulesEnabled: true,
		moduleLoggers:     make(map[string]*Logger),
	}
	logger.Hooks.Add(NewModuleLogHook(""))
	logger.moduleLoggers[""] = logger
	return
}

func LoggerForModules(enabledModules []string) *Logger {
	logger := NewLogger()
	logger.allModulesEnabled = false
	for _, moduleName := range enabledModules {
		// Lazy instantiation
		logger.moduleLoggers[moduleName] = nil
	}
	return logger
}

func (l *Logger) cloneForModule(moduleName string) (logger *Logger) {
	logger = &Logger{
		Logger: &logrus.Logger{
			Out:       l.Out,
			Formatter: l.Formatter,
			Hooks:     make(logrus.LevelHooks),
			Level:     l.Level,
		},
		allModulesEnabled: l.allModulesEnabled,
		moduleLoggers:     l.moduleLoggers,
	}
	logger.Hooks.Add(NewReplayHook(l.moduleLoggers[""]))
	logger.Hooks.Add(NewModuleLogHook(moduleName))
	return
}

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

func (logger *Logger) DisableAllModules() {
	logger.allModulesEnabled = false
	for k := range logger.moduleLoggers {
		delete(logger.moduleLoggers, k)
	}
}

func (logger *Logger) EnableModules(modules ...string) {
	if logger.allModulesEnabled {
		return
	}
	for _, moduleName := range modules {
		if _, isBound := logger.moduleLoggers[moduleName]; isBound {
			// Lazy instantiation
			logger.moduleLoggers[moduleName] = nil
		}
	}
}

func (l *Logger) Criticalf(format string, args ...interface{}) {
	l.WithField("priority", "CRITICAL").Error(args...)
}

func (l *Logger) Critical(args ...interface{}) {
	l.WithField("priority", "CRITICAL").Error(args...)
}

func (l *Logger) Criticalln(args ...interface{}) {
	l.WithField("priority", "CRITICAL").Error(args...)
}

func (l *Logger) Noticef(format string, args ...interface{}) {
	l.WithField("priority", "CRITICAL").Info(args...)
}

func (l *Logger) Notice(args ...interface{}) {
	l.WithField("priority", "CRITICAL").Info(args...)
}

func (l *Logger) Noticeln(args ...interface{}) {
	l.WithField("priority", "CRITICAL").Info(args...)
}

func (l *Logger) Disable() {
	l.Out = ioutil.Discard
}
