package pick

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"math"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
	_ "time/tzdata"

	"github.com/moukoublen/pick/internal/tst"
	"github.com/moukoublen/pick/internal/tst/testdata"
	"github.com/stretchr/testify/require"
)

type PickerTestCase struct {
	AccessFn      any
	Selector      string // selector or path
	ExpectedValue any
	ErrorAsserter tst.ErrorAsserter
}

func (tc *PickerTestCase) Name() string {
	return fmt.Sprintf("selector(%s)", tc.Selector)
}

func (tc *PickerTestCase) Run(t *testing.T) {
	t.Helper()
	t.Parallel()
	pickerFunctionCall := reflect.ValueOf(tc.AccessFn)

	args := []reflect.Value{reflect.ValueOf(tc.Selector)}
	returned := pickerFunctionCall.Call(args)

	got := returned[0].Interface()
	tst.AssertEqual(t, got, tc.ExpectedValue)

	var receivedError error
	if len(returned) > 1 {
		gotErr := returned[1].Interface()
		receivedError, _ = gotErr.(error)
	}
	tc.ErrorAsserter(t, receivedError)
}

func TestWithMixedTypesMap(t *testing.T) {
	t.Parallel()

	p := Wrap(testdata.MixedTypesMap)

	tests := []PickerTestCase{
		{
			AccessFn:      p.String,
			Selector:      "stringField",
			ExpectedValue: "abcd",
			ErrorAsserter: tst.NoError,
		},
		{
			AccessFn:      p.Relaxed().Int,
			Selector:      "sliceOfAnyComplex[0]",
			ExpectedValue: int(2),
			ErrorAsserter: tst.NoError,
		},
		{
			AccessFn:      p.Int8,
			Selector:      "sliceOfAnyComplex[0]",
			ExpectedValue: int8(2),
			ErrorAsserter: tst.NoError,
		},
		{
			AccessFn:      p.Int16,
			Selector:      "sliceOfAnyComplex[0]",
			ExpectedValue: int16(2),
			ErrorAsserter: tst.NoError,
		},
		{
			AccessFn:      p.Int32,
			Selector:      "sliceOfAnyComplex[4]",
			ExpectedValue: int32(5),
			ErrorAsserter: tst.NoError,
		},
		{
			AccessFn:      p.Int64,
			Selector:      "sliceOfAnyComplex[4]",
			ExpectedValue: int64(5),
			ErrorAsserter: tst.NoError,
		},
		{
			AccessFn:      p.Int64,
			Selector:      "sliceOfAnyComplex[3].key3",
			ExpectedValue: int64(6565),
			ErrorAsserter: tst.NoError,
		},
		{
			AccessFn:      p.Int32,
			Selector:      "sliceOfAnyComplex[2].A",
			ExpectedValue: int32(3),
			ErrorAsserter: tst.NoError,
		},
		{
			AccessFn:      p.Int32,
			Selector:      "sliceOfAnyComplex[2].Foo",
			ExpectedValue: int32(0),
			ErrorAsserter: tst.ExpectedErrorIs(ErrFieldNotFound),
		},
		{
			AccessFn:      p.Uint,
			Selector:      "pointerMapStringAny.fieldInt32",
			ExpectedValue: uint(6),
			ErrorAsserter: tst.NoError,
		},
		{
			AccessFn:      p.Uint8,
			Selector:      "pointerMapStringAny.fieldInt32",
			ExpectedValue: uint8(6),
			ErrorAsserter: tst.NoError,
		},
		{
			AccessFn:      p.Uint16,
			Selector:      "pointerMapStringAny.fieldInt32",
			ExpectedValue: uint16(6),
			ErrorAsserter: tst.NoError,
		},
		{
			AccessFn:      p.Uint32,
			Selector:      "pointerMapStringAny.fieldInt32",
			ExpectedValue: uint32(6),
			ErrorAsserter: tst.NoError,
		},
		{
			AccessFn:      p.Uint64,
			Selector:      "pointerMapStringAny.fieldInt32",
			ExpectedValue: uint64(6),
			ErrorAsserter: tst.NoError,
		},
		{
			AccessFn:      p.Relaxed().Uint64Slice,
			Selector:      "pointerMapStringAny.int32Slice",
			ExpectedValue: []uint64{10, 11, 12, 13, 14},
			ErrorAsserter: tst.NoError,
		},
		{
			AccessFn:      p.Int64Slice,
			Selector:      "pointerMapStringAny.int32Slice",
			ExpectedValue: []int64{10, 11, 12, 13, 14},
			ErrorAsserter: tst.NoError,
		},
		{
			AccessFn:      p.Relaxed().StringSlice,
			Selector:      "pointerMapStringAny.int32Slice",
			ExpectedValue: []string{"10", "11", "12", "13", "14"},
			ErrorAsserter: tst.NoError,
		},
		{
			AccessFn:      p.Bool,
			Selector:      "pointerMapStringAny.fieldBool",
			ExpectedValue: true,
			ErrorAsserter: tst.NoError,
		},
	}

	for idx, tc := range tests {
		name := fmt.Sprintf("%d_%s", idx, tc.Name())
		t.Run(name, tc.Run)
	}
}

