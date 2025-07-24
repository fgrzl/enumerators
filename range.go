package enumerators

// rangeEnumerator generates a sequence of values using a factory function.
type rangeEnumerator[T any] struct {
	start   int
	end     int
	current T
	err     error
	factory func(int) T
}

// MoveNext advances the enumerator to the next value in the range.
// It returns true if there are more values to enumerate, and false otherwise.
func (e *rangeEnumerator[T]) MoveNext() bool {
	if e.start < e.end {
		e.current = e.factory(e.start)
		e.start++
		return true
	}
	return false
}

// Current returns the current value from the enumerator.
func (e *rangeEnumerator[T]) Current() (T, error) {
	return e.current, e.err
}

// Err returns any error encountered during enumeration.
func (e *rangeEnumerator[T]) Err() error {
	return e.err
}

// Dispose cleans up resources. For rangeEnumerator, this is a no-op.
func (e *rangeEnumerator[T]) Dispose() {
	// no-op
}

// Range creates an enumerator that generates a sequence of values using a factory function.
// It starts at 'seed' and generates 'count' values by calling factory(i) for i from seed to seed+count-1.
func Range[T any](seed int, count int, factory func(i int) T) Enumerator[T] {
	return &rangeEnumerator[T]{
		start:   seed,
		end:     seed + count,
		factory: factory,
	}
}
