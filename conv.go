package pick

import (
	"errors"
	"reflect"

	"github.com/moukoublen/pick/internal/errorsx"
	"github.com/moukoublen/pick/iter"
)

type DefaultConverter struct {
	directConvertFunctionsTypes directConvertFunctionsTypes
	intConverter                intConvert[int]
	int8Converter               intConvert[int8]
	int16Converter              intConvert[int16]
	int32Converter              intConvert[int32]
	int64Converter              intConvert[int64]
	uintConverter               intConvert[uint]
	uint8Converter              intConvert[uint8]
	uint16Converter             intConvert[uint16]
	uint32Converter             intConvert[uint32]
	uint64Converter             intConvert[uint64]
}

func NewDefaultConverter() DefaultConverter {
	return DefaultConverter{
		directConvertFunctionsTypes: convertFunctionTypes,
		intConverter:                newIntConvert[int](),
		int8Converter:               newIntConvert[int8](),
		int16Converter:              newIntConvert[int16](),
		int32Converter:              newIntConvert[int32](),
		int64Converter:              newIntConvert[int64](),
		uintConverter:               newIntConvert[uint](),
		uint8Converter:              newIntConvert[uint8](),
		uint16Converter:             newIntConvert[uint16](),
		uint32Converter:             newIntConvert[uint32](),
		uint64Converter:             newIntConvert[uint64](),
	}
}

// ByType attempts to convert the `input` to the type defined by `asType`. It returns error if the convert fails.
// It first attempts to convert using a quick flow (performance wise) when the target type is a basic type, without using reflect.
// Then it tries to handle basic type aliases.
// And then it falls back to reflect usage depending on the target type.
// If no error is returned then it is safe to use type assertion in the returned value, to the type given in `asType`.
// e.g.
//
//	i, err := c.ByType("123", reflect.TypeOf(int64(0)))
//	i.(int64) // safe
func (c DefaultConverter) ByType(input any, asType reflect.Type) (any, error) {
	// if target type is a basic type.
	switch asType {
	case c.directConvertFunctionsTypes.typeOfBool:
		return c.AsBool(input)
	// case c.directConvertFunctionsTypes.typeOfByte: // there is no distinguish type for byte. Its only uint8.
	// 	return c.AsByte(input)
	case c.directConvertFunctionsTypes.typeOfInt8:
		return c.AsInt8(input)
	case c.directConvertFunctionsTypes.typeOfInt16:
		return c.AsInt16(input)
	case c.directConvertFunctionsTypes.typeOfInt32:
		return c.AsInt32(input)
	case c.directConvertFunctionsTypes.typeOfInt64:
		return c.AsInt64(input)
	case c.directConvertFunctionsTypes.typeOfInt:
		return c.AsInt(input)
	case c.directConvertFunctionsTypes.typeOfUint8:
		return c.AsUint8(input)
	case c.directConvertFunctionsTypes.typeOfUint16:
		return c.AsUint16(input)
	case c.directConvertFunctionsTypes.typeOfUint32:
		return c.AsUint32(input)
	case c.directConvertFunctionsTypes.typeOfUint64:
		return c.AsUint64(input)
	case c.directConvertFunctionsTypes.typeOfUint:
		return c.AsUint(input)
	case c.directConvertFunctionsTypes.typeOfFloat32:
		return c.AsFloat32(input)
	case c.directConvertFunctionsTypes.typeOfFloat64:
		return c.AsFloat64(input)
	case c.directConvertFunctionsTypes.typeOfString:
		return c.AsString(input)
	case c.directConvertFunctionsTypes.typeOfTime:
		return c.AsTime(input)
	case c.directConvertFunctionsTypes.typeOfDuration:
		return c.AsDuration(input)

	case c.directConvertFunctionsTypes.typeOfSliceBool:
		return c.AsBoolSlice(input)
	// case c.directConvertFunctionsTypes.typeOfSliceByte: // there is no distinguish type for byte. Its only uint8.
	// 	return c.AsByteSlice(input)
	case c.directConvertFunctionsTypes.typeOfSliceInt8:
		return c.AsInt8Slice(input)
	case c.directConvertFunctionsTypes.typeOfSliceInt16:
		return c.AsInt16Slice(input)
	case c.directConvertFunctionsTypes.typeOfSliceInt32:
		return c.AsInt32Slice(input)
	case c.directConvertFunctionsTypes.typeOfSliceInt64:
		return c.AsInt64Slice(input)
	case c.directConvertFunctionsTypes.typeOfSliceInt:
		return c.AsIntSlice(input)
	case c.directConvertFunctionsTypes.typeOfSliceUint8:
		return c.AsUint8Slice(input)
	case c.directConvertFunctionsTypes.typeOfSliceUint16:
		return c.AsUint16Slice(input)
	case c.directConvertFunctionsTypes.typeOfSliceUint32:
		return c.AsUint32Slice(input)
	case c.directConvertFunctionsTypes.typeOfSliceUint64:
		return c.AsUint64Slice(input)
	case c.directConvertFunctionsTypes.typeOfSliceUint:
		return c.AsUintSlice(input)
	case c.directConvertFunctionsTypes.typeOfSliceFloat32:
		return c.AsFloat32Slice(input)
	case c.directConvertFunctionsTypes.typeOfSliceFloat64:
		return c.AsFloat64Slice(input)
	case c.directConvertFunctionsTypes.typeOfSliceString:
		return c.AsStringSlice(input)
	case c.directConvertFunctionsTypes.typeOfSliceTime:
		return c.AsTimeSlice(input)
	case c.directConvertFunctionsTypes.typeOfSliceDuration:
		return c.AsDurationSlice(input)
	}

	asKind := asType.Kind()

	// if target type is a basic type alias (e.g. type myString string).
	if _, isBasicKind := c.directConvertFunctionsTypes.basicKindTypeMap[asKind]; isBasicKind {
		v, err := c.As(input, asKind)
		if err != nil {
			return nil, err
		}

		val := reflect.ValueOf(v)
		if !val.CanConvert(asType) {
			return nil, ErrConvertInvalidType
		}
		return val.Convert(asType).Interface(), nil
	}

	switch asKind {
	// if target type is slice / array
	case reflect.Array, reflect.Slice:
		return c.toSliceByType(input, asType.Elem())

	// if target type is map
	case reflect.Map:
		return c.toMapByType(input, asType, asType.Key(), asType.Elem())

	// if target type is pointer
	case reflect.Pointer:
		return c.toPointerByType(input, asType.Elem())
	}

	// fallback attempt to reflect convert
	val := reflect.ValueOf(input)
	if !val.CanConvert(asType) {
		return nil, ErrConvertInvalidType
	}

	return val.Convert(asType).Interface(), nil
}

