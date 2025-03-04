package pick

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/moukoublen/pick/iter"
)

// This file contains the top-level functions that operates to `Picker` and `SelectorRelaxedAPI`

//
// Top level functions that use default API.
//

// Each iterates over all elements selected by the given selector and applies the provided operation function to each element.
// It returns An error if any step in the selection or operation process fails, otherwise nil.
func Each(p Picker, selector string, operation func(index int, item Picker, totalLength int) error) error {
	item, err := p.Any(selector)
	if err != nil {
		return err
	}

	return iter.ForEach(
		item,
		func(item any, meta iter.CollectionOpMeta) error {
			return operation(meta.Index, p.Wrap(item), meta.Length)
		},
	)
}

// EachField iterates over all fields of the object selected by the given selector and applies the provided operation function to each field's value.
// It returns An error if any step in the selection or operation process fails, otherwise nil.
func EachField(p Picker, selector string, operation func(field string, value Picker, numOfFields int) error) error {
	item, err := p.Any(selector)
	if err != nil {
		return err
	}

	return iter.ForEachField(
		item,
		func(item any, meta iter.FieldOpMeta) error {
			return operation(meta.Name, p.Wrap(item), meta.Length)
		},
	)
}

//nolint:ireturn
func Map[Output any](p Picker, selector string, transform func(Picker) (Output, error)) ([]Output, error) {
	item, err := p.Any(selector)
	if err != nil {
		return nil, err
	}

	return iter.Map(
		item,
		func(item any, _ iter.CollectionOpMeta) (Output, error) {
			return transform(p.Wrap(item))
		},
	)
}

//nolint:ireturn
func MapFilter[Output any](p Picker, selector string, transform func(Picker) (Output, bool, error)) ([]Output, error) {
	item, err := p.Any(selector)
	if err != nil {
		return nil, err
	}

	return iter.MapFilter(
		item,
		func(item any, _ iter.CollectionOpMeta) (Output, bool, error) {
			return transform(p.Wrap(item))
		},
	)
}

//nolint:ireturn
func FlatMap[Output any](p Picker, selector string, transform func(Picker) ([]Output, error)) ([]Output, error) {
	item, err := p.Any(selector)
	if err != nil {
		return nil, err
	}

	doubleSlice, err := iter.Map(
		item,
		func(item any, _ iter.CollectionOpMeta) ([]Output, error) {
			return transform(p.Wrap(item))
		},
	)

	return flatten[Output](doubleSlice), err
}

// Path traverses with the provided path and if found,
// it resolves the convert type from the generic type.
func Path[Output any](p Picker, path ...Key) (Output, error) { //nolint:ireturn
	item, err := p.Path(path)
	if err != nil {
		var o Output
		return o, err
	}

	var defaultValue Output
	return convertAs(p.Converter, item, defaultValue)
}

// OrDefault will return the default value if any error occurs. If the error is ErrFieldNotFound the error will not be returned.
func OrDefault[Output any](p Picker, selector string, defaultValue Output) (Output, error) { //nolint:ireturn
	item, err := p.Any(selector)
	if err != nil {
		if errors.Is(err, ErrFieldNotFound) {
			return defaultValue, nil
		}

		return defaultValue, err
	}

	return convertAs(p.Converter, item, defaultValue)
}

// Get parses the selector, traverses with the provided path and if found,
// it resolves the convert type from the generic type.
func Get[Output any](p Picker, selector string) (Output, error) { //nolint:ireturn
	var defaultValue Output

	item, err := p.Any(selector)
	if err != nil {
		return defaultValue, err
	}

	return convertAs(p.Converter, item, defaultValue)
}

//
// Top level functions that use Relaxed API. (They have the `Relaxed` prefix in name)
//

// RelaxedEach applies operation function to each element of the given selector.
// The operation functions receives the index of the element, a SelectorRelaxedAPI
// and the total length of the slice (or 1 if input is a single element and not a slice).
func RelaxedEach(a RelaxedAPI, selector string, operation func(index int, item RelaxedAPI, length int) error) {
	item, path, err := parseSelectorAndTraverse(a.Picker, selector)
	if err != nil {
		a.gather(selector, err)
		return
	}

	err = iter.ForEach(item, func(item any, meta iter.CollectionOpMeta) error {
		opErr := operation(meta.Index, a.Wrap(item), meta.Length)
		if opErr != nil {
			path = append(path, Index(meta.Index))
			a.gather(a.notation.Format(path...), opErr)
			path = path[:len(path)-1]
		}
		return nil
	})
	if err != nil {
		a.gather(selector, err)
	}
}

