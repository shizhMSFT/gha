package set

type Set[T comparable] map[T]struct{}

func New[T comparable](items ...T) Set[T] {
	s := Set[T](make(map[T]struct{}))
	for _, item := range items {
		s.Add(item)
	}
	return s
}

func (s Set[T]) Add(item T) {
	s[item] = struct{}{}
}

func (s Set[T]) Len() int {
	return len(s)
}

func (s Set[T]) Contains(v T) bool {
	_, ok := s[v]
	return ok
}
