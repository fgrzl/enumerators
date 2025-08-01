package enumerators

// skipIfEnumerator skips elements from the base enumerator that satisfy a condition.
type skipIfEnumerator[T any] struct {
	base      Enumerator[T]
	condition func(T) bool
	current   T
	err       error
}

// MoveNext advances to the next element that doesn't satisfy the skip condition.
func (e *skipIfEnumerator[T]) MoveNext() bool {
	for {
		if !e.base.MoveNext() {
			return false
		}

		item, err := e.base.Current()
		if err != nil {
			e.err = err
			return false
		}

		if e.condition(item) {
			continue
		}

		e.current = item
		return true
	}
}

// Current returns the current element that was not skipped.
func (e *skipIfEnumerator[T]) Current() (T, error) {
	return e.current, e.err
}

// Err returns any error encountered during enumeration.
func (e *skipIfEnumerator[T]) Err() error {
	return e.err
}

// Dispose cleans up resources by disposing the underlying enumerator.
func (e *skipIfEnumerator[T]) Dispose() {
	e.base.Dispose()
}

// SkipIf creates an enumerator that skips elements satisfying the specified condition.
// Elements are skipped when condition(element) returns true.
func SkipIf[T any](enumerator Enumerator[T], condition func(T) bool) Enumerator[T] {
	return &skipIfEnumerator[T]{
		base:      enumerator,
		condition: condition,
	}
}
