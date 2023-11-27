package errors

import (
	"errors"
	"fmt"
	"strings"
)

func RecoverPanicToError(out *error) {
	if r := recover(); r != nil {
		panicError := &RecoverPanicError{recovered: r}
		if *out != nil {
			*out = JoinErrors(panicError, *out)
		} else {
			*out = panicError
		}
	}
}

func JoinErrors(errs ...error) error {
	e := &joinError{errs: make([]error, 0, len(errs))}
	for _, err := range errs {
		if err != nil {
			e.errs = append(e.errs, err)
		}
	}
	return e
}

type joinError struct {
	errs []error
}

func (e *joinError) Error() string {
	sb := strings.Builder{}
	for i, err := range e.errs {
		if i > 0 {
			sb.WriteString(" :: ")
		}
		sb.WriteString(err.Error())
	}
	return sb.String()
}

func (e *joinError) Unwrap() []error {
	return e.errs
}

// As checks if any of joined errors (.errs) can be casted as target.
// Implement this to be backwards compatible with go 1.19.
func (e *joinError) As(target any) bool {
	if err, is := target.(**joinError); is {
		*err = e
		return true
	}

	for _, err := range e.errs {
		if errors.As(err, target) {
			return true
		}
	}

	return false
}

// Is checks if any of joined errors (.errs) "is" the target error.
// Implement this to be backwards compatible with go 1.19.
func (e *joinError) Is(target error) bool {
	//nolint:errorlint
	if _, is := target.(*joinError); is {
		return true
	}

	for _, err := range e.errs {
		if errors.Is(err, target) {
			return true
		}
	}

	return false
}

type RecoverPanicError struct {
	recovered any
}

func (r *RecoverPanicError) Error() string {
	if err, is := r.recovered.(error); is {
		return err.Error()
	}

	return fmt.Sprintf("recovered panic: %#v", r.recovered)
}

func (r *RecoverPanicError) Recovered() any {
	return r.recovered
}

func (r *RecoverPanicError) Unwrap() error {
	err, _ := r.recovered.(error)
	return err
}
