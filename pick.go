package pick

import (
	"bytes"
	"encoding/json"
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

func (p *Picker) Path() PathAPI {
	return PathAPI{Picker: p}
}

func (p *Picker) PathMust(onErr ...func(string, error)) PathMustAPI {
	return PathMustAPI{Picker: p, onErr: onErr}
}

//nolint:ireturn
func Map[Output any](p *Picker, selector string, mapFn func(*Picker) (Output, error)) ([]Output, error) {
	item, err := pickSelector(p.data, p.notation, p.traverser, selector, omitCast)
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
	item, err := pickSelector(p.data, p.notation, p.traverser, selector, omitCast)
	if err != nil {
		return nil, err
	}

	gatherErrors := GatherErrorsFn(&returnedError)

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
	item, err := pickSelector(p.data, p.notation, p.traverser, selector, omitCast)
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
	return pickSelector(p.data, p.notation, p.traverser, selector, castFn)
}

//nolint:ireturn
func Path[Output any](p *Picker, path []Key, castFn func(any) (Output, error)) (Output, error) {
	return pickPath(p.data, p.traverser, path, castFn)
}

//nolint:ireturn
func SelectorMust[Output any](p *Picker, selector string, castFn func(any) (Output, error), onErr ...func(selector string, err error)) Output {
	return pickSelectorMust(p.data, p.notation, p.traverser, selector, castFn, onErr...)
}

//nolint:ireturn
func PathMust[Output any](p *Picker, path []Key, castFn func(any) (Output, error), onErr ...func(selector string, err error)) Output {
	return pickPathMust(p.data, p.traverser, path, castFn, onErr...)
}

//nolint:ireturn
func pickSelector[Output any](data any, notation Notation, traverser Traverser, selector string, castFn func(any) (Output, error)) (Output, error) {
	path, err := notation.Parse(selector)
	if err != nil {
		var d Output
		return d, err
	}

	return pickPath(data, traverser, path, castFn)
}

//nolint:ireturn
func pickPath[Output any](data any, traverser Traverser, path []Key, castFn func(any) (Output, error)) (Output, error) {
	item, err := traverser.Retrieve(data, path)
	if err != nil {
		var d Output
		return d, err
	}

	casted, err := castFn(item)
	return casted, err
}

//nolint:ireturn
func pickSelectorMust[Output any](data any, notation Notation, traverser Traverser, selector string, castFn func(any) (Output, error), onErr ...func(selector string, err error)) Output {
	casted, err := pickSelector(data, notation, traverser, selector, castFn)
	if err != nil {
		for _, fn := range onErr {
			fn(selector, err)
		}
	}
	return casted
}

//nolint:ireturn
func pickPathMust[Output any](data any, traverser Traverser, path []Key, castFn func(any) (Output, error), onErr ...func(selector string, err error)) Output {
	casted, err := pickPath(data, traverser, path, castFn)
	if err != nil {
		selector := DotNotation{}.Format(path...)
		for _, fn := range onErr {
			fn(selector, err)
		}
	}
	return casted
}

func omitCast(a any) (any, error) { return a, nil }
