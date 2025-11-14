package enumerators_test

import (
	"testing"

	"github.com/fgrzl/enumerators"
	"github.com/stretchr/testify/assert"
)

func TestShouldReturnMinimumElement(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{3, 1, 4, 1, 5})

	// Act
	result, err := enumerators.Min(input)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 1, result)
}

func TestShouldReturnMaximumElement(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{3, 1, 4, 1, 5})

	// Act
	result, err := enumerators.Max(input)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 5, result)
}

func TestShouldReturnZeroValueWhenMinOnEmpty(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{})

	// Act
	result, err := enumerators.Min(input)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 0, result)
}

func TestShouldReturnZeroValueWhenMaxOnEmpty(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{})

	// Act
	result, err := enumerators.Max(input)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 0, result)
}
