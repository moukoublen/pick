package pick

import (
	"encoding/json"
	"math"
	"reflect"
	"strconv"
	"strings"

	"github.com/moukoublen/pick/iter"
)

func (c DefaultConverter) AsInt(input any) (int, error) {
	return c.intConverter.convert(input)
}

func (c DefaultConverter) AsInt8(input any) (int8, error) {
	return c.int8Converter.convert(input)
}

func (c DefaultConverter) AsInt16(input any) (int16, error) {
	return c.int16Converter.convert(input)
}

func (c DefaultConverter) AsInt32(input any) (int32, error) {
	return c.int32Converter.convert(input)
}

func (c DefaultConverter) AsInt64(input any) (int64, error) {
	return c.int64Converter.convert(input)
}

func (c DefaultConverter) AsUint(input any) (uint, error) {
	return c.uintConverter.convert(input)
}

func (c DefaultConverter) AsUint8(input any) (uint8, error) {
	return c.uint8Converter.convert(input)
}

func (c DefaultConverter) AsUint16(input any) (uint16, error) {
	return c.uint16Converter.convert(input)
}

func (c DefaultConverter) AsUint32(input any) (uint32, error) {
	return c.uint32Converter.convert(input)
}

func (c DefaultConverter) AsUint64(input any) (uint64, error) {
	return c.uint64Converter.convert(input)
}

func (c DefaultConverter) AsIntSlice(input any) ([]int, error) {
	return iter.Map(input, iter.MapOpFn(c.AsInt))
}

func (c DefaultConverter) AsInt8Slice(input any) ([]int8, error) {
	return iter.Map(input, iter.MapOpFn(c.AsInt8))
}

func (c DefaultConverter) AsInt16Slice(input any) ([]int16, error) {
	return iter.Map(input, iter.MapOpFn(c.AsInt16))
}

func (c DefaultConverter) AsInt32Slice(input any) ([]int32, error) {
	return iter.Map(input, iter.MapOpFn(c.AsInt32))
}

func (c DefaultConverter) AsInt64Slice(input any) ([]int64, error) {
	return iter.Map(input, iter.MapOpFn(c.AsInt64))
}

func (c DefaultConverter) AsUintSlice(input any) ([]uint, error) {
	return iter.Map(input, iter.MapOpFn(c.AsUint))
}

func (c DefaultConverter) AsUint8Slice(input any) ([]uint8, error) {
	return iter.Map(input, iter.MapOpFn(c.AsUint8))
}

func (c DefaultConverter) AsUint16Slice(input any) ([]uint16, error) {
	return iter.Map(input, iter.MapOpFn(c.AsUint16))
}

func (c DefaultConverter) AsUint32Slice(input any) ([]uint32, error) {
	return iter.Map(input, iter.MapOpFn(c.AsUint32))
}

func (c DefaultConverter) AsUint64Slice(input any) ([]uint64, error) {
	return iter.Map(input, iter.MapOpFn(c.AsUint64))
}

type intConvert[T Integer] struct {
	signed bool
	kind   reflect.Kind
}

func newIntConvert[T Integer]() intConvert[T] {
	ic := intConvert[T]{}
	var t T
	ic.kind = reflect.TypeOf(t).Kind()
	ic.signed = ic.kind >= reflect.Int && ic.kind <= reflect.Int64

	return ic
}

//nolint:ireturn
func (ic intConvert[T]) fromInt64(origin int64) (T, error) {
	t := T(origin)
	if !int64ConvertValid(origin, ic.kind) {
		return t, newConvertError(ErrConvertOverFlow, origin)
	}
	return t, nil
}

//nolint:ireturn
func (ic intConvert[T]) fromUint64(origin uint64) (T, error) {
	t := T(origin)
	if !uint64ConvertValid(origin, ic.kind) {
		return t, newConvertError(ErrConvertOverFlow, origin)
	}
	return t, nil
}

//nolint:ireturn
func (ic intConvert[T]) fromFloat64(origin float64) (T, error) {
	converted, err := float64ToInt64(origin)
	if err != nil {
		return T(converted), err
	}

	return ic.fromInt64(converted)
}

//nolint:ireturn
func (ic intConvert[T]) fromString(origin string) (T, error) {
	if strings.ContainsAny(origin, ".e") {
		v, err := strconv.ParseFloat(origin, 64)
		if err != nil {
			return T(v), newConvertError(err, origin)
		}

		return ic.fromFloat64(v)
	}

	if ic.signed {
		v, err := strconv.ParseInt(origin, 10, 64)
		if err != nil {
			return T(v), newConvertError(err, origin)
		}

		return ic.fromInt64(v)
	}

	v, err := strconv.ParseUint(origin, 10, 64)
	if err != nil {
		return T(v), newConvertError(err, origin)
	}

	return ic.fromUint64(v)
}

