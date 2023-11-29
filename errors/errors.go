package errors

import (
	"errors"
	"fmt"
)

func RecoverPanicToError(out *error) {
	if r := recover(); r != nil {
		panicError := &RecoverPanicError{recovered: r}
		if *out != nil {
			*out = errors.Join(panicError, *out)
		} else {
			*out = panicError
		}
	}
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
