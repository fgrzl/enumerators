package enumerators_test

import (
	"testing"

	"github.com/fgrzl/enumerators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSkipIf_BasicSkipping(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{1, 2, 3, 4, 5, 6})
	isEven := func(x int) bool { return x%2 == 0 }

	// Act
	skipped := enumerators.SkipIf(input, isEven)
	result, err := enumerators.ToSlice(skipped)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 3, 5}, result) // Skip even numbers
}

func TestSkipIf_EmptyInput(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{})
	alwaysTrue := func(x int) bool { return true }

	// Act
	skipped := enumerators.SkipIf(input, alwaysTrue)
	result, err := enumerators.ToSlice(skipped)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestSkipIf_SkipAll(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{1, 2, 3, 4, 5})
	alwaysTrue := func(x int) bool { return true }

	// Act
	skipped := enumerators.SkipIf(input, alwaysTrue)
	result, err := enumerators.ToSlice(skipped)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result) // All elements skipped
}

func TestSkipIf_SkipNone(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{1, 2, 3, 4, 5})
	alwaysFalse := func(x int) bool { return false }

	// Act
	skipped := enumerators.SkipIf(input, alwaysFalse)
	result, err := enumerators.ToSlice(skipped)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, result) // No elements skipped
}

func TestSkipIf_StringSkipping(t *testing.T) {
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
	skipped := enumerators.SkipIf(input, hasA)
	result, err := enumerators.ToSlice(skipped)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []string{"cherry"}, result) // Only "cherry" has no 'a'
}

func TestSkipIf_StepByStep(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{1, 2, 3, 4, 5, 6})
	isEven := func(x int) bool { return x%2 == 0 }
	skipped := enumerators.SkipIf(input, isEven)

	// Act & Assert
	assert.True(t, skipped.MoveNext())
	current, err := skipped.Current()
	require.NoError(t, err)
	assert.Equal(t, 1, current)

	assert.True(t, skipped.MoveNext())
	current, err = skipped.Current()
	require.NoError(t, err)
	assert.Equal(t, 3, current)

	assert.True(t, skipped.MoveNext())
	current, err = skipped.Current()
	require.NoError(t, err)
	assert.Equal(t, 5, current)

	assert.False(t, skipped.MoveNext())
	assert.NoError(t, skipped.Err())
}

func TestSkipIf_SingleElement(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{5})
	isEven := func(x int) bool { return x%2 == 0 }

	// Act
	skipped := enumerators.SkipIf(input, isEven)
	result, err := enumerators.ToSlice(skipped)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []int{5}, result) // 5 is odd, so not skipped
}

func TestSkipIf_SingleElementSkipped(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{4})
	isEven := func(x int) bool { return x%2 == 0 }

	// Act
	skipped := enumerators.SkipIf(input, isEven)
	result, err := enumerators.ToSlice(skipped)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result) // 4 is even, so skipped
}

func TestSkipIf_Dispose(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{1, 2, 3})
	alwaysFalse := func(x int) bool { return false }
	skipped := enumerators.SkipIf(input, alwaysFalse)

	// Act
	skipped.Dispose() // Should not panic

	// Assert - should still work
	assert.True(t, skipped.MoveNext())
	current, err := skipped.Current()
	require.NoError(t, err)
	assert.Equal(t, 1, current)
}