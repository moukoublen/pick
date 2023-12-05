package pick

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type NotationFieldType int

const (
	NotationFieldTypeUnknown NotationFieldType = iota
	NotationFieldTypeName
	NotationFieldTypeIndex
)

func (s NotationFieldType) String() string {
	switch s {
	case NotationFieldTypeUnknown:
		return "unknown"
	case NotationFieldTypeName:
		return "name"
	case NotationFieldTypeIndex:
		return "index"
	default:
		return ""
	}
}

type Field struct {
	Name  string
	Index int
	Type  NotationFieldType
}

func (s Field) IsIndex() bool { return s.Type == NotationFieldTypeIndex }
func (s Field) IsName() bool  { return s.Type == NotationFieldTypeName }

func Name(name string) Field {
	return Field{
		Name:  name,
		Index: 0,
		Type:  NotationFieldTypeName,
	}
}

func Index(idx int) Field {
	return Field{
		Name:  "",
		Index: idx,
		Type:  NotationFieldTypeIndex,
	}
}

const (
	nameSeparator       rune = '.'
	indexSeparatorStart rune = '['
	indexSeparatorEnd   rune = ']'
)

type DotNotation struct {
	dotNotationFormatter
	dotNotationParser
}

type dotNotationFormatter struct{}

func (d dotNotationFormatter) formatField(s Field) string {
	switch s.Type {
	case NotationFieldTypeIndex:
		return fmt.Sprintf("[%d]", s.Index)
	case NotationFieldTypeName:
		return s.Name
	default:
		return ""
	}
}

func (d dotNotationFormatter) Format(s ...Field) string {
	sb := strings.Builder{}
	for i, c := range s {
		if i > 0 && c.IsName() {
			sb.WriteRune(nameSeparator)
		}
		sb.WriteString(d.formatField(c))
	}

	return sb.String()
}

type dotNotationParser struct{}

func (d dotNotationParser) Parse(selector string) ([]Field, error) {
	if len(selector) == 0 {
		return nil, nil
	}

	runeSlice := []rune(selector)

	sl := make([]Field, 0, d.estimateCount(runeSlice))

	for idx := 0; idx < len(runeSlice); {
		switch {
		case runeSlice[idx] == indexSeparatorStart:
			s, i, err := d.parseNextIndex(runeSlice, idx)
			if err != nil {
				return nil, err
			}
			sl = append(sl, s)
			idx = i

		case runeSlice[idx] == nameSeparator || idx == 0:
			s, i, err := d.parseNextName(runeSlice, idx)
			if err != nil {
				return nil, err
			}
			sl = append(sl, s)
			idx = i

		default:
			return nil, ErrInvalidFormat
		}
	}

	return sl, nil
}

func (d dotNotationParser) estimateCount(rns []rune) int {
	cnt := 0
	for _, r := range rns {
		switch r {
		case nameSeparator, indexSeparatorStart:
			cnt++
		}
	}
	cnt++

	return cnt
}

func (d dotNotationParser) parseNextName(rns []rune, idx int) (Field, int, error) {
	if idx >= len(rns) {
		return Field{}, len(rns), ErrInvalidFormatForName
	}

	// omit single leading dot if any.
	if rns[idx] == nameSeparator {
		idx++
		if idx >= len(rns) {
			return Field{}, len(rns), ErrInvalidFormatForName
		}
	}

	var i int
	for i = idx; i < len(rns); i++ {
		if rns[i] == indexSeparatorStart || rns[i] == nameSeparator {
			break
		}
		if !d.isRuneValidForName(rns[i]) {
			return Field{}, i, ErrInvalidFormatForName
		}
	}

	// this is the case of continues keys e.g. `..` or `.[`.
	if i == idx {
		return Field{}, i, ErrInvalidFormatForName
	}

	return Field{Type: NotationFieldTypeName, Name: string(rns[idx:i])}, i, nil
}

func (d dotNotationParser) isRuneValidForName(r rune) bool {
	return !unicode.IsControl(r) && r != nameSeparator && r != indexSeparatorStart && r != indexSeparatorEnd
}

func (d dotNotationParser) parseNextIndex(rns []rune, idx int) (Field, int, error) {
	if idx >= len(rns)-2 {
		return Field{}, len(rns), ErrInvalidFormatForIndex
	}

	// should start with '[' character.
	if rns[idx] != indexSeparatorStart {
		return Field{}, idx, ErrInvalidFormatForIndex
	}

	// omit the leading '[' character.
	idx++

	var i int
	for i = idx; i < len(rns); i++ {
		if rns[i] == indexSeparatorEnd {
			break
		}
		if !d.isRuneValidForIndex(rns[i]) {
			return Field{}, i, ErrInvalidFormatForIndex
		}
	}

	if i == idx || rns[i] != indexSeparatorEnd {
		return Field{}, 0, ErrInvalidFormatForIndex
	}

	str := string(rns[idx:i])

	n, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return Field{}, 0, err
	}

	// omit the trailing ']' character.
	i++

	return Field{Type: NotationFieldTypeIndex, Index: int(n)}, i, nil
}

func (d dotNotationParser) isRuneValidForIndex(r rune) bool {
	return unicode.IsNumber(r)
}

var (
	ErrInvalidFormatForName  = errors.New("invalid format for name field")
	ErrInvalidFormatForIndex = errors.New("invalid format for index field")
	ErrInvalidFormat         = errors.New("invalid selector format")
)
