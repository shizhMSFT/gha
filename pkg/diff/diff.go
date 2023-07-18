package diff

import "fmt"

// Change is a change to a field
type Change struct {
	Field string
	Old   string
	New   string
}

// Diff is the difference between two items
type Diff[T any] struct {
	Item    T
	Changes []Change
}

// DiffString returns the difference between two strings
func DiffString(field string, old, new string) (Change, bool) {
	if old == new {
		return Change{}, false
	}
	return Change{Field: field, Old: old, New: new}, true
}

// DiffSet returns the difference between two slice-based sets
func DiffSet[T comparable](field string, old, new []T) (Change, bool) {
	if len(old) != len(new) {
		return Change{Field: field, Old: fmt.Sprint(old), New: fmt.Sprint(new)}, true
	}
	oldSet := make(map[T]struct{})
	for _, item := range old {
		oldSet[item] = struct{}{}
	}
	newSet := make(map[T]struct{})
	for _, item := range new {
		newSet[item] = struct{}{}
	}
	for item := range oldSet {
		if _, ok := newSet[item]; !ok {
			return Change{Field: field, Old: fmt.Sprint(old), New: fmt.Sprint(new)}, true
		}
	}
	return Change{}, false
}
