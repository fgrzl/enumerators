package enumerators

import (
	"golang.org/x/exp/constraints"
)

// Min returns the minimum element in the enumerator.
// If the enumerator is empty, returns the zero value.
// The enumerator is disposed after finding the minimum.
func Min[T constraints.Ordered](enumerator Enumerator[T]) (T, error) {
	defer enumerator.Dispose()
	var min T
	found := false
	for enumerator.MoveNext() {
		item, err := enumerator.Current()
		if err != nil {
			return min, err
		}
		if !found || item < min {
			min = item
			found = true
		}
	}
	if !found {
		var zero T
		return zero, enumerator.Err()
	}
	return min, enumerator.Err()
}

// Max returns the maximum element in the enumerator.
// If the enumerator is empty, returns the zero value.
// The enumerator is disposed after finding the maximum.
func Max[T constraints.Ordered](enumerator Enumerator[T]) (T, error) {
	defer enumerator.Dispose()
	var max T
	found := false
	for enumerator.MoveNext() {
		item, err := enumerator.Current()
		if err != nil {
			return max, err
		}
		if !found || item > max {
			max = item
			found = true
		}
	}
	if !found {
		var zero T
		return zero, enumerator.Err()
	}
	return max, enumerator.Err()
}
