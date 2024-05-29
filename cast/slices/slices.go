package slices

import (
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

	if input == nil {
		return operation(input, OpMeta{Index: 0, Length: 1})
	}

	typeOfInput := reflect.TypeOf(input)
	kindOfInput := typeOfInput.Kind()

	// if not slice or array => single operation call attempt
	if kindOfInput != reflect.Array && kindOfInput != reflect.Slice {
		return operation(input, OpMeta{Index: 0, Length: 1})
	}

	// slow/costly attempt with reflect
	valueOfInput := reflect.ValueOf(input)
	length := valueOfInput.Len()
	for i := 0; i < length; i++ {
		item := valueOfInput.Index(i)
		if err := operation(item.Interface(), OpMeta{Index: i, Length: length}); err != nil {
			return err
		}
	}

	return rErr
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
