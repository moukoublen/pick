package pick

import (
	"testing"
	"time"

	"github.com/moukoublen/pick/internal/tst"
)

func TestDurationConverter(t *testing.T) {
	t.Parallel()

	converter := NewDefaultConverter()

	testCases := []singleConvertTestCase[time.Duration]{
		{
			input:         "8ns",
			expected:      8 * time.Nanosecond,
			errorAsserter: tst.NoError,
			directConvertFn: func(input any) (time.Duration, error) {
				return converter.AsDurationWithConfig(DurationConvertConfig{}, input)
			},
		},
		{
			input:         "8μs",
			expected:      8 * time.Microsecond,
			errorAsserter: tst.NoError,
			directConvertFn: func(input any) (time.Duration, error) {
				return converter.AsDurationWithConfig(DurationConvertConfig{}, input)
			},
		},
		{
			input:         "8ms",
			expected:      8 * time.Millisecond,
			errorAsserter: tst.NoError,
			directConvertFn: func(input any) (time.Duration, error) {
				return converter.AsDurationWithConfig(DurationConvertConfig{}, input)
			},
		},
		{
			input:         "8s",
			expected:      8 * time.Second,
			errorAsserter: tst.NoError,
			directConvertFn: func(input any) (time.Duration, error) {
				return converter.AsDurationWithConfig(DurationConvertConfig{}, input)
			},
		},
		{
			input:         "8m",
			expected:      8 * time.Minute,
			errorAsserter: tst.NoError,
			directConvertFn: func(input any) (time.Duration, error) {
				return converter.AsDurationWithConfig(DurationConvertConfig{}, input)
			},
		},
		{
			input:         "8h",
			expected:      8 * time.Hour,
			errorAsserter: tst.NoError,
			directConvertFn: func(input any) (time.Duration, error) {
				return converter.AsDurationWithConfig(DurationConvertConfig{}, input)
			},
		},

		{
			input:         8,
			expected:      8 * time.Nanosecond,
			errorAsserter: tst.NoError,
			directConvertFn: func(input any) (time.Duration, error) {
				return converter.AsDurationWithConfig(DurationConvertConfig{}, input)
			},
		},
		{
			input:         8,
			expected:      8 * time.Nanosecond,
			errorAsserter: tst.NoError,
			directConvertFn: func(input any) (time.Duration, error) {
				return converter.AsDurationWithConfig(DurationConvertConfig{DurationConvertNumberFormat: DurationNumberNanoseconds}, input)
			},
		},
		{
			input:         8,
			expected:      8 * time.Microsecond,
			errorAsserter: tst.NoError,
			directConvertFn: func(input any) (time.Duration, error) {
				return converter.AsDurationWithConfig(DurationConvertConfig{DurationConvertNumberFormat: DurationNumberMicroseconds}, input)
			},
		},
		{
			input:         8,
			expected:      8 * time.Millisecond,
			errorAsserter: tst.NoError,
			directConvertFn: func(input any) (time.Duration, error) {
				return converter.AsDurationWithConfig(DurationConvertConfig{DurationConvertNumberFormat: DurationNumberMilliseconds}, input)
			},
		},
		{
			input:         8,
			expected:      8 * time.Second,
			errorAsserter: tst.NoError,
			directConvertFn: func(input any) (time.Duration, error) {
				return converter.AsDurationWithConfig(DurationConvertConfig{DurationConvertNumberFormat: DurationNumberSeconds}, input)
			},
		},
		{
			input:         8,
			expected:      8 * time.Minute,
			errorAsserter: tst.NoError,
			directConvertFn: func(input any) (time.Duration, error) {
				return converter.AsDurationWithConfig(DurationConvertConfig{DurationConvertNumberFormat: DurationNumberMinutes}, input)
			},
		},
		{
			input:         8,
			expected:      8 * time.Hour,
			errorAsserter: tst.NoError,
			directConvertFn: func(input any) (time.Duration, error) {
				return converter.AsDurationWithConfig(DurationConvertConfig{DurationConvertNumberFormat: DurationNumberHours}, input)
			},
		},
	}

	runSingleConvertTestCases[time.Duration](t, testCases, converter.AsDuration)
}
