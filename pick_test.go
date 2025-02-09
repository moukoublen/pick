package pick

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"math"
	"path/filepath"
	"reflect"
	"testing"
	"time"
	_ "time/tzdata"

	"github.com/moukoublen/pick/internal/tst"
	"github.com/moukoublen/pick/internal/tst/testdata"
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
			AccessFn:      p.Must().Int,
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
			AccessFn:      p.Must().Uint64Slice,
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
			AccessFn:      p.Must().StringSlice,
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
			AccessFn:      p.Must().Bool,
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
			AccessFn:      p.Must().BoolSlice,
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
			AccessFn:      p.Must().String,
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
			AccessFn:      p.Must().StringSlice,
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
			AccessFn:      p.Must().Int64,
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
			AccessFn:      p.Must().Int64Slice,
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
			AccessFn:      p.Must().Float64,
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
			AccessFn:      p.Must().Float64Slice,
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
			AccessFn:      p.Must().Float32,
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
			AccessFn:      p.Must().Float32Slice,
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
			AccessFn:      p.Must().Byte,
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
			AccessFn:      p.Must().ByteSlice,
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
			AccessFn:      p.Must().Time,
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
			AccessFn: p.Must().TimeSlice,
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
			AccessFn:      p.Must().Duration,
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
			AccessFn:      p.Must().DurationSlice,
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
			ErrorAsserter: tst.ExpectedErrorIs(ErrCastOverFlow),
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
			ErrorAsserter: tst.ExpectedErrorIs(ErrCastOverFlow),
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
				return p.TimeWithConfig(TimeCastConfig{StringFormat: timeFormat1}, selector)
			},
			Selector:      "near_earth_objects.2023-01-01[1].close_approach_data[0].close_approach_date",
			ExpectedValue: time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
			ErrorAsserter: tst.NoError,
		},
		{
			AccessFn: func(selector string) (time.Time, error) {
				return p.TimeWithConfig(TimeCastConfig{StringFormat: timeFormat2}, selector)
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
				return MustMap(p.Must(), selector, func(a SelectorMustAPI) (string, error) { return a.String("id"), nil })
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
				return MustFlatMap(p.Must(), "near_earth_objects.2023-01-01", func(a SelectorMustAPI) ([]string, error) {
					return MustMap(a, "close_approach_data", func(a SelectorMustAPI) (string, error) {
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

func TestMapMust(t *testing.T) {
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
		itemsSlice := MustMap(p.Must(errSink), "near_earth_objects.2023-01-07", func(sm SelectorMustAPI) (Item, error) {
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
		_ = MustMap(p.Must(errSink), "near_earth_objects.2023-01-07", func(sm SelectorMustAPI) (Item, error) {
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
		err := Each(p, "near_earth_objects.2023-01-07", func(index int, _ Picker, length int) error {
			tst.AssertEqual(t, length, 17)
			if index == 4 {
				return errors.New("error")
			}
			return nil
		})
		tst.AssertEqual(t, err, errors.New("error"))
	})

	t.Run("EachM happy path", func(t *testing.T) {
		t.Parallel()
		errSink := &ErrorsSink{}
		MustEach(p.Must(errSink), "near_earth_objects.2023-01-07", func(index int, a SelectorMustAPI, length int) error {
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
		MustEach(p.Must(errSink), "near_earth_objects.2023-01-07", func(index int, a SelectorMustAPI, length int) error {
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

func TestReadme(t *testing.T) {
	assert := tst.AssertEqualFn(t)

	j := `{
    "item": {
        "one": 1,
        "two": "ok",
        "three": ["element 1", 2, "element 3"]
    },
    "float": 2.12
}`
	p1, _ := WrapJSON([]byte(j))
	{
		got, err := p1.String("item.three[1]")
		assert(got, "2")
		assert(err, nil)
	}
	{
		got, err := p1.Uint64("item.three[1]")
		assert(got, uint64(2))
		assert(err, nil)
	}
	{
		got, err := p1.String("item.three[-1]") // access the last element
		assert(got, "element 3")
		assert(err, nil)
	}
	{
		got, err := p1.Float32("float")
		assert(got, float32(2.12))
		assert(err, nil)
	}
	{
		got, err := p1.Int64("float")
		assert(got, int64(2))
		assert(err, ErrCastLostDecimals)
	}
	m := p1.Must()
	{
		got := m.Int32("item.one")
		assert(got, int32(1))
	}
	{
		got := m.Int32("non-existing")
		assert(got, int32(0))
	}
	{
		got, err := Get[int64](p1, "item.three[1]")
		assert(got, int64(2))
		assert(err, nil)
	}
	{
		got, err := Get[string](p1, "item.three[1]")
		assert(got, "2")
		assert(err, nil)
	}
	{
		got := MustGet[string](m, "item.three[1]")
		assert(got, "2")
	}
	{
		got, err := Path[string](p1, Field("item"), Field("three"), Index(1))
		assert(got, "2")
		assert(err, nil)
	}
	{
		got := MustPath[float32](m, Field("item"), Field("one"))
		assert(got, float32(1))
	}

	// Map examples
	j2 := `{
    "items": [
        {"id": 34, "name": "test1", "array": [1,2,3]},
        {"id": 35, "name": "test2", "array": [4,5,6]},
        {"id": 36, "name": "test3", "array": [7,8,9]}
    ]
}`
	p2, _ := WrapJSON([]byte(j2))

	{
		got, err := Map(p2, "items", func(p Picker) (int16, error) {
			n, _ := p.Int16("id")
			return n, nil
		})
		assert(got, []int16{34, 35, 36})
		assert(err, nil)

		got2, err2 := FlatMap(p2, "items", func(p Picker) ([]int16, error) {
			return p.Int16Slice("array")
		})
		assert(got2, []int16{1, 2, 3, 4, 5, 6, 7, 8, 9})
		assert(err2, nil)

		got3, err3 := MapFilter(p2, "items", func(p Picker) (int32, bool, error) {
			i, err := p.Int32("id")
			return i, i%2 == 0, err
		})
		assert(got3, []int32{34, 36})
		assert(err3, nil)
	}

	// Selector Must API
	assert(p1.Must().String("item.three[1]"), "2")
	sm := p1.Must()
	assert(sm.Uint64("item.three[1]"), uint64(2))
	assert(sm.Int32("item.one"), int32(1))
	assert(sm.Float32("float"), float32(2.12))
	assert(sm.Int64("float"), int64(2))

	// Selector Must API With ErrorSink
	sink := ErrorsSink{}
	sm2 := p1.Must(&sink)
	assert(sm2.String("item.three"), "")
	assert(sm2.String("item.three[1]"), "2")
	assert(sm2.Uint64("item.three[1]"), uint64(2))
	assert(sm2.Int32("item.one"), int32(1))
	assert(sm2.Float32("float"), float32(2.12))
	assert(sm2.Int64("float"), int64(2))
	assert(sink.Outcome() != nil, true)
	// es.Outcome() = picker error with selector `item.three` ... invalid type | picker error with selector `float` missing decimals error

	// time API
	dateData := map[string]any{
		"time1":     "1977-05-25T22:30:00Z",
		"time2":     "Wed, 25 May 1977 18:30:00 -0400",
		"timeSlice": []string{"1977-05-25T18:30:00Z", "1977-05-25T20:30:00Z", "1977-05-25T22:30:00Z"},
	}
	p3 := Wrap(dateData)
	{
		got, err := p3.Time("time1")
		assert(got, time.Date(1977, time.May, 25, 22, 30, 0, 0, time.UTC))
		assert(err, nil)
	}
	{
		loc, _ := time.LoadLocation("America/New_York")
		got, err := p3.TimeWithConfig(TimeCastConfig{StringFormat: time.RFC1123Z}, "time2")
		assert(got, time.Date(1977, time.May, 25, 18, 30, 0, 0, loc))
		assert(err, nil)
	}
	{
		got, err := p3.TimeSlice("timeSlice")
		assert(
			got,
			[]time.Time{
				time.Date(1977, time.May, 25, 18, 30, 0, 0, time.UTC),
				time.Date(1977, time.May, 25, 20, 30, 0, 0, time.UTC),
				time.Date(1977, time.May, 25, 22, 30, 0, 0, time.UTC),
			},
		)
		assert(err, nil)
	}
}
