package pick

import (
	"encoding/json"
	"fmt"
	"math"
	"runtime"
	"testing"
	"time"

	"github.com/moukoublen/pick/internal/tst"
)

type ConverterTester interface {
	Test(t *testing.T)
	SetInput(i any)
}

func matrixTestConstructorFn[Output any](c DefaultConverter) func(expected Output, errorAsserter tst.ErrorAsserter) *converterTestCase[Output] {
	return func(expected Output, errorAsserter tst.ErrorAsserter) *converterTestCase[Output] {
		return &converterTestCase[Output]{
			Converter:                c,
			Input:                    nil,
			Expected:                 expected,
			ErrorAsserter:            errorAsserter,
			OverwriteDirectConvertFn: nil,
			OmitConvertByDirectFn:    false,
			OmitConvertByKind:        false,
			OmitConvertByType:        false,
		}
	}
}

func splitBasedOnArch[Output any](for32bit, for64bit *converterTestCase[Output]) *converterTestCase[Output] {
	switch runtime.GOARCH {
	case "arm", "386":
		return for32bit
	default:
		return for64bit
	}
}

//nolint:maintidx
func TestConverterMatrix(t *testing.T) {
	converter := NewDefaultConverter()

	type stringAlias string

	// matrixExpectedResult constructor function aliases.
	expectByte := func(expected byte, errorAssertFn tst.ErrorAsserter) *converterTestCase[byte] {
		return &converterTestCase[byte]{
			Converter:                converter,
			Input:                    nil,
			Expected:                 expected,
			ErrorAsserter:            errorAssertFn,
			OverwriteDirectConvertFn: converter.AsByte,
			OmitConvertByDirectFn:    false,
			OmitConvertByKind:        true,
			OmitConvertByType:        true,
		}
	}
	expectInt8 := matrixTestConstructorFn[int8](converter)
	expectInt16 := matrixTestConstructorFn[int16](converter)
	expectInt32 := matrixTestConstructorFn[int32](converter)
	expectInt64 := matrixTestConstructorFn[int64](converter)
	expectInt := matrixTestConstructorFn[int](converter)
	expectUInt8 := func(expected uint8, errorAssertFn tst.ErrorAsserter) *converterTestCase[uint8] {
		return &converterTestCase[uint8]{
			Converter:                converter,
			Input:                    nil,
			Expected:                 expected,
			ErrorAsserter:            errorAssertFn,
			OverwriteDirectConvertFn: converter.AsUint8,
			OmitConvertByDirectFn:    false,
			OmitConvertByKind:        false,
			OmitConvertByType:        false,
		}
	}
	expectUint16 := matrixTestConstructorFn[uint16](converter)
	expectUint32 := matrixTestConstructorFn[uint32](converter)
	expectUint64 := matrixTestConstructorFn[uint64](converter)
	expectUint := matrixTestConstructorFn[uint](converter)
	expectFloat32 := matrixTestConstructorFn[float32](converter)
	expectFloat64 := matrixTestConstructorFn[float64](converter)
	expectString := matrixTestConstructorFn[string](converter)
	expectBool := matrixTestConstructorFn[bool](converter)
	expectTime := matrixTestConstructorFn[time.Time](converter)
	expectDuration := matrixTestConstructorFn[time.Duration](converter)

	testCases := map[string]struct {
		Input     any
		Asserters []ConverterTester
	}{
		"tc#000": {
			Input: nil,
			Asserters: []ConverterTester{
				expectByte(0, nil),
				expectInt8(0, nil),
				expectInt16(0, nil),
				expectInt32(0, nil),
				expectInt64(0, nil),
				expectInt(0, nil),
				expectUInt8(0, nil),
				expectUint16(0, nil),
				expectUint32(0, nil),
				expectUint64(0, nil),
				expectUint(0, nil),
				expectFloat32(0, nil),
				expectFloat64(0, nil),
				expectString("", nil),
				expectBool(false, nil),
				expectTime(time.Time{}, nil),
				expectDuration(time.Duration(0), nil),
			},
		},
		"tc#001": {
			Input: int8(12),
			Asserters: []ConverterTester{
				expectByte(12, nil),
				expectInt8(12, nil),
				expectInt16(12, nil),
				expectInt32(12, nil),
				expectInt64(12, nil),
				expectInt(12, nil),
				expectUInt8(12, nil),
				expectUint16(12, nil),
				expectUint32(12, nil),
				expectUint64(12, nil),
				expectUint(12, nil),
				expectFloat32(12, nil),
				expectFloat64(12, nil),
				expectString("12", nil),
				expectBool(true, nil),
				expectTime(time.Date(1970, time.January, 1, 0, 0, 12, 0, time.UTC), nil),
				expectDuration(time.Duration(12), nil),
			},
		},
		"tc#002": {
			Input: int8(math.MaxInt8),
			Asserters: []ConverterTester{
				expectByte(127, nil),
				expectInt8(127, nil),
				expectInt16(127, nil),
				expectInt32(127, nil),
				expectInt64(127, nil),
				expectInt(127, nil),
				expectUInt8(127, nil),
				expectUint16(127, nil),
				expectUint32(127, nil),
				expectUint64(127, nil),
				expectUint(127, nil),
				expectFloat32(127, nil),
				expectFloat64(127, nil),
				expectString("127", nil),
				expectBool(true, nil),
				expectTime(time.Date(1970, time.January, 1, 0, 2, 7, 0, time.UTC), nil),
				expectDuration(time.Duration(127), nil),
			},
		},
		"tc#003": {
			Input: int8(math.MinInt8),
			Asserters: []ConverterTester{
				expectByte(0x80, expectOverFlowError),
				expectInt8(-128, nil),
				expectInt16(-128, nil),
				expectInt32(-128, nil),
				expectInt64(-128, nil),
				expectInt(-128, nil),
				expectUInt8(0x80, expectOverFlowError),
				expectUint16(0xff80, expectOverFlowError),
				expectUint32(0xffffff80, expectOverFlowError),
				expectUint64(0xffffffffffffff80, expectOverFlowError),
				expectUint(0xffffffffffffff80, expectOverFlowError),
				expectFloat32(-128, nil),
				expectFloat64(-128, nil),
				expectString("-128", nil),
				expectBool(true, nil),
				expectTime(time.Date(1969, time.December, 31, 23, 57, 52, 0, time.UTC), nil),
				expectDuration(time.Duration(-128), nil),
			},
		},
		"tc#004": {
			Input: int16(math.MaxInt16),
			Asserters: []ConverterTester{
				expectByte(255, expectOverFlowError),
				expectInt8(-1, expectOverFlowError),
				expectInt16(32767, nil),
				expectInt32(32767, nil),
				expectInt64(32767, nil),
				expectInt(32767, nil),
				expectUInt8(255, expectOverFlowError),
				expectUint16(32767, nil),
				expectUint32(32767, nil),
				expectUint64(32767, nil),
				expectUint(32767, nil),
				expectFloat32(32767, nil),
				expectFloat64(32767, nil),
				expectString("32767", nil),
				expectBool(true, nil),
				expectTime(time.Date(1970, time.January, 1, 9, 6, 7, 0, time.UTC), nil),
				expectDuration(time.Duration(32767), nil),
			},
		},
		"tc#005": {
			Input: int16(math.MinInt16),
			Asserters: []ConverterTester{
				expectByte(0, expectOverFlowError),
				expectInt8(0, expectOverFlowError),
				expectInt16(-32768, nil),
				expectInt32(-32768, nil),
				expectInt64(-32768, nil),
				expectInt(-32768, nil),
				expectUInt8(0, expectOverFlowError),
				expectUint16(0x8000, expectOverFlowError),
				expectUint32(0xffff8000, expectOverFlowError),
				expectUint64(0xffffffffffff8000, expectOverFlowError),
				expectUint(0xffffffffffff8000, expectOverFlowError),
				expectFloat32(-32768, nil),
				expectFloat64(-32768, nil),
				expectString("-32768", nil),
				expectBool(true, nil),
				expectTime(time.Date(1969, time.December, 31, 14, 53, 52, 0, time.UTC), nil),
				expectDuration(time.Duration(-32768), nil),
			},
		},
		"tc#006": {
			Input: int32(math.MaxInt32),
			Asserters: []ConverterTester{
				expectByte(255, expectOverFlowError),
				expectInt8(-1, expectOverFlowError),
				expectInt16(-1, expectOverFlowError),
				expectInt32(2147483647, nil),
				expectInt64(2147483647, nil),
				expectInt(2147483647, nil),
				expectUInt8(255, expectOverFlowError),
				expectUint16(0xffff, expectOverFlowError),
				expectUint32(2147483647, nil),
				expectUint64(2147483647, nil),
				expectUint(2147483647, nil),
				expectFloat32(2147483647, nil),
				expectFloat64(2147483647, nil),
				expectString("2147483647", nil),
				expectBool(true, nil),
				expectTime(time.Date(2038, time.January, 19, 3, 14, 7, 0, time.UTC), nil),
				expectDuration(time.Duration(2147483647), nil),
			},
		},
		"tc#007": {
			Input: int32(math.MinInt32),
			Asserters: []ConverterTester{
				expectByte(0, expectOverFlowError),
				expectInt8(0, expectOverFlowError),
				expectInt16(0, expectOverFlowError),
				expectInt32(-2147483648, nil),
				expectInt64(-2147483648, nil),
				expectInt(-2147483648, nil),
				expectUInt8(0, expectOverFlowError),
				expectUint16(0, expectOverFlowError),
				expectUint32(0x80000000, expectOverFlowError),
				expectUint64(0xffffffff80000000, expectOverFlowError),
				expectUint(0xffffffff80000000, expectOverFlowError),
				expectFloat32(-2147483648, nil),
				expectFloat64(-2147483648, nil),
				expectString("-2147483648", nil),
				expectBool(true, nil),
				expectTime(time.Date(1901, time.December, 13, 20, 45, 52, 0, time.UTC), nil),
				expectDuration(time.Duration(-2147483648), nil),
			},
		},
		"tc#008": {
			Input: int64(math.MaxInt64),
			Asserters: []ConverterTester{
				expectByte(255, expectOverFlowError),
				expectInt8(-1, expectOverFlowError),
				expectInt16(-1, expectOverFlowError),
				expectInt32(-1, expectOverFlowError),
				expectInt64(math.MaxInt64, nil),
				splitBasedOnArch(expectInt(-1, expectOverFlowError), expectInt(math.MaxInt64, nil)),
				expectUInt8(255, expectOverFlowError),
				expectUint16(0xffff, expectOverFlowError),
				expectUint32(0xffffffff, expectOverFlowError),
				expectUint64(0x7fffffffffffffff, nil),
				splitBasedOnArch(expectUint(0xffffffff, expectOverFlowError), expectUint(0x7fffffffffffffff, nil)),
				expectFloat32(math.MaxInt64, nil),
				expectFloat64(math.MaxInt64, nil),
				expectString("9223372036854775807", nil),
				expectBool(true, nil),
				expectTime(time.Date(292277026596, time.December, 4, 15, 30, 7, 0, time.UTC), nil), // the largest int64 value does not have a corresponding time value.
				expectDuration(time.Duration(math.MaxInt64), nil),
			},
		},
		"tc#009": {
			Input: int64(math.MinInt64),
			Asserters: []ConverterTester{
				expectByte(0, expectOverFlowError),
				expectInt8(0, expectOverFlowError),
				expectInt16(0, expectOverFlowError),
				expectInt32(0, expectOverFlowError),
				expectInt64(math.MinInt64, nil),
				splitBasedOnArch(expectInt(0, expectOverFlowError), expectInt(math.MinInt64, nil)),
				expectUInt8(0, expectOverFlowError),
				expectUint16(0, expectOverFlowError),
				expectUint32(0, expectOverFlowError),
				expectUint64(0x8000000000000000, expectOverFlowError),
				splitBasedOnArch(expectUint(0, expectOverFlowError), expectUint(0x8000000000000000, expectOverFlowError)),
				expectFloat32(math.MinInt64, nil),
				expectFloat64(math.MinInt64, nil),
				expectString("-9223372036854775808", nil),
				expectBool(true, nil),
				expectTime(time.Date(292277026596, time.December, 4, 15, 30, 8, 0, time.UTC), nil), // the min int64 value does not have a corresponding time value.
				expectDuration(time.Duration(math.MinInt64), nil),
			},
		},
		"tc#010": {
			Input: uint8(math.MaxUint8),
			Asserters: []ConverterTester{
				expectByte(math.MaxUint8, nil),
				expectInt8(-1, expectOverFlowError),
				expectInt16(math.MaxUint8, nil),
				expectInt32(math.MaxUint8, nil),
				expectInt64(math.MaxUint8, nil),
				expectInt(math.MaxUint8, nil),
				expectUInt8(math.MaxUint8, nil),
				expectUint16(math.MaxUint8, nil),
				expectUint32(math.MaxUint8, nil),
				expectUint64(math.MaxUint8, nil),
				expectUint(math.MaxUint8, nil),
				expectFloat32(math.MaxUint8, nil),
				expectFloat64(math.MaxUint8, nil),
				expectString("255", nil),
				expectBool(true, nil),
				expectTime(time.Date(1970, time.January, 1, 0, 4, 15, 0, time.UTC), nil),
				expectDuration(time.Duration(255), nil),
			},
		},
		"tc#011": {
			Input: uint16(math.MaxUint16),
			Asserters: []ConverterTester{
				expectByte(255, expectOverFlowError),
				expectInt8(-1, expectOverFlowError),
				expectInt16(-1, expectOverFlowError),
				expectInt32(math.MaxUint16, nil),
				expectInt64(math.MaxUint16, nil),
				expectInt(math.MaxUint16, nil),
				expectUInt8(255, expectOverFlowError),
				expectUint16(math.MaxUint16, nil),
				expectUint32(math.MaxUint16, nil),
				expectUint64(math.MaxUint16, nil),
				expectUint(math.MaxUint16, nil),
				expectFloat32(math.MaxUint16, nil),
				expectFloat64(math.MaxUint16, nil),
				expectString("65535", nil),
				expectBool(true, nil),
				expectTime(time.Date(1970, time.January, 1, 18, 12, 15, 0, time.UTC), nil),
				expectDuration(time.Duration(65535), nil),
			},
		},
		"tc#012": {
			Input: uint32(math.MaxUint32),
			Asserters: []ConverterTester{
				expectByte(255, expectOverFlowError),
				expectInt8(-1, expectOverFlowError),
				expectInt16(-1, expectOverFlowError),
				expectInt32(-1, expectOverFlowError),
				expectInt64(4294967295, nil),
				splitBasedOnArch(expectInt(0, expectOverFlowError), expectInt(math.MaxUint32, nil)),
				expectUInt8(255, expectOverFlowError),
				expectUint16(65535, expectOverFlowError),
				expectUint32(math.MaxUint32, nil),
				expectUint64(math.MaxUint32, nil),
				expectUint(math.MaxUint32, nil),
				expectFloat32(math.MaxUint32, nil),
				expectFloat64(math.MaxUint32, nil),
				expectString("4294967295", nil),
				expectBool(true, nil),
				expectTime(time.Date(2106, time.February, 7, 6, 28, 15, 0, time.UTC), nil),
				expectDuration(time.Duration(math.MaxUint32), nil),
			},
		},
		"tc#013": {
			Input: uint64(math.MaxUint64),
			Asserters: []ConverterTester{
				expectByte(255, expectOverFlowError),
				expectInt8(-1, expectOverFlowError),
				expectInt16(-1, expectOverFlowError),
				expectInt32(-1, expectOverFlowError),
				expectInt64(-1, expectOverFlowError),
				splitBasedOnArch(expectInt(-1, expectOverFlowError), expectInt(-1, expectOverFlowError)),
				expectUInt8(255, expectOverFlowError),
				expectUint16(65535, expectOverFlowError),
				expectUint32(math.MaxUint32, expectOverFlowError),
				expectUint64(math.MaxUint64, nil),
				splitBasedOnArch(expectUint(0, expectOverFlowError), expectUint(math.MaxUint64, nil)),
				expectFloat32(math.MaxUint64, nil),
				expectFloat64(math.MaxUint64, nil),
				expectString("18446744073709551615", nil),
				expectBool(true, nil),
				expectTime(time.Date(1969, time.December, 31, 23, 59, 59, 0, time.UTC), expectOverFlowError), // max uint64 is not converted to valid date.
				expectDuration(time.Duration(-1), expectOverFlowError),
			},
		},
		"tc#014": {
			Input: byte(12),
			Asserters: []ConverterTester{
				expectByte(12, nil),
				expectInt8(12, nil),
				expectInt16(12, nil),
				expectInt32(12, nil),
				expectInt64(12, nil),
				expectInt(12, nil),
				expectUInt8(12, nil),
				expectUint16(12, nil),
				expectUint32(12, nil),
				expectUint64(12, nil),
				expectUint(12, nil),
				expectFloat32(12, nil),
				expectFloat64(12, nil),
				expectString("12", nil),
				expectBool(true, nil),
				expectTime(time.Date(1970, time.January, 1, 0, 0, 12, 0, time.UTC), nil),
				expectDuration(time.Duration(12), nil),
			},
		},
		"tc#015": {
			Input: "123",
			Asserters: []ConverterTester{
				expectByte(0, expectInvalidType),
				expectInt8(123, nil),
				expectInt16(123, nil),
				expectInt32(123, nil),
				expectInt64(123, nil),
				expectInt(123, nil),
				expectUInt8(123, nil),
				expectUint16(123, nil),
				expectUint32(123, nil),
				expectUint64(123, nil),
				expectUint(123, nil),
				expectFloat32(123, nil),
				expectFloat64(123, nil),
				expectString("123", nil),
				expectBool(false, expectMalformedSyntax),
				expectTime(time.Time{}, tst.ExpectedErrorOfType[*time.ParseError]()),
				expectDuration(time.Duration(0), tst.ExpectedErrorStringContains("time: missing unit in duration")),
			},
		},
		"tc#016": {
			Input: []byte("123"),
			Asserters: []ConverterTester{
				expectByte(0, expectInvalidType),
				expectInt8(123, nil),
				expectInt16(123, nil),
				expectInt32(123, nil),
				expectInt64(123, nil),
				expectInt(123, nil),
				expectUInt8(123, nil),
				expectUint16(123, nil),
				expectUint32(123, nil),
				expectUint64(123, nil),
				expectUint(123, nil),
				expectFloat32(123, nil),
				expectFloat64(123, nil),
				expectString("123", nil),
				expectBool(false, expectMalformedSyntax),
				expectTime(time.Time{}, tst.ExpectedErrorOfType[*time.ParseError]()),
				expectDuration(time.Duration(0), tst.ExpectedErrorStringContains("time: missing unit in duration")),
			},
		},
		"tc#017": {
			Input: "123.321",
			Asserters: []ConverterTester{
				expectByte(0, expectInvalidType),
				expectInt8(123, expectLostDecimals),
				expectInt16(123, expectLostDecimals),
				expectInt32(123, expectLostDecimals),
				expectInt64(123, expectLostDecimals),
				expectInt(123, expectLostDecimals),
				expectUInt8(123, expectLostDecimals),
				expectUint16(123, expectLostDecimals),
				expectUint32(123, expectLostDecimals),
				expectUint64(123, expectLostDecimals),
				expectUint(123, expectLostDecimals),
				expectFloat32(123.321, nil),
				expectFloat64(123.321, nil),
				expectString("123.321", nil),
				expectBool(false, expectMalformedSyntax),
				expectTime(time.Time{}, tst.ExpectedErrorOfType[*time.ParseError]()),
				expectDuration(time.Duration(0), tst.ExpectedErrorStringContains("time: missing unit in duration")),
			},
		},
		"tc#018": {
			Input: stringAlias("23"),
			Asserters: []ConverterTester{
				expectByte(0, expectInvalidType),
				expectInt8(23, nil),
				expectInt16(23, nil),
				expectInt32(23, nil),
				expectInt64(23, nil),
				expectInt(23, nil),
				expectUInt8(23, nil),
				expectUint16(23, nil),
				expectUint32(23, nil),
				expectUint64(23, nil),
				expectUint(23, nil),
				expectFloat32(23, nil),
				expectFloat64(23, nil),
				expectString("23", nil),
				expectBool(false, expectMalformedSyntax),
				expectTime(time.Time{}, tst.ExpectedErrorOfType[*time.ParseError]()),
				expectDuration(time.Duration(0), tst.ExpectedErrorStringContains("time: missing unit in duration")),
			},
		},
		"tc#019": {
			Input: "just string",
			Asserters: []ConverterTester{
				expectByte(0, expectInvalidType),
				expectInt8(0, expectMalformedSyntax),
				expectInt16(0, expectMalformedSyntax),
				expectInt32(0, expectMalformedSyntax),
				expectInt64(0, expectMalformedSyntax),
				expectInt(0, expectMalformedSyntax),
				expectUInt8(0, expectMalformedSyntax),
				expectUint16(0, expectMalformedSyntax),
				expectUint32(0, expectMalformedSyntax),
				expectUint64(0, expectMalformedSyntax),
				expectUint(0, expectMalformedSyntax),
				expectFloat32(0, expectMalformedSyntax),
				expectFloat64(0, expectMalformedSyntax),
				expectString("just string", nil),
				expectBool(false, expectMalformedSyntax),
				expectTime(time.Time{}, tst.ExpectedErrorOfType[*time.ParseError]()),
				expectDuration(time.Duration(0), tst.ExpectedErrorStringContains("time: invalid duration")),
			},
		},
		"tc#020": {
			Input: []byte("byte slice"),
			Asserters: []ConverterTester{
				expectByte(0, expectInvalidType),
				expectInt8(0, expectMalformedSyntax),
				expectInt16(0, expectMalformedSyntax),
				expectInt32(0, expectMalformedSyntax),
				expectInt64(0, expectMalformedSyntax),
				expectInt(0, expectMalformedSyntax),
				expectUInt8(0, expectMalformedSyntax),
				expectUint16(0, expectMalformedSyntax),
				expectUint32(0, expectMalformedSyntax),
				expectUint64(0, expectMalformedSyntax),
				expectUint(0, expectMalformedSyntax),
				expectFloat32(0, expectMalformedSyntax),
				expectFloat64(0, expectMalformedSyntax),
				expectString("byte slice", nil),
				expectBool(false, expectMalformedSyntax),
				expectTime(time.Time{}, tst.ExpectedErrorOfType[*time.ParseError]()),
				expectDuration(time.Duration(0), tst.ExpectedErrorStringContains("time: invalid duration")),
			},
		},
		"tc#021": {
			Input: float32(123),
			Asserters: []ConverterTester{
				expectByte(123, nil),
				expectInt8(123, nil),
				expectInt16(123, nil),
				expectInt32(123, nil),
				expectInt64(123, nil),
				expectInt(123, nil),
				expectUInt8(123, nil),
				expectUint16(123, nil),
				expectUint32(123, nil),
				expectUint64(123, nil),
				expectUint(123, nil),
				expectFloat32(123, nil),
				expectFloat64(123, nil),
				expectString("123", nil),
				expectBool(true, nil),
				expectTime(time.Date(1970, time.January, 1, 0, 2, 3, 0, time.UTC), nil),
				expectDuration(time.Duration(123), nil),
			},
		},
		"tc#022": {
			Input: float64(123),
			Asserters: []ConverterTester{
				expectByte(123, nil),
				expectInt8(123, nil),
				expectInt16(123, nil),
				expectInt32(123, nil),
				expectInt64(123, nil),
				expectInt(123, nil),
				expectUInt8(123, nil),
				expectUint16(123, nil),
				expectUint32(123, nil),
				expectUint64(123, nil),
				expectUint(123, nil),
				expectFloat32(123, nil),
				expectFloat64(123, nil),
				expectString("123", nil),
				expectTime(time.Date(1970, time.January, 1, 0, 2, 3, 0, time.UTC), nil),
				expectDuration(time.Duration(123), nil),
			},
		},
		"tc#023": {
			Input: float64(123.12),
			Asserters: []ConverterTester{
				expectByte(123, expectLostDecimals),
				expectInt8(123, expectLostDecimals),
				expectInt16(123, expectLostDecimals),
				expectInt32(123, expectLostDecimals),
				expectInt64(123, expectLostDecimals),
				expectInt(123, expectLostDecimals),
				expectUInt8(123, expectLostDecimals),
				expectUint16(123, expectLostDecimals),
				expectUint32(123, expectLostDecimals),
				expectUint64(123, expectLostDecimals),
				expectUint(123, expectLostDecimals),
				expectFloat32(123.12, nil),
				expectFloat64(123.12, nil),
				expectString("123.12", nil),
				expectTime(time.Date(1970, time.January, 1, 0, 2, 3, 0, time.UTC), expectLostDecimals),
				expectDuration(time.Duration(123), expectLostDecimals),
			},
		},
		"tc#024": {
			Input: float64(math.MaxFloat64),
			Asserters: []ConverterTester{
				expectByte(0, expectOverFlowError),
				expectInt8(0, expectOverFlowError),
				expectInt16(0, expectOverFlowError),
				expectInt32(0, expectOverFlowError),
				expectInt64(-9223372036854775808, expectOverFlowError),
				splitBasedOnArch(expectInt(0, expectOverFlowError), expectInt(-9223372036854775808, expectOverFlowError)),
				expectUInt8(0, expectOverFlowError),
				expectUint16(0, expectOverFlowError),
				expectUint32(0, expectOverFlowError),
				expectUint64(9223372036854775808, expectOverFlowError),
				splitBasedOnArch(expectUint(0, expectOverFlowError), expectUint(9223372036854775808, expectOverFlowError)),
				expectFloat32(float32(math.Inf(1)), expectOverFlowError),
				expectFloat64(math.MaxFloat64, nil),
				expectString("1.7977E+308", nil),
				expectBool(true, nil),
				expectTime(time.Date(292277026596, time.December, 4, 15, 30, 8, 0, time.UTC), expectOverFlowError),
				expectDuration(time.Duration(-9223372036854775808), expectOverFlowError),
			},
		},
		"tc#025": {
			Input: struct{}{},
			Asserters: []ConverterTester{
				expectByte(0, expectInvalidType),
				expectInt8(0, expectInvalidType),
				expectInt16(0, expectInvalidType),
				expectInt32(0, expectInvalidType),
				expectInt64(0, expectInvalidType),
				expectInt(0, expectInvalidType),
				expectUInt8(0, expectInvalidType),
				expectUint16(0, expectInvalidType),
				expectUint32(0, expectInvalidType),
				expectUint64(0, expectInvalidType),
				expectUint(0, expectInvalidType),
				expectFloat32(0, expectInvalidType),
				expectFloat64(0, expectInvalidType),
				expectString("", expectInvalidType),
				expectBool(false, expectInvalidType),
				expectTime(time.Time{}, expectInvalidType),
				expectDuration(time.Duration(0), expectInvalidType),
			},
		},
		"tc#026": {
			Input: json.RawMessage(`{"a":"b"}`),
			Asserters: []ConverterTester{
				expectByte(0, expectInvalidType),
				expectInt8(0, expectInvalidType),
				expectInt16(0, expectInvalidType),
				expectInt32(0, expectInvalidType),
				expectInt64(0, expectInvalidType),
				expectInt(0, expectInvalidType),
				expectUInt8(0, expectInvalidType),
				expectUint16(0, expectInvalidType),
				expectUint32(0, expectInvalidType),
				expectUint64(0, expectInvalidType),
				expectUint(0, expectInvalidType),
				expectFloat32(0, expectInvalidType),
				expectFloat64(0, expectInvalidType),
				expectString(`{"a":"b"}`, nil),
				expectBool(false, expectInvalidType),
				expectTime(time.Time{}, expectInvalidType),
				expectDuration(time.Duration(0), expectInvalidType),
			},
		},
		"tc#027": {
			Input: json.Number("123"),
			Asserters: []ConverterTester{
				expectByte(123, nil),
				expectInt8(123, nil),
				expectInt16(123, nil),
				expectInt32(123, nil),
				expectInt64(123, nil),
				expectInt(123, nil),
				expectUInt8(123, nil),
				expectUint16(123, nil),
				expectUint32(123, nil),
				expectUint64(123, nil),
				expectUint(123, nil),
				expectFloat32(123, nil),
				expectFloat64(123, nil),
				expectString("123", nil),
				expectBool(true, nil),
				expectTime(time.Date(1970, time.January, 1, 0, 2, 3, 0, time.UTC), nil),
				expectDuration(time.Duration(123), nil),
			},
		},
		"tc#028": {
			Input: json.Number("56782"),
			Asserters: []ConverterTester{
				expectByte(206, expectOverFlowError),
				expectInt8(-50, expectOverFlowError),
				expectInt16(-8754, expectOverFlowError),
				expectInt32(56782, nil),
				expectInt64(56782, nil),
				expectInt(56782, nil),
				expectUInt8(206, expectOverFlowError),
				expectUint16(56782, nil),
				expectUint32(56782, nil),
				expectUint64(56782, nil),
				expectUint(56782, nil),
				expectFloat32(56782, nil),
				expectFloat64(56782, nil),
				expectString("56782", nil),
				expectBool(true, nil),
				expectTime(time.Date(1970, time.January, 1, 15, 46, 22, 0, time.UTC), nil),
				expectDuration(time.Duration(56782), nil),
			},
		},
		"tc#029": {
			Input: "1.79769313486231570814527423731704356798070e+308",
			Asserters: []ConverterTester{
				expectByte(0, expectInvalidType),
				expectInt8(0, expectOverFlowError),
				expectInt16(0, expectOverFlowError),
				expectInt32(0, expectOverFlowError),
				expectInt64(-9223372036854775808, expectOverFlowError),
				splitBasedOnArch(expectInt(0, expectOverFlowError), expectInt(-9223372036854775808, expectOverFlowError)),
				expectUInt8(0, expectOverFlowError),
				expectUint16(0, expectOverFlowError),
				expectUint32(0, expectOverFlowError),
				expectUint64(9223372036854775808, expectOverFlowError),
				splitBasedOnArch(expectUint(0, expectOverFlowError), expectUint(9223372036854775808, expectOverFlowError)),
				expectFloat32(float32(math.Inf(1)), expectOverFlowError),
				expectFloat64(math.MaxFloat64, nil),
				expectString("1.79769313486231570814527423731704356798070e+308", nil),
				expectBool(false, expectMalformedSyntax),
				expectTime(time.Time{}, tst.ExpectedErrorOfType[*time.ParseError]()),
				expectDuration(time.Duration(0), tst.ExpectedErrorStringContains("time: unknown unit")),
			},
		},
		"tc#030": {
			Input: true,
			Asserters: []ConverterTester{
				expectByte(1, nil),
				expectInt8(1, nil),
				expectInt16(1, nil),
				expectInt32(1, nil),
				expectInt64(1, nil),
				expectInt(1, nil),
				expectUInt8(1, nil),
				expectUint16(1, nil),
				expectUint32(1, nil),
				expectUint64(1, nil),
				expectUint(1, nil),
				expectFloat32(1, nil),
				expectFloat64(1, nil),
				expectString("true", nil),
				expectBool(true, nil),
				expectTime(time.Time{}, expectInvalidType),
				expectDuration(time.Duration(0), expectInvalidType),
			},
		},
		"tc#031": {
			Input: false,
			Asserters: []ConverterTester{
				expectByte(0, nil),
				expectInt8(0, nil),
				expectInt16(0, nil),
				expectInt32(0, nil),
				expectInt64(0, nil),
				expectInt(0, nil),
				expectUInt8(0, nil),
				expectUint16(0, nil),
				expectUint32(0, nil),
				expectUint64(0, nil),
				expectUint(0, nil),
				expectFloat32(0, nil),
				expectFloat64(0, nil),
				expectString("false", nil),
				expectBool(false, nil),
				expectTime(time.Time{}, expectInvalidType),
				expectDuration(time.Duration(0), expectInvalidType),
			},
		},
		"tc#032": {
			Input: int(12),
			Asserters: []ConverterTester{
				expectByte(12, nil),
				expectInt8(12, nil),
				expectInt16(12, nil),
				expectInt32(12, nil),
				expectInt64(12, nil),
				expectInt(12, nil),
				expectUInt8(12, nil),
				expectUint16(12, nil),
				expectUint32(12, nil),
				expectUint64(12, nil),
				expectUint(12, nil),
				expectFloat32(12, nil),
				expectFloat64(12, nil),
				expectString("12", nil),
				expectBool(true, nil),
				expectTime(time.Date(1970, time.January, 1, 0, 0, 12, 0, time.UTC), nil),
				expectDuration(time.Duration(12), nil),
			},
		},
		"tc#033": {
			Input: uint(12),
			Asserters: []ConverterTester{
				expectByte(12, nil),
				expectInt8(12, nil),
				expectInt16(12, nil),
				expectInt32(12, nil),
				expectInt64(12, nil),
				expectInt(12, nil),
				expectUInt8(12, nil),
				expectUint16(12, nil),
				expectUint32(12, nil),
				expectUint64(12, nil),
				expectUint(12, nil),
				expectFloat32(12, nil),
				expectFloat64(12, nil),
				expectString("12", nil),
				expectBool(true, nil),
				expectTime(time.Date(1970, time.January, 1, 0, 0, 12, 0, time.UTC), nil),
				expectDuration(time.Duration(12), nil),
			},
		},
	}

	for k, tc := range testCases {
		name := fmt.Sprintf("%s__%T(%#v)", k, tc.Input, tc.Input)
		t.Run(name, func(t *testing.T) {
			for _, a := range tc.Asserters {
				a.SetInput(tc.Input)
				a.Test(t)
			}
		})
	}
}

