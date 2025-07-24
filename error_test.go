package enumerators_test

import (
	"errors"
	"testing"

	"github.com/fgrzl/enumerators"
	"github.com/stretchr/testify/assert"
)

func TestError_ReturnsError(t *testing.T) {
	// Arrange
	expectedError := errors.New("test error")
	
	// Act
	errorEnum := enumerators.Error[int](expectedError)

	// Assert
	assert.False(t, errorEnum.MoveNext())
	assert.Equal(t, expectedError, errorEnum.Err())
	
	current, err := errorEnum.Current()
	assert.Equal(t, 0, current)      // zero value for int
	assert.Equal(t, expectedError, err)
}

func TestError_NilError(t *testing.T) {
	// Arrange & Act
	errorEnum := enumerators.Error[string](nil)

	// Assert
	assert.False(t, errorEnum.MoveNext())
	assert.NoError(t, errorEnum.Err())
	
	current, err := errorEnum.Current()
	assert.Equal(t, "", current) // zero value for string
	assert.NoError(t, err)
}

func TestError_Dispose(t *testing.T) {
	// Arrange
	expectedError := errors.New("dispose test")
	errorEnum := enumerators.Error[float64](expectedError)

	// Act
	errorEnum.Dispose() // Should not panic

	// Assert - behavior should be unchanged
	assert.False(t, errorEnum.MoveNext())
	assert.Equal(t, expectedError, errorEnum.Err())
}

func TestError_ToSlice(t *testing.T) {
	// Arrange
	expectedError := errors.New("slice test")
	errorEnum := enumerators.Error[int](expectedError)

	// Act
	result, err := enumerators.ToSlice(errorEnum)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Empty(t, result)
}

func TestError_WithForEach(t *testing.T) {
	// Arrange
	expectedError := errors.New("foreach test")
	errorEnum := enumerators.Error[int](expectedError)
	callCount := 0

	// Act
	err := enumerators.ForEach(errorEnum, func(x int) error {
		callCount++
		return nil
	})

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, 0, callCount) // Should not call the action since MoveNext() returns false
}

func TestError_WithChain(t *testing.T) {
	// Arrange
	expectedError := errors.New("chain test")
	errorEnum := enumerators.Error[int](expectedError)
	normalEnum := enumerators.Slice([]int{1, 2, 3})

	// Act
	chained := enumerators.Chain(errorEnum, normalEnum)
	result, err := enumerators.ToSlice(chained)

	// Assert
	assert.NoError(t, err) // Chain should proceed to the second enumerator
	assert.Equal(t, []int{1, 2, 3}, result)
}