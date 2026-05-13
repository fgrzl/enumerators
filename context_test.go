package enumerators_test

import (
	"context"
	"testing"
	"time"

	"github.com/fgrzl/enumerators"
	"github.com/stretchr/testify/assert"
)

func TestShouldMapWithContext(t *testing.T) {
	// Arrange
	ctx := context.Background()
	input := enumerators.Slice([]int{1, 2, 3})

	// Act
	result, err := enumerators.ToSlice(enumerators.MapWithContext(ctx, input, func(ctx context.Context, x int) (int, error) {
		return x * 2, nil
	}))

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []int{2, 4, 6}, result)
}

func TestShouldCancelMapWithContext(t *testing.T) {
	// Arrange
	ctx, cancel := context.WithCancel(context.Background())
	input := enumerators.Slice([]int{1, 2, 3, 4, 5})

	// Cancel after first element
	go func() {
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()

	// Act
	result, err := enumerators.ToSlice(enumerators.MapWithContext(ctx, input, func(ctx context.Context, x int) (int, error) {
		time.Sleep(20 * time.Millisecond) // Slow operation
		return x * 2, nil
	}))

	// Assert
	assert.Error(t, err)
	assert.Equal(t, context.Canceled, err)
	assert.GreaterOrEqual(t, len(result), 1)
}

func TestShouldFilterWithContext(t *testing.T) {
	// Arrange
	ctx := context.Background()
	input := enumerators.Slice([]int{1, 2, 3, 4, 5})

	// Act
	result, err := enumerators.ToSlice(enumerators.FilterWithContext(ctx, input, func(ctx context.Context, x int) bool {
		return x%2 == 0
	}))

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []int{2, 4}, result)
}