// TestMixedTypesMapBool makes an extensive test in Bool/BoolSlice functions using all APIs.
func TestWithMixedTypesMapUsingBoolAPI(t *testing.T) {
	t.Parallel()

	p := Wrap(testdata.MixedTypesMap)

	tests := []PickerTestCase{
		{
			AccessFn:      p.Bool,
			Selector:      "pointerMapStringAny.fieldBool",
			ExpectedValue: true,
			ErrorAsserter: tst.NoError,
		},
		{
			AccessFn:      p.Relaxed().Bool,
			Selector:      "pointerMapStringAny.fieldBool",
			ExpectedValue: true,
			ErrorAsserter: tst.NoError,
		},

		{
			AccessFn:      p.BoolSlice,
			Selector:      "sliceOfAnyComplex[5]",
			ExpectedValue: []bool{true, true, true, false, true, true},
			ErrorAsserter: tst.NoError,
		},
		{
			AccessFn:      p.Relaxed().BoolSlice,
			Selector:      "sliceOfAnyComplex[5]",
			ExpectedValue: []bool{true, true, true, false, true, true},
			ErrorAsserter: tst.NoError,
		},
	}

	for idx, tc := range tests {
		name := fmt.Sprintf("%d_%s", idx, tc.Name())
		t.Run(name, tc.Run)
	}
}

// TestMixedTypesMapString makes an extensive test in String/StringSlice functions using all APIs.
func TestWithMixedTypesMapUsingStringAPI(t *testing.T) {
	t.Parallel()

	p := Wrap(testdata.MixedTypesMap)

	tests := []PickerTestCase{
		{
			AccessFn:      p.String,
			Selector:      "sliceOfAnyComplex[1]",
			ExpectedValue: "stringElement",
			ErrorAsserter: tst.NoError,
		},
		{
			AccessFn:      p.Relaxed().String,
			Selector:      "sliceOfAnyComplex[1]",
			ExpectedValue: "stringElement",
			ErrorAsserter: tst.NoError,
		},

		{
			AccessFn:      p.StringSlice,
			Selector:      "sliceOfAnyComplex[6]",
			ExpectedValue: []string{"abc", "def", "ghi"},
			ErrorAsserter: tst.NoError,
		},
		{
			AccessFn:      p.Relaxed().StringSlice,
			Selector:      "sliceOfAnyComplex[6]",
			ExpectedValue: []string{"abc", "def", "ghi"},
			ErrorAsserter: tst.NoError,
		},
	}

	for idx, tc := range tests {
		name := fmt.Sprintf("%d_%s", idx, tc.Name())
		t.Run(name, tc.Run)
	}
}

// TestMixedTypesMapInt64 makes an extensive test in Int64/Int64Slice functions using all APIs.
func TestWithMixedTypesMapUsingInt64(t *testing.T) {
	t.Parallel()

	p := Wrap(testdata.MixedTypesMap)

	tests := []PickerTestCase{
		{
			AccessFn:      p.Int64,
			Selector:      "sliceOfAnyComplex[0]",
			ExpectedValue: int64(2),
			ErrorAsserter: tst.NoError,
		},
		{
			AccessFn:      p.Relaxed().Int64,
			Selector:      "sliceOfAnyComplex[0]",
			ExpectedValue: int64(2),
			ErrorAsserter: tst.NoError,
		},

		{
			AccessFn:      p.Int64Slice,
			Selector:      "pointerMapStringAny.int32Slice",
			ExpectedValue: []int64{10, 11, 12, 13, 14},
			ErrorAsserter: tst.NoError,
		},
		{
			AccessFn:      p.Relaxed().Int64Slice,
			Selector:      "pointerMapStringAny.int32Slice",
			ExpectedValue: []int64{10, 11, 12, 13, 14},
			ErrorAsserter: tst.NoError,
		},
	}

	for idx, tc := range tests {
		name := fmt.Sprintf("%d_%s", idx, tc.Name())
		t.Run(name, tc.Run)
	}
}

