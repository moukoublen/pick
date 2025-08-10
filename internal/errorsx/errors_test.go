package errorsx

import (
	"errors"
	"testing"

	"github.com/ifnotnil/x/tst"
	"github.com/stretchr/testify/assert"
)

func expectJoinedErrorWithLen(expectedLen int) func(t tst.TestingT, err error) bool {
	return func(t tst.TestingT, err error) bool {
		t.Helper()
		joined, is := err.(interface{ Unwrap() []error })
		if !is {
			t.Errorf("error is expected to implements Unwrap() []error")
			return false
		}

		ers := joined.Unwrap()
		if len(ers) != expectedLen {
			t.Errorf("expected join of %d errors got %d errors", expectedLen, len(ers))
			return false
		}

		return true
	}
}

func ptrNilErr() *error {
	var err error
	return &err
}

// ptrErr pointer of copy of the given error.
func ptrErr(err error) *error {
	return &err
}

var errSample = errors.New("sample error")

func TestRecoverPanicToError(t *testing.T) {
	tests := map[string]struct {
		out           *error
		fnPanics      func()
		errorAsserter tst.ErrorAssertionFunc
	}{
		"no error": {
			out:           ptrNilErr(),
			fnPanics:      func() {},
			errorAsserter: tst.NoError(),
		},
		"panic to nil out": {
			out:      ptrNilErr(),
			fnPanics: func() { panic("panic!") },
			errorAsserter: tst.All(
				tst.ErrorStringContains(`recovered panic: "panic!"`),
				tst.ErrorOfType[*recoveredPanicError](
					func(t tst.TestingT, rpe *recoveredPanicError) { //nolint:thelper
						assert.Equal(t, "panic!", rpe.Recovered())
					},
				),
			),
		},
		"panic to not nil out": {
			out:      ptrErr(errSample),
			fnPanics: func() { panic("panic!") },
			errorAsserter: tst.All(
				expectJoinedErrorWithLen(2),
				tst.ErrorStringContains("recovered panic: \"panic!\"\nsample error"),
				tst.ErrorOfType[*recoveredPanicError](),
			),
		},
		"panic error to nil out": {
			out:      ptrNilErr(),
			fnPanics: func() { panic(errSample) },
			errorAsserter: tst.All(
				tst.ErrorStringContains(`recovered panic: sample error`),
				tst.ErrorOfType[*recoveredPanicError](
					func(t tst.TestingT, rpe *recoveredPanicError) { //nolint:thelper
						assert.ErrorIs(t, rpe, errSample)
					},
				),
			),
		},
		"panic error to not nil out": {
			out:      ptrErr(errSample),
			fnPanics: func() { panic(errSample) },
			errorAsserter: tst.All(
				expectJoinedErrorWithLen(2),
				tst.ErrorStringContains("recovered panic: sample error\nsample error"),
				tst.ErrorOfType[*recoveredPanicError](),
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// stacked position 2, assert.
			defer func() {
				tc.errorAsserter(t, *tc.out)
			}()

			// stacked position 1, recover.
			defer RecoverPanicToError(tc.out)

			// panic (or not).
			if tc.fnPanics != nil {
				tc.fnPanics()
			}
		})
	}
}
