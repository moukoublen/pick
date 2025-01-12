package iter

import (
	"errors"
	"reflect"

	"github.com/moukoublen/pick/internal/errorsx"
)

type FieldOpMeta struct {
	Field  any
	Length int
}

func ForEachField(input any, operation func(item any, meta FieldOpMeta) error) (rErr error) { //nolint:gocyclo
	defer errorsx.RecoverPanicToError(&rErr)

	// attempt to quick return on map of basic types by avoiding reflect.
	switch cc := input.(type) {
	case map[string]any:
		return forEachMap(cc, operation)
	case map[string]string:
		return forEachMap(cc, operation)
	}

	typeOfInput := reflect.TypeOf(input)

	if typeOfInput == nil {
		return nil // no op
	}

	kindOfInput := typeOfInput.Kind()
	switch kindOfInput {
	case reflect.Array, reflect.Slice:
		return ErrNoFields

	case reflect.Pointer, reflect.Interface:
		valueOfInput := reflect.ValueOf(input)
		if valueOfInput.IsNil() {
			return rErr
		}

		// deref
		el := valueOfInput.Elem()
		if !el.IsValid() {
			return rErr
		}

		return ForEachField(el.Interface(), operation)

	case reflect.Struct:
		valueOfInput := reflect.ValueOf(input)
		num := valueOfInput.NumField()
		for i := range num {
			fieldVal := valueOfInput.Field(i)
			fieldType := typeOfInput.Field(i)
			err := operation(fieldVal.Interface(), FieldOpMeta{Field: fieldType.Name, Length: num})
			if err != nil {
				return err
			}
		}

	case reflect.Map:
		valueOfInput := reflect.ValueOf(input)
		itr := valueOfInput.MapRange()
		num := valueOfInput.Len()
		for itr.Next() {
			k := itr.Key()
			v := itr.Value()
			err := operation(v.Interface(), FieldOpMeta{Field: k.Interface(), Length: num})
			if err != nil {
				return err
			}
		}

	default:
		// single operation call attempt
		return ErrNoFields
	}

	return rErr
}

var ErrNoFields = errors.New("type does not have fields")

func forEachMap[V any](m map[string]V, operation func(item any, meta FieldOpMeta) error) error {
	l := len(m)
	for k, v := range m {
		if err := operation(v, FieldOpMeta{Field: k, Length: l}); err != nil {
			return err
		}
	}

	return nil
}

type OpMeta struct {
	Index  int
	Length int
}

// ForEach applies the given operation to each element of the input if it is a collection
// (slice or array), or directly if it is a single item (pointer, interface, or other type).
// If the operation returns a non-nil error it will cause the entire ForEach function to terminate without applying the
// operation to the rest of the elements (if any).
//
// The function first tries to handle slices of basic types directly by avoiding reflection for performance reasons.
// If the input is not one of the directly handled types, it uses reflection to determine the input type and
// iterates over elements if it is a slice or array. For pointers or interfaces, it dereferences the input and
// applies the operation to the dereferenced value. For other types, it applies the operation directly.
//
// The function uses deferred recovery to capture and return any panic as an error.
func ForEach(input any, operation func(item any, meta OpMeta) error) (rErr error) { //nolint:gocyclo
	defer errorsx.RecoverPanicToError(&rErr)

	// attempt to quick return on slice of basic types by avoiding reflect.
	switch cc := input.(type) {
	case []any:
		return forEachSlice(cc, operation)
	case []string:
		return forEachSlice(cc, operation)
	case []int:
		return forEachSlice(cc, operation)
	case []int8:
		return forEachSlice(cc, operation)
	case []int16:
		return forEachSlice(cc, operation)
	case []int32:
		return forEachSlice(cc, operation)
	case []int64:
		return forEachSlice(cc, operation)
	case []uint:
		return forEachSlice(cc, operation)
	case []uint8:
		return forEachSlice(cc, operation)
	case []uint16:
		return forEachSlice(cc, operation)
	case []uint32:
		return forEachSlice(cc, operation)
	case []uint64:
		return forEachSlice(cc, operation)
	case []float32:
		return forEachSlice(cc, operation)
	case []float64:
		return forEachSlice(cc, operation)
	case []bool:
		return forEachSlice(cc, operation)
	}

	typeOfInput := reflect.TypeOf(input)

	if typeOfInput == nil {
		return nil // no op
	}

	kindOfInput := typeOfInput.Kind()
	switch kindOfInput {
	case reflect.Array, reflect.Slice:
		valueOfInput := reflect.ValueOf(input)
		length := valueOfInput.Len()
		for i := range length {
			item := valueOfInput.Index(i)
			if err := operation(item.Interface(), OpMeta{Index: i, Length: length}); err != nil {
				return err
			}
		}

		return rErr

	case reflect.Pointer, reflect.Interface:
		valueOfInput := reflect.ValueOf(input)
		if valueOfInput.IsNil() {
			return rErr
		}

		// deref
		el := valueOfInput.Elem()
		if !el.IsValid() {
			return rErr
		}

		// single operation call attempt
		return operation(el.Interface(), OpMeta{Index: 0, Length: 1})

	default:
		// single operation call attempt
		return operation(input, OpMeta{Index: 0, Length: 1})
	}
}

