package sort

import "slices"

type MapEntry[K comparable, V any] struct {
	Key   K
	Value V
}

type MapEntrySlice[K comparable, V any] []MapEntry[K, V]

func (s MapEntrySlice[K, V]) Sort(cmp func(a, b MapEntry[K, V]) int) MapEntrySlice[K, V] {
	slices.SortFunc(s, cmp)
	return s
}

func SliceFromMap[K comparable, V any](m map[K]V) MapEntrySlice[K, V] {
	s := make([]MapEntry[K, V], 0, len(m))
	for k, v := range m {
		s = append(s, MapEntry[K, V]{k, v})
	}
	return s
}
