package pick

import (
	"testing"

	"github.com/ifnotnil/x/tst"
)

func TestBoolConverter(t *testing.T) {
	t.Parallel()

	testCases := []singleConvertTestCase[bool]{
		{
			input:         "true",
			expected:      true,
			errorAsserter: tst.NoError(),
		},
		{
			input:         "false",
			expected:      false,
			errorAsserter: tst.NoError(),
		},
	}

	converter := NewDefaultConverter()
	runSingleConvertTestCases[bool](t, testCases, converter.AsBool)
}
