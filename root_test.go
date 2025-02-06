package pick

import (
	"testing"

	"github.com/moukoublen/pick/internal/tst"
)

func TestOrDefault(t *testing.T) {
	type stringAlias string

	tests := map[string]struct {
		data any
		call func(*Picker) (any, error)

		expectedValue any
		errorAsserter tst.ErrorAsserter
	}{
		"exists - no cast": {
			data: map[string]any{"one": "value"},
			call: func(p *Picker) (any, error) {
				v, err := OrDefault(p, "one", "default")
				return v, err
			},
			expectedValue: "value",
			errorAsserter: tst.NoError,
		},
		"not exists - return default": {
			data: map[string]any{"one": "value"},
			call: func(p *Picker) (any, error) {
				v, err := OrDefault(p, "two", "default")
				return v, err
			},
			expectedValue: "default",
			errorAsserter: tst.NoError,
		},
		"exists - with cast": {
			data: map[string]any{"one": 123},
			call: func(p *Picker) (any, error) {
				v, err := OrDefault(p, "one", "default")
				return v, err
			},
			expectedValue: "123",
			errorAsserter: tst.NoError,
		},
		"exists - with cast to alias": {
			data: map[string]any{"one": "value"},
			call: func(p *Picker) (any, error) {
				v, err := OrDefault[stringAlias](p, "one", stringAlias("default"))
				return v, err
			},
			expectedValue: stringAlias("value"),
			errorAsserter: tst.NoError,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			p := Wrap(tc.data)
			got, gotErr := tc.call(p)
			tc.errorAsserter(t, gotErr)
			tst.AssertEqual(t, got, tc.expectedValue)
		})
	}
}

func TestMustOrDefault(t *testing.T) {
	type stringAlias string

	tests := map[string]struct {
		data any
		call func(SelectorMustAPI) any

		expectedValue any
	}{
		"exists - no cast": {
			data: map[string]any{"one": "value"},
			call: func(a SelectorMustAPI) any {
				return MustOrDefault(a, "one", "default")
			},
			expectedValue: "value",
		},
		"not exists - default": {
			data: map[string]any{"one": "value"},
			call: func(a SelectorMustAPI) any {
				return MustOrDefault(a, "two", "default")
			},
			expectedValue: "default",
		},
		"exists - cast": {
			data: map[string]any{"one": 123},
			call: func(a SelectorMustAPI) any {
				return MustOrDefault(a, "one", "default")
			},
			expectedValue: "123",
		},
		"exists - with cast to alias": {
			data: map[string]any{"one": "value"},
			call: func(a SelectorMustAPI) any {
				return MustOrDefault[stringAlias](a, "one", stringAlias("default"))
			},
			expectedValue: stringAlias("value"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			p := Wrap(tc.data)
			got := tc.call(p.Must())
			tst.AssertEqual(t, got, tc.expectedValue)
		})
	}
}

func TestGet(t *testing.T) {
	type stringAlias string

	tests := map[string]struct {
		data any
		call func(*Picker) (any, error)

		expectedValue any
		errorAsserter tst.ErrorAsserter
	}{
		"exists - no cast": {
			data: map[string]any{"one": "value"},
			call: func(p *Picker) (any, error) {
				v, err := Get[string](p, "one")
				return v, err
			},
			expectedValue: "value",
			errorAsserter: tst.NoError,
		},
		"not exists": {
			data: map[string]any{"one": "value"},
			call: func(p *Picker) (any, error) {
				v, err := Get[string](p, "two")
				return v, err
			},
			expectedValue: "",
			errorAsserter: tst.ExpectedErrorIs(ErrFieldNotFound),
		},
		"exists - with cast": {
			data: map[string]any{"one": 123},
			call: func(p *Picker) (any, error) {
				v, err := Get[string](p, "one")
				return v, err
			},
			expectedValue: "123",
			errorAsserter: tst.NoError,
		},
		"exists - with cast to alias": {
			data: map[string]any{"one": "value"},
			call: func(p *Picker) (any, error) {
				v, err := Get[stringAlias](p, "one")
				return v, err
			},
			expectedValue: stringAlias("value"),
			errorAsserter: tst.NoError,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			p := Wrap(tc.data)
			got, gotErr := tc.call(p)
			tc.errorAsserter(t, gotErr)
			tst.AssertEqual(t, got, tc.expectedValue)
		})
	}
}

func TestMustGet(t *testing.T) {
	type stringAlias string

	tests := map[string]struct {
		data          any
		call          func(SelectorMustAPI) any
		expectedValue any
	}{
		"exists - no cast": {
			data: map[string]any{"one": "value"},
			call: func(a SelectorMustAPI) any {
				return MustGet[string](a, "one")
			},
			expectedValue: "value",
		},
		"not exists": {
			data: map[string]any{"one": "value"},
			call: func(a SelectorMustAPI) any {
				return MustGet[string](a, "two")
			},
			expectedValue: "",
		},
		"exists - cast": {
			data: map[string]any{"one": 123},
			call: func(a SelectorMustAPI) any {
				return MustGet[string](a, "one")
			},
			expectedValue: "123",
		},
		"not exists with type alias": {
			data: map[string]any{"one": "value"},
			call: func(a SelectorMustAPI) any {
				return MustGet[stringAlias](a, "one")
			},
			expectedValue: stringAlias("value"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			p := Wrap(tc.data)
			got := tc.call(p.Must())
			tst.AssertEqual(t, got, tc.expectedValue)
		})
	}
}
