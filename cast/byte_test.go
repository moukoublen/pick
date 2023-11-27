package cast

import (
	"encoding/json"
	"math"
	"testing"
)

func TestByteCaster(t *testing.T) {
	t.Parallel()

	testCases := []casterTestCase[byte]{
		{
			input:       byte(12),
			expected:    12,
			expectedErr: nil,
		},
		{
			input:       int32(12),
			expected:    12,
			expectedErr: nil,
		},
		{
			input:       int8(math.MaxInt8),
			expected:    127,
			expectedErr: nil,
		},
		{
			input:       int16(128),
			expected:    128,
			expectedErr: nil,
		},
		{
			input:       uint64(math.MaxUint64),
			expected:    math.MaxUint8,
			expectedErr: expectOverFlowError,
		},
		{
			input:       float64(123),
			expected:    123,
			expectedErr: nil,
		},
		{
			input:       float32(123.001),
			expected:    123,
			expectedErr: expectLostDecimals,
		},
		{
			input:       string("1"),
			expected:    0,
			expectedErr: expectInvalidType,
		},
	}

	caster := newByteCaster()
	casterTest[byte](t, testCases, caster.AsByte)
}

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
