package cast

import "encoding/json"

type byteCaster struct {
	uint8Caster intCast[uint8]
}

func newByteCaster() byteCaster {
	return byteCaster{
		uint8Caster: newIntCast[uint8](),
	}
}

func (bc byteCaster) AsByte(input any) (byte, error) {
	switch cc := input.(type) {
	case byte:
		return cc, nil
	case string:
		return byte(0), newCastError(ErrInvalidType, input)
	case json.RawMessage:
		return byte(0), newCastError(ErrInvalidType, input)
	}

	return bc.uint8Caster.cast(input)
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

	return castToSlice[byte](input, bc.AsByte)
}
