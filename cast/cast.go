package cast

import (
	"errors"
	"reflect"

	"github.com/moukoublen/pick/cast/slices"
	"github.com/moukoublen/pick/internal/errorsx"
)

type Caster struct {
	directCastFunctionsTypes directCastFunctionsTypes
	intCaster                intCast[int]
	int8Caster               intCast[int8]
	int16Caster              intCast[int16]
	int32Caster              intCast[int32]
	int64Caster              intCast[int64]
	uintCaster               intCast[uint]
	uint8Caster              intCast[uint8]
	uint16Caster             intCast[uint16]
	uint32Caster             intCast[uint32]
	uint64Caster             intCast[uint64]
}

func NewCaster() Caster {
	return Caster{
		directCastFunctionsTypes: castFunctionTypes,
		intCaster:                newIntCast[int](),
		int8Caster:               newIntCast[int8](),
		int16Caster:              newIntCast[int16](),
		int32Caster:              newIntCast[int32](),
		int64Caster:              newIntCast[int64](),
		uintCaster:               newIntCast[uint](),
		uint8Caster:              newIntCast[uint8](),
		uint16Caster:             newIntCast[uint16](),
		uint32Caster:             newIntCast[uint32](),
		uint64Caster:             newIntCast[uint64](),
	}
}

//nolint:gocyclo
func (c Caster) ByType(input any, asType reflect.Type) (any, error) {
	// check if direct function is available
	switch asType {
	case c.directCastFunctionsTypes.typeOfBool:
		return c.AsBool(input)
	// case c.directCastFunctionsTypes.typeOfByte: // there is no distinguish type for byte. Its only uint8.
	// 	return c.AsByte(input)
	case c.directCastFunctionsTypes.typeOfInt8:
		return c.AsInt8(input)
	case c.directCastFunctionsTypes.typeOfInt16:
		return c.AsInt16(input)
	case c.directCastFunctionsTypes.typeOfInt32:
		return c.AsInt32(input)
	case c.directCastFunctionsTypes.typeOfInt64:
		return c.AsInt64(input)
	case c.directCastFunctionsTypes.typeOfInt:
		return c.AsInt(input)
	case c.directCastFunctionsTypes.typeOfUint8:
		return c.AsUint8(input)
	case c.directCastFunctionsTypes.typeOfUint16:
		return c.AsUint16(input)
	case c.directCastFunctionsTypes.typeOfUint32:
		return c.AsUint32(input)
	case c.directCastFunctionsTypes.typeOfUint64:
		return c.AsUint64(input)
	case c.directCastFunctionsTypes.typeOfUint:
		return c.AsUint(input)
	case c.directCastFunctionsTypes.typeOfFloat32:
		return c.AsFloat32(input)
	case c.directCastFunctionsTypes.typeOfFloat64:
		return c.AsFloat64(input)
	case c.directCastFunctionsTypes.typeOfString:
		return c.AsString(input)
	case c.directCastFunctionsTypes.typeOfTime:
		return c.AsTime(input)
	case c.directCastFunctionsTypes.typeOfDuration:
		return c.AsDuration(input)

	case c.directCastFunctionsTypes.typeOfSliceBool:
		return c.AsBoolSlice(input)
	// case c.directCastFunctionsTypes.typeOfSliceByte: // there is no distinguish type for byte. Its only uint8.
	// 	return c.AsByteSlice(input)
	case c.directCastFunctionsTypes.typeOfSliceInt8:
		return c.AsInt8Slice(input)
	case c.directCastFunctionsTypes.typeOfSliceInt16:
		return c.AsInt16Slice(input)
	case c.directCastFunctionsTypes.typeOfSliceInt32:
		return c.AsInt32Slice(input)
	case c.directCastFunctionsTypes.typeOfSliceInt64:
		return c.AsInt64Slice(input)
	case c.directCastFunctionsTypes.typeOfSliceInt:
		return c.AsIntSlice(input)
	case c.directCastFunctionsTypes.typeOfSliceUint8:
		return c.AsUint8Slice(input)
	case c.directCastFunctionsTypes.typeOfSliceUint16:
		return c.AsUint16Slice(input)
	case c.directCastFunctionsTypes.typeOfSliceUint32:
		return c.AsUint32Slice(input)
	case c.directCastFunctionsTypes.typeOfSliceUint64:
		return c.AsUint64Slice(input)
	case c.directCastFunctionsTypes.typeOfSliceUint:
		return c.AsUintSlice(input)
	case c.directCastFunctionsTypes.typeOfSliceFloat32:
		return c.AsFloat32Slice(input)
	case c.directCastFunctionsTypes.typeOfSliceFloat64:
		return c.AsFloat64Slice(input)
	case c.directCastFunctionsTypes.typeOfSliceString:
		return c.AsStringSlice(input)
	case c.directCastFunctionsTypes.typeOfSliceTime:
		return c.AsTimeSlice(input)
	case c.directCastFunctionsTypes.typeOfSliceDuration:
		return c.AsDurationSlice(input)
	}

	// check basic types aliases
	asKind := asType.Kind()
	_, isBasicKind := c.directCastFunctionsTypes.basicKindTypeMap[asKind]
	if isBasicKind {
		v, err := c.As(input, asKind)
		if err != nil {
			return nil, err
		}

		val := reflect.ValueOf(v)
		if !val.CanConvert(asType) {
			return nil, ErrInvalidType
		}
		return val.Convert(asType).Interface(), nil
	}

	// slice / array
	if asKind == reflect.Array || asKind == reflect.Slice {
		return c.sliceByType(input, asType, asType.Elem())
	}

	// TODO: reflect.Map
	// TODO: reflect.Pointer
	// TODO: reflect.Interface

	// fallback attempt to reflect convert
	val := reflect.ValueOf(input)
	if !val.CanConvert(asType) {
		return nil, ErrInvalidType
	}

	return val.Convert(asType).Interface(), nil
}

