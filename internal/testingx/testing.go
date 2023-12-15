package testingx

import (
	"errors"
	"fmt"
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
	if !reflect.DeepEqual(subject, expected) {
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
