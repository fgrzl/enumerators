package enumerators

// Zip combines two enumerators into one, yielding pairs of elements.
// Stops when either enumerator is exhausted.
func Zip[T, U any](first Enumerator[T], second Enumerator[U]) Enumerator[Pair[T, U]] {
	return &zipEnumerator[T, U]{
		first:  first,
		second: second,
	}
}

type zipEnumerator[T, U any] struct {
	first   Enumerator[T]
	second  Enumerator[U]
	current Pair[T, U]
	err     error
}

func (z *zipEnumerator[T, U]) MoveNext() bool {
	if !z.first.MoveNext() {
		z.err = z.first.Err()
		return false
	}
	if !z.second.MoveNext() {
		z.err = z.second.Err()
		return false
	}
	firstVal, err := z.first.Current()
	if err != nil {
		z.err = err
		return false
	}
	secondVal, err := z.second.Current()
	if err != nil {
		z.err = err
		return false
	}
	z.current = Pair[T, U]{First: firstVal, Second: secondVal}
	return true
}

func (z *zipEnumerator[T, U]) Current() (Pair[T, U], error) {
	return z.current, z.err
}

func (z *zipEnumerator[T, U]) Err() error {
	return z.err
}

func (z *zipEnumerator[T, U]) Dispose() {
	z.first.Dispose()
	z.second.Dispose()
}

// Pair represents a pair of values from Zip operation.
type Pair[T, U any] struct {
	First  T
	Second U
}
