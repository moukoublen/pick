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

			gotVal := returnedVals[0]
			var pass bool
			if gotVal.Type().Comparable() {
				pass = gotVal.Equal(reflect.ValueOf(tc.expected))
			} else {
				a := gotVal.Interface()
				pass = reflect.DeepEqual(a, tc.expected)
			}
			if !pass {
				t.Fatalf("value comparison failed. Expected %#v got %#v", tc.expected, gotVal.Interface())
			}

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

func TestReadme(t *testing.T) {
	eq := testingx.AssertEqualFn(t)

	c := NewCaster()

	{
		got, err := c.AsInt8(int32(10))
		eq(got, int8(10))
		eq(err, nil)
	}

	{
		got, err := c.AsInt8("10")
		eq(got, int8(10))
		eq(err, nil)
	}

	{
		got, err := c.AsInt8(128)
		eq(got, int8(-128))
		eq(errors.Is(err, ErrCastOverFlow), true)
	}

	{
		got, err := c.AsInt8(10.12)
		eq(got, int8(10))
		eq(errors.Is(err, ErrCastLostDecimals), true)
	}

	{
		got, err := c.AsInt8(float64(10.00))
		eq(got, int8(10))
		eq(err, nil)
	}
}
