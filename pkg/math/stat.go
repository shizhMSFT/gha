package math

import "golang.org/x/exp/constraints"

// Min returns the minimum value in s.
func Min[T constraints.Ordered](s []T) T {
	if len(s) == 0 {
		panic("empty slice")
	}
	min := s[0]
	for _, x := range s[1:] {
		if x < min {
			min = x
		}
	}
	return min
}

// Max returns the maximum value in s.
func Max[T constraints.Ordered](s []T) T {
	if len(s) == 0 {
		panic("empty slice")
	}
	max := s[0]
	for _, x := range s[1:] {
		if x > max {
			max = x
		}
	}
	return max
}

// Mean returns the mean value in s.
func Mean[T constraints.Integer | constraints.Float](s []T) T {
	if len(s) == 0 {
		panic("empty slice")
	}
	var sum T
	for _, x := range s {
		sum += x
	}
	return sum / T(len(s))
}

// Percentile returns the pth percentile value in s.
// s must be sorted.
func Percentile[T constraints.Integer | constraints.Float](s []T, p float64) T {
	if len(s) == 0 {
		panic("empty slice")
	}
	if p < 0 || p > 1 {
		panic("invalid percentile")
	}
	i := int(p * float64(len(s)-1))
	return s[i]
}

// Median returns the median value in s.
// s must be sorted.
func Median[T constraints.Integer | constraints.Float](s []T) T {
	if len(s) == 0 {
		panic("empty slice")
	}
	if len(s)%2 == 1 {
		return s[len(s)/2]
	}
	return Mean(s[len(s)/2-1 : len(s)/2+1])
}
