package pick

import (
	"errors"
	"fmt"
	"testing"

	"github.com/moukoublen/pick/internal/testingx"
)

func TestMultiError(t *testing.T) {
	errOne := errors.New("one")
	errTwo := errors.New("two")
	errThree := errors.New("three")

	m := &multiError{}

	testingx.AssertEqual(t, m.Error(), "")

	m.Add(errOne)
	m.Add(errTwo)
	m.Add(errThree)

	testingx.AssertEqual(t, m.Error(), "one | two | three")
	testingx.AssertEqual(t, errors.Is(m, errOne), true)
	testingx.AssertEqual(t, errors.Is(m, errTwo), true)
	testingx.AssertEqual(t, errors.Is(m, errThree), true)
}

func ptr(e error) *error { return &e }

func TestGather(t *testing.T) {
	tests := []struct {
		destination *error
		newError    error
		expect      func(t *testing.T, dst *error)
	}{
		{
			destination: nil,
			newError:    nil,
			expect: func(t *testing.T, dst *error) { //nolint:thelper
				testingx.AssertEqual(t, dst, (*error)(nil))
			},
		},
		{
			destination: nil,
			newError:    errors.New("error"),
			expect: func(t *testing.T, dst *error) { //nolint:thelper
				testingx.AssertEqual(t, dst, (*error)(nil))
			},
		},
		{
			destination: new(error),
			newError:    errors.New("error"),
			expect: func(t *testing.T, dst *error) { //nolint:thelper
				m, is := (*dst).(*multiError) //nolint:errorlint
				testingx.AssertEqual(t, is, true)
				testingx.AssertEqual(t, m.Error(), "error")
			},
		},
		{
			destination: new(error),
			newError:    errors.New("error"),
			expect: func(t *testing.T, dst *error) { //nolint:thelper
				m, is := (*dst).(*multiError) //nolint:errorlint
				testingx.AssertEqual(t, is, true)
				testingx.AssertEqual(t, m.Error(), "error")
			},
		},
		{
			destination: ptr(&multiError{}),
			newError:    errors.New("error"),
			expect: func(t *testing.T, dst *error) { //nolint:thelper
				m, is := (*dst).(*multiError) //nolint:errorlint
				testingx.AssertEqual(t, is, true)
				testingx.AssertEqual(t, m.Error(), "error")
			},
		},
		{
			destination: ptr(errors.New("one")),
			newError:    errors.New("two"),
			expect: func(t *testing.T, dst *error) { //nolint:thelper
				m, is := (*dst).(*multiError) //nolint:errorlint
				testingx.AssertEqual(t, is, true)
				testingx.AssertEqual(t, len(m.errors), 2)
				testingx.AssertEqual(t, m.Error(), "one | two")
			},
		},
	}

	for i, tc := range tests {
		tc := tc
		t.Run(fmt.Sprintf("[%d]", i), func(t *testing.T) {
			gather(tc.destination, tc.newError)
			tc.expect(t, tc.destination)
		})
	}
}
