package enumerators

// distinctEnumerator yields unique elements from the base enumerator.
type distinctEnumerator[T comparable] struct {
	base    Enumerator[T]
	seen    map[T]bool
	current T
	err     error
}

// MoveNext advances to the next unique element.
// Returns true if more unique elements are available, false otherwise.
func (e *distinctEnumerator[T]) MoveNext() bool {
	for {
		if !e.base.MoveNext() {
			e.err = e.base.Err()
			return false
		}

		item, err := e.base.Current()
		if err != nil {
			e.err = err
			return false
		}

		if !e.seen[item] {
			e.seen[item] = true
			e.current = item
			return true
		}
	}
}

// Current returns the current unique element and any error encountered.
func (e *distinctEnumerator[T]) Current() (T, error) {
	return e.current, e.err
}

// Err returns any error encountered during enumeration.
func (e *distinctEnumerator[T]) Err() error {
	return e.err
}

// Dispose cleans up resources by disposing the underlying enumerator.
func (e *distinctEnumerator[T]) Dispose() {
	e.base.Dispose()
}

// Distinct creates an enumerator that yields unique elements.
// Elements are compared using ==, so T must be comparable.
func Distinct[T comparable](enumerator Enumerator[T]) Enumerator[T] {
	return &distinctEnumerator[T]{
		base: enumerator,
		seen: make(map[T]bool),
	}
}
