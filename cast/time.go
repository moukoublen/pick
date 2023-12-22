package cast

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"
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

type timeCaster struct{}

func newTimeCaster() timeCaster { return timeCaster{} }

func (tc timeCaster) AsTime(input any) (time.Time, error) {
	return tc.AsTimeWithConfig(input, TimeCastConfig{})
}

func (tc timeCaster) AsTimeWithConfig(input any, config TimeCastConfig) (time.Time, error) {
	switch origin := input.(type) {
	case int:
		return tc.AsTime(int64(origin))
	case int8:
		return tc.AsTime(int64(origin))
	case int16:
		return tc.AsTime(int64(origin))
	case int32:
		return tc.AsTime(int64(origin))
	case int64:
		var tm time.Time
		switch config.NumberFormat {
		case TimeCastNumberFormatUnix:
			tm = time.Unix(origin, 0).UTC()
		case TimeCastNumberFormatUnixMilli:
			tm = time.UnixMilli(origin).UTC()
		case TimeCastNumberFormatUnixMicro:
			tm = time.UnixMicro(origin).UTC()
		default:
			return tm, newCastError(ErrInvalidType, input)
		}
		return tm, nil

	case uint:
		return tc.AsTime(int64(origin))
	case uint8:
		return tc.AsTime(int64(origin))
	case uint16:
		return tc.AsTime(int64(origin))
	case uint32:
		return tc.AsTime(int64(origin))
	case uint64:
		return tc.AsTime(int64(origin))

	case float32:
		return tc.AsTime(float64(origin))
	case float64:
		casted, err := float64ToInt64(origin)
		tm, _ := tc.AsTime(casted) // best effort
		return tm, err

	case string:
		if config.PraseStringAsNumber {
			n, err := strconv.ParseInt(origin, 10, 64)
			if err != nil {
				return time.Time{}, newCastError(err, fmt.Errorf("error converting string to number: %w", err))
			}
			return tc.AsTimeWithConfig(n, config)
		}
		var tm time.Time
		var err error
		if config.ParseInLocation != nil {
			tm, err = time.ParseInLocation(config.getStringFormat(), origin, config.ParseInLocation)
		} else {
			tm, err = time.Parse(config.getStringFormat(), origin)
		}
		if err != nil {
			return time.Time{}, newCastError(err, input)
		}
		return tm, nil

	case json.Number:
		n, err := origin.Int64()
		if err != nil {
			return time.Time{}, newCastError(err, fmt.Errorf("error converting json number to number: %w", err))
		}
		return tc.AsTime(n)

	case []byte:
		switch config.ByteSliceFormat {
		case TimeCastByteSliceFormatBinary:
			tm := time.Time{}
			err := tm.UnmarshalBinary(origin)
			if err != nil {
				return tm, newCastError(err, input)
			}
		case TimeCastByteSliceFormatString:
			return tc.AsTime(string(origin))
		}
		return time.Time{}, newCastError(ErrInvalidType, input)

	case bool:
		return time.Time{}, newCastError(ErrInvalidType, input)

	case nil:
		return time.Time{}, nil

	default:
		// try to cast to basic (in case input is ~basic)
		if basic, err := tryCastToBasicType(input); err == nil {
			return tc.AsTimeWithConfig(basic, config)
		}

		return tryCastUsingReflect[time.Time](input)
	}
}

func (tc timeCaster) AsTimeSlice(input any) ([]time.Time, error) {
	return tc.AsTimeSliceWithConfig(input, TimeCastConfig{})
}

func (tc timeCaster) AsTimeSliceWithConfig(input any, config TimeCastConfig) ([]time.Time, error) {
	return ToSlice[time.Time](input, func(a any) (time.Time, error) {
		return tc.AsTimeWithConfig(a, config)
	})
}

var ErrTimeCastConfig = errors.New("invalid time caster config")
