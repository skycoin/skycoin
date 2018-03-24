package logging

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

type TestLogScenario struct {
	Entry    *logrus.Entry
	Sequence []int
}

func (s TestLogScenario) Reset() {
	s.Entry = nil
	s.Sequence = nil
}

func (s TestLogScenario) NewHook(seq int) logrus.Hook {
	return &TestHook{
		Scenario:    &s,
		SeqNo:       seq,
		ReturnError: false,
	}
}

type TestHook struct {
	ReturnError bool
	Scenario    *TestLogScenario
	SeqNo       int
}

func (h *TestHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *TestHook) Fire(entry *logrus.Entry) (err error) {
	h.Scenario.Entry = entry
	seqno := h.SeqNo
	err = nil
	if h.ReturnError {
		seqno = -seqno
		err = fmt.Errorf("%d", seqno)
	}
	h.Scenario.Sequence = append(h.Scenario.Sequence, seqno)
	return
}

func setupForTest(logger *Logger, s *TestLogScenario) *Logger {
	logger.Out = ioutil.Discard
	logger.Hooks.Add(s.NewHook(0))
	return logger
}

func TestRootLoggerStd(t *testing.T) {
	var scenario TestLogScenario

	log := setupForTest(NewLogger(), &scenario)

	// Events logged for tracked log level
	log.Info("TestRootLoggerStd 1")
	require.NotNil(t, scenario.Entry)
	require.Equal(t, scenario.Entry.Level, logrus.InfoLevel)
	require.Equal(t, []int{0}, scenario.Sequence)
	scenario.Reset()
	/*
		// Add()-ing hooks at the middle
		hooks.Add(scenario.NewHook(4))
		require.NoError(t, hooks.Fire(logrus.InfoLevel, entry))
		require.Equal(t, []int{1, 2, 4, 3}, scenario.Sequence)
		scenario.Reset()
	*/
}

func TestRootLoggerPriority(t *testing.T) {
	var scenario TestLogScenario

	log := setupForTest(NewLogger(), &scenario)

	// Events logged for tracked log level
	log.Notice("TestRootLoggerStd 1")
	require.NotNil(t, scenario.Entry)
	require.Equal(t, scenario.Entry.Level, logrus.InfoLevel)
	require.Equal(t, []int{0}, scenario.Sequence)
	scenario.Reset()
	/*
		// Add()-ing hooks at the middle
		hooks.Add(scenario.NewHook(4))
		require.NoError(t, hooks.Fire(logrus.InfoLevel, entry))
		require.Equal(t, []int{1, 2, 4, 3}, scenario.Sequence)
		scenario.Reset()
	*/
}
