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
