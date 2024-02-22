package cast

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/moukoublen/pick/internal/testingx"
)

type castTestExpectedResult[Output any] struct {
	expectedResult Output
	errorAssertFn  func(*testing.T, error)
	shouldRun      bool
}

// constructors for shortage.
func castTestExpectedResultFn[Output any]() func(result Output, errorAssertFn func(*testing.T, error)) castTestExpectedResult[Output] {
	return func(result Output, errorAssertFn func(*testing.T, error)) castTestExpectedResult[Output] {
		return castTestExpectedResult[Output]{expectedResult: result, errorAssertFn: errorAssertFn, shouldRun: true}
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

//nolint:maintidx
func TestCasterMatrix(t *testing.T) {
	t.Parallel()

	type stringAlias string

	// matrixExpectedResult constructor function aliases.
	expectByte := castTestExpectedResultFn[byte]()
	expectInt8 := castTestExpectedResultFn[int8]()
	expectInt16 := castTestExpectedResultFn[int16]()
	expectInt32 := castTestExpectedResultFn[int32]()
	expectInt64 := castTestExpectedResultFn[int64]()
	expectInt := castTestExpectedResultFn[int]()
	expectUInt8 := castTestExpectedResultFn[uint8]()
	expectUint16 := castTestExpectedResultFn[uint16]()
	expectUint32 := castTestExpectedResultFn[uint32]()
	expectUint64 := castTestExpectedResultFn[uint64]()
	expectUint := castTestExpectedResultFn[uint]()
	expectFloat32 := castTestExpectedResultFn[float32]()
	expectFloat64 := castTestExpectedResultFn[float64]()
	expectString := castTestExpectedResultFn[string]()
	expectBool := castTestExpectedResultFn[bool]()
	expectTime := castTestExpectedResultFn[time.Time]()
	expectDuration := castTestExpectedResultFn[time.Duration]()

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
		Bool    castTestExpectedResult[bool]
		Time    castTestExpectedResult[time.Time]
		Dur     castTestExpectedResult[time.Duration]
	}{
		{
			Input:   nil,
			Byte:    expectByte(0, nil),
			Int8:    expectInt8(0, nil),
			Int16:   expectInt16(0, nil),
			Int32:   expectInt32(0, nil),
			Int64:   expectInt64(0, nil),
			Int:     expectInt(0, nil),
			Uint8:   expectUInt8(0, nil),
			Uint16:  expectUint16(0, nil),
			Uint32:  expectUint32(0, nil),
			Uint64:  expectUint64(0, nil),
			Uint:    expectUint(0, nil),
			Float32: expectFloat32(0, nil),
			Float64: expectFloat64(0, nil),
			String:  expectString("", nil),
			Bool:    expectBool(false, nil),
			Time:    expectTime(time.Time{}, nil),
			Dur:     expectDuration(time.Duration(0), nil),
		},
		{
			Input:   int8(12),
			Byte:    expectByte(12, nil),
			Int8:    expectInt8(12, nil),
			Int16:   expectInt16(12, nil),
			Int32:   expectInt32(12, nil),
			Int64:   expectInt64(12, nil),
			Int:     expectInt(12, nil),
			Uint8:   expectUInt8(12, nil),
			Uint16:  expectUint16(12, nil),
			Uint32:  expectUint32(12, nil),
			Uint64:  expectUint64(12, nil),
			Uint:    expectUint(12, nil),
			Float32: expectFloat32(12, nil),
			Float64: expectFloat64(12, nil),
			String:  expectString("12", nil),
			Bool:    expectBool(true, nil),
			Time:    expectTime(time.Date(1970, time.January, 1, 0, 0, 12, 0, time.UTC), nil),
			Dur:     expectDuration(time.Duration(12_000_000), nil),
		},
		{
			Input:   int8(math.MaxInt8),
			Byte:    expectByte(127, nil),
			Int8:    expectInt8(127, nil),
			Int16:   expectInt16(127, nil),
			Int32:   expectInt32(127, nil),
			Int64:   expectInt64(127, nil),
			Int:     expectInt(127, nil),
			Uint8:   expectUInt8(127, nil),
			Uint16:  expectUint16(127, nil),
			Uint32:  expectUint32(127, nil),
			Uint64:  expectUint64(127, nil),
			Uint:    expectUint(127, nil),
			Float32: expectFloat32(127, nil),
			Float64: expectFloat64(127, nil),
			String:  expectString("127", nil),
			Bool:    expectBool(true, nil),
			Time:    expectTime(time.Date(1970, time.January, 1, 0, 2, 7, 0, time.UTC), nil),
			Dur:     expectDuration(time.Duration(127_000_000), nil),
		},
		{
			Input:   int8(math.MinInt8),
			Byte:    expectByte(0x80, expectOverFlowError),
			Int8:    expectInt8(-128, nil),
			Int16:   expectInt16(-128, nil),
			Int32:   expectInt32(-128, nil),
			Int64:   expectInt64(-128, nil),
			Int:     expectInt(-128, nil),
			Uint8:   expectUInt8(0x80, expectOverFlowError),
			Uint16:  expectUint16(0xff80, expectOverFlowError),
			Uint32:  expectUint32(0xffffff80, expectOverFlowError),
			Uint64:  expectUint64(0xffffffffffffff80, expectOverFlowError),
			Uint:    expectUint(0xffffffffffffff80, expectOverFlowError),
			Float32: expectFloat32(-128, nil),
			Float64: expectFloat64(-128, nil),
			String:  expectString("-128", nil),
			Bool:    expectBool(true, nil),
			Time:    expectTime(time.Date(1969, time.December, 31, 23, 57, 52, 0, time.UTC), nil),
			Dur:     expectDuration(time.Duration(-128_000_000), nil),
		},
		{
			Input:   int16(math.MaxInt16),
			Byte:    expectByte(255, expectOverFlowError),
			Int8:    expectInt8(-1, expectOverFlowError),
			Int16:   expectInt16(32767, nil),
			Int32:   expectInt32(32767, nil),
			Int64:   expectInt64(32767, nil),
			Int:     expectInt(32767, nil),
			Uint8:   expectUInt8(255, expectOverFlowError),
			Uint16:  expectUint16(32767, nil),
			Uint32:  expectUint32(32767, nil),
			Uint64:  expectUint64(32767, nil),
			Uint:    expectUint(32767, nil),
			Float32: expectFloat32(32767, nil),
			Float64: expectFloat64(32767, nil),
			String:  expectString("32767", nil),
			Bool:    expectBool(true, nil),
			Time:    expectTime(time.Date(1970, time.January, 1, 9, 6, 7, 0, time.UTC), nil),
			Dur:     expectDuration(time.Duration(32767_000_000), nil),
		},
		{
			Input:   int16(math.MinInt16),
			Byte:    expectByte(0, expectOverFlowError),
			Int8:    expectInt8(0, expectOverFlowError),
			Int16:   expectInt16(-32768, nil),
			Int32:   expectInt32(-32768, nil),
			Int64:   expectInt64(-32768, nil),
			Int:     expectInt(-32768, nil),
			Uint8:   expectUInt8(0, expectOverFlowError),
			Uint16:  expectUint16(0x8000, expectOverFlowError),
			Uint32:  expectUint32(0xffff8000, expectOverFlowError),
			Uint64:  expectUint64(0xffffffffffff8000, expectOverFlowError),
			Uint:    expectUint(0xffffffffffff8000, expectOverFlowError),
			Float32: expectFloat32(-32768, nil),
			Float64: expectFloat64(-32768, nil),
			String:  expectString("-32768", nil),
			Bool:    expectBool(true, nil),
			Time:    expectTime(time.Date(1969, time.December, 31, 14, 53, 52, 0, time.UTC), nil),
			Dur:     expectDuration(time.Duration(-32768_000_000), nil),
		},
		{
			Input:   int32(math.MaxInt32),
			Byte:    expectByte(255, expectOverFlowError),
			Int8:    expectInt8(-1, expectOverFlowError),
			Int16:   expectInt16(-1, expectOverFlowError),
			Int32:   expectInt32(2147483647, nil),
			Int64:   expectInt64(2147483647, nil),
			Int:     expectInt(2147483647, nil),
			Uint8:   expectUInt8(255, expectOverFlowError),
			Uint16:  expectUint16(0xffff, expectOverFlowError),
			Uint32:  expectUint32(2147483647, nil),
			Uint64:  expectUint64(2147483647, nil),
			Uint:    expectUint(2147483647, nil),
			Float32: expectFloat32(2147483647, nil),
			Float64: expectFloat64(2147483647, nil),
			String:  expectString("2147483647", nil),
			Bool:    expectBool(true, nil),
			Time:    expectTime(time.Date(2038, time.January, 19, 3, 14, 7, 0, time.UTC), nil),
			Dur:     expectDuration(time.Duration(2147483647_000_000), nil),
		},
		{
			Input:   int32(math.MinInt32),
			Byte:    expectByte(0, expectOverFlowError),
			Int8:    expectInt8(0, expectOverFlowError),
			Int16:   expectInt16(0, expectOverFlowError),
			Int32:   expectInt32(-2147483648, nil),
			Int64:   expectInt64(-2147483648, nil),
			Int:     expectInt(-2147483648, nil),
			Uint8:   expectUInt8(0, expectOverFlowError),
			Uint16:  expectUint16(0, expectOverFlowError),
			Uint32:  expectUint32(0x80000000, expectOverFlowError),
			Uint64:  expectUint64(0xffffffff80000000, expectOverFlowError),
			Uint:    expectUint(0xffffffff80000000, expectOverFlowError),
			Float32: expectFloat32(-2147483648, nil),
			Float64: expectFloat64(-2147483648, nil),
			String:  expectString("-2147483648", nil),
			Bool:    expectBool(true, nil),
			Time:    expectTime(time.Date(1901, time.December, 13, 20, 45, 52, 0, time.UTC), nil),
			Dur:     expectDuration(time.Duration(-2147483648_000_000), nil),
		},
		{
			Input:   int64(math.MaxInt64),
			Byte:    expectByte(255, expectOverFlowError),
			Int8:    expectInt8(-1, expectOverFlowError),
			Int16:   expectInt16(-1, expectOverFlowError),
			Int32:   expectInt32(-1, expectOverFlowError),
			Int64:   expectInt64(math.MaxInt64, nil),
			Int:     splitBasedOnArch(expectInt(-1, expectOverFlowError), expectInt(math.MaxInt64, nil)),
			Uint8:   expectUInt8(255, expectOverFlowError),
			Uint16:  expectUint16(0xffff, expectOverFlowError),
			Uint32:  expectUint32(0xffffffff, expectOverFlowError),
			Uint64:  expectUint64(0x7fffffffffffffff, nil),
			Uint:    splitBasedOnArch(expectUint(0xffffffff, expectOverFlowError), expectUint(0x7fffffffffffffff, nil)),
			Float32: expectFloat32(math.MaxInt64, nil),
			Float64: expectFloat64(math.MaxInt64, nil),
			String:  expectString("9223372036854775807", nil),
			Bool:    expectBool(true, nil),
			Time:    expectTime(time.Date(292277026596, time.December, 4, 15, 30, 7, 0, time.UTC), nil), // the largest int64 value does not have a corresponding time value.
			Dur:     expectDuration(time.Duration(-1000000), expectOverFlowError),
		},
		{
			Input:   int64(math.MinInt64),
			Byte:    expectByte(0, expectOverFlowError),
			Int8:    expectInt8(0, expectOverFlowError),
			Int16:   expectInt16(0, expectOverFlowError),
			Int32:   expectInt32(0, expectOverFlowError),
			Int64:   expectInt64(math.MinInt64, nil),
			Int:     splitBasedOnArch(expectInt(0, expectOverFlowError), expectInt(math.MinInt64, nil)),
			Uint8:   expectUInt8(0, expectOverFlowError),
			Uint16:  expectUint16(0, expectOverFlowError),
			Uint32:  expectUint32(0, expectOverFlowError),
			Uint64:  expectUint64(0x8000000000000000, expectOverFlowError),
			Uint:    splitBasedOnArch(expectUint(0, expectOverFlowError), expectUint(0x8000000000000000, expectOverFlowError)),
			Float32: expectFloat32(math.MinInt64, nil),
			Float64: expectFloat64(math.MinInt64, nil),
			String:  expectString("-9223372036854775808", nil),
			Bool:    expectBool(true, nil),
			Time:    expectTime(time.Date(292277026596, time.December, 4, 15, 30, 8, 0, time.UTC), nil), // the min int64 value does not have a corresponding time value.
			Dur:     expectDuration(time.Duration(0), expectOverFlowError),
		},
		{
			Input:   uint8(math.MaxUint8),
			Byte:    expectByte(math.MaxUint8, nil),
			Int8:    expectInt8(-1, expectOverFlowError),
			Int16:   expectInt16(math.MaxUint8, nil),
			Int32:   expectInt32(math.MaxUint8, nil),
			Int64:   expectInt64(math.MaxUint8, nil),
			Int:     expectInt(math.MaxUint8, nil),
			Uint8:   expectUInt8(math.MaxUint8, nil),
			Uint16:  expectUint16(math.MaxUint8, nil),
			Uint32:  expectUint32(math.MaxUint8, nil),
			Uint64:  expectUint64(math.MaxUint8, nil),
			Uint:    expectUint(math.MaxUint8, nil),
			Float32: expectFloat32(math.MaxUint8, nil),
			Float64: expectFloat64(math.MaxUint8, nil),
			String:  expectString("255", nil),
			Bool:    expectBool(true, nil),
			Time:    expectTime(time.Date(1970, time.January, 1, 0, 4, 15, 0, time.UTC), nil),
			Dur:     expectDuration(time.Duration(255_000_000), nil),
		},
		{
			Input:   uint16(math.MaxUint16),
			Byte:    expectByte(255, expectOverFlowError),
			Int8:    expectInt8(-1, expectOverFlowError),
			Int16:   expectInt16(-1, expectOverFlowError),
			Int32:   expectInt32(math.MaxUint16, nil),
			Int64:   expectInt64(math.MaxUint16, nil),
			Int:     expectInt(math.MaxUint16, nil),
			Uint8:   expectUInt8(255, expectOverFlowError),
			Uint16:  expectUint16(math.MaxUint16, nil),
			Uint32:  expectUint32(math.MaxUint16, nil),
			Uint64:  expectUint64(math.MaxUint16, nil),
			Uint:    expectUint(math.MaxUint16, nil),
			Float32: expectFloat32(math.MaxUint16, nil),
			Float64: expectFloat64(math.MaxUint16, nil),
			String:  expectString("65535", nil),
			Bool:    expectBool(true, nil),
			Time:    expectTime(time.Date(1970, time.January, 1, 18, 12, 15, 0, time.UTC), nil),
			Dur:     expectDuration(time.Duration(65535_000_000), nil),
		},
		{
			Input:   uint32(math.MaxUint32),
			Byte:    expectByte(255, expectOverFlowError),
			Int8:    expectInt8(-1, expectOverFlowError),
			Int16:   expectInt16(-1, expectOverFlowError),
			Int32:   expectInt32(-1, expectOverFlowError),
			Int64:   expectInt64(4294967295, nil),
			Int:     splitBasedOnArch(expectInt(0, expectOverFlowError), expectInt(math.MaxUint32, nil)),
			Uint8:   expectUInt8(255, expectOverFlowError),
			Uint16:  expectUint16(65535, expectOverFlowError),
			Uint32:  expectUint32(math.MaxUint32, nil),
			Uint64:  expectUint64(math.MaxUint32, nil),
			Uint:    expectUint(math.MaxUint32, nil),
			Float32: expectFloat32(math.MaxUint32, nil),
			Float64: expectFloat64(math.MaxUint32, nil),
			String:  expectString("4294967295", nil),
			Bool:    expectBool(true, nil),
			Time:    expectTime(time.Date(2106, time.February, 7, 6, 28, 15, 0, time.UTC), nil),
			Dur:     expectDuration(time.Duration(4294967295_000_000), nil),
		},
		{
			Input:   uint64(math.MaxUint64),
			Byte:    expectByte(255, expectOverFlowError),
			Int8:    expectInt8(-1, expectOverFlowError),
			Int16:   expectInt16(-1, expectOverFlowError),
			Int32:   expectInt32(-1, expectOverFlowError),
			Int64:   expectInt64(-1, expectOverFlowError),
			Int:     splitBasedOnArch(expectInt(-1, expectOverFlowError), expectInt(-1, expectOverFlowError)),
			Uint8:   expectUInt8(255, expectOverFlowError),
			Uint16:  expectUint16(65535, expectOverFlowError),
			Uint32:  expectUint32(math.MaxUint32, expectOverFlowError),
			Uint64:  expectUint64(math.MaxUint64, nil),
			Uint:    splitBasedOnArch(expectUint(0, expectOverFlowError), expectUint(math.MaxUint64, nil)),
			Float32: expectFloat32(math.MaxUint64, nil),
			Float64: expectFloat64(math.MaxUint64, nil),
			String:  expectString("18446744073709551615", nil),
			Bool:    expectBool(true, nil),
			Time:    expectTime(time.Date(1969, time.December, 31, 23, 59, 59, 0, time.UTC), expectOverFlowError), // max uint64 is not converted to valid date.
			Dur:     expectDuration(time.Duration(-1000000), expectOverFlowError),
		},
		{
			Input:   byte(12),
			Byte:    expectByte(12, nil),
			Int8:    expectInt8(12, nil),
			Int16:   expectInt16(12, nil),
			Int32:   expectInt32(12, nil),
			Int64:   expectInt64(12, nil),
			Int:     expectInt(12, nil),
			Uint8:   expectUInt8(12, nil),
			Uint16:  expectUint16(12, nil),
			Uint32:  expectUint32(12, nil),
			Uint64:  expectUint64(12, nil),
			Uint:    expectUint(12, nil),
			Float32: expectFloat32(12, nil),
			Float64: expectFloat64(12, nil),
			String:  expectString("12", nil),
			Bool:    expectBool(true, nil),
			Time:    expectTime(time.Date(1970, time.January, 1, 0, 0, 12, 0, time.UTC), nil),
			Dur:     expectDuration(time.Duration(12_000_000), nil),
		},
		{
			Input:   "123",
			Byte:    expectByte(0, expectInvalidType),
			Int8:    expectInt8(123, nil),
			Int16:   expectInt16(123, nil),
			Int32:   expectInt32(123, nil),
			Int64:   expectInt64(123, nil),
			Int:     expectInt(123, nil),
			Uint8:   expectUInt8(123, nil),
			Uint16:  expectUint16(123, nil),
			Uint32:  expectUint32(123, nil),
			Uint64:  expectUint64(123, nil),
			Uint:    expectUint(123, nil),
			Float32: expectFloat32(123, nil),
			Float64: expectFloat64(123, nil),
			String:  expectString("123", nil),
			Bool:    expectBool(false, expectMalformedSyntax),
			Time:    expectTime(time.Time{}, testingx.ExpectedErrorIsOfType(&time.ParseError{})),
			Dur:     expectDuration(time.Duration(0), testingx.ExpectedErrorStringContains("time: missing unit in duration")),
		},
		{
			Input:   []byte("123"),
			Byte:    expectByte(0, expectInvalidType),
			Int8:    expectInt8(123, nil),
			Int16:   expectInt16(123, nil),
			Int32:   expectInt32(123, nil),
			Int64:   expectInt64(123, nil),
			Int:     expectInt(123, nil),
			Uint8:   expectUInt8(123, nil),
			Uint16:  expectUint16(123, nil),
			Uint32:  expectUint32(123, nil),
			Uint64:  expectUint64(123, nil),
			Uint:    expectUint(123, nil),
			Float32: expectFloat32(123, nil),
			Float64: expectFloat64(123, nil),
			String:  expectString("123", nil),
			Bool:    expectBool(false, expectMalformedSyntax),
			Time:    expectTime(time.Time{}, testingx.ExpectedErrorIsOfType(&time.ParseError{})),
			Dur:     expectDuration(time.Duration(0), testingx.ExpectedErrorStringContains("time: missing unit in duration")),
		},
		{
			Input:   "123.321",
			Byte:    expectByte(0, expectInvalidType),
			Int8:    expectInt8(123, expectLostDecimals),
			Int16:   expectInt16(123, expectLostDecimals),
			Int32:   expectInt32(123, expectLostDecimals),
			Int64:   expectInt64(123, expectLostDecimals),
			Int:     expectInt(123, expectLostDecimals),
			Uint8:   expectUInt8(123, expectLostDecimals),
			Uint16:  expectUint16(123, expectLostDecimals),
			Uint32:  expectUint32(123, expectLostDecimals),
			Uint64:  expectUint64(123, expectLostDecimals),
			Uint:    expectUint(123, expectLostDecimals),
			Float32: expectFloat32(123.321, nil),
			Float64: expectFloat64(123.321, nil),
			String:  expectString("123.321", nil),
			Bool:    expectBool(false, expectMalformedSyntax),
			Time:    expectTime(time.Time{}, testingx.ExpectedErrorIsOfType(&time.ParseError{})),
			Dur:     expectDuration(time.Duration(0), testingx.ExpectedErrorStringContains("time: missing unit in duration")),
		},
		{
			Input:   stringAlias("23"),
			Byte:    expectByte(0, expectInvalidType),
			Int8:    expectInt8(23, nil),
			Int16:   expectInt16(23, nil),
			Int32:   expectInt32(23, nil),
			Int64:   expectInt64(23, nil),
			Int:     expectInt(23, nil),
			Uint8:   expectUInt8(23, nil),
			Uint16:  expectUint16(23, nil),
			Uint32:  expectUint32(23, nil),
			Uint64:  expectUint64(23, nil),
			Uint:    expectUint(23, nil),
			Float32: expectFloat32(23, nil),
			Float64: expectFloat64(23, nil),
			String:  expectString("23", nil),
			Bool:    expectBool(false, expectMalformedSyntax),
			Time:    expectTime(time.Time{}, testingx.ExpectedErrorIsOfType(&time.ParseError{})),
			Dur:     expectDuration(time.Duration(0), testingx.ExpectedErrorStringContains("time: missing unit in duration")),
		},
		{
			Input:   "just string",
			Byte:    expectByte(0, expectInvalidType),
			Int8:    expectInt8(0, expectMalformedSyntax),
			Int16:   expectInt16(0, expectMalformedSyntax),
			Int32:   expectInt32(0, expectMalformedSyntax),
			Int64:   expectInt64(0, expectMalformedSyntax),
			Int:     expectInt(0, expectMalformedSyntax),
			Uint8:   expectUInt8(0, expectMalformedSyntax),
			Uint16:  expectUint16(0, expectMalformedSyntax),
			Uint32:  expectUint32(0, expectMalformedSyntax),
			Uint64:  expectUint64(0, expectMalformedSyntax),
			Uint:    expectUint(0, expectMalformedSyntax),
			Float32: expectFloat32(0, expectMalformedSyntax),
			Float64: expectFloat64(0, expectMalformedSyntax),
			String:  expectString("just string", nil),
			Bool:    expectBool(false, expectMalformedSyntax),
			Time:    expectTime(time.Time{}, testingx.ExpectedErrorIsOfType(&time.ParseError{})),
			Dur:     expectDuration(time.Duration(0), testingx.ExpectedErrorStringContains("time: invalid duration")),
		},
		{
			Input:   []byte("byte slice"),
			Byte:    expectByte(0, expectInvalidType),
			Int8:    expectInt8(0, expectMalformedSyntax),
			Int16:   expectInt16(0, expectMalformedSyntax),
			Int32:   expectInt32(0, expectMalformedSyntax),
			Int64:   expectInt64(0, expectMalformedSyntax),
			Int:     expectInt(0, expectMalformedSyntax),
			Uint8:   expectUInt8(0, expectMalformedSyntax),
			Uint16:  expectUint16(0, expectMalformedSyntax),
			Uint32:  expectUint32(0, expectMalformedSyntax),
			Uint64:  expectUint64(0, expectMalformedSyntax),
			Uint:    expectUint(0, expectMalformedSyntax),
			Float32: expectFloat32(0, expectMalformedSyntax),
			Float64: expectFloat64(0, expectMalformedSyntax),
			String:  expectString("byte slice", nil),
			Bool:    expectBool(false, expectMalformedSyntax),
			Time:    expectTime(time.Time{}, testingx.ExpectedErrorIsOfType(&time.ParseError{})),
			Dur:     expectDuration(time.Duration(0), testingx.ExpectedErrorStringContains("time: invalid duration")),
		},
		{
			Input:   float32(123),
			Byte:    expectByte(123, nil),
			Int8:    expectInt8(123, nil),
			Int16:   expectInt16(123, nil),
			Int32:   expectInt32(123, nil),
			Int64:   expectInt64(123, nil),
			Int:     expectInt(123, nil),
			Uint8:   expectUInt8(123, nil),
			Uint16:  expectUint16(123, nil),
			Uint32:  expectUint32(123, nil),
			Uint64:  expectUint64(123, nil),
			Uint:    expectUint(123, nil),
			Float32: expectFloat32(123, nil),
			Float64: expectFloat64(123, nil),
			String:  expectString("123", nil),
			Bool:    expectBool(true, nil),
			Time:    expectTime(time.Date(1970, time.January, 1, 0, 2, 3, 0, time.UTC), nil),
			Dur:     expectDuration(time.Duration(123_000_000), nil),
		},
		{
			Input:   float64(123),
			Byte:    expectByte(123, nil),
			Int8:    expectInt8(123, nil),
			Int16:   expectInt16(123, nil),
			Int32:   expectInt32(123, nil),
			Int64:   expectInt64(123, nil),
			Int:     expectInt(123, nil),
			Uint8:   expectUInt8(123, nil),
			Uint16:  expectUint16(123, nil),
			Uint32:  expectUint32(123, nil),
			Uint64:  expectUint64(123, nil),
			Uint:    expectUint(123, nil),
			Float32: expectFloat32(123, nil),
			Float64: expectFloat64(123, nil),
			String:  expectString("123", nil),
			Time:    expectTime(time.Date(1970, time.January, 1, 0, 2, 3, 0, time.UTC), nil),
			Dur:     expectDuration(time.Duration(123_000_000), nil),
		},
		{
			Input:   float64(123.12),
			Byte:    expectByte(123, expectLostDecimals),
			Int8:    expectInt8(123, expectLostDecimals),
			Int16:   expectInt16(123, expectLostDecimals),
			Int32:   expectInt32(123, expectLostDecimals),
			Int64:   expectInt64(123, expectLostDecimals),
			Int:     expectInt(123, expectLostDecimals),
			Uint8:   expectUInt8(123, expectLostDecimals),
			Uint16:  expectUint16(123, expectLostDecimals),
			Uint32:  expectUint32(123, expectLostDecimals),
			Uint64:  expectUint64(123, expectLostDecimals),
			Uint:    expectUint(123, expectLostDecimals),
			Float32: expectFloat32(123.12, nil),
			Float64: expectFloat64(123.12, nil),
			String:  expectString("123.12", nil),
			Time:    expectTime(time.Date(1970, time.January, 1, 0, 2, 3, 0, time.UTC), expectLostDecimals),
			Dur:     expectDuration(time.Duration(123_000_000), expectLostDecimals),
		},
		{
			Input:   float64(math.MaxFloat64),
			Byte:    expectByte(0, expectOverFlowError),
			Int8:    expectInt8(0, expectOverFlowError),
			Int16:   expectInt16(0, expectOverFlowError),
			Int32:   expectInt32(0, expectOverFlowError),
			Int64:   expectInt64(-9223372036854775808, expectOverFlowError),
			Int:     splitBasedOnArch(expectInt(0, expectOverFlowError), expectInt(-9223372036854775808, expectOverFlowError)),
			Uint8:   expectUInt8(0, expectOverFlowError),
			Uint16:  expectUint16(0, expectOverFlowError),
			Uint32:  expectUint32(0, expectOverFlowError),
			Uint64:  expectUint64(9223372036854775808, expectOverFlowError),
			Uint:    splitBasedOnArch(expectUint(0, expectOverFlowError), expectUint(9223372036854775808, expectOverFlowError)),
			Float32: expectFloat32(float32(math.Inf(1)), expectOverFlowError),
			Float64: expectFloat64(math.MaxFloat64, nil),
			String:  expectString("1.7977E+308", nil),
			Bool:    expectBool(true, nil),
			Time:    expectTime(time.Date(292277026596, time.December, 4, 15, 30, 8, 0, time.UTC), expectOverFlowError),
			Dur:     expectDuration(time.Duration(0), expectOverFlowError),
		},
		{
			Input:   struct{}{},
			Byte:    expectByte(0, expectInvalidType),
			Int8:    expectInt8(0, expectInvalidType),
			Int16:   expectInt16(0, expectInvalidType),
			Int32:   expectInt32(0, expectInvalidType),
			Int64:   expectInt64(0, expectInvalidType),
			Int:     expectInt(0, expectInvalidType),
			Uint8:   expectUInt8(0, expectInvalidType),
			Uint16:  expectUint16(0, expectInvalidType),
			Uint32:  expectUint32(0, expectInvalidType),
			Uint64:  expectUint64(0, expectInvalidType),
			Uint:    expectUint(0, expectInvalidType),
			Float32: expectFloat32(0, expectInvalidType),
			Float64: expectFloat64(0, expectInvalidType),
			String:  expectString("", expectInvalidType),
			Bool:    expectBool(false, expectInvalidType),
			Time:    expectTime(time.Time{}, expectInvalidType),
			Dur:     expectDuration(time.Duration(0), expectInvalidType),
		},
		{
			Input:   json.RawMessage(`{"a":"b"}`),
			Byte:    expectByte(0, expectInvalidType),
			Int8:    expectInt8(0, expectInvalidType),
			Int16:   expectInt16(0, expectInvalidType),
			Int32:   expectInt32(0, expectInvalidType),
			Int64:   expectInt64(0, expectInvalidType),
			Int:     expectInt(0, expectInvalidType),
			Uint8:   expectUInt8(0, expectInvalidType),
			Uint16:  expectUint16(0, expectInvalidType),
			Uint32:  expectUint32(0, expectInvalidType),
			Uint64:  expectUint64(0, expectInvalidType),
			Uint:    expectUint(0, expectInvalidType),
			Float32: expectFloat32(0, expectInvalidType),
			Float64: expectFloat64(0, expectInvalidType),
			String:  expectString(`{"a":"b"}`, nil),
			Bool:    expectBool(false, expectInvalidType),
			Time:    expectTime(time.Time{}, expectInvalidType),
			Dur:     expectDuration(time.Duration(0), expectInvalidType),
		},
		{
			Input:   json.Number("123"),
			Byte:    expectByte(123, nil),
			Int8:    expectInt8(123, nil),
			Int16:   expectInt16(123, nil),
			Int32:   expectInt32(123, nil),
			Int64:   expectInt64(123, nil),
			Int:     expectInt(123, nil),
			Uint8:   expectUInt8(123, nil),
			Uint16:  expectUint16(123, nil),
			Uint32:  expectUint32(123, nil),
			Uint64:  expectUint64(123, nil),
			Uint:    expectUint(123, nil),
			Float32: expectFloat32(123, nil),
			Float64: expectFloat64(123, nil),
			String:  expectString("123", nil),
			Bool:    expectBool(true, nil),
			Time:    expectTime(time.Date(1970, time.January, 1, 0, 2, 3, 0, time.UTC), nil),
			Dur:     expectDuration(time.Duration(123_000_000), nil),
		},
		{
			Input:   json.Number("56782"),
			Byte:    expectByte(206, expectOverFlowError),
			Int8:    expectInt8(-50, expectOverFlowError),
			Int16:   expectInt16(-8754, expectOverFlowError),
			Int32:   expectInt32(56782, nil),
			Int64:   expectInt64(56782, nil),
			Int:     expectInt(56782, nil),
			Uint8:   expectUInt8(206, expectOverFlowError),
			Uint16:  expectUint16(56782, nil),
			Uint32:  expectUint32(56782, nil),
			Uint64:  expectUint64(56782, nil),
			Uint:    expectUint(56782, nil),
			Float32: expectFloat32(56782, nil),
			Float64: expectFloat64(56782, nil),
			String:  expectString("56782", nil),
			Bool:    expectBool(true, nil),
			Time:    expectTime(time.Date(1970, time.January, 1, 15, 46, 22, 0, time.UTC), nil),
			Dur:     expectDuration(time.Duration(56782_000_000), nil),
		},
		{
			Input:   "1.79769313486231570814527423731704356798070e+308",
			Byte:    expectByte(0, expectInvalidType),
			Int8:    expectInt8(0, expectOverFlowError),
			Int16:   expectInt16(0, expectOverFlowError),
			Int32:   expectInt32(0, expectOverFlowError),
			Int64:   expectInt64(-9223372036854775808, expectOverFlowError),
			Int:     splitBasedOnArch(expectInt(0, expectOverFlowError), expectInt(-9223372036854775808, expectOverFlowError)),
			Uint8:   expectUInt8(0, expectOverFlowError),
			Uint16:  expectUint16(0, expectOverFlowError),
			Uint32:  expectUint32(0, expectOverFlowError),
			Uint64:  expectUint64(9223372036854775808, expectOverFlowError),
			Uint:    splitBasedOnArch(expectUint(0, expectOverFlowError), expectUint(9223372036854775808, expectOverFlowError)),
			Float32: expectFloat32(float32(math.Inf(1)), expectOverFlowError),
			Float64: expectFloat64(math.MaxFloat64, nil),
			String:  expectString("1.79769313486231570814527423731704356798070e+308", nil),
			Bool:    expectBool(false, expectMalformedSyntax),
			Time:    expectTime(time.Time{}, testingx.ExpectedErrorIsOfType(&time.ParseError{})),
			Dur:     expectDuration(time.Duration(0), testingx.ExpectedErrorStringContains("time: unknown unit")),
		},
		{
			Input:   true,
			Byte:    expectByte(1, nil),
			Int8:    expectInt8(1, nil),
			Int16:   expectInt16(1, nil),
			Int32:   expectInt32(1, nil),
			Int64:   expectInt64(1, nil),
			Int:     expectInt(1, nil),
			Uint8:   expectUInt8(1, nil),
			Uint16:  expectUint16(1, nil),
			Uint32:  expectUint32(1, nil),
			Uint64:  expectUint64(1, nil),
			Uint:    expectUint(1, nil),
			Float32: expectFloat32(1, nil),
			Float64: expectFloat64(1, nil),
			String:  expectString("true", nil),
			Bool:    expectBool(true, nil),
			Time:    expectTime(time.Time{}, expectInvalidType),
			Dur:     expectDuration(time.Duration(0), expectInvalidType),
		},
		{
			Input:   false,
			Byte:    expectByte(0, nil),
			Int8:    expectInt8(0, nil),
			Int16:   expectInt16(0, nil),
			Int32:   expectInt32(0, nil),
			Int64:   expectInt64(0, nil),
			Int:     expectInt(0, nil),
			Uint8:   expectUInt8(0, nil),
			Uint16:  expectUint16(0, nil),
			Uint32:  expectUint32(0, nil),
			Uint64:  expectUint64(0, nil),
			Uint:    expectUint(0, nil),
			Float32: expectFloat32(0, nil),
			Float64: expectFloat64(0, nil),
			String:  expectString("false", nil),
			Bool:    expectBool(false, nil),
			Time:    expectTime(time.Time{}, expectInvalidType),
			Dur:     expectDuration(time.Duration(0), expectInvalidType),
		},
	}

	caster := NewCaster()
	for idx, tc := range testCases {
		tc := tc

		name := fmt.Sprintf("index(%d)__%T(%#v)", idx, tc.Input, tc.Input)
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
			t.Run("caster_bool", matrixSubTest[bool](tc.Input, caster.AsBool, tc.Bool))
			t.Run("caster_time", matrixSubTest[time.Time](tc.Input, caster.AsTime, tc.Time))
			t.Run("caster_duration", matrixSubTest[time.Duration](tc.Input, caster.AsDuration, tc.Dur))

			t.Run("caster_as_int8", matrixSubAs(caster, reflect.Int8, tc.Input, tc.Int8))
			t.Run("caster_as_int16", matrixSubAs(caster, reflect.Int16, tc.Input, tc.Int16))
			t.Run("caster_as_int32", matrixSubAs(caster, reflect.Int32, tc.Input, tc.Int32))
			t.Run("caster_as_int64", matrixSubAs(caster, reflect.Int64, tc.Input, tc.Int64))
			t.Run("caster_as_uint8", matrixSubAs(caster, reflect.Uint8, tc.Input, tc.Uint8))
			t.Run("caster_as_uint16", matrixSubAs(caster, reflect.Uint16, tc.Input, tc.Uint16))
			t.Run("caster_as_uint32", matrixSubAs(caster, reflect.Uint32, tc.Input, tc.Uint32))
			t.Run("caster_as_uint64", matrixSubAs(caster, reflect.Uint64, tc.Input, tc.Uint64))
			t.Run("caster_as_bool", matrixSubAs(caster, reflect.Bool, tc.Input, tc.Bool))
			t.Run("caster_as_float32", matrixSubAs(caster, reflect.Float32, tc.Input, tc.Float32))
			t.Run("caster_as_float64", matrixSubAs(caster, reflect.Float64, tc.Input, tc.Float64))
			t.Run("caster_as_string", matrixSubAs(caster, reflect.String, tc.Input, tc.String))
		})
	}
}

