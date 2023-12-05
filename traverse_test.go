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
		keys        []Field
	}{
		"nil": {
			input:       nil,
			keys:        nil,
			expected:    nil,
			expectedErr: nil,
		},

		"index access level 1": {
			input:       []any{"one", "two"},
			keys:        []Field{Index(1)},
			expected:    "two",
			expectedErr: nil,
		},

		"index access level 1 out of range": {
			input:       []any{"one", "two"},
			keys:        []Field{Index(6)},
			expected:    nil,
			expectedErr: testingx.ExpectedErrorIs(ErrIndexOutOfRange),
		},

		"index access slice of string level 1": {
			input:       []string{"one", "two"},
			keys:        []Field{Index(1)},
			expected:    "two",
			expectedErr: nil,
		},

		"name access level 1": {
			input:       map[string]any{"one": "value"},
			keys:        []Field{Name("one")},
			expected:    "value",
			expectedErr: nil,
		},

		"name access level 1 happy path": {
			input:       map[string]any{"one": "value"},
			keys:        []Field{Name("one")},
			expected:    "value",
			expectedErr: nil,
		},

		"name access level 1 not found": {
			input:       map[string]any{"one": "value"},
			keys:        []Field{Name("two")},
			expected:    nil,
			expectedErr: testingx.ExpectedErrorIs(ErrFieldNotFound),
		},

		"name access level 2 happy path": {
			input:       map[string]any{"one": map[string]any{"two": "value"}},
			keys:        []Field{Name("one"), Name("two")},
			expected:    "value",
			expectedErr: nil,
		},

		"name access level 2 renamed happy path": {
			input:       map[string]any{"one": renamed{"two": "value"}},
			keys:        []Field{Name("one"), Name("two")},
			expected:    "value",
			expectedErr: nil,
		},

		"mixed access level 2": {
			input:       []any{"one", map[string]any{"two": "value"}},
			keys:        []Field{Index(1), Name("two")},
			expected:    "value",
			expectedErr: nil,
		},

		"mixed access level 2 with cast": {
			input:       []any{"one", map[string]any{"4": "value"}},
			keys:        []Field{Index(1), Index(4)},
			expected:    "value",
			expectedErr: nil,
		},

		"mixed access level 2 with cast 2": {
			input:       map[string]any{"one": []string{"s0", "s1", "s2"}},
			keys:        []Field{Name("one"), Name("1")},
			expected:    "s1",
			expectedErr: nil,
		},

		"mixed access level 3 with struct": {
			input:       map[string]any{"one": renamed{"two": itemOne{FieldOne: "test", FieldTwo: 123}}},
			keys:        []Field{Name("one"), Name("two"), Name("FieldOne")},
			expected:    "test",
			expectedErr: nil,
		},

		"mixed access level 3 with struct using index": {
			input:       map[string]any{"one": renamed{"two": itemOne{FieldOne: "test", FieldTwo: 123}}},
			keys:        []Field{Name("one"), Name("two"), Index(0)},
			expected:    "test",
			expectedErr: nil,
		},

		"mixed access level 3 with struct using wrong index": {
			input:       map[string]any{"one": renamed{"two": itemOne{FieldOne: "test", FieldTwo: 123}}},
			keys:        []Field{Name("one"), Name("two"), Index(12)},
			expected:    nil,
			expectedErr: testingx.ExpectedErrorStringContains("reflect: Field index out of range"),
		},

		"mixed access level 3 with struct using wrong field": {
			input:       map[string]any{"one": renamed{"two": itemOne{FieldOne: "test", FieldTwo: 123}}},
			keys:        []Field{Name("one"), Name("two"), Name("Wrong")},
			expected:    nil,
			expectedErr: testingx.ExpectedErrorIs(ErrFieldNotFound),
		},

		"mixed access level 3 with pointer struct": {
			input:       map[string]any{"one": renamed{"two": &itemOne{FieldOne: "test", FieldTwo: 123}}},
			keys:        []Field{Name("one"), Name("two"), Name("FieldOne")},
			expected:    "test",
			expectedErr: nil,
		},

		"mixed access level 3 with map with int32 key": {
			input:       map[string]any{"one": map[int32]itemOne{42: {FieldOne: "test", FieldTwo: 123}}},
			keys:        []Field{Name("one"), Name("42"), Index(0)},
			expected:    "test",
			expectedErr: nil,
		},

		"mixed access level 3 with map with int32 key and index fields": {
			input:       map[string]any{"one": map[int32]itemOne{42: {FieldOne: "test", FieldTwo: 123}}},
			keys:        []Field{Name("one"), Index(42), Index(0)},
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
			got, err := dt.Get(tc.input, tc.keys)

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
		keys  []Field
	}{
		"index access level 1": {
			input: []any{"one", "two"},
			keys:  []Field{Index(1)},
		},

		"index access level 1 out of range": {
			input: []any{"one", "two"},
			keys:  []Field{Index(6)},
		},

		"index access level 1 slice of string ": {
			input: []string{"one", "two"},
			keys:  []Field{Index(1)},
		},

		"index access level 2": {
			input: []any{"one", []any{"two", "three"}},
			keys:  []Field{Index(1), Index(1)},
		},

		"index access level 3": {
			input: []any{"one", []any{"two", []any{"three", "four"}}},
			keys:  []Field{Index(1), Index(1), Index(1)},
		},

		"index access level 3 mixed slice of string": {
			input: []any{"one", []any{"two", []string{"three", "four"}}},
			keys:  []Field{Index(1), Index(1), Index(1)},
		},

		"name access level 1": {
			input: map[string]any{"one": "value"},
			keys:  []Field{Name("one")},
		},

		"name access level 1 not found": {
			input: map[string]any{"one": "value"},
			keys:  []Field{Name("two")},
		},

		"name access level 2": {
			input: map[string]any{"one": map[string]any{"two": "value"}},
			keys:  []Field{Name("one"), Name("two")},
		},

		"name access level 2 renamed": {
			input: map[string]any{"one": renamed{"two": "value"}},
			keys:  []Field{Name("one"), Name("two")},
		},

		"name access level 3": {
			input: map[string]any{"one": map[string]any{"two": map[string]any{"three": "value"}}},
			keys:  []Field{Name("one"), Name("two"), Name("three")},
		},

		"name access level 3 renamed": {
			input: map[string]any{"one": renamed{"two": map[string]any{"three": "value"}}},
			keys:  []Field{Name("one"), Name("two")},
		},

		"name access to struct level 1": {
			input: itemOne{FieldOne: "one", FieldTwo: 2},
			keys:  []Field{Name("FieldOne")},
		},

		"mixed access level 2": {
			input: []any{"one", map[string]any{"two": "value"}},
			keys:  []Field{Index(1), Name("two")},
		},

		"mixed access level 2 with key cast index to name": {
			input: []any{"one", map[string]any{"4": "value"}},
			keys:  []Field{Index(1), Index(4)},
		},

		"mixed access level 2 with key cast name to index": {
			input: map[string]any{"one": []string{"s0", "s1", "s2"}},
			keys:  []Field{Name("one"), Name("1")},
		},

		"mixed access level 3": {
			input: []any{"one", map[string]any{"two": []any{"a", "b", "c"}}},
			keys:  []Field{Index(1), Name("two"), Index(2)},
		},

		"mixed access level 3 with struct": {
			input: map[string]any{"one": renamed{"two": itemOne{FieldOne: "test", FieldTwo: 123}}},
			keys:  []Field{Name("one"), Name("two"), Name("FieldOne")},
		},

		"mixed access level 3 with struct using index": {
			input: map[string]any{"one": renamed{"two": itemOne{FieldOne: "test", FieldTwo: 123}}},
			keys:  []Field{Name("one"), Name("two"), Index(0)},
		},

		"mixed access level 3 with struct using wrong index": {
			input: map[string]any{"one": renamed{"two": itemOne{FieldOne: "test", FieldTwo: 123}}},
			keys:  []Field{Name("one"), Name("two"), Index(12)},
		},

		"mixed access level 3 with struct using wrong field": {
			input: map[string]any{"one": renamed{"two": itemOne{FieldOne: "test", FieldTwo: 123}}},
			keys:  []Field{Name("one"), Name("two"), Name("Wrong")},
		},

		"mixed access level 3 with pointer struct": {
			input: map[string]any{"one": renamed{"two": &itemOne{FieldOne: "test", FieldTwo: 123}}},
			keys:  []Field{Name("one"), Name("two"), Name("FieldOne")},
		},

		"mixed access level 3 with key cast name to int32 and struct index field": {
			input: map[string]any{"one": map[int32]itemOne{42: {FieldOne: "test", FieldTwo: 123}}},
			keys:  []Field{Name("one"), Name("42"), Index(0)},
		},

		"mixed access level 3 with key cast int to int32 and struct index field": {
			input: map[string]any{"one": map[int32]itemOne{42: {FieldOne: "test", FieldTwo: 123}}},
			keys:  []Field{Name("one"), Index(42), Index(0)},
		},
	}

	dt := DefaultTraverser{
		caster: cast.NewCaster(),
	}

	for name, tc := range tests {
		tc := tc
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = dt.Get(tc.input, tc.keys)
			}
		})
	}
}
