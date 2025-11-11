package pick

import (
	"reflect"
	"time"
)

// lets use this as global for performance reasons.
//
//nolint:gochecknoglobals
var convertFunctionTypes = newDirectConvertFunctionsTypes()

type directConvertFunctionsTypes struct {
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

func newDirectConvertFunctionsTypes() directConvertFunctionsTypes {
	return directConvertFunctionsTypes{
		typeOfBool:          reflect.TypeFor[bool](),
		typeOfInt8:          reflect.TypeFor[int8](),
		typeOfInt16:         reflect.TypeFor[int16](),
		typeOfInt32:         reflect.TypeFor[int32](),
		typeOfInt64:         reflect.TypeFor[int64](),
		typeOfInt:           reflect.TypeFor[int](),
		typeOfUint8:         reflect.TypeFor[uint8](),
		typeOfUint16:        reflect.TypeFor[uint16](),
		typeOfUint32:        reflect.TypeFor[uint32](),
		typeOfUint64:        reflect.TypeFor[uint64](),
		typeOfUint:          reflect.TypeFor[uint](),
		typeOfFloat32:       reflect.TypeFor[float32](),
		typeOfFloat64:       reflect.TypeFor[float64](),
		typeOfString:        reflect.TypeFor[string](),
		typeOfTime:          reflect.TypeFor[time.Time](),
		typeOfDuration:      reflect.TypeFor[time.Duration](),
		typeOfSliceBool:     reflect.TypeFor[[]bool](),
		typeOfSliceInt8:     reflect.TypeFor[[]int8](),
		typeOfSliceInt16:    reflect.TypeFor[[]int16](),
		typeOfSliceInt32:    reflect.TypeFor[[]int32](),
		typeOfSliceInt64:    reflect.TypeFor[[]int64](),
		typeOfSliceInt:      reflect.TypeFor[[]int](),
		typeOfSliceUint8:    reflect.TypeFor[[]uint8](),
		typeOfSliceUint16:   reflect.TypeFor[[]uint16](),
		typeOfSliceUint32:   reflect.TypeFor[[]uint32](),
		typeOfSliceUint64:   reflect.TypeFor[[]uint64](),
		typeOfSliceUint:     reflect.TypeFor[[]uint](),
		typeOfSliceFloat32:  reflect.TypeFor[[]float32](),
		typeOfSliceFloat64:  reflect.TypeFor[[]float64](),
		typeOfSliceString:   reflect.TypeFor[[]string](),
		typeOfSliceTime:     reflect.TypeFor[[]time.Time](),
		typeOfSliceDuration: reflect.TypeFor[[]time.Duration](),

		basicKindTypeMap: map[reflect.Kind]reflect.Type{
			reflect.Bool:    reflect.TypeFor[bool](),
			reflect.Int:     reflect.TypeFor[int](),
			reflect.Int8:    reflect.TypeFor[int8](),
			reflect.Int16:   reflect.TypeFor[int16](),
			reflect.Int32:   reflect.TypeFor[int32](),
			reflect.Int64:   reflect.TypeFor[int64](),
			reflect.Uint:    reflect.TypeFor[uint](),
			reflect.Uint8:   reflect.TypeFor[uint8](),
			reflect.Uint16:  reflect.TypeFor[uint16](),
			reflect.Uint32:  reflect.TypeFor[uint32](),
			reflect.Uint64:  reflect.TypeFor[uint64](),
			reflect.Float32: reflect.TypeFor[float32](),
			reflect.Float64: reflect.TypeFor[float64](),
			reflect.String:  reflect.TypeFor[string](),
		},
	}
}
