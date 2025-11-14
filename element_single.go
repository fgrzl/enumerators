package enumerators

import "errors"

// ElementAt returns the element at the specified zero-based index.
// If index is out of bounds, returns zero value and error.
// The enumerator is disposed after finding the element.
func ElementAt[T any](enumerator Enumerator[T], index int) (T, error) {
	defer enumerator.Dispose()
	i := 0
	for enumerator.MoveNext() {
		if i == index {
			return enumerator.Current()
		}
		i++
	}
	var zero T
	if err := enumerator.Err(); err != nil {
		return zero, err
	}
	return zero, errors.New("index out of bounds")
}

// Single returns the single element if the enumerator has exactly one element.
// If empty, returns ErrEmptySequence.
// If more than one, returns the first and error.
// The enumerator is disposed after checking.
func Single[T any](enumerator Enumerator[T]) (T, error) {
	defer enumerator.Dispose()
	if !enumerator.MoveNext() {
		var zero T
		if err := enumerator.Err(); err != nil {
			return zero, err
		}
		return zero, ErrEmptySequence
	}
	first, err := enumerator.Current()
	if err != nil {
		return first, err
	}
	if enumerator.MoveNext() {
		return first, errors.New("enumerator has more than one element")
	}
	return first, enumerator.Err()
}
