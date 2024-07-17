package slices

import (
	"errors"
	"reflect"

	"github.com/moukoublen/pick/internal/errorsx"
)

type OpMeta struct {
	Index  int
	Length int
}

type Op func(item any, meta OpMeta) error

type CastOp[T any] func(item any, meta OpMeta) (T, error)

type CastFilterOp[T any] func(item any, meta OpMeta) (T, bool, error)

func CastOpFn[T any](fn func(item any) (T, error)) CastOp[T] {
	return func(item any, _ OpMeta) (T, error) {
		return fn(item)
	}
}

func ForEach(input any, operation Op) (rErr error) {
	defer errorsx.RecoverPanicToError(&rErr)

	// attempt to quick return on slice of basic types by avoiding reflect.
	switch cc := input.(type) {
	case []any:
		return each(cc, operation)
	case []string:
		return each(cc, operation)
	case []int:
		return each(cc, operation)
	case []int8:
		return each(cc, operation)
	case []int16:
		return each(cc, operation)
	case []int32:
		return each(cc, operation)
	case []int64:
		return each(cc, operation)
	case []uint:
		return each(cc, operation)
	case []uint8:
		return each(cc, operation)
	case []uint16:
		return each(cc, operation)
	case []uint32:
		return each(cc, operation)
	case []uint64:
		return each(cc, operation)
	case []float32:
		return each(cc, operation)
	case []float64:
		return each(cc, operation)
	case []bool:
		return each(cc, operation)
	}

	typeOfInput := reflect.TypeOf(input)

	if typeOfInput == nil {
		return nil // no op
	}

	// if not slice or array => single operation call attempt
	kindOfInput := typeOfInput.Kind()
	switch kindOfInput {
	case reflect.Array, reflect.Slice:
		valueOfInput := reflect.ValueOf(input)
		length := valueOfInput.Len()
		for i := 0; i < length; i++ {
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
		// single operation call attempt
		return operation(input, OpMeta{Index: 0, Length: 1})

	default:
		// single operation call attempt
		return operation(input, OpMeta{Index: 0, Length: 1})
	}
}

func each[T any](s []T, operation Op) error {
	l := len(s)
	for i := range s {
		if err := operation(s[i], OpMeta{Index: i, Length: l}); err != nil {
			return err
		}
	}

	return nil
}

func AsSlice[T any](input any, castOp CastOp[T]) ([]T, error) {
	return AsSliceFilter(input, func(item any, meta OpMeta) (T, bool, error) {
		casted, err := castOp(item, meta)
		return casted, true, err
	})
}

func AsSliceFilter[T any](input any, castFilterOp CastFilterOp[T]) ([]T, error) {
	// quick returns just in case its already slice of T.
	if ss, is := input.([]T); is {
		return ss, nil
	}

	var castedSlice []T

	if input == nil {
		return castedSlice, nil
	}

	err := ForEach(input, func(item any, meta OpMeta) error {
		casted, keep, err := castFilterOp(item, meta)
		switch {
		case err != nil:
			return err
		case !keep:
			return nil
		default:
			if meta.Index == 0 && meta.Length > 0 {
				castedSlice = make([]T, 0, meta.Length)
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
