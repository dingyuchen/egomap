package queue

import "errors"

var ErrEmpty = errors.New("popping from empty queue")

type Queue[V any] struct {
	arr []V
}

func (q *Queue[V]) Enqueue(value V) {
}

func (q *Queue[V]) Dequeue() (V, error) {
	return q.arr[0], nil
}

func (q *Queue[V]) Peek() *V {
	return nil
}

func (q *Queue[V]) Len() int {
	return 0
}
