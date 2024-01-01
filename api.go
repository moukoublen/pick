package pick

import (
	"time"

	"github.com/moukoublen/pick/cast"
)

// Default Selector API (embedded into Picker)

func (p *Picker) Bool(selector string) (bool, error) {
	return pickSelector(p, selector, p.Caster.AsBool)
}

func (p *Picker) BoolSlice(selector string) ([]bool, error) {
	return pickSelector(p, selector, p.Caster.AsBoolSlice)
}

func (p *Picker) Byte(selector string) (byte, error) {
	return pickSelector(p, selector, p.Caster.AsByte)
}

func (p *Picker) ByteSlice(selector string) ([]byte, error) {
	return pickSelector(p, selector, p.Caster.AsByteSlice)
}

func (p *Picker) Float32(selector string) (float32, error) {
	return pickSelector(p, selector, p.Caster.AsFloat32)
}

func (p *Picker) Float32Slice(selector string) ([]float32, error) {
	return pickSelector(p, selector, p.Caster.AsFloat32Slice)
}

func (p *Picker) Float64(selector string) (float64, error) {
	return pickSelector(p, selector, p.Caster.AsFloat64)
}

func (p *Picker) Float64Slice(selector string) ([]float64, error) {
	return pickSelector(p, selector, p.Caster.AsFloat64Slice)
}

func (p *Picker) Int(selector string) (int, error) {
	return pickSelector(p, selector, p.Caster.AsInt)
}

func (p *Picker) IntSlice(selector string) ([]int, error) {
	return pickSelector(p, selector, p.Caster.AsIntSlice)
}

func (p *Picker) Int8(selector string) (int8, error) {
	return pickSelector(p, selector, p.Caster.AsInt8)
}

func (p *Picker) Int8Slice(selector string) ([]int8, error) {
	return pickSelector(p, selector, p.Caster.AsInt8Slice)
}

func (p *Picker) Int16(selector string) (int16, error) {
	return pickSelector(p, selector, p.Caster.AsInt16)
}

func (p *Picker) Int16Slice(selector string) ([]int16, error) {
	return pickSelector(p, selector, p.Caster.AsInt16Slice)
}

func (p *Picker) Int32(selector string) (int32, error) {
	return pickSelector(p, selector, p.Caster.AsInt32)
}

func (p *Picker) Int32Slice(selector string) ([]int32, error) {
	return pickSelector(p, selector, p.Caster.AsInt32Slice)
}

func (p *Picker) Int64(selector string) (int64, error) {
	return pickSelector(p, selector, p.Caster.AsInt64)
}

func (p *Picker) Int64Slice(selector string) ([]int64, error) {
	return pickSelector(p, selector, p.Caster.AsInt64Slice)
}

func (p *Picker) Uint(selector string) (uint, error) {
	return pickSelector(p, selector, p.Caster.AsUint)
}

func (p *Picker) UintSlice(selector string) ([]uint, error) {
	return pickSelector(p, selector, p.Caster.AsUintSlice)
}

func (p *Picker) Uint8(selector string) (uint8, error) {
	return pickSelector(p, selector, p.Caster.AsUint8)
}

func (p *Picker) Uint8Slice(selector string) ([]uint8, error) {
	return pickSelector(p, selector, p.Caster.AsUint8Slice)
}

func (p *Picker) Uint16(selector string) (uint16, error) {
	return pickSelector(p, selector, p.Caster.AsUint16)
}

func (p *Picker) Uint16Slice(selector string) ([]uint16, error) {
	return pickSelector(p, selector, p.Caster.AsUint16Slice)
}

func (p *Picker) Uint32(selector string) (uint32, error) {
	return pickSelector(p, selector, p.Caster.AsUint32)
}

func (p *Picker) Uint32Slice(selector string) ([]uint32, error) {
	return pickSelector(p, selector, p.Caster.AsUint32Slice)
}

func (p *Picker) Uint64(selector string) (uint64, error) {
	return pickSelector(p, selector, p.Caster.AsUint64)
}

func (p *Picker) Uint64Slice(selector string) ([]uint64, error) {
	return pickSelector(p, selector, p.Caster.AsUint64Slice)
}

func (p *Picker) String(selector string) (string, error) {
	return pickSelector(p, selector, p.Caster.AsString)
}

func (p *Picker) StringSlice(selector string) ([]string, error) {
	return pickSelector(p, selector, p.Caster.AsStringSlice)
}

func (p *Picker) Time(selector string) (time.Time, error) {
	return pickSelector(p, selector, p.Caster.AsTime)
}

