package enumerators

// DefaultIfEmpty returns the elements of the enumerator, or a singleton sequence
// containing the default value if the enumerator is empty.
func DefaultIfEmpty[T any](enumerator Enumerator[T], defaultValue T) Enumerator[T] {
	return &defaultIfEmptyEnumerator[T]{
		base:         enumerator,
		defaultValue: defaultValue,
		state:        stateChecking,
	}
}

type defaultIfEmptyState int

const (
	stateChecking defaultIfEmptyState = iota
	stateReturningBase
	stateReturningDefault
	stateDone
)

type defaultIfEmptyEnumerator[T any] struct {
	base         Enumerator[T]
	defaultValue T
	state        defaultIfEmptyState
	current      T
	err          error
}

func (e *defaultIfEmptyEnumerator[T]) MoveNext() bool {
	switch e.state {
	case stateChecking:
		// First call: check if base has any elements
		if e.base.MoveNext() {
			// Base has elements, switch to returning base
			e.state = stateReturningBase
			current, err := e.base.Current()
			if err != nil {
				e.err = err
				return false
			}
			e.current = current
			return true
		}
		// Base is empty, check for error
		if err := e.base.Err(); err != nil {
			e.err = err
			return false
		}
		// Base is empty, return default
		e.state = stateReturningDefault
		e.current = e.defaultValue
		return true

	case stateReturningBase:
		// Continue returning base elements
		if e.base.MoveNext() {
			current, err := e.base.Current()
			if err != nil {
				e.err = err
				return false
			}
			e.current = current
			return true
		}
		e.err = e.base.Err()
		e.state = stateDone
		return false

	case stateReturningDefault:
		// Already returned default, done
		e.state = stateDone
		return false

	case stateDone:
		return false
	}

	return false
}

func (e *defaultIfEmptyEnumerator[T]) Current() (T, error) {
	return e.current, e.err
}

func (e *defaultIfEmptyEnumerator[T]) Err() error {
	return e.err
}

func (e *defaultIfEmptyEnumerator[T]) Dispose() {
	e.base.Dispose()
}
