package pick

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/moukoublen/pick/cast"
	"github.com/moukoublen/pick/internal/testingx"
)

var expectInvalidType = testingx.ExpectedErrorIs(&cast.Error{}, cast.ErrInvalidType)

func ptr[T any](input T) *T { return &input }

func TestString(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		inner       any
		expectedErr func(*testing.T, error)
		expected    string
		selector    string
	}{
		"empty map": {
			inner:    map[string]any{},
			selector: "a.b",

			expected:    "",
			expectedErr: testingx.ExpectedErrorIs(ErrFieldNotFound),
		},
		"1 level found": {
			inner: map[string]any{
				"a": "b",
			},
			selector: "a",

			expected:    "b",
			expectedErr: nil,
		},
		"1 level found pointer": {
			inner: map[string]any{
				"ptr": ptr("b"),
			},
			selector: "ptr",

			expected:    "b",
			expectedErr: nil,
		},
		"2 level found": {
			inner: map[string]any{
				"a": map[string]any{"b": "c"},
			},
			selector: "a.b",

			expected:    "c",
			expectedErr: nil,
		},
		"2 level found int to string": {
			inner: map[string]any{
				"a": map[string]any{"b": 123},
			},
			selector: "a.b",

			expected:    "123",
			expectedErr: nil,
		},
		"2 level not found": {
			inner: map[string]any{
				"a": map[string]any{"b": 123},
			},
			selector: "a.h",

			expected:    "",
			expectedErr: testingx.ExpectedErrorIs(ErrFieldNotFound),
		},
		"2 level found cast error": {
			inner: map[string]any{
				"a": map[string]any{"b": struct{}{}},
			},
			selector: "a.b",

			expected:    "",
			expectedErr: expectInvalidType,
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			m := Wrap(tc.inner)
			s, err := m.String(tc.selector)

			switch {
			case tc.expectedErr != nil && err != nil:
				tc.expectedErr(t, err)
			case err != nil:
				t.Errorf("unexpected error received: %s [%#v]", err.Error(), err)
			case tc.expectedErr != nil:
				t.Errorf("expected error but none received")
			}

			if tc.expected != s {
				t.Errorf("expected result to be %v got %v", tc.expected, s)
			}
		})
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

func TestReadme(t *testing.T) {
	j := `{
    "item": {
        "one": 1,
        "two": "ok",
        "three": ["element 1", 2, "element 3"]
    },
    "float": 2.12
}`
	p, _ := WrapJSON([]byte(j))

	{
		returned, err := p.String("item.three[1]")
		assert(t, "2", returned, nil, err)
	}
	{
		returned, err := p.Uint64("item.three[1]")
		assert(t, uint64(2), returned, nil, err)
	}
	{
		returned, err := p.Int32("item.one")
		assert(t, int32(1), returned, nil, err)
	}
	{
		returned, err := p.Float32("float")
		assert(t, float32(2.12), returned, nil, err)
	}
	{
		returned, err := p.Int64("float")
		assert(t, int64(2), returned, cast.ErrCastLostDecimals, err)
	}

	j2 := `{
    "items": [
        {"id": 34, "name": "test1"},
        {"id": 35, "name": "test2"},
        {"id": 36, "name": "test3"}
    ]
}`
	p2, _ := WrapJSON([]byte(j2))

	type Foo struct{ ID int16 }

	slice, err := Map(p2, "items", func(p *Picker) (Foo, error) {
		f := Foo{}
		f.ID, _ = p.Int16("id")
		return f, nil
	})
	assert(t, []Foo{{ID: 34}, {ID: 35}, {ID: 36}}, slice, nil, err)
}

func assert(t *testing.T, a, b any, errA, errB error) {
	t.Helper()
	if !reflect.DeepEqual(a, b) {
		t.Errorf("expected %T(%#v) got %T(%#v)", a, a, b, b)
	}

	if !errors.Is(errB, errA) {
		t.Errorf("expected %T(%#v) got %T(%#v)", errA, errA, errB, errB)
	}
}
