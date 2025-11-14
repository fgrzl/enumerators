package enumerators

// skipEnumerator skips the first N elements from the base enumerator.
type skipEnumerator[T any] struct {
	base    Enumerator[T]
	toSkip  int
	skipped int
	current T
	err     error
}

// MoveNext advances to the next element after skipping the specified number.
// Returns true if more elements are available, false otherwise.
func (e *skipEnumerator[T]) MoveNext() bool {
	for e.skipped < e.toSkip {
		if !e.base.MoveNext() {
			e.err = e.base.Err()
			return false
		}
		e.skipped++
	}

	if !e.base.MoveNext() {
		e.err = e.base.Err()
		return false
	}

	e.current, e.err = e.base.Current()
	return e.err == nil
}

// Current returns the current element and any error encountered.
func (e *skipEnumerator[T]) Current() (T, error) {
	return e.current, e.err
}

// Err returns any error encountered during enumeration.
func (e *skipEnumerator[T]) Err() error {
	return e.err
}

// Dispose cleans up resources by disposing the underlying enumerator.
func (e *skipEnumerator[T]) Dispose() {
	e.base.Dispose()
}

// Skip creates an enumerator that skips the first N elements.
// If N is negative or zero, returns the original enumerator without wrapping.
func Skip[T any](enumerator Enumerator[T], n int) Enumerator[T] {
	if n <= 0 {
		return enumerator
	}
	return &skipEnumerator[T]{
		base:   enumerator,
		toSkip: n,
	}
}
