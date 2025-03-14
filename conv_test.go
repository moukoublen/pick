package pick

import (
	"fmt"
	"math"
	"reflect"
	"testing"

	"github.com/moukoublen/pick/internal/tst"
)

var (
	expectOverFlowError = tst.ExpectedErrorChecks(
		tst.ExpectedErrorOfType[*ConvertError](),
		tst.ExpectedErrorIs(ErrConvertOverFlow),
	)
	expectLostDecimals = tst.ExpectedErrorChecks(
		tst.ExpectedErrorOfType[*ConvertError](),
		tst.ExpectedErrorIs(ErrConvertLostDecimals),
	)
	expectMalformedSyntax = tst.ExpectedErrorChecks(
		tst.ExpectedErrorOfType[*ConvertError](),
		tst.ExpectedErrorIs(ErrConvertInvalidSyntax),
	)
	expectInvalidType = tst.ExpectedErrorChecks(
		tst.ExpectedErrorOfType[*ConvertError](),
		tst.ExpectedErrorIs(ErrConvertInvalidType),
	)
)

type converterTestCase[T any] struct {
	Input                    any
	Expected                 T
	ErrorAsserter            tst.ErrorAsserter
	OverwriteDirectConvertFn func(any) (T, error)
	Converter                DefaultConverter
	OmitConvertByDirectFn    bool
	OmitConvertByKind        bool
	OmitConvertByType        bool
}

func (c *converterTestCase[T]) SetInput(i any) {
	c.Input = i
}

