package linkedlist

type Node[T any] struct {
	Value T
	Next  *Node[T]
}

type LinkedList[T any] struct {
	head *Node[T]
	size int
}

func NewList[T any]() *LinkedList[T] {
	return &LinkedList[T]{}
}

func (l *LinkedList[T]) Insert(v T) {
	newNode := &Node[T]{
		Value: v,
		Next:  nil,
	}

	if l.head == nil {
		l.head = newNode
		l.size++
		return
	}

	cur := l.head
	for cur.Next != nil {
		cur = cur.Next
	}

	cur.Next = newNode
	l.size++
}

func (l *LinkedList[T]) Head() *Node[T] {
	return l.head
}

func (l *LinkedList[T]) Len() int {
	return l.size
}

func (l *LinkedList[T]) Remove(pred func(T) bool) bool {
	var prev *Node[T]
	cur := l.head

	for cur != nil {
		if pred(cur.Value) {
			if prev == nil {
				l.head = cur.Next
			} else {
				prev.Next = cur.Next
			}

			l.size--
			return true
		}

		prev = cur
		cur = cur.Next
	}

	return false
}
