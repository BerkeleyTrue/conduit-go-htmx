package maybe

type Maybe[T any] struct {
	value    T
	hasValue bool
}

func New[T any](value T) Maybe[T] {
	return Maybe[T]{
		value:    value,
		hasValue: true,
	}
}

func (m Maybe[T]) Get() (T, bool) {
	return m.value, m.hasValue
}
