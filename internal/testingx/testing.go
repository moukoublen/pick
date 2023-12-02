package testingx

import (
	"errors"
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
