package enumerators

// SliceEnumerator provides enumeration over a Go slice.
type SliceEnumerator[T any] struct {
	slice   []T
	cursor  int
	current T
	err     error
}

// MoveNext advances the enumerator to the next element in the slice.
// Returns true if more elements are available, false otherwise.
func (e *SliceEnumerator[T]) MoveNext() bool {
	e.cursor++
	if e.cursor >= len(e.slice) {
		return false
	}
	e.current = e.slice[e.cursor]
	return true
}

// Current returns the current element and any error encountered.
func (e *SliceEnumerator[T]) Current() (T, error) {
	return e.current, e.err
}

// Err returns any error encountered during enumeration.
func (e *SliceEnumerator[T]) Err() error {
	return e.err
}

// Dispose cleans up resources. For SliceEnumerator, this is a no-op.
func (enumerator *SliceEnumerator[T]) Dispose() {
	// no-op
}

// Slice creates a new enumerator that iterates over the provided slice.
func Slice[T any](slice []T) Enumerator[T] {
	return &SliceEnumerator[T]{
		slice:  slice,
		cursor: -1,
	}
}

// ToSlice converts an enumerator to a slice by consuming all its elements.
// The enumerator is automatically disposed after consumption.
func ToSlice[T any](enumerator Enumerator[T]) ([]T, error) {
	defer enumerator.Dispose()

	var slice []T
	for enumerator.MoveNext() {
		item, err := enumerator.Current()
		if err != nil {
			return slice, err
		}

		slice = append(slice, item)
	}

	return slice, enumerator.Err()
}
