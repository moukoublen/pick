package pick

import (
	"reflect"
	"testing"

	"github.com/moukoublen/pick/cast"
	"github.com/moukoublen/pick/internal/testingx"
)

func TestDefaultTraverser(t *testing.T) {
	t.Parallel()

	// type mapAlias map[string]any
	name := NameSelectorKey
	index := IndexSelectorKey

	type renamed map[string]any
	type itemOne struct {
		FieldOne string
		FieldTwo int
	}

	tests := map[string]struct {
		input         any
		expected      any
		expectedErr   func(*testing.T, error)
		keys          []SelectorKey
		expectedFound bool
	}{
		"nil": {
			input:         nil,
			keys:          nil,
			expected:      nil,
			expectedFound: false,
			expectedErr:   nil,
		},

		"index access level 1": {
			input:         []any{"one", "two"},
			keys:          []SelectorKey{index(1)},
			expected:      "two",
			expectedFound: true,
			expectedErr:   nil,
		},

		"index access level 1 out of range": {
			input:         []any{"one", "two"},
			keys:          []SelectorKey{index(6)},
			expected:      nil,
			expectedFound: false,
			expectedErr:   testingx.ExpectedErrorIs(ErrIndexOutOfRange),
		},

		"index access slice of string level 1": {
			input:         []string{"one", "two"},
			keys:          []SelectorKey{index(1)},
			expected:      "two",
			expectedFound: true,
			expectedErr:   nil,
		},

		"name access level 1": {
			input:         map[string]any{"one": "value"},
			keys:          []SelectorKey{name("one")},
			expected:      "value",
			expectedFound: true,
			expectedErr:   nil,
		},

		"name access level 1 happy path": {
			input:         map[string]any{"one": "value"},
			keys:          []SelectorKey{name("one")},
			expected:      "value",
			expectedFound: true,
			expectedErr:   nil,
		},

		"name access level 1 not found": {
			input:         map[string]any{"one": "value"},
			keys:          []SelectorKey{name("two")},
			expected:      nil,
			expectedFound: false,
			expectedErr:   nil,
		},

		"name access level 2 happy path": {
			input:         map[string]any{"one": map[string]any{"two": "value"}},
			keys:          []SelectorKey{name("one"), name("two")},
			expected:      "value",
			expectedFound: true,
			expectedErr:   nil,
		},

		"name access level 2 renamed happy path": {
			input:         map[string]any{"one": renamed{"two": "value"}},
			keys:          []SelectorKey{name("one"), name("two")},
			expected:      "value",
			expectedFound: true,
			expectedErr:   nil,
		},

		"mixed access level 2": {
			input:         []any{"one", map[string]any{"two": "value"}},
			keys:          []SelectorKey{index(1), name("two")},
			expected:      "value",
			expectedFound: true,
			expectedErr:   nil,
		},

		"mixed access level 2 with cast": {
			input:         []any{"one", map[string]any{"4": "value"}},
			keys:          []SelectorKey{index(1), index(4)},
			expected:      "value",
			expectedFound: true,
			expectedErr:   nil,
		},

		"mixed access level 2 with cast 2": {
			input:         map[string]any{"one": []string{"s0", "s1", "s2"}},
			keys:          []SelectorKey{name("one"), name("1")},
			expected:      "s1",
			expectedFound: true,
			expectedErr:   nil,
		},

		"mixed access level 3 with struct": {
			input:         map[string]any{"one": renamed{"two": itemOne{FieldOne: "test", FieldTwo: 123}}},
			keys:          []SelectorKey{name("one"), name("two"), name("FieldOne")},
			expected:      "test",
			expectedFound: true,
			expectedErr:   nil,
		},

		"mixed access level 3 with struct using index": {
			input:         map[string]any{"one": renamed{"two": itemOne{FieldOne: "test", FieldTwo: 123}}},
			keys:          []SelectorKey{name("one"), name("two"), index(0)},
			expected:      "test",
			expectedFound: true,
			expectedErr:   nil,
		},

		"mixed access level 3 with struct using wrong index": {
			input:         map[string]any{"one": renamed{"two": itemOne{FieldOne: "test", FieldTwo: 123}}},
			keys:          []SelectorKey{name("one"), name("two"), index(12)},
			expected:      nil,
			expectedFound: false,
			expectedErr:   testingx.ExpectedErrorStringContains("reflect: Field index out of range"),
		},

		"mixed access level 3 with struct using wrong field": {
			input:         map[string]any{"one": renamed{"two": itemOne{FieldOne: "test", FieldTwo: 123}}},
			keys:          []SelectorKey{name("one"), name("two"), name("Wrong")},
			expected:      nil,
			expectedFound: false,
			expectedErr:   nil,
		},

		"mixed access level 3 with pointer struct": {
			input:         map[string]any{"one": renamed{"two": &itemOne{FieldOne: "test", FieldTwo: 123}}},
			keys:          []SelectorKey{name("one"), name("two"), name("FieldOne")},
			expected:      "test",
			expectedFound: true,
			expectedErr:   nil,
		},

		"mixed access level 3 with map with int32 key": {
			input:         map[string]any{"one": map[int32]itemOne{42: {FieldOne: "test", FieldTwo: 123}}},
			keys:          []SelectorKey{name("one"), name("42"), index(0)},
			expected:      "test",
			expectedFound: true,
			expectedErr:   nil,
		},

		"mixed access level 3 with map with int32 key and index selector": {
			input:         map[string]any{"one": map[int32]itemOne{42: {FieldOne: "test", FieldTwo: 123}}},
			keys:          []SelectorKey{name("one"), index(42), index(0)},
			expected:      "test",
			expectedFound: true,
			expectedErr:   nil,
		},
	}

	dt := DefaultTraverser{
		caster:            cast.NewCaster(),
		selectorFormatter: defaultSelectorFormatter{},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got, found, err := dt.Get(tc.input, tc.keys)

			// check error
			testingx.AssertError(t, tc.expectedErr, err)

			// check found bool
			if found != tc.expectedFound {
				t.Errorf("expected found bool to be %v got %v", tc.expectedFound, found)
			}

			// check returned item
			if !reflect.DeepEqual(tc.expected, got) {
				t.Errorf("expected %#v got %#v", tc.expected, got)
			}
		})
	}
}

