package iter

import (
	"errors"
	"fmt"
	"testing"

	"github.com/moukoublen/pick/internal/testingx"
)

func TestIterMapErrorScenarios(t *testing.T) {
	errMock1 := errors.New("mock error")

	type testCase struct {
		input                 any
		inputSingleItemCastFn func(any) (int, error)
		expectedErr           func(*testing.T, error)
	}

	testsCases := []testCase{
		{
			input:                 []any{1, 2, 3},
			inputSingleItemCastFn: func(any) (int, error) { return 0, errMock1 },
			expectedErr:           testingx.ExpectedErrorIs(errMock1),
		},
		{
			input:                 []any{1, 2, 3},
			inputSingleItemCastFn: func(any) (int, error) { panic("panic") },
			expectedErr:           testingx.ExpectedErrorStringContains(`recovered panic: "panic"`),
		},
	}

	for idx, tc := range testsCases {
		name := fmt.Sprintf("test_%d_(%v)", idx, tc.input)
		t.Run(name, func(t *testing.T) {
			_, gotErr := Map(tc.input, MapOpFn(tc.inputSingleItemCastFn))
			testingx.AssertError(t, tc.expectedErr, gotErr)
		})
	}
}

type expectedOpCall struct {
	Item        any
	ReturnError error
	Meta        OpMeta
}

func generateExpectedCalls[T any](input []T) []expectedOpCall {
	e := make([]expectedOpCall, 0, len(input))

	for i, n := range input {
		e = append(e, expectedOpCall{
			Meta:        OpMeta{Index: i, Length: len(input)},
			Item:        n,
			ReturnError: nil,
		})
	}

	return e
}

