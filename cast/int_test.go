package cast

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"runtime"
	"testing"

	"github.com/moukoublen/pick/internal/testingx"
	"golang.org/x/exp/constraints"
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
		to := to
		for _, tc := range perKindTestCases {
			tc := tc
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
		to := to
		for _, tc := range perKindTestCases {
			tc := tc
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

type intCastExpectedResults[T constraints.Integer] struct {
	result  T
	errorFn func(*testing.T, error)
}

// constructors for shortage.
func newIntCastExpectedResultsFn[T constraints.Integer]() func(result T, errorFn func(*testing.T, error)) intCastExpectedResults[T] {
	return func(result T, errorFn func(*testing.T, error)) intCastExpectedResults[T] {
		return intCastExpectedResults[T]{result: result, errorFn: errorFn}
	}
}

func TestIntegerCaster(t *testing.T) {
	t.Parallel()

	i08 := newIntCastExpectedResultsFn[int8]()
	i16 := newIntCastExpectedResultsFn[int16]()
	i32 := newIntCastExpectedResultsFn[int32]()
	i64 := newIntCastExpectedResultsFn[int64]()
	u08 := newIntCastExpectedResultsFn[uint8]()
	u16 := newIntCastExpectedResultsFn[uint16]()
	u32 := newIntCastExpectedResultsFn[uint32]()
	u64 := newIntCastExpectedResultsFn[uint64]()

	type TestCase struct {
		input          any
		ExpectedInt8   intCastExpectedResults[int8]
		ExpectedInt16  intCastExpectedResults[int16]
		ExpectedInt32  intCastExpectedResults[int32]
		ExpectedInt64  intCastExpectedResults[int64]
		ExpectedUint8  intCastExpectedResults[uint8]
		ExpectedUint16 intCastExpectedResults[uint16]
		ExpectedUint32 intCastExpectedResults[uint32]
		ExpectedUint64 intCastExpectedResults[uint64]
	}

	tests := []TestCase{
		{
			input:          int8(12),
			ExpectedInt8:   i08(12, nil),
			ExpectedInt16:  i16(12, nil),
			ExpectedInt32:  i32(12, nil),
			ExpectedInt64:  i64(12, nil),
			ExpectedUint8:  u08(12, nil),
			ExpectedUint16: u16(12, nil),
			ExpectedUint32: u32(12, nil),
			ExpectedUint64: u64(12, nil),
		},
		{
			input:          int8(math.MaxInt8),
			ExpectedInt8:   i08(127, nil),
			ExpectedInt16:  i16(127, nil),
			ExpectedInt32:  i32(127, nil),
			ExpectedInt64:  i64(127, nil),
			ExpectedUint8:  u08(127, nil),
			ExpectedUint16: u16(127, nil),
			ExpectedUint32: u32(127, nil),
			ExpectedUint64: u64(127, nil),
		},
		{
			input:          int16(128),
			ExpectedInt8:   i08(-128, expectOverFlowError),
			ExpectedInt16:  i16(128, nil),
			ExpectedInt32:  i32(128, nil),
			ExpectedInt64:  i64(128, nil),
			ExpectedUint8:  u08(128, nil),
			ExpectedUint16: u16(128, nil),
			ExpectedUint32: u32(128, nil),
			ExpectedUint64: u64(128, nil),
		},
		{
			input:          uint8(math.MaxUint8),
			ExpectedInt8:   i08(-1, expectOverFlowError),
			ExpectedInt16:  i16(255, nil),
			ExpectedInt32:  i32(255, nil),
			ExpectedInt64:  i64(255, nil),
			ExpectedUint8:  u08(255, nil),
			ExpectedUint16: u16(255, nil),
			ExpectedUint32: u32(255, nil),
			ExpectedUint64: u64(255, nil),
		},
		{
			input:          uint64(math.MaxUint64),
			ExpectedInt8:   i08(-1, expectOverFlowError),
			ExpectedInt16:  i16(-1, expectOverFlowError),
			ExpectedInt32:  i32(-1, expectOverFlowError),
			ExpectedInt64:  i64(-1, expectOverFlowError),
			ExpectedUint8:  u08(math.MaxUint8, expectOverFlowError),
			ExpectedUint16: u16(math.MaxUint16, expectOverFlowError),
			ExpectedUint32: u32(math.MaxUint32, expectOverFlowError),
			ExpectedUint64: u64(math.MaxUint64, nil),
		},

		// float
		{
			input:          float32(123),
			ExpectedInt8:   i08(123, nil),
			ExpectedInt16:  i16(123, nil),
			ExpectedInt32:  i32(123, nil),
			ExpectedInt64:  i64(123, nil),
			ExpectedUint8:  u08(123, nil),
			ExpectedUint16: u16(123, nil),
			ExpectedUint32: u32(123, nil),
			ExpectedUint64: u64(123, nil),
		},
		{
			input:          float64(123),
			ExpectedInt8:   i08(123, nil),
			ExpectedInt16:  i16(123, nil),
			ExpectedInt32:  i32(123, nil),
			ExpectedInt64:  i64(123, nil),
			ExpectedUint8:  u08(123, nil),
			ExpectedUint16: u16(123, nil),
			ExpectedUint32: u32(123, nil),
			ExpectedUint64: u64(123, nil),
		},
		{
			input:          float32(123.001),
			ExpectedInt8:   i08(123, expectLostDecimals),
			ExpectedInt16:  i16(123, expectLostDecimals),
			ExpectedInt32:  i32(123, expectLostDecimals),
			ExpectedInt64:  i64(123, expectLostDecimals),
			ExpectedUint8:  u08(123, expectLostDecimals),
			ExpectedUint16: u16(123, expectLostDecimals),
			ExpectedUint32: u32(123, expectLostDecimals),
			ExpectedUint64: u64(123, expectLostDecimals),
		},
		{
			input:          float64(123.000001),
			ExpectedInt8:   i08(123, expectLostDecimals),
			ExpectedInt16:  i16(123, expectLostDecimals),
			ExpectedInt32:  i32(123, expectLostDecimals),
			ExpectedInt64:  i64(123, expectLostDecimals),
			ExpectedUint8:  u08(123, expectLostDecimals),
			ExpectedUint16: u16(123, expectLostDecimals),
			ExpectedUint32: u32(123, expectLostDecimals),
			ExpectedUint64: u64(123, expectLostDecimals),
		},
		{
			input:          float64(123.9999999),
			ExpectedInt8:   i08(123, expectLostDecimals),
			ExpectedInt16:  i16(123, expectLostDecimals),
			ExpectedInt32:  i32(123, expectLostDecimals),
			ExpectedInt64:  i64(123, expectLostDecimals),
			ExpectedUint8:  u08(123, expectLostDecimals),
			ExpectedUint16: u16(123, expectLostDecimals),
			ExpectedUint32: u32(123, expectLostDecimals),
			ExpectedUint64: u64(123, expectLostDecimals),
		},

		// string
		{
			input:          string("123"),
			ExpectedInt8:   i08(123, nil),
			ExpectedInt16:  i16(123, nil),
			ExpectedInt32:  i32(123, nil),
			ExpectedInt64:  i64(123, nil),
			ExpectedUint8:  u08(123, nil),
			ExpectedUint16: u16(123, nil),
			ExpectedUint32: u32(123, nil),
			ExpectedUint64: u64(123, nil),
		},
		{
			input:          string("123.123"),
			ExpectedInt8:   i08(123, expectLostDecimals),
			ExpectedInt16:  i16(123, expectLostDecimals),
			ExpectedInt32:  i32(123, expectLostDecimals),
			ExpectedInt64:  i64(123, expectLostDecimals),
			ExpectedUint8:  u08(123, expectLostDecimals),
			ExpectedUint16: u16(123, expectLostDecimals),
			ExpectedUint32: u32(123, expectLostDecimals),
			ExpectedUint64: u64(123, expectLostDecimals),
		},
		{
			input:          string("bad input"),
			ExpectedInt8:   i08(0, expectMalformedSyntax),
			ExpectedInt16:  i16(0, expectMalformedSyntax),
			ExpectedInt32:  i32(0, expectMalformedSyntax),
			ExpectedInt64:  i64(0, expectMalformedSyntax),
			ExpectedUint8:  u08(0, expectMalformedSyntax),
			ExpectedUint16: u16(0, expectMalformedSyntax),
			ExpectedUint32: u32(0, expectMalformedSyntax),
			ExpectedUint64: u64(0, expectMalformedSyntax),
		},

		// json number
		{
			input:          json.Number("123"),
			ExpectedInt8:   i08(123, nil),
			ExpectedInt16:  i16(123, nil),
			ExpectedInt32:  i32(123, nil),
			ExpectedInt64:  i64(123, nil),
			ExpectedUint8:  u08(123, nil),
			ExpectedUint16: u16(123, nil),
			ExpectedUint32: u32(123, nil),
			ExpectedUint64: u64(123, nil),
		},
		{
			input:          json.Number("56782"),
			ExpectedInt8:   i08(-50, expectOverFlowError),
			ExpectedInt16:  i16(-8754, expectOverFlowError),
			ExpectedInt32:  i32(56782, nil),
			ExpectedInt64:  i64(56782, nil),
			ExpectedUint8:  u08(206, expectOverFlowError),
			ExpectedUint16: u16(56782, nil),
			ExpectedUint32: u32(56782, nil),
			ExpectedUint64: u64(56782, nil),
		},

		// bool
		{
			input:          false,
			ExpectedInt8:   i08(0, nil),
			ExpectedInt16:  i16(0, nil),
			ExpectedInt32:  i32(0, nil),
			ExpectedInt64:  i64(0, nil),
			ExpectedUint8:  u08(0, nil),
			ExpectedUint16: u16(0, nil),
			ExpectedUint32: u32(0, nil),
			ExpectedUint64: u64(0, nil),
		},
		{
			input:          true,
			ExpectedInt8:   i08(1, nil),
			ExpectedInt16:  i16(1, nil),
			ExpectedInt32:  i32(1, nil),
			ExpectedInt64:  i64(1, nil),
			ExpectedUint8:  u08(1, nil),
			ExpectedUint16: u16(1, nil),
			ExpectedUint32: u32(1, nil),
			ExpectedUint64: u64(1, nil),
		},

		// nil
		{
			input:          nil,
			ExpectedInt8:   i08(0, nil),
			ExpectedInt16:  i16(0, nil),
			ExpectedInt32:  i32(0, nil),
			ExpectedInt64:  i64(0, nil),
			ExpectedUint8:  u08(0, nil),
			ExpectedUint16: u16(0, nil),
			ExpectedUint32: u32(0, nil),
			ExpectedUint64: u64(0, nil),
		},

		// Unknown
		{
			input:          struct{}{},
			ExpectedInt8:   i08(0, expectInvalidType),
			ExpectedInt16:  i16(0, expectInvalidType),
			ExpectedInt32:  i32(0, expectInvalidType),
			ExpectedInt64:  i64(0, expectInvalidType),
			ExpectedUint8:  u08(0, expectInvalidType),
			ExpectedUint16: u16(0, expectInvalidType),
			ExpectedUint32: u32(0, expectInvalidType),
			ExpectedUint64: u64(0, expectInvalidType),
		},
	}

	ic := newIntegerCaster()

	for _, testCase := range tests {
		tc := testCase
		typeName := "nil"
		if tc.input != nil {
			tp := reflect.TypeOf(tc.input)
			typeName = tp.Name()
		}
		name := fmt.Sprintf("%s(%v)", typeName, testCase.input)
		t.Run(
			name,
			func(t *testing.T) {
				t.Parallel()
				t.Run("caster[int8]", testIntegerCast(ic.int8Caster, tc.input, tc.ExpectedInt8))
				t.Run("caster[int16]", testIntegerCast(ic.int16Caster, tc.input, tc.ExpectedInt16))
				t.Run("caster[int32]", testIntegerCast(ic.int32Caster, tc.input, tc.ExpectedInt32))
				t.Run("caster[int64]", testIntegerCast(ic.int64Caster, tc.input, tc.ExpectedInt64))
				t.Run("caster[uint8]", testIntegerCast(ic.uint8Caster, tc.input, tc.ExpectedUint8))
				t.Run("caster[uint16]", testIntegerCast(ic.uint16Caster, tc.input, tc.ExpectedUint16))
				t.Run("caster[uint32]", testIntegerCast(ic.uint32Caster, tc.input, tc.ExpectedUint32))
				t.Run("caster[uint64]", testIntegerCast(ic.uint64Caster, tc.input, tc.ExpectedUint64))
			},
		)
	}
}

func testIntegerCast[T constraints.Integer](caster intCast[T], input any, castExpectedResults intCastExpectedResults[T]) func(t *testing.T) {
	return func(t *testing.T) {
		gotResult, gotError := caster.cast(input)
		testingx.AssertError(t, castExpectedResults.errorFn, gotError)
		if gotResult != castExpectedResults.result {
			t.Errorf("wrong returned value. Expected %d got %d", castExpectedResults.result, gotResult)
		}
	}
}

func BenchmarkIntCaster(b *testing.B) {
	ic := newIntegerCaster()

	tests := []any{
		int8(123),
		int16(123),
		int32(123),
		int64(123),
		int(123),

		uint8(8),
		uint16(8),
		uint32(8),
		uint64(8),
		uint(8),

		float32(123),
		float64(123),

		"123",
		json.Number("123"),

		true,
		false,

		nil,
	}

	for idx, tc := range tests {
		tc := tc

		typeOfTC := reflect.TypeOf(tc)
		name := "nil"
		if tc != nil {
			name = typeOfTC.String()
		}
		name = fmt.Sprintf("test_%d_(%s)", idx, name)

		b.Run(fmt.Sprintf("int8_caster_%s", name), benchmarkIntegerCaster(ic.int8Caster.cast, tc))
		b.Run(fmt.Sprintf("int16_caster_%s", name), benchmarkIntegerCaster(ic.int16Caster.cast, tc))
		b.Run(fmt.Sprintf("int32_caster_%s", name), benchmarkIntegerCaster(ic.int32Caster.cast, tc))
		b.Run(fmt.Sprintf("int64_caster_%s", name), benchmarkIntegerCaster(ic.int64Caster.cast, tc))
		b.Run(fmt.Sprintf("int_caster_%s", name), benchmarkIntegerCaster(ic.intCaster.cast, tc))
		b.Run(fmt.Sprintf("uint8_caster_%s", name), benchmarkIntegerCaster(ic.uint8Caster.cast, tc))
		b.Run(fmt.Sprintf("uint16_caster_%s", name), benchmarkIntegerCaster(ic.uint16Caster.cast, tc))
		b.Run(fmt.Sprintf("uint32_caster_%s", name), benchmarkIntegerCaster(ic.uint32Caster.cast, tc))
		b.Run(fmt.Sprintf("uint64_caster_%s", name), benchmarkIntegerCaster(ic.uint64Caster.cast, tc))
		b.Run(fmt.Sprintf("uint_caster_%s", name), benchmarkIntegerCaster(ic.uintCaster.cast, tc))
	}
}

func benchmarkIntegerCaster[T constraints.Integer](caster func(any) (T, error), input any) func(*testing.B) {
	return func(b *testing.B) {
		b.Helper()
		for i := 0; i < b.N; i++ {
			_, _ = caster(input)
		}
	}
}
