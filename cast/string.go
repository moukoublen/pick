package cast

import (
	"encoding/json"
	"strconv"
)

type stringCaster struct{}

func newStringCaster() stringCaster {
	return stringCaster{}
}

func (fc stringCaster) AsString(input any) (string, error) {
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
		return strconv.FormatFloat(float64(origin), 'E', 5, 32), nil
	case float64:
		return strconv.FormatFloat(origin, 'E', 5, 64), nil

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
		return castAttemptUsingReflect[string](input)
	}
}

func (fc stringCaster) AsStringSlice(input any) ([]string, error) {
	return castToSlice[string](input, fc.AsString)
}
