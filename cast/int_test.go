package cast

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"runtime"
	"testing"

	"github.com/moukoublen/pick/internal/testingx"
	"github.com/moukoublen/pick/numbers"
)

func TestInt64CastValid(t *testing.T) {
	t.Parallel()

	type testCase struct {
		Origin   int64
		Expected bool
	}

	// map[To][]{Origin, Expected}

	testCasesFor32bit := map[reflect.Kind][]testCase{
		reflect.Int: {
			{Origin: math.MinInt64, Expected: false},
			{Origin: math.MinInt32 - 1, Expected: false},
			{Origin: math.MinInt32, Expected: true},
			{Origin: 0, Expected: true},
			{Origin: math.MaxInt32, Expected: true},
			{Origin: math.MaxInt32 + 1, Expected: false},
			{Origin: math.MaxInt64, Expected: false},
		},
	}

	testCasesFor64bit := map[reflect.Kind][]testCase{
		reflect.Int: {
			{Origin: math.MinInt64, Expected: true},
			{Origin: math.MinInt32 - 1, Expected: true},
			{Origin: math.MinInt32, Expected: true},
			{Origin: 0, Expected: true},
			{Origin: math.MaxInt32, Expected: true},
			{Origin: math.MaxInt32 + 1, Expected: true},
			{Origin: math.MaxInt64, Expected: true},
		},
	}

	testCases := map[reflect.Kind][]testCase{
		reflect.Uint8: {
			{Origin: -1, Expected: false},
			{Origin: 0, Expected: true},
			{Origin: 10, Expected: true},
			{Origin: math.MaxUint8, Expected: true},
			{Origin: math.MaxUint8 + 1, Expected: false},
		},

		reflect.Uint16: {
			{Origin: -1, Expected: false},
			{Origin: 0, Expected: true},
			{Origin: math.MaxUint16, Expected: true},
			{Origin: math.MaxUint16 + 1, Expected: false},
		},

		reflect.Uint32: {
			{Origin: -1, Expected: false},
			{Origin: 0, Expected: true},
			{Origin: math.MaxUint32, Expected: true},
			{Origin: math.MaxUint32 + 1, Expected: false},
		},

		reflect.Uint64: {
			{Origin: -1, Expected: false},
			{Origin: 0, Expected: true},
			{Origin: math.MaxUint32 + 1, Expected: true},
			{Origin: math.MaxInt64, Expected: true},
		},

		reflect.Int8: {
			{Origin: math.MinInt8 - 1, Expected: false},
			{Origin: math.MinInt8, Expected: true},
			{Origin: 0, Expected: true},
			{Origin: math.MaxInt8, Expected: true},
			{Origin: math.MaxInt8 + 1, Expected: false},
		},

		reflect.Int16: {
			{Origin: math.MinInt16 - 1, Expected: false},
			{Origin: math.MinInt16, Expected: true},
			{Origin: 0, Expected: true},
			{Origin: math.MaxInt16, Expected: true},
			{Origin: math.MaxInt16 + 1, Expected: false},
		},

		reflect.Int32: {
			{Origin: math.MinInt32 - 1, Expected: false},
			{Origin: math.MinInt32, Expected: true},
			{Origin: 0, Expected: true},
			{Origin: math.MaxInt32, Expected: true},
			{Origin: math.MaxInt32 + 1, Expected: false},
		},

		reflect.Int64: {
			{Origin: math.MinInt64, Expected: true},
			{Origin: 0, Expected: true},
			{Origin: math.MaxInt64, Expected: true},
		},
	}

	var toAppend map[reflect.Kind][]testCase
	switch runtime.GOARCH {
	case "arm", "386":
		toAppend = testCasesFor32bit
	default:
		toAppend = testCasesFor64bit
	}
	for k, v := range toAppend {
		testCases[k] = append(testCases[k], v...)
	}

	for to, perKindTestCases := range testCases {
		for _, tc := range perKindTestCases {
			t.Run(
				fmt.Sprintf("%d to %s", tc.Origin, to.String()),
				func(t *testing.T) {
					t.Parallel()
					got := int64CastValid(tc.Origin, to)
					if got != tc.Expected {
						t.Errorf("expected %#v got %#v", tc.Expected, got)
					}
				},
			)
		}
	}
}

