package pick

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"

	"github.com/moukoublen/pick/cast"
	"github.com/moukoublen/pick/internal"
)

type Notation interface {
	Parse(selector string) ([]Key, error)
	Format(path ...Key) string
}

type Traverser interface {
	Retrieve(data any, path []Key) (any, error)
}

func WrapJSON(js []byte) (*Picker, error) {
	return WrapReaderJSON(bytes.NewReader(js))
}

func WrapReaderJSON(r io.Reader) (*Picker, error) {
	d := json.NewDecoder(r)
	return WrapDecoder(d)
}

func WrapDecoder(decoder interface{ Decode(destination any) error }) (*Picker, error) {
	var m any
	if err := decoder.Decode(&m); err != nil {
		return nil, err
	}

	return Wrap(m), nil
}

func Wrap(data any) *Picker {
	caster := cast.NewCaster()
	return NewPicker(data, NewDefaultTraverser(caster), caster, DotNotation{})
}

type Picker struct {
	data      any
	traverser Traverser
	Caster    Caster
	notation  Notation
}

func NewPicker(data any, t Traverser, c Caster, n Notation) *Picker {
	return &Picker{
		data:      data,
		traverser: t,
		Caster:    c,
		notation:  n,
	}
}

// Wrap returns a new Picker using the same traverser, caster and notation.
func (p *Picker) Wrap(data any) *Picker {
	return NewPicker(data, p.traverser, p.Caster, p.notation)
}

func (p *Picker) Data() any { return p.data }

func (p *Picker) Must(onErr ...func(string, error)) SelectorMustAPI {
	return SelectorMustAPI{Picker: p, onErr: onErr}
}

func (p *Picker) Any(selector string) (any, error) {
	path, err := p.notation.Parse(selector)
	if err != nil {
		return nil, err
	}

	return p.Traverse(path)
}

func (p *Picker) Traverse(path []Key) (any, error) {
	return p.traverser.Retrieve(p.data, path)
}

//
// Top level functions that use default API.
//

func Each(p *Picker, selector string, operation func(index int, p *Picker, length int) error) error {
	item, err := p.Any(selector)
	if err != nil {
		return err
	}

	return internal.TraverseSlice(
		item,
		func(i int, a any, l int) error {
			return operation(i, p.Wrap(a), l)
		},
	)
}

//nolint:ireturn
func Map[Output any](p *Picker, selector string, transform func(*Picker) (Output, error)) ([]Output, error) {
	item, err := p.Any(selector)
	if err != nil {
		return nil, err
	}

	return cast.ToSlice(
		item,
		func(_ int, a any, _ int) (Output, error) {
			return transform(p.Wrap(a))
		},
	)
}

//nolint:ireturn
func MapFilter[Output any](p *Picker, selector string, transform func(*Picker) (Output, bool, error)) ([]Output, error) {
	item, err := p.Any(selector)
	if err != nil {
		return nil, err
	}

	return cast.ToSliceFilter(
		item,
		func(_ int, a any, _ int) (Output, bool, error) {
			return transform(p.Wrap(a))
		},
	)
}

//nolint:ireturn
func FlatMap[Output any](p *Picker, selector string, transform func(*Picker) ([]Output, error)) ([]Output, error) {
	item, err := p.Any(selector)
	if err != nil {
		return nil, err
	}

	doubleSlice, err := cast.ToSlice(
		item,
		func(_ int, a any, _ int) ([]Output, error) {
			return transform(p.Wrap(a))
		},
	)

	return flatten[Output](doubleSlice), err
}

//nolint:ireturn
func pickSelector[Output any](p *Picker, selector string, castFn func(any) (Output, error)) (Output, error) {
	item, err := p.Any(selector)
	if err != nil {
		var o Output
		return o, err
	}

	return castFn(item)
}

//nolint:ireturn
func pickSelectorMust[Output any](p *Picker, selector string, castFn func(any) (Output, error), onErr ...func(selector string, err error)) Output {
	casted, err := pickSelector(p, selector, castFn)
	if err != nil {
		for _, fn := range onErr {
			fn(selector, err)
		}
	}
	return casted
}

//nolint:ireturn
func Path[Output any](p *Picker, path []Key, castFn func(any) (Output, error)) (Output, error) {
	item, err := p.Traverse(path)
	if err != nil {
		var o Output
		return o, err
	}

	return castFn(item)
}

