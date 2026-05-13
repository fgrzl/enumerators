package enumerators_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/fgrzl/enumerators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChannel_BasicPublishConsume(t *testing.T) {
	// Arrange
	ctx := context.Background()
	ch := enumerators.Channel[int](ctx, 10)

	// Act - publish values
	assert.True(t, ch.Publish(1))
	assert.True(t, ch.Publish(2))
	assert.True(t, ch.Publish(3))
	ch.Complete()

	// Act - consume values
	result, err := enumerators.ToSlice(ch)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3}, result)
}

func TestChannel_EmptyChannel(t *testing.T) {
	// Arrange
	ctx := context.Background()
	ch := enumerators.Channel[int](ctx, 10)

	// Act - complete without publishing
	ch.Complete()
	result, err := enumerators.ToSlice(ch)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestChannel_PublishAfterComplete(t *testing.T) {
	// Arrange
	ctx := context.Background()
	ch := enumerators.Channel[int](ctx, 10)

	// Act
	ch.Complete()

	// Assert - publishing after complete should fail gracefully (no panic)
	assert.False(t, ch.Publish(1), "Publishing after complete should return false")
}

func TestChannel_PublishWithError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	ch := enumerators.Channel[int](ctx, 10)
	expectedError := errors.New("test error")

	// Act
	assert.True(t, ch.Publish(1))
	ch.Error(expectedError)

	// Consume - ToSlice will stop when error is encountered.
	result, err := enumerators.ToSlice(ch)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.True(t, len(result) == 0 || (len(result) == 1 && result[0] == 1))
}

func TestChannel_StepByStep(t *testing.T) {
	// Arrange
	ctx := context.Background()
	ch := enumerators.Channel[string](ctx, 5)

	// Publish first value
	assert.True(t, ch.Publish("hello"))

	// Act & Assert
	assert.True(t, ch.MoveNext())
	current, err := ch.Current()
	require.NoError(t, err)
	assert.Equal(t, "hello", current)

	// Publish second value
	assert.True(t, ch.Publish("world"))

	assert.True(t, ch.MoveNext())
	current, err = ch.Current()
	require.NoError(t, err)
	assert.Equal(t, "world", current)

	// Complete and verify no more values
	ch.Complete()
	assert.False(t, ch.MoveNext())
	assert.NoError(t, ch.Err())
}

func TestChannel_ContextCancellation(t *testing.T) {
	// Arrange
	ctx, cancel := context.WithCancel(context.Background())
	ch := enumerators.Channel[int](ctx, 0)

	// Act
	go func() {
		time.Sleep(10 * time.Millisecond)
		ch.Publish(1)
	}()

	assert.True(t, ch.MoveNext())
	value, err := ch.Current()
	assert.NoError(t, err)
	assert.Equal(t, 1, value)

	cancel()

	// Assert
	assert.False(t, ch.MoveNext())
	assert.False(t, ch.Publish(2))
	assert.ErrorIs(t, ch.Err(), context.Canceled)
}

func TestChannel_ContextTimeout(t *testing.T) {
	// Arrange
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	ch := enumerators.Channel[int](ctx, 1)

	// Act
	assert.True(t, ch.Publish(1))
	assert.True(t, ch.MoveNext())
	value, err := ch.Current()
	assert.NoError(t, err)
	assert.Equal(t, 1, value)

	time.Sleep(100 * time.Millisecond)

	select {
	case <-ctx.Done():
	default:
		t.Fatal("Context should be timed out")
	}

	// Assert
	assert.False(t, ch.MoveNext())
	assert.False(t, ch.Publish(2), "Publishing after context timeout should return false")
	assert.ErrorIs(t, ch.Err(), context.DeadlineExceeded)
}

func TestChannel_BufferedChannel(t *testing.T) {
	// Arrange
	ctx := context.Background()
	ch := enumerators.Channel[int](ctx, 2)

	// Act - fill buffer exactly
	assert.True(t, ch.Publish(1))
	assert.True(t, ch.Publish(2))
	ch.Complete()

	// Consume all
	result, err := enumerators.ToSlice(ch)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2}, result)
}

