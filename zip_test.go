package enumerators_test

import (
	"testing"

	"github.com/fgrzl/enumerators"
	"github.com/stretchr/testify/assert"
)

func TestShouldZipTwoEnumerators(t *testing.T) {
	// Arrange
	first := enumerators.Slice([]int{1, 2, 3})
	second := enumerators.Slice([]string{"a", "b", "c"})

	// Act
	zipped := enumerators.Zip(first, second)
	result, err := enumerators.ToSlice(zipped)

	// Assert
	assert.NoError(t, err)
	expected := []enumerators.Pair[int, string]{
		{First: 1, Second: "a"},
		{First: 2, Second: "b"},
		{First: 3, Second: "c"},
	}
	assert.Equal(t, expected, result)
}

func TestShouldZipUntilShorterEnumerator(t *testing.T) {
	// Arrange
	first := enumerators.Slice([]int{1, 2})
	second := enumerators.Slice([]string{"a", "b", "c"})

	// Act
	zipped := enumerators.Zip(first, second)
	result, err := enumerators.ToSlice(zipped)

	// Assert
	assert.NoError(t, err)
	expected := []enumerators.Pair[int, string]{
		{First: 1, Second: "a"},
		{First: 2, Second: "b"},
	}
	assert.Equal(t, expected, result)
}
