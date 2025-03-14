package pick

import (
	"reflect"
	"time"
)

type Notation interface {
	Parse(selector string) ([]Key, error)
	Format(path ...Key) string
}

type Traverser interface {
	Retrieve(data any, path []Key) (any, error)
}

type ErrorGatherer interface {
	GatherSelector(selector string, err error)
}

type Converter interface {
	signedIntegerConverter
	unsignedIntegerConverter
	floatConverter
	stringConverter
	boolConverter
	byteConverter
	timeConverter
	durationConverter
	As(input any, asKind reflect.Kind) (any, error)
	ByType(input any, asType reflect.Type) (any, error)
}

type signedIntegerConverter interface {
	AsInt(item any) (int, error)
	AsInt8(item any) (int8, error)
	AsInt16(item any) (int16, error)
	AsInt32(item any) (int32, error)
	AsInt64(item any) (int64, error)

	AsIntSlice(item any) ([]int, error)
	AsInt8Slice(item any) ([]int8, error)
	AsInt16Slice(item any) ([]int16, error)
	AsInt32Slice(item any) ([]int32, error)
	AsInt64Slice(item any) ([]int64, error)
}

type unsignedIntegerConverter interface {
	AsUint(item any) (uint, error)
	AsUint8(item any) (uint8, error)
	AsUint16(item any) (uint16, error)
	AsUint32(item any) (uint32, error)
	AsUint64(item any) (uint64, error)

	AsUintSlice(item any) ([]uint, error)
	AsUint8Slice(item any) ([]uint8, error)
	AsUint16Slice(item any) ([]uint16, error)
	AsUint32Slice(item any) ([]uint32, error)
	AsUint64Slice(item any) ([]uint64, error)
}

type floatConverter interface {
	AsFloat32(item any) (float32, error)
	AsFloat64(item any) (float64, error)
	AsFloat32Slice(item any) ([]float32, error)
	AsFloat64Slice(item any) ([]float64, error)
}

type stringConverter interface {
	AsString(item any) (string, error)
	AsStringSlice(item any) ([]string, error)
}

type boolConverter interface {
	AsBool(item any) (bool, error)
	AsBoolSlice(input any) ([]bool, error)
}

type byteConverter interface {
	AsByte(item any) (byte, error)
	AsByteSlice(input any) ([]byte, error)
}

type timeConverter interface {
	AsTime(input any) (time.Time, error)
	AsTimeWithConfig(config TimeConvertConfig, input any) (time.Time, error)
	AsTimeSlice(input any) ([]time.Time, error)
	AsTimeSliceWithConfig(config TimeConvertConfig, input any) ([]time.Time, error)
}

type durationConverter interface {
	AsDuration(input any) (time.Duration, error)
	AsDurationWithConfig(config DurationConvertConfig, input any) (time.Duration, error)
	AsDurationSlice(input any) ([]time.Duration, error)
	AsDurationSliceWithConfig(config DurationConvertConfig, input any) ([]time.Duration, error)
}
