package pick

import (
	"errors"
	"reflect"

	"github.com/moukoublen/pick/internal/errorsx"
)

type iterationOpMeta struct {
	Index  int
	Length int
}

// forEach applies the given operation to each element of the input if it is a collection
// (slice or array), or directly if it is a single item (pointer, interface, or other type).
//
// The function first tries to handle slices of basic types directly by avoiding reflection for performance reasons.
// If the input is not one of the directly handled types, it uses reflection to determine the input type and
// iterates over elements if it is a slice or array. For pointers or interfaces, it dereferences the input and
// applies the operation to the dereferenced value. For other types, it applies the operation directly.
//
// The function uses deferred recovery to capture and return any panic as an error.
func forEach(input any, operation func(item any, meta iterationOpMeta) error) (rErr error) { //nolint:gocyclo
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
			if err := operation(item.Interface(), iterationOpMeta{Index: i, Length: length}); err != nil {
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
		return operation(el.Interface(), iterationOpMeta{Index: 0, Length: 1})

	default:
		// single operation call attempt
		return operation(input, iterationOpMeta{Index: 0, Length: 1})
	}
}

func forEachSlice[T any](s []T, operation func(item any, meta iterationOpMeta) error) error {
	l := len(s)
	for i := range s {
		if err := operation(s[i], iterationOpMeta{Index: i, Length: l}); err != nil {
			return err
		}
	}

	return nil
}

func mapTo[T any](input any, castOp func(item any, meta iterationOpMeta) (T, error)) ([]T, error) {
	// quick returns just in case its already slice of T.
	if ss, is := input.([]T); is {
		return ss, nil
	}

	return mapFilterTo(input, func(item any, meta iterationOpMeta) (T, bool, error) {
		casted, err := castOp(item, meta)
		return casted, err == nil, err
	})
}

func mapFilterTo[T any](input any, castFilterOp func(item any, meta iterationOpMeta) (T, bool, error)) ([]T, error) {
	var castedSlice []T

	if input == nil {
		return castedSlice, nil
	}

	err := forEach(input, func(item any, meta iterationOpMeta) error {
		casted, keep, err := castFilterOp(item, meta)
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

func mapOpFn[T any](fn func(item any) (T, error)) func(item any, meta iterationOpMeta) (T, error) {
	return func(item any, _ iterationOpMeta) (T, error) {
		return fn(item)
	}
}

// itemLen returns the result of built in len function if the input is of type slice, array, map, string or channel.
// If the input is pointer and not nil, it dereferences the destination.
func itemLen(input any) (l int, rErr error) {
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
		return itemLen(elemValue.Interface())
	}

	return -1, ErrNoLength
}

var ErrNoLength = errors.New("type does not have length")