// TestMixedTypesMapFloat64 makes an extensive test in Float64/Float64Slice functions using all APIs.
func TestUsingMixedTypesMapUsingFloat64API(t *testing.T) {
	t.Parallel()

	p := Wrap(testdata.MixedTypesMap)

	tests := []PickerTestCase{
		{
			AccessFn:      p.Float64,
			Selector:      "pointerMapStringAny.float64Slice[3]",
			ExpectedValue: float64(0.4),
			ErrorAsserter: tst.NoError,
		},
		{
			AccessFn:      p.Relaxed().Float64,
			Selector:      "pointerMapStringAny.float64Slice[3]",
			ExpectedValue: float64(0.4),
			ErrorAsserter: tst.NoError,
		},

		{
			AccessFn:      p.Float64Slice,
			Selector:      "pointerMapStringAny.float64Slice",
			ExpectedValue: []float64{0.1, 0.2, 0.3, 0.4},
			ErrorAsserter: tst.NoError,
		},
		{
			AccessFn:      p.Relaxed().Float64Slice,
			Selector:      "pointerMapStringAny.float64Slice",
			ExpectedValue: []float64{0.1, 0.2, 0.3, 0.4},
			ErrorAsserter: tst.NoError,
		},
	}

	for idx, tc := range tests {
		name := fmt.Sprintf("%d_%s", idx, tc.Name())
		t.Run(name, tc.Run)
	}
}

// TestMixedTypesMapFloat32 makes an extensive test in Float32/Float32Slice functions using all APIs.
func TestWithMixedTypesMapUsingFloat32API(t *testing.T) {
	t.Parallel()

	p := Wrap(testdata.MixedTypesMap)

	tests := []PickerTestCase{
		{
			AccessFn:      p.Float32,
			Selector:      "pointerMapStringAny.float64Slice[3]",
			ExpectedValue: float32(0.4),
			ErrorAsserter: tst.NoError,
		},
		{
			AccessFn:      p.Relaxed().Float32,
			Selector:      "pointerMapStringAny.float64Slice[3]",
			ExpectedValue: float32(0.4),
			ErrorAsserter: tst.NoError,
		},

		{
			AccessFn:      p.Float32Slice,
			Selector:      "pointerMapStringAny.float64Slice",
			ExpectedValue: []float32{0.1, 0.2, 0.3, 0.4},
			ErrorAsserter: tst.NoError,
		},
		{
			AccessFn:      p.Relaxed().Float32Slice,
			Selector:      "pointerMapStringAny.float64Slice",
			ExpectedValue: []float32{0.1, 0.2, 0.3, 0.4},
			ErrorAsserter: tst.NoError,
		},
	}

	for idx, tc := range tests {
		name := fmt.Sprintf("%d_%s", idx, tc.Name())
		t.Run(name, tc.Run)
	}
}

// TestMixedTypesMapByte makes an extensive test in Byte/ByteSlice functions using all APIs.
func TestWithMixedTypesMapUsingByteAPI(t *testing.T) {
	t.Parallel()

	p := Wrap(testdata.MixedTypesMap)

	tests := []PickerTestCase{
		{
			AccessFn:      p.Byte,
			Selector:      "sliceOfAnyComplex[7]",
			ExpectedValue: byte(math.MaxUint8),
			ErrorAsserter: tst.NoError,
		},
		{
			AccessFn:      p.Relaxed().Byte,
			Selector:      "sliceOfAnyComplex[7]",
			ExpectedValue: byte(math.MaxUint8),
			ErrorAsserter: tst.NoError,
		},

		{
			AccessFn:      p.ByteSlice,
			Selector:      "pointerMapStringAny.int32Slice",
			ExpectedValue: []byte{10, 11, 12, 13, 14},
			ErrorAsserter: tst.NoError,
		},
		{
			AccessFn:      p.Relaxed().ByteSlice,
			Selector:      "pointerMapStringAny.int32Slice",
			ExpectedValue: []byte{10, 11, 12, 13, 14},
			ErrorAsserter: tst.NoError,
		},
	}

	for idx, tc := range tests {
		name := fmt.Sprintf("%d_%s", idx, tc.Name())
		t.Run(name, tc.Run)
	}
}

