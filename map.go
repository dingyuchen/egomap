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
	Reader() Reader[K, V]
	Set(K, V)
	Refresh()
}

type pair[K comparable, V any] struct {
	key   K
	value V
}

type mapHandle[K comparable, V any] struct {
	readerHandle ReadHandler[K, V]
	writer       Writer[K, V]
	writeChan    chan (pair[K, V])
	freq         int
	mu           *sync.Mutex
}

func (mh *mapHandle[K, V]) bgWrite() {
	count := 0
	for p := range mh.writeChan {
		mh.writer.Set(p.key, p.value)
		count++
		if count%mh.freq == 0 {
			mh.writer.Refresh()
			mh.freq = 0
		}
	}
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
	mh.writeChan <- pair[K, V]{k, v}
}

func NewHandle[K comparable, V any](freq int) MapHandle[K, V] {
	if freq < 1 {
		panic("freq must be at least 1")
	}
	readerHandle, writer := New[K, V]()
	mh := &mapHandle[K, V]{
		readerHandle: readerHandle,
		writer:       writer,
		writeChan:    make(chan pair[K, V], freq),
		mu:           new(sync.Mutex),
		freq:         freq,
	}
	go mh.bgWrite()
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
