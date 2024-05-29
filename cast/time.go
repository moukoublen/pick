package cast

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/moukoublen/pick/cast/slices"
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
	return tc.AsTimeWithConfig(TimeCastConfig{}, input)
}

func (tc timeCaster) AsTimeWithConfig(config TimeCastConfig, input any) (time.Time, error) {
	switch origin := input.(type) {
	case int:
		return tc.AsTimeWithConfig(config, int64(origin))
	case int8:
		return tc.AsTimeWithConfig(config, int64(origin))
	case int16:
		return tc.AsTimeWithConfig(config, int64(origin))
	case int32:
		return tc.AsTimeWithConfig(config, int64(origin))
	case int64:
		return tc.fromInt64(config, origin)

	case uint:
		return tc.AsTimeWithConfig(config, int64(origin))
	case uint8:
		return tc.AsTimeWithConfig(config, int64(origin))
	case uint16:
		return tc.AsTimeWithConfig(config, int64(origin))
	case uint32:
		return tc.AsTimeWithConfig(config, int64(origin))
	case uint64:
		return tc.fromUint64(config, origin)

	case float32:
		return tc.AsTimeWithConfig(config, float64(origin))
	case float64:
		casted, err := float64ToInt64(origin)
		tm, _ := tc.AsTimeWithConfig(config, casted) // best effort
		return tm, err

	case string:
		return tc.fromString(config, origin)

	case json.Number:
		n, err := origin.Int64()
		if err != nil {
			return time.Time{}, newCastError(err, fmt.Errorf("error converting json number to number: %w", err))
		}
		return tc.AsTimeWithConfig(config, n)

	case []byte:
		return tc.fromByteSlice(config, origin)

	case bool:
		return time.Time{}, newCastError(ErrInvalidType, input)

	case nil:
		return time.Time{}, nil

	default:
		// try to cast to basic (in case input is ~basic)
		if basic, err := tryCastToBasicType(input); err == nil {
			return tc.AsTimeWithConfig(config, basic)
		}

		return tryReflectConvert[time.Time](input)
	}
}

func (tc timeCaster) fromInt64(config TimeCastConfig, origin int64) (time.Time, error) {
	var tm time.Time
	switch config.NumberFormat {
	case TimeCastNumberFormatUnix:
		tm = time.Unix(origin, 0).UTC()
	case TimeCastNumberFormatUnixMilli:
		tm = time.UnixMilli(origin).UTC()
	case TimeCastNumberFormatUnixMicro:
		tm = time.UnixMicro(origin).UTC()
	default:
		return tm, newCastError(ErrInvalidType, origin)
	}
	return tm, nil
}

func (tc timeCaster) fromUint64(config TimeCastConfig, origin uint64) (time.Time, error) {
	if !uint64CastValid(origin, reflect.Int64) {
		d, _ := tc.AsTimeWithConfig(config, int64(origin))
		return d, newCastError(ErrCastOverFlow, origin)
	}
	return tc.AsTimeWithConfig(config, int64(origin))
}

func (tc timeCaster) fromString(config TimeCastConfig, origin string) (time.Time, error) {
	if config.PraseStringAsNumber {
		n, err := strconv.ParseInt(origin, 10, 64)
		if err != nil {
			return time.Time{}, newCastError(err, fmt.Errorf("error converting string to number: %w", err))
		}
		return tc.AsTimeWithConfig(config, n)
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

func (tc timeCaster) fromByteSlice(config TimeCastConfig, origin []byte) (time.Time, error) {
	switch config.ByteSliceFormat {
	case TimeCastByteSliceFormatBinary:
		tm := time.Time{}
		err := tm.UnmarshalBinary(origin)
		if err != nil {
			return tm, newCastError(err, origin)
		}
	case TimeCastByteSliceFormatString:
		return tc.AsTimeWithConfig(config, string(origin))
	}
	return time.Time{}, newCastError(ErrInvalidType, origin)
}

func (tc timeCaster) AsTimeSlice(input any) ([]time.Time, error) {
	return tc.AsTimeSliceWithConfig(TimeCastConfig{}, input)
}

func (tc timeCaster) AsTimeSliceWithConfig(config TimeCastConfig, input any) ([]time.Time, error) {
	return slices.AsSlice[time.Time](input, func(item any, _ slices.OpMeta) (time.Time, error) {
		return tc.AsTimeWithConfig(config, item)
	})
}

var ErrTimeCastConfig = errors.New("invalid time caster config")
