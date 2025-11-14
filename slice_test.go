package enumerators_test

import (
	"testing"

	"github.com/fgrzl/enumerators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShouldIterateAllElementsWhenBasicSliceProvided(t *testing.T) {
	// Arrange
	input := []int{1, 2, 3, 4, 5}
	enumerator := enumerators.Slice(input)

	// Act
	var results []int
	for enumerator.MoveNext() {
		current, err := enumerator.Current()
		require.NoError(t, err)
		results = append(results, current)
	}

	// Assert
	assert.Equal(t, input, results)
	assert.NoError(t, enumerator.Err())
}

func TestShouldReturnFalseWhenSliceIsEmpty(t *testing.T) {
	// Arrange
	enumerator := enumerators.Slice([]int{})

	// Act
	hasNext := enumerator.MoveNext()

	// Assert
	assert.False(t, hasNext)
	assert.NoError(t, enumerator.Err())
}

func TestShouldReturnSingleElementWhenSliceHasOneItem(t *testing.T) {
	// Arrange
	enumerator := enumerators.Slice([]string{"hello"})

	// Act
	hasFirst := enumerator.MoveNext()
	current, err := enumerator.Current()
	hasSecond := enumerator.MoveNext()

	// Assert
	assert.True(t, hasFirst)
	require.NoError(t, err)
	assert.Equal(t, "hello", current)
	assert.False(t, hasSecond)
	assert.NoError(t, enumerator.Err())
}

func TestShouldStillFunctionWhenDisposeCalled(t *testing.T) {
	// Arrange
	enumerator := enumerators.Slice([]int{1, 2, 3})

	// Act
	enumerator.Dispose() // Should not panic
	hasNext := enumerator.MoveNext()
	current, err := enumerator.Current()

	// Assert - should still work after dispose
	assert.True(t, hasNext)
	require.NoError(t, err)
	assert.Equal(t, 1, current)
}

func TestShouldReturnZeroValueWhenCurrentCalledBeforeMoveNext(t *testing.T) {
	// Arrange
	enumerator := enumerators.Slice([]int{1, 2, 3})

	// Act - calling Current before MoveNext should return zero value
	current, err := enumerator.Current()

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 0, current) // zero value for int
}

func TestShouldReturnOriginalSliceWhenConvertingSliceEnumeratorToSlice(t *testing.T) {
	// Arrange
	original := []int{1, 2, 3, 4, 5}
	enumerator := enumerators.Slice(original)

	// Act
	result, err := enumerators.ToSlice(enumerator)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, original, result)
}

func TestShouldReturnSliceOfResultsWhenConvertingRangeEnumeratorToSlice(t *testing.T) {
	// Arrange - use a non-slice enumerator
	enumerator := enumerators.Range(0, 3, func(i int) int { return i * 2 })

	// Act
	result, err := enumerators.ToSlice(enumerator)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []int{0, 2, 4}, result)
}

func TestShouldReturnEmptySliceWhenConvertingEmptyEnumeratorToSlice(t *testing.T) {
	// Arrange
	enumerator := enumerators.Slice([]int{})

	// Act
	result, err := enumerators.ToSlice(enumerator)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func BenchmarkSliceIteration(b *testing.B) {
	input := make([]int, 10000)
	for i := range input {
		input[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		enumerator := enumerators.Slice(input)
		for enumerator.MoveNext() {
			_, _ = enumerator.Current()
		}
		enumerator.Dispose()
	}
}

func BenchmarkSliceToSlice(b *testing.B) {
	input := make([]int, 10000)
	for i := range input {
		input[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		enumerator := enumerators.Slice(input)
		_, _ = enumerators.ToSlice(enumerator)
	}
}

func BenchmarkDirectSliceIteration(b *testing.B) {
	input := make([]int, 10000)
	for i := range input {
		input[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sum := 0
		for _, v := range input {
			sum += v
		}
		_ = sum
	}
}

func TestShouldConsumeRemainingElementsWhenPartiallyConsumed(t *testing.T) {
	// Arrange
	enumerator := enumerators.Slice([]int{1, 2, 3, 4, 5})

	// Act - consume first two
	assert.True(t, enumerator.MoveNext())
	first, _ := enumerator.Current()
	assert.Equal(t, 1, first)
	assert.True(t, enumerator.MoveNext())
	second, _ := enumerator.Current()
	assert.Equal(t, 2, second)

	// Now consume remaining
	result, err := enumerators.ToSlice(enumerator)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []int{3, 4, 5}, result)
}

func FuzzToSlice(f *testing.F) {
	f.Add([]byte{1, 2, 3})
	f.Fuzz(func(t *testing.T, inputBytes []byte) {
		input := make([]int, len(inputBytes))
		for i, b := range inputBytes {
			input[i] = int(b)
		}
		enumerator := enumerators.Slice(input)
		result, err := enumerators.ToSlice(enumerator)
		assert.NoError(t, err)
		assert.Equal(t, input, result)
	})
}
