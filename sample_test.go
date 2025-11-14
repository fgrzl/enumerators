package enumerators_test

import (
	"testing"

	"github.com/fgrzl/enumerators"
	"github.com/stretchr/testify/assert"
)

func TestShouldSampleElements(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{1, 2, 3, 4, 5})

	// Act
	sampled := enumerators.Sample(input, 3)
	result, err := enumerators.ToSlice(sampled)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 3)
	// Check all elements are from original
	for _, v := range result {
		assert.Contains(t, []int{1, 2, 3, 4, 5}, v)
	}
}

func TestShouldSampleAllWhenNExceedsLength(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{1, 2, 3})

	// Act
	sampled := enumerators.Sample(input, 5)
	result, err := enumerators.ToSlice(sampled)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 3)
	assert.ElementsMatch(t, []int{1, 2, 3}, result)
}
