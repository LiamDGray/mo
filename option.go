package mo

import (
	"fmt"
	"reflect"
)

var optionNoSuchElement = fmt.Errorf("no such element")

// Some builds an Option when value is present.
func Some[T any](value T) Option[T] {
	return Option[T]{
		isPresent: true,
		value:     value,
	}
}

// None builds an Option when value is absent.
func None[T any]() Option[T] {
	return Option[T]{
		isPresent: false,
	}
}

// TupleToOption builds a Some Option when second argument is true, or None.
func TupleToOption[T any](value T, ok bool) Option[T] {
	if ok {
		return Some(value)
	}
	return None[T]()
}

// EmptyableToOption builds a Some Option when value is not empty, or None.
func EmptyableToOption[T any](value T) Option[T] {
	// 🤮
	isZero := reflect.ValueOf(&value).Elem().IsZero()
	if isZero {
		return None[T]()
	}

	return Some(value)
}

// Option is a container for an optional value of type T. If value exists, Option is
// of type Some. If the value is absent, Option is of type None.
type Option[T any] struct {
	isPresent bool
	value     T
}

// IsPresent returns true when value is absent.
func (o Option[T]) IsPresent() bool {
	return o.isPresent
}

// IsAbsent returns true when value is present.
func (o Option[T]) IsAbsent() bool {
	return !o.isPresent
}

// Size returns 1 when value is present or 0 instead.
func (o Option[T]) Size() int {
	if o.isPresent {
		return 1
	}

	return 0
}

// Get returns value and presence.
func (o Option[T]) Get() (T, bool) {
	if !o.isPresent {
		return empty[T](), false
	}

	return o.value, true
}

// MustGet returns value if present or panics instead.
func (o Option[T]) MustGet() T {
	if !o.isPresent {
		panic(optionNoSuchElement)
	}

	return o.value
}

// OrElse returns value if present or default value.
func (o Option[T]) OrElse(fallback T) T {
	if !o.isPresent {
		return fallback
	}

	return o.value
}

// OrEmpty returns value if present or empty value.
func (o Option[T]) OrEmpty() T {
	return o.value
}

// ForEach executes the given side-effecting function of value is present.
func (o Option[T]) ForEach(onValue func(value T)) {
	if o.isPresent {
		onValue(o.value)
	}
}

// Match executes the first function if value is present and second function if absent.
// It returns a new Option.
func (o Option[T]) Match(onValue func(value T) (T, bool), onNone func() (T, bool)) Option[T] {
	if o.isPresent {
		return TupleToOption(onValue(o.value))
	}
	return TupleToOption(onNone())
}

// Map executes the mapper function if value is present or returns None if absent.
func (o Option[T]) Map(mapper func(value T) (T, bool)) Option[T] {
	if o.isPresent {
		return TupleToOption(mapper(o.value))
	}

	return None[T]()
}

// MapNone executes the mapper function if value is absent or returns Option.
func (o Option[T]) MapNone(mapper func() (T, bool)) Option[T] {
	if o.isPresent {
		return Some(o.value)
	}

	return TupleToOption(mapper())
}

// FlatMap executes the mapper function if value is present or returns None if absent.
func (o Option[T]) FlatMap(mapper func(value T) Option[T]) Option[T] {
	if o.isPresent {
		return mapper(o.value)
	}

	return None[T]()
}
