package testingx

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"
)

// TestValueExpected GOOS - GOARCH - Expected Value.
type TestValueExpected[T any] map[string]map[string]T

func Expected[T any]() *TestValueExpected[T] {
	t := TestValueExpected[T](map[string]map[string]T{})
	return &t
}

func (tv *TestValueExpected[T]) Any(val T) *TestValueExpected[T] {
	return tv.At("*", "*", val)
}

func (tv *TestValueExpected[T]) Arch(goarch string, val T) *TestValueExpected[T] {
	return tv.At("*", goarch, val)
}

func (tv *TestValueExpected[T]) At(goos, goarch string, val T) *TestValueExpected[T] {
	atOs, osExists := (*tv)[goos]
	if !osExists || atOs == nil {
		atOs = map[string]T{}
	}

	atOs[goarch] = val
	(*tv)[goos] = atOs

	return tv
}

func (tv *TestValueExpected[T]) Get() T {
	atOs, osExists := (*tv)[runtime.GOOS]
	if !osExists || atOs == nil {
		atOs, osExists = (*tv)["*"]
		if !osExists || atOs == nil {
			panic(fmt.Sprintf("expected value is not set for os (%s %s)", runtime.GOOS, runtime.GOARCH))
		}
	}

	v, exists := atOs[runtime.GOARCH]
	if !exists {
		v, exists = atOs["*"]
		if !exists {
			panic(fmt.Sprintf("expected value is not set for arch (%s %s)", runtime.GOOS, runtime.GOARCH))
		}
	}

	return v
}

func AssertError(t *testing.T, assertErrFn func(*testing.T, error), err error) {
	t.Helper()
	switch {
	case err != nil && assertErrFn != nil:
		assertErrFn(t, err)
	case err != nil && assertErrFn == nil:
		t.Errorf("unexpected error returned.\nError: %T(%s)", err, err.Error())
	case err == nil && assertErrFn != nil:
		t.Errorf("expected error but none received")
	}
}

func ExpectedErrorChecks(expected ...func(*testing.T, error)) func(*testing.T, error) {
	return func(t *testing.T, err error) {
		t.Helper()
		for _, fn := range expected {
			fn(t, err)
		}
	}
}

func ExpectedErrorIs(allExpectedErrors ...error) func(*testing.T, error) {
	return func(t *testing.T, err error) {
		t.Helper()
		for _, expected := range allExpectedErrors {
			if !errors.Is(err, expected) {
				t.Errorf("error unexpected.\nExpected error: %T(%s) \nGot           : %T(%s)", expected, expected.Error(), err, err.Error())
			}
		}
	}
}

func ExpectedErrorOfType[T error](assertsOfType ...func(*testing.T, T)) func(*testing.T, error) {
	return func(t *testing.T, err error) {
		t.Helper()

		var wantErr T
		if !errors.As(err, &wantErr) {
			var tErr T
			t.Errorf("Error type check failed.\nExpected error type: %T\nGot                : %T(%s)", tErr, err, err)
		} else {
			for _, e := range assertsOfType {
				e(t, wantErr)
			}
		}
	}
}

func ExpectedErrorStringContains(s string) func(*testing.T, error) {
	return func(t *testing.T, err error) {
		t.Helper()
		if !strings.Contains(err.Error(), s) {
			t.Errorf("error string check failed. \nExpected to contain: %s\nGot                : %s\n", s, err.Error())
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
			t.Errorf("Expected error mismatch:\nExpected: %s\nGot     : %s", Format(expectedErr), Format(gotErr))
		}
		return
	}

	if reflect.TypeOf(subject) != reflect.TypeOf(expected) {
		t.Errorf("Expected type mismatch:\nExpected: %s\nGot     : %s", Format(expected), Format(subject))
		return
	}

	compFn := compareFn(expected)
	if !compFn(subject, expected) {
		t.Errorf("Expected value mismatch:\nExpected: %s\nGot     : %s", Format(expected), Format(subject))
	}
}

func compareFn(expected any) func(any, any) bool {
	switch expected.(type) {
	case float32:
		return compareFloat32
	case float64:
		return compareFloat64
	case time.Time:
		return compareTime
	case []float32:
		return compareSlicesFn[float32](compareFloat32)
	case []float64:
		return compareSlicesFn[float64](compareFloat64)
	case []time.Time:
		return compareSlicesFn[time.Time](compareTime)
	default:
		return reflect.DeepEqual
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
	case time.Time:
		val = t.Format(time.RFC3339Nano)
	case []int, []int8, []int16, []int32, []int64, []uint, []uint8, []uint16, []uint32, []uint64:
		val = formatSlice(t, Format)
	case []float32, []float64:
		val = formatSlice(t, Format)
	case []bool:
		val = formatSlice(t, Format)
	case []string:
		val = formatSlice(t, Format)
	case []any:
		val = formatSlice(t, Format)
	default:
		val = fmt.Sprintf("%#v", a)
	}

	return fmt.Sprintf("%T(%s)", a, val)
}

func formatSlice(sl any, elementFormatFn func(any) string) string {
	s := strings.Builder{}
	s.WriteRune('[')

	value := reflect.ValueOf(sl)
	for i := range value.Len() {
		item := value.Index(i)
		ifc := item.Interface()
		if i != 0 {
			s.WriteRune(',')
		}
		s.WriteString(elementFormatFn(ifc))
	}

	s.WriteRune(']')
	return s.String()
}

func compareFloat64(a, b any) bool {
	var (
		fx      float64
		fy      float64
		isFloat bool
	)
	fx, isFloat = a.(float64)
	if !isFloat {
		return false
	}
	fy, isFloat = b.(float64)
	if !isFloat {
		return false
	}

	if math.IsInf(fx, 1) && math.IsInf(fy, 1) {
		return true
	}

	if math.IsInf(fx, -1) && math.IsInf(fy, -1) {
		return true
	}

	const thr = float64(1e-10)
	return math.Abs(fx-fy) <= thr
}

func compareFloat32(a, b any) bool {
	var (
		fx      float32
		fy      float32
		isFloat bool
	)
	fx, isFloat = a.(float32)
	if !isFloat {
		return false
	}
	fy, isFloat = b.(float32)
	if !isFloat {
		return false
	}

	return compareFloat64(float64(fx), float64(fy))
}

func compareTime(x, y any) bool {
	var t1, t2 time.Time
	var isTime bool
	t1, isTime = x.(time.Time)
	if !isTime {
		return false
	}
	t2, isTime = y.(time.Time)
	if !isTime {
		return false
	}

	s1 := t1.Format(time.RFC3339Nano)
	s2 := t2.Format(time.RFC3339Nano)
	return s1 == s2 // t1.Equal(t2)
}

func compareSlicesFn[T any](compareFn func(any, any) bool) func(x, y any) bool {
	return func(x, y any) bool {
		var (
			sx, sy  []T
			isSlice bool
		)
		sx, isSlice = x.([]T)
		if !isSlice {
			return false
		}
		sy, isSlice = y.([]T)
		if !isSlice {
			return false
		}

		if len(sx) != len(sy) {
			return false
		}

		for i := range sx {
			if !compareFn(sx[i], sy[i]) {
				return false
			}
		}

		return true
	}
}
