package enumerators_test

import (
	"testing"

	"github.com/fgrzl/enumerators"
	"github.com/stretchr/testify/assert"
)

func TestShouldReturnUniqueElementsWhenDistinctApplied(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{1, 2, 2, 3, 1, 4})

	// Act
	distinct := enumerators.Distinct(input)
	result, err := enumerators.ToSlice(distinct)

	// Assert
	assert.NoError(t, err)
	assert.ElementsMatch(t, []int{1, 2, 3, 4}, result)
}

func TestShouldReturnEmptyWhenDistinctOnEmpty(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{})

	// Act
	distinct := enumerators.Distinct(input)
	result, err := enumerators.ToSlice(distinct)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result)
}
