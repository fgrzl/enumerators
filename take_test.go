package enumerators_test

import (
	"testing"

	"github.com/fgrzl/enumerators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTake_BasicTake(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{1, 2, 3, 4, 5})

	// Act
	taken := enumerators.Take(input, 3)
	result, err := enumerators.ToSlice(taken)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3}, result)
}

func TestTake_TakeZero(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{1, 2, 3, 4, 5})

	// Act
	taken := enumerators.Take(input, 0)
	result, err := enumerators.ToSlice(taken)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestTake_TakeMoreThanAvailable(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{1, 2, 3})

	// Act
	taken := enumerators.Take(input, 10)
	result, err := enumerators.ToSlice(taken)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3}, result)
}

func TestTake_EmptyInput(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{})

	// Act
	taken := enumerators.Take(input, 5)
	result, err := enumerators.ToSlice(taken)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestTake_NegativeCount(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{1, 2, 3, 4, 5})

	// Act
	taken := enumerators.Take(input, -1)
	result, err := enumerators.ToSlice(taken)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result) // negative count should result in empty
}

func TestTake_StepByStep(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]string{"a", "b", "c", "d", "e"})
	taken := enumerators.Take(input, 3)

	// Act & Assert
	assert.True(t, taken.MoveNext())
	current, err := taken.Current()
	require.NoError(t, err)
	assert.Equal(t, "a", current)

	assert.True(t, taken.MoveNext())
	current, err = taken.Current()
	require.NoError(t, err)
	assert.Equal(t, "b", current)

	assert.True(t, taken.MoveNext())
	current, err = taken.Current()
	require.NoError(t, err)
	assert.Equal(t, "c", current)

	assert.False(t, taken.MoveNext()) // Should stop after 3 items
	assert.NoError(t, taken.Err())
}

func TestTake_SingleElement(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{42})

	// Act
	taken := enumerators.Take(input, 1)
	result, err := enumerators.ToSlice(taken)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []int{42}, result)
}

func TestTake_Dispose(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{1, 2, 3})
	taken := enumerators.Take(input, 2)

	// Act
	taken.Dispose() // Should not panic

	// Assert - should still work
	assert.True(t, taken.MoveNext())
	current, err := taken.Current()
	require.NoError(t, err)
	assert.Equal(t, 1, current)
}