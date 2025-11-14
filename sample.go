package enumerators

import (
	"math/rand"
	"time"
)

var globalRand = rand.New(rand.NewSource(time.Now().UnixNano()))

// Sample returns a random sample of n elements from the enumerator.
// If n >= length, returns all elements in random order.
// The enumerator is consumed eagerly.
func Sample[T any](enumerator Enumerator[T], n int) Enumerator[T] {
	slice, err := ToSlice(enumerator)
	if err != nil {
		return Generate(func() (T, bool, error) {
			return *new(T), false, err
		})
	}
	if n >= len(slice) {
		// Shuffle all
		globalRand.Shuffle(len(slice), func(i, j int) {
			slice[i], slice[j] = slice[j], slice[i]
		})
		return Slice(slice)
	}
	// Reservoir sampling for n items
	result := make([]T, n)
	copy(result, slice[:n])
	for i := n; i < len(slice); i++ {
		j := globalRand.Intn(i + 1)
		if j < n {
			result[j] = slice[i]
		}
	}
	return Slice(result)
}
