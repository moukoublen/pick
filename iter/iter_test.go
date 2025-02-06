package iter

import (
	"errors"
	"fmt"
	"testing"

	"github.com/moukoublen/pick/internal/tst"
	"github.com/stretchr/testify/mock"
)

type expectedOpCall[M CollectionOpMeta | FieldOpMeta] struct {
	Item        any
	ReturnError error
	Meta        M
}

type MockOp[M CollectionOpMeta | FieldOpMeta] struct {
	mock.Mock
}

func (m *MockOp[M]) Operation(item any, meta M) error {
	return m.Called(item, meta).Error(0)
}

func (m *MockOp[M]) init(ex []expectedOpCall[M], ordered bool) {
	sl := make([]*mock.Call, 0, len(ex))
	for _, e := range ex {
		c := m.On("Operation", e.Item, e.Meta).Return(e.ReturnError)
		sl = append(sl, c)
	}

	if ordered {
		mock.InOrder(sl...)
	}
}

func generateExpectedCalls[Input any](input []Input) []expectedOpCall[CollectionOpMeta] {
	e := make([]expectedOpCall[CollectionOpMeta], 0, len(input))

	for i, n := range input {
		e = append(e, expectedOpCall[CollectionOpMeta]{
			Meta:        CollectionOpMeta{Index: i, Length: len(input)},
			Item:        n,
			ReturnError: nil,
		})
	}

	return e
}

