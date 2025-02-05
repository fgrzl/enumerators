package enumerators

import "errors"

func Chain[T any](enumerators []Enumerator[T]) Enumerator[T] {
	return &chainEnumerator[T]{enumerators: enumerators, index: 0}
}

type chainEnumerator[T any] struct {
	enumerators []Enumerator[T]
	index       int
}

// Current implements Enumerator.
func (c *chainEnumerator[T]) Current() (T, error) {
	if c.index >= len(c.enumerators) {
		var zero T
		return zero, errors.New("no current element")
	}
	return c.enumerators[c.index].Current()
}

// Dispose implements Enumerator.
func (c *chainEnumerator[T]) Dispose() {
	for _, e := range c.enumerators {
		e.Dispose()
	}
	c.enumerators = nil
}

// Err implements Enumerator.
func (c *chainEnumerator[T]) Err() error {
	if c.index >= len(c.enumerators) {
		return nil
	}
	return c.enumerators[c.index].Err()
}

// MoveNext implements Enumerator.
func (c *chainEnumerator[T]) MoveNext() bool {
	for c.index < len(c.enumerators) {
		if c.enumerators[c.index].MoveNext() {
			return true
		}
		c.index++ // Move to the next enumerator
	}
	return false
}
