package enumerators

import "errors"

type KeyedItem[K comparable, V any] struct {
	Key  K
	Item V
}

func ChunkWhen[V any, K comparable](
	in Enumerator[V],
	seed K,
	split func(prev K, item V) (K, bool, error),
) Enumerator[Enumerator[KeyedItem[K, V]]] {
	if in == nil {
		return &chunkWhenEnumerator[V, K]{exhausted: true}
	}
	return &chunkWhenEnumerator[V, K]{
		base:  in,
		prev:  seed,
		split: split,
	}
}

type chunkWhenEnumerator[V any, K comparable] struct {
	base         Enumerator[V]
	currentChunk *innerChunkWhenEnumerator[V, K]
	split        func(prev K, item V) (K, bool, error)
	prev         K
	pending      *V
	err          error
	exhausted    bool
}

func (e *chunkWhenEnumerator[V, K]) Dispose() {
	if e.currentChunk != nil {
		e.currentChunk.Dispose()
	}
	if e.base != nil {
		e.base.Dispose()
	}
}

func (e *chunkWhenEnumerator[V, K]) MoveNext() bool {
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

	var first V
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

	e.currentChunk = &innerChunkWhenEnumerator[V, K]{
		base:       e.base,
		first:      first,
		split:      e.split,
		state:      e.prev,
		setPending: func(v V) { e.pending = &v },
		setState:   func(v K) { e.prev = v },
		key:        e.prev,
	}
	return true
}

func (e *chunkWhenEnumerator[V, K]) Current() (Enumerator[KeyedItem[K, V]], error) {
	if e.currentChunk == nil {
		return nil, errors.New("no current chunk")
	}
	return &keyedChunkEnumerator[V, K]{inner: e.currentChunk}, nil
}

func (e *chunkWhenEnumerator[V, K]) Err() error {
	return e.err
}

type keyedChunkEnumerator[V any, K comparable] struct {
	inner *innerChunkWhenEnumerator[V, K]
}

func (e *keyedChunkEnumerator[V, K]) MoveNext() bool {
	return e.inner.MoveNext()
}

func (e *keyedChunkEnumerator[V, K]) Current() (KeyedItem[K, V], error) {
	v, err := e.inner.Current()
	return KeyedItem[K, V]{Key: e.inner.key, Item: v}, err
}

func (e *keyedChunkEnumerator[V, K]) Err() error {
	return e.inner.Err()
}

func (e *keyedChunkEnumerator[V, K]) Dispose() {
	e.inner.Dispose()
}

type innerChunkWhenEnumerator[V any, K comparable] struct {
	base       Enumerator[V]
	split      func(K, V) (K, bool, error)
	first      V
	state      K
	setState   func(K)
	setPending func(V)
	key        K

	current   V
	started   bool
	exhausted bool
	err       error
}

func (e *innerChunkWhenEnumerator[V, K]) MoveNext() bool {
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

	newState, split, err := e.split(e.state, item)
	if err != nil {
		e.exhausted = true
		e.err = err
		return false
	}

	if split {
		e.exhausted = true
		e.setState(newState)
		e.setPending(item)
		return false
	}

	e.state = newState
	e.current = item
	return true
}

func (e *innerChunkWhenEnumerator[V, K]) Current() (V, error) {
	return e.current, e.err
}

func (e *innerChunkWhenEnumerator[V, K]) Err() error {
	return e.err
}

func (e *innerChunkWhenEnumerator[V, K]) Dispose() {}
