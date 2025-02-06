package pick

import (
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/moukoublen/pick/internal/tst"
)

func TestMultiError(t *testing.T) {
	errOne := errors.New("one")
	errTwo := errors.New("two")
	errThree := errors.New("three")

	m := &multiError{}

	tst.AssertEqual(t, m.Error(), "")

	m.Add(errOne)
	m.Add(errTwo)
	m.Add(errThree)

	tst.AssertEqual(t, m.Error(), "one | two | three")
	tst.AssertEqual(t, errors.Is(m, errOne), true)
	tst.AssertEqual(t, errors.Is(m, errTwo), true)
	tst.AssertEqual(t, errors.Is(m, errThree), true)
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
				tst.AssertEqual(t, dst, (*error)(nil))
			},
		},
		{
			destination: nil,
			newError:    errors.New("error"),
			expect: func(t *testing.T, dst *error) {
				tst.AssertEqual(t, dst, (*error)(nil))
			},
		},
		{
			destination: new(error),
			newError:    errors.New("error"),
			expect: func(t *testing.T, dst *error) {
				m, is := (*dst).(*multiError)
				tst.AssertEqual(t, is, true)
				tst.AssertEqual(t, m.Error(), "error")
			},
		},
		{
			destination: new(error),
			newError:    errors.New("error"),
			expect: func(t *testing.T, dst *error) {
				m, is := (*dst).(*multiError)
				tst.AssertEqual(t, is, true)
				tst.AssertEqual(t, m.Error(), "error")
			},
		},
		{
			destination: ptr[error](&multiError{}),
			newError:    errors.New("error"),
			expect: func(t *testing.T, dst *error) {
				m, is := (*dst).(*multiError)
				tst.AssertEqual(t, is, true)
				tst.AssertEqual(t, m.Error(), "error")
			},
		},
		{
			destination: ptr[error](errors.New("one")),
			newError:    errors.New("two"),
			expect: func(t *testing.T, dst *error) {
				m, is := (*dst).(*multiError)
				tst.AssertEqual(t, is, true)
				tst.AssertEqual(t, len(m.errors), 2)
				tst.AssertEqual(t, m.Error(), "one | two")
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
		errorAsserter tst.ErrorAsserter
	}{
		"empty": {
			op:            func(_ *ErrorsSink) {},
			errorAsserter: tst.NoError,
		},

		"one Gather": {
			op: func(es *ErrorsSink) {
				es.Gather(errors.New("error"))
			},
			errorAsserter: tst.ExpectedErrorChecks(
				tst.ExpectedErrorStringContains("error"),
			),
		},

		"one GatherSelector": {
			op: func(es *ErrorsSink) {
				es.GatherSelector("one.two", io.ErrUnexpectedEOF)
			},
			errorAsserter: tst.ExpectedErrorChecks(
				tst.ExpectedErrorStringContains("unexpected EOF"),
				tst.ExpectedErrorOfType[*PickerError](
					func(t *testing.T, pe *PickerError) {
						tst.AssertEqual(t, "one.two", pe.Selector())
					},
				),
				tst.ExpectedErrorIs(io.ErrUnexpectedEOF),
			),
		},

		"mixed": {
			op: func(es *ErrorsSink) {
				es.GatherSelector("one.two", io.ErrUnexpectedEOF)
				es.GatherSelector("one.three", io.ErrClosedPipe)
			},
			errorAsserter: tst.ExpectedErrorChecks(
				tst.ExpectedErrorStringContains("picker error with selector `one.two` error: `unexpected EOF` |"),
				tst.ExpectedErrorStringContains("| picker error with selector `one.three` error: `io: read/write on closed pipe`"),
				tst.ExpectedErrorOfType[*multiError](
					func(t *testing.T, pe *multiError) {
						tst.ExpectedErrorOfType[*PickerError](
							func(t *testing.T, pe *PickerError) {
								tst.AssertEqual(t, "one.two", pe.Selector())
							},
						)(t, pe.errors[0])

						tst.ExpectedErrorOfType[*PickerError](
							func(t *testing.T, pe *PickerError) {
								tst.AssertEqual(t, "one.three", pe.Selector())
							},
						)(t, pe.errors[1])
					},
				),
				tst.ExpectedErrorIs(io.ErrUnexpectedEOF),
				tst.ExpectedErrorIs(io.ErrClosedPipe),
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			es := &ErrorsSink{}
			tc.op(es)
			tc.errorAsserter(t, es.Outcome())
		})
	}
}
