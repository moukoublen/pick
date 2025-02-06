package pick

import (
	"testing"

	"github.com/moukoublen/pick/internal/tst"
)

func TestBoolCaster(t *testing.T) {
	t.Parallel()

	testCases := []singleCastTestCase[bool]{
		{
			input:         "true",
			expected:      true,
			errorAsserter: tst.NoError,
		},
		{
			input:         "false",
			expected:      false,
			errorAsserter: tst.NoError,
		},
	}

	caster := NewDefaultCaster()
	runSingleCastTestCases[bool](t, testCases, caster.AsBool)
}
