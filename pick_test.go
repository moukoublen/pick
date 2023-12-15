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

	ob, err := WrapReaderJSON(file)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		accessFn      any
		expectedValue any
		expectedError func(*testing.T, error)
		selector      string
	}{
		{
			selector:      "near_earth_objects.2023-01-01[4].neo_reference_id",
			accessFn:      ob.String,
			expectedValue: "3703782",
			expectedError: nil,
		},
		{
			selector:      "near_earth_objects.2023-01-01[5].estimated_diameter.meters.estimated_diameter_max",
			accessFn:      ob.Float64,
			expectedValue: float64(68.2401509401),
			expectedError: nil,
		},
		{
			selector:      "near_earth_objects.2023-01-01[5].id",
			accessFn:      ob.Uint64,
			expectedValue: uint64(3720918),
			expectedError: nil,
		},
		{
			selector:      "near_earth_objects.2023-01-01[1].is_potentially_hazardous_asteroid",
			accessFn:      ob.Bool,
			expectedValue: true,
			expectedError: nil,
		},
		{
			selector: "near_earth_objects.2023-01-01",
			accessFn: func(selector string) ([]string, error) {
				return Map(ob, selector, func(p *Picker) (string, error) { return p.String("id") })
			},
			expectedValue: []string{"2154347", "2385186", "2453309", "3683468", "3703782", "3720918", "3767936", "3792438", "3824981", "3836251", "3837605", "3959234", "3986848", "54104550", "54105994", "54166175", "54202993", "54290862", "54335607", "54337027", "54337425", "54340039", "54341664"},
			expectedError: nil,
		},
		{
			selector: "",
			accessFn: func(selector string) ([]string, error) {
				return FlatMap(ob, "near_earth_objects.2023-01-01", func(p *Picker) ([]string, error) {
					return Map(p, "close_approach_data", func(p *Picker) (string, error) {
						return p.String("close_approach_date_full")
					})
				})
			},
			expectedValue: []string{"2023-Jan-01 18:44", "2023-Jan-01 19:45", "2023-Jan-01 20:20", "2023-Jan-01 13:38", "2023-Jan-01 00:59", "2023-Jan-01 17:33", "2023-Jan-01 09:38", "2023-Jan-01 09:49", "2023-Jan-01 03:04", "2023-Jan-01 22:31", "2023-Jan-01 04:15", "2023-Jan-01 02:10", "2023-Jan-01 10:47", "2023-Jan-01 16:46", "2023-Jan-01 12:02", "2023-Jan-01 16:03", "2023-Jan-01 13:39", "2023-Jan-01 12:50", "2023-Jan-01 20:45", "2023-Jan-01 07:16", "2023-Jan-01 01:15", "2023-Jan-01 23:21", "2023-Jan-01 09:02"},
			expectedError: nil,
		},
	}

	for idx, tc := range tests {
		tc := tc
		name := fmt.Sprintf("%d_%s", idx, tc.selector)
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			val := reflect.ValueOf(tc.accessFn)
			ret := val.Call([]reflect.Value{reflect.ValueOf(tc.selector)})
			got := ret[0].Interface()

			var err error
			if i := ret[1].Interface(); i != nil {
				err, _ = i.(error)
			}

			testingx.AssertError(t, tc.expectedError, err)

			if !reflect.DeepEqual(tc.expectedValue, got) {
				t.Errorf("wrong returned value, expected %T(%#v) found %T(%#v)", tc.expectedValue, tc.expectedValue, got, got)
			}
		})
	}
}

func TestMixedTypesMap(t *testing.T) {
	t.Parallel()

	ob := Wrap(testdata.MixedTypesMap)

	assert := testingx.AssertCompareFn(t)

	assert(ob.PathMust().String(Field("sliceOfAnyComplex"), Index(2), Field("C")), "asdf")
	assert(ob.PathMust().Int16(Field("sliceOfAnyComplex"), Index(2), Field("B")), int16(12))
	{
		got, err := ob.Path().String(Field("sliceOfAnyComplex"), Index(3))
		testingx.AssertError(t, testingx.ExpectedErrorIs(cast.ErrInvalidType), err)
		assert(got, "")
	}
	{
		got, err := ob.Path().String(Field("sliceOfAnyComplex"), Index(3), Field("key2"))
		assert(nil, err)
		assert(got, "value2")
	}
	{
		got, err := ob.Path().String(Field("sliceOfAnyComplex"), Index(3), Field("key2"))
		assert(nil, err)
		assert(got, "value2")
	}
	assert(ob.PathMust().Int64(Field("int32Number")), int64(12954))
	assert(ob.PathMust().Int32(Field("int32Number")), int32(12954))
	assert(ob.PathMust().Int16(Field("int32Number")), int16(12954))
	assert(ob.PathMust().Int8(Field("int32Number")), int8(-102))
	{
		got, err := ob.Path().Int8(Field("int32Number"))
		testingx.AssertError(t, testingx.ExpectedErrorIs(cast.ErrCastOverFlow), err)
		assert(got, int8(-102))
	}
	assert(ob.PathMust().Uint32(Field("sliceOfAnyComplex"), Index(4)), uint32(555))
	assert(ob.PathMust().Bool(Field("pointerMapStringAny"), Field("fieldBool")), true)
	assert(ob.PathMust().Byte(Field("pointerMapStringAny"), Field("fieldByte")), byte('.'))
}

func TestReadme(t *testing.T) {
	assert := testingx.AssertCompareFn(t)

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
