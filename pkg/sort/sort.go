package sort

import (
	"sort"

	"golang.org/x/exp/constraints"
)

// Sort sorts data.
func Sort[T constraints.Ordered](data []T) {
	sort.Slice(data, func(i, j int) bool { return data[i] < data[j] })
}
