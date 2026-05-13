package enumerators_test

import (
	"testing"

	"github.com/fgrzl/enumerators"
	"github.com/stretchr/testify/assert"
)

func TestShouldReturnElementAtIndex(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{10, 20, 30, 40})

	// Act
	result, err := enumerators.ElementAt(input, 2)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 30, result)
}

func TestShouldReturnErrorWhenIndexOutOfBounds(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{10, 20})

	// Act
	_, err := enumerators.ElementAt(input, 5)

	// Assert
	assert.Error(t, err)
}

func TestShouldReturnSingleElement(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{42})

	// Act
	result, err := enumerators.Single(input)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 42, result)
}

func TestShouldReturnErrorWhenSingleHasMultipleElements(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{1, 2})

	// Act
	_, err := enumerators.Single(input)

	// Assert
	assert.Error(t, err)
}

func TestShouldReturnErrorWhenSingleIsEmpty(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{})

	// Act
	_, err := enumerators.Single(input)

	// Assert
	assert.Error(t, err)
}
