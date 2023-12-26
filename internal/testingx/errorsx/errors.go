package errorsx

import (
	"errors"
	"fmt"
)

func RecoverPanicToError(out *error) {
	if r := recover(); r != nil {
		panicError := &recoveredPanicError{recovered: r}
		if *out != nil {
			*out = errors.Join(panicError, *out)
		} else {
			*out = panicError
		}
	}
}

type recoveredPanicError struct {
	recovered any
}

func (r *recoveredPanicError) Error() string {
	if err, is := r.recovered.(error); is {
		return fmt.Sprintf("recovered panic: %s", err.Error())
	}

	return fmt.Sprintf("recovered panic: %#v", r.recovered)
}

func (r *recoveredPanicError) Recovered() any {
	return r.recovered
}

func (r *recoveredPanicError) Unwrap() error {
	err, _ := r.recovered.(error)
	return err
}
