package pick

import (
	"encoding/json"

	"github.com/moukoublen/pick/cast"
)

type SelectorFormat interface {
	Parse(s string) ([]SelectorKey, error)
	Format(s []SelectorKey) string
}

type Traverser interface {
	Get(m any, selector []SelectorKey) (any, bool, error)
}

func Wrap(data any) *Picker {
	caster := cast.NewCaster()
	formatter := DefaultSelectorFormat{}
	return NewPicker(data, NewDefaultTraverser(caster), caster, formatter)
}

func WrapJSON(js []byte) (*Picker, error) {
	m := map[string]any{}
	err := json.Unmarshal(js, &m)
	if err != nil {
		return nil, err
	}

	caster := cast.NewCaster()
	formatter := DefaultSelectorFormat{}
	return NewPicker(m, NewDefaultTraverser(caster), caster, formatter), nil
}

type Picker struct {
	inner          any
	traverser      Traverser
	caster         Caster
	selectorFormat SelectorFormat
}

func NewPicker(inner any, t Traverser, c Caster, selectorFormat SelectorFormat) *Picker {
	return &Picker{
		inner:          inner,
		traverser:      t,
		caster:         c,
		selectorFormat: selectorFormat,
	}
}

func (o *Picker) Bool(selector string) (bool, bool, error) {
	return get(o, o.caster.AsBool, selector)
}

func (o *Picker) BoolSlice(selector string) ([]bool, bool, error) {
	return get(o, o.caster.AsBoolSlice, selector)
}

func (o *Picker) Byte(selector string) (byte, bool, error) {
	return get(o, o.caster.AsByte, selector)
}

func (o *Picker) ByteSlice(selector string) ([]byte, bool, error) {
	return get(o, o.caster.AsByteSlice, selector)
}

func (o *Picker) Float32(selector string) (float32, bool, error) {
	return get(o, o.caster.AsFloat32, selector)
}

func (o *Picker) Float32Slice(selector string) ([]float32, bool, error) {
	return get(o, o.caster.AsFloat32Slice, selector)
}

func (o *Picker) Float64(selector string) (float64, bool, error) {
	return get(o, o.caster.AsFloat64, selector)
}

func (o *Picker) Float64Slice(selector string) ([]float64, bool, error) {
	return get(o, o.caster.AsFloat64Slice, selector)
}

func (o *Picker) Int(selector string) (int, bool, error) {
	return get(o, o.caster.AsInt, selector)
}

func (o *Picker) IntSlice(selector string) ([]int, bool, error) {
	return get(o, o.caster.AsIntSlice, selector)
}

func (o *Picker) Int8(selector string) (int8, bool, error) {
	return get(o, o.caster.AsInt8, selector)
}

func (o *Picker) Int8Slice(selector string) ([]int8, bool, error) {
	return get(o, o.caster.AsInt8Slice, selector)
}

func (o *Picker) Int16(selector string) (int16, bool, error) {
	return get(o, o.caster.AsInt16, selector)
}

func (o *Picker) Int16Slice(selector string) ([]int16, bool, error) {
	return get(o, o.caster.AsInt16Slice, selector)
}

func (o *Picker) Int32(selector string) (int32, bool, error) {
	return get(o, o.caster.AsInt32, selector)
}

func (o *Picker) Int32Slice(selector string) ([]int32, bool, error) {
	return get(o, o.caster.AsInt32Slice, selector)
}

func (o *Picker) Int64(selector string) (int64, bool, error) {
	return get(o, o.caster.AsInt64, selector)
}

func (o *Picker) Int64Slice(selector string) ([]int64, bool, error) {
	return get(o, o.caster.AsInt64Slice, selector)
}

func (o *Picker) Uint(selector string) (uint, bool, error) {
	return get(o, o.caster.AsUint, selector)
}

func (o *Picker) UintSlice(selector string) ([]uint, bool, error) {
	return get(o, o.caster.AsUintSlice, selector)
}

func (o *Picker) Uint8(selector string) (uint8, bool, error) {
	return get(o, o.caster.AsUint8, selector)
}

func (o *Picker) Uint8Slice(selector string) ([]uint8, bool, error) {
	return get(o, o.caster.AsUint8Slice, selector)
}

func (o *Picker) Uint16(selector string) (uint16, bool, error) {
	return get(o, o.caster.AsUint16, selector)
}

func (o *Picker) Uint16Slice(selector string) ([]uint16, bool, error) {
	return get(o, o.caster.AsUint16Slice, selector)
}

func (o *Picker) Uint32(selector string) (uint32, bool, error) {
	return get(o, o.caster.AsUint32, selector)
}

func (o *Picker) Uint32Slice(selector string) ([]uint32, bool, error) {
	return get(o, o.caster.AsUint32Slice, selector)
}

func (o *Picker) Uint64(selector string) (uint64, bool, error) {
	return get(o, o.caster.AsUint64, selector)
}

func (o *Picker) Uint64Slice(selector string) ([]uint64, bool, error) {
	return get(o, o.caster.AsUint64Slice, selector)
}

func (o *Picker) String(selector string) (string, bool, error) {
	return get(o, o.caster.AsString, selector)
}

func (o *Picker) StringSlice(selector string) ([]string, bool, error) {
	return get(o, o.caster.AsStringSlice, selector)
}

// get is a necessary "cheat" because we cannot have a single generic receiver function.
//
//nolint:ireturn
func get[Output any](o *Picker, selectedCastFn func(any) (Output, error), selector string) (Output, bool, error) {
	s, err := o.selectorFormat.Parse(selector)
	if err != nil {
		var d Output
		return d, false, err
	}

	va, found, err := o.traverser.Get(o.inner, s)
	if err != nil || !found {
		var d Output
		return d, found, err
	}

	casted, err := selectedCastFn(va)
	return casted, true, err
}
