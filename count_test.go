package enumerators_test

import (
	"testing"

	"github.com/fgrzl/enumerators"
	"github.com/stretchr/testify/assert"
)

func TestShouldReturnCountOfElements(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{1, 2, 3, 4, 5})

	// Act
	result, err := enumerators.Count(input)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 5, result)
}

func TestShouldReturnZeroWhenCountOnEmpty(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{})

	// Act
	result, err := enumerators.Count(input)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 0, result)
}
