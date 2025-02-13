package enumerators

import (
	"errors"
	"sync"
)

// Generator generates values continuously.
type Generator[T any] struct {
	Enumerator[T]
	onNext    func() (T, bool, error)
	onDispose func()
	current   T
	hasNext   bool
	err       error
	dispose   sync.Once
	disposed  bool
}

// Create a new generator.
func Generate[T any](next func() (T, bool, error)) Enumerator[T] {
	return &Generator[T]{onNext: next}
}

func GenerateAndDispose[T any](next func() (T, bool, error), dispose func()) Enumerator[T] {
	return &Generator[T]{
		onNext:    next,
		onDispose: dispose,
	}
}

// Dispose cleans up the enumerator.
func (ce *Generator[T]) Dispose() {
	ce.dispose.Do(func() {
		ce.hasNext = false
		if ce.onDispose != nil {
			ce.onDispose()
		}
		ce.disposed = true
	})
}

// MoveNext generates the next value.
func (ce *Generator[T]) MoveNext() bool {
	// Generate the next value
	ce.current, ce.hasNext, ce.err = ce.onNext()

	if !ce.hasNext {
		return false
	}

	return ce.err == nil
}

// Current returns the current value or an error if disposed.
func (ce *Generator[T]) Current() (T, error) {
	if ce.disposed {
		var zero T
		return zero, errors.New("enumerator disposed")
	}
	return ce.current, ce.err
}

// Err returns the last error.
func (ce *Generator[T]) Err() error {
	return ce.err
}

type KeyValuePair[K comparable, V any] struct {
	Key   K
	Value V
}

func GenerateFromMap[K comparable, V any](m map[K]V) Enumerator[*KeyValuePair[K, V]] {
	keys := make([]K, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}

	index := 0
	return Generate(func() (*KeyValuePair[K, V], bool, error) {
		if index >= len(keys) {
			return nil, false, nil
		}

		key := keys[index]
		index++

		return &KeyValuePair[K, V]{Key: key, Value: m[key]}, true, nil
	})
}
