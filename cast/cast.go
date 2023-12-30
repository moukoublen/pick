package cast

import (
	"errors"
	"reflect"

	"github.com/moukoublen/pick/internal/testingx/errorsx"
)

type Caster struct {
	floatCaster
	stringCaster
	boolCaster
	timeCaster
	durationCaster
	integerCaster
	byteCaster
}

func NewCaster() Caster {
	return Caster{
		floatCaster:    newFloatCaster(),
		stringCaster:   newStringCaster(),
		boolCaster:     newBoolCaster(),
		timeCaster:     newTimeCaster(),
		durationCaster: newDurationCaster(),
		integerCaster:  newIntegerCaster(),
		byteCaster:     newByteCaster(),
	}
}

func (c Caster) As(input any, asKind reflect.Kind) (any, error) {
	//nolint:exhaustive
	switch asKind {
	case reflect.Float32:
		return c.AsFloat32(input)
	case reflect.Float64:
		return c.AsFloat64(input)
	case reflect.Int:
		return c.AsInt(input)
	case reflect.Int8:
		return c.AsInt8(input)
	case reflect.Int16:
		return c.AsInt16(input)
	case reflect.Int32:
		return c.AsInt32(input)
	case reflect.Int64:
		return c.AsInt64(input)
	case reflect.Uint:
		return c.AsUint(input)
	case reflect.Uint8:
		return c.AsUint8(input)
	case reflect.Uint16:
		return c.AsUint16(input)
	case reflect.Uint32:
		return c.AsUint32(input)
	case reflect.Uint64:
		return c.AsUint64(input)
	case reflect.Bool:
		return c.AsBool(input)
	case reflect.String:
		return c.AsString(input)
	}

	return nil, newCastError(ErrInvalidType, input)
}

func ToSlice[T any](input any, singleItemCastFn func(any) (T, error)) (_ []T, rErr error) {
	defer errorsx.RecoverPanicToError(&rErr)

	// quick returns just in case its already slice of T.
	if ss, is := input.([]T); is {
		return ss, nil
	}

	// attempt to quick return on slice of basic types by avoiding reflect.
	switch cc := input.(type) {
	case []any:
		return sliceToSlice(cc, singleItemCastFn)
	case []string:
		return sliceToSlice(cc, singleItemCastFn)
	case []int:
		return sliceToSlice(cc, singleItemCastFn)
	case []int8:
		return sliceToSlice(cc, singleItemCastFn)
	case []int16:
		return sliceToSlice(cc, singleItemCastFn)
	case []int32:
		return sliceToSlice(cc, singleItemCastFn)
	case []int64:
		return sliceToSlice(cc, singleItemCastFn)
	case []uint:
		return sliceToSlice(cc, singleItemCastFn)
	case []uint8:
		return sliceToSlice(cc, singleItemCastFn)
	case []uint16:
		return sliceToSlice(cc, singleItemCastFn)
	case []uint32:
		return sliceToSlice(cc, singleItemCastFn)
	case []uint64:
		return sliceToSlice(cc, singleItemCastFn)
	case []float32:
		return sliceToSlice(cc, singleItemCastFn)
	case []float64:
		return sliceToSlice(cc, singleItemCastFn)
	case []bool:
		return sliceToSlice(cc, singleItemCastFn)
	}

	typeOfInput := reflect.TypeOf(input)
	kindOfInput := typeOfInput.Kind()

	// if not slice or array then single cast attempt
	if kindOfInput != reflect.Array && kindOfInput != reflect.Slice {
		asT, err := singleItemCastFn(input)
		if err != nil {
			return nil, err
		}
		return []T{asT}, nil
	}

	// slow/costly attempt with reflect
	valueOfInput := reflect.ValueOf(input)
	castedSlice := make([]T, 0, valueOfInput.Len())
	for i := 0; i < valueOfInput.Len(); i++ {
		item := valueOfInput.Index(i)
		ifc := item.Interface()
		asStr, err := singleItemCastFn(ifc)
		if err != nil {
			return nil, err
		}
		castedSlice = append(castedSlice, asStr)
	}

	return castedSlice, nil
}

func sliceToSlice[In any, Out any](input []In, singleItemCastFn func(any) (Out, error)) (_ []Out, rErr error) {
	defer errorsx.RecoverPanicToError(&rErr)

	castedSlice := make([]Out, 0, len(input))
	for _, a := range input {
		casted, err := singleItemCastFn(a)
		if err != nil {
			return nil, err
		}
		castedSlice = append(castedSlice, casted)
	}
	return castedSlice, nil
}

// tryCastToBasicType checks input's Kind to identify if it can be casted as a basic type.
// If it can, it casts it and returns it.
// If not, it returns `ErrCannotBeCastedToBasic`.
func tryCastToBasicType(input any) (any, error) {
	if input == nil {
		return nil, newCastError(ErrCannotBeCastedToBasic, input)
	}

	t := reflect.TypeOf(input)
	k := t.Kind()

	if t.String() == k.String() {
		return input, newCastError(ErrAlreadyBasicType, input)
	}

	switch k {
	case reflect.Bool:
		return tryCastUsingReflect[bool](input)
	case reflect.Int:
		return tryCastUsingReflect[int](input)
	case reflect.Int8:
		return tryCastUsingReflect[int8](input)
	case reflect.Int16:
		return tryCastUsingReflect[int16](input)
	case reflect.Int32:
		return tryCastUsingReflect[int32](input)
	case reflect.Int64:
		return tryCastUsingReflect[int64](input)
	case reflect.Uint:
		return tryCastUsingReflect[uint](input)
	case reflect.Uint8:
		return tryCastUsingReflect[uint8](input)
	case reflect.Uint16:
		return tryCastUsingReflect[uint16](input)
	case reflect.Uint32:
		return tryCastUsingReflect[uint32](input)
	case reflect.Uint64:
		return tryCastUsingReflect[uint64](input)
	case reflect.Float32:
		return tryCastUsingReflect[float32](input)
	case reflect.Float64:
		return tryCastUsingReflect[float64](input)
	case reflect.String:
		return tryCastUsingReflect[string](input)
	}

	return nil, newCastError(ErrCannotBeCastedToBasic, input)
}

//nolint:ireturn
func tryCastUsingReflect[Out any](input any) (output Out, err error) {
	defer errorsx.RecoverPanicToError(&err)

	if input == nil {
		return output, newCastError(ErrInvalidType, input)
	}

	typeOfInput := reflect.TypeOf(input)
	valueOfInput := reflect.ValueOf(input)

	typeOfOutput := reflect.TypeOf(output)

	if !typeOfInput.ConvertibleTo(typeOfOutput) {
		return output, newCastError(ErrInvalidType, input)
	}

	convertedValue := valueOfInput.Convert(typeOfOutput)

	//nolint:forcetypeassert // if we get here we can safely assert.
	return convertedValue.Interface().(Out), nil
}

var (
	ErrCannotBeCastedToBasic = errors.New("value cannot be casted to basic type")
	ErrAlreadyBasicType      = errors.New("value is already basic type")
)
