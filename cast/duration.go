package cast

import (
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/moukoublen/pick/slices"
)

type DurationCastNumberFormat int

const (
	DurationNumberNanoseconds DurationCastNumberFormat = iota // default value
	DurationNumberMilliseconds
	DurationNumberMicroseconds
	DurationNumberSeconds
	DurationNumberMinutes
	DurationNumberHours
)

type DurationCastConfig struct {
	DurationCastNumberFormat DurationCastNumberFormat
}

func (c Caster) AsDuration(input any) (time.Duration, error) {
	return c.AsDurationWithConfig(DurationCastConfig{}, input)
}

func (c Caster) AsDurationWithConfig(config DurationCastConfig, input any) (time.Duration, error) {
	switch origin := input.(type) {
	case int:
		return c.durationFromInt64(config, int64(origin))
	case int8:
		return c.durationFromInt64(config, int64(origin))
	case int16:
		return c.durationFromInt64(config, int64(origin))
	case int32:
		return c.durationFromInt64(config, int64(origin))
	case int64:
		return c.durationFromInt64(config, origin)

	case uint:
		return c.AsDurationWithConfig(config, uint64(origin))
	case uint8:
		return c.AsDurationWithConfig(config, uint64(origin))
	case uint16:
		return c.AsDurationWithConfig(config, uint64(origin))
	case uint32:
		return c.AsDurationWithConfig(config, uint64(origin))
	case uint64:
		asInt64, err := c.AsInt64(origin)
		if err != nil {
			d, _ := c.AsDurationWithConfig(config, asInt64) // best effort
			return d, err
		}
		return c.AsDurationWithConfig(config, asInt64)

	case float32:
		return c.AsDurationWithConfig(config, float64(origin))
	case float64:
		casted, err := float64ToInt64(origin)
		d, _ := c.AsDurationWithConfig(config, casted) // best effort
		return d, err

	case string:
		d, err := time.ParseDuration(origin)
		if err != nil {
			return 0, newCastError(err, origin)
		}
		return d, nil

	case json.Number:
		n, err := origin.Int64()
		if err != nil {
			return time.Duration(0), newCastError(err, fmt.Errorf("error converting json number to number: %w", err))
		}
		return c.AsDurationWithConfig(config, n)

	case []byte:
		return c.AsDurationWithConfig(config, string(origin))

	case bool:
		return time.Duration(0), newCastError(ErrInvalidType, input)

	case nil:
		return time.Duration(0), nil

	case time.Duration:
		return origin, nil

	default:
		// try to cast to basic (in case input is ~basic)
		if basic, err := tryCastToBasicType(input); err == nil {
			return c.AsDurationWithConfig(config, basic)
		}

		return tryReflectConvert[time.Duration](input)
	}
}

func (c Caster) durationFromInt64(config DurationCastConfig, origin int64) (time.Duration, error) {
	limitCheck := func(d time.Duration) error {
		if origin >= (math.MinInt64/int64(d)) && origin <= (math.MaxInt64/int64(d)) {
			return nil
		}
		return newCastError(ErrCastOverFlow, origin)
	}

	var dr time.Duration
	var err error

	switch config.DurationCastNumberFormat {
	case DurationNumberSeconds:
		dr = time.Duration(origin) * time.Second
		err = limitCheck(time.Second)
	case DurationNumberMilliseconds:
		dr = time.Duration(origin) * time.Millisecond
		err = limitCheck(time.Millisecond)
	case DurationNumberMicroseconds:
		dr = time.Duration(origin) * time.Microsecond
		err = limitCheck(time.Microsecond)
	case DurationNumberNanoseconds:
		dr = time.Duration(origin) * time.Nanosecond
		err = limitCheck(time.Nanosecond)
	case DurationNumberMinutes:
		dr = time.Duration(origin) * time.Minute
		err = limitCheck(time.Minute)
	case DurationNumberHours:
		dr = time.Duration(origin) * time.Hour
		err = limitCheck(time.Hour)
	default:
		return dr, newCastError(ErrInvalidType, origin)
	}
	return dr, err
}

func (c Caster) AsDurationSlice(input any) ([]time.Duration, error) {
	return c.AsDurationSliceWithConfig(DurationCastConfig{}, input)
}

func (c Caster) AsDurationSliceWithConfig(config DurationCastConfig, input any) ([]time.Duration, error) {
	return slices.Map(input, func(item any, _ slices.OpMeta) (time.Duration, error) {
		return c.AsDurationWithConfig(config, item)
	})
}
