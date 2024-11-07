package pick

import (
	"encoding/json"
	"math"
	"reflect"
	"strconv"
	"strings"

	"github.com/moukoublen/pick/numbers"
	"github.com/moukoublen/pick/slices"
)

func (c DefaultCaster) AsInt(input any) (int, error) {
	return c.intCaster.cast(input)
}

func (c DefaultCaster) AsInt8(input any) (int8, error) {
	return c.int8Caster.cast(input)
}

func (c DefaultCaster) AsInt16(input any) (int16, error) {
	return c.int16Caster.cast(input)
}

func (c DefaultCaster) AsInt32(input any) (int32, error) {
	return c.int32Caster.cast(input)
}

func (c DefaultCaster) AsInt64(input any) (int64, error) {
	return c.int64Caster.cast(input)
}

func (c DefaultCaster) AsUint(input any) (uint, error) {
	return c.uintCaster.cast(input)
}

func (c DefaultCaster) AsUint8(input any) (uint8, error) {
	return c.uint8Caster.cast(input)
}

func (c DefaultCaster) AsUint16(input any) (uint16, error) {
	return c.uint16Caster.cast(input)
}

func (c DefaultCaster) AsUint32(input any) (uint32, error) {
	return c.uint32Caster.cast(input)
}

func (c DefaultCaster) AsUint64(input any) (uint64, error) {
	return c.uint64Caster.cast(input)
}

func (c DefaultCaster) AsIntSlice(input any) ([]int, error) {
	return slices.Map(input, slices.MapOpFn(c.AsInt))
}

func (c DefaultCaster) AsInt8Slice(input any) ([]int8, error) {
	return slices.Map(input, slices.MapOpFn(c.AsInt8))
}

func (c DefaultCaster) AsInt16Slice(input any) ([]int16, error) {
	return slices.Map(input, slices.MapOpFn(c.AsInt16))
}

func (c DefaultCaster) AsInt32Slice(input any) ([]int32, error) {
	return slices.Map(input, slices.MapOpFn(c.AsInt32))
}

func (c DefaultCaster) AsInt64Slice(input any) ([]int64, error) {
	return slices.Map(input, slices.MapOpFn(c.AsInt64))
}

func (c DefaultCaster) AsUintSlice(input any) ([]uint, error) {
	return slices.Map(input, slices.MapOpFn(c.AsUint))
}

func (c DefaultCaster) AsUint8Slice(input any) ([]uint8, error) {
	return slices.Map(input, slices.MapOpFn(c.AsUint8))
}

func (c DefaultCaster) AsUint16Slice(input any) ([]uint16, error) {
	return slices.Map(input, slices.MapOpFn(c.AsUint16))
}

func (c DefaultCaster) AsUint32Slice(input any) ([]uint32, error) {
	return slices.Map(input, slices.MapOpFn(c.AsUint32))
}

func (c DefaultCaster) AsUint64Slice(input any) ([]uint64, error) {
	return slices.Map(input, slices.MapOpFn(c.AsUint64))
}

type intCast[T numbers.Integer] struct {
	signed bool
	kind   reflect.Kind
}

func newIntCast[T numbers.Integer]() intCast[T] {
	ic := intCast[T]{}
	var t T
	ic.kind = reflect.TypeOf(t).Kind()
	ic.signed = ic.kind >= reflect.Int && ic.kind <= reflect.Int64

	return ic
}

//nolint:ireturn
func (ic intCast[T]) fromInt64(origin int64) (T, error) {
	t := T(origin)
	if !int64CastValid(origin, ic.kind) {
		return t, newCastError(ErrCastOverFlow, origin)
	}
	return t, nil
}

//nolint:ireturn
func (ic intCast[T]) fromUint64(origin uint64) (T, error) {
	t := T(origin)
	if !uint64CastValid(origin, ic.kind) {
		return t, newCastError(ErrCastOverFlow, origin)
	}
	return t, nil
}

//nolint:ireturn
func (ic intCast[T]) fromFloat64(origin float64) (T, error) {
	casted, err := float64ToInt64(origin)
	if err != nil {
		return T(casted), err
	}

	return ic.fromInt64(casted)
}

//nolint:ireturn
func (ic intCast[T]) fromString(origin string) (T, error) {
	if strings.ContainsAny(origin, ".e") {
		v, err := strconv.ParseFloat(origin, 64)
		if err != nil {
			return T(v), newCastError(err, origin)
		}

		return ic.fromFloat64(v)
	}

	if ic.signed {
		v, err := strconv.ParseInt(origin, 10, 64)
		if err != nil {
			return T(v), newCastError(err, origin)
		}

		return ic.fromInt64(v)
	}

	v, err := strconv.ParseUint(origin, 10, 64)
	if err != nil {
		return T(v), newCastError(err, origin)
	}

	return ic.fromUint64(v)
}

//nolint:ireturn
func (ic intCast[T]) fromBool(origin bool) T {
	if origin {
		return 1
	}

	return 0
}

//nolint:ireturn
func (ic intCast[T]) cast(input any) (T, error) {
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
		// try to cast to basic (in case input is ~basic)
		if basic, err := tryCastToBasicType(input); err == nil {
			return ic.cast(basic)
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
// Casts range checks
//

func int64CastValid(origin int64, to reflect.Kind) bool {
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

func uint64CastValid(origin uint64, to reflect.Kind) bool {
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
	casted := int64(origin)

	if origin > math.MaxInt64 {
		return casted, newCastError(ErrCastOverFlow, origin)
	}
	if origin < math.MinInt64 {
		return casted, newCastError(ErrCastOverFlow, origin)
	}

	if !floatIsWhole(origin) {
		return casted, newCastError(ErrCastLostDecimals, origin)
	}

	return casted, nil
}
