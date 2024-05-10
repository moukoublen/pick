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

type CasterTester interface {
	Test(t *testing.T)
	SetInput(i any)
}

func matrixTestConstructorFn[Output any](c Caster) func(expected Output, errorAssertFn func(*testing.T, error)) *casterTestCaseMel[Output] {
	return func(expected Output, errorAssertFn func(*testing.T, error)) *casterTestCaseMel[Output] {
		return &casterTestCaseMel[Output]{
			Caster:                c,
			Input:                 nil,
			Expected:              expected,
			ExpectedErr:           errorAssertFn,
			OverwriteDirectCastFn: nil,
			OmitCastByDirectFn:    false,
			OmitCastByKind:        false,
			OmitCastByType:        false,
		}
	}
}

func splitBasedOnArch[Output any](for32bit, for64bit *casterTestCaseMel[Output]) *casterTestCaseMel[Output] {
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

	caster := NewCaster()

	type stringAlias string

	// matrixExpectedResult constructor function aliases.
	expectByte := func(expected byte, errorAssertFn func(*testing.T, error)) *casterTestCaseMel[byte] {
		return &casterTestCaseMel[byte]{
			Caster:                caster,
			Input:                 nil,
			Expected:              expected,
			ExpectedErr:           errorAssertFn,
			OverwriteDirectCastFn: caster.AsByte,
			OmitCastByDirectFn:    false,
			OmitCastByKind:        true,
			OmitCastByType:        true,
		}
	}
	expectInt8 := matrixTestConstructorFn[int8](caster)
	expectInt16 := matrixTestConstructorFn[int16](caster)
	expectInt32 := matrixTestConstructorFn[int32](caster)
	expectInt64 := matrixTestConstructorFn[int64](caster)
	expectInt := matrixTestConstructorFn[int](caster)
	expectUInt8 := func(expected uint8, errorAssertFn func(*testing.T, error)) *casterTestCaseMel[uint8] {
		return &casterTestCaseMel[uint8]{
			Caster:                caster,
			Input:                 nil,
			Expected:              expected,
			ExpectedErr:           errorAssertFn,
			OverwriteDirectCastFn: caster.AsUint8,
			OmitCastByDirectFn:    false,
			OmitCastByKind:        false,
			OmitCastByType:        false,
		}
	}
	expectUint16 := matrixTestConstructorFn[uint16](caster)
	expectUint32 := matrixTestConstructorFn[uint32](caster)
	expectUint64 := matrixTestConstructorFn[uint64](caster)
	expectUint := matrixTestConstructorFn[uint](caster)
	expectFloat32 := matrixTestConstructorFn[float32](caster)
	expectFloat64 := matrixTestConstructorFn[float64](caster)
	expectString := matrixTestConstructorFn[string](caster)
	expectBool := matrixTestConstructorFn[bool](caster)
	expectTime := matrixTestConstructorFn[time.Time](caster)
	expectDuration := matrixTestConstructorFn[time.Duration](caster)

	testCases := map[string]struct {
		Input     any
		Asserters []CasterTester
	}{
		"tc#000": {
			Input: nil,
			Asserters: []CasterTester{
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
			Asserters: []CasterTester{
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
			Asserters: []CasterTester{
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
			Asserters: []CasterTester{
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
			Asserters: []CasterTester{
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
			Asserters: []CasterTester{
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
			Asserters: []CasterTester{
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
			Asserters: []CasterTester{
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
			Asserters: []CasterTester{
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
			Asserters: []CasterTester{
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
			Asserters: []CasterTester{
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
			Asserters: []CasterTester{
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
			Asserters: []CasterTester{
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
			Asserters: []CasterTester{
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
			Asserters: []CasterTester{
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
			Asserters: []CasterTester{
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
				expectTime(time.Time{}, testingx.ExpectedErrorIsOfType(&time.ParseError{})),
				expectDuration(time.Duration(0), testingx.ExpectedErrorStringContains("time: missing unit in duration")),
			},
		},
		"tc#016": {
			Input: []byte("123"),
			Asserters: []CasterTester{
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
				expectTime(time.Time{}, testingx.ExpectedErrorIsOfType(&time.ParseError{})),
				expectDuration(time.Duration(0), testingx.ExpectedErrorStringContains("time: missing unit in duration")),
			},
		},
		"tc#017": {
			Input: "123.321",
			Asserters: []CasterTester{
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
				expectTime(time.Time{}, testingx.ExpectedErrorIsOfType(&time.ParseError{})),
				expectDuration(time.Duration(0), testingx.ExpectedErrorStringContains("time: missing unit in duration")),
			},
		},
		"tc#018": {
			Input: stringAlias("23"),
			Asserters: []CasterTester{
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
				expectTime(time.Time{}, testingx.ExpectedErrorIsOfType(&time.ParseError{})),
				expectDuration(time.Duration(0), testingx.ExpectedErrorStringContains("time: missing unit in duration")),
			},
		},
		"tc#019": {
			Input: "just string",
			Asserters: []CasterTester{
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
				expectTime(time.Time{}, testingx.ExpectedErrorIsOfType(&time.ParseError{})),
				expectDuration(time.Duration(0), testingx.ExpectedErrorStringContains("time: invalid duration")),
			},
		},
		"tc#020": {
			Input: []byte("byte slice"),
			Asserters: []CasterTester{
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
				expectTime(time.Time{}, testingx.ExpectedErrorIsOfType(&time.ParseError{})),
				expectDuration(time.Duration(0), testingx.ExpectedErrorStringContains("time: invalid duration")),
			},
		},
		"tc#021": {
			Input: float32(123),
			Asserters: []CasterTester{
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
			Asserters: []CasterTester{
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
			Asserters: []CasterTester{
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
			Asserters: []CasterTester{
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
			Asserters: []CasterTester{
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
			Asserters: []CasterTester{
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
			Asserters: []CasterTester{
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
			Asserters: []CasterTester{
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
			Asserters: []CasterTester{
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
				expectTime(time.Time{}, testingx.ExpectedErrorIsOfType(&time.ParseError{})),
				expectDuration(time.Duration(0), testingx.ExpectedErrorStringContains("time: unknown unit")),
			},
		},
		"tc#030": {
			Input: true,
			Asserters: []CasterTester{
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
			Asserters: []CasterTester{
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
	}

	for k, tc := range testCases {
		tc := tc
		name := fmt.Sprintf("%s__%T(%#v)", k, tc.Input, tc.Input)
		t.Run(name, func(t *testing.T) {
			for _, a := range tc.Asserters {
				a.SetInput(tc.Input)
				a.Test(t)
			}
		})
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
