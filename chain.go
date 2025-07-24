package enumerators

import "errors"

// Chain creates an enumerator that sequentially chains multiple enumerators together.
// Elements from each enumerator are yielded in order until that enumerator is exhausted,
// then the next enumerator in the chain begins yielding elements.
func Chain[T any](enumerators ...Enumerator[T]) Enumerator[T] {
	return &chainEnumerator[T]{enumerators: enumerators, index: 0}
}

// chainEnumerator chains multiple enumerators together, yielding their elements sequentially.
type chainEnumerator[T any] struct {
	enumerators []Enumerator[T]
	index       int
}

// Current returns the current element from the active enumerator in the chain.
func (c *chainEnumerator[T]) Current() (T, error) {
	if c.index >= len(c.enumerators) {
		var zero T
		return zero, errors.New("no current element")
	}
	return c.enumerators[c.index].Current()
}

// Dispose cleans up resources by disposing all enumerators in the chain.
func (c *chainEnumerator[T]) Dispose() {
	for _, e := range c.enumerators {
		e.Dispose()
	}
	c.enumerators = nil
}

// Err returns any error from the currently active enumerator.
func (c *chainEnumerator[T]) Err() error {
	if c.index >= len(c.enumerators) {
		return nil
	}
	return c.enumerators[c.index].Err()
}

// MoveNext advances through the chain of enumerators, moving to the next enumerator when the current one is exhausted.
func (c *chainEnumerator[T]) MoveNext() bool {
	for c.index < len(c.enumerators) {
		if c.enumerators[c.index].MoveNext() {
			return true
		}
		c.index++ // Move to the next enumerator
	}
	return false
}
