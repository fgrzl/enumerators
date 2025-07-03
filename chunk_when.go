package enumerators

import "errors"

func ChunkWhen[T any, TCompare comparable](
	in Enumerator[T],
	seed TCompare,
	split func(val TCompare, item T) (TCompare, bool, error),
) Enumerator[Enumerator[T]] {
	if in == nil {
		return &chunkWhenEnumerator[T, TCompare]{exhausted: true}
	}
	return &chunkWhenEnumerator[T, TCompare]{
		base: in,
		prev: seed,
		when: split,
	}
}

type chunkWhenEnumerator[T any, TVal any] struct {
	base         Enumerator[T]
	currentChunk *innerChunkDiffEnumerator[T, TVal]
	when         func(val TVal, item T) (TVal, bool, error)
	prev         TVal
	pending      *T
	err          error
	exhausted    bool
}

func (e *chunkWhenEnumerator[T, TVal]) Dispose() {
	if e.currentChunk != nil {
		e.currentChunk.Dispose()
	}
	if e.base != nil {
		e.base.Dispose()
	}
}

func (e *chunkWhenEnumerator[T, TVal]) MoveNext() bool {
	if e.exhausted {
		return false
	}

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

	e.currentChunk = &innerChunkDiffEnumerator[T, TVal]{
		base:       e.base,
		first:      first,
		when:       e.when,
		state:      e.prev,
		setPending: func(v T) { e.pending = &v },
		setState:   func(v TVal) { e.prev = v },
	}
	return true
}

func (e *chunkWhenEnumerator[T, TVal]) Current() (Enumerator[T], error) {
	if e.currentChunk == nil {
		return nil, errors.New("no current chunk")
	}
	return e.currentChunk, e.currentChunk.err
}

func (e *chunkWhenEnumerator[T, TVal]) Err() error {
	return e.err
}

// --- Inner chunk enumerator ---

type innerChunkDiffEnumerator[T any, TVal any] struct {
	base       Enumerator[T]
	when       func(TVal, T) (TVal, bool, error)
	first      T
	state      TVal
	setState   func(TVal)
	setPending func(T)
	current    T
	started    bool
	exhausted  bool
	err        error
}

func (e *innerChunkDiffEnumerator[T, TVal]) MoveNext() bool {
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

	newState, split, err := e.when(e.state, item)
	if err != nil {
		e.exhausted = true
		e.err = err
		return false
	}

	if split {
		e.exhausted = true
		if e.setState != nil {
			e.setState(newState)
		}
		if e.setPending != nil {
			e.setPending(item)
		}
		return false
	}

	e.state = newState
	e.current = item
	return true
}

func (e *innerChunkDiffEnumerator[T, TVal]) Current() (T, error) {
	return e.current, e.err
}

func (e *innerChunkDiffEnumerator[T, TVal]) Err() error {
	return e.err
}

func (e *innerChunkDiffEnumerator[T, TVal]) Dispose() {}
