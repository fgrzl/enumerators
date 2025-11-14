package enumerators

// takeWhileEnumerator yields elements from the base enumerator while a condition is satisfied.
type takeWhileEnumerator[T any] struct {
	base      Enumerator[T]
	condition func(T) bool
	current   T
	err       error
}

// MoveNext advances to the next element while the condition is satisfied.
// Stops when the condition fails for any element.
func (e *takeWhileEnumerator[T]) MoveNext() bool {
	if !e.base.MoveNext() {
		e.err = e.base.Err() // Propagate any error from base enumerator
		return false
	}

	item, err := e.base.Current()
	if err != nil {
		e.err = err
		return false
	}

	if !e.condition(item) {
		return false // Stop taking when condition fails
	}

	e.current = item
	return true
}

// Current returns the current element that satisfied the condition.
func (e *takeWhileEnumerator[T]) Current() (T, error) {
	return e.current, e.err
}

// Err returns any error encountered during enumeration.
func (e *takeWhileEnumerator[T]) Err() error {
	return e.err
}

// Dispose cleans up resources by disposing the underlying enumerator.
func (e *takeWhileEnumerator[T]) Dispose() {
	e.base.Dispose()
}

// TakeWhile creates an enumerator that yields elements while the specified condition is satisfied.
// Enumeration stops as soon as any element fails the condition.
func TakeWhile[T any](enumerator Enumerator[T], condition func(T) bool) Enumerator[T] {
	return &takeWhileEnumerator[T]{
		base:      enumerator,
		condition: condition,
	}
}
