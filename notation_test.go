package pick

import (
	"reflect"
	"testing"

	"github.com/moukoublen/pick/internal/testingx"
)

func TestDotNotation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		expectedError     func(*testing.T, error)
		input             string
		expectedFields    []Field
		expectedFormatted string
	}{
		{
			input:          "",
			expectedFields: nil,
			expectedError:  nil,
		},
		{
			input:          "one.two",
			expectedFields: []Field{Name("one"), Name("two")},
			expectedError:  nil,
		},
		{
			input:          "one[1]",
			expectedFields: []Field{Name("one"), Index(1)},
			expectedError:  nil,
		},
		{
			input:          "one[1432]",
			expectedFields: []Field{Name("one"), Index(1432)},
			expectedError:  nil,
		},
		{
			input:          "[1][2][3]",
			expectedFields: []Field{Index(1), Index(2), Index(3)},
			expectedError:  nil,
		},
		{
			input:          "[154][34][376]",
			expectedFields: []Field{Index(154), Index(34), Index(376)},
			expectedError:  nil,
		},
		{
			input:          "[154].a[2].three",
			expectedFields: []Field{Index(154), Name("a"), Index(2), Name("three")},
			expectedError:  nil,
		},
		{
			input:          "r[154].a[2].three",
			expectedFields: []Field{Name("r"), Index(154), Name("a"), Index(2), Name("three")},
			expectedError:  nil,
		},
		{
			input:          "ελληνικά[154].a[2].three",
			expectedFields: []Field{Name("ελληνικά"), Index(154), Name("a"), Index(2), Name("three")},
			expectedError:  nil,
		},
		{
			input:          "[154].asd[",
			expectedFields: []Field(nil),
			expectedError:  testingx.ExpectedErrorIs(ErrInvalidFormatForIndex),
		},
		{
			input:          "[154].asd[r]",
			expectedFields: []Field(nil),
			expectedError:  testingx.ExpectedErrorIs(ErrInvalidFormatForIndex),
		},
		{
			input:          "..",
			expectedFields: []Field(nil),
			expectedError:  testingx.ExpectedErrorIs(ErrInvalidFormatForName),
		},
	}

	dsf := DotNotation{}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			t.Parallel()
			got, err := dsf.Parse(tc.input)
			testingx.AssertError(t, tc.expectedError, err)
			if !reflect.DeepEqual(tc.expectedFields, got) {
				t.Errorf("For input: %s \nExpected: %v\nGot     : %v\n", tc.input, tc.expectedFields, got)
			}

			if tc.expectedError != nil {
				return
			}

			// format and check formatted.
			gotFormatted := dsf.Format(got)
			expectedFormatted := tc.input
			if tc.expectedFormatted != "" {
				expectedFormatted = tc.expectedFormatted
			}
			if expectedFormatted != gotFormatted {
				t.Errorf("For input: %s \nExpected formatted: %s\nGot formatted     : %s\n", tc.input, expectedFormatted, gotFormatted)
			}
		})
	}
}

func BenchmarkDotNotation(b *testing.B) {
	tests := []string{
		"one.two",
		"one[1]",
		"[154][34][376]",
		"[154].a[2].three",
		"ελληνικά[154].a[2].three",
		"near_earth_objects.2023-01-01[1].is_potentially_hazardous_asteroid",
		"near_earth_objects.2023-01-01[5].estimated_diameter.meters.estimated_diameter_max",
		"near_earth_objects_estimated_diameter_meters_estimated_diameter_max",
		"one",
		"[123]",
	}

	d := DotNotation{}

	for _, tc := range tests {
		tc := tc
		b.Run(tc, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = d.Parse(tc)
			}
		})
	}
}
