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
}

func casterTest[T any](t *testing.T, testCases []casterTestCase[T], castFn func(any) (T, error)) {
	t.Helper()
	casterTestWithCompare[T](t, testCases, castFn, func(x, y T) bool { return reflect.DeepEqual(x, y) })
}

func casterTestWithCompare[T any](t *testing.T, testCases []casterTestCase[T], castFn func(any) (T, error), equalFn func(T, T) bool) {
	t.Helper()

	for idx, tc := range testCases {
		tc := tc

		typeName := "nil"
		if tc.input != nil {
			tp := reflect.TypeOf(tc.input)
			typeName = tp.Name()
		}

		name := fmt.Sprintf("index:%d input_type:%s input_value:(%#v)", idx, typeName, tc.input)
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got, gotErr := castFn(tc.input)
			testingx.AssertError(t, tc.expectedErr, gotErr)
			if !equalFn(tc.expected, got) {
				t.Errorf("wrong returned value. Expected %#v got %#v", tc.expected, got)
			}
		})
	}
}

func TestCastToSliceErrorScenarios(t *testing.T) {
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
			_, gotErr := castToSlice(tc.input, tc.inputSingleItemCastFn)
			testingx.AssertError(t, tc.expectedErr, gotErr)
		})
	}
}

func TestCastAttemptUsingReflect(t *testing.T) {
	t.Parallel()

	t.Run("string", func(t *testing.T) {
		t.Parallel()
		type stringAlias string
		type stringSecondAlias stringAlias

		testCases := []casterTestCase[string]{
			{
				input:       stringAlias("test"),
				expected:    "test",
				expectedErr: nil,
			},
			{
				input:       stringSecondAlias(stringAlias("test")),
				expected:    "test",
				expectedErr: nil,
			},
		}
		casterTest[string](t, testCases, castAttemptUsingReflect[string])
	})

	t.Run("map[string]string", func(t *testing.T) {
		t.Parallel()
		type mapAlias map[string]string

		testCases := []casterTestCase[map[string]string]{
			{
				input:       mapAlias{"abc": "cba"},
				expected:    map[string]string{"abc": "cba"},
				expectedErr: nil,
			},
		}
		casterTest[map[string]string](t, testCases, castAttemptUsingReflect[map[string]string])
	})
}
