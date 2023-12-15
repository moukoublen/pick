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

func (p *Picker) Bool(selector string) (bool, error) {
	return Selector(p.data, p.notation, p.traverser, selector, p.caster.AsBool)
}

func (p *Picker) BoolSlice(selector string) ([]bool, error) {
	return Selector(p.data, p.notation, p.traverser, selector, p.caster.AsBoolSlice)
}

func (p *Picker) Byte(selector string) (byte, error) {
	return Selector(p.data, p.notation, p.traverser, selector, p.caster.AsByte)
}

func (p *Picker) ByteSlice(selector string) ([]byte, error) {
	return Selector(p.data, p.notation, p.traverser, selector, p.caster.AsByteSlice)
}

func (p *Picker) Float32(selector string) (float32, error) {
	return Selector(p.data, p.notation, p.traverser, selector, p.caster.AsFloat32)
}

func (p *Picker) Float32Slice(selector string) ([]float32, error) {
	return Selector(p.data, p.notation, p.traverser, selector, p.caster.AsFloat32Slice)
}

func (p *Picker) Float64(selector string) (float64, error) {
	return Selector(p.data, p.notation, p.traverser, selector, p.caster.AsFloat64)
}

func (p *Picker) Float64Slice(selector string) ([]float64, error) {
	return Selector(p.data, p.notation, p.traverser, selector, p.caster.AsFloat64Slice)
}

func (p *Picker) Int(selector string) (int, error) {
	return Selector(p.data, p.notation, p.traverser, selector, p.caster.AsInt)
}

func (p *Picker) IntSlice(selector string) ([]int, error) {
	return Selector(p.data, p.notation, p.traverser, selector, p.caster.AsIntSlice)
}

func (p *Picker) Int8(selector string) (int8, error) {
	return Selector(p.data, p.notation, p.traverser, selector, p.caster.AsInt8)
}

func (p *Picker) Int8Slice(selector string) ([]int8, error) {
	return Selector(p.data, p.notation, p.traverser, selector, p.caster.AsInt8Slice)
}

func (p *Picker) Int16(selector string) (int16, error) {
	return Selector(p.data, p.notation, p.traverser, selector, p.caster.AsInt16)
}

func (p *Picker) Int16Slice(selector string) ([]int16, error) {
	return Selector(p.data, p.notation, p.traverser, selector, p.caster.AsInt16Slice)
}

func (p *Picker) Int32(selector string) (int32, error) {
	return Selector(p.data, p.notation, p.traverser, selector, p.caster.AsInt32)
}

func (p *Picker) Int32Slice(selector string) ([]int32, error) {
	return Selector(p.data, p.notation, p.traverser, selector, p.caster.AsInt32Slice)
}

func (p *Picker) Int64(selector string) (int64, error) {
	return Selector(p.data, p.notation, p.traverser, selector, p.caster.AsInt64)
}

func (p *Picker) Int64Slice(selector string) ([]int64, error) {
	return Selector(p.data, p.notation, p.traverser, selector, p.caster.AsInt64Slice)
}

func (p *Picker) Uint(selector string) (uint, error) {
	return Selector(p.data, p.notation, p.traverser, selector, p.caster.AsUint)
}

func (p *Picker) UintSlice(selector string) ([]uint, error) {
	return Selector(p.data, p.notation, p.traverser, selector, p.caster.AsUintSlice)
}

func (p *Picker) Uint8(selector string) (uint8, error) {
	return Selector(p.data, p.notation, p.traverser, selector, p.caster.AsUint8)
}

func (p *Picker) Uint8Slice(selector string) ([]uint8, error) {
	return Selector(p.data, p.notation, p.traverser, selector, p.caster.AsUint8Slice)
}

func (p *Picker) Uint16(selector string) (uint16, error) {
	return Selector(p.data, p.notation, p.traverser, selector, p.caster.AsUint16)
}

func (p *Picker) Uint16Slice(selector string) ([]uint16, error) {
	return Selector(p.data, p.notation, p.traverser, selector, p.caster.AsUint16Slice)
}

func (p *Picker) Uint32(selector string) (uint32, error) {
	return Selector(p.data, p.notation, p.traverser, selector, p.caster.AsUint32)
}

func (p *Picker) Uint32Slice(selector string) ([]uint32, error) {
	return Selector(p.data, p.notation, p.traverser, selector, p.caster.AsUint32Slice)
}

func (p *Picker) Uint64(selector string) (uint64, error) {
	return Selector(p.data, p.notation, p.traverser, selector, p.caster.AsUint64)
}

func (p *Picker) Uint64Slice(selector string) ([]uint64, error) {
	return Selector(p.data, p.notation, p.traverser, selector, p.caster.AsUint64Slice)
}

func (p *Picker) String(selector string) (string, error) {
	return Selector(p.data, p.notation, p.traverser, selector, p.caster.AsString)
}

func (p *Picker) StringSlice(selector string) ([]string, error) {
	return Selector(p.data, p.notation, p.traverser, selector, p.caster.AsStringSlice)
}

//nolint:ireturn
func Map[Output any](p *Picker, selector string, mapFn func(*Picker) (Output, error)) ([]Output, error) {
	item, err := Selector(p.data, p.notation, p.traverser, selector, omitCast)
	if err != nil {
		return nil, err
	}

	return cast.ToSlice(item, func(a any) (Output, error) { return mapFn(Wrap(a)) })
}

//nolint:ireturn
func FlatMap[Output any](p *Picker, selector string, mapFn func(*Picker) ([]Output, error)) ([]Output, error) {
	item, err := Selector(p.data, p.notation, p.traverser, selector, omitCast)
	if err != nil {
		return nil, err
	}

	doubleSlice, err := cast.ToSlice(item, func(a any) ([]Output, error) { return mapFn(Wrap(a)) })
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
func Selector[Output any](data any, notation Notation, traverser Traverser, selector string, castFn func(any) (Output, error)) (Output, error) {
	path, err := notation.Parse(selector)
	if err != nil {
		var d Output
		return d, err
	}

	return Path(data, traverser, path, castFn)
}

//nolint:ireturn
func Path[Output any](data any, traverser Traverser, path []Key, castFn func(any) (Output, error)) (Output, error) {
	item, err := traverser.Retrieve(data, path)
	if err != nil {
		var d Output
		return d, err
	}

	casted, err := castFn(item)
	return casted, err
}

//nolint:ireturn
func SelectorMust[Output any](data any, notation Notation, traverser Traverser, selector string, castFn func(any) (Output, error), onErr ...func(selector string, err error)) Output {
	casted, err := Selector(data, notation, traverser, selector, castFn)
	if err != nil {
		for _, fn := range onErr {
			fn(selector, err)
		}
	}
	return casted
}

//nolint:ireturn
func PathMust[Output any](data any, traverser Traverser, path []Key, castFn func(any) (Output, error), onErr ...func(selector string, err error)) Output {
	casted, err := Path(data, traverser, path, castFn)
	if err != nil {
		selector := DotNotation{}.Format(path...)
		for _, fn := range onErr {
			fn(selector, err)
		}
	}
	return casted
}

func omitCast(a any) (any, error) { return a, nil }
