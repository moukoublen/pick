package slices

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/moukoublen/pick/internal/testingx"
)

func TestToSliceErrorScenarios(t *testing.T) {
	t.Parallel()

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
		tc := tc
		name := fmt.Sprintf("test_%d_(%v)", idx, tc.input)
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			_, gotErr := AsSlice(tc.input, CastOpFn(tc.inputSingleItemCastFn))
			testingx.AssertError(t, tc.expectedErr, gotErr)
		})
	}
}

type expectedOpCall struct {
	Meta        OpMeta
	Item        any
	ReturnError error
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

func TestForEach(t *testing.T) {
	t.Parallel()

	mockOp := func(t *testing.T, expectedCalls []expectedOpCall) Op {
		t.Helper()
		idx := 0
		return func(item any, meta OpMeta) error {
			exp := expectedCalls[idx]

			if !reflect.DeepEqual(item, exp.Item) {
				t.Errorf("expected operation argument mismatch. At %d call, expected %T(%#v) got %T(%#v)", idx, exp.Item, exp.Item, item, item)
			}

			if meta != exp.Meta {
				t.Errorf("expected operation metadata argument mismatch. At %d call, expected %#v got %#v", idx, exp.Meta, meta)
			}

			idx++
			return exp.ReturnError
		}
	}

	mockError := errors.New("error")

	tests := map[string]struct {
		Input         any
		ExpectedErr   func(*testing.T, error)
		ExpectedCalls []expectedOpCall
	}{
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
			Input:       [8]int8{1, 2, 3, 4, 5, 6, 7, 8},
			ExpectedErr: nil,
			ExpectedCalls: []expectedOpCall{
				{
					Meta:        OpMeta{Index: 0, Length: 8},
					Item:        int8(1),
					ReturnError: nil,
				},
				{
					Meta:        OpMeta{Index: 1, Length: 8},
					Item:        int8(2),
					ReturnError: nil,
				},
				{
					Meta:        OpMeta{Index: 2, Length: 8},
					Item:        int8(3),
					ReturnError: nil,
				},
				{
					Meta:        OpMeta{Index: 3, Length: 8},
					Item:        int8(4),
					ReturnError: nil,
				},
				{
					Meta:        OpMeta{Index: 4, Length: 8},
					Item:        int8(5),
					ReturnError: nil,
				},
				{
					Meta:        OpMeta{Index: 5, Length: 8},
					Item:        int8(6),
					ReturnError: nil,
				},
				{
					Meta:        OpMeta{Index: 6, Length: 8},
					Item:        int8(7),
					ReturnError: nil,
				},
				{
					Meta:        OpMeta{Index: 7, Length: 8},
					Item:        int8(8),
					ReturnError: nil,
				},
			},
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
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			gotErr := ForEach(tc.Input, mockOp(t, tc.ExpectedCalls))
			testingx.AssertError(t, tc.ExpectedErr, gotErr)
		})
	}
}

func BenchmarkForEach(b *testing.B) {
	var noop Op = func(_ any, _ OpMeta) error { return nil }

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
		tc := tc
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = ForEach(tc.Input, noop)
			}
		})
	}
}
