package testdata

import "math"

type SampleStruct struct {
	A int32
	B uint32
	C string
}

var MixedTypesMap = map[string]any{
	"stringField": "abcd",
	"sliceOfAnyComplex": []any{
		int32(2),                            // [0]
		"stringElement",                     // [1]
		SampleStruct{A: 3, B: 4, C: "asdf"}, // [2]
		map[string]string{
			"key1": "value1",
			"key2": "value2",
			"key3": "6565",
		}, // [3]
		ptr(uint32(5)), // [4]
		[]bool{true, true, true, false, true, true}, // [5]
		[]string{"abc", "def", "ghi"},               // [6]
		byte(math.MaxUint8),                         // [7]
	},
	"pointerMapStringAny": &map[string]any{
		"fieldBool":    true,
		"fieldByte":    '.',
		"fieldInt32":   int32(6),
		"int32Slice":   []int32{10, 11, 12, 13, 14},
		"float64Slice": []float64{0.1, 0.2, 0.3, 0.4},
	},
	"float32":     float32(7.7),
	"bool":        true,
	"int32Number": int32(8),
}

func ptr[T any](x T) *T { return &x }
