package cast

import (
	"encoding/json"
	"math"
	"testing"
)

func TestStringCaster(t *testing.T) {
	t.Parallel()

	testCases := []casterTestCase[string]{
		{
			input:       uint64(math.MaxUint64),
			expected:    "18446744073709551615",
			expectedErr: nil,
		},
		{
			input:       int32(-123),
			expected:    "-123",
			expectedErr: nil,
		},
		{
			input:       int32(-123),
			expected:    "-123",
			expectedErr: nil,
		},
		{
			input:       float64(math.MaxFloat64),
			expected:    "1.7977E+308",
			expectedErr: nil,
		},
		{
			input:       true,
			expected:    "true",
			expectedErr: nil,
		},
		{
			input:       struct{}{},
			expected:    "",
			expectedErr: expectInvalidType,
		},
		{
			input:       "string",
			expected:    "string",
			expectedErr: nil,
		},
		{
			input:       json.RawMessage(`{"a":"b"}`),
			expected:    `{"a":"b"}`,
			expectedErr: nil,
		},
		{
			input:       json.Number(`5`),
			expected:    "5",
			expectedErr: nil,
		},
		{
			input:       []byte(`abc`),
			expected:    "abc",
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
			expected:    []string{},
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
