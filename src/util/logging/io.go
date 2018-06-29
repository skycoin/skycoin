package logging

import (
	"io"
	"io/ioutil"
)

// OutputRealm Proxy object controlling
type OutputRealm struct {
	io.Writer
}

func (w *OutputRealm) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// SetOutputTo redirects output to an io.Writer for all clients connected to this realm
func (w *OutputRealm) SetOutputTo(_w io.Writer) {
	w.Writer = _w
}

// Disable output of all clients bound to this realm
func (w *OutputRealm) Disable() {
	w.Writer = ioutil.Discard
}

func setLoggerOutput(logger *MasterLogger, w io.Writer) {
	if realm, ok := logger.Out.(*OutputRealm); ok {
		realm.SetOutputTo(w)
	} else {
		logger.Out = w
	}
}

// SetOutputTo sets the logger's output to an io.Writer
func SetOutputTo(w io.Writer) {
	setLoggerOutput(log, w)
}

func disableLoggerOutput(logger *MasterLogger) {
	if realm, ok := logger.Out.(*OutputRealm); ok {
		realm.Disable()
	} else {
		logger.Out = ioutil.Discard
	}
}

// Disable disables the logger completely
func Disable() {
	disableLoggerOutput(log)
}
