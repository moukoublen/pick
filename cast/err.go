package cast

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

var (
	ErrCastOverFlow     = errors.New("overflow error")
	ErrCastLostDecimals = errors.New("missing decimals error")
	ErrInvalidType      = errors.New("invalid type")
	ErrInvalidSyntax    = errors.New("invalid syntax")
)

type Error struct {
	inner         error
	originalValue any
	kind          reflect.Kind
}

func (c *Error) Error() string {
	innerStr := ""
	if c.inner != nil {
		innerStr = c.inner.Error()
	}

	return fmt.Sprintf("cast error, original value %s(%v), error: %s", c.kind.String(), c.originalValue, innerStr)
}

func (c *Error) Unwrap() error {
	return c.inner
}

func (c *Error) Is(e error) bool {
	//nolint:errorlint // non need to unwrap here.
	_, is := e.(*Error)

	return is
}

func newCastError(inner error, originalValue any) *Error {
	if errors.Is(inner, strconv.ErrRange) {
		inner = ErrCastOverFlow
	} else if errors.Is(inner, strconv.ErrSyntax) {
		inner = ErrInvalidSyntax
	}

	kind := reflect.Invalid
	if originalValue != nil {
		t := reflect.TypeOf(originalValue)
		kind = t.Kind()
	}

	return &Error{
		inner:         inner,
		originalValue: originalValue,
		kind:          kind,
	}
}
