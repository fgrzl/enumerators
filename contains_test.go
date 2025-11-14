package enumerators_test

import (
	"testing"

	"github.com/fgrzl/enumerators"
	"github.com/stretchr/testify/assert"
)

func TestShouldReturnTrueWhenElementIsContained(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{1, 2, 3, 4})

	// Act
	result, err := enumerators.Contains(input, 3)

	// Assert
	assert.NoError(t, err)
	assert.True(t, result)
}

func TestShouldReturnFalseWhenElementIsNotContained(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{1, 2, 3, 4})

	// Act
	result, err := enumerators.Contains(input, 5)

	// Assert
	assert.NoError(t, err)
	assert.False(t, result)
}
