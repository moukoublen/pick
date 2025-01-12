package pick

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/moukoublen/pick/iter"
)

type TimeCastNumberFormat int

const (
	TimeCastNumberFormatUnix TimeCastNumberFormat = iota
	TimeCastNumberFormatUnixMilli
	TimeCastNumberFormatUnixMicro
)

type TimeCastByteSliceFormat int

const (
	TimeCastByteSliceFormatString TimeCastByteSliceFormat = iota
	TimeCastByteSliceFormatBinary
)

type TimeCastConfig struct {
	ParseInLocation     *time.Location
	StringFormat        string
	PraseStringAsNumber bool
	NumberFormat        TimeCastNumberFormat
	ByteSliceFormat     TimeCastByteSliceFormat
}

func (cnf TimeCastConfig) getStringFormat() string {
	if cnf.StringFormat == "" {
		return time.RFC3339Nano
	}
	return cnf.StringFormat
}

func (c DefaultCaster) AsTime(input any) (time.Time, error) {
	return c.AsTimeWithConfig(TimeCastConfig{}, input)
}

func (c DefaultCaster) AsTimeWithConfig(config TimeCastConfig, input any) (time.Time, error) {
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
		casted, err := float64ToInt64(origin)
		tm, _ := c.AsTimeWithConfig(config, casted) // best effort
		return tm, err

	case string:
		return c.timeFromString(config, origin)

	case json.Number:
		n, err := origin.Int64()
		if err != nil {
			return time.Time{}, newCastError(err, fmt.Errorf("error converting json number to number: %w", err))
		}
		return c.AsTimeWithConfig(config, n)

	case []byte:
		return c.timeFromByteSlice(config, origin)

	case bool:
		return time.Time{}, newCastError(ErrCastInvalidType, input)

	case nil:
		return time.Time{}, nil

	default:
		// try to cast to basic (in case input is ~basic)
		if basic, err := tryCastToBasicType(input); err == nil {
			return c.AsTimeWithConfig(config, basic)
		}

		return tryReflectConvert[time.Time](input)
	}
}

func (c DefaultCaster) timeFromInt64(config TimeCastConfig, origin int64) (time.Time, error) {
	var tm time.Time
	switch config.NumberFormat {
	case TimeCastNumberFormatUnix:
		tm = time.Unix(origin, 0).UTC()
	case TimeCastNumberFormatUnixMilli:
		tm = time.UnixMilli(origin).UTC()
	case TimeCastNumberFormatUnixMicro:
		tm = time.UnixMicro(origin).UTC()
	default:
		return tm, newCastError(ErrCastInvalidType, origin)
	}
	return tm, nil
}

func (c DefaultCaster) timeFromString(config TimeCastConfig, origin string) (time.Time, error) {
	if config.PraseStringAsNumber {
		n, err := strconv.ParseInt(origin, 10, 64)
		if err != nil {
			return time.Time{}, newCastError(err, fmt.Errorf("error converting string to number: %w", err))
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
		return time.Time{}, newCastError(err, origin)
	}
	return tm, nil
}

func (c DefaultCaster) timeFromByteSlice(config TimeCastConfig, origin []byte) (time.Time, error) {
	switch config.ByteSliceFormat {
	case TimeCastByteSliceFormatBinary:
		tm := time.Time{}
		err := tm.UnmarshalBinary(origin)
		if err != nil {
			return tm, newCastError(err, origin)
		}
	case TimeCastByteSliceFormatString:
		return c.AsTimeWithConfig(config, string(origin))
	}
	return time.Time{}, newCastError(ErrCastInvalidType, origin)
}

func (c DefaultCaster) AsTimeSlice(input any) ([]time.Time, error) {
	return c.AsTimeSliceWithConfig(TimeCastConfig{}, input)
}

func (c DefaultCaster) AsTimeSliceWithConfig(config TimeCastConfig, input any) ([]time.Time, error) {
	return iter.Map[time.Time](input, func(item any, _ iter.OpMeta) (time.Time, error) {
		return c.AsTimeWithConfig(config, item)
	})
}

var ErrTimeCastConfig = errors.New("invalid time caster config")
