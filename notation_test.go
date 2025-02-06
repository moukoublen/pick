package pick

import (
	"strconv"
	"testing"

	"github.com/moukoublen/pick/internal/tst"
)

func TestDotNotation(t *testing.T) {
	// t.Parallel()

	tests := []struct {
		errorAsserter     tst.ErrorAsserter
		input             string
		expectedPath      []Key
		expectedFormatted string
	}{
		{
			input:         "",
			expectedPath:  nil,
			errorAsserter: tst.NoError,
		},
		{
			input:         "one.two",
			expectedPath:  []Key{Field("one"), Field("two")},
			errorAsserter: tst.NoError,
		},
		{
			input:         "one[1]",
			expectedPath:  []Key{Field("one"), Index(1)},
			errorAsserter: tst.NoError,
		},
		{
			input:         "one[1432]",
			expectedPath:  []Key{Field("one"), Index(1432)},
			errorAsserter: tst.NoError,
		},
		{
			input:         "[1][2][3]",
			expectedPath:  []Key{Index(1), Index(2), Index(3)},
			errorAsserter: tst.NoError,
		},
		{
			input:         "[1][-1].field",
			expectedPath:  []Key{Index(1), Index(-1), Field("field")},
			errorAsserter: tst.NoError,
		},
		{
			input:         "[154][34][376]",
			expectedPath:  []Key{Index(154), Index(34), Index(376)},
			errorAsserter: tst.NoError,
		},
		{
			input:         "[154].a[2].three",
			expectedPath:  []Key{Index(154), Field("a"), Index(2), Field("three")},
			errorAsserter: tst.NoError,
		},
		{
			input:         "r[154].a[2].three",
			expectedPath:  []Key{Field("r"), Index(154), Field("a"), Index(2), Field("three")},
			errorAsserter: tst.NoError,
		},
		{
			input:         "ελληνικά[154].a[2].three",
			expectedPath:  []Key{Field("ελληνικά"), Index(154), Field("a"), Index(2), Field("three")},
			errorAsserter: tst.NoError,
		},
		{
			input:         "start[3].ελληνικά.a[2].three",
			expectedPath:  []Key{Field("start"), Index(3), Field("ελληνικά"), Field("a"), Index(2), Field("three")},
			errorAsserter: tst.NoError,
		},
		{
			input:         "[154].asd[",
			expectedPath:  []Key(nil),
			errorAsserter: tst.ExpectedErrorIs(ErrInvalidSelectorFormatForIndex),
		},
		{
			input:         "[154].asd.",
			expectedPath:  []Key(nil),
			errorAsserter: tst.ExpectedErrorIs(ErrInvalidSelectorFormatForName),
		},
		{
			input:         "[154].asd[r]",
			expectedPath:  []Key(nil),
			errorAsserter: tst.ExpectedErrorIs(ErrInvalidSelectorFormatForIndex),
		},
		{
			input:         "..",
			expectedPath:  []Key(nil),
			errorAsserter: tst.ExpectedErrorIs(ErrInvalidSelectorFormatForName),
		},
	}

	dsf := DotNotation{}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			// t.Parallel()
			got, err := dsf.Parse(tc.input)
			tc.errorAsserter(t, err)
			tst.AssertEqual(t, got, tc.expectedPath)

			if err != nil {
				return
			}

			// format and check formatted.
			gotFormatted := dsf.Format(got...)
			expectedFormatted := tc.input
			if tc.expectedFormatted != "" {
				expectedFormatted = tc.expectedFormatted
			}
			tst.AssertEqual(t, expectedFormatted, gotFormatted)
		})
	}
}

func BenchmarkDotNotation(b *testing.B) {
	tests := []string{
		0:  "",
		1:  "one[1]",
		2:  "one.two",
		3:  "[154][34][376]",
		4:  "[1][-1].field",
		5:  "[154].a[2].three",
		6:  "ελληνικά[154].a[2].three",
		7:  "start[3].ελληνικά.a[2].three",
		8:  "near_earth_objects.2023-01-01[1].is_potentially_hazardous_asteroid",
		9:  "near_earth_objects.2023-01-01[5].estimated_diameter.meters.estimated_diameter_max",
		10: "near_earth_objects_estimated_diameter_meters_estimated_diameter_max",
		11: "one",
		12: "[123]",
	}

	d := DotNotation{}

	for i, tc := range tests {
		name := "0000" + strconv.Itoa(i)
		b.Run(name[(len(name)-4):], func(b *testing.B) {
			for range b.N {
				_, _ = d.Parse(tc)
			}
		})
	}
}
