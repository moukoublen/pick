package pick

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/moukoublen/pick/iter"
)

func WrapJSON(js []byte) (Picker, error) {
	return WrapReaderJSON(bytes.NewReader(js))
}

func WrapReaderJSON(r io.Reader) (Picker, error) {
	d := json.NewDecoder(r)
	return WrapDecoder(d)
}

func WrapDecoder(decoder interface{ Decode(destination any) error }) (Picker, error) {
	var m any
	if err := decoder.Decode(&m); err != nil {
		return Picker{}, err
	}

	return Wrap(m), nil
}

// WrapJSONRequest reads the JSON request body from an HTTP request and wraps it into a Picker.
// It ensures proper cleanup of the request body to prevent resource leaks.
// Important note: After this function is called the body will be drained and closed.
func WrapJSONRequest(r *http.Request) (p Picker, rErr error) {
	if r == nil || r.Body == nil || r.Body == http.NoBody {
		return Wrap(nil), nil
	}

	defer drainAndClose(r.Body, &rErr)

	return WrapDecoder(json.NewDecoder(r.Body))
}

// WrapJSONResponse reads the JSON response body from an HTTP response and wraps it into a Picker.
// It ensures proper cleanup of the response body to prevent resource leaks.
// Important note: After this function is called the body will be drained and closed.
func WrapJSONResponse(r *http.Response) (p Picker, rErr error) {
	if r == nil || r.Body == nil || r.Body == http.NoBody {
		return Wrap(nil), nil
	}

	defer drainAndClose(r.Body, &rErr)

	return WrapDecoder(json.NewDecoder(r.Body))
}

func drainAndClose(b io.ReadCloser, outErr *error) {
	_, discardErr := io.Copy(io.Discard, b)
	*outErr = errors.Join(*outErr, discardErr, b.Close())
}

func Wrap(data any) Picker {
	caster := NewDefaultCaster()
	return NewPicker(data, NewDefaultTraverser(caster), caster, DotNotation{})
}

type Picker struct {
	data      any
	traverser Traverser
	Caster    Caster
	notation  Notation
}

func NewPicker(data any, t Traverser, c Caster, n Notation) Picker {
	return Picker{
		data:      data,
		traverser: t,
		Caster:    c,
		notation:  n,
	}
}

func (p Picker) Data() any { return p.data }

func (p Picker) Relaxed(onErr ...ErrorGatherer) RelaxedAPI {
	return RelaxedAPI{Picker: p, errorGatherers: onErr}
}

func (p Picker) Any(selector string) (any, error) {
	path, err := p.notation.Parse(selector)
	if err != nil {
		return nil, err
	}

	return p.Path(path)
}

func (p Picker) Path(path []Key) (any, error) {
	return p.traverser.Retrieve(p.data, path)
}

func (p Picker) Len(selector string) (int, error) {
	path, err := p.notation.Parse(selector)
	if err != nil {
		return 0, err
	}

	a, err := p.Path(path)
	if err != nil {
		return 0, err
	}

	return iter.Len(a)
}

// Default API (embedded into Picker)

func (p Picker) Each(selector string, operation func(index int, p Picker, length int) error) error {
	return Each(p, selector, operation)
}

func (p Picker) Bool(selector string) (bool, error) {
	return pickSelector(p, selector, p.Caster.AsBool)
}

func (p Picker) BoolSlice(selector string) ([]bool, error) {
	return pickSelector(p, selector, p.Caster.AsBoolSlice)
}

func (p Picker) Byte(selector string) (byte, error) {
	return pickSelector(p, selector, p.Caster.AsByte)
}

func (p Picker) ByteSlice(selector string) ([]byte, error) {
	return pickSelector(p, selector, p.Caster.AsByteSlice)
}

func (p Picker) Float32(selector string) (float32, error) {
	return pickSelector(p, selector, p.Caster.AsFloat32)
}

func (p Picker) Float32Slice(selector string) ([]float32, error) {
	return pickSelector(p, selector, p.Caster.AsFloat32Slice)
}

func (p Picker) Float64(selector string) (float64, error) {
	return pickSelector(p, selector, p.Caster.AsFloat64)
}

