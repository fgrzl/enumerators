package enumerators_test

import (
	"testing"

	"github.com/fgrzl/enumerators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShouldReturnEvenNumbersWhenFilteringByEven(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{1, 2, 3, 4, 5, 6})
	isEven := func(x int) bool { return x%2 == 0 }

	// Act
	filtered := enumerators.Filter(input, isEven)
	result, err := enumerators.ToSlice(filtered)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []int{2, 4, 6}, result)
}

func TestShouldReturnEmptyWhenFilteringEmptyInput(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{})
	alwaysTrue := func(x int) bool { return true }

	// Act
	filtered := enumerators.Filter(input, alwaysTrue)
	result, err := enumerators.ToSlice(filtered)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestShouldReturnEmptyWhenNoElementsMatchFilter(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{1, 3, 5, 7})
	isEven := func(x int) bool { return x%2 == 0 }

	// Act
	filtered := enumerators.Filter(input, isEven)
	result, err := enumerators.ToSlice(filtered)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestShouldReturnAllElementsWhenAllElementsMatchFilter(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{2, 4, 6, 8})
	isEven := func(x int) bool { return x%2 == 0 }

	// Act
	filtered := enumerators.Filter(input, isEven)
	result, err := enumerators.ToSlice(filtered)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []int{2, 4, 6, 8}, result)
}

func TestFilter_StringFiltering(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]string{"apple", "banana", "cherry", "date"})
	hasA := func(s string) bool {
		for _, r := range s {
			if r == 'a' {
				return true
			}
		}
		return false
	}

	// Act
	filtered := enumerators.Filter(input, hasA)
	result, err := enumerators.ToSlice(filtered)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []string{"apple", "banana", "date"}, result)
}

func TestFilter_StepByStep(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{1, 2, 3, 4, 5})
	isOdd := func(x int) bool { return x%2 == 1 }
	filtered := enumerators.Filter(input, isOdd)

	// Act & Assert
	assert.True(t, filtered.MoveNext())
	current, err := filtered.Current()
	require.NoError(t, err)
	assert.Equal(t, 1, current)

	assert.True(t, filtered.MoveNext())
	current, err = filtered.Current()
	require.NoError(t, err)
	assert.Equal(t, 3, current)

	assert.True(t, filtered.MoveNext())
	current, err = filtered.Current()
	require.NoError(t, err)
	assert.Equal(t, 5, current)

	assert.False(t, filtered.MoveNext())
	assert.NoError(t, filtered.Err())
}

func TestFilter_Dispose(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{1, 2, 3})
	alwaysTrue := func(x int) bool { return true }
	filtered := enumerators.Filter(input, alwaysTrue)

	// Act
	filtered.Dispose() // Should not panic

	// Assert - should still work (slice enumerator doesn't care about dispose)
	assert.True(t, filtered.MoveNext())
}

func TestShouldPropagateErrorWhenChainedWithMap(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{1, 2, 3, 4, 5})
	mapper := func(x int) (int, error) {
		if x == 3 {
			return 0, assert.AnError
		}
		return x * 2, nil
	}
	isEven := func(x int) bool { return x%2 == 0 }

	// Act
	mapped := enumerators.Map(input, mapper)
	filtered := enumerators.Filter(mapped, isEven)
	result, err := enumerators.ToSlice(filtered)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, assert.AnError, err)
	assert.Equal(t, []int{2, 4}, result) // 1->2 (even), 2->4 (even), 3->error
}
