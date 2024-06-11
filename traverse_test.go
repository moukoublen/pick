package pick

import (
	"testing"

	"github.com/moukoublen/pick/internal/tst"
)

func TestDefaultTraverser(t *testing.T) {
	// t.Parallel()

	type renamed map[string]any
	type itemOne struct {
		FieldOne string
		FieldTwo int
	}

	tests := map[string]struct {
		input         any
		expected      any
		errorAsserter tst.ErrorAsserter
		keys          []Key
	}{
		"nil": {
			input:         nil,
			keys:          nil,
			expected:      nil,
			errorAsserter: tst.NoError,
		},

		"access zero level": {
			input:         []any{"one", "two"},
			keys:          []Key{},
			expected:      []any{"one", "two"},
			errorAsserter: tst.NoError,
		},

		"index access level 1": {
			input:         []any{"one", "two"},
			keys:          []Key{Index(1)},
			expected:      "two",
			errorAsserter: tst.NoError,
		},

		"index access level 1 negative index": {
			input:         []any{"one", "two", "three"},
			keys:          []Key{Index(-1)},
			expected:      "three",
			errorAsserter: tst.NoError,
		},

		"index access level 1 negative index 2": {
			input:         []any{"one", "two", "three"},
			keys:          []Key{Index(-2)},
			expected:      "two",
			errorAsserter: tst.NoError,
		},

		"index access level 1 out of range": {
			input:    []any{"one", "two"},
			keys:     []Key{Index(6)},
			expected: nil,
			errorAsserter: tst.ExpectedErrorChecks(
				tst.ExpectedErrorIs(ErrIndexOutOfRange),
				tst.ExpectedErrorOfType[*TraverseError](
					func(t *testing.T, te *TraverseError) { //nolint:thelper
						tst.AssertEqual(t, "selector: [6] : error trying to traverse: field not found: index out of range", te.Error())
						tst.AssertEqual(t, te.Path(), []Key{Index(6)})
					},
				),
			),
		},

		"index access slice of string level 1": {
			input:         []string{"one", "two"},
			keys:          []Key{Index(1)},
			expected:      "two",
			errorAsserter: tst.NoError,
		},

		"name access level 1": {
			input:         map[string]any{"one": "value"},
			keys:          []Key{Field("one")},
			expected:      "value",
			errorAsserter: tst.NoError,
		},

		"name access level 1 happy path": {
			input:         map[string]any{"one": "value"},
			keys:          []Key{Field("one")},
			expected:      "value",
			errorAsserter: tst.NoError,
		},

		"name access level 1 not found": {
			input:         map[string]any{"one": "value"},
			keys:          []Key{Field("two")},
			expected:      nil,
			errorAsserter: tst.ExpectedErrorIs(ErrFieldNotFound),
		},

		"name access level 2 happy path": {
			input:         map[string]any{"one": map[string]any{"two": "value"}},
			keys:          []Key{Field("one"), Field("two")},
			expected:      "value",
			errorAsserter: tst.NoError,
		},

		"name access level 2 but nil": {
			input:    map[string]any{"one": nil},
			keys:     []Key{Field("one"), Field("two")},
			expected: nil,
			errorAsserter: tst.ExpectedErrorChecks(
				tst.ExpectedErrorIs(ErrFieldNotFound),
				tst.ExpectedErrorOfType[*TraverseError](
					func(t *testing.T, te *TraverseError) { //nolint:thelper
						tst.AssertEqual(t, te.Error(), "selector: one.two : error trying to traverse: field not found")
						tst.AssertEqual(t, te.Path(), []Key{Field("one"), Field("two")})
					},
				),
			),
		},

		"name access level 3 but not exists": {
			input:    map[string]any{"one": 12},
			keys:     []Key{Field("one"), Field("two"), Field("tree")},
			expected: nil,
			errorAsserter: tst.ExpectedErrorChecks(
				tst.ExpectedErrorIs(ErrFieldNotFound),
				tst.ExpectedErrorOfType[*TraverseError](
					func(t *testing.T, te *TraverseError) { //nolint:thelper
						tst.AssertEqual(t, te.Error(), "selector: one.two : error trying to traverse: field not found")
						tst.AssertEqual(t, te.Path(), []Key{Field("one"), Field("two")})
					},
				),
			),
		},

		"name access level 2 renamed happy path": {
			input:         map[string]any{"one": renamed{"two": "value"}},
			keys:          []Key{Field("one"), Field("two")},
			expected:      "value",
			errorAsserter: tst.NoError,
		},

		"mixed access level 2": {
			input:         []any{"one", map[string]any{"two": "value"}},
			keys:          []Key{Index(1), Field("two")},
			expected:      "value",
			errorAsserter: tst.NoError,
		},

		"mixed access level 2 with cast": {
			input:         []any{"one", map[string]any{"4": "value"}},
			keys:          []Key{Index(1), Index(4)},
			expected:      "value",
			errorAsserter: tst.NoError,
		},

		"mixed access level 2 with cast 2": {
			input:         map[string]any{"one": []string{"s0", "s1", "s2"}},
			keys:          []Key{Field("one"), Field("1")},
			expected:      "s1",
			errorAsserter: tst.NoError,
		},

		"mixed access level 3 with struct": {
			input:         map[string]any{"one": renamed{"two": itemOne{FieldOne: "test", FieldTwo: 123}}},
			keys:          []Key{Field("one"), Field("two"), Field("FieldOne")},
			expected:      "test",
			errorAsserter: tst.NoError,
		},

		"mixed access level 3 with struct using index": {
			input:         map[string]any{"one": renamed{"two": itemOne{FieldOne: "test", FieldTwo: 123}}},
			keys:          []Key{Field("one"), Field("two"), Index(0)},
			expected:      "test",
			errorAsserter: tst.NoError,
		},

		"mixed access level 3 with struct using wrong index": {
			input:    map[string]any{"one": renamed{"two": itemOne{FieldOne: "test", FieldTwo: 123}}},
			keys:     []Key{Field("one"), Field("two"), Index(12)},
			expected: nil,
			errorAsserter: tst.ExpectedErrorChecks(
				tst.ExpectedErrorOfType[*TraverseError](),
				tst.ExpectedErrorStringContains("reflect: Field index out of range"),
			),
		},

		"mixed access level 3 with struct using wrong field": {
			input:    map[string]any{"one": renamed{"two": itemOne{FieldOne: "test", FieldTwo: 123}}},
			keys:     []Key{Field("one"), Field("two"), Field("Wrong")},
			expected: nil,
			errorAsserter: tst.ExpectedErrorChecks(
				tst.ExpectedErrorOfType[*TraverseError](),
				tst.ExpectedErrorIs(ErrFieldNotFound),
			),
		},

		"mixed access level 3 with pointer struct": {
			input:         map[string]any{"one": renamed{"two": &itemOne{FieldOne: "test", FieldTwo: 123}}},
			keys:          []Key{Field("one"), Field("two"), Field("FieldOne")},
			expected:      "test",
			errorAsserter: tst.NoError,
		},

		"mixed access level 3 with map with int32 key": {
			input:         map[string]any{"one": map[int32]itemOne{42: {FieldOne: "test", FieldTwo: 123}}},
			keys:          []Key{Field("one"), Field("42"), Index(0)},
			expected:      "test",
			errorAsserter: tst.NoError,
		},

		"mixed access level 3 with map with int32 key and index fields": {
			input:         map[string]any{"one": map[int32]itemOne{42: {FieldOne: "test", FieldTwo: 123}}},
			keys:          []Key{Field("one"), Index(42), Index(0)},
			expected:      "test",
			errorAsserter: tst.NoError,
		},

		"index access slice of int level 1": {
			input:         []int{1, 2, 3},
			keys:          []Key{Index(1)},
			expected:      2,
			errorAsserter: tst.NoError,
		},
		"index access slice of int8 level 1": {
			input:         []int8{1, 2, 3},
			keys:          []Key{Index(1)},
			expected:      int8(2),
			errorAsserter: tst.NoError,
		},
		"index access slice of int16 level 1": {
			input:         []int16{1, 2, 3},
			keys:          []Key{Index(1)},
			expected:      int16(2),
			errorAsserter: tst.NoError,
		},
		"index access slice of int32 level 1": {
			input:         []int32{1, 2, 3},
			keys:          []Key{Index(1)},
			expected:      int32(2),
			errorAsserter: tst.NoError,
		},
		"index access slice of int64 level 1": {
			input:         []int64{1, 2, 3},
			keys:          []Key{Index(1)},
			expected:      int64(2),
			errorAsserter: tst.NoError,
		},
		"index access slice of uint level 1": {
			input:         []uint{1, 2, 3},
			keys:          []Key{Index(1)},
			expected:      uint(2),
			errorAsserter: tst.NoError,
		},
		"index access slice of uint8 level 1": {
			input:         []uint8{1, 2, 3},
			keys:          []Key{Index(1)},
			expected:      uint8(2),
			errorAsserter: tst.NoError,
		},
		"index access slice of uint16 level 1": {
			input:         []uint16{1, 2, 3},
			keys:          []Key{Index(1)},
			expected:      uint16(2),
			errorAsserter: tst.NoError,
		},
		"index access slice of uint32 level 1": {
			input:         []uint32{1, 2, 3},
			keys:          []Key{Index(1)},
			expected:      uint32(2),
			errorAsserter: tst.NoError,
		},
		"index access slice of uint64 level 1": {
			input:         []uint64{1, 2, 3},
			keys:          []Key{Index(1)},
			expected:      uint64(2),
			errorAsserter: tst.NoError,
		},
		"index access slice of float32 level 1": {
			input:         []float32{1, 2, 3},
			keys:          []Key{Index(1)},
			expected:      float32(2),
			errorAsserter: tst.NoError,
		},
		"index access slice of float64 level 1": {
			input:         []float64{1, 2, 3},
			keys:          []Key{Index(1)},
			expected:      float64(2),
			errorAsserter: tst.NoError,
		},
		"index access slice of bool level 1": {
			input:         []bool{false, true, true},
			keys:          []Key{Index(1)},
			expected:      true,
			errorAsserter: tst.NoError,
		},
		"index access slice of foo level 1": {
			input:         []foo{{A: 1}, {A: 2}, {A: 3}},
			keys:          []Key{Index(1)},
			expected:      foo{A: 2},
			errorAsserter: tst.NoError,
		},
	}

	dt := DefaultTraverser{
		keyCaster: NewDefaultCaster(),
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// t.Parallel()
			got, err := dt.Retrieve(tc.input, tc.keys)

			// check error
			tc.errorAsserter(t, err)

			// check returned item
			tst.AssertEqual(t, got, tc.expected)
		})
	}
}

