package pick

import (
	"reflect"
	"testing"

	"github.com/moukoublen/pick/internal/testingx"
)

func TestDefaultSelectorFormat(t *testing.T) {
	t.Parallel()
	// type mapAlias map[string]any
	name := NameSelectorKey
	index := IndexSelectorKey

	tests := []struct {
		expectedError     func(*testing.T, error)
		input             string
		expectedSelector  []SelectorKey
		expectedFormatted string
	}{
		{
			input:            "",
			expectedSelector: nil,
			expectedError:    nil,
		},
		{
			input:            "one.two",
			expectedSelector: []SelectorKey{name("one"), name("two")},
			expectedError:    nil,
		},
		{
			input:            "one[1]",
			expectedSelector: []SelectorKey{name("one"), index(1)},
			expectedError:    nil,
		},
		{
			input:            "one[1432]",
			expectedSelector: []SelectorKey{name("one"), index(1432)},
			expectedError:    nil,
		},
		{
			input:            "[1][2][3]",
			expectedSelector: []SelectorKey{index(1), index(2), index(3)},
			expectedError:    nil,
		},
		{
			input:            "[154][34][376]",
			expectedSelector: []SelectorKey{index(154), index(34), index(376)},
			expectedError:    nil,
		},
		{
			input:            "[154].a[2].three",
			expectedSelector: []SelectorKey{index(154), name("a"), index(2), name("three")},
			expectedError:    nil,
		},
		{
			input:            "r[154].a[2].three",
			expectedSelector: []SelectorKey{name("r"), index(154), name("a"), index(2), name("three")},
			expectedError:    nil,
		},
		{
			input:            "ελληνικά[154].a[2].three",
			expectedSelector: []SelectorKey{name("ελληνικά"), index(154), name("a"), index(2), name("three")},
			expectedError:    nil,
		},
		{
			input:            "[154].asd[",
			expectedSelector: []SelectorKey(nil),
			expectedError:    testingx.ExpectedErrorIs(ErrInvalidFormatForIndex),
		},
		{
			input:            "[154].asd[r]",
			expectedSelector: []SelectorKey(nil),
			expectedError:    testingx.ExpectedErrorIs(ErrInvalidFormatForIndex),
		},
		{
			input:            "..",
			expectedSelector: []SelectorKey(nil),
			expectedError:    testingx.ExpectedErrorIs(ErrInvalidFormatForName),
		},
	}

	dsf := DefaultSelectorFormat{}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			t.Parallel()
			got, err := dsf.Parse(tc.input)
			testingx.AssertError(t, tc.expectedError, err)
			if !reflect.DeepEqual(tc.expectedSelector, got) {
				t.Errorf("For input: %s \nExpected: %v\nGot     : %v\n", tc.input, tc.expectedSelector, got)
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
