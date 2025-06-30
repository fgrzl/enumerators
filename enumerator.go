// Package enumerators provides generic iterator utilities such as chunking,
// mapping, and cleanup behavior for stream-like processing.
package enumerators

// Enumerator represents a generic iterator over a sequence of values of type T.
// It should be disposed when no longer needed.
type Enumerator[T any] interface {
	Disposable
	MoveNext() bool
	Current() (T, error)
	Err() error
}

// Consume advances the enumerator to completion and disposes it.
// It returns any error encountered during iteration.
func Consume[T any](e Enumerator[T]) error {
	defer e.Dispose()
	for e.MoveNext() {
		// Intentionally empty; just consume.
	}
	return e.Err()
}

// Empty returns an enumerator that yields no elements.
func Empty[T any]() Enumerator[T] {
	return &emptyEnumerator[T]{}
}

type emptyEnumerator[T any] struct{}

func (e *emptyEnumerator[T]) MoveNext() bool      { return false }
func (e *emptyEnumerator[T]) Current() (T, error) { var zero T; return zero, nil }
func (e *emptyEnumerator[T]) Err() error          { return nil }
func (e *emptyEnumerator[T]) Dispose()            {}

// Cleanup wraps an enumerator with a cleanup function that runs on Dispose.
func Cleanup[T any](enumerator Enumerator[T], cleanup func()) *cleanupEnumerator[T] {
	return &cleanupEnumerator[T]{base: enumerator, cleanup: cleanup}
}

type cleanupEnumerator[T any] struct {
	base        Enumerator[T]
	cleanup     func()
	cleanupDone bool
}

func (e *cleanupEnumerator[T]) MoveNext() bool      { return e.base.MoveNext() }
func (e *cleanupEnumerator[T]) Current() (T, error) { return e.base.Current() }
func (e *cleanupEnumerator[T]) Err() error          { return e.base.Err() }

// Dispose ensures the underlying enumerator is disposed and the cleanup function is run once.
func (e *cleanupEnumerator[T]) Dispose() {
	if !e.cleanupDone {
		e.cleanupDone = true
		e.base.Dispose()
		if e.cleanup != nil {
			e.cleanup()
		}
	}
}
