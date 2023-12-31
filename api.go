package pick

import (
	"time"

	"github.com/moukoublen/pick/cast"
)

// Default Selector API (embedded into Picker)

func (p *Picker) Bool(selector string) (bool, error) {
	return Selector(p, selector, p.caster.AsBool)
}

func (p *Picker) BoolSlice(selector string) ([]bool, error) {
	return Selector(p, selector, p.caster.AsBoolSlice)
}

func (p *Picker) Byte(selector string) (byte, error) {
	return Selector(p, selector, p.caster.AsByte)
}

func (p *Picker) ByteSlice(selector string) ([]byte, error) {
	return Selector(p, selector, p.caster.AsByteSlice)
}

func (p *Picker) Float32(selector string) (float32, error) {
	return Selector(p, selector, p.caster.AsFloat32)
}

func (p *Picker) Float32Slice(selector string) ([]float32, error) {
	return Selector(p, selector, p.caster.AsFloat32Slice)
}

func (p *Picker) Float64(selector string) (float64, error) {
	return Selector(p, selector, p.caster.AsFloat64)
}

func (p *Picker) Float64Slice(selector string) ([]float64, error) {
	return Selector(p, selector, p.caster.AsFloat64Slice)
}

func (p *Picker) Int(selector string) (int, error) {
	return Selector(p, selector, p.caster.AsInt)
}

func (p *Picker) IntSlice(selector string) ([]int, error) {
	return Selector(p, selector, p.caster.AsIntSlice)
}

func (p *Picker) Int8(selector string) (int8, error) {
	return Selector(p, selector, p.caster.AsInt8)
}

func (p *Picker) Int8Slice(selector string) ([]int8, error) {
	return Selector(p, selector, p.caster.AsInt8Slice)
}

func (p *Picker) Int16(selector string) (int16, error) {
	return Selector(p, selector, p.caster.AsInt16)
}

func (p *Picker) Int16Slice(selector string) ([]int16, error) {
	return Selector(p, selector, p.caster.AsInt16Slice)
}

func (p *Picker) Int32(selector string) (int32, error) {
	return Selector(p, selector, p.caster.AsInt32)
}

func (p *Picker) Int32Slice(selector string) ([]int32, error) {
	return Selector(p, selector, p.caster.AsInt32Slice)
}

func (p *Picker) Int64(selector string) (int64, error) {
	return Selector(p, selector, p.caster.AsInt64)
}

func (p *Picker) Int64Slice(selector string) ([]int64, error) {
	return Selector(p, selector, p.caster.AsInt64Slice)
}

func (p *Picker) Uint(selector string) (uint, error) {
	return Selector(p, selector, p.caster.AsUint)
}

func (p *Picker) UintSlice(selector string) ([]uint, error) {
	return Selector(p, selector, p.caster.AsUintSlice)
}

func (p *Picker) Uint8(selector string) (uint8, error) {
	return Selector(p, selector, p.caster.AsUint8)
}

func (p *Picker) Uint8Slice(selector string) ([]uint8, error) {
	return Selector(p, selector, p.caster.AsUint8Slice)
}

func (p *Picker) Uint16(selector string) (uint16, error) {
	return Selector(p, selector, p.caster.AsUint16)
}

func (p *Picker) Uint16Slice(selector string) ([]uint16, error) {
	return Selector(p, selector, p.caster.AsUint16Slice)
}

func (p *Picker) Uint32(selector string) (uint32, error) {
	return Selector(p, selector, p.caster.AsUint32)
}

func (p *Picker) Uint32Slice(selector string) ([]uint32, error) {
	return Selector(p, selector, p.caster.AsUint32Slice)
}

func (p *Picker) Uint64(selector string) (uint64, error) {
	return Selector(p, selector, p.caster.AsUint64)
}

func (p *Picker) Uint64Slice(selector string) ([]uint64, error) {
	return Selector(p, selector, p.caster.AsUint64Slice)
}

func (p *Picker) String(selector string) (string, error) {
	return Selector(p, selector, p.caster.AsString)
}

func (p *Picker) StringSlice(selector string) ([]string, error) {
	return Selector(p, selector, p.caster.AsStringSlice)
}

func (p *Picker) Time(selector string) (time.Time, error) {
	return Selector(p, selector, p.caster.AsTime)
}

