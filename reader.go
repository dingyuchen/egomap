package egomap

import "sync/atomic"

type Reader[K comparable, V any] interface {
	Get(K) (V, bool)
	Unregister()
}

type reader[K comparable, V any] struct {
	id         int
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

func (r *reader[K, V]) Unregister() {
	r.removeSelf()
}

type ReadHandler[K comparable, V any] interface {
	NewReader() Reader[K, V]
}

type readhandler[K comparable, V any] struct {
	counter  int
	innerMap *leftRightMap[K, V]
	writer   writeHandler[K, V]
}

func (rh *readhandler[K, V]) NewReader() Reader[K, V] {
	reader := &reader[K, V]{
		id:         rh.counter,
		epoch:      &atomic.Uint32{},
		innerMap:   rh.innerMap,
		removeSelf: func() { rh.writer.unregister(rh.counter) },
	}
	rh.counter++
	rh.writer.register(reader)
	return reader
}

func NewReadHandler[K comparable, V any](innerMap *leftRightMap[K, V], writer *writer[K, V]) ReadHandler[K, V] {
	return &readhandler[K, V]{
		innerMap: innerMap,
		writer:   writer,
	}
}
