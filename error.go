package pick

import (
	"fmt"
	"strings"
)

type multiError struct {
	errors []error
}

func (e *multiError) Error() string {
	if e == nil || len(e.errors) == 0 {
		return ""
	}

	s := strings.Builder{}
	s.WriteRune('{')
	for i, err := range e.errors {
		if i != 0 {
			s.WriteRune('|')
		}
		s.WriteString(err.Error())
	}
	s.WriteRune('}')

	return s.String()
}

func (e *multiError) Add(err error) {
	e.errors = append(e.errors, err)
}

func (e *multiError) Unwrap() []error {
	return e.errors
}

func gather(dst *error, newErr error) {
	if newErr == nil {
		return
	}
	if dst == nil {
		return
	}

	var gatherer *multiError
	if *dst == nil {
		gatherer = &multiError{}
		*dst = gatherer
	} else if g, is := (*dst).(*multiError); is { //nolint:errorlint // this is specifically this way.
		gatherer = g
	} else {
		gatherer = &multiError{}
		gatherer.Add(*dst)
		*dst = gatherer
	}

	gatherer.Add(newErr)
}

func gatherErrorsFn(dst *error) func(string, error) {
	return func(selector string, err error) {
		gather(dst, &PickerError{selector: selector, inner: err})
	}
}

type PickerError struct {
	inner    error
	selector string
}

func (e *PickerError) Selector() string {
	return e.selector
}

func (e *PickerError) Error() string {
	return fmt.Sprintf("picker error with selector `%s` error: `%s`", e.selector, e.inner.Error())
}

func (e *PickerError) Unwrap() error {
	return e.inner
}