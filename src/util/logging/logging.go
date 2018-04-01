package logging

import (
	"io"
)

var log = NewLogger(LogPriorityKey, LogPriorityCritical)

const (
	// LogModuleKey is the key used for the module name data entry
	LogModuleKey = "module"
	// LogPriorityKey is the log entry key for priority log statements
	LogPriorityKey = "priority"
	// LogPriorityCritical is the log entry value for priority log statements
	LogPriorityCritical = "CRITICAL"
)

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

// RedirectTo redirects log to the given io.Wirter
func RedirectTo(w io.Writer) {
	for k := range log.moduleLoggers {
		log.moduleLoggers[k].Out = w
	}
}
