package pick

import (
	"encoding/json"
	"math"
	"strconv"

	"github.com/moukoublen/pick/iter"
)

func (c DefaultConverter) AsFloat64(input any) (float64, error) {
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
			return v, newConvertError(err, origin)
		}
		return v, nil
	case json.Number:
		return c.AsFloat64(string(origin))
	case []byte:
		return c.AsFloat64(string(origin))

	case bool:
		if origin {
			return 1, nil
		}
		return 0, nil

	case nil:
		return 0, nil

	default:
		// try to convert to basic (in case is ~basic)
		if basic, err := tryConvertToBasicType(input); err == nil {
			return c.AsFloat64(basic)
		}

		return tryReflectConvert[float64](input)
	}
}

func (c DefaultConverter) AsFloat32(input any) (float32, error) {
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
			return float32(origin), newConvertError(ErrConvertOverFlow, input)
		}
		return float32(origin), nil

	case string:
		v, err := strconv.ParseFloat(origin, 32)
		if err != nil {
			return float32(v), newConvertError(err, origin)
		}
		return float32(v), nil
	case json.Number:
		return c.AsFloat32(string(origin))
	case []byte:
		return c.AsFloat32(string(origin))

	case bool:
		if origin {
			return 1, nil
		}
		return 0, nil

	case nil:
		return 0, nil

	default:
		// try to convert to basic (in case input is ~basic)
		if basic, err := tryConvertToBasicType(input); err == nil {
			return c.AsFloat32(basic)
		}

		return tryReflectConvert[float32](input)
	}
}

func (c DefaultConverter) AsFloat32Slice(input any) ([]float32, error) {
	return iter.Map(input, iter.MapOpFn(c.AsFloat32))
}

func (c DefaultConverter) AsFloat64Slice(input any) ([]float64, error) {
	return iter.Map(input, iter.MapOpFn(c.AsFloat64))
}