// TestWithMixedTypesMapUsingTimeAPI makes an extensive test in Time/TimeSlice functions using all APIs.
func TestWithMixedTypesMapUsingTimeAPI(t *testing.T) {
	t.Parallel()

	p := Wrap(testdata.MixedTypesMap)

	tests := []PickerTestCase{
		{
			AccessFn:      p.Time,
			Selector:      "times.timeRFC3339Nano",
			ExpectedValue: time.Date(1977, time.May, 25, 22, 30, 0, 0, time.UTC),
			ErrorAsserter: tst.NoError,
		},
		{
			AccessFn:      p.Relaxed().Time,
			Selector:      "times.timeRFC3339Nano",
			ExpectedValue: time.Date(1977, time.May, 25, 22, 30, 0, 0, time.UTC),
			ErrorAsserter: tst.NoError,
		},

		{
			AccessFn: p.TimeSlice,
			Selector: "times.timeUnixSecondsSlice",
			ExpectedValue: []time.Time{
				time.Date(1977, time.May, 25, 18, 30, 0, 0, time.UTC),
				time.Date(1977, time.May, 25, 18, 30, 1, 0, time.UTC),
				time.Date(1977, time.May, 25, 18, 30, 2, 0, time.UTC),
			},
			ErrorAsserter: tst.NoError,
		},
		{
			AccessFn: p.Relaxed().TimeSlice,
			Selector: "times.timeUnixSecondsSlice",
			ExpectedValue: []time.Time{
				time.Date(1977, time.May, 25, 18, 30, 0, 0, time.UTC),
				time.Date(1977, time.May, 25, 18, 30, 1, 0, time.UTC),
				time.Date(1977, time.May, 25, 18, 30, 2, 0, time.UTC),
			},
			ErrorAsserter: tst.NoError,
		},
	}

	for idx, tc := range tests {
		name := fmt.Sprintf("%d_%s", idx, tc.Name())
		t.Run(name, tc.Run)
	}
}

// TestWithMixedTypesMapUsingDurationAPI makes an extensive test in Duration/DurationSlice functions using all APIs.
func TestWithMixedTypesMapUsingDurationAPI(t *testing.T) {
	t.Parallel()

	p := Wrap(testdata.MixedTypesMap)

	tests := []PickerTestCase{
		{
			AccessFn:      p.Duration,
			Selector:      "durations.single",
			ExpectedValue: time.Duration(4) * time.Second,
			ErrorAsserter: tst.NoError,
		},
		{
			AccessFn:      p.Relaxed().Duration,
			Selector:      "durations.single",
			ExpectedValue: time.Duration(4) * time.Second,
			ErrorAsserter: tst.NoError,
		},

		{
			AccessFn:      p.DurationSlice,
			Selector:      "durations.slice",
			ExpectedValue: []time.Duration{5 * time.Second, 6 * time.Second, 7 * time.Second},
			ErrorAsserter: tst.NoError,
		},
		{
			AccessFn:      p.Relaxed().DurationSlice,
			Selector:      "durations.slice",
			ExpectedValue: []time.Duration{5 * time.Second, 6 * time.Second, 7 * time.Second},
			ErrorAsserter: tst.NoError,
		},
	}

	for idx, tc := range tests {
		name := fmt.Sprintf("%d_%s", idx, tc.Name())
		t.Run(name, tc.Run)
	}
}

//go:embed internal/tst/testdata
var testData embed.FS

func loadTestData(t *testing.T, filename string) fs.File {
	t.Helper()

	path := filepath.Join("internal", "tst", "testdata", filename)
	f, err := testData.Open(path)
	if err != nil {
		t.Fatalf("error during testdate file opening %s", err.Error())
	}

	return f
}

