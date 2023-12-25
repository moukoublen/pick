package pick

import (
	"time"

	"github.com/moukoublen/pick/cast"
)

// Default Selector API (embedded into Picker)

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

func (p *Picker) Time(selector string) (time.Time, error) {
	return Selector(p.data, p.notation, p.traverser, selector, p.caster.AsTime)
}

func (p *Picker) TimeWithConfig(config cast.TimeCastConfig, selector string) (time.Time, error) {
	return Selector(p.data, p.notation, p.traverser, selector, func(input any) (time.Time, error) {
		return p.caster.AsTimeWithConfig(config, input)
	})
}

func (p *Picker) TimeSlice(selector string) ([]time.Time, error) {
	return Selector(p.data, p.notation, p.traverser, selector, p.caster.AsTimeSlice)
}

func (p *Picker) TimeSliceWithConfig(config cast.TimeCastConfig, selector string) ([]time.Time, error) {
	return Selector(p.data, p.notation, p.traverser, selector, func(input any) ([]time.Time, error) {
		return p.caster.AsTimeSliceWithConfig(config, input)
	})
}

// Selector Must API

type SelectorMustAPI struct {
	*Picker
	onErr []func(selector string, err error)
}

func (a SelectorMustAPI) Bool(selector string) bool {
	return SelectorMust(a.data, a.notation, a.traverser, selector, a.caster.AsBool, a.onErr...)
}

func (a SelectorMustAPI) BoolSlice(selector string) []bool {
	return SelectorMust(a.data, a.notation, a.traverser, selector, a.caster.AsBoolSlice, a.onErr...)
}

func (a SelectorMustAPI) Byte(selector string) byte {
	return SelectorMust(a.data, a.notation, a.traverser, selector, a.caster.AsByte, a.onErr...)
}

func (a SelectorMustAPI) ByteSlice(selector string) []byte {
	return SelectorMust(a.data, a.notation, a.traverser, selector, a.caster.AsByteSlice, a.onErr...)
}

func (a SelectorMustAPI) Float32(selector string) float32 {
	return SelectorMust(a.data, a.notation, a.traverser, selector, a.caster.AsFloat32, a.onErr...)
}

func (a SelectorMustAPI) Float32Slice(selector string) []float32 {
	return SelectorMust(a.data, a.notation, a.traverser, selector, a.caster.AsFloat32Slice, a.onErr...)
}

func (a SelectorMustAPI) Float64(selector string) float64 {
	return SelectorMust(a.data, a.notation, a.traverser, selector, a.caster.AsFloat64, a.onErr...)
}

func (a SelectorMustAPI) Float64Slice(selector string) []float64 {
	return SelectorMust(a.data, a.notation, a.traverser, selector, a.caster.AsFloat64Slice, a.onErr...)
}

func (a SelectorMustAPI) Int(selector string) int {
	return SelectorMust(a.data, a.notation, a.traverser, selector, a.caster.AsInt, a.onErr...)
}

func (a SelectorMustAPI) IntSlice(selector string) []int {
	return SelectorMust(a.data, a.notation, a.traverser, selector, a.caster.AsIntSlice, a.onErr...)
}

func (a SelectorMustAPI) Int8(selector string) int8 {
	return SelectorMust(a.data, a.notation, a.traverser, selector, a.caster.AsInt8, a.onErr...)
}

func (a SelectorMustAPI) Int8Slice(selector string) []int8 {
	return SelectorMust(a.data, a.notation, a.traverser, selector, a.caster.AsInt8Slice, a.onErr...)
}

func (a SelectorMustAPI) Int16(selector string) int16 {
	return SelectorMust(a.data, a.notation, a.traverser, selector, a.caster.AsInt16, a.onErr...)
}

func (a SelectorMustAPI) Int16Slice(selector string) []int16 {
	return SelectorMust(a.data, a.notation, a.traverser, selector, a.caster.AsInt16Slice, a.onErr...)
}