func (p *Picker) TimeWithConfig(config cast.TimeCastConfig, selector string) (time.Time, error) {
	return pickSelector(p, selector, func(input any) (time.Time, error) {
		return p.Caster.AsTimeWithConfig(config, input)
	})
}

func (p *Picker) TimeSlice(selector string) ([]time.Time, error) {
	return pickSelector(p, selector, p.Caster.AsTimeSlice)
}

func (p *Picker) TimeSliceWithConfig(config cast.TimeCastConfig, selector string) ([]time.Time, error) {
	return pickSelector(p, selector, func(input any) ([]time.Time, error) {
		return p.Caster.AsTimeSliceWithConfig(config, input)
	})
}

func (p *Picker) Duration(selector string) (time.Duration, error) {
	return pickSelector(p, selector, p.Caster.AsDuration)
}

func (p *Picker) DurationWithConfig(config cast.DurationCastConfig, selector string) (time.Duration, error) {
	return pickSelector(p, selector, func(input any) (time.Duration, error) {
		return p.Caster.AsDurationWithConfig(config, input)
	})
}

func (p *Picker) DurationSlice(selector string) ([]time.Duration, error) {
	return pickSelector(p, selector, p.Caster.AsDurationSlice)
}

func (p *Picker) DurationSliceWithConfig(config cast.DurationCastConfig, selector string) ([]time.Duration, error) {
	return pickSelector(p, selector, func(input any) ([]time.Duration, error) {
		return p.Caster.AsDurationSliceWithConfig(config, input)
	})
}

// Selector Must API

type SelectorMustAPI struct {
	*Picker
	onErr []func(selector string, err error)
}

func (a SelectorMustAPI) Bool(selector string) bool {
	return pickSelectorMust(a.Picker, selector, a.Caster.AsBool, a.onErr...)
}

func (a SelectorMustAPI) BoolSlice(selector string) []bool {
	return pickSelectorMust(a.Picker, selector, a.Caster.AsBoolSlice, a.onErr...)
}

func (a SelectorMustAPI) Byte(selector string) byte {
	return pickSelectorMust(a.Picker, selector, a.Caster.AsByte, a.onErr...)
}

func (a SelectorMustAPI) ByteSlice(selector string) []byte {
	return pickSelectorMust(a.Picker, selector, a.Caster.AsByteSlice, a.onErr...)
}

func (a SelectorMustAPI) Float32(selector string) float32 {
	return pickSelectorMust(a.Picker, selector, a.Caster.AsFloat32, a.onErr...)
}

func (a SelectorMustAPI) Float32Slice(selector string) []float32 {
	return pickSelectorMust(a.Picker, selector, a.Caster.AsFloat32Slice, a.onErr...)
}

func (a SelectorMustAPI) Float64(selector string) float64 {
	return pickSelectorMust(a.Picker, selector, a.Caster.AsFloat64, a.onErr...)
}

func (a SelectorMustAPI) Float64Slice(selector string) []float64 {
	return pickSelectorMust(a.Picker, selector, a.Caster.AsFloat64Slice, a.onErr...)
}

func (a SelectorMustAPI) Int(selector string) int {
	return pickSelectorMust(a.Picker, selector, a.Caster.AsInt, a.onErr...)
}

func (a SelectorMustAPI) IntSlice(selector string) []int {
	return pickSelectorMust(a.Picker, selector, a.Caster.AsIntSlice, a.onErr...)
}

func (a SelectorMustAPI) Int8(selector string) int8 {
	return pickSelectorMust(a.Picker, selector, a.Caster.AsInt8, a.onErr...)
}

func (a SelectorMustAPI) Int8Slice(selector string) []int8 {
	return pickSelectorMust(a.Picker, selector, a.Caster.AsInt8Slice, a.onErr...)
}

func (a SelectorMustAPI) Int16(selector string) int16 {
	return pickSelectorMust(a.Picker, selector, a.Caster.AsInt16, a.onErr...)
}

func (a SelectorMustAPI) Int16Slice(selector string) []int16 {
	return pickSelectorMust(a.Picker, selector, a.Caster.AsInt16Slice, a.onErr...)
}

func (a SelectorMustAPI) Int32(selector string) int32 {
	return pickSelectorMust(a.Picker, selector, a.Caster.AsInt32, a.onErr...)
}

func (a SelectorMustAPI) Int32Slice(selector string) []int32 {
	return pickSelectorMust(a.Picker, selector, a.Caster.AsInt32Slice, a.onErr...)
}

