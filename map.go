package egomap

import (
	"sync"
	"sync/atomic"
)

type leftRightMap[K comparable, V any] struct {
	backingMaps [2]map[K]V
	ptr         *atomic.Int32
}

func (m *leftRightMap[K, V]) readable() map[K]V {
	return m.backingMaps[m.ptr.Load()]
}

func (m *leftRightMap[K, V]) writeable() map[K]V {
	return m.backingMaps[1-m.ptr.Load()]
}

func (m *leftRightMap[K, V]) swap() {
	m.ptr.Store(1 - m.ptr.Load())
}

type MapHandle[K comparable, V any] interface {
	Writer[K, V]
	Reader() Reader[K, V]
}

type mapHandle[K comparable, V any] struct {
	readerHandle ReadHandler[K, V]
	writer       Writer[K, V]
	freq         int
	count        int
	mu           *sync.Mutex
}

func (mh *mapHandle[K, V]) Refresh() {
	mh.mu.Lock()
	mh.writer.Refresh()
	mh.mu.Unlock()
}

func (mh *mapHandle[K, V]) Reader() Reader[K, V] {
	return mh.readerHandle.Reader()
}

func (mh *mapHandle[K, V]) Set(k K, v V) {
	mh.mu.Lock()
	mh.writer.Set(k, v)
	mh.tick()
	mh.mu.Unlock()
}

func (mh *mapHandle[K, V]) tick() {
	mh.count++
	if mh.count%mh.freq == 0 {
		mh.writer.Refresh()
		mh.count = 0
	}
}

func (mh *mapHandle[K, V]) Delete(k K) {
	mh.mu.Lock()
	mh.writer.Delete(k)
	mh.tick()
	mh.mu.Unlock()
}

func NewHandle[K comparable, V any](freq int) MapHandle[K, V] {
	if freq < 1 {
		panic("freq must be at least 1")
	}
	readerHandle, writer := New[K, V]()
	mh := &mapHandle[K, V]{
		readerHandle: readerHandle,
		writer:       writer,
		mu:           new(sync.Mutex),
		freq:         freq,
		count:        0,
	}
	return mh
}

func New[K comparable, V any]() (ReadHandler[K, V], Writer[K, V]) {
	innerMap := &leftRightMap[K, V]{
		backingMaps: [2]map[K]V{
			{},
			{},
		},
		ptr: &atomic.Int32{},
	}
	writer := NewWriter(innerMap)
	return NewReadHandler(innerMap, writer), writer
}
