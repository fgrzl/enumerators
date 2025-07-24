package enumerators_test

import (
	"errors"
	"strconv"
	"testing"

	"github.com/fgrzl/enumerators"
	"github.com/stretchr/testify/assert"
)

func TestFilterMap_BasicTransformation(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]string{"1", "2", "not-a-number", "4", "5"})
	parseToInt := func(s string) (int, bool, error) {
		num, err := strconv.Atoi(s)
		if err != nil {
			return 0, false, nil // Skip invalid numbers
		}
		return num, true, nil
	}

	// Act
	filtered := enumerators.FilterMap(input, parseToInt)
	result, err := enumerators.ToSlice(filtered)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 4, 5}, result)
}

func TestFilterMap_EmptyInput(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]string{})
	transform := func(s string) (int, bool, error) {
		return len(s), true, nil
	}

	// Act
	filtered := enumerators.FilterMap(input, transform)
	result, err := enumerators.ToSlice(filtered)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestFilterMap_AllFiltered(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]string{"abc", "def", "ghi"})
	rejectAll := func(s string) (int, bool, error) {
		return 0, false, nil // Reject all
	}

	// Act
	filtered := enumerators.FilterMap(input, rejectAll)
	result, err := enumerators.ToSlice(filtered)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestFilterMap_AllAccepted(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]string{"hello", "world", "test"})
	toLength := func(s string) (int, bool, error) {
		return len(s), true, nil
	}

	// Act
	filtered := enumerators.FilterMap(input, toLength)
	result, err := enumerators.ToSlice(filtered)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []int{5, 5, 4}, result)
}

func TestFilterMap_TransformationError(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]string{"ok", "error", "ok"})
	expectedError := errors.New("transformation error")
	transformWithError := func(s string) (int, bool, error) {
		if s == "error" {
			return 0, false, expectedError
		}
		return len(s), true, nil
	}

	// Act
	filtered := enumerators.FilterMap(input, transformWithError)
	result, err := enumerators.ToSlice(filtered)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, []int{2}, result) // Should get "ok" before error
}

func TestFilterMap_ConditionalTransformation(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]int{1, 2, 3, 4, 5, 6})
	evenDoubler := func(x int) (int, bool, error) {
		if x%2 == 0 {
			return x * 2, true, nil // Double even numbers
		}
		return 0, false, nil // Skip odd numbers
	}

	// Act
	filtered := enumerators.FilterMap(input, evenDoubler)
	result, err := enumerators.ToSlice(filtered)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []int{4, 8, 12}, result) // 2*2, 4*2, 6*2
}

func TestFilterMap_StepByStep(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]string{"1", "abc", "3"})
	parseToInt := func(s string) (int, bool, error) {
		num, err := strconv.Atoi(s)
		if err != nil {
			return 0, false, nil
		}
		return num, true, nil
	}
	filtered := enumerators.FilterMap(input, parseToInt)

	// Act & Assert
	assert.True(t, filtered.MoveNext())
	current, err := filtered.Current()
	assert.NoError(t, err)
	assert.Equal(t, 1, current)

	assert.True(t, filtered.MoveNext()) // Should skip "abc"
	current, err = filtered.Current()
	assert.NoError(t, err)
	assert.Equal(t, 3, current)

	assert.False(t, filtered.MoveNext())
	assert.NoError(t, filtered.Err())
}

func TestFilterMap_SingleElement(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]string{"42"})
	parseToInt := func(s string) (int, bool, error) {
		num, err := strconv.Atoi(s)
		return num, err == nil, err
	}

	// Act
	filtered := enumerators.FilterMap(input, parseToInt)
	result, err := enumerators.ToSlice(filtered)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []int{42}, result)
}

func TestFilterMap_Dispose(t *testing.T) {
	// Arrange
	input := enumerators.Slice([]string{"1", "2", "3"})
	transform := func(s string) (int, bool, error) {
		num, _ := strconv.Atoi(s)
		return num, true, nil
	}
	filtered := enumerators.FilterMap(input, transform)

	// Act
	filtered.Dispose() // Should not panic

	// Assert - should still work
	assert.True(t, filtered.MoveNext())
	current, err := filtered.Current()
	assert.NoError(t, err)
	assert.Equal(t, 1, current)
}