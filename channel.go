package enumerators

import (
	"context"
	"sync"
)

// ChannelEnumerator provides enumeration over channels with context support.
// It supports publishing values, error handling, and graceful termination.
type ChannelEnumerator[T any] struct {
	context      context.Context
	dataCh       chan T
	errCh        chan error
	completeCh   chan struct{}
	stopCh       chan struct{}
	doneCh       chan struct{}
	current      T
	err          error
	stopOnce     sync.Once
	completeOnce sync.Once
	disposeOnce  sync.Once
}

// MoveNext advances the enumerator to the next value in the range.
// Returns true if more values are available, false otherwise.
func (e *ChannelEnumerator[T]) MoveNext() bool {
	if e.err != nil {
		return false
	}

	for {
		select {
		case <-e.doneCh:
			return false
		case <-e.context.Done():
			e.err = e.context.Err()
			e.stopPublishing()
			return false
		default:
		}

		// Drain published items before observing completion so buffered values remain readable.
		select {
		case data := <-e.dataCh:
			e.current = data
			return true
		default:
		}

		select {
		case err := <-e.errCh:
			if err != nil {
				e.err = err
			}
			e.stopPublishing()
			return false
		default:
		}

		select {
		case <-e.doneCh:
			return false
		case <-e.context.Done():
			e.err = e.context.Err()
			e.stopPublishing()
			return false
		case data := <-e.dataCh:
			e.current = data
			return true
		case err := <-e.errCh:
			if err != nil {
				e.err = err
			}
			e.stopPublishing()
			return false
		case <-e.completeCh:
			select {
			case data := <-e.dataCh:
				e.current = data
				return true
			default:
			}

			select {
			case err := <-e.errCh:
				if err != nil {
					e.err = err
				}
				e.stopPublishing()
				return false
			default:
			}

			return false
		}
	}
}

// Current returns the current value and any error encountered.
func (e *ChannelEnumerator[T]) Current() (T, error) {
	return e.current, e.err
}

// Err returns any error encountered during enumeration.
func (e *ChannelEnumerator[T]) Err() error {
	return e.err
}

// Dispose cleans up resources and signals termination.
func (e *ChannelEnumerator[T]) Dispose() {
	e.Complete()
	e.disposeOnce.Do(func() {
		e.stopPublishing()
		close(e.doneCh)
	})
}

// Publish sends a value to the enumerator for consumption.
func (e *ChannelEnumerator[T]) Publish(msg T) bool {
	select {
	case <-e.context.Done():
		return false
	case <-e.doneCh:
		return false
	case <-e.stopCh:
		return false
	default:
	}

	select {
	case <-e.context.Done():
		return false
	case <-e.doneCh:
		return false
	case <-e.stopCh:
		return false
	case e.dataCh <- msg:
		return true
	}
}

// Error signals an error to the enumerator.
func (e *ChannelEnumerator[T]) Error(err error) {
	if err == nil {
		return
	}

	select {
	case <-e.context.Done():
		return
	case <-e.doneCh:
		return
	case <-e.stopCh:
		return
	default:
	}

	select {
	case <-e.context.Done():
	case <-e.doneCh:
	case <-e.stopCh:
	case e.errCh <- err:
		e.stopPublishing()
	}
}

// Complete signals that no more values will be published.
func (e *ChannelEnumerator[T]) Complete() {
	e.completeOnce.Do(func() {
		e.stopPublishing()
		close(e.completeCh)
	})
}

func (e *ChannelEnumerator[T]) stopPublishing() {
	e.stopOnce.Do(func() {
		close(e.stopCh)
	})
}

// Channel creates a new channel-based enumerator with the specified buffer size.
// The enumerator respects the provided context for cancellation.
func Channel[T any](ctx context.Context, size int) *ChannelEnumerator[T] {
	return &ChannelEnumerator[T]{
		context:    ctx,
		dataCh:     make(chan T, size),
		errCh:      make(chan error, 1),
		completeCh: make(chan struct{}),
		stopCh:     make(chan struct{}),
		doneCh:     make(chan struct{}),
	}
}
