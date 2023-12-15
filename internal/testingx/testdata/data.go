package testdata

type SampleStruct struct {
	A int32
	B uint32
	C string
}

var MixedTypesMap = map[string]any{
	"stringField": "abcd",
	"sliceOfAnyComplex": []any{
		int32(1),
		"stringElement",
		SampleStruct{A: 1, B: 12, C: "asdf"},
		map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
	},
	"mapStringAny": map[string]any{},
	"float32":      float32(12.12),
	"bool":         true,
	"int32Number":  int32(12954),
}