func TestConverterSliceMatrix(t *testing.T) {
	converter := NewDefaultConverter()

	// matrixExpectedResult constructor function aliases.
	expectByte := func(expected []byte, errorAssertFn tst.ErrorAsserter) *converterTestCase[[]byte] {
		return &converterTestCase[[]byte]{
			Converter:                converter,
			Input:                    nil,
			Expected:                 expected,
			ErrorAsserter:            errorAssertFn,
			OverwriteDirectConvertFn: converter.AsByteSlice,
			OmitConvertByDirectFn:    false,
			OmitConvertByKind:        true,
			OmitConvertByType:        true,
		}
	}
	expectInt8 := matrixTestConstructorFn[[]int8](converter)
	expectInt16 := matrixTestConstructorFn[[]int16](converter)
	expectInt32 := matrixTestConstructorFn[[]int32](converter)
	expectInt64 := matrixTestConstructorFn[[]int64](converter)
	expectInt := matrixTestConstructorFn[[]int](converter)
	expectUInt8 := func(expected []uint8, errorAssertFn tst.ErrorAsserter) *converterTestCase[[]uint8] {
		return &converterTestCase[[]uint8]{
			Converter:                converter,
			Input:                    nil,
			Expected:                 expected,
			ErrorAsserter:            errorAssertFn,
			OverwriteDirectConvertFn: converter.AsUint8Slice,
			OmitConvertByDirectFn:    false,
			OmitConvertByKind:        false,
			OmitConvertByType:        false,
		}
	}
	expectUint16 := matrixTestConstructorFn[[]uint16](converter)
	expectUint32 := matrixTestConstructorFn[[]uint32](converter)
	expectUint64 := matrixTestConstructorFn[[]uint64](converter)
	expectUint := matrixTestConstructorFn[[]uint](converter)
	expectFloat32 := matrixTestConstructorFn[[]float32](converter)
	expectFloat64 := matrixTestConstructorFn[[]float64](converter)
	expectString := matrixTestConstructorFn[[]string](converter)
	expectBool := matrixTestConstructorFn[[]bool](converter)
	expectTime := matrixTestConstructorFn[[]time.Time](converter)
	expectDuration := matrixTestConstructorFn[[]time.Duration](converter)

	_ = []ConverterTester{
		expectByte([]byte{}, nil),
		expectInt8([]int8{}, nil),
		expectInt16([]int16{}, nil),
		expectInt32([]int32{}, nil),
		expectInt64([]int64{}, nil),
		expectInt([]int{}, nil),
		expectUInt8([]uint8{}, nil),
		expectUint16([]uint16{}, nil),
		expectUint32([]uint32{}, nil),
		expectUint64([]uint64{}, nil),
		expectUint([]uint{}, nil),
		expectFloat32([]float32{}, nil),
		expectFloat64([]float64{}, nil),
		expectString([]string{}, nil),
		expectBool([]bool{}, nil),
		expectTime([]time.Time{}, nil),
		expectDuration([]time.Duration{}, nil),
	}

	testCases := map[string]struct {
		Input     any
		Asserters []ConverterTester
	}{
		"tc#000": {
			Input: nil,
			Asserters: []ConverterTester{
				expectByte([]byte(nil), nil),
				expectInt8([]int8(nil), nil),
				expectInt16([]int16(nil), nil),
				expectInt32([]int32(nil), nil),
				expectInt64([]int64(nil), nil),
				expectInt([]int(nil), nil),
				expectUInt8([]uint8(nil), nil),
				expectUint16([]uint16(nil), nil),
				expectUint32([]uint32(nil), nil),
				expectUint64([]uint64(nil), nil),
				expectUint([]uint(nil), nil),
				expectFloat32([]float32(nil), nil),
				expectFloat64([]float64(nil), nil),
				expectString([]string(nil), nil),
				expectBool([]bool(nil), nil),
				expectTime([]time.Time(nil), nil),
				expectDuration([]time.Duration(nil), nil),
			},
		},
		"tc#001": {
			Input: []int8{1, 2, 3},
			Asserters: []ConverterTester{
				expectByte([]byte{1, 2, 3}, nil),
				expectInt8([]int8{1, 2, 3}, nil),
				expectInt16([]int16{1, 2, 3}, nil),
				expectInt32([]int32{1, 2, 3}, nil),
				expectInt64([]int64{1, 2, 3}, nil),
				expectInt([]int{1, 2, 3}, nil),
				expectUInt8([]uint8{1, 2, 3}, nil),
				expectUint16([]uint16{1, 2, 3}, nil),
				expectUint32([]uint32{1, 2, 3}, nil),
				expectUint64([]uint64{1, 2, 3}, nil),
				expectUint([]uint{1, 2, 3}, nil),
				expectFloat32([]float32{1, 2, 3}, nil),
				expectFloat64([]float64{1, 2, 3}, nil),
				expectString([]string{"1", "2", "3"}, nil),
				expectBool([]bool{true, true, true}, nil),
				expectTime([]time.Time{time.Unix(1, 0).UTC(), time.Unix(2, 0).UTC(), time.Unix(3, 0).UTC()}, nil),
				expectDuration([]time.Duration{1, 2, 3}, nil),
			},
		},
		"tc#002": {
			Input: []int16{1, 2, 3},
			Asserters: []ConverterTester{
				expectByte([]byte{1, 2, 3}, nil),
				expectInt8([]int8{1, 2, 3}, nil),
				expectInt16([]int16{1, 2, 3}, nil),
				expectInt32([]int32{1, 2, 3}, nil),
				expectInt64([]int64{1, 2, 3}, nil),
				expectInt([]int{1, 2, 3}, nil),
				expectUInt8([]uint8{1, 2, 3}, nil),
				expectUint16([]uint16{1, 2, 3}, nil),
				expectUint32([]uint32{1, 2, 3}, nil),
				expectUint64([]uint64{1, 2, 3}, nil),
				expectUint([]uint{1, 2, 3}, nil),
				expectFloat32([]float32{1, 2, 3}, nil),
				expectFloat64([]float64{1, 2, 3}, nil),
				expectString([]string{"1", "2", "3"}, nil),
				expectBool([]bool{true, true, true}, nil),
				expectTime([]time.Time{time.Unix(1, 0).UTC(), time.Unix(2, 0).UTC(), time.Unix(3, 0).UTC()}, nil),
				expectDuration([]time.Duration{1, 2, 3}, nil),
			},
		},
		"tc#003": {
			Input: []int32{1, 2, 3},
			Asserters: []ConverterTester{
				expectByte([]byte{1, 2, 3}, nil),
				expectInt8([]int8{1, 2, 3}, nil),
				expectInt16([]int16{1, 2, 3}, nil),
				expectInt32([]int32{1, 2, 3}, nil),
				expectInt64([]int64{1, 2, 3}, nil),
				expectInt([]int{1, 2, 3}, nil),
				expectUInt8([]uint8{1, 2, 3}, nil),
				expectUint16([]uint16{1, 2, 3}, nil),
				expectUint32([]uint32{1, 2, 3}, nil),
				expectUint64([]uint64{1, 2, 3}, nil),
				expectUint([]uint{1, 2, 3}, nil),
				expectFloat32([]float32{1, 2, 3}, nil),
				expectFloat64([]float64{1, 2, 3}, nil),
				expectString([]string{"1", "2", "3"}, nil),
				expectBool([]bool{true, true, true}, nil),
				expectTime([]time.Time{time.Unix(1, 0).UTC(), time.Unix(2, 0).UTC(), time.Unix(3, 0).UTC()}, nil),
				expectDuration([]time.Duration{1, 2, 3}, nil),
			},
		},
		"tc#004": {
			Input: []int64{1, 2, 3},
			Asserters: []ConverterTester{
				expectByte([]byte{1, 2, 3}, nil),
				expectInt8([]int8{1, 2, 3}, nil),
				expectInt16([]int16{1, 2, 3}, nil),
				expectInt32([]int32{1, 2, 3}, nil),
				expectInt64([]int64{1, 2, 3}, nil),
				expectInt([]int{1, 2, 3}, nil),
				expectUInt8([]uint8{1, 2, 3}, nil),
				expectUint16([]uint16{1, 2, 3}, nil),
				expectUint32([]uint32{1, 2, 3}, nil),
				expectUint64([]uint64{1, 2, 3}, nil),
				expectUint([]uint{1, 2, 3}, nil),
				expectFloat32([]float32{1, 2, 3}, nil),
				expectFloat64([]float64{1, 2, 3}, nil),
				expectString([]string{"1", "2", "3"}, nil),
				expectBool([]bool{true, true, true}, nil),
				expectTime([]time.Time{time.Unix(1, 0).UTC(), time.Unix(2, 0).UTC(), time.Unix(3, 0).UTC()}, nil),
				expectDuration([]time.Duration{1, 2, 3}, nil),
			},
		},
		"tc#005": {
			Input: []int{1, 2, 3},
			Asserters: []ConverterTester{
				expectByte([]byte{1, 2, 3}, nil),
				expectInt8([]int8{1, 2, 3}, nil),
				expectInt16([]int16{1, 2, 3}, nil),
				expectInt32([]int32{1, 2, 3}, nil),
				expectInt64([]int64{1, 2, 3}, nil),
				expectInt([]int{1, 2, 3}, nil),
				expectUInt8([]uint8{1, 2, 3}, nil),
				expectUint16([]uint16{1, 2, 3}, nil),
				expectUint32([]uint32{1, 2, 3}, nil),
				expectUint64([]uint64{1, 2, 3}, nil),
				expectUint([]uint{1, 2, 3}, nil),
				expectFloat32([]float32{1, 2, 3}, nil),
				expectFloat64([]float64{1, 2, 3}, nil),
				expectString([]string{"1", "2", "3"}, nil),
				expectBool([]bool{true, true, true}, nil),
				expectTime([]time.Time{time.Unix(1, 0).UTC(), time.Unix(2, 0).UTC(), time.Unix(3, 0).UTC()}, nil),
				expectDuration([]time.Duration{1, 2, 3}, nil),
			},
		},
		"tc#006": {
			Input: []uint8{1, 2, 3},
			Asserters: []ConverterTester{
				expectByte([]byte{1, 2, 3}, nil),
				expectInt8([]int8{1, 2, 3}, nil),
				expectInt16([]int16{1, 2, 3}, nil),
				expectInt32([]int32{1, 2, 3}, nil),
				expectInt64([]int64{1, 2, 3}, nil),
				expectInt([]int{1, 2, 3}, nil),
				expectUInt8([]uint8{1, 2, 3}, nil),
				expectUint16([]uint16{1, 2, 3}, nil),
				expectUint32([]uint32{1, 2, 3}, nil),
				expectUint64([]uint64{1, 2, 3}, nil),
				expectUint([]uint{1, 2, 3}, nil),
				expectFloat32([]float32{1, 2, 3}, nil),
				expectFloat64([]float64{1, 2, 3}, nil),
				expectString([]string{"1", "2", "3"}, nil),
				expectBool([]bool{true, true, true}, nil),
				expectTime([]time.Time{time.Unix(1, 0).UTC(), time.Unix(2, 0).UTC(), time.Unix(3, 0).UTC()}, nil),
				expectDuration([]time.Duration{1, 2, 3}, nil),
			},
		},
		"tc#007": {
			Input: []uint16{1, 2, 3},
			Asserters: []ConverterTester{
				expectByte([]byte{1, 2, 3}, nil),
				expectInt8([]int8{1, 2, 3}, nil),
				expectInt16([]int16{1, 2, 3}, nil),
				expectInt32([]int32{1, 2, 3}, nil),
				expectInt64([]int64{1, 2, 3}, nil),
				expectInt([]int{1, 2, 3}, nil),
				expectUInt8([]uint8{1, 2, 3}, nil),
				expectUint16([]uint16{1, 2, 3}, nil),
				expectUint32([]uint32{1, 2, 3}, nil),
				expectUint64([]uint64{1, 2, 3}, nil),
				expectUint([]uint{1, 2, 3}, nil),
				expectFloat32([]float32{1, 2, 3}, nil),
				expectFloat64([]float64{1, 2, 3}, nil),
				expectString([]string{"1", "2", "3"}, nil),
				expectBool([]bool{true, true, true}, nil),
				expectTime([]time.Time{time.Unix(1, 0).UTC(), time.Unix(2, 0).UTC(), time.Unix(3, 0).UTC()}, nil),
				expectDuration([]time.Duration{1, 2, 3}, nil),
			},
		},
		"tc#008": {
			Input: []uint32{1, 2, 3},
			Asserters: []ConverterTester{
				expectByte([]byte{1, 2, 3}, nil),
				expectInt8([]int8{1, 2, 3}, nil),
				expectInt16([]int16{1, 2, 3}, nil),
				expectInt32([]int32{1, 2, 3}, nil),
				expectInt64([]int64{1, 2, 3}, nil),
				expectInt([]int{1, 2, 3}, nil),
				expectUInt8([]uint8{1, 2, 3}, nil),
				expectUint16([]uint16{1, 2, 3}, nil),
				expectUint32([]uint32{1, 2, 3}, nil),
				expectUint64([]uint64{1, 2, 3}, nil),
				expectUint([]uint{1, 2, 3}, nil),
				expectFloat32([]float32{1, 2, 3}, nil),
				expectFloat64([]float64{1, 2, 3}, nil),
				expectString([]string{"1", "2", "3"}, nil),
				expectBool([]bool{true, true, true}, nil),
				expectTime([]time.Time{time.Unix(1, 0).UTC(), time.Unix(2, 0).UTC(), time.Unix(3, 0).UTC()}, nil),
				expectDuration([]time.Duration{1, 2, 3}, nil),
			},
		},
		"tc#009": {
			Input: []uint64{1, 2, 3},
			Asserters: []ConverterTester{
				expectByte([]byte{1, 2, 3}, nil),
				expectInt8([]int8{1, 2, 3}, nil),
				expectInt16([]int16{1, 2, 3}, nil),
				expectInt32([]int32{1, 2, 3}, nil),
				expectInt64([]int64{1, 2, 3}, nil),
				expectInt([]int{1, 2, 3}, nil),
				expectUInt8([]uint8{1, 2, 3}, nil),
				expectUint16([]uint16{1, 2, 3}, nil),
				expectUint32([]uint32{1, 2, 3}, nil),
				expectUint64([]uint64{1, 2, 3}, nil),
				expectUint([]uint{1, 2, 3}, nil),
				expectFloat32([]float32{1, 2, 3}, nil),
				expectFloat64([]float64{1, 2, 3}, nil),
				expectString([]string{"1", "2", "3"}, nil),
				expectBool([]bool{true, true, true}, nil),
				expectTime([]time.Time{time.Unix(1, 0).UTC(), time.Unix(2, 0).UTC(), time.Unix(3, 0).UTC()}, nil),
				expectDuration([]time.Duration{1, 2, 3}, nil),
			},
		},
		"tc#010": {
			Input: []uint{1, 2, 3},
			Asserters: []ConverterTester{
				expectByte([]byte{1, 2, 3}, nil),
				expectInt8([]int8{1, 2, 3}, nil),
				expectInt16([]int16{1, 2, 3}, nil),
				expectInt32([]int32{1, 2, 3}, nil),
				expectInt64([]int64{1, 2, 3}, nil),
				expectInt([]int{1, 2, 3}, nil),
				expectUInt8([]uint8{1, 2, 3}, nil),
				expectUint16([]uint16{1, 2, 3}, nil),
				expectUint32([]uint32{1, 2, 3}, nil),
				expectUint64([]uint64{1, 2, 3}, nil),
				expectUint([]uint{1, 2, 3}, nil),
				expectFloat32([]float32{1, 2, 3}, nil),
				expectFloat64([]float64{1, 2, 3}, nil),
				expectString([]string{"1", "2", "3"}, nil),
				expectBool([]bool{true, true, true}, nil),
				expectTime([]time.Time{time.Unix(1, 0).UTC(), time.Unix(2, 0).UTC(), time.Unix(3, 0).UTC()}, nil),
				expectDuration([]time.Duration{1, 2, 3}, nil),
			},
		},
		"tc#011": {
			Input: []string{"1", "2", "3"},
			Asserters: []ConverterTester{
				expectByte([]byte(nil), expectInvalidType),
				expectInt8([]int8{1, 2, 3}, nil),
				expectInt16([]int16{1, 2, 3}, nil),
				expectInt32([]int32{1, 2, 3}, nil),
				expectInt64([]int64{1, 2, 3}, nil),
				expectInt([]int{1, 2, 3}, nil),
				expectUInt8([]uint8{1, 2, 3}, nil),
				expectUint16([]uint16{1, 2, 3}, nil),
				expectUint32([]uint32{1, 2, 3}, nil),
				expectUint64([]uint64{1, 2, 3}, nil),
				expectUint([]uint{1, 2, 3}, nil),
				expectFloat32([]float32{1, 2, 3}, nil),
				expectFloat64([]float64{1, 2, 3}, nil),
				expectString([]string{"1", "2", "3"}, nil),
				expectBool([]bool(nil), expectMalformedSyntax),
				expectTime([]time.Time(nil), tst.ExpectedErrorStringContains("error: parsing time")),
				expectDuration([]time.Duration(nil), tst.ExpectedErrorStringContains("time: missing unit in duration")),
			},
		},
	}

	for k, tc := range testCases {
		name := fmt.Sprintf("%s__%T(%#v)", k, tc.Input, tc.Input)
		t.Run(name, func(t *testing.T) {
			for _, a := range tc.Asserters {
				a.SetInput(tc.Input)
				a.Test(t)
			}
		})
	}
}

