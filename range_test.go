package enumerators_test

import (
	"testing"

	"github.com/fgrzl/enumerators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShouldGenerateSequenceWhenCreatingBasicRange(t *testing.T) {
	// Arrange 
	generator := func(i int) int { return i * 2 }

	// Act
	rangeEnum := enumerators.Range(0, 5, generator)
	result, err := enumerators.ToSlice(rangeEnum)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []int{0, 2, 4, 6, 8}, result)
}

func TestShouldReturnEmptyWhenRangeCountIsZero(t *testing.T) {
	// Arrange
	generator := func(i int) int { return i }

	// Act
	rangeEnum := enumerators.Range(0, 0, generator)
	result, err := enumerators.ToSlice(rangeEnum)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestShouldGenerateOneItemWhenRangeCountIsOne(t *testing.T) {
	// Arrange
	generator := func(i int) string { return "item-" + string(rune('0'+i)) }

	// Act
	rangeEnum := enumerators.Range(5, 1, generator)
	result, err := enumerators.ToSlice(rangeEnum)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []string{"item-5"}, result)
}

func TestShouldReturnEmptyWhenRangeCountIsNegative(t *testing.T) {
	// Arrange 
	generator := func(i int) int { return i }

	// Act - negative count should produce empty result
	rangeEnum := enumerators.Range(0, -5, generator)
	result, err := enumerators.ToSlice(rangeEnum)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestRange_StringGeneration(t *testing.T) {
	// Arrange & Act
	rangeEnum := enumerators.Range(1, 3, func(i int) string {
		return "value-" + string(rune('0'+i))
	})
	result, err := enumerators.ToSlice(rangeEnum)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []string{"value-1", "value-2", "value-3"}, result)
}

func TestRange_StepByStep(t *testing.T) {
	// Arrange
	rangeEnum := enumerators.Range(10, 3, func(i int) int { return i + 100 })

	// Act & Assert
	assert.True(t, rangeEnum.MoveNext())
	current, err := rangeEnum.Current()
	require.NoError(t, err)
	assert.Equal(t, 110, current)

	assert.True(t, rangeEnum.MoveNext())
	current, err = rangeEnum.Current()
	require.NoError(t, err)
	assert.Equal(t, 111, current)

	assert.True(t, rangeEnum.MoveNext())
	current, err = rangeEnum.Current()
	require.NoError(t, err)
	assert.Equal(t, 112, current)

	assert.False(t, rangeEnum.MoveNext())
	assert.NoError(t, rangeEnum.Err())
}

func TestRange_CurrentBeforeMoveNext(t *testing.T) {
	// Arrange
	rangeEnum := enumerators.Range(0, 3, func(i int) int { return i })

	// Act - calling Current before MoveNext should return zero value
	current, err := rangeEnum.Current()

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 0, current) // zero value for int
}

func TestRange_Dispose(t *testing.T) {
	// Arrange
	rangeEnum := enumerators.Range(0, 3, func(i int) int { return i })

	// Act
	rangeEnum.Dispose() // Should not panic

	// Assert - should still work
	assert.True(t, rangeEnum.MoveNext())
	current, err := rangeEnum.Current()
	require.NoError(t, err)
	assert.Equal(t, 0, current)
}