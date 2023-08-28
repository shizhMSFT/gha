package set

type Set[T comparable] map[T]struct{}

func New[T comparable]() Set[T] {
	return make(map[T]struct{})
}

func (s Set[T]) Add(v T) {
	s[v] = struct{}{}
}

func (s Set[T]) Len() int {
	return len(s)
}

func (s Set[T]) Contains(v T) bool {
	_, ok := s[v]
	return ok
}
