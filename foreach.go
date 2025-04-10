package enumerators

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
