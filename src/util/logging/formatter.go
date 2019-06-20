package logging

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/mgutz/ansi"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"
)

const defaultTimestampFormat = time.RFC3339

var (
	baseTimestamp      = time.Now()
	defaultColorScheme = &ColorScheme{
		InfoLevelStyle:   "green",
		WarnLevelStyle:   "yellow",
		ErrorLevelStyle:  "red",
		FatalLevelStyle:  "red",
		PanicLevelStyle:  "red",
		DebugLevelStyle:  "blue",
		PrefixStyle:      "cyan",
		TimestampStyle:   "black+h",
		CallContextStyle: "black+h",
		CriticalStyle:    "magenta+h",
	}
	noColorsColorScheme = &compiledColorScheme{
		InfoLevelColor:   ansi.ColorFunc(""),
		WarnLevelColor:   ansi.ColorFunc(""),
		ErrorLevelColor:  ansi.ColorFunc(""),
		FatalLevelColor:  ansi.ColorFunc(""),
		PanicLevelColor:  ansi.ColorFunc(""),
		DebugLevelColor:  ansi.ColorFunc(""),
		PrefixColor:      ansi.ColorFunc(""),
		TimestampColor:   ansi.ColorFunc(""),
		CallContextColor: ansi.ColorFunc(""),
		CriticalColor:    ansi.ColorFunc(""),
	}
	defaultCompiledColorScheme = compileColorScheme(defaultColorScheme)
)

func miniTS() int {
	return int(time.Since(baseTimestamp) / time.Second)
}

// ColorScheme configures the logging output colors
type ColorScheme struct {
	InfoLevelStyle   string
	WarnLevelStyle   string
	ErrorLevelStyle  string
	FatalLevelStyle  string
	PanicLevelStyle  string
	DebugLevelStyle  string
	PrefixStyle      string
	TimestampStyle   string
	CallContextStyle string
	CriticalStyle    string
}

type compiledColorScheme struct {
	InfoLevelColor   func(string) string
	WarnLevelColor   func(string) string
	ErrorLevelColor  func(string) string
	FatalLevelColor  func(string) string
	PanicLevelColor  func(string) string
	DebugLevelColor  func(string) string
	PrefixColor      func(string) string
	TimestampColor   func(string) string
	CallContextColor func(string) string
	CriticalColor    func(string) string
}

// TextFormatter formats log output
type TextFormatter struct {
	// Set to true to bypass checking for a TTY before outputting colors.
	ForceColors bool

	// Force disabling colors. For a TTY colors are enabled by default.
	DisableColors bool

	// Force formatted layout, even for non-TTY output.
	ForceFormatting bool

	// Disable timestamp logging. useful when output is redirected to logging
	// system that already adds timestamps.
	DisableTimestamp bool

	// Disable the conversion of the log levels to uppercase
	DisableUppercase bool

	// Enable logging the full timestamp when a TTY is attached instead of just
	// the time passed since beginning of execution.
	FullTimestamp bool

	// Timestamp format to use for display when a full timestamp is printed.
	TimestampFormat string

	// The fields are sorted by default for a consistent output. For applications
	// that log extremely frequently and don't use the JSON formatter this may not
	// be desired.
	DisableSorting bool

	// Wrap empty fields in quotes if true.
	QuoteEmptyFields bool

	// Can be set to the override the default quoting character "
	// with something else. For example: ', or `.
	QuoteCharacter string

	// Pad msg field with spaces on the right for display.
	// The value for this parameter will be the size of padding.
	// Its default value is zero, which means no padding will be applied for msg.
	SpacePadding int

	// Always use quotes for string values (except for empty fields)
	AlwaysQuoteStrings bool

	// Color scheme to use.
	colorScheme *compiledColorScheme

	// Whether the logger's out is to a terminal.
	isTerminal bool

	sync.Once
}

func getCompiledColor(main string, fallback string) func(string) string {
	var style string
	if main != "" {
		style = main
	} else {
		style = fallback
	}
	return ansi.ColorFunc(style)
}

func compileColorScheme(s *ColorScheme) *compiledColorScheme {
	return &compiledColorScheme{
		InfoLevelColor:   getCompiledColor(s.InfoLevelStyle, defaultColorScheme.InfoLevelStyle),
		WarnLevelColor:   getCompiledColor(s.WarnLevelStyle, defaultColorScheme.WarnLevelStyle),
		ErrorLevelColor:  getCompiledColor(s.ErrorLevelStyle, defaultColorScheme.ErrorLevelStyle),
		FatalLevelColor:  getCompiledColor(s.FatalLevelStyle, defaultColorScheme.FatalLevelStyle),
		PanicLevelColor:  getCompiledColor(s.PanicLevelStyle, defaultColorScheme.PanicLevelStyle),
		DebugLevelColor:  getCompiledColor(s.DebugLevelStyle, defaultColorScheme.DebugLevelStyle),
		PrefixColor:      getCompiledColor(s.PrefixStyle, defaultColorScheme.PrefixStyle),
		TimestampColor:   getCompiledColor(s.TimestampStyle, defaultColorScheme.TimestampStyle),
		CallContextColor: getCompiledColor(s.CallContextStyle, defaultColorScheme.CallContextStyle),
		CriticalColor:    getCompiledColor(s.CriticalStyle, defaultColorScheme.CriticalStyle),
	}
}