func (c *converterTestCase[T]) Test(t *testing.T) {
	t.Helper()
	tps := newDirectConvertFunctionsTypes()

	typeOfExpected := reflect.TypeOf(c.Expected)

	if c.ErrorAsserter == nil {
		c.ErrorAsserter = tst.NoError
	}

	if c.OverwriteDirectConvertFn != nil {
		t.Run(fmt.Sprintf("to_(%s)_custom_direct", typeOfExpected.String()), func(t *testing.T) {
			got, gotErr := c.OverwriteDirectConvertFn(c.Input)
			c.ErrorAsserter(t, gotErr)
			tst.AssertEqual(t, got, c.Expected)
		})
	} else if !c.OmitConvertByDirectFn {
		t.Run(fmt.Sprintf("to_(%s)_direct", typeOfExpected.String()), func(t *testing.T) {
			var got any
			var gotErr error
			switch typeOfExpected {
			case tps.typeOfBool:
				got, gotErr = c.Converter.AsBool(c.Input)
			// case tps.typeOfByte: // there is no distinguish type for byte. Its only uint8.
			// 	got, gotErr = c.Converter.AsByte(c.Input)
			case tps.typeOfInt8:
				got, gotErr = c.Converter.AsInt8(c.Input)
			case tps.typeOfInt16:
				got, gotErr = c.Converter.AsInt16(c.Input)
			case tps.typeOfInt32:
				got, gotErr = c.Converter.AsInt32(c.Input)
			case tps.typeOfInt64:
				got, gotErr = c.Converter.AsInt64(c.Input)
			case tps.typeOfInt:
				got, gotErr = c.Converter.AsInt(c.Input)
			case tps.typeOfUint8:
				got, gotErr = c.Converter.AsUint8(c.Input)
			case tps.typeOfUint16:
				got, gotErr = c.Converter.AsUint16(c.Input)
			case tps.typeOfUint32:
				got, gotErr = c.Converter.AsUint32(c.Input)
			case tps.typeOfUint64:
				got, gotErr = c.Converter.AsUint64(c.Input)
			case tps.typeOfUint:
				got, gotErr = c.Converter.AsUint(c.Input)
			case tps.typeOfFloat32:
				got, gotErr = c.Converter.AsFloat32(c.Input)
			case tps.typeOfFloat64:
				got, gotErr = c.Converter.AsFloat64(c.Input)
			case tps.typeOfString:
				got, gotErr = c.Converter.AsString(c.Input)
			case tps.typeOfTime:
				got, gotErr = c.Converter.AsTime(c.Input)
			case tps.typeOfDuration:
				got, gotErr = c.Converter.AsDuration(c.Input)

			case tps.typeOfSliceBool:
				got, gotErr = c.Converter.AsBoolSlice(c.Input)
			// case tps.typeOfSliceByte:
			// 	got, gotErr = c.Converter.AsByteSlice(c.Input)
			case tps.typeOfSliceInt8:
				got, gotErr = c.Converter.AsInt8Slice(c.Input)
			case tps.typeOfSliceInt16:
				got, gotErr = c.Converter.AsInt16Slice(c.Input)
			case tps.typeOfSliceInt32:
				got, gotErr = c.Converter.AsInt32Slice(c.Input)
			case tps.typeOfSliceInt64:
				got, gotErr = c.Converter.AsInt64Slice(c.Input)
			case tps.typeOfSliceInt:
				got, gotErr = c.Converter.AsIntSlice(c.Input)
			case tps.typeOfSliceUint8:
				got, gotErr = c.Converter.AsUint8Slice(c.Input)
			case tps.typeOfSliceUint16:
				got, gotErr = c.Converter.AsUint16Slice(c.Input)
			case tps.typeOfSliceUint32:
				got, gotErr = c.Converter.AsUint32Slice(c.Input)
			case tps.typeOfSliceUint64:
				got, gotErr = c.Converter.AsUint64Slice(c.Input)
			case tps.typeOfSliceUint:
				got, gotErr = c.Converter.AsUintSlice(c.Input)
			case tps.typeOfSliceFloat32:
				got, gotErr = c.Converter.AsFloat32Slice(c.Input)
			case tps.typeOfSliceFloat64:
				got, gotErr = c.Converter.AsFloat64Slice(c.Input)
			case tps.typeOfSliceString:
				got, gotErr = c.Converter.AsStringSlice(c.Input)
			case tps.typeOfSliceTime:
				got, gotErr = c.Converter.AsTimeSlice(c.Input)
			case tps.typeOfSliceDuration:
				got, gotErr = c.Converter.AsDurationSlice(c.Input)

			default:
				t.SkipNow()
			}

			c.ErrorAsserter(t, gotErr)
			tst.AssertEqual(t, got, c.Expected)
		})
	}

	if !c.OmitConvertByKind {
		t.Run(fmt.Sprintf("to_(%s)_by_kind", typeOfExpected.String()), func(t *testing.T) {
			var got any
			var gotErr error
			switch typeOfExpected {
			case tps.typeOfBool:
				got, gotErr = c.Converter.As(c.Input, reflect.Bool)
			case tps.typeOfInt8:
				got, gotErr = c.Converter.As(c.Input, reflect.Int8)
			case tps.typeOfInt16:
				got, gotErr = c.Converter.As(c.Input, reflect.Int16)
			case tps.typeOfInt32:
				got, gotErr = c.Converter.As(c.Input, reflect.Int32)
			case tps.typeOfInt64:
				got, gotErr = c.Converter.As(c.Input, reflect.Int64)
			case tps.typeOfInt:
				got, gotErr = c.Converter.As(c.Input, reflect.Int)
			case tps.typeOfUint8:
				got, gotErr = c.Converter.As(c.Input, reflect.Uint8)
			case tps.typeOfUint16:
				got, gotErr = c.Converter.As(c.Input, reflect.Uint16)
			case tps.typeOfUint32:
				got, gotErr = c.Converter.As(c.Input, reflect.Uint32)
			case tps.typeOfUint64:
				got, gotErr = c.Converter.As(c.Input, reflect.Uint64)
			case tps.typeOfUint:
				got, gotErr = c.Converter.As(c.Input, reflect.Uint)
			case tps.typeOfFloat32:
				got, gotErr = c.Converter.As(c.Input, reflect.Float32)
			case tps.typeOfFloat64:
				got, gotErr = c.Converter.As(c.Input, reflect.Float64)
			case tps.typeOfString:
				got, gotErr = c.Converter.As(c.Input, reflect.String)
			default:
				t.SkipNow()
			}

			c.ErrorAsserter(t, gotErr)
			tst.AssertEqual(t, got, c.Expected)
		})
	}

	if !c.OmitConvertByType {
		t.Run(fmt.Sprintf("to_(%s)_by_type", typeOfExpected.String()), func(t *testing.T) {
			got, gotErr := c.Converter.ByType(c.Input, typeOfExpected)
			c.ErrorAsserter(t, gotErr)
			tst.AssertEqual(t, got, c.Expected)
		})
	}
}

