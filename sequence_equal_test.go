package enumerators_test

import (
	"testing"

	"github.com/fgrzl/enumerators"
	"github.com/stretchr/testify/assert"
)

func TestShouldReturnTrueForEqualSequences(t *testing.T) {
	// Arrange
	first := enumerators.Slice([]int{1, 2, 3})
	second := enumerators.Slice([]int{1, 2, 3})

	// Act
	equal, err := enumerators.SequenceEqual(first, second)

	// Assert
	assert.NoError(t, err)
	assert.True(t, equal)
}

func TestShouldReturnFalseForDifferentSequences(t *testing.T) {
	// Arrange
	first := enumerators.Slice([]int{1, 2, 3})
	second := enumerators.Slice([]int{1, 2, 4})

	// Act
	equal, err := enumerators.SequenceEqual(first, second)

	// Assert
	assert.NoError(t, err)
	assert.False(t, equal)
}

func TestShouldReturnFalseForDifferentLengths(t *testing.T) {
	// Arrange
	first := enumerators.Slice([]int{1, 2, 3})
	second := enumerators.Slice([]int{1, 2})

	// Act
	equal, err := enumerators.SequenceEqual(first, second)

	// Assert
	assert.NoError(t, err)
	assert.False(t, equal)
}
