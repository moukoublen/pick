package pick

import (
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/moukoublen/pick/cast"
	"github.com/moukoublen/pick/internal/testingx"
	"github.com/moukoublen/pick/internal/testingx/testdata"
)

type PickerTestCase struct {
	AccessFn      any
	Selector      string // selector or path
	Path          []Key  // selector or path
	ExpectedValue any
	ExpectedError func(*testing.T, error)
}

func (tc *PickerTestCase) Name() string {
	if tc.Selector != "" {
		return fmt.Sprintf("selector(%s)", tc.Selector)
	}

	return fmt.Sprintf("path(%s)", DotNotation{}.Format(tc.Path...))
}

func (tc *PickerTestCase) Run(t *testing.T) {
	t.Helper()
	t.Parallel()
	pickerFunctionCall := reflect.ValueOf(tc.AccessFn)

	var args []reflect.Value
	if len(tc.Path) > 0 {
		for _, k := range tc.Path {
			args = append(args, reflect.ValueOf(k))
		}
	} else {
		args = []reflect.Value{reflect.ValueOf(tc.Selector)}
	}

	returned := pickerFunctionCall.Call(args)

	got := returned[0].Interface()
	testingx.AssertEqual(t, got, tc.ExpectedValue)

	var receivedError error
	if len(returned) > 1 {
		gotErr := returned[1].Interface()
		receivedError, _ = gotErr.(error)
	}
	testingx.AssertError(t, tc.ExpectedError, receivedError)
}

func TestMixedTypesMap(t *testing.T) {
	t.Parallel()

	p := Wrap(testdata.MixedTypesMap)

	tests := []PickerTestCase{
		{
			AccessFn:      p.Path().String,
			Path:          []Key{Field("stringField")},
			ExpectedValue: "abcd",
			ExpectedError: nil,
		},
		{
			AccessFn:      p.Path().Int,
			Path:          []Key{Field("sliceOfAnyComplex"), Index(0)},
			ExpectedValue: int(2),
			ExpectedError: nil,
		},
		{
			AccessFn:      p.Path().Int8,
			Path:          []Key{Field("sliceOfAnyComplex"), Index(0)},
			ExpectedValue: int8(2),
			ExpectedError: nil,
		},
		{
			AccessFn:      p.Path().Int16,
			Path:          []Key{Field("sliceOfAnyComplex"), Index(0)},
			ExpectedValue: int16(2),
			ExpectedError: nil,
		},
		{
			AccessFn:      p.Path().Int32,
			Path:          []Key{Field("sliceOfAnyComplex"), Index(4)},
			ExpectedValue: int32(5),
			ExpectedError: nil,
		},
		{
			AccessFn:      p.Path().Int64,
			Path:          []Key{Field("sliceOfAnyComplex"), Index(4)},
			ExpectedValue: int64(5),
			ExpectedError: nil,
		},
		{
			AccessFn:      p.Path().Int64,
			Path:          []Key{Field("sliceOfAnyComplex"), Index(3), Field("key3")},
			ExpectedValue: int64(6565),
			ExpectedError: nil,
		},
		{
			AccessFn:      p.Path().Int32,
			Path:          []Key{Field("sliceOfAnyComplex"), Index(2), Field("A")},
			ExpectedValue: int32(3),
			ExpectedError: nil,
		},
		{
			AccessFn:      p.Path().Int32,
			Path:          []Key{Field("sliceOfAnyComplex"), Index(2), Field("Foo")},
			ExpectedValue: int32(0),
			ExpectedError: testingx.ExpectedErrorIs(ErrFieldNotFound),
		},
		{
			AccessFn:      p.Path().Uint,
			Path:          []Key{Field("pointerMapStringAny"), Field("fieldInt32")},
			ExpectedValue: uint(6),
			ExpectedError: nil,
		},
		{
			AccessFn:      p.Path().Uint8,
			Path:          []Key{Field("pointerMapStringAny"), Field("fieldInt32")},
			ExpectedValue: uint8(6),
			ExpectedError: nil,
		},
		{
			AccessFn:      p.Path().Uint16,
			Path:          []Key{Field("pointerMapStringAny"), Field("fieldInt32")},
			ExpectedValue: uint16(6),
			ExpectedError: nil,
		},
		{
			AccessFn:      p.Path().Uint32,
			Path:          []Key{Field("pointerMapStringAny"), Field("fieldInt32")},
			ExpectedValue: uint32(6),
			ExpectedError: nil,
		},
		{
			AccessFn:      p.Path().Uint64,
			Path:          []Key{Field("pointerMapStringAny"), Field("fieldInt32")},
			ExpectedValue: uint64(6),
			ExpectedError: nil,
		},
		{
			AccessFn:      p.Path().Uint64Slice,
			Path:          []Key{Field("pointerMapStringAny"), Field("int32Slice")},
			ExpectedValue: []uint64{10, 11, 12, 13, 14},
			ExpectedError: nil,
		},
		{
			AccessFn:      p.Path().Int64Slice,
			Path:          []Key{Field("pointerMapStringAny"), Field("int32Slice")},
			ExpectedValue: []int64{10, 11, 12, 13, 14},
			ExpectedError: nil,
		},
		{
			AccessFn:      p.Path().StringSlice,
			Path:          []Key{Field("pointerMapStringAny"), Field("int32Slice")},
			ExpectedValue: []string{"10", "11", "12", "13", "14"},
			ExpectedError: nil,
		},
		{
			AccessFn:      p.Path().Bool,
			Path:          []Key{Field("pointerMapStringAny"), Field("fieldBool")},
			ExpectedValue: true,
			ExpectedError: nil,
		},
	}

	for idx, tc := range tests {
		tc := tc
		name := fmt.Sprintf("%d_%s", idx, tc.Name())
		t.Run(name, tc.Run)
	}
}

