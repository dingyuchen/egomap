package oplog

import "github.com/dingyuchen/egomap/internal/queue"

type Log[K comparable, V any] interface {
	AddWrite(K, V)
	AddDelete(K)
	Poll() (Operation[K, V], error)
}

const (
	Write = iota
	Delete
)

type op int

type Operation[K comparable, V any] struct {
	Inst    op
	Payload opsData[K, V]
}

type opsData[K comparable, V any] struct {
	Key   K
	Value V
}

type oplog[K comparable, V any] struct {
	queue *queue.Queue[Operation[K, V]]
}

func (l *oplog[K, V]) AddWrite(key K, value V) {
	op := Operation[K, V]{
		Inst: Write,
		Payload: opsData[K, V]{
			key,
			value,
		},
	}
	l.queue.Enqueue(op)
}

func (l *oplog[K, V]) AddDelete(key K) {
	op := Operation[K, V]{
		Inst: Write,
		Payload: opsData[K, V]{
			Key: key,
		},
	}
	l.queue.Enqueue(op)
}

func (l *oplog[K, V]) Poll() (Operation[K, V], error) {
	return l.queue.Dequeue()
}

func New[K comparable, V any]() Log[K, V] {
	return &oplog[K, V]{
		queue: queue.New[Operation[K, V]](),
	}
}
