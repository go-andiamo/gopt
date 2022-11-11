package gopt

import (
	"database/sql"
	"encoding/json"
	"errors"
	"reflect"
)

type String *string

var (
	_EmptyString    = Optional[string]{}
	_EmptyInterface = Optional[interface{}]{}
	_EmptyInt       = Optional[int]{}
	_EmptyInt8      = Optional[int8]{}
	_EmptyInt16     = Optional[int16]{}
	_EmptyInt32     = Optional[int32]{}
	_EmptyInt64     = Optional[int64]{}
	_EmptyUint      = Optional[uint]{}
	_EmptyUint8     = Optional[uint8]{}
	_EmptyUint16    = Optional[uint16]{}
	_EmptyUint32    = Optional[uint32]{}
	_EmptyUint64    = Optional[uint64]{}
	_EmptyBool      = Optional[bool]{}
	_EmptyFloat32   = Optional[float32]{}
	_EmptyFloat64   = Optional[float64]{}
	_EmptyByte      = Optional[byte]{}
	_EmptyRune      = Optional[rune]{}
)

var (
	EmptyString    = _EmptyString
	EmptyInterface = _EmptyInterface
	EmptyInt       = _EmptyInt
	EmptyInt8      = _EmptyInt8
	EmptyInt16     = _EmptyInt16
	EmptyInt32     = _EmptyInt32
	EmptyInt64     = _EmptyInt64
	EmptyUint      = _EmptyUint
	EmptyUint8     = _EmptyUint8
	EmptyUint16    = _EmptyUint16
	EmptyUint32    = _EmptyUint32
	EmptyUint64    = _EmptyUint64
	EmptyBool      = _EmptyBool
	EmptyFloat32   = _EmptyFloat32
	EmptyFloat64   = _EmptyFloat64
	EmptyByte      = _EmptyByte
	EmptyRune      = _EmptyRune
)

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
}

// Get returns the value and an error if the value is not present
func (o Optional[T]) Get() (T, error) {
	if !o.present {
		return o.value, errors.New("not present")
	}
	return o.value, nil
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

// IfPresent if the value is present, calls the supplied function with the value, otherwise does nothing
func (o Optional[T]) IfPresent(f func(v T)) Optional[T] {
	if o.present {
		f(o.value)
	}
	return o
}

// IfPresentOtherwise if the value is present, calls the supplied function with the value, otherwise calls the other function
func (o Optional[T]) IfPresentOtherwise(f func(v T), other func()) Optional[T] {
	if o.present {
		f(o.value)
	} else {
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
	if !o.present && isPresent(v) {
		o.present = true
		o.value = v
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

// DoWith if the value is present, calls the supplied function with the value
//
// Returns the original optional
func (o Optional[T]) DoWith(f func(v T)) Optional[T] {
	if o.present {
		f(o.value)
	}
	return o
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
		return nil
	}
	v := o.value
	err := json.Unmarshal(data, &v)
	if err == nil {
		o.present = true
		o.value = v
	} else {
		o.present = false
	}
	return err
}

// Scan implements sql.Scan
func (o *Optional[T]) Scan(value interface{}) error {
	if value == nil {
		o.present = false
	} else if av, ok := value.(T); ok {
		o.present = true
		o.value = av
	} else if ok, err := o.callScannable(value); ok {
		return err
	} else if bd, ok := value.([]byte); ok {
		var uv T
		if err := json.Unmarshal(bd, &uv); err == nil {
			o.present = true
			o.value = uv
		} else {
			o.present = false
		}
	} else {
		o.present = false
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
			o.present = true
		}
		return true, err
	}
	return false, nil
}
