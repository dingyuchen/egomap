package egomap

import (
	"sync"

	"github.com/dingyuchen/egomap/internal/oplog"
	"github.com/dingyuchen/egomap/internal/queue"
)

type Writer[K comparable, V any] interface {
	Set(K, V)
	Refresh()
}

type writeHandler[K comparable, V any] interface {
	register(*reader[K, V])
	unregister(int)
}

type writer[K comparable, V any] struct {
	mu       *sync.RWMutex
	innerMap *leftRightMap[K, V]
	oplog    oplog.Log[K, V]
	readers  map[int]*reader[K, V]
	scan     queue.Queue[*reader[K, V]]
}

func (w *writer[K, V]) Set(key K, value V) {
	w.oplog.AddWrite(key, value)
}

func (w *writer[K, V]) Refresh() {
	w.mu.RLock()
	for w.scan.Len() > 0 {
		r, _ := w.scan.Dequeue()
		if epoch := r.epoch.Load(); epoch%2 != 0 {
			w.scan.Enqueue(r)
		}
	}
	w.applyWrites()
	w.innerMap.swap()
	for _, r := range w.readers {
		if epoch := r.epoch.Load(); epoch%2 != 0 {
			w.scan.Enqueue(r)
		}
	}
	w.mu.RUnlock()
}

func (w *writer[K, V]) applyWrites() {
	m := w.innerMap.writeable()
	for op, err := w.oplog.Poll(); err != nil; op, err = w.oplog.Poll() {
		switch op.Inst {
		case oplog.Write:
			m[op.Payload.Key] = op.Payload.Value
		case oplog.Delete:
			delete(m, op.Payload.Key)
		}

	}
}

func NewWriter[K comparable, V any](innerMap *leftRightMap[K, V]) *writer[K, V] {
	return &writer[K, V]{
		innerMap: innerMap,
		mu:       new(sync.RWMutex),
		oplog:    oplog.New[K, V](),
	}
}

func (w *writer[K, V]) register(r *reader[K, V]) {
	w.mu.Lock()
	w.readers[r.id] = r
	w.mu.Unlock()
}

func (w *writer[K, V]) unregister(id int) {
	w.mu.Lock()
	delete(w.readers, id)
	w.mu.Unlock()
}
