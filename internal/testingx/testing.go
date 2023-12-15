package testingx

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func AssertError(t *testing.T, assertErrFn func(*testing.T, error), err error) {
	t.Helper()
	switch {
	case err != nil && assertErrFn != nil:
		assertErrFn(t, err)
	case err != nil && assertErrFn == nil:
		t.Errorf("unexpected error returned: %s", err.Error())
	case err == nil && assertErrFn != nil:
		t.Errorf("expected error but none received")
	}
}

func ExpectedErrorIs(allExpectedErrors ...error) func(*testing.T, error) {
	return func(t *testing.T, err error) {
		t.Helper()
		for _, expected := range allExpectedErrors {
			if is := errors.Is(err, expected); !is {
				t.Errorf("expected error [%T]{%s} but got: [%T]{%s}", expected, expected.Error(), err, err.Error())
			}
		}
	}
}

func ExpectedErrorStringContains(s string) func(*testing.T, error) {
	return func(t *testing.T, err error) {
		t.Helper()
		if !strings.Contains(err.Error(), s) {
			t.Errorf("error expected to contain \n%s\n but is \n%s\n", s, err.Error())
		}
	}
}

func AssertEqualFn(t *testing.T) func(subject, expected any) {
	t.Helper()
	return func(subject, expected any) {
		t.Helper()
		AssertEqual(t, subject, expected)
	}
}

func AssertEqual(t *testing.T, subject, expected any) {
	t.Helper()
	if expectedErr, is := expected.(error); is {
		gotErr, _ := subject.(error)
		if errors.Is(expectedErr, gotErr) {
			t.Errorf("expected error: %s got: %s", Format(expectedErr), Format(gotErr))
		}
		return
	}

	var compFn func(any, any) bool

	switch expected.(type) {
	case float32:
		compFn = CompareFloat32
	case float64:
		compFn = CompareFloat64
	case []float32:
		compFn = CompareFloat32Slices
	case []float64:
		compFn = CompareFloat64Slices
	default:
		compFn = reflect.DeepEqual
	}

	if !compFn(subject, expected) {
		t.Errorf("Assert error:\nExpected:%s\nGot     : %s", Format(expected), Format(subject))
	}
}

func Format(a any) string {
	var val string
	switch t := a.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		val = fmt.Sprintf("%d", a)
	case string:
		val = t
	case bool:
		val = strconv.FormatBool(t)
	case float32, float64:
		val = fmt.Sprintf("%g", a)
	default:
		val = fmt.Sprintf("%v", a)
	}

	return fmt.Sprintf("%T(%s)", a, val)
}

func CompareFloat64(a, b any) bool {
	fx := a.(float64) //nolint:forcetypeassert
	fy := b.(float64) //nolint:forcetypeassert

	if math.IsInf(fx, 1) && math.IsInf(fy, 1) {
		return true
	}

	if math.IsInf(fx, -1) && math.IsInf(fy, -1) {
		return true
	}

	const thr = float64(1e-10)
	return math.Abs(fx-fy) <= thr
}

func CompareFloat32(a, b any) bool {
	fx := a.(float32) //nolint:forcetypeassert
	fy := b.(float32) //nolint:forcetypeassert

	if math.IsInf(float64(fx), 1) && math.IsInf(float64(fy), 1) {
		return true
	}

	if math.IsInf(float64(fx), -1) && math.IsInf(float64(fy), -1) {
		return true
	}

	const thr = float64(1e-7)
	return math.Abs(float64(fx-fy)) <= thr
}

func CompareFloat64Slices(a, b any) bool {
	fx := a.([]float64) //nolint:forcetypeassert
	fy := b.([]float64) //nolint:forcetypeassert

	if len(fx) != len(fy) {
		return false
	}

	for i := range fx {
		if !CompareFloat64(fx[i], fy[i]) {
			return false
		}
	}

	return true
}

func CompareFloat32Slices(a, b any) bool {
	fx := a.([]float32) //nolint:forcetypeassert
	fy := b.([]float32) //nolint:forcetypeassert

	if len(fx) != len(fy) {
		return false
	}

	for i := range fx {
		if !CompareFloat32(fx[i], fy[i]) {
			return false
		}
	}

	return true
}
