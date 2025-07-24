package enumerators_test

import (
	"testing"

	"github.com/fgrzl/enumerators"
	"github.com/stretchr/testify/assert"
)

func TestForEach_BasicIteration(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{1, 2, 3, 4, 5})
	var collected []int

	// Act
	err := enumerators.ForEach(input, func(x int) error {
		collected = append(collected, x)
		return nil
	})

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, collected)
}

func TestForEach_EmptyEnumerator(t *testing.T) {
	// Arrange
	input := enumerators.Empty[int]()
	callCount := 0

	// Act
	err := enumerators.ForEach(input, func(x int) error {
		callCount++
		return nil
	})

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 0, callCount)
}

func TestForEach_SingleElement(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]string{"hello"})
	var result string

	// Act
	err := enumerators.ForEach(input, func(s string) error {
		result = s + "-processed"
		return nil
	})

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "hello-processed", result)
}

func TestForEach_ActionError(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{1, 2, 3, 4, 5})
	var collected []int
	expectedError := assert.AnError

	// Act
	err := enumerators.ForEach(input, func(x int) error {
		collected = append(collected, x)
		if x == 3 {
			return expectedError
		}
		return nil
	})

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, []int{1, 2, 3}, collected) // Should process up to the error
}

func TestForEach_WithSideEffects(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{1, 2, 3})
	sum := 0
	count := 0

	// Act
	err := enumerators.ForEach(input, func(x int) error {
		sum += x
		count++
		return nil
	})

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 6, sum)   // 1 + 2 + 3
	assert.Equal(t, 3, count) // 3 elements
}

func TestForEach_DisposesEnumerator(t *testing.T) {
	// Arrange
	disposed := false
	input := enumerators.Cleanup(
		enumerators.Slice([]int{1, 2, 3}),
		func() { disposed = true },
	)

	// Act
	err := enumerators.ForEach(input, func(x int) error { return nil })

	// Assert
	assert.NoError(t, err)
	assert.True(t, disposed, "Enumerator should be disposed after ForEach")
}

func TestForEach_DisposesOnError(t *testing.T) {
	// Arrange
	disposed := false
	input := enumerators.Cleanup(
		enumerators.Slice([]int{1, 2, 3}),
		func() { disposed = true },
	)

	// Act
	err := enumerators.ForEach(input, func(x int) error {
		if x == 2 {
			return assert.AnError
		}
		return nil
	})

	// Assert
	assert.Error(t, err)
	assert.True(t, disposed, "Enumerator should be disposed even when error occurs")
}

func TestForEach_ProcessingOrder(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]string{"first", "second", "third"})
	var order []string

	// Act
	err := enumerators.ForEach(input, func(s string) error {
		order = append(order, s)
		return nil
	})

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []string{"first", "second", "third"}, order)
}