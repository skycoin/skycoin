package require

import (
	assert "github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"

	_assert "github.com/skycoin/skycoin/src/testutil/assert"
)

// PanicsWithCondition asserts that the code inside the specified PanicTestFunc panics, and that
// the recovered panic value meets a given condition.
//
//   assert.PanicsWithCondition(t, func(value){ return assert.True(t, isCrazy(value)) }, func(){ GoCrazy() })
//
// Returns whether the assertion was successful (true) or not (false).
func PanicsWithCondition(t require.TestingT, condition _assert.TestValuePredicate, f assert.PanicTestFunc, msgAndArgs ...interface{}) {
	if !_assert.PanicsWithCondition(t, condition, f, msgAndArgs...) {
		t.FailNow()
	}
}

// PanicsWithLogMessage asserts that the code inside the specified PanicTestFunc panics, and that
// an expected string is included in log message.
//
//   assert.PanicsWithLogMessage(t, "Log msg", func(){ log.Panic("Log msg for X") })
//
// Returns whether the assertion was successful (true) or not (false).
func PanicsWithLogMessage(t require.TestingT, expectedMessage string, f assert.PanicTestFunc, msgAndArgs ...interface{}) {
	if !_assert.PanicsWithLogMessage(t, expectedMessage, f, msgAndArgs) {
		t.FailNow()
	}
}
