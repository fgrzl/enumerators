package enumerators

// Any returns true if at least one element satisfies the predicate.
// If the enumerator is empty, returns false.
// The enumerator is disposed after checking.
func Any[T any](enumerator Enumerator[T], predicate func(T) bool) (bool, error) {
	defer enumerator.Dispose()
	for enumerator.MoveNext() {
		item, err := enumerator.Current()
		if err != nil {
			return false, err
		}
		if predicate(item) {
			return true, nil
		}
	}
	return false, enumerator.Err()
}

// All returns true if all elements satisfy the predicate.
// If the enumerator is empty, returns true.
// The enumerator is disposed after checking.
func All[T any](enumerator Enumerator[T], predicate func(T) bool) (bool, error) {
	defer enumerator.Dispose()
	for enumerator.MoveNext() {
		item, err := enumerator.Current()
		if err != nil {
			return false, err
		}
		if !predicate(item) {
			return false, nil
		}
	}
	return true, enumerator.Err()
}