func (f *TextFormatter) init(entry *logrus.Entry) {
	if len(f.QuoteCharacter) == 0 {
		f.QuoteCharacter = "\""
	}
	if entry.Logger != nil {
		f.isTerminal = f.checkIfTerminal(entry.Logger.Out)
	}
}

func (f *TextFormatter) checkIfTerminal(w io.Writer) bool {
	switch v := w.(type) {
	case *os.File:
		return terminal.IsTerminal(int(v.Fd()))
	default:
		return false
	}
}

// SetColorScheme sets the TextFormatter's color scheme configuration
func (f *TextFormatter) SetColorScheme(colorScheme *ColorScheme) {
	f.colorScheme = compileColorScheme(colorScheme)
}

// Format formats a logrus.Entry
func (f *TextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	keys := make([]string, 0, len(entry.Data))
	for k := range entry.Data {
		keys = append(keys, k)
	}
	lastKeyIdx := len(keys) - 1

	if !f.DisableSorting {
		sort.Strings(keys)
	}
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	f.Do(func() { f.init(entry) })

	isFormatted := f.ForceFormatting || f.isTerminal

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = defaultTimestampFormat
	}
	if isFormatted {
		isColored := (f.ForceColors || f.isTerminal) && !f.DisableColors
		var colorScheme *compiledColorScheme
		if isColored {
			if f.colorScheme == nil {
				colorScheme = defaultCompiledColorScheme
			} else {
				colorScheme = f.colorScheme
			}
		} else {
			colorScheme = noColorsColorScheme
		}
		f.printColored(b, entry, keys, timestampFormat, colorScheme)
	} else {
		if !f.DisableTimestamp {
			f.appendKeyValue(b, "time", entry.Time.Format(timestampFormat), true)
		}
		f.appendKeyValue(b, "level", entry.Level.String(), true)
		if entry.Message != "" {
			f.appendKeyValue(b, "msg", entry.Message, lastKeyIdx >= 0)
		}
		for i, key := range keys {
			f.appendKeyValue(b, key, entry.Data[key], lastKeyIdx != i)
		}
	}

	b.WriteByte('\n') //nolint:gosec
	return b.Bytes(), nil
}