func (c DefaultConverter) toSliceByType(input any, asSliceElemType reflect.Type) (any, error) {
	inputValue := reflect.ValueOf(input)

	sc := 1
	switch inputValue.Kind() {
	case reflect.Array, reflect.Slice:
		sc = inputValue.Len()
	}
	sliceValue := reflect.MakeSlice(reflect.SliceOf(asSliceElemType), sc, sc)

	err := iter.ForEach(input, func(item any, meta iter.CollectionOpMeta) error {
		converted, cerr := c.ByType(item, asSliceElemType)
		if cerr != nil {
			return cerr
		}

		convertedValue := reflect.ValueOf(converted)
		sliceValue.Index(meta.Index).Set(convertedValue)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return sliceValue.Interface(), nil
}

// toMapByType can convert to map if the input is either slice, array or map of any type.
func (c DefaultConverter) toMapByType(input any, dstType, dstKeyType, dstValueType reflect.Type) (r any, rErr error) { //nolint:revive
	defer errorsx.RecoverPanicToError(&rErr)

	srcValue := reflect.ValueOf(input)
	srcKind := srcValue.Kind()

	if srcKind != reflect.Map && srcKind != reflect.Array && srcKind != reflect.Slice {
		return nil, ErrConvertInvalidType
	}

	if srcKind == reflect.Map { // quick return if already the type.
		srcType := srcValue.Type()
		if srcType.Key() == dstKeyType && srcType.Elem() == dstValueType {
			return input, nil
		}
	}

	switch srcKind {
	// if source kind is map
	case reflect.Map:
		srcType := srcValue.Type()
		if srcType.Key() == dstKeyType && srcType.Elem() == dstValueType {
			return input, nil
		}
		dstValue := reflect.MakeMapWithSize(dstType, srcValue.Len())
		iter := srcValue.MapRange()
		for iter.Next() {
			keyConverted, kErr := c.ByType(iter.Key().Interface(), dstKeyType)
			if kErr != nil {
				return nil, kErr
			}

			valueConverted, vErr := c.ByType(iter.Value().Interface(), dstValueType)
			if vErr != nil {
				return nil, vErr
			}

			dstValue.SetMapIndex(reflect.ValueOf(keyConverted), reflect.ValueOf(valueConverted))
		}

		return dstValue.Interface(), nil

	// if source kind is array or slice
	case reflect.Array, reflect.Slice:
		dstValue := reflect.MakeMapWithSize(dstType, srcValue.Len())
		length := srcValue.Len()
		for i := range length {
			item := srcValue.Index(i)

			keyConverted, kErr := c.ByType(i, dstKeyType)
			if kErr != nil {
				return nil, kErr
			}

			valueConverted, vErr := c.ByType(item.Interface(), dstValueType)
			if vErr != nil {
				return nil, vErr
			}

			dstValue.SetMapIndex(reflect.ValueOf(keyConverted), reflect.ValueOf(valueConverted))
		}

		return dstValue.Interface(), nil
	}

	return nil, ErrConvertInvalidType
}

func (c DefaultConverter) toPointerByType(input any, pointerTargetType reflect.Type) (any, error) { //nolint:revive
	// pointerValue := reflect.New(pointerTargetType)

	// TODO: implement.
	return nil, ErrConvertInvalidType
}

func (c DefaultConverter) As(input any, asKind reflect.Kind) (any, error) {
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

	return nil, newConvertError(ErrConvertInvalidType, input)
}

// tryConvertToBasicType checks input's Kind to identify if it can be converted as a basic type.
// If it can, it converts it and returns it.
// If not, it returns `ErrCannotBeConvertedToBasic`.
func tryConvertToBasicType(input any) (any, error) {
	if input == nil {
		return nil, newConvertError(ErrCannotConvertToBasic, input)
	}

	t := reflect.TypeOf(input)
	k := t.Kind()

	if t.String() == k.String() {
		return input, newConvertError(ErrAlreadyBasicType, input)
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

	return nil, newConvertError(ErrCannotConvertToBasic, input)
}

//nolint:ireturn
func tryReflectConvert[Out any](input any) (output Out, err error) {
	defer errorsx.RecoverPanicToError(&err)

	if input == nil {
		return output, newConvertError(ErrConvertInvalidType, input)
	}

	typeOfInput := reflect.TypeOf(input)
	valueOfInput := reflect.ValueOf(input)

	typeOfOutput := reflect.TypeOf(output)

	if !typeOfInput.ConvertibleTo(typeOfOutput) {
		return output, newConvertError(ErrConvertInvalidType, input)
	}

	convertedValue := valueOfInput.Convert(typeOfOutput)

	//nolint:forcetypeassert // if we get here we can safely assert.
	return convertedValue.Interface().(Out), nil
}

var (
	ErrCannotConvertToBasic = errors.New("value cannot be converted to basic type")
	ErrAlreadyBasicType     = errors.New("value is already basic type")
)

//nolint:gochecknoglobals
var defaultConverterGlobal = NewDefaultConverter()

// Convert attempts to convert the input value to the specified (generic) Output type.
// It returns the converted value of type Output and an error if the converting fails.
// e.g.
//
//	Convert[string](123) // "123", nil
func Convert[Output any](input any) (Output, error) { //nolint:ireturn
	var o Output

	n, err := defaultConverterGlobal.ByType(input, reflect.TypeOf(o))
	if err != nil {
		return o, err
	}

	converted, is := n.(Output)
	if !is {
		return o, newConvertError(ErrConvertInvalidType, input)
	}

	return converted, nil
}
