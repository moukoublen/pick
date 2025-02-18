package pick

import (
	"encoding/json"

	"github.com/moukoublen/pick/iter"
)

func (c DefaultCaster) AsByte(input any) (byte, error) {
	switch origin := input.(type) {
	case byte:
		return origin, nil

	case int, int8, int16, int32, int64, uint, uint16, uint32, uint64, float32, float64, bool:
		return c.uint8Caster.cast(input)

	case string:
		return byte(0), newCastError(ErrCastInvalidType, input)
	case json.Number:
		n, err := origin.Float64()
		if err != nil {
			return 0, newCastError(err, input)
		}
		return c.AsByte(n)
	case []byte:
		return c.AsByte(string(origin))

	case nil:
		return 0, nil

	default:
		// try to cast to basic (in case input is ~basic)
		if basic, err := tryCastToBasicType(input); err == nil {
			return c.AsByte(basic)
		}

		return tryReflectConvert[byte](input)
	}
}

func (c DefaultCaster) AsByteSlice(input any) ([]byte, error) {
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