func TestNasaDataFile(t *testing.T) {
	t.Parallel()

	file := loadTestData(t, "nasa.json")

	timeFormat1 := "2006-01-02"
	timeFormat2 := "2006-Jan-02 15:04"

	p, err := WrapReaderJSON(file)
	if err != nil {
		t.Fatal(err)
	}

	tests := []PickerTestCase{
		{
			AccessFn:      p.String,
			Selector:      "near_earth_objects.2023-01-01[4].neo_reference_id",
			ExpectedValue: "3703782",
			ErrorAsserter: tst.NoError,
		},
		{
			AccessFn:      p.Float32,
			Selector:      "near_earth_objects.2023-01-01[5].estimated_diameter.meters.estimated_diameter_max",
			ExpectedValue: float32(68.2401509401),
			ErrorAsserter: tst.NoError,
		},
		{
			AccessFn:      p.Float64,
			Selector:      "near_earth_objects.2023-01-01[5].estimated_diameter.meters.estimated_diameter_max",
			ExpectedValue: float64(68.2401509401),
			ErrorAsserter: tst.NoError,
		},
		{
			AccessFn:      p.Uint8,
			Selector:      "near_earth_objects.2023-01-01[5].id",
			ExpectedValue: uint8(214),
			ErrorAsserter: tst.ExpectedErrorIs(ErrConvertOverFlow),
		},
		{
			AccessFn:      p.Len,
			Selector:      "near_earth_objects.2023-01-01",
			ExpectedValue: 23,
			ErrorAsserter: tst.NoError,
		},
		{
			AccessFn:      p.Uint16,
			Selector:      "near_earth_objects.2023-01-01[5].id",
			ExpectedValue: uint16(50902),
			ErrorAsserter: tst.ExpectedErrorIs(ErrConvertOverFlow),
		},
		{
			AccessFn:      p.Uint32,
			Selector:      "near_earth_objects.2023-01-01[5].id",
			ExpectedValue: uint32(3720918),
			ErrorAsserter: tst.NoError,
		},
		{
			AccessFn:      p.Uint64,
			Selector:      "near_earth_objects.2023-01-01[5].id",
			ExpectedValue: uint64(3720918),
			ErrorAsserter: tst.NoError,
		},
		{
			AccessFn:      p.Uint,
			Selector:      "near_earth_objects.2023-01-01[5].id",
			ExpectedValue: uint(3720918),
			ErrorAsserter: tst.NoError,
		},
		{
			AccessFn:      p.Bool,
			Selector:      "near_earth_objects.2023-01-01[1].is_potentially_hazardous_asteroid",
			ExpectedValue: true,
			ErrorAsserter: tst.NoError,
		},
		{
			AccessFn: func(selector string) (time.Time, error) {
				return p.TimeWithConfig(TimeConvertConfig{StringFormat: timeFormat1}, selector)
			},
			Selector:      "near_earth_objects.2023-01-01[1].close_approach_data[0].close_approach_date",
			ExpectedValue: time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
			ErrorAsserter: tst.NoError,
		},
		{
			AccessFn: func(selector string) (time.Time, error) {
				return p.TimeWithConfig(TimeConvertConfig{StringFormat: timeFormat2}, selector)
			},
			Selector:      "near_earth_objects.2023-01-01[1].close_approach_data[0].close_approach_date_full",
			ExpectedValue: time.Date(2023, time.January, 1, 19, 45, 0, 0, time.UTC),
			ErrorAsserter: tst.NoError,
		},
		{
			AccessFn: func(selector string) ([]string, error) {
				return Map(p, selector, func(p Picker) (string, error) { return p.String("id") })
			},
			Selector:      "near_earth_objects.2023-01-01",
			ExpectedValue: []string{"2154347", "2385186", "2453309", "3683468", "3703782", "3720918", "3767936", "3792438", "3824981", "3836251", "3837605", "3959234", "3986848", "54104550", "54105994", "54166175", "54202993", "54290862", "54335607", "54337027", "54337425", "54340039", "54341664"},
			ErrorAsserter: tst.NoError,
		},
		{
			AccessFn: func(selector string) []string {
				return RelaxedMap(p.Relaxed(), selector, func(a RelaxedAPI) (string, error) { return a.String("id"), nil })
			},
			Selector:      "near_earth_objects.2023-01-01",
			ExpectedValue: []string{"2154347", "2385186", "2453309", "3683468", "3703782", "3720918", "3767936", "3792438", "3824981", "3836251", "3837605", "3959234", "3986848", "54104550", "54105994", "54166175", "54202993", "54290862", "54335607", "54337027", "54337425", "54340039", "54341664"},
			ErrorAsserter: tst.NoError,
		},
		{
			AccessFn: func(_ string) ([]string, error) {
				return FlatMap(p, "near_earth_objects.2023-01-01", func(p Picker) ([]string, error) {
					return Map(p, "close_approach_data", func(p Picker) (string, error) {
						return p.String("close_approach_date_full")
					})
				})
			},
			Selector:      "",
			ExpectedValue: []string{"2023-Jan-01 18:44", "2023-Jan-01 19:45", "2023-Jan-01 20:20", "2023-Jan-01 13:38", "2023-Jan-01 00:59", "2023-Jan-01 17:33", "2023-Jan-01 09:38", "2023-Jan-01 09:49", "2023-Jan-01 03:04", "2023-Jan-01 22:31", "2023-Jan-01 04:15", "2023-Jan-01 02:10", "2023-Jan-01 10:47", "2023-Jan-01 16:46", "2023-Jan-01 12:02", "2023-Jan-01 16:03", "2023-Jan-01 13:39", "2023-Jan-01 12:50", "2023-Jan-01 20:45", "2023-Jan-01 07:16", "2023-Jan-01 01:15", "2023-Jan-01 23:21", "2023-Jan-01 09:02"},
			ErrorAsserter: tst.NoError,
		},
		{
			AccessFn: func(_ string) []string {
				return RelaxedFlatMap(p.Relaxed(), "near_earth_objects.2023-01-01", func(a RelaxedAPI) ([]string, error) {
					return RelaxedMap(a, "close_approach_data", func(a RelaxedAPI) (string, error) {
						return a.String("close_approach_date_full"), nil
					}), nil
				})
			},
			Selector:      "",
			ExpectedValue: []string{"2023-Jan-01 18:44", "2023-Jan-01 19:45", "2023-Jan-01 20:20", "2023-Jan-01 13:38", "2023-Jan-01 00:59", "2023-Jan-01 17:33", "2023-Jan-01 09:38", "2023-Jan-01 09:49", "2023-Jan-01 03:04", "2023-Jan-01 22:31", "2023-Jan-01 04:15", "2023-Jan-01 02:10", "2023-Jan-01 10:47", "2023-Jan-01 16:46", "2023-Jan-01 12:02", "2023-Jan-01 16:03", "2023-Jan-01 13:39", "2023-Jan-01 12:50", "2023-Jan-01 20:45", "2023-Jan-01 07:16", "2023-Jan-01 01:15", "2023-Jan-01 23:21", "2023-Jan-01 09:02"},
			ErrorAsserter: tst.NoError,
		},
	}

	for idx, tc := range tests {
		name := fmt.Sprintf("%d_%s", idx, tc.Name())
		t.Run(name, tc.Run)
	}
}

