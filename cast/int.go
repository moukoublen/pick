package cast

import (
	"encoding/json"
	"math"
	"reflect"
	"strconv"
	"strings"

	"golang.org/x/exp/constraints"
)

type integerCaster struct {
	intCaster    intCast[int]
	int8Caster   intCast[int8]
	int16Caster  intCast[int16]
	int32Caster  intCast[int32]
	int64Caster  intCast[int64]
	uintCaster   intCast[uint]
	uint8Caster  intCast[uint8]
	uint16Caster intCast[uint16]
	uint32Caster intCast[uint32]
	uint64Caster intCast[uint64]
}

func newIntegerCaster() integerCaster {
	return integerCaster{
		intCaster:    newIntCast[int](),
		int8Caster:   newIntCast[int8](),
		int16Caster:  newIntCast[int16](),
		int32Caster:  newIntCast[int32](),
		int64Caster:  newIntCast[int64](),
		uintCaster:   newIntCast[uint](),
		uint8Caster:  newIntCast[uint8](),
		uint16Caster: newIntCast[uint16](),
		uint32Caster: newIntCast[uint32](),
		uint64Caster: newIntCast[uint64](),
	}
}

func (c integerCaster) AsInt(input any) (int, error) {
	return c.intCaster.cast(input)
}

func (c integerCaster) AsInt8(input any) (int8, error) {
	return c.int8Caster.cast(input)
}

func (c integerCaster) AsInt16(input any) (int16, error) {
	return c.int16Caster.cast(input)
}

func (c integerCaster) AsInt32(input any) (int32, error) {
	return c.int32Caster.cast(input)
}

func (c integerCaster) AsInt64(input any) (int64, error) {
	return c.int64Caster.cast(input)
}

func (c integerCaster) AsUint(input any) (uint, error) {
	return c.uintCaster.cast(input)
}

func (c integerCaster) AsUint8(input any) (uint8, error) {
	return c.uint8Caster.cast(input)
}

func (c integerCaster) AsUint16(input any) (uint16, error) {
	return c.uint16Caster.cast(input)
}

func (c integerCaster) AsUint32(input any) (uint32, error) {
	return c.uint32Caster.cast(input)
}

func (c integerCaster) AsUint64(input any) (uint64, error) {
	return c.uint64Caster.cast(input)
}

func (c integerCaster) AsIntSlice(input any) ([]int, error) {
	return castToSlice[int](input, c.AsInt)
}

func (c integerCaster) AsInt8Slice(input any) ([]int8, error) {
	return castToSlice[int8](input, c.AsInt8)
}

func (c integerCaster) AsInt16Slice(input any) ([]int16, error) {
	return castToSlice[int16](input, c.AsInt16)
}

func (c integerCaster) AsInt32Slice(input any) ([]int32, error) {
	return castToSlice[int32](input, c.AsInt32)
}

func (c integerCaster) AsInt64Slice(input any) ([]int64, error) {
	return castToSlice[int64](input, c.AsInt64)
}

func (c integerCaster) AsUintSlice(input any) ([]uint, error) {
	return castToSlice[uint](input, c.AsUint)
}

func (c integerCaster) AsUint8Slice(input any) ([]uint8, error) {
	return castToSlice[uint8](input, c.AsUint8)
}

func (c integerCaster) AsUint16Slice(input any) ([]uint16, error) {
	return castToSlice[uint16](input, c.AsUint16)
}

func (c integerCaster) AsUint32Slice(input any) ([]uint32, error) {
	return castToSlice[uint32](input, c.AsUint32)
}

func (c integerCaster) AsUint64Slice(input any) ([]uint64, error) {
	return castToSlice[uint64](input, c.AsUint64)
}

type intCast[T constraints.Integer] struct {
	signed  bool
	kind    reflect.Kind
	intSize sizeOfInteger
}

func newIntCast[T constraints.Integer]() intCast[T] {
	ic := intCast[T]{}
	var t T
	ic.kind = reflect.TypeOf(t).Kind()
	ic.signed = ic.kind >= reflect.Int && ic.kind <= reflect.Int64
	ic.intSize = getIntSize(ic.kind)

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
	casted := int64(origin)
	if origin-float64(casted) > 0 {
		return T(casted), newCastError(ErrCastLostDecimals, origin)
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

	case bool:
		return ic.fromBool(origin), nil

	case nil:
		return 0, nil

	default:
		return castAttemptUsingReflect[T](input)
	}
}

type sizeOfInteger int

const (
	//nolint:gomnd // there are not magic numbers to be fair
	sizeOfInt    sizeOfInteger = 32 << (^uint(0) >> 63) // 32 or 64
	sizeOfInt8   sizeOfInteger = 8
	sizeOfInt16  sizeOfInteger = 16
	sizeOfInt32  sizeOfInteger = 32
	sizeOfInt64  sizeOfInteger = 64
	sizeOfUInt8  sizeOfInteger = 8
	sizeOfUInt16 sizeOfInteger = 16
	sizeOfUInt32 sizeOfInteger = 32
	sizeOfUInt64 sizeOfInteger = 64
)

func getIntSize(k reflect.Kind) sizeOfInteger {
	switch k {
	case reflect.Int:
		return sizeOfInt
	case reflect.Int8:
		return sizeOfInt8
	case reflect.Int16:
		return sizeOfInt16
	case reflect.Int32:
		return sizeOfInt32
	case reflect.Int64:
		return sizeOfInt64
	case reflect.Uint:
		return sizeOfInt
	case reflect.Uint8:
		return sizeOfUInt8
	case reflect.Uint16:
		return sizeOfUInt16
	case reflect.Uint32:
		return sizeOfUInt32
	case reflect.Uint64:
		return sizeOfUInt64
	default:
		return sizeOfInt
	}
}

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
