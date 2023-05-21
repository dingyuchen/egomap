package oplog

type Log[K comparable, V any] interface {
	AddWrite(K, V)
	AddDelete(K)
	Apply(map[K]V)
}

type Op int

const (
	Write Op = iota
	Delete
)

type Operation[K comparable, V any] struct {
	Inst    Op
	Payload opsData[K, V]
}

type opsData[K comparable, V any] struct {
	Key   K
	Value V
}

type oplog[K comparable, V any] struct {
	queue   []Operation[K, V]
	backLog []Operation[K, V]
}

func (l *oplog[K, V]) AddWrite(key K, value V) {
	op := Operation[K, V]{
		Inst: Write,
		Payload: opsData[K, V]{
			key,
			value,
		},
	}
	l.queue = append(l.queue, op)
}

func (l *oplog[K, V]) AddDelete(key K) {
	op := Operation[K, V]{
		Inst: Write,
		Payload: opsData[K, V]{
			Key: key,
		},
	}
	l.queue = append(l.queue, op)
}

func (l *oplog[K, V]) Apply(m map[K]V) {
	// pop backlog
	for _, op := range l.backLog {
		switch op.Inst {
		case Write:
			m[op.Payload.Key] = op.Payload.Value
		case Delete:
			delete(m, op.Payload.Key)
		}
	}
	l.backLog = l.backLog[:0]

	for _, op := range l.queue {
		switch op.Inst {
		case Write:
			m[op.Payload.Key] = op.Payload.Value
		case Delete:
			delete(m, op.Payload.Key)
		}
	}
	l.backLog, l.queue = l.queue, l.backLog
}

func New[K comparable, V any]() Log[K, V] {
	return &oplog[K, V]{
		queue:   make([]Operation[K, V], 0, 2),
		backLog: make([]Operation[K, V], 0, 2),
	}
}
