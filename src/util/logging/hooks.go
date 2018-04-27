package logging

import (
	"io"

	"github.com/sirupsen/logrus"
)

// WriteHook is a logrus.Hook that logs to an io.Writer
type WriteHook struct {
	w         io.Writer
	formatter logrus.Formatter
}

// NewWriteHook returns a new WriteHook
func NewWriteHook(w io.Writer) *WriteHook {
	return &WriteHook{
		w: w,
		formatter: &TextFormatter{
			DisableColors:      true,
			FullTimestamp:      true,
			AlwaysQuoteStrings: true,
			QuoteEmptyFields:   true,
			ForceFormatting:    true,
		},
	}
}

// Levels returns Levels accepted by the WriteHook.
// All logrus.Levels are returned.
func (f *WriteHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Fire writes a logrus.Entry to the file
func (f *WriteHook) Fire(e *logrus.Entry) error {
	b, err := f.formatter.Format(e)
	if err != nil {
		return err
	}

	_, err = f.w.Write(b)
	return err
}
