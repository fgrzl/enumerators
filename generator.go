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
	done      bool
	err       error
	dispose   sync.Once
	disposed  bool
}

// Create a new generator.
func NewGenerator[T any](next func() (T, bool, error)) Enumerator[T] {
	return &Generator[T]{onNext: next}
}

func NewGeneratorWithDispose[T any](next func() (T, bool, error), dispose func()) Enumerator[T] {
	return &Generator[T]{
		onNext:    next,
		onDispose: dispose,
	}
}

// Dispose cleans up the enumerator.
func (ce *Generator[T]) Dispose() {
	ce.done = true
	ce.dispose.Do(func() {
		if ce.onDispose != nil {
			ce.onDispose()
		}
		ce.disposed = true
	})
}

// MoveNext generates the next value.
func (ce *Generator[T]) MoveNext() bool {
	if ce.done {
		return false
	}

	// Generate the next value
	ce.current, ce.done, ce.err = ce.onNext()

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
