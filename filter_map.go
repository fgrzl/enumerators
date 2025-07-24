package enumerators

// filterMapper applies a transformation and filtering function to each element from the base enumerator.
type filterMapper[TIn any, TOut any] struct {
	base    Enumerator[TIn]
	apply   func(TIn) (TOut, bool, error)
	current TOut
	err     error
}

// MoveNext advances to the next element that is both transformed and accepted by the filter.
func (e *filterMapper[TIn, TOut]) MoveNext() bool {
	for {
		if !e.base.MoveNext() {
			return false
		}

		item, err := e.base.Current()
		if err != nil {
			e.err = err
			return false
		}

		u, ok, err := e.apply(item)

		if err != nil {
			e.err = err
			return false
		}

		if !ok {
			continue
		}

		e.err = err
		e.current = u
		return true
	}
}

// Current returns the transformed current element and any error encountered.
func (e *filterMapper[TIn, TOut]) Current() (TOut, error) {
	return e.current, e.err
}

// Err returns any error encountered during enumeration or transformation.
func (e *filterMapper[TIn, TOut]) Err() error {
	return e.err
}

// Dispose cleans up resources by disposing the underlying enumerator.
func (e *filterMapper[TIn, TOut]) Dispose() {
	e.base.Dispose()
}

// FilterMap creates an enumerator that applies both transformation and filtering in a single pass.
// The apply function receives an element and returns (transformedValue, shouldInclude, error).
// Only elements where shouldInclude is true are yielded after transformation.
func FilterMap[TIn any, TOut any](enumerator Enumerator[TIn], apply func(TIn) (TOut, bool, error)) Enumerator[TOut] {
	return &filterMapper[TIn, TOut]{
		base:  enumerator,
		apply: apply,
	}
}
