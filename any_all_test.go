package enumerators_test

import (
	"testing"

	"github.com/fgrzl/enumerators"
	"github.com/stretchr/testify/assert"
)

func TestShouldReturnTrueWhenAnyElementMatches(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{1, 2, 3, 4})

	// Act
	result, err := enumerators.Any(input, func(x int) bool { return x > 2 })

	// Assert
	assert.NoError(t, err)
	assert.True(t, result)
}

func TestShouldReturnFalseWhenNoElementMatches(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{1, 2, 3})

	// Act
	result, err := enumerators.Any(input, func(x int) bool { return x > 10 })

	// Assert
	assert.NoError(t, err)
	assert.False(t, result)
}

func TestShouldReturnTrueWhenAllElementsMatch(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{2, 4, 6})

	// Act
	result, err := enumerators.All(input, func(x int) bool { return x%2 == 0 })

	// Assert
	assert.NoError(t, err)
	assert.True(t, result)
}

func TestShouldReturnFalseWhenNotAllElementsMatch(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{2, 3, 4})

	// Act
	result, err := enumerators.All(input, func(x int) bool { return x%2 == 0 })

	// Assert
	assert.NoError(t, err)
	assert.False(t, result)
}
