package stack

type Stack[T any] struct {
	val []T
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{
		val: make([]T, 0),
	}
}

func (s *Stack[T]) Push(val T) {
	s.val = append(s.val, val)
}

func (s *Stack[T]) Pop() T {
	last := s.Last()
	s.val = s.val[:len(s.val)-1]
	return last
}

func (s *Stack[T]) Empty() bool {
	return len(s.val) == 0
}

func (s *Stack[T]) Len() int {
	return len(s.val)
}

func (s *Stack[T]) Last() T {
	return s.val[len(s.val)-1]
}
