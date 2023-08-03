package sort

import "sort"

type MapEntry[K comparable, V any] struct {
	Key   K
	Value V
}

type MapEntrySlice[K comparable, V any] []MapEntry[K, V]

func (s MapEntrySlice[K, V]) Sort(less func(s []MapEntry[K, V], i, j int) bool) MapEntrySlice[K, V] {
	sort.Slice(s, func(i, j int) bool {
		return less(s, i, j)
	})
	return s
}

func SliceFromMap[K comparable, V any](m map[K]V) MapEntrySlice[K, V] {
	s := make([]MapEntry[K, V], 0, len(m))
	for k, v := range m {
		s = append(s, MapEntry[K, V]{k, v})
	}
	return s
}
