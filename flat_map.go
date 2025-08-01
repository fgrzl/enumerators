package enumerators

import "fmt"

// flatMapEnumerator applies a function that returns an enumerator for each element,
// then flattens the results into a single sequence.
type flatMapEnumerator[T any, U any] struct {
	base    Enumerator[T]
	mapper  func(T) Enumerator[U]
	current Enumerator[U]
	err     error
}

// MoveNext advances through the flattened sequence of enumerators.
// Returns true if more elements are available, false otherwise.
func (e *flatMapEnumerator[T, U]) MoveNext() bool {
	// If we have a current enumerator, try to advance it
	for {
		if e.current != nil {
			if e.current.MoveNext() {
				return true
			}
			e.current.Dispose()
			e.current = nil
		}

		// Move to the next item in the base enumerator
		if !e.base.MoveNext() {
			e.err = e.base.Err()
			return false
		}

		// Get the next enumerator from the mapper
		item, err := e.base.Current()
		if err != nil {
			return false
		}
		e.current = e.mapper(item)
	}
}

// Current returns the current element from the active sub-enumerator.
func (e *flatMapEnumerator[T, U]) Current() (U, error) {
	if e.current == nil {
		var zero U
		return zero, fmt.Errorf("no current item")
	}
	return e.current.Current()
}

// Err returns any error encountered during enumeration.
func (e *flatMapEnumerator[T, U]) Err() error {
	return e.err
}

// Dispose cleans up resources by disposing both the base and current enumerators.
func (e *flatMapEnumerator[T, U]) Dispose() {
	e.base.Dispose()
	if e.current != nil {
		e.current.Dispose()
	}
}

// FlatMap creates an enumerator that applies a function to each element and flattens the results.
// The mapper function receives an element of type T and returns an enumerator of type U.
// All elements from each returned enumerator are yielded in sequence.
func FlatMap[T any, U any](parent Enumerator[T], mapper func(T) Enumerator[U]) Enumerator[U] {
	return &flatMapEnumerator[T, U]{
		base:   parent,
		mapper: mapper,
	}
}
