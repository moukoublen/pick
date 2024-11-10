package pick

import (
	"errors"
	"fmt"
	"strconv"
)

var (
	ErrCastOverFlow      = errors.New("overflow error")
	ErrCastLostDecimals  = errors.New("missing decimals error")
	ErrCastInvalidType   = errors.New("invalid type")
	ErrCastInvalidSyntax = errors.New("invalid syntax")
)

type Error struct {
	inner         error
	originalValue any
}

func (c *Error) Error() string {
	innerStr := ""
	if c.inner != nil {
		innerStr = c.inner.Error()
	}

	return fmt.Sprintf("cast error, original value %T(%v), error: %s", c.originalValue, c.originalValue, innerStr)
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
		inner = ErrCastInvalidSyntax
	}

	return &Error{
		inner:         inner,
		originalValue: originalValue,
	}
}
