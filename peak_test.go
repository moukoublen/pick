package pick

import (
	"embed"
	"encoding/json"
	"fmt"
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
		inner         any
		expectedErr   func(*testing.T, error)
		expected      string
		selector      string
		expectedFound bool
	}{
		"empty map": {
			inner:    map[string]any{},
			selector: "a.b",

			expected:      "",
			expectedFound: false,
			expectedErr:   nil,
		},
		"1 level found": {
			inner: map[string]any{
				"a": "b",
			},
			selector: "a",

			expected:      "b",
			expectedFound: true,
			expectedErr:   nil,
		},
		"1 level found pointer": {
			inner: map[string]any{
				"ptr": ptr("b"),
			},
			selector: "ptr",

			expected:      "b",
			expectedFound: true,
			expectedErr:   nil,
		},
		"2 level found": {
			inner: map[string]any{
				"a": map[string]any{"b": "c"},
			},
			selector: "a.b",

			expected:      "c",
			expectedFound: true,
			expectedErr:   nil,
		},
		"2 level found int to string": {
			inner: map[string]any{
				"a": map[string]any{"b": 123},
			},
			selector: "a.b",

			expected:      "123",
			expectedFound: true,
			expectedErr:   nil,
		},
		"2 level not found": {
			inner: map[string]any{
				"a": map[string]any{"b": 123},
			},
			selector: "a.h",

			expected:      "",
			expectedFound: false,
			expectedErr:   nil,
		},
		"2 level found cast error": {
			inner: map[string]any{
				"a": map[string]any{"b": struct{}{}},
			},
			selector: "a.b",

			expected:      "",
			expectedFound: true,
			expectedErr:   expectInvalidType,
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			m := Wrap(tc.inner)
			s, found, err := m.String(tc.selector)

			switch {
			case tc.expectedErr != nil && err != nil:
				tc.expectedErr(t, err)
			case err != nil:
				t.Errorf("unexpected error received: %s [%#v]", err.Error(), err)
			case tc.expectedErr != nil:
				t.Errorf("expected error but none received")
			}

			if tc.expectedFound != found {
				t.Errorf("expected found to be %v got %v", tc.expectedFound, found)
			}

			if tc.expected != s {
				t.Errorf("expected result to be %v got %v", tc.expected, s)
			}
		})
	}
}

//go:embed internal/testingx/testdata
var testData embed.FS

func loadTestData(t *testing.T, filename string, decodeInto any) {
	t.Helper()

	path := filepath.Join("internal", "testingx", "testdata", filename)
	f, err := testData.Open(path)
	if err != nil {
		t.Fatalf("error during testdate file opening %s", err.Error())
	}

	if err := json.NewDecoder(f).Decode(decodeInto); err != nil {
		t.Fatalf("error during testdate file decoding %s", err.Error())
	}
}

func TestNasaDataFile(t *testing.T) {
	t.Parallel()
	inner := map[string]any{}
	loadTestData(t, "nasa.json", &inner)

	ob := Wrap(inner)

	tests := []struct {
		accessFn      any
		expectedValue any
		expectedError func(*testing.T, error)
		selector      string
		expectedFound bool
	}{
		{
			selector:      "near_earth_objects.2023-01-01[4].neo_reference_id",
			accessFn:      ob.String,
			expectedValue: "3703782",
			expectedFound: true,
			expectedError: nil,
		},
		{
			selector:      "near_earth_objects.2023-01-01[5].estimated_diameter.meters.estimated_diameter_max",
			accessFn:      ob.Float64,
			expectedValue: float64(68.2401509401),
			expectedFound: true,
			expectedError: nil,
		},
		{
			selector:      "near_earth_objects.2023-01-01[5].id",
			accessFn:      ob.Uint64,
			expectedValue: uint64(3720918),
			expectedFound: true,
			expectedError: nil,
		},
		{
			selector:      "near_earth_objects.2023-01-01[1].is_potentially_hazardous_asteroid",
			accessFn:      ob.Bool,
			expectedValue: true,
			expectedFound: true,
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
			found := ret[1].Bool()

			var err error
			if i := ret[2].Interface(); i != nil {
				err, _ = i.(error)
			}

			testingx.AssertError(t, tc.expectedError, err)

			if found != tc.expectedFound {
				t.Errorf("wrong returned value found, expected %t found %t", tc.expectedFound, found)
			}
			if !reflect.DeepEqual(tc.expectedValue, got) {
				typeOfExpected := reflect.TypeOf(tc.expectedValue)
				typeOfGot := reflect.TypeOf(got)
				t.Errorf("wrong returned value, expected %s(%#v) found %s(%#v)", typeOfExpected.String(), tc.expectedValue, typeOfGot.String(), got)
			}
		})
	}
}
