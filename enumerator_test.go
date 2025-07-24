package enumerators_test

import (
	"testing"

	"github.com/fgrzl/enumerators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmpty_NoElements(t *testing.T) {
	// Arrange
	enumerator := enumerators.Empty[int]()

	// Act & Assert
	assert.False(t, enumerator.MoveNext())
	assert.NoError(t, enumerator.Err())
	
	current, err := enumerator.Current()
	assert.NoError(t, err)
	assert.Equal(t, 0, current) // zero value for int
}

func TestEmpty_Dispose(t *testing.T) {
	// Arrange
	enumerator := enumerators.Empty[string]()

	// Act
	enumerator.Dispose() // Should not panic

	// Assert
	assert.False(t, enumerator.MoveNext())
}

func TestEmpty_ToSlice(t *testing.T) {
	// Arrange
	enumerator := enumerators.Empty[float64]()

	// Act
	result, err := enumerators.ToSlice(enumerator)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestConsume_WithSliceEnumerator(t *testing.T) {
	// Arrange
	enumerator := enumerators.Slice([]int{1, 2, 3, 4, 5})

	// Act
	err := enumerators.Consume(enumerator)

	// Assert
	assert.NoError(t, err)
	// Enumerator should be disposed automatically
}

func TestConsume_EmptyEnumerator(t *testing.T) {
	// Arrange
	enumerator := enumerators.Empty[int]()

	// Act
	err := enumerators.Consume(enumerator)

	// Assert
	assert.NoError(t, err)
}

func TestCleanup_CallsCleanupFunction(t *testing.T) {
	// Arrange
	base := enumerators.Slice([]int{1, 2, 3})
	cleanupCalled := false
	cleanup := func() {
		cleanupCalled = true
	}

	// Act
	wrapper := enumerators.Cleanup(base, cleanup)
	wrapper.Dispose()

	// Assert
	assert.True(t, cleanupCalled)
}

func TestCleanup_CallsCleanupOnlyOnce(t *testing.T) {
	// Arrange
	base := enumerators.Slice([]int{1, 2, 3})
	cleanupCallCount := 0
	cleanup := func() {
		cleanupCallCount++
	}

	// Act
	wrapper := enumerators.Cleanup(base, cleanup)
	wrapper.Dispose()
	wrapper.Dispose() // Second call should not call cleanup again

	// Assert
	assert.Equal(t, 1, cleanupCallCount)
}

func TestCleanup_ForwardsEnumeratorMethods(t *testing.T) {
	// Arrange
	base := enumerators.Slice([]int{10, 20, 30})
	cleanup := func() {} // no-op
	wrapper := enumerators.Cleanup(base, cleanup)

	// Act & Assert
	assert.True(t, wrapper.MoveNext())
	current, err := wrapper.Current()
	require.NoError(t, err)
	assert.Equal(t, 10, current)

	assert.True(t, wrapper.MoveNext())
	current, err = wrapper.Current()
	require.NoError(t, err)
	assert.Equal(t, 20, current)

	assert.True(t, wrapper.MoveNext())
	current, err = wrapper.Current()
	require.NoError(t, err)
	assert.Equal(t, 30, current)

	assert.False(t, wrapper.MoveNext())
	assert.NoError(t, wrapper.Err())
}

func TestCleanup_WithNilCleanupFunction(t *testing.T) {
	// Arrange
	base := enumerators.Slice([]int{1})
	wrapper := enumerators.Cleanup(base, nil)

	// Act
	wrapper.Dispose() // Should not panic

	// Assert - should still work
	assert.True(t, wrapper.MoveNext())
}