package cast

import (
	"encoding/json"
	"testing"
)

func TestByteCasterSlice(t *testing.T) {
	t.Parallel()

	testCases := []casterTestCase[[]byte]{
		{
			input:       byte(12),
			expected:    []byte{12},
			expectedErr: nil,
		},
		{
			input:       []byte{12, 13, 14},
			expected:    []byte{12, 13, 14},
			expectedErr: nil,
		},
		{
			input:       "Yoda",
			expected:    []byte{0x59, 0x6f, 0x64, 0x61},
			expectedErr: nil,
		},
		{
			input:       json.RawMessage(`{}`),
			expected:    []byte{0x7b, 0x7d},
			expectedErr: nil,
		},
	}

	caster := newByteCaster()
	casterTest[[]byte](t, testCases, caster.AsByteSlice)
}
