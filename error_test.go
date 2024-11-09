package pick

import (
	"errors"
	"fmt"
	"io"
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

func ptr[T any](o T) *T {
	return &o
}

func TestGather(t *testing.T) {
	tests := []struct {
		destination *error
		newError    error
		expect      func(t *testing.T, dst *error)
	}{
		{
			destination: nil,
			newError:    nil,
			expect: func(t *testing.T, dst *error) {
				testingx.AssertEqual(t, dst, (*error)(nil))
			},
		},
		{
			destination: nil,
			newError:    errors.New("error"),
			expect: func(t *testing.T, dst *error) {
				testingx.AssertEqual(t, dst, (*error)(nil))
			},
		},
		{
			destination: new(error),
			newError:    errors.New("error"),
			expect: func(t *testing.T, dst *error) {
				m, is := (*dst).(*multiError)
				testingx.AssertEqual(t, is, true)
				testingx.AssertEqual(t, m.Error(), "error")
			},
		},
		{
			destination: new(error),
			newError:    errors.New("error"),
			expect: func(t *testing.T, dst *error) {
				m, is := (*dst).(*multiError)
				testingx.AssertEqual(t, is, true)
				testingx.AssertEqual(t, m.Error(), "error")
			},
		},
		{
			destination: ptr[error](&multiError{}),
			newError:    errors.New("error"),
			expect: func(t *testing.T, dst *error) {
				m, is := (*dst).(*multiError)
				testingx.AssertEqual(t, is, true)
				testingx.AssertEqual(t, m.Error(), "error")
			},
		},
		{
			destination: ptr[error](errors.New("one")),
			newError:    errors.New("two"),
			expect: func(t *testing.T, dst *error) {
				m, is := (*dst).(*multiError)
				testingx.AssertEqual(t, is, true)
				testingx.AssertEqual(t, len(m.errors), 2)
				testingx.AssertEqual(t, m.Error(), "one | two")
			},
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("[%d]", i), func(t *testing.T) {
			gather(tc.destination, tc.newError)
			tc.expect(t, tc.destination)
		})
	}
}

func TestErrorsSink(t *testing.T) {
	tests := map[string]struct {
		op            func(es *ErrorsSink)
		expectedError func(*testing.T, error)
	}{
		"empty": {
			op:            func(_ *ErrorsSink) {},
			expectedError: nil,
		},

		"one Gather": {
			op: func(es *ErrorsSink) {
				es.Gather(errors.New("error"))
			},
			expectedError: testingx.ExpectedErrorChecks(
				testingx.ExpectedErrorStringContains("error"),
			),
		},

		"one GatherSelector": {
			op: func(es *ErrorsSink) {
				es.GatherSelector("one.two", io.ErrUnexpectedEOF)
			},
			expectedError: testingx.ExpectedErrorChecks(
				testingx.ExpectedErrorStringContains("unexpected EOF"),
				testingx.ExpectedErrorOfType[*PickerError](
					func(t *testing.T, pe *PickerError) {
						testingx.AssertEqual(t, "one.two", pe.Selector())
					},
				),
				testingx.ExpectedErrorIs(io.ErrUnexpectedEOF),
			),
		},

		"mixed": {
			op: func(es *ErrorsSink) {
				es.GatherSelector("one.two", io.ErrUnexpectedEOF)
				es.GatherSelector("one.three", io.ErrClosedPipe)
			},
			expectedError: testingx.ExpectedErrorChecks(
				testingx.ExpectedErrorStringContains("picker error with selector `one.two` error: `unexpected EOF` |"),
				testingx.ExpectedErrorStringContains("| picker error with selector `one.three` error: `io: read/write on closed pipe`"),
				testingx.ExpectedErrorOfType[*multiError](
					func(t *testing.T, pe *multiError) {
						testingx.ExpectedErrorOfType[*PickerError](
							func(t *testing.T, pe *PickerError) {
								testingx.AssertEqual(t, "one.two", pe.Selector())
							},
						)(t, pe.errors[0])

						testingx.ExpectedErrorOfType[*PickerError](
							func(t *testing.T, pe *PickerError) {
								testingx.AssertEqual(t, "one.three", pe.Selector())
							},
						)(t, pe.errors[1])
					},
				),
				testingx.ExpectedErrorIs(io.ErrUnexpectedEOF),
				testingx.ExpectedErrorIs(io.ErrClosedPipe),
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			es := &ErrorsSink{}
			tc.op(es)
			testingx.AssertError(t, tc.expectedError, es.Outcome())
		})
	}
}
