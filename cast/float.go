package cast

import (
	"encoding/json"
	"math"
	"strconv"
)

type floatCaster struct{}

func newFloatCaster() floatCaster {
	return floatCaster{}
}

func (fc floatCaster) AsFloat64(input any) (float64, error) {
	switch origin := input.(type) {
	case int:
		return float64(origin), nil
	case int8:
		return float64(origin), nil
	case int16:
		return float64(origin), nil
	case int32:
		return float64(origin), nil
	case int64:
		return float64(origin), nil

	case uint:
		return float64(origin), nil
	case uint8:
		return float64(origin), nil
	case uint16:
		return float64(origin), nil
	case uint32:
		return float64(origin), nil
	case uint64:
		return float64(origin), nil

	case float32:
		return float64(origin), nil
	case float64:
		return origin, nil

	case string:
		v, err := strconv.ParseFloat(origin, 64)
		if err != nil {
			return v, newCastError(err, origin)
		}
		return v, nil
	case json.Number:
		return fc.AsFloat64(string(origin))

	case bool:
		if origin {
			return 1, nil
		}
		return 0, nil

	case nil:
		return 0, nil

	default:
		return castAttemptUsingReflect[float64](input)
	}
}

func (fc floatCaster) AsFloat32(input any) (float32, error) {
	switch origin := input.(type) {
	case int:
		return float32(origin), nil
	case int8:
		return float32(origin), nil
	case int16:
		return float32(origin), nil
	case int32:
		return float32(origin), nil
	case int64:
		return float32(origin), nil

	case uint:
		return float32(origin), nil
	case uint8:
		return float32(origin), nil
	case uint16:
		return float32(origin), nil
	case uint32:
		return float32(origin), nil
	case uint64:
		return float32(origin), nil

	case float32:
		return origin, nil
	case float64:
		if origin > float64(math.MaxFloat32) || origin < (-float64(math.MaxFloat32)) {
			return float32(origin), newCastError(ErrCastOverFlow, input)
		}
		return float32(origin), nil

	case string:
		v, err := strconv.ParseFloat(origin, 32)
		if err != nil {
			return float32(v), newCastError(err, origin)
		}
		return float32(v), nil
	case json.Number:
		return fc.AsFloat32(string(origin))

	case bool:
		if origin {
			return 1, nil
		}
		return 0, nil

	case nil:
		return 0, nil

	default:
		return castAttemptUsingReflect[float32](input)
	}
}

func (fc floatCaster) AsFloat32Slice(input any) ([]float32, error) {
	return castToSlice[float32](input, fc.AsFloat32)
}

func (fc floatCaster) AsFloat64Slice(input any) ([]float64, error) {
	return castToSlice[float64](input, fc.AsFloat64)
}