func (a SelectorMustAPI) Int64(selector string) int64 {
	return pickSelectorMust(a.Picker, selector, a.Caster.AsInt64, a.onErr...)
}

func (a SelectorMustAPI) Int64Slice(selector string) []int64 {
	return pickSelectorMust(a.Picker, selector, a.Caster.AsInt64Slice, a.onErr...)
}

func (a SelectorMustAPI) Uint(selector string) uint {
	return pickSelectorMust(a.Picker, selector, a.Caster.AsUint, a.onErr...)
}

func (a SelectorMustAPI) UintSlice(selector string) []uint {
	return pickSelectorMust(a.Picker, selector, a.Caster.AsUintSlice, a.onErr...)
}

func (a SelectorMustAPI) Uint8(selector string) uint8 {
	return pickSelectorMust(a.Picker, selector, a.Caster.AsUint8, a.onErr...)
}

func (a SelectorMustAPI) Uint8Slice(selector string) []uint8 {
	return pickSelectorMust(a.Picker, selector, a.Caster.AsUint8Slice, a.onErr...)
}

func (a SelectorMustAPI) Uint16(selector string) uint16 {
	return pickSelectorMust(a.Picker, selector, a.Caster.AsUint16, a.onErr...)
}

func (a SelectorMustAPI) Uint16Slice(selector string) []uint16 {
	return pickSelectorMust(a.Picker, selector, a.Caster.AsUint16Slice, a.onErr...)
}

func (a SelectorMustAPI) Uint32(selector string) uint32 {
	return pickSelectorMust(a.Picker, selector, a.Caster.AsUint32, a.onErr...)
}

func (a SelectorMustAPI) Uint32Slice(selector string) []uint32 {
	return pickSelectorMust(a.Picker, selector, a.Caster.AsUint32Slice, a.onErr...)
}

func (a SelectorMustAPI) Uint64(selector string) uint64 {
	return pickSelectorMust(a.Picker, selector, a.Caster.AsUint64, a.onErr...)
}

func (a SelectorMustAPI) Uint64Slice(selector string) []uint64 {
	return pickSelectorMust(a.Picker, selector, a.Caster.AsUint64Slice, a.onErr...)
}

func (a SelectorMustAPI) String(selector string) string {
	return pickSelectorMust(a.Picker, selector, a.Caster.AsString, a.onErr...)
}

func (a SelectorMustAPI) StringSlice(selector string) []string {
	return pickSelectorMust(a.Picker, selector, a.Caster.AsStringSlice, a.onErr...)
}

func (a SelectorMustAPI) Time(selector string) time.Time {
	return pickSelectorMust(a.Picker, selector, a.Caster.AsTime, a.onErr...)
}

func (a SelectorMustAPI) TimeWithConfig(config cast.TimeCastConfig, selector string) time.Time {
	return pickSelectorMust(a.Picker, selector, func(input any) (time.Time, error) {
		return a.Caster.AsTimeWithConfig(config, input)
	}, a.onErr...)
}

func (a SelectorMustAPI) TimeSlice(selector string) []time.Time {
	return pickSelectorMust(a.Picker, selector, a.Caster.AsTimeSlice, a.onErr...)
}

func (a SelectorMustAPI) TimeSliceWithConfig(config cast.TimeCastConfig, selector string) []time.Time {
	return pickSelectorMust(a.Picker, selector, func(input any) ([]time.Time, error) {
		return a.Caster.AsTimeSliceWithConfig(config, input)
	}, a.onErr...)
}

func (a SelectorMustAPI) Duration(selector string) time.Duration {
	return pickSelectorMust(a.Picker, selector, a.Caster.AsDuration, a.onErr...)
}

func (a SelectorMustAPI) DurationWithConfig(config cast.DurationCastConfig, selector string) time.Duration {
	return pickSelectorMust(a.Picker, selector, func(input any) (time.Duration, error) {
		return a.Caster.AsDurationWithConfig(config, input)
	}, a.onErr...)
}

func (a SelectorMustAPI) DurationSlice(selector string) []time.Duration {
	return pickSelectorMust(a.Picker, selector, a.Caster.AsDurationSlice, a.onErr...)
}

func (a SelectorMustAPI) DurationSliceWithConfig(config cast.DurationCastConfig, selector string) []time.Duration {
	return pickSelectorMust(a.Picker, selector, func(input any) ([]time.Duration, error) {
		return a.Caster.AsDurationSliceWithConfig(config, input)
	}, a.onErr...)
}
