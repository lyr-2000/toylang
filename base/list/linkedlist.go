package list

import (
	"bytes"
	"fmt"
	"strings"
)

type Node = DoubleNode
type DoubleNode struct {
	Value      interface{}
	Next, Prev *DoubleNode
}

func (n *Node) String() string {
	return fmt.Sprintf("Node{Value: %v, Next: %v}", n.Value, n.Next)
}

type LinkedList = DualLinkedList
type DualLinkedList struct {
	Head, Tail *Node
	ListCnt    int
}

func (it *LinkedList) Len() int {
	return it.ListCnt
}
func (l *LinkedList) PopLast() interface{} {
	return ListRemoveLast(l)
}
func (l *LinkedList) PopFirst() interface{} {

	return ListRemoveFirst(l)

}

func (l *LinkedList) String() string {
	if l == nil {
		return "Node[null]"
	}
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("Node(cnt=%d)[", l.ListCnt))

	var temp = l.Head
	for temp != nil {
		buf.Write([]byte(fmt.Sprintf("%v", temp.Value)))
		if temp.Next != nil {
			buf.WriteRune(',')
		}
		temp = temp.Next
	}
	buf.WriteString("]")
	return buf.String()
}
func ListAppend(l *LinkedList, v interface{}) *LinkedList {
	if l == nil {
		l = &DualLinkedList{}
	}
	if l.Head == nil {
		l.Head = &Node{Value: v}
		l.Tail = l.Head
	} else {
		l.Tail.Next = &Node{Value: v}

		l.Tail = l.Tail.Next
	}
	l.ListCnt++
	return l
}

func ListPrepend(l *LinkedList, v interface{}) *LinkedList {
	if l == nil {
		l = &DualLinkedList{}
	}

	if l.Head == nil {
		l.Head = &Node{Value: v}
		l.Tail = l.Head
	} else {

		l.Head.Prev = &Node{Value: v}
		l.Head = l.Head.Prev
	}
	l.ListCnt++
	return l
}

func ListPeekFirst(l *LinkedList) interface{} {
	if l == nil {
		return nil
	}
	if l.Head == nil {
		return nil
	}
	return l.Head.Value
}

func ListPeekLast(l *LinkedList) interface{} {
	if l == nil || l.Head == nil {
		return nil
	}
	return l.Tail.Value
}
func ListRemoveFirst(l *LinkedList) interface{} {
	if l == nil {
		return nil
	}
	if l.Head == nil {
		return nil
	}
	v := l.Head.Value
	removed := l.Head

	l.Head = l.Head.Next
	l.Head.Prev = nil  //remove prev pointer
	removed.Next = nil //remove next pointer , help gc
	l.ListCnt--
	return v
}

func ListRemoveLast(l *LinkedList) interface{} {
	if l == nil || l.Head == nil {
		return nil
	}

	v := l.Tail.Value
	removed := l.Tail
	l.Tail = l.Tail.Prev
	l.Tail.Next = nil
	removed.Prev = nil //remove prev pointer,help gc
	l.ListCnt--
	return v
}

type SingleNode struct {
	Value interface{}
	Next  *SingleNode
}

type Queue struct {
	Head, Tail *SingleNode
	ListCnt    int //list count
}

func (q *Queue) String() string {
	if q == nil {
		return "node{}"
	}
	s := strings.Builder{}
	s.WriteString("node{}")
	node := q.Head
	for node != nil {
		s.WriteString(fmt.Sprintf("%+v,", node.Value))
		node = node.Next
	}
	return s.String()
}
func QueueAppend(q *Queue, v interface{}) *Queue {
	if q == nil || q.Head == nil {
		q = &Queue{}
	}
	if q.Head == nil {
		q.Head = &SingleNode{Value: v}
		q.Tail = q.Head
	} else {
		q.Tail.Next = &SingleNode{Value: v}
		q.Tail = q.Tail.Next
	}
	q.ListCnt++
	return q
}
func QueuePoll(q *Queue) interface{} {
	if q == nil || q.Head == nil {
		return nil
	}
	remove := q.Head
	v := q.Head.Value
	q.Head = q.Head.Next
	q.ListCnt--
	remove.Next = nil //help gc
	return v
}
func QueueClear(q *Queue) {
	if q == nil {
		return
	}
	if q.ListCnt < 0 {
		panic("illegal queue state")
	}
	for q.ListCnt > 0 {
		QueuePoll(q)
	}
}
func QueuePeek(q *Queue) interface{} {
	if q == nil || q.Head == nil {
		return nil
	}

	return q.Head.Value
}
func QueueSize(q *Queue) int {
	if q == nil {
		return 0
	}
	return q.ListCnt
}

func QueuePrepend(q *Queue, v interface{}) *Queue {
	if q == nil {
		q = &Queue{}
	}
	if q.Head == nil {
		q.Head = &SingleNode{Value: v}
		q.Tail = q.Head
	} else {
		v := &SingleNode{Value: v}
		v.Next = q.Head
		q.Head = v
	}

	q.ListCnt++
	return q
}
