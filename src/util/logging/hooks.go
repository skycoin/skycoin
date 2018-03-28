package logging

import (
	"reflect"

	"github.com/sirupsen/logrus"
)

func NewReplayHook(logger *Logger) ReplayHook {
	return ExclusiveReplayHook(logger, []reflect.Type{
		reflect.TypeOf(ReplayHook{}),
		reflect.TypeOf(ModuleLogHook{}),
	})
}

// Do not replay hooks of given exclude types
func ExclusiveReplayHook(logger *Logger, exclude []reflect.Type) (h ReplayHook) {
	h = ReplayHook{
		Logger:       logger,
		excludeTypes: make(map[reflect.Type]struct{}, len(exclude)),
	}
	for _, _type := range exclude {
		h.excludeTypes[_type] = struct{}{}
	}
	return
}

// Hook for replaying hooks bound to another logger
type ReplayHook struct {
	Logger       *Logger
	excludeTypes map[reflect.Type]struct{}
}

func (h ReplayHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

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

// Tag log entries with module information
type ModuleLogHook struct {
	FieldKey    string
	PriorityKey string
	ModuleName  string
}

func NewModuleLogHook(moduleName string) logrus.Hook {
	return ModuleLogHook{
		FieldKey:    "module",
		PriorityKey: "priority",
		ModuleName:  moduleName,
	}
}

func (h ModuleLogHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

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