func (a SelectorMustAPI) Int32(selector string) int32 {
	return SelectorMust(a.data, a.notation, a.traverser, selector, a.caster.AsInt32, a.onErr...)
}

func (a SelectorMustAPI) Int32Slice(selector string) []int32 {
	return SelectorMust(a.data, a.notation, a.traverser, selector, a.caster.AsInt32Slice, a.onErr...)
}

func (a SelectorMustAPI) Int64(selector string) int64 {
	return SelectorMust(a.data, a.notation, a.traverser, selector, a.caster.AsInt64, a.onErr...)
}

func (a SelectorMustAPI) Int64Slice(selector string) []int64 {
	return SelectorMust(a.data, a.notation, a.traverser, selector, a.caster.AsInt64Slice, a.onErr...)
}

func (a SelectorMustAPI) Uint(selector string) uint {
	return SelectorMust(a.data, a.notation, a.traverser, selector, a.caster.AsUint, a.onErr...)
}

func (a SelectorMustAPI) UintSlice(selector string) []uint {
	return SelectorMust(a.data, a.notation, a.traverser, selector, a.caster.AsUintSlice, a.onErr...)
}

func (a SelectorMustAPI) Uint8(selector string) uint8 {
	return SelectorMust(a.data, a.notation, a.traverser, selector, a.caster.AsUint8, a.onErr...)
}

func (a SelectorMustAPI) Uint8Slice(selector string) []uint8 {
	return SelectorMust(a.data, a.notation, a.traverser, selector, a.caster.AsUint8Slice, a.onErr...)
}

func (a SelectorMustAPI) Uint16(selector string) uint16 {
	return SelectorMust(a.data, a.notation, a.traverser, selector, a.caster.AsUint16, a.onErr...)
}

func (a SelectorMustAPI) Uint16Slice(selector string) []uint16 {
	return SelectorMust(a.data, a.notation, a.traverser, selector, a.caster.AsUint16Slice, a.onErr...)
}

func (a SelectorMustAPI) Uint32(selector string) uint32 {
	return SelectorMust(a.data, a.notation, a.traverser, selector, a.caster.AsUint32, a.onErr...)
}

func (a SelectorMustAPI) Uint32Slice(selector string) []uint32 {
	return SelectorMust(a.data, a.notation, a.traverser, selector, a.caster.AsUint32Slice, a.onErr...)
}

func (a SelectorMustAPI) Uint64(selector string) uint64 {
	return SelectorMust(a.data, a.notation, a.traverser, selector, a.caster.AsUint64, a.onErr...)
}

func (a SelectorMustAPI) Uint64Slice(selector string) []uint64 {
	return SelectorMust(a.data, a.notation, a.traverser, selector, a.caster.AsUint64Slice, a.onErr...)
}

func (a SelectorMustAPI) String(selector string) string {
	return SelectorMust(a.data, a.notation, a.traverser, selector, a.caster.AsString, a.onErr...)
}

func (a SelectorMustAPI) StringSlice(selector string) []string {
	return SelectorMust(a.data, a.notation, a.traverser, selector, a.caster.AsStringSlice, a.onErr...)
}

func (a *SelectorMustAPI) Time(selector string) time.Time {
	return SelectorMust(a.data, a.notation, a.traverser, selector, a.caster.AsTime, a.onErr...)
}

func (a *SelectorMustAPI) TimeWithConfig(config cast.TimeCastConfig, selector string) time.Time {
	return SelectorMust(a.data, a.notation, a.traverser, selector, func(input any) (time.Time, error) {
		return a.caster.AsTimeWithConfig(config, input)
	}, a.onErr...)
}

func (a *SelectorMustAPI) TimeSlice(selector string) []time.Time {
	return SelectorMust(a.data, a.notation, a.traverser, selector, a.caster.AsTimeSlice, a.onErr...)
}

func (a *SelectorMustAPI) TimeSliceWithConfig(config cast.TimeCastConfig, selector string) []time.Time {
	return SelectorMust(a.data, a.notation, a.traverser, selector, func(input any) ([]time.Time, error) {
		return a.caster.AsTimeSliceWithConfig(config, input)
	}, a.onErr...)
}

