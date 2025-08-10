package pick

import (
	"testing"
	"time"
	_ "time/tzdata"

	"github.com/ifnotnil/x/tst"
)

func TestTimeConverter(t *testing.T) {
	t.Parallel()

	converter := NewDefaultConverter()

	tzPlus4, _ := time.LoadLocation("Etc/GMT-4")
	tzMinus7, _ := time.LoadLocation("Etc/GMT+7")
	tzPlus8, _ := time.LoadLocation("Etc/GMT-8")
	tzAthens, _ := time.LoadLocation("Europe/Athens")

	type int64Alias int64
	type stringAlias string
	testCases := []singleConvertTestCase[time.Time]{
		{
			input:         int64Alias(1700000000),
			expected:      time.Date(2023, time.November, 14, 22, 13, 20, 0, time.UTC),
			errorAsserter: tst.NoError(),
		},
		{
			input:         int64(1700000000 * 1000),
			expected:      time.Date(2023, time.November, 14, 22, 13, 20, 0, time.UTC),
			errorAsserter: tst.NoError(),
			directConvertFn: func(input any) (time.Time, error) {
				return converter.AsTimeWithConfig(TimeConvertConfig{NumberFormat: TimeConvertNumberFormatUnixMilli}, input)
			},
		},
		{
			input:         int64(1700000000 * 1000 * 1000),
			expected:      time.Date(2023, time.November, 14, 22, 13, 20, 0, time.UTC),
			errorAsserter: tst.NoError(),
			directConvertFn: func(input any) (time.Time, error) {
				return converter.AsTimeWithConfig(TimeConvertConfig{NumberFormat: TimeConvertNumberFormatUnixMicro}, input)
			},
		},
		{
			input:         int32(12),
			expected:      time.Date(1970, time.January, 1, 0, 0, 12, 0, time.UTC),
			errorAsserter: tst.NoError(),
		},
		{
			input:         int8(12),
			expected:      time.Date(1970, time.January, 1, 0, 0, 12, 0, time.UTC),
			errorAsserter: tst.NoError(),
		},
		{
			input:         stringAlias("abcd"),
			expected:      time.Time{},
			errorAsserter: tst.ErrorOfType[*time.ParseError](),
		},
		{
			input:         stringAlias("2023-11-14T15:04:05+04:00"),
			expected:      time.Date(2023, time.November, 14, 15, 4, 5, 0, tzPlus4),
			errorAsserter: tst.NoError(),
		},
		{
			input:         "2023-11-14T15:04:05+08:00",
			expected:      time.Date(2023, time.November, 14, 15, 4, 5, 0, tzPlus8),
			errorAsserter: tst.NoError(),
		},
		{
			input:         "2023-11-14T15:04:05Z",
			expected:      time.Date(2023, time.November, 14, 15, 4, 5, 0, time.UTC),
			errorAsserter: tst.NoError(),
		},
		{
			input:         "2023-11-14T15:04:05.12Z",
			expected:      time.Date(2023, time.November, 14, 15, 4, 5, 120000000, time.UTC),
			errorAsserter: tst.NoError(),
		},
		{
			input:         "1700000000000",
			expected:      time.Date(2023, time.November, 14, 22, 13, 20, 0, time.UTC),
			errorAsserter: tst.NoError(),
			directConvertFn: func(input any) (time.Time, error) {
				return converter.AsTimeWithConfig(TimeConvertConfig{PraseStringAsNumber: true, NumberFormat: TimeConvertNumberFormatUnixMilli}, input)
			},
		},
		{
			input:         "Mon, 02 Jan 2006 15:04:05 -0700",
			expected:      time.Date(2006, time.January, 2, 15, 4, 5, 0, tzMinus7),
			errorAsserter: tst.NoError(),
			directConvertFn: func(input any) (time.Time, error) {
				return converter.AsTimeWithConfig(TimeConvertConfig{StringFormat: time.RFC1123Z}, input)
			},
		},
		{
			input:         "Mon Jan 2 15:04:05 2016",
			expected:      time.Date(2016, time.January, 2, 15, 4, 5, 0, tzAthens),
			errorAsserter: tst.NoError(),
			directConvertFn: func(input any) (time.Time, error) {
				return converter.AsTimeWithConfig(TimeConvertConfig{StringFormat: time.ANSIC, ParseInLocation: tzAthens}, input)
			},
		},
	}

	runSingleConvertTestCases[time.Time](t, testCases, converter.AsTime)
}

func TestTimeSliceConverter(t *testing.T) {
	t.Parallel()
	converter := NewDefaultConverter()

	tzPlus4, _ := time.LoadLocation("Etc/GMT-4")
	tzPlus8, _ := time.LoadLocation("Etc/GMT-8")

	testCases := []singleConvertTestCase[[]time.Time]{
		{
			input:         int64(1700000000 * 1000),
			expected:      []time.Time{time.Date(2023, time.November, 14, 22, 13, 20, 0, time.UTC)},
			errorAsserter: tst.NoError(),
			directConvertFn: func(input any) ([]time.Time, error) {
				return converter.AsTimeSliceWithConfig(TimeConvertConfig{NumberFormat: TimeConvertNumberFormatUnixMilli}, input)
			},
		},
		{
			input:         int64(1700000000),
			expected:      []time.Time{time.Date(2023, time.November, 14, 22, 13, 20, 0, time.UTC)},
			errorAsserter: tst.NoError(),
		},
		{
			input: []int64{int64(1700000000), int64(1700000001), int64(1700000002), int64(1700000003)},
			expected: []time.Time{
				time.Date(2023, time.November, 14, 22, 13, 20, 0, time.UTC),
				time.Date(2023, time.November, 14, 22, 13, 21, 0, time.UTC),
				time.Date(2023, time.November, 14, 22, 13, 22, 0, time.UTC),
				time.Date(2023, time.November, 14, 22, 13, 23, 0, time.UTC),
			},
			errorAsserter: tst.NoError(),
		},
		{
			input: []any{int64(1700000000), int64(1700000001), int64(1700000002), int64(1700000003), "2023-11-14T15:04:05+08:00"},
			expected: []time.Time{
				time.Date(2023, time.November, 14, 22, 13, 20, 0, time.UTC),
				time.Date(2023, time.November, 14, 22, 13, 21, 0, time.UTC),
				time.Date(2023, time.November, 14, 22, 13, 22, 0, time.UTC),
				time.Date(2023, time.November, 14, 22, 13, 23, 0, time.UTC),
				time.Date(2023, time.November, 14, 15, 4, 5, 0, tzPlus8),
			},
			errorAsserter: tst.NoError(),
		},
		{
			input: []string{"2023-11-14T15:04:05+08:00", "2023-11-14T15:04:05+04:00", "2023-11-14T15:04:05.65Z"},
			expected: []time.Time{
				time.Date(2023, time.November, 14, 15, 4, 5, 0, tzPlus8),
				time.Date(2023, time.November, 14, 15, 4, 5, 0, tzPlus4),
				time.Date(2023, time.November, 14, 15, 4, 5, 650000000, time.UTC),
			},
			errorAsserter: tst.NoError(),
		},
	}

	runSingleConvertTestCases(t, testCases, converter.AsTimeSlice)
}
