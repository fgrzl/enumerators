package enumerators

import (
	"cmp"
)

// Min returns the minimum element in the enumerator.
// If the enumerator is empty, returns the zero value.
// The enumerator is disposed after finding the minimum.
func Min[T cmp.Ordered](enumerator Enumerator[T]) (T, error) {
	defer enumerator.Dispose()
	var minVal T
	found := false
	for enumerator.MoveNext() {
		item, err := enumerator.Current()
		if err != nil {
			return minVal, err
		}
		if !found || item < minVal {
			minVal = item
			found = true
		}
	}
	if !found {
		var zero T
		return zero, enumerator.Err()
	}
	return minVal, enumerator.Err()
}

// Max returns the maximum element in the enumerator.
// If the enumerator is empty, returns the zero value.
// The enumerator is disposed after finding the maximum.
func Max[T cmp.Ordered](enumerator Enumerator[T]) (T, error) {
	defer enumerator.Dispose()
	var maxVal T
	found := false
	for enumerator.MoveNext() {
		item, err := enumerator.Current()
		if err != nil {
			return maxVal, err
		}
		if !found || item > maxVal {
			maxVal = item
			found = true
		}
	}
	if !found {
		var zero T
		return zero, enumerator.Err()
	}
	return maxVal, enumerator.Err()
}
