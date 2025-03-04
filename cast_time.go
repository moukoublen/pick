package pick

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/moukoublen/pick/iter"
)

type TimeConvertNumberFormat int

const (
	TimeConvertNumberFormatUnix TimeConvertNumberFormat = iota
	TimeConvertNumberFormatUnixMilli
	TimeConvertNumberFormatUnixMicro
)

type TimeConvertByteSliceFormat int

const (
	TimeConvertByteSliceFormatString TimeConvertByteSliceFormat = iota
	TimeConvertByteSliceFormatBinary
)

type TimeConvertConfig struct {
	ParseInLocation     *time.Location
	StringFormat        string
	PraseStringAsNumber bool
	NumberFormat        TimeConvertNumberFormat
	ByteSliceFormat     TimeConvertByteSliceFormat
}

func (cnf TimeConvertConfig) getStringFormat() string {
	if cnf.StringFormat == "" {
		return time.RFC3339Nano
	}
	return cnf.StringFormat
}

func (c DefaultConverter) AsTime(input any) (time.Time, error) {
	return c.AsTimeWithConfig(TimeConvertConfig{}, input)
}

func (c DefaultConverter) AsTimeWithConfig(config TimeConvertConfig, input any) (time.Time, error) {
	switch origin := input.(type) {
	case int:
		return c.timeFromInt64(config, int64(origin))
	case int8:
		return c.timeFromInt64(config, int64(origin))
	case int16:
		return c.timeFromInt64(config, int64(origin))
	case int32:
		return c.timeFromInt64(config, int64(origin))
	case int64:
		return c.timeFromInt64(config, origin)

	case uint:
		return c.AsTimeWithConfig(config, uint64(origin))
	case uint8:
		return c.AsTimeWithConfig(config, uint64(origin))
	case uint16:
		return c.AsTimeWithConfig(config, uint64(origin))
	case uint32:
		return c.AsTimeWithConfig(config, uint64(origin))
	case uint64:
		asInt64, err := c.AsInt64(origin)
		if err != nil {
			t, _ := c.timeFromInt64(config, asInt64) // best effort
			return t, err
		}
		return c.timeFromInt64(config, asInt64)

	case float32:
		return c.AsTimeWithConfig(config, float64(origin))
	case float64:
		converted, err := float64ToInt64(origin)
		tm, _ := c.AsTimeWithConfig(config, converted) // best effort
		return tm, err

	case string:
		return c.timeFromString(config, origin)

	case json.Number:
		n, err := origin.Int64()
		if err != nil {
			return time.Time{}, newConvertError(err, fmt.Errorf("error converting json number to number: %w", err))
		}
		return c.AsTimeWithConfig(config, n)

	case []byte:
		return c.timeFromByteSlice(config, origin)

	case bool:
		return time.Time{}, newConvertError(ErrConvertInvalidType, input)

	case nil:
		return time.Time{}, nil

	default:
		// try to convert to basic (in case input is ~basic)
		if basic, err := tryConvertToBasicType(input); err == nil {
			return c.AsTimeWithConfig(config, basic)
		}

		return tryReflectConvert[time.Time](input)
	}
}

func (c DefaultConverter) timeFromInt64(config TimeConvertConfig, origin int64) (time.Time, error) {
	var tm time.Time
	switch config.NumberFormat {
	case TimeConvertNumberFormatUnix:
		tm = time.Unix(origin, 0).UTC()
	case TimeConvertNumberFormatUnixMilli:
		tm = time.UnixMilli(origin).UTC()
	case TimeConvertNumberFormatUnixMicro:
		tm = time.UnixMicro(origin).UTC()
	default:
		return tm, newConvertError(ErrConvertInvalidType, origin)
	}
	return tm, nil
}

func (c DefaultConverter) timeFromString(config TimeConvertConfig, origin string) (time.Time, error) {
	if config.PraseStringAsNumber {
		n, err := strconv.ParseInt(origin, 10, 64)
		if err != nil {
			return time.Time{}, newConvertError(err, fmt.Errorf("error converting string to number: %w", err))
		}
		return c.AsTimeWithConfig(config, n)
	}
	var tm time.Time
	var err error
	if config.ParseInLocation != nil {
		tm, err = time.ParseInLocation(config.getStringFormat(), origin, config.ParseInLocation)
	} else {
		tm, err = time.Parse(config.getStringFormat(), origin)
	}
	if err != nil {
		return time.Time{}, newConvertError(err, origin)
	}
	return tm, nil
}

func (c DefaultConverter) timeFromByteSlice(config TimeConvertConfig, origin []byte) (time.Time, error) {
	switch config.ByteSliceFormat {
	case TimeConvertByteSliceFormatBinary:
		tm := time.Time{}
		err := tm.UnmarshalBinary(origin)
		if err != nil {
			return tm, newConvertError(err, origin)
		}
	case TimeConvertByteSliceFormatString:
		return c.AsTimeWithConfig(config, string(origin))
	}
	return time.Time{}, newConvertError(ErrConvertInvalidType, origin)
}

func (c DefaultConverter) AsTimeSlice(input any) ([]time.Time, error) {
	return c.AsTimeSliceWithConfig(TimeConvertConfig{}, input)
}

func (c DefaultConverter) AsTimeSliceWithConfig(config TimeConvertConfig, input any) ([]time.Time, error) {
	return iter.Map[time.Time](input, func(item any, _ iter.CollectionOpMeta) (time.Time, error) {
		return c.AsTimeWithConfig(config, item)
	})
}

var ErrTimeConvertConfig = errors.New("invalid time converter config")
