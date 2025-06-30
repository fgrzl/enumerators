package enumerators

import "errors"

// ChunkWhen returns an enumerator of enumerators, where a new chunk starts
// whenever the provided `when` function returns true for an item.
// The triggering item starts the new chunk.
func ChunkWhen[T any](
	in Enumerator[T],
	when func(item T) (bool, error),
) Enumerator[Enumerator[T]] {
	if in == nil {
		return &chunkWhenEnumerator[T]{exhausted: true}
	}
	return &chunkWhenEnumerator[T]{
		base: in,
		when: when,
	}
}

// chunkWhenEnumerator yields inner enumerators, each representing a chunk
// of items from the base enumerator, split according to the `when` predicate.
type chunkWhenEnumerator[T any] struct {
	base         Enumerator[T]
	currentChunk *innerWhenChunkEnumerator[T]
	when         func(item T) (bool, error)
	pending      *T
	err          error
	exhausted    bool
}

// Dispose releases resources held by the chunkWhenEnumerator and its current chunk.
func (e *chunkWhenEnumerator[T]) Dispose() {
	if e.currentChunk != nil {
		e.currentChunk.Dispose()
	}
	if e.base != nil {
		e.base.Dispose()
	}
}

// MoveNext advances to the next chunk, draining the previous chunk fully if needed.
func (e *chunkWhenEnumerator[T]) MoveNext() bool {
	if e.exhausted {
		return false
	}

	// Drain the previous chunk before advancing
	if e.currentChunk != nil && !e.currentChunk.exhausted {
		for e.currentChunk.MoveNext() {
		}
		if e.currentChunk.err != nil {
			e.err = e.currentChunk.err
			e.exhausted = true
			return false
		}
	}

	var first T
	var err error

	// Use pending item from previous chunk's terminator
	if e.pending != nil {
		first = *e.pending
		e.pending = nil
	} else {
		if !e.base.MoveNext() {
			e.exhausted = true
			e.err = e.base.Err()
			return false
		}
		first, err = e.base.Current()
		if err != nil {
			e.err = err
			e.exhausted = true
			return false
		}
	}

	e.currentChunk = &innerWhenChunkEnumerator[T]{
		base:  e.base,
		first: first,
		when:  e.when,
		setPending: func(v T) {
			e.pending = &v
		},
	}
	return true
}

// Current returns the current chunk enumerator.
func (e *chunkWhenEnumerator[T]) Current() (Enumerator[T], error) {
	if e.currentChunk == nil {
		return nil, errors.New("no current chunk")
	}
	return e.currentChunk, e.currentChunk.err
}

// Err returns any error encountered during chunk enumeration.
func (e *chunkWhenEnumerator[T]) Err() error {
	return e.err
}

// innerWhenChunkEnumerator yields a single chunk of items until `when` returns true again.
type innerWhenChunkEnumerator[T any] struct {
	base       Enumerator[T]
	when       func(item T) (bool, error)
	first      T
	setPending func(T)
	current    T
	exhausted  bool
	started    bool
	err        error
}

// MoveNext advances to the next item in the chunk.
// The chunk ends when `when(item)` returns true, and that item will be used as the start of the next chunk.
func (e *innerWhenChunkEnumerator[T]) MoveNext() bool {
	if e.exhausted {
		return false
	}

	if !e.started {
		e.current = e.first
		e.started = true
		return true
	}

	if !e.base.MoveNext() {
		e.exhausted = true
		e.err = e.base.Err()
		return false
	}

	item, err := e.base.Current()
	if err != nil {
		e.exhausted = true
		e.err = err
		return false
	}

	startNew, err := e.when(item)
	if err != nil {
		e.exhausted = true
		e.err = err
		return false
	}

	if startNew {
		e.exhausted = true
		if e.setPending != nil {
			e.setPending(item)
		}
		return false
	}

	e.current = item
	return true
}

// Current returns the current item in the chunk.
func (e *innerWhenChunkEnumerator[T]) Current() (T, error) {
	return e.current, e.err
}

// Dispose is a no-op for innerWhenChunkEnumerator.
func (e *innerWhenChunkEnumerator[T]) Dispose() {}

// Err returns any error encountered during chunk iteration.
func (e *innerWhenChunkEnumerator[T]) Err() error {
	return e.err
}
