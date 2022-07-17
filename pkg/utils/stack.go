package utils

type Stack struct {
	Items []interface{}
}

func NewStack() *Stack {
	return &Stack{
		Items: []interface{}{},
	}
}

func (s *Stack) Push(item interface{}) {
	s.Items = append(s.Items, item)
}

func (s *Stack) Pop() interface{} {
	var item interface{} = nil

	if len(s.Items) > 0 {
		item = s.Items[len(s.Items)-1:][0]
		s.Items = s.Items[:len(s.Items)-1]
	}

	return item
}
