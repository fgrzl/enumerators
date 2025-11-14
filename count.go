package enumerators

// Count returns the number of elements in the enumerator.
// The enumerator is disposed after counting.
func Count[T any](enumerator Enumerator[T]) (int, error) {
	defer enumerator.Dispose()
	count := 0
	for enumerator.MoveNext() {
		count++
	}
	return count, enumerator.Err()
}
