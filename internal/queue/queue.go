package queue

type Queue[V any] interface {
	Enqueue(V)
	Dequeue() V
	Len() int
	Peek() V
	Iter() QueueIter[V]
}

type QueueIter[T any] interface {
	HasNext() bool
	Next() T
}

type queue[V any] struct {
	buf  []V
	head int
	iter int
	tail int
	len  int // total size
}

func (q *queue[V]) Enqueue(value V) {
	if q.len == len(q.buf) {
		q.grow()
	}
	q.buf[q.tail] = value
	q.tail = (q.tail + 1) % len(q.buf)
	q.len++
}

func (q *queue[V]) grow() {
	newBuf := make([]V, len(q.buf)<<1)
	if q.head < q.tail {
		copy(newBuf, q.buf[q.head:q.tail])
	} else {
		n := copy(newBuf, q.buf[q.head:])
		copy(newBuf[n:], q.buf[:q.tail])
	}
	q.head = 0
	q.tail = q.len
	q.buf = newBuf
}

func (q *queue[V]) Dequeue() V {
	var empty V
	if q.head == q.tail && q.len == 0 {
		return empty
	}

	v := q.buf[q.head]
	q.buf[q.head] = empty
	// advance pointer
	q.head = q.wrap(q.head)
	q.len--
	return v
}

func (q *queue[V]) wrap(i int) int {
	return (i + 1) & (len(q.buf) - 1)
}

func (q *queue[V]) Peek() V {
	if q.len == 0 {
		var empty V
		return empty
	}
	return q.buf[q.head]
}

func (q *queue[V]) Iter() QueueIter[V] {
	q.iter = 0
	return q
}

func (q *queue[V]) Len() int {
	return q.len
}

func (q *queue[V]) HasNext() bool {
	return q.iter != q.len
}

func (q *queue[V]) Next() V {
	v := q.buf[(q.head+q.iter)&(len(q.buf)-1)]
	q.iter = q.wrap(q.iter)
	return v
}

func New[V any]() Queue[V] {
	return &queue[V]{
		buf: make([]V, 2),
	}
}
