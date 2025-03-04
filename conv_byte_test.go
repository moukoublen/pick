package pick

import (
	"encoding/json"
	"testing"

	"github.com/moukoublen/pick/internal/tst"
)

func TestByteConverterSlice(t *testing.T) {
	t.Parallel()

	testCases := []singleConvertTestCase[[]byte]{
		{
			input:         byte(12),
			expected:      []byte{12},
			errorAsserter: tst.NoError,
		},
		{
			input:         []byte{12, 13, 14},
			expected:      []byte{12, 13, 14},
			errorAsserter: tst.NoError,
		},
		{
			input:         "Yoda",
			expected:      []byte{0x59, 0x6f, 0x64, 0x61},
			errorAsserter: tst.NoError,
		},
		{
			input:         json.RawMessage(`{}`),
			expected:      []byte{0x7b, 0x7d},
			errorAsserter: tst.NoError,
		},
	}

	converter := NewDefaultConverter()
	runSingleConvertTestCases[[]byte](t, testCases, converter.AsByteSlice)
}
