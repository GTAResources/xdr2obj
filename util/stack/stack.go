package stack

type Item struct {
	Value interface{}
}

type Interface interface {
	Push(item *Item)
	Pop() *Item
}

type Stack struct {
	stack []*Item
	count int
}

func (s *Stack) Allocate(size int) {
	s.stack = make([]*Item, size)
}

func (s *Stack) Push(item *Item) {
	if s.count >= len(s.stack) {
		stack := make([]*Item, len(s.stack)*2)
		copy(stack, s.stack)
		s.stack = stack
	}
	s.stack[s.count] = item
	s.count++
}

func (s *Stack) Pop() *Item {
	if s.count == 0 {
		return nil
	}
	item := s.stack[s.count-1]
	s.count--
	return item
}
