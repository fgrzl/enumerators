package enumerators_test

import (
	"testing"

	"github.com/fgrzl/enumerators"
	"github.com/stretchr/testify/assert"
)

func TestShouldUnionTwoEnumerators(t *testing.T) {
	// Arrange
	first := enumerators.Slice([]int{1, 2, 3})
	second := enumerators.Slice([]int{3, 4, 5})

	// Act
	unioned := enumerators.Union(first, second)
	result, err := enumerators.ToSlice(unioned)

	// Assert
	assert.NoError(t, err)
	assert.ElementsMatch(t, []int{1, 2, 3, 4, 5}, result)
}

func TestShouldIntersectTwoEnumerators(t *testing.T) {
	// Arrange
	first := enumerators.Slice([]int{1, 2, 3, 4})
	second := enumerators.Slice([]int{3, 4, 5, 6})

	// Act
	intersected := enumerators.Intersect(first, second)
	result, err := enumerators.ToSlice(intersected)

	// Assert
	assert.NoError(t, err)
	assert.ElementsMatch(t, []int{3, 4}, result)
}

func TestShouldExceptElements(t *testing.T) {
	// Arrange
	first := enumerators.Slice([]int{1, 2, 3, 4})
	second := enumerators.Slice([]int{3, 4, 5})

	// Act
	excepted := enumerators.Except(first, second)
	result, err := enumerators.ToSlice(excepted)

	// Assert
	assert.NoError(t, err)
	assert.ElementsMatch(t, []int{1, 2}, result)
}
