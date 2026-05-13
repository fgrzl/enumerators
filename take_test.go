package enumerators_test

import (
	"testing"

	"github.com/fgrzl/enumerators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShouldReturnFirstNElementsWhenTakingFromSlice(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{1, 2, 3, 4, 5})

	// Act
	taken := enumerators.Take(input, 3)
	result, err := enumerators.ToSlice(taken)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3}, result)
}

func TestShouldReturnEmptyWhenTakeCountIsZero(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{1, 2, 3, 4, 5})

	// Act
	taken := enumerators.Take(input, 0)
	result, err := enumerators.ToSlice(taken)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestShouldReturnAllElementsWhenTakeCountExceedsAvailable(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{1, 2, 3})

	// Act
	taken := enumerators.Take(input, 10)
	result, err := enumerators.ToSlice(taken)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3}, result)
}

func TestShouldReturnEmptyWhenInputIsEmpty(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{})

	// Act
	taken := enumerators.Take(input, 5)
	result, err := enumerators.ToSlice(taken)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestShouldReturnEmptyWhenTakeCountIsNegative(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{1, 2, 3, 4, 5})

	// Act
	taken := enumerators.Take(input, -1)
	result, err := enumerators.ToSlice(taken)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result) // negative count should result in empty
}

func TestShouldStopAfterSpecifiedCountWhenIteratingStepByStep(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]string{"a", "b", "c", "d", "e"})
	taken := enumerators.Take(input, 3)

	// Act - First element
	hasFirst := taken.MoveNext()
	firstCurrent, firstErr := taken.Current()

	// Act - Second element
	hasSecond := taken.MoveNext()
	secondCurrent, secondErr := taken.Current()

	// Act - Third element
	hasThird := taken.MoveNext()
	thirdCurrent, thirdErr := taken.Current()

	// Act - Should stop after 3 items
	hasFourth := taken.MoveNext()

	// Assert
	assert.True(t, hasFirst)
	require.NoError(t, firstErr)
	assert.Equal(t, "a", firstCurrent)

	assert.True(t, hasSecond)
	require.NoError(t, secondErr)
	assert.Equal(t, "b", secondCurrent)

	assert.True(t, hasThird)
	require.NoError(t, thirdErr)
	assert.Equal(t, "c", thirdCurrent)

	assert.False(t, hasFourth) // Should stop after 3 items
	assert.NoError(t, taken.Err())
}

func TestShouldReturnSingleElementWhenTakingOneFromSingleElementSlice(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{42})

	// Act
	taken := enumerators.Take(input, 1)
	result, err := enumerators.ToSlice(taken)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []int{42}, result)
}

func TestShouldStillFunctionWhenDisposeCalledOnTake(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{1, 2, 3})
	taken := enumerators.Take(input, 2)

	// Act
	taken.Dispose() // Should not panic
	hasNext := taken.MoveNext()
	current, err := taken.Current()

	// Assert - should still work
	assert.True(t, hasNext)
	require.NoError(t, err)
	assert.Equal(t, 1, current)
}