func TestIterForEachField(t *testing.T) {
	tests := map[string]struct {
		Input         any
		ErrorAsserter tst.ErrorAsserter
		ExpectedCalls []expectedOpCall[FieldOpMeta]
	}{
		"nil": {
			Input:         nil,
			ErrorAsserter: tst.NoError,
			ExpectedCalls: []expectedOpCall[FieldOpMeta]{},
		},
		"map[string]any": {
			Input:         map[string]any{"one": 1, "two": 2},
			ErrorAsserter: tst.NoError,
			ExpectedCalls: []expectedOpCall[FieldOpMeta]{
				{
					Meta:        FieldOpMeta{Field: "one", Length: 2},
					Item:        1,
					ReturnError: nil,
				},
				{
					Meta:        FieldOpMeta{Field: "two", Length: 2},
					Item:        2,
					ReturnError: nil,
				},
			},
		},
		"string": {
			Input:         "string",
			ErrorAsserter: tst.ExpectedErrorIs(ErrNoFields),
			ExpectedCalls: []expectedOpCall[FieldOpMeta]{},
		},
		"map[string]string": {
			Input:         map[string]any{"one": "1", "two": "2"},
			ErrorAsserter: tst.NoError,
			ExpectedCalls: []expectedOpCall[FieldOpMeta]{
				{
					Meta:        FieldOpMeta{Field: "one", Length: 2},
					Item:        "1",
					ReturnError: nil,
				},
				{
					Meta:        FieldOpMeta{Field: "two", Length: 2},
					Item:        "2",
					ReturnError: nil,
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			m := &MockOp[FieldOpMeta]{}
			m.Test(t)
			m.init(tc.ExpectedCalls, false)
			gotErr := ForEachField(tc.Input, m.Operation)
			tc.ErrorAsserter(t, gotErr)
			m.AssertExpectations(t)
		})
	}
}

func TestIterMapErrorScenarios(t *testing.T) {
	errMock1 := errors.New("mock error")

	type testCase struct {
		input                 any
		inputSingleItemCastFn func(any) (int, error)
		errorAsserter         tst.ErrorAsserter
	}

	testsCases := []testCase{
		{
			input:                 []any{1, 2, 3},
			inputSingleItemCastFn: func(any) (int, error) { return 0, errMock1 },
			errorAsserter:         tst.ExpectedErrorIs(errMock1),
		},
		{
			input:                 []any{1, 2, 3},
			inputSingleItemCastFn: func(any) (int, error) { panic("panic") },
			errorAsserter:         tst.ExpectedErrorStringContains(`recovered panic: "panic"`),
		},
	}

	for idx, tc := range testsCases {
		name := fmt.Sprintf("test_%d_(%v)", idx, tc.input)
		t.Run(name, func(t *testing.T) {
			_, gotErr := Map(tc.input, MapOpFn(tc.inputSingleItemCastFn))
			tc.errorAsserter(t, gotErr)
		})
	}
}

func TestIterForEach(t *testing.T) {
	mockError := errors.New("error")

	ptrStr := ptr("test")

	tests := map[string]struct {
		Input         any
		ErrorAsserter tst.ErrorAsserter
		ExpectedCalls []expectedOpCall[CollectionOpMeta]
	}{
		"nil": {
			Input:         nil,
			ErrorAsserter: tst.NoError,
			ExpectedCalls: []expectedOpCall[CollectionOpMeta]{},
		},
		"string": {
			Input:         "abc",
			ErrorAsserter: tst.NoError,
			ExpectedCalls: []expectedOpCall[CollectionOpMeta]{
				{
					Meta:        CollectionOpMeta{Index: 0, Length: 1},
					Item:        "abc",
					ReturnError: nil,
				},
			},
		},
		"string error": {
			Input:         "abc",
			ErrorAsserter: tst.ExpectedErrorIs(mockError),
			ExpectedCalls: []expectedOpCall[CollectionOpMeta]{
				{
					Meta:        CollectionOpMeta{Index: 0, Length: 1},
					Item:        "abc",
					ReturnError: mockError,
				},
			},
		},
		"struct{}": {
			Input:         struct{}{},
			ErrorAsserter: tst.NoError,
			ExpectedCalls: []expectedOpCall[CollectionOpMeta]{
				{
					Meta:        CollectionOpMeta{Index: 0, Length: 1},
					Item:        struct{}{},
					ReturnError: nil,
				},
			},
		},

		"[]any:0": {
			Input:         []any{},
			ErrorAsserter: tst.NoError,
			ExpectedCalls: []expectedOpCall[CollectionOpMeta]{},
		},
		"[]any:8": {
			Input:         []any{1, 2, 3, 4, 5, 6, 7, 8},
			ErrorAsserter: tst.NoError,
			ExpectedCalls: generateExpectedCalls([]any{1, 2, 3, 4, 5, 6, 7, 8}),
		},
		"[]string:8": {
			Input: []string{
				"Named must your fear be before banish it you can.",
				"When you look at the dark side, careful you must be. For the dark side looks back.",
				"abc",
				"abc",
				"abc",
				"abc",
				"abc",
				"abc",
			},
			ErrorAsserter: tst.NoError,
			ExpectedCalls: generateExpectedCalls([]string{
				"Named must your fear be before banish it you can.",
				"When you look at the dark side, careful you must be. For the dark side looks back.",
				"abc",
				"abc",
				"abc",
				"abc",
				"abc",
				"abc",
			}),
		},
		"[]int8:8": {
			Input:         []int8{1, 2, 3, 4, 5, 6, 7, 8},
			ErrorAsserter: tst.NoError,
			ExpectedCalls: generateExpectedCalls([]int8{1, 2, 3, 4, 5, 6, 7, 8}),
		},
		"[]int16:8": {
			Input:         []int16{1, 2, 3, 4, 5, 6, 7, 8},
			ErrorAsserter: tst.NoError,
			ExpectedCalls: generateExpectedCalls([]int16{1, 2, 3, 4, 5, 6, 7, 8}),
		},
		"[]int32:8": {
			Input:         []int32{1, 2, 3, 4, 5, 6, 7, 8},
			ErrorAsserter: tst.NoError,
			ExpectedCalls: generateExpectedCalls([]int32{1, 2, 3, 4, 5, 6, 7, 8}),
		},
		"[]int64:8": {
			Input:         []int64{1, 2, 3, 4, 5, 6, 7, 8},
			ErrorAsserter: tst.NoError,
			ExpectedCalls: generateExpectedCalls([]int64{1, 2, 3, 4, 5, 6, 7, 8}),
		},
		"[]int:8": {
			Input:         []int{1, 2, 3, 4, 5, 6, 7, 8},
			ErrorAsserter: tst.NoError,
			ExpectedCalls: generateExpectedCalls([]int{1, 2, 3, 4, 5, 6, 7, 8}),
		},
		"[]uint8:8": {
			Input:         []uint8{1, 2, 3, 4, 5, 6, 7, 8},
			ErrorAsserter: tst.NoError,
			ExpectedCalls: generateExpectedCalls([]uint8{1, 2, 3, 4, 5, 6, 7, 8}),
		},
		"[]uint16:8": {
			Input:         []uint16{1, 2, 3, 4, 5, 6, 7, 8},
			ErrorAsserter: tst.NoError,
			ExpectedCalls: generateExpectedCalls([]uint16{1, 2, 3, 4, 5, 6, 7, 8}),
		},
		"[]uint32:8": {
			Input:         []uint32{1, 2, 3, 4, 5, 6, 7, 8},
			ErrorAsserter: tst.NoError,
			ExpectedCalls: generateExpectedCalls([]uint32{1, 2, 3, 4, 5, 6, 7, 8}),
		},
		"[]uint64:8": {
			Input:         []uint64{1, 2, 3, 4, 5, 6, 7, 8},
			ErrorAsserter: tst.NoError,
			ExpectedCalls: generateExpectedCalls([]uint64{1, 2, 3, 4, 5, 6, 7, 8}),
		},
		"[]uint:8": {
			Input:         []uint{1, 2, 3, 4, 5, 6, 7, 8},
			ErrorAsserter: tst.NoError,
			ExpectedCalls: generateExpectedCalls([]uint{1, 2, 3, 4, 5, 6, 7, 8}),
		},
		"[]float32:8": {
			Input:         []float32{1, 2, 3, 4, 5, 6, 7, 8},
			ErrorAsserter: tst.NoError,
			ExpectedCalls: generateExpectedCalls([]float32{1, 2, 3, 4, 5, 6, 7, 8}),
		},
		"[]float64:8": {
			Input:         []float64{1, 2, 3, 4, 5, 6, 7, 8},
			ErrorAsserter: tst.NoError,
			ExpectedCalls: generateExpectedCalls([]float64{1, 2, 3, 4, 5, 6, 7, 8}),
		},
		"[]bool:8": {
			Input:         []bool{false, false, false, false, false, false, false, false},
			ErrorAsserter: tst.NoError,
			ExpectedCalls: generateExpectedCalls([]bool{false, false, false, false, false, false, false, false}),
		},
		"[]struct{}:4": {
			Input:         []struct{}{{}, {}, {}, {}},
			ErrorAsserter: tst.NoError,
			ExpectedCalls: generateExpectedCalls([]struct{}{{}, {}, {}, {}}),
		},

		"[8]int8": {
			Input:         [8]int8{1, 2, 3, 4, 5, 6, 7, 8},
			ErrorAsserter: tst.NoError,
			ExpectedCalls: generateExpectedCalls([]int8{1, 2, 3, 4, 5, 6, 7, 8}),
		},

		"[8]int8 error": {
			Input:         [8]int8{1, 2, 3, 4, 5, 6, 7, 8},
			ErrorAsserter: tst.ExpectedErrorIs(mockError),
			ExpectedCalls: []expectedOpCall[CollectionOpMeta]{
				{
					Meta:        CollectionOpMeta{Index: 0, Length: 8},
					Item:        int8(1),
					ReturnError: nil,
				},
				{
					Meta:        CollectionOpMeta{Index: 1, Length: 8},
					Item:        int8(2),
					ReturnError: mockError,
				},
			},
		},

		"*string nil": {
			Input:         (*string)(nil),
			ErrorAsserter: tst.NoError,
			ExpectedCalls: []expectedOpCall[CollectionOpMeta]{},
		},

		"*string not nil": {
			Input:         ptrStr,
			ErrorAsserter: tst.NoError,
			ExpectedCalls: []expectedOpCall[CollectionOpMeta]{
				{
					Meta:        CollectionOpMeta{Index: 0, Length: 1},
					Item:        *ptrStr,
					ReturnError: nil,
				},
			},
		},

		"**string not nil": {
			Input:         &ptrStr,
			ErrorAsserter: tst.NoError,
			ExpectedCalls: []expectedOpCall[CollectionOpMeta]{
				{
					Meta:        CollectionOpMeta{Index: 0, Length: 1},
					Item:        ptrStr,
					ReturnError: nil,
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			m := &MockOp[CollectionOpMeta]{}
			m.Test(t)
			m.init(tc.ExpectedCalls, true)
			gotErr := ForEach(tc.Input, m.Operation)
			tc.ErrorAsserter(t, gotErr)
			m.AssertExpectations(t)
		})
	}
}

func BenchmarkIterForEach(b *testing.B) {
	noop := func(_ any, _ CollectionOpMeta) error { return nil }

	tests := map[string]struct {
		Input any
	}{
		"string": {
			Input: "abc",
		},
		"struct{}": {
			Input: struct{}{},
		},

		"[]string:8": {
			Input: []string{"Named must your fear be before banish it you can.", "When you look at the dark side, careful you must be. For the dark side looks back.", "abc", "abc", "abc", "abc", "abc", "abc"},
		},
		"[]int8:8": {
			Input: []int8{1, 2, 3, 4, 5, 6, 7, 8},
		},
		"[]int16:8": {
			Input: []int16{1, 2, 3, 4, 5, 6, 7, 8},
		},
		"[]int32:8": {
			Input: []int32{1, 2, 3, 4, 5, 6, 7, 8},
		},
		"[]int64:8": {
			Input: []int64{1, 2, 3, 4, 5, 6, 7, 8},
		},
		"[]int:8": {
			Input: []int{1, 2, 3, 4, 5, 6, 7, 8},
		},
		"[]uint8:8": {
			Input: []uint8{1, 2, 3, 4, 5, 6, 7, 8},
		},
		"[]uint16:8": {
			Input: []uint16{1, 2, 3, 4, 5, 6, 7, 8},
		},
		"[]uint32:8": {
			Input: []uint32{1, 2, 3, 4, 5, 6, 7, 8},
		},
		"[]uint64:8": {
			Input: []uint64{1, 2, 3, 4, 5, 6, 7, 8},
		},
		"[]uint:8": {
			Input: []uint{1, 2, 3, 4, 5, 6, 7, 8},
		},
		"[]float32:8": {
			Input: []float32{1, 2, 3, 4, 5, 6, 7, 8},
		},
		"[]float64:8": {
			Input: []float64{1, 2, 3, 4, 5, 6, 7, 8},
		},
		"[]bool:8": {
			Input: []bool{false, false, false, false, false, false, false, false},
		},
		"[]struct{}:4": {
			Input: []struct{}{{}, {}, {}, {}},
		},

		"[8]string": {
			Input: [8]string{"Named must your fear be before banish it you can.", "When you look at the dark side, careful you must be. For the dark side looks back.", "abc", "abc", "abc", "abc", "abc", "abc"},
		},
		"[8]int8": {
			Input: [8]int8{1, 2, 3, 4, 5, 6, 7, 8},
		},
		"[8]int16": {
			Input: [8]int16{1, 2, 3, 4, 5, 6, 7, 8},
		},
		"[8]int32": {
			Input: [8]int32{1, 2, 3, 4, 5, 6, 7, 8},
		},
		"[8]int64": {
			Input: [8]int64{1, 2, 3, 4, 5, 6, 7, 8},
		},
		"[8]int": {
			Input: [8]int{1, 2, 3, 4, 5, 6, 7, 8},
		},
		"[8]uint8": {
			Input: [8]uint8{1, 2, 3, 4, 5, 6, 7, 8},
		},
		"[8]uint16": {
			Input: [8]uint16{1, 2, 3, 4, 5, 6, 7, 8},
		},
		"[8]uint32": {
			Input: [8]uint32{1, 2, 3, 4, 5, 6, 7, 8},
		},
		"[8]uint64": {
			Input: [8]uint64{1, 2, 3, 4, 5, 6, 7, 8},
		},
		"[8]uint": {
			Input: [8]uint{1, 2, 3, 4, 5, 6, 7, 8},
		},
		"[8]float32": {
			Input: [8]float32{1, 2, 3, 4, 5, 6, 7, 8},
		},
		"[8]float64": {
			Input: [8]float64{1, 2, 3, 4, 5, 6, 7, 8},
		},
		"[8]bool": {
			Input: [8]bool{false, false, false, false, false, false, false, false},
		},
		"[8]struct{}": {
			Input: [8]struct{}{{}, {}, {}, {}, {}, {}, {}, {}},
		},
	}

	for name, tc := range tests {
		b.Run(name, func(b *testing.B) {
			for range b.N {
				_ = ForEach(tc.Input, noop)
			}
		})
	}
}

type avgInterface interface {
	Avg() int
}

type implementsAvgInterface []int

func (s implementsAvgInterface) Avg() int {
	var sum int
	for _, n := range s {
		sum += n
	}

	return sum / len(s)
}

var noLength = tst.ExpectedErrorIs(ErrNoLength)

type (
	sliceIntAlias []int
	arrayIntAlias [5]int
	stringAlias   string
)

var lenTests = map[string]struct {
	Input         any
	ErrorAsserter tst.ErrorAsserter
	Expected      int
}{
	"nil any int nil": {
		Input:         nil,
		ErrorAsserter: noLength,
		Expected:      -1,
	},
	"slice any": {
		Input:         []any{1, 2, "3"},
		ErrorAsserter: tst.NoError,
		Expected:      3,
	},
	"slice map[string]any": {
		Input:         []map[string]any{{}, {}, {}},
		ErrorAsserter: tst.NoError,
		Expected:      3,
	},
	"slice int8": {
		Input:         []int8{1, 2, 3},
		ErrorAsserter: tst.NoError,
		Expected:      3,
	},
	"slice int16": {
		Input:         []int16{1, 2, 3},
		ErrorAsserter: tst.NoError,
		Expected:      3,
	},
	"slice int32": {
		Input:         []int32{1, 2, 3},
		ErrorAsserter: tst.NoError,
		Expected:      3,
	},
	"slice int64": {
		Input:         []int64{1, 2, 3},
		ErrorAsserter: tst.NoError,
		Expected:      3,
	},
	"slice uint": {
		Input:         []uint{1, 2, 3},
		ErrorAsserter: tst.NoError,
		Expected:      3,
	},
	"slice uint8": {
		Input:         []uint8{1, 2, 3},
		ErrorAsserter: tst.NoError,
		Expected:      3,
	},
	"slice uint16": {
		Input:         []uint16{1, 2, 3},
		ErrorAsserter: tst.NoError,
		Expected:      3,
	},
	"slice uint32": {
		Input:         []uint32{1, 2, 3},
		ErrorAsserter: tst.NoError,
		Expected:      3,
	},
	"slice uint64": {
		Input:         []uint64{1, 2, 3},
		ErrorAsserter: tst.NoError,
		Expected:      3,
	},
	"slice float32": {
		Input:         []float32{1, 2, 3},
		ErrorAsserter: tst.NoError,
		Expected:      3,
	},
	"slice float64": {
		Input:         []float32{1, 2, 3},
		ErrorAsserter: tst.NoError,
		Expected:      3,
	},
	"slice bool": {
		Input:         []bool{true, true, false},
		ErrorAsserter: tst.NoError,
		Expected:      3,
	},
	"slice int": {
		Input:         []int{1, 2, 3},
		ErrorAsserter: tst.NoError,
		Expected:      3,
	},
	"slice int nil": {
		Input:         []int(nil),
		ErrorAsserter: tst.NoError,
		Expected:      0,
	},
	"array int 4": {
		Input:         [4]int{1, 2, 3, 4},
		ErrorAsserter: tst.NoError,
		Expected:      4,
	},
	"array int32 3": {
		Input:         [3]int32{1, 2, 3},
		ErrorAsserter: tst.NoError,
		Expected:      3,
	},
	"sliceIntAlias int": {
		Input:         sliceIntAlias{1, 2},
		ErrorAsserter: tst.NoError,
		Expected:      2,
	},
	"sliceIntAlias int nil": {
		Input:         sliceIntAlias(nil),
		ErrorAsserter: tst.NoError,
		Expected:      0,
	},
	"arrayIntAlias int": {
		Input:         arrayIntAlias{1, 2, 3, 4, 5},
		ErrorAsserter: tst.NoError,
		Expected:      5,
	},
	"struct slice": {
		Input:         []struct{}{{}, {}, {}, {}, {}},
		ErrorAsserter: tst.NoError,
		Expected:      5,
	},
	"string": {
		Input:         "abcd",
		ErrorAsserter: tst.NoError,
		Expected:      4,
	},
	"string slice": {
		Input:         []string{"abcd", "abc", "ab", "a"},
		ErrorAsserter: tst.NoError,
		Expected:      4,
	},
	"stringAlias": {
		Input:         stringAlias("abcd"),
		ErrorAsserter: tst.NoError,
		Expected:      4,
	},
	"string pointer": {
		Input:         ptr("test"),
		ErrorAsserter: tst.NoError,
		Expected:      4,
	},
	"string pointer nil": {
		Input:         (*string)(nil),
		ErrorAsserter: noLength,
		Expected:      -1,
	},
	"slice pointer  bool": {
		Input:         []*bool{ptr(true), ptr(true), ptr(true)},
		ErrorAsserter: tst.NoError,
		Expected:      3,
	},
}

func TestLen(t *testing.T) {
	for name, tc := range lenTests {
		t.Run(name, func(t *testing.T) {
			got, gotErr := Len(tc.Input)
			tc.ErrorAsserter(t, gotErr)
			tst.AssertEqual(t, got, tc.Expected)
		})
	}

	t.Run("avgInterface wraps implementsAvgInterface", func(t *testing.T) {
		var a avgInterface = implementsAvgInterface{1, 2, 3, 4, 5, 6, 7}
		func(a avgInterface) {
			got, gotErr := Len(a)
			tst.NoError(t, gotErr)
			tst.AssertEqual(t, got, 7)
		}(a)
	})
}

func BenchmarkLen(b *testing.B) {
	for name, tc := range lenTests {
		b.Run(name, func(b *testing.B) {
			for range b.N {
				_, _ = Len(tc.Input)
			}
		})
	}
}

func ptr[T any](x T) *T { return &x }