func TestIterForEach(t *testing.T) {
	// t.Parallel()

	mockOp := func(t *testing.T, expectedCalls []expectedOpCall) func(item any, meta OpMeta) error {
		t.Helper()
		idx := 0
		t.Cleanup(func() {
			if idx != len(expectedCalls) {
				t.Errorf("mockOp not all expected calls were performed. Expected %d calls made %d", len(expectedCalls), idx)
			}
		})
		return func(item any, meta OpMeta) error {
			t.Helper()
			exp := expectedCalls[idx]

			testingx.AssertEqual(t, item, exp.Item)
			testingx.AssertEqual(t, meta, exp.Meta)

			idx++
			return exp.ReturnError
		}
	}

	mockError := errors.New("error")

	ptrStr := ptr("test")

	tests := map[string]struct {
		Input         any
		ExpectedErr   func(*testing.T, error)
		ExpectedCalls []expectedOpCall
	}{
		"nil": {
			Input:         nil,
			ExpectedErr:   nil,
			ExpectedCalls: []expectedOpCall{},
		},
		"string": {
			Input:       "abc",
			ExpectedErr: nil,
			ExpectedCalls: []expectedOpCall{
				{
					Meta:        OpMeta{Index: 0, Length: 1},
					Item:        "abc",
					ReturnError: nil,
				},
			},
		},
		"string error": {
			Input:       "abc",
			ExpectedErr: testingx.ExpectedErrorIs(mockError),
			ExpectedCalls: []expectedOpCall{
				{
					Meta:        OpMeta{Index: 0, Length: 1},
					Item:        "abc",
					ReturnError: mockError,
				},
			},
		},
		"struct{}": {
			Input:       struct{}{},
			ExpectedErr: nil,
			ExpectedCalls: []expectedOpCall{
				{
					Meta:        OpMeta{Index: 0, Length: 1},
					Item:        struct{}{},
					ReturnError: nil,
				},
			},
		},

		"[]any:0": {
			Input:         []any{},
			ExpectedErr:   nil,
			ExpectedCalls: []expectedOpCall{},
		},
		"[]any:8": {
			Input:         []any{1, 2, 3, 4, 5, 6, 7, 8},
			ExpectedErr:   nil,
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
			ExpectedErr: nil,
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
			ExpectedErr:   nil,
			ExpectedCalls: generateExpectedCalls([]int8{1, 2, 3, 4, 5, 6, 7, 8}),
		},
		"[]int16:8": {
			Input:         []int16{1, 2, 3, 4, 5, 6, 7, 8},
			ExpectedErr:   nil,
			ExpectedCalls: generateExpectedCalls([]int16{1, 2, 3, 4, 5, 6, 7, 8}),
		},
		"[]int32:8": {
			Input:         []int32{1, 2, 3, 4, 5, 6, 7, 8},
			ExpectedErr:   nil,
			ExpectedCalls: generateExpectedCalls([]int32{1, 2, 3, 4, 5, 6, 7, 8}),
		},
		"[]int64:8": {
			Input:         []int64{1, 2, 3, 4, 5, 6, 7, 8},
			ExpectedErr:   nil,
			ExpectedCalls: generateExpectedCalls([]int64{1, 2, 3, 4, 5, 6, 7, 8}),
		},
		"[]int:8": {
			Input:         []int{1, 2, 3, 4, 5, 6, 7, 8},
			ExpectedErr:   nil,
			ExpectedCalls: generateExpectedCalls([]int{1, 2, 3, 4, 5, 6, 7, 8}),
		},
		"[]uint8:8": {
			Input:         []uint8{1, 2, 3, 4, 5, 6, 7, 8},
			ExpectedErr:   nil,
			ExpectedCalls: generateExpectedCalls([]uint8{1, 2, 3, 4, 5, 6, 7, 8}),
		},
		"[]uint16:8": {
			Input:         []uint16{1, 2, 3, 4, 5, 6, 7, 8},
			ExpectedErr:   nil,
			ExpectedCalls: generateExpectedCalls([]uint16{1, 2, 3, 4, 5, 6, 7, 8}),
		},
		"[]uint32:8": {
			Input:         []uint32{1, 2, 3, 4, 5, 6, 7, 8},
			ExpectedErr:   nil,
			ExpectedCalls: generateExpectedCalls([]uint32{1, 2, 3, 4, 5, 6, 7, 8}),
		},
		"[]uint64:8": {
			Input:         []uint64{1, 2, 3, 4, 5, 6, 7, 8},
			ExpectedErr:   nil,
			ExpectedCalls: generateExpectedCalls([]uint64{1, 2, 3, 4, 5, 6, 7, 8}),
		},
		"[]uint:8": {
			Input:         []uint{1, 2, 3, 4, 5, 6, 7, 8},
			ExpectedErr:   nil,
			ExpectedCalls: generateExpectedCalls([]uint{1, 2, 3, 4, 5, 6, 7, 8}),
		},
		"[]float32:8": {
			Input:         []float32{1, 2, 3, 4, 5, 6, 7, 8},
			ExpectedErr:   nil,
			ExpectedCalls: generateExpectedCalls([]float32{1, 2, 3, 4, 5, 6, 7, 8}),
		},
		"[]float64:8": {
			Input:         []float64{1, 2, 3, 4, 5, 6, 7, 8},
			ExpectedErr:   nil,
			ExpectedCalls: generateExpectedCalls([]float64{1, 2, 3, 4, 5, 6, 7, 8}),
		},
		"[]bool:8": {
			Input:         []bool{false, false, false, false, false, false, false, false},
			ExpectedErr:   nil,
			ExpectedCalls: generateExpectedCalls([]bool{false, false, false, false, false, false, false, false}),
		},
		"[]struct{}:4": {
			Input:         []struct{}{{}, {}, {}, {}},
			ExpectedErr:   nil,
			ExpectedCalls: generateExpectedCalls([]struct{}{{}, {}, {}, {}}),
		},

		"[8]int8": {
			Input:         [8]int8{1, 2, 3, 4, 5, 6, 7, 8},
			ExpectedErr:   nil,
			ExpectedCalls: generateExpectedCalls([]int8{1, 2, 3, 4, 5, 6, 7, 8}),
		},

		"[8]int8 error": {
			Input:       [8]int8{1, 2, 3, 4, 5, 6, 7, 8},
			ExpectedErr: testingx.ExpectedErrorIs(mockError),
			ExpectedCalls: []expectedOpCall{
				{
					Meta:        OpMeta{Index: 0, Length: 8},
					Item:        int8(1),
					ReturnError: nil,
				},
				{
					Meta:        OpMeta{Index: 1, Length: 8},
					Item:        int8(2),
					ReturnError: mockError,
				},
			},
		},

		"*string nil": {
			Input:         (*string)(nil),
			ExpectedErr:   nil,
			ExpectedCalls: []expectedOpCall{},
		},

		"*string not nil": {
			Input:       ptrStr,
			ExpectedErr: nil,
			ExpectedCalls: []expectedOpCall{
				{
					Meta:        OpMeta{Index: 0, Length: 1},
					Item:        *ptrStr,
					ReturnError: nil,
				},
			},
		},

		"**string not nil": {
			Input:       &ptrStr,
			ExpectedErr: nil,
			ExpectedCalls: []expectedOpCall{
				{
					Meta:        OpMeta{Index: 0, Length: 1},
					Item:        ptrStr,
					ReturnError: nil,
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			gotErr := ForEach(tc.Input, mockOp(t, tc.ExpectedCalls))
			testingx.AssertError(t, tc.ExpectedErr, gotErr)
		})
	}
}

func BenchmarkIterForEach(b *testing.B) {
	noop := func(_ any, _ OpMeta) error { return nil }

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

var noLength = testingx.ExpectedErrorIs(ErrNoLength)

type (
	sliceIntAlias []int
	arrayIntAlias [5]int
	stringAlias   string
)

var lenTests = map[string]struct {
	Input       any
	ExpectedErr func(*testing.T, error)
	Expected    int
}{
	"nil any int nil": {
		Input:       nil,
		ExpectedErr: noLength,
		Expected:    -1,
	},
	"slice any": {
		Input:       []any{1, 2, "3"},
		ExpectedErr: nil,
		Expected:    3,
	},
	"slice map[string]any": {
		Input:       []map[string]any{{}, {}, {}},
		ExpectedErr: nil,
		Expected:    3,
	},
	"slice int8": {
		Input:       []int8{1, 2, 3},
		ExpectedErr: nil,
		Expected:    3,
	},
	"slice int16": {
		Input:       []int16{1, 2, 3},
		ExpectedErr: nil,
		Expected:    3,
	},
	"slice int32": {
		Input:       []int32{1, 2, 3},
		ExpectedErr: nil,
		Expected:    3,
	},
	"slice int64": {
		Input:       []int64{1, 2, 3},
		ExpectedErr: nil,
		Expected:    3,
	},
	"slice uint": {
		Input:       []uint{1, 2, 3},
		ExpectedErr: nil,
		Expected:    3,
	},
	"slice uint8": {
		Input:       []uint8{1, 2, 3},
		ExpectedErr: nil,
		Expected:    3,
	},
	"slice uint16": {
		Input:       []uint16{1, 2, 3},
		ExpectedErr: nil,
		Expected:    3,
	},
	"slice uint32": {
		Input:       []uint32{1, 2, 3},
		ExpectedErr: nil,
		Expected:    3,
	},
	"slice uint64": {
		Input:       []uint64{1, 2, 3},
		ExpectedErr: nil,
		Expected:    3,
	},
	"slice float32": {
		Input:       []float32{1, 2, 3},
		ExpectedErr: nil,
		Expected:    3,
	},
	"slice float64": {
		Input:       []float32{1, 2, 3},
		ExpectedErr: nil,
		Expected:    3,
	},
	"slice bool": {
		Input:       []bool{true, true, false},
		ExpectedErr: nil,
		Expected:    3,
	},
	"slice int": {
		Input:       []int{1, 2, 3},
		ExpectedErr: nil,
		Expected:    3,
	},
	"slice int nil": {
		Input:       []int(nil),
		ExpectedErr: nil,
		Expected:    0,
	},
	"array int 4": {
		Input:       [4]int{1, 2, 3, 4},
		ExpectedErr: nil,
		Expected:    4,
	},
	"array int32 3": {
		Input:       [3]int32{1, 2, 3},
		ExpectedErr: nil,
		Expected:    3,
	},
	"sliceIntAlias int": {
		Input:       sliceIntAlias{1, 2},
		ExpectedErr: nil,
		Expected:    2,
	},
	"sliceIntAlias int nil": {
		Input:       sliceIntAlias(nil),
		ExpectedErr: nil,
		Expected:    0,
	},
	"arrayIntAlias int": {
		Input:       arrayIntAlias{1, 2, 3, 4, 5},
		ExpectedErr: nil,
		Expected:    5,
	},
	"struct slice": {
		Input:       []struct{}{{}, {}, {}, {}, {}},
		ExpectedErr: nil,
		Expected:    5,
	},
	"string": {
		Input:       "abcd",
		ExpectedErr: nil,
		Expected:    4,
	},
	"string slice": {
		Input:       []string{"abcd", "abc", "ab", "a"},
		ExpectedErr: nil,
		Expected:    4,
	},
	"stringAlias": {
		Input:       stringAlias("abcd"),
		ExpectedErr: nil,
		Expected:    4,
	},
	"string pointer": {
		Input:       ptr("test"),
		ExpectedErr: nil,
		Expected:    4,
	},
	"string pointer nil": {
		Input:       (*string)(nil),
		ExpectedErr: noLength,
		Expected:    -1,
	},
	"slice pointer  bool": {
		Input:       []*bool{ptr(true), ptr(true), ptr(true)},
		ExpectedErr: nil,
		Expected:    3,
	},
}

func TestLen(t *testing.T) {
	for name, tc := range lenTests {
		t.Run(name, func(t *testing.T) {
			got, gotErr := Len(tc.Input)
			testingx.AssertError(t, tc.ExpectedErr, gotErr)
			testingx.AssertEqual(t, got, tc.Expected)
		})
	}

	t.Run("avgInterface wraps implementsAvgInterface", func(t *testing.T) {
		var a avgInterface = implementsAvgInterface{1, 2, 3, 4, 5, 6, 7}
		func(a avgInterface) {
			got, gotErr := Len(a)
			testingx.AssertError(t, nil, gotErr)
			testingx.AssertEqual(t, got, 7)
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
