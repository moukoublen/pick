package pick

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	errorsx "github.com/moukoublen/pick/errors"
)

type DefaultTraverser struct {
	caster              Caster
	skipItemDereference bool
}

func NewDefaultTraverser(caster Caster) DefaultTraverser {
	return DefaultTraverser{
		caster:              caster,
		skipItemDereference: false,
	}
}

func (d DefaultTraverser) Get(data any, fields []Field) (any, error) {
	return d.get(data, fields)
}

func (d DefaultTraverser) get(obj any, fields []Field) (any, error) {
	if len(fields) == 0 {
		return obj, nil
	}

	var (
		currentItem any
		err         error
	)

	currentItem = obj
	for i, field := range fields {
		currentItem, err = d.getField(currentItem, field)
		if err != nil {
			return currentItem, newTraverseError("error trying to traverse", fields, i, err)
		}
	}

	if d.skipItemDereference {
		return currentItem, nil
	}

	if currentItem == nil {
		return currentItem, nil
	}

	typeOfItem := reflect.TypeOf(currentItem)
	kindOfItem := typeOfItem.Kind()
	if kindOfItem == reflect.Pointer || kindOfItem == reflect.Interface {
		return d.deref(currentItem), nil
	}

	return currentItem, nil
}

func (d DefaultTraverser) getField(item any, field Field) (any, error) {
	// attempts to fast return without reflect.
	switch field.Type {
	case NotationFieldTypeName:
		// fast return if item is map[string]any.
		if m, isMap := item.(map[string]any); isMap {
			val, found := m[field.Name]
			if !found {
				return val, ErrFieldNotFound
			}
			return val, nil
		}

	case NotationFieldTypeIndex:
		// fast return if item is []any.
		if s, isSlice := item.([]any); isSlice {
			if field.Index >= len(s) || field.Index < 0 {
				return nil, ErrIndexOutOfRange
			}
			return s[field.Index], nil
		}
	}

	// "slow" return.
	typeOfItem := reflect.TypeOf(item)
	kindOfItem := typeOfItem.Kind()

	switch kindOfItem {
	case reflect.Map:
		valueOfItem := reflect.ValueOf(item)
		return d.accessMap(typeOfItem, kindOfItem, valueOfItem, field)

	case reflect.Struct:
		valueOfItem := reflect.ValueOf(item)
		return d.accessStruct(typeOfItem, kindOfItem, valueOfItem, field)

	case reflect.Array, reflect.Slice:
		valueOfItem := reflect.ValueOf(item)
		return d.accessSlice(typeOfItem, kindOfItem, valueOfItem, field)

	case reflect.Pointer, reflect.Interface: // if pointer/interface get target and re-call.
		derefItem := d.deref(item)
		return d.getField(derefItem, field)
	}

	return nil, ErrFieldNotFound
}

func (d DefaultTraverser) accessMap(typeOfItem reflect.Type, _ reflect.Kind, valueOfItem reflect.Value, field Field) (returnValue any, err error) {
	defer errorsx.RecoverPanicToError(&err)

	kindOfMapKey := typeOfItem.Key().Kind()

	var resultValue reflect.Value

	switch {
	case kindOfMapKey == reflect.String && field.IsName():
		resultValue = valueOfItem.MapIndex(reflect.ValueOf(field.Name))
	case kindOfMapKey == reflect.String && field.IsIndex():
		k := strconv.Itoa(field.Index)
		resultValue = valueOfItem.MapIndex(reflect.ValueOf(k))
	case field.IsName():
		key, err := d.caster.As(field.Name, kindOfMapKey)
		if err != nil {
			return nil, errors.Join(ErrKeyCast, err)
		}
		resultValue = valueOfItem.MapIndex(reflect.ValueOf(key))
	case field.IsIndex():
		key, err := d.caster.As(field.Index, kindOfMapKey)
		if err != nil {
			return nil, errors.Join(ErrKeyCast, err)
		}
		resultValue = valueOfItem.MapIndex(reflect.ValueOf(key))
	}

	if !resultValue.IsValid() {
		return nil, ErrFieldNotFound
	}

	return resultValue.Interface(), nil
}

func (d DefaultTraverser) accessSlice(_ reflect.Type, _ reflect.Kind, valueOfItem reflect.Value, field Field) (returnValue any, err error) {
	defer errorsx.RecoverPanicToError(&err)

	var resultValue reflect.Value

	if field.IsIndex() {
		if field.Index >= valueOfItem.Len() {
			return nil, ErrIndexOutOfRange
		}
		resultValue = valueOfItem.Index(field.Index)
	} else if field.IsName() {
		// try to cast to int
		i, err := d.caster.AsInt(field.Name)
		if err != nil {
			return nil, errors.Join(ErrKeyCast, err)
		}
		if i >= valueOfItem.Len() {
			return nil, ErrIndexOutOfRange
		}
		resultValue = valueOfItem.Index(i)
	}

	if !resultValue.IsValid() {
		return nil, ErrIndexOutOfRange
	}

	return resultValue.Interface(), nil
}

func (d DefaultTraverser) accessStruct(_ reflect.Type, _ reflect.Kind, valueOfItem reflect.Value, field Field) (returnValue any, err error) {
	defer errorsx.RecoverPanicToError(&err)

	var resultValue reflect.Value

	if field.IsIndex() {
		resultValue = valueOfItem.Field(field.Index)
	} else if field.IsName() {
		resultValue = valueOfItem.FieldByName(field.Name)
	}

	if !resultValue.IsValid() {
		return nil, ErrFieldNotFound
	}

	return resultValue.Interface(), nil
}

func (d DefaultTraverser) deref(item any) any {
	valueOfItem := reflect.ValueOf(item)
	targetValue := valueOfItem.Elem()
	return targetValue.Interface()
}

type TraverseError struct {
	inner      error
	msg        string
	fields     []Field
	fieldIndex int
}

func newTraverseError(msg string, fields []Field, fieldIndex int, inner error) *TraverseError {
	return &TraverseError{
		msg:        msg,
		fields:     fields,
		fieldIndex: fieldIndex,
		inner:      inner,
	}
}

func (t *TraverseError) Unwrap() error {
	return t.inner
}

func (t *TraverseError) Error() string {
	if t.inner != nil {
		return fmt.Sprintf("selector: %s - %s: %s", formatErrorAt(t.fields, t.fieldIndex), t.msg, t.inner.Error())
	}
	return fmt.Sprintf("selector: %s - %s", formatErrorAt(t.fields, t.fieldIndex), t.msg)
}

var (
	ErrIndexOutOfRange = errors.New("index out of range")
	ErrKeyCast         = errors.New("key cast error")
)

func formatErrorAt(s []Field, idx int) string {
	sb := strings.Builder{}
	for i, c := range s {
		if c.IsIndex() {
			if i == idx {
				sb.WriteString(">")
			}
			sb.WriteString(fmt.Sprintf("[%d]", c.Index))
		} else {
			if i > 0 {
				sb.WriteRune(nameSeparator)
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

var ErrFieldNotFound = errors.New("field not found")