// Path API

type PathAPI struct {
	*Picker
}

func (a PathAPI) Bool(path ...Key) (bool, error) {
	return Path(a.data, a.traverser, path, a.caster.AsBool)
}

func (a PathAPI) BoolSlice(path ...Key) ([]bool, error) {
	return Path(a.data, a.traverser, path, a.caster.AsBoolSlice)
}

func (a PathAPI) Byte(path ...Key) (byte, error) {
	return Path(a.data, a.traverser, path, a.caster.AsByte)
}

func (a PathAPI) ByteSlice(path ...Key) ([]byte, error) {
	return Path(a.data, a.traverser, path, a.caster.AsByteSlice)
}

func (a PathAPI) Float32(path ...Key) (float32, error) {
	return Path(a.data, a.traverser, path, a.caster.AsFloat32)
}

func (a PathAPI) Float32Slice(path ...Key) ([]float32, error) {
	return Path(a.data, a.traverser, path, a.caster.AsFloat32Slice)
}

func (a PathAPI) Float64(path ...Key) (float64, error) {
	return Path(a.data, a.traverser, path, a.caster.AsFloat64)
}

func (a PathAPI) Float64Slice(path ...Key) ([]float64, error) {
	return Path(a.data, a.traverser, path, a.caster.AsFloat64Slice)
}

func (a PathAPI) Int(path ...Key) (int, error) {
	return Path(a.data, a.traverser, path, a.caster.AsInt)
}

func (a PathAPI) IntSlice(path ...Key) ([]int, error) {
	return Path(a.data, a.traverser, path, a.caster.AsIntSlice)
}

func (a PathAPI) Int8(path ...Key) (int8, error) {
	return Path(a.data, a.traverser, path, a.caster.AsInt8)
}

func (a PathAPI) Int8Slice(path ...Key) ([]int8, error) {
	return Path(a.data, a.traverser, path, a.caster.AsInt8Slice)
}

func (a PathAPI) Int16(path ...Key) (int16, error) {
	return Path(a.data, a.traverser, path, a.caster.AsInt16)
}

func (a PathAPI) Int16Slice(path ...Key) ([]int16, error) {
	return Path(a.data, a.traverser, path, a.caster.AsInt16Slice)
}

func (a PathAPI) Int32(path ...Key) (int32, error) {
	return Path(a.data, a.traverser, path, a.caster.AsInt32)
}

func (a PathAPI) Int32Slice(path ...Key) ([]int32, error) {
	return Path(a.data, a.traverser, path, a.caster.AsInt32Slice)
}

func (a PathAPI) Int64(path ...Key) (int64, error) {
	return Path(a.data, a.traverser, path, a.caster.AsInt64)
}

func (a PathAPI) Int64Slice(path ...Key) ([]int64, error) {
	return Path(a.data, a.traverser, path, a.caster.AsInt64Slice)
}

func (a PathAPI) Uint(path ...Key) (uint, error) {
	return Path(a.data, a.traverser, path, a.caster.AsUint)
}

func (a PathAPI) UintSlice(path ...Key) ([]uint, error) {
	return Path(a.data, a.traverser, path, a.caster.AsUintSlice)
}

func (a PathAPI) Uint8(path ...Key) (uint8, error) {
	return Path(a.data, a.traverser, path, a.caster.AsUint8)
}

func (a PathAPI) Uint8Slice(path ...Key) ([]uint8, error) {
	return Path(a.data, a.traverser, path, a.caster.AsUint8Slice)
}

func (a PathAPI) Uint16(path ...Key) (uint16, error) {
	return Path(a.data, a.traverser, path, a.caster.AsUint16)
}

func (a PathAPI) Uint16Slice(path ...Key) ([]uint16, error) {
	return Path(a.data, a.traverser, path, a.caster.AsUint16Slice)
}

func (a PathAPI) Uint32(path ...Key) (uint32, error) {
	return Path(a.data, a.traverser, path, a.caster.AsUint32)
}

