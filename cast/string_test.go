package cast

import (
	"testing"
)

func TestStringCaster(t *testing.T) {
	t.Parallel()

	type int32Alias int32
	type stringAlias string
	type float32alias float32
	testCases := []casterTestCase[string]{
		{
			input:       int32Alias(123456),
			expected:    "123456",
			expectedErr: nil,
		},
		{
			input:       stringAlias("abcd"),
			expected:    "abcd",
			expectedErr: nil,
		},
		{
			input:       float32alias(12.123),
			expected:    "12.123",
			expectedErr: nil,
		},
	}

	caster := newStringCaster()
	casterTest[string](t, testCases, caster.AsString)
}

func TestStringSliceCaster(t *testing.T) {
	t.Parallel()

	testCases := []casterTestCase[[]string]{
		{
			input:       "singe string",
			expected:    []string{"singe string"},
			expectedErr: nil,
		},
		{
			input:       444,
			expected:    []string{"444"},
			expectedErr: nil,
		},
		{
			input:       []string{"string", "slice"},
			expected:    []string{"string", "slice"},
			expectedErr: nil,
		},
		{
			input:       [2]string{"string", "array"},
			expected:    []string{"string", "array"},
			expectedErr: nil,
		},
		{
			input:       []any{"slice", "of", 3},
			expected:    []string{"slice", "of", "3"},
			expectedErr: nil,
		},
		{
			input:       []any{"slice", int32(12), float64(1.456)},
			expected:    []string{"slice", "12", "1.456"},
			expectedErr: nil,
		},
		{
			input:       []int32{1, 2, 3, 4},
			expected:    []string{"1", "2", "3", "4"},
			expectedErr: nil,
		},
		{
			input:       []int32(nil),
			expected:    []string(nil),
			expectedErr: nil,
		},
		{
			input:       map[string]string{"": ""},
			expected:    nil,
			expectedErr: expectInvalidType,
		},
		{
			input:       []any{"slice", map[string]string{"": ""}, 3},
			expected:    nil,
			expectedErr: expectInvalidType,
		},
	}

	caster := newStringCaster()
	casterTest(t, testCases, caster.AsStringSlice)
}
