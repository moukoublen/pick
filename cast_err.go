package pick

import (
	"errors"
	"fmt"
	"strconv"
)

var (
	ErrConvertOverFlow      = errors.New("overflow error")
	ErrConvertLostDecimals  = errors.New("missing decimals error")
	ErrConvertInvalidType   = errors.New("invalid type")
	ErrConvertInvalidSyntax = errors.New("invalid syntax")
)

type ConvertError struct {
	inner         error
	originalValue any
}

func (c *ConvertError) Error() string {
	innerStr := ""
	if c.inner != nil {
		innerStr = c.inner.Error()
	}

	return fmt.Sprintf("convert error, original value %T(%v), error: %s", c.originalValue, c.originalValue, innerStr)
}

func (c *ConvertError) Unwrap() error {
	return c.inner
}

func newConvertError(inner error, originalValue any) *ConvertError {
	if errors.Is(inner, strconv.ErrRange) {
		inner = ErrConvertOverFlow
	} else if errors.Is(inner, strconv.ErrSyntax) {
		inner = ErrConvertInvalidSyntax
	}

	return &ConvertError{
		inner:         inner,
		originalValue: originalValue,
	}
}
