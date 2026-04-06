package enumerators_test

import (
	"fmt"
	"slices"
	"sync"
	"testing"

	"github.com/fgrzl/enumerators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShouldSampleElements(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{1, 2, 3, 4, 5})

	// Act
	sampled := enumerators.Sample(input, 3)
	result, err := enumerators.ToSlice(sampled)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 3)
	for _, value := range result {
		assert.Contains(t, []int{1, 2, 3, 4, 5}, value)
	}
}

func TestShouldSampleAllWhenNExceedsLength(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{1, 2, 3})

	// Act
	sampled := enumerators.Sample(input, 5)
	result, err := enumerators.ToSlice(sampled)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 3)
	assert.ElementsMatch(t, []int{1, 2, 3}, result)
}

func TestShouldReturnEmptySampleWhenNIsZero(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{1, 2, 3})

	// Act
	result, err := enumerators.ToSlice(enumerators.Sample(input, 0))

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestShouldReturnEmptySampleWhenNIsNegative(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{1, 2, 3})

	// Act
	result, err := enumerators.ToSlice(enumerators.Sample(input, -1))

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestShouldReturnEmptySampleForEmptyInput(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{})

	// Act
	result, err := enumerators.ToSlice(enumerators.Sample(input, 3))

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestShouldSampleConcurrently(t *testing.T) {
	// Arrange
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	const goroutineCount = 12
	const iterationsPerGoroutine = 20

	var wg sync.WaitGroup
	errCh := make(chan error, goroutineCount)

	// Act
	for idx := 0; idx < goroutineCount; idx++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for iteration := 0; iteration < iterationsPerGoroutine; iteration++ {
				result, err := enumerators.ToSlice(enumerators.Sample(enumerators.Slice(input), 4))
				if err != nil {
					errCh <- err
					return
				}
				if len(result) != 4 {
					errCh <- fmt.Errorf("unexpected sample length: %d", len(result))
					return
				}
				for _, value := range result {
					if !slices.Contains(input, value) {
						errCh <- fmt.Errorf("unexpected sampled value: %d", value)
						return
					}
				}
			}
		}()
	}

	wg.Wait()
	close(errCh)

	// Assert
	for err := range errCh {
		require.NoError(t, err)
	}
}
