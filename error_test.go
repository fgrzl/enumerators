package enumerators_test

import (
	"errors"
	"testing"

	"github.com/fgrzl/enumerators"
	"github.com/stretchr/testify/assert"
)

func TestShouldReturnProvidedError_WhenCreatedWithError(t *testing.T) {
	// Arrange
	expectedError := errors.New("test error")
	
	// Act
	errorEnum := enumerators.Error[int](expectedError)
	hasMoved := errorEnum.MoveNext()
	current, currentErr := errorEnum.Current()

	// Assert
	assert.False(t, hasMoved)
	assert.Equal(t, expectedError, errorEnum.Err())
	assert.Equal(t, 0, current)      // zero value for int
	assert.Equal(t, expectedError, currentErr)
}

func TestShouldNotHaveError_WhenCreatedWithNilError(t *testing.T) {
	// Arrange & Act
	errorEnum := enumerators.Error[string](nil)
	hasMoved := errorEnum.MoveNext()
	current, currentErr := errorEnum.Current()

	// Assert
	assert.False(t, hasMoved)
	assert.NoError(t, errorEnum.Err())
	assert.Equal(t, "", current) // zero value for string
	assert.NoError(t, currentErr)
}

func TestShouldMaintainBehavior_WhenDisposeCalledOnError(t *testing.T) {
	// Arrange
	expectedError := errors.New("dispose test")
	errorEnum := enumerators.Error[float64](expectedError)

	// Act
	errorEnum.Dispose() // Should not panic
	hasMoved := errorEnum.MoveNext()

	// Assert - behavior should be unchanged
	assert.False(t, hasMoved)
	assert.Equal(t, expectedError, errorEnum.Err())
}

func TestShouldReturnError_WhenConvertingErrorEnumeratorToSlice(t *testing.T) {
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

func TestShouldNotExecuteAction_WhenUsingErrorEnumeratorWithForEach(t *testing.T) {
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

func TestShouldProceedToSecondEnumerator_WhenChainedWithNormalEnumerator(t *testing.T) {
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