//nolint:ireturn
func PathMust[Output any](p *Picker, path []Key, castFn func(any) (Output, error), onErr ...func(selector string, err error)) Output {
	casted, err := Path(p, path, castFn)
	if err != nil {
		selector := DotNotation{}.Format(path...)
		for _, fn := range onErr {
			fn(selector, err)
		}
	}
	return casted
}

//
// Top level functions that use must API. (They have the `M` postfix in name)
//

// EachM applies operation function to each element of the given selector.
// The operation functions receives the index of the element, a SelectorMustAPI
// and the total length of the slice (or 1 if input is a single element and not a slice).
func EachM(p *Picker, selector string, operation func(index int, item SelectorMustAPI, length int) error) (returnedError error) {
	item, err := p.Any(selector)
	if err != nil {
		if errors.Is(err, ErrFieldNotFound) {
			return nil
		}
		return err
	}

	gatherErrors := gatherErrorsFn(&returnedError)

	err = internal.TraverseSlice(
		item,
		func(i int, a any, l int) error {
			opErr := operation(i, p.Wrap(a).Must(gatherErrors), l)
			if opErr != nil {
				gatherErrors(selector, opErr)
			}
			return nil
		},
	)
	if err != nil {
		gather(&returnedError, err)
	}

	return returnedError
}

// MapM transform each element of a slice (or a the single element if selector leads to not slice)
// by applying the transform function.
// It also gathers any possible error of Must API to `multipleError` and returns it.
// Example:
//
//	itemsSlice, err := MapM(p, "near_earth_objects.2023-01-07", func(sm SelectorMustAPI) Item {
//		return Item{
//			Name:   sm.String("name"),
//			Sentry: sm.Bool("is_sentry_object"),
//		}
//	})
//
//nolint:ireturn
func MapM[Output any](p *Picker, selector string, transform func(SelectorMustAPI) (Output, error)) (_ []Output, returnedError error) {
	item, err := p.Any(selector)
	if err != nil {
		if errors.Is(err, ErrFieldNotFound) {
			return nil, nil
		}
		return nil, err
	}

	gatherErrors := gatherErrorsFn(&returnedError)

	sl, err := cast.ToSlice(
		item,
		func(_ int, a any, _ int) (Output, error) {
			t, opErr := transform(p.Wrap(a).Must(gatherErrors))
			if opErr != nil {
				gatherErrors(selector, opErr)
			}
			return t, nil
		},
	)
	if err != nil {
		gather(&returnedError, err)
	}

	return sl, returnedError
}

//nolint:ireturn
func MapFilterM[Output any](p *Picker, selector string, transform func(SelectorMustAPI) (Output, bool, error)) (_ []Output, returnedError error) {
	item, err := p.Any(selector)
	if err != nil {
		if errors.Is(err, ErrFieldNotFound) {
			return nil, nil
		}
		return nil, err
	}

	gatherErrors := gatherErrorsFn(&returnedError)

	sl, err := cast.ToSliceFilter(
		item,
		func(_ int, a any, _ int) (Output, bool, error) {
			t, keep, opErr := transform(p.Wrap(a).Must(gatherErrors))
			if opErr != nil {
				gatherErrors(selector, opErr)
			}
			return t, keep, nil
		},
	)
	if err != nil {
		gather(&returnedError, err)
	}

	return sl, returnedError
}

//nolint:ireturn
func FlatMapM[Output any](p *Picker, selector string, transform func(SelectorMustAPI) ([]Output, error)) (_ []Output, returnedError error) {
	item, err := p.Any(selector)
	if err != nil {
		if errors.Is(err, ErrFieldNotFound) {
			return nil, nil
		}
		return nil, err
	}

	gatherErrors := gatherErrorsFn(&returnedError)

	doubleSlice, err := cast.ToSlice(
		item,
		func(_ int, a any, _ int) ([]Output, error) {
			t, opErr := transform(p.Wrap(a).Must(gatherErrors))
			if opErr != nil {
				gatherErrors(selector, opErr)
			}
			return t, nil
		},
	)
	if err != nil {
		gather(&returnedError, err)
	}

	return flatten[Output](doubleSlice), returnedError
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
