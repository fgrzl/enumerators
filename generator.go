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
	if !ce.hasNext {
		return false
	}

	// Generate the next value
	ce.current, ce.hasNext, ce.err = ce.onNext()

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
