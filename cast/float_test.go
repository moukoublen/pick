package cast

import (
	"encoding/json"
	"fmt"
	"math"
	"testing"

	"github.com/moukoublen/pick/internal/testingx"
)

func TestFloatCaster(t *testing.T) {
	t.Parallel()
	type testCase struct {
		input           any
		float32CastErr  func(*testing.T, error)
		float64CastErr  func(*testing.T, error)
		expectedFloat64 float64
		expectedFloat32 float32
	}

	testsCases := []testCase{
		{
			input:           uint64(math.MaxUint64),
			expectedFloat32: float32(math.MaxUint64),
			float32CastErr:  nil,
			expectedFloat64: float64(math.MaxUint64),
			float64CastErr:  nil,
		},
		{
			input:           "1.79769313486231570814527423731704356798070e+308",
			expectedFloat32: float32(math.Inf(1)),
			float32CastErr:  expectOverFlowError,
			expectedFloat64: float64(math.MaxFloat64),
			float64CastErr:  nil,
		},
		{
			input:           json.Number("-12.4"),
			expectedFloat32: -12.4,
			float32CastErr:  nil,
			expectedFloat64: -12.4,
			float64CastErr:  nil,
		},
		{
			input:           true,
			expectedFloat32: 1,
			float32CastErr:  nil,
			expectedFloat64: 1,
			float64CastErr:  nil,
		},
		{
			input:           "Bad Input",
			expectedFloat32: 0,
			float32CastErr:  expectMalformedSyntax,
			expectedFloat64: 0,
			float64CastErr:  expectMalformedSyntax,
		},
		{
			input:           float64(math.MaxFloat64),
			expectedFloat32: float32(math.Inf(+1)),
			float32CastErr:  expectOverFlowError,
			expectedFloat64: float64(math.MaxFloat64),
			float64CastErr:  nil,
		},
	}

	caster := newFloatCaster()
	for _, tc := range testsCases {
		tc := tc

		name := fmt.Sprintf("%T(%v)", tc.input, tc.input)

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got64, err64 := caster.AsFloat64(tc.input)
			testingx.AssertError(t, tc.float64CastErr, err64)
			if !compareFloat64(got64, tc.expectedFloat64) {
				t.Errorf("wrong returned value. Expected %f got %f", tc.expectedFloat64, got64)
			}

			got32, err32 := caster.AsFloat32(tc.input)
			testingx.AssertError(t, tc.float32CastErr, err32)
			if !compareFloat32(got32, tc.expectedFloat32) {
				t.Errorf("wrong returned value. Expected %f got %f", tc.expectedFloat32, got32)
			}
		})
	}
}
