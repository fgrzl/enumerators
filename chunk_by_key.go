package enumerators

func ChunkByKey[T any, TCompare comparable](
	in Enumerator[T],
	keyFunc func(T) TCompare,
) Enumerator[Enumerator[T]] {
	return ChunkWhen(in, nil, func(prev *TCompare, item T) (*TCompare, bool, error) {
		curr := keyFunc(item)
		if prev == nil {
			return &curr, false, nil
		}
		split := *prev != curr
		return &curr, split, nil
	})
}
