package enumerators_test

import (
	"testing"

	"github.com/fgrzl/enumerators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPeekable_BasicPeek(t *testing.T) {
	// Arrange
	base := enumerators.Slice([]int{1, 2, 3})
	peekable := enumerators.Peekable(base)

	// Act & Assert
	// Peek without advancing
	value, hasNext, err := peekable.Peek()
	assert.NoError(t, err)
	assert.True(t, hasNext)
	assert.Equal(t, 1, value)

	// Peek again should return same value
	value, hasNext, err = peekable.Peek()
	assert.NoError(t, err)
	assert.True(t, hasNext)
	assert.Equal(t, 1, value)

	// MoveNext should advance to the peeked value
	assert.True(t, peekable.MoveNext())
	current, err := peekable.Current()
	require.NoError(t, err)
	assert.Equal(t, 1, current)
}

func TestPeekable_PeekAndAdvance(t *testing.T) {
	// Arrange
	base := enumerators.Slice([]int{10, 20, 30})
	peekable := enumerators.Peekable(base)

	// Act & Assert
	// Peek first value
	value, hasNext, err := peekable.Peek()
	assert.NoError(t, err)
	assert.True(t, hasNext)
	assert.Equal(t, 10, value)

	// Move to first value
	assert.True(t, peekable.MoveNext())
	current, err := peekable.Current()
	require.NoError(t, err)
	assert.Equal(t, 10, current)

	// Peek second value
	value, hasNext, err = peekable.Peek()
	assert.NoError(t, err)
	assert.True(t, hasNext)
	assert.Equal(t, 20, value)

	// Move to second value
	assert.True(t, peekable.MoveNext())
	current, err = peekable.Current()
	require.NoError(t, err)
	assert.Equal(t, 20, current)
}

func TestPeekable_HasNext(t *testing.T) {
	// Arrange
	base := enumerators.Slice([]int{1, 2})
	peekable := enumerators.Peekable(base)

	// Act & Assert
	assert.True(t, peekable.HasNext())
	
	assert.True(t, peekable.MoveNext())
	assert.True(t, peekable.HasNext())
	
	assert.True(t, peekable.MoveNext())
	assert.False(t, peekable.HasNext())
}

func TestPeekable_EmptyEnumerator(t *testing.T) {
	// Arrange
	base := enumerators.Empty[int]()
	peekable := enumerators.Peekable(base)

	// Act & Assert
	value, hasNext, err := peekable.Peek()
	assert.NoError(t, err)
	assert.False(t, hasNext)
	assert.Equal(t, 0, value) // zero value for int

	assert.False(t, peekable.HasNext())
	assert.False(t, peekable.MoveNext())
}

func TestPeekable_SingleElement(t *testing.T) {
	// Arrange
	base := enumerators.Slice([]string{"hello"})
	peekable := enumerators.Peekable(base)

	// Act & Assert
	value, hasNext, err := peekable.Peek()
	assert.NoError(t, err)
	assert.True(t, hasNext)
	assert.Equal(t, "hello", value)

	assert.True(t, peekable.HasNext())
	
	assert.True(t, peekable.MoveNext())
	current, err := peekable.Current()
	require.NoError(t, err)
	assert.Equal(t, "hello", current)

	assert.False(t, peekable.HasNext())
	assert.False(t, peekable.MoveNext())
}

func TestPeekable_ToSlice(t *testing.T) {
	// Arrange
	base := enumerators.Slice([]int{1, 2, 3, 4, 5})
	peekable := enumerators.Peekable(base)

	// Peek first value to test interaction
	value, hasNext, err := peekable.Peek()
	assert.NoError(t, err)
	assert.True(t, hasNext)
	assert.Equal(t, 1, value)

	// Act
	result, err := enumerators.ToSlice(peekable)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, result)
}

func TestPeekable_WithFilter(t *testing.T) {
	// Arrange
	base := enumerators.Slice([]int{1, 2, 3, 4, 5, 6})
	filtered := enumerators.Filter(base, func(x int) bool { return x%2 == 0 })
	peekable := enumerators.Peekable(filtered)

	// Act & Assert
	value, hasNext, err := peekable.Peek()
	assert.NoError(t, err)
	assert.True(t, hasNext)
	assert.Equal(t, 2, value) // First even number

	assert.True(t, peekable.MoveNext())
	current, err := peekable.Current()
	require.NoError(t, err)
	assert.Equal(t, 2, current)

	value, hasNext, err = peekable.Peek()
	assert.NoError(t, err)
	assert.True(t, hasNext)
	assert.Equal(t, 4, value) // Next even number

	result, err := enumerators.ToSlice(peekable)
	assert.NoError(t, err)
	assert.Equal(t, []int{4, 6}, result) // Remaining even numbers
}

func TestPeekable_PeekAfterComplete(t *testing.T) {
	// Arrange
	base := enumerators.Slice([]int{1})
	peekable := enumerators.Peekable(base)

	// Move through the single element
	assert.True(t, peekable.MoveNext())
	assert.False(t, peekable.MoveNext())

	// Act & Assert - peek after completion
	value, hasNext, err := peekable.Peek()
	assert.NoError(t, err)
	assert.False(t, hasNext)
	assert.Equal(t, 0, value) // zero value for int
}

func TestPeekable_Dispose(t *testing.T) {
	// Arrange
	disposed := false
	base := enumerators.Cleanup(
		enumerators.Slice([]int{1, 2, 3}),
		func() { disposed = true },
	)
	peekable := enumerators.Peekable(base)

	// Act
	peekable.Dispose()

	// Assert
	assert.True(t, disposed, "Base enumerator should be disposed")
}

func TestPeekable_ErrorHandling(t *testing.T) {
	// Arrange
	base := enumerators.Slice([]int{1, 2, 3})
	peekable := enumerators.Peekable(base)

	// Act & Assert - Err should delegate to base
	assert.NoError(t, peekable.Err())
	
	// Peek and move should still work
	assert.True(t, peekable.HasNext())
	assert.True(t, peekable.MoveNext())
}