package enumerators_test

import (
	"testing"

	"github.com/fgrzl/enumerators"
	"github.com/stretchr/testify/assert"
)

func TestShouldCalculateTotalWhenSummingIntegers(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{1, 2, 3, 4, 5})

	// Act
	result, err := enumerators.Sum(input, func(x int) (int, error) { return x, nil })

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 15, result)
}

func TestShouldReturnZeroWhenSummingEmptyEnumerator(t *testing.T) {
	// Arrange
	input := enumerators.Empty[int]()

	// Act
	result, err := enumerators.Sum(input, func(x int) (int, error) { return x, nil })

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 0, result) // zero value for int
}

func TestShouldApplyTransformWhenSummingSingleElement(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{42})

	// Act
	result, err := enumerators.Sum(input, func(x int) (int, error) { return x * 2, nil })

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 84, result)
}

func TestShouldSumFloatsWhenFloatSliceProvided(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]float64{1.1, 2.2, 3.3})

	// Act
	result, err := enumerators.Sum(input, func(x float64) (float64, error) { return x, nil })

	// Assert
	assert.NoError(t, err)
	assert.InDelta(t, 6.6, result, 0.0001)
}

func TestShouldSumStringLengthsWhenStringSliceProvided(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]string{"hello", "world", "test"})

	// Act
	result, err := enumerators.Sum(input, func(s string) (int, error) { return len(s), nil })

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 14, result) // 5 + 5 + 4
}

func TestShouldApplyTransformationWhenSummingSquares(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{1, 2, 3, 4})

	// Act - sum of squares
	result, err := enumerators.Sum(input, func(x int) (int, error) { return x * x, nil })

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 30, result) // 1 + 4 + 9 + 16
}

func TestShouldReturnErrorWhenSelectorFails(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{1, 2, 3})
	selectorError := assert.AnError

	// Act
	result, err := enumerators.Sum(input, func(x int) (int, error) {
		if x == 2 {
			return 0, selectorError
		}
		return x, nil
	})

	// Assert
	assert.Error(t, err)
	assert.Equal(t, selectorError, err)
	assert.Equal(t, 0, result) // zero value when error occurs
}

func TestShouldSumNegativeNumbersWhenNegativeValuesProvided(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{-1, -2, -3, 4, 5})

	// Act
	result, err := enumerators.Sum(input, func(x int) (int, error) { return x, nil })

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 3, result) // -1 + -2 + -3 + 4 + 5 = 3
}

func TestShouldDisposeEnumeratorWhenSumCompletes(t *testing.T) {
	// Arrange
	disposed := false
	input := enumerators.Cleanup(
		enumerators.Slice([]int{1, 2, 3}),
		func() { disposed = true },
	)

	// Act
	_, err := enumerators.Sum(input, func(x int) (int, error) { return x, nil })

	// Assert
	assert.NoError(t, err)
	assert.True(t, disposed, "Enumerator should be disposed after sum")
}

func BenchmarkSum(b *testing.B) {
	input := make([]int, 1000)
	for i := range input {
		input[i] = i + 1
	}
	enumerator := enumerators.Slice(input)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result, err := enumerators.Sum(enumerator, func(x int) (int, error) { return x, nil })
		if err != nil || result != 500500 {
			b.Fatal("unexpected result")
		}
	}
}
