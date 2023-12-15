package testdata

type SampleStruct struct {
	A int32
	B uint32
	C string
}

var MixedTypesMap = map[string]any{
	"stringField": "abcd",
	"sliceOfAnyComplex": []any{
		int32(2),
		"stringElement",
		SampleStruct{A: 3, B: 4, C: "asdf"},
		map[string]string{
			"key1": "value1",
			"key2": "value2",
			"key3": "6565",
		},
		ptr(uint32(5)),
		[]bool{true, true, true, false, true, true},
		[]string{"abc", "def", "ghi"},
	},
	"pointerMapStringAny": &map[string]any{
		"fieldBool":  true,
		"fieldByte":  '.',
		"fieldInt32": int32(6),
		"int32Slice": []int32{10, 11, 12, 13, 14},
	},
	"float32":     float32(7.7),
	"bool":        true,
	"int32Number": int32(8),
}

func ptr[T any](x T) *T { return &x }