func (a PathAPI) Uint32Slice(path ...Key) ([]uint32, error) {
	return Path(a.data, a.traverser, path, a.caster.AsUint32Slice)
}

func (a PathAPI) Uint64(path ...Key) (uint64, error) {
	return Path(a.data, a.traverser, path, a.caster.AsUint64)
}

func (a PathAPI) Uint64Slice(path ...Key) ([]uint64, error) {
	return Path(a.data, a.traverser, path, a.caster.AsUint64Slice)
}

func (a PathAPI) String(path ...Key) (string, error) {
	return Path(a.data, a.traverser, path, a.caster.AsString)
}

func (a PathAPI) StringSlice(path ...Key) ([]string, error) {
	return Path(a.data, a.traverser, path, a.caster.AsStringSlice)
}

func (a *PathAPI) Time(path ...Key) (time.Time, error) {
	return Path(a.data, a.traverser, path, a.caster.AsTime)
}

func (a *PathAPI) TimeWithConfig(config cast.TimeCastConfig, path ...Key) (time.Time, error) {
	return Path(a.data, a.traverser, path, func(input any) (time.Time, error) {
		return a.caster.AsTimeWithConfig(config, input)
	})
}

func (a *PathAPI) TimeSlice(path ...Key) ([]time.Time, error) {
	return Path(a.data, a.traverser, path, a.caster.AsTimeSlice)
}

func (a *PathAPI) TimeSliceWithConfig(config cast.TimeCastConfig, path ...Key) ([]time.Time, error) {
	return Path(a.data, a.traverser, path, func(input any) ([]time.Time, error) {
		return a.caster.AsTimeSliceWithConfig(config, input)
	})
}

// Path Must API

type PathMustAPI struct {
	*Picker
	onErr []func(selector string, err error)
}

func (a PathMustAPI) Bool(path ...Key) bool {
	return PathMust(a.data, a.traverser, path, a.caster.AsBool, a.onErr...)
}

func (a PathMustAPI) BoolSlice(path ...Key) []bool {
	return PathMust(a.data, a.traverser, path, a.caster.AsBoolSlice, a.onErr...)
}

func (a PathMustAPI) Byte(path ...Key) byte {
	return PathMust(a.data, a.traverser, path, a.caster.AsByte, a.onErr...)
}

func (a PathMustAPI) ByteSlice(path ...Key) []byte {
	return PathMust(a.data, a.traverser, path, a.caster.AsByteSlice, a.onErr...)
}

func (a PathMustAPI) Float32(path ...Key) float32 {
	return PathMust(a.data, a.traverser, path, a.caster.AsFloat32, a.onErr...)
}

func (a PathMustAPI) Float32Slice(path ...Key) []float32 {
	return PathMust(a.data, a.traverser, path, a.caster.AsFloat32Slice, a.onErr...)
}

func (a PathMustAPI) Float64(path ...Key) float64 {
	return PathMust(a.data, a.traverser, path, a.caster.AsFloat64, a.onErr...)
}

func (a PathMustAPI) Float64Slice(path ...Key) []float64 {
	return PathMust(a.data, a.traverser, path, a.caster.AsFloat64Slice, a.onErr...)
}

func (a PathMustAPI) Int(path ...Key) int {
	return PathMust(a.data, a.traverser, path, a.caster.AsInt, a.onErr...)
}

func (a PathMustAPI) IntSlice(path ...Key) []int {
	return PathMust(a.data, a.traverser, path, a.caster.AsIntSlice, a.onErr...)
}

func (a PathMustAPI) Int8(path ...Key) int8 {
	return PathMust(a.data, a.traverser, path, a.caster.AsInt8, a.onErr...)
}

func (a PathMustAPI) Int8Slice(path ...Key) []int8 {
	return PathMust(a.data, a.traverser, path, a.caster.AsInt8Slice, a.onErr...)
}

func (a PathMustAPI) Int16(path ...Key) int16 {
	return PathMust(a.data, a.traverser, path, a.caster.AsInt16, a.onErr...)
}

