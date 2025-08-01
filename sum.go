package enumerators

import (
	"golang.org/x/exp/constraints"
)

// Sum consumes an enumerator and returns the sum of all elements after applying a selector function.
// The enumerator is automatically disposed after computation.
// Returns an error if enumeration fails or if the selector function returns an error.
func Sum[T any, TSum constraints.Ordered](enumerator Enumerator[T], selector func(item T) (TSum, error)) (TSum, error) {
	defer enumerator.Dispose()
	var sum TSum
	for enumerator.MoveNext() {

		item, err := enumerator.Current()
		if err != nil {
			var zero TSum
			return zero, err
		}

		value, err := selector(item)
		if err != nil {
			var zero TSum
			return zero, err
		}
		sum += value
	}

	return sum, nil
}
