package gopt

import (
	"database/sql"
	"encoding/json"
	"errors"
	"reflect"
)

type String *string

// Of creates a new optional with the supplied value
func Of[T any](value T) Optional[T] {
	return Optional[T]{
		present: isPresent(value),
		value:   value,
	}
}

// OfNillable creates a new optional with the supplied value
//
// If the supplied value is nil, an empty (not present) optional is returned
func OfNillable[T any](value T) Optional[T] {
	if isPresent(value) {
		return Optional[T]{
			present: true,
			value:   value,
		}
	}
	return Optional[T]{
		present: false,
	}
}

// OfNillableString creates a new string optional with the supplied value
//
// If the supplied value is an empty string, an empty (not-present) optional is returned
func OfNillableString(value string) Optional[string] {
	return Optional[string]{
		present: value != "",
		value:   value,
	}
}

// Empty creates a new empty (not-present) optional of the specified type
func Empty[T any]() Optional[T] {
	return Optional[T]{
		present: false,
	}
}

func isPresent(v any) bool {
	vo := reflect.ValueOf(v)
	switch vk := vo.Kind(); vk {
	case reflect.Ptr:
		vk = vo.Elem().Kind()
		if vk == reflect.Invalid {
			return false
		}
	case reflect.Map, reflect.Slice, reflect.Interface:
		return !vo.IsNil()
	}
	return v != nil
}

type Optional[T any] struct {
	present bool
	value   T
	set     bool
}

// Get returns the value and an error if the value is not present
func (o Optional[T]) Get() (T, error) {
	if !o.present {
		return o.emptyValue(), errors.New("not present")
	}
	return o.value, nil
}

func (o Optional[T]) emptyValue() T {
	return (Optional[T]{}).value
}

func (o *Optional[T]) clear(set bool) {
	o.value = o.emptyValue()
	o.present = false
	o.set = set
}

// AsEmpty returns a new empty optional of the same type
func (o Optional[T]) AsEmpty() Optional[T] {
	return Optional[T]{
		present: false,
	}
}

// IsPresent returns true if the value is present, otherwise false
func (o Optional[T]) IsPresent() bool {
	return o.present
}

// WasSet returns true if the last setting operation set the value, otherwise false
//
// Setting operations are UnmarshalJSON, Scan and OrElseSet
//
// Use UnSet() to clear this flag alone
func (o Optional[T]) WasSet() bool {
	return o.set
}

// UnSet clears the set flag (see WasSet)
func (o *Optional[T]) UnSet() Optional[T] {
	o.set = false
	return *o
}

// Clear clears the optional
//
// Clearing sets the present to false, the set flag to false and the value to an empty value
func (o *Optional[T]) Clear() Optional[T] {
	o.clear(false)
	return *o
}

// IfPresent if the value is present, calls the supplied function with the value, otherwise does nothing
func (o Optional[T]) IfPresent(f func(v T)) Optional[T] {
	if o.present && f != nil {
		f(o.value)
	}
	return o
}

// IfPresentOtherwise if the value is present, calls the supplied function with the value, otherwise calls the other function
func (o Optional[T]) IfPresentOtherwise(f func(v T), other func()) Optional[T] {
	if o.present && f != nil {
		f(o.value)
	} else if !o.present && other != nil {
		other()
	}
	return o
}

// IfSet if the value was set and is present, calls the supplied function with the value
//
// if the value was set but is not present, calls the supplied notPresent function
//
// otherwise, does nothing
func (o Optional[T]) IfSet(f func(v T), notPresent func()) Optional[T] {
	if o.set && o.present && f != nil {
		f(o.value)
	} else if o.set && !o.present && notPresent != nil {
		notPresent()
	}
	return o
}

// IfSetOtherwise if the value was set and is present, calls the supplied function with the value
//
// if the value was set but is not present, calls the supplied notPresent function
//
// otherwise, calls the other func
func (o Optional[T]) IfSetOtherwise(f func(v T), notPresent func(), other func()) Optional[T] {
	if o.set && o.present && f != nil {
		f(o.value)
	} else if o.set && !o.present && notPresent != nil {
		notPresent()
	} else if !o.set && !o.present && other != nil {
		other()
	}
	return o
}

// OrElse returns the value if present, otherwise returns other
func (o Optional[T]) OrElse(other T) T {
	if o.present {
		return o.value
	}
	return other
}

// OrElseGet returns the value if present, otherwise returns the result of calling the supplied function
func (o Optional[T]) OrElseGet(f func() T) T {
	if o.present {
		return o.value
	}
	return f()
}

// OrElseSet if the value is not present it is set to the supplied value
func (o *Optional[T]) OrElseSet(v T) Optional[T] {
	if !o.present {
		if isPresent(v) {
			o.present = true
			o.value = v
		} else {
			o.present = false
			o.value = o.emptyValue()
		}
		o.set = true
	}
	return *o
}

// OrElseError returns the supplied error if the value is not present, otherwise returns nil
func (o Optional[T]) OrElseError(err error) error {
	if !o.present {
		return err
	}
	return nil
}

// OrElsePanic if the value is not present, panics with the supplied value, otherwise does nothing
func (o Optional[T]) OrElsePanic(v any) {
	if !o.present {
		panic(v)
	}
}

