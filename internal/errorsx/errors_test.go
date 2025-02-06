package errorsx

import (
	"errors"
	"testing"

	"github.com/moukoublen/pick/internal/tst"
)

func expectJoinedErrorWithLen(expectedLen int) func(t *testing.T, err error) {
	return func(t *testing.T, err error) {
		t.Helper()
		joined, is := err.(interface{ Unwrap() []error })
		if !is {
			t.Errorf("error is expected to implements Unwrap() []error")
			return
		}

		ers := joined.Unwrap()
		if len(ers) != expectedLen {
			t.Errorf("expected join of %d errors got %d errors", expectedLen, len(ers))
		}
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
		errorAsserter tst.ErrorAsserter
	}{
		"no error": {
			out:           ptrNilErr(),
			fnPanics:      func() {},
			errorAsserter: tst.NoError,
		},
		"panic to nil out": {
			out:      ptrNilErr(),
			fnPanics: func() { panic("panic!") },
			errorAsserter: tst.ExpectedErrorChecks(
				tst.ExpectedErrorStringContains(`recovered panic: "panic!"`),
				tst.ExpectedErrorOfType[*recoveredPanicError](
					func(t *testing.T, rpe *recoveredPanicError) { //nolint:thelper
						tst.AssertEqual(t, rpe.Recovered(), "panic!")
					},
				),
			),
		},
		"panic to not nil out": {
			out:      ptrErr(errSample),
			fnPanics: func() { panic("panic!") },
			errorAsserter: tst.ExpectedErrorChecks(
				expectJoinedErrorWithLen(2),
				tst.ExpectedErrorStringContains("recovered panic: \"panic!\"\nsample error"),
				tst.ExpectedErrorOfType[*recoveredPanicError](),
			),
		},
		"panic error to nil out": {
			out:      ptrNilErr(),
			fnPanics: func() { panic(errSample) },
			errorAsserter: tst.ExpectedErrorChecks(
				tst.ExpectedErrorStringContains(`recovered panic: sample error`),
				tst.ExpectedErrorOfType[*recoveredPanicError](
					func(t *testing.T, rpe *recoveredPanicError) { //nolint:thelper
						tst.ExpectedErrorIs(errSample)(t, rpe.Unwrap())
					},
				),
			),
		},
		"panic error to not nil out": {
			out:      ptrErr(errSample),
			fnPanics: func() { panic(errSample) },
			errorAsserter: tst.ExpectedErrorChecks(
				expectJoinedErrorWithLen(2),
				tst.ExpectedErrorStringContains("recovered panic: sample error\nsample error"),
				tst.ExpectedErrorOfType[*recoveredPanicError](),
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
