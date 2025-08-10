package pick

import (
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/ifnotnil/x/tst"
	"github.com/moukoublen/pick/internal/testingx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMultiError(t *testing.T) {
	errOne := errors.New("one")
	errTwo := errors.New("two")
	errThree := errors.New("three")

	m := &multiError{}

	require.Empty(t, m.Error())

	m.Add(errOne)
	m.Add(errTwo)
	m.Add(errThree)

	assert.Equal(t, "one | two | three", m.Error())
	require.ErrorIs(t, m, errOne)
	require.ErrorIs(t, m, errTwo)
	require.ErrorIs(t, m, errThree)
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
		errorAsserter tst.ErrorAssertionFunc
	}{
		"empty": {
			op:            func(_ *ErrorsSink) {},
			errorAsserter: tst.NoError(),
		},

		"one Gather": {
			op: func(es *ErrorsSink) {
				es.Gather(errors.New("error"))
			},
			errorAsserter: tst.All(
				tst.ErrorStringContains("error"),
			),
		},

		"one GatherSelector": {
			op: func(es *ErrorsSink) {
				es.GatherSelector("one.two", io.ErrUnexpectedEOF)
			},
			errorAsserter: tst.All(
				tst.ErrorStringContains("unexpected EOF"),
				tst.ErrorOfType[*PickerError](
					func(t tst.TestingT, pe *PickerError) {
						assert.Equal(t, "one.two", pe.Selector())
					},
				),
				tst.ErrorIs(io.ErrUnexpectedEOF),
			),
		},

		"mixed": {
			op: func(es *ErrorsSink) {
				es.GatherSelector("one.two", io.ErrUnexpectedEOF)
				es.GatherSelector("one.three", io.ErrClosedPipe)
			},
			errorAsserter: tst.All(
				tst.ErrorStringContains("picker error with selector `one.two` error: `unexpected EOF` |"),
				tst.ErrorStringContains("| picker error with selector `one.three` error: `io: read/write on closed pipe`"),
				tst.ErrorOfType[*multiError](
					func(t tst.TestingT, pe *multiError) {
						tst.ErrorOfType[*PickerError](
							func(t tst.TestingT, pe *PickerError) {
								assert.Equal(t, "one.two", pe.Selector())
							},
						)(t, pe.errors[0])

						tst.ErrorOfType[*PickerError](
							func(t tst.TestingT, pe *PickerError) {
								assert.Equal(t, "one.three", pe.Selector())
							},
						)(t, pe.errors[1])
					},
				),
				tst.ErrorIs(io.ErrUnexpectedEOF),
				tst.ErrorIs(io.ErrClosedPipe),
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
