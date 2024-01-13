package pick

import (
	"bytes"
	"encoding/json"
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

	return p.Path(path)
}

func (p *Picker) Path(path []Key) (any, error) {
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
func Path[Output any](p *Picker, path []Key, castFn func(any) (Output, error)) (Output, error) {
	item, err := p.Path(path)
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
func EachM(a SelectorMustAPI, selector string, operation func(index int, item SelectorMustAPI, length int) error) {
	item, path, err := parseSelectorAndTraverse(a.Picker, selector)
	if err != nil {
		a.gather(selector, err)
		return
	}

	sliceOp := func(idx int, item any, l int) error {
		opErr := operation(idx, a.Wrap(item), l)
		if opErr != nil {
			path = append(path, Index(idx))
			a.gather(a.notation.Format(path...), opErr)
			path = path[:len(path)-1]
		}
		return nil
	}

	err = internal.TraverseSlice(item, sliceOp)
	if err != nil {
		a.gather(selector, err)
	}
}

// MapM transform each element of a slice (or a the single element if selector leads to not slice)
// by applying the transform function.
// It also gathers any possible error of Must API to `multipleError` and returns it.
// Example:
//
//	itemsSlice, err := MapM(p.Must(), "near_earth_objects.2023-01-07", func(sm SelectorMustAPI) Item {
//		return Item{
//			Name:   sm.String("name"),
//			Sentry: sm.Bool("is_sentry_object"),
//		}
//	})
//
//nolint:ireturn
func MapM[Output any](a SelectorMustAPI, selector string, transform func(SelectorMustAPI) (Output, error)) []Output {
	return MapFilterM[Output](a, selector, func(sma SelectorMustAPI) (Output, bool, error) {
		o, err := transform(sma)
		return o, true, err
	})
}

//nolint:ireturn
func MapFilterM[Output any](a SelectorMustAPI, selector string, transform func(SelectorMustAPI) (Output, bool, error)) []Output {
	item, path, err := parseSelectorAndTraverse(a.Picker, selector)
	if err != nil {
		a.gather(selector, err)
		return nil
	}

	sliceOp := func(idx int, item any, _ int) (Output, bool, error) {
		t, keep, opErr := transform(a.Wrap(item))
		if opErr != nil {
			path = append(path, Index(idx))
			a.gather(a.notation.Format(path...), opErr)
			path = path[:len(path)-1]
		}
		return t, keep, nil
	}

	sl, err := cast.ToSliceFilter(item, sliceOp)
	if err != nil {
		a.gather(selector, err)
	}

	return sl
}

//nolint:ireturn
func FlatMapM[Output any](a SelectorMustAPI, selector string, transform func(SelectorMustAPI) ([]Output, error)) []Output {
	item := MapM[[]Output](a, selector, transform)
	return flatten[Output](item)
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
