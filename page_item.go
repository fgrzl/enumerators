package enumerators

import (
	"fmt"
)

func PageItemEnumerator[T any](fetchPage func() ([]T, bool, error)) Enumerator[T] {
	return &pageItemEnumerator[T]{fetchPage: fetchPage}
}

type pageItemEnumerator[T any] struct {
	fetchPage func() ([]T, bool, error) // returns next page of items, hasMore, error
	items     []T
	index     int
	hasMore   bool
	err       error
	disposed  bool
}

func (e *pageItemEnumerator[T]) MoveNext() bool {
	if e.disposed || e.err != nil {
		return false
	}

	e.index++

	// current page has more items
	if e.index < len(e.items) {
		return true
	}

	// need next page
	if !e.hasMore && len(e.items) > 0 {
		return false
	}

	e.items, e.hasMore, e.err = e.fetchPage()
	e.index = 0

	return len(e.items) > 0 && e.err == nil
}

func (e *pageItemEnumerator[T]) Current() (T, error) {
	if e.index >= len(e.items) {
		var zero T
		return zero, fmt.Errorf("no current item")
	}
	return e.items[e.index], nil
}

func (e *pageItemEnumerator[T]) Err() error {
	return e.err
}

func (e *pageItemEnumerator[T]) Dispose() {
	e.disposed = true
}