// TestMixedTypesMapBool makes an extensive test in Bool/BoolSlice functions using all APIs.
func TestMixedTypesMapBool(t *testing.T) {
	t.Parallel()

	p := Wrap(testdata.MixedTypesMap)

	tests := []PickerTestCase{
		{
			AccessFn:      p.Bool,
			Selector:      "pointerMapStringAny.fieldBool",
			ExpectedValue: true,
			ExpectedError: nil,
		},
		{
			AccessFn:      p.Must().Bool,
			Selector:      "pointerMapStringAny.fieldBool",
			ExpectedValue: true,
			ExpectedError: nil,
		},
		{
			AccessFn:      p.Path().Bool,
			Path:          []Key{Field("pointerMapStringAny"), Field("fieldBool")},
			ExpectedValue: true,
			ExpectedError: nil,
		},
		{
			AccessFn:      p.PathMust().Bool,
			Path:          []Key{Field("pointerMapStringAny"), Field("fieldBool")},
			ExpectedValue: true,
			ExpectedError: nil,
		},

		{
			AccessFn:      p.BoolSlice,
			Selector:      "sliceOfAnyComplex[5]",
			ExpectedValue: []bool{true, true, true, false, true, true},
			ExpectedError: nil,
		},
		{
			AccessFn:      p.Must().BoolSlice,
			Selector:      "sliceOfAnyComplex[5]",
			ExpectedValue: []bool{true, true, true, false, true, true},
			ExpectedError: nil,
		},
		{
			AccessFn:      p.Path().BoolSlice,
			Path:          []Key{Field("sliceOfAnyComplex"), Index(5)},
			ExpectedValue: []bool{true, true, true, false, true, true},
			ExpectedError: nil,
		},
		{
			AccessFn:      p.PathMust().BoolSlice,
			Path:          []Key{Field("sliceOfAnyComplex"), Index(5)},
			ExpectedValue: []bool{true, true, true, false, true, true},
			ExpectedError: nil,
		},
	}

	for idx, tc := range tests {
		tc := tc
		name := fmt.Sprintf("%d_%s", idx, tc.Name())
		t.Run(name, tc.Run)
	}
}

