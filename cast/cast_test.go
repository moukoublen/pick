package cast

import (
	"errors"
	"fmt"
	"reflect"
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

func TestTryCastUsingReflect(t *testing.T) {
	type intAlias int
	tests := map[string]struct {
		fn          any
		input       any
		expected    any
		expectedErr func(*testing.T, error)
	}{
		"intAlias to int16": {
			fn:          tryCastUsingReflect[int16],
			input:       intAlias(13),
			expected:    int16(13),
			expectedErr: nil,
		},
		"struct to int16 expect error": {
			fn:          tryCastUsingReflect[int16],
			input:       struct{}{},
			expected:    int16(0),
			expectedErr: testingx.ExpectedErrorIs(ErrInvalidType),
		},
		"string to []byte": {
			fn:          tryCastUsingReflect[[]byte],
			input:       "str",
			expected:    []byte("str"),
			expectedErr: nil,
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			fnVal := reflect.ValueOf(tc.fn)
			returnedVals := fnVal.Call([]reflect.Value{
				reflect.ValueOf(tc.input),
			})

			if len(returnedVals) != 2 {
				t.Fatalf("number of returned values %d", len(returnedVals))
			}

			reflect.DeepEqual(returnedVals[0], reflect.ValueOf(tc.expected))

			errVal := returnedVals[1]
			errInf := errVal.Interface()
			if errInf == nil {
				testingx.AssertError(t, tc.expectedErr, nil)
			} else {
				err, is := errInf.(error)
				if !is {
					t.Errorf("second returned item is not of type error. Type: %s", errVal.Type().String())
				}
				testingx.AssertError(t, tc.expectedErr, err)
			}
		})
	}
}
