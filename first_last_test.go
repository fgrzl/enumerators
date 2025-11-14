package enumerators_test

import (
	"testing"

	"github.com/fgrzl/enumerators"
	"github.com/stretchr/testify/assert"
)

func TestShouldReturnFirstElementWhenSliceHasElements(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{1, 2, 3})

	// Act
	result, err := enumerators.First(input)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 1, result)
}

func TestShouldReturnLastElementWhenSliceHasElements(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{1, 2, 3})

	// Act
	result, err := enumerators.Last(input)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 3, result)
}

func TestShouldReturnErrorWhenFirstOnEmpty(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{})

	// Act
	_, err := enumerators.First(input)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, enumerators.ErrEmptySequence, err)
}

func TestShouldReturnErrorWhenLastOnEmpty(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{})

	// Act
	_, err := enumerators.Last(input)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, enumerators.ErrEmptySequence, err)
}