// TestMixedTypesMapString makes an extensive test in String/StringSlice functions using all APIs.
func TestMixedTypesMapString(t *testing.T) {
	t.Parallel()

	p := Wrap(testdata.MixedTypesMap)

	tests := []PickerTestCase{
		{
			AccessFn:      p.String,
			Selector:      "sliceOfAnyComplex[1]",
			ExpectedValue: "stringElement",
			ExpectedError: nil,
		},
		{
			AccessFn:      p.Must().String,
			Selector:      "sliceOfAnyComplex[1]",
			ExpectedValue: "stringElement",
			ExpectedError: nil,
		},
		{
			AccessFn:      p.Path().String,
			Path:          []Key{Field("sliceOfAnyComplex"), Index(1)},
			ExpectedValue: "stringElement",
			ExpectedError: nil,
		},
		{
			AccessFn:      p.PathMust().String,
			Path:          []Key{Field("sliceOfAnyComplex"), Index(1)},
			ExpectedValue: "stringElement",
			ExpectedError: nil,
		},

		{
			AccessFn:      p.StringSlice,
			Selector:      "sliceOfAnyComplex[6]",
			ExpectedValue: []string{"abc", "def", "ghi"},
			ExpectedError: nil,
		},
		{
			AccessFn:      p.Must().StringSlice,
			Selector:      "sliceOfAnyComplex[6]",
			ExpectedValue: []string{"abc", "def", "ghi"},
			ExpectedError: nil,
		},
		{
			AccessFn:      p.Path().StringSlice,
			Path:          []Key{Field("sliceOfAnyComplex"), Index(6)},
			ExpectedValue: []string{"abc", "def", "ghi"},
			ExpectedError: nil,
		},
		{
			AccessFn:      p.PathMust().StringSlice,
			Path:          []Key{Field("sliceOfAnyComplex"), Index(6)},
			ExpectedValue: []string{"abc", "def", "ghi"},
			ExpectedError: nil,
		},
	}

	for idx, tc := range tests {
		tc := tc
		name := fmt.Sprintf("%d_%s", idx, tc.Name())
		t.Run(name, tc.Run)
	}
}

// TestMixedTypesMapInt64 makes an extensive test in Int64/Int64Slice functions using all APIs.
func TestMixedTypesMapInt64(t *testing.T) {
	t.Parallel()

	p := Wrap(testdata.MixedTypesMap)

	tests := []PickerTestCase{
		{
			AccessFn:      p.Int64,
			Selector:      "sliceOfAnyComplex[0]",
			ExpectedValue: int64(2),
			ExpectedError: nil,
		},
		{
			AccessFn:      p.Must().Int64,
			Selector:      "sliceOfAnyComplex[0]",
			ExpectedValue: int64(2),
			ExpectedError: nil,
		},
		{
			AccessFn:      p.Path().Int64,
			Path:          []Key{Field("sliceOfAnyComplex"), Index(0)},
			ExpectedValue: int64(2),
			ExpectedError: nil,
		},
		{
			AccessFn:      p.PathMust().Int64,
			Path:          []Key{Field("sliceOfAnyComplex"), Index(0)},
			ExpectedValue: int64(2),
			ExpectedError: nil,
		},

		{
			AccessFn:      p.Int64Slice,
			Selector:      "pointerMapStringAny.int32Slice",
			ExpectedValue: []int64{10, 11, 12, 13, 14},
			ExpectedError: nil,
		},
		{
			AccessFn:      p.Must().Int64Slice,
			Selector:      "pointerMapStringAny.int32Slice",
			ExpectedValue: []int64{10, 11, 12, 13, 14},
			ExpectedError: nil,
		},
		{
			AccessFn:      p.Path().Int64Slice,
			Path:          []Key{Field("pointerMapStringAny"), Field("int32Slice")},
			ExpectedValue: []int64{10, 11, 12, 13, 14},
			ExpectedError: nil,
		},
		{
			AccessFn:      p.PathMust().Int64Slice,
			Path:          []Key{Field("pointerMapStringAny"), Field("int32Slice")},
			ExpectedValue: []int64{10, 11, 12, 13, 14},
			ExpectedError: nil,
		},
	}

	for idx, tc := range tests {
		tc := tc
		name := fmt.Sprintf("%d_%s", idx, tc.Name())
		t.Run(name, tc.Run)
	}
}

