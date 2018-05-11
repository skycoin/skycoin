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
func PanicsWithLogMessage(t require.TestingT, expectedMessage string, f assert.PanicTestFunc, msgAndArgs ...interface{}) {
	if !_assert.PanicsWithLogMessage(t, expectedMessage, f, msgAndArgs) {
		t.FailNow()
	}
}

// IsPermutation asserts that two containers (slices, arrays, ...) have exactly the same items.
//
//    assert.Same(t, "123", "132")
//    assert.Same(t, []int{1, 2, 3}, []int{1, 3, 2})
//
// Returns whether the assertion was successful (true) or not (false).
//
// Pointer variable equality is determined based on the equality of the
// referenced values (as opposed to the memory addresses). Function equality
// cannot be determined and will always fail.
//
// Equivalent to Subset(expected, actual) && Subset(actual,expected) but in
// a single pass.
//
func IsPermutation(t require.TestingT, expected, actual interface{}, msgAndArgs ...interface{}) {
	if !_assert.IsPermutation(t, expected, actual, msgAndArgs...) {
		t.FailNow()
	}
}
