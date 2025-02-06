package pick

import (
	"testing"
	"time"

	"github.com/moukoublen/pick/internal/tst"
)

func TestDurationCaster(t *testing.T) {
	t.Parallel()

	caster := NewDefaultCaster()

	testCases := []singleCastTestCase[time.Duration]{
		{
			input:         "8ns",
			expected:      8 * time.Nanosecond,
			errorAsserter: tst.NoError,
			directCastFn: func(input any) (time.Duration, error) {
				return caster.AsDurationWithConfig(DurationCastConfig{}, input)
			},
		},
		{
			input:         "8μs",
			expected:      8 * time.Microsecond,
			errorAsserter: tst.NoError,
			directCastFn: func(input any) (time.Duration, error) {
				return caster.AsDurationWithConfig(DurationCastConfig{}, input)
			},
		},
		{
			input:         "8ms",
			expected:      8 * time.Millisecond,
			errorAsserter: tst.NoError,
			directCastFn: func(input any) (time.Duration, error) {
				return caster.AsDurationWithConfig(DurationCastConfig{}, input)
			},
		},
		{
			input:         "8s",
			expected:      8 * time.Second,
			errorAsserter: tst.NoError,
			directCastFn: func(input any) (time.Duration, error) {
				return caster.AsDurationWithConfig(DurationCastConfig{}, input)
			},
		},
		{
			input:         "8m",
			expected:      8 * time.Minute,
			errorAsserter: tst.NoError,
			directCastFn: func(input any) (time.Duration, error) {
				return caster.AsDurationWithConfig(DurationCastConfig{}, input)
			},
		},
		{
			input:         "8h",
			expected:      8 * time.Hour,
			errorAsserter: tst.NoError,
			directCastFn: func(input any) (time.Duration, error) {
				return caster.AsDurationWithConfig(DurationCastConfig{}, input)
			},
		},

		{
			input:         8,
			expected:      8 * time.Nanosecond,
			errorAsserter: tst.NoError,
			directCastFn: func(input any) (time.Duration, error) {
				return caster.AsDurationWithConfig(DurationCastConfig{}, input)
			},
		},
		{
			input:         8,
			expected:      8 * time.Nanosecond,
			errorAsserter: tst.NoError,
			directCastFn: func(input any) (time.Duration, error) {
				return caster.AsDurationWithConfig(DurationCastConfig{DurationCastNumberFormat: DurationNumberNanoseconds}, input)
			},
		},
		{
			input:         8,
			expected:      8 * time.Microsecond,
			errorAsserter: tst.NoError,
			directCastFn: func(input any) (time.Duration, error) {
				return caster.AsDurationWithConfig(DurationCastConfig{DurationCastNumberFormat: DurationNumberMicroseconds}, input)
			},
		},
		{
			input:         8,
			expected:      8 * time.Millisecond,
			errorAsserter: tst.NoError,
			directCastFn: func(input any) (time.Duration, error) {
				return caster.AsDurationWithConfig(DurationCastConfig{DurationCastNumberFormat: DurationNumberMilliseconds}, input)
			},
		},
		{
			input:         8,
			expected:      8 * time.Second,
			errorAsserter: tst.NoError,
			directCastFn: func(input any) (time.Duration, error) {
				return caster.AsDurationWithConfig(DurationCastConfig{DurationCastNumberFormat: DurationNumberSeconds}, input)
			},
		},
		{
			input:         8,
			expected:      8 * time.Minute,
			errorAsserter: tst.NoError,
			directCastFn: func(input any) (time.Duration, error) {
				return caster.AsDurationWithConfig(DurationCastConfig{DurationCastNumberFormat: DurationNumberMinutes}, input)
			},
		},
		{
			input:         8,
			expected:      8 * time.Hour,
			errorAsserter: tst.NoError,
			directCastFn: func(input any) (time.Duration, error) {
				return caster.AsDurationWithConfig(DurationCastConfig{DurationCastNumberFormat: DurationNumberHours}, input)
			},
		},
	}

	runSingleCastTestCases[time.Duration](t, testCases, caster.AsDuration)
}