// Filter if the value is present and calling the supplied filter function returns true, returns a new optional describing the value
//
// Otherwise returns an empty optional
func (o Optional[T]) Filter(f func(v T) bool) Optional[T] {
	if o.present && f(o.value) {
		return Optional[T]{
			present: true,
			value:   o.value,
		}
	}
	return Optional[T]{
		present: false,
	}
}

// Map if the value is present and the result of calling the supplied mapping function returns non-nil, returns
// an optional describing that returned value
//
// Otherwise returns an empty optional
func (o Optional[T]) Map(f func(v T) any) Optional[any] {
	if o.present {
		v := f(o.value)
		if isPresent(v) {
			return Of(v)
		}
	}
	return Optional[any]{
		present: false,
	}
}

// MarshalJSON implements JSON marshal
//
// If the value is present, returns the marshalled data for the value
//
// Otherwise, returns the marshalled data for null
func (o Optional[T]) MarshalJSON() ([]byte, error) {
	if !o.present {
		return []byte("null"), nil
	}
	return json.Marshal(o.value)
}

// UnmarshalJSON implements JSON unmarshal
//
// if the supplied data is null representation, sets the present to false
//
// Otherwise, unmarshal the data as the value and sets the optional to present (unless the result of
// unmarshalling the value returns an error - in which case the present is set to false)
func (o *Optional[T]) UnmarshalJSON(data []byte) error {
	if len(data) == 4 && data[0] == 'n' && data[1] == 'u' && data[2] == 'l' && data[3] == 'l' {
		o.present = false
		o.value = o.emptyValue()
		o.set = true
		return nil
	}
	v := o.value
	err := json.Unmarshal(data, &v)
	if err == nil && isPresent(v) {
		o.present = true
		o.value = v
	} else {
		o.present = false
		o.value = o.emptyValue()
	}
	o.set = true
	return err
}

// Scan implements sql.Scan
func (o *Optional[T]) Scan(value interface{}) error {
	if value == nil {
		o.clear(true)
	} else if av, ok := value.(T); ok {
		o.present = true
		o.value = av
		o.set = true
	} else if ok, err := o.callScannable(value); ok {
		return err
	} else if bd, ok := value.([]byte); ok {
		var uv T
		if unErr := json.Unmarshal(bd, &uv); unErr == nil {
			if isPresent(uv) {
				o.present = true
				o.value = uv
			} else {
				o.present = false
				o.value = o.emptyValue()
			}
			o.set = true
		} else {
			o.clear(true)
			return unErr
		}
	} else {
		o.clear(true)
	}
	return nil
}

func (o *Optional[T]) callScannable(value interface{}) (bool, error) {
	var nv reflect.Value
	if !isPresent(o.value) {
		rt := reflect.TypeOf(o.value)
		if rt.Kind() == reflect.Pointer {
			rt = rt.Elem()
		}
		nv = reflect.New(rt)
	} else {
		nv = reflect.ValueOf(o.value)
	}
	anv := nv.Interface()
	if sanv, ok := anv.(sql.Scanner); ok {
		err := sanv.Scan(value)
		if err == nil {
			o.value = anv.(T)
			o.present = isPresent(anv)
			o.set = true
		} else {
			o.clear(true)
		}
		return true, err
	}
	return false, nil
}

// EmptyString returns an empty optional of type string
func EmptyString() Optional[string] {
	return Empty[string]()
}

// EmptyInterface returns an empty optional of type interface{}
func EmptyInterface() Optional[interface{}] {
	return Empty[interface{}]()
}

// EmptyInt returns an empty optional of type int
func EmptyInt() Optional[int] {
	return Empty[int]()
}

// EmptyInt8 returns an empty optional of type int8
func EmptyInt8() Optional[int8] {
	return Empty[int8]()
}

// EmptyInt16 returns an empty optional of type int16
func EmptyInt16() Optional[int16] {
	return Empty[int16]()
}

// EmptyInt32 returns an empty optional of type int32
func EmptyInt32() Optional[int32] {
	return Empty[int32]()
}

// EmptyInt64 returns an empty optional of type int64
func EmptyInt64() Optional[int64] {
	return Empty[int64]()
}

// EmptyUint returns an empty optional of type uint
func EmptyUint() Optional[uint] {
	return Empty[uint]()
}

// EmptyUint8 returns an empty optional of type uint8
func EmptyUint8() Optional[uint8] {
	return Empty[uint8]()
}

// EmptyUint16 returns an empty optional of type uint16
func EmptyUint16() Optional[uint16] {
	return Empty[uint16]()
}

// EmptyUint32 returns an empty optional of type uint32
func EmptyUint32() Optional[uint32] {
	return Empty[uint32]()
}

// EmptyUint64 returns an empty optional of type uint64
func EmptyUint64() Optional[uint64] {
	return Empty[uint64]()
}

// EmptyBool returns an empty optional of type bool
func EmptyBool() Optional[bool] {
	return Empty[bool]()
}

// EmptyFloat32 returns an empty optional of type float32
func EmptyFloat32() Optional[float32] {
	return Empty[float32]()
}

// EmptyFloat64 returns an empty optional of type float64
func EmptyFloat64() Optional[float64] {
	return Empty[float64]()
}

// EmptyByte returns an empty optional of type byte
func EmptyByte() Optional[byte] {
	return Empty[byte]()
}

// EmptyRune returns an empty optional of type rune
func EmptyRune() Optional[rune] {
	return Empty[rune]()
}
