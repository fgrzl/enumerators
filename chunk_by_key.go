package enumerators

func ChunkByKey[T any, K comparable](
	in Enumerator[T],
	keyFunc func(T) K,
) Enumerator[Enumerator[KeyedItem[K, T]]] {
	var last K
	first := true
	return ChunkWhen(in, last, func(prev K, item T) (K, bool, error) {
		curr := keyFunc(item)
		if first {
			first = false
			return curr, false, nil // first group starts without split
		}
		return curr, curr != prev, nil
	})
}
