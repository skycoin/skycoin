package kvstorage

// Error wraps key-value storage related errors.
// It wraps errors caused by user input, but not errors caused by
// programmer input or internal issues.
type Error struct {
	error
}

// NewError creates an Error
func NewError(err error) error {
	if err == nil {
		return nil
	}
	return Error{err}
}
