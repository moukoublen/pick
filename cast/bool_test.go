package cast

import (
	"math"
	"testing"

	"github.com/moukoublen/pick/internal/testingx"
)

func TestBoolCaster(t *testing.T) {
	t.Parallel()

	testCases := []casterTestCase[bool]{
		{
			input:       uint64(math.MaxUint64),
			expected:    true,
			expectedErr: nil,
		},
		{
			input:       int32(-123),
			expected:    true,
			expectedErr: nil,
		},
		{
			input:       float64(math.MaxFloat64),
			expected:    true,
			expectedErr: nil,
		},
		{
			input:       true,
			expected:    true,
			expectedErr: nil,
		},
		{
			input:       false,
			expected:    false,
			expectedErr: nil,
		},
		{
			input:       struct{}{},
			expected:    false,
			expectedErr: expectInvalidType,
		},
		{
			input:       "string",
			expected:    false,
			expectedErr: testingx.ExpectedErrorStringContains("invalid syntax"),
		},
		{
			input:       "true",
			expected:    true,
			expectedErr: nil,
		},
		{
			input:       "false",
			expected:    false,
			expectedErr: nil,
		},
	}

	caster := newBoolCaster()
	casterTest[bool](t, testCases, caster.AsBool)
}
