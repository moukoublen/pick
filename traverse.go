package pick

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"github.com/moukoublen/pick/internal/errorsx"
)

type KeyCaster interface {
	ByType(input any, asType reflect.Type) (any, error)
	AsInt(item any) (int, error)
}

type DefaultTraverser struct {
	keyCaster           KeyCaster
	nilVal              reflect.Value
	skipItemDereference bool
}

func NewDefaultTraverser(keyCaster KeyCaster) DefaultTraverser {
	return DefaultTraverser{
		keyCaster:           keyCaster,
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
			return currentItem, NewTraverseError("error trying to traverse", path, i, err)
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
		// fast return if item is a slice of basic type.
		handled, val, err := attemptAccessSliceOfBasicType(item, key.Index)
		if handled || err != nil {
			return val, err
		}
	}

	// "slow" return.
	typeOfItem := reflect.TypeOf(item)
	kindOfItem := typeOfItem.Kind()

	var resultValue reflect.Value
	var resultError error
	switch kindOfItem {
	case reflect.Map:
		resultValue, resultError = d.getValueFromMap(typeOfItem, kindOfItem, item, key)

	case reflect.Struct:
		resultValue, resultError = d.getValueFromStruct(item, key)

	case reflect.Array, reflect.Slice:
		resultValue, resultError = d.getValueFromSlice(item, key)

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

func (d DefaultTraverser) getValueFromMap(typeOfItem reflect.Type, _ reflect.Kind, item any, key Key) (returnValue reflect.Value, err error) {
	defer errorsx.RecoverPanicToError(&err)

	valueOfItem := reflect.ValueOf(item)

	kindOfMapKey := typeOfItem.Key().Kind()

	var resultValue reflect.Value

	switch {
	case kindOfMapKey == reflect.String && key.IsField():
		resultValue = valueOfItem.MapIndex(reflect.ValueOf(key.Name))
	case kindOfMapKey == reflect.String && key.IsIndex():
		k := strconv.Itoa(key.Index)
		resultValue = valueOfItem.MapIndex(reflect.ValueOf(k))
	default:
		key, err := d.keyCaster.ByType(key.Any(), typeOfItem.Key())
		if err != nil {
			return d.nilVal, errors.Join(ErrKeyCast, err)
		}
		resultValue = valueOfItem.MapIndex(reflect.ValueOf(key))
	}

	return resultValue, nil
}

func (d DefaultTraverser) getValueFromSlice(item any, key Key) (returnValue reflect.Value, err error) {
	defer errorsx.RecoverPanicToError(&err)

	valueOfItem := reflect.ValueOf(item)

	var index int
	switch key.Type {
	case KeyTypeIndex:
		index = key.Index
	case KeyTypeField:
		i, er := d.keyCaster.AsInt(key.Name)
		if er != nil {
			return d.nilVal, errors.Join(ErrKeyCast, er)
		}
		index = i
	default:
		return d.nilVal, ErrKeyUnknown
	}

	idx, err := Index(index).calculateIndex(valueOfItem.Len())
	if err != nil {
		return d.nilVal, err
	}

	return valueOfItem.Index(idx), nil
}

func (d DefaultTraverser) getValueFromStruct(item any, key Key) (returnValue reflect.Value, err error) {
	defer errorsx.RecoverPanicToError(&err)

	valueOfItem := reflect.ValueOf(item)

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

func attemptAccessSliceOfBasicType(sl any, index int) (handled bool, val any, err error) {
	key := Index(index)

	switch sl := sl.(type) {
	case []any:
		val, err = accessSlice(sl, key)
	case []int:
		val, err = accessSlice(sl, key)
	case []int8:
		val, err = accessSlice(sl, key)
	case []int16:
		val, err = accessSlice(sl, key)
	case []int32:
		val, err = accessSlice(sl, key)
	case []int64:
		val, err = accessSlice(sl, key)
	case []uint:
		val, err = accessSlice(sl, key)
	case []uint8:
		val, err = accessSlice(sl, key)
	case []uint16:
		val, err = accessSlice(sl, key)
	case []uint32:
		val, err = accessSlice(sl, key)
	case []uint64:
		val, err = accessSlice(sl, key)
	case []float32:
		val, err = accessSlice(sl, key)
	case []float64:
		val, err = accessSlice(sl, key)
	case []bool:
		val, err = accessSlice(sl, key)
	case []string:
		val, err = accessSlice(sl, key)
	default:
		return false, nil, nil
	}

	return true, val, err
}

func accessSlice[T any](sl []T, key Key) (any, error) {
	i, err := key.calculateIndex(len(sl))
	if err != nil {
		return nil, err
	}

	return sl[i], nil
}

type TraverseError struct {
	inner      error
	msg        string
	path       []Key
	fieldIndex int
}

func NewTraverseError(msg string, path []Key, fieldIndex int, inner error) *TraverseError {
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
		return fmt.Sprintf("selector: %s : %s: %s", formatPath(t.path[:t.fieldIndex+1]), t.msg, t.inner.Error())
	}
	return fmt.Sprintf("selector: %s : %s", formatPath(t.path[:t.fieldIndex+1]), t.msg)
}

func (t *TraverseError) Path() []Key {
	return t.path[:t.fieldIndex+1]
}

var (
	ErrFieldNotFound   = errors.New("field not found")
	ErrIndexOutOfRange = fmt.Errorf("%w: index out of range", ErrFieldNotFound)
	ErrKeyCast         = errors.New("key cast error")
	ErrKeyUnknown      = errors.New("key type unknown")
)