func BenchmarkConverterSlice(b *testing.B) {
	testCases := []any{
		[]any{"abc", "def"},
		[]string{"abc", "def"},
		[]any{1, 2, 3, 4},
		[]int32{1, 2, 3, 4},
		[]any{"1", "2", "3", "4"},
		[]string{"1", "2", "3", "4"},
	}

	c := NewDefaultConverter()

	b.Run("AsBoolSlice", converterSubBenchmarks(testCases, c.AsBoolSlice))
	b.Run("AsByteSlice", converterSubBenchmarks(testCases, c.AsByteSlice))
	b.Run("AsFloat32Slice", converterSubBenchmarks(testCases, c.AsFloat32Slice))
	b.Run("AsFloat64Slice", converterSubBenchmarks(testCases, c.AsFloat64Slice))
	b.Run("AsIntSlice", converterSubBenchmarks(testCases, c.AsIntSlice))
	b.Run("AsInt8Slice", converterSubBenchmarks(testCases, c.AsInt8Slice))
	b.Run("AsInt16Slice", converterSubBenchmarks(testCases, c.AsInt16Slice))
	b.Run("AsInt32Slice", converterSubBenchmarks(testCases, c.AsInt32Slice))
	b.Run("AsInt64Slice", converterSubBenchmarks(testCases, c.AsInt64Slice))
	b.Run("AsUintSlice", converterSubBenchmarks(testCases, c.AsUintSlice))
	b.Run("AsUint8Slice", converterSubBenchmarks(testCases, c.AsUint8Slice))
	b.Run("AsUint16Slice", converterSubBenchmarks(testCases, c.AsUint16Slice))
	b.Run("AsUint32Slice", converterSubBenchmarks(testCases, c.AsUint32Slice))
	b.Run("AsUint64Slice", converterSubBenchmarks(testCases, c.AsUint64Slice))
	b.Run("AsStringSlice", converterSubBenchmarks(testCases, c.AsStringSlice))
}

func converterSubBenchmarks[Output any](testCases []any, convertFn func(any) (Output, error)) func(b *testing.B) {
	return func(b *testing.B) {
		b.Helper()
		for i, tc := range testCases {
			name := fmt.Sprintf("%d %s", i, tst.Format(tc))
			b.Run(name, matrixSubBenchmark(tc, convertFn))
		}
	}
}

func matrixSubBenchmark[Output any](input any, convertFn func(any) (Output, error)) func(b *testing.B) {
	return func(b *testing.B) {
		b.Helper()
		for range b.N {
			_, err := convertFn(input)
			if err != nil {
				b.Skipf("skipped because of error %s", err.Error())
			}
		}
	}
}
