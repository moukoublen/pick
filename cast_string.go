package pick

import (
	"encoding/json"
	"strconv"

	"github.com/moukoublen/pick/iter"
)

func (c DefaultCaster) AsString(input any) (string, error) {
	switch origin := input.(type) {
	case int:
		return strconv.FormatInt(int64(origin), 10), nil
	case int8:
		return strconv.FormatInt(int64(origin), 10), nil
	case int16:
		return strconv.FormatInt(int64(origin), 10), nil
	case int32:
		return strconv.FormatInt(int64(origin), 10), nil
	case int64:
		return strconv.FormatInt(origin, 10), nil

	case uint:
		return strconv.FormatUint(uint64(origin), 10), nil
	case uint8:
		return strconv.FormatUint(uint64(origin), 10), nil
	case uint16:
		return strconv.FormatUint(uint64(origin), 10), nil
	case uint32:
		return strconv.FormatUint(uint64(origin), 10), nil
	case uint64:
		return strconv.FormatUint(origin, 10), nil

	case float32:
		return strconv.FormatFloat(float64(origin), 'G', 5, 32), nil
	case float64:
		return strconv.FormatFloat(origin, 'G', 5, 64), nil

	case string:
		return origin, nil
	case json.Number:
		return string(origin), nil
	case json.RawMessage:
		return string(origin), nil
	case []byte:
		return string(origin), nil

	case bool:
		return strconv.FormatBool(origin), nil

	case nil:
		return "", nil

	default:
		// try to cast to basic (in case input is ~basic)
		if basic, err := tryCastToBasicType(input); err == nil {
			return c.AsString(basic)
		}

		return tryReflectConvert[string](input)
	}
}

func (c DefaultCaster) AsStringSlice(input any) ([]string, error) {
	return iter.Map(input, iter.MapOpFn(c.AsString))
}