type foo struct{ A int }

func BenchmarkDefaultTraverser(b *testing.B) {
	type renamed map[string]any
	type itemOne struct {
		FieldOne string
		FieldTwo int
	}

	tests := map[string]struct {
		input any
		keys  []Key
	}{
		"[]any": {
			input: []any{"one", "two"},
			keys:  []Key{Index(1)},
		},

		"[]any(index out of range)": {
			input: []any{"one", "two"},
			keys:  []Key{Index(6)},
		},

		"[]string": {
			input: []string{"one", "two"},
			keys:  []Key{Index(1)},
		},
		"[]int": {
			input: []int{1, 2, 3},
			keys:  []Key{Index(1)},
		},
		"[]int8": {
			input: []int8{1, 2, 3},
			keys:  []Key{Index(1)},
		},
		"[]int16": {
			input: []int16{1, 2, 3},
			keys:  []Key{Index(1)},
		},
		"[]int32": {
			input: []int32{1, 2, 3},
			keys:  []Key{Index(1)},
		},
		"[]int64": {
			input: []int64{1, 2, 3},
			keys:  []Key{Index(1)},
		},
		"[]uint": {
			input: []uint{1, 2, 3},
			keys:  []Key{Index(1)},
		},
		"[]uint8": {
			input: []uint8{1, 2, 3},
			keys:  []Key{Index(1)},
		},
		"[]uint16": {
			input: []uint16{1, 2, 3},
			keys:  []Key{Index(1)},
		},
		"[]uint32": {
			input: []uint32{1, 2, 3},
			keys:  []Key{Index(1)},
		},
		"[]uint64": {
			input: []uint64{1, 2, 3},
			keys:  []Key{Index(1)},
		},
		"[]float32": {
			input: []float32{1, 2, 3},
			keys:  []Key{Index(1)},
		},
		"[]float64": {
			input: []float64{1, 2, 3},
			keys:  []Key{Index(1)},
		},
		"[]bool": {
			input: []bool{false, true, false},
			keys:  []Key{Index(1)},
		},

		"[]map[string]string": {
			input: []map[string]string{{}, {}, {}},
			keys:  []Key{Index(1)},
		},
		"[]foo": {
			input: []foo{{A: 1}, {A: 2}, {A: 3}},
			keys:  []Key{Index(1)},
		},

		"[]any -> []any": {
			input: []any{"one", []any{"two", "three"}},
			keys:  []Key{Index(1), Index(1)},
		},

		"[]any -> []any -> []any": {
			input: []any{"one", []any{"two", []any{"three", "four"}}},
			keys:  []Key{Index(1), Index(1), Index(1)},
		},

		"[]any -> []any -> []string": {
			input: []any{"one", []any{"two", []string{"three", "four"}}},
			keys:  []Key{Index(1), Index(1), Index(1)},
		},

		"map[string]any": {
			input: map[string]any{"one": "value"},
			keys:  []Key{Field("one")},
		},

		"map[string]any(field not found)": {
			input: map[string]any{"one": "value"},
			keys:  []Key{Field("two")},
		},

		"map[string]any -> map[string]any": {
			input: map[string]any{"one": map[string]any{"two": "value"}},
			keys:  []Key{Field("one"), Field("two")},
		},

		"map[string]any -> renamed(map[string]any)": {
			input: map[string]any{"one": renamed{"two": "value"}},
			keys:  []Key{Field("one"), Field("two")},
		},

		"map[string]any -> map[string]any -> map[string]any": {
			input: map[string]any{"one": map[string]any{"two": map[string]any{"three": "value"}}},
			keys:  []Key{Field("one"), Field("two"), Field("three")},
		},

		"map[string]any -> renamed(map[string]any) -> map[string]any": {
			input: map[string]any{"one": renamed{"two": map[string]any{"three": "value"}}},
			keys:  []Key{Field("one"), Field("two")},
		},

		"struct": {
			input: itemOne{FieldOne: "one", FieldTwo: 2},
			keys:  []Key{Field("FieldOne")},
		},

		"[]any -> map[string]any": {
			input: []any{"one", map[string]any{"two": "value"}},
			keys:  []Key{Index(1), Field("two")},
		},

		"[]any -> map[string]any(using index)": {
			input: []any{"one", map[string]any{"4": "value"}},
			keys:  []Key{Index(1), Index(4)},
		},

		"map[string]any -> []string(using field)": {
			input: map[string]any{"one": []string{"s0", "s1", "s2"}},
			keys:  []Key{Field("one"), Field("1")},
		},

		"[]any -> map[string]any -> []any": {
			input: []any{"one", map[string]any{"two": []any{"a", "b", "c"}}},
			keys:  []Key{Index(1), Field("two"), Index(2)},
		},

		"map[string]any -> renamed(map[string]any) -> struct": {
			input: map[string]any{"one": renamed{"two": itemOne{FieldOne: "test", FieldTwo: 123}}},
			keys:  []Key{Field("one"), Field("two"), Field("FieldOne")},
		},

		"map[string]any -> renamed(map[string]any) -> struct(using index)": {
			input: map[string]any{"one": renamed{"two": itemOne{FieldOne: "test", FieldTwo: 123}}},
			keys:  []Key{Field("one"), Field("two"), Index(0)},
		},

		"map[string]any -> renamed(map[string]any) -> struct(using wrong index)": {
			input: map[string]any{"one": renamed{"two": itemOne{FieldOne: "test", FieldTwo: 123}}},
			keys:  []Key{Field("one"), Field("two"), Index(12)},
		},

		"map[string]any -> renamed(map[string]any) -> struct(using wrong field)": {
			input: map[string]any{"one": renamed{"two": itemOne{FieldOne: "test", FieldTwo: 123}}},
			keys:  []Key{Field("one"), Field("two"), Field("Wrong")},
		},

		"map[string]any -> renamed(map[string]any) -> &struct": {
			input: map[string]any{"one": renamed{"two": &itemOne{FieldOne: "test", FieldTwo: 123}}},
			keys:  []Key{Field("one"), Field("two"), Field("FieldOne")},
		},

		"map[string]any -> map[int32]struct(using field) -> struct(using index)": {
			input: map[string]any{"one": map[int32]itemOne{42: {FieldOne: "test", FieldTwo: 123}}},
			keys:  []Key{Field("one"), Field("42"), Index(0)},
		},

		"map[string]any -> map[int32]struct(using index) -> struct(using index)": {
			input: map[string]any{"one": map[int32]itemOne{42: {FieldOne: "test", FieldTwo: 123}}},
			keys:  []Key{Field("one"), Index(42), Index(0)},
		},
	}

	dt := DefaultTraverser{
		keyCaster: NewDefaultCaster(),
	}

	for name, tc := range tests {
		b.Run(name, func(b *testing.B) {
			for range b.N {
				_, _ = dt.Retrieve(tc.input, tc.keys)
			}
		})
	}
}

