package enumerators

import (
	"errors"
)

// NewPeekableEnumerator constructs a PeekableEnumerator
func Peekable[T any](inner Enumerator[T]) *PeekableEnumerator[T] {
	return &PeekableEnumerator[T]{
		inner: inner,
	}
}

// PeekableEnumerator wraps an Enumerator and adds peek functionality
type PeekableEnumerator[T any] struct {
	inner     Enumerator[T]
	hasPeeked bool
	peekValue T
	peekErr   error
}

// Dispose releases resources
func (p *PeekableEnumerator[T]) Dispose() {
	p.inner.Dispose()
}

// MoveNext advances to the next element
func (p *PeekableEnumerator[T]) MoveNext() bool {
	if p.hasPeeked {
		p.hasPeeked = false
		return true
	}
	return p.inner.MoveNext()
}

// Current returns the current element
func (p *PeekableEnumerator[T]) Current() (T, error) {
	if p.hasPeeked {
		return p.peekValue, p.peekErr
	}
	return p.inner.Current()
}

// Err returns the last error from the underlying enumerator
func (p *PeekableEnumerator[T]) Err() error {
	return p.inner.Err()
}

// Peek allows looking at the next element without advancing the enumerator
func (p *PeekableEnumerator[T]) Peek() (T, error) {
	if !p.hasPeeked {
		if !p.inner.MoveNext() {
			var zero T
			return zero, errors.New("no more elements to peek")
		}
		p.peekValue, p.peekErr = p.inner.Current()
		p.hasPeeked = true
	}
	return p.peekValue, p.peekErr
}

// HasNext checks if there is a next element available without advancing the enumerator
func (p *PeekableEnumerator[T]) HasNext() bool {
	_, err := p.Peek()
	return err == nil
}
