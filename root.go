package pick

import (
	"github.com/moukoublen/pick/cast/slices"
)

// This file contains the top-level functions that operates to `Picker` and `SelectorMustAPI`

//
// Top level functions that use default API.
//

// Each iterates over all elements selected by the given selector and applies the provided operation function to each element.
// It returns An error if any step in the selection or operation process fails, otherwise nil.
func Each(p *Picker, selector string, operation func(index int, p *Picker, totalLength int) error) error {
	item, err := p.Any(selector)
	if err != nil {
		return err
	}

	return slices.ForEach(
		item,
		func(item any, meta slices.OpMeta) error {
			return operation(meta.Index, p.Wrap(item), meta.Length)
		},
	)
}

//nolint:ireturn
func Map[Output any](p *Picker, selector string, transform func(*Picker) (Output, error)) ([]Output, error) {
	item, err := p.Any(selector)
	if err != nil {
		return nil, err
	}

	return slices.Map(
		item,
		func(item any, _ slices.OpMeta) (Output, error) {
			return transform(p.Wrap(item))
		},
	)
}

//nolint:ireturn
func MapFilter[Output any](p *Picker, selector string, transform func(*Picker) (Output, bool, error)) ([]Output, error) {
	item, err := p.Any(selector)
	if err != nil {
		return nil, err
	}

	return slices.MapFilter(
		item,
		func(item any, _ slices.OpMeta) (Output, bool, error) {
			return transform(p.Wrap(item))
		},
	)
}

//nolint:ireturn
func FlatMap[Output any](p *Picker, selector string, transform func(*Picker) ([]Output, error)) ([]Output, error) {
	item, err := p.Any(selector)
	if err != nil {
		return nil, err
	}

	doubleSlice, err := slices.Map(
		item,
		func(item any, _ slices.OpMeta) ([]Output, error) {
			return transform(p.Wrap(item))
		},
	)

	return flatten[Output](doubleSlice), err
}

//nolint:ireturn
func Path[Output any](p *Picker, path []Key, castFn func(any) (Output, error)) (Output, error) {
	item, err := p.Path(path)
	if err != nil {
		var o Output
		return o, err
	}

	return castFn(item)
}

//
// Top level functions that use must API. (They have the `Must` prefix in name)
//

// MustEach applies operation function to each element of the given selector.
// The operation functions receives the index of the element, a SelectorMustAPI
// and the total length of the slice (or 1 if input is a single element and not a slice).
func MustEach(a SelectorMustAPI, selector string, operation func(index int, item SelectorMustAPI, length int) error) {
	item, path, err := parseSelectorAndTraverse(a.Picker, selector)
	if err != nil {
		a.gather(selector, err)
		return
	}

	err = slices.ForEach(item, func(item any, meta slices.OpMeta) error {
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

// MustMap transform each element of a slice (or a the single element if selector leads to not slice)
// by applying the transform function.
// It also gathers any possible error of Must API to `multipleError` and returns it.
// Example:
//
//	itemsSlice, err := MustMap(p.Must(), "near_earth_objects.2023-01-07", func(sm SelectorMustAPI) Item {
//		return Item{
//			Name:   sm.String("name"),
//			Sentry: sm.Bool("is_sentry_object"),
//		}
//	})
//
//nolint:ireturn
func MustMap[Output any](a SelectorMustAPI, selector string, transform func(SelectorMustAPI) (Output, error)) []Output {
	return MustMapFilter[Output](a, selector, func(sma SelectorMustAPI) (Output, bool, error) {
		o, err := transform(sma)
		return o, true, err
	})
}

//nolint:ireturn
func MustMapFilter[Output any](a SelectorMustAPI, selector string, transform func(SelectorMustAPI) (Output, bool, error)) []Output {
	item, path, err := parseSelectorAndTraverse(a.Picker, selector)
	if err != nil {
		a.gather(selector, err)
		return nil
	}

	sl, err := slices.MapFilter(item, func(item any, meta slices.OpMeta) (Output, bool, error) {
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
func MustFlatMap[Output any](a SelectorMustAPI, selector string, transform func(SelectorMustAPI) ([]Output, error)) []Output {
	item := MustMap[[]Output](a, selector, transform)
	return flatten[Output](item)
}

//nolint:ireturn
func MustPath[Output any](a SelectorMustAPI, path []Key, castFn func(any) (Output, error)) Output {
	casted, err := Path(a.Picker, path, castFn)
	if err != nil {
		selector := DotNotation{}.Format(path...)
		a.gather(selector, err)
	}
	return casted
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

func parseSelectorAndTraverse(p *Picker, selector string) (any, []Key, error) {
	path, err := p.notation.Parse(selector)
	if err != nil {
		return nil, path, err
	}

	item, err := p.Path(path)
	return item, path, err
}