func TestTraverseError(t *testing.T) { //nolint:thelper
	err1 := NewTraverseError("not good", []Key{Field("one")}, 0, nil)
	tst.AssertEqual(t, err1.Error(), "selector: one : not good")
	tst.AssertEqual(t, err1.Unwrap(), error(nil))
	tst.AssertEqual(t, err1.Path(), []Key{Field("one")})

	err2 := NewTraverseError("not good", []Key{Field("one"), Field("two")}, 0, nil)
	tst.AssertEqual(t, err2.Error(), "selector: one : not good")
	tst.AssertEqual(t, err2.Unwrap(), error(nil))
	tst.AssertEqual(t, err2.Path(), []Key{Field("one")})

	err3 := NewTraverseError("not good", []Key{Field("one"), Field("two"), Index(2)}, 1, ErrFieldNotFound)
	tst.AssertEqual(t, err3.Error(), "selector: one.two : not good: field not found")
	tst.ExpectedErrorIs(ErrFieldNotFound)(t, err3)
	tst.AssertEqual(t, err3.Path(), []Key{Field("one"), Field("two")})
}

func TestSet(t *testing.T) {
	dt := DefaultTraverser{
		caster: cast.NewCaster(),
	}

	tests := map[string]struct {
		Destination   any
		Path          []Key
		ValueToSet    any
		ExpectedError func(*testing.T, error)
		Expected      any
	}{
		"1 level map add new key": {
			Destination: map[string]string{
				"one": "a",
				"two": "b",
			},
			Path:          []Key{Field("three")},
			ValueToSet:    "c",
			ExpectedError: nil,
			Expected: map[string]string{
				"one":   "a",
				"two":   "b",
				"three": "c",
			},
		},
		"2 level map add new key": {
			Destination: map[string]any{
				"one":   "a",
				"two":   "b",
				"three": map[string]int{"inner": 0},
			},
			Path:          []Key{Field("three"), Field("inner2")},
			ValueToSet:    12,
			ExpectedError: nil,
			Expected: map[string]any{
				"one": "a",
				"two": "b",
				"three": map[string]int{
					"inner":  0,
					"inner2": 12,
				},
			},
		},
		"replace value in slice": {
			Destination: []string{
				"one",
				"two",
			},
			Path:          []Key{Index(1)},
			ValueToSet:    "c",
			ExpectedError: nil,
			Expected: []string{
				"one",
				"c",
			},
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			err := dt.Set(tc.Destination, tc.Path, tc.ValueToSet)
			testingx.AssertError(t, tc.ExpectedError, err)
			testingx.AssertEqual(t, tc.Destination, tc.Expected)
		})
	}
}