type singleConvertTestCase[T any] struct {
	input           any
	errorAsserter   tst.ErrorAsserter
	expected        T
	directConvertFn func(any) (T, error)
}

func runSingleConvertTestCases[T any](t *testing.T, testCases []singleConvertTestCase[T], defaultConvertFn func(any) (T, error)) {
	t.Helper()
	for idx, tc := range testCases {
		name := fmt.Sprintf("index:%d %s", idx, tst.Format(tc.input))
		t.Run(name, func(t *testing.T) {
			// t.Helper()
			// t.Parallel()
			var (
				got    T
				gotErr error
			)
			if tc.directConvertFn != nil {
				got, gotErr = tc.directConvertFn(tc.input)
			} else {
				got, gotErr = defaultConvertFn(tc.input)
			}
			tc.errorAsserter(t, gotErr)

			tst.AssertEqual(t, got, tc.expected)
		})
	}
}

func TestTryConvertUsingReflect(t *testing.T) {
	type intAlias int
	tests := map[string]struct {
		fn            any
		input         any
		expected      any
		errorAsserter tst.ErrorAsserter
	}{
		"intAlias to int16": {
			fn:            tryReflectConvert[int16],
			input:         intAlias(13),
			expected:      int16(13),
			errorAsserter: tst.NoError,
		},
		"struct to int16 expect error": {
			fn:            tryReflectConvert[int16],
			input:         struct{}{},
			expected:      int16(0),
			errorAsserter: tst.ExpectedErrorIs(ErrConvertInvalidType),
		},
		"string to []byte": {
			fn:            tryReflectConvert[[]byte],
			input:         "str",
			expected:      []byte("str"),
			errorAsserter: tst.NoError,
		},
	}

	for name, tc := range tests {
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
				pass = tst.Compare(a, tc.expected)
			}
			if !pass {
				t.Fatalf("value comparison failed. Expected %#v got %#v", tc.expected, gotVal.Interface())
			}

			errVal := returnedVals[1]
			errInf := errVal.Interface()
			if errInf == nil {
				tc.errorAsserter(t, nil)
			} else {
				err, is := errInf.(error)
				if !is {
					t.Errorf("second returned item is not of type error. Type: %s", errVal.Type().String())
				}
				tc.errorAsserter(t, err)
			}
		})
	}
}

