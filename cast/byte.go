package cast

import (
	"encoding/json"
)

type byteCaster struct {
	uint8Caster intCast[uint8]
}

func newByteCaster() byteCaster {
	return byteCaster{
		uint8Caster: newIntCast[uint8](),
	}
}

func (bc byteCaster) AsByte(input any) (byte, error) {
	switch origin := input.(type) {
	case byte:
		return origin, nil

	case int, int8, int16, int32, int64, uint, uint16, uint32, uint64, float32, float64, bool:
		return bc.uint8Caster.cast(input)

	case string:
		return byte(0), newCastError(ErrInvalidType, input)
	case json.Number:
		n, err := origin.Float64()
		if err != nil {
			return 0, newCastError(err, input)
		}
		return bc.AsByte(n)
	case []byte:
		return bc.AsByte(string(origin))

	case nil:
		return 0, nil

	default:
		// try to cast to basic (in case input is ~basic)
		if basic, err := tryCastToBasicType(input); err == nil {
			return bc.AsByte(basic)
		}

		return tryCastUsingReflect[byte](input)
	}
}

func (bc byteCaster) AsByteSlice(input any) ([]byte, error) {
	switch cc := input.(type) {
	case []byte:
		return cc, nil
	case string:
		return []byte(cc), nil
	case json.RawMessage:
		return []byte(cc), nil
	}

	return ToSlice[byte](input, bc.AsByte)
}
