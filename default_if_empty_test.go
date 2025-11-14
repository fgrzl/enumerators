package enumerators_test

import (
	"testing"

	"github.com/fgrzl/enumerators"
	"github.com/stretchr/testify/assert"
)

func TestShouldReturnOriginalElementsWhenNotEmpty(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{1, 2, 3})

	// Act
	result, err := enumerators.ToSlice(enumerators.DefaultIfEmpty(input, 42))

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3}, result)
}

func TestShouldReturnDefaultValueWhenEmpty(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{})

	// Act
	result, err := enumerators.ToSlice(enumerators.DefaultIfEmpty(input, 42))

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []int{42}, result)
}
