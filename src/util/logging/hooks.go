package logging

import (
	"reflect"

	"github.com/sirupsen/logrus"
)

// NewReplayHook creates a ReplayHook
func NewReplayHook(logger *Logger) ReplayHook {
	return newExclusiveReplayHook(logger, []reflect.Type{
		reflect.TypeOf(ReplayHook{}),
		reflect.TypeOf(ModuleLogHook{}),
	})
}

// newExclusiveReplayHook creates a ReplayHook that does not replay hooks of given excluded types
func newExclusiveReplayHook(logger *Logger, exclude []reflect.Type) (h ReplayHook) {
	h = ReplayHook{
		Logger:       logger,
		excludeTypes: make(map[reflect.Type]struct{}, len(exclude)),
	}
	for _, _type := range exclude {
		h.excludeTypes[_type] = struct{}{}
	}
	return
}

// ReplayHook is a hook for replaying hooks bound to another logger
type ReplayHook struct {
	Logger       *Logger
	excludeTypes map[reflect.Type]struct{}
}

// Levels returns all levels
func (h ReplayHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Fire fires other hooks that are not excluded from replay
func (h ReplayHook) Fire(entry *logrus.Entry) error {
	hooks := h.Logger.Hooks
	level := entry.Level
	for _, hook := range hooks[level] {
		if _, ok := h.excludeTypes[reflect.TypeOf(hook)]; ok {
			continue
		}
		if err := hook.Fire(entry); err != nil {
			return err
		}
	}
	return hooks.Fire(entry.Level, entry)
}

// ModuleLogHook tags log entries with module information
type ModuleLogHook struct {
	FieldKey    string
	PriorityKey string
	ModuleName  string
}

// NewModuleLogHook creates a ModuleLogHook
func NewModuleLogHook(moduleName string) ModuleLogHook {
	return ModuleLogHook{
		FieldKey:    LogModuleKey,
		PriorityKey: LogPriorityKey,
		ModuleName:  moduleName,
	}
}

// Levels returns all levels
func (h ModuleLogHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Fire adds module prefix to the logrus.Entry data
func (h ModuleLogHook) Fire(entry *logrus.Entry) error {
	entry.Data[h.FieldKey] = h.ModuleName
	prefix := h.ModuleName
	if value, hasField := entry.Data[h.PriorityKey]; hasField && value.(string) != "" {
		prefix += ":" + value.(string)
	}
	if prefix != "" {
		entry.Data["prefix"] = "[" + prefix + "]"
	}

	return nil
}
