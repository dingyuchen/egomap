package queue

import "errors"

var ErrEmpty = errors.New("popping from empty queue")

const (
	maxInternalSlice    = 128
	initialInternalSize = 16
)

type Queue[V any] struct {
	head *node[V]
	tail *node[V]
	ptr  int
	len  int // total size
}

func (q *Queue[V]) Enqueue(value V) {
	if len(q.tail.arr) < maxInternalSlice {
		q.tail.arr = append(q.tail.arr, value)
	} else {
		arr := make([]V, 0, maxInternalSlice)
		arr = append(arr, value)
		q.tail.next = &node[V]{
			arr: arr,
		}
		q.tail = q.tail.next
	}
	q.len++
}

func (q *Queue[V]) Dequeue() (V, error) {
	var empty V
	if q.len == 0 {
		return empty, ErrEmpty
	}

	v := q.head.arr[q.ptr]
	q.head.arr[q.ptr] = empty // allow value to be freed
	// advance pointer
	q.ptr++
	q.len--
	if q.ptr >= len(q.head.arr) {
		q.ptr = 0
		q.head = q.head.next
	}
	if q.head == nil {
		start := new(node[V])
		q.head = start
		q.tail = start
	}
	return v, nil
}

func (q *Queue[V]) Peek() *V {
	if q.len == 0 {
		return nil
	}
	return &q.head.arr[q.ptr]
}

func (q *Queue[V]) Iter() QueueIter[V] {
	return &queueIter[V]{
		q.head,
		q.ptr,
	}
}

func (q *Queue[V]) Len() int {
	return q.len
}

func New[V any]() *Queue[V] {
	start := new(node[V])
	return &Queue[V]{
		head: start,
		tail: start,
	}
}

type QueueIter[T any] interface {
	HasNext() bool
	Next() T
}

type queueIter[T any] struct {
	node *node[T]
	ptr  int
}

func (i *queueIter[T]) HasNext() bool {
	return i.node != nil && i.ptr < len(i.node.arr)
}

func (i *queueIter[T]) Next() T {
	if i.node == nil {
		var empty T
		return empty
	}
	v := i.node.arr[i.ptr]
	i.ptr++
	if i.ptr >= len(i.node.arr) {
		i.ptr = 0
		i.node = i.node.next
	}
	return v
}

type node[V any] struct {
	arr  []V
	next *node[V]
}
