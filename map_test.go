package enumerators

import (
	"errors"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShouldCreateMapFromStringLengths_WhenMappingStringsToLengths(t *testing.T) {
	// Arrange
	input := Slice([]string{"a", "bb", "ccc"})

	// Act
	result, err := ToMap(
		input,
		func(s string) int { return len(s) }, // key: length of string
		func(s string) string { return s + "-x" }, // value: original + "-x"
	)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 3)
	assert.Equal(t, "a-x", result[1])
	assert.Equal(t, "bb-x", result[2])
	assert.Equal(t, "ccc-x", result[3])
}

func TestShouldTransformAllElements_WhenMappingWithBasicFunction(t *testing.T) {
	// Arrange
	input := Slice([]int{1, 2, 3, 4, 5})
	doubler := func(x int) (int, error) { return x * 2, nil }

	// Act
	mapped := Map(input, doubler)
	result, err := ToSlice(mapped)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []int{2, 4, 6, 8, 10}, result)
}

func TestShouldReturnEmpty_WhenMappingEmptyInput(t *testing.T) {
	// Arrange
	input := Slice([]int{})
	transform := func(x int) (string, error) { return strconv.Itoa(x), nil }

	// Act
	mapped := Map(input, transform)
	result, err := ToSlice(mapped)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestShouldChangeTypes_WhenMappingIntegersToStrings(t *testing.T) {
	// Arrange
	input := Slice([]int{1, 2, 3})
	toString := func(x int) (string, error) { return "num-" + strconv.Itoa(x), nil }

	// Act
	mapped := Map(input, toString)
	result, err := ToSlice(mapped)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []string{"num-1", "num-2", "num-3"}, result)
}

func TestShouldStopAndReturnError_WhenTransformationFails(t *testing.T) {
	// Arrange
	input := Slice([]int{1, 2, 3, 4, 5})
	expectedError := errors.New("transformation error")
	errorOnThree := func(x int) (int, error) {
		if x == 3 {
			return 0, expectedError
		}
		return x * 10, nil
	}

	// Act
	mapped := Map(input, errorOnThree)
	result, err := ToSlice(mapped)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, []int{10, 20}, result) // Should get results before error
}

func TestShouldProcessElementsSequentially_WhenIteratingStepByStep(t *testing.T) {
	// Arrange
	input := Slice([]string{"hello", "world"})
	toUpper := func(s string) (string, error) { return strings.ToUpper(s), nil }
	mapped := Map(input, toUpper)

	// Act - First element
	hasFirst := mapped.MoveNext()
	firstCurrent, firstErr := mapped.Current()

	// Act - Second element  
	hasSecond := mapped.MoveNext()
	secondCurrent, secondErr := mapped.Current()

	// Act - No more elements
	hasThird := mapped.MoveNext()

	// Assert
	assert.True(t, hasFirst)
	require.NoError(t, firstErr)
	assert.Equal(t, "HELLO", firstCurrent)

	assert.True(t, hasSecond)
	require.NoError(t, secondErr)
	assert.Equal(t, "WORLD", secondCurrent)

	assert.False(t, hasThird)
	assert.NoError(t, mapped.Err())
}

func TestMap_Dispose(t *testing.T) {
	// Arrange
	input := Slice([]int{1, 2, 3})
	identity := func(x int) (int, error) { return x, nil }
	mapped := Map(input, identity)

	// Act
	mapped.Dispose() // Should not panic

	// Assert - should still work
	assert.True(t, mapped.MoveNext())
	current, err := mapped.Current()
	require.NoError(t, err)
	assert.Equal(t, 1, current)
}

func TestToMap_NilEnumerator(t *testing.T) {
	// Arrange & Act
	result, err := ToMap[int, int, string](nil, func(x int) int { return x }, func(x int) string { return strconv.Itoa(x) })

	// Assert
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestToMap_EmptyEnumerator(t *testing.T) {
	// Arrange
	input := Empty[int]()

	// Act
	result, err := ToMap(input, func(x int) int { return x }, func(x int) string { return strconv.Itoa(x) })

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestToMap_DuplicateKeys(t *testing.T) {
	// Arrange
	input := Slice([]string{"a", "b", "c", "aa"}) // "a" and "aa" both have key 1 when using first char

	// Act
	result, err := ToMap(
		input,
		func(s string) rune { return rune(s[0]) }, // key: first character
		func(s string) string { return s },
	)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 3) // 'a', 'b', 'c'
	assert.Contains(t, result, 'a')
	assert.Contains(t, result, 'b')
	assert.Contains(t, result, 'c')
	// Last value with same key should win
	assert.Equal(t, "aa", result['a'])
}

func TestToMapWithKey_BasicUsage(t *testing.T) {
	// Arrange - use strings with unique lengths to avoid conflicts
	input := Slice([]string{"cat", "banana", "hi"})

	// Act
	result, err := ToMapWithKey(input, func(s string) int { return len(s) })

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 3) // lengths 2, 3, and 6
	assert.Equal(t, "hi", result[2])
	assert.Equal(t, "cat", result[3])
	assert.Equal(t, "banana", result[6])
}

func TestToMapWithKey_NilEnumerator(t *testing.T) {
	// Arrange & Act
	result, err := ToMapWithKey[int, string](nil, func(s string) int { return len(s) })

	// Assert
	assert.NoError(t, err)
	assert.Nil(t, result)
}