func BenchmarkDefaultTraverser(b *testing.B) {
	// type mapAlias map[string]any
	n := NameSelectorKey
	i := IndexSelectorKey

	type renamed map[string]any
	type itemOne struct {
		FieldOne string
		FieldTwo int
	}

	tests := map[string]struct {
		input any
		keys  []SelectorKey
	}{
		"index access level 1": {
			input: []any{"one", "two"},
			keys:  []SelectorKey{i(1)},
		},

		"index access level 1 out of range": {
			input: []any{"one", "two"},
			keys:  []SelectorKey{i(6)},
		},

		"index access level 1 slice of string ": {
			input: []string{"one", "two"},
			keys:  []SelectorKey{i(1)},
		},

		"index access level 2": {
			input: []any{"one", []any{"two", "three"}},
			keys:  []SelectorKey{i(1), i(1)},
		},

		"index access level 3": {
			input: []any{"one", []any{"two", []any{"three", "four"}}},
			keys:  []SelectorKey{i(1), i(1), i(1)},
		},

		"index access level 3 mixed slice of string": {
			input: []any{"one", []any{"two", []string{"three", "four"}}},
			keys:  []SelectorKey{i(1), i(1), i(1)},
		},

		"name access level 1": {
			input: map[string]any{"one": "value"},
			keys:  []SelectorKey{n("one")},
		},

		"name access level 1 not found": {
			input: map[string]any{"one": "value"},
			keys:  []SelectorKey{n("two")},
		},

		"name access level 2": {
			input: map[string]any{"one": map[string]any{"two": "value"}},
			keys:  []SelectorKey{n("one"), n("two")},
		},

		"name access level 2 renamed": {
			input: map[string]any{"one": renamed{"two": "value"}},
			keys:  []SelectorKey{n("one"), n("two")},
		},

		"name access level 3": {
			input: map[string]any{"one": map[string]any{"two": map[string]any{"three": "value"}}},
			keys:  []SelectorKey{n("one"), n("two"), n("three")},
		},

		"name access level 3 renamed": {
			input: map[string]any{"one": renamed{"two": map[string]any{"three": "value"}}},
			keys:  []SelectorKey{n("one"), n("two")},
		},

		"name access to struct level 1": {
			input: itemOne{FieldOne: "one", FieldTwo: 2},
			keys:  []SelectorKey{n("FieldOne")},
		},

		"mixed access level 2": {
			input: []any{"one", map[string]any{"two": "value"}},
			keys:  []SelectorKey{i(1), n("two")},
		},

		"mixed access level 2 with key cast index to name": {
			input: []any{"one", map[string]any{"4": "value"}},
			keys:  []SelectorKey{i(1), i(4)},
		},

		"mixed access level 2 with key cast name to index": {
			input: map[string]any{"one": []string{"s0", "s1", "s2"}},
			keys:  []SelectorKey{n("one"), n("1")},
		},

		"mixed access level 3": {
			input: []any{"one", map[string]any{"two": []any{"a", "b", "c"}}},
			keys:  []SelectorKey{i(1), n("two"), i(2)},
		},

		"mixed access level 3 with struct": {
			input: map[string]any{"one": renamed{"two": itemOne{FieldOne: "test", FieldTwo: 123}}},
			keys:  []SelectorKey{n("one"), n("two"), n("FieldOne")},
		},

		"mixed access level 3 with struct using index": {
			input: map[string]any{"one": renamed{"two": itemOne{FieldOne: "test", FieldTwo: 123}}},
			keys:  []SelectorKey{n("one"), n("two"), i(0)},
		},

		"mixed access level 3 with struct using wrong index": {
			input: map[string]any{"one": renamed{"two": itemOne{FieldOne: "test", FieldTwo: 123}}},
			keys:  []SelectorKey{n("one"), n("two"), i(12)},
		},

		"mixed access level 3 with struct using wrong field": {
			input: map[string]any{"one": renamed{"two": itemOne{FieldOne: "test", FieldTwo: 123}}},
			keys:  []SelectorKey{n("one"), n("two"), n("Wrong")},
		},

		"mixed access level 3 with pointer struct": {
			input: map[string]any{"one": renamed{"two": &itemOne{FieldOne: "test", FieldTwo: 123}}},
			keys:  []SelectorKey{n("one"), n("two"), n("FieldOne")},
		},

		"mixed access level 3 with key cast name to int32 and struct index selector": {
			input: map[string]any{"one": map[int32]itemOne{42: {FieldOne: "test", FieldTwo: 123}}},
			keys:  []SelectorKey{n("one"), n("42"), i(0)},
		},

		"mixed access level 3 with key cast int to int32 and struct index selector": {
			input: map[string]any{"one": map[int32]itemOne{42: {FieldOne: "test", FieldTwo: 123}}},
			keys:  []SelectorKey{n("one"), i(42), i(0)},
		},
	}

	dt := DefaultTraverser{
		caster:            cast.NewCaster(),
		selectorFormatter: defaultSelectorFormatter{},
	}

	for name, tc := range tests {
		tc := tc
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _, _ = dt.Get(tc.input, tc.keys)
			}
		})
	}
}