func TestMapRelaxed(t *testing.T) {
	t.Parallel()
	file := loadTestData(t, "nasa.json")
	p, pErr := WrapReaderJSON(file)
	if pErr != nil {
		t.Fatal(pErr)
	}

	type Item struct {
		Name   string
		Sentry bool
	}

	t.Run("happy path", func(t *testing.T) {
		t.Parallel()
		errSink := &ErrorsSink{}
		itemsSlice := RelaxedMap(p.Relaxed(errSink), "near_earth_objects.2023-01-07", func(sm RelaxedAPI) (Item, error) {
			return Item{
				Name:   sm.String("name"),
				Sentry: sm.Bool("is_sentry_object"),
			}, nil
		})
		tst.AssertEqual(t, errSink.Outcome(), nil)
		tst.AssertEqual(t, itemsSlice, []Item{
			{Name: "344133 (2000 AD6)", Sentry: false},
			{Name: "369454 (2010 NZ1)", Sentry: false},
			{Name: "452334 (2001 LB)", Sentry: false},
			{Name: "(2003 MK4)", Sentry: false},
			{Name: "(2006 HJ18)", Sentry: false},
			{Name: "(2008 BX2)", Sentry: false},
			{Name: "(2013 QM10)", Sentry: false},
			{Name: "(2013 YL2)", Sentry: false},
			{Name: "(2016 AE166)", Sentry: false},
			{Name: "(2019 AR3)", Sentry: false},
			{Name: "(2019 AJ8)", Sentry: false},
			{Name: "(2021 RH)", Sentry: false},
			{Name: "(2022 BB3)", Sentry: false},
			{Name: "(2022 OF)", Sentry: false},
			{Name: "(2022 YD6)", Sentry: false},
			{Name: "(2023 AF)", Sentry: false},
			{Name: "(2023 AG1)", Sentry: false},
		})
	})

	t.Run("gather errors", func(t *testing.T) {
		t.Parallel()
		errSink := &ErrorsSink{}
		_ = RelaxedMap(p.Relaxed(errSink), "near_earth_objects.2023-01-07", func(sm RelaxedAPI) (Item, error) {
			return Item{
				Name:   sm.String("name"),
				Sentry: sm.Bool("wrong.path"),
			}, nil
		})

		var g *multiError
		tst.AssertEqual(t, errors.As(errSink.Outcome(), &g), true)
		tst.AssertEqual(t, len(g.errors), 17)
		for _, e := range g.errors {
			tst.ExpectedErrorIs(ErrFieldNotFound)(t, e)
		}
	})
}