// RelaxedEachField applies operation function to each field of the object of the given selector.
// The operation functions receives the name of the field of the element, a SelectorRelaxedAPI
// and the total length of the slice (or 1 if input is a single element and not a slice).
func RelaxedEachField(a RelaxedAPI, selector string, operation func(field string, value RelaxedAPI, numOfFields int) error) {
	item, path, err := parseSelectorAndTraverse(a.Picker, selector)
	if err != nil {
		a.gather(selector, err)
		return
	}

	err = iter.ForEachField(item, func(item any, meta iter.FieldOpMeta) error {
		opErr := operation(meta.Name, a.Wrap(item), meta.Length)
		if opErr != nil {
			path = append(path, Field(meta.Name))
			a.gather(a.notation.Format(path...), opErr)
			path = path[:len(path)-1]
		}
		return nil
	})
	if err != nil {
		a.gather(selector, err)
	}
}

// RelaxedMap transform each element of a slice (or a the single element if selector leads to not slice)
// by applying the transform function.
// It also gathers any possible error of Relaxed API to `multipleError` and returns it.
// Example:
//
//	itemsSlice, err := RelaxedMap(p.Relaxed(), "near_earth_objects.2023-01-07", func(sm SelectorRelaxedAPI) Item {
//		return Item{
//			Name:   sm.String("name"),
//			Sentry: sm.Bool("is_sentry_object"),
//		}
//	})
//
//nolint:ireturn
func RelaxedMap[Output any](a RelaxedAPI, selector string, transform func(RelaxedAPI) (Output, error)) []Output {
	return RelaxedMapFilter[Output](a, selector, func(sma RelaxedAPI) (Output, bool, error) {
		o, err := transform(sma)
		return o, true, err
	})
}

//nolint:ireturn
func RelaxedMapFilter[Output any](a RelaxedAPI, selector string, transform func(RelaxedAPI) (Output, bool, error)) []Output {
	item, path, err := parseSelectorAndTraverse(a.Picker, selector)
	if err != nil {
		a.gather(selector, err)
		return nil
	}

	sl, err := iter.MapFilter(item, func(item any, meta iter.CollectionOpMeta) (Output, bool, error) {
		t, keep, opErr := transform(a.Wrap(item))
		if opErr != nil {
			path = append(path, Index(meta.Index))
			a.gather(a.notation.Format(path...), opErr)
			path = path[:len(path)-1]
		}
		return t, keep, nil
	})
	if err != nil {
		a.gather(selector, err)
	}

	return sl
}

//nolint:ireturn
func RelaxedFlatMap[Output any](a RelaxedAPI, selector string, transform func(RelaxedAPI) ([]Output, error)) []Output {
	item := RelaxedMap[[]Output](a, selector, transform)
	return flatten[Output](item)
}

// RelaxedPath is the version of [Path] that uses SelectorRelaxedAPI.
func RelaxedPath[Output any](a RelaxedAPI, path ...Key) Output { //nolint:ireturn
	converted, err := Path[Output](a.Picker, path...)
	if err != nil {
		selector := DotNotation{}.Format(path...)
		a.gather(selector, err)
	}
	return converted
}

// RelaxedOrDefault will return the default value if any error occurs. Version of [OrDefault] that uses SelectorRelaxedAPI.
func RelaxedOrDefault[Output any](a RelaxedAPI, selector string, defaultValue Output) Output { //nolint:ireturn
	item, err := OrDefault(a.Picker, selector, defaultValue)
	if err != nil {
		a.gather(selector, err)
	}

	return item
}

// RelaxedGet resolves the convert type from the generic type. Version of [Get] that uses SelectorRelaxedAPI.
func RelaxedGet[Output any](a RelaxedAPI, selector string) Output { //nolint:ireturn
	item, err := Get[Output](a.Picker, selector)
	if err != nil {
		a.gather(selector, err)
	}

	return item
}

func flatten[Output any](doubleSlice [][]Output) []Output {
	// calculate total capacity
	l := 0
	for i := range doubleSlice {
		l += len(doubleSlice[i])
	}

	// flatten
	outputSlice := make([]Output, 0, l)
	for _, ds := range doubleSlice {
		outputSlice = append(outputSlice, ds...)
	}

	return outputSlice
}

func parseSelectorAndTraverse(p Picker, selector string) (any, []Key, error) {
	path, err := p.notation.Parse(selector)
	if err != nil {
		return nil, path, err
	}

	item, err := p.Path(path)
	return item, path, err
}

func convertAs[Output any](converter Converter, data any, defaultValue Output) (Output, error) { //nolint:ireturn
	c, err := converter.ByType(data, reflect.TypeOf(defaultValue))
	if err != nil {
		return defaultValue, err
	}

	asOutput, is := c.(Output)
	if !is {
		return defaultValue, fmt.Errorf("converted value cannot be asserted to type: %w", ErrConvertInvalidType) // this is not possible
	}

	return asOutput, nil
}
