package enumerators

// Union returns an enumerator with unique elements from both enumerators.
// Preserves order from first enumerator, then adds from second.
func Union[T comparable](first, second Enumerator[T]) Enumerator[T] {
	firstSlice, err := ToSlice(first)
	if err != nil {
		return Generate(func() (T, bool, error) {
			return *new(T), false, err
		})
	}
	secondSlice, err := ToSlice(second)
	if err != nil {
		return Generate(func() (T, bool, error) {
			return *new(T), false, err
		})
	}
	seen := make(map[T]bool)
	var result []T
	for _, item := range firstSlice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	for _, item := range secondSlice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return Slice(result)
}

// Intersect returns an enumerator with elements present in both enumerators.
// Preserves order from first enumerator.
func Intersect[T comparable](first, second Enumerator[T]) Enumerator[T] {
	firstSlice, err := ToSlice(first)
	if err != nil {
		return Generate(func() (T, bool, error) {
			return *new(T), false, err
		})
	}
	secondSlice, err := ToSlice(second)
	if err != nil {
		return Generate(func() (T, bool, error) {
			return *new(T), false, err
		})
	}
	secondSet := make(map[T]bool)
	for _, item := range secondSlice {
		secondSet[item] = true
	}
	var result []T
	seen := make(map[T]bool)
	for _, item := range firstSlice {
		if secondSet[item] && !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return Slice(result)
}

// Except returns an enumerator with elements from first enumerator not in second.
// Preserves order from first enumerator.
func Except[T comparable](first, second Enumerator[T]) Enumerator[T] {
	firstSlice, err := ToSlice(first)
	if err != nil {
		return Generate(func() (T, bool, error) {
			return *new(T), false, err
		})
	}
	secondSlice, err := ToSlice(second)
	if err != nil {
		return Generate(func() (T, bool, error) {
			return *new(T), false, err
		})
	}
	secondSet := make(map[T]bool)
	for _, item := range secondSlice {
		secondSet[item] = true
	}
	var result []T
	for _, item := range firstSlice {
		if !secondSet[item] {
			result = append(result, item)
		}
	}
	return Slice(result)
}