func TestChannel_Dispose(t *testing.T) {
	// Arrange
	ctx := context.Background()
	ch := enumerators.Channel[int](ctx, 10)

	// Publish some values
	assert.True(t, ch.Publish(1))
	assert.True(t, ch.Publish(2))

	// Act
	ch.Dispose()

	// Assert - publishing after dispose should fail gracefully
	assert.False(t, ch.Publish(3))
}

func TestChannel_DisposeAfterComplete(t *testing.T) {
	// Arrange
	ctx := context.Background()
	ch := enumerators.Channel[int](ctx, 10)

	// Act
	ch.Complete()
	ch.Dispose()

	// Assert - publishing after dispose should fail gracefully
	assert.False(t, ch.Publish(1))
}

func TestChannel_CurrentBeforeMoveNext(t *testing.T) {
	// Arrange
	ctx := context.Background()
	ch := enumerators.Channel[int](ctx, 10)

	// Act - calling Current before MoveNext should return zero value
	current, err := ch.Current()

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 0, current)
}

func TestChannel_ErrorHandling(t *testing.T) {
	// Arrange
	ctx := context.Background()
	ch := enumerators.Channel[int](ctx, 10)
	expectedError := errors.New("channel error")

	// Act
	ch.Error(expectedError)

	// Assert
	assert.False(t, ch.MoveNext())
	assert.Equal(t, expectedError, ch.Err())
}

func TestChannel_MultipleErrors(t *testing.T) {
	// Arrange
	ctx := context.Background()
	ch := enumerators.Channel[int](ctx, 10)
	firstError := errors.New("first error")
	secondError := errors.New("second error")

	// Act
	ch.Error(firstError)
	ch.Error(secondError)

	// Assert
	assert.False(t, ch.MoveNext())
	assert.Equal(t, firstError, ch.Err())
}

func TestChannel_ErrorAfterComplete(t *testing.T) {
	// Arrange
	ctx := context.Background()
	ch := enumerators.Channel[int](ctx, 1)

	// Act
	ch.Complete()
	assert.NotPanics(t, func() {
		ch.Error(errors.New("late error"))
	})

	// Assert
	assert.False(t, ch.MoveNext())
	assert.NoError(t, ch.Err())
}

func TestChannel_ErrorAfterDispose(t *testing.T) {
	// Arrange
	ctx := context.Background()
	ch := enumerators.Channel[int](ctx, 1)

	// Act
	ch.Dispose()
	assert.NotPanics(t, func() {
		ch.Error(errors.New("late error"))
	})

	// Assert
	assert.False(t, ch.MoveNext())
	assert.NoError(t, ch.Err())
}

func TestChannel_ConcurrentPublish(t *testing.T) {
	// Arrange
	ctx := context.Background()
	const publisherCount = 8
	const itemsPerPublisher = 25
	ch := enumerators.Channel[int](ctx, 8)
	expected := make([]int, 0, publisherCount*itemsPerPublisher)
	publishErrors := make(chan error, publisherCount)

	for publisher := 0; publisher < publisherCount; publisher++ {
		for item := 0; item < itemsPerPublisher; item++ {
			expected = append(expected, publisher*itemsPerPublisher+item)
		}
	}

	var wg sync.WaitGroup
	for publisher := 0; publisher < publisherCount; publisher++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for item := 0; item < itemsPerPublisher; item++ {
				value := publisher*itemsPerPublisher + item
				if !ch.Publish(value) {
					publishErrors <- errors.New("publish failed")
					return
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(publishErrors)
		ch.Complete()
	}()

	// Act
	result, err := enumerators.ToSlice(ch)

	// Assert
	require.NoError(t, err)
	for publishErr := range publishErrors {
		require.NoError(t, publishErr)
	}
	assert.Len(t, result, publisherCount*itemsPerPublisher)
	assert.ElementsMatch(t, expected, result)
}
