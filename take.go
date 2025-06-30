package enumerators

type takeEnumerator[T any] struct {
	base    Enumerator[T]
	size    int
	count   int
	current T
	err     error
}

func (e *takeEnumerator[T]) MoveNext() bool {
	if !e.base.MoveNext() {
		return false
	}

	item, err := e.base.Current()
	if err != nil {
		e.err = err
		return false
	}

	if e.count >= e.size {
		return false
	}

	e.current = item
	e.count++
	return true
}

func (e *takeEnumerator[T]) Current() (T, error) {
	return e.current, e.err
}

func (e *takeEnumerator[T]) Err() error {
	return e.err
}

func (e *takeEnumerator[T]) Dispose() {
	e.base.Dispose()
}

// Take takes the specified number of items
func Take[T any](enumerator Enumerator[T], take int) Enumerator[T] {
	return &takeEnumerator[T]{
		base: enumerator,
		size: take,
	}
}
