package tst

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

type ErrorAsserter func(*testing.T, error)

func errorExists(t *testing.T, err error) bool {
	t.Helper()
	if err == nil {
		t.Errorf("expected error but none received")
		return false
	}

	return true
}

func NoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("nil error expected but received : %T(%s)", err, err.Error())
	}
}

func ExpectedErrorChecks(expected ...func(*testing.T, error)) ErrorAsserter {
	return func(t *testing.T, err error) {
		t.Helper()
		if !errorExists(t, err) {
			return
		}
		for _, fn := range expected {
			fn(t, err)
		}
	}
}

func ExpectedErrorIs(allExpectedErrors ...error) ErrorAsserter {
	return func(t *testing.T, err error) {
		t.Helper()
		if !errorExists(t, err) {
			return
		}
		for _, expected := range allExpectedErrors {
			if !errors.Is(err, expected) {
				s := stringBuilder{}
				s.WriteStringf("error unexpected.\n")
				s.WriteStringf("Expected error: %T(%s)\n", expected, expected.Error())
				s.WriteStringf("Got           : %T(%s)", err, err.Error())
				t.Error(s.String())
			}
		}
	}
}

func ExpectedErrorOfType[T error](assertsOnType ...func(*testing.T, T)) ErrorAsserter {
	return func(t *testing.T, err error) {
		t.Helper()
		if !errorExists(t, err) {
			return
		}
		var wantErr T
		if !errors.As(err, &wantErr) {
			var tErr T
			s := stringBuilder{}
			s.WriteStringf("error type check failed.\n")
			s.WriteStringf("Expected error type: %T\n", tErr)
			s.WriteStringf("Got                : %T(%s)", err, err.Error())
			t.Error(s.String())
		} else {
			for _, e := range assertsOnType {
				e(t, wantErr)
			}
		}
	}
}

func ExpectedErrorStringContains(str string) ErrorAsserter {
	return func(t *testing.T, err error) {
		t.Helper()
		if !errorExists(t, err) {
			return
		}
		if !strings.Contains(err.Error(), str) {
			s := stringBuilder{}
			s.WriteStringf("error string check failed.\n")
			s.WriteStringf("Expected to contain: %s\n", str)
			s.WriteStringf("Got                : %T(%s)", err, err.Error())
			t.Error(s.String())
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
			s := stringBuilder{}
			s.WriteStringf("Expected error mismatch:\n")
			s.WriteStringf("Expected: %s\n", Format(expectedErr))
			s.WriteStringf("Got     : %s\n", Format(gotErr))
			t.Error(s.String())
		}
		return
	}

	if reflect.TypeOf(subject) != reflect.TypeOf(expected) {
		s := stringBuilder{}
		s.WriteStringf("Expected type mismatch:\n")
		s.WriteStringf("Expected: %s\n", Format(expected))
		s.WriteStringf("Got     : %s\n", Format(subject))
		t.Error(s.String())
		return
	}

	if !cmp.Equal(subject, expected, comparerOptions...) {
		diff := cmp.Diff(subject, expected, comparerOptions...)
		s := stringBuilder{}
		s.WriteStringf("Expected value mismatch:\n")
		s.WriteStringf("Expected: %s\n", Format(expected))
		s.WriteStringf("Got     : %s\n", Format(subject))
		s.WriteStringf("Diff    : %s", diff)
		t.Error(s.String())
	}
}

func Compare(a, b any) bool {
	return cmp.Equal(a, b, comparerOptions...)
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

//nolint:gochecknoglobals
var comparerOptions = []cmp.Option{
	cmp.Comparer(compareFloat64),
	cmp.Comparer(compareFloat32),
	cmp.Comparer(compareTime),
}

func compareFloat64(fx, fy float64) bool {
	if math.IsInf(fx, 1) && math.IsInf(fy, 1) {
		return true
	}

	if math.IsInf(fx, -1) && math.IsInf(fy, -1) {
		return true
	}

	const thr = float64(1e-10)
	return math.Abs(fx-fy) <= thr
}

func compareFloat32(fx, fy float32) bool {
	return compareFloat64(float64(fx), float64(fy))
}

func compareTime(t1, t2 time.Time) bool {
	s1 := t1.Format(time.RFC3339Nano)
	s2 := t2.Format(time.RFC3339Nano)
	return s1 == s2 // t1.Equal(t2)
}

type stringBuilder struct {
	strings.Builder
}

func (s *stringBuilder) WriteStringf(format string, a ...any) {
	fmt.Fprintf(s, format, a...)
}
