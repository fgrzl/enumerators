package enumerators_test

import (
	"testing"

	"github.com/fgrzl/enumerators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSlice_BasicIteration(t *testing.T) {
	// Arrange
	input := []int{1, 2, 3, 4, 5}
	enumerator := enumerators.Slice(input)

	// Act
	var results []int
	for enumerator.MoveNext() {
		current, err := enumerator.Current()
		require.NoError(t, err)
		results = append(results, current)
	}

	// Assert
	assert.Equal(t, input, results)
	assert.NoError(t, enumerator.Err())
}

func TestSlice_EmptySlice(t *testing.T) {
	// Arrange
	enumerator := enumerators.Slice([]int{})

	// Act & Assert
	assert.False(t, enumerator.MoveNext())
	assert.NoError(t, enumerator.Err())
}

func TestSlice_SingleElement(t *testing.T) {
	// Arrange
	enumerator := enumerators.Slice([]string{"hello"})

	// Act & Assert
	assert.True(t, enumerator.MoveNext())
	current, err := enumerator.Current()
	require.NoError(t, err)
	assert.Equal(t, "hello", current)
	
	assert.False(t, enumerator.MoveNext())
	assert.NoError(t, enumerator.Err())
}

func TestSlice_Dispose(t *testing.T) {
	// Arrange
	enumerator := enumerators.Slice([]int{1, 2, 3})

	// Act
	enumerator.Dispose() // Should not panic

	// Assert - should still work after dispose
	assert.True(t, enumerator.MoveNext())
	current, err := enumerator.Current()
	require.NoError(t, err)
	assert.Equal(t, 1, current)
}

func TestSlice_CurrentBeforeMoveNext(t *testing.T) {
	// Arrange
	enumerator := enumerators.Slice([]int{1, 2, 3})

	// Act - calling Current before MoveNext should return zero value
	current, err := enumerator.Current()

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 0, current) // zero value for int
}

func TestToSlice_WithSliceEnumerator(t *testing.T) {
	// Arrange
	original := []int{1, 2, 3, 4, 5}
	enumerator := enumerators.Slice(original)

	// Act
	result, err := enumerators.ToSlice(enumerator)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, original, result)
}

func TestToSlice_WithOtherEnumerator(t *testing.T) {
	// Arrange - use a non-slice enumerator
	enumerator := enumerators.Range(0, 3, func(i int) int { return i * 2 })

	// Act
	result, err := enumerators.ToSlice(enumerator)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []int{0, 2, 4}, result)
}

func TestToSlice_EmptyEnumerator(t *testing.T) {
	// Arrange
	enumerator := enumerators.Slice([]int{})

	// Act
	result, err := enumerators.ToSlice(enumerator)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result)
}