//nolint:ireturn
func (ic intConvert[T]) fromBool(origin bool) T {
	if origin {
		return 1
	}

	return 0
}

//nolint:ireturn
func (ic intConvert[T]) convert(input any) (T, error) {
	switch origin := input.(type) {
	case int:
		return ic.fromInt64(int64(origin))
	case int8:
		return ic.fromInt64(int64(origin))
	case int16:
		return ic.fromInt64(int64(origin))
	case int32:
		return ic.fromInt64(int64(origin))
	case int64:
		return ic.fromInt64(origin)

	case uint:
		return ic.fromUint64(uint64(origin))
	case uint8:
		return ic.fromUint64(uint64(origin))
	case uint16:
		return ic.fromUint64(uint64(origin))
	case uint32:
		return ic.fromUint64(uint64(origin))
	case uint64:
		return ic.fromUint64(origin)

	case float32:
		return ic.fromFloat64(float64(origin))
	case float64:
		return ic.fromFloat64(origin)

	case string:
		return ic.fromString(origin)
	case json.Number:
		return ic.fromString(string(origin))
	case []byte:
		return ic.fromString(string(origin))

	case bool:
		return ic.fromBool(origin), nil

	case nil:
		return 0, nil

	default:
		// try to convert to basic (in case input is ~basic)
		if basic, err := tryConvertToBasicType(input); err == nil {
			return ic.convert(basic)
		}

		return tryReflectConvert[T](input)
	}
}

type sizeOfInteger int

// The int, uint types are usually 32 bits wide on 32-bit systems and 64 bits wide on 64-bit systems.
// https://go.dev/ref/spec#Numeric_types

const (
	//nolint:mnd,gomnd // there are not magic numbers to be fair
	sizeOfInt    sizeOfInteger = 32 << (^uint(0) >> 63) // 32 or 64 // https://github.com/golang/go/blob/3d33437c450aa74014ea1d41cd986b6ee6266984/src/math/const.go#L40
	sizeOfInt8   sizeOfInteger = 8
	sizeOfInt16  sizeOfInteger = 16
	sizeOfInt32  sizeOfInteger = 32
	sizeOfInt64  sizeOfInteger = 64
	sizeOfUInt8  sizeOfInteger = 8
	sizeOfUInt16 sizeOfInteger = 16
	sizeOfUInt32 sizeOfInteger = 32
	sizeOfUInt64 sizeOfInteger = 64
)

//
// Converts range checks
//

func int64ConvertValid(origin int64, to reflect.Kind) bool {
	switch to {
	case reflect.Int:
		return origin >= math.MinInt && origin <= math.MaxInt
	case reflect.Int8:
		return origin >= math.MinInt8 && origin <= math.MaxInt8
	case reflect.Int16:
		return origin >= math.MinInt16 && origin <= math.MaxInt16
	case reflect.Int32:
		return origin >= math.MinInt32 && origin <= math.MaxInt32
	case reflect.Int64:
		return true

	case reflect.Uint:
		if sizeOfInt == sizeOfInt32 {
			return origin >= 0 && origin <= math.MaxUint32
		}
		return origin >= 0
	case reflect.Uint8:
		return origin >= 0 && origin <= math.MaxUint8
	case reflect.Uint16:
		return origin >= 0 && origin <= math.MaxUint16
	case reflect.Uint32:
		return origin >= 0 && origin <= math.MaxUint32
	case reflect.Uint64:
		return origin >= 0

	default:
		return false
	}
}

func uint64ConvertValid(origin uint64, to reflect.Kind) bool {
	switch to {
	case reflect.Int:
		return origin <= math.MaxInt
	case reflect.Int8:
		return origin <= math.MaxInt8
	case reflect.Int16:
		return origin <= math.MaxInt16
	case reflect.Int32:
		return origin <= math.MaxInt32
	case reflect.Int64:
		return origin <= math.MaxInt64

	case reflect.Uint:
		if sizeOfInt == sizeOfInt32 {
			return origin <= math.MaxUint32
		}
		return true
	case reflect.Uint8:
		return origin <= math.MaxUint8
	case reflect.Uint16:
		return origin <= math.MaxUint16
	case reflect.Uint32:
		return origin <= math.MaxUint32
	case reflect.Uint64:
		return true

	default:
		return false
	}
}

func floatIsWhole(num float64) bool {
	return num == math.Trunc(num)
}

func float64ToInt64(origin float64) (int64, error) {
	converted := int64(origin)

	if origin > math.MaxInt64 {
		return converted, newConvertError(ErrConvertOverFlow, origin)
	}
	if origin < math.MinInt64 {
		return converted, newConvertError(ErrConvertOverFlow, origin)
	}

	if !floatIsWhole(origin) {
		return converted, newConvertError(ErrConvertLostDecimals, origin)
	}

	return converted, nil
}
