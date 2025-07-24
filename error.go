package enumerators

// errorEnumerator represents an enumerator that immediately yields an error.
type errorEnumerator[T any] struct {
	err error
}

// MoveNext always returns false as error enumerators have no elements.
func (e *errorEnumerator[T]) MoveNext() bool {
	return false
}

// Current returns the zero value and the stored error.
func (e *errorEnumerator[T]) Current() (T, error) {
	var zero T
	return zero, e.err
}

// Dispose cleans up resources. For errorEnumerator, this is a no-op.
func (e *errorEnumerator[T]) Dispose() {
	// No resources to clean up
}

// Err returns the stored error.
func (e *errorEnumerator[T]) Err() error {
	return e.err
}

// Error creates an enumerator that immediately returns the specified error.
// This is useful for creating error states in enumerator chains.
func Error[T any](err error) Enumerator[T] {
	return &errorEnumerator[T]{err: err}
}
