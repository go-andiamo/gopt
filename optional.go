package gopt

import (
	"encoding/json"
	"errors"
	"reflect"
)

type String *string

func Of[T any](value T) Optional[T] {
	return Optional[T]{
		present: isPresent(value),
		value:   value,
	}
}

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

func Empty[T any]() Optional[T] {
	return Optional[T]{
		present: false,
	}
}

func isPresent(v any) bool {
	vo := reflect.ValueOf(v)
	vk := vo.Kind()
	if vk == reflect.Ptr {
		vk = vo.Elem().Kind()
		if vk == reflect.Invalid {
			return false
		}
	}
	return v != nil
}

type Optional[T any] struct {
	present bool
	value   T
}

func (o Optional[T]) Get() (T, error) {
	if !o.present {
		return o.value, errors.New("not present")
	}
	return o.value, nil
}

func (o Optional[T]) Empty() Optional[T] {
	return Optional[T]{
		present: false,
	}
}

func (o Optional[T]) IsPresent() bool {
	return o.present
}

func (o Optional[T]) IfPresent(f func(v T)) {
	if o.present {
		f(o.value)
	}
}

func (o Optional[T]) OrElse(other T) T {
	if o.present {
		return o.value
	}
	return other
}

func (o Optional[T]) OrElseGet(f func() T) T {
	if o.present {
		return o.value
	}
	return f()
}

func (o Optional[T]) OrElseError(f func() error) error {
	if !o.present {
		return f()
	}
	return nil
}

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

func (o Optional[T]) Map(f func(v T) any) Optional[any] {
	if o.present {
		v := f(o.value)
		if v != nil {
			return Of(v)
		}
	}
	return Optional[any]{
		present: false,
	}
}

func (o Optional[T]) MarshalJSON() ([]byte, error) {
	if !o.present {
		return []byte("null"), nil
	}
	return json.Marshal(o.value)
}

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