func (p *Picker) TimeWithConfig(config cast.TimeCastConfig, selector string) (time.Time, error) {
	return Selector(p, selector, func(input any) (time.Time, error) {
		return p.caster.AsTimeWithConfig(config, input)
	})
}

func (p *Picker) TimeSlice(selector string) ([]time.Time, error) {
	return Selector(p, selector, p.caster.AsTimeSlice)
}

func (p *Picker) TimeSliceWithConfig(config cast.TimeCastConfig, selector string) ([]time.Time, error) {
	return Selector(p, selector, func(input any) ([]time.Time, error) {
		return p.caster.AsTimeSliceWithConfig(config, input)
	})
}

func (p *Picker) Duration(selector string) (time.Duration, error) {
	return Selector(p, selector, p.caster.AsDuration)
}

func (p *Picker) DurationWithConfig(config cast.DurationCastConfig, selector string) (time.Duration, error) {
	return Selector(p, selector, func(input any) (time.Duration, error) {
		return p.caster.AsDurationWithConfig(config, input)
	})
}

func (p *Picker) DurationSlice(selector string) ([]time.Duration, error) {
	return Selector(p, selector, p.caster.AsDurationSlice)
}

func (p *Picker) DurationSliceWithConfig(config cast.DurationCastConfig, selector string) ([]time.Duration, error) {
	return Selector(p, selector, func(input any) ([]time.Duration, error) {
		return p.caster.AsDurationSliceWithConfig(config, input)
	})
}

// Selector Must API

type SelectorMustAPI struct {
	*Picker
	onErr []func(selector string, err error)
}

func (a SelectorMustAPI) Bool(selector string) bool {
	return SelectorMust(a.Picker, selector, a.caster.AsBool, a.onErr...)
}

func (a SelectorMustAPI) BoolSlice(selector string) []bool {
	return SelectorMust(a.Picker, selector, a.caster.AsBoolSlice, a.onErr...)
}

func (a SelectorMustAPI) Byte(selector string) byte {
	return SelectorMust(a.Picker, selector, a.caster.AsByte, a.onErr...)
}

func (a SelectorMustAPI) ByteSlice(selector string) []byte {
	return SelectorMust(a.Picker, selector, a.caster.AsByteSlice, a.onErr...)
}

func (a SelectorMustAPI) Float32(selector string) float32 {
	return SelectorMust(a.Picker, selector, a.caster.AsFloat32, a.onErr...)
}

func (a SelectorMustAPI) Float32Slice(selector string) []float32 {
	return SelectorMust(a.Picker, selector, a.caster.AsFloat32Slice, a.onErr...)
}

func (a SelectorMustAPI) Float64(selector string) float64 {
	return SelectorMust(a.Picker, selector, a.caster.AsFloat64, a.onErr...)
}

func (a SelectorMustAPI) Float64Slice(selector string) []float64 {
	return SelectorMust(a.Picker, selector, a.caster.AsFloat64Slice, a.onErr...)
}

func (a SelectorMustAPI) Int(selector string) int {
	return SelectorMust(a.Picker, selector, a.caster.AsInt, a.onErr...)
}

func (a SelectorMustAPI) IntSlice(selector string) []int {
	return SelectorMust(a.Picker, selector, a.caster.AsIntSlice, a.onErr...)
}

func (a SelectorMustAPI) Int8(selector string) int8 {
	return SelectorMust(a.Picker, selector, a.caster.AsInt8, a.onErr...)
}

func (a SelectorMustAPI) Int8Slice(selector string) []int8 {
	return SelectorMust(a.Picker, selector, a.caster.AsInt8Slice, a.onErr...)
}

func (a SelectorMustAPI) Int16(selector string) int16 {
	return SelectorMust(a.Picker, selector, a.caster.AsInt16, a.onErr...)
}

func (a SelectorMustAPI) Int16Slice(selector string) []int16 {
	return SelectorMust(a.Picker, selector, a.caster.AsInt16Slice, a.onErr...)
}

func (a SelectorMustAPI) Int32(selector string) int32 {
	return SelectorMust(a.Picker, selector, a.caster.AsInt32, a.onErr...)
}

func (a SelectorMustAPI) Int32Slice(selector string) []int32 {
	return SelectorMust(a.Picker, selector, a.caster.AsInt32Slice, a.onErr...)
}

