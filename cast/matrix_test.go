package cast

import (
	"encoding/json"
	"fmt"
	"math"
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
func newCastTestExpectedResultConstructor[Output any]() func(result Output, errorAssertFn func(*testing.T, error)) castTestExpectedResult[Output] {
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

	// matrixExpectedResult constructor function aliases (prefix `ex` from expected result)
	exByte := newCastTestExpectedResultConstructor[byte]()
	exInt8 := newCastTestExpectedResultConstructor[int8]()
	exInt16 := newCastTestExpectedResultConstructor[int16]()
	exInt32 := newCastTestExpectedResultConstructor[int32]()
	exInt64 := newCastTestExpectedResultConstructor[int64]()
	exInt := newCastTestExpectedResultConstructor[int]()
	exUInt8 := newCastTestExpectedResultConstructor[uint8]()
	exUint16 := newCastTestExpectedResultConstructor[uint16]()
	exUint32 := newCastTestExpectedResultConstructor[uint32]()
	exUint64 := newCastTestExpectedResultConstructor[uint64]()
	exUint := newCastTestExpectedResultConstructor[uint]()
	exFloat32 := newCastTestExpectedResultConstructor[float32]()
	exFloat64 := newCastTestExpectedResultConstructor[float64]()
	exString := newCastTestExpectedResultConstructor[string]()
	exBool := newCastTestExpectedResultConstructor[bool]()
	exTime := newCastTestExpectedResultConstructor[time.Time]()
	exDuration := newCastTestExpectedResultConstructor[time.Duration]()

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
			Bool:    exBool(false, nil),
			Time:    exTime(time.Time{}, nil),
			Dur:     exDuration(time.Duration(0), nil),
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
			Bool:    exBool(true, nil),
			Time:    exTime(time.Date(1970, time.January, 1, 0, 0, 12, 0, time.UTC), nil),
			Dur:     exDuration(time.Duration(12_000_000), nil),
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
			Bool:    exBool(true, nil),
			Time:    exTime(time.Date(1970, time.January, 1, 0, 2, 7, 0, time.UTC), nil),
			Dur:     exDuration(time.Duration(127_000_000), nil),
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
			Bool:    exBool(true, nil),
			Time:    exTime(time.Date(1969, time.December, 31, 23, 57, 52, 0, time.UTC), nil),
			Dur:     exDuration(time.Duration(-128_000_000), nil),
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
			Bool:    exBool(true, nil),
			Time:    exTime(time.Date(1970, time.January, 1, 9, 6, 7, 0, time.UTC), nil),
			Dur:     exDuration(time.Duration(32767_000_000), nil),
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
			Bool:    exBool(true, nil),
			Time:    exTime(time.Date(1969, time.December, 31, 14, 53, 52, 0, time.UTC), nil),
			Dur:     exDuration(time.Duration(-32768_000_000), nil),
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
			Bool:    exBool(true, nil),
			Time:    exTime(time.Date(2038, time.January, 19, 3, 14, 7, 0, time.UTC), nil),
			Dur:     exDuration(time.Duration(2147483647_000_000), nil),
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
			Bool:    exBool(true, nil),
			Time:    exTime(time.Date(1901, time.December, 13, 20, 45, 52, 0, time.UTC), nil),
			Dur:     exDuration(time.Duration(-2147483648_000_000), nil),
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
			Bool:    exBool(true, nil),
			Time:    exTime(time.Date(292277026596, time.December, 4, 15, 30, 7, 0, time.UTC), nil), // the largest int64 value does not have a corresponding time value.
			Dur:     exDuration(time.Duration(-1000000), expectOverFlowError),
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
			Bool:    exBool(true, nil),
			Time:    exTime(time.Date(292277026596, time.December, 4, 15, 30, 8, 0, time.UTC), nil), // the min int64 value does not have a corresponding time value.
			Dur:     exDuration(time.Duration(0), expectOverFlowError),
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
			Bool:    exBool(true, nil),
			Time:    exTime(time.Date(1970, time.January, 1, 0, 4, 15, 0, time.UTC), nil),
			Dur:     exDuration(time.Duration(255_000_000), nil),
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
			Bool:    exBool(true, nil),
			Time:    exTime(time.Date(1970, time.January, 1, 18, 12, 15, 0, time.UTC), nil),
			Dur:     exDuration(time.Duration(65535_000_000), nil),
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
			Bool:    exBool(true, nil),
			Time:    exTime(time.Date(2106, time.February, 7, 6, 28, 15, 0, time.UTC), nil),
			Dur:     exDuration(time.Duration(4294967295_000_000), nil),
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
			Bool:    exBool(true, nil),
			Time:    exTime(time.Date(1969, time.December, 31, 23, 59, 59, 0, time.UTC), expectOverFlowError), // max uint64 is not converted to valid date.
			Dur:     exDuration(time.Duration(-1000000), expectOverFlowError),
		},
		{
			Input:   byte(12),
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
			Bool:    exBool(true, nil),
			Time:    exTime(time.Date(1970, time.January, 1, 0, 0, 12, 0, time.UTC), nil),
			Dur:     exDuration(time.Duration(12_000_000), nil),
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
			Bool:    exBool(false, expectMalformedSyntax),
			Time:    exTime(time.Time{}, testingx.ExpectedErrorIsOfType(&time.ParseError{})),
			Dur:     exDuration(time.Duration(0), testingx.ExpectedErrorStringContains("time: missing unit in duration")),
		},
		{
			Input:   []byte("123"),
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
			Bool:    exBool(false, expectMalformedSyntax),
			Time:    exTime(time.Time{}, testingx.ExpectedErrorIsOfType(&time.ParseError{})),
			Dur:     exDuration(time.Duration(0), testingx.ExpectedErrorStringContains("time: missing unit in duration")),
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
			Bool:    exBool(false, expectMalformedSyntax),
			Time:    exTime(time.Time{}, testingx.ExpectedErrorIsOfType(&time.ParseError{})),
			Dur:     exDuration(time.Duration(0), testingx.ExpectedErrorStringContains("time: missing unit in duration")),
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
			Bool:    exBool(false, expectMalformedSyntax),
			Time:    exTime(time.Time{}, testingx.ExpectedErrorIsOfType(&time.ParseError{})),
			Dur:     exDuration(time.Duration(0), testingx.ExpectedErrorStringContains("time: missing unit in duration")),
		},
		{
			Input:   "just string",
			Byte:    exByte(0, expectInvalidType),
			Int8:    exInt8(0, expectMalformedSyntax),
			Int16:   exInt16(0, expectMalformedSyntax),
			Int32:   exInt32(0, expectMalformedSyntax),
			Int64:   exInt64(0, expectMalformedSyntax),
			Int:     exInt(0, expectMalformedSyntax),
			Uint8:   exUInt8(0, expectMalformedSyntax),
			Uint16:  exUint16(0, expectMalformedSyntax),
			Uint32:  exUint32(0, expectMalformedSyntax),
			Uint64:  exUint64(0, expectMalformedSyntax),
			Uint:    exUint(0, expectMalformedSyntax),
			Float32: exFloat32(0, expectMalformedSyntax),
			Float64: exFloat64(0, expectMalformedSyntax),
			String:  exString("just string", nil),
			Bool:    exBool(false, expectMalformedSyntax),
			Time:    exTime(time.Time{}, testingx.ExpectedErrorIsOfType(&time.ParseError{})),
			Dur:     exDuration(time.Duration(0), testingx.ExpectedErrorStringContains("time: invalid duration")),
		},
		{
			Input:   []byte("byte slice"),
			Byte:    exByte(0, expectInvalidType),
			Int8:    exInt8(0, expectMalformedSyntax),
			Int16:   exInt16(0, expectMalformedSyntax),
			Int32:   exInt32(0, expectMalformedSyntax),
			Int64:   exInt64(0, expectMalformedSyntax),
			Int:     exInt(0, expectMalformedSyntax),
			Uint8:   exUInt8(0, expectMalformedSyntax),
			Uint16:  exUint16(0, expectMalformedSyntax),
			Uint32:  exUint32(0, expectMalformedSyntax),
			Uint64:  exUint64(0, expectMalformedSyntax),
			Uint:    exUint(0, expectMalformedSyntax),
			Float32: exFloat32(0, expectMalformedSyntax),
			Float64: exFloat64(0, expectMalformedSyntax),
			String:  exString("byte slice", nil),
			Bool:    exBool(false, expectMalformedSyntax),
			Time:    exTime(time.Time{}, testingx.ExpectedErrorIsOfType(&time.ParseError{})),
			Dur:     exDuration(time.Duration(0), testingx.ExpectedErrorStringContains("time: invalid duration")),
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
			Bool:    exBool(true, nil),
			Time:    exTime(time.Date(1970, time.January, 1, 0, 2, 3, 0, time.UTC), nil),
			Dur:     exDuration(time.Duration(123_000_000), nil),
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
			Time:    exTime(time.Date(1970, time.January, 1, 0, 2, 3, 0, time.UTC), nil),
			Dur:     exDuration(time.Duration(123_000_000), nil),
		},
		{
			Input:   float64(123.12),
			Byte:    exByte(123, expectLostDecimals),
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
			Float32: exFloat32(123.12, nil),
			Float64: exFloat64(123.12, nil),
			String:  exString("123.12", nil),
			Time:    exTime(time.Date(1970, time.January, 1, 0, 2, 3, 0, time.UTC), expectLostDecimals),
			Dur:     exDuration(time.Duration(123_000_000), expectLostDecimals),
		},
		{
			Input:   float64(math.MaxFloat64),
			Byte:    exByte(0, expectOverFlowError),
			Int8:    exInt8(0, expectOverFlowError),
			Int16:   exInt16(0, expectOverFlowError),
			Int32:   exInt32(0, expectOverFlowError),
			Int64:   exInt64(-9223372036854775808, expectOverFlowError),
			Int:     splitBasedOnArch(exInt(0, expectOverFlowError), exInt(-9223372036854775808, expectOverFlowError)),
			Uint8:   exUInt8(0, expectOverFlowError),
			Uint16:  exUint16(0, expectOverFlowError),
			Uint32:  exUint32(0, expectOverFlowError),
			Uint64:  exUint64(9223372036854775808, expectOverFlowError),
			Uint:    splitBasedOnArch(exUint(0, expectOverFlowError), exUint(9223372036854775808, expectOverFlowError)),
			Float32: exFloat32(float32(math.Inf(1)), expectOverFlowError),
			Float64: exFloat64(math.MaxFloat64, nil),
			String:  exString("1.7977E+308", nil),
			Bool:    exBool(true, nil),
			Time:    exTime(time.Date(292277026596, time.December, 4, 15, 30, 8, 0, time.UTC), expectOverFlowError),
			Dur:     exDuration(time.Duration(0), expectOverFlowError),
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
			Bool:    exBool(false, expectInvalidType),
			Time:    exTime(time.Time{}, expectInvalidType),
			Dur:     exDuration(time.Duration(0), expectInvalidType),
		},
		{
			Input:   json.RawMessage(`{"a":"b"}`),
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
			String:  exString(`{"a":"b"}`, nil),
			Bool:    exBool(false, expectInvalidType),
			Time:    exTime(time.Time{}, expectInvalidType),
			Dur:     exDuration(time.Duration(0), expectInvalidType),
		},
		{
			Input:   json.Number("123"),
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
			Bool:    exBool(true, nil),
			Time:    exTime(time.Date(1970, time.January, 1, 0, 2, 3, 0, time.UTC), nil),
			Dur:     exDuration(time.Duration(123_000_000), nil),
		},
		{
			Input:   json.Number("56782"),
			Byte:    exByte(206, expectOverFlowError),
			Int8:    exInt8(-50, expectOverFlowError),
			Int16:   exInt16(-8754, expectOverFlowError),
			Int32:   exInt32(56782, nil),
			Int64:   exInt64(56782, nil),
			Int:     exInt(56782, nil),
			Uint8:   exUInt8(206, expectOverFlowError),
			Uint16:  exUint16(56782, nil),
			Uint32:  exUint32(56782, nil),
			Uint64:  exUint64(56782, nil),
			Uint:    exUint(56782, nil),
			Float32: exFloat32(56782, nil),
			Float64: exFloat64(56782, nil),
			String:  exString("56782", nil),
			Bool:    exBool(true, nil),
			Time:    exTime(time.Date(1970, time.January, 1, 15, 46, 22, 0, time.UTC), nil),
			Dur:     exDuration(time.Duration(56782_000_000), nil),
		},
		{
			Input:   "1.79769313486231570814527423731704356798070e+308",
			Byte:    exByte(0, expectInvalidType),
			Int8:    exInt8(0, expectOverFlowError),
			Int16:   exInt16(0, expectOverFlowError),
			Int32:   exInt32(0, expectOverFlowError),
			Int64:   exInt64(-9223372036854775808, expectOverFlowError),
			Int:     splitBasedOnArch(exInt(0, expectOverFlowError), exInt(-9223372036854775808, expectOverFlowError)),
			Uint8:   exUInt8(0, expectOverFlowError),
			Uint16:  exUint16(0, expectOverFlowError),
			Uint32:  exUint32(0, expectOverFlowError),
			Uint64:  exUint64(9223372036854775808, expectOverFlowError),
			Uint:    splitBasedOnArch(exUint(0, expectOverFlowError), exUint(9223372036854775808, expectOverFlowError)),
			Float32: exFloat32(float32(math.Inf(1)), expectOverFlowError),
			Float64: exFloat64(math.MaxFloat64, nil),
			String:  exString("1.79769313486231570814527423731704356798070e+308", nil),
			Bool:    exBool(false, expectMalformedSyntax),
			Time:    exTime(time.Time{}, testingx.ExpectedErrorIsOfType(&time.ParseError{})),
			Dur:     exDuration(time.Duration(0), testingx.ExpectedErrorStringContains("time: unknown unit")),
		},
		{
			Input:   true,
			Byte:    exByte(1, nil),
			Int8:    exInt8(1, nil),
			Int16:   exInt16(1, nil),
			Int32:   exInt32(1, nil),
			Int64:   exInt64(1, nil),
			Int:     exInt(1, nil),
			Uint8:   exUInt8(1, nil),
			Uint16:  exUint16(1, nil),
			Uint32:  exUint32(1, nil),
			Uint64:  exUint64(1, nil),
			Uint:    exUint(1, nil),
			Float32: exFloat32(1, nil),
			Float64: exFloat64(1, nil),
			String:  exString("true", nil),
			Bool:    exBool(true, nil),
			Time:    exTime(time.Time{}, expectInvalidType),
			Dur:     exDuration(time.Duration(0), expectInvalidType),
		},
		{
			Input:   false,
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
			String:  exString("false", nil),
			Bool:    exBool(false, nil),
			Time:    exTime(time.Time{}, expectInvalidType),
			Dur:     exDuration(time.Duration(0), expectInvalidType),
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
