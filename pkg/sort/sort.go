package sort

import (
	"sort"

	"golang.org/x/exp/constraints"
)

// GenericSlice is a generic slice type.
type GenericSlice[T constraints.Ordered] []T

func (s GenericSlice[T]) Len() int           { return len(s) }
func (s GenericSlice[T]) Less(i, j int) bool { return s[i] < s[j] }
func (s GenericSlice[T]) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// IsSorted reports whether data is sorted.
func IsSorted[T constraints.Ordered](data []T) bool {
	return sort.IsSorted(GenericSlice[T](data))
}

// Sort sorts data.
func Sort[T constraints.Ordered](data []T) {
	sort.Sort(GenericSlice[T](data))
}