// TestMixedTypesMapFloat64 makes an extensive test in Float64/Float64Slice functions using all APIs.
func TestMixedTypesMapFloat64(t *testing.T) {
	t.Parallel()

	p := Wrap(testdata.MixedTypesMap)

	tests := []PickerTestCase{
		{
			AccessFn:      p.Float64,
			Selector:      "pointerMapStringAny.float64Slice[3]",
			ExpectedValue: float64(0.4),
			ExpectedError: nil,
		},
		{
			AccessFn:      p.Must().Float64,
			Selector:      "pointerMapStringAny.float64Slice[3]",
			ExpectedValue: float64(0.4),
			ExpectedError: nil,
		},
		{
			AccessFn:      p.Path().Float64,
			Path:          []Key{Field("pointerMapStringAny"), Field("float64Slice"), Index(3)},
			ExpectedValue: float64(0.4),
			ExpectedError: nil,
		},
		{
			AccessFn:      p.PathMust().Float64,
			Path:          []Key{Field("pointerMapStringAny"), Field("float64Slice"), Index(3)},
			ExpectedValue: float64(0.4),
			ExpectedError: nil,
		},

		{
			AccessFn:      p.Float64Slice,
			Selector:      "pointerMapStringAny.float64Slice",
			ExpectedValue: []float64{0.1, 0.2, 0.3, 0.4},
			ExpectedError: nil,
		},
		{
			AccessFn:      p.Must().Float64Slice,
			Selector:      "pointerMapStringAny.float64Slice",
			ExpectedValue: []float64{0.1, 0.2, 0.3, 0.4},
			ExpectedError: nil,
		},
		{
			AccessFn:      p.Path().Float64Slice,
			Path:          []Key{Field("pointerMapStringAny"), Field("float64Slice")},
			ExpectedValue: []float64{0.1, 0.2, 0.3, 0.4},
			ExpectedError: nil,
		},
		{
			AccessFn:      p.PathMust().Float64Slice,
			Path:          []Key{Field("pointerMapStringAny"), Field("float64Slice")},
			ExpectedValue: []float64{0.1, 0.2, 0.3, 0.4},
			ExpectedError: nil,
		},
	}

	for idx, tc := range tests {
		tc := tc
		name := fmt.Sprintf("%d_%s", idx, tc.Name())
		t.Run(name, tc.Run)
	}
}

//go:embed internal/testingx/testdata
var testData embed.FS

func loadTestData(t *testing.T, filename string) fs.File {
	t.Helper()

	path := filepath.Join("internal", "testingx", "testdata", filename)
	f, err := testData.Open(path)
	if err != nil {
		t.Fatalf("error during testdate file opening %s", err.Error())
	}

	return f
}

