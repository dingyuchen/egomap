package oplog

import "github.com/dingyuchen/egomap/internal/queue"

type Log[K comparable, V any] interface {
	AddWrite(K, V)
	AddDelete(K)
	Poll() []Operation[K, V]
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
	queue   *queue.Queue[Operation[K, V]]
	backLog *queue.Queue[Operation[K, V]]
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

func (l *oplog[K, V]) Poll() []Operation[K, V] {
	ops := make([]Operation[K, V], 0, l.backLog.Len()+l.queue.Len())
	for op, err := l.backLog.Dequeue(); err == nil; op, err = l.backLog.Dequeue() {
		ops = append(ops, op)
	}
	for op, err := l.queue.Dequeue(); err == nil; op, err = l.queue.Dequeue() {
		ops = append(ops, op)
		l.backLog.Enqueue(op)
	}
	return ops
}

func New[K comparable, V any]() Log[K, V] {
	return &oplog[K, V]{
		queue:   queue.New[Operation[K, V]](),
		backLog: queue.New[Operation[K, V]](),
	}
}
