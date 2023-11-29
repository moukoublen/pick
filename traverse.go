package pick

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"

	errorsx "github.com/moukoublen/pick/errors"
)

type DefaultTraverser struct {
	caster              Caster
	selectorFormatter   SelectorFormatter
	skipItemDereference bool
}

func NewDefaultTraverser(caster Caster, selectorFormatter SelectorFormatter) DefaultTraverser {
	return DefaultTraverser{
		caster:              caster,
		selectorFormatter:   selectorFormatter,
		skipItemDereference: false,
	}
}

func (d DefaultTraverser) Get(o any, selector []SelectorKey) (any, bool, error) {
	return d.get(o, selector)
}

func (d DefaultTraverser) get(obj any, selector []SelectorKey) (any, bool, error) {
	if len(selector) == 0 {
		return nil, false, nil
	}

	var (
		currentItem any
		found       bool
		err         error
	)

	currentItem = obj
	for i, curSelector := range selector {
		currentItem, found, err = d.getSingleSelector(currentItem, curSelector)
		if err != nil {
			return currentItem, found, newTraverseError("error trying to traverse", selector, i, d.selectorFormatter, err)
		}
		if !found {
			return nil, false, nil
		}
	}

	if found {
		if d.skipItemDereference {
			return currentItem, found, nil
		}

		if currentItem == nil {
			return currentItem, found, nil
		}
		typeOfItem := reflect.TypeOf(currentItem)
		kindOfItem := typeOfItem.Kind()
		if kindOfItem == reflect.Pointer || kindOfItem == reflect.Interface {
			return d.deref(currentItem), found, nil
		}

		return currentItem, found, nil
	}

	return nil, false, nil
}

func (d DefaultTraverser) getSingleSelector(item any, selector SelectorKey) (any, bool, error) {
	// attempts to fast return without reflect.
	switch selector.SelectorType {
	case SelectorKeyTypeName:
		// fast return if item is map[string]any.
		if m, isMap := item.(map[string]any); isMap {
			val, found := m[selector.Name]
			return val, found, nil
		}

	case SelectorKeyTypeIndex:
		// fast return if item is []any.
		if s, isSlice := item.([]any); isSlice {
			if selector.Index >= len(s) || selector.Index < 0 {
				return nil, false, ErrIndexOutOfRange
			}
			return s[selector.Index], true, nil
		}
	}

	// "slow" return.
	typeOfItem := reflect.TypeOf(item)
	kindOfItem := typeOfItem.Kind()

	switch kindOfItem {
	case reflect.Map:
		valueOfItem := reflect.ValueOf(item)
		return d.accessMap(typeOfItem, kindOfItem, valueOfItem, selector)

	case reflect.Struct:
		valueOfItem := reflect.ValueOf(item)
		return d.accessStruct(typeOfItem, kindOfItem, valueOfItem, selector)

	case reflect.Array, reflect.Slice:
		valueOfItem := reflect.ValueOf(item)
		return d.accessSlice(typeOfItem, kindOfItem, valueOfItem, selector)

	case reflect.Pointer, reflect.Interface: // if pointer/interface get target and re-call.
		derefItem := d.deref(item)
		return d.getSingleSelector(derefItem, selector)
	}

	return nil, false, nil
}

func (d DefaultTraverser) accessMap(typeOfItem reflect.Type, _ reflect.Kind, valueOfItem reflect.Value, selector SelectorKey) (returnValue any, found bool, err error) {
	defer errorsx.RecoverPanicToError(&err)

	kindOfMapKey := typeOfItem.Key().Kind()

	var resultValue reflect.Value

	switch {
	case kindOfMapKey == reflect.String && selector.IsName():
		resultValue = valueOfItem.MapIndex(reflect.ValueOf(selector.Name))
	case kindOfMapKey == reflect.String && selector.IsIndex():
		k := strconv.Itoa(selector.Index)
		resultValue = valueOfItem.MapIndex(reflect.ValueOf(k))
	case selector.IsName():
		key, err := d.caster.As(selector.Name, kindOfMapKey)
		if err != nil {
			return nil, false, errors.Join(ErrKeyCast, err)
		}
		resultValue = valueOfItem.MapIndex(reflect.ValueOf(key))
	case selector.IsIndex():
		key, err := d.caster.As(selector.Index, kindOfMapKey)
		if err != nil {
			return nil, false, errors.Join(ErrKeyCast, err)
		}
		resultValue = valueOfItem.MapIndex(reflect.ValueOf(key))
	}

	if !resultValue.IsValid() {
		return nil, false, nil
	}

	return resultValue.Interface(), true, nil
}

func (d DefaultTraverser) accessSlice(_ reflect.Type, _ reflect.Kind, valueOfItem reflect.Value, selector SelectorKey) (returnValue any, found bool, err error) {
	defer errorsx.RecoverPanicToError(&err)

	var resultValue reflect.Value

	if selector.IsIndex() {
		if selector.Index >= valueOfItem.Len() {
			return nil, false, ErrIndexOutOfRange
		}
		resultValue = valueOfItem.Index(selector.Index)
	} else if selector.IsName() {
		// try to cast to int
		i, err := d.caster.AsInt(selector.Name)
		if err != nil {
			return nil, false, errors.Join(ErrKeyCast, err)
		}
		if i >= valueOfItem.Len() {
			return nil, false, ErrIndexOutOfRange
		}
		resultValue = valueOfItem.Index(i)
	}

	if !resultValue.IsValid() {
		return nil, false, nil
	}

	return resultValue.Interface(), true, nil
}

func (d DefaultTraverser) accessStruct(_ reflect.Type, _ reflect.Kind, valueOfItem reflect.Value, selector SelectorKey) (returnValue any, found bool, err error) {
	defer errorsx.RecoverPanicToError(&err)

	var resultValue reflect.Value

	if selector.IsIndex() {
		resultValue = valueOfItem.Field(selector.Index)
	} else if selector.IsName() {
		resultValue = valueOfItem.FieldByName(selector.Name)
	}

	if !resultValue.IsValid() {
		return nil, false, nil
	}

	return resultValue.Interface(), true, nil
}

func (d DefaultTraverser) deref(item any) any {
	valueOfItem := reflect.ValueOf(item)
	targetValue := valueOfItem.Elem()
	return targetValue.Interface()
}

type TraverseError struct {
	selectorFormatter SelectorFormatter
	inner             error
	msg               string
	selector          []SelectorKey
	selectorIndex     int
}

func newTraverseError(msg string, selector []SelectorKey, selectorIndex int, selectorFormatter SelectorFormatter, inner error) *TraverseError {
	return &TraverseError{
		msg:               msg,
		selector:          selector,
		selectorIndex:     selectorIndex,
		selectorFormatter: selectorFormatter,
		inner:             inner,
	}
}

func (t *TraverseError) Unwrap() error {
	return t.inner
}

func (t *TraverseError) Error() string {
	if t.inner != nil {
		return fmt.Sprintf("selector: %s - %s: %s", t.selectorFormatter.FormatWithIndex(t.selector, t.selectorIndex), t.msg, t.inner.Error())
	}
	return fmt.Sprintf("selector: %s - %s", t.selectorFormatter.FormatWithIndex(t.selector, t.selectorIndex), t.msg)
}

var (
	ErrIndexOutOfRange = errors.New("index out of range")
	ErrKeyCast         = errors.New("key cast error")
)
