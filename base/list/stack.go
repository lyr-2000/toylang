package list

type Stack struct {
	Queue *Queue
}

func NewStack() *Stack {
	q := new(Queue)
	return &Stack{Queue: q}
}
func StackPush(s *Stack, v interface{}) {
	QueuePrepend(s.Queue, v)

}

func (s *Stack) Top() interface{} {
	return QueuePeek(s.Queue)
}
func (s *Stack) Len() int {
	return s.StackSize()
}
func (s *Stack) StackSize() int {
	return QueueSize(s.Queue)
}
func StackPop(s *Stack) interface{} {
	return QueuePoll(s.Queue)
}

func (s *Stack) Pop() interface{} {
	return StackPop(s)
}
func (s *Stack) Push(v interface{}) {
	StackPush(s, v)
}
