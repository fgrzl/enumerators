package enumerators

// ChunkByKey splits an input stream into chunks based on a computed key.
// Each chunk yields items sharing the same key, wrapped in a KeyedChunk.
func ChunkByKey[V any, K comparable](
	in Enumerator[V],
	keyFunc func(V) K,
) Enumerator[KeyedChunk[K, V]] {
	var last *K
	return ChunkWhen(in, func(item V) (K, bool, error) {
		curr := keyFunc(item)
		if last == nil {
			last = &curr
			return curr, false, nil
		}
		split := *last != curr
		*last = curr
		return curr, split, nil
	})
}
