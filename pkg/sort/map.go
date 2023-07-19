package sort

import (
	"sort"

	"golang.org/x/exp/constraints"
)

type MapEntry[K comparable, V any] struct {
	Key   K
	Value V
}

type MapEntrySlice[K comparable, V constraints.Ordered] []MapEntry[K, V]

func (s MapEntrySlice[K, V]) Len() int           { return len(s) }
func (s MapEntrySlice[K, V]) Less(i, j int) bool { return s[i].Value < s[j].Value }
func (s MapEntrySlice[K, V]) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func SliceFromMap[K comparable, V constraints.Ordered](m map[K]V, ascending bool) []MapEntry[K, V] {
	s := make([]MapEntry[K, V], 0, len(m))
	for k, v := range m {
		s = append(s, MapEntry[K, V]{k, v})
	}
	if ascending {
		sort.Sort(MapEntrySlice[K, V](s))
	} else {
		sort.Sort(sort.Reverse(MapEntrySlice[K, V](s)))
	}
	return s
}
