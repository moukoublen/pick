package pick

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type SelectorKeyType int

const (
	SelectorKeyTypeUnknown SelectorKeyType = iota
	SelectorKeyTypeName
	SelectorKeyTypeIndex
)

func (s SelectorKeyType) String() string {
	switch s {
	case SelectorKeyTypeUnknown:
		return "unknown"
	case SelectorKeyTypeName:
		return "name"
	case SelectorKeyTypeIndex:
		return "index"
	default:
		return ""
	}
}

type SelectorKey struct {
	Name         string
	Index        int
	SelectorType SelectorKeyType
}

func (s SelectorKey) IsIndex() bool { return s.SelectorType == SelectorKeyTypeIndex }
func (s SelectorKey) IsName() bool  { return s.SelectorType == SelectorKeyTypeName }

func NameSelectorKey(name string) SelectorKey {
	return SelectorKey{
		Name:         name,
		Index:        0,
		SelectorType: SelectorKeyTypeName,
	}
}

func IndexSelectorKey(idx int) SelectorKey {
	return SelectorKey{
		Name:         "",
		Index:        idx,
		SelectorType: SelectorKeyTypeIndex,
	}
}

type DefaultSelectorFormat struct {
	defaultSelectorFormatter
	defaultSelectorParser
}

type defaultSelectorFormatter struct{}

func (d defaultSelectorFormatter) formatSelectorKey(s SelectorKey) string {
	switch s.SelectorType {
	case SelectorKeyTypeIndex:
		return fmt.Sprintf("[%d]", s.Index)
	case SelectorKeyTypeName:
		return s.Name
	default:
		return ""
	}
}

func (d defaultSelectorFormatter) Format(s []SelectorKey) string {
	sb := strings.Builder{}
	for i, c := range s {
		if i > 0 && c.IsName() {
			sb.WriteRune(nameSeparator)
		}
		sb.WriteString(d.formatSelectorKey(c))
	}

	return sb.String()
}

const (
	nameSeparator       rune = '.'
	indexSeparatorStart rune = '['
	indexSeparatorEnd   rune = ']'
)

type defaultSelectorParser struct{}

func (d defaultSelectorParser) isRuneValidForName(r rune) bool {
	return !unicode.IsControl(r) && r != nameSeparator && r != indexSeparatorStart && r != indexSeparatorEnd
}

func (d defaultSelectorParser) isRuneValidForIndex(r rune) bool {
	return unicode.IsNumber(r)
}

func (d defaultSelectorParser) estimateCount(s string) int {
	cnt := 0
	for _, r := range s {
		switch r {
		case nameSeparator, indexSeparatorStart:
			cnt++
		}
	}
	cnt++

	return cnt
}

func (d defaultSelectorParser) parseNextName(rns []rune, idx int) (SelectorKey, int, error) {
	if idx >= len(rns) {
		return SelectorKey{}, len(rns), ErrInvalidFormatForName
	}

	// omit single leading dot if any.
	if rns[idx] == nameSeparator {
		idx++
		if idx >= len(rns) {
			return SelectorKey{}, len(rns), ErrInvalidFormatForName
		}
	}

	var i int
	for i = idx; i < len(rns); i++ {
		if rns[i] == indexSeparatorStart || rns[i] == nameSeparator {
			break
		}
		if !d.isRuneValidForName(rns[i]) {
			return SelectorKey{}, i, ErrInvalidFormatForName
		}
	}

	// this is the case of continues keys e.g. `..` or `.[`.
	if i == idx {
		return SelectorKey{}, i, ErrInvalidFormatForName
	}

	return SelectorKey{SelectorType: SelectorKeyTypeName, Name: string(rns[idx:i])}, i, nil
}

func (d defaultSelectorParser) parseNextIndex(rns []rune, idx int) (SelectorKey, int, error) {
	if idx >= len(rns)-2 {
		return SelectorKey{}, len(rns), ErrInvalidFormatForIndex
	}

	// omit the leading '[' character.
	if rns[idx] != indexSeparatorStart {
		return SelectorKey{}, idx, ErrInvalidFormatForIndex
	}
	idx++

	var i int
	for i = idx; i < len(rns); i++ {
		if rns[i] == indexSeparatorEnd {
			break
		}
		if !d.isRuneValidForIndex(rns[i]) {
			return SelectorKey{}, i, ErrInvalidFormatForName
		}
	}

	if i == idx || rns[i] != indexSeparatorEnd {
		return SelectorKey{}, 0, ErrInvalidFormatForIndex
	}

	str := string(rns[idx:i])

	n, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return SelectorKey{}, 0, err
	}

	// omit the trailing ']'.
	i++
	return SelectorKey{SelectorType: SelectorKeyTypeIndex, Index: int(n)}, i, nil
}

func (d defaultSelectorParser) Parse(s string) ([]SelectorKey, error) {
	if len(s) == 0 {
		return nil, nil
	}

	selector := make([]SelectorKey, 0, d.estimateCount(s))

	for idx, rns := 0, []rune(s); idx < len(rns); {
		switch {
		case rns[idx] == indexSeparatorStart:
			s, i, err := d.parseNextIndex(rns, idx)
			if err != nil {
				return nil, err
			}
			selector = append(selector, s)
			idx = i
		case rns[idx] == nameSeparator || idx == 0:
			s, i, err := d.parseNextName(rns, idx)
			if err != nil {
				return nil, err
			}
			selector = append(selector, s)
			idx = i
		default:
			return nil, ErrInvalidFormat
		}
	}

	return selector, nil
}

var (
	ErrInvalidFormatForName  = errors.New("invalid format for name selector")
	ErrInvalidFormatForIndex = errors.New("invalid format for index selector")
	ErrInvalidFormat         = errors.New("invalid selector format")
)
