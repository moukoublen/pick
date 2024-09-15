package pick

import (
	"reflect"
	"testing"

	"github.com/moukoublen/pick/cast"
	"github.com/moukoublen/pick/internal/testingx"
)

func TestDefaultTraverser(t *testing.T) {
	t.Parallel()

	type renamed map[string]any
	type itemOne struct {
		FieldOne string
		FieldTwo int
	}

	tests := map[string]struct {
		input       any
		expected    any
		expectedErr func(*testing.T, error)
		keys        []Key
	}{
		"nil": {
			input:       nil,
			keys:        nil,
			expected:    nil,
			expectedErr: nil,
		},

		"access zero level": {
			input:       []any{"one", "two"},
			keys:        []Key{},
			expected:    []any{"one", "two"},
			expectedErr: nil,
		},

		"index access level 1": {
			input:       []any{"one", "two"},
			keys:        []Key{Index(1)},
			expected:    "two",
			expectedErr: nil,
		},

		"index access level 1 negative index": {
			input:       []any{"one", "two", "three"},
			keys:        []Key{Index(-1)},
			expected:    "three",
			expectedErr: nil,
		},

		"index access level 1 negative index 2": {
			input:       []any{"one", "two", "three"},
			keys:        []Key{Index(-2)},
			expected:    "two",
			expectedErr: nil,
		},

		"index access level 1 out of range": {
			input:       []any{"one", "two"},
			keys:        []Key{Index(6)},
			expected:    nil,
			expectedErr: testingx.ExpectedErrorIs(ErrIndexOutOfRange),
		},

		"index access slice of string level 1": {
			input:       []string{"one", "two"},
			keys:        []Key{Index(1)},
			expected:    "two",
			expectedErr: nil,
		},

		"name access level 1": {
			input:       map[string]any{"one": "value"},
			keys:        []Key{Field("one")},
			expected:    "value",
			expectedErr: nil,
		},

		"name access level 1 happy path": {
			input:       map[string]any{"one": "value"},
			keys:        []Key{Field("one")},
			expected:    "value",
			expectedErr: nil,
		},

		"name access level 1 not found": {
			input:       map[string]any{"one": "value"},
			keys:        []Key{Field("two")},
			expected:    nil,
			expectedErr: testingx.ExpectedErrorIs(ErrFieldNotFound),
		},

		"name access level 2 happy path": {
			input:       map[string]any{"one": map[string]any{"two": "value"}},
			keys:        []Key{Field("one"), Field("two")},
			expected:    "value",
			expectedErr: nil,
		},

		"name access level 2 but nil": {
			input:       map[string]any{"one": nil},
			keys:        []Key{Field("one"), Field("two")},
			expected:    nil,
			expectedErr: testingx.ExpectedErrorIs(ErrFieldNotFound),
		},

		"name access level 2 renamed happy path": {
			input:       map[string]any{"one": renamed{"two": "value"}},
			keys:        []Key{Field("one"), Field("two")},
			expected:    "value",
			expectedErr: nil,
		},

		"mixed access level 2": {
			input:       []any{"one", map[string]any{"two": "value"}},
			keys:        []Key{Index(1), Field("two")},
			expected:    "value",
			expectedErr: nil,
		},

		"mixed access level 2 with cast": {
			input:       []any{"one", map[string]any{"4": "value"}},
			keys:        []Key{Index(1), Index(4)},
			expected:    "value",
			expectedErr: nil,
		},

		"mixed access level 2 with cast 2": {
			input:       map[string]any{"one": []string{"s0", "s1", "s2"}},
			keys:        []Key{Field("one"), Field("1")},
			expected:    "s1",
			expectedErr: nil,
		},

		"mixed access level 3 with struct": {
			input:       map[string]any{"one": renamed{"two": itemOne{FieldOne: "test", FieldTwo: 123}}},
			keys:        []Key{Field("one"), Field("two"), Field("FieldOne")},
			expected:    "test",
			expectedErr: nil,
		},

		"mixed access level 3 with struct using index": {
			input:       map[string]any{"one": renamed{"two": itemOne{FieldOne: "test", FieldTwo: 123}}},
			keys:        []Key{Field("one"), Field("two"), Index(0)},
			expected:    "test",
			expectedErr: nil,
		},

		"mixed access level 3 with struct using wrong index": {
			input:       map[string]any{"one": renamed{"two": itemOne{FieldOne: "test", FieldTwo: 123}}},
			keys:        []Key{Field("one"), Field("two"), Index(12)},
			expected:    nil,
			expectedErr: testingx.ExpectedErrorStringContains("reflect: Field index out of range"),
		},

		"mixed access level 3 with struct using wrong field": {
			input:       map[string]any{"one": renamed{"two": itemOne{FieldOne: "test", FieldTwo: 123}}},
			keys:        []Key{Field("one"), Field("two"), Field("Wrong")},
			expected:    nil,
			expectedErr: testingx.ExpectedErrorIs(ErrFieldNotFound),
		},

		"mixed access level 3 with pointer struct": {
			input:       map[string]any{"one": renamed{"two": &itemOne{FieldOne: "test", FieldTwo: 123}}},
			keys:        []Key{Field("one"), Field("two"), Field("FieldOne")},
			expected:    "test",
			expectedErr: nil,
		},

		"mixed access level 3 with map with int32 key": {
			input:       map[string]any{"one": map[int32]itemOne{42: {FieldOne: "test", FieldTwo: 123}}},
			keys:        []Key{Field("one"), Field("42"), Index(0)},
			expected:    "test",
			expectedErr: nil,
		},

		"mixed access level 3 with map with int32 key and index fields": {
			input:       map[string]any{"one": map[int32]itemOne{42: {FieldOne: "test", FieldTwo: 123}}},
			keys:        []Key{Field("one"), Index(42), Index(0)},
			expected:    "test",
			expectedErr: nil,
		},
	}

	dt := DefaultTraverser{
		caster: cast.NewCaster(),
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got, err := dt.Retrieve(tc.input, tc.keys)

			// check error
			testingx.AssertError(t, tc.expectedErr, err)

			// check returned item
			if !reflect.DeepEqual(tc.expected, got) {
				t.Errorf("expected %#v got %#v", tc.expected, got)
			}
		})
	}
}

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
		caster: cast.NewCaster(),
	}

	for name, tc := range tests {
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = dt.Retrieve(tc.input, tc.keys)
			}
		})
	}
}