func forEachSlice[T any](s []T, operation func(item any, meta OpMeta) error) error {
	l := len(s)
	for i := range s {
		if err := operation(s[i], OpMeta{Index: i, Length: l}); err != nil {
			return err
		}
	}

	return nil
}

// Map applies a transformation operation to each element of the provided input if it is a collection
// (slice or array), or directly to it, if it is a single item (pointer, interface, or other type).
func Map[T any](input any, operation func(item any, meta OpMeta) (T, error)) ([]T, error) {
	// quick returns just in case its already slice of T.
	if ss, is := input.([]T); is {
		return ss, nil
	}

	return MapFilter(input, func(item any, meta OpMeta) (T, bool, error) {
		casted, err := operation(item, meta)
		return casted, true, err
	})
}

// MapFilter applies a transformation and filtering operation to each element of the provided input if it is a collection
// (slice or array), or directly to it, if it is a single item (pointer, interface, or other type).
// The operation function returns three values:
//  1. A transformed value of type T
//  2. A boolean indicating whether the transformed value should be included (If the boolean is false for a given element, that element is excluded)
//  3. An error, which if non-nil will cause the entire MapFilter operation to terminate
//
// It returns a slice containing all included transformed elements, or an error if any operation fails.
func MapFilter[T any](input any, operation func(item any, meta OpMeta) (T, bool, error)) ([]T, error) {
	var castedSlice []T

	if input == nil {
		return castedSlice, nil
	}

	err := ForEach(input, func(item any, meta OpMeta) error {
		casted, keep, err := operation(item, meta)
		switch {
		case err != nil:
			return err
		case !keep:
			return nil
		default:
			if meta.Index == 0 && meta.Length > 0 {
				castedSlice = make([]T, 0, meta.Length)
			} else if meta.Length == 0 {
				return nil
			}
			castedSlice = append(castedSlice, casted)
			return nil
		}
	})
	if err != nil {
		return nil, err
	}

	return castedSlice, nil
}

func MapOpFn[T any](fn func(item any) (T, error)) func(item any, meta OpMeta) (T, error) {
	return func(item any, _ OpMeta) (T, error) {
		return fn(item)
	}
}

// Len returns the result of built in len function if the input is of type slice, array, map, string or channel.
// If the input is pointer and not nil, it dereferences the destination.
func Len(input any) (l int, rErr error) {
	defer errorsx.RecoverPanicToError(&rErr)

	// attempt to quick return on slice of basic types by avoiding reflect.
	switch cc := input.(type) {
	case []any:
		return len(cc), nil
	case []map[string]any:
		return len(cc), nil
	case []string:
		return len(cc), nil
	case []int:
		return len(cc), nil
	case []int8:
		return len(cc), nil
	case []int16:
		return len(cc), nil
	case []int32:
		return len(cc), nil
	case []int64:
		return len(cc), nil
	case []uint:
		return len(cc), nil
	case []uint8:
		return len(cc), nil
	case []uint16:
		return len(cc), nil
	case []uint32:
		return len(cc), nil
	case []uint64:
		return len(cc), nil
	case []float32:
		return len(cc), nil
	case []float64:
		return len(cc), nil
	case []bool:
		return len(cc), nil
	case string:
		return len(cc), nil
	case map[string]any:
		return len(cc), nil
	}

	typeOfInput := reflect.TypeOf(input)

	if typeOfInput == nil {
		return -1, ErrNoLength
	}

	kindOfInput := typeOfInput.Kind()

	switch kindOfInput {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		valueOfInput := reflect.ValueOf(input)
		return valueOfInput.Len(), nil
	case reflect.Pointer, reflect.Interface:
		valueOfInput := reflect.ValueOf(input)
		if valueOfInput.IsNil() {
			return -1, ErrNoLength
		}

		elemValue := valueOfInput.Elem()
		return Len(elemValue.Interface())
	}

	return -1, ErrNoLength
}

var ErrNoLength = errors.New("type does not have length")