func (a PathMustAPI) Int16Slice(path ...Key) []int16 {
	return PathMust(a.data, a.traverser, path, a.caster.AsInt16Slice, a.onErr...)
}

func (a PathMustAPI) Int32(path ...Key) int32 {
	return PathMust(a.data, a.traverser, path, a.caster.AsInt32, a.onErr...)
}

func (a PathMustAPI) Int32Slice(path ...Key) []int32 {
	return PathMust(a.data, a.traverser, path, a.caster.AsInt32Slice, a.onErr...)
}

func (a PathMustAPI) Int64(path ...Key) int64 {
	return PathMust(a.data, a.traverser, path, a.caster.AsInt64, a.onErr...)
}

func (a PathMustAPI) Int64Slice(path ...Key) []int64 {
	return PathMust(a.data, a.traverser, path, a.caster.AsInt64Slice, a.onErr...)
}

func (a PathMustAPI) Uint(path ...Key) uint {
	return PathMust(a.data, a.traverser, path, a.caster.AsUint, a.onErr...)
}

func (a PathMustAPI) UintSlice(path ...Key) []uint {
	return PathMust(a.data, a.traverser, path, a.caster.AsUintSlice, a.onErr...)
}

func (a PathMustAPI) Uint8(path ...Key) uint8 {
	return PathMust(a.data, a.traverser, path, a.caster.AsUint8, a.onErr...)
}

func (a PathMustAPI) Uint8Slice(path ...Key) []uint8 {
	return PathMust(a.data, a.traverser, path, a.caster.AsUint8Slice, a.onErr...)
}

func (a PathMustAPI) Uint16(path ...Key) uint16 {
	return PathMust(a.data, a.traverser, path, a.caster.AsUint16, a.onErr...)
}

func (a PathMustAPI) Uint16Slice(path ...Key) []uint16 {
	return PathMust(a.data, a.traverser, path, a.caster.AsUint16Slice, a.onErr...)
}

func (a PathMustAPI) Uint32(path ...Key) uint32 {
	return PathMust(a.data, a.traverser, path, a.caster.AsUint32, a.onErr...)
}

func (a PathMustAPI) Uint32Slice(path ...Key) []uint32 {
	return PathMust(a.data, a.traverser, path, a.caster.AsUint32Slice, a.onErr...)
}

func (a PathMustAPI) Uint64(path ...Key) uint64 {
	return PathMust(a.data, a.traverser, path, a.caster.AsUint64, a.onErr...)
}

func (a PathMustAPI) Uint64Slice(path ...Key) []uint64 {
	return PathMust(a.data, a.traverser, path, a.caster.AsUint64Slice, a.onErr...)
}

func (a PathMustAPI) String(path ...Key) string {
	return PathMust(a.data, a.traverser, path, a.caster.AsString, a.onErr...)
}

func (a PathMustAPI) StringSlice(path ...Key) []string {
	return PathMust(a.data, a.traverser, path, a.caster.AsStringSlice, a.onErr...)
}

func (a *PathMustAPI) Time(path ...Key) time.Time {
	return PathMust(a.data, a.traverser, path, a.caster.AsTime, a.onErr...)
}

func (a *PathMustAPI) TimeWithConfig(config cast.TimeCastConfig, path ...Key) time.Time {
	return PathMust(a.data, a.traverser, path, func(input any) (time.Time, error) {
		return a.caster.AsTimeWithConfig(config, input)
	}, a.onErr...)
}

func (a *PathMustAPI) TimeSlice(path ...Key) []time.Time {
	return PathMust(a.data, a.traverser, path, a.caster.AsTimeSlice, a.onErr...)
}

func (a *PathMustAPI) TimeSliceWithConfig(config cast.TimeCastConfig, path ...Key) []time.Time {
	return PathMust(a.data, a.traverser, path, func(input any) ([]time.Time, error) {
		return a.caster.AsTimeSliceWithConfig(config, input)
	}, a.onErr...)
}
