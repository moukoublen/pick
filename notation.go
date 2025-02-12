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

func (s Key) calculateIndex(length int) (int, error) {
	i := s.Index
	if i < 0 {
		i = length + i
	}

	if i >= length {
		return i, ErrIndexOutOfRange
	}

	return i, nil
}

func (s Key) IsIndex() bool { return s.Type == KeyTypeIndex }
func (s Key) IsField() bool { return s.Type == KeyTypeField }

func (s Key) Any() any {
	switch s.Type {
	case KeyTypeIndex:
		return s.Index
	case KeyTypeField:
		return s.Name
	}

	return nil
}

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

type fsmState int

const (
	fsmStateReset fsmState = iota
	fsmStateIndexStarted
	fsmStateIndex
	fsmStateIndexEnded
	fsmStateFieldSeparated
	fsmStateField
)

func (s fsmState) Input(received rune) (fsmState, error) {
	switch s {
	case fsmStateReset:
		return s.stateReset(received)
	case fsmStateIndexStarted:
		return s.stateIndexStarted(received)
	case fsmStateIndex:
		return s.stateIndex(received)
	case fsmStateIndexEnded:
		return s.stateIndexEnded(received)
	case fsmStateFieldSeparated:
		return s.stateFieldSeparated(received)
	case fsmStateField:
		return s.stateField(received)
	}

	return s, ErrInvalidSelectorFormat
}

func (s fsmState) stateReset(received rune) (fsmState, error) {
	switch {
	case received == fieldSeparator:
		return fsmStateFieldSeparated, nil

	case received == indexSeparatorStart:
		return fsmStateIndexStarted, nil

	case unicode.IsControl(received):
		return fsmStateReset, ErrInvalidSelectorFormatForName

	default:
		return fsmStateField, nil
	}
}

func (s fsmState) stateIndexStarted(received rune) (fsmState, error) {
	switch {
	case unicode.IsNumber(received):
		return fsmStateIndex, nil
	case received == '-':
		return fsmStateIndex, nil

	default:
		return fsmStateIndexStarted, ErrInvalidSelectorFormatForIndex
	}
}

func (s fsmState) stateIndex(received rune) (fsmState, error) {
	switch {
	case unicode.IsNumber(received):
		return fsmStateIndex, nil

	case received == indexSeparatorEnd:
		return fsmStateIndexEnded, nil

	default:
		return fsmStateIndex, ErrInvalidSelectorFormatForIndex
	}
}

func (s fsmState) stateIndexEnded(received rune) (fsmState, error) {
	switch received {
	case fieldSeparator:
		return fsmStateFieldSeparated, nil

	case indexSeparatorStart:
		return fsmStateIndexStarted, nil

	default:
		return fsmStateIndexEnded, ErrInvalidSelectorFormat
	}
}

func (s fsmState) stateField(received rune) (fsmState, error) {
	switch {
	case received == fieldSeparator:
		return fsmStateFieldSeparated, nil

	case received == indexSeparatorStart:
		return fsmStateIndexStarted, nil

	case unicode.IsControl(received):
		return fsmStateField, ErrInvalidSelectorFormatForName

	default:
		return fsmStateField, nil
	}
}

func (s fsmState) stateFieldSeparated(received rune) (fsmState, error) {
	switch {
	case received == fieldSeparator:
		return fsmStateFieldSeparated, ErrInvalidSelectorFormatForName

	case unicode.IsControl(received):
		return fsmStateField, ErrInvalidSelectorFormatForName

	default:
		return fsmStateField, nil
	}
}

func (s fsmState) oneOf(states ...fsmState) bool {
	for _, st := range states {
		if s == st {
			return true
		}
	}
	return false
}

type dotNotationParser struct{}

func (d dotNotationParser) Parse(selector string) ([]Key, error) {
	if len(selector) == 0 {
		return nil, nil
	}

	keys := make([]Key, 0, d.estimatePathSize(selector))

	var (
		lastState  fsmState
		tokenStart int
	)
	for i, r := range selector {
		newState, err := lastState.Input(r)
		if err != nil {
			return nil, err
		}

		if lastState == newState {
			continue
		}

		// <reset> -> <...> or <field separated> -> <...> or <index started> -> <...>
		//   => new token starts
		if lastState.oneOf(fsmStateReset, fsmStateFieldSeparated, fsmStateIndexStarted) {
			tokenStart = i
			lastState = newState
			continue
		}

		// <index> -> <index ended>
		//   => new index token
		if lastState == fsmStateIndex && newState == fsmStateIndexEnded {
			k, err := d.parseIndexToken(selector[tokenStart:i])
			if err != nil {
				return nil, err
			}
			keys = append(keys, k)
			lastState = newState
			continue
		}

		// <field> -> <field separated || index started>
		//   => new field token
		if lastState == fsmStateField && newState.oneOf(fsmStateFieldSeparated, fsmStateIndexStarted) {
			k, err := d.parseFieldToken(selector[tokenStart:i])
			if err != nil {
				return nil, err
			}
			keys = append(keys, k)
		}

		lastState = newState
	}

	// ends with <reset || index ended>, valid
	if lastState.oneOf(fsmStateReset, fsmStateIndexEnded) {
		return keys, nil
	}

	// ends with incomplete index, invalid.
	if lastState.oneOf(fsmStateIndexStarted, fsmStateIndex) {
		return nil, ErrInvalidSelectorFormatForIndex
	}

	// ends with incomplete field, invalid.
	if lastState == fsmStateFieldSeparated {
		return nil, ErrInvalidSelectorFormatForName
	}

	// ends with field, valid but must parse.
	if lastState == fsmStateField {
		k, err := d.parseFieldToken(selector[tokenStart:])
		if err != nil {
			return nil, err
		}
		keys = append(keys, k)

		return keys, nil
	}

	return keys, nil
}

func (d dotNotationParser) estimatePathSize(selector string) int {
	numOfKeys := 1
	for _, r := range selector {
		if r == fieldSeparator || r == indexSeparatorStart {
			numOfKeys++
		}
	}

	return numOfKeys
}

func (d dotNotationParser) parseIndexToken(token string) (Key, error) {
	k := Key{Type: KeyTypeIndex}
	i, err := strconv.Atoi(token)
	if err != nil {
		return Key{}, errors.Join(ErrInvalidSelectorFormatForName, err)
	}
	k.Index = int(i)

	return k, nil
}

func (d dotNotationParser) parseFieldToken(token string) (Key, error) {
	k := Key{Type: KeyTypeField}
	if len(token) == 0 {
		return k, ErrInvalidSelectorFormatForName
	}

	k.Name = token
	return k, nil
}

var (
	ErrInvalidSelectorFormatForName  = errors.New("invalid format for name key")
	ErrInvalidSelectorFormatForIndex = errors.New("invalid format for index key")
	ErrInvalidSelectorFormat         = errors.New("invalid selector format")
)

// formatPath uses default dot notation formatter to format a path to string.
func formatPath(path []Key) string {
	f := dotNotationFormatter{}
	return f.Format(path...)
}
