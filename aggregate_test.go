package enumerators_test

import (
	"testing"

	"github.com/fgrzl/enumerators"
	"github.com/stretchr/testify/assert"
)

func TestShouldAggregateElementsWithSeed(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{1, 2, 3, 4})

	// Act
	result, err := enumerators.Aggregate(input, 0, func(acc, x int) (int, error) {
		return acc + x, nil
	})

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 10, result)
}

func TestShouldReturnErrorForEmptySequenceInAggregate(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{})

	// Act
	_, err := enumerators.Aggregate(input, 0, func(acc, x int) (int, error) {
		return acc + x, nil
	})

	// Assert
	assert.Error(t, err)
	assert.Equal(t, enumerators.ErrEmptySequence, err)
}

func TestShouldAggregateWithSeedEvenForEmptySequence(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{})

	// Act
	result, err := enumerators.AggregateWithSeed(input, 42, func(acc, x int) int {
		return acc + x
	})

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 42, result)
}
