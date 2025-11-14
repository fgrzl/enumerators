package enumerators

// Partition splits the enumerator into two based on a predicate.
// Returns two enumerators: one for elements where predicate is true, one for false.
// The original enumerator is consumed eagerly.
func Partition[T any](enumerator Enumerator[T], predicate func(T) bool) (Enumerator[T], Enumerator[T]) {
	slice, err := ToSlice(enumerator)
	if err != nil {
		errorEnum := Generate(func() (T, bool, error) {
			return *new(T), false, err
		})
		return errorEnum, errorEnum
	}
	var trueSlice, falseSlice []T
	for _, item := range slice {
		if predicate(item) {
			trueSlice = append(trueSlice, item)
		} else {
			falseSlice = append(falseSlice, item)
		}
	}
	return Slice(trueSlice), Slice(falseSlice)
}
