package cast

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"runtime"
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

type castTestExpectedResult[Output any] struct {
	expectedResult Output
	compareFn      func(x, y any) bool
	errorAssertFn  func(*testing.T, error)
	shouldRun      bool
}

// constructors for shortage.
func newCastTestExpectedResultConstructor[Output any](compareFn func(x, y any) bool) func(result Output, errorAssertFn func(*testing.T, error)) castTestExpectedResult[Output] {
	return func(result Output, errorAssertFn func(*testing.T, error)) castTestExpectedResult[Output] {
		return castTestExpectedResult[Output]{expectedResult: result, compareFn: compareFn, errorAssertFn: errorAssertFn, shouldRun: true}
	}
}

func splitBasedOnArch[Output any](for32bit, for64bit castTestExpectedResult[Output]) castTestExpectedResult[Output] {
	switch runtime.GOARCH {
	case "arm", "386":
		return for32bit
	default:
		return for64bit
	}
}

func TestCasterMatrix(t *testing.T) {
	t.Parallel()

	// matrixExpectedResult constructor function aliases (prefix `ex` from expected result)
	exByte := newCastTestExpectedResultConstructor[byte](reflect.DeepEqual)
	exInt8 := newCastTestExpectedResultConstructor[int8](reflect.DeepEqual)
	exInt16 := newCastTestExpectedResultConstructor[int16](reflect.DeepEqual)
	exInt32 := newCastTestExpectedResultConstructor[int32](reflect.DeepEqual)
	exInt64 := newCastTestExpectedResultConstructor[int64](reflect.DeepEqual)
	exInt := newCastTestExpectedResultConstructor[int](reflect.DeepEqual)
	exUInt8 := newCastTestExpectedResultConstructor[uint8](reflect.DeepEqual)
	exUint16 := newCastTestExpectedResultConstructor[uint16](reflect.DeepEqual)
	exUint32 := newCastTestExpectedResultConstructor[uint32](reflect.DeepEqual)
	exUint64 := newCastTestExpectedResultConstructor[uint64](reflect.DeepEqual)
	exUint := newCastTestExpectedResultConstructor[uint](reflect.DeepEqual)
	exFloat32 := newCastTestExpectedResultConstructor[float32](compareFloat32)
	exFloat64 := newCastTestExpectedResultConstructor[float64](compareFloat64)
	exString := newCastTestExpectedResultConstructor[string](reflect.DeepEqual)

	testCases := []struct {
		Input any

		Byte    castTestExpectedResult[byte]
		Int8    castTestExpectedResult[int8]
		Int16   castTestExpectedResult[int16]
		Int32   castTestExpectedResult[int32]
		Int64   castTestExpectedResult[int64]
		Int     castTestExpectedResult[int]
		Uint8   castTestExpectedResult[uint8]
		Uint16  castTestExpectedResult[uint16]
		Uint32  castTestExpectedResult[uint32]
		Uint64  castTestExpectedResult[uint64]
		Uint    castTestExpectedResult[uint]
		Float32 castTestExpectedResult[float32]
		Float64 castTestExpectedResult[float64]
		String  castTestExpectedResult[string]
	}{
		{
			Input:   nil,
			Byte:    exByte(0, nil),
			Int8:    exInt8(0, nil),
			Int16:   exInt16(0, nil),
			Int32:   exInt32(0, nil),
			Int64:   exInt64(0, nil),
			Int:     exInt(0, nil),
			Uint8:   exUInt8(0, nil),
			Uint16:  exUint16(0, nil),
			Uint32:  exUint32(0, nil),
			Uint64:  exUint64(0, nil),
			Uint:    exUint(0, nil),
			Float32: exFloat32(0, nil),
			Float64: exFloat64(0, nil),
			String:  exString("", nil),
		},
		{
			Input:   int8(12),
			Byte:    exByte(12, nil),
			Int8:    exInt8(12, nil),
			Int16:   exInt16(12, nil),
			Int32:   exInt32(12, nil),
			Int64:   exInt64(12, nil),
			Int:     exInt(12, nil),
			Uint8:   exUInt8(12, nil),
			Uint16:  exUint16(12, nil),
			Uint32:  exUint32(12, nil),
			Uint64:  exUint64(12, nil),
			Uint:    exUint(12, nil),
			Float32: exFloat32(12, nil),
			Float64: exFloat64(12, nil),
			String:  exString("12", nil),
		},
		{
			Input:   int8(math.MaxInt8),
			Byte:    exByte(127, nil),
			Int8:    exInt8(127, nil),
			Int16:   exInt16(127, nil),
			Int32:   exInt32(127, nil),
			Int64:   exInt64(127, nil),
			Int:     exInt(127, nil),
			Uint8:   exUInt8(127, nil),
			Uint16:  exUint16(127, nil),
			Uint32:  exUint32(127, nil),
			Uint64:  exUint64(127, nil),
			Uint:    exUint(127, nil),
			Float32: exFloat32(127, nil),
			Float64: exFloat64(127, nil),
			String:  exString("127", nil),
		},
		{
			Input:   int8(math.MinInt8),
			Byte:    exByte(0x80, expectOverFlowError),
			Int8:    exInt8(-128, nil),
			Int16:   exInt16(-128, nil),
			Int32:   exInt32(-128, nil),
			Int64:   exInt64(-128, nil),
			Int:     exInt(-128, nil),
			Uint8:   exUInt8(0x80, expectOverFlowError),
			Uint16:  exUint16(0xff80, expectOverFlowError),
			Uint32:  exUint32(0xffffff80, expectOverFlowError),
			Uint64:  exUint64(0xffffffffffffff80, expectOverFlowError),
			Uint:    exUint(0xffffffffffffff80, expectOverFlowError),
			Float32: exFloat32(-128, nil),
			Float64: exFloat64(-128, nil),
			String:  exString("-128", nil),
		},
		{
			Input:   int16(math.MaxInt16),
			Byte:    exByte(255, expectOverFlowError),
			Int8:    exInt8(-1, expectOverFlowError),
			Int16:   exInt16(32767, nil),
			Int32:   exInt32(32767, nil),
			Int64:   exInt64(32767, nil),
			Int:     exInt(32767, nil),
			Uint8:   exUInt8(255, expectOverFlowError),
			Uint16:  exUint16(32767, nil),
			Uint32:  exUint32(32767, nil),
			Uint64:  exUint64(32767, nil),
			Uint:    exUint(32767, nil),
			Float32: exFloat32(32767, nil),
			Float64: exFloat64(32767, nil),
			String:  exString("32767", nil),
		},
		{
			Input:   int16(math.MinInt16),
			Byte:    exByte(0, expectOverFlowError),
			Int8:    exInt8(0, expectOverFlowError),
			Int16:   exInt16(-32768, nil),
			Int32:   exInt32(-32768, nil),
			Int64:   exInt64(-32768, nil),
			Int:     exInt(-32768, nil),
			Uint8:   exUInt8(0, expectOverFlowError),
			Uint16:  exUint16(0x8000, expectOverFlowError),
			Uint32:  exUint32(0xffff8000, expectOverFlowError),
			Uint64:  exUint64(0xffffffffffff8000, expectOverFlowError),
			Uint:    exUint(0xffffffffffff8000, expectOverFlowError),
			Float32: exFloat32(-32768, nil),
			Float64: exFloat64(-32768, nil),
			String:  exString("-32768", nil),
		},
		{
			Input:   int32(math.MaxInt32),
			Byte:    exByte(255, expectOverFlowError),
			Int8:    exInt8(-1, expectOverFlowError),
			Int16:   exInt16(-1, expectOverFlowError),
			Int32:   exInt32(2147483647, nil),
			Int64:   exInt64(2147483647, nil),
			Int:     exInt(2147483647, nil),
			Uint8:   exUInt8(255, expectOverFlowError),
			Uint16:  exUint16(0xffff, expectOverFlowError),
			Uint32:  exUint32(2147483647, nil),
			Uint64:  exUint64(2147483647, nil),
			Uint:    exUint(2147483647, nil),
			Float32: exFloat32(2147483647, nil),
			Float64: exFloat64(2147483647, nil),
			String:  exString("2147483647", nil),
		},
		{
			Input:   int32(math.MinInt32),
			Byte:    exByte(0, expectOverFlowError),
			Int8:    exInt8(0, expectOverFlowError),
			Int16:   exInt16(0, expectOverFlowError),
			Int32:   exInt32(-2147483648, nil),
			Int64:   exInt64(-2147483648, nil),
			Int:     exInt(-2147483648, nil),
			Uint8:   exUInt8(0, expectOverFlowError),
			Uint16:  exUint16(0, expectOverFlowError),
			Uint32:  exUint32(0x80000000, expectOverFlowError),
			Uint64:  exUint64(0xffffffff80000000, expectOverFlowError),
			Uint:    exUint(0xffffffff80000000, expectOverFlowError),
			Float32: exFloat32(-2147483648, nil),
			Float64: exFloat64(-2147483648, nil),
			String:  exString("-2147483648", nil),
		},
		{
			Input:   int64(math.MaxInt64),
			Byte:    exByte(255, expectOverFlowError),
			Int8:    exInt8(-1, expectOverFlowError),
			Int16:   exInt16(-1, expectOverFlowError),
			Int32:   exInt32(-1, expectOverFlowError),
			Int64:   exInt64(math.MaxInt64, nil),
			Int:     splitBasedOnArch(exInt(-1, expectOverFlowError), exInt(math.MaxInt64, nil)),
			Uint8:   exUInt8(255, expectOverFlowError),
			Uint16:  exUint16(0xffff, expectOverFlowError),
			Uint32:  exUint32(0xffffffff, expectOverFlowError),
			Uint64:  exUint64(0x7fffffffffffffff, nil),
			Uint:    splitBasedOnArch(exUint(0xffffffff, expectOverFlowError), exUint(0x7fffffffffffffff, nil)),
			Float32: exFloat32(math.MaxInt64, nil),
			Float64: exFloat64(math.MaxInt64, nil),
			String:  exString("9223372036854775807", nil),
		},
		{
			Input:   int64(math.MinInt64),
			Byte:    exByte(0, expectOverFlowError),
			Int8:    exInt8(0, expectOverFlowError),
			Int16:   exInt16(0, expectOverFlowError),
			Int32:   exInt32(0, expectOverFlowError),
			Int64:   exInt64(math.MinInt64, nil),
			Int:     splitBasedOnArch(exInt(0, expectOverFlowError), exInt(math.MinInt64, nil)),
			Uint8:   exUInt8(0, expectOverFlowError),
			Uint16:  exUint16(0, expectOverFlowError),
			Uint32:  exUint32(0, expectOverFlowError),
			Uint64:  exUint64(0x8000000000000000, expectOverFlowError),
			Uint:    splitBasedOnArch(exUint(0, expectOverFlowError), exUint(0x8000000000000000, expectOverFlowError)),
			Float32: exFloat32(math.MinInt64, nil),
			Float64: exFloat64(math.MinInt64, nil),
			String:  exString("-9223372036854775808", nil),
		},
	}

	caster := NewCaster()
	for idx, tc := range testCases {
		tc := tc
		typeName := "nil"
		if tc.Input != nil {
			tp := reflect.TypeOf(tc.Input)
			typeName = tp.Name()
		}

		name := fmt.Sprintf("index[%d]__inputType[%s]__inputValue[%#v]", idx, typeName, tc.Input)
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			t.Run("caster_byte", matrixSubTest[byte](tc.Input, caster.AsByte, tc.Byte))
			t.Run("caster_int8", matrixSubTest[int8](tc.Input, caster.AsInt8, tc.Int8))
			t.Run("caster_int16", matrixSubTest[int16](tc.Input, caster.AsInt16, tc.Int16))
			t.Run("caster_int32", matrixSubTest[int32](tc.Input, caster.AsInt32, tc.Int32))
			t.Run("caster_int64", matrixSubTest[int64](tc.Input, caster.AsInt64, tc.Int64))
			t.Run("caster_int", matrixSubTest[int](tc.Input, caster.AsInt, tc.Int))
			t.Run("caster_uint8", matrixSubTest[uint8](tc.Input, caster.AsUint8, tc.Uint8))
			t.Run("caster_uint16", matrixSubTest[uint16](tc.Input, caster.AsUint16, tc.Uint16))
			t.Run("caster_uint32", matrixSubTest[uint32](tc.Input, caster.AsUint32, tc.Uint32))
			t.Run("caster_uint64", matrixSubTest[uint64](tc.Input, caster.AsUint64, tc.Uint64))
			t.Run("caster_uint", matrixSubTest[uint](tc.Input, caster.AsUint, tc.Uint))
			t.Run("caster_float32", matrixSubTest[float32](tc.Input, caster.AsFloat32, tc.Float32))
			t.Run("caster_float64", matrixSubTest[float64](tc.Input, caster.AsFloat64, tc.Float64))
			t.Run("caster_string", matrixSubTest[string](tc.Input, caster.AsString, tc.String))
		})
	}
}