func (c Caster) sliceByType(input any, inputType, sliceElemType reflect.Type) (any, error) {
	inputValue := reflect.ValueOf(input)

	sc := 1
	switch inputType.Kind() {
	case reflect.Array, reflect.Slice:
		sc = inputValue.Len()
	}
	sliceValue := reflect.MakeSlice(reflect.SliceOf(sliceElemType), sc, sc)

	err := slices.ForEach(input, func(item any, meta slices.OpMeta) error {
		casted, cerr := c.ByType(item, sliceElemType)
		if cerr != nil {
			return cerr
		}

		castedValue := reflect.ValueOf(casted)
		sliceValue.Index(meta.Index).Set(castedValue)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return sliceValue.Interface(), nil
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
		return tryReflectConvert[bool](input)
	case reflect.Int:
		return tryReflectConvert[int](input)
	case reflect.Int8:
		return tryReflectConvert[int8](input)
	case reflect.Int16:
		return tryReflectConvert[int16](input)
	case reflect.Int32:
		return tryReflectConvert[int32](input)
	case reflect.Int64:
		return tryReflectConvert[int64](input)
	case reflect.Uint:
		return tryReflectConvert[uint](input)
	case reflect.Uint8:
		return tryReflectConvert[uint8](input)
	case reflect.Uint16:
		return tryReflectConvert[uint16](input)
	case reflect.Uint32:
		return tryReflectConvert[uint32](input)
	case reflect.Uint64:
		return tryReflectConvert[uint64](input)
	case reflect.Float32:
		return tryReflectConvert[float32](input)
	case reflect.Float64:
		return tryReflectConvert[float64](input)
	case reflect.String:
		return tryReflectConvert[string](input)
	}

	return nil, newCastError(ErrCannotBeCastedToBasic, input)
}

//nolint:ireturn
func tryReflectConvert[Out any](input any) (output Out, err error) {
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
