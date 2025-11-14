package enumerators_test

import (
	"testing"

	"github.com/fgrzl/enumerators"
	"github.com/stretchr/testify/assert"
)

func TestShouldPartitionElementsBasedOnPredicate(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{1, 2, 3, 4, 5, 6})

	// Act
	even, odd := enumerators.Partition(input, func(x int) bool { return x%2 == 0 })
	evenSlice, err1 := enumerators.ToSlice(even)
	oddSlice, err2 := enumerators.ToSlice(odd)

	// Assert
	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.Equal(t, []int{2, 4, 6}, evenSlice)
	assert.Equal(t, []int{1, 3, 5}, oddSlice)
}
