package enumerators_test

import (
	"testing"

	"github.com/fgrzl/enumerators"
	"github.com/stretchr/testify/assert"
)

func TestChainEnumerator_Err(t *testing.T) {
	// Arrange
	enumerator1 := enumerators.Slice([]int{1, 2, 3})
	enumerator2 := enumerators.Slice([]int{4, 5, 6})

	// Act
	chainEnum := enumerators.Chain(enumerator1, enumerator2)

	// Assert
	slice, err := enumerators.ToSlice(chainEnum)
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, slice)
}
