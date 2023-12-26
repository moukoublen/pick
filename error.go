package pick

import (
	"fmt"
	"strings"
)

type multipleError struct {
	errors []error
}

func (e *multipleError) Error() string {
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

func (e *multipleError) Add(err error) {
	e.errors = append(e.errors, err)
}

func (e *multipleError) Unwrap() []error {
	return e.errors
}

func gather(dst *error, newErr error) {
	if newErr == nil {
		return
	}
	if dst == nil {
		return
	}

	var gatherer *multipleError
	if *dst == nil {
		gatherer = &multipleError{}
		*dst = gatherer
	} else if g, is := (*dst).(*multipleError); is { //nolint:errorlint // this is specifically this way.
		gatherer = g
	} else {
		gatherer = &multipleError{}
		gatherer.Add(*dst)
		*dst = gatherer
	}

	gatherer.Add(newErr)
}

func GatherErrorsFn(dst *error) func(string, error) {
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
