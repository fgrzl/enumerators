package enumerators

import (
	"math/rand"
	"sync"
	"time"
)

var globalRand = rand.New(rand.NewSource(time.Now().UnixNano()))
var globalRandMu sync.Mutex

// Sample returns a random sample of n elements from the enumerator.
// If n <= 0, it returns an empty enumerator.
// If n >= length, it returns all elements in random order.
// The enumerator is consumed eagerly.
func Sample[T any](enumerator Enumerator[T], n int) Enumerator[T] {
	slice, err := ToSlice(enumerator)
	if err != nil {
		return Generate(func() (T, bool, error) {
			return *new(T), false, err
		})
	}

	if n <= 0 {
		return Empty[T]()
	}

	sampleRand := nextSampleRand()
	if n >= len(slice) {
		// Shuffle all
		sampleRand.Shuffle(len(slice), func(i, j int) {
			slice[i], slice[j] = slice[j], slice[i]
		})
		return Slice(slice)
	}

	// Reservoir sampling for n items.
	result := make([]T, n)
	copy(result, slice[:n])
	for i := n; i < len(slice); i++ {
		j := sampleRand.Intn(i + 1)
		if j < n {
			result[j] = slice[i]
		}
	}

	return Slice(result)
}

func nextSampleRand() *rand.Rand {
	globalRandMu.Lock()
	defer globalRandMu.Unlock()

	return rand.New(rand.NewSource(globalRand.Int63()))
}
