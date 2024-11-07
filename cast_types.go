package pick

import (
	"reflect"
	"time"
)

// lets use this as global for performance reasons.
//
//nolint:gochecknoglobals
var castFunctionTypes = newDirectCastFunctionsTypes()

type directCastFunctionsTypes struct {
	typeOfBool reflect.Type
	// typeOfByte          reflect.Type // there is no distinguish type for byte. Its only uint8.
	typeOfInt8      reflect.Type
	typeOfInt16     reflect.Type
	typeOfInt32     reflect.Type
	typeOfInt64     reflect.Type
	typeOfInt       reflect.Type
	typeOfUint8     reflect.Type
	typeOfUint16    reflect.Type
	typeOfUint32    reflect.Type
	typeOfUint64    reflect.Type
	typeOfUint      reflect.Type
	typeOfFloat32   reflect.Type
	typeOfFloat64   reflect.Type
	typeOfString    reflect.Type
	typeOfTime      reflect.Type
	typeOfDuration  reflect.Type
	typeOfSliceBool reflect.Type
	// typeOfSliceByte     reflect.Type // there is no distinguish type for byte. Its only uint8.
	typeOfSliceInt8     reflect.Type
	typeOfSliceInt16    reflect.Type
	typeOfSliceInt32    reflect.Type
	typeOfSliceInt64    reflect.Type
	typeOfSliceInt      reflect.Type
	typeOfSliceUint8    reflect.Type
	typeOfSliceUint16   reflect.Type
	typeOfSliceUint32   reflect.Type
	typeOfSliceUint64   reflect.Type
	typeOfSliceUint     reflect.Type
	typeOfSliceFloat32  reflect.Type
	typeOfSliceFloat64  reflect.Type
	typeOfSliceString   reflect.Type
	typeOfSliceTime     reflect.Type
	typeOfSliceDuration reflect.Type

	basicKindTypeMap map[reflect.Kind]reflect.Type
}

func newDirectCastFunctionsTypes() directCastFunctionsTypes {
	return directCastFunctionsTypes{
		typeOfBool:          reflect.TypeOf(false),
		typeOfInt8:          reflect.TypeOf(int8(0)),
		typeOfInt16:         reflect.TypeOf(int16(0)),
		typeOfInt32:         reflect.TypeOf(int32(0)),
		typeOfInt64:         reflect.TypeOf(int64(0)),
		typeOfInt:           reflect.TypeOf(int(0)),
		typeOfUint8:         reflect.TypeOf(uint8(0)),
		typeOfUint16:        reflect.TypeOf(uint16(0)),
		typeOfUint32:        reflect.TypeOf(uint32(0)),
		typeOfUint64:        reflect.TypeOf(uint64(0)),
		typeOfUint:          reflect.TypeOf(uint(0)),
		typeOfFloat32:       reflect.TypeOf(float32(0)),
		typeOfFloat64:       reflect.TypeOf(float64(0)),
		typeOfString:        reflect.TypeOf(""),
		typeOfTime:          reflect.TypeOf(time.Time{}),
		typeOfDuration:      reflect.TypeOf(time.Duration(0)),
		typeOfSliceBool:     reflect.TypeOf([]bool{}),
		typeOfSliceInt8:     reflect.TypeOf([]int8{}),
		typeOfSliceInt16:    reflect.TypeOf([]int16{}),
		typeOfSliceInt32:    reflect.TypeOf([]int32{}),
		typeOfSliceInt64:    reflect.TypeOf([]int64{}),
		typeOfSliceInt:      reflect.TypeOf([]int{}),
		typeOfSliceUint8:    reflect.TypeOf([]uint8{}),
		typeOfSliceUint16:   reflect.TypeOf([]uint16{}),
		typeOfSliceUint32:   reflect.TypeOf([]uint32{}),
		typeOfSliceUint64:   reflect.TypeOf([]uint64{}),
		typeOfSliceUint:     reflect.TypeOf([]uint{}),
		typeOfSliceFloat32:  reflect.TypeOf([]float32{}),
		typeOfSliceFloat64:  reflect.TypeOf([]float64{}),
		typeOfSliceString:   reflect.TypeOf([]string{}),
		typeOfSliceTime:     reflect.TypeOf([]time.Time{}),
		typeOfSliceDuration: reflect.TypeOf([]time.Duration{}),

		basicKindTypeMap: map[reflect.Kind]reflect.Type{
			reflect.Bool:    reflect.TypeOf(false),
			reflect.Int:     reflect.TypeOf(int(0)),
			reflect.Int8:    reflect.TypeOf(int8(0)),
			reflect.Int16:   reflect.TypeOf(int16(0)),
			reflect.Int32:   reflect.TypeOf(int32(0)),
			reflect.Int64:   reflect.TypeOf(int64(0)),
			reflect.Uint:    reflect.TypeOf(uint(0)),
			reflect.Uint8:   reflect.TypeOf(uint8(0)),
			reflect.Uint16:  reflect.TypeOf(uint16(0)),
			reflect.Uint32:  reflect.TypeOf(uint32(0)),
			reflect.Uint64:  reflect.TypeOf(uint64(0)),
			reflect.Float32: reflect.TypeOf(float32(0)),
			reflect.Float64: reflect.TypeOf(float64(0)),
			reflect.String:  reflect.TypeOf(""),
		},
	}
}
