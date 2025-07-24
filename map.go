package enumerators

// mapEnumerator applies a transformation function to each element from the base enumerator.
type mapEnumerator[T any, U any] struct {
	base    Enumerator[T]
	mapper  func(T) (U, error)
	current U
	err     error
}

// MoveNext advances to the next element and applies the transformation.
// Returns true if more elements are available, false otherwise.
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

// Current returns the transformed current element and any error encountered.
func (e *mapEnumerator[T, U]) Current() (U, error) {
	return e.current, e.err
}

// Err returns any error encountered during enumeration or transformation.
func (e *mapEnumerator[T, U]) Err() error {
	return e.err
}

// Dispose cleans up resources by disposing the underlying enumerator.
func (e *mapEnumerator[T, U]) Dispose() {
	e.base.Dispose()
}

// Map creates an enumerator that applies a transformation function to each element.
// The mapper function receives an element of type T and returns a transformed element of type U.
// If the mapper function returns an error, enumeration stops and the error is propagated.
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
		current, err := enumerator.Current()
		if err != nil {
			return nil, err
		}
		result[keyFn(current)] = valFn(current)
	}
	return result, nil
}

// ToMapWithKey iterates over the provided enumerator and builds a map using each item as the value.
// The key for each entry is computed by applying keyFn(item).
//
// The enumerator is disposed automatically at the end of processing.
//
// Returns an error if iteration fails. If the enumerator is nil, returns nil, nil.
func ToMapWithKey[TKey comparable, TValue any](enumerator Enumerator[TValue], keyFn func(TValue) TKey) (map[TKey]TValue, error) {
	return ToMap(enumerator, keyFn, func(v TValue) TValue { return v })
}
