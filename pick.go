package pick

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/moukoublen/pick/cast"
)

type SelectorFormat interface {
	Parse(s string) ([]SelectorKey, error)
	Format(s []SelectorKey) string
}

type Traverser interface {
	Get(m any, selector []SelectorKey) (any, error)
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
	return NewPicker(data, NewDefaultTraverser(caster), caster, DefaultSelectorFormat{})
}

type Picker struct {
	inner          any
	traverser      Traverser
	caster         Caster
	selectorFormat SelectorFormat
}

func NewPicker(data any, t Traverser, c Caster, s SelectorFormat) *Picker {
	return &Picker{
		inner:          data,
		traverser:      t,
		caster:         c,
		selectorFormat: s,
	}
}

func (p *Picker) Bool(selector string) (bool, error) {
	return Pick(p, selector, p.caster.AsBool)
}

func (p *Picker) BoolSlice(selector string) ([]bool, error) {
	return Pick(p, selector, p.caster.AsBoolSlice)
}

func (p *Picker) Byte(selector string) (byte, error) {
	return Pick(p, selector, p.caster.AsByte)
}

func (p *Picker) ByteSlice(selector string) ([]byte, error) {
	return Pick(p, selector, p.caster.AsByteSlice)
}

func (p *Picker) Float32(selector string) (float32, error) {
	return Pick(p, selector, p.caster.AsFloat32)
}

func (p *Picker) Float32Slice(selector string) ([]float32, error) {
	return Pick(p, selector, p.caster.AsFloat32Slice)
}

func (p *Picker) Float64(selector string) (float64, error) {
	return Pick(p, selector, p.caster.AsFloat64)
}

func (p *Picker) Float64Slice(selector string) ([]float64, error) {
	return Pick(p, selector, p.caster.AsFloat64Slice)
}

func (p *Picker) Int(selector string) (int, error) {
	return Pick(p, selector, p.caster.AsInt)
}

func (p *Picker) IntSlice(selector string) ([]int, error) {
	return Pick(p, selector, p.caster.AsIntSlice)
}

func (p *Picker) Int8(selector string) (int8, error) {
	return Pick(p, selector, p.caster.AsInt8)
}

func (p *Picker) Int8Slice(selector string) ([]int8, error) {
	return Pick(p, selector, p.caster.AsInt8Slice)
}

func (p *Picker) Int16(selector string) (int16, error) {
	return Pick(p, selector, p.caster.AsInt16)
}

func (p *Picker) Int16Slice(selector string) ([]int16, error) {
	return Pick(p, selector, p.caster.AsInt16Slice)
}

func (p *Picker) Int32(selector string) (int32, error) {
	return Pick(p, selector, p.caster.AsInt32)
}

func (p *Picker) Int32Slice(selector string) ([]int32, error) {
	return Pick(p, selector, p.caster.AsInt32Slice)
}

func (p *Picker) Int64(selector string) (int64, error) {
	return Pick(p, selector, p.caster.AsInt64)
}

func (p *Picker) Int64Slice(selector string) ([]int64, error) {
	return Pick(p, selector, p.caster.AsInt64Slice)
}

func (p *Picker) Uint(selector string) (uint, error) {
	return Pick(p, selector, p.caster.AsUint)
}

func (p *Picker) UintSlice(selector string) ([]uint, error) {
	return Pick(p, selector, p.caster.AsUintSlice)
}

func (p *Picker) Uint8(selector string) (uint8, error) {
	return Pick(p, selector, p.caster.AsUint8)
}

func (p *Picker) Uint8Slice(selector string) ([]uint8, error) {
	return Pick(p, selector, p.caster.AsUint8Slice)
}

func (p *Picker) Uint16(selector string) (uint16, error) {
	return Pick(p, selector, p.caster.AsUint16)
}

func (p *Picker) Uint16Slice(selector string) ([]uint16, error) {
	return Pick(p, selector, p.caster.AsUint16Slice)
}

func (p *Picker) Uint32(selector string) (uint32, error) {
	return Pick(p, selector, p.caster.AsUint32)
}

func (p *Picker) Uint32Slice(selector string) ([]uint32, error) {
	return Pick(p, selector, p.caster.AsUint32Slice)
}

func (p *Picker) Uint64(selector string) (uint64, error) {
	return Pick(p, selector, p.caster.AsUint64)
}

func (p *Picker) Uint64Slice(selector string) ([]uint64, error) {
	return Pick(p, selector, p.caster.AsUint64Slice)
}

func (p *Picker) String(selector string) (string, error) {
	return Pick(p, selector, p.caster.AsString)
}

func (p *Picker) StringSlice(selector string) ([]string, error) {
	return Pick(p, selector, p.caster.AsStringSlice)
}

//nolint:ireturn
func Map[Output any](p *Picker, selector string, mapFn func(*Picker) (Output, error)) ([]Output, error) {
	item, err := Pick(p, selector, castPassThrough)
	if err != nil {
		return nil, err
	}

	return cast.ToSlice(item, func(a any) (Output, error) { return mapFn(Wrap(a)) })
}

//nolint:ireturn
func FlatMap[Output any](p *Picker, selector string, mapFn func(*Picker) ([]Output, error)) ([]Output, error) {
	item, err := Pick(p, selector, castPassThrough)
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
func Pick[Output any](p *Picker, selector string, selectedCastFn func(any) (Output, error)) (Output, error) {
	s, err := p.selectorFormat.Parse(selector)
	if err != nil {
		var d Output
		return d, err
	}

	item, err := p.traverser.Get(p.inner, s)
	if err != nil {
		var d Output
		return d, err
	}

	casted, err := selectedCastFn(item)
	return casted, err
}

func castPassThrough(a any) (any, error) { return a, nil }
