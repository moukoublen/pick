package pick

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type KeyType int

const (
	KeyTypeUnknown KeyType = iota
	KeyTypeField
	KeyTypeIndex
)

func (s KeyType) String() string {
	switch s {
	case KeyTypeUnknown:
		return "unknown"
	case KeyTypeField:
		return "field"
	case KeyTypeIndex:
		return "index"
	default:
		return ""
	}
}

type Key struct {
	Name  string
	Index int
	Type  KeyType
}

func (s Key) IsIndex() bool { return s.Type == KeyTypeIndex }
func (s Key) IsField() bool { return s.Type == KeyTypeField }

func Field(name string) Key {
	return Key{
		Name:  name,
		Index: 0,
		Type:  KeyTypeField,
	}
}

func Index(idx int) Key {
	return Key{
		Name:  "",
		Index: idx,
		Type:  KeyTypeIndex,
	}
}

const (
	fieldSeparator      rune = '.'
	indexSeparatorStart rune = '['
	indexSeparatorEnd   rune = ']'
)

type DotNotation struct {
	dotNotationFormatter
	dotNotationParser
}

type dotNotationFormatter struct{}

func (d dotNotationFormatter) formatKey(k Key) string {
	switch k.Type {
	case KeyTypeIndex:
		return fmt.Sprintf("[%d]", k.Index)
	case KeyTypeField:
		return k.Name
	default:
		return ""
	}
}

func (d dotNotationFormatter) Format(path ...Key) string {
	sb := strings.Builder{}
	for i, c := range path {
		if i > 0 && c.IsField() {
			sb.WriteRune(fieldSeparator)
		}
		sb.WriteString(d.formatKey(c))
	}

	return sb.String()
}

type dotNotationParser struct{}

func (d dotNotationParser) Parse(selector string) ([]Key, error) {
	if len(selector) == 0 {
		return nil, nil
	}

	runeSlice := []rune(selector)

	path := make([]Key, 0, d.estimateCount(runeSlice))

	for idx := 0; idx < len(runeSlice); {
		switch {
		case runeSlice[idx] == indexSeparatorStart:
			s, i, err := d.parseNextIndex(runeSlice, idx)
			if err != nil {
				return nil, err
			}
			path = append(path, s)
			idx = i

		case runeSlice[idx] == fieldSeparator || idx == 0:
			s, i, err := d.parseNextName(runeSlice, idx)
			if err != nil {
				return nil, err
			}
			path = append(path, s)
			idx = i

		default:
			return nil, ErrInvalidFormat
		}
	}

	return path, nil
}

func (d dotNotationParser) estimateCount(rns []rune) int {
	cnt := 0
	for _, r := range rns {
		switch r {
		case fieldSeparator, indexSeparatorStart:
			cnt++
		}
	}
	cnt++

	return cnt
}

func (d dotNotationParser) parseNextName(rns []rune, idx int) (Key, int, error) {
	if idx >= len(rns) {
		return Key{}, len(rns), ErrInvalidFormatForName
	}

	// omit single leading dot if any.
	if rns[idx] == fieldSeparator {
		idx++
		if idx >= len(rns) {
			return Key{}, len(rns), ErrInvalidFormatForName
		}
	}

	var i int
	for i = idx; i < len(rns); i++ {
		if rns[i] == indexSeparatorStart || rns[i] == fieldSeparator {
			break
		}
		if !d.isRuneValidForName(rns[i]) {
			return Key{}, i, ErrInvalidFormatForName
		}
	}

	// this is the case of continues keys e.g. `..` or `.[`.
	if i == idx {
		return Key{}, i, ErrInvalidFormatForName
	}

	return Key{Type: KeyTypeField, Name: string(rns[idx:i])}, i, nil
}

func (d dotNotationParser) isRuneValidForName(r rune) bool {
	return !unicode.IsControl(r) && r != fieldSeparator && r != indexSeparatorStart && r != indexSeparatorEnd
}

func (d dotNotationParser) parseNextIndex(rns []rune, idx int) (Key, int, error) {
	if idx >= len(rns)-2 {
		return Key{}, len(rns), ErrInvalidFormatForIndex
	}

	// should start with '[' character.
	if rns[idx] != indexSeparatorStart {
		return Key{}, idx, ErrInvalidFormatForIndex
	}

	// omit the leading '[' character.
	idx++

	var i int
	for i = idx; i < len(rns); i++ {
		if rns[i] == indexSeparatorEnd {
			break
		}
		if !d.isRuneValidForIndex(rns[i]) {
			return Key{}, i, ErrInvalidFormatForIndex
		}
	}

	if i == idx || rns[i] != indexSeparatorEnd {
		return Key{}, 0, ErrInvalidFormatForIndex
	}

	str := string(rns[idx:i])

	n, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return Key{}, 0, err
	}

	// omit the trailing ']' character.
	i++

	return Key{Type: KeyTypeIndex, Index: int(n)}, i, nil
}

func (d dotNotationParser) isRuneValidForIndex(r rune) bool {
	return unicode.IsNumber(r)
}

var (
	ErrInvalidFormatForName  = errors.New("invalid format for name key")
	ErrInvalidFormatForIndex = errors.New("invalid format for index key")
	ErrInvalidFormat         = errors.New("invalid selector format")
)
