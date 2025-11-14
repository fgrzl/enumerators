package enumerators_test

import (
	"testing"

	"github.com/fgrzl/enumerators"
	"github.com/stretchr/testify/assert"
)

func TestShouldSortElementsInAscendingOrder(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{3, 1, 4, 1, 5})

	// Act
	sorted := enumerators.Sort(input, func(a, b int) bool { return a < b })
	result, err := enumerators.ToSlice(sorted)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 1, 3, 4, 5}, result)
}

func TestShouldReverseElements(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{1, 2, 3, 4, 5})

	// Act
	reversed := enumerators.Reverse(input)
	result, err := enumerators.ToSlice(reversed)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []int{5, 4, 3, 2, 1}, result)
}
