package enumerators

// First returns the first element of the enumerator.
// If the enumerator is empty, returns ErrEmptySequence.
// The enumerator is disposed after getting the first element.
func First[T any](enumerator Enumerator[T]) (T, error) {
	defer enumerator.Dispose()
	if enumerator.MoveNext() {
		return enumerator.Current()
	}
	var zero T
	if err := enumerator.Err(); err != nil {
		return zero, err
	}
	return zero, ErrEmptySequence
}

// Last returns the last element of the enumerator.
// If the enumerator is empty, returns ErrEmptySequence.
// The enumerator is disposed after consuming all elements.
func Last[T any](enumerator Enumerator[T]) (T, error) {
	defer enumerator.Dispose()
	var last T
	found := false
	for enumerator.MoveNext() {
		item, err := enumerator.Current()
		if err != nil {
			return last, err
		}
		last = item
		found = true
	}
	if !found {
		var zero T
		if err := enumerator.Err(); err != nil {
			return zero, err
		}
		return zero, ErrEmptySequence
	}
	return last, enumerator.Err()
}