func TestUint64CastValid(t *testing.T) {
	t.Parallel()

	type testCase struct {
		Origin   uint64
		Expected bool
	}

	// map[To][]{Origin, Expected}

	testCasesFor32bit := map[reflect.Kind][]testCase{
		reflect.Int: {
			{Origin: 0, Expected: true},
			{Origin: math.MaxInt32, Expected: true},
			{Origin: math.MaxInt32 + 1, Expected: false},
			{Origin: math.MaxInt64, Expected: false},
		},
	}

	testCasesFor64bit := map[reflect.Kind][]testCase{
		reflect.Int: {
			{Origin: 0, Expected: true},
			{Origin: math.MaxInt32, Expected: true},
			{Origin: math.MaxInt32 + 1, Expected: true},
			{Origin: math.MaxInt64, Expected: true},
		},
	}

	testCases := map[reflect.Kind][]testCase{
		reflect.Uint8: {
			{Origin: 0, Expected: true},
			{Origin: 10, Expected: true},
			{Origin: math.MaxUint8, Expected: true},
			{Origin: math.MaxUint8 + 1, Expected: false},
		},

		reflect.Uint16: {
			{Origin: 0, Expected: true},
			{Origin: math.MaxUint16, Expected: true},
			{Origin: math.MaxUint16 + 1, Expected: false},
		},

		reflect.Uint32: {
			{Origin: 0, Expected: true},
			{Origin: math.MaxUint32, Expected: true},
			{Origin: math.MaxUint32 + 1, Expected: false},
		},

		reflect.Uint64: {
			{Origin: 0, Expected: true},
			{Origin: math.MaxUint64, Expected: true},
		},

		reflect.Int8: {
			{Origin: 0, Expected: true},
			{Origin: math.MaxInt8, Expected: true},
			{Origin: math.MaxInt8 + 1, Expected: false},
		},

		reflect.Int16: {
			{Origin: 0, Expected: true},
			{Origin: math.MaxInt16, Expected: true},
			{Origin: math.MaxInt16 + 1, Expected: false},
		},

		reflect.Int32: {
			{Origin: 0, Expected: true},
			{Origin: math.MaxInt32, Expected: true},
			{Origin: math.MaxInt32 + 1, Expected: false},
		},

		reflect.Int64: {
			{Origin: 0, Expected: true},
			{Origin: math.MaxInt64, Expected: true},
			{Origin: math.MaxInt64 + 1, Expected: false},
		},
	}

	var toAppend map[reflect.Kind][]testCase
	switch runtime.GOARCH {
	case "arm", "386":
		toAppend = testCasesFor32bit
	default:
		toAppend = testCasesFor64bit
	}
	for k, v := range toAppend {
		testCases[k] = append(testCases[k], v...)
	}

	for to, perKindTestCases := range testCases {
		for _, tc := range perKindTestCases {
			t.Run(
				fmt.Sprintf("%d to %s", tc.Origin, to.String()),
				func(t *testing.T) {
					t.Parallel()
					got := uint64CastValid(tc.Origin, to)
					if got != tc.Expected {
						t.Errorf("expected %#v got %#v", tc.Expected, got)
					}
				},
			)
		}
	}
}

func BenchmarkIntCaster(b *testing.B) {
	ic := NewCaster()

	tests := []any{
		int8(123),
		int16(123),
		int32(123),
		int64(123),
		int(123),

		uint8(123),
		uint16(123),
		uint32(123),
		uint64(123),
		uint(123),

		float32(123),
		float64(123),

		"-128",                 // math.MinInt8
		"-32768",               // math.MinInt16
		"-2147483648",          // math.MinInt32
		"-9223372036854775808", // math.MinInt64
		"127",                  // math.MaxInt8
		"32767",                // math.MaxInt16
		"2147483647",           // math.MaxInt32
		"9223372036854775807",  // math.MaxInt64
		"255",                  // math.MaxUint8
		"65535",                // math.MaxUint16
		"4294967295",           // math.MaxUint32
		"18446744073709551615", // math.MaxUint64
		json.Number("123"),

		true,
		false,

		nil,
	}

	for idx, tc := range tests {
		input := fmt.Sprintf("%d:%s", idx, testingx.Format(tc))
		b.Run("caster{int8}   "+input, benchmarkIntegerCaster(ic.int8Caster.cast, tc))
		b.Run("caster{int16}  "+input, benchmarkIntegerCaster(ic.int16Caster.cast, tc))
		b.Run("caster{int32}  "+input, benchmarkIntegerCaster(ic.int32Caster.cast, tc))
		b.Run("caster{int64}  "+input, benchmarkIntegerCaster(ic.int64Caster.cast, tc))
		b.Run("caster{int}    "+input, benchmarkIntegerCaster(ic.intCaster.cast, tc))
		b.Run("caster{uint8}  "+input, benchmarkIntegerCaster(ic.uint8Caster.cast, tc))
		b.Run("caster{uint16} "+input, benchmarkIntegerCaster(ic.uint16Caster.cast, tc))
		b.Run("caster{uint32} "+input, benchmarkIntegerCaster(ic.uint32Caster.cast, tc))
		b.Run("caster{uint64} "+input, benchmarkIntegerCaster(ic.uint64Caster.cast, tc))
		b.Run("caster{uint}   "+input, benchmarkIntegerCaster(ic.uintCaster.cast, tc))
	}
}

func benchmarkIntegerCaster[T numbers.Integer](caster func(any) (T, error), input any) func(*testing.B) {
	return func(b *testing.B) {
		b.Helper()
		for i := 0; i < b.N; i++ {
			_, _ = caster(input)
		}
	}
}
