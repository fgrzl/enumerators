package enumerators_test

import (
	"testing"

	"github.com/fgrzl/enumerators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTakeWhile_BasicTakeWhile(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{1, 2, 3, 4, 5, 1, 2})
	lessThanFour := func(x int) bool { return x < 4 }

	// Act
	taken := enumerators.TakeWhile(input, lessThanFour)
	result, err := enumerators.ToSlice(taken)

	// Assert
	assert.NoError(t, err)
	// Should stop at 4 and return only elements before it
	assert.Equal(t, []int{1, 2, 3}, result)
}

func TestTakeWhile_EmptyInput(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{})
	alwaysTrue := func(x int) bool { return true }

	// Act
	taken := enumerators.TakeWhile(input, alwaysTrue)
	result, err := enumerators.ToSlice(taken)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestTakeWhile_NoMatches(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{5, 6, 7, 8})
	lessThanFive := func(x int) bool { return x < 5 }

	// Act
	taken := enumerators.TakeWhile(input, lessThanFive)
	result, err := enumerators.ToSlice(taken)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result) // First element doesn't match, so none taken
}

func TestTakeWhile_AllMatch(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{1, 2, 3, 4})
	lessThanTen := func(x int) bool { return x < 10 }

	// Act
	taken := enumerators.TakeWhile(input, lessThanTen)
	result, err := enumerators.ToSlice(taken)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3, 4}, result)
}

func TestTakeWhile_StringCondition(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]string{"a", "b", "c", "d", "A", "B"})
	isLowercase := func(s string) bool { return len(s) == 1 && s[0] >= 'a' && s[0] <= 'z' }

	// Act
	taken := enumerators.TakeWhile(input, isLowercase)
	result, err := enumerators.ToSlice(taken)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []string{"a", "b", "c", "d"}, result) // Stops at first uppercase
}

func TestTakeWhile_StepByStep(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{2, 4, 6, 7, 8})
	isEven := func(x int) bool { return x%2 == 0 }
	taken := enumerators.TakeWhile(input, isEven)

	// Act & Assert
	assert.True(t, taken.MoveNext())
	current, err := taken.Current()
	require.NoError(t, err)
	assert.Equal(t, 2, current)

	assert.True(t, taken.MoveNext())
	current, err = taken.Current()
	require.NoError(t, err)
	assert.Equal(t, 4, current)

	assert.True(t, taken.MoveNext())
	current, err = taken.Current()
	require.NoError(t, err)
	assert.Equal(t, 6, current)

	// Should stop at 7 (odd number)
	assert.False(t, taken.MoveNext())
	assert.NoError(t, taken.Err())
}

func TestTakeWhile_SingleElement(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{5})
	alwaysTrue := func(x int) bool { return true }

	// Act
	taken := enumerators.TakeWhile(input, alwaysTrue)
	result, err := enumerators.ToSlice(taken)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []int{5}, result)
}

func TestTakeWhile_Dispose(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{1, 2, 3})
	alwaysTrue := func(x int) bool { return true }
	taken := enumerators.TakeWhile(input, alwaysTrue)

	// Act
	taken.Dispose() // Should not panic

	// Assert - should still work
	assert.True(t, taken.MoveNext())
	current, err := taken.Current()
	require.NoError(t, err)
	assert.Equal(t, 1, current)
}