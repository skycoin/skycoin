package testutil

import (
	"fmt"

	logrus "github.com/sirupsen/logrus"
	assert "github.com/stretchr/testify/assert"
)

// didPanic returns true if the function passed to it panics. Otherwise, it returns false.
func didPanic(f assert.PanicTestFunc) (bool, interface{}) {

	didPanic := false
	var message interface{}
	func() {

		defer func() {
			if message = recover(); message != nil {
				didPanic = true
			}
		}()

		// call the target function
		f()

	}()

	return didPanic, message
}

// TestValuePredicate checks that a value meets certain condition
// This function type may be seen as a continuation of a test scenario.
// Assertions executed inside of it may be recorded to the calling context
// by accessing free-vars in function closure.
type TestValuePredicate func(value interface{}) (success bool)

// PanicsWithCondition asserts that the code inside the specified PanicTestFunc panics, and that
// the recovered panic value meets a given condition.
//
//   assert.PanicsWithCondition(t, func(value){ return assert.True(t, isCrazy(value)) }, func(){ GoCrazy() })
//
// Returns whether the assertion was successful (true) or not (false).
func PanicsWithCondition(t assert.TestingT, condition TestValuePredicate, f assert.PanicTestFunc, msgAndArgs ...interface{}) bool {

	funcDidPanic, panicValue := didPanic(f)
	if !funcDidPanic {
		return assert.Fail(t, fmt.Sprintf("func %#v should panic\n\r\tPanic value:\t%v", f, panicValue), msgAndArgs...)
	}
	return condition(panicValue)
}

// PanicsWithLogMessage asserts that the code inside the specified PanicTestFunc panics, and that
// an expected string is included in log message.
//
//   assert.PanicsWithLogMessage(t, "Log msg", func(){ log.Panic("Log msg for X") })
//
// Returns whether the assertion was successful (true) or not (false).
func PanicsWithLogMessage(t assert.TestingT, expectedMessage string, f assert.PanicTestFunc, msgAndArgs ...interface{}) bool {
	return PanicsWithCondition(t, func(logValue interface{}) bool {
		gotMessage, gotIt := "", false
		if entry, isEntry := logValue.(*logrus.Entry); isEntry {
			gotMessage, gotIt = entry.Message, true
		} else {
			if msg, isString := logValue.(string); isString {
				gotMessage, gotIt = msg, true
			}
		}
		if gotIt {
			return assert.Contains(t, gotMessage, expectedMessage)
		}
		return assert.Fail(t, "expected string or log entry but got %T", logValue)
	}, f, msgAndArgs)
}
