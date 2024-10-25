package errorsx

import (
	"errors"
	"testing"

	"github.com/moukoublen/pick/internal/testingx"
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
		out         *error
		fnPanics    func()
		expectedErr func(*testing.T, error)
	}{
		"no error": {
			out:         ptrNilErr(),
			fnPanics:    func() {},
			expectedErr: nil,
		},
		"panic to nil out": {
			out:      ptrNilErr(),
			fnPanics: func() { panic("panic!") },
			expectedErr: testingx.ExpectedErrorChecks(
				testingx.ExpectedErrorStringContains(`recovered panic: "panic!"`),
				testingx.ExpectedErrorOfType[*recoveredPanicError](
					func(t *testing.T, rpe *recoveredPanicError) { //nolint:thelper
						testingx.AssertEqual(t, rpe.Recovered(), "panic!")
					},
				),
			),
		},
		"panic to not nil out": {
			out:      ptrErr(errSample),
			fnPanics: func() { panic("panic!") },
			expectedErr: testingx.ExpectedErrorChecks(
				expectJoinedErrorWithLen(2),
				testingx.ExpectedErrorStringContains("recovered panic: \"panic!\"\nsample error"),
				testingx.ExpectedErrorOfType[*recoveredPanicError](),
			),
		},
		"panic error to nil out": {
			out:      ptrNilErr(),
			fnPanics: func() { panic(errSample) },
			expectedErr: testingx.ExpectedErrorChecks(
				testingx.ExpectedErrorStringContains(`recovered panic: sample error`),
				testingx.ExpectedErrorOfType[*recoveredPanicError](
					func(t *testing.T, rpe *recoveredPanicError) { //nolint:thelper
						testingx.ExpectedErrorIs(errSample)(t, rpe.Unwrap())
					},
				),
			),
		},
		"panic error to not nil out": {
			out:      ptrErr(errSample),
			fnPanics: func() { panic(errSample) },
			expectedErr: testingx.ExpectedErrorChecks(
				expectJoinedErrorWithLen(2),
				testingx.ExpectedErrorStringContains("recovered panic: sample error\nsample error"),
				testingx.ExpectedErrorOfType[*recoveredPanicError](),
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			defer func() {
				testingx.AssertError(t, tc.expectedErr, *tc.out)
			}()

			defer RecoverPanicToError(tc.out)
			if tc.fnPanics != nil {
				tc.fnPanics()
			}
		})
	}
}
