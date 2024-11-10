package pick

import (
	"testing"
)

func TestBoolCaster(t *testing.T) {
	t.Parallel()

	testCases := []singleCastTestCase[bool]{
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

	caster := NewDefaultCaster()
	runSingleCastTestCases[bool](t, testCases, caster.AsBool)
}