func (p Picker) Float64Slice(selector string) ([]float64, error) {
	return pickSelector(p, selector, p.Caster.AsFloat64Slice)
}

func (p Picker) Int(selector string) (int, error) {
	return pickSelector(p, selector, p.Caster.AsInt)
}

func (p Picker) IntSlice(selector string) ([]int, error) {
	return pickSelector(p, selector, p.Caster.AsIntSlice)
}

func (p Picker) Int8(selector string) (int8, error) {
	return pickSelector(p, selector, p.Caster.AsInt8)
}

func (p Picker) Int8Slice(selector string) ([]int8, error) {
	return pickSelector(p, selector, p.Caster.AsInt8Slice)
}

func (p Picker) Int16(selector string) (int16, error) {
	return pickSelector(p, selector, p.Caster.AsInt16)
}

func (p Picker) Int16Slice(selector string) ([]int16, error) {
	return pickSelector(p, selector, p.Caster.AsInt16Slice)
}

func (p Picker) Int32(selector string) (int32, error) {
	return pickSelector(p, selector, p.Caster.AsInt32)
}

func (p Picker) Int32Slice(selector string) ([]int32, error) {
	return pickSelector(p, selector, p.Caster.AsInt32Slice)
}

func (p Picker) Int64(selector string) (int64, error) {
	return pickSelector(p, selector, p.Caster.AsInt64)
}

func (p Picker) Int64Slice(selector string) ([]int64, error) {
	return pickSelector(p, selector, p.Caster.AsInt64Slice)
}

func (p Picker) Uint(selector string) (uint, error) {
	return pickSelector(p, selector, p.Caster.AsUint)
}

func (p Picker) UintSlice(selector string) ([]uint, error) {
	return pickSelector(p, selector, p.Caster.AsUintSlice)
}

func (p Picker) Uint8(selector string) (uint8, error) {
	return pickSelector(p, selector, p.Caster.AsUint8)
}

func (p Picker) Uint8Slice(selector string) ([]uint8, error) {
	return pickSelector(p, selector, p.Caster.AsUint8Slice)
}

func (p Picker) Uint16(selector string) (uint16, error) {
	return pickSelector(p, selector, p.Caster.AsUint16)
}

func (p Picker) Uint16Slice(selector string) ([]uint16, error) {
	return pickSelector(p, selector, p.Caster.AsUint16Slice)
}

func (p Picker) Uint32(selector string) (uint32, error) {
	return pickSelector(p, selector, p.Caster.AsUint32)
}

func (p Picker) Uint32Slice(selector string) ([]uint32, error) {
	return pickSelector(p, selector, p.Caster.AsUint32Slice)
}

func (p Picker) Uint64(selector string) (uint64, error) {
	return pickSelector(p, selector, p.Caster.AsUint64)
}

func (p Picker) Uint64Slice(selector string) ([]uint64, error) {
	return pickSelector(p, selector, p.Caster.AsUint64Slice)
}

func (p Picker) String(selector string) (string, error) {
	return pickSelector(p, selector, p.Caster.AsString)
}

func (p Picker) StringSlice(selector string) ([]string, error) {
	return pickSelector(p, selector, p.Caster.AsStringSlice)
}

func (p Picker) Time(selector string) (time.Time, error) {
	return pickSelector(p, selector, p.Caster.AsTime)
}

func (p Picker) TimeWithConfig(config TimeCastConfig, selector string) (time.Time, error) {
	return pickSelector(p, selector, func(input any) (time.Time, error) {
		return p.Caster.AsTimeWithConfig(config, input)
	})
}

func (p Picker) TimeSlice(selector string) ([]time.Time, error) {
	return pickSelector(p, selector, p.Caster.AsTimeSlice)
}

func (p Picker) TimeSliceWithConfig(config TimeCastConfig, selector string) ([]time.Time, error) {
	return pickSelector(p, selector, func(input any) ([]time.Time, error) {
		return p.Caster.AsTimeSliceWithConfig(config, input)
	})
}

