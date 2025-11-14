package enumerators

// Contains returns true if the enumerator contains the specified value.
// The enumerator is disposed after checking.
func Contains[T comparable](enumerator Enumerator[T], value T) (bool, error) {
	defer enumerator.Dispose()
	for enumerator.MoveNext() {
		item, err := enumerator.Current()
		if err != nil {
			return false, err
		}
		if item == value {
			return true, nil
		}
	}
	return false, enumerator.Err()
}
