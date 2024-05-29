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

type casterTestCaseMel[T any] struct {
	Input                 any
	Expected              T
	ExpectedErr           func(*testing.T, error)
	OverwriteDirectCastFn func(any) (T, error)
	Caster                Caster
	OmitCastByDirectFn    bool
	OmitCastByKind        bool
	OmitCastByType        bool
}

func (c *casterTestCaseMel[T]) SetInput(i any) {
	c.Input = i
}

func (c *casterTestCaseMel[T]) Test(t *testing.T) {
	t.Helper()
	tps := newDirectCastFunctionsTypes()

	typeOfExpected := reflect.TypeOf(c.Expected)

	if c.OverwriteDirectCastFn != nil {
		t.Run(fmt.Sprintf("to_(%s)_custom_direct", typeOfExpected.String()), func(t *testing.T) {
			got, gotErr := c.OverwriteDirectCastFn(c.Input)
			testingx.AssertError(t, c.ExpectedErr, gotErr)
			testingx.AssertEqual(t, got, c.Expected)
		})
	} else if !c.OmitCastByDirectFn {
		t.Run(fmt.Sprintf("to_(%s)_direct", typeOfExpected.String()), func(t *testing.T) {
			var got any
			var gotErr error
			switch typeOfExpected {
			case tps.typeOfBool:
				got, gotErr = c.Caster.AsBool(c.Input)
			// case tps.typeOfByte: // there is no distinguish type for byte. Its only uint8.
			// 	got, gotErr = c.Caster.AsByte(c.Input)
			case tps.typeOfInt8:
				got, gotErr = c.Caster.AsInt8(c.Input)
			case tps.typeOfInt16:
				got, gotErr = c.Caster.AsInt16(c.Input)
			case tps.typeOfInt32:
				got, gotErr = c.Caster.AsInt32(c.Input)
			case tps.typeOfInt64:
				got, gotErr = c.Caster.AsInt64(c.Input)
			case tps.typeOfInt:
				got, gotErr = c.Caster.AsInt(c.Input)
			case tps.typeOfUint8:
				got, gotErr = c.Caster.AsUint8(c.Input)
			case tps.typeOfUint16:
				got, gotErr = c.Caster.AsUint16(c.Input)
			case tps.typeOfUint32:
				got, gotErr = c.Caster.AsUint32(c.Input)
			case tps.typeOfUint64:
				got, gotErr = c.Caster.AsUint64(c.Input)
			case tps.typeOfUint:
				got, gotErr = c.Caster.AsUint(c.Input)
			case tps.typeOfFloat32:
				got, gotErr = c.Caster.AsFloat32(c.Input)
			case tps.typeOfFloat64:
				got, gotErr = c.Caster.AsFloat64(c.Input)
			case tps.typeOfString:
				got, gotErr = c.Caster.AsString(c.Input)
			case tps.typeOfTime:
				got, gotErr = c.Caster.AsTime(c.Input)
			case tps.typeOfDuration:
				got, gotErr = c.Caster.AsDuration(c.Input)

			case tps.typeOfSliceBool:
				got, gotErr = c.Caster.AsBoolSlice(c.Input)
			// case tps.typeOfSliceByte:
			// 	got, gotErr = c.Caster.AsByteSlice(c.Input)
			case tps.typeOfSliceInt8:
				got, gotErr = c.Caster.AsInt8Slice(c.Input)
			case tps.typeOfSliceInt16:
				got, gotErr = c.Caster.AsInt16Slice(c.Input)
			case tps.typeOfSliceInt32:
				got, gotErr = c.Caster.AsInt32Slice(c.Input)
			case tps.typeOfSliceInt64:
				got, gotErr = c.Caster.AsInt64Slice(c.Input)
			case tps.typeOfSliceInt:
				got, gotErr = c.Caster.AsIntSlice(c.Input)
			case tps.typeOfSliceUint8:
				got, gotErr = c.Caster.AsUint8Slice(c.Input)
			case tps.typeOfSliceUint16:
				got, gotErr = c.Caster.AsUint16Slice(c.Input)
			case tps.typeOfSliceUint32:
				got, gotErr = c.Caster.AsUint32Slice(c.Input)
			case tps.typeOfSliceUint64:
				got, gotErr = c.Caster.AsUint64Slice(c.Input)
			case tps.typeOfSliceUint:
				got, gotErr = c.Caster.AsUintSlice(c.Input)
			case tps.typeOfSliceFloat32:
				got, gotErr = c.Caster.AsFloat32Slice(c.Input)
			case tps.typeOfSliceFloat64:
				got, gotErr = c.Caster.AsFloat64Slice(c.Input)
			case tps.typeOfSliceString:
				got, gotErr = c.Caster.AsStringSlice(c.Input)
			case tps.typeOfSliceTime:
				got, gotErr = c.Caster.AsTimeSlice(c.Input)
			case tps.typeOfSliceDuration:
				got, gotErr = c.Caster.AsDurationSlice(c.Input)

			default:
				t.SkipNow()
			}

			testingx.AssertError(t, c.ExpectedErr, gotErr)
			testingx.AssertEqual(t, got, c.Expected)
		})
	}

	if !c.OmitCastByKind {
		t.Run(fmt.Sprintf("to_(%s)_by_kind", typeOfExpected.String()), func(t *testing.T) {
			var got any
			var gotErr error
			switch typeOfExpected {
			case tps.typeOfBool:
				got, gotErr = c.Caster.As(c.Input, reflect.Bool)
			case tps.typeOfInt8:
				got, gotErr = c.Caster.As(c.Input, reflect.Int8)
			case tps.typeOfInt16:
				got, gotErr = c.Caster.As(c.Input, reflect.Int16)
			case tps.typeOfInt32:
				got, gotErr = c.Caster.As(c.Input, reflect.Int32)
			case tps.typeOfInt64:
				got, gotErr = c.Caster.As(c.Input, reflect.Int64)
			case tps.typeOfInt:
				got, gotErr = c.Caster.As(c.Input, reflect.Int)
			case tps.typeOfUint8:
				got, gotErr = c.Caster.As(c.Input, reflect.Uint8)
			case tps.typeOfUint16:
				got, gotErr = c.Caster.As(c.Input, reflect.Uint16)
			case tps.typeOfUint32:
				got, gotErr = c.Caster.As(c.Input, reflect.Uint32)
			case tps.typeOfUint64:
				got, gotErr = c.Caster.As(c.Input, reflect.Uint64)
			case tps.typeOfUint:
				got, gotErr = c.Caster.As(c.Input, reflect.Uint)
			case tps.typeOfFloat32:
				got, gotErr = c.Caster.As(c.Input, reflect.Float32)
			case tps.typeOfFloat64:
				got, gotErr = c.Caster.As(c.Input, reflect.Float64)
			case tps.typeOfString:
				got, gotErr = c.Caster.As(c.Input, reflect.String)
			default:
				t.SkipNow()
			}

			testingx.AssertError(t, c.ExpectedErr, gotErr)
			testingx.AssertEqual(t, got, c.Expected)
		})
	}

	if !c.OmitCastByType {
		t.Run(fmt.Sprintf("to_(%s)_by_type", typeOfExpected.String()), func(t *testing.T) {
			got, gotErr := c.Caster.ByType(c.Input, typeOfExpected)
			testingx.AssertError(t, c.ExpectedErr, gotErr)
			testingx.AssertEqual(t, got, c.Expected)
		})
	}
}