func TestByType(t *testing.T) {
	t.Parallel()

	c := NewDefaultConverter()

	type aliasInt int
	type aliasString string
	type aliasString2 aliasString

	type Foo struct{ A int32 }
	type AliasFoo Foo

	tests := []struct {
		input         any
		expected      any
		errorAsserter tst.ErrorAsserter
	}{
		{
			input:         int(123),
			expected:      uint(123),
			errorAsserter: tst.NoError,
		},
		{
			input:         aliasInt(123),
			expected:      uint(123),
			errorAsserter: tst.NoError,
		},
		{
			input:         []aliasInt{1, 2, 3, 4},
			expected:      []int{1, 2, 3, 4},
			errorAsserter: tst.NoError,
		},
		{
			input:         []int16{1, 2, 3, 4},
			expected:      []int32{1, 2, 3, 4},
			errorAsserter: tst.NoError,
		},
		{
			input:         []int{1, 2, 3, 4},
			expected:      []aliasInt{1, 2, 3, 4},
			errorAsserter: tst.NoError,
		},
		{
			input:         123,
			expected:      []int{123},
			errorAsserter: tst.NoError,
		},
		{
			input:         aliasString("123"),
			expected:      uint(123),
			errorAsserter: tst.NoError,
		},
		{
			input:         aliasString2("123"),
			expected:      uint(123),
			errorAsserter: tst.NoError,
		},
		{
			input:         byte(123),
			expected:      uint16(123),
			errorAsserter: tst.NoError,
		},
		{
			input:         Foo{A: 123},
			expected:      AliasFoo{A: 123},
			errorAsserter: tst.NoError,
		},
		{
			input:         AliasFoo{A: 123},
			expected:      Foo{A: 123},
			errorAsserter: tst.NoError,
		},
		{
			input:         []AliasFoo{{A: 121}, {A: 122}, {A: 123}},
			expected:      []Foo{{A: 121}, {A: 122}, {A: 123}},
			errorAsserter: tst.NoError,
		},
	}

	for idx, tc := range tests {
		fromType := reflect.TypeOf(tc.input).String()
		toType := reflect.TypeOf(tc.expected).String()
		t.Run(fmt.Sprintf("#%d#%s->%s", idx, fromType, toType), func(t *testing.T) {
			t.Parallel()
			got, err := c.ByType(tc.input, reflect.TypeOf(tc.expected))
			tc.errorAsserter(t, err)
			tst.AssertEqual(t, got, tc.expected)
		})
	}
}

