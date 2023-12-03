package cast

import (
	"encoding/json"
	"strconv"
)

type boolCaster struct{}

func newBoolCaster() boolCaster { return boolCaster{} }

func (bc boolCaster) AsBool(input any) (bool, error) {
	switch origin := input.(type) {
	case int:
		return origin != 0, nil
	case int8:
		return origin != 0, nil
	case int16:
		return origin != 0, nil
	case int32:
		return origin != 0, nil
	case int64:
		return origin != 0, nil

	case uint:
		return origin != 0, nil
	case uint8:
		return origin != 0, nil
	case uint16:
		return origin != 0, nil
	case uint32:
		return origin != 0, nil
	case uint64:
		return origin != 0, nil

	case float32:
		return origin != 0, nil
	case float64:
		return origin != 0, nil

	case string:
		b, err := strconv.ParseBool(origin)
		if err != nil {
			return false, newCastError(err, input)
		}
		return b, nil
	case json.Number:
		n, err := origin.Float64()
		if err != nil {
			return false, newCastError(err, input)
		}
		return bc.AsBool(n)
	case []byte:
		return bc.AsBool(string(origin))

	case bool:
		return origin, nil

	case nil:
		return false, nil

	default:
		// try to cast to basic (in case input is ~basic)
		if basic, err := tryCastToBasicType(input); err == nil {
			return bc.AsBool(basic)
		}

		return tryCastUsingReflect[bool](input)
	}
}

func (bc boolCaster) AsBoolSlice(input any) ([]bool, error) {
	return ToSlice[bool](input, bc.AsBool)
}
