package queue

type node[T any] struct {
	data T
	next *node[T]
}

type Queue[T any] struct {
	head *node[T]
	tail *node[T]
}

func New[T any]() Queue[T] {
	return Queue[T]{
		head: nil,
		tail: nil,
	}
}

func (q *Queue[T]) Push(data T) {
	node := &node[T]{data: data, next: nil}

	if q.head == nil {
		q.head = node
		q.tail = q.head
		return
	}

	q.tail.next = node
	q.tail = q.tail.next
}

func (q *Queue[T]) Pop() (T, bool) {
	var result T

	if q.head == nil {
		return result, false
	}

	result = q.head.data
	q.head = q.head.next
	return result, true
}