func (p Picker) Duration(selector string) (time.Duration, error) {
	return pickSelector(p, selector, p.Caster.AsDuration)
}

func (p Picker) DurationWithConfig(config DurationCastConfig, selector string) (time.Duration, error) {
	return pickSelector(p, selector, func(input any) (time.Duration, error) {
		return p.Caster.AsDurationWithConfig(config, input)
	})
}

func (p Picker) DurationSlice(selector string) ([]time.Duration, error) {
	return pickSelector(p, selector, p.Caster.AsDurationSlice)
}

func (p Picker) DurationSliceWithConfig(config DurationCastConfig, selector string) ([]time.Duration, error) {
	return pickSelector(p, selector, func(input any) ([]time.Duration, error) {
		return p.Caster.AsDurationSliceWithConfig(config, input)
	})
}

// Wrap returns a new Picker using the same traverser, caster and notation.
func (p Picker) Wrap(data any) Picker {
	return NewPicker(data, p.traverser, p.Caster, p.notation)
}

// Relaxed API

type RelaxedAPI struct {
	Picker
	errorGatherers []ErrorGatherer
}

func (a RelaxedAPI) gather(selector string, err error) {
	for _, eg := range a.errorGatherers {
		eg.GatherSelector(selector, err)
	}
}

func (a RelaxedAPI) Each(selector string, operation func(index int, item RelaxedAPI, length int) error) {
	RelaxedEach(a, selector, operation)
}

func (a RelaxedAPI) Bool(selector string) bool {
	return pickRelaxed(a, selector, a.Caster.AsBool)
}

func (a RelaxedAPI) BoolSlice(selector string) []bool {
	return pickRelaxed(a, selector, a.Caster.AsBoolSlice)
}

func (a RelaxedAPI) Byte(selector string) byte {
	return pickRelaxed(a, selector, a.Caster.AsByte)
}

func (a RelaxedAPI) ByteSlice(selector string) []byte {
	return pickRelaxed(a, selector, a.Caster.AsByteSlice)
}

func (a RelaxedAPI) Float32(selector string) float32 {
	return pickRelaxed(a, selector, a.Caster.AsFloat32)
}

func (a RelaxedAPI) Float32Slice(selector string) []float32 {
	return pickRelaxed(a, selector, a.Caster.AsFloat32Slice)
}

func (a RelaxedAPI) Float64(selector string) float64 {
	return pickRelaxed(a, selector, a.Caster.AsFloat64)
}

func (a RelaxedAPI) Float64Slice(selector string) []float64 {
	return pickRelaxed(a, selector, a.Caster.AsFloat64Slice)
}

func (a RelaxedAPI) Int(selector string) int {
	return pickRelaxed(a, selector, a.Caster.AsInt)
}

func (a RelaxedAPI) IntSlice(selector string) []int {
	return pickRelaxed(a, selector, a.Caster.AsIntSlice)
}

func (a RelaxedAPI) Int8(selector string) int8 {
	return pickRelaxed(a, selector, a.Caster.AsInt8)
}

func (a RelaxedAPI) Int8Slice(selector string) []int8 {
	return pickRelaxed(a, selector, a.Caster.AsInt8Slice)
}

func (a RelaxedAPI) Int16(selector string) int16 {
	return pickRelaxed(a, selector, a.Caster.AsInt16)
}

func (a RelaxedAPI) Int16Slice(selector string) []int16 {
	return pickRelaxed(a, selector, a.Caster.AsInt16Slice)
}

func (a RelaxedAPI) Int32(selector string) int32 {
	return pickRelaxed(a, selector, a.Caster.AsInt32)
}

func (a RelaxedAPI) Int32Slice(selector string) []int32 {
	return pickRelaxed(a, selector, a.Caster.AsInt32Slice)
}

func (a RelaxedAPI) Int64(selector string) int64 {
	return pickRelaxed(a, selector, a.Caster.AsInt64)
}

func (a RelaxedAPI) Int64Slice(selector string) []int64 {
	return pickRelaxed(a, selector, a.Caster.AsInt64Slice)
}

