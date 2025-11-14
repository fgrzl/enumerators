package enumerators

import (
	"sort"
)

// Sort sorts the elements of the enumerator using the provided less function.
// The enumerator is consumed and sorted eagerly, then returned as a new enumerator.
// The less function should return true if a < b.
func Sort[T any](enumerator Enumerator[T], less func(a, b T) bool) Enumerator[T] {
	slice, err := ToSlice(enumerator)
	if err != nil {
		return Generate(func() (T, bool, error) {
			return *new(T), false, err
		})
	}
	sort.Slice(slice, func(i, j int) bool {
		return less(slice[i], slice[j])
	})
	return Slice(slice)
}

// Reverse reverses the order of elements in the enumerator.
// The enumerator is consumed eagerly, then returned as a new enumerator in reverse order.
func Reverse[T any](enumerator Enumerator[T]) Enumerator[T] {
	slice, err := ToSlice(enumerator)
	if err != nil {
		return Generate(func() (T, bool, error) {
			return *new(T), false, err
		})
	}
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
	return Slice(slice)
}
