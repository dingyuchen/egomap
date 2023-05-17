package egomap

import "sync/atomic"

type Reader[K comparable, V any] interface {
	Get(K) (V, bool)
	Close()
}

type reader[K comparable, V any] struct {
	id         uint32
	innerMap   *leftRightMap[K, V]
	epoch      *atomic.Uint32
	removeSelf func()
}

func (r *reader[K, V]) Get(key K) (V, bool) {
	r.epoch.Add(1)
	m := r.innerMap.readable()
	v, ok := m[key]
	r.epoch.Add(1)
	return v, ok
}

func (r *reader[K, V]) Close() {
	r.removeSelf()
}

type ReadHandler[K comparable, V any] interface {
	Reader() Reader[K, V]
}

type readhandler[K comparable, V any] struct {
	counter  *atomic.Uint32
	innerMap *leftRightMap[K, V]
	writer   writeHandler[K, V]
}

func (rh *readhandler[K, V]) Reader() Reader[K, V] {
	id := rh.counter.Load()
	reader := &reader[K, V]{
		id:         id,
		epoch:      &atomic.Uint32{},
		innerMap:   rh.innerMap,
		removeSelf: func() { rh.writer.unregister(id) },
	}
	rh.counter.Add(1)
	rh.writer.register(reader)
	return reader
}

func NewReadHandler[K comparable, V any](innerMap *leftRightMap[K, V], writer *writer[K, V]) ReadHandler[K, V] {
	return &readhandler[K, V]{
		innerMap: innerMap,
		writer:   writer,
	}
}
