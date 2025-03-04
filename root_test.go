package pick

import (
	"errors"
	"testing"

	"github.com/moukoublen/pick/internal/tst"
	"github.com/stretchr/testify/mock"
)

func TestOrDefault(t *testing.T) {
	type stringAlias string

	tests := map[string]struct {
		data any
		call func(Picker) (any, error)

		expectedValue any
		errorAsserter tst.ErrorAsserter
	}{
		"exists - no convert": {
			data: map[string]any{"one": "value"},
			call: func(p Picker) (any, error) {
				v, err := OrDefault(p, "one", "default")
				return v, err
			},
			expectedValue: "value",
			errorAsserter: tst.NoError,
		},
		"not exists - return default": {
			data: map[string]any{"one": "value"},
			call: func(p Picker) (any, error) {
				v, err := OrDefault(p, "two", "default")
				return v, err
			},
			expectedValue: "default",
			errorAsserter: tst.NoError,
		},
		"exists - with convert": {
			data: map[string]any{"one": 123},
			call: func(p Picker) (any, error) {
				v, err := OrDefault(p, "one", "default")
				return v, err
			},
			expectedValue: "123",
			errorAsserter: tst.NoError,
		},
		"exists - with convert to alias": {
			data: map[string]any{"one": "value"},
			call: func(p Picker) (any, error) {
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

func TestRelaxedOrDefault(t *testing.T) {
	type stringAlias string

	tests := map[string]struct {
		data any
		call func(RelaxedAPI) any

		expectedValue any
	}{
		"exists - no convert": {
			data: map[string]any{"one": "value"},
			call: func(a RelaxedAPI) any {
				return RelaxedOrDefault(a, "one", "default")
			},
			expectedValue: "value",
		},
		"not exists - default": {
			data: map[string]any{"one": "value"},
			call: func(a RelaxedAPI) any {
				return RelaxedOrDefault(a, "two", "default")
			},
			expectedValue: "default",
		},
		"exists - convert": {
			data: map[string]any{"one": 123},
			call: func(a RelaxedAPI) any {
				return RelaxedOrDefault(a, "one", "default")
			},
			expectedValue: "123",
		},
		"exists - with convert to alias": {
			data: map[string]any{"one": "value"},
			call: func(a RelaxedAPI) any {
				return RelaxedOrDefault[stringAlias](a, "one", stringAlias("default"))
			},
			expectedValue: stringAlias("value"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			p := Wrap(tc.data)
			got := tc.call(p.Relaxed())
			tst.AssertEqual(t, got, tc.expectedValue)
		})
	}
}

func TestGet(t *testing.T) {
	type stringAlias string

	tests := map[string]struct {
		data any
		call func(Picker) (any, error)

		expectedValue any
		errorAsserter tst.ErrorAsserter
	}{
		"exists - no convert": {
			data: map[string]any{"one": "value"},
			call: func(p Picker) (any, error) {
				v, err := Get[string](p, "one")
				return v, err
			},
			expectedValue: "value",
			errorAsserter: tst.NoError,
		},
		"not exists": {
			data: map[string]any{"one": "value"},
			call: func(p Picker) (any, error) {
				v, err := Get[string](p, "two")
				return v, err
			},
			expectedValue: "",
			errorAsserter: tst.ExpectedErrorIs(ErrFieldNotFound),
		},
		"exists - with convert": {
			data: map[string]any{"one": 123},
			call: func(p Picker) (any, error) {
				v, err := Get[string](p, "one")
				return v, err
			},
			expectedValue: "123",
			errorAsserter: tst.NoError,
		},
		"exists - with convert to alias": {
			data: map[string]any{"one": "value"},
			call: func(p Picker) (any, error) {
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

func TestRelaxedGet(t *testing.T) {
	type stringAlias string

	tests := map[string]struct {
		data          any
		call          func(RelaxedAPI) any
		expectedValue any
	}{
		"exists - no convert": {
			data: map[string]any{"one": "value"},
			call: func(a RelaxedAPI) any {
				return RelaxedGet[string](a, "one")
			},
			expectedValue: "value",
		},
		"not exists": {
			data: map[string]any{"one": "value"},
			call: func(a RelaxedAPI) any {
				return RelaxedGet[string](a, "two")
			},
			expectedValue: "",
		},
		"exists - convert": {
			data: map[string]any{"one": 123},
			call: func(a RelaxedAPI) any {
				return RelaxedGet[string](a, "one")
			},
			expectedValue: "123",
		},
		"not exists with type alias": {
			data: map[string]any{"one": "value"},
			call: func(a RelaxedAPI) any {
				return RelaxedGet[stringAlias](a, "one")
			},
			expectedValue: stringAlias("value"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			p := Wrap(tc.data)
			got := tc.call(p.Relaxed())
			tst.AssertEqual(t, got, tc.expectedValue)
		})
	}
}

func TestEachField(t *testing.T) {
	mockErr := errors.New("error")
	type stringAlias string

	type Foo struct {
		A string
		B int
	}

	type M map[string]any

	tests := map[string]struct {
		data          any
		selector      string
		expectedCalls func(*mockOp)
		errorAsserter tst.ErrorAsserter
	}{
		"map[string]any": {
			data:     map[string]any{"one": "v1", "two": 42},
			selector: "",
			expectedCalls: func(mo *mockOp) {
				mo.
					ExpectOp(2, "one", "v1", nil).
					ExpectOp(2, "two", 42, nil)
			},
			errorAsserter: tst.NoError,
		},

		"map[string]any - inner": {
			data:     map[string]any{"one": map[string]any{"two": map[string]any{"three": 42, "four": "abc"}}},
			selector: "one.two",
			expectedCalls: func(mo *mockOp) {
				mo.
					ExpectOp(2, "three", 42, nil).
					ExpectOp(2, "four", "abc", nil)
			},
			errorAsserter: tst.NoError,
		},

		"map[string]any - alias": {
			data:     M{"one": M{"two": M{"three": 42, "four": "abc"}}},
			selector: "one.two",
			expectedCalls: func(mo *mockOp) {
				mo.
					ExpectOp(2, "three", 42, nil).
					ExpectOp(2, "four", "abc", nil)
			},
			errorAsserter: tst.NoError,
		},

		"Foo struct": {
			data:     Foo{A: "a", B: 42},
			selector: "",
			expectedCalls: func(mo *mockOp) {
				mo.
					ExpectOp(2, "A", "a", nil).
					ExpectOp(2, "B", 42, nil)
			},
			errorAsserter: tst.NoError,
		},

		"map[stringAlias]any": {
			data:     map[stringAlias]any{"first": 10, "second": "abc"},
			selector: "",
			expectedCalls: func(mo *mockOp) {
				mo.
					ExpectOp(2, "first", 10, nil).
					ExpectOp(2, "second", "abc", nil)
			},
			errorAsserter: tst.NoError,
		},

		"map[int]int": {
			data:     map[int]int{100: 10, 200: 20},
			selector: "",
			expectedCalls: func(mo *mockOp) {
				mo.
					ExpectOp(2, "100", 10, nil).
					ExpectOp(2, "200", 20, nil)
			},
			errorAsserter: tst.NoError,
		},

		"map[string]any - error": {
			data: map[string]any{
				"one":   "v1",
				"two":   "v2",
				"three": "v3",
			},
			selector: "",
			expectedCalls: func(mo *mockOp) {
				maybe := func(op, opRelaxed *mock.Call) {
					op.Maybe()
					opRelaxed.Maybe()
				}
				mo.
					ExpectOp(3, "one", "v1", mockErr).
					ExpectOp(3, "two", "v2", nil, maybe).
					ExpectOp(3, "three", "v3", nil, maybe)
			},
			errorAsserter: tst.ExpectedErrorIs(mockErr),
		},
	}

	// func(field string, p Picker, numOfFields int) error
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			p := Wrap(tc.data)
			m := &mockOp{}
			m.Test(t)
			tc.expectedCalls(m)
			err := EachField(p, tc.selector, m.Operation)
			tc.errorAsserter(t, err)

			eg := &ErrorsSink{}
			RelaxedEachField(p.Relaxed(eg), tc.selector, m.OperationRelaxed)
			tc.errorAsserter(t, eg.Outcome())

			m.AssertExpectations(t)
		})
	}
}

type mockOp struct {
	mock.Mock
}

func (m *mockOp) ExpectOp(numOfFields int, field string, item any, returnError error, ext ...func(op, opRelaxed *mock.Call)) *mockOp {
	c1 := m.On("Operation", field, numOfFields, item).Return(returnError)
	c2 := m.On("OperationRelaxed", field, numOfFields, item).Return(returnError)
	for _, f := range ext {
		f(c1, c2)
	}
	return m
}

func (m *mockOp) Operation(field string, item Picker, numOfFields int) error {
	return m.Called(field, numOfFields, item.Data()).Error(0)
}

func (m *mockOp) OperationRelaxed(field string, item RelaxedAPI, numOfFields int) error {
	return m.Called(field, numOfFields, item.Data()).Error(0)
}
