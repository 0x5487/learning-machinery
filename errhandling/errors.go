package errhandling

import "errors"

// ThrowErr returns error
func ThrowErr(name string) error {
	return errors.New("oops")
}

// ThrowPanic throw panic
func ThrowPanic(name string) error {
	err := errors.New("dead")
	panic(err)
}
