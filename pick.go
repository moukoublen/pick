package pick

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"

	"github.com/moukoublen/pick/cast"
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
	m := map[string]any{}
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
	caster    Caster
	notation  Notation
}

func NewPicker(data any, t Traverser, c Caster, n Notation) *Picker {
	return &Picker{
		data:      data,
		traverser: t,
		caster:    c,
		notation:  n,
	}
}

// Wrap returns a new Picker using the same traverser, caster and notation.
func (p *Picker) Wrap(data any) *Picker {
	return NewPicker(data, p.traverser, p.caster, p.notation)
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

// Each applies operation function to each element of the given selector.
// The operation functions receives the index of the element, a SelectorMustAPI
// and the total length of the slice (or 1 if input is a single element and not a slice).
func (p *Picker) Each(selector string, operation func(index int, item any, length int)) (returnedError error) {
	item, err := p.Any(selector)
	if err != nil {
		if errors.Is(err, ErrFieldNotFound) {
			return nil
		}
		return err
	}

	gatherErrors := gatherErrorsFn(&returnedError)

	err = traverseEach(item, func(i int, a any, l int) {
		operation(i, p.Wrap(a).Must(gatherErrors), l)
	})
	if err != nil {
		gather(&returnedError, err)
	}

	return returnedError
}

//nolint:ireturn
func Map[Output any](p *Picker, selector string, mapFn func(*Picker) (Output, error)) ([]Output, error) {
	item, err := p.Any(selector)
	if err != nil {
		return nil, err
	}

	return cast.ToSlice(item, func(a any) (Output, error) { return mapFn(p.Wrap(a)) })
}

// MapMust is like Map but provides a SelectorMust into the map function's argument.
// It also gathers any possible error of Must API to `multipleError` and returns it.
// It's helpful when a clean field-to-field mapping is preferred, but a possible error
// for each field must also be perceived.
// Example:
//
//	itemsSlice, err := MapMust(p, "near_earth_objects.2023-01-07", func(sm SelectorMustAPI) Item {
//		return Item{
//			Name:   sm.String("name"),
//			Sentry: sm.Bool("is_sentry_object"),
//		}
//	})
//
//nolint:ireturn
func MapMust[Output any](p *Picker, selector string, mapFn func(SelectorMustAPI) Output) (_ []Output, returnedError error) {
	item, err := p.Any(selector)
	if err != nil {
		return nil, err
	}

	gatherErrors := gatherErrorsFn(&returnedError)

	sl, err := cast.ToSlice(item, func(a any) (Output, error) {
		return mapFn(p.Wrap(a).Must(gatherErrors)), nil
	})
	if err != nil {
		gather(&returnedError, err)
	}

	return sl, returnedError
}

//nolint:ireturn
func FlatMap[Output any](p *Picker, selector string, mapFn func(*Picker) ([]Output, error)) ([]Output, error) {
	item, err := p.Any(selector)
	if err != nil {
		return nil, err
	}

	doubleSlice, err := cast.ToSlice(item, func(a any) ([]Output, error) { return mapFn(p.Wrap(a)) })
	if err != nil {
		return nil, err
	}
	l := 0
	for i := range doubleSlice {
		l += len(doubleSlice[i])
	}

	outputSlice := make([]Output, 0, l)
	for _, ds := range doubleSlice {
		outputSlice = append(outputSlice, ds...)
	}

	return outputSlice, nil
}

//nolint:ireturn
func Selector[Output any](p *Picker, selector string, castFn func(any) (Output, error)) (Output, error) {
	item, err := p.Any(selector)
	if err != nil {
		var o Output
		return o, err
	}

	return castFn(item)
}

//nolint:ireturn
func SelectorMust[Output any](p *Picker, selector string, castFn func(any) (Output, error), onErr ...func(selector string, err error)) Output {
	casted, err := Selector(p, selector, castFn)
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
