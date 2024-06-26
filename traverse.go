package pick

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/moukoublen/pick/internal/errorsx"
)

type DefaultTraverser struct {
	caster              Caster
	nilVal              reflect.Value
	skipItemDereference bool
}

func NewDefaultTraverser(caster Caster) DefaultTraverser {
	return DefaultTraverser{
		caster:              caster,
		skipItemDereference: false,
		nilVal:              reflect.Value{},
	}
}

func (d DefaultTraverser) Retrieve(data any, path []Key) (any, error) {
	if len(path) == 0 {
		return data, nil
	}

	var (
		currentItem any
		err         error
	)

	currentItem = data
	for i, field := range path {
		currentItem, err = d.accessKey(currentItem, field)
		if err != nil {
			return currentItem, newTraverseError("error trying to traverse", path, i, err)
		}
	}

	if currentItem == nil || d.skipItemDereference {
		return currentItem, nil
	}

	// try dereference if pointer or interface
	typeOfItem := reflect.TypeOf(currentItem)
	kindOfItem := typeOfItem.Kind()
	if kindOfItem == reflect.Pointer || kindOfItem == reflect.Interface {
		return d.deref(currentItem), nil
	}

	return currentItem, nil
}

func (d DefaultTraverser) accessKey(item any, key Key) (any, error) {
	if item == nil {
		return nil, ErrFieldNotFound
	}

	// attempts to fast return without reflect.
	switch key.Type {
	case KeyTypeField:
		// fast return if item is map[string]any.
		if m, isMap := item.(map[string]any); isMap {
			val, found := m[key.Name]
			if !found {
				return val, ErrFieldNotFound
			}
			return val, nil
		}

	case KeyTypeIndex:
		// fast return if item is []any.
		if s, isSlice := item.([]any); isSlice {
			i, err := key.calculateIndex(len(s))
			if err != nil {
				return nil, err
			}
			return s[i], nil
		}
	}

	// "slow" return.
	typeOfItem := reflect.TypeOf(item)
	kindOfItem := typeOfItem.Kind()

	var resultValue reflect.Value
	var resultError error
	switch kindOfItem {
	case reflect.Map:
		valueOfItem := reflect.ValueOf(item)
		resultValue, resultError = d.getValueFromMap(typeOfItem, kindOfItem, valueOfItem, key)

	case reflect.Struct:
		valueOfItem := reflect.ValueOf(item)
		resultValue, resultError = d.getValueFromStruct(typeOfItem, kindOfItem, valueOfItem, key)

	case reflect.Array, reflect.Slice:
		valueOfItem := reflect.ValueOf(item)
		resultValue, resultError = d.getValueFromSlice(valueOfItem, key)

	case reflect.Pointer, reflect.Interface: // if pointer/interface get target and re-call.
		derefItem := d.deref(item)
		return d.accessKey(derefItem, key)

	default:
		return nil, ErrFieldNotFound
	}

	if resultValue.IsValid() {
		return resultValue.Interface(), resultError
	}

	return nil, resultError
}

func (d DefaultTraverser) getValueFromMap(typeOfItem reflect.Type, _ reflect.Kind, valueOfItem reflect.Value, key Key) (returnValue reflect.Value, err error) {
	defer errorsx.RecoverPanicToError(&err)

	kindOfMapKey := typeOfItem.Key().Kind()

	var resultValue reflect.Value

	switch {
	case kindOfMapKey == reflect.String && key.IsField():
		resultValue = valueOfItem.MapIndex(reflect.ValueOf(key.Name))
	case kindOfMapKey == reflect.String && key.IsIndex():
		k := strconv.Itoa(key.Index)
		resultValue = valueOfItem.MapIndex(reflect.ValueOf(k))
	case key.IsField():
		key, err := d.caster.As(key.Name, kindOfMapKey)
		if err != nil {
			return d.nilVal, errors.Join(ErrKeyCast, err)
		}
		resultValue = valueOfItem.MapIndex(reflect.ValueOf(key))
	case key.IsIndex():
		key, err := d.caster.As(key.Index, kindOfMapKey)
		if err != nil {
			return d.nilVal, errors.Join(ErrKeyCast, err)
		}
		resultValue = valueOfItem.MapIndex(reflect.ValueOf(key))
	}

	return resultValue, nil
}

func (d DefaultTraverser) getValueFromSlice(valueOfItem reflect.Value, key Key) (returnValue reflect.Value, err error) {
	defer errorsx.RecoverPanicToError(&err)

	var resultValue reflect.Value

	if key.IsIndex() {
		idx, err := key.calculateIndex(valueOfItem.Len())
		if err != nil {
			return d.nilVal, err
		}
		resultValue = valueOfItem.Index(idx)
	} else if key.IsField() {
		// try to cast to int
		i, err := d.caster.AsInt(key.Name)
		if err != nil {
			return d.nilVal, errors.Join(ErrKeyCast, err)
		}
		k := Index(i)
		return d.getValueFromSlice(valueOfItem, k)
	}

	return resultValue, nil
}

func (d DefaultTraverser) getValueFromStruct(_ reflect.Type, _ reflect.Kind, valueOfItem reflect.Value, key Key) (returnValue reflect.Value, err error) {
	defer errorsx.RecoverPanicToError(&err)

	var resultValue reflect.Value

	if key.IsIndex() {
		resultValue = valueOfItem.Field(key.Index)
	} else if key.IsField() {
		resultValue = valueOfItem.FieldByName(key.Name)
	}

	if !resultValue.IsValid() {
		return d.nilVal, ErrFieldNotFound
	}

	return resultValue, nil
}

func (d DefaultTraverser) deref(item any) any {
	valueOfItem := reflect.ValueOf(item)
	targetValue := valueOfItem.Elem()
	return targetValue.Interface()
}

type TraverseError struct {
	inner      error
	msg        string
	path       []Key
	fieldIndex int
}

func newTraverseError(msg string, path []Key, fieldIndex int, inner error) *TraverseError {
	return &TraverseError{
		msg:        msg,
		path:       path,
		fieldIndex: fieldIndex,
		inner:      inner,
	}
}

func (t *TraverseError) Unwrap() error {
	return t.inner
}

func (t *TraverseError) Error() string {
	if t.inner != nil {
		return fmt.Sprintf("selector: %s - %s: %s", formatErrorAt(t.path, t.fieldIndex), t.msg, t.inner.Error())
	}
	return fmt.Sprintf("selector: %s - %s", formatErrorAt(t.path, t.fieldIndex), t.msg)
}

var (
	ErrFieldNotFound   = errors.New("field not found")
	ErrIndexOutOfRange = fmt.Errorf("%w: index out of range", ErrFieldNotFound)
	ErrKeyCast         = errors.New("key cast error")
)

func formatErrorAt(path []Key, idx int) string {
	sb := strings.Builder{}
	for i, c := range path {
		if c.IsIndex() {
			if i == idx {
				sb.WriteString(">")
			}
			sb.WriteString(fmt.Sprintf("[%d]", c.Index))
		} else {
			if i > 0 {
				sb.WriteRune(fieldSeparator)
			}
			if i == idx {
				sb.WriteString(">")
			}
			sb.WriteString(c.Name)
		}
		if i == idx {
			sb.WriteString("<")
		}
	}

	return sb.String()
}
