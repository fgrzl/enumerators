package enumerators

// filterEnumerator filters elements from the base enumerator based on a predicate function.
type filterEnumerator[T any] struct {
	base    Enumerator[T]
	filter  func(T) bool
	current T
	err     error
}

// MoveNext advances to the next element that satisfies the filter condition.
// Returns true if such an element is found, false otherwise.
func (e *filterEnumerator[T]) MoveNext() bool {
	for {
		if !e.base.MoveNext() {
			e.err = e.base.Err() // Propagate any error from base enumerator
			return false
		}

		item, err := e.base.Current()
		if err != nil {
			e.err = err
			return false
		}

		if !e.filter(item) {
			continue
		}

		e.current = item
		return true
	}
}

// Current returns the current filtered element and any error encountered.
func (e *filterEnumerator[T]) Current() (T, error) {
	return e.current, e.err
}

// Err returns any error encountered during enumeration.
func (e *filterEnumerator[T]) Err() error {
	return e.err
}

// Dispose cleans up resources by disposing the underlying enumerator.
func (e *filterEnumerator[T]) Dispose() {
	e.base.Dispose()
}

// Filter creates an enumerator that only yields elements satisfying the predicate function.
// The filter function receives an element of type T and returns true if the element should be included.
func Filter[T any](parent Enumerator[T], filter func(T) bool) Enumerator[T] {
	return &filterEnumerator[T]{
		base:   parent,
		filter: filter,
	}
}