func (a RelaxedAPI) Uint(selector string) uint {
	return pickRelaxed(a, selector, a.Caster.AsUint)
}

func (a RelaxedAPI) UintSlice(selector string) []uint {
	return pickRelaxed(a, selector, a.Caster.AsUintSlice)
}

func (a RelaxedAPI) Uint8(selector string) uint8 {
	return pickRelaxed(a, selector, a.Caster.AsUint8)
}

func (a RelaxedAPI) Uint8Slice(selector string) []uint8 {
	return pickRelaxed(a, selector, a.Caster.AsUint8Slice)
}

func (a RelaxedAPI) Uint16(selector string) uint16 {
	return pickRelaxed(a, selector, a.Caster.AsUint16)
}

func (a RelaxedAPI) Uint16Slice(selector string) []uint16 {
	return pickRelaxed(a, selector, a.Caster.AsUint16Slice)
}

func (a RelaxedAPI) Uint32(selector string) uint32 {
	return pickRelaxed(a, selector, a.Caster.AsUint32)
}

func (a RelaxedAPI) Uint32Slice(selector string) []uint32 {
	return pickRelaxed(a, selector, a.Caster.AsUint32Slice)
}

func (a RelaxedAPI) Uint64(selector string) uint64 {
	return pickRelaxed(a, selector, a.Caster.AsUint64)
}

func (a RelaxedAPI) Uint64Slice(selector string) []uint64 {
	return pickRelaxed(a, selector, a.Caster.AsUint64Slice)
}

func (a RelaxedAPI) String(selector string) string {
	return pickRelaxed(a, selector, a.Caster.AsString)
}

func (a RelaxedAPI) StringSlice(selector string) []string {
	return pickRelaxed(a, selector, a.Caster.AsStringSlice)
}

func (a RelaxedAPI) Time(selector string) time.Time {
	return pickRelaxed(a, selector, a.Caster.AsTime)
}

func (a RelaxedAPI) TimeWithConfig(config TimeCastConfig, selector string) time.Time {
	return pickRelaxed(a, selector, func(input any) (time.Time, error) {
		return a.Caster.AsTimeWithConfig(config, input)
	})
}

func (a RelaxedAPI) TimeSlice(selector string) []time.Time {
	return pickRelaxed(a, selector, a.Caster.AsTimeSlice)
}

func (a RelaxedAPI) TimeSliceWithConfig(config TimeCastConfig, selector string) []time.Time {
	return pickRelaxed(a, selector, func(input any) ([]time.Time, error) {
		return a.Caster.AsTimeSliceWithConfig(config, input)
	})
}

func (a RelaxedAPI) Duration(selector string) time.Duration {
	return pickRelaxed(a, selector, a.Caster.AsDuration)
}

func (a RelaxedAPI) DurationWithConfig(config DurationCastConfig, selector string) time.Duration {
	return pickRelaxed(a, selector, func(input any) (time.Duration, error) {
		return a.Caster.AsDurationWithConfig(config, input)
	})
}

func (a RelaxedAPI) DurationSlice(selector string) []time.Duration {
	return pickRelaxed(a, selector, a.Caster.AsDurationSlice)
}

func (a RelaxedAPI) DurationSliceWithConfig(config DurationCastConfig, selector string) []time.Duration {
	return pickRelaxed(a, selector, func(input any) ([]time.Duration, error) {
		return a.Caster.AsDurationSliceWithConfig(config, input)
	})
}

func (a RelaxedAPI) Wrap(data any) RelaxedAPI {
	return NewPicker(data, a.traverser, a.Caster, a.notation).Relaxed(a.errorGatherers...)
}

//nolint:ireturn
func pickSelector[Output any](p Picker, selector string, castFn func(any) (Output, error)) (Output, error) {
	item, err := p.Any(selector)
	if err != nil {
		var o Output
		return o, err
	}

	return castFn(item)
}

//nolint:ireturn
func pickRelaxed[Output any](a RelaxedAPI, selector string, castFn func(any) (Output, error)) Output {
	casted, err := pickSelector(a.Picker, selector, castFn)
	if err != nil {
		a.gather(selector, err)
	}
	return casted
}