func TestEach(t *testing.T) {
	t.Parallel()
	file := loadTestData(t, "nasa.json")
	p, pErr := WrapReaderJSON(file)
	if pErr != nil {
		t.Fatal(pErr)
	}

	t.Run("Each happy path", func(t *testing.T) {
		t.Parallel()
		err := Each(p, "near_earth_objects.2023-01-07", func(index int, p Picker, length int) error {
			tst.AssertEqual(t, length, 17)
			if index == 4 {
				s, err := p.String("name")
				tst.AssertEqual(t, s, "(2006 HJ18)")
				tst.AssertEqual(t, err, nil)
			}
			return nil
		})
		tst.AssertEqual(t, err, nil)
	})

	t.Run("Each error", func(t *testing.T) {
		t.Parallel()
		mockErr := errors.New("error")
		err := Each(p, "near_earth_objects.2023-01-07", func(index int, _ Picker, length int) error {
			tst.AssertEqual(t, length, 17)
			if index == 4 {
				return mockErr
			}
			return nil
		})
		tst.ExpectedErrorIs(mockErr)(t, err)
	})

	t.Run("EachM happy path", func(t *testing.T) {
		t.Parallel()
		errSink := &ErrorsSink{}
		RelaxedEach(p.Relaxed(errSink), "near_earth_objects.2023-01-07", func(index int, a RelaxedAPI, length int) error {
			tst.AssertEqual(t, length, 17)
			if index == 4 {
				s := a.String("name")
				tst.AssertEqual(t, s, "(2006 HJ18)")
			}
			return nil
		})
		tst.AssertEqual(t, errSink.Outcome(), nil)
	})

	t.Run("EachM error", func(t *testing.T) {
		t.Parallel()
		errSink := &ErrorsSink{}
		RelaxedEach(p.Relaxed(errSink), "near_earth_objects.2023-01-07", func(index int, a RelaxedAPI, length int) error {
			tst.AssertEqual(t, length, 17)
			if index == 4 {
				s := a.String("name")
				tst.AssertEqual(t, s, "(2006 HJ18)")
			}
			return nil
		})
		tst.AssertEqual(t, errSink.Outcome(), nil)
	})
}

func TestHTTP(t *testing.T) {
	t.Run("request", func(t *testing.T) {
		b := strings.NewReader(`{"one": 1}`)
		r, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "http://localhost", b)
		require.NoError(t, err)

		p, err := WrapJSONRequest(r)
		require.NoError(t, err)
		require.Equal(t, 1, p.Relaxed().Int("one"))
	})

	t.Run("request nil body", func(t *testing.T) {
		r, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "http://localhost", http.NoBody)
		require.NoError(t, err)

		p, err := WrapJSONRequest(r)
		require.NoError(t, err)
		require.Nil(t, p.Data())
	})

	t.Run("response", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		recorder.WriteHeader(http.StatusOK)
		_, _ = recorder.WriteString(`{"one": "two"}`)

		p, err := WrapJSONResponse(recorder.Result())
		require.NoError(t, err)
		require.Equal(t, "two", p.Relaxed().String("one"))
	})

	t.Run("response nil body", func(t *testing.T) {
		p, err := WrapJSONResponse(&http.Response{Body: http.NoBody})
		require.NoError(t, err)
		require.Nil(t, p.Data())
	})
}