type singleCastTestCase[T any] struct {
	input        any
	expectedErr  func(*testing.T, error)
	expected     T
	directCastFn func(any) (T, error)
}

func runSingleCastTestCases[T any](t *testing.T, testCases []singleCastTestCase[T], defaultCastFn func(any) (T, error)) {
	t.Helper()
	for idx, tc := range testCases {
		tc := tc
		name := fmt.Sprintf("index:%d %s", idx, testingx.Format(tc.input))
		t.Run(name, func(t *testing.T) {
			// t.Helper()
			// t.Parallel()
			var (
				got    T
				gotErr error
			)
			if tc.directCastFn != nil {
				got, gotErr = tc.directCastFn(tc.input)
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
			fn:          tryReflectConvert[int16],
			input:       intAlias(13),
			expected:    int16(13),
			expectedErr: nil,
		},
		"struct to int16 expect error": {
			fn:          tryReflectConvert[int16],
			input:       struct{}{},
			expected:    int16(0),
			expectedErr: testingx.ExpectedErrorIs(ErrInvalidType),
		},
		"string to []byte": {
			fn:          tryReflectConvert[[]byte],
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

func TestByType(t *testing.T) {
	t.Parallel()

	c := NewCaster()

	type aliasInt int
	type aliasString string
	type aliasString2 aliasString

	type Foo struct{ A int32 }
	type AliasFoo Foo

	tests := []struct {
		input         any
		expected      any
		expectedError func(*testing.T, error)
	}{
		{
			input:         int(123),
			expected:      uint(123),
			expectedError: nil,
		},
		{
			input:         aliasInt(123),
			expected:      uint(123),
			expectedError: nil,
		},
		{
			input:         []aliasInt{1, 2, 3, 4},
			expected:      []int{1, 2, 3, 4},
			expectedError: nil,
		},
		{
			input:         []int16{1, 2, 3, 4},
			expected:      []int32{1, 2, 3, 4},
			expectedError: nil,
		},
		{
			input:         []int{1, 2, 3, 4},
			expected:      []aliasInt{1, 2, 3, 4},
			expectedError: nil,
		},
		{
			input:         aliasString("123"),
			expected:      uint(123),
			expectedError: nil,
		},
		{
			input:         aliasString2("123"),
			expected:      uint(123),
			expectedError: nil,
		},
		{
			input:         byte(123),
			expected:      uint16(123),
			expectedError: nil,
		},
		{
			input:         Foo{A: 123},
			expected:      AliasFoo{A: 123},
			expectedError: nil,
		},
		{
			input:         AliasFoo{A: 123},
			expected:      Foo{A: 123},
			expectedError: nil,
		},
		{
			input:         []AliasFoo{{A: 121}, {A: 122}, {A: 123}},
			expected:      []Foo{{A: 121}, {A: 122}, {A: 123}},
			expectedError: nil,
		},
	}

	for idx, tc := range tests {
		tc := tc
		fromType := reflect.TypeOf(tc.input).String()
		toType := reflect.TypeOf(tc.expected).String()
		t.Run(fmt.Sprintf("#%d#%s->%s", idx, fromType, toType), func(t *testing.T) {
			t.Parallel()
			got, err := c.ByType(tc.input, reflect.TypeOf(tc.expected))
			testingx.AssertError(t, tc.expectedError, err)
			testingx.AssertEqual(t, got, tc.expected)
		})
	}
}