func TestNasaDataFile(t *testing.T) {
	t.Parallel()

	file := loadTestData(t, "nasa.json")

	p, err := WrapReaderJSON(file)
	if err != nil {
		t.Fatal(err)
	}

	tests := []PickerTestCase{
		{
			AccessFn:      p.String,
			Selector:      "near_earth_objects.2023-01-01[4].neo_reference_id",
			ExpectedValue: "3703782",
			ExpectedError: nil,
		},
		{
			AccessFn:      p.Float32,
			Selector:      "near_earth_objects.2023-01-01[5].estimated_diameter.meters.estimated_diameter_max",
			ExpectedValue: float32(68.2401509401),
			ExpectedError: nil,
		},
		{
			AccessFn:      p.Float64,
			Selector:      "near_earth_objects.2023-01-01[5].estimated_diameter.meters.estimated_diameter_max",
			ExpectedValue: float64(68.2401509401),
			ExpectedError: nil,
		},
		{
			AccessFn:      p.Uint8,
			Selector:      "near_earth_objects.2023-01-01[5].id",
			ExpectedValue: uint8(214),
			ExpectedError: testingx.ExpectedErrorIs(cast.ErrCastOverFlow),
		},
		{
			AccessFn:      p.Uint16,
			Selector:      "near_earth_objects.2023-01-01[5].id",
			ExpectedValue: uint16(50902),
			ExpectedError: testingx.ExpectedErrorIs(cast.ErrCastOverFlow),
		},
		{
			AccessFn:      p.Uint32,
			Selector:      "near_earth_objects.2023-01-01[5].id",
			ExpectedValue: uint32(3720918),
			ExpectedError: nil,
		},
		{
			AccessFn:      p.Uint64,
			Selector:      "near_earth_objects.2023-01-01[5].id",
			ExpectedValue: uint64(3720918),
			ExpectedError: nil,
		},
		{
			AccessFn:      p.Uint,
			Selector:      "near_earth_objects.2023-01-01[5].id",
			ExpectedValue: uint(3720918),
			ExpectedError: nil,
		},
		{
			AccessFn:      p.Bool,
			Selector:      "near_earth_objects.2023-01-01[1].is_potentially_hazardous_asteroid",
			ExpectedValue: true,
			ExpectedError: nil,
		},
		{
			AccessFn: func(selector string) ([]string, error) {
				return Map(p, selector, func(p *Picker) (string, error) { return p.String("id") })
			},
			Selector:      "near_earth_objects.2023-01-01",
			ExpectedValue: []string{"2154347", "2385186", "2453309", "3683468", "3703782", "3720918", "3767936", "3792438", "3824981", "3836251", "3837605", "3959234", "3986848", "54104550", "54105994", "54166175", "54202993", "54290862", "54335607", "54337027", "54337425", "54340039", "54341664"},
			ExpectedError: nil,
		},
		{
			AccessFn: func(selector string) ([]string, error) {
				return FlatMap(p, "near_earth_objects.2023-01-01", func(p *Picker) ([]string, error) {
					return Map(p, "close_approach_data", func(p *Picker) (string, error) {
						return p.String("close_approach_date_full")
					})
				})
			},
			Selector:      "",
			ExpectedValue: []string{"2023-Jan-01 18:44", "2023-Jan-01 19:45", "2023-Jan-01 20:20", "2023-Jan-01 13:38", "2023-Jan-01 00:59", "2023-Jan-01 17:33", "2023-Jan-01 09:38", "2023-Jan-01 09:49", "2023-Jan-01 03:04", "2023-Jan-01 22:31", "2023-Jan-01 04:15", "2023-Jan-01 02:10", "2023-Jan-01 10:47", "2023-Jan-01 16:46", "2023-Jan-01 12:02", "2023-Jan-01 16:03", "2023-Jan-01 13:39", "2023-Jan-01 12:50", "2023-Jan-01 20:45", "2023-Jan-01 07:16", "2023-Jan-01 01:15", "2023-Jan-01 23:21", "2023-Jan-01 09:02"},
			ExpectedError: nil,
		},
	}

	for idx, tc := range tests {
		tc := tc
		name := fmt.Sprintf("%d_%s", idx, tc.Name())
		t.Run(name, tc.Run)
	}
}

func TestReadme(t *testing.T) {
	assert := testingx.AssertEqualFn(t)

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
		got := p1.Must().Int32("item.one")
		assert(got, int32(1))
	}
	{
		got, err := p1.Float32("float")
		assert(got, float32(2.12))
		assert(err, nil)
	}
	{
		got, err := p1.Int64("float")
		assert(got, int64(2))
		assert(err, cast.ErrCastLostDecimals)
	}

	// Map examples
	j2 := `{
    "items": [
        {"id": 34, "name": "test1"},
        {"id": 35, "name": "test2"},
        {"id": 36, "name": "test3"}
    ]
}`
	p2, _ := WrapJSON([]byte(j2))

	type Foo struct{ ID int16 }

	{
		got, err := Map(p2, "items", func(p *Picker) (Foo, error) {
			f := Foo{}
			f.ID, _ = p.Int16("id")
			return f, nil
		})
		assert(got, []Foo{{ID: 34}, {ID: 35}, {ID: 36}})
		assert(err, nil)
	}

	// Selector Must API
	assert(p1.Must().String("item.three[1]"), "2")
	assert(p1.Must().Uint64("item.three[1]"), uint64(2))
	sm := p1.Must()
	assert(sm.Int32("item.one"), int32(1))
	assert(sm.Float32("float"), float32(2.12))
	assert(sm.Int64("float"), int64(2))

	// Path API
	{
		got, err := p1.Path().String(Field("item"), Field("three"), Index(1))
		assert(got, "2")
		assert(err, nil)
	}
	pa := p1.Path()
	{
		got, err := pa.Uint64(Field("item"), Field("three"), Index(1))
		assert(got, uint64(2))
		assert(err, nil)
	}
	{
		got, err := pa.Int32(Field("item"), Field("one"))
		assert(got, int32(1))
		assert(err, nil)
	}
	pm := p1.PathMust()
	assert(pm.Float32(Field("float")), float32(2.12))
	assert(pm.Int64(Field("float")), int64(2))
}
