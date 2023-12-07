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

		"index access level 1": {
			input:       []any{"one", "two"},
			keys:        []Key{Index(1)},
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
		tc := tc
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
		"index access level 1": {
			input: []any{"one", "two"},
			keys:  []Key{Index(1)},
		},

		"index access level 1 out of range": {
			input: []any{"one", "two"},
			keys:  []Key{Index(6)},
		},

		"index access level 1 slice of string ": {
			input: []string{"one", "two"},
			keys:  []Key{Index(1)},
		},

		"index access level 2": {
			input: []any{"one", []any{"two", "three"}},
			keys:  []Key{Index(1), Index(1)},
		},

		"index access level 3": {
			input: []any{"one", []any{"two", []any{"three", "four"}}},
			keys:  []Key{Index(1), Index(1), Index(1)},
		},

		"index access level 3 mixed slice of string": {
			input: []any{"one", []any{"two", []string{"three", "four"}}},
			keys:  []Key{Index(1), Index(1), Index(1)},
		},

		"name access level 1": {
			input: map[string]any{"one": "value"},
			keys:  []Key{Field("one")},
		},

		"name access level 1 not found": {
			input: map[string]any{"one": "value"},
			keys:  []Key{Field("two")},
		},

		"name access level 2": {
			input: map[string]any{"one": map[string]any{"two": "value"}},
			keys:  []Key{Field("one"), Field("two")},
		},

		"name access level 2 renamed": {
			input: map[string]any{"one": renamed{"two": "value"}},
			keys:  []Key{Field("one"), Field("two")},
		},

		"name access level 3": {
			input: map[string]any{"one": map[string]any{"two": map[string]any{"three": "value"}}},
			keys:  []Key{Field("one"), Field("two"), Field("three")},
		},

		"name access level 3 renamed": {
			input: map[string]any{"one": renamed{"two": map[string]any{"three": "value"}}},
			keys:  []Key{Field("one"), Field("two")},
		},

		"name access to struct level 1": {
			input: itemOne{FieldOne: "one", FieldTwo: 2},
			keys:  []Key{Field("FieldOne")},
		},

		"mixed access level 2": {
			input: []any{"one", map[string]any{"two": "value"}},
			keys:  []Key{Index(1), Field("two")},
		},

		"mixed access level 2 with key cast index to name": {
			input: []any{"one", map[string]any{"4": "value"}},
			keys:  []Key{Index(1), Index(4)},
		},

		"mixed access level 2 with key cast name to index": {
			input: map[string]any{"one": []string{"s0", "s1", "s2"}},
			keys:  []Key{Field("one"), Field("1")},
		},

		"mixed access level 3": {
			input: []any{"one", map[string]any{"two": []any{"a", "b", "c"}}},
			keys:  []Key{Index(1), Field("two"), Index(2)},
		},

		"mixed access level 3 with struct": {
			input: map[string]any{"one": renamed{"two": itemOne{FieldOne: "test", FieldTwo: 123}}},
			keys:  []Key{Field("one"), Field("two"), Field("FieldOne")},
		},

		"mixed access level 3 with struct using index": {
			input: map[string]any{"one": renamed{"two": itemOne{FieldOne: "test", FieldTwo: 123}}},
			keys:  []Key{Field("one"), Field("two"), Index(0)},
		},

		"mixed access level 3 with struct using wrong index": {
			input: map[string]any{"one": renamed{"two": itemOne{FieldOne: "test", FieldTwo: 123}}},
			keys:  []Key{Field("one"), Field("two"), Index(12)},
		},

		"mixed access level 3 with struct using wrong field": {
			input: map[string]any{"one": renamed{"two": itemOne{FieldOne: "test", FieldTwo: 123}}},
			keys:  []Key{Field("one"), Field("two"), Field("Wrong")},
		},

		"mixed access level 3 with pointer struct": {
			input: map[string]any{"one": renamed{"two": &itemOne{FieldOne: "test", FieldTwo: 123}}},
			keys:  []Key{Field("one"), Field("two"), Field("FieldOne")},
		},

		"mixed access level 3 with key cast name to int32 and struct index field": {
			input: map[string]any{"one": map[int32]itemOne{42: {FieldOne: "test", FieldTwo: 123}}},
			keys:  []Key{Field("one"), Field("42"), Index(0)},
		},

		"mixed access level 3 with key cast int to int32 and struct index field": {
			input: map[string]any{"one": map[int32]itemOne{42: {FieldOne: "test", FieldTwo: 123}}},
			keys:  []Key{Field("one"), Index(42), Index(0)},
		},
	}

	dt := DefaultTraverser{
		caster: cast.NewCaster(),
	}

	for name, tc := range tests {
		tc := tc
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = dt.Retrieve(tc.input, tc.keys)
			}
		})
	}
}
