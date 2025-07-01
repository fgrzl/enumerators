package enumerators

type mapEnumerator[T any, U any] struct {
	base    Enumerator[T]
	mapper  func(T) (U, error)
	current U
	err     error
}

func (e *mapEnumerator[T, U]) MoveNext() bool {
	if !e.base.MoveNext() {
		return false
	}

	item, err := e.base.Current()
	if err != nil {
		e.err = err
		return false
	}

	u, err := e.mapper(item)
	if err != nil {
		e.err = err
		return false
	}
	e.current = u
	return true
}

func (e *mapEnumerator[T, U]) Current() (U, error) {
	return e.current, e.err
}

func (e *mapEnumerator[T, U]) Err() error {
	return e.err
}

func (e *mapEnumerator[T, U]) Dispose() {
	e.base.Dispose()
}

// Map creates a mapped enumerator
func Map[T any, U any](enumerator Enumerator[T], mapper func(T) (U, error)) Enumerator[U] {
	return &mapEnumerator[T, U]{
		base:   enumerator,
		mapper: mapper,
	}
}

// ToMap iterates over the provided enumerator and builds a map by applying
// the keyFn and valFn to each item. The resulting map uses the key from keyFn(item)
// and the value from valFn(item). If the enumerator yields an error during iteration,
// the function returns immediately with that error.
//
// The enumerator is disposed automatically at the end of processing.
//
// Returns an error if iteration fails. If the enumerator is nil, returns nil, nil.
func ToMap[T any, TKey comparable, TValue any](enumerator Enumerator[T], keyFn func(T) TKey, valFn func(T) TValue) (map[TKey]TValue, error) {
	if enumerator == nil {
		return nil, nil
	}
	defer enumerator.Dispose()

	result := make(map[TKey]TValue)
	for enumerator.MoveNext() {
		item, err := enumerator.Current()
		if err != nil {
			return nil, err
		}
		result[keyFn(item)] = valFn(item)
	}
	return result, nil
}
