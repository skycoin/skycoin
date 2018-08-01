package util

import (
	"errors"
	"fmt"
)

// ValueError an operation or function receives an argument
// that has the right type but an inappropriate value
type ValueError struct {
	ErrorData  error
	ParamName  string
	ParamValue interface{}
}

func (err ValueError) Error() string {
	return fmt.Sprintf("Invalid value for '%s' : %s",
		err.ParamName, err.ErrorData.Error())
}

// NewValueError instantiate value error from concrete error
func NewValueError(err error, paramName string, paramValue interface{}) ValueError {
	return ValueError{
		err, paramName, paramValue,
	}
}

// NewValueErrorFromString shortcut for new value errors given an error message
func NewValueErrorFromString(errorMessage string, paramName string, paramValue interface{}) ValueError {
	return NewValueError(errors.New(errorMessage), paramName, paramValue)
}

// SameError check for identical error conditions
func SameError(err1, err2 error) bool {
	if _err, isOfType := err1.(ValueError); isOfType {
		err1 = _err
	}
	if _err, isOfType := err2.(ValueError); isOfType {
		err2 = _err
	}
	return err1 == err2
}
