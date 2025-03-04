package pick

import (
	"encoding/json"

	"github.com/moukoublen/pick/iter"
)

func (c DefaultConverter) AsByte(input any) (byte, error) {
	switch origin := input.(type) {
	case byte:
		return origin, nil

	case int, int8, int16, int32, int64, uint, uint16, uint32, uint64, float32, float64, bool:
		return c.uint8Converter.convert(input)

	case string:
		return byte(0), newConvertError(ErrConvertInvalidType, input)
	case json.Number:
		n, err := origin.Float64()
		if err != nil {
			return 0, newConvertError(err, input)
		}
		return c.AsByte(n)
	case []byte:
		return c.AsByte(string(origin))

	case nil:
		return 0, nil

	default:
		// try to convert to basic (in case input is ~basic)
		if basic, err := tryConvertToBasicType(input); err == nil {
			return c.AsByte(basic)
		}

		return tryReflectConvert[byte](input)
	}
}

func (c DefaultConverter) AsByteSlice(input any) ([]byte, error) {
	switch cc := input.(type) {
	case []byte:
		return cc, nil
	case string:
		return []byte(cc), nil
	case json.RawMessage:
		return []byte(cc), nil
	}

	return iter.Map(input, iter.MapOpFn(c.AsByte))
}
