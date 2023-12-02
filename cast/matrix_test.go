package cast

import (
	"fmt"
	"math"
	"reflect"
	"runtime"
	"testing"

	"github.com/moukoublen/pick/internal/testingx"
)

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

	type stringAlias string

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
		{
			Input:   uint8(math.MaxUint8),
			Byte:    exByte(math.MaxUint8, nil),
			Int8:    exInt8(-1, expectOverFlowError),
			Int16:   exInt16(math.MaxUint8, nil),
			Int32:   exInt32(math.MaxUint8, nil),
			Int64:   exInt64(math.MaxUint8, nil),
			Int:     exInt(math.MaxUint8, nil),
			Uint8:   exUInt8(math.MaxUint8, nil),
			Uint16:  exUint16(math.MaxUint8, nil),
			Uint32:  exUint32(math.MaxUint8, nil),
			Uint64:  exUint64(math.MaxUint8, nil),
			Uint:    exUint(math.MaxUint8, nil),
			Float32: exFloat32(math.MaxUint8, nil),
			Float64: exFloat64(math.MaxUint8, nil),
			String:  exString("255", nil),
		},
		{
			Input:   uint16(math.MaxUint16),
			Byte:    exByte(255, expectOverFlowError),
			Int8:    exInt8(-1, expectOverFlowError),
			Int16:   exInt16(-1, expectOverFlowError),
			Int32:   exInt32(math.MaxUint16, nil),
			Int64:   exInt64(math.MaxUint16, nil),
			Int:     exInt(math.MaxUint16, nil),
			Uint8:   exUInt8(255, expectOverFlowError),
			Uint16:  exUint16(math.MaxUint16, nil),
			Uint32:  exUint32(math.MaxUint16, nil),
			Uint64:  exUint64(math.MaxUint16, nil),
			Uint:    exUint(math.MaxUint16, nil),
			Float32: exFloat32(math.MaxUint16, nil),
			Float64: exFloat64(math.MaxUint16, nil),
			String:  exString("65535", nil),
		},
		{
			Input:   uint32(math.MaxUint32),
			Byte:    exByte(255, expectOverFlowError),
			Int8:    exInt8(-1, expectOverFlowError),
			Int16:   exInt16(-1, expectOverFlowError),
			Int32:   exInt32(-1, expectOverFlowError),
			Int64:   exInt64(4294967295, nil),
			Int:     splitBasedOnArch(exInt(0, expectOverFlowError), exInt(math.MaxUint32, nil)),
			Uint8:   exUInt8(255, expectOverFlowError),
			Uint16:  exUint16(65535, expectOverFlowError),
			Uint32:  exUint32(math.MaxUint32, nil),
			Uint64:  exUint64(math.MaxUint32, nil),
			Uint:    exUint(math.MaxUint32, nil),
			Float32: exFloat32(math.MaxUint32, nil),
			Float64: exFloat64(math.MaxUint32, nil),
			String:  exString("4294967295", nil),
		},
		{
			Input:   uint64(math.MaxUint64),
			Byte:    exByte(255, expectOverFlowError),
			Int8:    exInt8(-1, expectOverFlowError),
			Int16:   exInt16(-1, expectOverFlowError),
			Int32:   exInt32(-1, expectOverFlowError),
			Int64:   exInt64(-1, expectOverFlowError),
			Int:     splitBasedOnArch(exInt(-1, expectOverFlowError), exInt(-1, expectOverFlowError)),
			Uint8:   exUInt8(255, expectOverFlowError),
			Uint16:  exUint16(65535, expectOverFlowError),
			Uint32:  exUint32(math.MaxUint32, expectOverFlowError),
			Uint64:  exUint64(math.MaxUint64, nil),
			Uint:    splitBasedOnArch(exUint(0, expectOverFlowError), exUint(math.MaxUint64, nil)),
			Float32: exFloat32(math.MaxUint64, nil),
			Float64: exFloat64(math.MaxUint64, nil),
			String:  exString("18446744073709551615", nil),
		},
		{
			Input:   "123",
			Byte:    exByte(0, expectInvalidType),
			Int8:    exInt8(123, nil),
			Int16:   exInt16(123, nil),
			Int32:   exInt32(123, nil),
			Int64:   exInt64(123, nil),
			Int:     exInt(123, nil),
			Uint8:   exUInt8(123, nil),
			Uint16:  exUint16(123, nil),
			Uint32:  exUint32(123, nil),
			Uint64:  exUint64(123, nil),
			Uint:    exUint(123, nil),
			Float32: exFloat32(123, nil),
			Float64: exFloat64(123, nil),
			String:  exString("123", nil),
		},
		{
			Input:   "123.321",
			Byte:    exByte(0, expectInvalidType),
			Int8:    exInt8(123, expectLostDecimals),
			Int16:   exInt16(123, expectLostDecimals),
			Int32:   exInt32(123, expectLostDecimals),
			Int64:   exInt64(123, expectLostDecimals),
			Int:     exInt(123, expectLostDecimals),
			Uint8:   exUInt8(123, expectLostDecimals),
			Uint16:  exUint16(123, expectLostDecimals),
			Uint32:  exUint32(123, expectLostDecimals),
			Uint64:  exUint64(123, expectLostDecimals),
			Uint:    exUint(123, expectLostDecimals),
			Float32: exFloat32(123.321, nil),
			Float64: exFloat64(123.321, nil),
			String:  exString("123.321", nil),
		},
		{
			Input:   stringAlias("23"),
			Byte:    exByte(0, expectInvalidType),
			Int8:    exInt8(23, nil),
			Int16:   exInt16(23, nil),
			Int32:   exInt32(23, nil),
			Int64:   exInt64(23, nil),
			Int:     exInt(23, nil),
			Uint8:   exUInt8(23, nil),
			Uint16:  exUint16(23, nil),
			Uint32:  exUint32(23, nil),
			Uint64:  exUint64(23, nil),
			Uint:    exUint(23, nil),
			Float32: exFloat32(23, nil),
			Float64: exFloat64(23, nil),
			String:  exString("23", nil),
		},
		{
			Input:   float32(123),
			Byte:    exByte(123, nil),
			Int8:    exInt8(123, nil),
			Int16:   exInt16(123, nil),
			Int32:   exInt32(123, nil),
			Int64:   exInt64(123, nil),
			Int:     exInt(123, nil),
			Uint8:   exUInt8(123, nil),
			Uint16:  exUint16(123, nil),
			Uint32:  exUint32(123, nil),
			Uint64:  exUint64(123, nil),
			Uint:    exUint(123, nil),
			Float32: exFloat32(123, nil),
			Float64: exFloat64(123, nil),
			String:  exString("123", nil),
		},
		{
			Input:   float64(123),
			Byte:    exByte(123, nil),
			Int8:    exInt8(123, nil),
			Int16:   exInt16(123, nil),
			Int32:   exInt32(123, nil),
			Int64:   exInt64(123, nil),
			Int:     exInt(123, nil),
			Uint8:   exUInt8(123, nil),
			Uint16:  exUint16(123, nil),
			Uint32:  exUint32(123, nil),
			Uint64:  exUint64(123, nil),
			Uint:    exUint(123, nil),
			Float32: exFloat32(123, nil),
			Float64: exFloat64(123, nil),
			String:  exString("123", nil),
		},
		{
			Input:   struct{}{},
			Byte:    exByte(0, expectInvalidType),
			Int8:    exInt8(0, expectInvalidType),
			Int16:   exInt16(0, expectInvalidType),
			Int32:   exInt32(0, expectInvalidType),
			Int64:   exInt64(0, expectInvalidType),
			Int:     exInt(0, expectInvalidType),
			Uint8:   exUInt8(0, expectInvalidType),
			Uint16:  exUint16(0, expectInvalidType),
			Uint32:  exUint32(0, expectInvalidType),
			Uint64:  exUint64(0, expectInvalidType),
			Uint:    exUint(0, expectInvalidType),
			Float32: exFloat32(0, expectInvalidType),
			Float64: exFloat64(0, expectInvalidType),
			String:  exString("", expectInvalidType),
		},
	}

	caster := NewCaster()
	for idx, tc := range testCases {
		tc := tc

		name := fmt.Sprintf("index[%d]__inputType[%T]__inputValue[%#v]", idx, tc.Input, tc.Input)
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
