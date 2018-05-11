package testutil

import (
	"fmt"
	"reflect"
	"strings"

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

// containsElement try loop over the list check if the list includes the element.
// return (false, false) if impossible.
// return (true, false) if element was not found.
// return (true, true) if element was found.
func containsElement(list interface{}, element interface{}, excludeIdxs map[int]struct{}) (ok bool, found int) {

	listValue := reflect.ValueOf(list)
	elementValue := reflect.ValueOf(element)
	defer func() {
		if e := recover(); e != nil {
			ok = false
			found = -1
		}
	}()

	if reflect.TypeOf(list).Kind() == reflect.String {
		found = 0
		isExcluded := true
		str, substr := listValue.String(), elementValue.String()
		// Ensure -1 is not excluded, just in case.
		delete(excludeIdxs, -1)
		for isExcluded {
			match := strings.Index(str, substr)
			if match != -1 {
				str = str[match+1:]
				found += match
			} else {
				found = -1
			}
			_, isExcluded = excludeIdxs[found]
		}
		ok = true
		return
	}

	if reflect.TypeOf(list).Kind() == reflect.Map {
		mapKeys := listValue.MapKeys()
		for i := 0; i < len(mapKeys); i++ {
			if _, isExcluded := excludeIdxs[i]; isExcluded {
				continue
			}
			if assert.ObjectsAreEqual(mapKeys[i].Interface(), element) {
				return true, i
			}
		}
		return true, -1
	}

	for i := 0; i < listValue.Len(); i++ {
		if _, isExcluded := excludeIdxs[i]; isExcluded {
			continue
		}
		if assert.ObjectsAreEqual(listValue.Index(i).Interface(), element) {
			return true, i
		}
	}
	return true, -1

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
// FIXME: O(n^2) but can be O(n)
func IsPermutation(t assert.TestingT, expected, actual interface{}, msgAndArgs ...interface{}) (ok bool) {
	if expected == nil {
		return actual == nil
	}

	actualValue := reflect.ValueOf(actual)
	defer func() {
		if e := recover(); e != nil {
			ok = false
		}
	}()

	expectedKind := reflect.TypeOf(expected).Kind()
	actualKind := reflect.TypeOf(actual).Kind()

	if expectedKind != reflect.Array && expectedKind != reflect.Slice {
		return assert.Fail(t, fmt.Sprintf("%q has an unsupported type %s", expected, expectedKind), msgAndArgs...)
	}

	if actualKind != reflect.Array && actualKind != reflect.Slice {
		return assert.Fail(t, fmt.Sprintf("%q has an unsupported type %s", actual, actualKind), msgAndArgs...)
	}

	expectedValue := reflect.ValueOf(expected)
	if expectedValue.Len() != actualValue.Len() {
		return assert.Fail(t,
			fmt.Sprintf(
				"Expected collection of %d items but actual length is %d",
				expectedValue.Len(), actualValue.Len()),
			msgAndArgs...)
	}

	visitedIdxs := make(map[int]struct{})
	for i := 0; i < actualValue.Len(); i++ {
		element := actualValue.Index(i).Interface()
		ok, found := containsElement(expected, element, visitedIdxs)
		if !ok {
			return assert.Fail(t, fmt.Sprintf("\"%s\" could not be applied builtin len()", expected), msgAndArgs...)
		}
		if found == -1 {
			return assert.Fail(t, fmt.Sprintf("\"%s\" does not contain \"%s\"", expected, element), msgAndArgs...)
		} else {
			visitedIdxs[found] = struct{}{}
		}
	}

	return true
}
