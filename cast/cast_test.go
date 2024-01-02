package cast

import (
	"errors"
	"fmt"
	"testing"

	"github.com/moukoublen/pick/internal/testingx"
)

var (
	expectOverFlowError   = testingx.ExpectedErrorIs(&Error{}, ErrCastOverFlow)
	expectLostDecimals    = testingx.ExpectedErrorIs(&Error{}, ErrCastLostDecimals)
	expectMalformedSyntax = testingx.ExpectedErrorIs(&Error{}, ErrInvalidSyntax)
	expectInvalidType     = testingx.ExpectedErrorIs(&Error{}, ErrInvalidType)
)

type casterTestCase[T any] struct {
	input       any
	expectedErr func(*testing.T, error)
	expected    T
	castFn      func(any) (T, error)
}

func casterTest[T any](t *testing.T, testCases []casterTestCase[T], defaultCastFn func(any) (T, error)) {
	t.Helper()

	for idx, tc := range testCases {
		tc := tc
		name := fmt.Sprintf("index:%d %s", idx, testingx.Format(tc.input))
		t.Run(name, func(t *testing.T) {
			t.Helper()
			t.Parallel()
			var (
				got    T
				gotErr error
			)
			if tc.castFn != nil {
				got, gotErr = tc.castFn(tc.input)
			} else {
				got, gotErr = defaultCastFn(tc.input)
			}
			testingx.AssertError(t, tc.expectedErr, gotErr)

			testingx.AssertEqual(t, got, tc.expected)
		})
	}
}

func TestToSliceErrorScenarios(t *testing.T) {
	t.Parallel()

	errMock1 := errors.New("mock error")

	type testCase struct {
		input                 any
		inputSingleItemCastFn func(any) (int, error)
		expectedErr           func(*testing.T, error)
	}

	testsCases := []testCase{
		{
			input:                 []any{1, 2, 3},
			inputSingleItemCastFn: func(any) (int, error) { return 0, errMock1 },
			expectedErr:           testingx.ExpectedErrorIs(errMock1),
		},
		{
			input:                 []any{1, 2, 3},
			inputSingleItemCastFn: func(any) (int, error) { panic("panic") },
			expectedErr:           testingx.ExpectedErrorStringContains(`recovered panic: "panic"`),
		},
	}

	for idx, tc := range testsCases {
		tc := tc
		name := fmt.Sprintf("test_%d_(%v)", idx, tc.input)
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			_, gotErr := ToSlice(tc.input, sliceOp(tc.inputSingleItemCastFn))
			testingx.AssertError(t, tc.expectedErr, gotErr)
		})
	}
}
