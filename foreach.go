package enumerators

// ForEach iterates over all elements in the enumerator, calling the provided action function for each element.
// The enumerator is automatically disposed after iteration completes or when an error occurs.
// Returns the first error encountered during iteration or from the action function.
func ForEach[T any](enumerator Enumerator[T], do func(T) error) error {
	defer enumerator.Dispose()
	for enumerator.MoveNext() {
		item, err := enumerator.Current()
		if err != nil {
			return err
		}

		if err := do(item); err != nil {
			return err
		}
	}

	return enumerator.Err()
}