func matrixSubAs[Output any](caster Caster, asKind reflect.Kind, input any, subTestCase castTestExpectedResult[Output]) func(*testing.T) {
	return func(t *testing.T) {
		t.Helper()
		t.Parallel()
		if !subTestCase.shouldRun {
			t.SkipNow()
		}

		got, gotErr := caster.As(input, asKind)
		testingx.AssertError(t, subTestCase.errorAssertFn, gotErr)
		testingx.AssertEqual(t, got, subTestCase.expectedResult)
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
		testingx.AssertEqual(t, got, subTestCase.expectedResult)
	}
}

func BenchmarkCasterSlice(b *testing.B) {
	testCases := []any{
		[]any{"abc", "def"},
		[]string{"abc", "def"},
		[]any{1, 2, 3, 4},
		[]int32{1, 2, 3, 4},
		[]any{"1", "2", "3", "4"},
		[]string{"1", "2", "3", "4"},
	}

	c := NewCaster()

	b.Run("AsBoolSlice", casterSubBenchmarks(testCases, c.AsBoolSlice))
	b.Run("AsByteSlice", casterSubBenchmarks(testCases, c.AsByteSlice))
	b.Run("AsFloat32Slice", casterSubBenchmarks(testCases, c.AsFloat32Slice))
	b.Run("AsFloat64Slice", casterSubBenchmarks(testCases, c.AsFloat64Slice))
	b.Run("AsIntSlice", casterSubBenchmarks(testCases, c.AsIntSlice))
	b.Run("AsInt8Slice", casterSubBenchmarks(testCases, c.AsInt8Slice))
	b.Run("AsInt16Slice", casterSubBenchmarks(testCases, c.AsInt16Slice))
	b.Run("AsInt32Slice", casterSubBenchmarks(testCases, c.AsInt32Slice))
	b.Run("AsInt64Slice", casterSubBenchmarks(testCases, c.AsInt64Slice))
	b.Run("AsUintSlice", casterSubBenchmarks(testCases, c.AsUintSlice))
	b.Run("AsUint8Slice", casterSubBenchmarks(testCases, c.AsUint8Slice))
	b.Run("AsUint16Slice", casterSubBenchmarks(testCases, c.AsUint16Slice))
	b.Run("AsUint32Slice", casterSubBenchmarks(testCases, c.AsUint32Slice))
	b.Run("AsUint64Slice", casterSubBenchmarks(testCases, c.AsUint64Slice))
	b.Run("AsStringSlice", casterSubBenchmarks(testCases, c.AsStringSlice))
}

func casterSubBenchmarks[Output any](testCases []any, castFn func(any) (Output, error)) func(b *testing.B) {
	return func(b *testing.B) {
		b.Helper()
		for i, tc := range testCases {
			tc := tc
			name := fmt.Sprintf("%d %s", i, testingx.Format(tc))
			b.Run(name, matrixSubBenchmark(tc, castFn))
		}
	}
}

func matrixSubBenchmark[Output any](input any, castFn func(any) (Output, error)) func(b *testing.B) {
	return func(b *testing.B) {
		b.Helper()
		for i := 0; i < b.N; i++ {
			_, err := castFn(input)
			if err != nil {
				b.Skipf("skipped because of error %s", err.Error())
			}
		}
	}
}