func (a SelectorMustAPI) Int64(selector string) int64 {
	return SelectorMust(a.Picker, selector, a.caster.AsInt64, a.onErr...)
}

func (a SelectorMustAPI) Int64Slice(selector string) []int64 {
	return SelectorMust(a.Picker, selector, a.caster.AsInt64Slice, a.onErr...)
}

func (a SelectorMustAPI) Uint(selector string) uint {
	return SelectorMust(a.Picker, selector, a.caster.AsUint, a.onErr...)
}

func (a SelectorMustAPI) UintSlice(selector string) []uint {
	return SelectorMust(a.Picker, selector, a.caster.AsUintSlice, a.onErr...)
}

func (a SelectorMustAPI) Uint8(selector string) uint8 {
	return SelectorMust(a.Picker, selector, a.caster.AsUint8, a.onErr...)
}

func (a SelectorMustAPI) Uint8Slice(selector string) []uint8 {
	return SelectorMust(a.Picker, selector, a.caster.AsUint8Slice, a.onErr...)
}

func (a SelectorMustAPI) Uint16(selector string) uint16 {
	return SelectorMust(a.Picker, selector, a.caster.AsUint16, a.onErr...)
}

func (a SelectorMustAPI) Uint16Slice(selector string) []uint16 {
	return SelectorMust(a.Picker, selector, a.caster.AsUint16Slice, a.onErr...)
}

func (a SelectorMustAPI) Uint32(selector string) uint32 {
	return SelectorMust(a.Picker, selector, a.caster.AsUint32, a.onErr...)
}

func (a SelectorMustAPI) Uint32Slice(selector string) []uint32 {
	return SelectorMust(a.Picker, selector, a.caster.AsUint32Slice, a.onErr...)
}

func (a SelectorMustAPI) Uint64(selector string) uint64 {
	return SelectorMust(a.Picker, selector, a.caster.AsUint64, a.onErr...)
}

func (a SelectorMustAPI) Uint64Slice(selector string) []uint64 {
	return SelectorMust(a.Picker, selector, a.caster.AsUint64Slice, a.onErr...)
}

func (a SelectorMustAPI) String(selector string) string {
	return SelectorMust(a.Picker, selector, a.caster.AsString, a.onErr...)
}

func (a SelectorMustAPI) StringSlice(selector string) []string {
	return SelectorMust(a.Picker, selector, a.caster.AsStringSlice, a.onErr...)
}

func (a SelectorMustAPI) Time(selector string) time.Time {
	return SelectorMust(a.Picker, selector, a.caster.AsTime, a.onErr...)
}

func (a SelectorMustAPI) TimeWithConfig(config cast.TimeCastConfig, selector string) time.Time {
	return SelectorMust(a.Picker, selector, func(input any) (time.Time, error) {
		return a.caster.AsTimeWithConfig(config, input)
	}, a.onErr...)
}

func (a SelectorMustAPI) TimeSlice(selector string) []time.Time {
	return SelectorMust(a.Picker, selector, a.caster.AsTimeSlice, a.onErr...)
}

func (a SelectorMustAPI) TimeSliceWithConfig(config cast.TimeCastConfig, selector string) []time.Time {
	return SelectorMust(a.Picker, selector, func(input any) ([]time.Time, error) {
		return a.caster.AsTimeSliceWithConfig(config, input)
	}, a.onErr...)
}

func (a SelectorMustAPI) Duration(selector string) time.Duration {
	return SelectorMust(a.Picker, selector, a.caster.AsDuration, a.onErr...)
}

func (a SelectorMustAPI) DurationWithConfig(config cast.DurationCastConfig, selector string) time.Duration {
	return SelectorMust(a.Picker, selector, func(input any) (time.Duration, error) {
		return a.caster.AsDurationWithConfig(config, input)
	}, a.onErr...)
}

func (a SelectorMustAPI) DurationSlice(selector string) []time.Duration {
	return SelectorMust(a.Picker, selector, a.caster.AsDurationSlice, a.onErr...)
}

func (a SelectorMustAPI) DurationSliceWithConfig(config cast.DurationCastConfig, selector string) []time.Duration {
	return SelectorMust(a.Picker, selector, func(input any) ([]time.Duration, error) {
		return a.caster.AsDurationSliceWithConfig(config, input)
	}, a.onErr...)
}
