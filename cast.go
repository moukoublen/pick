package pick

import (
	"reflect"
	"time"

	"github.com/moukoublen/pick/cast"
)

type signedIntegerCaster interface {
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

type unsignedIntegerCaster interface {
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

type floatCaster interface {
	AsFloat32(item any) (float32, error)
	AsFloat64(item any) (float64, error)
	AsFloat32Slice(item any) ([]float32, error)
	AsFloat64Slice(item any) ([]float64, error)
}

type stringCaster interface {
	AsString(item any) (string, error)
	AsStringSlice(item any) ([]string, error)
}

type boolCaster interface {
	AsBool(item any) (bool, error)
	AsBoolSlice(input any) ([]bool, error)
}

type byteCaster interface {
	AsByte(item any) (byte, error)
	AsByteSlice(input any) ([]byte, error)
}

type timeCaster interface {
	AsTime(input any) (time.Time, error)
	AsTimeWithConfig(input any, config cast.TimeCastConfig) (time.Time, error)
	AsTimeSlice(input any) ([]time.Time, error)
	AsTimeSliceWithConfig(input any, config cast.TimeCastConfig) ([]time.Time, error)
}

type Caster interface {
	signedIntegerCaster
	unsignedIntegerCaster
	floatCaster
	stringCaster
	boolCaster
	byteCaster
	timeCaster
	As(input any, asKind reflect.Kind) (any, error)
}
