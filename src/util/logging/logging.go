package logging

import (
	"io"
)

var log = NewLogger()

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
