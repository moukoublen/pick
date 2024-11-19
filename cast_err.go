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

type CastError struct {
	inner         error
	originalValue any
}

func (c *CastError) Error() string {
	innerStr := ""
	if c.inner != nil {
		innerStr = c.inner.Error()
	}

	return fmt.Sprintf("cast error, original value %T(%v), error: %s", c.originalValue, c.originalValue, innerStr)
}

func (c *CastError) Unwrap() error {
	return c.inner
}

func newCastError(inner error, originalValue any) *CastError {
	if errors.Is(inner, strconv.ErrRange) {
		inner = ErrCastOverFlow
	} else if errors.Is(inner, strconv.ErrSyntax) {
		inner = ErrCastInvalidSyntax
	}

	return &CastError{
		inner:         inner,
		originalValue: originalValue,
	}
}
