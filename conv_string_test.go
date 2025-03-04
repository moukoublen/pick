package pick

import (
	"testing"

	"github.com/moukoublen/pick/internal/tst"
)

func TestStringConverter(t *testing.T) {
	t.Parallel()

	type int32Alias int32
	type stringAlias string
	type float32alias float32
	testCases := []singleConvertTestCase[string]{
		{
			input:         int32Alias(123456),
			expected:      "123456",
			errorAsserter: tst.NoError,
		},
		{
			input:         stringAlias("abcd"),
			expected:      "abcd",
			errorAsserter: tst.NoError,
		},
		{
			input:         float32alias(12.123),
			expected:      "12.123",
			errorAsserter: tst.NoError,
		},
	}

	converter := NewDefaultConverter()
	runSingleConvertTestCases[string](t, testCases, converter.AsString)
}

func TestStringSliceConverter(t *testing.T) {
	t.Parallel()

	testCases := []singleConvertTestCase[[]string]{
		{
			input:         "singe string",
			expected:      []string{"singe string"},
			errorAsserter: tst.NoError,
		},
		{
			input:         444,
			expected:      []string{"444"},
			errorAsserter: tst.NoError,
		},
		{
			input:         []string{"string", "slice"},
			expected:      []string{"string", "slice"},
			errorAsserter: tst.NoError,
		},
		{
			input:         [2]string{"string", "array"},
			expected:      []string{"string", "array"},
			errorAsserter: tst.NoError,
		},
		{
			input:         []any{"slice", "of", 3},
			expected:      []string{"slice", "of", "3"},
			errorAsserter: tst.NoError,
		},
		{
			input:         []any{"slice", int32(12), float64(1.456)},
			expected:      []string{"slice", "12", "1.456"},
			errorAsserter: tst.NoError,
		},
		{
			input:         []int32{1, 2, 3, 4},
			expected:      []string{"1", "2", "3", "4"},
			errorAsserter: tst.NoError,
		},
		{
			input:         []int32(nil),
			expected:      []string(nil),
			errorAsserter: tst.NoError,
		},
		{
			input:         map[string]string{"": ""},
			expected:      nil,
			errorAsserter: expectInvalidType,
		},
		{
			input:         []any{"slice", map[string]string{"": ""}, 3},
			expected:      nil,
			errorAsserter: expectInvalidType,
		},
	}

	converter := NewDefaultConverter()
	runSingleConvertTestCases(t, testCases, converter.AsStringSlice)
}
