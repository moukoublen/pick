package internal

import (
	"reflect"

	"github.com/moukoublen/pick/internal/testingx/errorsx"
)

func TraverseSlice(input any, operation func(index int, item any, length int) error) (rErr error) {
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
	kindOfInput := typeOfInput.Kind()

	// if not slice or array => single operation call attempt
	if kindOfInput != reflect.Array && kindOfInput != reflect.Slice {
		return operation(0, input, 1)
	}

	// slow/costly attempt with reflect
	valueOfInput := reflect.ValueOf(input)
	length := valueOfInput.Len()
	for i := 0; i < length; i++ {
		item := valueOfInput.Index(i)
		if err := operation(i, item.Interface(), length); err != nil {
			return err
		}
	}

	return rErr
}

func each[T any](s []T, operation func(index int, item any, length int) error) error {
	l := len(s)
	for i := range s {
		if err := operation(i, s[i], l); err != nil {
			return err
		}
	}

	return nil
}
