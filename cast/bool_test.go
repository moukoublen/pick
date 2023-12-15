package cast

import (
	"testing"
)

func TestBoolCaster(t *testing.T) {
	t.Parallel()

	testCases := []casterTestCase[bool]{
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
