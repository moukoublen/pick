package pick

import (
	"errors"
	"reflect"

	"github.com/moukoublen/pick/internal/errorsx"
	"github.com/moukoublen/pick/iter"
)

type DefaultCaster struct {
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

func NewDefaultCaster() DefaultCaster {
	return DefaultCaster{
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

// ByType attempts to cast the `input` to the type defined by `asType`. It returns error if the cast fails.
// It first attempts to cast using a quick flow (performance wise) when the target type is a basic type, without using reflect.
// Then it tries to handle basic type aliases.
// And then it falls back to reflect usage depending on the target type.
// If no error is returned then it is safe to use type assertion in the returned value, to the type given in `asType`.
// e.g.
//
//	i, err := c.ByType("123", reflect.TypeOf(int64(0)))
//	i.(int64) // safe
func (c DefaultCaster) ByType(input any, asType reflect.Type) (any, error) {
	// if target type is a basic type.
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

	asKind := asType.Kind()

	// if target type is a basic type alias (e.g. type myString string).
	if _, isBasicKind := c.directCastFunctionsTypes.basicKindTypeMap[asKind]; isBasicKind {
		v, err := c.As(input, asKind)
		if err != nil {
			return nil, err
		}

		val := reflect.ValueOf(v)
		if !val.CanConvert(asType) {
			return nil, ErrCastInvalidType
		}
		return val.Convert(asType).Interface(), nil
	}

	switch asKind {
	// if target type is slice / array
	case reflect.Array, reflect.Slice:
		return c.toSliceByType(input, asType.Elem())

	// if target type is map
	case reflect.Map:
		return c.toMapByType(input, asType.Key(), asType.Elem())

	// if target type is pointer
	case reflect.Pointer:
		return c.toPointerByType(input, asType.Elem())
	}

	// fallback attempt to reflect convert
	val := reflect.ValueOf(input)
	if !val.CanConvert(asType) {
		return nil, ErrCastInvalidType
	}

	return val.Convert(asType).Interface(), nil
}

func (c DefaultCaster) toSliceByType(input any, asSliceElemType reflect.Type) (any, error) {
	inputValue := reflect.ValueOf(input)

	sc := 1
	switch inputValue.Kind() {
	case reflect.Array, reflect.Slice:
		sc = inputValue.Len()
	}
	sliceValue := reflect.MakeSlice(reflect.SliceOf(asSliceElemType), sc, sc)

	err := iter.ForEach(input, func(item any, meta iter.CollectionOpMeta) error {
		casted, cerr := c.ByType(item, asSliceElemType)
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

// toMapByType can cast to map if the input is either slice, array or map of any type.
func (c DefaultCaster) toMapByType(input any, keyType, valueType reflect.Type) (any, error) { //nolint:revive
	// TODO: implement.
	return nil, ErrCastInvalidType
}

func (c DefaultCaster) toPointerByType(input any, pointerTargetType reflect.Type) (any, error) { //nolint:revive
	// pointerValue := reflect.New(pointerTargetType)

	// TODO: implement.
	return nil, ErrCastInvalidType
}

func (c DefaultCaster) As(input any, asKind reflect.Kind) (any, error) {
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

	return nil, newCastError(ErrCastInvalidType, input)
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
		return output, newCastError(ErrCastInvalidType, input)
	}

	typeOfInput := reflect.TypeOf(input)
	valueOfInput := reflect.ValueOf(input)

	typeOfOutput := reflect.TypeOf(output)

	if !typeOfInput.ConvertibleTo(typeOfOutput) {
		return output, newCastError(ErrCastInvalidType, input)
	}

	convertedValue := valueOfInput.Convert(typeOfOutput)

	//nolint:forcetypeassert // if we get here we can safely assert.
	return convertedValue.Interface().(Out), nil
}

var (
	ErrCannotBeCastedToBasic = errors.New("value cannot be casted to basic type")
	ErrAlreadyBasicType      = errors.New("value is already basic type")
)

//nolint:gochecknoglobals
var defaultCasterGlobal = NewDefaultCaster()

// Cast attempts to convert the input value to the specified (generic) Output type.
// It returns the converted value of type Output and an error if the casting fails.
// e.g.
//
//	Cast[string](123) // "123", nil
func Cast[Output any](input any) (Output, error) { //nolint:ireturn
	var o Output

	n, err := defaultCasterGlobal.ByType(input, reflect.TypeOf(o))
	if err != nil {
		return o, err
	}

	casted, is := n.(Output)
	if !is {
		return o, newCastError(ErrCastInvalidType, input)
	}

	return casted, nil
}
