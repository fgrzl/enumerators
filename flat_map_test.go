package enumerators_test

import (
	"testing"

	"github.com/fgrzl/enumerators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFlatMap_BasicFlattening(t *testing.T) {
	// Arrange
	input := enumerators.Slice([][]int{{1, 2}, {3, 4}, {5}})

	// Act
	flattened := enumerators.FlatMap(input, enumerators.Slice[int])
	result, err := enumerators.ToSlice(flattened)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, result)
}

func TestFlatMap_EmptyInput(t *testing.T) {
	// Arrange
	input := enumerators.Slice([][]int{})

	// Act
	flattened := enumerators.FlatMap(input, enumerators.Slice[int])
	result, err := enumerators.ToSlice(flattened)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestFlatMap_EmptySubSequences(t *testing.T) {
	// Arrange
	input := enumerators.Slice([][]int{{1, 2}, {}, {3, 4}})

	// Act
	flattened := enumerators.FlatMap(input, enumerators.Slice[int])
	result, err := enumerators.ToSlice(flattened)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3, 4}, result) // Empty slice is skipped
}

func TestFlatMap_SingleElement(t *testing.T) {
	// Arrange
	input := enumerators.Slice([][]string{{"hello", "world"}})

	// Act
	flattened := enumerators.FlatMap(input, enumerators.Slice[string])
	result, err := enumerators.ToSlice(flattened)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []string{"hello", "world"}, result)
}

func TestFlatMap_StringToChars(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]string{"hi", "go"})
	stringToChars := func(s string) enumerators.Enumerator[rune] {
		runes := []rune(s)
		return enumerators.Slice(runes)
	}

	// Act
	flattened := enumerators.FlatMap(input, stringToChars)
	result, err := enumerators.ToSlice(flattened)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []rune{'h', 'i', 'g', 'o'}, result)
}

func TestFlatMap_NumberToRange(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{2, 3})
	numberToRange := func(n int) enumerators.Enumerator[int] {
		return enumerators.Range(0, n, func(i int) int { return i })
	}

	// Act
	flattened := enumerators.FlatMap(input, numberToRange)
	result, err := enumerators.ToSlice(flattened)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []int{0, 1, 0, 1, 2}, result) // Range(0,2) + Range(0,3)
}

func TestFlatMap_StepByStep(t *testing.T) {
	// Arrange
	input := enumerators.Slice([][]int{{10, 20}, {30}})
	flattened := enumerators.FlatMap(input, enumerators.Slice[int])

	// Act & Assert
	assert.True(t, flattened.MoveNext())
	current, err := flattened.Current()
	require.NoError(t, err)
	assert.Equal(t, 10, current)

	assert.True(t, flattened.MoveNext())
	current, err = flattened.Current()
	require.NoError(t, err)
	assert.Equal(t, 20, current)

	assert.True(t, flattened.MoveNext())
	current, err = flattened.Current()
	require.NoError(t, err)
	assert.Equal(t, 30, current)

	assert.False(t, flattened.MoveNext())
	assert.NoError(t, flattened.Err())
}

func TestFlatMap_DisposesSubEnumerators(t *testing.T) {
	// Arrange
	disposed1 := false
	disposed2 := false

	input := enumerators.Slice([]int{1, 2})
	createEnum := func(n int) enumerators.Enumerator[int] {
		var cleanup func()
		if n == 1 {
			cleanup = func() { disposed1 = true }
		} else {
			cleanup = func() { disposed2 = true }
		}
		return enumerators.Cleanup(
			enumerators.Slice([]int{n * 10}),
			cleanup,
		)
	}

	// Act
	flattened := enumerators.FlatMap(input, createEnum)
	result, err := enumerators.ToSlice(flattened)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []int{10, 20}, result)
	assert.True(t, disposed1, "First sub-enumerator should be disposed")
	assert.True(t, disposed2, "Second sub-enumerator should be disposed")
}

func TestFlatMap_Dispose(t *testing.T) {
	// Arrange
	input := enumerators.Slice([][]int{{1, 2}, {3, 4}})
	flattened := enumerators.FlatMap(input, enumerators.Slice[int])

	// Move to start working with sub-enumerator
	assert.True(t, flattened.MoveNext())

	// Act
	flattened.Dispose() // Should dispose base and current sub-enumerator

	// Assert - should not panic, and the functionality should still work for testing
	// (though in practice, using after dispose is not recommended)
	current, err := flattened.Current()
	require.NoError(t, err)
	assert.Equal(t, 1, current)
}

func TestFlatMap_CurrentBeforeMoveNext(t *testing.T) {
	// Arrange
	input := enumerators.Slice([][]int{{1, 2}})
	flattened := enumerators.FlatMap(input, enumerators.Slice[int])

	// Act - calling Current before MoveNext should return error
	current, err := flattened.Current()

	// Assert
	assert.Error(t, err)
	assert.Equal(t, 0, current) // zero value for int
	assert.Contains(t, err.Error(), "no current item")
}
