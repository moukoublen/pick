package testingx

import (
	"errors"
	"reflect"
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
				expectedType := reflect.TypeOf(expected)
				gotType := reflect.TypeOf(err)
				t.Errorf("expected error [%s] but got: [%s] with error: [%s]", expectedType.String(), gotType.String(), err.Error())
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
