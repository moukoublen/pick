package pick

import (
	"errors"
	"testing"

	"github.com/moukoublen/pick/internal/testingx"
)

func TestMultiError(t *testing.T) {
	errOne := errors.New("one")
	errTwo := errors.New("two")
	errThree := errors.New("three")

	m := &multiError{}
	m.Add(errOne)
	m.Add(errTwo)
	m.Add(errThree)

	testingx.AssertEqual(t, m.Error(), "one | two | three")
	testingx.AssertEqual(t, errors.Is(m, errOne), true)
	testingx.AssertEqual(t, errors.Is(m, errTwo), true)
	testingx.AssertEqual(t, errors.Is(m, errThree), true)
}
