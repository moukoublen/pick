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
	AsString(item any) (string, error)
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

func (d DefaultTraverser) Set(data any, path []Key, newValue any) (rErr error) {
	defer errorsx.RecoverPanicToError(&rErr)

	// optimized for map[string]any
	if ma, is := data.(map[string]any); is {
		return d.setToMapStringAny(ma, path, newValue)
	}

	// optimized for *map[string]any
	if ma, is := data.(*map[string]any); is {
		if ma == nil {
			ma = &map[string]any{}
		}
		return d.setToMapStringAny(*ma, path, newValue)
	}

	var (
		valueOfDest = reflect.ValueOf(data)
		typeOfDest  = reflect.TypeOf(data)
		kindOfDest  = typeOfDest.Kind()
	)

	if len(path) == 0 {
		// instant set
		castedVal, err := d.keyCaster.ByType(newValue, typeOfDest)
		if err != nil {
			return err
		}
		valueOfDest.Set(reflect.ValueOf(castedVal))
		return nil
	}

	refresh := func(val reflect.Value) {
		valueOfDest = val
		typeOfDest = val.Type()
		kindOfDest = val.Kind()
	}

	// traverse
	lastKey := path[len(path)-1]
	if len(path) > 1 {
		excludeLast := path[:len(path)-1]
		for _, field := range excludeLast {
			resultValue, err := d.accessWithReflect(typeOfDest, kindOfDest, valueOfDest, field)
			if err != nil {
				return err
			}
			refresh(resultValue)
		}
	}

	// set value using lastKey
	return d.setWithReflect(valueOfDest, typeOfDest, kindOfDest, lastKey, newValue)
}

// setToMapStringAny is the optimistic flow that assumes only map[string]any
func (d DefaultTraverser) setToMapStringAny(ma map[string]any, path []Key, newValue any) error {
	if len(path) == 0 {
		// tbd
	}

	// for _, p := range path {
	// }

	return nil
}

func (d DefaultTraverser) accessWithReflect(typeOfItem reflect.Type, kindOfItem reflect.Kind, valueOfItem reflect.Value, key Key) (reflect.Value, error) {
	switch kindOfItem {
	case reflect.Map:
		return d.getValueFromMap(typeOfItem, kindOfItem, valueOfItem, key)
	case reflect.Struct:
		// return d.getValueFromStruct(typeOfItem, kindOfItem, valueOfItem, key)
	case reflect.Array, reflect.Slice:
		return d.getValueFromSlice(valueOfItem, key)
	case reflect.Pointer, reflect.Interface:
		// deref TODO: reconsider
		// check valueOfItem.IsNil()
		v := valueOfItem.Elem()
		return d.accessWithReflect(v.Type(), v.Kind(), v, key)
	default:
		return reflect.Value{}, ErrFieldNotFound
	}

	return reflect.Value{}, ErrFieldNotFound
}

func (d DefaultTraverser) setWithReflect(valueOfDestItem reflect.Value, typeOfDestItem reflect.Type, kindOfDestItem reflect.Kind, key Key, valueToSet any) error {
	switch kindOfDestItem {
	case reflect.Map:
		keyType := typeOfDestItem.Key()
		elemType := typeOfDestItem.Elem()

		keyCastedValue, err := d.keyAsReflectValue(key, keyType)
		if err != nil {
			return err
		}

		valCasted, err := d.keyCaster.ByType(valueToSet, elemType)
		if err != nil {
			return err
		}

		valueOfDestItem.SetMapIndex(keyCastedValue, reflect.ValueOf(valCasted))
		return nil

	case reflect.Struct:
		fieldName, err := d.keyCaster.AsString(key.Any())
		if err != nil {
			return errors.Join(ErrKeyCast, err)
		}
		dst := valueOfDestItem.FieldByName(fieldName)
		return d.setReflectValue(dst, valueToSet)

	case reflect.Array, reflect.Slice:
		itemIndex, err := d.keyCaster.AsInt(key.Any())
		if err != nil {
			return errors.Join(ErrKeyCast, err)
		}

		elemType := typeOfDestItem.Elem()
		valCasted, err := d.keyCaster.ByType(valueToSet, elemType)
		if err != nil {
			return err
		}

		if itemIndex > valueOfDestItem.Len() {
			return ErrIndexOutOfRange
		}

		v := valueOfDestItem.Index(itemIndex)
		if !v.IsValid() {
			return ErrDestinationValueNotValid
		}

		v.Set(reflect.ValueOf(valCasted))

		return nil

	case reflect.Pointer:
		// return ErrDestinationValueNotValid
		v := valueOfDestItem.Elem()
		return d.setWithReflect(v, v.Type(), v.Kind(), key, valueToSet)
	case reflect.Interface:
		v := valueOfDestItem.Elem()
		return d.setWithReflect(v, v.Type(), v.Kind(), key, valueToSet)

	default:
		return ErrFieldNotFound
	}
}

func (d DefaultTraverser) setReflectValue(dst reflect.Value, newVal any) (err error) {
	defer errorsx.RecoverPanicToError(&err)

	if !dst.IsValid() {
		err = ErrDestinationValueNotValid
		return
	}

	dstType := dst.Type()

	var casted any
	casted, err = d.keyCaster.ByType(newVal, dstType)
	if err != nil {
		return err
	}

	dst.Set(reflect.ValueOf(casted))
	return
}

func (d DefaultTraverser) keyAsReflectValue(key Key, asType reflect.Type) (reflect.Value, error) {
	v, err := d.keyCaster.ByType(key.Any(), asType)
	if err != nil {
		return reflect.Value{}, err
	}

	return reflect.ValueOf(v), nil
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
	ErrFieldNotFound            = errors.New("field not found")
	ErrIndexOutOfRange          = fmt.Errorf("%w: index out of range", ErrFieldNotFound)
	ErrKeyCast                  = errors.New("key cast error")
	ErrKeyUnknown               = errors.New("key type unknown")
	ErrDestinationValueNotValid = errors.New("destination value is not valid")
)
