package enumerators_test

import (
	"testing"

	"github.com/fgrzl/enumerators"
	"github.com/stretchr/testify/assert"
)

func TestShouldSkipFirstNElementsWhenSkippingFromSlice(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{1, 2, 3, 4, 5})

	// Act
	skipped := enumerators.Skip(input, 2)
	result, err := enumerators.ToSlice(skipped)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []int{3, 4, 5}, result)
}

func TestShouldReturnAllElementsWhenSkipCountExceedsLength(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{1, 2, 3})

	// Act
	skipped := enumerators.Skip(input, 10)
	result, err := enumerators.ToSlice(skipped)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestShouldReturnAllElementsWhenSkipCountIsZero(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{1, 2, 3})

	// Act
	skipped := enumerators.Skip(input, 0)
	result, err := enumerators.ToSlice(skipped)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3}, result)
}