func (f *TextFormatter) printColored(b *bytes.Buffer, entry *logrus.Entry, keys []string, timestampFormat string, colorScheme *compiledColorScheme) {
	var levelColor func(string) string
	var levelText string
	switch entry.Level {
	case logrus.InfoLevel:
		levelColor = colorScheme.InfoLevelColor
	case logrus.WarnLevel:
		levelColor = colorScheme.WarnLevelColor
	case logrus.ErrorLevel:
		levelColor = colorScheme.ErrorLevelColor
	case logrus.FatalLevel:
		levelColor = colorScheme.FatalLevelColor
	case logrus.PanicLevel:
		levelColor = colorScheme.PanicLevelColor
	default:
		levelColor = colorScheme.DebugLevelColor
	}

	priority, ok := entry.Data[logPriorityKey]
	hasPriority := ok && priority == logPriorityCritical

	if entry.Level != logrus.WarnLevel {
		levelText = entry.Level.String()
	} else {
		levelText = "warn"
	}

	if !f.DisableUppercase {
		levelText = strings.ToUpper(levelText)
	}

	level := levelColor(levelText)
	message := entry.Message
	prefix := ""

	prefixText := extractPrefix(entry)
	if prefixText != "" {
		prefixText = " " + prefixText + ":"
		prefix = colorScheme.PrefixColor(prefixText)
	}

	messageFormat := "%s"
	if f.SpacePadding != 0 {
		messageFormat = fmt.Sprintf("%%-%ds", f.SpacePadding)
	}
	if message != "" {
		messageFormat = " " + messageFormat
	}

	callContextParts := []string{}
	if ifile, ok := entry.Data["file"]; ok {
		if sfile, ok := ifile.(string); ok && sfile != "" {
			callContextParts = append(callContextParts, sfile)
		}
	}
	if ifunc, ok := entry.Data["func"]; ok {
		if sfunc, ok := ifunc.(string); ok && sfunc != "" {
			callContextParts = append(callContextParts, sfunc)
		}
	}
	if iline, ok := entry.Data["line"]; ok {
		sline := ""
		switch iline := iline.(type) {
		case string:
			sline = iline
		case int, uint, int32, int64, uint32, uint64:
			sline = fmt.Sprint(iline)
		}
		if sline != "" {
			callContextParts = append(callContextParts, fmt.Sprint(sline))
		}
	}
	callContextText := strings.Join(callContextParts, ":")
	callContext := colorScheme.CallContextColor(callContextText)
	if callContext != "" {
		callContext = " " + callContext
	}

	if f.DisableTimestamp {
		if hasPriority {
			str := fmt.Sprintf("%s%s%s"+messageFormat, levelText, callContextText, prefixText, message)
			fmt.Fprint(b, colorScheme.CriticalColor(str))
		} else {
			fmt.Fprintf(b, "%s%s%s"+messageFormat, level, callContext, prefix, message)
		}
	} else {
		var timestamp string
		if !f.FullTimestamp {
			timestamp = fmt.Sprintf("[%04d]", miniTS())
		} else {
			timestamp = fmt.Sprintf("[%s]", entry.Time.Format(timestampFormat))
		}

		coloredTimestamp := colorScheme.TimestampColor(timestamp)

		if hasPriority {
			str := fmt.Sprintf("%s %s%s%s"+messageFormat, timestamp, levelText, callContextText, prefixText, message)
			fmt.Fprint(b, colorScheme.CriticalColor(str))
		} else {
			fmt.Fprintf(b, "%s %s%s%s"+messageFormat, coloredTimestamp, level, callContext, prefix, message)
		}
	}

	for _, k := range keys {
		if k != "prefix" && k != "file" && k != "func" && k != "line" && k != logPriorityKey && k != logModuleKey {
			v := entry.Data[k]
			fmt.Fprintf(b, " %s", f.formatKeyValue(levelColor(k), v))
		}
	}
}

func (f *TextFormatter) needsQuoting(text string) bool {
	if len(text) == 0 {
		return f.QuoteEmptyFields
	}

	if f.AlwaysQuoteStrings {
		return true
	}

	for _, ch := range text {
		if !((ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') ||
			ch == '-' || ch == '.') {
			return true
		}
	}

	return false
}

func extractPrefix(e *logrus.Entry) string {
	var module string
	if iModule, ok := e.Data[logModuleKey]; ok {
		module, _ = iModule.(string)
	}

	var priority string
	if iPriority, ok := e.Data[logPriorityKey]; ok {
		priority, _ = iPriority.(string)
	}

	switch {
	case priority == "":
		return fmt.Sprintf("[%s]", module)
	case module == "":
		return fmt.Sprintf("[%s]", priority)
	default:
		return fmt.Sprintf("[%s:%s]", module, priority)
	}
}

func (f *TextFormatter) formatKeyValue(key string, value interface{}) string {
	return fmt.Sprintf("%s=%s", key, f.formatValue(value))
}

func (f *TextFormatter) formatValue(value interface{}) string {
	switch value := value.(type) {
	case string:
		if f.needsQuoting(value) {
			return fmt.Sprintf("%s%+v%s", f.QuoteCharacter, value, f.QuoteCharacter)
		}
		return value
	case error:
		errmsg := value.Error()
		if f.needsQuoting(errmsg) {
			return fmt.Sprintf("%s%+v%s", f.QuoteCharacter, errmsg, f.QuoteCharacter)
		}
		return errmsg
	default:
		return fmt.Sprintf("%+v", value)
	}
}

func (f *TextFormatter) appendKeyValue(b *bytes.Buffer, key string, value interface{}, appendSpace bool) {
	b.WriteString(key)      //nolint:gosec
	b.WriteByte('=')        //nolint:gosec
	f.appendValue(b, value) //nolint:gosec

	if appendSpace {
		b.WriteByte(' ') //nolint:gosec
	}
}

func (f *TextFormatter) appendValue(b *bytes.Buffer, value interface{}) {
	switch value := value.(type) {
	case string:
		if f.needsQuoting(value) {
			fmt.Fprintf(b, "%s%+v%s", f.QuoteCharacter, value, f.QuoteCharacter)
		} else {
			b.WriteString(value) //nolint:gosec
		}
	case error:
		errmsg := value.Error()
		if f.needsQuoting(errmsg) {
			fmt.Fprintf(b, "%s%+v%s", f.QuoteCharacter, errmsg, f.QuoteCharacter)
		} else {
			b.WriteString(errmsg) //nolint:gosec
		}
	default:
		fmt.Fprint(b, value)
	}
}