func TestConvertGeneric(t *testing.T) {
	type stringAlias string
	type mapStringAny map[string]any

	tests := []struct {
		convertFn     func() (any, error)
		expected      any
		errorAsserter tst.ErrorAsserter
	}{
		0: {
			convertFn: func() (any, error) {
				return Convert[string](1)
			},
			expected:      string("1"),
			errorAsserter: tst.NoError,
		},
		1: {
			convertFn: func() (any, error) {
				return Convert[int64]("1234567")
			},
			expected:      int64(1234567),
			errorAsserter: tst.NoError,
		},
		2: {
			convertFn: func() (any, error) {
				return Convert[[]uint8]([]int64{1, 2, 3, 4, 5, 6, 7})
			},
			expected:      []uint8{1, 2, 3, 4, 5, 6, 7},
			errorAsserter: tst.NoError,
		},
		3: {
			convertFn: func() (any, error) {
				return Convert[[]stringAlias]([]string{"one", "two"})
			},
			expected:      []stringAlias{"one", "two"},
			errorAsserter: tst.NoError,
		},
		4: {
			convertFn: func() (any, error) {
				return Convert[[]string]([]stringAlias{"one", "two"})
			},
			expected:      []string{"one", "two"},
			errorAsserter: tst.NoError,
		},
		5: {
			convertFn: func() (any, error) {
				m := map[string]any{"one": 1, "two": 2}
				return Convert[map[string]string](m)
			},
			expected:      map[string]string{"one": "1", "two": "2"},
			errorAsserter: tst.NoError,
		},
		6: {
			convertFn: func() (any, error) {
				m := map[string]int{"one": 1, "two": 2}
				return Convert[map[string]string](m)
			},
			expected:      map[string]string{"one": "1", "two": "2"},
			errorAsserter: tst.NoError,
		},
		7: {
			convertFn: func() (any, error) {
				m := map[string]bool{"one": true, "two": false}
				return Convert[map[string]int](m)
			},
			expected:      map[string]int{"one": 1, "two": 0},
			errorAsserter: tst.NoError,
		},
		8: {
			convertFn: func() (any, error) {
				m := map[string]string{"1": "1", "2": "2"}
				return Convert[map[int]string](m)
			},
			expected:      map[int]string{1: "1", 2: "2"},
			errorAsserter: tst.NoError,
		},
		9: {
			convertFn: func() (any, error) {
				return Convert[map[string]any](int(42))
			},
			expected:      map[string]any(nil),
			errorAsserter: tst.ExpectedErrorIs(ErrConvertInvalidType),
		},
		10: {
			convertFn: func() (any, error) {
				m := map[string]any{"one": "a", "two": "b"}
				return Convert[map[string]string](m)
			},
			expected:      map[string]string{"one": "a", "two": "b"},
			errorAsserter: tst.NoError,
		},
		11: {
			convertFn: func() (any, error) {
				m := map[string]any{"one": []any{"a", "b"}, "two": []any{"c", "d"}}
				return Convert[map[string][]string](m)
			},
			expected:      map[string][]string{"one": {"a", "b"}, "two": {"c", "d"}},
			errorAsserter: tst.NoError,
		},
		12: {
			convertFn: func() (any, error) {
				m := map[string]any{"one": "a", "two": []string{"abc", "def"}}
				return Convert[map[string]string](m)
			},
			expected:      map[string]string(nil),
			errorAsserter: tst.ExpectedErrorIs(ErrConvertInvalidType),
		},
		13: {
			convertFn: func() (any, error) {
				m := map[string]any{"one": 1, "two": 2}
				return Convert[map[string]any](m)
			},
			expected:      map[string]any{"one": 1, "two": 2},
			errorAsserter: tst.NoError,
		},
		14: { // key cast error
			convertFn: func() (any, error) {
				m := map[int32]any{math.MaxInt32: "abc"}
				return Convert[map[int8]any](m)
			},
			expected:      map[int8]any(nil),
			errorAsserter: tst.ExpectedErrorIs(ErrConvertOverFlow),
		},
		15: { // value cast error
			convertFn: func() (any, error) {
				m := map[string]int32{"abc": math.MaxInt32}
				return Convert[map[string]int8](m)
			},
			expected:      map[string]int8(nil),
			errorAsserter: tst.ExpectedErrorIs(ErrConvertOverFlow),
		},
		16: {
			convertFn: func() (any, error) {
				s := []string{"a", "b", "c"}
				return Convert[map[int]string](s)
			},
			expected:      map[int]string{0: "a", 1: "b", 2: "c"},
			errorAsserter: tst.NoError,
		},
		17: {
			convertFn: func() (any, error) {
				s := []string{"a", "b", "c"}
				return Convert[map[int]any](s)
			},
			expected:      map[int]any{0: "a", 1: "b", 2: "c"},
			errorAsserter: tst.NoError,
		},
		18: {
			convertFn: func() (any, error) {
				s := []string{"a", "b", "c"}
				return Convert[map[string]string](s)
			},
			expected:      map[string]string{"0": "a", "1": "b", "2": "c"},
			errorAsserter: tst.NoError,
		},
		19: { // from slice key convert error
			convertFn: func() (any, error) {
				s := make([]string, 1000)
				return Convert[map[int8]string](s)
			},
			expected:      map[int8]string(nil),
			errorAsserter: tst.ExpectedErrorIs(ErrConvertOverFlow),
		},
		20: { // from slice value convert error
			convertFn: func() (any, error) {
				s := []string{"1000"}
				return Convert[map[int8]int8](s)
			},
			expected:      map[int8]int8(nil),
			errorAsserter: tst.ExpectedErrorIs(ErrConvertOverFlow),
		},
		21: {
			convertFn: func() (any, error) {
				m := mapStringAny{"a": "b"}
				return Convert[map[string]any](m)
			},
			expected:      map[string]any{"a": "b"},
			errorAsserter: tst.NoError,
		},
	}

	for idx, tc := range tests {
		t.Run(fmt.Sprintf("[%d]_%T", idx, tc.expected), func(t *testing.T) {
			got, err := tc.convertFn()
			tc.errorAsserter(t, err)
			tst.AssertEqual(t, got, tc.expected)
		})
	}
}
