package enumerators

// Aggregate applies an accumulator function over the enumerator.
// Returns the final accumulator value and any error encountered.
// If the enumerator is empty, returns the zero value of TAccumulate and an error.
func Aggregate[TSource, TAccumulate any](enumerator Enumerator[TSource], seed TAccumulate, accumulator func(TAccumulate, TSource) (TAccumulate, error)) (TAccumulate, error) {
	defer enumerator.Dispose()

	result := seed
	hasElements := false

	for enumerator.MoveNext() {
		hasElements = true
		current, err := enumerator.Current()
		if err != nil {
			var zero TAccumulate
			return zero, err
		}

		result, err = accumulator(result, current)
		if err != nil {
			var zero TAccumulate
			return zero, err
		}
	}

	if err := enumerator.Err(); err != nil {
		var zero TAccumulate
		return zero, err
	}

	if !hasElements {
		var zero TAccumulate
		return zero, ErrEmptySequence
	}

	return result, nil
}

// AggregateWithSeed applies an accumulator function over the enumerator starting with a seed value.
// Unlike Aggregate, this always succeeds even for empty enumerators.
func AggregateWithSeed[TSource, TAccumulate any](enumerator Enumerator[TSource], seed TAccumulate, accumulator func(TAccumulate, TSource) TAccumulate) (TAccumulate, error) {
	defer enumerator.Dispose()

	result := seed
	for enumerator.MoveNext() {
		current, err := enumerator.Current()
		if err != nil {
			var zero TAccumulate
			return zero, err
		}
		result = accumulator(result, current)
	}

	return result, enumerator.Err()
}
