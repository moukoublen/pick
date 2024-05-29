package cast

import (
	"testing"
	"time"
)

func TestDurationCaster(t *testing.T) {
	t.Parallel()

	caster := newDurationCaster()

	testCases := []singleCastTestCase[time.Duration]{
		{
			input:       "8ns",
			expected:    8 * time.Nanosecond,
			expectedErr: nil,
			directCastFn: func(input any) (time.Duration, error) {
				return caster.AsDurationWithConfig(DurationCastConfig{}, input)
			},
		},
		{
			input:       "8μs",
			expected:    8 * time.Microsecond,
			expectedErr: nil,
			directCastFn: func(input any) (time.Duration, error) {
				return caster.AsDurationWithConfig(DurationCastConfig{}, input)
			},
		},
		{
			input:       "8ms",
			expected:    8 * time.Millisecond,
			expectedErr: nil,
			directCastFn: func(input any) (time.Duration, error) {
				return caster.AsDurationWithConfig(DurationCastConfig{}, input)
			},
		},
		{
			input:       "8s",
			expected:    8 * time.Second,
			expectedErr: nil,
			directCastFn: func(input any) (time.Duration, error) {
				return caster.AsDurationWithConfig(DurationCastConfig{}, input)
			},
		},
		{
			input:       "8m",
			expected:    8 * time.Minute,
			expectedErr: nil,
			directCastFn: func(input any) (time.Duration, error) {
				return caster.AsDurationWithConfig(DurationCastConfig{}, input)
			},
		},
		{
			input:       "8h",
			expected:    8 * time.Hour,
			expectedErr: nil,
			directCastFn: func(input any) (time.Duration, error) {
				return caster.AsDurationWithConfig(DurationCastConfig{}, input)
			},
		},

		{
			input:       8,
			expected:    8 * time.Nanosecond,
			expectedErr: nil,
			directCastFn: func(input any) (time.Duration, error) {
				return caster.AsDurationWithConfig(DurationCastConfig{}, input)
			},
		},
		{
			input:       8,
			expected:    8 * time.Nanosecond,
			expectedErr: nil,
			directCastFn: func(input any) (time.Duration, error) {
				return caster.AsDurationWithConfig(DurationCastConfig{DurationCastNumberFormat: DurationNumberNanoseconds}, input)
			},
		},
		{
			input:       8,
			expected:    8 * time.Microsecond,
			expectedErr: nil,
			directCastFn: func(input any) (time.Duration, error) {
				return caster.AsDurationWithConfig(DurationCastConfig{DurationCastNumberFormat: DurationNumberMicroseconds}, input)
			},
		},
		{
			input:       8,
			expected:    8 * time.Millisecond,
			expectedErr: nil,
			directCastFn: func(input any) (time.Duration, error) {
				return caster.AsDurationWithConfig(DurationCastConfig{DurationCastNumberFormat: DurationNumberMilliseconds}, input)
			},
		},
		{
			input:       8,
			expected:    8 * time.Second,
			expectedErr: nil,
			directCastFn: func(input any) (time.Duration, error) {
				return caster.AsDurationWithConfig(DurationCastConfig{DurationCastNumberFormat: DurationNumberSeconds}, input)
			},
		},
		{
			input:       8,
			expected:    8 * time.Minute,
			expectedErr: nil,
			directCastFn: func(input any) (time.Duration, error) {
				return caster.AsDurationWithConfig(DurationCastConfig{DurationCastNumberFormat: DurationNumberMinutes}, input)
			},
		},
		{
			input:       8,
			expected:    8 * time.Hour,
			expectedErr: nil,
			directCastFn: func(input any) (time.Duration, error) {
				return caster.AsDurationWithConfig(DurationCastConfig{DurationCastNumberFormat: DurationNumberHours}, input)
			},
		},
	}

	runSingleCastTestCases[time.Duration](t, testCases, caster.AsDuration)
}
