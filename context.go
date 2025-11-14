package enumerators

import "context"

// MapWithContext applies a transformation function to each element with context support.
// The mapper function receives context for cancellation and the element to transform.
// If context is cancelled, enumeration stops and returns the cancellation error.
func MapWithContext[T any, U any](ctx context.Context, enumerator Enumerator[T], mapper func(context.Context, T) (U, error)) Enumerator[U] {
	return &mapContextEnumerator[T, U]{
		ctx:    ctx,
		base:   enumerator,
		mapper: mapper,
	}
}

type mapContextEnumerator[T any, U any] struct {
	ctx     context.Context
	base    Enumerator[T]
	mapper  func(context.Context, T) (U, error)
	current U
	err     error
}

func (e *mapContextEnumerator[T, U]) MoveNext() bool {
	// Check context cancellation first
	select {
	case <-e.ctx.Done():
		e.err = e.ctx.Err()
		return false
	default:
	}

	if !e.base.MoveNext() {
		e.err = e.base.Err()
		return false
	}

	item, err := e.base.Current()
	if err != nil {
		e.err = err
		return false
	}

	u, err := e.mapper(e.ctx, item)
	if err != nil {
		e.err = err
		return false
	}
	e.current = u
	return true
}

func (e *mapContextEnumerator[T, U]) Current() (U, error) {
	return e.current, e.err
}

func (e *mapContextEnumerator[T, U]) Err() error {
	return e.err
}

func (e *mapContextEnumerator[T, U]) Dispose() {
	e.base.Dispose()
}

// FilterWithContext filters elements based on a predicate with context support.
// The predicate function receives context for cancellation.
// If context is cancelled, enumeration stops and returns the cancellation error.
func FilterWithContext[T any](ctx context.Context, enumerator Enumerator[T], predicate func(context.Context, T) bool) Enumerator[T] {
	return &filterContextEnumerator[T]{
		ctx:       ctx,
		base:      enumerator,
		predicate: predicate,
	}
}

type filterContextEnumerator[T any] struct {
	ctx       context.Context
	base      Enumerator[T]
	predicate func(context.Context, T) bool
	current   T
	err       error
}

func (e *filterContextEnumerator[T]) MoveNext() bool {
	for {
		// Check context cancellation
		select {
		case <-e.ctx.Done():
			e.err = e.ctx.Err()
			return false
		default:
		}

		if !e.base.MoveNext() {
			e.err = e.base.Err()
			return false
		}

		item, err := e.base.Current()
		if err != nil {
			e.err = err
			return false
		}

		if e.predicate(e.ctx, item) {
			e.current = item
			return true
		}
	}
}

func (e *filterContextEnumerator[T]) Current() (T, error) {
	return e.current, e.err
}

func (e *filterContextEnumerator[T]) Err() error {
	return e.err
}

func (e *filterContextEnumerator[T]) Dispose() {
	e.base.Dispose()
}
