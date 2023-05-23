package egomap

import (
	"sync"
	"sync/atomic"

	"github.com/dingyuchen/egomap/internal/oplog"
)

type Writer[K comparable, V any] interface {
	Set(K, V)
	Delete(K)
	Refresh()
}

type writeHandler[K comparable, V any] interface {
	register(*reader[K, V])
	unregister(uint32)
}

type writer[K comparable, V any] struct {
	mu       *sync.RWMutex
	innerMap *leftRightMap[K, V]
	oplog    oplog.Log[K, V]
	readers  map[uint32]*reader[K, V]
	seen     []scan
}

type scan struct {
	past  uint32
	epoch *atomic.Uint32
}

func (w *writer[K, V]) Set(key K, value V) {
	w.oplog.AddWrite(key, value)
}

func (w *writer[K, V]) Delete(key K) {
	w.oplog.AddDelete(key)
}

func (w *writer[K, V]) Refresh() {
	w.mu.RLock()
	for len(w.seen) > 0 {
		i := 0
		for _, r := range w.seen {
			if epoch := r.epoch.Load(); r.past == epoch {
				w.seen[i] = r
				i++
			}
		}
		w.seen = w.seen[:i]
	}
	w.applyWrites()
	w.innerMap.swap()
	for _, r := range w.readers {
		if epoch := r.epoch.Load(); epoch%2 != 0 {
			w.seen = append(w.seen, scan{
				past:  epoch,
				epoch: r.epoch,
			})
		}
	}
	w.mu.RUnlock()
}

func (w *writer[K, V]) applyWrites() {
	m := w.innerMap.writeable()
	w.oplog.Apply(m)
}

func newWriter[K comparable, V any](innerMap *leftRightMap[K, V]) *writer[K, V] {
	return &writer[K, V]{
		innerMap: innerMap,
		mu:       new(sync.RWMutex),
		oplog:    oplog.New[K, V](),
		readers:  map[uint32]*reader[K, V]{},
		seen:     []scan{},
	}
}

func (w *writer[K, V]) register(r *reader[K, V]) {
	w.mu.Lock()
	w.readers[r.id] = r
	w.mu.Unlock()
}

func (w *writer[K, V]) unregister(id uint32) {
	w.mu.Lock()
	delete(w.readers, id)
	w.mu.Unlock()
}