func matrixSubTest[Output any](input any, castFn func(any) (Output, error), subTestCase castTestExpectedResult[Output]) func(*testing.T) {
	return func(t *testing.T) {
		t.Helper()
		t.Parallel()
		if !subTestCase.shouldRun {
			t.SkipNow()
		}

		got, gotErr := castFn(input)
		testingx.AssertError(t, subTestCase.errorAssertFn, gotErr)

		compareFn := subTestCase.compareFn
		if compareFn == nil {
			compareFn = reflect.DeepEqual
		}

		if !compareFn(subTestCase.expectedResult, got) {
			t.Errorf("wrong returned value. Expected %#v got %#v", subTestCase.expectedResult, got)
		}
	}
}

func compareFloat64(a, b any) bool {
	fx := a.(float64) //nolint:forcetypeassert
	fy := b.(float64) //nolint:forcetypeassert

	if math.IsInf(fx, 1) && math.IsInf(fy, 1) {
		return true
	}

	if math.IsInf(fx, -1) && math.IsInf(fy, -1) {
		return true
	}

	const thr = float64(1e-10)
	return math.Abs(fx-fy) <= thr
}

func compareFloat32(a, b any) bool {
	fx := a.(float32) //nolint:forcetypeassert
	fy := b.(float32) //nolint:forcetypeassert

	if math.IsInf(float64(fx), 1) && math.IsInf(float64(fy), 1) {
		return true
	}

	if math.IsInf(float64(fx), -1) && math.IsInf(float64(fy), -1) {
		return true
	}

	const thr = float64(1e-7)
	return math.Abs(float64(fx-fy)) <= thr
}
