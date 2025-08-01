package enumerators

// takeEnumerator limits the number of elements yielded from the base enumerator.
type takeEnumerator[T any] struct {
	base    Enumerator[T]
	size    int
	count   int
	current T
	err     error
}

// MoveNext advances to the next element, up to the specified limit.
func (e *takeEnumerator[T]) MoveNext() bool {
	if !e.base.MoveNext() {
		return false
	}

	item, err := e.base.Current()
	if err != nil {
		e.err = err
		return false
	}

	if e.count >= e.size {
		return false
	}

	e.current = item
	e.count++
	return true
}

// Current returns the current element and any error encountered.
func (e *takeEnumerator[T]) Current() (T, error) {
	return e.current, e.err
}

// Err returns any error encountered during enumeration.
func (e *takeEnumerator[T]) Err() error {
	return e.err
}

// Dispose cleans up resources by disposing the underlying enumerator.
func (e *takeEnumerator[T]) Dispose() {
	e.base.Dispose()
}

// Take creates an enumerator that yields at most the specified number of elements.
// If take is negative or zero, the enumerator will yield no elements.
func Take[T any](enumerator Enumerator[T], take int) Enumerator[T] {
	return &takeEnumerator[T]{
		base: enumerator,
		size: take,
	}
}
